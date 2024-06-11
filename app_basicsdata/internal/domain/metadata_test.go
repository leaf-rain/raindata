package domain

import (
	"context"
	pb_metadata "github.com/leaf-rain/raindata/app_basicsdata/api/grpc"
	"testing"
)

func TestMetadata_GetMetadata(t *testing.T) {
	result, err := domain.metadata.GetMetadata(context.Background(), &pb_metadata.GetMetadataRequest{
		App:       "test_app",
		EventName: "test_event",
	})
	if err != nil || result == nil {
		t.Fatal(err)
	}
	for _, item := range result.Metadata {
		t.Log(item)
	}
}

func TestMetadata_PutMetadata(t *testing.T) {
	result, err := domain.metadata.PutMetadata(context.Background(), &pb_metadata.MetadataRequest{
		App:       "test_app",
		EventName: "test_event",
		Fields: []*pb_metadata.Field{
			{
				Name: "tf1",
				Type: "String",
			},
		},
	})
	if err != nil || result == nil {
		t.Fatal(err)
	}
	for _, item := range result.Metadata {
		t.Log(item)
	}
}
