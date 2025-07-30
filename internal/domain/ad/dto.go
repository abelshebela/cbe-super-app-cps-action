package ad

import (
	"mime/multipart"
)

type AdvertRequest struct {
	ID          string                `json:"id"`
	Title       string                `form:"title" json:"title"`
	Description string                `form:"description" json:"description"`
	BannerImage *multipart.FileHeader `form:"banner_image" json:"banner_image"`
	AdvertFor   AdvertFor             `form:"advert_for" json:"advert_for"`
	Date        AdvertDate            `form:"date" json:"date"`
}
