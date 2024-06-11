package infrastructure

import "go.uber.org/zap"

//go:generate wire

type Infrastructure struct {
	logger *zap.Logger
}

func NewInfrastructure(logger *zap.Logger) *Infrastructure {
	return &Infrastructure{
		logger: logger,
	}
}
