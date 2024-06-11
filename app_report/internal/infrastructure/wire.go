//go:build wireinject
// +build wireinject

package infrastructure

import (
	"github.com/google/wire"
	"github.com/leaf-rain/raindata/app_report/internal/infrastructure/config"
	"github.com/leaf-rain/raindata/app_report/pkg/logger"
)

var WireInfrastructureSet = wire.NewSet(
	config.NewCmdArgs,
	config.InitConfig,
	config.GetLogCfgByConfig,
	logger.InitLogger,
	NewInfrastructure,
)

func Initialize() (*Infrastructure, error) {
	wire.Build(WireInfrastructureSet)
	return &Infrastructure{}, nil
}
