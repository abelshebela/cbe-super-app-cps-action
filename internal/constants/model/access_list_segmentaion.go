package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type AccessListSegmentation struct {
	ID               bson.ObjectID `json:"_id" bson:"_id"`
	Type             string        `json:"type" bson:"type"`
	SegmentationType string        `json:"segmentation_type" bson:"segmentation_type"`
	ServiceID        bson.ObjectID `json:"service_id" bson:"service_id"`
	ServiceName      string        `json:"service_name" bson:"service_name"`
	SegmentedID      bson.ObjectID `json:"segmented_id" bson:"segmented_id"`
	SegmentationCode string        `json:"segmentation_code" bson:"segmentation_code"`
	SegmentationName string        `json:"segmentation_name" bson:"segmentation_name"`
	Enabled          bool          `json:"enabled" bson:"enabled"`
	CreatedAt        time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at" bson:"updated_at"`
	DeletedAt        time.Time     `json:"deleted_at" bson:"deleted_at"`
}
