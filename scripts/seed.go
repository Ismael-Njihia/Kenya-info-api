package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/Ismael-Njihia/Kenya-info-api/internal/config"
	"github.com/Ismael-Njihia/Kenya-info-api/internal/database"
	"github.com/Ismael-Njihia/Kenya-info-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dataDir = "scripts/data"

func main() {
	log.Println("🚀 Starting production-ready seeder...")

	// load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := database.Connect(&cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer func() {
		_ = db.Client.Disconnect(context.Background())
	}()

	ctx := context.Background()

	// optionally clear collections first (commented by default)
	// clearCollections(ctx, db)

	// Load JSON files
	counties, err := loadCounties(dataDir + "/counties.json")
	if err != nil {
		log.Fatalf("loadCounties: %v", err)
	}
	constituencies, err := loadConstituencies(dataDir + "/constituencies.json")
	if err != nil {
		log.Fatalf("loadConstituencies: %v", err)
	}
	wards, err := loadWards(dataDir + "/wards.json")
	if err != nil {
		log.Fatalf("loadWards: %v", err)
	}
	leaders, err := loadLeaders(dataDir + "/leaders.json")
	if err != nil {
		log.Fatalf("loadLeaders: %v", err)
	}

	// Upsert counties (unique by Code OR Name)
	countyMap, err := upsertCounties(ctx, db, counties)
	if err != nil {
		log.Fatalf("upsertCounties: %v", err)
	}

	// Upsert constituencies (link to counties)
	constMap, err := upsertConstituencies(ctx, db, constituencies, countyMap)
	if err != nil {
		log.Fatalf("upsertConstituencies: %v", err)
	}

	// Upsert wards (link to constituency & county)
	wardMap, err := upsertWards(ctx, db, wards, countyMap, constMap)
	if err != nil {
		log.Fatalf("upsertWards: %v", err)
	}

	// Upsert leaders (link to county/constituency/ward where provided)
	if err := upsertLeaders(ctx, db, leaders, countyMap, constMap, wardMap); err != nil {
		log.Fatalf("upsertLeaders: %v", err)
	}

	log.Println("✅ Seeder finished successfully.")
}

/* ---------- Helpers to load JSON ---------- */

func readJSONFile(path string, v interface{}) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

func loadCounties(path string) ([]models.County, error) {
	var c []models.County
	if err := readJSONFile(path, &c); err != nil {
		return nil, fmt.Errorf("read counties json: %w", err)
	}
	return c, nil
}

func loadConstituencies(path string) ([]models.Constituency, error) {
	var c []models.Constituency
	if err := readJSONFile(path, &c); err != nil {
		return nil, fmt.Errorf("read constituencies json: %w", err)
	}
	return c, nil
}

func loadWards(path string) ([]models.Ward, error) {
	var w []models.Ward
	if err := readJSONFile(path, &w); err != nil {
		return nil, fmt.Errorf("read wards json: %w", err)
	}
	return w, nil
}

func loadLeaders(path string) ([]models.Leader, error) {
	var l []models.Leader
	if err := readJSONFile(path, &l); err != nil {
		return nil, fmt.Errorf("read leaders json: %w", err)
	}
	return l, nil
}

/* ---------- Upsert functions (idempotent) ---------- */

// upsertCounties returns map[name]ObjectID
func upsertCounties(ctx context.Context, db *database.Database, counties []models.County) (map[string]primitive.ObjectID, error) {
	coll := db.GetCollection("counties")
	out := make(map[string]primitive.ObjectID)
	now := time.Now()

	for _, county := range counties {
		// Set timestamps
		if county.CreatedAt.IsZero() {
			county.CreatedAt = now
		}
		county.UpdatedAt = now

		// Generate ID if not provided
		if county.ID.IsZero() {
			county.ID = primitive.NewObjectID()
		}

		// Unique filter: prefer Code, fallback to Name
		var filter bson.M
		if county.Code != 0 {
			filter = bson.M{"code": county.Code}
		} else {
			filter = bson.M{"name": county.Name}
		}

		// ✅ Exclude _id from being updated
		update := bson.M{
			"$set": bson.M{
				"name":       county.Name,
				"code":       county.Code,
				"created_at": county.CreatedAt,
				"updated_at": county.UpdatedAt,
			},
			"$setOnInsert": bson.M{
				"_id": county.ID,
			},
		}

		opts := options.Update().SetUpsert(true)
		if _, err := coll.UpdateOne(ctx, filter, update, opts); err != nil {
			return nil, fmt.Errorf("upsert county %s failed: %w", county.Name, err)
		}

		out[county.Name] = county.ID
		log.Printf("✅ Upserted county: %s (code=%d)", county.Name, county.Code)
	}

	return out, nil
}

