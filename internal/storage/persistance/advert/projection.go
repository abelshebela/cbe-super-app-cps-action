package advert

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// AdvertMapper maps model.Advert to bson.M for MongoDB operations
func AdvertMapper(advert model.Advert) bson.M {
	result := bson.M{}

	if advert.Title != "" {
		result["title"] = advert.Title
	}
	if advert.Description != "" {
		result["description"] = advert.Description
	}

	if advert.BannerImage != "" {
		result["banner_image"] = advert.BannerImage
	}

	if advert.AdvertFor != "" {
		result["advert_for"] = advert.AdvertFor
	}

	// Use dot notation for partial updates on nested advert_date
	if !advert.Date.StartedAt.IsZero() {
		result["advert_date.started_at"] = advert.Date.StartedAt
	}

	if !advert.Date.ExpiredAt.IsZero() {
		result["advert_date.expired_at"] = advert.Date.ExpiredAt
	}

	if !advert.DeletedAt.IsZero() {
		result["deleted_at"] = advert.DeletedAt
	}

	result["enabled"] = advert.Enabled
	result["is_deleted"] = advert.IsDeleted
	result["last_updated_at"] = time.Now()

	// Wrap with $set to be used directly in UpdateOne
	return result
}
