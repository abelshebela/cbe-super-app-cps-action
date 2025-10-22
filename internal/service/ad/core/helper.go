package core

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
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

func DuplicateAdvertChecker(ctx context.Context, advert model.Advert, repoAd storage.AdvertRepository, isCreate bool, id string) (bool, error) {
	existingAdverts, err := repoAd.FindAllWithPagination(ctx, types.Filter{
		Search: advert.Title,
		Filters: map[string]interface{}{
			"is_deleted": false,
			"created_at": map[string]interface{}{"$gte": advert.CreatedAt},
			"title":      advert.Title,
		},
	})
	if err != nil {
		return false, err
	}
	if existingAdverts.Meta.TotalDocs > 0 {
		return true, nil
	}
	if !isCreate && id != "" {
		prevAdvert, err := repoAd.FindByID(ctx, id)
		if err != nil {
			return false, err
		}
		if prevAdvert.Title == advert.Title {
			return false, errors.New(localization.ErrorAdvertTitleNotChanged.Code)
		}
	}
	return false, nil
}
