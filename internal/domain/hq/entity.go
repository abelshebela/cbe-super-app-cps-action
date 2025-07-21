package hq

import (
	"time"

	utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type HQ struct {
	ID            bson.ObjectID `json:"id" bson:"_id"`
	Name          string        `json:"name" bson:"name"`
	Enabled       bool          `json:"enabled" bson:"enabled"`
	IsDeleted     bool          `json:"is_deleted" bson:"is_deleted"`
	ArchiveExpiry uint          `json:"archive_expiry" bson:"archive_expiry"`
	BlockTime     uint          `json:"block_time" bson:"block_time"`
	ArchiveTime   uint          `json:"archive_time" bson:"archive_time"`
	CreatedAt     time.Time     `json:"created_at" bson:"created_at"`
	LastModified  time.Time     `json:"last_modified" bson:"last_modified"`
}

type HQRespose struct {
	Page  int   `json:"page"`
	HQ    []*HQ `json:"hq"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

type UpdateBlockTimeRequest struct {
	ID         string `json:"id"`
	BlockTime  uint   `json:"block_time"`
	MakerID    string `json:"maker_id"`
	MakerName  string `json:"maker_name,omitempty"`
	MakerPhone string `json:"maker_phone,omitempty"`
	Department string `json:"department,omitempty"`
}

type UpdateArchiveTimeRequest struct {
	ID          string `json:"id"`
	ArchiveTime uint   `json:"archive_time"`
	MakerID     string `json:"maker_id"`
	MakerName   string `json:"maker_name,omitempty"`
	MakerPhone  string `json:"maker_phone,omitempty"`
	Department  string `json:"department,omitempty"`
}

type ApproveRejectRequest struct {
	ActionCode     string            `json:"action_code"`
	CheckerID      string            `json:"checker_id"`
	CheckerName    string            `json:"checker_name,omitempty"`
	CheckerPhone   string            `json:"checker_phone,omitempty"`
	Decision       utils.DecisonEnum `json:"decision"`
	RejectedReason string            `json:"rejected_reason"`
}
