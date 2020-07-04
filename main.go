package main

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/gabe-ochoa/grpc-server-base/pkg/config"
	"github.com/gabe-ochoa/grpc-server-base/pkg/middleware"
	"github.com/gabe-ochoa/grpc-server-base/pkg/serve"
)

// Version is set via linker vaiable at build time
// Version conforms to CalVer in the form YYYY.MM.DD.HH.MM.SS
var (
	Version string
)

func main() {
	log.Infof("Starting grpc-server version %s 😇", Version)

	// Get server runtime configuration
	cfg := config.MustLoadConfig()

	// Logger
	middleware.SetupLogger(cfg.LogFormat, cfg.LogLevel)

	// Start the http proxy server in a goroutine
	go func() {
		// Set up a context
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		// Enable cors debug logging but only if we're at a debug log level
		corsDebug := false
		if log.GetLevel() == log.DebugLevel {
			corsDebug = true
		}

		if err := serve.HTTPGateway(ctx, cfg.Port, cfg.GRPCPort, corsDebug); err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Fatal("gRPC Server encountered an error. Shutting down 😔")
		}
	}()

	if err := serve.GRPC(cfg); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("gRPC Server encountered an error. Shutting down 😔")
	}
}
