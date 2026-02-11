package core

import (
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
)

// generateAdvert creates a copy of an advert
func GenerateAdvert(advert model.Advert) *model.Advert {
	return &model.Advert{
		ID:            advert.ID,
		Title:         advert.Title,
		Description:   advert.Description,
		BannerImage:   advert.BannerImage,
		AdvertFor:     advert.AdvertFor,
		Enabled:       advert.Enabled,
		IsDeleted:     advert.IsDeleted,
		CreatedAt:     advert.CreatedAt,
		LastUpdatedAt: advert.LastUpdatedAt,
		DeletedAt:     advert.DeletedAt,
	}
}
