package ad

import (
	"mime/multipart"
	"time"
)

// AdvertRequest represents the request payload for creating or updating an advert
type AdvertRequest struct {
	ID          string                `form:"id" json:"id"`
	Title       string                `form:"title" json:"title"`
	Description string                `form:"description" json:"description"`
	BannerImage *multipart.FileHeader `form:"banner_image" json:"banner_image"`
	AdvertFor   string                `form:"advert_for" json:"advert_for"`
	Date        AdvertDate            `form:"date" json:"date"`
}

// AdvertDate represents the start and end dates for an advert
type AdvertDate struct {
	StartedAt time.Time `form:"started_at" json:"started_at"`
	ExpiredAt time.Time `form:"expired_at" json:"expired_at"`
}
