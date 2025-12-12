package ad

import (
	shared_constant "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/constants"
	"time"
)

// AdvertResponse represents the response format for an advert
type AdvertResponse struct {
	ID            string              `json:"id"`
	Title         string              `json:"title"`
	Description   string              `json:"description"`
	BannerImage   string              `json:"banner_image"`
	AdvertFor     shared_constant.AdvertFor `json:"advert_for"`
	Enabled       bool                `json:"enabled"`
	CreatedAt     time.Time           `json:"created_at"`
	LastUpdatedAt time.Time           `json:"last_updated_at"`
}
