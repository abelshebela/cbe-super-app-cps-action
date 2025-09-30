package core

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/dto/ad"
	// "cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"time"
	// local_util "cbe-super-app-cps-action/pkgs/utils"
)

// ToDomainAdvertRequest converts an HTTP AdvertRequest to a domain-level dto.AdvertRequest
func ToAdvert(httpRequest ad.AdvertRequest) (model.Advert, error) {
	advertFor := constants.AdvertFor(httpRequest.AdvertFor)

	return model.Advert{
		Title:       httpRequest.Title,
		Description: httpRequest.Description,
		AdvertFor:   advertFor,
		Date: types.AdvertDate{
			StartedAt: httpRequest.Date.StartedAt,
			ExpiredAt: httpRequest.Date.ExpiredAt,
		},
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
		Enabled:       false,
		IsDeleted:     false,
	}, nil
}

// ToAdvertResponse converts an entity.Advert to a dto.AdvertResponse
func ToAdvertResponse(advert model.Advert) ad.AdvertResponse {
	return ad.AdvertResponse{
		ID:            advert.ID.Hex(),
		Title:         advert.Title,
		Description:   advert.Description,
		BannerImage:   advert.BannerImage,
		AdvertFor:     advert.AdvertFor,
		Date:          ad.AdvertDate{StartedAt: advert.Date.StartedAt, ExpiredAt: advert.Date.ExpiredAt},
		Enabled:       advert.Enabled,
		CreatedAt:     advert.CreatedAt,
		LastUpdatedAt: advert.LastUpdatedAt,
	}
}

// ToAdvertResponses converts a slice of entity.Advert to a slice of dto.AdvertResponse
func ToAdvertResponses(adverts []*model.Advert) []*ad.AdvertResponse {
	responses := make([]*ad.AdvertResponse, len(adverts))
	for i, advert := range adverts {
		responses[i] = &ad.AdvertResponse{
			ID:            advert.ID.Hex(),
			Title:         advert.Title,
			Description:   advert.Description,
			BannerImage:   advert.BannerImage,
			AdvertFor:     advert.AdvertFor,
			Date:          ad.AdvertDate{StartedAt: advert.Date.StartedAt, ExpiredAt: advert.Date.ExpiredAt},
			Enabled:       advert.Enabled,
			CreatedAt:     advert.CreatedAt,
			LastUpdatedAt: advert.LastUpdatedAt,
		}
	}
	return responses
}
