package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	pb "github.com/tribecart/proto/tribecart/v1"
	"github.com/tribecart/services/products/internal/repository"
	"github.com/tribecart/services/products/internal/service"
)

type server struct {
	pb.UnimplementedProductServiceServer
	db *sql.DB
}

func NewServer(db *sql.DB) *server {
	return &server{
		db: db,
	}
}

func (s *server) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	product := &pb.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	}

	query := "INSERT INTO products (name, description, price) VALUES ($1, $2, $3) RETURNING id"
	err := s.db.QueryRowContext(ctx, query, product.Name, product.Description, product.Price).Scan(&product.Id)
	if err != nil {
		log.Printf("Failed to insert product: %v", err)
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	log.Printf("Product created: %v", product)
	return &pb.CreateProductResponse{Product: product}, nil
}

func (s *server) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	product := &pb.Product{}
	query := "SELECT id, name, description, price FROM products WHERE id = $1"
	err := s.db.QueryRowContext(ctx, query, req.Id).Scan(&product.Id, &product.Name, &product.Description, &product.Price)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("product not found")
	}
	if err != nil {
		log.Printf("Failed to get product: %v", err)
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	log.Printf("Product retrieved: %v", product)
	return &pb.GetProductResponse{Product: product}, nil
}

func (s *server) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	product := req.Product
	query := "UPDATE products SET name = $1, description = $2, price = $3 WHERE id = $4"
	_, err := s.db.ExecContext(ctx, query, product.Name, product.Description, product.Price, product.Id)
	if err != nil {
		log.Printf("Failed to update product: %v", err)
		return nil, fmt.Errorf("failed to update product: %w", err)
	}
	log.Printf("Product updated: %v", product)
	return &pb.UpdateProductResponse{Product: product}, nil
}

func (s *server) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	query := "DELETE FROM products WHERE id = $1"
	res, err := s.db.ExecContext(ctx, query, req.Id)
	if err != nil {
		log.Printf("Failed to delete product: %v", err)
		return nil, fmt.Errorf("failed to delete product: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("Failed to get rows affected: %v", err)
		return nil, fmt.Errorf("failed to delete product: %w", err)
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("product not found")
	}
	log.Printf("Product deleted: %s", req.Id)
	return &pb.DeleteProductResponse{Success: true}, nil
}

func (s *server) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	var products []*pb.Product
	query := "SELECT id, name, description, price FROM products"
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("Failed to list products: %v", err)
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		product := &pb.Product{}
		if err := rows.Scan(&product.Id, &product.Name, &product.Description, &product.Price); err != nil {
			log.Printf("Failed to scan product: %v", err)
			return nil, fmt.Errorf("failed to list products: %w", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error during rows iteration: %v", err)
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	log.Printf("Listed %d products", len(products))
	return &pb.ListProductsResponse{Products: products}, nil
}

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	GRPCPort   string
}

func loadConfig() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "tribecart_products"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		GRPCPort:   getEnv("GRPC_PORT", "50051"),
	}, nil
}

func initDB(cfg *Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to database")
	return db, nil
}

func runMigrations(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS products (
			id VARCHAR(36) PRIMARY KEY,
			seller_id VARCHAR(36),
			name VARCHAR(255) NOT NULL,
			description TEXT,
			sku VARCHAR(100),
			barcode VARCHAR(100),
			price DECIMAL(10, 2) NOT NULL,
			sale_price DECIMAL(10, 2),
			cost_price DECIMAL(10, 2),
			stock_quantity INTEGER NOT NULL DEFAULT 0,
			track_inventory BOOLEAN NOT NULL DEFAULT true,
			min_stock_level INTEGER,
			weight DECIMAL(10, 2),
			length DECIMAL(10, 2),
			width DECIMAL(10, 2),
			height DECIMAL(10, 2),
			type INTEGER NOT NULL DEFAULT 1, -- 1: physical, 2: digital, 3: service
			status INTEGER NOT NULL DEFAULT 1, -- 1: draft, 2: active, 3: archived, 4: out of stock, 5: discontinued
			is_featured BOOLEAN NOT NULL DEFAULT false,
			is_visible BOOLEAN NOT NULL DEFAULT true,
			requires_shipping BOOLEAN NOT NULL DEFAULT true,
			is_taxable BOOLEAN NOT NULL DEFAULT true,
			tax_class_id VARCHAR(36),
			seo_title VARCHAR(255),
			seo_description TEXT,
			seo_keywords VARCHAR(255),
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMP WITH TIME ZONE
		);

		CREATE INDEX IF NOT EXISTS idx_products_seller_id ON products(seller_id);
		CREATE INDEX IF NOT EXISTS idx_products_sku ON products(sku);
		CREATE INDEX IF NOT EXISTS idx_products_barcode ON products(barcode);
		CREATE INDEX IF NOT EXISTS idx_products_status ON products(status);
		CREATE INDEX IF NOT EXISTS idx_products_created_at ON products(created_at);
	`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS product_categories (
			product_id VARCHAR(36) NOT NULL,
			category_id VARCHAR(36) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			PRIMARY KEY (product_id, category_id)
		);

		CREATE INDEX IF NOT EXISTS idx_product_categories_product_id ON product_categories(product_id);
		CREATE INDEX IF NOT EXISTS idx_product_categories_category_id ON product_categories(category_id);
	`)
	if err != nil {
		return fmt.Errorf("failed to create product_categories table: %w", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS stock_movements (
			id VARCHAR(36) PRIMARY KEY,
			product_id VARCHAR(36) NOT NULL,
			variant_id VARCHAR(36),
			quantity INTEGER NOT NULL,
			reference_id VARCHAR(100),
			reason VARCHAR(100),
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
		);

		CREATE INDEX IF NOT EXISTS idx_stock_movements_product_id ON stock_movements(product_id);
		CREATE INDEX IF NOT EXISTS idx_stock_movements_variant_id ON stock_movements(variant_id);
	`)
	if err != nil {
		return fmt.Errorf("failed to create stock_movements table: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := initDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	if err := runMigrations(db); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	repo := repository.NewPostgresProductRepository(db)
	productSvc := service.NewProductService(repo)

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(UnaryServerInterceptor()),
	}

	srv := grpc.NewServer(opts...)

	pb.RegisterProductServiceServer(srv, productSvc)

	healthSrv := health.NewServer()
	healthSrv.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(srv, healthSrv)

	reflection.Register(srv)

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.GRPCPort, err)
	}

	log.Printf("Starting gRPC server on port %s...", cfg.GRPCPort)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	<-quit

	log.Println("Shutting down server...")

	srv.GracefulStop()
	healthSrv.SetServingStatus("", grpc_health_v1.HealthCheckResponse_NOT_SERVING)

	log.Println("Server stopped")
}

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		log.Printf("gRPC method %s called", info.FullMethod)

		resp, err := handler(ctx, req)

		if err != nil {
			log.Printf("gRPC method %s failed: %v", info.FullMethod, err)
		}

		return resp, err
	}
}
