package handlers

import (
	"net/http"
	"strconv"

	"admin-dashboard/internal/models"
	"admin-dashboard/internal/services"

	"github.com/gin-gonic/gin"
)

// RoleHandler handles role-related HTTP requests
type RoleHandler struct {
	roleService *services.RoleService
}

// NewRoleHandler creates a new role handler
func NewRoleHandler(roleService *services.RoleService) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
	}
}

// Get gets a role by ID
// @Summary Get a role by ID
// @Description Get a role by its ID
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Success 200 {object} models.Role "Role details"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 404 {object} map[string]string "Role not found"
// @Failure 500 {object} map[string]string "Server error"
// @Router /roles/{id} [get]
func (h *RoleHandler) Get(c *gin.Context) {
	// Parse ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}
	
	// Get role
	role, err := h.roleService.Get(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}
	
	c.JSON(http.StatusOK, role)
}

// Create creates a new role
// @Summary Create a new role
// @Description Create a new role with the provided details
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param role body models.RoleRequest true "Role details"
// @Success 201 {object} models.Role "Created role"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Server error"
// @Router /roles [post]
func (h *RoleHandler) Create(c *gin.Context) {
	var request models.RoleRequest
	
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
	
	// Create role
	role, err := h.roleService.Create(&request, employeeID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, role)
}

// Update updates a role
// @Summary Update a role
// @Description Update a role with the provided details
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Param role body models.RoleRequest true "Role details"
// @Success 200 {object} models.Role "Updated role"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 404 {object} map[string]string "Role not found"
// @Failure 500 {object} map[string]string "Server error"
// @Router /roles/{id} [put]
func (h *RoleHandler) Update(c *gin.Context) {
	// Parse ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}
	
	var request models.RoleRequest
	
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
	
	// Update role
	role, err := h.roleService.Update(uint(id), &request, employeeID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, role)
}

// Delete deletes a role
// @Summary Delete a role
// @Description Delete a role by its ID
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Server error"
// @Router /roles/{id} [delete]
func (h *RoleHandler) Delete(c *gin.Context) {
	// Parse ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}
	
	// Delete role
	err = h.roleService.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
}

// List lists all roles with pagination
// @Summary List all roles
// @Description List all roles with pagination
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Page size (default: 10)"
// @Param search query string false "Search term"
// @Success 200 {object} models.PaginatedResponse "List of roles"
// @Failure 500 {object} map[string]string "Server error"
// @Router /roles [get]
func (h *RoleHandler) List(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	
	// Get roles
	roles, err := h.roleService.List(page, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, roles)
}

// ListAll lists all active roles without pagination
// @Summary List all active roles
// @Description List all active roles without pagination
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Role "List of roles"
// @Failure 500 {object} map[string]string "Server error"
// @Router /roles/all [get]
func (h *RoleHandler) ListAll(c *gin.Context) {
	// Get all roles
	roles, err := h.roleService.ListAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, roles)
}

// RegisterRoutes registers the role routes
func (h *RoleHandler) RegisterRoutes(router *gin.RouterGroup, authMiddleware *gin.HandlerFunc) {
	roleGroup := router.Group("/roles")
	roleGroup.Use(*authMiddleware) // Apply auth middleware
	{
		roleGroup.GET("", h.List)
		roleGroup.GET("/all", h.ListAll)
		roleGroup.POST("", h.Create)
		roleGroup.GET("/:id", h.Get)
		roleGroup.PUT("/:id", h.Update)
		roleGroup.DELETE("/:id", h.Delete)
	}
}