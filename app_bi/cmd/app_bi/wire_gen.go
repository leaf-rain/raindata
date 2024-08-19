// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/leaf-rain/raindata/app_bi/internal/biz"
	"github.com/leaf-rain/raindata/app_bi/internal/conf"
	"github.com/leaf-rain/raindata/app_bi/internal/data/entity"
	"github.com/leaf-rain/raindata/app_bi/internal/server"
	"github.com/leaf-rain/raindata/app_bi/internal/service"
	"go.uber.org/zap"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, confData *conf.Data, logger *zap.Logger, logLogger log.Logger) (*kratos.App, func(), error) {
	dataData, cleanup, err := entity.NewData(confData, logger)
	if err != nil {
		return nil, nil, err
	}
	greeterRepo := entity.NewGreeterRepo(dataData, logger)
	greeterUsecase := biz.NewGreeterUsecase(greeterRepo, logger)
	greeterService := service.NewGreeterService(greeterUsecase, logger)
	grpcServer := server.NewGRPCServer(confServer, greeterService, logger)
	httpServer := server.NewHTTPServer(confServer, logger, greeterService)
	app := newApp(logLogger, grpcServer, httpServer)
	return app, func() {
		cleanup()
	}, nil
}
