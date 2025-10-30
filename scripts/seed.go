package main

import (
	"context"
	"log"
	"time"

	"github.com/Ismael-Njihia/Kenya-info-api/internal/config"
	"github.com/Ismael-Njihia/Kenya-info-api/internal/database"
	"github.com/Ismael-Njihia/Kenya-info-api/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// -------- MAIN --------

func main() {
	log.Println("🚀 Starting database seeding...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Failed to load configuration: %v", err)
	}

	// Connect to DB
	db, err := database.Connect(&cfg.Database)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer db.Client.Disconnect(context.Background())

	ctx := context.Background()

	// Clear previous data
	if err := clearCollections(ctx, db); err != nil {
		log.Printf("⚠️ Warning clearing collections: %v", err)
	}

	// Seed data
	countyIDs, err := seedCounties(ctx, db)
	if err != nil {
		log.Fatalf("❌ Failed to seed counties: %v", err)
	}

	constituencyIDs, err := seedConstituencies(ctx, db, countyIDs)
	if err != nil {
		log.Fatalf("❌ Failed to seed constituencies: %v", err)
	}

	wardIDs, err := seedWards(ctx, db, countyIDs, constituencyIDs)
	if err != nil {
		log.Fatalf("❌ Failed to seed wards: %v", err)
	}

	if err := seedLeaders(ctx, db, countyIDs, constituencyIDs, wardIDs); err != nil {
		log.Fatalf("❌ Failed to seed leaders: %v", err)
	}

	log.Println("✅ Database seeding completed successfully!")
}

// -------- CLEAR --------

func clearCollections(ctx context.Context, db *database.Database) error {
	collections := []string{"counties", "constituencies", "wards", "leaders"}
	for _, name := range collections {
		if err := db.GetCollection(name).Drop(ctx); err != nil {
			return err
		}
		log.Printf("🧹 Cleared %s collection", name)
	}
	return nil
}

// -------- SEED COUNTIES --------

func seedCounties(ctx context.Context, db *database.Database) (map[string]primitive.ObjectID, error) {
	now := time.Now()
	counties := []models.County{
		{ID: primitive.NewObjectID(), Code: 47, Name: "Nairobi", Capital: "Nairobi", Population: 4397073, Area: 696, Governor: "Johnson Sakaja", CreatedAt: now, UpdatedAt: now},
		{ID: primitive.NewObjectID(), Code: 22, Name: "Kisumu", Capital: "Kisumu City", Population: 1155574, Area: 567, Governor: "Anyang' Nyong’o", CreatedAt: now, UpdatedAt: now},
		{ID: primitive.NewObjectID(), Code: 1, Name: "Mombasa", Capital: "Mombasa", Population: 1208333, Area: 219, Governor: "Abdulswamad Nassir", CreatedAt: now, UpdatedAt: now},
	}

	coll := db.GetCollection("counties")
	docs := make([]interface{}, len(counties))
	countyIDs := make(map[string]primitive.ObjectID)

	for i, county := range counties {
		docs[i] = county
		countyIDs[county.Name] = county.ID
	}

	if _, err := coll.InsertMany(ctx, docs); err != nil {
		return nil, err
	}

	log.Printf("🌍 Seeded %d counties", len(counties))
	return countyIDs, nil
}

// -------- SEED CONSTITUENCIES --------

func seedConstituencies(ctx context.Context, db *database.Database, countyIDs map[string]primitive.ObjectID) (map[string]primitive.ObjectID, error) {
	now := time.Now()
	constituencies := []models.Constituency{
		{ID: primitive.NewObjectID(), Name: "Westlands", CountyID: countyIDs["Nairobi"], CountyName: "Nairobi", MP: "Timothy Wanyonyi", Population: 308854, Area: 72.4, CreatedAt: now, UpdatedAt: now},
		{ID: primitive.NewObjectID(), Name: "Kisumu Central", CountyID: countyIDs["Kisumu"], CountyName: "Kisumu", MP: "Joshua Oron", Population: 221417, Area: 50.1, CreatedAt: now, UpdatedAt: now},
		{ID: primitive.NewObjectID(), Name: "Mvita", CountyID: countyIDs["Mombasa"], CountyName: "Mombasa", MP: "Mohammed Machele", Population: 154063, Area: 14.8, CreatedAt: now, UpdatedAt: now},
	}

	coll := db.GetCollection("constituencies")
	docs := make([]interface{}, len(constituencies))
	constituencyIDs := make(map[string]primitive.ObjectID)

	for i, c := range constituencies {
		docs[i] = c
		constituencyIDs[c.Name] = c.ID
	}

	if _, err := coll.InsertMany(ctx, docs); err != nil {
		return nil, err
	}

	log.Printf("🏛️ Seeded %d constituencies", len(constituencies))
	return constituencyIDs, nil
}

