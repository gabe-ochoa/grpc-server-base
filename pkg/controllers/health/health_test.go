package controllers

import (
	"testing"

	healthv1 "github.com/grpc-serverchat/grpc-server-server/protos/gen/v1/health"
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

func TestConnectivityCheck(t *testing.T) {
	testServer := NewServer()
	testCtx := context.Background()
	request := healthv1.ConnectivityCheckRequest{}

	response, err := testServer.ConnectivityCheck(testCtx, &request)

	assert.NoError(t, err)

	assert.Equal(t, true, response.Success)
}
