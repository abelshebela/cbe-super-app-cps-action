package avatar

import (
	"cbe-super-app-cps-action/internal/constants/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// AvatarMapper maps Avatar model to BSON for database operations
func AvatarMapper(data model.Avatar) bson.M {
	result := bson.M{}
	if data.Label != "" {
		result["lable"] = data.Label
	}
	if data.Avatar != "" {
		result["description"] = data.Avatar
	}

	result["enabled"] = data.Enable
	return result
}
