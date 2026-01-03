package access_list_segmentation_dto

import "time"

type CreateAccessListSegmentationRequest struct {
	AccessListKeys  []string `json:"access_list_keys" bson:"access_list_keys" validate:"required"`
	SegmentType     string   `json:"segment_type" bson:"segment_type" validate:"required,oneof=R D C U B"`
	SegmentCode     string   `json:"segment_code" bson:"segment_code" validate:"required"`
	SegmentName     string   `json:"segment_name" bson:"segment_name" validate:"required"`
	AccessListNames []string `json:"access_list_names" bson:"access_list_names" validate:"required"`
	SegmentedID     string   `json:"segmented_id" bson:"segmented_id" validate:"required"`
	Type            string   `json:"type" bson:"type" validate:"required,oneof=R D C U"`
}

type UpdateAccessListSegmentationRequest struct {
	ID                string `json:"id" bson:"id"`
	NewAccessListKey  string `json:"new_access_list_key" bson:"new_access_list_key"`
	NewAccessListName string `json:"new_access_list_name" bson:"new_access_list_name"`
	NewSegmentedID    string `json:"new_segmented_id" bson:"new_segmented_id"`
	Type              string `json:"type" bson:"type" validate:"oneof=R D C U B"`
	SegmentType       string `json:"segment_type" bson:"segment_type"`
	SegmentCode       string `json:"segment_code" bson:"segment_code"`
	SegmentName       string `json:"segment_name" bson:"segment_name"`
}

type AccessListSegmentationResponse struct {
	ID             string    `json:"id" bson:"id"`
	AccessListKey  string    `json:"access_list_key" bson:"access_list_key"`
	AccessListName string    `json:"access_list_name" bson:"access_list_name"`
	Type           string    `json:"type" bson:"type"`
	SegmentType    string    `json:"segment_type" bson:"segment_type"`
	SegmentCode    string    `json:"segment_code" bson:"segment_code"`
	SegmentName    string    `json:"segment_name" bson:"segment_name"`
	SegmentedID    string    `json:"segmented_id" bson:"segmented_id"`
	Enabled        bool      `json:"enabled" bson:"enabled"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" bson:"updated_at"`
}
