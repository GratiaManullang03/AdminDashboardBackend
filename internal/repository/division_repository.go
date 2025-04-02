package repository

import (
	"time"

	"admin-dashboard/internal/models"

	"gorm.io/gorm"
)

// DivisionRepository handles division-related database operations
type DivisionRepository struct {
	db *gorm.DB
}

// NewDivisionRepository creates a new division repository
func NewDivisionRepository(db *gorm.DB) *DivisionRepository {
	return &DivisionRepository{
		db: db,
	}
}

// FindByID finds a division by ID
func (r *DivisionRepository) FindByID(id uint) (*models.Division, error) {
	var division models.Division
	result := r.db.First(&division, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &division, nil
}

// FindByCode finds a division by code
func (r *DivisionRepository) FindByCode(code string) (*models.Division, error) {
	var division models.Division
	result := r.db.Where("div_code = ?", code).First(&division)
	if result.Error != nil {
		return nil, result.Error
	}
	return &division, nil
}

// Create creates a new division
func (r *DivisionRepository) Create(division *models.Division, createdBy string) error {
	// Set creation info
	now := time.Now()
	division.CreatedAt = now
	division.UpdatedAt = now
	division.CreatedBy = createdBy
	division.UpdatedBy = createdBy
	
	return r.db.Create(division).Error
}

// Update updates a division
func (r *DivisionRepository) Update(division *models.Division, updatedBy string) error {
	// Set update info
	division.UpdatedAt = time.Now()
	division.UpdatedBy = updatedBy
	
	return r.db.Model(&models.Division{}).Where("div_id = ?", division.ID).Updates(map[string]interface{}{
		"div_code":       division.Code,
		"div_name":       division.Name,
		"div_is_active":  division.IsActive,
		"div_updated_at": division.UpdatedAt,
		"div_updated_by": division.UpdatedBy,
	}).Error
}

// Delete deletes a division
func (r *DivisionRepository) Delete(id uint) error {
	// Check if the division exists
	var division models.Division
	if err := r.db.First(&division, id).Error; err != nil {
		return err
	}
	
	// Check if there are any users in this division
	var count int64
	if err := r.db.Model(&models.User{}).Where("u_division_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	
	if count > 0 {
		return r.db.Model(&models.Division{}).Where("div_id = ?", id).Update("div_is_active", false).Error
	}
	
	// Delete the division
	return r.db.Delete(&models.Division{}, id).Error
}

// List lists all divisions with pagination
func (r *DivisionRepository) List(page, limit int, search string) (*models.PaginatedResponse, error) {
	var divisions []models.Division
	var totalItems int64
	
	// Base query
	query := r.db.Model(&models.Division{})
	
	// Apply search if provided
	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("div_name ILIKE ? OR div_code ILIKE ?", searchTerm, searchTerm)
	}
	
	// Get total count
	if err := query.Count(&totalItems).Error; err != nil {
		return nil, err
	}
	
	// Apply pagination
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&divisions).Error; err != nil {
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
		Items:       divisions,
	}
	
	return response, nil
}

// ListAll lists all active divisions without pagination
func (r *DivisionRepository) ListAll() ([]models.Division, error) {
	var divisions []models.Division
	
	// Get all active divisions
	err := r.db.Where("div_is_active = ?", true).Find(&divisions).Error
	if err != nil {
		return nil, err
	}
	
	return divisions, nil
}