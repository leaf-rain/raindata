package server

import (
	"context"
	"github.com/google/wire"
	"github.com/leaf-rain/raindata/app_bi/internal/data/data"
	"github.com/leaf-rain/raindata/app_bi/internal/service"
	"go.uber.org/zap"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewServer, NewHTTPServer)

type Server struct {
	Ctx    context.Context
	svc    *service.Service
	data   *data.Data
	logger *zap.Logger
}

func NewServer(ctx context.Context, svc *service.Service, data *data.Data, logger *zap.Logger) *Server {
	server := &Server{
		Ctx:    ctx,
		svc:    svc,
		data:   data,
		logger: logger,
	}
	return server
}
