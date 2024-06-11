package entity

import "go.uber.org/zap"

type FieldInfo struct {
	FieldName   string
	DbFieldName string
	FieldType   string
}

type EventManager struct {
	logger *zap.Logger
}

func NewEventManager(logger *zap.Logger) *EventManager {
	return &EventManager{
		logger: logger,
	}
}
