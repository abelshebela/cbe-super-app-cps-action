package color

import (
	"cbe-super-app-cps-action/internal/constants/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// ColorMapper maps Color model to BSON for database operations
func ColorMapper(data model.Color) bson.M {
	result := bson.M{}
	if data.Color != "" {
		result["color"] = data.Color
	}
	result["enabled"] = data.Enabled

	return result
}
