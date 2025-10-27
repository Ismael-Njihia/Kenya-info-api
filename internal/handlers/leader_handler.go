package handlers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/Ismael-Njihia/Kenya-info-api/internal/models"
	"github.com/Ismael-Njihia/Kenya-info-api/internal/services"
	"github.com/gin-gonic/gin"
)

type LeaderHandler struct {
	service *services.LeaderService
}

func NewLeaderHandler(service *services.LeaderService) *LeaderHandler {
	return &LeaderHandler{service: service}
}

// GetAll godoc
// @Summary Get all leaders
// @Description Get a list of all leaders with pagination
// @Tags leaders
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(50)
// @Success 200 {object} models.PaginatedResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/v1/leaders [get]
func (h *LeaderHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 50
	}

	leaders, total, err := h.service.GetAll(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to fetch leaders",
		})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success:    true,
		Data:       leaders,
		Page:       page,
		PageSize:   pageSize,
		TotalCount: total,
		TotalPages: totalPages,
	})
}

// GetByID godoc
// @Summary Get leader by ID
// @Description Get a single leader by their ID
// @Tags leaders
// @Accept json
// @Produce json
// @Param id path string true "Leader ID"
// @Success 200 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/v1/leaders/{id} [get]
func (h *LeaderHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	leader, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "leader not found" {
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
		Data:    leader,
	})
}

// GetByPosition godoc
// @Summary Get leaders by position
// @Description Get leaders filtered by position (e.g., Governor, Senator, MCA)
// @Tags leaders
// @Accept json
// @Produce json
// @Param position query string true "Position"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(50)
// @Success 200 {object} models.PaginatedResponse
// @Failure 400 {object} models.APIResponse
// @Router /api/v1/leaders/position [get]
func (h *LeaderHandler) GetByPosition(c *gin.Context) {
	position := c.Query("position")
	if position == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Position is required",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 50
	}

	leaders, total, err := h.service.GetByPosition(c.Request.Context(), position, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to fetch leaders",
		})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Success:    true,
		Data:       leaders,
		Page:       page,
		PageSize:   pageSize,
		TotalCount: total,
		TotalPages: totalPages,
	})
}

// GetByCountyID godoc
// @Summary Get leaders by county ID
// @Description Get all leaders belonging to a specific county
// @Tags leaders
// @Accept json
// @Produce json
// @Param county_id path string true "County ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(50)
// @Success 200 {object} models.PaginatedResponse
// @Failure 400 {object} models.APIResponse
// @Router /api/v1/counties/{county_id}/leaders [get]
func (h *LeaderHandler) GetByCountyID(c *gin.Context) {
	countyID := c.Param("county_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 50
	}

	leaders, total, err := h.service.GetByCountyID(c.Request.Context(), countyID, page, pageSize)
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
		Data:       leaders,
		Page:       page,
		PageSize:   pageSize,
		TotalCount: total,
		TotalPages: totalPages,
	})
}

// Create godoc
// @Summary Create a new leader
// @Description Create a new leader record
// @Tags leaders
// @Accept json
// @Produce json
// @Param leader body models.Leader true "Leader object"
// @Success 201 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Router /api/v1/leaders [post]
func (h *LeaderHandler) Create(c *gin.Context) {
	var leader models.Leader
	if err := c.ShouldBindJSON(&leader); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if err := h.service.Create(c.Request.Context(), &leader); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to create leader",
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Leader created successfully",
		Data:    leader,
	})
}

// Update godoc
// @Summary Update a leader
// @Description Update an existing leader record
// @Tags leaders
// @Accept json
// @Produce json
// @Param id path string true "Leader ID"
// @Param leader body models.Leader true "Leader object"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/v1/leaders/{id} [put]
func (h *LeaderHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var leader models.Leader
	if err := c.ShouldBindJSON(&leader); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if err := h.service.Update(c.Request.Context(), id, &leader); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "leader not found" {
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
		Message: "Leader updated successfully",
	})
}

// Delete godoc
// @Summary Delete a leader
// @Description Delete a leader record
// @Tags leaders
// @Accept json
// @Produce json
// @Param id path string true "Leader ID"
// @Success 200 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/v1/leaders/{id} [delete]
func (h *LeaderHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "leader not found" {
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
		Message: "Leader deleted successfully",
	})
}
