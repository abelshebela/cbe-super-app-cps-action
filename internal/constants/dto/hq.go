package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type HQ struct {
	ID                   bson.ObjectID `bson:"_id" json:"id"`
	UniqueID             string        `bson:"unique_id" json:"unique_id"`
	LatestiOSVersion     string        `bson:"latest_ios_version" json:"latest_ios_version"`
	LatestAndroidVersion string        `bson:"latest_android_version" json:"latest_android_version"`
	ArchiveExpiry        uint          `json:"archive_expiry" bson:"archive_expiry"`
	BlockTime            uint          `json:"block_time" bson:"block_time"`
	BlockTimeStatus      string        `bson:"block_time_status" json:"block_time_status"`
	ArchiveTime          uint          `bson:"archive_time" json:"archive_time"`
	ArchiveTimeStatus    string        `bson:"archive_time_status" json:"archive_time_status"`
	Enabled              bool          `bson:"enabled" json:"enabled"`
	IsDeleted            bool          `bson:"is_deleted" json:"is_deleted"`
	CreatedAt            time.Time     `bson:"created_at" json:"created_at"`
	LastModified         time.Time     `bson:"last_modified" json:"last_modified"`
}
