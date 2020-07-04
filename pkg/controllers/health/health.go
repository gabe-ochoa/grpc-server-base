package controllers

import (
	healthv1 "github.com/gabe-ochoa/grpc-server-base/protos/gen/v1/health"

	"context"
)

// Interface provided for mocks and testing
type Server interface {
	healthv1.HealthAPIServer
}

type server struct{}

func NewServer() Server {
	return &server{}
}

// API Server implementation / routes to other controllers

func (s *server) Health(ctx context.Context, request *healthv1.HealthRequest) (*healthv1.HealthResponse, error) {
	return &healthv1.HealthResponse{
		Status: 1,
	}, nil
}
