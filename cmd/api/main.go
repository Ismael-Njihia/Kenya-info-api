package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ismael-Njihia/Kenya-info-api/config"
	"github.com/Ismael-Njihia/Kenya-info-api/internal/handlers"
	"github.com/Ismael-Njihia/Kenya-info-api/internal/logger"
	"github.com/Ismael-Njihia/Kenya-info-api/internal/middleware"
	"github.com/Ismael-Njihia/Kenya-info-api/internal/repository"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	_ "github.com/Ismael-Njihia/Kenya-info-api/docs"
)

// @title Kenya Info API
// @version 1.0
// @description Fast and reliable API for Kenya's counties, wards, and leaders data
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url https://github.com/Ismael-Njihia/Kenya-info-api
// @contact.email support@kenyainfo.api

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http https

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	if err := logger.InitLogger(cfg.Logger.Level, cfg.Logger.OutputPath); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	logger.Log.Info("Starting Kenya Info API")

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Database.Timeout)*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.Database.URI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Log.Fatal("Failed to connect to MongoDB", zap.Error(err))
	}
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			logger.Log.Error("Failed to disconnect from MongoDB", zap.Error(err))
		}
	}()

	// Ping database
	if err = client.Ping(ctx, nil); err != nil {
		logger.Log.Fatal("Failed to ping MongoDB", zap.Error(err))
	}
	logger.Log.Info("Connected to MongoDB successfully")

	db := client.Database(cfg.Database.Database)

	// Initialize repositories
	countyRepo := repository.NewCountyRepository(db)
	wardRepo := repository.NewWardRepository(db)
	leaderRepo := repository.NewLeaderRepository(db)

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler(db)
	countyHandler := handlers.NewCountyHandler(countyRepo)
	wardHandler := handlers.NewWardHandler(wardRepo)
	leaderHandler := handlers.NewLeaderHandler(leaderRepo)

	// Setup Gin
	gin.SetMode(cfg.Server.Mode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())
	router.Use(middleware.RequestLogger())

	// Health check
	router.GET("/health", healthHandler.HealthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// County routes
		counties := v1.Group("/counties")
		{
			counties.GET("", countyHandler.GetAllCounties)
			counties.GET("/:id", countyHandler.GetCountyByID)
			counties.GET("/code/:code", countyHandler.GetCountyByCode)
			counties.POST("", countyHandler.CreateCounty)
			counties.PUT("/:id", countyHandler.UpdateCounty)
			counties.DELETE("/:id", countyHandler.DeleteCounty)
		}

		// Ward routes
		wards := v1.Group("/wards")
		{
			wards.GET("", wardHandler.GetAllWards)
			wards.GET("/:id", wardHandler.GetWardByID)
			wards.GET("/county/:county_id", wardHandler.GetWardsByCountyID)
			wards.POST("", wardHandler.CreateWard)
			wards.PUT("/:id", wardHandler.UpdateWard)
			wards.DELETE("/:id", wardHandler.DeleteWard)
		}

		// Leader routes
		leaders := v1.Group("/leaders")
		{
			leaders.GET("", leaderHandler.GetAllLeaders)
			leaders.GET("/:id", leaderHandler.GetLeaderByID)
			leaders.GET("/position/:position", leaderHandler.GetLeadersByPosition)
			leaders.GET("/county/:county", leaderHandler.GetLeadersByCounty)
			leaders.POST("", leaderHandler.CreateLeader)
			leaders.PUT("/:id", leaderHandler.UpdateLeader)
			leaders.DELETE("/:id", leaderHandler.DeleteLeader)
		}
	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		logger.Log.Info("Server starting", zap.String("address", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Log.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Log.Info("Server exited")
}
