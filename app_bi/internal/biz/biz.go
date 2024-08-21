package biz

import (
	"context"
	"github.com/google/wire"
	"github.com/leaf-rain/raindata/app_bi/internal/data/data"
	"go.uber.org/zap"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewBusiness)

type Business struct {
	Ctx    context.Context
	data   *data.Data
	logger *zap.Logger
}

func NewBusiness(ctx context.Context, data *data.Data, logger *zap.Logger) *Business {
	return &Business{
		Ctx:    ctx,
		data:   data,
		logger: logger,
	}
}
