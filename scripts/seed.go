package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Ismael-Njihia/Kenya-info-api/internal/config"
	"github.com/Ismael-Njihia/Kenya-info-api/internal/database"
	"github.com/Ismael-Njihia/Kenya-info-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	log.Println("Starting database seeding...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := database.Connect(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := db.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting from database: %v", err)
		}
	}()

	ctx := context.Background()

	// Clear existing data (optional - comment out if you want to keep existing data)
	log.Println("Clearing existing data...")
	if err := clearCollections(ctx, db); err != nil {
		log.Printf("Warning: Failed to clear collections: %v", err)
	}

	// Seed counties
	log.Println("Seeding counties...")
	countyIDs, err := seedCounties(ctx, db)
	if err != nil {
		log.Fatalf("Failed to seed counties: %v", err)
	}

	// Seed wards
	log.Println("Seeding wards...")
	wardIDs, err := seedWards(ctx, db, countyIDs)
	if err != nil {
		log.Fatalf("Failed to seed wards: %v", err)
	}

	// Seed leaders
	log.Println("Seeding leaders...")
	if err := seedLeaders(ctx, db, countyIDs, wardIDs); err != nil {
		log.Fatalf("Failed to seed leaders: %v", err)
	}

	log.Println("Database seeding completed successfully!")
}

func clearCollections(ctx context.Context, db *database.Database) error {
	collections := []string{"counties", "wards", "leaders"}
	for _, coll := range collections {
		if err := db.GetCollection(coll).Drop(ctx); err != nil {
			return err
		}
	}
	return nil
}

