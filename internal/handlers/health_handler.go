package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	startTime time.Time
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{
		startTime: time.Now(),
	}
}

// HealthCheck godoc
// @Summary Health check
// @Description Check if the API is running
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	uptime := time.Since(h.startTime)
	
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"uptime":  uptime.String(),
		"service": "Kenya Info API",
		"version": "1.0.0",
	})
}
