package data

import (
	"context"
	"github.com/leaf-rain/raindata/app_bi/internal/conf"
	"github.com/leaf-rain/raindata/app_bi/third_party/rredis"
	"github.com/leaf-rain/raindata/common/rgorm"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo)

// Data .
type Data struct {
	Ctx       context.Context
	RdClient  *rredis.Client
	SqlClient *gorm.DB
}

// NewData .
func NewData(c *conf.Data, logger *zap.Logger) (*Data, func(), error) {
	data := &Data{}
	var gormLogger = rgorm.NewGormZapLogger(logger)
	data.SqlClient = rgorm.NewRGrom(rgorm.DtGromConfig{
		DriverName:   c.Database.DriverName,
		DbSource:     c.Database.DbSource,
		MaxOpenConns: int(c.Database.MaxOpenConns),
		MaxIdleConns: int(c.Database.MaxIdleConns),
		IdleTimeOut:  int(c.Database.IdleTimeOut),
		Debug:        c.Database.Debug,
		Logger:       gormLogger,
	})
	var err error
	data.RdClient, err = rredis.NewRedis(rredis.Config{
		PoolSize:     0,
		Addr:         nil,
		Pwd:          "",
		DialTimeout:  0,
		ReadTimeout:  0,
		WriteTimeout: 0,
		DB:           0,
	}, context.Background())
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		err := data.RdClient.Close()
		if err != nil {
			logger.Error("[cleanup] close the data resources", zap.Error(err))
		}
		logger.Info("close the data resources")
	}
	return data, cleanup, nil
}
