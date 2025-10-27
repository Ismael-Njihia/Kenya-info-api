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

// WardRepository handles ward data operations
type WardRepository struct {
	collection *mongo.Collection
}

// NewWardRepository creates a new ward repository
func NewWardRepository(db *mongo.Database) *WardRepository {
	return &WardRepository{
		collection: db.Collection("wards"),
	}
}

// GetAll retrieves all wards
func (r *WardRepository) GetAll(ctx context.Context) ([]models.Ward, error) {
	cursor, err := r.collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "name", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var wards []models.Ward
	if err = cursor.All(ctx, &wards); err != nil {
		return nil, err
	}
	return wards, nil
}

// GetByID retrieves a ward by ID
func (r *WardRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Ward, error) {
	var ward models.Ward
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&ward)
	if err != nil {
		return nil, err
	}
	return &ward, nil
}

// GetByCountyID retrieves wards by county ID
func (r *WardRepository) GetByCountyID(ctx context.Context, countyID primitive.ObjectID) ([]models.Ward, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"county_id": countyID}, options.Find().SetSort(bson.D{{Key: "name", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var wards []models.Ward
	if err = cursor.All(ctx, &wards); err != nil {
		return nil, err
	}
	return wards, nil
}

// Create creates a new ward
func (r *WardRepository) Create(ctx context.Context, ward *models.Ward) error {
	ward.CreatedAt = time.Now()
	ward.UpdatedAt = time.Now()
	result, err := r.collection.InsertOne(ctx, ward)
	if err != nil {
		return err
	}
	ward.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// Update updates an existing ward
func (r *WardRepository) Update(ctx context.Context, id primitive.ObjectID, ward *models.Ward) error {
	ward.UpdatedAt = time.Now()
	update := bson.M{
		"$set": bson.M{
			"name":              ward.Name,
			"county_id":         ward.CountyID,
			"county_name":       ward.CountyName,
			"constituency_name": ward.ConstituencyName,
			"population":        ward.Population,
			"updated_at":        ward.UpdatedAt,
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// Delete deletes a ward
func (r *WardRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// Count returns the total number of wards
func (r *WardRepository) Count(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}
