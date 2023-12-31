package grpc

import (
	"context"
	"fmt"
	"net"
	"os"
	"runtime/debug"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"learning/grpc-project-service/pkg/provider"
)

// CustomOpts when create a grpc server, you can custom yourself interceptor.
type CustomOpts struct {
	UnaryInterceptor  []grpc.UnaryServerInterceptor
	StreamInterceptor []grpc.StreamServerInterceptor
	ServerOption      []grpc.ServerOption
}

// Server grpc server provider.
// Provides a server that listens for grpc traffic and forwards them to the configured handlers.
// Also adds a bunch of useful interceptors for tracing, metrics and so on.
// Relies heavily on https://github.com/grpc-ecosystem packages.
type Server struct {
	provider.AbstractRunProvider

	Config   *Config
	Listener net.Listener
	Server   *grpc.Server
	Opts     []CustomOpts
}

func recoverHandler(ctx context.Context, p interface{}) (err error) {
	fmt.Fprintf(os.Stderr, "Service_Recover_Handler: %v, stack: \n%s", p, debug.Stack())
	//provider_logrus.WithError(errors.WithStack(status.Errorf(codes.Internal, "%v", p))).Errorln("service panic")
	return status.Errorf(codes.Internal, "%v", p)
}

// New creates a grpc server provider.
func New(config *Config, customOpts ...CustomOpts) *Server {
	if config == nil {
		config = NewConfigFromEnv()
	}

	return &Server{
		Config: config,
		Opts:   customOpts,
	}
}

// Init creates the grpc server (doesn't start it yet) and adds useful interceptors.
func (p *Server) Init() error {
	logger := logrus.NewEntry(logrus.StandardLogger())

	grpc_logrus.JsonPbMarshaller = NewJsonPbMarshaller()
	opts := []grpc_logrus.Option{
		grpc_logrus.WithDurationField(func(duration time.Duration) (key string, value interface{}) {
			return "grpc.time_ns", duration.Nanoseconds()
		}),
	}

	// Unary and streaming have the same interceptors.
	unaryInterceptors := []grpc.UnaryServerInterceptor{
		grpc_ctxtags.UnaryServerInterceptor(),
		// @review need combine otelgrpc.UnaryServerInterceptor into tracing.OtelInterceptor
		//otelgrpc.UnaryServerInterceptor(),
		// insert trace id into context
		//tracing.OtelInterceptor(),
		grpc_prometheus.UnaryServerInterceptor,
		grpc_auth.UnaryServerInterceptor(p.authFunc),
		grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandlerContext(recoverHandler)),
	}
	streamInterceptors := []grpc.StreamServerInterceptor{
		grpc_ctxtags.StreamServerInterceptor(),
		otelgrpc.StreamServerInterceptor(),
		grpc_prometheus.StreamServerInterceptor,
		grpc_auth.StreamServerInterceptor(p.authFunc),
		grpc_recovery.StreamServerInterceptor(grpc_recovery.WithRecoveryHandlerContext(recoverHandler)),
	}

	if p.Config.LogInterceptor {
		unaryInterceptors = append(unaryInterceptors, grpc_logrus.UnaryServerInterceptor(logger, opts...))
		streamInterceptors = append(streamInterceptors, grpc_logrus.StreamServerInterceptor(logger, opts...))
	}

	// Payload is only logged by the Server if it was configured to do so.
	if p.Config.LogPayload {
		unaryInterceptors = append(unaryInterceptors, grpc_logrus.PayloadUnaryServerInterceptor(logger, p.logDeciderFunc))
		streamInterceptors = append(streamInterceptors, grpc_logrus.PayloadStreamServerInterceptor(logger, p.logDeciderFunc))
	}

	var serverOpts []grpc.ServerOption
	for _, opt := range p.Opts {
		unaryInterceptors = append(unaryInterceptors, opt.UnaryInterceptor...)
		streamInterceptors = append(streamInterceptors, opt.StreamInterceptor...)
		serverOpts = append(serverOpts, opt.ServerOption...)
	}

	serverOpts = append(
		serverOpts,
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(unaryInterceptors...)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(streamInterceptors...)),
	)

	p.Server = grpc.NewServer(serverOpts...)

	return nil
}

// Run creates a grpc listener on the configured port which is used to start the grpc server.
// Uses the grpc server reflection functionality find the available handlers.
func (p *Server) Run() error {
	addr := fmt.Sprintf(":%d", p.Config.Port)
	//logEntry := logrus.WithField("addr", addr)

	reflection.Register(p.Server)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		//logEntry.WithError(err).Error("GRPC Server Listener could not be created")
		return err
	}
	p.Listener = listener
	p.SetRunning(true)
	p.registerHealthEndpoint()

	//logEntry.Info("GRPC Server Provider launched")
	if err := p.Server.Serve(listener); err != nil {
		//logEntry.WithError(err).Error("GRPC Server Provider launch failed")
		return err
	}

	return nil
}

// Close shuts down the grpc server.
func (p *Server) Close() error {
	p.Server.GracefulStop()

	return p.AbstractRunProvider.Close()
}

func (p *Server) authFunc(ctx context.Context) (context.Context, error) {
	// TODO: Add support for authentication.
	return ctx, nil
}

func (p *Server) logDeciderFunc(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
	// TODO: Should we really log everything?
	return true
}

func (p *Server) registerHealthEndpoint() {
	if !p.Config.EnableHealth {
		//logrus.Debug("GRPC Server health endpoint disabled")
		return
	}
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(p.Server, healthServer)
	//logrus.Debug("GRPC Server health endpoint registered")
}
