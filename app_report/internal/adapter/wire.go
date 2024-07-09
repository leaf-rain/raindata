//go:build wireinject
// +build wireinject

package adapter

import (
	"github.com/google/wire"
	"github.com/leaf-rain/raindata/app_report/internal/adapter/interface_app"
	"github.com/leaf-rain/raindata/app_report/internal/application"
	"github.com/leaf-rain/raindata/app_report/internal/application/interface_domain"
	"github.com/leaf-rain/raindata/app_report/internal/domain"
	"github.com/leaf-rain/raindata/app_report/internal/domain/interface_repo"
	"github.com/leaf-rain/raindata/app_report/internal/infrastructure"
)

var WireAdapterSet = wire.NewSet(
	NewGrpcServer,
	NewAdapter,
)

func Initialize() (*Adapter, error) {
	wire.Build(
		wire.Bind(new(interface_app.InterfaceStream), new(*application.AppStream)),
		wire.Bind(new(interface_repo.InterfaceEventManager), new(*domain.EventManager)),
		wire.Bind(new(interface_domain.InterfaceWriter), new(*domain.CkWriter)),
		infrastructure.WireInfrastructureSet,
		domain.WireDomainSet,
		application.WireApplicationSet,
		WireAdapterSet,
	)
	return &Adapter{}, nil
}
