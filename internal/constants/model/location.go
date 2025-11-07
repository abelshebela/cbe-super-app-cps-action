package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type AccountBlockType string

const (
	TypeRegion   AccountBlockType = "R"
	TypeDistrict AccountBlockType = "D"
	TypeCity     AccountBlockType = "C"
	TypeBranch   AccountBlockType = "B"
)

type AccountBlock struct {
	ID        bson.ObjectID    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string           `bson:"name" json:"name"`
	Code      string           `bson:"code" json:"code"`
	Address   string           `bson:"address" json:"address"`
	ParentID  *bson.ObjectID   `bson:"parent_id,omitempty" json:"parent_id,omitempty"`
	Parent    *AccountBlock    `bson:"parent,omitempty" json:"parent,omitempty"`
	Slug      string           `bson:"slug" json:"slug"`
	Type      AccountBlockType `bson:"type" json:"type"`
	IsEnabled bool             `bson:"is_enabled" json:"is_enabled"`
	IsDeleted bool             `bson:"is_deleted,omitempty" json:"-"`
	CreatedAt time.Time        `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time        `bson:"updated_at" json:"updated_at"`
}

type Branch struct {
	ID                    bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	BranchCode            string        `bson:"branch_code" json:"branch_code"`
	BranchName            string        `bson:"branch_name" json:"branch_name"`
	BranchAddress         string        `bson:"branch_address" json:"branch_address"`
	DistrictCode          string        `bson:"district_code" json:"district_code"`
	DistrictName          string        `bson:"district_name" json:"district_name"`
	RegionName            string        `bson:"region_name" json:"region_name"`
	RecordStat            string        `bson:"record_stat" json:"record_stat"`
	Enabled               bool          `bson:"enabled" json:"enabled"`
	EnableOrDisableReason string        `bson:"enable_or_disable_reason,omitempty" json:"enable_or_disable_reason,omitempty"`
	IsDeleted             bool          `bson:"is_deleted,omitempty" json:"-"`
	CreatedAt             time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt             time.Time     `bson:"updated_at" json:"updated_at"`
}

type Region struct {
	ID                    bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	RegionCode            string        `json:"region_code" bson:"region_code"`
	RegionName            string        `json:"region_name" bson:"region_name"`
	RegionAddress         string        `json:"region_address" bson:"region_address"`
	Enabled               bool          `json:"enabled" bson:"enabled"`
	EnableOrDisableReason string        `bson:"enable_or_disable_reason,omitempty" json:"enable_or_disable_reason,omitempty"`
	IsDeleted             bool          `bson:"is_deleted,omitempty" json:"-"`
	CreatedAt             time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt             time.Time     `json:"updated_at" bson:"updated_at"`
}

type District struct {
	ID                    bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	DistrictCode          string        `json:"district_code" bson:"district_code"`
	DistrictName          string        `json:"district_name" bson:"district_name"`
	DistrictAddress       string        `json:"district_address" bson:"district_address"`
	RegionID              string        `json:"region_id" bson:"region_id"`
	RegionName            string        `json:"region_name" bson:"region_name"`
	Enabled               bool          `json:"enabled" bson:"enabled"`
	EnableOrDisableReason string        `bson:"enable_or_disable_reason,omitempty" json:"enable_or_disable_reason,omitempty"`
	IsDeleted             bool          `bson:"is_deleted,omitempty" json:"-"`
	CreatedAt             time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt             time.Time     `json:"updated_at" bson:"updated_at"`
}

type City struct {
	ID                    bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CityCode              string        `json:"city_code" bson:"city_code,omitempty"`
	CityName              string        `json:"city_name" bson:"city_name,omitempty"`
	CityAddress           string        `json:"city_address" bson:"city_address,omitempty"`
	DistrictID            string        `json:"district_id" bson:"district_id,omitempty"`
	DistrictName          string        `json:"district_name" bson:"district_name,omitempty"`
	RegionID              string        `json:"region_id" bson:"region_id,omitempty"`
	RegionName            string        `json:"region_name" bson:"region_name,omitempty"`
	Enabled               bool          `json:"enabled" bson:"enabled"`
	EnableOrDisableReason string        `bson:"enable_or_disable_reason,omitempty" json:"enable_or_disable_reason,omitempty"`
	IsDeleted             bool          `bson:"is_deleted,omitempty" json:"-"`
	CreatedAt             time.Time     `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt             time.Time     `json:"updated_at" bson:"updated_at,omitempty"`
}
