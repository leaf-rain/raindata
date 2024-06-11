//go:build wireinject
// +build wireinject

package application

import (
	"github.com/google/wire"
	"github.com/leaf-rain/raindata/app_basicsdata/internal/application/interface_domain"
	"github.com/leaf-rain/raindata/app_basicsdata/internal/domain"
	"github.com/leaf-rain/raindata/app_basicsdata/internal/infrastructure"
)

var WireApplicationSet = wire.NewSet(
	NewAppMetadata,
	NewApplications,
)

func Initialize() (*Applications, error) {
	wire.Build(
		wire.Bind(new(interface_domain.InterfaceDomain), new(*domain.Metadata)),
		infrastructure.WireInfrastructureSet,
		domain.WireDomainSet,
		WireApplicationSet,
	)
	return &Applications{}, nil
}
