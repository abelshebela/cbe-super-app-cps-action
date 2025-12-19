package access_list_segmentation_dto

import "time"

type CreateAccessListSegmentationRequest struct {
	ServiceID   string   `json:"service_id" bson:"service_id" validate:"required"`
	ServiceName string   `json:"service_name" bson:"service_name" validate:"required"`
	SegmentedID []string `json:"segmented_id" bson:"segmented_id" validate:"required"`
	Type        string   `json:"type" bson:"type" validate:"required,oneof=R D C U"`
}

type UpdateAccessListSegmentationRequest struct {
	ID             string `json:"id" bson:"id"`
	NewServiceID   string `json:"new_service_id" bson:"new_service_id"`
	NewServiceName string `json:"new_service_name" bson:"new_service_name"`
	NewSegmentedID string `json:"new_segmented_id" bson:"new_segmented_id"`
	Type           string `json:"type" bson:"type" validate:"oneof=R D C U"`
}

type AccessListSegmentationResponse struct {
	ID          string    `json:"id" bson:"id"`
	ServiceID   string    `json:"service_id" bson:"service_id"`
	ServiceName string    `json:"service_name" bson:"service_name"`
	Type        string    `json:"type" bson:"type"`
	SegmentedID string    `json:"segmented_id" bson:"segmented_id"`
	Enabled     bool      `json:"enabled" bson:"enabled"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}
