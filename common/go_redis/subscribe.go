package go_redis

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type RedisSubscribe struct {
	ctx        context.Context
	conn       *redis.Client
	channel    string
	msgChannel chan string
	cancel     func()
}

func NewRedisSubscribe(conn *redis.Client, channel string, msgChannel chan string) *RedisSubscribe {
	ctx, cancel := context.WithCancel(context.TODO())
	return &RedisSubscribe{ctx: ctx, cancel: cancel, conn: conn, channel: channel, msgChannel: msgChannel}
}

func (r *RedisSubscribe) PubMessage(ctx context.Context, msg string) {
	r.conn.Publish(ctx, r.channel, msg)
}

func (r *RedisSubscribe) SubMessage(ctx context.Context) {
	pubsub := r.conn.Subscribe(ctx, r.channel)
	_, err := pubsub.Receive(ctx)
	if err != nil {
		panic(err)
	}
	ch := pubsub.Channel()
	go func() {
		for {
			select {
			case <-ctx.Done():
				break
			case msg := <-ch:
				r.msgChannel <- msg.Payload
			}
		}
	}()
}

func (r *RedisSubscribe) Close() {
	r.cancel()
}
