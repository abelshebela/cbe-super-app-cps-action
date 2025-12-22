package model

import "time"

type CustomerSegmentation struct {
	ID                 string    `json:"id" bson:"_id,omitempty"`
	CustomerRole       string    `json:"customer_role" bson:"customer_role"`
	CustomerSegment    string    `json:"customer_segment" bson:"customer_segment"`
	CustomerSubSegment string    `json:"customer_sub_segment" bson:"customer_sub_segment"`
	CustomerGroup      string    `json:"customer_group" bson:"customer_group"`
	CreatedAt          time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" bson:"updated_at"`
}
