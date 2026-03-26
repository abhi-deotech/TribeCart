package auth

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// ContextKey is a type for context keys
type ContextKey string

const (
	// UserIDKey is the context key for user ID
	UserIDKey ContextKey = "userID"
	// UserRoleKey is the context key for user role
	UserRoleKey ContextKey = "userRole"
)

// AuthInterceptor is a gRPC interceptor for authentication
type AuthInterceptor struct {
	jwtMgr *JWTManager
}

// NewAuthInterceptor creates a new auth interceptor
func NewAuthInterceptor(jwtMgr *JWTManager) *AuthInterceptor {
	return &AuthInterceptor{jwtMgr: jwtMgr}
}

// Unary returns a server interceptor function to authenticate and authorize unary RPC
func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
		// Skip authentication for public endpoints
		if isPublicEndpoint(info.FullMethod) {
			return handler(ctx, req)
		}

		// Extract token from metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
		}

		// Get authorization header
		authHeader, ok := md["authorization"]
		if !ok || len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "authorization token is not provided")
		}

		// Validate the token format (Bearer <token>)
		tokenParts := strings.Fields(authHeader[0])
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			return nil, status.Error(codes.Unauthenticated, "invalid authorization header format")
		}

		// Validate the token
		claims, err := i.jwtMgr.ValidateToken(tokenParts[1])
		if err != nil {
			if err == ErrTokenExpired {
				return nil, status.Error(codes.Unauthenticated, "token has expired")
			}
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		// Check if the token is an access token
		if claims.Type != AccessToken {
			return nil, status.Error(codes.Unauthenticated, "invalid token type")
		}

		// Add user info to context
		ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserRoleKey, claims.Role)

		// Check if the user has the required role
		if !hasRequiredRole(claims.Role, info.FullMethod) {
			return nil, status.Error(codes.PermissionDenied, "insufficient permissions")
		}

		// Call the handler
		return handler(ctx, req)
	}
}

// isPublicEndpoint checks if the endpoint is public
func isPublicEndpoint(method string) bool {
	// Define public endpoints (without authentication)
	publicEndpoints := map[string]bool{
		"/tribecart.v1.UserService/Login":        true,
		"/tribecart.v1.UserService/RefreshToken": true,
		"/tribecart.v1.UserService/ForgotPassword": true,
		"/tribecart.v1.UserService/ResetPassword":  true,
		"/tribecart.v1.UserService/CreateUser":    true, // Typically public for self-registration
	}

	return publicEndpoints[method]
}

// hasRequiredRole checks if the user has the required role for the endpoint
func hasRequiredRole(userRole string, method string) bool {
	// Define role-based access control (RBAC) rules
	rbacRules := map[string][]string{
		// Admin-only endpoints
		"/tribecart.v1.UserService/ListUsers": {"admin", "super_admin"},
		"/tribecart.v1.UserService/UpdateUser": {"admin", "super_admin"},
		"/tribecart.v1.UserService/DeleteUser": {"admin", "super_admin"},

		// Seller endpoints
		"/tribecart.v1.UserService/SomeSellerEndpoint": {"seller", "admin", "super_admin"},

		// Default: all authenticated users can access
		"default": {"user", "seller", "admin", "super_admin"},
	}

	// Check if there's a specific rule for this endpoint
	requiredRoles, exists := rbacRules[method]
	if !exists {
		// Use default rule if no specific rule exists
		requiredRoles = rbacRules["default"]
	}

	// Check if the user's role is in the required roles
	for _, role := range requiredRoles {
		if role == userRole {
			return true
		}
	}

	return false
}

// GetUserIDFromContext gets the user ID from the context
func GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok || userID == "" {
		return "", status.Error(codes.Unauthenticated, "user ID not found in context")
	}
	return userID, nil
}

// GetUserRoleFromContext gets the user role from the context
func GetUserRoleFromContext(ctx context.Context) (string, error) {
	role, ok := ctx.Value(UserRoleKey).(string)
	if !ok || role == "" {
		return "", status.Error(codes.Unauthenticated, "user role not found in context")
	}
	return role, nil
}
