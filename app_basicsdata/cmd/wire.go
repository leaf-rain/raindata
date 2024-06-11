//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/leaf-rain/raindata/app_basicsdata/internal/adapter"
	"github.com/leaf-rain/raindata/app_basicsdata/internal/adapter/interface_app"
	"github.com/leaf-rain/raindata/app_basicsdata/internal/application"
	"github.com/leaf-rain/raindata/app_basicsdata/internal/application/interface_domain"
	"github.com/leaf-rain/raindata/app_basicsdata/internal/domain"
	"github.com/leaf-rain/raindata/app_basicsdata/internal/infrastructure"
)

func Initialize() (*adapter.Adapter, error) {
	wire.Build(
		wire.Bind(new(interface_domain.InterfaceDomain), new(*domain.Metadata)),
		wire.Bind(new(interface_app.InterfaceApplication), new(*application.AppMetadata)),
		infrastructure.WireInfrastructureSet,
		domain.WireDomainSet,
		application.WireApplicationSet,
		adapter.WireAdapterSet,
	)
	return &adapter.Adapter{}, nil
}
