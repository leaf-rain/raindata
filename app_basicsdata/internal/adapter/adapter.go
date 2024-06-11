package adapter

//go:generate wire

type Adapter struct {
	GrpcServer *GrpcServer
}

func NewAdapter(grpcServer *GrpcServer) *Adapter {
	return &Adapter{
		GrpcServer: grpcServer,
	}
}

func (ada *Adapter) Run() {
	ada.GrpcServer.Run()
}

func (ada *Adapter) Close() {
	ada.GrpcServer.Close()
}
