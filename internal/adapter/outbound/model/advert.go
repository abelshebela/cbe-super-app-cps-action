package model

import (
	"fmt"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Advert struct {
	ID            bson.ObjectID `json:"id,omitempty" bson:"_id"`
	Title         string        `json:"title" bson:"title"`
	Description   string        `json:"description" bson:"description"`
	BannerImage   string        `json:"banner_image" bson:"banner_image"`
	AdvertFor     AdvertFor     `json:"advert_for" bson:"advert_for"`
	Date          AdvertDate    `json:"advert_date" bson:"advert_date"`
	Enabled       bool          `json:"enabled" bson:"enabled"`
	IsDeleted     bool          `json:"is_deleted" bson:"is_deleted"`
	DeletedAt     time.Time     `json:"deleted_at" bson:"deleted_at"`
	CreatedAt     time.Time     `json:"created_at" bson:"created_at"`
	LastUpdatedAt time.Time     `json:"last_updated_at" bson:"last_updated_at"`
}

func ToAdvertDomain(ad Advert) *entity.Advert {
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

func ToAdvertDocument(ad *entity.Advert) (Advert, error) {
	var objectID bson.ObjectID
	if ad.ID != "" {
		objID, err := bson.ObjectIDFromHex(ad.ID)
		if err != nil {
			return Advert{}, fmt.Errorf("INVALID_ID")
		}
		objectID = objID

	} else {
		objectID = bson.NewObjectID()
	}

	return Advert{
		ID:            objectID,
		Title:         ad.Title,
		Description:   ad.Description,
		BannerImage:   ad.BannerImage,
		AdvertFor:     AdvertFor(ad.AdvertFor),
		Date:          AdvertDate(ad.Date),
		Enabled:       ad.Enabled,
		IsDeleted:     ad.IsDeleted,
		DeletedAt:     ad.DeletedAt,
		CreatedAt:     ad.CreatedAt,
		LastUpdatedAt: ad.LastUpdatedAt,
	}, nil
}
