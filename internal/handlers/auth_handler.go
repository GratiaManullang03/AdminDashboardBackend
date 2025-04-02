package handlers

import (
	"net/http"

	"admin-dashboard/internal/models"
	"admin-dashboard/internal/services"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login handles user login
// @Summary Login a user
// @Description Authenticate a user and return a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.UserLoginRequest true "User credentials"
// @Success 200 {object} models.LoginResponse "Login successful"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Authentication failed"
// @Failure 500 {object} map[string]string "Server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var request models.UserLoginRequest
	
	// Bind JSON to request struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Authenticate user
	response, err := h.authService.Login(request.Email, request.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, response)
}

// Profile gets the current user's profile
// @Summary Get current user profile
// @Description Get the profile of the currently authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.UserResponse "User profile"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Server error"
// @Router /auth/profile [get]
func (h *AuthHandler) Profile(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Get user data from service
	user, err := h.authService.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user profile"})
		return
	}
	
	c.JSON(http.StatusOK, user)
}

// RegisterRoutes registers the auth routes
func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup, authMiddleware *gin.HandlerFunc) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", h.Login) // Login tanpa autentikasi
		
		// Gunakan middleware untuk endpoint profile
		if authMiddleware != nil {
			authGroup.GET("/profile", *authMiddleware, h.Profile)
		} else {
			authGroup.GET("/profile", h.Profile)
		}
	}
}