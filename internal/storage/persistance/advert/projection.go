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

	result["updated_at"] = time.Now()

	return result
}
