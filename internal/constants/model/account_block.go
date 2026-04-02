package model

import (
	"time"
)

type AccountBlockType string

const (
	TypeRegion   AccountBlockType = "R"
	TypeDistrict AccountBlockType = "D"
	TypeCity     AccountBlockType = "C"
	TypeBranch   AccountBlockType = "B"
)

type AccountBlock struct {
	ID            string             `json:"id,omitempty"`
	Name          string             `json:"name"`
	Code          string             `json:"code"`
	Address       string             `json:"address"`
	ParentID      *string            `json:"parent_id,omitempty"`
	Parent        *AccountBlock      `json:"parent,omitempty"`
	Slug          string             `json:"slug"`
	Type          AccountBlockType   `json:"type"`
	IsEnabled     bool               `json:"is_enabled"`
	CityID        *string            `json:"city_id,omitempty"`
	RegionID      *string            `json:"region_id,omitempty"`
	DistrictID    *string            `json:"district_id,omitempty"`
	DisableReason AccountBlockReason `json:"disable_reason"`
	IsDeleted     bool               `json:"-"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
}

type AccountBlockReason struct {
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
}
