package services

import (
	"errors"
	"time"

	"admin-dashboard/internal/models"
	"admin-dashboard/internal/repository"

	"gorm.io/gorm"
)

// UserService handles user-related operations
type UserService struct {
	userRepository     *repository.UserRepository
	roleRepository     *repository.RoleRepository
	divisionRepository *repository.DivisionRepository
	positionRepository *repository.PositionRepository
}

// NewUserService creates a new user service
func NewUserService(
	userRepository *repository.UserRepository,
	roleRepository *repository.RoleRepository,
	divisionRepository *repository.DivisionRepository,
	positionRepository *repository.PositionRepository,
) *UserService {
	return &UserService{
		userRepository:     userRepository,
		roleRepository:     roleRepository,
		divisionRepository: divisionRepository,
		positionRepository: positionRepository,
	}
}

// Get gets a user by ID
func (s *UserService) Get(id uint) (*models.UserResponse, error) {
	// Get user
	user, err := s.userRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	
	// Get user roles
	roles, err := s.roleRepository.GetUserRoles(user.ID)
	if err != nil {
		return nil, err
	}
	
	// Convert roles to string array
	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.Name
	}
	
	// Format birthdate and join date
	var birthdateStr string
	if user.Birthdate != nil {
		birthdateStr = user.Birthdate.Format("2006-01-02")
	}
	
	// Create user response
	userResponse := &models.UserResponse{
		ID:           user.ID,
		UID:          user.UID,
		EmployeeID:   user.EmployeeID,
		Name:         user.Name,
		Email:        user.Email,
		Phone:        user.Phone,
		Address:      user.Address,
		Birthdate:    birthdateStr,
		JoinDate:     user.JoinDate.Format("2006-01-02"),
		ProfileImage: user.ProfileImage,
		IsManager:    user.IsManager,
		IsActive:     user.IsActive,
		Roles:        roleNames,
	}
	
	// Add related information if available
	if user.Division != nil {
		userResponse.Division = user.Division.Name
	}
	
	if user.Position != nil {
		userResponse.Position = user.Position.Name
	}
	
	if user.Manager != nil {
		userResponse.Manager = user.Manager.Name
	}
	
	return userResponse, nil
}

