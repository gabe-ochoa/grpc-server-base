package controllers

import (
	"testing"

	healthv1 "github.com/gabe-ochoa/grpc-server-base/protos/gen/v1/health"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestHealth(t *testing.T) {
	testServer := NewServer()
	testCtx := context.Background()
	request := healthv1.HealthRequest{}

	response, err := testServer.Health(testCtx, &request)

	assert.NoError(t, err)

	assert.Equal(t, "SERVING", response.Status.String())
}
