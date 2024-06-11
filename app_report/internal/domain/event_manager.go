package domain

import (
	"context"
	"github.com/leaf-rain/fastjson"
	"go.uber.org/zap"
)

type EventManager struct {
	logger *zap.Logger
}

func NewEventManager(logger *zap.Logger) *EventManager {
	return &EventManager{
		logger: logger,
	}
}

func (em *EventManager) StorageEvent(ctx context.Context, msg string) error {
	v, err := fastjson.Parse(msg)
	if err != nil || v == nil {
		em.logger.Error("[StorageEvent] fastjson.Parse failed", zap.String("msg", msg), zap.Error(err))
		return err
	}

	panic("implement me")
}
