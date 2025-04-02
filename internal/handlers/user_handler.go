package handlers

import (
	"net/http"
	"strconv"

	"admin-dashboard/internal/models"
	"admin-dashboard/internal/services"

	"github.com/gin-gonic/gin"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Get gets a user by ID
// @Summary Get a user by ID
// @Description Get a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} models.UserResponse "User details"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Server error"
// @Router /users/{id} [get]
func (h *UserHandler) Get(c *gin.Context) {
	// Parse ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	
	// Get user
	user, err := h.userService.Get(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	
	c.JSON(http.StatusOK, user)
}

// Create creates a new user
// @Summary Create a new user
// @Description Create a new user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body models.CreateUserRequest true "User details"
// @Success 201 {object} models.UserResponse "Created user"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Server error"
// @Router /users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var request models.CreateUserRequest
	
	// Bind JSON to request struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Get creator ID from context
	employeeID, exists := c.Get("employeeID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Create user
	user, err := h.userService.Create(&request, employeeID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, user)
}

// Update updates a user
// @Summary Update a user
// @Description Update a user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param user body models.UpdateUserRequest true "User details"
// @Success 200 {object} models.UserResponse "Updated user"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Server error"
// @Router /users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	// Parse ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	
	var request models.UpdateUserRequest
	
	// Bind JSON to request struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Get updater ID from context
	employeeID, exists := c.Get("employeeID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	
	// Update user
	user, err := h.userService.Update(uint(id), &request, employeeID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, user)
}

// Delete deletes a user
// @Summary Delete a user
// @Description Delete a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Server error"
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	// Parse ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	
	// Delete user
	err = h.userService.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// List lists all users with pagination
// @Summary List all users
// @Description List all users with pagination
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Page size (default: 10)"
// @Param search query string false "Search term"
// @Success 200 {object} models.PaginatedResponse "List of users"
// @Failure 500 {object} map[string]string "Server error"
// @Router /users [get]
func (h *UserHandler) List(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	
	// Get users
	users, err := h.userService.List(page, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, users)
}

// RegisterRoutes registers the user routes
func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup, authMiddleware *gin.HandlerFunc) {
	userGroup := router.Group("/users")
	userGroup.Use(*authMiddleware) // Apply auth middleware
	{
		userGroup.GET("", h.List)
		userGroup.POST("", h.Create)
		userGroup.GET("/:id", h.Get)
		userGroup.PUT("/:id", h.Update)
		userGroup.DELETE("/:id", h.Delete)
	}
}