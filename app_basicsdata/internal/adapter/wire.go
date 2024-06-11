//go:build wireinject
// +build wireinject

package adapter

import (
	"github.com/google/wire"
	"github.com/leaf-rain/raindata/app_basicsdata/internal/adapter/interface_app"
	"github.com/leaf-rain/raindata/app_basicsdata/internal/application"
	"github.com/leaf-rain/raindata/app_basicsdata/internal/application/interface_domain"
	"github.com/leaf-rain/raindata/app_basicsdata/internal/domain"
	"github.com/leaf-rain/raindata/app_basicsdata/internal/infrastructure"
)

var WireAdapterSet = wire.NewSet(
	NewGrpcServer,
	NewAdapter,
)

func Initialize() (*Adapter, error) {
	wire.Build(
		wire.Bind(new(interface_domain.InterfaceDomain), new(*domain.Metadata)),
		wire.Bind(new(interface_app.InterfaceApplication), new(*application.AppMetadata)),
		infrastructure.WireInfrastructureSet,
		domain.WireDomainSet,
		application.WireApplicationSet,
		WireAdapterSet,
	)
	return &Adapter{}, nil
}