func seedCounties(ctx context.Context, db *database.Database) (map[string]interface{}, error) {
	counties := []models.County{
		{Code: 1, Name: "Mombasa", Capital: "Mombasa City", Population: 1208333, Area: 212.5, Governor: "Abdulswamad Nassir"},
		{Code: 2, Name: "Kwale", Capital: "Kwale", Population: 866820, Area: 8270.0, Governor: "Fatuma Achani"},
		{Code: 3, Name: "Kilifi", Capital: "Kilifi", Population: 1453787, Area: 12245.9, Governor: "Gideon Mung'aro"},
		{Code: 4, Name: "Tana River", Capital: "Hola", Population: 315943, Area: 35375.8, Governor: "Dhadho Godhana"},
		{Code: 5, Name: "Lamu", Capital: "Lamu", Population: 143920, Area: 6497.7, Governor: "Issa Timamy"},
		{Code: 6, Name: "Taita-Taveta", Capital: "Voi", Population: 340671, Area: 17083.9, Governor: "Andrew Mwadime"},
		{Code: 7, Name: "Garissa", Capital: "Garissa", Population: 841353, Area: 45720.2, Governor: "Nathif Jama"},
		{Code: 8, Name: "Wajir", Capital: "Wajir", Population: 781263, Area: 55840.6, Governor: "Ahmed Muktar"},
		{Code: 9, Name: "Mandera", Capital: "Mandera", Population: 1025756, Area: 25797.7, Governor: "Mohamed Khalif"},
		{Code: 10, Name: "Marsabit", Capital: "Marsabit", Population: 459785, Area: 66923.1, Governor: "Mohamud Ali"},
		{Code: 11, Name: "Isiolo", Capital: "Isiolo", Population: 268002, Area: 25336.1, Governor: "Abdi Guyo"},
		{Code: 12, Name: "Meru", Capital: "Meru", Population: 1545714, Area: 5127.0, Governor: "Kawira Mwangaza"},
		{Code: 13, Name: "Tharaka-Nithi", Capital: "Chuka", Population: 393177, Area: 2609.5, Governor: "Muthomi Njuki"},
		{Code: 14, Name: "Embu", Capital: "Embu", Population: 608599, Area: 2555.9, Governor: "Cecily Mbarire"},
		{Code: 15, Name: "Kitui", Capital: "Kitui", Population: 1136187, Area: 24385.1, Governor: "Julius Malombe"},
		{Code: 16, Name: "Machakos", Capital: "Machakos", Population: 1421932, Area: 5952.9, Governor: "Wavinya Ndeti"},
		{Code: 17, Name: "Makueni", Capital: "Wote", Population: 987653, Area: 8008.9, Governor: "Mutula Kilonzo Jr"},
		{Code: 18, Name: "Nyandarua", Capital: "Ol Kalou", Population: 638289, Area: 3107.7, Governor: "Moses Kiarie"},
		{Code: 19, Name: "Nyeri", Capital: "Nyeri", Population: 759164, Area: 2361.0, Governor: "Mutahi Kahiga"},
		{Code: 20, Name: "Kirinyaga", Capital: "Kerugoya", Population: 610411, Area: 1205.4, Governor: "Anne Waiguru"},
		{Code: 21, Name: "Murang'a", Capital: "Murang'a", Population: 1056640, Area: 2325.8, Governor: "Irungu Kang'ata"},
		{Code: 22, Name: "Kiambu", Capital: "Kiambu", Population: 2417735, Area: 2449.2, Governor: "Kimani Wamatangi"},
		{Code: 23, Name: "Turkana", Capital: "Lodwar", Population: 926976, Area: 68680.3, Governor: "Jeremiah Lomorukai"},
		{Code: 24, Name: "West Pokot", Capital: "Kapenguria", Population: 621241, Area: 9169.4, Governor: "Simon Kachapin"},
		{Code: 25, Name: "Samburu", Capital: "Maralal", Population: 310327, Area: 20182.5, Governor: "Jonathan Lelelit"},
		{Code: 26, Name: "Trans Nzoia", Capital: "Kitale", Population: 990341, Area: 2469.9, Governor: "George Natembeya"},
		{Code: 27, Name: "Uasin Gishu", Capital: "Eldoret", Population: 1163186, Area: 2955.3, Governor: "Jonathan Bii"},
		{Code: 28, Name: "Elgeyo-Marakwet", Capital: "Iten", Population: 454480, Area: 3029.9, Governor: "Wisley Rotich"},
		{Code: 29, Name: "Nandi", Capital: "Kapsabet", Population: 885711, Area: 2884.4, Governor: "Stephen Sang"},
		{Code: 30, Name: "Baringo", Capital: "Kabarnet", Population: 666763, Area: 11015.3, Governor: "Benjamin Cheboi"},
		{Code: 31, Name: "Laikipia", Capital: "Rumuruti", Population: 518560, Area: 9229.4, Governor: "Joshua Irungu"},
		{Code: 32, Name: "Nakuru", Capital: "Nakuru", Population: 2162202, Area: 7509.5, Governor: "Susan Kihika"},
		{Code: 33, Name: "Narok", Capital: "Narok", Population: 1157873, Area: 17944.1, Governor: "Patrick Ntutu"},
		{Code: 34, Name: "Kajiado", Capital: "Kajiado", Population: 1117840, Area: 21292.7, Governor: "Joseph Ole Lenku"},
		{Code: 35, Name: "Kericho", Capital: "Kericho", Population: 901777, Area: 2111.0, Governor: "Eric Mutai"},
		{Code: 36, Name: "Bomet", Capital: "Bomet", Population: 875689, Area: 1882.5, Governor: "Hillary Barchok"},
		{Code: 37, Name: "Kakamega", Capital: "Kakamega", Population: 1867579, Area: 3224.5, Governor: "Fernandes Barasa"},
		{Code: 38, Name: "Vihiga", Capital: "Vihiga", Population: 590013, Area: 563.0, Governor: "Wilber Ottichilo"},
		{Code: 39, Name: "Bungoma", Capital: "Bungoma", Population: 1670570, Area: 2206.9, Governor: "Kenneth Lusaka"},
		{Code: 40, Name: "Busia", Capital: "Busia", Population: 893681, Area: 1694.5, Governor: "Paul Otuoma"},
		{Code: 41, Name: "Siaya", Capital: "Siaya", Population: 993183, Area: 2530.5, Governor: "James Orengo"},
		{Code: 42, Name: "Kisumu", Capital: "Kisumu", Population: 1155574, Area: 2009.5, Governor: "Anyang' Nyong'o"},
		{Code: 43, Name: "Homa Bay", Capital: "Homa Bay", Population: 1131950, Area: 3154.7, Governor: "Gladys Wanga"},
		{Code: 44, Name: "Migori", Capital: "Migori", Population: 1116436, Area: 2586.4, Governor: "Ochilo Ayacko"},
		{Code: 45, Name: "Kisii", Capital: "Kisii", Population: 1266860, Area: 1317.9, Governor: "Simba Arati"},
		{Code: 46, Name: "Nyamira", Capital: "Nyamira", Population: 605576, Area: 912.5, Governor: "Amos Nyaribo"},
		{Code: 47, Name: "Nairobi", Capital: "Nairobi City", Population: 4397073, Area: 696.0, Governor: "Johnson Sakaja"},
	}

	countyIDs := make(map[string]interface{})
	collection := db.GetCollection("counties")

	for i := range counties {
		counties[i].CreatedAt = time.Now()
		counties[i].UpdatedAt = time.Now()

		result, err := collection.InsertOne(ctx, counties[i])
		if err != nil {
			return nil, fmt.Errorf("failed to insert county %s: %w", counties[i].Name, err)
		}
		countyIDs[counties[i].Name] = result.InsertedID
	}

	// Create index on county code
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "code", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	if _, err := collection.Indexes().CreateOne(ctx, indexModel); err != nil {
		log.Printf("Warning: Failed to create index: %v", err)
	}

	log.Printf("Successfully seeded %d counties", len(counties))
	return countyIDs, nil
}

