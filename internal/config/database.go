package config

import (
	"fmt"
	"log"
	"time"

	"admin-dashboard/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database represents the database connection
type Database struct {
	DB *gorm.DB
}

// NewDatabase creates a new database connection
func NewDatabase(config *Config) (*Database, error) {
	// Configure GORM logger
	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	fmt.Println("DSN:", config.DBConfig.DSN)

	// Connect to database
	db, err := gorm.Open(postgres.Open(config.DBConfig.DSN), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Connected to PostgreSQL database!")

	return &Database{DB: db}, nil
}

// Migrate automigrates all models
func (d *Database) Migrate() error {
	log.Println("Running database migrations...")

	// Auto-migrate models
	err := d.DB.AutoMigrate(
		&models.Division{},
		&models.Position{},
		&models.Role{},
		&models.User{},
		&models.UserRole{},
	)
	if err != nil {
		return fmt.Errorf("failed to automigrate: %w", err)
	}

	log.Println("Database migrations completed successfully!")
	return nil
}