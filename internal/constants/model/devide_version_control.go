package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type DeviceVersionControl struct {
	ID             bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	LatestVersion  string        `bson:"latest_version" json:"latest_version"`
	Platform       string        `bson:"platform" json:"platform"`
	CreatedBy      string        `bson:"created_by" json:"created_by"`
	CreatedAt      time.Time     `bson:"created_at" json:"created_at"`
	UpdatedBy      string        `bson:"updated_by" json:"updated_by"`
	UpdatedAt      time.Time     `bson:"updated_at" json:"updated_at"`
	ForceUpdate    bool          `bson:"force_update" json:"force_update"`
	ReleaseNotes   string        `bson:"release_notes" json:"release_notes"`
	Enabled        bool          `bson:"enabled" json:"enabled"`
	LastModifiedAt time.Time     `bson:"last_modified_at" json:"last_modified_at"`
}
