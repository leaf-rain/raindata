package interface_app

import (
	"context"

	pb_metadata "github.com/leaf-rain/raindata/app_basicsdata/api/grpc"
)

type InterfaceApplication interface {
	GetMetadata(ctx context.Context, in *pb_metadata.GetMetadataRequest) (*pb_metadata.MetadataResponse, error)
	PutMetadata(ctx context.Context, in *pb_metadata.MetadataRequest) (*pb_metadata.MetadataResponse, error)
}
