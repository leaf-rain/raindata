package domain

import (
	"context"
	pb_metadata "github.com/leaf-rain/raindata/app_basicsdata/api/grpc"
	"github.com/leaf-rain/raindata/app_basicsdata/internal/infrastructure/config"
	"github.com/leaf-rain/raindata/app_basicsdata/internal/infrastructure/entity"
	"github.com/leaf-rain/raindata/common/clickhouse_sqlx"
	"github.com/leaf-rain/raindata/common/consts"
	"github.com/leaf-rain/raindata/common/ecode"
	"github.com/leaf-rain/raindata/common/etcd"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"sync"
	"time"
)

type Metadata struct {
	logger     *zap.Logger
	repo       *entity.Repository
	locks      sync.Map
	cfg        *config.Config
	etcdClient *clientv3.Client
	ck         *clickhouse_sqlx.Clickhouse
	ttl        int64
}

func (domain *Metadata) keyMetadata(app, eventName string) string {
	return consts.ETCDKeyLockPre + "/" + app + "/" + eventName
}

func NewMetadata(logger *zap.Logger, repo *entity.Repository, cfg *config.Config, etcdClient *clientv3.Client, ck *clickhouse_sqlx.Clickhouse) (*Metadata, error) {
	return &Metadata{
		logger:     logger.Named("domain.Metadata"),
		repo:       repo,
		cfg:        cfg,
		etcdClient: etcdClient,
		ck:         ck,
		locks:      sync.Map{},
		ttl:        5,
	}, nil
}

func (domain *Metadata) GetMetadata(ctx context.Context, in *pb_metadata.GetMetadataRequest) (*pb_metadata.MetadataResponse, error) {
	entityBody := domain.repo.NewMetadata(ctx, domain.logger, &pb_metadata.MetadataRequest{
		App:       in.App,
		EventName: in.EventName,
	}, domain.ck)
	return entityBody.GetMetadata()
}

func (domain *Metadata) PutMetadata(ctx context.Context, in *pb_metadata.MetadataRequest) (*pb_metadata.MetadataResponse, error) {
	lockKey := domain.keyMetadata(in.App, in.EventName)
	if _, ok := domain.locks.Load(lockKey); ok {
		return nil, ecode.ERR_METADATA_UPDATEING
	}
	var now = time.Now().Unix()
	domain.locks.Store(lockKey, now)
	defer domain.locks.Delete(lockKey)
	var err error
	// 加载一次内存，如果内存中数据相同则，返回
	entityBody := domain.repo.NewMetadata(ctx, domain.logger, in, domain.ck)
	// 存在修改在获取分布式锁
	if domain.etcdClient != nil {
		var lock *etcd.Lock
		lock, err = etcd.NewLock(ctx, domain.etcdClient, lockKey, domain.ttl)
		if err != nil {
			domain.logger.Error("[StorageMetadata] init lock error", zap.Error(err))
			return nil, err
		}
		err = lock.Lock()
		if err != nil {
			domain.logger.Error("[StorageMetadata] lock error", zap.Error(err))
			return nil, err
		}
		defer func() {
			if err = lock.Unlock(); err != nil {
				domain.logger.Error("[StorageMetadata] unlock error", zap.Error(err))
			}
		}()
	}
	// 校验内存
	var entityCache *pb_metadata.MetadataResponse
	entityCache, err = entityBody.GetMetadata()
	if err != nil {
		domain.logger.Error("[StorageMetadata] entityBody.GetMetadata failed", zap.String("name", in.EventName), zap.Error(err))
		return nil, err
	}
	var result *pb_metadata.MetadataResponse
	var ty int
	if entityCache == nil {
		entityBody.SetFields(&pb_metadata.MetadataRequest{
			App:       in.App,
			Fields:    nil,
			EventName: in.EventName,
		})
		result = entityBody.GroupMetadata(in.Fields)
		ty = 1
	} else {
		key := entity.GetMetadataKey(in.Fields)
		cacheKey := entity.GetMetadataKey(entityCache.Metadata)
		if key == cacheKey {
			return entityCache, nil
		}
		entityBody.SetFields(&pb_metadata.MetadataRequest{
			App:       in.App,
			Fields:    entityCache.Metadata,
			EventName: in.EventName,
		})
		result = entityBody.GroupMetadata(in.Fields)
		ty = 2
	}
	result, err = entityBody.PutMetadata(ty)
	if err != nil {
		domain.logger.Error("[StorageMetadata] entityBody.PutMetadata failed", zap.String("name", in.EventName), zap.Error(err))
	}
	return result, err
}
