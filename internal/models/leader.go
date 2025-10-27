package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Leader represents a Kenyan political leader
type Leader struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Position    string             `json:"position" bson:"position"`
	County      string             `json:"county" bson:"county"`
	Party       string             `json:"party" bson:"party"`
	Email       string             `json:"email,omitempty" bson:"email,omitempty"`
	Phone       string             `json:"phone,omitempty" bson:"phone,omitempty"`
	StartDate   time.Time          `json:"start_date" bson:"start_date"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}
