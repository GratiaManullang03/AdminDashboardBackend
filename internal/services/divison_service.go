package services

import (
	"errors"

	"admin-dashboard/internal/models"
	"admin-dashboard/internal/repository"

	"gorm.io/gorm"
)

// DivisionService handles division-related operations
type DivisionService struct {
	divisionRepository *repository.DivisionRepository
}

// NewDivisionService creates a new division service
func NewDivisionService(divisionRepository *repository.DivisionRepository) *DivisionService {
	return &DivisionService{
		divisionRepository: divisionRepository,
	}
}

// Get gets a division by ID
func (s *DivisionService) Get(id uint) (*models.Division, error) {
	return s.divisionRepository.FindByID(id)
}

// Create creates a new division
func (s *DivisionService) Create(request *models.DivisionRequest, createdBy string) (*models.Division, error) {
	// Check if code already exists
	_, err := s.divisionRepository.FindByCode(request.Code)
	if err == nil {
		return nil, errors.New("division code already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	
	// Create division object
	division := &models.Division{
		Code:     request.Code,
		Name:     request.Name,
		IsActive: true, // Default to active
	}
	
	// Create division in database
	err = s.divisionRepository.Create(division, createdBy)
	if err != nil {
		return nil, err
	}
	
	// Get created division
	return s.divisionRepository.FindByID(division.ID)
}

// Update updates a division
func (s *DivisionService) Update(id uint, request *models.DivisionRequest, updatedBy string) (*models.Division, error) {
	// Get division
	division, err := s.divisionRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	
	// Check if code is changed and already exists
	if request.Code != division.Code {
		existing, err := s.divisionRepository.FindByCode(request.Code)
		if err == nil && existing.ID != id {
			return nil, errors.New("division code already exists")
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		
		division.Code = request.Code
	}
	
	// Update division fields
	division.Name = request.Name
	
	// Update division in database
	err = s.divisionRepository.Update(division, updatedBy)
	if err != nil {
		return nil, err
	}
	
	// Get updated division
	return s.divisionRepository.FindByID(division.ID)
}

// Delete deletes a division
func (s *DivisionService) Delete(id uint) error {
	return s.divisionRepository.Delete(id)
}

// List lists all divisions with pagination
func (s *DivisionService) List(page, pageSize int, search string) (*models.PaginatedResponse, error) {
	// Validate page and pageSize
	if page < 1 {
		page = 1
	}
	
	if pageSize < 1 {
		pageSize = 10
	}
	
	// Get paginated list of divisions
	return s.divisionRepository.List(page, pageSize, search)
}

// ListAll lists all active divisions without pagination
func (s *DivisionService) ListAll() ([]models.Division, error) {
	return s.divisionRepository.ListAll()
}