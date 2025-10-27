package repository

import (
	"context"
	"time"

	"github.com/Ismael-Njihia/Kenya-info-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CountyRepository handles county data operations
type CountyRepository struct {
	collection *mongo.Collection
}

// NewCountyRepository creates a new county repository
func NewCountyRepository(db *mongo.Database) *CountyRepository {
	return &CountyRepository{
		collection: db.Collection("counties"),
	}
}

// GetAll retrieves all counties
func (r *CountyRepository) GetAll(ctx context.Context) ([]models.County, error) {
	cursor, err := r.collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "code", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var counties []models.County
	if err = cursor.All(ctx, &counties); err != nil {
		return nil, err
	}
	return counties, nil
}

// GetByID retrieves a county by ID
func (r *CountyRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.County, error) {
	var county models.County
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&county)
	if err != nil {
		return nil, err
	}
	return &county, nil
}

// GetByCode retrieves a county by code
func (r *CountyRepository) GetByCode(ctx context.Context, code int) (*models.County, error) {
	var county models.County
	err := r.collection.FindOne(ctx, bson.M{"code": code}).Decode(&county)
	if err != nil {
		return nil, err
	}
	return &county, nil
}

// Create creates a new county
func (r *CountyRepository) Create(ctx context.Context, county *models.County) error {
	county.CreatedAt = time.Now()
	county.UpdatedAt = time.Now()
	result, err := r.collection.InsertOne(ctx, county)
	if err != nil {
		return err
	}
	county.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// Update updates an existing county
func (r *CountyRepository) Update(ctx context.Context, id primitive.ObjectID, county *models.County) error {
	county.UpdatedAt = time.Now()
	update := bson.M{
		"$set": bson.M{
			"code":       county.Code,
			"name":       county.Name,
			"capital":    county.Capital,
			"governor":   county.Governor,
			"population": county.Population,
			"area":       county.Area,
			"updated_at": county.UpdatedAt,
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// Delete deletes a county
func (r *CountyRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// Count returns the total number of counties
func (r *CountyRepository) Count(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}
