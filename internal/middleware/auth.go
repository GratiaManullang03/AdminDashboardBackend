package middleware

import (
	"log"
	"net/http"
	"strings"

	"admin-dashboard/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware represents the authentication middleware
type AuthMiddleware struct {
	jwtManager *utils.JWTManager
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(jwtManager *utils.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
	}
}

// Authenticate authenticates the user using JWT
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		log.Printf("Auth Header received: %s", authHeader) // Log header yang diterima

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Check if the Authorization header has the Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		// Extract the token
		tokenString := parts[1]
		log.Printf("Token extracted: %s", tokenString) // Log token yang diekstrak

		// Validate the token
		claims, err := m.jwtManager.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set the user in the context
		c.Set("userID", claims.UserID)
		c.Set("uid", claims.UID)
		c.Set("employeeID", claims.EmployeeID)
		c.Set("email", claims.Email)
		c.Set("roles", claims.Roles)

		c.Next()
	}
}

// RequireRole requires the user to have at least one of the specified roles
func (m *AuthMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user roles from context
		userRoles, exists := c.Get("roles")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User roles not found in context"})
			c.Abort()
			return
		}

		// Check if the user has at least one of the required roles
		for _, role := range roles {
			for _, userRole := range userRoles.([]string) {
				if role == userRole {
					c.Next()
					return
				}
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
		c.Abort()
	}
}