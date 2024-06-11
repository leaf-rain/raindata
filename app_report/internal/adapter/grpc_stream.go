package adapter

import (
	"context"
	raindata_pb "github.com/leaf-rain/raindata/app_report/api/grpc"
	"github.com/leaf-rain/raindata/app_report/internal/adapter/interface_app"
	config "github.com/leaf-rain/raindata/app_report/internal/infrastructure/config"
	"github.com/leaf-rain/raindata/common/recuperate"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"io"
)

type GrpcServer struct {
	config       *config.Config
	logger       *zap.Logger
	methodStream interface_app.InterfaceStream
	raindata_pb.UnimplementedLogServerServer
}

func NewGrpcServer(config *config.Config, logger *zap.Logger, methodStream interface_app.InterfaceStream) *GrpcServer {
	return &GrpcServer{
		config:       config,
		logger:       logger,
		methodStream: methodStream,
	}
}

func (gs *GrpcServer) Run() {
	var err error
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	raindata_pb.RegisterLogServerServer(grpcServer, gs)
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

func (gs *GrpcServer) StreamReport(stream raindata_pb.LogServer_StreamReportServer) error {
	// 1.for循环接收客户端发送的消息
	for {
		// 2. 通过 Recv() 不断获取客户端 send()推送的消息
		req, err := stream.Recv() // Recv内部也是调用RecvMsg
		// 3. err == io.EOF表示已经获取全部数据
		if err == io.EOF {
			// 4.SendAndClose 返回并关闭连接
			// 在客户端发送完毕后服务端即可返回响应
			gs.logger.Info("[StreamReport] conn colse.", zap.Any("stream", stream))
			return stream.SendAndClose(&raindata_pb.StreamReportResponse{})
		}
		if err != nil {
			return err
		}
		gs.logger.Info("[StreamReport] received.", zap.Any("msg", req.GetMessage()))
		err = gs.methodStream.Stream(context.TODO(), req.GetMessage())
		if err != nil {
			gs.logger.Error("[StreamReport] methodStream.Stream failed", zap.String("message", req.GetMessage()), zap.Error(err))
		}
	}
}
