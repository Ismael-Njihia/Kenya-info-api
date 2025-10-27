package handlers

import (
	"net/http"

	"github.com/Ismael-Njihia/Kenya-info-api/internal/repository"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	db *mongo.Database
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *mongo.Database) *HealthHandler {
	return &HealthHandler{db: db}
}

// HealthCheck godoc
// @Summary Health check
// @Description Check if the API is running and database is connected
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	// Try to ping the database
	ctx := c.Request.Context()
	err := h.db.Client().Ping(ctx, nil)
	
	status := "healthy"
	dbStatus := "connected"
	statusCode := http.StatusOK
	
	if err != nil {
		status = "unhealthy"
		dbStatus = "disconnected"
		statusCode = http.StatusServiceUnavailable
	}

	countyRepo := repository.NewCountyRepository(h.db)
	wardRepo := repository.NewWardRepository(h.db)
	leaderRepo := repository.NewLeaderRepository(h.db)

	countyCount, _ := countyRepo.Count(ctx)
	wardCount, _ := wardRepo.Count(ctx)
	leaderCount, _ := leaderRepo.Count(ctx)

	c.JSON(statusCode, gin.H{
		"status": status,
		"database": dbStatus,
		"data": gin.H{
			"counties": countyCount,
			"wards":    wardCount,
			"leaders":  leaderCount,
		},
	})
}
