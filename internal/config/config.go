package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for our application
type Config struct {
	DBConfig  DBConfig
	JWTConfig JWTConfig
	Server    ServerConfig
}

// DBConfig holds database related configuration
type DBConfig struct {
	Connection string
	Host       string
	Port       string
	Database   string
	Username   string
	Password   string
	SSLMode    string
	DSN        string
}

// JWTConfig holds JWT related configuration
type JWTConfig struct {
	Secret string
	Expiry int // in hours
}

// ServerConfig holds server related configuration
type ServerConfig struct {
	Host string
	Port string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: .env file not found, using environment variables")
	}

	// Database config
	dbConfig := DBConfig{
		Connection: getEnv("DB_CONNECTION", "postgresql"),
		Host:       getEnv("DB_HOST", "localhost"),
		Port:       getEnv("DB_PORT", "5432"),
		Database:   getEnv("DB_DATABASE", "postgres"),
		Username:   getEnv("DB_USERNAME", "postgres"),
		Password:   getEnv("DB_PASSWORD", ""),
		SSLMode:    getEnv("DB_SSL_MODE", "disable"),
	}

	// Generate DSN (Data Source Name)
	dbConfig.DSN = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.Username, dbConfig.Password, dbConfig.Database, dbConfig.SSLMode)

	// JWT config
	jwtExpiry, err := strconv.Atoi(getEnv("JWT_EXPIRY", "24"))
	if err != nil {
		jwtExpiry = 24 // Default to 24 hours
	}
	jwtConfig := JWTConfig{
		Secret: getEnv("JWT_SECRET", "default-jwt-secret-key"),
		Expiry: jwtExpiry,
	}

	// Server config
	serverConfig := ServerConfig{
		Host: getEnv("SERVER_HOST", "0.0.0.0"),
		Port: getEnv("SERVER_PORT", "8080"),
	}

	config := &Config{
		DBConfig:  dbConfig,
		JWTConfig: jwtConfig,
		Server:    serverConfig,
	}

	if os.Getenv("RAILWAY_ENVIRONMENT") == "production" {
		dbConfig.SSLMode = "require"
	}

	return config, nil
}

// Helper function to get an environment variable or return a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}