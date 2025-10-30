package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// County represents a Kenyan county
type County struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Code       int                `json:"code" bson:"code" binding:"required"`
	Name       string             `json:"name" bson:"name" binding:"required"`
	Capital    string             `json:"capital" bson:"capital"`
	Population int64              `json:"population" bson:"population"`
	Area       float64            `json:"area" bson:"area"`
	Governor   string             `json:"governor,omitempty" bson:"governor,omitempty"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
}

// Constituency represents a constituency within a county
type Constituency struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name" bson:"name" binding:"required"`
	CountyID   primitive.ObjectID `json:"county_id" bson:"county_id" binding:"required"`
	CountyName string             `json:"county_name" bson:"county_name"`
	MP         string             `json:"mp,omitempty" bson:"mp,omitempty"`
	Population int64              `json:"population" bson:"population"`
	Area       float64            `json:"area" bson:"area"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
}

// Ward represents a ward within a constituency
type Ward struct {
	ID               primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name             string             `json:"name" bson:"name" binding:"required"`
	CountyID         primitive.ObjectID `json:"county_id" bson:"county_id" binding:"required"`
	CountyName       string             `json:"county_name" bson:"county_name"`
	ConstituencyID   primitive.ObjectID `json:"constituency_id" bson:"constituency_id" binding:"required"`
	ConstituencyName string             `json:"constituency_name" bson:"constituency_name"`
	MCA              string             `json:"mca,omitempty" bson:"mca,omitempty"`
	Population       int64              `json:"population" bson:"population"`
	Area             float64            `json:"area,omitempty" bson:"area,omitempty"`
	CreatedAt        time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at" bson:"updated_at"`
}

// Leader represents a political or administrative leader
type Leader struct {
	ID               primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name             string             `json:"name" bson:"name" binding:"required"`
	Position         string             `json:"position" bson:"position" binding:"required"`
	CountyID         primitive.ObjectID `json:"county_id,omitempty" bson:"county_id,omitempty"`
	CountyName       string             `json:"county_name,omitempty" bson:"county_name,omitempty"`
	ConstituencyID   primitive.ObjectID `json:"constituency_id,omitempty" bson:"constituency_id,omitempty"`
	ConstituencyName string             `json:"constituency_name,omitempty" bson:"constituency_name,omitempty"`
	WardID           primitive.ObjectID `json:"ward_id,omitempty" bson:"ward_id,omitempty"`
	WardName         string             `json:"ward_name,omitempty" bson:"ward_name,omitempty"`
	Party            string             `json:"party,omitempty" bson:"party,omitempty"`
	Email            string             `json:"email,omitempty" bson:"email,omitempty"`
	Phone            string             `json:"phone,omitempty" bson:"phone,omitempty"`
	CreatedAt        time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at" bson:"updated_at"`
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalCount int64       `json:"total_count"`
	TotalPages int         `json:"total_pages"`
}
