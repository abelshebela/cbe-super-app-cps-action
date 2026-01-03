package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type HQ struct {
	ID                      bson.ObjectID `bson:"_id" json:"id"`
	UniqueID                string        `bson:"unique_id" json:"unique_id"`
	LatestiOSVersion        string        `bson:"latest_ios_version" json:"latest_ios_version"`
	LatestAndroidVersion    string        `bson:"latest_android_version" json:"latest_android_version"`
	ArchiveExpiry           uint32        `json:"archive_expiry" bson:"archive_expiry"`
	BlockTime               uint32        `json:"block_time" bson:"block_time"`
	BlockTimeStatus         string        `bson:"block_time_status" json:"block_time_status"`
	ArchiveTime             uint32        `json:"archive_time" bson:"archive_time"`
	ArchiveTimeStatus       string        `bson:"archive_time_status" json:"archive_time_status"`
	PasswordExpiry          uint32        `json:"password_expiry" bson:"password_expiry"`
	CreatedAtPasswordExpiry time.Time     `json:"created_at_password_expiry" bson:"created_at_password_expiry"`
	UpdatedAtPasswordExpiry time.Time     `json:"updated_at_password_expiry" bson:"updated_at_password_expiry"`
	CreatedAtBlock          time.Time     `json:"created_at_block" bson:"created_at_block"`
	UpdatedAtBlock          time.Time     `json:"updated_at_block" bson:"updated_at_block"`
	CreatedAtArchive        time.Time     `json:"created_at_archive" bson:"created_at_archive"`
	UpdatedAtArchive        time.Time     `json:"updated_at_archive" bson:"updated_at_archive"`
	TotalCap                uint64        `json:"total_cap" bson:"total_cap"`
	CreatedAtTotalCap       time.Time     `json:"created_at_total_cap" bson:"created_at_total_cap"`
	UpdatedAtTotalCap       time.Time     `json:"updated_at_total_cap" bson:"updated_at_total_cap"`
	Enabled                 bool          `bson:"enabled" json:"enabled"`
	IsDeleted               bool          `bson:"is_deleted" json:"is_deleted"`
	CreatedAt               time.Time     `bson:"created_at" json:"created_at"`
	LastModified            time.Time     `bson:"last_modified" json:"last_modified"`
}
