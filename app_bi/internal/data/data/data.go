package data

import (
	"context"
	"github.com/leaf-rain/raindata/app_bi/internal/conf"
	"github.com/leaf-rain/raindata/app_bi/third_party/rredis"
	"github.com/leaf-rain/raindata/common/rgorm"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"

	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData)

// Data .
type Data struct {
	Ctx               context.Context
	RdClient          *rredis.Client
	SqlClient         *gorm.DB
	SingleflightGroup *singleflight.Group
	Config            *conf.Bootstrap
	logger            *zap.Logger
}

// NewData .
func NewData(c *conf.Bootstrap, logger *zap.Logger) (*Data, func(), error) {
	data := &Data{
		Config: c,
		logger: logger,
	}
	var gormLogger = rgorm.NewGormZapLogger(logger)
	data.SqlClient = rgorm.NewRGrom(rgorm.DtGromConfig{
		DriverName:   c.Data.Database.DriverName,
		DbSource:     c.Data.Database.DbSource,
		MaxOpenConns: int(c.Data.Database.MaxOpenConns),
		MaxIdleConns: int(c.Data.Database.MaxIdleConns),
		IdleTimeOut:  int(c.Data.Database.IdleTimeOut),
		Debug:        c.Data.Database.Debug,
		Logger:       gormLogger,
	})
	var err error
	data.RdClient, err = rredis.NewRedis(rredis.Config{
		PoolSize:     int(c.Data.Redis.PoolSize),
		Addr:         c.Data.Redis.Addr,
		Pwd:          c.Data.Redis.Pwd,
		DialTimeout:  c.Data.Redis.DialTimeout,
		ReadTimeout:  c.Data.Redis.ReadTimeout,
		WriteTimeout: c.Data.Redis.WriteTimeout,
		DB:           int(c.Data.Redis.Db),
	}, context.Background())
	if err != nil {
		return nil, nil, err
	}
	data.SingleflightGroup = &singleflight.Group{}
	cleanup := func() {
		err = data.RdClient.Close()
		if err != nil {
			logger.Error("[cleanup] close the data resources", zap.Error(err))
		}
		logger.Info("close the data resources")
	}
	data.initDb()
	return data, cleanup, nil
}

type initDb interface {
	MigrateTable(ctx context.Context) error
	TableCreated(ctx context.Context) bool
	InitializeData(ctx context.Context) error
}

func (data *Data) initDb() {
	data.createTable(NewEntityExaFileUploadAndDownload(data))
	data.createTable(NewEntitySysApi(data))
	data.createTable(NewEntitySysIgnoreApi(data))
	data.createTable(NewEntitySysAuthority(data))
	data.createTable(NewEntitySysCasbin(data))
	data.createTable(NewJWT(data))
	data.createTable(NewEntitySysOperationRecord(data))
	data.createTable(NewEntitySysUser(data))
	data.createTable(NewEntitySysUserAuthority(data))
}

func (data *Data) createTable(i initDb) {
	if !i.TableCreated(data.Ctx) {
		err := i.MigrateTable(data.Ctx)
		if err != nil {
			panic("init db failed. err:" + err.Error())
		}
		err = i.InitializeData(data.Ctx)
		if err != nil {
			panic("init data failed. err:" + err.Error())
		}
	}
}
