package domain

import (
	"go.uber.org/zap"
)

//go:generate wire

type Domain struct {
	logger *zap.Logger
	// 写ck
	ckWriter *CkWriter
	// 事件管理
	eventManager *EventManager
}

func NewDomain(
	logger *zap.Logger,
	ckWriter *CkWriter,
	eventManager *EventManager,
) *Domain {
	return &Domain{
		logger:       logger,
		ckWriter:     ckWriter,
		eventManager: eventManager,
	}
}
