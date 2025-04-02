package main

import (
	// "fmt"
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
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create database connection
	db, err := config.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	err = db.Migrate()
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

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
	router := gin.Default()

	// Apply global middleware
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())
	router.Use(middleware.ErrorHandler())

	// API routes
	api := router.Group("/api")
	{
		// Auth routes dengan middleware untuk profile
		authHandler.RegisterRoutes(api, &authenticate)

		// Protected routes (authentication required)
		userHandler.RegisterRoutes(api, &authenticate)
		roleHandler.RegisterRoutes(api, &authenticate)
		divisionHandler.RegisterRoutes(api, &authenticate)
		positionHandler.RegisterRoutes(api, &authenticate)
		dashboardHandler.RegisterRoutes(api, &authenticate)
	}

	// Add health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.String(200, "OK")
	})

	// Start the server
	// serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	// log.Printf("Server is running on %s", serverAddr)
	// if err := router.Run(serverAddr); err != nil {
	// 	log.Fatalf("Failed to start server: %v", err)
	// }

	// router.SetTrustedProxies([]string{"0.0.0.0/0", "::/0"})
	// router.Run(":" + cfg.Server.Port) // This binds to all interfaces, IPv4 and IPv6

	// Add right before router.Run
	envs := os.Environ()
	log.Println("Environment variables:")
	for _, env := range envs {
		log.Println(env)
	}
	log.Println("Attempting to start HTTP server...")

	log.Printf("Starting server with explicit binding to 0.0.0.0:3000")
	if err := router.Run("0.0.0.0:3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}