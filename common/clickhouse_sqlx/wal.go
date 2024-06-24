package clickhouse_sqlx

import "context"

type WALForTimeout interface {
	Put(ctx context.Context, subkey, data string, stage int) error
	Del(ctx context.Context, subkey string) error
	GetForTimeout(ctx context.Context, t int) ([]string, error)
}
