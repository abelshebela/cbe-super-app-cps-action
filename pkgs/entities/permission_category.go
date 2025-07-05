package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type PermissionCategory struct {
	ID           bson.ObjectID   `bson:"_id,omitempty" json:"_id,omitempty"`
	CategoryName string          `bson:"categoryName" json:"categoryName"`
	Access       string          `bson:"access" json:"access"`
	Permissions  []bson.ObjectID `bson:"permissions" json:"permissions"`
	Enabled      bool            `bson:"enabled,omitempty" json:"enabled,omitempty"`
	IsDeleted    bool            `bson:"isDeleted,omitempty" json:"isDeleted,omitempty"`
	CreatedAt    time.Time       `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt    time.Time       `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}
