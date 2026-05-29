package ad

import (
	"html"
	"mime/multipart"
	"strings"
)

// AdvertRequest represents the request payload for creating or updating an advert
type AdvertRequest struct {
	ID          string                `form:"id" json:"id"`
	Title       string                `form:"title" json:"title"`
	Description string                `form:"description" json:"description"`
	BannerImage *multipart.FileHeader `form:"banner_image" json:"banner_image"`
	AdvertFor   string                `form:"advert_for" json:"advert_for"`
}

func (r *AdvertRequest) Clean() {
	r.Title = strings.TrimSpace(r.Title)
	r.Description = strings.TrimSpace(r.Description)
	r.Title = html.EscapeString(r.Title)
	r.Description = html.EscapeString(r.Description)
}
