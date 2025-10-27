package repository

import (
	"context"
	"testing"
	"time"

	"github.com/Ismael-Njihia/Kenya-info-api/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestDB(t *testing.T) (*mongo.Database, func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	require.NoError(t, err)

	db := client.Database("kenya_info_test_" + primitive.NewObjectID().Hex())

	cleanup := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		db.Drop(ctx)
		client.Disconnect(ctx)
	}

	return db, cleanup
}

func TestCountyRepository_Create(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewCountyRepository(db)
	county := &models.County{
		Code:       1,
		Name:       "Mombasa",
		Capital:    "Mombasa City",
		Governor:   "Test Governor",
		Population: 1000000,
		Area:       200.0,
	}

	err := repo.Create(context.Background(), county)
	assert.NoError(t, err)
	assert.NotEqual(t, primitive.NilObjectID, county.ID)
}

func TestCountyRepository_GetByCode(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewCountyRepository(db)
	county := &models.County{
		Code:       1,
		Name:       "Mombasa",
		Capital:    "Mombasa City",
		Governor:   "Test Governor",
		Population: 1000000,
		Area:       200.0,
	}

	err := repo.Create(context.Background(), county)
	require.NoError(t, err)

	found, err := repo.GetByCode(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, "Mombasa", found.Name)
}

func TestCountyRepository_GetAll(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewCountyRepository(db)
	
	counties := []models.County{
		{Code: 1, Name: "Mombasa", Capital: "Mombasa City", Governor: "Gov1", Population: 1000000, Area: 200.0},
		{Code: 2, Name: "Kwale", Capital: "Kwale", Governor: "Gov2", Population: 800000, Area: 8000.0},
	}

	for i := range counties {
		err := repo.Create(context.Background(), &counties[i])
		require.NoError(t, err)
	}

	all, err := repo.GetAll(context.Background())
	assert.NoError(t, err)
	assert.Len(t, all, 2)
}

func TestCountyRepository_Update(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewCountyRepository(db)
	county := &models.County{
		Code:       1,
		Name:       "Mombasa",
		Capital:    "Mombasa City",
		Governor:   "Old Governor",
		Population: 1000000,
		Area:       200.0,
	}

	err := repo.Create(context.Background(), county)
	require.NoError(t, err)

	county.Governor = "New Governor"
	err = repo.Update(context.Background(), county.ID, county)
	assert.NoError(t, err)

	updated, err := repo.GetByID(context.Background(), county.ID)
	assert.NoError(t, err)
	assert.Equal(t, "New Governor", updated.Governor)
}

func TestCountyRepository_Delete(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewCountyRepository(db)
	county := &models.County{
		Code:       1,
		Name:       "Mombasa",
		Capital:    "Mombasa City",
		Governor:   "Test Governor",
		Population: 1000000,
		Area:       200.0,
	}

	err := repo.Create(context.Background(), county)
	require.NoError(t, err)

	err = repo.Delete(context.Background(), county.ID)
	assert.NoError(t, err)

	_, err = repo.GetByID(context.Background(), county.ID)
	assert.Error(t, err)
}
