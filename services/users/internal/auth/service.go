package auth

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/tribecart/proto/tribecart/v1"
)

var (
	// ErrUserNotFound is returned when a user is not found
	ErrUserNotFound = errors.New("user not found")
	// ErrInvalidCredentials is returned when the provided credentials are invalid
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrEmailAlreadyExists is returned when a user with the given email already exists
	ErrEmailAlreadyExists = errors.New("email already exists")
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	CreateUser(ctx context.Context, user *pb.User) error
	GetUserByID(ctx context.Context, id string) (*pb.User, error)
	GetUserByEmail(ctx context.Context, email string) (*pb.User, error)
	UpdateUser(ctx context.Context, user *pb.User) error
	UpdateUserPassword(ctx context.Context, userID, hashedPassword string) error
}

// AuthService handles authentication and user management
type AuthService struct {
	repo      UserRepository
	jwtMgr    *JWTManager
	jwtSecret string
}

// NewAuthService creates a new authentication service
func NewAuthService(repo UserRepository, privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey, jwtSecret string) *AuthService {
	jwtMgr := NewJWTManager(
		privateKey,
		publicKey,
		time.Hour*24,       // Access token expires in 24 hours
		time.Hour*24*7,     // Refresh token expires in 7 days
		"tribecart-auth",   // Issuer
	)

	return &AuthService{
		repo:      repo,
		jwtMgr:    jwtMgr,
		jwtSecret: jwtSecret,
	}
}

// Register creates a new user account
func (s *AuthService) Register(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	// Check if user with this email already exists
	existingUser, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, status.Error(codes.AlreadyExists, ErrEmailAlreadyExists.Error())
	}

	// Hash the password
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}

	// Create the user
	now := timestamppb.Now()
	user := &pb.User{
		Id:         uuid.New().String(),
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Email:      req.Email,
		Password:   hashedPassword,
		PhoneNumber: req.PhoneNumber,
		Role:       pb.UserRole_USER_ROLE_CUSTOMER, // Default role
		Status:     pb.UserStatus_USER_STATUS_ACTIVE,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// Save the user to the database
	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	// Don't return the hashed password
	user.Password = ""
	return user, nil
}

// Login authenticates a user and returns JWT tokens
func (s *AuthService) Login(ctx context.Context, email, password string) (*pb.LoginResponse, error) {
	// Get the user by email
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, ErrInvalidCredentials.Error())
	}

	// Check if the account is active
	if user.Status != pb.UserStatus_USER_STATUS_ACTIVE {
		return nil, status.Error(codes.PermissionDenied, "account is not active")
	}

	// Verify the password
	valid, err := VerifyPassword(password, user.Password)
	if err != nil || !valid {
		return nil, status.Error(codes.Unauthenticated, ErrInvalidCredentials.Error())
	}

	// Generate JWT tokens
	accessToken, refreshToken, expiresAt, refreshExpiresAt, err := s.jwtMgr.GenerateTokenPair(
		user.Id,
		user.Email,
		user.Role.String(),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate tokens: %v", err)
	}

	// Update last login time
	user.LastLoginAt = timestamppb.Now()
	if err := s.repo.UpdateUser(ctx, user); err != nil {
		// Log the error but don't fail the login
		fmt.Printf("failed to update last login time: %v\n", err)
	}

	// Don't return the hashed password
	user.Password = ""

	return &pb.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
		ExpiresAt:    timestamppb.New(expiresAt),
	}, nil
}

// RefreshToken generates a new access token using a refresh token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*pb.RefreshTokenResponse, error) {
	// Validate the refresh token
	claims, err := s.jwtMgr.ValidateToken(refreshToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
	}

	// Check if it's a refresh token
	if claims.Type != RefreshToken {
		return nil, status.Error(codes.Unauthenticated, "not a refresh token")
	}

	// Get the user
	user, err := s.repo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "user not found")
	}

	// Check if the account is active
	if user.Status != pb.UserStatus_USER_STATUS_ACTIVE {
		return nil, status.Error(codes.PermissionDenied, "account is not active")
	}

	// Generate a new access token
	accessToken, expiresAt, err := s.jwtMgr.GenerateToken(
		user.Id,
		user.Email,
		user.Role.String(),
		AccessToken,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate access token: %v", err)
	}

	return &pb.RefreshTokenResponse{
		AccessToken: accessToken,
		ExpiresAt:   timestamppb.New(expiresAt),
	}, nil
}

// ChangePassword changes a user's password
func (s *AuthService) ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error {
	// Get the user
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return status.Error(codes.NotFound, "user not found")
	}

	// Verify the current password
	valid, err := VerifyPassword(currentPassword, user.Password)
	if err != nil || !valid {
		return status.Error(codes.Unauthenticated, "invalid current password")
	}

	// Hash the new password
	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}

	// Update the password
	if err := s.repo.UpdateUserPassword(ctx, userID, hashedPassword); err != nil {
		return status.Errorf(codes.Internal, "failed to update password: %v", err)
	}

	return nil
}

// ForgotPassword initiates the password reset process
func (s *AuthService) ForgotPassword(ctx context.Context, email string) (string, error) {
	// Get the user by email
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		// Don't reveal that the email doesn't exist
		return "", nil
	}

	// Generate a password reset token (JWT with short expiration)
	token, expiresAt, err := s.jwtMgr.GenerateToken(
		user.Id,
		user.Email,
		user.Role.String(),
		AccessToken, // Reusing token type for password reset
	)
	if err != nil {
		return "", status.Errorf(codes.Internal, "failed to generate reset token: %v", err)
	}

	// TODO: Send email with reset link containing the token
	resetLink := fmt.Sprintf("https://tribecart.app/reset-password?token=%s", token)

	// In a real application, you would send an email with the reset link
	fmt.Printf("Password reset link for %s: %s\n", email, resetLink)

	return resetLink, nil
}

// ResetPassword resets a user's password using a reset token
func (s *AuthService) ResetPassword(ctx context.Context, token, newPassword string) error {
	// Validate the reset token
	claims, err := s.jwtMgr.ValidateToken(token)
	if err != nil {
		return status.Error(codes.Unauthenticated, "invalid or expired reset token")
	}

	// Hash the new password
	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}

	// Update the password
	if err := s.repo.UpdateUserPassword(ctx, claims.UserID, hashedPassword); err != nil {
		return status.Errorf(codes.Internal, "failed to update password: %v", err)
	}

	return nil
}
