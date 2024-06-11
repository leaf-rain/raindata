package application

import (
	"context"
	"go.uber.org/zap"

	"github.com/leaf-rain/raindata/app_report/internal/application/interface_domain"
	"github.com/leaf-rain/raindata/app_report/internal/infrastructure/config"
)

type AppStream struct {
	logger       *zap.Logger
	config       *config.Config
	eventManager interface_domain.InterfaceEventManager
	ckWriter     interface_domain.InterfaceWriter
}

func NewAppStream(logger *zap.Logger, config *config.Config, eventManager interface_domain.InterfaceEventManager, ckWriter interface_domain.InterfaceWriter) *AppStream {
	return &AppStream{
		logger:       logger,
		config:       config,
		eventManager: eventManager,
		ckWriter:     ckWriter,
	}
}

func (app *AppStream) Stream(ctx context.Context, msg string) error {
	var err error
	err = app.eventManager.StorageEvent(ctx, msg)
	if err != nil {
		app.logger.Error("[Stream] eventManager.UpdateEventManager failed", zap.String("msg", msg), zap.Error(err))
		return err
	}
	// todo:直接写es
	err = app.ckWriter.WriterMsg(ctx, msg)
	if err != nil {
		app.logger.Error("[Stream] writer.WriterMsg failed", zap.String("msg", msg), zap.Error(err))
	}
	return err
}
