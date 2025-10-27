package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// County represents a Kenyan county
type County struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Code        int                `json:"code" bson:"code"`
	Name        string             `json:"name" bson:"name"`
	Capital     string             `json:"capital" bson:"capital"`
	Governor    string             `json:"governor" bson:"governor"`
	Population  int64              `json:"population" bson:"population"`
	Area        float64            `json:"area" bson:"area"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}
