package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// type CustomerSegmentation struct {
// 	ID                 bson.ObjectID `json:"id" bson:"_id,omitempty"`
// 	CustomerRole       string        `json:"customer_role" bson:"customer_role"`
// 	CustomerSegment    string        `json:"customer_segment" bson:"customer_segment"`
// 	CustomerSubSegment string        `json:"customer_sub_segment" bson:"customer_sub_segment"`
// 	CustomerGroup      string        `json:"customer_group" bson:"customer_group"`
// 	IsDeleted          bool          `json:"is_deleted" bson:"is_deleted"`
// 	CreatedAt          time.Time     `json:"created_at" bson:"created_at"`
// 	UpdatedAt          time.Time     `json:"updated_at" bson:"updated_at"`
// }

type CustomerRoleInfo struct {
	ID   bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Name string        `json:"name" bson:"name"`
}

type CustomerSubSegments struct {
	Name            string `json:"name" bson:"name"`
	CustomerGroup   string `json:"cust_group" bson:"cust_group"`
	CustomerSegment string `json:"cust_segment" bson:"cust_segment"`
}

type CustomerSegmentation struct {
	ID                  bson.ObjectID         `json:"id" bson:"_id,omitempty"`
	CustomerRole        CustomerRoleInfo      `json:"customer_role" bson:"customer_role"`
	CustomerSubSegments []CustomerSubSegments `json:"t24_customer_sub_segments" bson:"t24_customer_sub_segments"`
	IsEnabled           bool                  `json:"is_enabled" bson:"is_enabled"`
	IsDeleted           bool                  `json:"is_deleted" bson:"is_deleted"`
	CreatedAt           time.Time             `json:"created_at" bson:"created_at"`
	UpdatedAt           time.Time             `json:"updated_at" bson:"updated_at"`
}
