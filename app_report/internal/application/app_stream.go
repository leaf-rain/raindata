package application

import (
	"context"
	"go.uber.org/zap"

	"github.com/leaf-rain/raindata/app_report/internal/application/interface_domain"
	"github.com/leaf-rain/raindata/app_report/internal/infrastructure/config"
)

type AppStream struct {
	logger *zap.Logger
	config *config.Config
	writer interface_domain.InterfaceWriter
}

func NewAppStream(logger *zap.Logger, config *config.Config, writer interface_domain.InterfaceWriter) *AppStream {
	return &AppStream{
		logger: logger,
		config: config,
		writer: writer,
	}
}

func (app *AppStream) Stream(ctx context.Context, msg string) error {
	var err error
	err = app.writer.WriterMsg(ctx, msg)
	if err != nil {
		app.logger.Error("[Stream] writer.WriterMsg failed", zap.String("msg", msg), zap.Error(err))
	}
	return err
}
