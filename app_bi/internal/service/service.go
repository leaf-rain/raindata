package service

import (
	"context"
	"github.com/google/wire"
	"github.com/leaf-rain/raindata/app_bi/internal/biz"
	"github.com/leaf-rain/raindata/app_bi/internal/data/data"
	"go.uber.org/zap"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewService)

type Service struct {
	Ctx    context.Context
	biz    *biz.Business
	data   *data.Data
	logger *zap.Logger
}

func NewService(ctx context.Context, biz *biz.Business, data *data.Data, logger *zap.Logger) *Service {
	svc := &Service{
		Ctx:    ctx,
		biz:    biz,
		data:   data,
		logger: logger,
	}
	return svc
}
