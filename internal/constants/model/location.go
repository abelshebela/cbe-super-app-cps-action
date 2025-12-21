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
	ID         bson.ObjectID    `bson:"_id,omitempty" json:"id,omitempty"`
	Name       string           `bson:"name" json:"name"`
	Code       string           `bson:"code" json:"code"`
	Address    string           `bson:"address" json:"address"`
	ParentID   *bson.ObjectID   `bson:"parent_id,omitempty" json:"parent_id,omitempty"`
	Parent     *AccountBlock    `bson:"parent,omitempty" json:"parent,omitempty"`
	Slug       string           `bson:"slug" json:"slug"`
	Type       AccountBlockType `bson:"type" json:"type"`
	IsEnabled  bool             `bson:"is_enabled" json:"is_enabled"`
	CityID     string           `bson:"city_id,omitempty" json:"city_id,omitempty"`
	RegionID   string           `bson:"region_id,omitempty" json:"region_id,omitempty"`
	DistrictID string           `bson:"district_id,omitempty" json:"district_id,omitempty"`
	IsDeleted  bool             `bson:"is_deleted,omitempty" json:"-"`
	CreatedAt  time.Time        `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time        `bson:"updated_at" json:"updated_at"`
}
