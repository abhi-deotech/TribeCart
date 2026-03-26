.PHONY: build test run clean

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=orderservice

# Build the application
build:
	$(GOBUILD) -o bin/$(BINARY_NAME) -v

# Run tests
test:
	$(GOTEST) -v ./... -coverprofile=coverage.out

# Run the application
run:
	$(GOBUILD) -o bin/$(BINARY_NAME) -v
	./bin/$(BINARY_NAME)

# Clean build files
clean:
	$(GOCLEAN)
	rm -f bin/$(BINARY_NAME)

# Generate protobuf code
proto:
	protoc --go_out=. --go-grpc_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_opt=paths=source_relative \
    -I../../proto \
    ../../proto/tribecart/v1/orders.proto

# Install dependencies
deps:
	$(GOGET) google.golang.org/grpc
	$(GOGET) google.golang.org/protobuf
	$(GOGET) github.com/stretchr/testify/assert

# Run with live reload using air (install with: go install github.com/cosmtrek/air@latest)
dev:
	air

# Run linter
lint:
	golangci-lint run

# Run all checks
check: test lint
