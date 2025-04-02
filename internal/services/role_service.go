package services

import (
	"errors"

	"admin-dashboard/internal/models"
	"admin-dashboard/internal/repository"

	"gorm.io/gorm"
)

// RoleService handles role-related operations
type RoleService struct {
	roleRepository *repository.RoleRepository
}

// NewRoleService creates a new role service
func NewRoleService(roleRepository *repository.RoleRepository) *RoleService {
	return &RoleService{
		roleRepository: roleRepository,
	}
}

// Get gets a role by ID
func (s *RoleService) Get(id uint) (*models.Role, error) {
	return s.roleRepository.FindByID(id)
}

// Create creates a new role
func (s *RoleService) Create(request *models.RoleRequest, createdBy string) (*models.Role, error) {
	// Check if name already exists
	_, err := s.roleRepository.FindByName(request.Name)
	if err == nil {
		return nil, errors.New("role name already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	
	// Create role object
	role := &models.Role{
		Name:     request.Name,
		Level:    request.Level,
		IsActive: true, // Default to active
	}
	
	// Create role in database
	err = s.roleRepository.Create(role, createdBy)
	if err != nil {
		return nil, err
	}
	
	// Get created role
	return s.roleRepository.FindByID(role.ID)
}

// Update updates a role
func (s *RoleService) Update(id uint, request *models.RoleRequest, updatedBy string) (*models.Role, error) {
	// Get role
	role, err := s.roleRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	
	// Check if name is changed and already exists
	if request.Name != role.Name {
		existing, err := s.roleRepository.FindByName(request.Name)
		if err == nil && existing.ID != id {
			return nil, errors.New("role name already exists")
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		
		role.Name = request.Name
	}
	
	// Update role fields
	role.Level = request.Level
	
	// Update role in database
	err = s.roleRepository.Update(role, updatedBy)
	if err != nil {
		return nil, err
	}
	
	// Get updated role
	return s.roleRepository.FindByID(role.ID)
}

// Delete deletes a role
func (s *RoleService) Delete(id uint) error {
	return s.roleRepository.Delete(id)
}

// List lists all roles with pagination
func (s *RoleService) List(page, pageSize int, search string) (*models.PaginatedResponse, error) {
	// Validate page and pageSize
	if page < 1 {
		page = 1
	}
	
	if pageSize < 1 {
		pageSize = 10
	}
	
	// Get paginated list of roles
	return s.roleRepository.List(page, pageSize, search)
}

// ListAll lists all active roles without pagination
func (s *RoleService) ListAll() ([]models.Role, error) {
	return s.roleRepository.ListAll()
}

// GetUserRoles gets all roles for a user
func (s *RoleService) GetUserRoles(userID uint) ([]models.Role, error) {
	return s.roleRepository.GetUserRoles(userID)
}