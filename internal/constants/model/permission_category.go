package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type PermissionCategory struct {
	ID           bson.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	CategoryName string        `bson:"category_name,omitempty" json:"category_name,omitempty"`
	PortalCard   string        `bson:"portal_card,omitempty" json:"portal_card,omitempty"`
	Access       string        `bson:"access" json:"access"`
	Permissions  interface{}   `bson:"permissions" json:"permissions"`
	Enabled      bool          `bson:"enabled,omitempty" json:"enabled,omitempty"`
	IsDeleted    bool          `bson:"isDeleted,omitempty" json:"is_deleted,omitempty"`
	CreatedAt    time.Time     `bson:"createdAt,omitempty" json:"created_at,omitempty"`
	UpdatedAt    time.Time     `bson:"updatedAt,omitempty" json:"updated_at,omitempty"`
}
