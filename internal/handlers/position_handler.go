package handlers

import (
	"net/http"
	"strconv"

	"admin-dashboard/internal/models"
	"admin-dashboard/internal/services"

	"github.com/gin-gonic/gin"
)

// PositionHandler handles position-related HTTP requests
type PositionHandler struct {
	positionService *services.PositionService
}

// NewPositionHandler creates a new position handler
func NewPositionHandler(positionService *services.PositionService) *PositionHandler {
	return &PositionHandler{
		positionService: positionService,
	}
}

// Get gets a position by ID
// @Summary Get a position by ID
// @Description Get a position by its ID
// @Tags positions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Position ID"
// @Success 200 {object} models.Position "Position details"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 404 {object} map[string]string "Position not found"
// @Failure 500 {object} map[string]string "Server error"
// @Router /positions/{id} [get]
func (h *PositionHandler) Get(c *gin.Context) {
	// Parse ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid position ID"})
		return
	}
	
	// Get position
	position, err := h.positionService.Get(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Position not found"})
		return
	}
	
	c.JSON(http.StatusOK, position)
}

// Create creates a new position
// @Summary Create a new position
// @Description Create a new position with the provided details
// @Tags positions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param position body models.PositionRequest true "Position details"
// @Success 201 {object} models.Position "Created position"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Server error"
// @Router /positions [post]
func (h *PositionHandler) Create(c *gin.Context) {
	var request models.PositionRequest
	
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
	
	// Create position
	position, err := h.positionService.Create(&request, employeeID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, position)
}

// Update updates a position
// @Summary Update a position
// @Description Update a position with the provided details
// @Tags positions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Position ID"
// @Param position body models.PositionRequest true "Position details"
// @Success 200 {object} models.Position "Updated position"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 404 {object} map[string]string "Position not found"
// @Failure 500 {object} map[string]string "Server error"
// @Router /positions/{id} [put]
func (h *PositionHandler) Update(c *gin.Context) {
	// Parse ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid position ID"})
		return
	}
	
	var request models.PositionRequest
	
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
	
	// Update position
	position, err := h.positionService.Update(uint(id), &request, employeeID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, position)
}

// Delete deletes a position
// @Summary Delete a position
// @Description Delete a position by its ID
// @Tags positions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Position ID"
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Server error"
// @Router /positions/{id} [delete]
func (h *PositionHandler) Delete(c *gin.Context) {
	// Parse ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid position ID"})
		return
	}
	
	// Delete position
	err = h.positionService.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Position deleted successfully"})
}

// List lists all positions with pagination
// @Summary List all positions
// @Description List all positions with pagination
// @Tags positions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Page size (default: 10)"
// @Param search query string false "Search term"
// @Success 200 {object} models.PaginatedResponse "List of positions"
// @Failure 500 {object} map[string]string "Server error"
// @Router /positions [get]
func (h *PositionHandler) List(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	
	// Get positions
	positions, err := h.positionService.List(page, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, positions)
}

// ListAll lists all active positions without pagination
// @Summary List all active positions
// @Description List all active positions without pagination
// @Tags positions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Position "List of positions"
// @Failure 500 {object} map[string]string "Server error"
// @Router /positions/all [get]
func (h *PositionHandler) ListAll(c *gin.Context) {
	// Get all positions
	positions, err := h.positionService.ListAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, positions)
}

// RegisterRoutes registers the position routes
func (h *PositionHandler) RegisterRoutes(router *gin.RouterGroup, authMiddleware *gin.HandlerFunc) {
	positionGroup := router.Group("/positions")
	positionGroup.Use(*authMiddleware) // Apply auth middleware
	{
		positionGroup.GET("", h.List)
		positionGroup.GET("/all", h.ListAll)
		positionGroup.POST("", h.Create)
		positionGroup.GET("/:id", h.Get)
		positionGroup.PUT("/:id", h.Update)
		positionGroup.DELETE("/:id", h.Delete)
	}
}