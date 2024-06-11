package application

import (
	"go.uber.org/zap"
)

//go:generate wire

type Applications struct {
	logger      *zap.Logger
	appMetadata *AppMetadata
}

func NewApplications(logger *zap.Logger, am *AppMetadata) *Applications {
	return &Applications{
		logger:      logger,
		appMetadata: am,
	}
}