func seedWards(ctx context.Context, db *database.Database, countyIDs map[string]interface{}) (map[string]interface{}, error) {
	// Sample wards for Nairobi county
	wards := []models.Ward{
		{Name: "Westlands", CountyID: countyIDs["Nairobi"].(interface{}), CountyName: "Nairobi", MCA: "Tim Wanyonyi", Population: 120000},
		{Name: "Kilimani", CountyID: countyIDs["Nairobi"].(interface{}), CountyName: "Nairobi", MCA: "Moses Ogeto", Population: 110000},
		{Name: "Kangemi", CountyID: countyIDs["Nairobi"].(interface{}), CountyName: "Nairobi", MCA: "Jayendra Malde", Population: 95000},
		{Name: "Karura", CountyID: countyIDs["Nairobi"].(interface{}), CountyName: "Nairobi", MCA: "Mary Mwami", Population: 80000},
		{Name: "Parklands", CountyID: countyIDs["Nairobi"].(interface{}), CountyName: "Nairobi", MCA: "David Mberia", Population: 75000},
		{Name: "Embakasi Central", CountyID: countyIDs["Nairobi"].(interface{}), CountyName: "Nairobi", MCA: "Michael Ogada", Population: 150000},
		{Name: "Embakasi East", CountyID: countyIDs["Nairobi"].(interface{}), CountyName: "Nairobi", MCA: "Francis Otieno", Population: 145000},
		{Name: "Embakasi North", CountyID: countyIDs["Nairobi"].(interface{}), CountyName: "Nairobi", MCA: "James Ndung'u", Population: 140000},
		{Name: "Embakasi South", CountyID: countyIDs["Nairobi"].(interface{}), CountyName: "Nairobi", MCA: "Michael Ochieng", Population: 135000},
		{Name: "Embakasi West", CountyID: countyIDs["Nairobi"].(interface{}), CountyName: "Nairobi", MCA: "Mark Mugambi", Population: 130000},
	}

	// Add sample wards for other counties
	otherCountyWards := []struct {
		county string
		wards  []string
	}{
		{"Mombasa", []string{"Mvita", "Likoni", "Nyali", "Kisauni", "Changamwe", "Jomvu"}},
		{"Kisumu", []string{"Kolwa Central", "Kolwa East", "Manyatta B", "Nyalenda A", "Nyalenda B"}},
		{"Nakuru", []string{"Bahati", "Gilgil", "Nakuru East", "Nakuru West", "Njoro", "Molo"}},
		{"Kiambu", []string{"Gatundu North", "Gatundu South", "Juja", "Kiambaa", "Kiambu", "Kikuyu"}},
	}

	wardIDs := make(map[string]interface{})
	collection := db.GetCollection("wards")

	// Insert Nairobi wards
	for i := range wards {
		wards[i].CreatedAt = time.Now()
		wards[i].UpdatedAt = time.Now()

		result, err := collection.InsertOne(ctx, wards[i])
		if err != nil {
			return nil, fmt.Errorf("failed to insert ward %s: %w", wards[i].Name, err)
		}
		wardIDs[wards[i].Name] = result.InsertedID
	}

	// Insert wards for other counties
	for _, countyWards := range otherCountyWards {
		if countyID, ok := countyIDs[countyWards.county]; ok {
			for _, wardName := range countyWards.wards {
				ward := models.Ward{
					Name:       wardName,
					CountyID:   countyID.(interface{}),
					CountyName: countyWards.county,
					Population: 50000,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				}
				result, err := collection.InsertOne(ctx, ward)
				if err != nil {
					return nil, fmt.Errorf("failed to insert ward %s: %w", wardName, err)
				}
				wardIDs[wardName] = result.InsertedID
			}
		}
	}

	log.Printf("Successfully seeded wards")
	return wardIDs, nil
}

