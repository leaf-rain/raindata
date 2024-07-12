package application

import (
	"github.com/leaf-rain/raindata/app_report/internal/application/interface_domain"
	"github.com/leaf-rain/raindata/app_report/internal/infrastructure/config"
	"go.uber.org/zap"
)

//go:generate wire

type Applications struct {
	config *config.Config
	logger *zap.Logger
	writer interface_domain.InterfaceWriter
}

func NewApplications(config *config.Config, logger *zap.Logger, writer interface_domain.InterfaceWriter) *Applications {
	return &Applications{
		config: config,
		logger: logger,
		writer: writer,
	}
}
