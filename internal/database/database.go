package database

import (
	"context"
	"fmt"
	"time"

	"github.com/Ismael-Njihia/Kenya-info-api/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func Connect(cfg *config.DatabaseConfig) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.URI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the DB
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	fmt.Println("✅ Connected successfully to MongoDB Atlas")

	return &Database{
		Client:   client,
		Database: client.Database(cfg.Name),
	}, nil
}

// GetCollection returns a MongoDB collection
func (d *Database) GetCollection(name string) *mongo.Collection {
	return d.Database.Collection(name)
}

// Close cleanly disconnects the MongoDB client
func (d *Database) Close(ctx context.Context) error {
	return d.Client.Disconnect(ctx)
}
