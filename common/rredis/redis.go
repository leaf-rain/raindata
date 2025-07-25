package rredis

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"strconv"
	"strings"
	"sync"
	"time"
)

var ProviderSet = wire.NewSet(NewRedis)

var clientMap = sync.Map{}

type Client struct {
	redis.Cmdable
}

func NewRedis(o *RedisCfg, ctx context.Context) (client *Client, err error) {
	if o == nil || len(o.Addr) == 0 {
		return nil, errors.New("redis config error")
	}

	key := strings.Join(o.Addr, "_") + "_" + strconv.FormatInt(o.DB, 10)
	c, _ := clientMap.Load(key)
	if c != nil {
		client = c.(*Client)
		return
	}
	var redisCli redis.Cmdable
	if len(o.Addr) > 1 {
		redisCli = redis.NewClusterClient(
			&redis.ClusterOptions{
				Addrs:        o.Addr,
				PoolSize:     int(o.PoolSize),
				DialTimeout:  time.Second * time.Duration(o.DialTimeout),
				ReadTimeout:  time.Second * time.Duration(o.ReadTimeout),
				WriteTimeout: time.Second * time.Duration(o.WriteTimeout),
				Password:     o.Pwd,
			},
		)
	} else {
		redisCli = redis.NewClient(
			&redis.Options{
				Addr:         o.Addr[0],
				DialTimeout:  time.Second * time.Duration(o.DialTimeout),
				ReadTimeout:  time.Second * time.Duration(o.ReadTimeout),
				WriteTimeout: time.Second * time.Duration(o.WriteTimeout),
				Password:     o.Pwd,
				PoolSize:     int(o.PoolSize),
				DB:           int(o.DB),
			},
		)
	}
	err = redisCli.Ping(ctx).Err()
	if nil != err {
		return nil, err
	}

	client = new(Client)
	client.Cmdable = redisCli
	clientMap.Store(key, client)
	return client, nil
}

func (c *Client) Process(cmd redis.Cmder) error {
	switch redisCli := c.Cmdable.(type) {
	case *redis.ClusterClient:
		return redisCli.Process(context.TODO(), cmd)
	case *redis.Client:
		return redisCli.Process(context.TODO(), cmd)
	default:
		return nil
	}
}

func (c *Client) Close() error {
	switch redisCli := c.Cmdable.(type) {
	case *redis.ClusterClient:
		return redisCli.Close()
	case *redis.Client:
		return redisCli.Close()
	default:
		return nil
	}
}
