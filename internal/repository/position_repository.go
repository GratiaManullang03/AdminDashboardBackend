package repository

import (
	"time"

	"admin-dashboard/internal/models"

	"gorm.io/gorm"
)

// PositionRepository handles position-related database operations
type PositionRepository struct {
	db *gorm.DB
}

// NewPositionRepository creates a new position repository
func NewPositionRepository(db *gorm.DB) *PositionRepository {
	return &PositionRepository{
		db: db,
	}
}

// FindByID finds a position by ID
func (r *PositionRepository) FindByID(id uint) (*models.Position, error) {
	var position models.Position
	result := r.db.First(&position, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &position, nil
}

// FindByCode finds a position by code
func (r *PositionRepository) FindByCode(code string) (*models.Position, error) {
	var position models.Position
	result := r.db.Where("pos_code = ?", code).First(&position)
	if result.Error != nil {
		return nil, result.Error
	}
	return &position, nil
}

// Create creates a new position
func (r *PositionRepository) Create(position *models.Position, createdBy string) error {
	// Set creation info
	now := time.Now()
	position.CreatedAt = now
	position.UpdatedAt = now
	position.CreatedBy = createdBy
	position.UpdatedBy = createdBy
	
	return r.db.Create(position).Error
}

// Update updates a position
func (r *PositionRepository) Update(position *models.Position, updatedBy string) error {
	// Set update info
	position.UpdatedAt = time.Now()
	position.UpdatedBy = updatedBy
	
	return r.db.Model(&models.Position{}).Where("pos_id = ?", position.ID).Updates(map[string]interface{}{
		"pos_code":       position.Code,
		"pos_name":       position.Name,
		"pos_is_active":  position.IsActive,
		"pos_updated_at": position.UpdatedAt,
		"pos_updated_by": position.UpdatedBy,
	}).Error
}

// Delete deletes a position
func (r *PositionRepository) Delete(id uint) error {
	// Check if the position exists
	var position models.Position
	if err := r.db.First(&position, id).Error; err != nil {
		return err
	}
	
	// Check if there are any users with this position
	var count int64
	if err := r.db.Model(&models.User{}).Where("u_position_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	
	if count > 0 {
		return r.db.Model(&models.Position{}).Where("pos_id = ?", id).Update("pos_is_active", false).Error
	}
	
	// Delete the position
	return r.db.Delete(&models.Position{}, id).Error
}

// List lists all positions with pagination
func (r *PositionRepository) List(page, limit int, search string) (*models.PaginatedResponse, error) {
	var positions []models.Position
	var totalItems int64
	
	// Base query
	query := r.db.Model(&models.Position{})
	
	// Apply search if provided
	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("pos_name ILIKE ? OR pos_code ILIKE ?", searchTerm, searchTerm)
	}
	
	// Get total count
	if err := query.Count(&totalItems).Error; err != nil {
		return nil, err
	}
	
	// Apply pagination
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&positions).Error; err != nil {
		return nil, err
	}
	
	// Calculate total pages
	totalPages := (totalItems + int64(limit) - 1) / int64(limit)
	
	// Create response
	response := &models.PaginatedResponse{
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		CurrentPage: int64(page),
		PageSize:    int64(limit),
		Items:       positions,
	}
	
	return response, nil
}

// ListAll lists all active positions without pagination
func (r *PositionRepository) ListAll() ([]models.Position, error) {
	var positions []models.Position
	
	// Get all active positions
	err := r.db.Where("pos_is_active = ?", true).Find(&positions).Error
	if err != nil {
		return nil, err
	}
	
	return positions, nil
}