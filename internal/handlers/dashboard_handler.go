package handlers

import (
	"net/http"

	"admin-dashboard/internal/services"

	"github.com/gin-gonic/gin"
)

// DashboardHandler handles dashboard-related HTTP requests
type DashboardHandler struct {
	dashboardService *services.DashboardService
}

// NewDashboardHandler creates a new dashboard handler
func NewDashboardHandler(dashboardService *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

// GetStatistics gets dashboard statistics
// @Summary Get dashboard statistics
// @Description Get statistics for the dashboard
// @Tags dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.Statistics "Dashboard statistics"
// @Failure 500 {object} map[string]string "Server error"
// @Router /dashboard/statistics [get]
func (h *DashboardHandler) GetStatistics(c *gin.Context) {
	// Get statistics
	stats, err := h.dashboardService.GetStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, stats)
}

// RegisterRoutes registers the dashboard routes
func (h *DashboardHandler) RegisterRoutes(router *gin.RouterGroup, authMiddleware *gin.HandlerFunc) {
	dashboardGroup := router.Group("/dashboard")
	dashboardGroup.Use(*authMiddleware) // Apply auth middleware
	{
		dashboardGroup.GET("/statistics", h.GetStatistics)
	}
}