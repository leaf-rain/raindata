package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	var ctx = context.Background() // 67ea9730230709ab47e9b089

	// MongoDB连接
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(ctx)

	// 选择数据库
	database := mongoClient.Database("testdb")

	// 自增序列集合
	sequenceCollection := database.Collection("sequences")
	var rankMap = new(sync.Map)
	var wg = new(sync.WaitGroup)
	for i := 0; i < 20000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := "xxx" + strconv.Itoa(i)
			id, err := getNextSequenceValue(ctx, sequenceCollection, key)
			if err != nil {
				log.Fatalf("Failed to get next sequence value: %v", err)
			}
			if id <= 0 {
				panic("id <= 0")
			}
			if ddd, ok := rankMap.Load(id); ok {
				panic(fmt.Sprintf("id=%d, key=%s, ddd=%v", id, key, ddd))
			}
			rankMap.Store(id, key)
			fmt.Printf("Document value set successfully, key: %s; id: %d\n", key, id)
			//switch key {
			//case "xxx19590":
			//	if id != 18994 {
			//		panic("xxx19590 != 18994")
			//	}
			//case "xxx18185":
			//	if id != 17944 {
			//		panic("xxx18185 != 17944")
			//	}
			//case "xxx19862":
			//	if id != 18106 {
			//		panic("xxx19862 != 18106")
			//	}
			//case "xxx15516":
			//	if id != 15624 {
			//		panic("xxx15516 != 15624")
			//	}
			//}
		}(i)
	}
	wg.Wait()
	fmt.Println("Done")
}

// getNextSequenceValue 获取并更新自增序列的下一个值，使用事务保证原子性
func getNextSequenceValue(ctx context.Context, collection *mongo.Collection, value string) (int64, error) {
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
