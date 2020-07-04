// +build tools

package tools

// These are imports for the tools we use to build grpc-server-base
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module

import (
	_ "github.com/fullstorydev/grpcurl/cmd/grpcurl"
	_ "github.com/golang/protobuf/protoc-gen-go"
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway"
	_ "github.com/vektra/mockery"
)
