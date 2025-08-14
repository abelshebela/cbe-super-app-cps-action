package ad

import (
	entity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad"
)

// ToDomainAdvertRequest converts an HTTP AdvertRequest to a domain-level dto.AdvertRequest
func ToDomainAdvertRequest(httpRequest AdvertRequest) (entity.AdvertRequest, error) {
	advertFor := entity.AdvertFor(httpRequest.AdvertFor)

	return entity.AdvertRequest{
		ID:          httpRequest.ID,
		Title:       httpRequest.Title,
		Description: httpRequest.Description,
		BannerImage: httpRequest.BannerImage,
		AdvertFor:   advertFor,
		Date: entity.AdvertDate{
			StartedAt: httpRequest.Date.StartedAt,
			ExpiredAt: httpRequest.Date.ExpiredAt,
		},
	}, nil
}

// ToAdvertResponse converts an entity.Advert to a dto.AdvertResponse
func ToAdvertResponse(advert entity.Advert) AdvertResponse {
	return AdvertResponse{
		ID:            advert.ID,
		Title:         advert.Title,
		Description:   advert.Description,
		BannerImage:   advert.BannerImage,
		AdvertFor:     advert.AdvertFor,
		Date:          AdvertDate{StartedAt: advert.Date.StartedAt, ExpiredAt: advert.Date.ExpiredAt},
		Enabled:       advert.Enabled,
		CreatedAt:     advert.CreatedAt,
		LastUpdatedAt: advert.LastUpdatedAt,
	}
}

// ToAdvertResponses converts a slice of entity.Advert to a slice of dto.AdvertResponse
func ToAdvertResponses(adverts []*entity.Advert) []*AdvertResponse {
	responses := make([]*AdvertResponse, len(adverts))
	for i, advert := range adverts {
		responses[i] = &AdvertResponse{
			ID:            advert.ID,
			Title:         advert.Title,
			Description:   advert.Description,
			BannerImage:   advert.BannerImage,
			AdvertFor:     advert.AdvertFor,
			Date:          AdvertDate{StartedAt: advert.Date.StartedAt, ExpiredAt: advert.Date.ExpiredAt},
			Enabled:       advert.Enabled,
			CreatedAt:     advert.CreatedAt,
			LastUpdatedAt: advert.LastUpdatedAt,
		}
	}
	return responses
}
