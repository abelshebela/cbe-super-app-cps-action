package hq

import (
	"time"

	utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
)

type HQ struct {
	ID        string
	Name      string
	Enabled   bool
	IsDeleted bool

	ArchiveExpiry uint
	BlockTime     uint
	ArchiveTime   uint
	CreatedAt     time.Time
	LastModified  time.Time
}

type UpdateBlockTimeRequest struct {
	ID         string `json:"id"`
	BlockTime  uint   `json:"block_time"`
	MakerID    string `json:"maker_id"`
	MakerName  string `json:"maker_name,omitempty"`
	MakerPhone string `json:"maker_phone,omitempty"`
}

type UpdateArchiveTimeRequest struct {
	ID          string `json:"id"`
	ArchiveTime uint   `json:"archive_time"`
	MakerID     string `json:"maker_id"`
	MakerName   string `json:"maker_name,omitempty"`
	MakerPhone  string `json:"maker_phone,omitempty"`
}

type ApproveRejectRequest struct {
	ActionCode     string            `json:"action_code"`
	CheckerID      string            `json:"checker_id"`
	CheckerName    string            `json:"checker_name,omitempty"`
	CheckerPhone   string            `json:"checker_phone,omitempty"`
	Decision       utils.DecisonEnum `json:"decision"`
	RejectedReason string            `json:"rejected_reason"`
}
