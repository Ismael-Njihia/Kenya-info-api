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

type CountyService struct {
	collection *mongo.Collection
}

func NewCountyService(db *database.Database) *CountyService {
	return &CountyService{
		collection: db.GetCollection("counties"),
	}
}

// GetAll retrieves all counties with pagination
func (s *CountyService) GetAll(ctx context.Context, page, pageSize int) ([]models.County, int64, error) {
	skip := (page - 1) * pageSize
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(pageSize)).SetSort(bson.D{{Key: "code", Value: 1}})

	cursor, err := s.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		logger.Error("Failed to fetch counties", zap.Error(err))
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var counties []models.County
	if err := cursor.All(ctx, &counties); err != nil {
		logger.Error("Failed to decode counties", zap.Error(err))
		return nil, 0, err
	}

	total, err := s.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		logger.Error("Failed to count counties", zap.Error(err))
		return nil, 0, err
	}

	return counties, total, nil
}

// GetByID retrieves a county by ID
func (s *CountyService) GetByID(ctx context.Context, id string) (*models.County, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid county ID: %w", err)
	}

	var county models.County
	err = s.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&county)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("county not found")
		}
		logger.Error("Failed to fetch county", zap.Error(err), zap.String("id", id))
		return nil, err
	}

	return &county, nil
}

// GetByCode retrieves a county by code
func (s *CountyService) GetByCode(ctx context.Context, code int) (*models.County, error) {
	var county models.County
	err := s.collection.FindOne(ctx, bson.M{"code": code}).Decode(&county)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("county not found")
		}
		logger.Error("Failed to fetch county by code", zap.Error(err), zap.Int("code", code))
		return nil, err
	}

	return &county, nil
}

// GetByName retrieves a county by name
func (s *CountyService) GetByName(ctx context.Context, name string) (*models.County, error) {
	var county models.County
	filter := bson.M{"name": bson.M{"$regex": primitive.Regex{Pattern: name, Options: "i"}}}
	err := s.collection.FindOne(ctx, filter).Decode(&county)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("county not found")
		}
		logger.Error("Failed to fetch county by name", zap.Error(err), zap.String("name", name))
		return nil, err
	}

	return &county, nil
}

// Create creates a new county
func (s *CountyService) Create(ctx context.Context, county *models.County) error {
	county.CreatedAt = time.Now()
	county.UpdatedAt = time.Now()

	result, err := s.collection.InsertOne(ctx, county)
	if err != nil {
		logger.Error("Failed to create county", zap.Error(err))
		return err
	}

	county.ID = result.InsertedID.(primitive.ObjectID)
	logger.Info("County created", zap.String("id", county.ID.Hex()), zap.String("name", county.Name))
	return nil
}

// Update updates an existing county
func (s *CountyService) Update(ctx context.Context, id string, county *models.County) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid county ID: %w", err)
	}

	county.UpdatedAt = time.Now()
	update := bson.M{"$set": county}

	result, err := s.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		logger.Error("Failed to update county", zap.Error(err), zap.String("id", id))
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("county not found")
	}

	logger.Info("County updated", zap.String("id", id))
	return nil
}

// Delete deletes a county
func (s *CountyService) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid county ID: %w", err)
	}

	result, err := s.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		logger.Error("Failed to delete county", zap.Error(err), zap.String("id", id))
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("county not found")
	}

	logger.Info("County deleted", zap.String("id", id))
	return nil
}
