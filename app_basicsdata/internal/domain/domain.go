package domain

import (
	"go.uber.org/zap"
)

//go:generate wire

type Domain struct {
	logger *zap.Logger
	// 元数据管理
	metadata *Metadata
}

func NewDomain(
	logger *zap.Logger,
	metadata *Metadata,
) *Domain {
	return &Domain{
		logger:   logger,
		metadata: metadata,
	}
}
