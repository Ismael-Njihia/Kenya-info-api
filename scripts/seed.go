package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Ismael-Njihia/Kenya-info-api/config"
	"github.com/Ismael-Njihia/Kenya-info-api/internal/models"
	"github.com/Ismael-Njihia/Kenya-info-api/internal/repository"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Database.Timeout)*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.Database.URI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database(cfg.Database.Database)
	log.Println("Connected to MongoDB successfully")

	// Initialize repositories
	countyRepo := repository.NewCountyRepository(db)
	wardRepo := repository.NewWardRepository(db)
	leaderRepo := repository.NewLeaderRepository(db)

	// Seed counties (sample of Kenya's 47 counties)
	counties := []models.County{
		{Code: 1, Name: "Mombasa", Capital: "Mombasa City", Governor: "Abdulswamad Nassir", Population: 1208333, Area: 229.7},
		{Code: 2, Name: "Kwale", Capital: "Kwale", Governor: "Fatuma Achani", Population: 866820, Area: 8270.3},
		{Code: 3, Name: "Kilifi", Capital: "Kilifi", Governor: "Gideon Mung'aro", Population: 1453787, Area: 12245.9},
		{Code: 22, Name: "Nairobi", Capital: "Nairobi City", Governor: "Johnson Sakaja", Population: 4397073, Area: 696.1},
		{Code: 30, Name: "Narok", Capital: "Narok", Governor: "Patrick Ntutu", Population: 1157873, Area: 17944.1},
		{Code: 47, Name: "Nairobi", Capital: "Nairobi City", Governor: "Johnson Sakaja", Population: 4397073, Area: 696.1},
	}

	log.Println("Seeding counties...")
	for _, county := range counties {
		if err := countyRepo.Create(ctx, &county); err != nil {
			log.Printf("Warning: Could not create county %s: %v\n", county.Name, err)
		} else {
			log.Printf("Created county: %s (Code: %d)\n", county.Name, county.Code)
		}
	}

	// Seed wards
	wards := []models.Ward{
		{Name: "Mvita", CountyName: "Mombasa", ConstituencyName: "Mvita", Population: 90000},
		{Name: "Kilindini", CountyName: "Mombasa", ConstituencyName: "Mvita", Population: 85000},
		{Name: "Changamwe", CountyName: "Mombasa", ConstituencyName: "Changamwe", Population: 120000},
		{Name: "Karen", CountyName: "Nairobi", ConstituencyName: "Langata", Population: 95000},
		{Name: "Westlands", CountyName: "Nairobi", ConstituencyName: "Westlands", Population: 110000},
	}

	log.Println("Seeding wards...")
	for _, ward := range wards {
		if err := wardRepo.Create(ctx, &ward); err != nil {
			log.Printf("Warning: Could not create ward %s: %v\n", ward.Name, err)
		} else {
			log.Printf("Created ward: %s (%s County)\n", ward.Name, ward.CountyName)
		}
	}

	// Seed leaders
	leaders := []models.Leader{
		{
			Name:      "William Ruto",
			Position:  "President",
			County:    "National",
			Party:     "UDA",
			StartDate: time.Date(2022, 9, 13, 0, 0, 0, 0, time.UTC),
		},
		{
			Name:      "Rigathi Gachagua",
			Position:  "Deputy President",
			County:    "National",
			Party:     "UDA",
			StartDate: time.Date(2022, 9, 13, 0, 0, 0, 0, time.UTC),
		},
		{
			Name:      "Johnson Sakaja",
			Position:  "Governor",
			County:    "Nairobi",
			Party:     "UDA",
			StartDate: time.Date(2022, 8, 25, 0, 0, 0, 0, time.UTC),
		},
		{
			Name:      "Abdulswamad Nassir",
			Position:  "Governor",
			County:    "Mombasa",
			Party:     "ODM",
			StartDate: time.Date(2022, 8, 25, 0, 0, 0, 0, time.UTC),
		},
	}

	log.Println("Seeding leaders...")
	for _, leader := range leaders {
		if err := leaderRepo.Create(ctx, &leader); err != nil {
			log.Printf("Warning: Could not create leader %s: %v\n", leader.Name, err)
		} else {
			log.Printf("Created leader: %s (%s, %s)\n", leader.Name, leader.Position, leader.County)
		}
	}

	fmt.Println("\nSeeding completed successfully!")
	fmt.Printf("Total counties: %d\n", len(counties))
	fmt.Printf("Total wards: %d\n", len(wards))
	fmt.Printf("Total leaders: %d\n", len(leaders))
}
