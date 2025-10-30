package handlers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/Ismael-Njihia/Kenya-info-api/internal/models"
	"github.com/Ismael-Njihia/Kenya-info-api/internal/services"
	"github.com/gin-gonic/gin"
)

type ConstituencyHandler struct {
	service *services.ConstituencyService
}

func NewConstituencyHandler(service *services.ConstituencyService) *ConstituencyHandler {
	return &ConstituencyHandler{service: service}
}

// GetAll godoc
// @Summary Get all constituencies
// @Description Get a list of all constituencies with pagination
// @Tags constituencies
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(50)
// @Success 200 {object} models.PaginatedResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/v1/constituencies [get]
func (h *ConstituencyHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 50
	}

	constituencies, total, err := h.service.GetAll(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to fetch constituencies",
		})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success:    true,
		Data:       constituencies,
		Page:       page,
		PageSize:   pageSize,
		TotalCount: total,
		TotalPages: totalPages,
	})
}

// GetByID godoc
// @Summary Get constituency by ID
// @Description Get a single constituency by its ID
// @Tags constituencies
// @Accept json
// @Produce json
// @Param id path string true "Constituency ID"
// @Success 200 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/v1/constituencies/{id} [get]
func (h *ConstituencyHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	constituency, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "constituency not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    constituency,
	})
}

// GetByCountyID godoc
// @Summary Get constituencies by county ID
// @Description Get all constituencies belonging to a specific county
// @Tags constituencies
// @Accept json
// @Produce json
// @Param county_id path string true "County ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(50)
// @Success 200 {object} models.PaginatedResponse
// @Failure 400 {object} models.APIResponse
// @Router /api/v1/counties/{county_id}/constituencies [get]
func (h *ConstituencyHandler) GetByCountyID(c *gin.Context) {
	countyID := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 50
	}

	constituencies, total, err := h.service.GetByCountyID(c.Request.Context(), countyID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success:    true,
		Data:       constituencies,
		Page:       page,
		PageSize:   pageSize,
		TotalCount: total,
		TotalPages: totalPages,
	})
}

// Create godoc
// @Summary Create a new constituency
// @Description Create a new constituency record
// @Tags constituencies
// @Accept json
// @Produce json
// @Param constituency body models.Constituency true "Constituency object"
// @Success 201 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Router /api/v1/constituencies [post]
func (h *ConstituencyHandler) Create(c *gin.Context) {
	var constituency models.Constituency
	if err := c.ShouldBindJSON(&constituency); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if err := h.service.Create(c.Request.Context(), &constituency); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to create constituency",
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Constituency created successfully",
		Data:    constituency,
	})
}

// Update godoc
// @Summary Update a constituency
// @Description Update an existing constituency record
// @Tags constituencies
// @Accept json
// @Produce json
// @Param id path string true "Constituency ID"
// @Param constituency body models.Constituency true "Constituency object"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/v1/constituencies/{id} [put]
func (h *ConstituencyHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var constituency models.Constituency
	if err := c.ShouldBindJSON(&constituency); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if err := h.service.Update(c.Request.Context(), id, &constituency); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "constituency not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Constituency updated successfully",
	})
}

// Delete godoc
// @Summary Delete a constituency
// @Description Delete a constituency record
// @Tags constituencies
// @Accept json
// @Produce json
// @Param id path string true "Constituency ID"
// @Success 200 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/v1/constituencies/{id} [delete]
func (h *ConstituencyHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "constituency not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Constituency deleted successfully",
	})
}
