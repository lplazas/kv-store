package server

import (
	"context"
	"errors"
	"github.com/gc-plazas/kv-store/internal"
	"github.com/gc-plazas/kv-store/internal/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ ServerServiceServer = (*ClusterServer)(nil)

type ClusterServer struct {
	UnimplementedServerServiceServer
	cluster *internal.Cluster
}

func (c ClusterServer) GetValue(ctx context.Context, request *GetValueRequest) (*GetValueResponse, error) {
	if request.Key == "" {
		return nil, status.Error(codes.InvalidArgument, "key cannot be empty")
	}

	value, err := c.cluster.GetValue(ctx, request.Key)
	if err != nil {
		if errors.As(err, &errs.ValueNotFound{}) {
			return nil, status.Error(codes.NotFound, "no value found for key")
		}
	}

	return &GetValueResponse{
		Value: value,
	}, nil
}

func (c ClusterServer) mustEmbedUnimplementedServerServiceServer() {
	panic("implement me")
}
