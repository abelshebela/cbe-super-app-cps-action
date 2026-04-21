package model

import (
	"time"
)

type CustomerRoleInfo struct {
	ID   string `json:"id" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name"`
}

type CustSegment struct {
	Id                 string `json:"id,omitempty" bson:"id"`
	CustomerGroup      string `json:"cust_group" bson:"cust_group"`
	CustomerSegment    string `json:"cust_segment" bson:"cust_segment"`
	CustomerSubSegment string `json:"cus_sub_segment" bson:"cust_sub_segment"`
}

type CustomerSegmentation struct {
	ID                        string           `json:"id" bson:"_id,omitempty"`
	CustomerRole              CustomerRoleInfo `json:"customer_role" bson:"customer_role"`
	RemovedCustomerSegmentIDs []string         `json:"removed_customer_segment_ids,omitempty"`
	CustomerSegments          []CustSegment    `json:"t24_customer_sub_segments" bson:"t24_customer_sub_segments"`
	IsEnabled                 bool             `json:"is_enabled" bson:"is_enabled"`
	IsDeleted                 bool             `json:"is_deleted" bson:"is_deleted"`
	CreatedAt                 time.Time        `json:"created_at" bson:"created_at"`
	UpdatedAt                 time.Time        `json:"updated_at" bson:"updated_at"`
}
