package rclickhouse

import (
	"context"
	"github.com/leaf-rain/raindata/common/snowflake"
)

type Data interface {
	GetId() int64
	GetData() string
	GetStage() int64
}

type WALForTimeout interface {
	Put(ctx context.Context, data Data) (int64, error)
	Del(ctx context.Context, id int64) (err error)
	GetForTimeout(ctx context.Context, t int64) (chan Data, error)
}

var _ WALForTimeout = (*NotWAL)(nil)

// NotWAL 不使用WAL
type NotWAL struct {
}

func (n NotWAL) Put(ctx context.Context, data Data) (int64, error) {
	return snowflake.SnowflakeInt64(), nil
}

func (n NotWAL) Del(ctx context.Context, id int64) (err error) {
	return nil
}

func (n NotWAL) GetForTimeout(ctx context.Context, t int64) (chan Data, error) {
	return nil, nil
}
