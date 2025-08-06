package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Branch struct {
	ID            bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	BranchCode    string        `json:"branchCode" bson:"branchCode"`
	BranchName    string        `json:"branchName" bson:"branchName"`
	BranchAddress string        `json:"branchAddress" bson:"branchAddress"`
	DistrictCode  string        `json:"districtCode" bson:"districtCode"`
	DistrictName  string        `json:"districtName" bson:"districtName"`
	BranchRegion  string        `json:"branchRegion" bson:"branchRegion"`
	RecordStat    string        `json:"RecordStat" bson:"RecordStat"`
	CreatedAt     time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt" bson:"updatedAt"`
	Version       int           `json:"__v" bson:"__v"`
	Enabled       bool          `json:"enabled" bson:"enabled"`
}
