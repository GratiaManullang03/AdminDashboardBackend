package services

import (
	"time"

	"admin-dashboard/internal/models"

	"gorm.io/gorm"
)

// DashboardService handles dashboard-related operations
type DashboardService struct {
	db *gorm.DB
}

// NewDashboardService creates a new dashboard service
func NewDashboardService(db *gorm.DB) *DashboardService {
	return &DashboardService{
		db: db,
	}
}

// GetStatistics gets dashboard statistics
func (s *DashboardService) GetStatistics() (*models.Statistics, error) {
	var stats models.Statistics
	
	// Get total users
	if err := s.db.Model(&models.User{}).Count(&stats.TotalUsers).Error; err != nil {
		return nil, err
	}
	
	// Get active users
	if err := s.db.Model(&models.User{}).Where("u_is_active = ?", true).Count(&stats.ActiveUsers).Error; err != nil {
		return nil, err
	}
	
	// Get total divisions
	if err := s.db.Model(&models.Division{}).Count(&stats.TotalDivisions).Error; err != nil {
		return nil, err
	}
	
	// Get total positions
	if err := s.db.Model(&models.Position{}).Count(&stats.TotalPositions).Error; err != nil {
		return nil, err
	}
	
	// Get users per division
	var usersPerDivision []struct {
		DivisionName string
		UserCount int64
	}

	divisionQuery := `
        SELECT d.div_name as division_name, COUNT(u.u_id) as user_count 
        FROM "user".users u 
        JOIN "user".divisions d ON u.u_division_id = d.div_id 
        WHERE u.u_is_active = true 
        GROUP BY d.div_name
    `
    
    if err := s.db.Raw(divisionQuery).Scan(&usersPerDivision).Error; err != nil {
        return nil, err
    }
	
	// Convert to map for response
	stats.UsersPerDivision = make([]map[string]interface{}, len(usersPerDivision))
	for i, item := range usersPerDivision {
		stats.UsersPerDivision[i] = map[string]interface{}{
			"division": item.DivisionName,
			"count":    item.UserCount,
		}
	}
	
	// Get users per position
	var usersPerPosition []struct {
		PositionName string
		UserCount    int64
	}
	
	positionQuery := `
        SELECT p.pos_name as position_name, COUNT(u.u_id) as user_count 
        FROM "user".users u 
        JOIN "user".positions p ON u.u_position_id = p.pos_id 
        WHERE u.u_is_active = true 
        GROUP BY p.pos_name
    `
    
    if err := s.db.Raw(positionQuery).Scan(&usersPerPosition).Error; err != nil {
        return nil, err
    }
	
	// Convert to map for response
	stats.UsersPerPosition = make([]map[string]interface{}, len(usersPerPosition))
	for i, item := range usersPerPosition {
		stats.UsersPerPosition[i] = map[string]interface{}{
			"position": item.PositionName,
			"count":    item.UserCount,
		}
	}
	
	// Get new users this month
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	
	if err := s.db.Model(&models.User{}).
		Where("u_join_date >= ?", startOfMonth).
		Count(&stats.NewUsersThisMonth).Error; err != nil {
		return nil, err
	}
	
	return &stats, nil
}