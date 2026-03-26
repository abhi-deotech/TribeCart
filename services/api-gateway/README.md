# API Gateway

The API Gateway is the entry point for all client requests in the TribeCart application. It handles request routing, authentication, and protocol translation between HTTP/REST and gRPC.

## Features

- HTTP/REST to gRPC protocol translation
- Request routing and load balancing
- CORS support
- Request/response validation
- Error handling and transformation
- Health check endpoint

## API Endpoints

### Health Check
- `GET /health` - Check if the API Gateway is running

### User Service
- `POST /api/v1/users/register` - Register a new user

### Order Service
- `POST /api/v1/orders` - Create a new order
- `GET /api/v1/orders` - List orders with optional filtering
- `GET /api/v1/orders/{id}` - Get order by ID
- `POST /api/v1/orders/{id}/cancel` - Cancel an order
- `PATCH /api/v1/orders/{id}/status` - Update order status
- `GET /api/v1/orders/user/{user_id}` - Get order history for a user

## Development

### Prerequisites

- Go 1.24 or later
- Protocol Buffers compiler (protoc)
- Go plugins for protoc

### Building

```bash
# Build the application
go build -o api-gateway
```

### Running Locally

```bash
# Run the API Gateway
./api-gateway
```

The API Gateway will be available at `http://localhost:8000`

### Environment Variables

- `USER_SERVICE_ADDR` - Address of the User gRPC service (default: `users-service:8080`)
- `ORDER_SERVICE_ADDR` - Address of the Order gRPC service (default: `orders-service:8080`)
- `PORT` - Port to run the HTTP server on (default: `8000`)

## Testing

Run the tests with:

```bash
go test -v ./...
```

## Deployment

The API Gateway is containerized and can be deployed using the provided Dockerfile or as part of the docker-compose setup.

### Building the Docker Image

```bash
docker build -t tribecart/api-gateway .
```

### Running with Docker Compose

The API Gateway is included in the main `docker-compose.yml` file. To start all services:

```bash
docker-compose up -d
```

## Error Handling

The API Gateway provides consistent error responses in the following format:

```json
{
  "error": "Error message describing the issue"
}
```

## Logging

The API Gateway logs all requests and errors to stdout. In production, these logs should be collected and analyzed using a logging service.
