package application

import (
	"context"
	pb_metadata "github.com/leaf-rain/raindata/app_basicsdata/api/grpc"
	"go.uber.org/zap"

	"github.com/leaf-rain/raindata/app_basicsdata/internal/application/interface_domain"
	"github.com/leaf-rain/raindata/app_basicsdata/internal/infrastructure/config"
)

type AppMetadata struct {
	logger   *zap.Logger
	config   *config.Config
	metadata interface_domain.InterfaceDomain
}

func NewAppMetadata(logger *zap.Logger, config *config.Config, metadata interface_domain.InterfaceDomain) *AppMetadata {
	return &AppMetadata{
		logger:   logger,
		config:   config,
		metadata: metadata,
	}
}

func (app *AppMetadata) GetMetadata(ctx context.Context, in *pb_metadata.GetMetadataRequest) (*pb_metadata.MetadataResponse, error) {
	// 事件属性管理器, 获取字段字段值，以及表名称
	result, err := app.metadata.GetMetadata(ctx, in)
	if err != nil {
		app.logger.Error("[GetMetadata] app metadata.StorageMetadata failed", zap.String("fields", in.EventName), zap.Error(err))
	}
	return result, err
}

func (app *AppMetadata) PutMetadata(ctx context.Context, in *pb_metadata.MetadataRequest) (*pb_metadata.MetadataResponse, error) {
	return nil, nil
}
