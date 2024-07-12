// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package domain

import (
	"github.com/google/wire"
	"github.com/leaf-rain/raindata/app_report/internal/domain/interface_repo"
	"github.com/leaf-rain/raindata/app_report/internal/infrastructure/config"
	"github.com/leaf-rain/raindata/app_report/internal/infrastructure/repository"
	"github.com/leaf-rain/raindata/app_report/pkg/logger"
	"github.com/leaf-rain/raindata/common/clickhouse_sqlx"
)

// Injectors from wire.go:

func Initialize() (*Domain, error) {
	cmdArgs := config.NewCmdArgs()
	configConfig, err := config.InitConfig(cmdArgs)
	if err != nil {
		return nil, err
	}
	logConfig := config.GetLogCfgByConfig(configConfig)
	zapLogger, err := logger.InitLogger(logConfig)
	if err != nil {
		return nil, err
	}
	defaultMetadata := interface_repo.NewMetadata()
	clickhouseConfig := config.GetCKCfgByConfig(configConfig)
	clickhouseCluster, err := clickhouse_sqlx.InitClusterConn(clickhouseConfig)
	if err != nil {
		return nil, err
	}
	ckWriter := repository.NewCkWriter(zapLogger, clickhouseCluster)
	snowflakeId := interface_repo.NewSnowflakeId()
	domain := NewDomain(zapLogger, defaultMetadata, ckWriter, snowflakeId)
	return domain, nil
}

// wire.go:

var WireDomainSet = wire.NewSet(interface_repo.NewSnowflakeId, interface_repo.NewMetadata, NewDomain,
	NewCkWriter,
)
