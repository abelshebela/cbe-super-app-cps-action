package mappers

import (
	"fmt"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	entity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func ToAdvertDomain(ad model.Advert) *entity.Advert {
	return &entity.Advert{
		ID:            ad.ID.Hex(),
		Title:         ad.Title,
		Description:   ad.Description,
		BannerImage:   ad.BannerImage,
		AdvertFor:     entity.AdvertFor(ad.AdvertFor),
		Date:          entity.AdvertDate(ad.Date),
		Enabled:       ad.Enabled,
		IsDeleted:     ad.IsDeleted,
		DeletedAt:     ad.DeletedAt,
		CreatedAt:     ad.CreatedAt,
		LastUpdatedAt: ad.LastUpdatedAt,
	}
}

func ToAdvertDocument(ad *entity.Advert) (*model.Advert, error) {
	var objectID bson.ObjectID
	if ad.ID != "" {
		objID, err := bson.ObjectIDFromHex(ad.ID)
		if err != nil {
			return nil, fmt.Errorf("INVALID_ID")
		}
		objectID = objID

	} else {
		objectID = bson.NewObjectID()
	}

	return &model.Advert{
		ID:            objectID,
		Title:         ad.Title,
		Description:   ad.Description,
		BannerImage:   ad.BannerImage,
		AdvertFor:     model.AdvertFor(ad.AdvertFor),
		Date:          model.AdvertDate(ad.Date),
		Enabled:       ad.Enabled,
		IsDeleted:     ad.IsDeleted,
		DeletedAt:     ad.DeletedAt,
		CreatedAt:     ad.CreatedAt,
		LastUpdatedAt: ad.LastUpdatedAt,
	}, nil
}
