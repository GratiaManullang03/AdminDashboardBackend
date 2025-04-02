package services

import (
	"errors"

	"admin-dashboard/internal/models"
	"admin-dashboard/internal/repository"

	"gorm.io/gorm"
)

// PositionService handles position-related operations
type PositionService struct {
	positionRepository *repository.PositionRepository
}

// NewPositionService creates a new position service
func NewPositionService(positionRepository *repository.PositionRepository) *PositionService {
	return &PositionService{
		positionRepository: positionRepository,
	}
}

// Get gets a position by ID
func (s *PositionService) Get(id uint) (*models.Position, error) {
	return s.positionRepository.FindByID(id)
}

// Create creates a new position
func (s *PositionService) Create(request *models.PositionRequest, createdBy string) (*models.Position, error) {
	// Check if code already exists
	_, err := s.positionRepository.FindByCode(request.Code)
	if err == nil {
		return nil, errors.New("position code already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	
	// Create position object
	position := &models.Position{
		Code:     request.Code,
		Name:     request.Name,
		IsActive: true, // Default to active
	}
	
	// Create position in database
	err = s.positionRepository.Create(position, createdBy)
	if err != nil {
		return nil, err
	}
	
	// Get created position
	return s.positionRepository.FindByID(position.ID)
}

// Update updates a position
func (s *PositionService) Update(id uint, request *models.PositionRequest, updatedBy string) (*models.Position, error) {
	// Get position
	position, err := s.positionRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	
	// Check if code is changed and already exists
	if request.Code != position.Code {
		existing, err := s.positionRepository.FindByCode(request.Code)
		if err == nil && existing.ID != id {
			return nil, errors.New("position code already exists")
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		
		position.Code = request.Code
	}
	
	// Update position fields
	position.Name = request.Name
	
	// Update position in database
	err = s.positionRepository.Update(position, updatedBy)
	if err != nil {
		return nil, err
	}
	
	// Get updated position
	return s.positionRepository.FindByID(position.ID)
}

// Delete deletes a position
func (s *PositionService) Delete(id uint) error {
	return s.positionRepository.Delete(id)
}

// List lists all positions with pagination
func (s *PositionService) List(page, pageSize int, search string) (*models.PaginatedResponse, error) {
	// Validate page and pageSize
	if page < 1 {
		page = 1
	}
	
	if pageSize < 1 {
		pageSize = 10
	}
	
	// Get paginated list of positions
	return s.positionRepository.List(page, pageSize, search)
}

// ListAll lists all active positions without pagination
func (s *PositionService) ListAll() ([]models.Position, error) {
	return s.positionRepository.ListAll()
}