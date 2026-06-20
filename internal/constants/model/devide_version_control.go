package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// DeviceState enum values
const (
	DeviceStateForceUpdate  = "FORCE_UPDATE"
	DeviceStateMaintenance  = "MAINTENANCE"
	DeviceStateStable       = "STABLE"
)

type DeviceVersionControl struct {
	ID                bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	LatestVersion     string        `bson:"latest_version" json:"latest_version"`
	Platform          string        `bson:"platform" json:"platform"`
	CreatedBy         string        `bson:"created_by" json:"created_by"`
	CreatedAt         time.Time     `bson:"created_at,omitempty" json:"created_at"`
	UpdatedBy         string        `bson:"updated_by" json:"updated_by"`
	UpdatedAt         time.Time     `bson:"updated_at,omitempty" json:"updated_at"`
	ForceUpdate       bool          `bson:"force_update" json:"force_update"`
	ReleaseNotes      string        `bson:"release_notes" json:"release_notes"`
	Enabled           bool          `bson:"enabled" json:"enabled"`
	IsMaintenanceMode bool          `bson:"is_maintenance_mode" json:"is_maintenance_mode"`
	DeviceState       string        `bson:"device_state" json:"device_state"`
	LastModifiedAt    time.Time     `bson:"last_modified_at,omitempty" json:"last_modified_at"`
}
