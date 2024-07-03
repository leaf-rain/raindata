package clickhouse_sqlx

import "context"

type WALForTimeout interface {
	Put(ctx context.Context, subkey, data string, stage int) error
	Del(ctx context.Context, subkey string) error
	GetForTimeout(ctx context.Context, t int) ([]string, error)
}

var _ WALForTimeout = (*NotWAL)(nil)

// NotWAL 不使用WAL
type NotWAL struct {
}

func (n NotWAL) Put(ctx context.Context, subkey, data string, stage int) error {
	return nil
}

func (n NotWAL) Del(ctx context.Context, subkey string) error {
	return nil
}

func (n NotWAL) GetForTimeout(ctx context.Context, t int) ([]string, error) {
	return nil, nil
}
