#!/bin/bash

set -e

export PATH=$PATH:$(go env GOPATH)/bin

# Create necessary directories
mkdir -p google/protobuf

# Download required proto files
PROTOBUF_VERSION=3.19.4
if [ ! -f "google/protobuf/field_mask.proto" ]; then
  curl -sSL https://raw.githubusercontent.com/protocolbuffers/protobuf/v${PROTOBUF_VERSION}/src/google/protobuf/field_mask.proto -o google/protobuf/field_mask.proto
  curl -sSL https://raw.githubusercontent.com/protocolbuffers/protobuf/v${PROTOBUF_VERSION}/src/google/protobuf/timestamp.proto -o google/protobuf/timestamp.proto
  curl -sSL https://raw.githubusercontent.com/protocolbuffers/protobuf/v${PROTOBUF_VERSION}/src/google/protobuf/empty.proto -o google/protobuf/empty.proto
fi

# Generate Go code
protoc \
  -I. \
  -I$(go env GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis \
  --go_out=. \
  --go_opt=module=github.com/tribecart/proto \
  --go_opt=Mgoogle/protobuf/any.proto=google.golang.org/protobuf/types/known/anypb \
  --go_opt=Mgoogle/protobuf/timestamp.proto=google.golang.org/protobuf/types/known/timestamppb \
  --go_opt=Mgoogle/protobuf/empty.proto=google.golang.org/protobuf/types/known/emptypb \
  --go_opt=Mgoogle/protobuf/field_mask.proto=google.golang.org/protobuf/types/known/fieldmaskpb \
  --go-grpc_out=. \
  --go-grpc_opt=module=github.com/tribecart/proto \
  --go-grpc_opt=require_unimplemented_servers=false \
  tribecart/v1/*.proto

# Fix import paths in generated files
find tribecart/v1 -name "*.pb.go" -exec sed -i 's|"google/protobuf/field_mask.proto"|"google.golang.org/protobuf/types/known/fieldmaskpb"|g' {} \;

echo "Protobuf files generated successfully"
