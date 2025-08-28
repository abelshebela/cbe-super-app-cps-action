package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Branch struct {
	ID            bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	BranchCode    string        `bson:"branch_code" json:"branch_code"`
	BranchName    string        `bson:"branch_name" json:"branch_name"`
	BranchAddress string        `bson:"branch_address" json:"branch_address"`
	DistrictCode  string        `bson:"district_code" json:"district_code"`
	DistrictName  string        `bson:"district_name" json:"district_name"`
	BranchRegion  string        `bson:"branch_region" json:"branch_region"`
	RecordStat    string        `bson:"record_stat" json:"record_stat"`
	CreatedAt     time.Time     `bson:"created_at" json:"created_at"`
	IsDeleted     bool          `bson:"is_deleted" json:"-"`
	UpdatedAt     time.Time     `bson:"updated_at" json:"updated_at"`
	Enabled       bool          `bson:"enabled" json:"enabled"`
}

type Region struct {
	ID            bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	RegionCode    string        `json:"region_code" bson:"region_code"`
	RegionName    string        `json:"region_name" bson:"region_name"`
	RegionAddress string        `json:"region_address" bson:"region_address"`
	CreatedAt     time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at" bson:"updated_at"`
	Enabled       bool          `json:"enabled" bson:"enabled"`
}

type District struct {
	ID              bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	DistrictCode    string        `json:"district_code" bson:"district_code"`
	DistrictName    string        `json:"district_name" bson:"district_name"`
	DistrictAddress string        `json:"district_address" bson:"district_address"`
	RegionID        string        `json:"region_id" bson:"region_id"`
	RegionName      string        `json:"region_name" bson:"region_name"`
	CreatedAt       time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at" bson:"updated_at"`
	Enabled         bool          `json:"enabled" bson:"enabled"`
}

type City struct {
	ID           bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CityCode     string        `json:"city_code" bson:"city_code"`
	CityName     string        `json:"city_name" bson:"city_name"`
	CityAddress  string        `json:"city_address" bson:"city_address"`
	DistrictID   string        `json:"district_id" bson:"district_id"`
	DistrictName string        `json:"district_name" bson:"district_name"`
	RegionID     string        `json:"region_id" bson:"region_id"`
	RegionName   string        `json:"region_name" bson:"region_name"`
	CreatedAt    time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at" bson:"updated_at"`
	Enabled      bool          `json:"enabled" bson:"enabled"`
}
