package rredis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type Client struct {
	redis.Cmdable
}

type Config struct {
	PoolSize     int      `yaml:"poolSize"`
	Addr         []string `yaml:"addr"`
	Pwd          string   `yaml:"pwd"`
	DialTimeout  int64    `yaml:"dialTimeout"`
	ReadTimeout  int64    `yaml:"readTimeout"`
	WriteTimeout int64    `yaml:"writeTimeout"`
	DB           int      `yaml:"db"`
}

func NewRedis(o Config, ctx context.Context) (client *Client, err error) {
	var redisCli redis.Cmdable
	if len(o.Addr) > 1 {
		redisCli = redis.NewClusterClient(
			&redis.ClusterOptions{
				Addrs:        o.Addr,
				PoolSize:     o.PoolSize,
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
				PoolSize:     o.PoolSize,
				DB:           o.DB,
			},
		)
	}
	err = redisCli.Ping(ctx).Err()
	if nil != err {
		panic(err)
	}

	client = new(Client)
	client.Cmdable = redisCli
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
