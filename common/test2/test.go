package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ScanEvent struct {
	CName string
	Key   string
	Field string
	Value string
}

type UserExtI interface {
	BaseDataGet(ctx context.Context, cName, key, field string) (string, error)
	BaseDataSet(ctx context.Context, cName, key, field, value string) error
	BaseDataGetAll() <-chan ScanEvent
}

func NewUserExt(mongoClient *mongo.Database, redisClient *redis.Client, logger log.Logger,
	userExt UserExtI) (*UserExt, error) {
	var defaultFieldNum int64 = 30
	ext := &UserExt{
		defaultFieldNum:    defaultFieldNum,
		mongoClient:        mongoClient,
		redisClient:        redisClient,
		mongoCollectionMap: new(sync.Map),
		menuCollection:     mongoClient.Collection("user_ext_menu"),
		logger:             log.NewHelper(logger, log.WithMessageKey("user_ext")),

		userExt: userExt,
	}
	var defaultCollectionKey = ext.getMongoCollectionIndex(0)
	var defaultCollection = mongoClient.Collection(defaultCollectionKey)
	ext.mongoCollectionMap.Store(defaultCollectionKey, defaultCollection)
	ext.defaultCollection = defaultCollection
	return ext, nil
}

type UserExt struct {
	defaultFieldNum    int64
	redisClient        *redis.Client
	mongoClient        *mongo.Database
	mongoCollectionMap *sync.Map
	defaultCollection  *mongo.Collection
	menuCollection     *mongo.Collection
	logger             *log.Helper

	userExt UserExtI
}

func (ext *UserExt) getMongoCollectionName(ctx context.Context, key string) (string, error) {
	index, err := ext.GetNextSequenceValue(ctx, ext.menuCollection, key)
	if err != nil {
		return "", err
	}
	return ext.getMongoCollectionIndex(index), nil
}

func (ext *UserExt) getMongoCollectionIndex(index int64) string {
	return "user_ext_data" + strconv.Itoa(int(index/ext.defaultFieldNum))
}

func (ext *UserExt) getCacheKey(key, field string) (string, time.Duration) {
	if field != "" {
		key += "_" + field
	}
	return key, time.Hour * 24 * 3
}

// GetNextSequenceValue 获取并更新自增序列的下一个值，使用事务保证原子性
func (ext *UserExt) GetNextSequenceValue(ctx context.Context, collection *mongo.Collection, value string) (int64, error) {
	filter := bson.M{"value": value}
	var seqDoc = struct {
		Value string `bson:"value"`
	}{
		Value: "",
	}
	// 文档不存在，进行自增操作
	update := bson.M{
		"$set": bson.M{"value": value},
	}
	opts := options.FindOneAndUpdate().SetUpsert(true)
	err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&seqDoc)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return -1, err
	}
	findOpt := options.Find().SetSort(bson.M{"_id": 1})
	cur, err := collection.Find(ctx, bson.D{}, findOpt)
	if err != nil {
		return -1, err
	}
	defer cur.Close(ctx)
	var rank int64
	for cur.Next(ctx) {
		var doc = struct {
			Value string `bson:"value"`
		}{
			Value: "",
		}
		err := cur.Decode(&doc)
		if err != nil || doc.Value == value {
			break
		}
		rank++
	}
	rank += 1
	return rank, nil
}

func (ext *UserExt) getMongoCollection(ctx context.Context, key string) (*mongo.Collection, error) {
	var result = ext.defaultCollection
	collectionI, ok := ext.mongoCollectionMap.Load(key)
	var err error
	if !ok {
		var cn string
		cn, err = ext.getMongoCollectionName(ctx, key)
		if err != nil {
			return result, err
		}
		result = ext.mongoClient.Collection(cn)
		ext.mongoCollectionMap.Store(key, result)
	} else {
		result, _ = collectionI.(*mongo.Collection)
	}
	return result, err
}

