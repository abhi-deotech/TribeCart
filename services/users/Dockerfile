# build stage
FROM golang:1.24-alpine AS build

WORKDIR /app

# Install build dependencies for CGO and PostgreSQL
RUN apk add --no-cache gcc musl-dev libpq-dev git

# Copy the entire monorepo workspace
COPY . .

# Run go mod download from the root of the workspace
RUN go mod download

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/server ./services/users/main.go

# final stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies for PostgreSQL
RUN apk add --no-cache libpq

COPY --from=build /app/server .

EXPOSE 8080

CMD ["./server"]
