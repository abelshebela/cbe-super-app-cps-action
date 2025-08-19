package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type CPSAction struct {
	ID                 bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ActionCode         string        `bson:"action_code" json:"action_code"`
	UniqueId           string        `bson:"unique_id" json:"unique_id,omitempty"`
	MakerID            string        `bson:"maker_id" json:"maker_id"`
	MakerName          string        `bson:"maker_name" json:"maker_name"`
	MakerPhoneNumber   string        `bson:"maker_phone_number" json:"maker_phone_number"`
	CheckerID          string        `bson:"checker_id" json:"checker_id,omitempty"`
	CheckerName        string        `bson:"checker_name" json:"checker_name,omitempty"`
	CheckerPhoneNumber string        `bson:"checker_phone_number" json:"checker_phone_number,omitempty"`
	Department         string        `bson:"department" json:"department"`
	RejectionReason    string        `bson:"rejection_reason" json:"rejection_reason,omitempty"`
	PreviousAction     interface{}   `bson:"previous_action" json:"previous_action,omitempty"`
	CurrentAction      interface{}   `bson:"current_action" json:"current_action,omitempty"`
	ActionStatus       string        `bson:"action_status" json:"action_status,omitempty"`
	ActionType         string        `bson:"action_type" json:"action_type,omitempty"`
	IsDeleted          bool          `bson:"is_deleted" json:"is_deleted,omitempty"`
	RequestAction      string        `bson:"request_action" json:"request_action"`
	// Branch fields
	BranchCodes        []string      `bson:"branch_codes" json:"branch_codes,omitempty"`
	BranchNames        []string      `bson:"branch_names" json:"branch_names,omitempty"`
	// City fields
	CityCodes          []string      `bson:"city_codes" json:"city_codes,omitempty"`
	CityNames          []string      `bson:"city_names" json:"city_names,omitempty"`
	// Region fields
	RegionCodes        []string      `bson:"region_codes" json:"region_codes,omitempty"`
	RegionNames        []string      `bson:"region_names" json:"region_names,omitempty"`
	// District fields
	DistrictCodes      []string      `bson:"district_codes" json:"district_codes,omitempty"`
	DistrictNames      []string      `bson:"district_names" json:"district_names,omitempty"`
	// Timestamps
	CreatedAt          time.Time     `bson:"created_at" json:"created_at,omitempty"`
	LastModifiedAt     time.Time     `bson:"last_modified_at" json:"last_modified_at,omitempty"`
	MakerActionTime    time.Time     `bson:"maker_action_time" json:"maker_action_time,omitempty"`
	CheckerActionTime  *time.Time    `bson:"checker_action_time" json:"checker_action_time,omitempty"`
	// Additional fields
	Reason             string        `bson:"reason" json:"reason,omitempty"`
	Status             string        `bson:"status" json:"status,omitempty"`
	LastUpdated        time.Time     `bson:"last_updated" json:"last_updated,omitempty"`
}
