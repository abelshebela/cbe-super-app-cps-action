package ad

import (
	"cbe-super-app-cps-action/internal/constants/model"
)

// generateAdvert creates a copy of an advert
func GenerateAdvert(advert model.Advert) *model.Advert {
	return &model.Advert{
		ID:            advert.ID,
		Title:         advert.Title,
		Description:   advert.Description,
		BannerImage:   advert.BannerImage,
		AdvertFor:     advert.AdvertFor,
		Date:          advert.Date,
		Enabled:       advert.Enabled,
		IsDeleted:     advert.IsDeleted,
		CreatedAt:     advert.CreatedAt,
		LastUpdatedAt: advert.LastUpdatedAt,
		DeletedAt:     advert.DeletedAt,
	}
}
