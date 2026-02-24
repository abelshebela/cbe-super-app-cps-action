package avatar

import (
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// AvatarMapper maps Avatar model to BSON for database operations
func AvatarMapper(data model.Avatar) bson.M {
	result := bson.M{}
	if data.Avatar != "" {
		result["avatar"] = data.Avatar
	}
	if data.Label != "" {
		result["label"] = data.Label
	}
	result["last_modified_at"] = time.Now()

	return result
}
