package router

import (
	"learning/grpc-project-service/pkg/provider"
	"learning/grpc-project-service/pkg/provider/grpc"
	"learning/grpc-project-service/pkg/provider/grpc/gateway"

	pb "learning/grpc-project-service/api/gen/go/platform/v1"
	"learning/grpc-project-service/internal/controller"
)

type Router struct {
	provider.AbstractRunProvider
	controller      controller.Controller
	grpcProvider    *grpc.Server
	gatewayProvider *gateway.Gateway
}

func NewRouter(grpcProvider *grpc.Server, gatewayProvider *gateway.Gateway, controller controller.Controller) *Router {
	return &Router{
		controller:      controller,
		grpcProvider:    grpcProvider,
		gatewayProvider: gatewayProvider,
	}
}

func (r *Router) Init() error {
	pb.RegisterProjectAPIServer(r.grpcProvider.Server, r.controller)
	return nil
}

func (r *Router) Run() error {
	if r.gatewayProvider.Config.Enabled {
		if err := provider.WaitForRunningProvider(r.gatewayProvider, 2); err != nil {
			return err
		}

		if err := r.gatewayProvider.RegisterServices(
			pb.RegisterProjectAPIHandler,
		); err != nil {
			//logging.WithError(err).Errorf("Could not register gateway service handlers")
			return err
		}
	}
	//logging.Infof("Scaffold router launched successful")
	r.SetRunning(true)
	return nil
}
