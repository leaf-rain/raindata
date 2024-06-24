package go_redis

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/leaf-rain/fastjson"
	"time"
)

var (
	ErrBusinessStaring = errors.New("business staring")
)

type WalForTimeout struct {
	conn       redis.Cmdable
	key        string
	receiveKey string
	stageKey   string
	parserPool fastjson.ParserPool
	arenaPool  fastjson.ArenaPool
	lock       *RedisLock
	ttl        time.Duration
}

func (w *WalForTimeout) Put(ctx context.Context, subkey, data string, stage int) error {
	p := w.parserPool.Get()
	v, err := p.Parse(data)
	if err != nil {
		return err
	}
	if v.GetInt64(w.receiveKey) == 0 {
		a := w.arenaPool.Get()
		v.Set(w.receiveKey, a.NewNumberInt(int(time.Now().Unix())))
		w.arenaPool.Put(a)
	}
	if v.GetInt64(w.stageKey) != 0 {
		a := w.arenaPool.Get()
		v.Set(w.stageKey, a.NewNumberInt(stage))
		w.arenaPool.Put(a)
	}
	err = w.conn.HSet(ctx, w.key, subkey, v.String()).Err()
	w.conn.Expire(ctx, w.key, w.ttl) // 避免无限制存储
	return err
}

func (w *WalForTimeout) Del(ctx context.Context, subkey string) (err error) {
	err = w.conn.HDel(ctx, w.key, subkey).Err()
	w.conn.Expire(ctx, w.key, w.ttl) // 避免无限制存储
	return err
}

func (w *WalForTimeout) GetForTimeout(ctx context.Context, t int) ([]string, error) {
	lockResult, err := w.lock.TryLock(ctx)
	defer w.lock.UnLock(ctx)
	if err != nil {
		return nil, err
	}
	if !lockResult {
		return nil, ErrBusinessStaring
	}
	var cursor uint64
	var keys []string
	var value string
	var receiveTime int
	var now = int(time.Now().Unix())
	var result []string
	for {
		cmd := w.conn.HScan(ctx, w.key, cursor, "*", 200)
		// 处理命令结果
		keys, cursor, err = cmd.Result()
		if err != nil {
			return nil, err
		}
		// 打印扫描到的键值对
		for _, item := range keys {
			value, err = w.conn.HGet(ctx, w.key, item).Result()
			if err != nil {
				return nil, err
			}
			receiveTime = fastjson.GetInt([]byte(value), w.receiveKey)
			if now-receiveTime > t {
				result = append(result, value)
			}
		}
		// 游标为0表示迭代结束
		if cursor == 0 {
			break
		}
	}
	return result, err
}

func NewWalForTimeout(conn redis.Cmdable, key, receiveKey, stageKey string) *WalForTimeout {
	return &WalForTimeout{
		conn:       conn,
		key:        key,
		receiveKey: receiveKey,
		stageKey:   stageKey,
		parserPool: fastjson.ParserPool{},
		arenaPool:  fastjson.ArenaPool{},
		lock:       NewRedisLock(conn, key, "redis_wal_lock:"+key, 20),
		ttl:        time.Hour * 24 * 30,
	}
}
