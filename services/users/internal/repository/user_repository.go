package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/tribecart/proto/tribecart/v1"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// User management
	CreateUser(ctx context.Context, user *pb.User) error
	GetUserByID(ctx context.Context, id string) (*pb.User, error)
	GetUserByEmail(ctx context.Context, email string) (*pb.User, error)
	UpdateUser(ctx context.Context, user *pb.User) error
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context, page, pageSize int32, filter string) ([]*pb.User, int32, error)

	// Authentication
	UpdateUserPassword(ctx context.Context, userID, hashedPassword string) error
	GetUserByCredentials(ctx context.Context, email, hashedPassword string) (*pb.User, error)

	// Address management
	CreateAddress(ctx context.Context, userID string, address *pb.Address) error
	GetAddress(ctx context.Context, id string) (*pb.Address, error)
	ListAddresses(ctx context.Context, userID string) ([]*pb.Address, error)
	UpdateAddress(ctx context.Context, address *pb.Address) error
	DeleteAddress(ctx context.Context, id string) error
}

// PostgresUserRepository implements UserRepository for PostgreSQL
type PostgresUserRepository struct {
	db *sql.DB
}

// NewPostgresUserRepository creates a new PostgresUserRepository
func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

// CreateUser creates a new user in the database
func (r *PostgresUserRepository) CreateUser(ctx context.Context, user *pb.User) error {
	query := `
		INSERT INTO users (
			id, first_name, last_name, email, password, phone_number, 
			role, status, email_verified, phone_verified, metadata, created_at, updated_at
		) VALUES ($1, $2, $3, LOWER($4), $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.Id,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password, // This should be the hashed password
		sql.NullString{String: user.PhoneNumber, Valid: user.PhoneNumber != ""},
		user.Role.String(),
		user.Status.String(),
		user.EmailVerified,
		user.PhoneVerified,
		nil, // TODO: Implement metadata handling
		time.Now(),
		time.Now(),
	)

	if err != nil {
		// Handle duplicate email error
		if strings.Contains(err.Error(), "users_email_key") {
			return fmt.Errorf("email already exists: %w", err)
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUserByID retrieves a user by ID
func (r *PostgresUserRepository) GetUserByID(ctx context.Context, id string) (*pb.User, error) {
	query := `
		SELECT 
			id, first_name, last_name, email, phone_number, role, status,
			email_verified, phone_verified, last_login_at, created_at, updated_at
		FROM users 
		WHERE id = $1 AND status != $2
	`

	var user pb.User
	var roleStr, statusStr string
	var lastLoginAt, createdAt, updatedAt time.Time
	var phoneNumber sql.NullString

	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
		pb.UserStatus_USER_STATUS_DELETED.String(),
	).Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&phoneNumber,
		&roleStr,
		&statusStr,
		&user.EmailVerified,
		&user.PhoneVerified,
		&lastLoginAt,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Convert string enums to protobuf enums
	user.Role = pb.UserRole(pb.UserRole_value[roleStr])
	user.Status = pb.UserStatus(pb.UserStatus_value[statusStr])
	user.PhoneNumber = phoneNumber.String
	user.LastLoginAt = timestamppb.New(lastLoginAt)
	user.CreatedAt = timestamppb.New(createdAt)
	user.UpdatedAt = timestamppb.New(updatedAt)

	return &user, nil
}

// GetUserByEmail retrieves a user by email
func (r *PostgresUserRepository) GetUserByEmail(ctx context.Context, email string) (*pb.User, error) {
	query := `
		SELECT 
			id, first_name, last_name, email, password, phone_number, role, status,
			email_verified, phone_verified, last_login_at, created_at, updated_at
		FROM users 
		WHERE email = LOWER($1) AND status != $2
	`

	var user pb.User
	var roleStr, statusStr string
	var lastLoginAt, createdAt, updatedAt time.Time
	var phoneNumber sql.NullString

	err := r.db.QueryRowContext(
		ctx,
		query,
		email,
		pb.UserStatus_USER_STATUS_DELETED.String(),
	).Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password, // Include password for authentication
		&phoneNumber,
		&roleStr,
		&statusStr,
		&user.EmailVerified,
		&user.PhoneVerified,
		&lastLoginAt,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	// Convert string enums to protobuf enums
	user.Role = pb.UserRole(pb.UserRole_value[roleStr])
	user.Status = pb.UserStatus(pb.UserStatus_value[statusStr])
	user.PhoneNumber = phoneNumber.String
	user.LastLoginAt = timestamppb.New(lastLoginAt)
	user.CreatedAt = timestamppb.New(createdAt)
	user.UpdatedAt = timestamppb.New(updatedAt)

	return &user, nil
}

// UpdateUser updates a user in the database
func (r *PostgresUserRepository) UpdateUser(ctx context.Context, user *pb.User) error {
	query := `
		UPDATE users 
		SET 
			first_name = $1,
			last_name = $2,
			phone_number = $3,
			role = $4,
			status = $5,
			email_verified = $6,
			phone_verified = $7,
			last_login_at = $8,
			updated_at = $9
		WHERE id = $10
		RETURNING updated_at
	`

	var lastLoginAt *time.Time
	if user.LastLoginAt != nil {
		t := user.LastLoginAt.AsTime()
		lastLoginAt = &t
	}

	var updatedAt time.Time
	err := r.db.QueryRowContext(
		ctx,
		query,
		user.FirstName,
		user.LastName,
		sql.NullString{String: user.PhoneNumber, Valid: user.PhoneNumber != ""},
		user.Role.String(),
		user.Status.String(),
		user.EmailVerified,
		user.PhoneVerified,
		lastLoginAt,
		time.Now(),
		user.Id,
	).Scan(&updatedAt)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	user.UpdatedAt = timestamppb.New(updatedAt)
	return nil
}

// DeleteUser soft deletes a user by setting status to DELETED
func (r *PostgresUserRepository) DeleteUser(ctx context.Context, id string) error {
	query := `
		UPDATE users 
		SET status = $1, updated_at = $2
		WHERE id = $3
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		pb.UserStatus_USER_STATUS_DELETED.String(),
		time.Now(),
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// UpdateUserPassword updates a user's password
func (r *PostgresUserRepository) UpdateUserPassword(ctx context.Context, userID, hashedPassword string) error {
	query := `
		UPDATE users 
		SET password = $1, updated_at = $2
		WHERE id = $3
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		hashedPassword,
		time.Now(),
		userID,
	)

	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// ListUsers retrieves a paginated list of users with optional filtering
func (r *PostgresUserRepository) ListUsers(
	ctx context.Context,
	page, pageSize int32,
	filter string,
) ([]*pb.User, int32, error) {
	// First, get the total count
	var totalCount int32
	countQuery := `SELECT COUNT(*) FROM users WHERE status != $1`
	if filter != "" {
		countQuery += " AND (first_name ILIKE $2 OR last_name ILIKE $2 OR email ILIKE $2)"
	}

	var err error
	if filter != "" {
		filterPattern := "%" + filter + "%"
		err = r.db.QueryRowContext(
			ctx,
			countQuery,
			pb.UserStatus_USER_STATUS_DELETED.String(),
			filterPattern,
		).Scan(&totalCount)
	} else {
		err = r.db.QueryRowContext(
			ctx,
			countQuery,
			pb.UserStatus_USER_STATUS_DELETED.String(),
		).Scan(&totalCount)
	}

	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Then get the paginated results
	offset := (page - 1) * pageSize
	query := `
		SELECT 
			id, first_name, last_name, email, phone_number, role, status,
			email_verified, phone_verified, last_login_at, created_at, updated_at
		FROM users
		WHERE status != $1
	`

	args := []interface{}{pb.UserStatus_USER_STATUS_DELETED.String()}
	argIndex := 2

	if filter != "" {
		query += fmt.Sprintf(" AND (first_name ILIKE $%d OR last_name ILIKE $%d OR email ILIKE $%d)", 
			argIndex, argIndex, argIndex)
		args = append(args, "%"+filter+"%")
		argIndex++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, pageSize, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*pb.User
	for rows.Next() {
		var user pb.User
		var roleStr, statusStr string
		var lastLoginAt, createdAt, updatedAt time.Time
		var phoneNumber sql.NullString

		err := rows.Scan(
			&user.Id,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&phoneNumber,
			&roleStr,
			&statusStr,
			&user.EmailVerified,
			&user.PhoneVerified,
			&lastLoginAt,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}

		// Convert string enums to protobuf enums
		user.Role = pb.UserRole(pb.UserRole_value[roleStr])
		user.Status = pb.UserStatus(pb.UserStatus_value[statusStr])
		user.PhoneNumber = phoneNumber.String
		user.LastLoginAt = timestamppb.New(lastLoginAt)
		user.CreatedAt = timestamppb.New(createdAt)
		user.UpdatedAt = timestamppb.New(updatedAt)

		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating users: %w", err)
	}

	return users, totalCount, nil
}
