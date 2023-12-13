package user

import (
	"context"
	pb "learning/grpc-project-service/api/gen/go/core/v1"
	"learning/grpc-project-service/pkg/provider"
)

type mockUserResource struct {
	provider.AbstractProvider
}

func (b *mockUserResource) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	return new(pb.GetUserResponse), nil
}
