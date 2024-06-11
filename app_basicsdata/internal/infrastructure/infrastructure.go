package infrastructure

import (
	"github.com/leaf-rain/raindata/app_basicsdata/internal/infrastructure/entity"
	"go.uber.org/zap"
)

//go:generate wire

type Infrastructure struct {
	logger *zap.Logger
	repo   *entity.Repository
}

func NewInfrastructure(logger *zap.Logger, repo *entity.Repository) *Infrastructure {
	return &Infrastructure{
		logger: logger,
		repo:   repo,
	}
}
