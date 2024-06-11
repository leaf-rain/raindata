//go:build wireinject
// +build wireinject

package domain

import (
	"github.com/google/wire"
	"github.com/leaf-rain/raindata/app_basicsdata/internal/infrastructure"
)

var WireDomainSet = wire.NewSet(
	NewMetadata,
	NewDomain,
)

func Initialize() (*Domain, error) {
	wire.Build(
		infrastructure.WireInfrastructureSet,
		WireDomainSet,
	)
	return &Domain{}, nil
}
