//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/leaf-rain/raindata/app_bi/internal/biz"
	"github.com/leaf-rain/raindata/app_bi/internal/conf"
	"github.com/leaf-rain/raindata/app_bi/internal/data/entity"
	"github.com/leaf-rain/raindata/app_bi/internal/server"
	"github.com/leaf-rain/raindata/app_bi/internal/service"
	"go.uber.org/zap"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, *zap.Logger, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, entity.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
