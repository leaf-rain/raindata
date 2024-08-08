package service

import (
	"context"
	"go.uber.org/zap"

	v1 "github.com/leaf-rain/raindata/app_bi/api/helloworld/v1"
	"github.com/leaf-rain/raindata/app_bi/internal/biz"
)

// GreeterService is a greeter service.
type GreeterService struct {
	v1.UnimplementedGreeterServer

	uc  *biz.GreeterUsecase
	log *zap.Logger
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUsecase, logger *zap.Logger) *GreeterService {
	return &GreeterService{uc: uc, log: logger}
}

// SayHello implements helloworld.GreeterServer.
func (s *GreeterService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	g, err := s.uc.CreateGreeter(ctx, &biz.Greeter{Hello: in.Name})
	if err != nil {
		return nil, err
	}
	return &v1.HelloReply{Message: "Hello " + g.Hello}, nil
}
