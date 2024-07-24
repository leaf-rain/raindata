//go:build wireinject
// +build wireinject

package main

//go:generate wire

import (
	"github.com/google/wire"
	"github.com/leaf-rain/raindata/app_report/internal/adapter"
	"github.com/leaf-rain/raindata/app_report/internal/adapter/interface_app"
	"github.com/leaf-rain/raindata/app_report/internal/application"
	"github.com/leaf-rain/raindata/app_report/internal/application/interface_domain"
	"github.com/leaf-rain/raindata/app_report/internal/domain"
	"github.com/leaf-rain/raindata/app_report/internal/domain/interface_repo"
	"github.com/leaf-rain/raindata/app_report/internal/infrastructure"
	"github.com/leaf-rain/raindata/app_report/internal/infrastructure/repository"
)

func Initialize() (*adapter.Adapter, error) {
	wire.Build(
		wire.Bind(new(interface_app.InterfaceStream), new(*application.AppStream)),
		wire.Bind(new(interface_domain.InterfaceWriter), new(*domain.Writer)),
		wire.Bind(new(interface_repo.InterfaceMetadataRepo), new(*interface_repo.DefaultMetadata)),
		wire.Bind(new(interface_repo.InterfaceWriterRepo), new(*repository.SRWriter)),
		wire.Bind(new(interface_repo.InterfaceIdRepo), new(*interface_repo.SnowflakeId)),
		infrastructure.WireInfrastructureSet,
		domain.WireDomainSet,
		application.WireApplicationSet,
		adapter.WireAdapterSet,
	)
	return &adapter.Adapter{}, nil
}
