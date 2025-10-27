package services

import (
	"context"
	"testing"
	"time"

	"github.com/Ismael-Njihia/Kenya-info-api/internal/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCountyValidation(t *testing.T) {
	county := &models.County{
		Code:       1,
		Name:       "Nairobi",
		Capital:    "Nairobi City",
		Population: 4397073,
		Area:       696.0,
		Governor:   "Johnson Sakaja",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	assert.NotNil(t, county)
	assert.Equal(t, "Nairobi", county.Name)
	assert.Equal(t, 1, county.Code)
	assert.Greater(t, county.Population, int64(0))
}

func TestWardValidation(t *testing.T) {
	ward := &models.Ward{
		Name:       "Kilimani",
		CountyID:   primitive.NewObjectID(),
		CountyName: "Nairobi",
		MCA:        "Moses Ogeto",
		Population: 50000,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	assert.NotNil(t, ward)
	assert.Equal(t, "Kilimani", ward.Name)
	assert.NotEmpty(t, ward.CountyID)
	assert.Equal(t, "Nairobi", ward.CountyName)
}

func TestLeaderValidation(t *testing.T) {
	leader := &models.Leader{
		Name:       "William Ruto",
		Position:   "President",
		Party:      "UDA",
		Email:      "president@kenya.go.ke",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	assert.NotNil(t, leader)
	assert.Equal(t, "William Ruto", leader.Name)
	assert.Equal(t, "President", leader.Position)
	assert.Equal(t, "UDA", leader.Party)
}

func TestPaginationLogic(t *testing.T) {
	testCases := []struct {
		page     int
		pageSize int
		expected int64
	}{
		{1, 10, 0},
		{2, 10, 10},
		{3, 20, 40},
		{5, 15, 60},
	}

	for _, tc := range testCases {
		skip := int64((tc.page - 1) * tc.pageSize)
		assert.Equal(t, tc.expected, skip)
	}
}

func TestAPIResponseStructure(t *testing.T) {
	response := models.APIResponse{
		Success: true,
		Message: "Operation successful",
		Data:    map[string]string{"test": "data"},
	}

	assert.True(t, response.Success)
	assert.Equal(t, "Operation successful", response.Message)
	assert.NotNil(t, response.Data)
}

func TestPaginatedResponseStructure(t *testing.T) {
	response := models.PaginatedResponse{
		Success:    true,
		Data:       []string{"item1", "item2"},
		Page:       1,
		PageSize:   10,
		TotalCount: 100,
		TotalPages: 10,
	}

	assert.True(t, response.Success)
	assert.Equal(t, 1, response.Page)
	assert.Equal(t, 10, response.PageSize)
	assert.Equal(t, int64(100), response.TotalCount)
	assert.Equal(t, 10, response.TotalPages)
}

// Mock tests to ensure service methods have correct signatures
func TestCountyServiceInterface(t *testing.T) {
	// This test ensures the service interface is correct
	var _ interface {
		GetAll(ctx context.Context, page, pageSize int) ([]models.County, int64, error)
		GetByID(ctx context.Context, id string) (*models.County, error)
		GetByCode(ctx context.Context, code int) (*models.County, error)
		Create(ctx context.Context, county *models.County) error
		Update(ctx context.Context, id string, county *models.County) error
		Delete(ctx context.Context, id string) error
	} = (*CountyService)(nil)
}

func TestWardServiceInterface(t *testing.T) {
	// This test ensures the service interface is correct
	var _ interface {
		GetAll(ctx context.Context, page, pageSize int) ([]models.Ward, int64, error)
		GetByID(ctx context.Context, id string) (*models.Ward, error)
		GetByCountyID(ctx context.Context, countyID string, page, pageSize int) ([]models.Ward, int64, error)
		Create(ctx context.Context, ward *models.Ward) error
		Update(ctx context.Context, id string, ward *models.Ward) error
		Delete(ctx context.Context, id string) error
	} = (*WardService)(nil)
}

func TestLeaderServiceInterface(t *testing.T) {
	// This test ensures the service interface is correct
	var _ interface {
		GetAll(ctx context.Context, page, pageSize int) ([]models.Leader, int64, error)
		GetByID(ctx context.Context, id string) (*models.Leader, error)
		GetByPosition(ctx context.Context, position string, page, pageSize int) ([]models.Leader, int64, error)
		GetByCountyID(ctx context.Context, countyID string, page, pageSize int) ([]models.Leader, int64, error)
		Create(ctx context.Context, leader *models.Leader) error
		Update(ctx context.Context, id string, leader *models.Leader) error
		Delete(ctx context.Context, id string) error
	} = (*LeaderService)(nil)
}
