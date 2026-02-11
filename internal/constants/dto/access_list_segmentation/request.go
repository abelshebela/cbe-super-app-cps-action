package access_list_segmentation_dto

import "time"

type CreateAccessListSegmentationRequest struct {
	AccessListKeys  []string `json:"access_list_keys" bson:"access_list_keys" validate:"required"`   //the service which is being segmented NB: we actually work on a separate collection called access_list_segmentation rather than services
	AccessListNames []string `json:"access_list_names" bson:"access_list_names" validate:"required"` //the service names which is being segmented
	SegmentCode     string   `json:"segment_code" bson:"segment_code" validate:"required"`           //for type account or by account type
	SegmentName     string   `json:"segment_name" bson:"segment_name" validate:"required"`           //for type account or by account type
	SegmentedID     string   `json:"segmented_id" bson:"segmented_id" validate:"required"`           //id of the region, district, country,branch
	SegmentType     string   `json:"segment_type" bson:"segment_type" validate:"required,oneof=Account Block"`
	Type            string   `json:"type" bson:"type" validate:"required,oneof=R D C B"`
}

type EnableDisableAccessListSegmentationRequest struct {
	AccessListKeys []string `json:"access_list_keys" bson:"access_list_keys" validate:"required"`
}

type BulkDisableAccessListSegmentationRequest struct {
	ID   string   `json:"segment_key" bson:"segment_key" validate:"required"`
	Keys []string `json:"keys" bson:"keys" validate:"required"`
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
