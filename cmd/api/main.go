package main

import (
	"log"
	"os"

	"admin-dashboard/internal/config"
	"admin-dashboard/internal/handlers"
	"admin-dashboard/internal/middleware"
	"admin-dashboard/internal/repository"
	"admin-dashboard/internal/services"
	"admin-dashboard/internal/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// Add extensive logging
	log.Println("========== STARTUP DIAGNOSTICS ==========")
	log.Println("Current working directory:", getWorkingDir())
	log.Println("Environment variables:")
	for _, env := range os.Environ() {
		log.Println("  ", env)
	}

	// Force production mode for Gin
	gin.SetMode(gin.ReleaseMode)
	
	// Load configuration
	log.Println("Loading configuration...")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	
	// Create database connection
	log.Println("Connecting to database...")
	db, err := config.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	log.Println("Running database migrations...")
	err = db.Migrate()
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Initializing application components...")
	// Initialize JWT manager
	jwtManager := utils.NewJWTManager(&cfg.JWTConfig)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.DB)
	roleRepo := repository.NewRoleRepository(db.DB)
	divisionRepo := repository.NewDivisionRepository(db.DB)
	positionRepo := repository.NewPositionRepository(db.DB)

	// Initialize services
	authService := services.NewAuthService(userRepo, roleRepo, jwtManager)
	userService := services.NewUserService(userRepo, roleRepo, divisionRepo, positionRepo)
	roleService := services.NewRoleService(roleRepo)
	divisionService := services.NewDivisionService(divisionRepo)
	positionService := services.NewPositionService(positionRepo)
	dashboardService := services.NewDashboardService(db.DB)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	roleHandler := handlers.NewRoleHandler(roleService)
	divisionHandler := handlers.NewDivisionHandler(divisionService)
	positionHandler := handlers.NewPositionHandler(positionService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtManager)
	authenticate := authMiddleware.Authenticate()

	// Set up Gin router
	log.Println("Setting up HTTP router...")
	router := gin.New() // Use New() instead of Default() to avoid default middleware
	
	// Explicitly add recovery middleware
	router.Use(gin.Recovery())

	// Apply global middleware
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())
	router.Use(middleware.ErrorHandler())
	
	// Set trusted proxies for Railway
	router.SetTrustedProxies(nil) // Trust all proxies in Railway
	
	// Add health check endpoint
	router.GET("/health", func(c *gin.Context) {
		log.Println("Health check endpoint called!")
		c.String(200, "OK")
	})

	// API routes
	api := router.Group("/api")
	{
		// Auth routes with middleware for profile
		authHandler.RegisterRoutes(api, &authenticate)

		// Protected routes (authentication required)
		userHandler.RegisterRoutes(api, &authenticate)
		roleHandler.RegisterRoutes(api, &authenticate)
		divisionHandler.RegisterRoutes(api, &authenticate)
		positionHandler.RegisterRoutes(api, &authenticate)
		dashboardHandler.RegisterRoutes(api, &authenticate)
	}

	// Get port from environment with fallback
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Log server startup
	log.Printf("Starting HTTP server on port %s", port)
	
	// Use an explicit binding address
	if err := router.Run("0.0.0.0:" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getWorkingDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "ERROR: " + err.Error()
	}
	return dir
}