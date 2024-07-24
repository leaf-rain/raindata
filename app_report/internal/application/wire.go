//go:build wireinject
// +build wireinject

package application

//go:generate wire

import (
	"github.com/google/wire"
	"github.com/leaf-rain/raindata/app_report/internal/application/interface_domain"
	"github.com/leaf-rain/raindata/app_report/internal/domain"
	"github.com/leaf-rain/raindata/app_report/internal/domain/interface_repo"
	"github.com/leaf-rain/raindata/app_report/internal/infrastructure"
	"github.com/leaf-rain/raindata/app_report/internal/infrastructure/repository"
)

var WireApplicationSet = wire.NewSet(
	NewApplications,
	NewAppStream,
)

func Initialize() (*Applications, error) {
	wire.Build(
		wire.Bind(new(interface_domain.InterfaceWriter), new(*domain.Writer)),
		wire.Bind(new(interface_repo.InterfaceMetadataRepo), new(*interface_repo.DefaultMetadata)),
		wire.Bind(new(interface_repo.InterfaceWriterRepo), new(*repository.SRWriter)),
		wire.Bind(new(interface_repo.InterfaceIdRepo), new(*interface_repo.SnowflakeId)),
		infrastructure.WireInfrastructureSet,
		domain.WireDomainSet,
		WireApplicationSet,
	)
	return &Applications{}, nil
}
