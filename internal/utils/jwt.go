package utils

import (
	"errors"
	"fmt"
	"log"
	"time"

	"admin-dashboard/internal/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// CustomClaims represents the claims in the JWT token
type CustomClaims struct {
	UserID     uint      `json:"user_id"`
	UID        uuid.UUID `json:"uid"`
	EmployeeID string    `json:"employee_id"`
	Email      string    `json:"email"`
	Roles      []string  `json:"roles"`
	jwt.RegisteredClaims
}

// JWTManager handles JWT token operations
type JWTManager struct {
	config *config.JWTConfig
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(config *config.JWTConfig) *JWTManager {
	return &JWTManager{
		config: config,
	}
}

// GenerateToken generates a new JWT token
func (m *JWTManager) GenerateToken(userID uint, uid uuid.UUID, employeeID, email string, roles []string) (string, error) {
	// Set expiration time
	expirationTime := time.Now().Add(time.Duration(m.config.Expiry) * time.Hour)

	// Create claims
	claims := &CustomClaims{
		UserID:     userID,
		UID:        uid,
		EmployeeID: employeeID,
		Email:      email,
		Roles:      roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate the token string
	tokenString, err := token.SignedString([]byte(m.config.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates the token and returns the claims
func (m *JWTManager) ValidateToken(tokenString string) (*CustomClaims, error) {
	log.Printf("Validating token: %s", tokenString) // Log token yang akan divalidasi
    log.Printf("Using secret: %s", m.config.Secret) // Log secret yang digunakan (hanya untuk debug!)

	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.config.Secret), nil
	})

	if err != nil {
		log.Printf("Token parse error: %v", err) // Log error parsing
		return nil, err
	}

	// Check if token is valid
	if !token.Valid {
		log.Printf("Token is invalid") // Log jika token tidak valid
		return nil, errors.New("invalid token")
	}

	// Get claims
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	log.Printf("Token validation successful for user: %s", claims.Email) // Log jika validasi berhasil

	return claims, nil
}