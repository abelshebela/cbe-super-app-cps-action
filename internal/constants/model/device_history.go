package model

import (
	"time"

	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type DeviceLinkHistroy struct {
	ID              bson.ObjectID      `json:"id,omitempty" bson:"_id,omitempty"`
	UserID          bson.ObjectID      `json:"user_id" bson:"user_id"`
	DeviceType      constants.Platform `json:"device_type" bson:"device_type"`
	DeviceUUID      string             `json:"device_uuid" bson:"device_uuid"`
	DeviceName      string             `json:"device_name" bson:"device_name"`
	DeviceOSVersion string             `json:"device_os_version" bson:"device_os_version"`
	LinkedAt        time.Time          `json:"linked_at" bson:"linked_at"`
	UnLinkedAt      time.Time          `json:"unlinked_at" bson:"unlinked_at"`
}
