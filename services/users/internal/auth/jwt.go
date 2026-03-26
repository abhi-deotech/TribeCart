package auth

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// ErrInvalidToken is returned when the token is invalid
	ErrInvalidToken = errors.New("invalid token")
	// ErrTokenExpired is returned when the token has expired
	ErrTokenExpired = errors.New("token has expired")
)

// TokenType represents the type of JWT token
type TokenType string

const (
	// AccessToken is used for regular API access
	AccessToken TokenType = "access"
	// RefreshToken is used to obtain a new access token
	RefreshToken TokenType = "refresh"
)

// CustomClaims contains the JWT claims
// See: https://www.iana.org/assignments/jwt/jwt.xhtml
// for registered claim names
type CustomClaims struct {
	UserID string    `json:"uid"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	Type   TokenType `json:"type"`
	jwt.RegisteredClaims
}

// JWTManager handles JWT token creation and validation
type JWTManager struct {
	privateKey     *rsa.PrivateKey
	publicKey      *rsa.PublicKey
	accessExpires  time.Duration
	refreshExpires time.Duration
	issuer         string
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey, accessExpires, refreshExpires time.Duration, issuer string) *JWTManager {
	return &JWTManager{
		privateKey:     privateKey,
		publicKey:      publicKey,
		accessExpires:  accessExpires,
		refreshExpires: refreshExpires,
		issuer:         issuer,
	}
}

// GenerateToken generates a new JWT token
func (m *JWTManager) GenerateToken(userID, email, role string, tokenType TokenType) (string, time.Time, error) {
	expiration := time.Now().Add(m.accessExpires)
	if tokenType == RefreshToken {
		expiration = time.Now().Add(m.refreshExpires)
	}

	claims := &CustomClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    m.issuer,
			Subject:   string(tokenType),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(m.privateKey)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, expiration, nil
}

// ValidateToken validates the JWT token and returns the claims
func (m *JWTManager) ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return m.publicKey, nil
		},
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// GenerateTokenPair generates both access and refresh tokens
func (m *JWTManager) GenerateTokenPair(userID, email, role string) (accessToken, refreshToken string, accessExp, refreshExp time.Time, err error) {
	accessToken, accessExp, err = m.GenerateToken(userID, email, role, AccessToken)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, refreshExp, err = m.GenerateToken(userID, email, role, RefreshToken)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, refreshToken, accessExp, refreshExp, nil
}
