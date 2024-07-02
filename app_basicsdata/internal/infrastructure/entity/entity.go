package entity

import (
	"context"
	pb_metadata "github.com/leaf-rain/raindata/app_basicsdata/api/grpc"
	"github.com/leaf-rain/raindata/app_basicsdata/internal/infrastructure/config"
	"github.com/leaf-rain/raindata/common/clickhouse_sqlx"
	commonConfig "github.com/leaf-rain/raindata/common/config"
	"github.com/leaf-rain/raindata/common/etcd"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type Repository struct {
	etcdClient     *clientv3.Client
	clickhouseSqlx *clickhouse_sqlx.Conn
	dynamicConfig  commonConfig.ConfigInterface
	logger         *zap.Logger
	cfg            *config.Config
}

func NewRepository(ctx context.Context, etcdClient *clientv3.Client, clickhouseSqlx *clickhouse_sqlx.Conn, logger *zap.Logger, cfg *config.Config) (*Repository, error) {
	var repo *Repository
	if etcdClient != nil {
		// 初始化etcd配置中心
		metadataCache, _, _ := etcd.NewFileConfigByCatalogue(ctx, logger, etcdClient, cfg.MetadataPath, marshal)
		repo = &Repository{
			etcdClient:     etcdClient,
			clickhouseSqlx: clickhouseSqlx,
			dynamicConfig:  metadataCache,
			logger:         logger,
			cfg:            cfg,
		}
	} else {
		// 初始化本地配置中心
		metadataCache := commonConfig.NewFileConfigByCatalogue(logger, cfg.MetadataPath, "yaml", marshal)
		repo = &Repository{
			clickhouseSqlx: clickhouseSqlx,
			dynamicConfig:  metadataCache,
			logger:         logger,
			cfg:            cfg,
		}
	}
	return repo, nil
}

func (repo *Repository) Close() {
	repo.dynamicConfig.Close()
}

func marshal(name string) interface{} {
	return &pb_metadata.MetadataResponse{}
}
