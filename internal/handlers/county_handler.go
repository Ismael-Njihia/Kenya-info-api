package handlers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/Ismael-Njihia/Kenya-info-api/internal/models"
	"github.com/Ismael-Njihia/Kenya-info-api/internal/services"
	"github.com/gin-gonic/gin"
)

type CountyHandler struct {
	service *services.CountyService
}

func NewCountyHandler(service *services.CountyService) *CountyHandler {
	return &CountyHandler{service: service}
}

// GetAll godoc
// @Summary Get all counties
// @Description Get a list of all counties with pagination
// @Tags counties
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(47)
// @Success 200 {object} models.PaginatedResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/v1/counties [get]
func (h *CountyHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "47"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 47
	}

	counties, total, err := h.service.GetAll(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to fetch counties",
		})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success:    true,
		Data:       counties,
		Page:       page,
		PageSize:   pageSize,
		TotalCount: total,
		TotalPages: totalPages,
	})
}

// GetByID godoc
// @Summary Get county by ID
// @Description Get a single county by its ID
// @Tags counties
// @Accept json
// @Produce json
// @Param id path string true "County ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/v1/counties/{id} [get]
func (h *CountyHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	county, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "county not found" {
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
		Data:    county,
	})
}

// GetByCode godoc
// @Summary Get county by code
// @Description Get a single county by its code number
// @Tags counties
// @Accept json
// @Produce json
// @Param code path int true "County Code"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/v1/counties/code/{code} [get]
func (h *CountyHandler) GetByCode(c *gin.Context) {
	code, err := strconv.Atoi(c.Param("code"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid county code",
		})
		return
	}

	county, err := h.service.GetByCode(c.Request.Context(), code)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "county not found" {
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
		Data:    county,
	})
}

// GetByName godoc
// @Summary Get county by name
// @Description Get a single county by its name (case-insensitive)
// @Tags counties
// @Accept json
// @Produce json
// @Param name query string true "County Name"
// @Success 200 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/v1/counties/search [get]
func (h *CountyHandler) GetByName(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "County name is required",
		})
		return
	}

	county, err := h.service.GetByName(c.Request.Context(), name)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "county not found" {
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
		Data:    county,
	})
}

// Create godoc
// @Summary Create a new county
// @Description Create a new county record
// @Tags counties
// @Accept json
// @Produce json
// @Param county body models.County true "County object"
// @Success 201 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Router /api/v1/counties [post]
func (h *CountyHandler) Create(c *gin.Context) {
	var county models.County
	if err := c.ShouldBindJSON(&county); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if err := h.service.Create(c.Request.Context(), &county); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to create county",
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "County created successfully",
		Data:    county,
	})
}

// Update godoc
// @Summary Update a county
// @Description Update an existing county record
// @Tags counties
// @Accept json
// @Produce json
// @Param id path string true "County ID"
// @Param county body models.County true "County object"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/v1/counties/{id} [put]
func (h *CountyHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var county models.County
	if err := c.ShouldBindJSON(&county); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if err := h.service.Update(c.Request.Context(), id, &county); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "county not found" {
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
		Message: "County updated successfully",
	})
}

// Delete godoc
// @Summary Delete a county
// @Description Delete a county record
// @Tags counties
// @Accept json
// @Produce json
// @Param id path string true "County ID"
// @Success 200 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/v1/counties/{id} [delete]
func (h *CountyHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "county not found" {
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
		Message: "County deleted successfully",
	})
}
