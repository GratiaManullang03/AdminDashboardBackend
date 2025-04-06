package repository

import (
	"time"

	"admin-dashboard/internal/models"

	"gorm.io/gorm"
)

// RoleRepository handles role-related database operations
type RoleRepository struct {
	db *gorm.DB
}

// NewRoleRepository creates a new role repository
func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{
		db: db,
	}
}

// FindByID finds a role by ID
func (r *RoleRepository) FindByID(id uint) (*models.Role, error) {
	var role models.Role
	result := r.db.First(&role, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &role, nil
}

// FindByName finds a role by name
func (r *RoleRepository) FindByName(name string) (*models.Role, error) {
	var role models.Role
	result := r.db.Where("role_name = ?", name).First(&role)
	if result.Error != nil {
		return nil, result.Error
	}
	return &role, nil
}

// Create creates a new role
func (r *RoleRepository) Create(role *models.Role, createdBy string) error {
	// Set creation info
	now := time.Now()
	role.CreatedAt = now
	role.UpdatedAt = now
	role.CreatedBy = createdBy
	role.UpdatedBy = createdBy
	
	return r.db.Create(role).Error
}

// Update updates a role
func (r *RoleRepository) Update(role *models.Role, updatedBy string) error {
	// Set update info
	role.UpdatedAt = time.Now()
	role.UpdatedBy = updatedBy
	
	return r.db.Model(&models.Role{}).Where("role_id = ?", role.ID).Updates(map[string]interface{}{
		"role_name":      role.Name,
		"role_level":     role.Level,
		"role_is_active": role.IsActive,
		"role_updated_at": role.UpdatedAt,
		"role_updated_by": role.UpdatedBy,
	}).Error
}

// Delete deletes a role
func (r *RoleRepository) Delete(id uint) error {
	// Check if the role exists
	var role models.Role
	if err := r.db.First(&role, id).Error; err != nil {
		return err
	}
	
	// Check if there are any users with this role
	var count int64
	if err := r.db.Model(&models.UserRole{}).Where("ur_role_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	
	if count > 0 {
		return r.db.Model(&models.Role{}).Where("role_id = ?", id).Update("role_is_active", false).Error
	}
	
	// Delete the role
	return r.db.Delete(&models.Role{}, id).Error
}

// List lists all roles with pagination
func (r *RoleRepository) List(page, limit int, search string) (*models.PaginatedResponse, error) {
	var roles []models.Role
	var totalItems int64
	
	// Base query
	query := r.db.Model(&models.Role{})
	
	// Apply search if provided
	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("role_name ILIKE ?", searchTerm)
	}
	
	// Get total count
	if err := query.Count(&totalItems).Error; err != nil {
		return nil, err
	}
	
	// Apply pagination
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&roles).Error; err != nil {
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
		Items:       roles,
	}
	
	return response, nil
}

// ListAll lists all active roles without pagination
func (r *RoleRepository) ListAll() ([]models.Role, error) {
	var roles []models.Role
	
	// Get all active roles
	err := r.db.Where("role_is_active = ?", true).Order("role_level DESC").Find(&roles).Error
	if err != nil {
		return nil, err
	}
	
	return roles, nil
}

// GetUserRoles gets all roles for a user
func (r *RoleRepository) GetUserRoles(userID uint) ([]models.Role, error) {
    var roles []models.Role
    // Gunakan Raw SQL untuk menghindari masalah dengan schema
    query := `
        SELECT r.* 
        FROM "user".roles r
        JOIN "user".user_roles ur ON r.role_id = ur.ur_role_id
        WHERE ur.ur_user_id = ?
    `
    err := r.db.Raw(query, userID).Scan(&roles).Error
    if err != nil {
        return nil, err
    }
    return roles, nil
}