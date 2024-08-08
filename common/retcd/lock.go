package retcd

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

// 锁的简单封装
type Lock struct {
	key     string
	value   string
	ttl     int64
	ctx     context.Context
	mutex   *concurrency.Mutex
	session *concurrency.Session
	cli     *clientv3.Client
}

func NewLock(ctx context.Context, cli *clientv3.Client, key string, ttl int64) (*Lock, error) {
	if cli == nil {
		return nil, nil
	}
	lock := &Lock{
		key: key,
		ctx: ctx,
		cli: cli,
		ttl: ttl,
	}
	return lock, nil
}

func (l *Lock) initSession() error {
	var err error
	l.session, err = concurrency.NewSession(l.cli, concurrency.WithTTL(int(l.ttl)))
	if err != nil {
		return err
	}
	l.mutex = concurrency.NewMutex(l.session, l.key)
	return nil
}

func (l *Lock) TryLock() error {
	if l.session == nil || l.mutex == nil {
		err := l.initSession()
		if err != nil {
			return err
		}
	}
	return l.mutex.TryLock(l.ctx)
}

func (l *Lock) Lock() error {
	if l.session == nil || l.mutex == nil {
		err := l.initSession()
		if err != nil {
			return err
		}
	}
	return l.mutex.Lock(l.ctx)
}

func (l *Lock) Unlock() error {
	err := l.mutex.Unlock(l.ctx)
	l.session.Close()
	l.session = nil
	l.mutex = nil
	return err
}
