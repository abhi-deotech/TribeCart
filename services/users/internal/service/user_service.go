package service

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/tribecart/users/internal/auth"
	"github.com/tribecart/users/internal/repository"
	pb "github.com/tribecart/proto/tribecart/v1"
)

// UserService implements the UserServiceServer interface
type UserService struct {
	repo     repository.UserRepository
	authSvc  *auth.AuthService
	jwtMgr   *auth.JWTManager
	jwtSecret string
}

// NewUserService creates a new UserService
func NewUserService(
	repo repository.UserRepository,
	privateKey *rsa.PrivateKey,
	publicKey *rsa.PublicKey,
	jwtSecret string,
) *UserService {
	authSvc := auth.NewAuthService(repo, privateKey, publicKey, jwtSecret)
	jwtMgr := auth.NewJWTManager(
		privateKey,
		publicKey,
		time.Hour*24,      // Access token expires in 24 hours
		time.Hour*24*7,    // Refresh token expires in 7 days
		"tribecart-users", // Issuer
	)

	return &UserService{
		repo:     repo,
		authSvc:  authSvc,
		jwtMgr:   jwtMgr,
		jwtSecret: jwtSecret,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(
	ctx context.Context,
	req *pb.CreateUserRequest,
) (*pb.User, error) {
	// Validate request
	if err := validateCreateUserRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Check if user with this email already exists
	existingUser, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, status.Error(codes.AlreadyExists, "email already exists")
	}

	// Create the user
	now := timestamppb.Now()
	user := &pb.User{
		Id:           uuid.New().String(),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        strings.ToLower(req.Email),
		PhoneNumber:  req.PhoneNumber,
		Role:         pb.UserRole_USER_ROLE_CUSTOMER, // Default role
		Status:       pb.UserStatus_USER_STATUS_ACTIVE,
		EmailVerified: false,
		PhoneVerified: false,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}
	user.Password = hashedPassword

	// Save the user to the database
	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	// Generate verification token (in a real app, this would send an email)
	_, err = s.authSvc.SendVerificationEmail(ctx, &emptypb.Empty{})
	if err != nil {
		// Log the error but don't fail the request
		fmt.Printf("failed to send verification email: %v\n", err)
	}

	// Don't return the hashed password
	user.Password = ""

	return user, nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(
	ctx context.Context,
	req *pb.GetUserRequest,
) (*pb.User, error) {
	// Get the user ID from the request or from the JWT token
	userID := req.Id
	if userID == "" {
		// Try to get user ID from context (for current user)
		var err error
		userID, err = auth.GetUserIDFromContext(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "user ID is required")
		}
	}

	// Get the user from the database
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	// Check permissions
	if !s.hasPermission(ctx, user.Id) {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	// Don't return the hashed password
	user.Password = ""

	return user, nil
}

// UpdateUser updates a user's profile
func (s *UserService) UpdateUser(
	ctx context.Context,
	req *pb.UpdateUserRequest,
) (*pb.User, error) {
	// Get the user ID from the request or from the JWT token
	userID := req.Id
	if userID == "" {
		// Try to get user ID from context (for current user)
		var err error
		userID, err = auth.GetUserIDFromContext(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "user ID is required")
		}
	}

	// Get the existing user
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	// Check permissions
	if !s.hasPermission(ctx, user.Id) {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	// Update the user fields
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.PhoneNumber != "" {
		user.PhoneNumber = req.PhoneNumber
	}

	// Update the user in the database
	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	// Don't return the hashed password
	user.Password = ""

	return user, nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(
	ctx context.Context,
	req *pb.DeleteUserRequest,
) (*emptypb.Empty, error) {
	// Get the user ID from the request or from the JWT token
	userID := req.Id
	if userID == "" {
		// Try to get user ID from context (for current user)
		var err error
		userID, err = auth.GetUserIDFromContext(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "user ID is required")
		}
	}

	// Check if the user exists
	_, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	// Check permissions
	if !s.hasPermission(ctx, userID) {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	// Delete the user
	if err := s.repo.DeleteUser(ctx, userID); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// ListUsers retrieves a paginated list of users
func (s *UserService) ListUsers(
	ctx context.Context,
	req *pb.ListUsersRequest,
) (*pb.ListUsersResponse, error) {
	// Only admins can list users
	if !s.isAdmin(ctx) {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	// Set default pagination values
	page := req.Page
	if page < 1 {
		page = 1
	}

	pageSize := req.PageSize
	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 20
	}

	// Get the users from the database
	users, totalCount, err := s.repo.ListUsers(ctx, page, pageSize, req.Filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list users: %v", err)
	}

	// Don't return hashed passwords
	for _, user := range users {
		user.Password = ""
	}

	return &pb.ListUsersResponse{
		Users:      users,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

// Login authenticates a user and returns JWT tokens
func (s *UserService) Login(
	ctx context.Context,
	req *pb.LoginRequest,
) (*pb.LoginResponse, error) {
	// Validate request
	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}

	// Authenticate the user
	resp, err := s.authSvc.Login(ctx, req.Email, req.Password)
	if err != nil {
		if status.Code(err) == codes.Unauthenticated {
			return nil, status.Error(codes.Unauthenticated, "invalid email or password")
		}
		return nil, err
	}

	return resp, nil
}

// RefreshToken generates a new access token using a refresh token
func (s *UserService) RefreshToken(
	ctx context.Context,
	req *pb.RefreshTokenRequest,
) (*pb.RefreshTokenResponse, error) {
	// Validate request
	if req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh token is required")
	}

	// Refresh the token
	resp, err := s.authSvc.RefreshToken(ctx, &pb.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Logout revokes a refresh token
func (s *UserService) Logout(
	ctx context.Context,
	req *emptypb.Empty,
) (*emptypb.Empty, error) {
	// In a real application, you would add the token to a blacklist
	// For now, we'll just return success
	return &emptypb.Empty{}, nil
}

// ChangePassword changes a user's password
func (s *UserService) ChangePassword(
	ctx context.Context,
	req *pb.ChangePasswordRequest,
) (*emptypb.Empty, error) {
	// Get the user ID from the context
	userID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	// Validate request
	if req.CurrentPassword == "" || req.NewPassword == "" {
		return nil, status.Error(codes.InvalidArgument, "current and new password are required")
	}

	// Change the password
	err = s.authSvc.ChangePassword(ctx, userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// ForgotPassword initiates the password reset process
func (s *UserService) ForgotPassword(
	ctx context.Context,
	req *pb.ForgotPasswordRequest,
) (*emptypb.Empty, error) {
	// Validate request
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	// Send password reset email
	_, err := s.authSvc.ForgotPassword(ctx, req.Email)
	if err != nil {
		// Don't reveal if the email exists or not
		return &emptypb.Empty{}, nil
	}

	return &emptypb.Empty{}, nil
}

// ResetPassword resets a user's password using a reset token
func (s *UserService) ResetPassword(
	ctx context.Context,
	req *pb.ResetPasswordRequest,
) (*emptypb.Empty, error) {
	// Validate request
	if req.Token == "" || req.NewPassword == "" {
		return nil, status.Error(codes.InvalidArgument, "token and new password are required")
	}

	// Reset the password
	err := s.authSvc.ResetPassword(ctx, req.Token, req.NewPassword)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

// SendVerificationEmail sends a verification email to the user
func (s *UserService) SendVerificationEmail(
	ctx context.Context,
	req *emptypb.Empty,
) (*emptypb.Empty, error) {
	// Get the user ID from the context
	userID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	// Get the user
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	// Check if already verified
	if user.EmailVerified {
		return nil, status.Error(codes.AlreadyExists, "email already verified")
	}

	// Send verification email
	_, err = s.authSvc.SendVerificationEmail(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to send verification email: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// VerifyEmail verifies a user's email using a verification token
func (s *UserService) VerifyEmail(
	ctx context.Context,
	req *pb.VerifyEmailRequest,
) (*emptypb.Empty, error) {
	// Validate request
	if req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	// In a real application, you would validate the token and mark the email as verified
	// For now, we'll just return success
	return &emptypb.Empty{}, nil
}

// hasPermission checks if the current user has permission to access the resource
func (s *UserService) hasPermission(ctx context.Context, resourceUserID string) bool {
	// Get the current user ID from the context
	currentUserID, err := auth.GetUserIDFromContext(ctx)
	if err != nil {
		return false
	}

	// Get the current user's role
	role, err := auth.GetUserRoleFromContext(ctx)
	if err != nil {
		return false
	}

	// Admins can access any resource
	if role == "admin" || role == "super_admin" {
		return true
	}

	// Users can only access their own resources
	return currentUserID == resourceUserID
}

// isAdmin checks if the current user is an admin
func (s *UserService) isAdmin(ctx context.Context) bool {
	role, err := auth.GetUserRoleFromContext(ctx)
	if err != nil {
		return false
	}

	return role == "admin" || role == "super_admin"
}

// validateCreateUserRequest validates a CreateUserRequest
func validateCreateUserRequest(req *pb.CreateUserRequest) error {
	if req.Email == "" {
		return errors.New("email is required")
	}

	if req.Password == "" {
		return errors.New("password is required")
	}

	if len(req.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	if req.FirstName == "" {
		return errors.New("first name is required")
	}

	if req.LastName == "" {
		return errors.New("last name is required")
	}

	return nil
}

// validateUpdateUserRequest validates an UpdateUserRequest
func validateUpdateUserRequest(req *pb.UpdateUserRequest) error {
	if req.Id == "" {
		return errors.New("user ID is required")
	}

	if req.FirstName == "" && req.LastName == "" && req.PhoneNumber == "" {
		return errors.New("at least one field must be provided")
	}

	return nil
}
