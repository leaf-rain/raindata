package adapter

import (
	"context"
	pb_metadata "github.com/leaf-rain/raindata/app_basicsdata/api/grpc"
	"github.com/leaf-rain/raindata/app_basicsdata/internal/adapter/interface_app"
	"github.com/leaf-rain/raindata/app_basicsdata/internal/infrastructure/config"
	"github.com/leaf-rain/raindata/common/recuperate"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	config      *config.Config
	logger      *zap.Logger
	application interface_app.InterfaceApplication
	pb_metadata.UnimplementedMetadataServerServer
}

func NewGrpcServer(config *config.Config, logger *zap.Logger, application interface_app.InterfaceApplication) *GrpcServer {
	return &GrpcServer{
		config:      config,
		logger:      logger,
		application: application,
	}
}

func (gs *GrpcServer) Run() {
	var err error
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb_metadata.RegisterMetadataServerServer(grpcServer, gs)
	recuperate.GoSafe(func() {
		gs.logger.Info("grpc server start success.", zap.String("addr", gs.config.GrpcAddr))
		err = grpcServer.Serve(gs.config.GrpcListener)
		if err != nil {
			gs.logger.Fatal("grpcServer Serve failed", zap.Error(err))
		}
	})
}

func (gs *GrpcServer) Close() {

}

func (gs *GrpcServer) GetMetadata(ctx context.Context, in *pb_metadata.GetMetadataRequest) (*pb_metadata.MetadataResponse, error) {
	return gs.application.GetMetadata(ctx, in)
}

func (gs *GrpcServer) PutMetadata(ctx context.Context, in *pb_metadata.MetadataRequest) (*pb_metadata.MetadataResponse, error) {
	return gs.application.PutMetadata(ctx, in)
}
