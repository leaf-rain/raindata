package rredis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisLock struct {
	conn    redis.Cmdable
	timeout time.Duration
	key     string
	val     string
}

func NewRedisLock(conn redis.Cmdable, key, val string, timeout time.Duration) *RedisLock {
	return &RedisLock{conn: conn, timeout: timeout, key: "lock:" + key, val: val}
}

// TryLock return true ===> Get the lock successfully
func (lock *RedisLock) TryLock(ctx context.Context) (bool, error) {
	return lock.conn.SetNX(ctx, lock.key, lock.val, lock.timeout).Result()
}

// redis的lock没有自旋功能
func (lock *RedisLock) Lock(ctx context.Context) error {
	return lock.conn.Set(ctx, lock.key, lock.val, lock.timeout).Err()
}

func (lock *RedisLock) UnLock(ctx context.Context) error {
	luaDel := redis.NewScript("if redis.call('get',KEYS[1]) == ARGV[1] then " +
		"return redis.call('del',KEYS[1]) else return 0 end")
	return luaDel.Run(ctx, lock.conn, []string{lock.key}, lock.val).Err()
}

func (lock *RedisLock) GetLockKey() string {
	return lock.key
}

func (lock *RedisLock) GetLockVal() string {
	return lock.val
}
