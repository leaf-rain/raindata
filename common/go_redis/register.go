package go_redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type RedisDiscoveryConfig struct {
	HartTime int64  `json:"HartTime,omitempty" yaml:"HartTime"`
	Key      string `json:"Key,omitempty" yaml:"Key"`
	Subkey   string `json:"Subkey,omitempty" yaml:"Subkey"`
	Value    string `json:"Value,omitempty" yaml:"Value"`
}

type RedisDiscovery struct {
	hartTime int64
	key      string
	subkey   string
	value    string
	ctx      context.Context
	conn     redis.Cmdable
	ticker   *time.Ticker
	logger   *zap.Logger
}

func NewRedisDiscovery(ctx context.Context, conn redis.Cmdable, logger *zap.Logger, conf RedisDiscoveryConfig) *RedisDiscovery {
	result := &RedisDiscovery{
		ctx:      ctx,
		conn:     conn,
		hartTime: conf.HartTime,
		key:      conf.Key,
		subkey:   conf.Subkey,
		value:    conf.Value,
		logger:   logger,
	}
	result.Init()
	return result
}

func (rd *RedisDiscovery) Init() {
	if rd.ticker != nil {
		rd.ticker.Stop()
	}
	rd.ticker = time.NewTicker(time.Duration(rd.hartTime) * time.Second)
	go func() {
		var nowStr, subkey string
		var count int64
		var cursor uint64
		for range rd.ticker.C {
			now := time.Now().Unix()
			nowStr = strconv.FormatInt(time.Now().Unix(), 10)
			subkey = nowStr + rd.subkey
			rd.conn.HSet(rd.ctx, rd.key, subkey, time.Now().Unix())
			rd.conn.Expire(rd.ctx, rd.key, time.Duration(rd.hartTime)*time.Second*5)
			count += 1
			if count%2 == 0 {
				count = 0
				// 清除哈希中过期key
				var keys []string
				var err error
				var value int64
				for {
					cmd := rd.conn.HScan(rd.ctx, rd.key, cursor, "*", 200)
					// 处理命令结果
					keys, cursor, err = cmd.Result()
					if err != nil {
						rd.logger.Error("[RedisDiscovery] conn.Hscan failed.", zap.Error(err), zap.String("key", rd.key), zap.String("subkey", subkey))
						continue
					}
					// 打印扫描到的键值对
					for _, item := range keys {
						value, _ = removePrefixDigitsAndUnderscores(item)
						if value < now-(rd.hartTime*2) {
							rd.conn.HDel(rd.ctx, rd.key, item)
						}
					}
					// 游标为0表示迭代结束
					if cursor == 0 {
						break
					}
				}
			}
		}
	}()
}

func removePrefixDigitsAndUnderscores(s string) (int64, string) {
	i := 0
	for ; i < len(s); i++ {
		if s[i] == '_' {
			break
		}
	}
	num := s[:i]
	t, _ := strconv.ParseInt(num, 10, 64)
	return t, s[i:]
}

func (rd *RedisDiscovery) Close() {
	rd.ticker.Stop()
}

func (rd *RedisDiscovery) SetValue(value string) {
	rd.value = value
}

func (rd *RedisDiscovery) GetAllValue() []string {
	var keys []string
	var err error
	var cursor uint64
	var result []string
	for {
		cmd := rd.conn.HScan(rd.ctx, rd.key, cursor, "*", 200)
		// 处理命令结果
		keys, cursor, err = cmd.Result()
		if err != nil {
			rd.logger.Error("[RedisDiscovery] conn.Hscan failed.", zap.Error(err), zap.String("key", rd.key), zap.String("subkey", rd.subkey))
			continue
		}
		// 打印扫描到的键值对
		for _, item := range keys {
			value, _ := rd.conn.HGet(rd.ctx, rd.key, item).Result()
			if value != "" {
				result = append(result, value)
			}
		}
		// 游标为0表示迭代结束
		if cursor == 0 {
			break
		}
	}
	return result
}
