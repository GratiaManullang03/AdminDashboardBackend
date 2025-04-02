package handlers

import (
	"net/http"
	"strconv"

	"admin-dashboard/internal/models"
	"admin-dashboard/internal/services"

	"github.com/gin-gonic/gin"
)

// DivisionHandler handles division-related HTTP requests
type DivisionHandler struct {
	divisionService *services.DivisionService
}

// NewDivisionHandler creates a new division handler
func NewDivisionHandler(divisionService *services.DivisionService) *DivisionHandler {
	return &DivisionHandler{
		divisionService: divisionService,
	}
}

// Get gets a division by ID
// @Summary Get a division by ID
// @Description Get a division by its ID
// @Tags divisions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Division ID"
// @Success 200 {object} models.Division "Division details"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 404 {object} map[string]string "Division not found"
// @Failure 500 {object} map[string]string "Server error"
// @Router /divisions/{id} [get]
func (h *DivisionHandler) Get(c *gin.Context) {
	// Parse ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid division ID"})
		return
	}
	
	// Get division
	division, err := h.divisionService.Get(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Division not found"})
		return
	}
	
	c.JSON(http.StatusOK, division)
}

// Create creates a new division
// @Summary Create a new division
// @Description Create a new division with the provided details
// @Tags divisions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param division body models.DivisionRequest true "Division details"
// @Success 201 {object} models.Division "Created division"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Server error"
// @Router /divisions [post]
func (h *DivisionHandler) Create(c *gin.Context) {
	var request models.DivisionRequest
	
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
	
	// Create division
	division, err := h.divisionService.Create(&request, employeeID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, division)
}

// Update updates a division
// @Summary Update a division
// @Description Update a division with the provided details
// @Tags divisions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Division ID"
// @Param division body models.DivisionRequest true "Division details"
// @Success 200 {object} models.Division "Updated division"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 404 {object} map[string]string "Division not found"
// @Failure 500 {object} map[string]string "Server error"
// @Router /divisions/{id} [put]
func (h *DivisionHandler) Update(c *gin.Context) {
	// Parse ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid division ID"})
		return
	}
	
	var request models.DivisionRequest
	
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
	
	// Update division
	division, err := h.divisionService.Update(uint(id), &request, employeeID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, division)
}

// Delete deletes a division
// @Summary Delete a division
// @Description Delete a division by its ID
// @Tags divisions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Division ID"
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Server error"
// @Router /divisions/{id} [delete]
func (h *DivisionHandler) Delete(c *gin.Context) {
	// Parse ID from URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid division ID"})
		return
	}
	
	// Delete division
	err = h.divisionService.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Division deleted successfully"})
}

// List lists all divisions with pagination
// @Summary List all divisions
// @Description List all divisions with pagination
// @Tags divisions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Page size (default: 10)"
// @Param search query string false "Search term"
// @Success 200 {object} models.PaginatedResponse "List of divisions"
// @Failure 500 {object} map[string]string "Server error"
// @Router /divisions [get]
func (h *DivisionHandler) List(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	
	// Get divisions
	divisions, err := h.divisionService.List(page, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, divisions)
}

// ListAll lists all active divisions without pagination
// @Summary List all active divisions
// @Description List all active divisions without pagination
// @Tags divisions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Division "List of divisions"
// @Failure 500 {object} map[string]string "Server error"
// @Router /divisions/all [get]
func (h *DivisionHandler) ListAll(c *gin.Context) {
	// Get all divisions
	divisions, err := h.divisionService.ListAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, divisions)
}

// RegisterRoutes registers the division routes
func (h *DivisionHandler) RegisterRoutes(router *gin.RouterGroup, authMiddleware *gin.HandlerFunc) {
	divisionGroup := router.Group("/divisions")
	divisionGroup.Use(*authMiddleware) // Apply auth middleware
	{
		divisionGroup.GET("", h.List)
		divisionGroup.GET("/all", h.ListAll)
		divisionGroup.POST("", h.Create)
		divisionGroup.GET("/:id", h.Get)
		divisionGroup.PUT("/:id", h.Update)
		divisionGroup.DELETE("/:id", h.Delete)
	}
}