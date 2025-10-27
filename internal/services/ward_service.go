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

type WardService struct {
	collection *mongo.Collection
}

func NewWardService(db *database.Database) *WardService {
	return &WardService{
		collection: db.GetCollection("wards"),
	}
}

// GetAll retrieves all wards with pagination
func (s *WardService) GetAll(ctx context.Context, page, pageSize int) ([]models.Ward, int64, error) {
	skip := (page - 1) * pageSize
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(pageSize)).SetSort(bson.D{{Key: "name", Value: 1}})

	cursor, err := s.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		logger.Error("Failed to fetch wards", zap.Error(err))
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var wards []models.Ward
	if err := cursor.All(ctx, &wards); err != nil {
		logger.Error("Failed to decode wards", zap.Error(err))
		return nil, 0, err
	}

	total, err := s.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		logger.Error("Failed to count wards", zap.Error(err))
		return nil, 0, err
	}

	return wards, total, nil
}

// GetByID retrieves a ward by ID
func (s *WardService) GetByID(ctx context.Context, id string) (*models.Ward, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ward ID: %w", err)
	}

	var ward models.Ward
	err = s.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&ward)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("ward not found")
		}
		logger.Error("Failed to fetch ward", zap.Error(err), zap.String("id", id))
		return nil, err
	}

	return &ward, nil
}

// GetByCountyID retrieves wards by county ID
func (s *WardService) GetByCountyID(ctx context.Context, countyID string, page, pageSize int) ([]models.Ward, int64, error) {
	objectID, err := primitive.ObjectIDFromHex(countyID)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid county ID: %w", err)
	}

	skip := (page - 1) * pageSize
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(pageSize)).SetSort(bson.D{{Key: "name", Value: 1}})

	filter := bson.M{"county_id": objectID}
	cursor, err := s.collection.Find(ctx, filter, opts)
	if err != nil {
		logger.Error("Failed to fetch wards by county", zap.Error(err), zap.String("county_id", countyID))
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var wards []models.Ward
	if err := cursor.All(ctx, &wards); err != nil {
		logger.Error("Failed to decode wards", zap.Error(err))
		return nil, 0, err
	}

	total, err := s.collection.CountDocuments(ctx, filter)
	if err != nil {
		logger.Error("Failed to count wards", zap.Error(err))
		return nil, 0, err
	}

	return wards, total, nil
}

// Create creates a new ward
func (s *WardService) Create(ctx context.Context, ward *models.Ward) error {
	ward.CreatedAt = time.Now()
	ward.UpdatedAt = time.Now()

	result, err := s.collection.InsertOne(ctx, ward)
	if err != nil {
		logger.Error("Failed to create ward", zap.Error(err))
		return err
	}

	ward.ID = result.InsertedID.(primitive.ObjectID)
	logger.Info("Ward created", zap.String("id", ward.ID.Hex()), zap.String("name", ward.Name))
	return nil
}

// Update updates an existing ward
func (s *WardService) Update(ctx context.Context, id string, ward *models.Ward) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ward ID: %w", err)
	}

	ward.UpdatedAt = time.Now()
	update := bson.M{"$set": ward}

	result, err := s.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		logger.Error("Failed to update ward", zap.Error(err), zap.String("id", id))
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("ward not found")
	}

	logger.Info("Ward updated", zap.String("id", id))
	return nil
}

// Delete deletes a ward
func (s *WardService) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ward ID: %w", err)
	}

	result, err := s.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		logger.Error("Failed to delete ward", zap.Error(err), zap.String("id", id))
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("ward not found")
	}

	logger.Info("Ward deleted", zap.String("id", id))
	return nil
}