// upsertConstituencies returns map[name]ObjectID
func upsertConstituencies(
	ctx context.Context,
	db *database.Database,
	list []models.Constituency,
	countyMap map[string]primitive.ObjectID,
) (map[string]primitive.ObjectID, error) {
	coll := db.GetCollection("constituencies")
	out := make(map[string]primitive.ObjectID)
	now := time.Now()

	for _, c := range list {
		// link county id if provided as name
		if c.CountyID.IsZero() && c.CountyName != "" {
			if id, ok := countyMap[c.CountyName]; ok {
				c.CountyID = id
			} else {
				log.Printf("⚠️ constituency %s references unknown county %s — skipping", c.Name, c.CountyName)
				continue
			}
		}

		// timestamps & id
		if c.CreatedAt.IsZero() {
			c.CreatedAt = now
		}
		c.UpdatedAt = now
		if c.ID.IsZero() {
			c.ID = primitive.NewObjectID()
		}

		filter := bson.M{"name": c.Name}

		// ✅ Use explicit field sets (never touch _id)
		update := bson.M{
			"$set": bson.M{
				"name":        c.Name,
				"county_id":   c.CountyID,
				"county_name": c.CountyName,
				"created_at":  c.CreatedAt,
				"updated_at":  c.UpdatedAt,
			},
			"$setOnInsert": bson.M{
				"_id": c.ID,
			},
		}

		opts := options.Update().SetUpsert(true)
		if _, err := coll.UpdateOne(ctx, filter, update, opts); err != nil {
			return nil, fmt.Errorf("upsert constituency %s: %w", c.Name, err)
		}
		out[c.Name] = c.ID
		log.Printf("⤴️ Upserted constituency: %s → county=%s", c.Name, c.CountyName)
	}

	return out, nil
}

// upsertWards returns map[name]ObjectID
func upsertWards(
	ctx context.Context,
	db *database.Database,
	list []models.Ward,
	countyMap map[string]primitive.ObjectID,
	constMap map[string]primitive.ObjectID,
) (map[string]primitive.ObjectID, error) {

	coll := db.GetCollection("wards")
	out := make(map[string]primitive.ObjectID)
	now := time.Now()

	for _, w := range list {
		// 🧭 Link County
		if w.CountyID.IsZero() && w.CountyName != "" {
			if id, ok := countyMap[w.CountyName]; ok {
				w.CountyID = id
			} else {
				log.Printf("⚠️ ward %s references unknown county %s — skipping", w.Name, w.CountyName)
				continue
			}
		}

		// 🧭 Link Constituency
		if w.ConstituencyID.IsZero() && w.ConstituencyName != "" {
			if id, ok := constMap[w.ConstituencyName]; ok {
				w.ConstituencyID = id
			} else {
				log.Printf("⚠️ ward %s references unknown constituency %s — skipping", w.Name, w.ConstituencyName)
				continue
			}
		}

		// 🕒 Handle timestamps & IDs
		if w.CreatedAt.IsZero() {
			w.CreatedAt = now
		}
		w.UpdatedAt = now
		if w.ID.IsZero() {
			w.ID = primitive.NewObjectID()
		}

		// 🎯 Unique filter: ward name + constituency
		filter := bson.M{
			"name":            w.Name,
			"constituency_id": w.ConstituencyID,
		}

		// ✅ Only set allowed fields (never _id)
		update := bson.M{
			"$set": bson.M{
				"name":              w.Name,
				"county_id":         w.CountyID,
				"county_name":       w.CountyName,
				"constituency_id":   w.ConstituencyID,
				"constituency_name": w.ConstituencyName,
				"created_at":        w.CreatedAt,
				"updated_at":        w.UpdatedAt,
			},
			"$setOnInsert": bson.M{
				"_id": w.ID,
			},
		}

		opts := options.Update().SetUpsert(true)
		if _, err := coll.UpdateOne(ctx, filter, update, opts); err != nil {
			return nil, fmt.Errorf("upsert ward %s: %w", w.Name, err)
		}

		out[w.Name] = w.ID
		log.Printf("⤴️ Upserted ward: %s → constituency=%s, county=%s", w.Name, w.ConstituencyName, w.CountyName)
	}

	return out, nil
}

