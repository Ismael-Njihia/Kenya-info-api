package services

import (
	"context"
	"fmt"
	"time"

	"github.com/Ismael-Njihia/Kenya-info-api/internal/database"
	"github.com/Ismael-Njihia/Kenya-info-api/internal/models"
	"github.com/Ismael-Njihia/Kenya-info-api/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type LeaderService struct {
	collection *mongo.Collection
}

func NewLeaderService(db *database.Database) *LeaderService {
	return &LeaderService{
		collection: db.GetCollection("leaders"),
	}
}

// GetAll retrieves all leaders with pagination
func (s *LeaderService) GetAll(ctx context.Context, page, pageSize int) ([]models.Leader, int64, error) {
	skip := (page - 1) * pageSize
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(pageSize)).SetSort(bson.D{{Key: "name", Value: 1}})

	cursor, err := s.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		logger.Error("Failed to fetch leaders", zap.Error(err))
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var leaders []models.Leader
	if err := cursor.All(ctx, &leaders); err != nil {
		logger.Error("Failed to decode leaders", zap.Error(err))
		return nil, 0, err
	}

	total, err := s.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		logger.Error("Failed to count leaders", zap.Error(err))
		return nil, 0, err
	}

	return leaders, total, nil
}

// GetByID retrieves a leader by ID
func (s *LeaderService) GetByID(ctx context.Context, id string) (*models.Leader, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid leader ID: %w", err)
	}

	var leader models.Leader
	err = s.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&leader)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("leader not found")
		}
		logger.Error("Failed to fetch leader", zap.Error(err), zap.String("id", id))
		return nil, err
	}

	return &leader, nil
}

// GetByPosition retrieves leaders by position
func (s *LeaderService) GetByPosition(ctx context.Context, position string, page, pageSize int) ([]models.Leader, int64, error) {
	skip := (page - 1) * pageSize
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(pageSize)).SetSort(bson.D{{Key: "name", Value: 1}})

	filter := bson.M{"position": bson.M{"$regex": primitive.Regex{Pattern: position, Options: "i"}}}
	cursor, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		logger.Error("Failed to fetch leaders by position", zap.Error(err), zap.String("position", position))
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var leaders []models.Leader
	if err := cursor.All(ctx, &leaders); err != nil {
		logger.Error("Failed to decode leaders", zap.Error(err))
		return nil, 0, err
	}

	total, err := s.collection.CountDocuments(ctx, filter)
	if err != nil {
		logger.Error("Failed to count leaders", zap.Error(err))
		return nil, 0, err
	}

	return leaders, total, nil
}

// GetByCountyID retrieves leaders by county ID
func (s *LeaderService) GetByCountyID(ctx context.Context, countyID string, page, pageSize int) ([]models.Leader, int64, error) {
	objectID, err := primitive.ObjectIDFromHex(countyID)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid county ID: %w", err)
	}

	skip := (page - 1) * pageSize
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(pageSize)).SetSort(bson.D{{Key: "name", Value: 1}})

	filter := bson.M{"county_id": objectID}
	cursor, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		logger.Error("Failed to fetch leaders by county", zap.Error(err), zap.String("county_id", countyID))
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var leaders []models.Leader
	if err := cursor.All(ctx, &leaders); err != nil {
		logger.Error("Failed to decode leaders", zap.Error(err))
		return nil, 0, err
	}

	total, err := s.collection.CountDocuments(ctx, filter)
	if err != nil {
		logger.Error("Failed to count leaders", zap.Error(err))
		return nil, 0, err
	}

	return leaders, total, nil
}

// Create creates a new leader
func (s *LeaderService) Create(ctx context.Context, leader *models.Leader) error {
	leader.CreatedAt = time.Now()
	leader.UpdatedAt = time.Now()

	result, err := s.collection.InsertOne(ctx, leader)
	if err != nil {
		logger.Error("Failed to create leader", zap.Error(err))
		return err
	}

	leader.ID = result.InsertedID.(primitive.ObjectID)
	logger.Info("Leader created", zap.String("id", leader.ID.Hex()), zap.String("name", leader.Name))
	return nil
}

// Update updates an existing leader
func (s *LeaderService) Update(ctx context.Context, id string, leader *models.Leader) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid leader ID: %w", err)
	}

	leader.UpdatedAt = time.Now()
	update := bson.M{"$set": leader}

	result, err := s.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		logger.Error("Failed to update leader", zap.Error(err), zap.String("id", id))
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("leader not found")
	}

	logger.Info("Leader updated", zap.String("id", id))
	return nil
}

// Delete deletes a leader
func (s *LeaderService) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid leader ID: %w", err)
	}

	result, err := s.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		logger.Error("Failed to delete leader", zap.Error(err), zap.String("id", id))
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("leader not found")
	}

	logger.Info("Leader deleted", zap.String("id", id))
	return nil
}
