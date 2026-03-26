# Products Service

The Products Service is a gRPC-based microservice responsible for managing products in the TribeCart e-commerce platform.

## Features

- Create, read, update, and delete products
- Manage product inventory and stock levels
- Support for product variants and categories
- Search and filter products
- Bulk import/export functionality
- Support for digital and physical products

## Prerequisites

- Go 1.21 or later
- PostgreSQL 13+
- Docker and Docker Compose (optional, for containerized deployment)

## Getting Started

### Local Development

1. **Clone the repository**
   ```bash
   git clone https://github.com/tribecart/tribecart.git
   cd services/products
   ```

2. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Start dependencies**
   ```bash
   docker-compose up -d postgres
   ```

4. **Run migrations**
   ```bash
   go run main.go migrate
   ```

5. **Start the service**
   ```bash
   go run main.go
   ```

### Using Docker Compose

```bash
docker-compose up -d
```

This will start:
- Products Service (gRPC on port 50051)
- PostgreSQL database
- pgAdmin (available at http://localhost:5050)

## API Documentation

The service implements the gRPC service defined in `proto/tribecart/v1/products.proto`.

### Example gRPC Calls

#### Create a Product

```protobuf
rpc CreateProduct(CreateProductRequest) returns (Product) {}
```

#### Get a Product

```protobuf
rpc GetProduct(GetProductRequest) returns (Product) {}
```

#### Update a Product

```protobuf
rpc UpdateProduct(UpdateProductRequest) returns (Product) {}
```

#### List Products

```protobuf
rpc ListProducts(ListProductsRequest) returns (ListProductsResponse) {}
```

### Health Check

The service includes a gRPC health check endpoint that can be used by container orchestration systems.

## Development

### Running Tests

```bash
go test -v ./...
```

### Code Generation

After updating the protobuf definitions, regenerate the Go code:

```bash
cd ../../proto
./generate.sh
```

### Linting

```bash
golangci-lint run
```

## Deployment

### Building the Docker Image

```bash
docker build -t tribecart/products-service:latest .
```

### Kubernetes

See the `deploy/kubernetes` directory for example Kubernetes manifests.

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| DB_HOST | localhost | PostgreSQL host |
| DB_PORT | 5432 | PostgreSQL port |
| DB_USER | postgres | Database user |
| DB_PASSWORD | postgres | Database password |
| DB_NAME | tribecart_products | Database name |
| DB_SSLMODE | disable | SSL mode for database connection |
| GRPC_PORT | 50051 | gRPC server port |
| LOG_LEVEL | info | Log level (debug, info, warn, error) |
| LOG_FORMAT | json | Log format (json, text) |

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
