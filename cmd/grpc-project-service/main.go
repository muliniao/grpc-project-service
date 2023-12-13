package main

import (
	"learning/grpc-project-service/internal/controller"
	"learning/grpc-project-service/internal/router"
	"learning/grpc-project-service/internal/service"
	"learning/grpc-project-service/pkg/provider/app"
	"learning/grpc-project-service/pkg/provider/grpc"
	"learning/grpc-project-service/pkg/provider/grpc/gateway"
	"learning/grpc-project-service/pkg/stack"
)

func main() {
	st := stack.New()
	defer st.MustClose()

	// Root app
	appConfig := app.NewConfigFromEnv()
	appProvider := app.New(appConfig)
	st.MustInit(appProvider)

	// grpc
	grpcConfig := grpc.NewConfigFromEnv()
	grpcProvider := grpc.New(grpcConfig)
	st.MustInit(grpcProvider)

	// grpc-gateway
	gatewayConfig := gateway.NewConfigFromEnv()
	gatewayProvider := gateway.New(gatewayConfig, grpcProvider, appProvider)
	st.MustInit(gatewayProvider)

	svc := service.New()
	st.MustInit(svc)

	rt := router.NewRouter(grpcProvider, gatewayProvider, controller.New(svc))
	st.MustInit(rt)

	st.MustRun()
}
