# Orders Service

The Orders service is a gRPC-based microservice responsible for managing orders in the TribeCart e-commerce platform. It handles order creation, retrieval, status updates, and history tracking.

## Features

- Create, retrieve, update, and cancel orders
- Track order status history
- Support for order items with quantities and pricing
- Integration with payment and shipping services
- Comprehensive error handling and validation
- Database migrations for schema management

## Prerequisites

- Go 1.21 or later
- PostgreSQL 14 or later
- Docker and Docker Compose (for containerized deployment)

## Configuration

The service can be configured using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `POSTGRES_USER` | Database user | `postgres` |
| `POSTGRES_PASSWORD` | Database password | `postgres` |
| `POSTGRES_DB` | Database name | `tribecart_orders` |
| `DB_SSLMODE` | SSL mode for database connection | `disable` |
| `PORT` | gRPC server port | `8080` |

## Local Development

### Database Setup

1. Start the PostgreSQL database using Docker Compose:
   ```bash
   docker-compose up -d postgres
   ```

2. Create the database (if it doesn't exist):
   ```bash
   createdb -h localhost -U postgres tribecart_orders
   ```

### Running the Service

1. Install dependencies:
   ```bash
   go mod download
   ```

2. Run migrations:
   ```bash
   go run cmd/migrate/main.go
   ```

3. Start the service:
   ```bash
   go run main.go
   ```

### Using Docker

Build and run the service using Docker Compose:

```bash
docker-compose up -d --build orders-service
```

## API Documentation

The service implements the following gRPC service:

```protobuf
service OrderService {
  rpc CreateOrder(CreateOrderRequest) returns (Order);
  rpc GetOrder(GetOrderRequest) returns (Order);
  rpc ListOrders(ListOrdersRequest) returns (ListOrdersResponse);
  rpc UpdateOrderStatus(UpdateOrderStatusRequest) returns (Order);
  rpc CancelOrder(CancelOrderRequest) returns (Order);
  rpc GetOrderHistory(GetOrderHistoryRequest) returns (OrderHistoryResponse);
}
```

For detailed message definitions, see the [protobuf definitions](../../proto/tribecart/v1/order_service.proto).

## Database Schema

The service uses the following database schema:

### Tables

#### orders
- `id` (UUID, PK): Unique order identifier
- `user_id` (UUID): ID of the user who placed the order
- `status` (enum): Current status of the order
- `subtotal` (decimal): Order subtotal (before tax and shipping)
- `tax_amount` (decimal): Tax amount
- `shipping_cost` (decimal): Shipping cost
- `total_amount` (decimal): Total order amount (subtotal + tax + shipping)
- `payment_method` (enum): Payment method used
- `payment_id` (string): External payment ID
- `shipping_address_id` (UUID): ID of the shipping address
- `tracking_number` (string): Shipping tracking number
- `carrier` (string): Shipping carrier
- `cancelled_at` (timestamp): When the order was cancelled (if applicable)
- `cancelled_reason` (text): Reason for cancellation
- `metadata` (jsonb): Additional order metadata
- `created_at` (timestamp): When the order was created
- `updated_at` (timestamp): When the order was last updated

#### order_items
- `id` (UUID, PK): Unique item identifier
- `order_id` (UUID, FK): Reference to the order
- `product_id` (UUID): ID of the product
- `quantity` (integer): Quantity ordered
- `unit_price` (decimal): Price per unit
- `subtotal` (decimal): Item subtotal (quantity * unit_price)
- `tax_amount` (decimal): Tax amount for this item
- `total_amount` (decimal): Total amount for this item (subtotal + tax)
- `metadata` (jsonb): Additional item metadata
- `created_at` (timestamp): When the item was created
- `updated_at` (timestamp): When the item was last updated

#### order_status_history
- `id` (UUID, PK): Unique history entry ID
- `order_id` (UUID, FK): Reference to the order
- `status` (enum): Status at this point in history
- `notes` (text): Additional notes about the status change
- `metadata` (jsonb): Additional metadata
- `created_at` (timestamp): When the status change occurred

## Testing

Run the test suite:

```bash
go test -v ./...
```

## Deployment

The service is designed to be deployed as a container. A sample Dockerfile is provided in the repository root.

## Monitoring and Observability

The service exposes the following endpoints for monitoring:

- `/healthz`: Health check endpoint
- `/metrics`: Prometheus metrics
- `/debug/pprof`: Go pprof endpoints

## License

This project is licensed under the MIT License - see the [LICENSE](../../LICENSE) file for details.
