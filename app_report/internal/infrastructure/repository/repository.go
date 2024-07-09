package repository

import (
	"github.com/leaf-rain/raindata/common/clickhouse_sqlx"
	"go.uber.org/zap"
)

type Repository struct {
	logger    *zap.Logger
	ckCluster *clickhouse_sqlx.ClickhouseCluster
}

func NewRepository(
	logger *zap.Logger,
	ckCluster *clickhouse_sqlx.ClickhouseCluster,
) *Repository {
	return &Repository{
		logger:    logger,
		ckCluster: ckCluster,
	}
}
