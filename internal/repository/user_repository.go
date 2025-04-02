package repository

import (
	"errors"
	"time"

	"admin-dashboard/internal/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserRepository handles user-related database operations
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	result := r.db.Preload("Division").
		Preload("Position").
		Preload("Manager").
		Preload("Roles").
		First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.db.Preload("Division").
		Preload("Position").
		Preload("Manager").
		Preload("Roles").
		Where("u_email = ?", email).
		First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// FindByEmployeeID finds a user by employee ID
func (r *UserRepository) FindByEmployeeID(employeeID string) (*models.User, error) {
	var user models.User
	result := r.db.Preload("Division").
		Preload("Position").
		Preload("Manager").
		Preload("Roles").
		Where("u_employee_id = ?", employeeID).
		First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// FindByUID finds a user by UUID
func (r *UserRepository) FindByUID(uid uuid.UUID) (*models.User, error) {
	var user models.User
	result := r.db.Preload("Division").
		Preload("Position").
		Preload("Manager").
		Preload("Roles").
		Where("u_uid = ?", uid).
		First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User, roleIDs []uint, createdBy string) error {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	
	// Set creation info
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	user.CreatedBy = createdBy
	user.UpdatedBy = createdBy

	// Start a transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Create the user
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Assign roles to the user
	for _, roleID := range roleIDs {
		userRole := models.UserRole{
			UserID:    user.ID,
			RoleID:    roleID,
			CreatedAt: now,
			CreatedBy: createdBy,
		}
		if err := tx.Create(&userRole).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	return tx.Commit().Error
}

// Update updates a user
func (r *UserRepository) Update(user *models.User, roleIDs []uint, updatedBy string) error {
	// Set update info
	user.UpdatedAt = time.Now()
	user.UpdatedBy = updatedBy

	// Start a transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Update user
	if err := tx.Model(&models.User{}).Where("u_id = ?", user.ID).Updates(map[string]interface{}{
		"u_name":          user.Name,
		"u_email":         user.Email,
		"u_phone":         user.Phone,
		"u_address":       user.Address,
		"u_birthdate":     user.Birthdate,
		"u_join_date":     user.JoinDate,
		"u_profile_image": user.ProfileImage,
		"u_division_id":   user.DivisionID,
		"u_position_id":   user.PositionID,
		"u_is_manager":    user.IsManager,
		"u_manager_id":    user.ManagerID,
		"u_is_active":     user.IsActive,
		"u_updated_at":    user.UpdatedAt,
		"u_updated_by":    user.UpdatedBy,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// If roleIDs are provided, update user roles
	if roleIDs != nil {
		// Delete existing user roles
		if err := tx.Where("ur_user_id = ?", user.ID).Delete(&models.UserRole{}).Error; err != nil {
			tx.Rollback()
			return err
		}

		// Assign new roles
		for _, roleID := range roleIDs {
			userRole := models.UserRole{
				UserID:    user.ID,
				RoleID:    roleID,
				CreatedAt: time.Now(),
				CreatedBy: updatedBy,
			}
			if err := tx.Create(&userRole).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// Commit the transaction
	return tx.Commit().Error
}

// UpdatePassword updates a user's password
func (r *UserRepository) UpdatePassword(userID uint, password string, updatedBy string) error {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update the password
	return r.db.Model(&models.User{}).Where("u_id = ?", userID).Updates(map[string]interface{}{
		"u_password":   string(hashedPassword),
		"u_updated_at": time.Now(),
		"u_updated_by": updatedBy,
	}).Error
}

// Delete deletes a user
func (r *UserRepository) Delete(id uint) error {
	// Check if the user exists
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return err
	}

	// Start a transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Delete user roles
	if err := tx.Where("ur_user_id = ?", id).Delete(&models.UserRole{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete the user
	if err := tx.Delete(&models.User{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	return tx.Commit().Error
}

// List lists all users with pagination
func (r *UserRepository) List(page, limit int, search string) (*models.PaginatedResponse, error) {
	var users []models.User
	var totalItems int64
	
	// Base query
	query := r.db.Model(&models.User{}).
		Preload("Division").
		Preload("Position").
		Preload("Manager").
		Preload("Roles")
	
	// Apply search if provided
	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("u_name ILIKE ? OR u_email ILIKE ? OR u_employee_id ILIKE ?", searchTerm, searchTerm, searchTerm)
	}
	
	// Get total count
	if err := query.Count(&totalItems).Error; err != nil {
		return nil, err
	}
	
	// Apply pagination
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
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
		Items:       users,
	}
	
	return response, nil
}

// Authenticate authenticates a user with email and password
func (r *UserRepository) Authenticate(email, password string) (*models.User, error) {
	// Find user by email
	user, err := r.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}
	
	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("your account is inactive")
	}
	
	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}
	
	return user, nil
}