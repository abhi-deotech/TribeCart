# Build stage
FROM golang:1.24.6-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev git

# Copy go mod files
COPY services/products/go.mod services/products/go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags='-w -s' -o /app/bin/products-service ./services/products/main.go

# Final stage
FROM alpine:3.19

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy the binary from builder
COPY --from=builder /app/bin/products-service .

# Copy configuration files
COPY --from=builder /app/.env.example .env

# Copy protobuf files
COPY --from=builder /app/proto/tribecart/v1/*.proto ./proto/tribecart/v1/

# Expose gRPC port
EXPOSE 50051

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:50051/healthz || exit 1

# Run the application
CMD ["/app/products-service"]