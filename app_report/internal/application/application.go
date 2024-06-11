package application

import (
	"go.uber.org/zap"
)

//go:generate wire

type Applications struct {
	logger    *zap.Logger
	appStream *AppStream
}

func NewApplications(logger *zap.Logger, appStream *AppStream) *Applications {
	return &Applications{
		logger:    logger,
		appStream: appStream,
	}
}