// GetFieldFromMongoAndCache 获取单个字段的值
func (ext *UserExt) GetFieldFromMongoAndCache(ctx context.Context, cName, key, field string) (string, error) {
	var cacheKey, ttl = ext.getCacheKey(key, field)
	// 从Redis缓存中获取字段值
	val, err := ext.redisClient.Get(ctx, cacheKey).Result()
	if errors.Is(err, redis.Nil) {
		// 如果Redis中没有缓存，则从MongoDB获取
		var collection, err = ext.getMongoCollection(ctx, key)
		if err != nil || collection == nil {
			return "", fmt.Errorf("[GetFieldFromMongoAndCache] failed to get MongoDB, key:%s, field:%s, collection: %v", key, field, err)
		}
		var filter = bson.D{{"_id", cName}}
		var projection = bson.D{{key, 1}}
		// 创建查询选项并设置投影
		findOptions := options.FindOne().SetProjection(projection)
		// 执行查询
		var result = make(map[string]string)
		err = collection.FindOne(context.Background(), filter, findOptions).Decode(&result)
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			return "", err
		} else if err != nil {
			// 获取原始数据
			err = nil
			val, err = ext.userExt.BaseDataGet(ctx, cName, key, field)
			if err != nil {
				return "", err
			}
			err = ext.SetFieldFromMongoAndCache(ctx, cName, key, field, val, false, true)
			if err != nil {
				return "", err
			}
		} else {
			// 将MongoDB中的字段值转换为字符串
			val, _ = result[key]
		}
		// 将字段值缓存到Redis
		err = ext.redisClient.Set(ctx, cacheKey, val, ttl).Err()
		if err != nil {
			return "", fmt.Errorf("[GetFieldFromMongoAndCache] failed to set field value in Redis: %v", err)
		}
	} else if err != nil {
		return "", err
	} else {
		ext.redisClient.Expire(ctx, cacheKey, ttl)
	}
	return val, nil
}

// SetFieldFromMongoAndCache 设置单个字段的值
func (ext *UserExt) SetFieldFromMongoAndCache(ctx context.Context, cName, key, field, value string, setOld bool, upset bool) error {
	var cacheKey, _ = ext.getCacheKey(key, field)
	var filter = bson.D{{"_id", cName}}
	var update = bson.M{
		"$set": bson.M{
			key: value,
		},
	}
	collection, err := ext.getMongoCollection(ctx, key)
	if err != nil {
		return fmt.Errorf("[SetFieldFromMongoAndCache] failed to get MongoDB, key:%s, field:%s, collection: %v", key, field, err)
	}
	var needDel bool
	// 将字段值更新到MongoDB
	if upset {
		var data = map[string]string{
			"_id": cName,
			key:   value,
		}
		_, err = collection.InsertOne(ctx, data)
		if err == nil {
			needDel = true
		}
	} else {
		var updateResult *mongo.UpdateResult
		updateResult, err = collection.UpdateOne(ctx, filter, update)
		if err != nil {
			return fmt.Errorf("[SetFieldFromMongoAndCache] failed to update document in MongoDB: %v", err)
		}
		if updateResult != nil && updateResult.ModifiedCount > 0 {
			needDel = true
		}
	}
	if needDel {
		// 将字段值缓存到Redis
		err = ext.redisClient.Del(ctx, cacheKey).Err()
		if err != nil {
			return fmt.Errorf("[SetFieldFromMongoAndCache] failed to set field value in Redis: %v", err)
		}
	}
	if setOld {
		// 更老缓存
		err = ext.userExt.BaseDataSet(ctx, cName, key, field, value)
	}
	return nil
}

// HScanFromMongoAndCache :将老缓存scan到db中
func (ext *UserExt) HScanFromMongoAndCache() {
	if ext.userExt == nil || ext.userExt.BaseDataGetAll == nil {
		ext.logger.Errorf("BaseDataGetAll is nil")
		return
	}
	ch := ext.userExt.BaseDataGetAll()
	if ch == nil {
		ext.logger.Errorf("cNext is nil")
		return
	}
	var ctx = context.Background()
	for data := range ch {
		// 如果新缓存已经有值则不更新
		err := ext.SetFieldFromMongoAndCache(ctx, data.CName, data.Key, data.Field, data.Value, false, true)
		if err != nil {
			ext.logger.Errorf("[HScanFromMongoAndCache] failed to set field value in user_ext, data:%+v, err:%v", data, err)
		}
		time.Sleep(time.Second / 10)
	}
}

func main() {
	var ctx = context.Background()
	// MongoDB连接
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	mongoClientDb := mongoClient.Database("testdb")
	// Redis连接
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "bj-crs-5p0lg7wj.sql.tencentcdb.com:20163", // Redis服务器地址
		Password: "cmd+zhM15DBMdO1afAEyOFi63",                // 密码，如果没有设置密码则为空
		DB:       8,                                          // 默认数据库
	})
	recentGameCount := &RecentGameCount{
		rd: redisClient,
	}
	ext, err := NewUserExt(mongoClientDb, redisClient, log.DefaultLogger, recentGameCount)
	if err != nil {
		log.Fatalf("Failed to create UserExt: %v", err)
	}
	result, err := ext.GetFieldFromMongoAndCache(ctx, "602096805", "recent_game_count", "602096805")
	if err != nil {
		log.Fatalf("Failed to get field value from MongoDB: %v", err)
	}
	log.Info("Field value:", result)
	ext.HScanFromMongoAndCache()
}
