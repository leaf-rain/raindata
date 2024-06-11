//go:build wireinject
// +build wireinject

package infrastructure

import (
	"github.com/google/wire"
	"github.com/leaf-rain/raindata/app_basicsdata/internal/infrastructure/config"
	"github.com/leaf-rain/raindata/app_basicsdata/internal/infrastructure/entity"
	"github.com/leaf-rain/raindata/app_report/pkg/logger"
	"github.com/leaf-rain/raindata/common/clickhouse_sqlx"
	"github.com/leaf-rain/raindata/common/etcd"
)

var WireInfrastructureSet = wire.NewSet(
	config.NewCmdArgs,
	config.InitConfig,
	config.GetLogCfgByConfig,
	config.GetCtx,
	config.GetEtcdConfig,
	config.GetClickhouseConfig,
	logger.InitLogger,
	etcd.NewEtcdClient,
	clickhouse_sqlx.NewClickhouse,
	entity.NewRepository,
	NewInfrastructure,
)

func Initialize() (*Infrastructure, error) {
	wire.Build(WireInfrastructureSet)
	return &Infrastructure{}, nil
}
