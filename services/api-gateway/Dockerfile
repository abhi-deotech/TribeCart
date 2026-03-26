# build stage
FROM golang:1.24-alpine AS build

WORKDIR /app

# Install git
RUN apk add --no-cache git

# Copy the entire monorepo workspace
COPY . .

# Run go mod download from the root of the workspace
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./services/api-gateway/main.go

# final stage
FROM alpine:latest

WORKDIR /app

COPY --from=build /app/server .

EXPOSE 8000

CMD ["./server"]
