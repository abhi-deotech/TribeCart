package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	v1 "github.com/tribecart/proto/tribecart/v1"
	"tribecart/orders/db"
	"tribecart/orders/models"
	"tribecart/orders/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	_ "github.com/lib/pq"
)

// server implements the OrderServiceServer interface
type server struct {
	v1.UnimplementedOrderServiceServer
	service *service.OrderService
}

// NewServer creates a new gRPC server with the given order service
func NewServer(svc *service.OrderService) *server {
	return &server{
		service: svc,
	}
}

func (s *server) CreateOrder(ctx context.Context, req *v1.CreateOrderRequest) (*v1.Order, error) {
	return s.service.CreateOrder(ctx, req)
}

func (s *server) GetOrder(ctx context.Context, req *v1.GetOrderRequest) (*v1.Order, error) {
	return s.service.GetOrder(ctx, req)
}

func (s *server) ListOrders(ctx context.Context, req *v1.ListOrdersRequest) (*v1.ListOrdersResponse, error) {
	return s.service.ListOrders(ctx, req)
}

func (s *server) UpdateOrderStatus(ctx context.Context, req *v1.UpdateOrderStatusRequest) (*v1.Order, error) {
	return s.service.UpdateOrderStatus(ctx, req)
}

func (s *server) CancelOrder(ctx context.Context, req *v1.CancelOrderRequest) (*v1.Order, error) {
	return s.service.CancelOrder(ctx, req)
}

func (s *server) GetOrderHistory(ctx context.Context, req *v1.GetOrderHistoryRequest) (*v1.OrderHistoryResponse, error) {
	return s.service.GetOrderHistory(ctx, req)
}

func loadConfig() (*db.Config, error) {
	cfg := &db.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("POSTGRES_USER", "postgres"),
		Password: getEnv("POSTGRES_PASSWORD", "postgres"),
		DBName:   getEnv("POSTGRES_DB", "tribecart_orders"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// Validate required fields
	if cfg.Host == "" || cfg.Port == "" || cfg.User == "" || cfg.DBName == "" {
		return nil, fmt.Errorf("missing required database configuration")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func runMigrations(db *sql.DB) error {
	migrator, err := db.NewMigrator()
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	// Run migrations
	if err := migrator.Up(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func newDB() (*sql.DB, error) {
	// Load configuration
	cfg, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Connect to database
	db, err := db.Connect(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run migrations
	if err := runMigrations(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

func main() {
	// Initialize logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Create context that listens for the interrupt signal
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Initialize database
	db, err := newDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	log.Println("Database connection established and migrations run successfully")

	// Create repository and service
	repo := models.NewDBOrderRepository(db)
	svc := service.NewOrderService(repo)

	// Create gRPC server
	srv := grpc.NewServer(
		grpc.UnaryInterceptor(serverInterceptor),
	)

	// Register services
	v1.RegisterOrderServiceServer(srv, NewServer(svc))
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(srv, healthServer)
	reflection.Register(srv)

	// Start gRPC server
	port := getEnv("PORT", "8080")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Starting gRPC server on port %s", port)

	// Set health check status to serving
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	// Start server in a goroutine
	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Listen for interrupt signal
	<-ctx.Done()

	log.Println("Shutting down server...")

	// Graceful shutdown
	done := make(chan bool)
	go func() {
		srv.GracefulStop()
		close(done)
	}()

	// Wait for graceful shutdown or timeout
	select {
	case <-done:
		log.Println("Server stopped gracefully")
	case <-time.After(10 * time.Second):
		log.Println("Forcing server to stop")
		srv.Stop()
	}
}

func serverInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	// Log the incoming request
	log.Printf("Received %s request", info.FullMethod)

	// Record start time
	start := time.Now()

	// Call the handler
	resp, err := handler(ctx, req)

	// Log the result
	duration := time.Since(start)

	if err != nil {
		// Extract gRPC status if available
		if st, ok := status.FromError(err); ok {
			log.Printf("Request %s completed with status %s (%s) in %v",
				info.FullMethod, st.Code(), st.Message(), duration)
		} else {
			log.Printf("Request %s failed: %v (%v)", info.FullMethod, err, duration)
		}
	} else {
		log.Printf("Request %s completed successfully in %v", info.FullMethod, duration)
	}

	return resp, err
}
