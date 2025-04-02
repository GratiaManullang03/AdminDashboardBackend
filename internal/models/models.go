package models

import (
	"time"

	"github.com/google/uuid"
)

// Division represents the division table
type Division struct {
	ID        uint      `gorm:"primaryKey;column:div_id" json:"id"`
	Code      string    `gorm:"unique;column:div_code" json:"code"`
	Name      string    `gorm:"column:div_name" json:"name"`
	IsActive  bool      `gorm:"default:true;column:div_is_active" json:"is_active"`
	CreatedAt time.Time `gorm:"column:div_created_at" json:"created_at"`
	CreatedBy string    `gorm:"column:div_created_by" json:"created_by"`
	UpdatedAt time.Time `gorm:"column:div_updated_at" json:"updated_at"`
	UpdatedBy string    `gorm:"column:div_updated_by" json:"updated_by"`
}

// TableName overrides the table name
func (Division) TableName() string {
	return "user.divisions"
}

// Position represents the positions table
type Position struct {
	ID        uint      `gorm:"primaryKey;column:pos_id" json:"id"`
	Code      string    `gorm:"unique;column:pos_code" json:"code"`
	Name      string    `gorm:"column:pos_name" json:"name"`
	IsActive  bool      `gorm:"default:true;column:pos_is_active" json:"is_active"`
	CreatedAt time.Time `gorm:"column:pos_created_at" json:"created_at"`
	CreatedBy string    `gorm:"column:pos_created_by" json:"created_by"`
	UpdatedAt time.Time `gorm:"column:pos_updated_at" json:"updated_at"`
	UpdatedBy string    `gorm:"column:pos_updated_by" json:"updated_by"`
}

// TableName overrides the table name
func (Position) TableName() string {
	return "user.positions"
}

// Role represents the roles table
type Role struct {
	ID        uint      `gorm:"primaryKey;column:role_id" json:"id"`
	Name      string    `gorm:"unique;column:role_name" json:"name"`
	Level     int       `gorm:"column:role_level" json:"level"`
	IsActive  bool      `gorm:"default:true;column:role_is_active" json:"is_active"`
	CreatedAt time.Time `gorm:"column:role_created_at" json:"created_at"`
	CreatedBy string    `gorm:"column:role_created_by" json:"created_by"`
	UpdatedAt time.Time `gorm:"column:role_updated_at" json:"updated_at"`
	UpdatedBy string    `gorm:"column:role_updated_by" json:"updated_by"`
	// Relations
	Users []*User `gorm:"many2many:user_roles;foreignKey:role_id;joinForeignKey:ur_role_id;References:u_id;joinReferences:ur_user_id" json:"users,omitempty"`
}

// TableName overrides the table name
func (Role) TableName() string {
	return "user.roles"
}

// User represents the users table
type User struct {
	ID           uint       `gorm:"primaryKey;column:u_id" json:"id"`
	UID          uuid.UUID  `gorm:"type:uuid;unique;column:u_uid;default:gen_random_uuid()" json:"uid"`
	EmployeeID   string     `gorm:"unique;column:u_employee_id" json:"employee_id"`
	Name         string     `gorm:"column:u_name" json:"name"`
	Email        string     `gorm:"unique;column:u_email" json:"email"`
	Password     string     `gorm:"column:u_password" json:"-"` // Never return password in JSON
	Phone        string     `gorm:"column:u_phone" json:"phone"`
	Address      string     `gorm:"column:u_address" json:"address"`
	Birthdate    *time.Time `gorm:"column:u_birthdate" json:"birthdate"`
	JoinDate     time.Time  `gorm:"column:u_join_date" json:"join_date"`
	ProfileImage string     `gorm:"column:u_profile_image" json:"profile_image"`
	DivisionID   *uint      `gorm:"column:u_division_id" json:"division_id"`
	PositionID   *uint      `gorm:"column:u_position_id" json:"position_id"`
	IsManager    bool       `gorm:"default:false;column:u_is_manager" json:"is_manager"`
	ManagerID    *uint      `gorm:"column:u_manager_id" json:"manager_id"`
	IsActive     bool       `gorm:"default:true;column:u_is_active" json:"is_active"`
	CreatedAt    time.Time  `gorm:"column:u_created_at" json:"created_at"`
	CreatedBy    string     `gorm:"column:u_created_by" json:"created_by"`
	UpdatedAt    time.Time  `gorm:"column:u_updated_at" json:"updated_at"`
	UpdatedBy    string     `gorm:"column:u_updated_by" json:"updated_by"`
	// Relations
	Division  *Division  `gorm:"foreignKey:u_division_id;references:div_id" json:"division,omitempty"`
	Position  *Position  `gorm:"foreignKey:u_position_id;references:pos_id" json:"position,omitempty"`
	Manager   *User      `gorm:"foreignKey:u_manager_id;references:u_id" json:"manager,omitempty"`
	Roles     []*Role    `gorm:"many2many:user_roles;foreignKey:u_id;joinForeignKey:ur_user_id;References:role_id;joinReferences:ur_role_id" json:"roles,omitempty"`
	Employees []*User    `gorm:"foreignKey:u_manager_id;references:u_id" json:"employees,omitempty"`
	UserRoles []UserRole `gorm:"foreignKey:ur_user_id;references:u_id" json:"user_roles,omitempty"`
}

