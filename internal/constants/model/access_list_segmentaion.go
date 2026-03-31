package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type AccessListSegmentation struct {
	ID               bson.ObjectID `json:"_id" bson:"_id"`
	Type             string        `json:"type" bson:"type"`
	SegmentationType string        `json:"segmentation_type" bson:"segmentation_type"`
	AccessListKey    string        `json:"access_list_key" bson:"access_list_key"`
	AccessListName   string        `json:"access_list_name" bson:"access_list_name"`
	SegmentedID      bson.ObjectID `json:"segmented_id" bson:"segmented_id"`
	SegmentationCode string        `json:"segmentation_code" bson:"segmentation_code"`
	SegmentationName string        `json:"segmentation_name" bson:"segmentation_name"`
	Enabled          bool          `json:"enabled" bson:"enabled"`
	CreatedAt        time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at" bson:"updated_at"`
	DeletedAt        time.Time     `json:"deleted_at" bson:"deleted_at"`
}

type AccessListSegmentationByCustomerRole struct {
	ID               string    `json:"_id" bson:"_id"`
	AccessListKey    string    `json:"access_list_key" bson:"access_list_key"`
	AccessListName   string    `json:"access_list_name" bson:"access_list_name"`
	Type             string    `json:"type" bson:"type"`                           //IFB,cbe
	SegmentationType string    `json:"segmentation_type" bson:"segmentation_type"` //account
	SegmentationCode string    `json:"segmentation_code" bson:"segmentation_code"` //customer role code eg. Mass retail
	SegmentationName string    `json:"segmentation_name" bson:"segmentation_name"` // name of the customer role eg. Mass
	Enabled          int       `json:"enabled" bson:"enabled"`
	CreatedAt        time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" bson:"updated_at"`
	DeletedAt        time.Time `json:"deleted_at" bson:"deleted_at"`
}
type AccessListSegmentationByGeography struct {
	ID               string    `json:"_id" bson:"_id"`
	AccessListKey    string    `json:"access_list_key" bson:"access_list_key"`
	AccessListName   string    `json:"access_list_name" bson:"access_list_name"`
	Type             string    `json:"type" bson:"type"`                           //B,D,R,C
	SegmentationType string    `json:"segmentation_type" bson:"segmentation_type"` //block
	SegmentedID      string    `json:"segmented_id" bson:"segmented_id"`           //block id, region id, city id, district id
	Enabled          int       `json:"enabled" bson:"enabled"`
	CreatedAt        time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" bson:"updated_at"`
	DeletedAt        time.Time `json:"deleted_at" bson:"deleted_at"`
}
