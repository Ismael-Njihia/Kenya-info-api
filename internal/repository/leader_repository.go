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

// LeaderRepository handles leader data operations
type LeaderRepository struct {
	collection *mongo.Collection
}

// NewLeaderRepository creates a new leader repository
func NewLeaderRepository(db *mongo.Database) *LeaderRepository {
	return &LeaderRepository{
		collection: db.Collection("leaders"),
	}
}

// GetAll retrieves all leaders
func (r *LeaderRepository) GetAll(ctx context.Context) ([]models.Leader, error) {
	cursor, err := r.collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "name", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var leaders []models.Leader
	if err = cursor.All(ctx, &leaders); err != nil {
		return nil, err
	}
	return leaders, nil
}

// GetByID retrieves a leader by ID
func (r *LeaderRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Leader, error) {
	var leader models.Leader
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&leader)
	if err != nil {
		return nil, err
	}
	return &leader, nil
}

// GetByPosition retrieves leaders by position
func (r *LeaderRepository) GetByPosition(ctx context.Context, position string) ([]models.Leader, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"position": position}, options.Find().SetSort(bson.D{{Key: "name", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var leaders []models.Leader
	if err = cursor.All(ctx, &leaders); err != nil {
		return nil, err
	}
	return leaders, nil
}

// GetByCounty retrieves leaders by county
func (r *LeaderRepository) GetByCounty(ctx context.Context, county string) ([]models.Leader, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"county": county}, options.Find().SetSort(bson.D{{Key: "name", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var leaders []models.Leader
	if err = cursor.All(ctx, &leaders); err != nil {
		return nil, err
	}
	return leaders, nil
}

// Create creates a new leader
func (r *LeaderRepository) Create(ctx context.Context, leader *models.Leader) error {
	leader.CreatedAt = time.Now()
	leader.UpdatedAt = time.Now()
	result, err := r.collection.InsertOne(ctx, leader)
	if err != nil {
		return err
	}
	leader.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// Update updates an existing leader
func (r *LeaderRepository) Update(ctx context.Context, id primitive.ObjectID, leader *models.Leader) error {
	leader.UpdatedAt = time.Now()
	update := bson.M{
		"$set": bson.M{
			"name":       leader.Name,
			"position":   leader.Position,
			"county":     leader.County,
			"party":      leader.Party,
			"email":      leader.Email,
			"phone":      leader.Phone,
			"start_date": leader.StartDate,
			"updated_at": leader.UpdatedAt,
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// Delete deletes a leader
func (r *LeaderRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// Count returns the total number of leaders
func (r *LeaderRepository) Count(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}
