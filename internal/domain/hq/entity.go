package hq

import "time"

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
	BlockTime  uint   `json:"block_time"`
	MakerID    string `json:"maker_id"`
	MakerName  string `json:"maker_name,omitempty"`
	MakerPhone string `json:"maker_phone,omitempty"`
}

type UpdateArchiveTimeRequest struct {
	ArchiveTime uint   `json:"archive_time"`
	MakerID     string `json:"maker_id"`
	MakerName   string `json:"maker_name,omitempty"`
	MakerPhone  string `json:"maker_phone,omitempty"`
}

type ApproveRejectRequest struct {
	ActionCode   string `json:"action_code"`
	CheckerID    string `json:"checker_id"`
	CheckerName  string `json:"checker_name,omitempty"`
	CheckerPhone string `json:"checker_phone,omitempty"`
	Approved     bool   `json:"approved"`
}
