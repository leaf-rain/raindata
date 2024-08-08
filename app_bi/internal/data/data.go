package data

import (
	"github.com/leaf-rain/raindata/app_bi/internal/conf"
	"go.uber.org/zap"

	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo)

// Data .
type Data struct {
	// TODO wrapped database client
}

// NewData .
func NewData(c *conf.Data, logger *zap.Logger) (*Data, func(), error) {
	cleanup := func() {
		logger.Info("close the data resources")
	}
	return &Data{}, cleanup, nil
}
