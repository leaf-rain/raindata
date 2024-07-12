package application

import (
	"context"
	"go.uber.org/zap"
	"strconv"
)

type AppStream struct {
	app *Applications
}

func NewAppStream(app *Applications) *AppStream {
	return &AppStream{
		app: app,
	}
}

func (app *AppStream) Stream(ctx context.Context, msg string) error {
	var appidStr = msg[:32]
	appid, err := strconv.ParseInt(appidStr, 10, 64)
	if err != nil {
		app.app.logger.Error("[Stream] parse appid failed", zap.String("msg", msg), zap.Error(err))
	}
	err = app.app.writer.WriterMsg(ctx, appid, msg[32:])
	if err != nil {
		app.app.logger.Error("[Stream] writer.WriterMsg failed", zap.String("msg", msg), zap.Error(err))
	}
	return err
}
