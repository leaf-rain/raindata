package interface_repo

import (
	"context"
	"github.com/leaf-rain/raindata/common/snowflake"
)

type InterfaceIdRepo interface {
	InitId(ctx context.Context, appid int64) (int64, error)
}

type SnowflakeId struct{}

func NewSnowflakeId() *SnowflakeId {
	return &SnowflakeId{}
}

func (s SnowflakeId) InitId(ctx context.Context, appid int64) (int64, error) {
	return snowflake.SnowflakeInt64(), nil
}
