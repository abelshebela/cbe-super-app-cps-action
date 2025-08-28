package budget

import (
	"time"

	"cbe-super-app-cps-action/internal/constants/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func IconMapper(icon model.Icon) bson.M {
	result := bson.M{}

	if icon.Icon != "" {
		result["icon"] = icon.Icon
	}
	result["enabled"] = icon.Enabled
	result["is_deleted"] = icon.IsDeleted

	if !icon.CreatedAt.IsZero() {
		result["created_at"] = icon.CreatedAt
	}
	result["last_modified"] = time.Now()

	return bson.M{"$set": result}
}

func ColorMapper(color model.Color) bson.M {
	result := bson.M{}

	if color.Color != "" {
		result["color"] = color.Color
	}
	result["enabled"] = color.Enabled
	result["is_deleted"] = color.IsDeleted

	if !color.CreatedAt.IsZero() {
		result["created_at"] = color.CreatedAt
	}
	if !color.UpdatedAt.IsZero() {
		result["updated_at"] = color.UpdatedAt
	}
	result["last_modified"] = time.Now()

	return bson.M{"$set": result}
}
