package donation

import (
	"mime/multipart"

	"time"
)

type DonationImage struct {
	ID        string `json:"id" bson:"id"`
	PhotoURL  string `json:"photo_url" bson:"photo_url"`
	CreatedAt string `json:"created_at" bson:"created_at"`
}

type DonationCPSRequest struct {
	ID                  string          `json:"id" bson:"id"`
	DonationCode        string          `json:"donation_code,omitempty" bson:"donation_code,omitempty"`
	CompanyID           string          `json:"company_id" bson:"company_id"`
	CategoryID          string          `json:"category_id" bson:"category_id"`
	Title               string          `json:"title" bson:"title"`
	IsFeatured          bool            `json:"is_featured" bson:"is_featured"`
	Target              int32           `json:"target" bson:"target"`
	DonationDescription string          `json:"donation_description" bson:"donation_description"`
	DonationImages      []DonationImage `json:"donation_images,omitempty" bson:"donation_images,omitempty"`
	CoverImage          string          `json:"cover_image,omitempty" bson:"cover_image,omitempty"`
	EndDate             string          `json:"end_date" bson:"end_date"`
	StartDate           string          `json:"start_date" bson:"start_date"`
	ImageIDToDelete     string          `json:"image_id_to_delete,omitempty" bson:"image_id_to_delete,omitempty"`
	Enabled             bool            `json:"enabled" bson:"enabled"`
}

type DonationImageUpdateCPSRequest struct {
	ID       string `json:"id" bson:"id"`               // Donation ID
	ImageID  string `json:"image_id" bson:"image_id"`   // Image ID to update
	PhotoURL string `json:"photo_url" bson:"photo_url"` // New photo URL
}

type DonationRequest struct {
	DonationCode        string                  `json:"donation_code,omitempty" bson:"donation_code,omitempty"`
	CompanyID           string                  `json:"company_id" bson:"company_id"`
	CategoryID          string                  `json:"category_id" bson:"category_id"`
	Title               string                  `json:"title" bson:"title"`
	IsFeatured          bool                    `json:"is_featured" bson:"is_featured"`
	Target              int32                   `json:"target" bson:"target"`
	DonationDescription string                  `json:"donation_description" bson:"donation_description"`
	RemovedImages		[]string`json:"removed_images" bson:"donation_images"`
	DonationImages      []*multipart.FileHeader `json:"donation_images" bson:"donation_images"`
	CoverImage          *multipart.FileHeader   `json:"cover_image,omitempty" bson:"cover_image,omitempty"`
	EndDate             time.Time               `json:"end_date" bson:"end_date"`
	StartDate           time.Time               `json:"start_date" bson:"start_date"`
	Enabled             bool                    `json:"enabled" bson:"enabled"`
}

type DonationImageUpdateRequest struct {
	ImageID string                `json:"image_id" bson:"image_id"` // Image ID to update
	Image   *multipart.FileHeader `json:"image" bson:"image"`       // New image file
}

type DonationImageDeleteRequest struct {
	ImageID string `json:"image_id"`
}
type DonationImageDeleteCPSRequest struct {
	ImageIDToDelete string `json:"image_id_to_delete" bson:"image_id_to_delete"`
}

type DonationImageAddRequest struct {
	DonationImages []DonationImage `json:"donation_images,omitempty" bson:"donation_images,omitempty"`
}
type EnableDonationRequest struct {
	Enabled bool `json:"enabled" bson:"enabled"`
}

