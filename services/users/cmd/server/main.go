package main

import (
	"context"
	"crypto/rsa"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/tribecart/users/internal/auth"
	"github.com/tribecart/users/internal/repository"
	"github.com/tribecart/users/internal/service"
	pb "github.com/tribecart/proto/tribecart/v1"
)

const (
	serviceName = "users-service"
	dbDriver   = "pgx"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using environment variables")
	}

	// Initialize logger
	logger := log.New(os.Stdout, "[users-service] ", log.LstdFlags|log.Lshortfile)

	// Load configuration
	cfg, err := loadConfig()
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	db, err := initDB(cfg)
	if err != nil {
		logger.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Run database migrations
	if err := runMigrations(db); err != nil {
		logger.Fatalf("Failed to run database migrations: %v", err)
	}

	// Initialize JWT keys
	privateKey, publicKey, err := loadJWTKeys(cfg)
	if err != nil {
		logger.Fatalf("Failed to load JWT keys: %v", err)
	}

	// Initialize repository
	userRepo := repository.NewPostgresUserRepository(db)

	// Initialize services
	userSvc := service.NewUserService(
		userRepo,
		privateKey,
		publicKey,
		cfg.JWTSecret,
	)

	// Create gRPC server with middleware
	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			auth.UnaryServerInterceptor(publicKey, cfg.JWTSecret),
			// Add other interceptors here (logging, metrics, etc.)
		),
	)

	// Register services
	pb.RegisterUserServiceServer(srv, userSvc)

	// Enable reflection for gRPC CLI tools (like grpcurl)
	reflection.Register(srv)

	// Start gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		logger.Fatalf("Failed to listen: %v", err)
	}

	logger.Printf("Starting %s on %s", serviceName, lis.Addr())

	// Graceful shutdown
	go func() {
		if err := srv.Serve(lis); err != nil {
			logger.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Println("Shutting down server...")

	// Gracefully stop the server
	srv.GracefulStop()

	logger.Println("Server stopped")
}

// Config holds the application configuration
type Config struct {
	// Database
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// Server
	GRPCPort int

	// JWT
	JWTPrivateKeyPath string
	JWTPublicKeyPath  string
	JWTSecret         string
}

// loadConfig loads configuration from environment variables
func loadConfig() (*Config, error) {
	cfg := &Config{
		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnvAsInt("DB_PORT", 5432),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "tribecart_users"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		// Server
		GRPCPort: getEnvAsInt("GRPC_PORT", 50051),

		// JWT
		JWTPrivateKeyPath: getEnv("JWT_PRIVATE_KEY_PATH", "config/keys/jwtRS256.key"),
		JWTPublicKeyPath:  getEnv("JWT_PUBLIC_KEY_PATH", "config/keys/jwtRS256.key.pub"),
		JWTSecret:         getEnv("JWT_SECRET", "change-me-to-a-secure-secret"),
	}

	return cfg, nil
}

// initDB initializes the database connection
func initDB(cfg *Config) (*sql.DB, error) {
	// Build connection string
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	// Open database connection
	db, err := sql.Open(dbDriver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

// runMigrations runs database migrations
func runMigrations(db *sql.DB) error {
	// In a real application, you would use a migration library like golang-migrate
	// For now, we'll just log that migrations would run here
	log.Println("Running database migrations...")
	return nil
}

// loadJWTKeys loads the JWT private and public keys
func loadJWTKeys(cfg *Config) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	// In a real application, you would load these from files or environment variables
	// For now, we'll generate a new key pair if they don't exist
	privateKey, err := auth.LoadOrGeneratePrivateKey(cfg.JWTPrivateKeyPath, 4096)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load private key: %w", err)
	}

	publicKey, err := auth.LoadOrGeneratePublicKey(cfg.JWTPublicKeyPath, privateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load public key: %w", err)
	}

	return privateKey, publicKey, nil
}

// Helper functions for environment variables
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	var value int
	_, err := fmt.Sscanf(valueStr, "%d", &value)
	if err != nil {
		return defaultValue
	}

	return value
}
