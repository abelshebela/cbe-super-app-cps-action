package icon

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// IconMapper maps an Icon model to a bson.M for updates
func IconMapper(icon model.Icon) bson.M {
	return bson.M{
		"$set": bson.M{
			"icon":          icon.Icon,
			"enabled":       icon.Enabled,
			"last_modified": time.Now(),
		},
	}
}
