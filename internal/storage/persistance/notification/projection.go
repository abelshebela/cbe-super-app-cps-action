package notification

import (
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// NotificationMapper maps a Notification model to a bson.M for updates
func NotificationMapper(notification model.Notification) bson.M {
	return bson.M{
		"title":              notification.Title,
		"notification_code":  notification.NotificationCode,
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
	}
}

func ToNotificationDocument(domain *model.Notification) (*model.NotificationDocument, error) {
	var ID bson.ObjectID

	if !domain.ID.IsZero() {
		ID = domain.ID
	} else {
		ID = bson.NewObjectID()
	}

	return &model.NotificationDocument{
		ID:                ID,
		NotificationCode:  domain.NotificationCode,
		NotificationType:  domain.NotificationType,
		NotificationBody:  domain.NotificationBody,
		IsPublic:          domain.IsPublic,
		For:               string(domain.For),
		CreatedBy:         domain.CreatedBy,
		NotificationParts: domain.NotificationParts,
		Seen:              domain.Seen,
		Title:             domain.Title,
		Status:            string(domain.Status),
		Enabled:           domain.Enabled,
		IsDeleted:         domain.IsDeleted,
		CreatedAt:         domain.CreatedAt,
		LastModified:      domain.LastModified,
	}, nil
}