func seedLeaders(ctx context.Context, db *database.Database, countyIDs, wardIDs map[string]interface{}) error {
	leaders := []models.Leader{
		// National Leaders
		{Name: "William Ruto", Position: "President", Party: "UDA", Email: "president@kenya.go.ke"},
		{Name: "Rigathi Gachagua", Position: "Deputy President", Party: "UDA", Email: "deputypresident@kenya.go.ke"},
		{Name: "Kithure Kindiki", Position: "Cabinet Secretary - Interior", Party: "UDA"},
		{Name: "Alfred Mutua", Position: "Cabinet Secretary - Foreign Affairs", Party: "UDA"},
		{Name: "Aden Duale", Position: "Cabinet Secretary - Defense", Party: "UDA"},

		// Senators (Sample)
		{Name: "Ledama Ole Kina", Position: "Senator", CountyID: countyIDs["Narok"].(interface{}), CountyName: "Narok", Party: "ODM"},
		{Name: "Edwin Sifuna", Position: "Senator", CountyID: countyIDs["Nairobi"].(interface{}), CountyName: "Nairobi", Party: "ODM"},
		{Name: "Karen Nyamu", Position: "Senator", CountyID: countyIDs["Nairobi"].(interface{}), CountyName: "Nairobi", Party: "UDA"},

		// Members of Parliament (Sample)
		{Name: "Kimani Ichung'wah", Position: "Member of Parliament", CountyID: countyIDs["Kiambu"].(interface{}), CountyName: "Kiambu", Party: "UDA"},
		{Name: "Raila Odinga", Position: "Member of Parliament", CountyID: countyIDs["Siaya"].(interface{}), CountyName: "Siaya", Party: "ODM"},
		{Name: "Martha Karua", Position: "Member of Parliament", CountyID: countyIDs["Kirinyaga"].(interface{}), CountyName: "Kirinyaga", Party: "NARC Kenya"},

		// County Assembly Members (Sample for Nairobi)
		{Name: "Tim Wanyonyi", Position: "MCA", CountyID: countyIDs["Nairobi"].(interface{}), CountyName: "Nairobi", WardID: wardIDs["Westlands"].(interface{}), WardName: "Westlands", Party: "ODM"},
		{Name: "Moses Ogeto", Position: "MCA", CountyID: countyIDs["Nairobi"].(interface{}), CountyName: "Nairobi", WardID: wardIDs["Kilimani"].(interface{}), WardName: "Kilimani", Party: "ODM"},
	}

	collection := db.GetCollection("leaders")

	for i := range leaders {
		leaders[i].CreatedAt = time.Now()
		leaders[i].UpdatedAt = time.Now()

		if _, err := collection.InsertOne(ctx, leaders[i]); err != nil {
			return fmt.Errorf("failed to insert leader %s: %w", leaders[i].Name, err)
		}
	}

	log.Printf("Successfully seeded %d leaders", len(leaders))
	return nil
}
