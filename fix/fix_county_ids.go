package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load environment variables
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("⚠️ No .env file found, using system environment variables")
	}

	uri := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("MONGODB_DATABASE")
	if uri == "" || dbName == "" {
		log.Fatal("❌ Missing MONGODB_URI or MONGODB_DATABASE in environment")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database(dbName)

	fixCollection(ctx, db.Collection("constituencies"))
	fixCollection(ctx, db.Collection("wards"))

	log.Println("✅ Migration complete: county_id fields fixed where necessary.")
}

func fixCollection(ctx context.Context, coll *mongo.Collection) {
	cursor, err := coll.Find(ctx, bson.M{"county_id": bson.M{"$type": "string"}})
	if err != nil {
		log.Fatalf("Error finding documents in %s: %v", coll.Name(), err)
	}
	defer cursor.Close(ctx)

	count := 0
	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			log.Printf("❌ Decode error: %v", err)
			continue
		}

		id, ok := doc["_id"].(primitive.ObjectID)
		if !ok {
			continue
		}

		countyIDStr, ok := doc["county_id"].(string)
		if !ok || countyIDStr == "" {
			continue
		}

		objID, err := primitive.ObjectIDFromHex(countyIDStr)
		if err != nil {
			log.Printf("❌ Invalid county_id (%s): %v", countyIDStr, err)
			continue
		}

		_, err = coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"county_id": objID}})
		if err != nil {
			log.Printf("❌ Failed to update %s _id=%v: %v", coll.Name(), id, err)
			continue
		}

		count++
	}

	log.Printf("✅ Updated %d documents in %s", count, coll.Name())
}
