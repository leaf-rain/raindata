package acl

import (
	"context"
	basicsdata_metadata "github.com/leaf-rain/raindata/app_basicsdata/api/grpc"
)

type BasicsData interface {
	GetMetadata(ctx context.Context, in *basicsdata_metadata.GetMetadataRequest) (*basicsdata_metadata.MetadataResponse, error)
	PutMetadata(ctx context.Context, in *basicsdata_metadata.MetadataRequest) (*basicsdata_metadata.MetadataResponse, error)
}

type mudBall struct {
}
