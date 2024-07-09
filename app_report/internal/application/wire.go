//go:build wireinject
// +build wireinject

package application

import (
	"github.com/google/wire"
	"github.com/leaf-rain/raindata/app_report/internal/application/interface_domain"
	"github.com/leaf-rain/raindata/app_report/internal/domain"
	"github.com/leaf-rain/raindata/app_report/internal/domain/interface_repo"
	"github.com/leaf-rain/raindata/app_report/internal/infrastructure"
)

var WireApplicationSet = wire.NewSet(
	NewAppStream,
	NewApplications,
)

func Initialize() (*Applications, error) {
	wire.Build(
		wire.Bind(new(interface_repo.InterfaceEventManager), new(*domain.EventManager)),
		wire.Bind(new(interface_domain.InterfaceWriter), new(*domain.CkWriter)),
		infrastructure.WireInfrastructureSet,
		domain.WireDomainSet,
		WireApplicationSet,
	)
	return &Applications{}, nil
}
