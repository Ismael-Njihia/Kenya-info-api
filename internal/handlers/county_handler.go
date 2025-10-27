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

// CountyHandler handles county-related requests
type CountyHandler struct {
	repo *repository.CountyRepository
}

// NewCountyHandler creates a new county handler
func NewCountyHandler(repo *repository.CountyRepository) *CountyHandler {
	return &CountyHandler{repo: repo}
}

// GetAllCounties godoc
// @Summary Get all counties
// @Description Get all counties in Kenya
// @Tags counties
// @Accept json
// @Produce json
// @Success 200 {array} models.County
// @Failure 500 {object} map[string]string
// @Router /api/v1/counties [get]
func (h *CountyHandler) GetAllCounties(c *gin.Context) {
	counties, err := h.repo.GetAll(c.Request.Context())
	if err != nil {
		logger.Log.Error("Failed to get counties", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve counties"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": counties, "count": len(counties)})
}

// GetCountyByID godoc
// @Summary Get county by ID
// @Description Get a specific county by its ID
// @Tags counties
// @Accept json
// @Produce json
// @Param id path string true "County ID"
// @Success 200 {object} models.County
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/counties/{id} [get]
func (h *CountyHandler) GetCountyByID(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid county ID"})
		return
	}

	county, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		logger.Log.Error("Failed to get county", zap.Error(err), zap.String("id", id.Hex()))
		c.JSON(http.StatusNotFound, gin.H{"error": "County not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": county})
}

// GetCountyByCode godoc
// @Summary Get county by code
// @Description Get a specific county by its code
// @Tags counties
// @Accept json
// @Produce json
// @Param code path int true "County Code"
// @Success 200 {object} models.County
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/counties/code/{code} [get]
func (h *CountyHandler) GetCountyByCode(c *gin.Context) {
	var params struct {
		Code int `uri:"code" binding:"required"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid county code"})
		return
	}

	county, err := h.repo.GetByCode(c.Request.Context(), params.Code)
	if err != nil {
		logger.Log.Error("Failed to get county by code", zap.Error(err), zap.Int("code", params.Code))
		c.JSON(http.StatusNotFound, gin.H{"error": "County not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": county})
}

// CreateCounty godoc
// @Summary Create a new county
// @Description Create a new county record
// @Tags counties
// @Accept json
// @Produce json
// @Param county body models.County true "County data"
// @Success 201 {object} models.County
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/counties [post]
func (h *CountyHandler) CreateCounty(c *gin.Context) {
	var county models.County
	if err := c.ShouldBindJSON(&county); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Create(c.Request.Context(), &county); err != nil {
		logger.Log.Error("Failed to create county", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create county"})
		return
	}

	logger.Log.Info("County created successfully", zap.String("id", county.ID.Hex()))
	c.JSON(http.StatusCreated, gin.H{"data": county})
}

// UpdateCounty godoc
// @Summary Update a county
// @Description Update an existing county record
// @Tags counties
// @Accept json
// @Produce json
// @Param id path string true "County ID"
// @Param county body models.County true "County data"
// @Success 200 {object} models.County
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/counties/{id} [put]
func (h *CountyHandler) UpdateCounty(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid county ID"})
		return
	}

	var county models.County
	if err := c.ShouldBindJSON(&county); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Update(c.Request.Context(), id, &county); err != nil {
		logger.Log.Error("Failed to update county", zap.Error(err), zap.String("id", id.Hex()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update county"})
		return
	}

	logger.Log.Info("County updated successfully", zap.String("id", id.Hex()))
	county.ID = id
	c.JSON(http.StatusOK, gin.H{"data": county})
}

// DeleteCounty godoc
// @Summary Delete a county
// @Description Delete a county record
// @Tags counties
// @Accept json
// @Produce json
// @Param id path string true "County ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/counties/{id} [delete]
func (h *CountyHandler) DeleteCounty(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid county ID"})
		return
	}

	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		logger.Log.Error("Failed to delete county", zap.Error(err), zap.String("id", id.Hex()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete county"})
		return
	}

	logger.Log.Info("County deleted successfully", zap.String("id", id.Hex()))
	c.JSON(http.StatusOK, gin.H{"message": "County deleted successfully"})
}
