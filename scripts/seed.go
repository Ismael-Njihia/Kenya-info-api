package main

import (
	"context"
	"encoding/json"
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
	log.Println("🚀 Starting clean Kenya data seeder...")

	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	// Connect DB
	db, err := database.Connect(&cfg.Database)
	if err != nil {
		log.Fatalf("❌ Failed to connect DB: %v", err)
	}
	defer func() {
		_ = db.Client.Disconnect(context.Background())
	}()

	ctx := context.Background()

	// Load datasets
	counties := mustLoadJSON[[]models.County](dataDir + "/counties.json")
	constituencies := mustLoadJSON[[]models.Constituency](dataDir + "/constituencies.json")
	wards := mustLoadJSON[[]models.Ward](dataDir + "/wards.json")
	leaders := mustLoadJSON[[]models.Leader](dataDir + "/leaders.json")

	// Optional — clear previous data
	clearCollections(ctx, db)

	// Upsert in order
	countyMap := mustUpsertCounties(ctx, db, counties)
	constMap := mustUpsertConstituencies(ctx, db, constituencies, countyMap)
	wardMap := mustUpsertWards(ctx, db, wards, countyMap, constMap)
	mustUpsertLeaders(ctx, db, leaders, countyMap, constMap, wardMap)

	log.Println("✅ Done! All data seeded successfully.")
}

/* ---------- Utilities ---------- */

func mustLoadJSON[T any](path string) T {
	var data T
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("❌ Read file %s: %v", path, err)
	}
	if err := json.Unmarshal(b, &data); err != nil {
		log.Fatalf("❌ Unmarshal %s: %v", path, err)
	}
	return data
}

/* ---------- Upsert Logic ---------- */

func mustUpsertCounties(ctx context.Context, db *database.Database, list []models.County) map[string]primitive.ObjectID {
	coll := db.GetCollection("counties")
	now := time.Now()
	out := map[string]primitive.ObjectID{}

	for _, c := range list {
		if c.ID.IsZero() {
			c.ID = primitive.NewObjectID()
		}
		c.CreatedAt = firstNonZero(c.CreatedAt, now)
		c.UpdatedAt = now

		filter := bson.M{"code": c.Code}
		update := bson.M{
			"$set": bson.M{
				"name":       c.Name,
				"code":       c.Code,
				"updated_at": c.UpdatedAt,
			},
			"$setOnInsert": bson.M{
				"_id":        c.ID,
				"created_at": c.CreatedAt,
			},
		}

		opts := options.Update().SetUpsert(true)
		if _, err := coll.UpdateOne(ctx, filter, update, opts); err != nil {
			log.Fatalf("❌ upsert county %s: %v", c.Name, err)
		}
		out[c.Name] = c.ID
		log.Printf("✅ County: %s (code %d)", c.Name, c.Code)
	}

	return out
}

func mustUpsertConstituencies(ctx context.Context, db *database.Database, list []models.Constituency, countyMap map[string]primitive.ObjectID) map[string]primitive.ObjectID {
	coll := db.GetCollection("constituencies")
	now := time.Now()
	out := map[string]primitive.ObjectID{}

	for _, c := range list {
		if c.CountyName != "" {
			if id, ok := countyMap[c.CountyName]; ok {
				c.CountyID = id
			} else {
				log.Printf("⚠️ Constituency %s unknown county %s", c.Name, c.CountyName)
				continue
			}
		}

		if c.ID.IsZero() {
			c.ID = primitive.NewObjectID()
		}
		c.CreatedAt = firstNonZero(c.CreatedAt, now)
		c.UpdatedAt = now

		filter := bson.M{"name": c.Name}
		update := bson.M{
			"$set": bson.M{
				"name":        c.Name,
				"county_id":   c.CountyID,
				"county_name": c.CountyName,
				"updated_at":  c.UpdatedAt,
			},
			"$setOnInsert": bson.M{
				"_id":        c.ID,
				"created_at": c.CreatedAt,
			},
		}

		opts := options.Update().SetUpsert(true)
		if _, err := coll.UpdateOne(ctx, filter, update, opts); err != nil {
			log.Fatalf("❌ upsert constituency %s: %v", c.Name, err)
		}
		out[c.Name] = c.ID
		log.Printf("⤴️ Constituency: %s → %s", c.Name, c.CountyName)
	}
	return out
}

