package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Ward represents a Kenyan ward
type Ward struct {
	ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name             string             `json:"name" bson:"name"`
	CountyID         primitive.ObjectID `json:"county_id" bson:"county_id"`
	CountyName       string             `json:"county_name" bson:"county_name"`
	ConstituencyName string             `json:"constituency_name" bson:"constituency_name"`
	Population       int64              `json:"population" bson:"population"`
	CreatedAt        time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at" bson:"updated_at"`
}
