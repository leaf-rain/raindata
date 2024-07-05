package clickhouse_sqlx

import (
	"context"
	"github.com/leaf-rain/raindata/common/snowflake"
)

type WALForTimeout interface {
	Put(ctx context.Context, subkey, data string, stage int) (int64, error)
	Del(ctx context.Context, subkey string) error
	GetForTimeout(ctx context.Context, t int) (chan string, error)
}

var _ WALForTimeout = (*NotWAL)(nil)

// NotWAL 不使用WAL
type NotWAL struct {
}

func (n NotWAL) Put(ctx context.Context, subkey, data string, stage int) (int64, error) {
	return snowflake.SnowflakeInt64(), nil
}

func (n NotWAL) Del(ctx context.Context, subkey string) error {
	return nil
}

func (n NotWAL) GetForTimeout(ctx context.Context, t int) (chan string, error) {
	return nil, nil
}
