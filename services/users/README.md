# Users Service

The Users Service is responsible for managing user accounts, authentication, and authorization in the TribeCart platform.

## Features

- User registration and profile management
- JWT-based authentication
- Role-based access control (RBAC)
- Email verification
- Password reset functionality
- User sessions management
- Audit logging

## Prerequisites

- Go 1.21 or later
- PostgreSQL 14 or later
- Make (optional, for development)

## Configuration

Create a `.env` file in the project root with the following variables:

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=tribecart_users
DB_SSLMODE=disable

# Server
GRPC_PORT=50051
HTTP_PORT=8080

# JWT
JWT_SECRET=your-secret-key
JWT_PRIVATE_KEY_PATH=config/keys/jwtRS256.key
JWT_PUBLIC_KEY_PATH=config/keys/jwtRS256.key.pub
```

## Database Setup

1. Create a PostgreSQL database:

```bash
createdb tribecart_users
```

2. Run migrations:

```bash
make migrate-up
```

## Running the Service

### Using Make (recommended for development)

```bash
make run
```

### Using Go directly

```bash
go run cmd/server/main.go
```

## API Documentation

The service exposes the following gRPC endpoints:

### User Management

- `CreateUser` - Register a new user
- `GetUser` - Get user details
- `UpdateUser` - Update user profile
- `DeleteUser` - Delete a user
- `ListUsers` - List all users (admin only)

### Authentication

- `Login` - Authenticate a user
- `Logout` - Invalidate a user's session
- `RefreshToken` - Get a new access token using a refresh token
- `ChangePassword` - Change a user's password
- `ForgotPassword` - Initiate password reset
- `ResetPassword` - Reset password using a token
- `SendVerificationEmail` - Send email verification
- `VerifyEmail` - Verify email using a token

## Development

### Code Generation

After updating the protobuf definitions, regenerate the Go code:

```bash
make generate
```

### Testing

Run the test suite:

```bash
make test
```

### Linting

Run the linter:

```bash
make lint
```

### Formatting

Format the code:

```bash
make fmt
```

## Deployment

### Building the Docker Image

```bash
docker build -t tribecart/users-service:latest .
```

### Running with Docker Compose

```yaml
version: '3.8'

services:
  users-service:
    image: tribecart/users-service:latest
    environment:
      - DB_HOST=postgres
      - DB_USER=postgres
      - DB_PASSWORD=postgres
      - DB_NAME=tribecart_users
      - DB_SSLMODE=disable
      - JWT_SECRET=your-secret-key
    ports:
      - "50051:50051"
    depends_on:
      - postgres
    restart: unless-stopped

  postgres:
    image: postgres:14-alpine
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
      - POSTGRES_DB=tribecart_users
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    restart: unless-stopped

volumes:
  postgres_data:
```

## Security Considerations

- Always use HTTPS in production
- Rotate JWT keys periodically
- Use strong database passwords
- Keep dependencies up to date
- Follow the principle of least privilege for database users

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
