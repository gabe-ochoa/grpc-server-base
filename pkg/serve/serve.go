package serve

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	// grpc
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rs/cors"

	// controllers
	"github.com/gabe-ochoa/grpc-server-base/pkg/config"
	health_controller "github.com/gabe-ochoa/grpc-server-base/pkg/controllers/health"

	// protos
	healthv1 "github.com/gabe-ochoa/grpc-server-base/protos/gen/v1/health"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func GRPC(cfg config.Config) error {
	log.WithFields(log.Fields{
		"grpcPort": cfg.GRPCPort,
	}).Info("Starting gRPC server")

	address := fmt.Sprintf("0.0.0.0:%s", cfg.GRPCPort)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	opts := grpcOpts(cfg)
	s := grpc.NewServer(opts)

	// TODO [gochoa] 2020-05-30: support the built in grpc health server
	// import "google.golang.org/grpc/health"
	// healthServer := health.NewServer()

	// Health server
	healthServer := health_controller.NewServer()
	healthv1.RegisterHealthAPIServer(s, healthServer)

	return s.Serve(listener)
}

func HTTPGateway(ctx context.Context, httpPort, grpcPort string, corsDebug bool) error {
	log.WithFields(log.Fields{
		"httpPort": httpPort,
		"grpcPort": grpcPort,
	}).Info("Starting http proxy server")

	// Register gRPC server endpoint and register HTTP Proxy
	// The grpc server has to be running and accesible
	mux := runtime.NewServeMux()

	// gRPC options
	opts := []grpc.DialOption{grpc.WithInsecure()}
	gRPCAddress := fmt.Sprintf("0.0.0.0:%s", grpcPort)

	// health
	err := healthv1.RegisterHealthAPIHandlerFromEndpoint(ctx, mux, gRPCAddress, opts)
	if err != nil {
		return err
	}

	handler := cors.New(CORSOpts(corsDebug)).Handler(mux)

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	httpAddress := fmt.Sprintf(":%s", httpPort)
	return http.ListenAndServe(httpAddress, handler)
}

func grpcOpts(cfg config.Config) grpc.ServerOption {
	// Setup Logging middleware
	// Logrus entry is used, allowing pre-definition of certain fields by the user.
	logrusEntry := log.NewEntry(log.StandardLogger())
	// Shared options for the logger, with a custom gRPC code to log level function.
	opts := []grpc_logrus.Option{
		grpc_logrus.WithLevels(grpc_logrus.DefaultCodeToLevel),
		grpc_logrus.WithDurationField(func(duration time.Duration) (key string, value interface{}) {
			return "grpc.time_us", duration.Microseconds()
		}),
	}

	// Make sure that log statements internal to gRPC library are logged using the logrus Logger as well.
	grpc_logrus.ReplaceGrpcLogger(logrusEntry)

	grpcOpts := grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_logrus.UnaryServerInterceptor(logrusEntry, opts...),
		grpc_prometheus.UnaryServerInterceptor,
		grpc_recovery.UnaryServerInterceptor(),
		// grpc_auth.UnaryServerInterceptor(middleware.APIAuth),
	))

	// Register server middleware without auth
	return grpcOpts
}

func CORSOpts(debug bool) cors.Options {
	return cors.Options{
		AllowedOrigins: []string{"http://localhost:8080"},
		AllowedMethods: []string{"GET", "POST"},
		Debug:          debug,
	}
}