// Create creates a new user
func (s *UserService) Create(request *models.CreateUserRequest, createdBy string) (*models.UserResponse, error) {
	// Check if employee ID already exists
	_, err := s.userRepository.FindByEmployeeID(request.EmployeeID)
	if err == nil {
		return nil, errors.New("employee ID already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	
	// Check if email already exists
	_, err = s.userRepository.FindByEmail(request.Email)
	if err == nil {
		return nil, errors.New("email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	
	// Parse birthdate
	var birthdate *time.Time
	if request.Birthdate != "" {
		t, err := time.Parse("2006-01-02", request.Birthdate)
		if err != nil {
			return nil, errors.New("invalid birthdate format, use YYYY-MM-DD")
		}
		birthdate = &t
	}
	
	// Parse join date
	joinDate, err := time.Parse("2006-01-02", request.JoinDate)
	if err != nil {
		return nil, errors.New("invalid join date format, use YYYY-MM-DD")
	}
	
	// Create user object
	user := &models.User{
		EmployeeID:   request.EmployeeID,
		Name:         request.Name,
		Email:        request.Email,
		Password:     request.Password,
		Phone:        request.Phone,
		Address:      request.Address,
		Birthdate:    birthdate,
		JoinDate:     joinDate,
		ProfileImage: request.ProfileImage,
		DivisionID:   request.DivisionID,
		PositionID:   request.PositionID,
		IsManager:    request.IsManager,
		ManagerID:    request.ManagerID,
		IsActive:     true, // Default to active
	}
	
	// Create user in database
	err = s.userRepository.Create(user, request.RoleIDs, createdBy)
	if err != nil {
		return nil, err
	}
	
	// Get created user
	return s.Get(user.ID)
}

// Update updates a user
func (s *UserService) Update(id uint, request *models.UpdateUserRequest, updatedBy string) (*models.UserResponse, error) {
	// Get user
	user, err := s.userRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	
	// Update user fields if provided
	if request.Name != "" {
		user.Name = request.Name
	}
	
	if request.Email != "" && request.Email != user.Email {
		// Check if email already exists
		_, err = s.userRepository.FindByEmail(request.Email)
		if err == nil {
			return nil, errors.New("email already exists")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		
		user.Email = request.Email
	}
	
	if request.Phone != "" {
		user.Phone = request.Phone
	}
	
	if request.Address != "" {
		user.Address = request.Address
	}
	
	if request.Birthdate != "" {
		t, err := time.Parse("2006-01-02", request.Birthdate)
		if err != nil {
			return nil, errors.New("invalid birthdate format, use YYYY-MM-DD")
		}
		user.Birthdate = &t
	}
	
	if request.JoinDate != "" {
		joinDate, err := time.Parse("2006-01-02", request.JoinDate)
		if err != nil {
			return nil, errors.New("invalid join date format, use YYYY-MM-DD")
		}
		user.JoinDate = joinDate
	}
	
	if request.ProfileImage != "" {
		user.ProfileImage = request.ProfileImage
	}
	
	if request.DivisionID != nil {
		user.DivisionID = request.DivisionID
	}
	
	if request.PositionID != nil {
		user.PositionID = request.PositionID
	}
	
	if request.IsManager != nil {
		user.IsManager = *request.IsManager
	}
	
	if request.ManagerID != nil {
		user.ManagerID = request.ManagerID
	}
	
	if request.IsActive != nil {
		user.IsActive = *request.IsActive
	}
	
	// Update user in database
	err = s.userRepository.Update(user, request.RoleIDs, updatedBy)
	if err != nil {
		return nil, err
	}
	
	// Get updated user
	return s.Get(user.ID)
}

// Delete deletes a user
func (s *UserService) Delete(id uint) error {
	return s.userRepository.Delete(id)
}

// List lists all users with pagination
func (s *UserService) List(page, pageSize int, search string) (*models.PaginatedResponse, error) {
	// Validate page and pageSize
	if page < 1 {
		page = 1
	}
	
	if pageSize < 1 {
		pageSize = 10
	}
	
	// Get paginated list of users
	paginatedResponse, err := s.userRepository.List(page, pageSize, search)
	if err != nil {
		return nil, err
	}
	
	// Convert users to user responses
	users := paginatedResponse.Items.([]models.User)
	userResponses := make([]models.UserResponse, len(users))
	
	for i, user := range users {
		// Format birthdate and join date
		var birthdateStr string
		if user.Birthdate != nil {
			birthdateStr = user.Birthdate.Format("2006-01-02")
		}
		
		// Get roles
		roleNames := make([]string, len(user.Roles))
		for j, role := range user.Roles {
			roleNames[j] = role.Name
		}
		
		// Create user response
		userResponse := models.UserResponse{
			ID:           user.ID,
			UID:          user.UID,
			EmployeeID:   user.EmployeeID,
			Name:         user.Name,
			Email:        user.Email,
			Phone:        user.Phone,
			Address:      user.Address,
			Birthdate:    birthdateStr,
			JoinDate:     user.JoinDate.Format("2006-01-02"),
			ProfileImage: user.ProfileImage,
			IsManager:    user.IsManager,
			IsActive:     user.IsActive,
			Roles:        roleNames,
		}
		
		// Add related information if available
		if user.Division != nil {
			userResponse.Division = user.Division.Name
		}
		
		if user.Position != nil {
			userResponse.Position = user.Position.Name
		}
		
		if user.Manager != nil {
			userResponse.Manager = user.Manager.Name
		}
		
		userResponses[i] = userResponse
	}
	
	// Update response items
	paginatedResponse.Items = userResponses
	
	return paginatedResponse, nil
}

// UpdatePassword updates a user's password
func (s *UserService) UpdatePassword(id uint, password string, updatedBy string) error {
	return s.userRepository.UpdatePassword(id, password, updatedBy)
}