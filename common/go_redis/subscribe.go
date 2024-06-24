package go_redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type RedisSubscribe struct {
	conn *redis.Client
}

func NewRedisSubscribe(conn *redis.Client) *RedisSubscribe {
	return &RedisSubscribe{conn: conn}
}

func (r *RedisSubscribe) PubMessage(ctx context.Context, channel, msg string) {
	r.conn.Publish(ctx, channel, msg)
}

func (r *RedisSubscribe) SubMessage(ctx context.Context, channel, msg string) {
	pubsub := r.conn.Subscribe(ctx, channel)
	_, err := pubsub.Receive(ctx)
	if err != nil {
		panic(err)
	}
	ch := pubsub.Channel()
	for msg := range ch {
		fmt.Println(msg.Channel, msg.Payload)
	}
}
