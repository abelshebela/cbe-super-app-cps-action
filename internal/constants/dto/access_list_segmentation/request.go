package access_list_segmentation_dto

import "time"

type CreateAccessListSegmentationRequest struct {
	AccessListKeys []string `json:"access_list_keys" bson:"access_list_keys" validate:"required"` //the service which is being segmented
	SegmentationID string   `json:"segmentation_id" bson:"segmentation_id" validate:"required"`   //id of the region, district, country,branch
	SegmentType    string   `json:"segment_type" bson:"segment_type" validate:"required,oneof=Account Block"`
	Type           string   `json:"type" bson:"type" validate:"required,oneof=R D C B"`
	Reason         string   `json:"reason" bson:"reason" validate:"required"`
}

type EnableDisableAccessListSegmentationRequest struct {
	AccessListKeys   []string `json:"access_list_keys" bson:"access_list_keys" validate:"required"`
	SegmentationType string   `json:"segmentation_type" bson:"segmentation_type" validate:"required,oneof=Account Block"`
	Reason           string   `json:"reason" bson:"reason" validate:"required"`
}

type BulkDisableAccessListSegmentationRequest struct {
	ID               string   `json:"segment_key" bson:"segment_key" validate:"required"`
	Keys             []string `json:"keys" bson:"keys" validate:"required"`
	Enabled          bool     `json:"enabled" bson:"enabled"`
	SegmentationType string   `json:"segmentation_type" bson:"segmentation_type" validate:"required,oneof=Account Block"`
}

type UpdateAccessListSegmentationRequest struct {
	ID                string `json:"id" bson:"id"`
	NewAccessListKey  string `json:"new_access_list_key" bson:"new_access_list_key"`
	NewAccessListName string `json:"new_access_list_name" bson:"new_access_list_name"`
	NewSegmentationID string `json:"new_segmentation_id" bson:"new_segmentation_id"`
	Type              string `json:"type" bson:"type" validate:"oneof=R D C U B"`
	SegmentType       string `json:"segment_type" bson:"segment_type"`
	Reason            string `json:"reason" bson:"reason" validate:"required"`
}

type AccessListSegmentationResponse struct {
	ID             string    `json:"id" bson:"id"`
	AccessListKey  string    `json:"access_list_key" bson:"access_list_key"`
	AccessListName string    `json:"access_list_name" bson:"access_list_name"`
	Reason         string    `json:"reason" bson:"reason"`
	Type           string    `json:"type" bson:"type"`
	SegmentType    string    `json:"segment_type" bson:"segment_type"`
	SegmentCode    string    `json:"segment_code" bson:"segment_code"`
	SegmentName    string    `json:"segment_name" bson:"segment_name"`
	SegmentedID    string    `json:"segmented_id" bson:"segmented_id"`
	Enabled        bool      `json:"enabled" bson:"enabled"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" bson:"updated_at"`
}
