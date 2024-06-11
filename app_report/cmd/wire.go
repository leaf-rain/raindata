//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/leaf-rain/raindata/app_report/internal/adapter"
	"github.com/leaf-rain/raindata/app_report/internal/adapter/interface_app"
	"github.com/leaf-rain/raindata/app_report/internal/application"
	"github.com/leaf-rain/raindata/app_report/internal/application/interface_domain"
	"github.com/leaf-rain/raindata/app_report/internal/domain"
	"github.com/leaf-rain/raindata/app_report/internal/infrastructure"
)

func Initialize() (*adapter.Adapter, error) {
	wire.Build(
		wire.Bind(new(interface_app.InterfaceStream), new(*application.AppStream)),
		wire.Bind(new(interface_domain.InterfaceEventManager), new(*domain.EventManager)),
		wire.Bind(new(interface_domain.InterfaceWriter), new(*domain.CkWriter)),
		infrastructure.WireInfrastructureSet,
		domain.WireDomainSet,
		application.WireApplicationSet,
		adapter.WireAdapterSet,
	)
	return &adapter.Adapter{}, nil
}