// upsertLeaders inserts leaders and links them to referenced entities when provided
func upsertLeaders(
	ctx context.Context,
	db *database.Database,
	list []models.Leader,
	countyMap map[string]primitive.ObjectID,
	constMap map[string]primitive.ObjectID,
	wardMap map[string]primitive.ObjectID,
) error {
	coll := db.GetCollection("leaders")
	now := time.Now()

	for _, l := range list {
		// 🧭 Resolve text references into ObjectIDs
		if l.CountyID.IsZero() && l.CountyName != "" {
			if id, ok := countyMap[l.CountyName]; ok {
				l.CountyID = id
			} else {
				log.Printf("⚠️ leader %s references unknown county %s — continuing", l.Name, l.CountyName)
			}
		}
		if l.ConstituencyID.IsZero() && l.ConstituencyName != "" {
			if id, ok := constMap[l.ConstituencyName]; ok {
				l.ConstituencyID = id
			} else {
				log.Printf("⚠️ leader %s references unknown constituency %s — continuing", l.Name, l.ConstituencyName)
			}
		}
		if l.WardID.IsZero() && l.WardName != "" {
			if id, ok := wardMap[l.WardName]; ok {
				l.WardID = id
			} else {
				log.Printf("⚠️ leader %s references unknown ward %s — continuing", l.Name, l.WardName)
			}
		}

		// 🕒 Handle timestamps
		if l.CreatedAt.IsZero() {
			l.CreatedAt = now
		}
		l.UpdatedAt = now

		// 🆔 Ensure ID exists
		if l.ID.IsZero() {
			l.ID = primitive.NewObjectID()
		}

		// 🎯 Unique filter: name + position (so the same person/role combination isn’t duplicated)
		filter := bson.M{
			"name":     l.Name,
			"position": l.Position,
		}

		// ✅ Explicitly control what gets updated (never update _id)
		update := bson.M{
			"$set": bson.M{
				"name":              l.Name,
				"position":          l.Position,
				"county_id":         l.CountyID,
				"county_name":       l.CountyName,
				"constituency_id":   l.ConstituencyID,
				"constituency_name": l.ConstituencyName,
				"ward_id":           l.WardID,
				"ward_name":         l.WardName,
				"party":             l.Party,
				"email":             l.Email,
				"phone":             l.Phone,
				"updated_at":        l.UpdatedAt,
			},
			"$setOnInsert": bson.M{
				"_id":        l.ID,
				"created_at": l.CreatedAt,
			},
		}

		opts := options.Update().SetUpsert(true)
		if _, err := coll.UpdateOne(ctx, filter, update, opts); err != nil {
			return fmt.Errorf("upsert leader %s: %w", l.Name, err)
		}

		log.Printf("👤 Upserted leader: %s (%s) → county=%s, constituency=%s, ward=%s",
			l.Name, l.Position, l.CountyName, l.ConstituencyName, l.WardName)
	}

	return nil
}

/* (optional) clear collections before seeding */
func clearCollections(ctx context.Context, db *database.Database) {
	names := []string{"counties", "constituencies", "wards", "leaders"}
	for _, n := range names {
		if err := db.GetCollection(n).Drop(ctx); err != nil {
			log.Printf("clear %s: %v", n, err)
		} else {
			log.Printf("cleared %s", n)
		}
	}
}