// TableName overrides the table name
func (User) TableName() string {
	return "user.users"
}

// UserRole represents the user_roles table (many-to-many relationship)
type UserRole struct {
	ID        uint      `gorm:"primaryKey;column:ur_id" json:"id"`
	UserID    uint      `gorm:"column:ur_user_id" json:"user_id"`
	RoleID    uint      `gorm:"column:ur_role_id" json:"role_id"`
	CreatedAt time.Time `gorm:"column:ur_created_at" json:"created_at"`
	CreatedBy string    `gorm:"column:ur_created_by" json:"created_by"`
	// Relations
	User *User `gorm:"foreignKey:ur_user_id;references:u_id" json:"user,omitempty"`
	Role *Role `gorm:"foreignKey:ur_role_id;references:role_id" json:"role,omitempty"`
}

// TableName overrides the table name
func (UserRole) TableName() string {
	return "user.user_roles"
}

// DTOs (Data Transfer Objects)

// UserLoginRequest represents login request payload
type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// UserResponse represents user data without sensitive information
type UserResponse struct {
	ID           uint      `json:"id"`
	UID          uuid.UUID `json:"uid"`
	EmployeeID   string    `json:"employee_id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone,omitempty"`
	Address      string    `json:"address,omitempty"`
	Birthdate    string    `json:"birthdate,omitempty"`
	JoinDate     string    `json:"join_date"`
	ProfileImage string    `json:"profile_image,omitempty"`
	Division     string    `json:"division,omitempty"`
	Position     string    `json:"position,omitempty"`
	IsManager    bool      `json:"is_manager"`
	Manager      string    `json:"manager,omitempty"`
	IsActive     bool      `json:"is_active"`
	Roles        []string  `json:"roles,omitempty"`
}

// CreateUserRequest represents payload for creating a new user
type CreateUserRequest struct {
	EmployeeID   string    `json:"employee_id" binding:"required"`
	Name         string    `json:"name" binding:"required"`
	Email        string    `json:"email" binding:"required,email"`
	Password     string    `json:"password" binding:"required,min=6"`
	Phone        string    `json:"phone"`
	Address      string    `json:"address"`
	Birthdate    string    `json:"birthdate"`
	JoinDate     string    `json:"join_date" binding:"required"`
	ProfileImage string    `json:"profile_image"`
	DivisionID   *uint     `json:"division_id"`
	PositionID   *uint     `json:"position_id"`
	IsManager    bool      `json:"is_manager"`
	ManagerID    *uint     `json:"manager_id"`
	RoleIDs      []uint    `json:"role_ids"`
}

// UpdateUserRequest represents payload for updating a user
type UpdateUserRequest struct {
	Name         string    `json:"name"`
	Email        string    `json:"email" binding:"email"`
	Phone        string    `json:"phone"`
	Address      string    `json:"address"`
	Birthdate    string    `json:"birthdate"`
	JoinDate     string    `json:"join_date"`
	ProfileImage string    `json:"profile_image"`
	DivisionID   *uint     `json:"division_id"`
	PositionID   *uint     `json:"position_id"`
	IsManager    *bool     `json:"is_manager"`
	ManagerID    *uint     `json:"manager_id"`
	IsActive     *bool     `json:"is_active"`
	RoleIDs      []uint    `json:"role_ids"`
}

// LoginResponse represents response after successful login
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// DivisionRequest represents payload for creating/updating division
type DivisionRequest struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

// PositionRequest represents payload for creating/updating position
type PositionRequest struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

// RoleRequest represents payload for creating/updating role
type RoleRequest struct {
	Name  string `json:"name" binding:"required"`
	Level int    `json:"level" binding:"required"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	TotalItems  int64       `json:"total_items"`
	TotalPages  int64       `json:"total_pages"`
	CurrentPage int64       `json:"current_page"`
	PageSize    int64       `json:"page_size"`
	Items       interface{} `json:"items"`
}

// Statistics represents dashboard statistics
type Statistics struct {
	TotalUsers        int64 `json:"total_users"`
	ActiveUsers       int64 `json:"active_users"`
	TotalDivisions    int64 `json:"total_divisions"`
	TotalPositions    int64 `json:"total_positions"`
	UsersPerDivision  []map[string]interface{} `json:"users_per_division"`
	UsersPerPosition  []map[string]interface{} `json:"users_per_position"`
	NewUsersThisMonth int64 `json:"new_users_this_month"`
}