package config

import (
	env "github.com/caarlos0/env/v6"

	log "github.com/sirupsen/logrus"
)

type Config struct {
	Port          string `env:"PORT" envDefault:"8081"`
	GRPCPort      string `env:"GRPC_PORT" envDefault:"9090"`
	envServerMode string `env:"SERVER_MODE" envDefault:"development"`
	ServerMode    ServerMode
	ServiceName   string `env:"SERVICE_NAME" envDefault:"grpc-server-api"`
	LogLevel      string `env:"LOG_LEVEL" envDefault:"debug"`
	LogFormat     string `env:"LOG_FORMAT" envDefault:"text"`
}

// MustLoadConfig is used to load the config from ENV and defaults
func MustLoadConfig() Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("%+v\n", err)
	}

	cfg.ServerMode = parseServerMode(cfg.envServerMode)

	// This log line must but INFO because we haven't configured the logger yet
	log.WithFields(log.Fields{
		"config": cfg,
	}).Info("Configuration loaded")
	return cfg
}

type ServerMode int

const (
	Development ServerMode = iota
	Staging
	Production
)

func parseServerMode(envServerMode string) ServerMode {
	switch envServerMode {
	case "staging":
		return Staging
	case "production":
		return Production
	}
	// Default to development
	return Development
}

func (s ServerMode) String() string {
	switch s {
	case Staging:
		return "staging"
	case Production:
		return "production"
	}
	return "development"
}
