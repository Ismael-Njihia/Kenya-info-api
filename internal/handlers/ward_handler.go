package handlers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/Ismael-Njihia/Kenya-info-api/internal/models"
	"github.com/Ismael-Njihia/Kenya-info-api/internal/services"
	"github.com/gin-gonic/gin"
)

type WardHandler struct {
	service *services.WardService
}

func NewWardHandler(service *services.WardService) *WardHandler {
	return &WardHandler{service: service}
}

// GetAll godoc
// @Summary Get all wards
// @Description Get a list of all wards with pagination
// @Tags wards
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(50)
// @Success 200 {object} models.PaginatedResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/v1/wards [get]
func (h *WardHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 50
	}

	wards, total, err := h.service.GetAll(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to fetch wards",
		})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success:    true,
		Data:       wards,
		Page:       page,
		PageSize:   pageSize,
		TotalCount: total,
		TotalPages: totalPages,
	})
}

// GetByID godoc
// @Summary Get ward by ID
// @Description Get a single ward by its ID
// @Tags wards
// @Accept json
// @Produce json
// @Param id path string true "Ward ID"
// @Success 200 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/v1/wards/{id} [get]
func (h *WardHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	ward, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "ward not found" {
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
		Data:    ward,
	})
}

// GetByCountyID godoc
// @Summary Get wards by county ID
// @Description Get all wards belonging to a specific county
// @Tags wards
// @Accept json
// @Produce json
// @Param county_id path string true "County ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(50)
// @Success 200 {object} models.PaginatedResponse
// @Failure 400 {object} models.APIResponse
// @Router /api/v1/counties/{county_id}/wards [get]
func (h *WardHandler) GetByCountyID(c *gin.Context) {
	countyID := c.Param("county_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 50
	}

	wards, total, err := h.service.GetByCountyID(c.Request.Context(), countyID, page, pageSize)
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
		Data:       wards,
		Page:       page,
		PageSize:   pageSize,
		TotalCount: total,
		TotalPages: totalPages,
	})
}

// Create godoc
// @Summary Create a new ward
// @Description Create a new ward record
// @Tags wards
// @Accept json
// @Produce json
// @Param ward body models.Ward true "Ward object"
// @Success 201 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Router /api/v1/wards [post]
func (h *WardHandler) Create(c *gin.Context) {
	var ward models.Ward
	if err := c.ShouldBindJSON(&ward); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if err := h.service.Create(c.Request.Context(), &ward); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to create ward",
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Ward created successfully",
		Data:    ward,
	})
}

// Update godoc
// @Summary Update a ward
// @Description Update an existing ward record
// @Tags wards
// @Accept json
// @Produce json
// @Param id path string true "Ward ID"
// @Param ward body models.Ward true "Ward object"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/v1/wards/{id} [put]
func (h *WardHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var ward models.Ward
	if err := c.ShouldBindJSON(&ward); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if err := h.service.Update(c.Request.Context(), id, &ward); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "ward not found" {
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
		Message: "Ward updated successfully",
	})
}

// Delete godoc
// @Summary Delete a ward
// @Description Delete a ward record
// @Tags wards
// @Accept json
// @Produce json
// @Param id path string true "Ward ID"
// @Success 200 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/v1/wards/{id} [delete]
func (h *WardHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "ward not found" {
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
		Message: "Ward deleted successfully",
	})
}