// -------- SEED WARDS --------

func seedWards(ctx context.Context, db *database.Database, countyIDs map[string]primitive.ObjectID, constituencyIDs map[string]primitive.ObjectID) (map[string]primitive.ObjectID, error) {
	now := time.Now()
	wards := []models.Ward{
		{ID: primitive.NewObjectID(), Name: "Parklands/Highridge", CountyID: countyIDs["Nairobi"], CountyName: "Nairobi", ConstituencyID: constituencyIDs["Westlands"], ConstituencyName: "Westlands", MCA: "Jayendra Malde", Population: 45000, Area: 15.4, CreatedAt: now, UpdatedAt: now},
		{ID: primitive.NewObjectID(), Name: "Railways Ward", CountyID: countyIDs["Kisumu"], CountyName: "Kisumu", ConstituencyID: constituencyIDs["Kisumu Central"], ConstituencyName: "Kisumu Central", MCA: "Samuel Onyango", Population: 38000, Area: 9.3, CreatedAt: now, UpdatedAt: now},
		{ID: primitive.NewObjectID(), Name: "Tudor Ward", CountyID: countyIDs["Mombasa"], CountyName: "Mombasa", ConstituencyID: constituencyIDs["Mvita"], ConstituencyName: "Mvita", MCA: "Abdallah Mbarak", Population: 29000, Area: 6.7, CreatedAt: now, UpdatedAt: now},
	}

	coll := db.GetCollection("wards")
	docs := make([]interface{}, len(wards))
	wardIDs := make(map[string]primitive.ObjectID)

	for i, w := range wards {
		docs[i] = w
		wardIDs[w.Name] = w.ID
	}

	if _, err := coll.InsertMany(ctx, docs); err != nil {
		return nil, err
	}

	log.Printf("🏘️ Seeded %d wards", len(wards))
	return wardIDs, nil
}

// -------- SEED LEADERS --------

func seedLeaders(ctx context.Context, db *database.Database, countyIDs map[string]primitive.ObjectID, constituencyIDs map[string]primitive.ObjectID, wardIDs map[string]primitive.ObjectID) error {
	now := time.Now()
	leaders := []models.Leader{
		{
			ID:         primitive.NewObjectID(),
			Name:       "Johnson Sakaja",
			Position:   "Governor",
			CountyID:   countyIDs["Nairobi"],
			CountyName: "Nairobi",
			Party:      "UDA",
			Email:      "governor@nairobi.go.ke",
			Phone:      "+254700000001",
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ID:               primitive.NewObjectID(),
			Name:             "Timothy Wanyonyi",
			Position:         "MP",
			CountyID:         countyIDs["Nairobi"],
			CountyName:       "Nairobi",
			ConstituencyID:   constituencyIDs["Westlands"],
			ConstituencyName: "Westlands",
			Party:            "ODM",
			Email:            "wanyonyi@parliament.go.ke",
			Phone:            "+254700000002",
			CreatedAt:        now,
			UpdatedAt:        now,
		},
		{
			ID:               primitive.NewObjectID(),
			Name:             "Jayendra Malde",
			Position:         "MCA",
			CountyID:         countyIDs["Nairobi"],
			CountyName:       "Nairobi",
			ConstituencyID:   constituencyIDs["Westlands"],
			ConstituencyName: "Westlands",
			WardID:           wardIDs["Parklands/Highridge"],
			WardName:         "Parklands/Highridge",
			Party:            "Jubilee",
			Email:            "jayendra@nairobi.go.ke",
			Phone:            "+254700000003",
			CreatedAt:        now,
			UpdatedAt:        now,
		},
	}

	coll := db.GetCollection("leaders")
	docs := make([]interface{}, len(leaders))
	for i, l := range leaders {
		docs[i] = l
	}

	if _, err := coll.InsertMany(ctx, docs); err != nil {
		return err
	}

	log.Printf("👔 Seeded %d leaders", len(leaders))
	return nil
}
