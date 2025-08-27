package ad

import (
	"cbe-super-app-cps-action/internal/constants"
	"time"
)

// AdvertResponse represents the response format for an advert
type AdvertResponse struct {
	ID            string           `json:"id"`
	Title         string           `json:"title"`
	Description   string           `json:"description"`
	BannerImage   string           `json:"banner_image"`
	AdvertFor     constants.AdvertFor `json:"advert_for"`
	Date          AdvertDate       `json:"date"`
	Enabled       bool             `json:"enabled"`
	CreatedAt     time.Time        `json:"created_at"`
	LastUpdatedAt time.Time        `json:"last_updated_at"`
}