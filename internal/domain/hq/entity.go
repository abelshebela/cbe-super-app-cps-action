package hq

import (
	"time"

	utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type HQ struct {
	ID                      bson.ObjectID `json:"id" bson:"_id"`
	BlockTime               uint32        `json:"block_time" bson:"block_time"`
	ArchiveTime             uint32        `json:"archive_time" bson:"archive_time"`
	PasswordExpiry          uint32        `json:"password_expiry" bson:"password_expiry"`
	CreatedAtPasswordExpiry time.Time     `json:"created_at_password_expiry" bson:"created_at_password_expiry"`
	UpdatedAtPasswordExpiry time.Time     `json:"updated_at_password_expiry" bson:"updated_at_password_expiry"`
	CreatedAtBlock          time.Time     `json:"created_at_block" bson:"created_at_block"`
	UpdatedAtBlock          time.Time     `json:"updated_at_block" bson:"updated_at_block"`
	CreatedAtArchive        time.Time     `json:"created_at_archive" bson:"created_at_archive"`
	UpdatedAtArchive        time.Time     `json:"updated_at_archive" bson:"updated_at_archive"`
}

type HQRespose struct {
	Page  int   `json:"page"`
	HQ    []*HQ `json:"hq"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

type BlockTimeResponse struct {
	BlockTime      uint32    `json:"block_time"`
	CreatedAtBlock time.Time `json:"created_at_block"`
	UpdatedAtBlock time.Time `json:"updated_at_block"`
}

type ArchiveTimeResponse struct {
	ArchiveTime      uint32    `json:"archive_time"`
	CreatedAtArchive time.Time `json:"created_at_archive"`
	UpdatedAtArchive time.Time `json:"updated_at_archive"`
}

type PasswordExpiryResponse struct {
	PasswordExpiry          uint32    `json:"password_expiry"`
	CreatedAtPasswordExpiry time.Time `json:"created_at_password_expiry"`
	UpdatedAtPasswordExpiry time.Time `json:"updated_at_password_expiry"`
}

type UpdateBlockTimeRequest struct {
	BlockTime  uint32 `json:"block_time"`
	MakerID    string `json:"maker_id"`
	MakerName  string `json:"maker_name,omitempty"`
	MakerPhone string `json:"maker_phone,omitempty"`
	Department string `json:"department,omitempty"`
}

type UpdateArchiveTimeRequest struct {
	ArchiveTime uint32 `json:"archive_time"`
	MakerID     string `json:"maker_id"`
	MakerName   string `json:"maker_name,omitempty"`
	MakerPhone  string `json:"maker_phone,omitempty"`
	Department  string `json:"department,omitempty"`
}

type UpdatePasswordExpiryRequest struct {
	PasswordExpiry uint32 `json:"password_expiry"`
	MakerID        string `json:"maker_id"`
	MakerName      string `json:"maker_name,omitempty"`
	MakerPhone     string `json:"maker_phone,omitempty"`
	Department     string `json:"department,omitempty"`
}

type ApproveRejectRequest struct {
	ActionCode     string            `json:"action_code"`
	CheckerID      string            `json:"checker_id"`
	CheckerName    string            `json:"checker_name,omitempty"`
	CheckerPhone   string            `json:"checker_phone,omitempty"`
	Decision       utils.DecisonEnum `json:"decision"`
	RejectedReason string            `json:"rejected_reason"`
}
