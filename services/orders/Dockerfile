# Build stage
FROM golang:1.24.6-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

# Copy go mod and sum files
COPY services/orders/go.mod services/orders/go.sum ./
COPY proto proto
COPY go.work go.work.sum ./
COPY services services
RUN go work sync

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o orderservice ./services/orders/main.go

# Use a smaller image for the final container
FROM alpine:3.18

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata postgresql-client

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/orderservice .

# Copy migration files
COPY --from=builder /app/migrations ./migrations

# Create a non-root user
RUN adduser -D -g '' appuser && \
    chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose the gRPC and health check ports
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/healthz || exit 1

# Command to run the executable
CMD ["/app/orderservice"]
