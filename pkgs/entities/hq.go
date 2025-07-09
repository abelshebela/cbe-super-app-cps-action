package entities

import "time"

type HQ struct {
	ID                   string          `bson:"_id" json:"id"`
	UniqueID             string          `bson:"unique_id" json:"unique_id"`
	Name                 string          `bson:"name" json:"name"`
	Address              string          `bson:"address" json:"address"`
	PhoneNumber          string          `bson:"phoneNumber" json:"phoneNumber"`
	Email                string          `bson:"email" json:"email"`
	LinkedAccounts       []LinkedAccount `bson:"linkedAccounts" json:"linkedAccounts"`
	LatestiOSVersion     string          `bson:"latestiOSVersion" json:"latestiOSVersion"`
	LatestAndroidVersion string          `bson:"latestAndroidVersion" json:"latestAndroidVersion"`
	ArchiveExpiry        uint            `json:"archive_expiry" bson:"archive_expiry"`
	BlockTime            uint            `json:"block_time" bson:"block_time"`
	BlockTimeStatus      string          `bson:"block_time_status" json:"block_time_status"`
	ArchiveTime          uint            `bson:"archive_time" json:"archive_time"`
	ArchiveTimeStatus    string          `bson:"archive_time_status" json:"archive_time_status"`
	Enabled              bool            `bson:"enabled" json:"enabled"`
	IsDeleted            bool            `bson:"isDeleted" json:"isDeleted"`
	CreatedAt            time.Time       `bson:"createdAt" json:"createdAt"`
	LastModified         time.Time       `bson:"lastModified" json:"lastModified"`
}
