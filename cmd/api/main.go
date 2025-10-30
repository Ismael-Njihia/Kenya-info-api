package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ismael-Njihia/Kenya-info-api/internal/config"
	"github.com/Ismael-Njihia/Kenya-info-api/internal/database"
	"github.com/Ismael-Njihia/Kenya-info-api/internal/handlers"
	"github.com/Ismael-Njihia/Kenya-info-api/internal/middleware"
	"github.com/Ismael-Njihia/Kenya-info-api/internal/services"
	"github.com/Ismael-Njihia/Kenya-info-api/pkg/logger"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	_ "github.com/Ismael-Njihia/Kenya-info-api/docs"
)

// @title Kenya Info API
// @version 1.0
// @description A fast, reliable API providing structured data on Kenya's counties, wards, and leaders
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@kenya-info-api.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /
// @schemes http https
func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	if err := logger.Init(cfg.Logger.Level); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting Kenya Info API", zap.String("env", cfg.Server.Env))

	// Connect to MongoDB
	db, err := database.Connect(&cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	logger.Info("✅ Connected successfully to MongoDB Atlas")

	// Ensure database disconnects gracefully
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := db.Close(ctx); err != nil {
			logger.Error("Error disconnecting from MongoDB", zap.Error(err))
		}
	}()

	// Initialize services
	countyService := services.NewCountyService(db)
	constituencyService := services.NewConstituencyService(db)
	wardService := services.NewWardService(db)
	leaderService := services.NewLeaderService(db)

	// Initialize handlers
	countyHandler := handlers.NewCountyHandler(countyService)
	constituencyHandler := handlers.NewConstituencyHandler(constituencyService)
	wardHandler := handlers.NewWardHandler(wardService)
	leaderHandler := handlers.NewLeaderHandler(leaderService)
	healthHandler := handlers.NewHealthHandler()

	// Setup Gin
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS())

	// Health check
	router.GET("/health", healthHandler.HealthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// County routes
		counties := v1.Group("/counties")
		{
			counties.GET("", countyHandler.GetAll)
			counties.GET("/:id", countyHandler.GetByID)
			counties.GET("/code/:code", countyHandler.GetByCode)
			counties.GET("/search", countyHandler.GetByName)
			counties.POST("", countyHandler.Create)
			counties.PUT("/:id", countyHandler.Update)
			counties.DELETE("/:id", countyHandler.Delete)

			// ✅ Fixed wildcard conflict — consistent use of :id
			counties.GET("/:id/wards", wardHandler.GetByCountyID)
			counties.GET("/:id/leaders", leaderHandler.GetByCountyID)
			counties.GET("/:id/constituencies", constituencyHandler.GetByCountyID)
		}
		// Constituency routes
		constituencies := v1.Group("/constituencies")
		{
			constituencies.GET("", constituencyHandler.GetAll)
			constituencies.GET("/:id", constituencyHandler.GetByID)
			constituencies.POST("", constituencyHandler.Create)
			constituencies.PUT("/:id", constituencyHandler.Update)
			constituencies.DELETE("/:id", constituencyHandler.Delete)

			// ✅ Fixed wildcard conflict — consistent use of :id
			constituencies.GET("/:id/wards", wardHandler.GetByConstituencyID)
		}

		// Ward routes
		wards := v1.Group("/wards")
		{
			wards.GET("", wardHandler.GetAll)
			wards.GET("/:id", wardHandler.GetByID)
			wards.POST("", wardHandler.Create)
			wards.PUT("/:id", wardHandler.Update)
			wards.DELETE("/:id", wardHandler.Delete)
		}

		// Leader routes
		leaders := v1.Group("/leaders")
		{
			leaders.GET("", leaderHandler.GetAll)
			leaders.GET("/:id", leaderHandler.GetByID)
			leaders.GET("/position", leaderHandler.GetByPosition)
			leaders.POST("", leaderHandler.Create)
			leaders.PUT("/:id", leaderHandler.Update)
			leaders.DELETE("/:id", leaderHandler.Delete)
		}
	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Create HTTP server
	srv := &http.Server{
		Addr:           ":" + cfg.Server.Port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Start server
	go func() {
		logger.Info("Server starting", zap.String("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited gracefully")
}
