package acl

import (
	"context"
	basicsdata_metadata "github.com/leaf-rain/raindata/app_basicsdata/api/grpc"
)

type BasicsData interface {
	GetMetadata(ctx context.Context, in *basicsdata_metadata.GetMetadataRequest) (*basicsdata_metadata.MetadataResponse, error)
	PutMetadata(ctx context.Context, in *basicsdata_metadata.MetadataRequest) (*basicsdata_metadata.MetadataResponse, error)
}

// 大泥球模式
var _ BasicsData = (*MudBallBasicsData)(nil)

type MudBallBasicsData struct {
}

func (m MudBallBasicsData) GetMetadata(ctx context.Context, in *basicsdata_metadata.GetMetadataRequest) (*basicsdata_metadata.MetadataResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (m MudBallBasicsData) PutMetadata(ctx context.Context, in *basicsdata_metadata.MetadataRequest) (*basicsdata_metadata.MetadataResponse, error) {
	//TODO implement me
	panic("implement me")
}
