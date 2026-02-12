package ad

import (
	"mime/multipart"
)

// AdvertRequest represents the request payload for creating or updating an advert
type AdvertRequest struct {
	ID          string                `form:"id" json:"id"`
	Title       string                `form:"title" json:"title"`
	Description string                `form:"description" json:"description"`
	BannerImage *multipart.FileHeader `form:"banner_image" json:"banner_image"`
	AdvertFor   string                `form:"advert_for" json:"advert_for"`
}
