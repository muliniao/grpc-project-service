package user

import (
	"context"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	grpcClient "google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	pb "learning/grpc-project-service/api/gen/go/core/v1"
	"learning/grpc-project-service/pkg/provider"
)

type (
	UserResource interface {
		provider.Provider
		GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error)
	}

	userResource struct {
		provider.AbstractProvider
		Config *Config
		client pb.UserAPIClient
	}
)

func (b *userResource) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Resource::GetUser")
	defer span.Finish()

	return b.client.GetUser(ctx, req)
}

func New(config *Config) UserResource {
	if config.MockEnabled {
		return new(mockUserResource)
	}
	return &userResource{
		Config: config,
	}
}

func (b *userResource) Init() error {
	//logEntry := logging.WithFields(logging.Fields{
	//	"user_mock_enabled": b.Config.MockEnabled,
	//	"user_addr":         b.Config.UserAddr,
	//})

	if b.Config.MockEnabled {
		//logEntry.Info("Initial User Resource With Mocked Stub")
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	unaryInterceptors := []grpcClient.UnaryClientInterceptor{
		grpc_opentracing.UnaryClientInterceptor(),
	}

	// don't use sdk
	//conn, err := grpcSdk.NewClientConn(
	//	grpcSdk.WithCtx(ctx),
	//	grpcSdk.WithServerAddr(b.Config.UserAddr),
	//	grpcSdk.WithTenantID(utils.GetTenantID()),
	//)
	conn, err := grpcClient.DialContext(ctx, b.Config.UserAddr,
		grpcClient.WithInsecure(),
		grpcClient.WithKeepaliveParams(keepalive.ClientParameters{
			PermitWithoutStream: true,
		}),
		grpcClient.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(unaryInterceptors...)),
	)
	if err != nil {
		//logEntry.WithError(err).Errorf("Billing Provider launch failed")
		return err
	}
	b.client = pb.NewUserAPIClient(conn)
	return nil
}
