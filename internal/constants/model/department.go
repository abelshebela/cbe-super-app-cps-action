package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Department struct {
	ID             bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	DepartmentCode string        `bson:"department_code" json:"department_code"`
	Department     string        `bson:"department" json:"department"`
	PortalCards    []string      `bson:"portal_cards" json:"portal_cards"`
	Enabled        bool          `bson:"enabled" json:"enabled"`
	IsDeleted      bool          `bson:"is_deleted" json:"is_deleted"`
	CreatedAt      time.Time     `bson:"created_at" json:"created_at"`
	LastModified   time.Time     `bson:"last_modified" json:"last_modified"`
}
