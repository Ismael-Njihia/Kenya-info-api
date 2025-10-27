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

// WardHandler handles ward-related requests
type WardHandler struct {
	repo *repository.WardRepository
}

// NewWardHandler creates a new ward handler
func NewWardHandler(repo *repository.WardRepository) *WardHandler {
	return &WardHandler{repo: repo}
}

// GetAllWards godoc
// @Summary Get all wards
// @Description Get all wards in Kenya
// @Tags wards
// @Accept json
// @Produce json
// @Success 200 {array} models.Ward
// @Failure 500 {object} map[string]string
// @Router /api/v1/wards [get]
func (h *WardHandler) GetAllWards(c *gin.Context) {
	wards, err := h.repo.GetAll(c.Request.Context())
	if err != nil {
		logger.Log.Error("Failed to get wards", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve wards"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": wards, "count": len(wards)})
}

// GetWardByID godoc
// @Summary Get ward by ID
// @Description Get a specific ward by its ID
// @Tags wards
// @Accept json
// @Produce json
// @Param id path string true "Ward ID"
// @Success 200 {object} models.Ward
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/wards/{id} [get]
func (h *WardHandler) GetWardByID(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ward ID"})
		return
	}

	ward, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		logger.Log.Error("Failed to get ward", zap.Error(err), zap.String("id", id.Hex()))
		c.JSON(http.StatusNotFound, gin.H{"error": "Ward not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": ward})
}

// GetWardsByCountyID godoc
// @Summary Get wards by county ID
// @Description Get all wards in a specific county
// @Tags wards
// @Accept json
// @Produce json
// @Param county_id path string true "County ID"
// @Success 200 {array} models.Ward
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/wards/county/{county_id} [get]
func (h *WardHandler) GetWardsByCountyID(c *gin.Context) {
	countyID, err := primitive.ObjectIDFromHex(c.Param("county_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid county ID"})
		return
	}

	wards, err := h.repo.GetByCountyID(c.Request.Context(), countyID)
	if err != nil {
		logger.Log.Error("Failed to get wards by county", zap.Error(err), zap.String("county_id", countyID.Hex()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve wards"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": wards, "count": len(wards)})
}

// CreateWard godoc
// @Summary Create a new ward
// @Description Create a new ward record
// @Tags wards
// @Accept json
// @Produce json
// @Param ward body models.Ward true "Ward data"
// @Success 201 {object} models.Ward
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/wards [post]
func (h *WardHandler) CreateWard(c *gin.Context) {
	var ward models.Ward
	if err := c.ShouldBindJSON(&ward); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Create(c.Request.Context(), &ward); err != nil {
		logger.Log.Error("Failed to create ward", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ward"})
		return
	}

	logger.Log.Info("Ward created successfully", zap.String("id", ward.ID.Hex()))
	c.JSON(http.StatusCreated, gin.H{"data": ward})
}

// UpdateWard godoc
// @Summary Update a ward
// @Description Update an existing ward record
// @Tags wards
// @Accept json
// @Produce json
// @Param id path string true "Ward ID"
// @Param ward body models.Ward true "Ward data"
// @Success 200 {object} models.Ward
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/wards/{id} [put]
func (h *WardHandler) UpdateWard(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ward ID"})
		return
	}

	var ward models.Ward
	if err := c.ShouldBindJSON(&ward); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Update(c.Request.Context(), id, &ward); err != nil {
		logger.Log.Error("Failed to update ward", zap.Error(err), zap.String("id", id.Hex()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ward"})
		return
	}

	logger.Log.Info("Ward updated successfully", zap.String("id", id.Hex()))
	ward.ID = id
	c.JSON(http.StatusOK, gin.H{"data": ward})
}

// DeleteWard godoc
// @Summary Delete a ward
// @Description Delete a ward record
// @Tags wards
// @Accept json
// @Produce json
// @Param id path string true "Ward ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/wards/{id} [delete]
func (h *WardHandler) DeleteWard(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ward ID"})
		return
	}

	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		logger.Log.Error("Failed to delete ward", zap.Error(err), zap.String("id", id.Hex()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete ward"})
		return
	}

	logger.Log.Info("Ward deleted successfully", zap.String("id", id.Hex()))
	c.JSON(http.StatusOK, gin.H{"message": "Ward deleted successfully"})
}
