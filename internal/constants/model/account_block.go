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
	ID                string               `json:"id,omitempty"`
	Name              string               `json:"name"`
	Code              string               `json:"code"`
	ParentID          *string              `json:"parent_id,omitempty"`
	Parent            *AccountBlock        `json:"parent,omitempty"`
	Type              AccountBlockType     `json:"type"`
	RegionName        string               `json:"region_name,omitempty"`
	FederalRegionName string               `json:"federal_region_name,omitempty"`
	DistrictName      string               `json:"district_name,omitempty"`
	DaoCode           int64                `json:"dao_code,omitempty"`
	AccountType       string               `json:"account_type,omitempty"`
	IsEnabled         bool                 `json:"is_enabled"`
	CityID            *string              `json:"city_id,omitempty"`
	RegionID          *string              `json:"region_id,omitempty"`
	DistrictID        *string              `json:"district_id,omitempty"`
	DisableReason     []AccountBlockReason `json:"disable_reason"`
	IsDeleted         bool                 `json:"-"`
	CreatedAt         time.Time            `json:"created_at"`
	UpdatedAt         time.Time            `json:"updated_at"`
}

// AccountBlockReason is one disable event (stored in account_block_disable_reasons).
type AccountBlockReason struct {
	ID        string    `json:"id,omitempty"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
}
