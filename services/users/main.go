package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"time"

	pb "github.com/tribecart/proto/tribecart/v1"
)

type server struct {
	pb.UnimplementedUserServiceServer
	db *sql.DB
}

func NewServer(db *sql.DB) *server {
	return &server{
		db: db,
	}
}

func (s *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	user := &pb.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		// Password is not in the User message in proto, so we only handle it for the DB
	}

	query := "INSERT INTO users (first_name, last_name, email, password) VALUES ($1, $2, $3, $4) RETURNING id"
	err := s.db.QueryRowContext(ctx, query, user.FirstName, user.LastName, user.Email, req.Password).Scan(&user.Id)
	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	log.Printf("User created: %v", user)
	return user, nil
}

func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	user := &pb.User{}
	query := "SELECT id, first_name, last_name, email FROM users WHERE id = $1"
	err := s.db.QueryRowContext(ctx, query, req.Id).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	log.Printf("User retrieved: %v", user)
	return user, nil
}

func initializeDatabase(db *sql.DB) error {
	// This line is inserted before createTableSQL by Serena
	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}
	log.Println("Users table ensured.")
	// New line inserted by Serena
	return nil
}

func main() {
	// Database connection
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbUser := os.Getenv("POSTGRES_USER")
		dbPassword := os.Getenv("POSTGRES_PASSWORD")
		dbName := os.Getenv("POSTGRES_DB")

		connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			dbHost, dbPort, dbUser, dbPassword, dbName)
	}

	// Retry logic for database connection (DNS propagation is sometimes slow on Render)
	var db *sql.DB
	var err error
	for i := 0; i < 5; i++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			err = db.Ping()
			if err == nil {
				log.Println("Successfully connected to database!")
				break
			}
		}
		log.Printf("Waiting for database... attempt %d/5: %v", i+1, err)
		time.Sleep(5 * time.Second)
		if i == 4 {
			log.Fatalf("failed to connect to database after 5 attempts: %v", err)
		}
	}
	defer db.Close()

	// Initialize database schema
	err = initializeDatabase(db)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	// gRPC server setup
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, NewServer(db))
	reflection.Register(s)
	log.Println("Server listening at", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
