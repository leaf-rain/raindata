//go:build wireinject
// +build wireinject

package domain

//go:generate wire

import (
	"github.com/google/wire"
	"github.com/leaf-rain/raindata/app_report/internal/domain/interface_repo"
	"github.com/leaf-rain/raindata/app_report/internal/infrastructure"
	"github.com/leaf-rain/raindata/app_report/internal/infrastructure/repository"
)

var WireDomainSet = wire.NewSet(
	interface_repo.NewSnowflakeId,
	interface_repo.NewMetadata,
	NewDomain,
	NewCkWriter,
)

func Initialize() (*Domain, error) {
	wire.Build(
		wire.Bind(new(interface_repo.InterfaceMetadataRepo), new(*interface_repo.DefaultMetadata)),
		wire.Bind(new(interface_repo.InterfaceWriterRepo), new(*repository.SRWriter)),
		wire.Bind(new(interface_repo.InterfaceIdRepo), new(*interface_repo.SnowflakeId)),
		infrastructure.WireInfrastructureSet,
		WireDomainSet,
	)
	return &Domain{}, nil
}
