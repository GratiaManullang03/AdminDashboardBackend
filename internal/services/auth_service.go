package services

import (
	"admin-dashboard/internal/models"
	"admin-dashboard/internal/repository"
	"admin-dashboard/internal/utils"
)

// AuthService handles authentication related operations
type AuthService struct {
	userRepository *repository.UserRepository
	roleRepository *repository.RoleRepository
	jwtManager     *utils.JWTManager
}

// NewAuthService creates a new auth service
func NewAuthService(
	userRepository *repository.UserRepository,
	roleRepository *repository.RoleRepository,
	jwtManager *utils.JWTManager,
) *AuthService {
	return &AuthService{
		userRepository: userRepository,
		roleRepository: roleRepository,
		jwtManager:     jwtManager,
	}
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(email, password string) (*models.LoginResponse, error) {
	// Authenticate user
	user, err := s.userRepository.Authenticate(email, password)
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
	
	// Generate JWT token
	token, err := s.jwtManager.GenerateToken(user.ID, user.UID, user.EmployeeID, user.Email, roleNames)
	if err != nil {
		return nil, err
	}
	
	// Format birthdate and join date
	var birthdateStr string
	if user.Birthdate != nil {
		birthdateStr = user.Birthdate.Format("2006-01-02")
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
	
	// Create login response
	response := &models.LoginResponse{
		Token: token,
		User:  userResponse,
	}
	
	return response, nil
}

// GetUserByID gets a user by ID and returns a UserResponse
func (s *AuthService) GetUserByID(id uint) (*models.UserResponse, error) {
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