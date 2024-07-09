package domain

import (
	"github.com/leaf-rain/raindata/app_report/internal/application/interface_domain"
	"go.uber.org/zap"
)

//go:generate wire

var _ interface_domain.InterfaceWriter = (*Domain)(nil)

type Domain struct {
	logger *zap.Logger
	// 写ck
	ckWriter *Writer
}

func NewDomain(
	logger *zap.Logger,
	ckWriter *Writer,
) *Domain {
	return &Domain{
		logger:   logger,
		ckWriter: ckWriter,
	}
}
