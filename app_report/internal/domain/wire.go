//go:build wireinject
// +build wireinject

package domain

import (
	"github.com/google/wire"
	"github.com/leaf-rain/raindata/app_report/internal/infrastructure"
)

var WireDomainSet = wire.NewSet(
	NewCkWriter,
	NewEventManager,
	NewDomain,
)

func Initialize() (*Domain, error) {
	wire.Build(
		infrastructure.WireInfrastructureSet,
		WireDomainSet,
	)
	return &Domain{}, nil
}
