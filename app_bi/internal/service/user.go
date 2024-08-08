package service

import (
	"go.uber.org/zap"
)

type UserService struct {
	log *zap.Logger
}

// NewUserService new a greeter service.
func NewUserService(logger *zap.Logger) *UserService {
	return &UserService{log: logger}
}
