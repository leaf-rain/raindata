// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/leaf-rain/raindata/app_basicsdata/internal/adapter"
	"github.com/leaf-rain/raindata/app_basicsdata/internal/application"
	"github.com/leaf-rain/raindata/app_basicsdata/internal/domain"
	"github.com/leaf-rain/raindata/app_basicsdata/internal/infrastructure/config"
	"github.com/leaf-rain/raindata/app_basicsdata/internal/infrastructure/entity"
	"github.com/leaf-rain/raindata/app_report/pkg/logger"
	"github.com/leaf-rain/raindata/common/clickhouse_sqlx"
	"github.com/leaf-rain/raindata/common/etcd"
)

import (
	_ "embed"
	_ "net/http/pprof"
)

// Injectors from wire.go:

func Initialize() (*adapter.Adapter, error) {
	cmdArgs := config.NewCmdArgs()
	configConfig, err := config.InitConfig(cmdArgs)
	if err != nil {
		return nil, err
	}
	logConfig := config.GetLogCfgByConfig(configConfig)
	zapLogger := logger.InitLogger(logConfig)
	context := config.GetCtx(configConfig)
	etcdConfig := config.GetEtcdConfig(configConfig)
	client, err := etcd.NewEtcdClient(etcdConfig, zapLogger)
	if err != nil {
		return nil, err
	}
	clickhouseConfig := config.GetClickhouseConfig(configConfig)
	clickhouse, err := clickhouse_sqlx.NewClickhouse(clickhouseConfig)
	if err != nil {
		return nil, err
	}
	repository, err := entity.NewRepository(context, client, clickhouse, zapLogger, configConfig)
	if err != nil {
		return nil, err
	}
	metadata, err := domain.NewMetadata(zapLogger, repository, configConfig, client, clickhouse)
	if err != nil {
		return nil, err
	}
	appMetadata := application.NewAppMetadata(zapLogger, configConfig, metadata)
	grpcServer := adapter.NewGrpcServer(configConfig, zapLogger, appMetadata)
	adapterAdapter := adapter.NewAdapter(grpcServer)
	return adapterAdapter, nil
}
