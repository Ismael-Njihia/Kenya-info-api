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

type ConstituencyService struct {
	collection *mongo.Collection
}

func NewConstituencyService(db *database.Database) *ConstituencyService {
	return &ConstituencyService{
		collection: db.GetCollection("constituencies"),
	}
}

// GetAll retrieves all constituencies with pagination
func (s *ConstituencyService) GetAll(ctx context.Context, page, pageSize int) ([]models.Constituency, int64, error) {
	skip := (page - 1) * pageSize
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(pageSize)).SetSort(bson.D{{Key: "name", Value: 1}})

	cursor, err := s.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		logger.Error("Failed to fetch constituencies", zap.Error(err))
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var constituencies []models.Constituency
	if err := cursor.All(ctx, &constituencies); err != nil {
		logger.Error("Failed to decode constituencies", zap.Error(err))
		return nil, 0, err
	}

	total, err := s.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		logger.Error("Failed to count constituencies", zap.Error(err))
		return nil, 0, err
	}

	return constituencies, total, nil
}

// GetByID retrieves a single constituency by ID
func (s *ConstituencyService) GetByID(ctx context.Context, id string) (*models.Constituency, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid constituency ID: %w", err)
	}

	var constituency models.Constituency
	err = s.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&constituency)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("constituency not found")
		}
		logger.Error("Failed to fetch constituency", zap.Error(err), zap.String("id", id))
		return nil, err
	}

	return &constituency, nil
}

// GetByCountyID retrieves constituencies belonging to a specific county
func (s *ConstituencyService) GetByCountyID(ctx context.Context, countyID string, page, pageSize int) ([]models.Constituency, int64, error) {
	objectID, err := primitive.ObjectIDFromHex(countyID)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid county ID: %w", err)
	}

	skip := (page - 1) * pageSize
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "name", Value: 1}})

	filter := bson.M{"county_id": objectID}

	cursor, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		logger.Error("Failed to fetch constituencies by county", zap.Error(err), zap.String("county_id", countyID))
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var constituencies []models.Constituency
	if err := cursor.All(ctx, &constituencies); err != nil {
		logger.Error("Failed to decode constituencies", zap.Error(err))
		return nil, 0, err
	}

	total, err := s.collection.CountDocuments(ctx, filter)
	if err != nil {
		logger.Error("Failed to count constituencies", zap.Error(err))
		return nil, 0, err
	}

	return constituencies, total, nil
}

// Create a new constituency
func (s *ConstituencyService) Create(ctx context.Context, constituency *models.Constituency) error {
	constituency.CreatedAt = time.Now()
	constituency.UpdatedAt = time.Now()

	result, err := s.collection.InsertOne(ctx, constituency)
	if err != nil {
		logger.Error("Failed to create constituency", zap.Error(err))
		return err
	}

	constituency.ID = result.InsertedID.(primitive.ObjectID)
	logger.Info("Constituency created", zap.String("id", constituency.ID.Hex()), zap.String("name", constituency.Name))
	return nil
}

// Update an existing constituency
func (s *ConstituencyService) Update(ctx context.Context, id string, constituency *models.Constituency) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid constituency ID: %w", err)
	}

	constituency.UpdatedAt = time.Now()
	update := bson.M{"$set": constituency}

	result, err := s.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		logger.Error("Failed to update constituency", zap.Error(err))
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("constituency not found")
	}

	logger.Info("Constituency updated", zap.String("id", id))
	return nil
}

// Delete a constituency
func (s *ConstituencyService) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid constituency ID: %w", err)
	}

	result, err := s.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		logger.Error("Failed to delete constituency", zap.Error(err))
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("constituency not found")
	}

	logger.Info("Constituency deleted", zap.String("id", id))
	return nil
}
