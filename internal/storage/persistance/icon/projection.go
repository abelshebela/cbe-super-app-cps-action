package icon

import (
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// IconMapper maps an Icon model to a bson.M for updates
func IconMapper(icon model.Icon) bson.M {
	update := bson.M{
		"enabled":       icon.Enabled,
		"last_modified": time.Now(),
	}

	if icon.Icon != "" {
		update["icon"] = icon.Icon
	}

	return update
}
