package notification

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// NotificationMapper maps a Notification model to a bson.M for updates
func NotificationMapper(notification model.Notification) bson.M {
	return bson.M{
		"$set": bson.M{
			"title":              notification.Title,
			"notification_type":  notification.NotificationType,
			"notification_body":  notification.NotificationBody,
			"is_public":          notification.IsPublic,
			"for":                notification.For,
			"created_by":         notification.CreatedBy,
			"notification_parts": notification.NotificationParts,
			"seen":               notification.Seen,
			"enabled":            notification.Enabled,
			"status":             notification.Status,
			"last_modified":      time.Now(),
		},
	}
}
