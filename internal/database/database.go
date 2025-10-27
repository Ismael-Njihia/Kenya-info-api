package database

import (
	"context"
	"fmt"
	"time"

	"github.com/Ismael-Njihia/Kenya-info-api/internal/config"
	"github.com/Ismael-Njihia/Kenya-info-api/pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type Database struct {
	Client *mongo.Client
	DB     *mongo.Database
}

// Connect establishes a connection to MongoDB
func Connect(cfg *config.DatabaseConfig) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
	defer cancel()

	logger.Info("Connecting to MongoDB", zap.String("database", cfg.Database))

	clientOptions := options.Client().ApplyURI(cfg.MongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	logger.Info("Successfully connected to MongoDB")

	return &Database{
		Client: client,
		DB:     client.Database(cfg.Database),
	}, nil
}

// Disconnect closes the database connection
func (db *Database) Disconnect(ctx context.Context) error {
	logger.Info("Disconnecting from MongoDB")
	return db.Client.Disconnect(ctx)
}

// GetCollection returns a MongoDB collection
func (db *Database) GetCollection(name string) *mongo.Collection {
	return db.DB.Collection(name)
}