func mustUpsertWards(ctx context.Context, db *database.Database, list []models.Ward, countyMap, constMap map[string]primitive.ObjectID) map[string]primitive.ObjectID {
	coll := db.GetCollection("wards")
	now := time.Now()
	out := map[string]primitive.ObjectID{}

	for _, w := range list {
		if w.CountyName != "" {
			if id, ok := countyMap[w.CountyName]; ok {
				w.CountyID = id
			}
		}
		if w.ConstituencyName != "" {
			if id, ok := constMap[w.ConstituencyName]; ok {
				w.ConstituencyID = id
			}
		}

		if w.ID.IsZero() {
			w.ID = primitive.NewObjectID()
		}
		w.CreatedAt = firstNonZero(w.CreatedAt, now)
		w.UpdatedAt = now

		filter := bson.M{"name": w.Name, "constituency_id": w.ConstituencyID}
		update := bson.M{
			"$set": bson.M{
				"name":              w.Name,
				"county_id":         w.CountyID,
				"county_name":       w.CountyName,
				"constituency_id":   w.ConstituencyID,
				"constituency_name": w.ConstituencyName,
				"updated_at":        w.UpdatedAt,
			},
			"$setOnInsert": bson.M{
				"_id":        w.ID,
				"created_at": w.CreatedAt,
			},
		}

		opts := options.Update().SetUpsert(true)
		if _, err := coll.UpdateOne(ctx, filter, update, opts); err != nil {
			log.Fatalf("❌ upsert ward %s: %v", w.Name, err)
		}
		out[w.Name] = w.ID
		log.Printf("🏙️ Ward: %s (%s → %s)", w.Name, w.ConstituencyName, w.CountyName)
	}
	return out
}

func mustUpsertLeaders(ctx context.Context, db *database.Database, list []models.Leader, countyMap, constMap, wardMap map[string]primitive.ObjectID) {
	coll := db.GetCollection("leaders")
	now := time.Now()

	for _, l := range list {
		if id, ok := countyMap[l.CountyName]; ok {
			l.CountyID = id
		}
		if id, ok := constMap[l.ConstituencyName]; ok {
			l.ConstituencyID = id
		}
		if id, ok := wardMap[l.WardName]; ok {
			l.WardID = id
		}

		if l.ID.IsZero() {
			l.ID = primitive.NewObjectID()
		}
		l.CreatedAt = firstNonZero(l.CreatedAt, now)
		l.UpdatedAt = now

		filter := bson.M{"name": l.Name, "position": l.Position}
		update := bson.M{
			"$set": bson.M{
				"name":              l.Name,
				"position":          l.Position,
				"party":             l.Party,
				"email":             l.Email,
				"phone":             l.Phone,
				"county_id":         l.CountyID,
				"county_name":       l.CountyName,
				"constituency_id":   l.ConstituencyID,
				"constituency_name": l.ConstituencyName,
				"ward_id":           l.WardID,
				"ward_name":         l.WardName,
				"updated_at":        l.UpdatedAt,
			},
			"$setOnInsert": bson.M{
				"_id":        l.ID,
				"created_at": l.CreatedAt,
			},
		}

		opts := options.Update().SetUpsert(true)
		if _, err := coll.UpdateOne(ctx, filter, update, opts); err != nil {
			log.Fatalf("❌ upsert leader %s: %v", l.Name, err)
		}
		log.Printf("👤 Leader: %s (%s) — %s/%s/%s", l.Name, l.Position, l.CountyName, l.ConstituencyName, l.WardName)
	}
}

/* ---------- Helpers ---------- */

func firstNonZero(t time.Time, fallback time.Time) time.Time {
	if t.IsZero() {
		return fallback
	}
	return t
}

func clearCollections(ctx context.Context, db *database.Database) {
	names := []string{"counties", "constituencies", "wards", "leaders"}
	for _, n := range names {
		if err := db.GetCollection(n).Drop(ctx); err != nil {
			log.Printf("⚠️ Drop %s: %v", n, err)
		} else {
			log.Printf("🧹 Cleared %s", n)
		}
	}
}
