package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/eikuma/stockle/backend/internal/config"
	"github.com/eikuma/stockle/backend/internal/controllers"
	"github.com/eikuma/stockle/backend/internal/database"
	"github.com/eikuma/stockle/backend/internal/middleware"
	"github.com/eikuma/stockle/backend/internal/repositories"
	"github.com/eikuma/stockle/backend/internal/services"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	if err := middleware.InitLogger(cfg); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Run database migrations
	log.Printf("Running database migrations...")
	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}
	log.Printf("Database migrations completed successfully")

	// Set gin mode
	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Initialize Gin router
	router := setupRouter(cfg)

	// Setup HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func setupRouter(cfg *config.Config) *gin.Engine {
	router := gin.New()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(database.GetDB())

	// Initialize services
	authService := services.NewAuthService(userRepo, &cfg.JWT)

	// Initialize controllers
	healthController := controllers.NewHealthController(cfg)
	authController := controllers.NewAuthController(authService)
	userController := controllers.NewUserController(userRepo)

	// Initialize rate limiter
	rateLimiter := middleware.NewRateLimiter()

	// Middleware
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.SecurityHeaders())
	router.Use(middleware.CORS(cfg))

	// API routes
	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			// Health endpoints
			v1.GET("/health", healthController.Health)
			v1.GET("/health/ready", healthController.Ready)
			v1.GET("/health/live", healthController.Live)

			// Authentication endpoints (with rate limiting) - Google OAuth専用
			auth := v1.Group("/auth")
			auth.Use(rateLimiter.RateLimit(20, time.Minute)) // Google OAuthのみなので制限を緩和
			{
				// Google OAuth専用エンドポイント
				auth.POST("/google", authController.GoogleAuth)
				
				// 認証後の操作用エンドポイント
				auth.POST("/refresh", authController.RefreshToken)
				auth.POST("/logout", authController.Logout)
				auth.GET("/me", middleware.AuthRequired(authService), authController.Me)
				
				// 従来のエンドポイントは無効化（404を返す）
				// auth.POST("/register", authController.Register) // 無効化
				// auth.POST("/login", authController.Login)       // 無効化
			}

			// User endpoints
			users := v1.Group("/users")
			users.Use(middleware.AuthRequired(authService)) // All user endpoints require authentication
			{
				users.GET("/me", userController.GetProfile)
				users.PUT("/me", userController.UpdateProfile)
			}
		}
	}

	// Root health endpoint
	router.GET("/health", healthController.Health)

	return router
}