package models

import (
	"time"
)

type HQ struct {
	ID                string    `bson:"_id" json:"id"`
	UniqueID          string    `bson:"unique_id" json:"unique_id"`
	Name              string    `bson:"name" json:"name"`
	BlockTime         uint      `bson:"block_time" json:"block_time"`
	BlockTimeStatus   string    `bson:"block_time_status" json:"block_time_status"`
	ArchiveTime       uint      `bson:"archive_time" json:"archive_time"`
	ArchiveTimeStatus string    `bson:"archive_time_status" json:"archive_time_status"`
	CreatedAt         time.Time `bson:"created_at" json:"created_at"`
	LastModifiedAt    time.Time `bson:"last_modified_at" json:"last_modified_at"`
}

type UpdateBlockTimeRequest struct {
	BlockTime uint   `json:"block_time"`
	MakerID   string `json:"maker_id"`
}

type UpdateArchiveTimeRequest struct {
	ArchiveTime uint   `json:"archive_time"`
	MakerID     string `json:"maker_id"`
}

type ApproveRejectRequest struct {
	ActionCode string `json:"action_code"`
	CheckerID  string `json:"checker_id"`
	Approved   bool   `json:"approved"`
}
