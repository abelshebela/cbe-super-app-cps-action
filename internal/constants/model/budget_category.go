package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type BudgetCategory struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string        `bson:"name" json:"name"`
	Color     string        `bson:"color" json:"color"`
	Icon      string        `bson:"icon" json:"icon"`
	Enabled   bool          `bson:"enabled" json:"enabled"`
	IsDeleted bool          `bson:"is_deleted" json:"is_deleted"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
}