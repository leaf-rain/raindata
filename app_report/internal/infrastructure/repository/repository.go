package repository

import (
	"github.com/leaf-rain/raindata/common/rclickhouse"
	"go.uber.org/zap"
)

type Repository struct {
	logger    *zap.Logger
	ckCluster *rclickhouse.ClickhouseCluster
}

func NewRepository(
	logger *zap.Logger,
	ckCluster *rclickhouse.ClickhouseCluster,
) *Repository {
	return &Repository{
		logger:    logger,
		ckCluster: ckCluster,
	}
}
