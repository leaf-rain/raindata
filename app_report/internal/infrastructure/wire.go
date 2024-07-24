//go:build wireinject
// +build wireinject

package infrastructure

//go:generate wire

import (
	"github.com/google/wire"
	"github.com/leaf-rain/raindata/app_report/internal/infrastructure/config"
	"github.com/leaf-rain/raindata/app_report/internal/infrastructure/repository"
	"github.com/leaf-rain/raindata/app_report/pkg/logger"
	"github.com/leaf-rain/raindata/common/go_sql_driver"
)

var WireInfrastructureSet = wire.NewSet(
	config.NewCmdArgs,
	config.InitConfig,
	config.GetLogCfgByConfig,
	//config.GetCKCfgByConfig,
	//clickhouse_sqlx.InitClusterConn,
	config.GetSqlCfgByConfig,
	go_sql_driver.NewSql,
	logger.InitLogger,
	//repository.NewCkWriter,
	repository.NewSRWriter,
	NewInfrastructure,
)

func Initialize() (*Infrastructure, error) {
	wire.Build(WireInfrastructureSet)
	return &Infrastructure{}, nil
}
