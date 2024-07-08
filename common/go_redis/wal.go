package go_redis

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/leaf-rain/raindata/common/recuperate"
	"go.uber.org/zap"
	"strconv"
	"time"
)

var (
	ErrBusinessStaring = errors.New("business staring")
)

type Data interface {
	GetId() int64
	GetData() string
	GetStage() int64
}

type dataInfo struct {
	id    int64
	data  string
	stage int64
}

func (d dataInfo) GetId() int64 {
	return d.id
}

func (d dataInfo) GetData() string {
	return d.data
}

func (d dataInfo) GetStage() int64 {
	return d.stage
}

type WalForTimeout struct {
	conn       redis.Cmdable
	key        string
	receiveKey string
	stageKey   string
	lock       *RedisLock
	ttl        time.Duration
	logger     *zap.Logger
}

func (w *WalForTimeout) Put(ctx context.Context, data Data) (int64, error) {
	var err error
	_id := data.GetId()
	if _id <= 0 {
		_id, err = w.conn.HIncrBy(ctx, w.key, "_id", 1).Result()
		if err != nil {
			return 0, err
		}
	}
	buf1 := make([]byte, 8)
	binary.BigEndian.PutUint64(buf1, uint64(time.Now().Unix()))
	buf2 := make([]byte, 8)
	binary.BigEndian.PutUint64(buf2, uint64(data.GetStage()))

	var buf bytes.Buffer
	buf.Write(buf1)
	buf.Write(buf2)
	buf.WriteString(data.GetData())

	subKey := strconv.FormatInt(_id, 10)

	err = w.conn.HSet(ctx, w.key, subKey, buf.String()).Err()
	w.conn.Expire(ctx, w.key, w.ttl) // 避免无限制存储
	return _id, err
}

func (w *WalForTimeout) Del(ctx context.Context, id int64) (err error) {
	subKey := strconv.FormatInt(id, 10)
	err = w.conn.HDel(ctx, w.key, subKey).Err()
	w.conn.Expire(ctx, w.key, w.ttl) // 避免无限制存储
	return err
}

func (w *WalForTimeout) GetForTimeout(ctx context.Context, t int64) (chan Data, error) {
	lockResult, err := w.lock.TryLock(ctx)
	defer w.lock.UnLock(ctx)
	if err != nil {
		return nil, err
	}
	if !lockResult {
		return nil, ErrBusinessStaring
	}
	var result = make(chan Data)
	recuperate.GoSafe(func() {
		defer close(result)
		var cursor uint64
		var keys []string
		var _id, receiveTime, stage int64
		var tmp string
		var now = time.Now().Unix()
		for {
			cmd := w.conn.HScan(ctx, w.key, cursor, "*", 200)
			// 处理命令结果
			keys, cursor, err = cmd.Result()
			if err != nil {
				w.logger.Error("redis wal get error", zap.Error(err), zap.String("key", w.key))
				return
			}
			// 打印扫描到的键值对
			for _, item := range keys {
				if item == "_id" || len(item) <= 16 {
					continue
				}
				value, err := w.conn.HGet(ctx, w.key, item).Result()
				if err != nil {
					w.logger.Error("redis wal get error", zap.Error(err), zap.String("key", w.key))
					return
				}
				tmp = value[:8]
				receiveTime, _ = strconv.ParseInt(tmp, 10, 64)
				if now-receiveTime > t {
					_id, _ = strconv.ParseInt(item, 10, 64)
					tmp = value[8:16]
					stage, _ = strconv.ParseInt(tmp, 10, 64)
					result <- dataInfo{
						id:    _id,
						data:  value[16:],
						stage: stage,
					}
				}
			}
			// 游标为0表示迭代结束
			if cursor == 0 {
				break
			}
		}
	})
	return result, err
}

func NewWalForTimeout(conn redis.Cmdable, key, receiveKey, stageKey string, logger *zap.Logger) *WalForTimeout {
	return &WalForTimeout{
		conn:       conn,
		key:        key,
		receiveKey: receiveKey,
		stageKey:   stageKey,
		lock:       NewRedisLock(conn, key, "redis_wal_lock:"+key, 20),
		ttl:        time.Hour * 24 * 30,
		logger:     logger,
	}
}
