# TODO move these commands into the package.json

filesToTest=`go list ./... | grep -v mocks | grep -v gen`
projectDir=$(shell pwd)
protocGenDoc=$(shell which protoc-gen-doc)
version=$(shell date +%Y.%m.%d.%H.%M.%S)

all: lint-docker test build-linux

.PHONY: build-linux
build-linux:
	GOOS=linux GOARCH=amd64 go build -o dist/grpc-server-api-linux

.PHONY: build-darwin
build-darwin:
	GOOS=darwin GOARCH=amd64 go build -o dist/grpc-server-api-darwin

.PHONY: build
build: protos
	go build -o dist/grpc-server-api -ldflags="-X 'main.Version=v$(version)'"

.PHONY: build-container
build-container:
	docker-compose build

.PHONY: clean
clean:
	rm -rf dist/

.PHONY: mod
mod:
	go mod tidy && go mod vendor

.PHONY: devtools
devtools:
	go get github.com/vektra/mockery/.../
	docker pull uber/prototool
	go install github.com/fullstorydev/grpcurl/cmd/grpcurl
	go get -u sourcegraph.com/sourcegraph/prototools/cmd/protoc-gen-doc

.PHONY: lint
lint: protos-lint
	golangci-lint run -v

.PHONY: lint-docker
lint-docker:
	docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:v1.24.0 golangci-lint run -v

# Test skips acceptance tests and generated files
.PHONY: test
test:
	go test -cover $(filesToTest)

# Coverage repor
.PHONY: coverage
coverage:
	go test -coverprofile=coverage.out $(filesToTest)
	go tool cover -html=coverage.out
	rm coverage.out

.PHONY: acceptance-test
acceptance-test:
	go test ./acceptancetests

## Local Testing

# Used for mimicing the production environment
.PHONY: start-api-container
start-api-container: build-linux
	docker-compose run --service-ports api

# Used for fast API development
.PHONY: start-api
start-api:
	LOG_LEVEL=trace \
	DB_PORT=54320 \
	DB_PASSWORD=pwd0123456789 \
	go run main.go

.PHONY: stop-local-server
stop-local-server:
	docker-compose down

.PHONY: start
start: build-linux docker-build
	docker-compose up

.PHONY: docker-build
docker-build:
	docker-compose build

.PHONY: protos
protos:
	docker run -v ${projectDir}/protos:/protos uber/prototool prototool all /protos

.PHONY: protos-lint
protos-lint:
	docker run -v ${projectDir}/protos:/protos uber/prototool prototool lint /protos
	docker run -v ${projectDir}/protos:/protos uber/prototool prototool format --overwrite /protos
