package handlers

import (
	"net/http"

	"github.com/Ismael-Njihia/Kenya-info-api/internal/logger"
	"github.com/Ismael-Njihia/Kenya-info-api/internal/models"
	"github.com/Ismael-Njihia/Kenya-info-api/internal/repository"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

// LeaderHandler handles leader-related requests
type LeaderHandler struct {
	repo *repository.LeaderRepository
}

// NewLeaderHandler creates a new leader handler
func NewLeaderHandler(repo *repository.LeaderRepository) *LeaderHandler {
	return &LeaderHandler{repo: repo}
}

// GetAllLeaders godoc
// @Summary Get all leaders
// @Description Get all leaders in Kenya
// @Tags leaders
// @Accept json
// @Produce json
// @Success 200 {array} models.Leader
// @Failure 500 {object} map[string]string
// @Router /api/v1/leaders [get]
func (h *LeaderHandler) GetAllLeaders(c *gin.Context) {
	leaders, err := h.repo.GetAll(c.Request.Context())
	if err != nil {
		logger.Log.Error("Failed to get leaders", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve leaders"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": leaders, "count": len(leaders)})
}

// GetLeaderByID godoc
// @Summary Get leader by ID
// @Description Get a specific leader by their ID
// @Tags leaders
// @Accept json
// @Produce json
// @Param id path string true "Leader ID"
// @Success 200 {object} models.Leader
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/leaders/{id} [get]
func (h *LeaderHandler) GetLeaderByID(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid leader ID"})
		return
	}

	leader, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		logger.Log.Error("Failed to get leader", zap.Error(err), zap.String("id", id.Hex()))
		c.JSON(http.StatusNotFound, gin.H{"error": "Leader not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": leader})
}

// GetLeadersByPosition godoc
// @Summary Get leaders by position
// @Description Get all leaders with a specific position
// @Tags leaders
// @Accept json
// @Produce json
// @Param position path string true "Position"
// @Success 200 {array} models.Leader
// @Failure 500 {object} map[string]string
// @Router /api/v1/leaders/position/{position} [get]
func (h *LeaderHandler) GetLeadersByPosition(c *gin.Context) {
	position := c.Param("position")
	leaders, err := h.repo.GetByPosition(c.Request.Context(), position)
	if err != nil {
		logger.Log.Error("Failed to get leaders by position", zap.Error(err), zap.String("position", position))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve leaders"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": leaders, "count": len(leaders)})
}

// GetLeadersByCounty godoc
// @Summary Get leaders by county
// @Description Get all leaders from a specific county
// @Tags leaders
// @Accept json
// @Produce json
// @Param county path string true "County name"
// @Success 200 {array} models.Leader
// @Failure 500 {object} map[string]string
// @Router /api/v1/leaders/county/{county} [get]
func (h *LeaderHandler) GetLeadersByCounty(c *gin.Context) {
	county := c.Param("county")
	leaders, err := h.repo.GetByCounty(c.Request.Context(), county)
	if err != nil {
		logger.Log.Error("Failed to get leaders by county", zap.Error(err), zap.String("county", county))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve leaders"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": leaders, "count": len(leaders)})
}

// CreateLeader godoc
// @Summary Create a new leader
// @Description Create a new leader record
// @Tags leaders
// @Accept json
// @Produce json
// @Param leader body models.Leader true "Leader data"
// @Success 201 {object} models.Leader
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/leaders [post]
func (h *LeaderHandler) CreateLeader(c *gin.Context) {
	var leader models.Leader
	if err := c.ShouldBindJSON(&leader); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Create(c.Request.Context(), &leader); err != nil {
		logger.Log.Error("Failed to create leader", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create leader"})
		return
	}

	logger.Log.Info("Leader created successfully", zap.String("id", leader.ID.Hex()))
	c.JSON(http.StatusCreated, gin.H{"data": leader})
}

// UpdateLeader godoc
// @Summary Update a leader
// @Description Update an existing leader record
// @Tags leaders
// @Accept json
// @Produce json
// @Param id path string true "Leader ID"
// @Param leader body models.Leader true "Leader data"
// @Success 200 {object} models.Leader
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/leaders/{id} [put]
func (h *LeaderHandler) UpdateLeader(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid leader ID"})
		return
	}

	var leader models.Leader
	if err := c.ShouldBindJSON(&leader); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Update(c.Request.Context(), id, &leader); err != nil {
		logger.Log.Error("Failed to update leader", zap.Error(err), zap.String("id", id.Hex()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update leader"})
		return
	}

	logger.Log.Info("Leader updated successfully", zap.String("id", id.Hex()))
	leader.ID = id
	c.JSON(http.StatusOK, gin.H{"data": leader})
}

// DeleteLeader godoc
// @Summary Delete a leader
// @Description Delete a leader record
// @Tags leaders
// @Accept json
// @Produce json
// @Param id path string true "Leader ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/leaders/{id} [delete]
func (h *LeaderHandler) DeleteLeader(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid leader ID"})
		return
	}

	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		logger.Log.Error("Failed to delete leader", zap.Error(err), zap.String("id", id.Hex()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete leader"})
		return
	}

	logger.Log.Info("Leader deleted successfully", zap.String("id", id.Hex()))
	c.JSON(http.StatusOK, gin.H{"message": "Leader deleted successfully"})
}
