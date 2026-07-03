package core

import (
	// "cbe-super-app-cps-action/internal/constants"
	notify "cbe-super-app-cps-action/internal/constants/dto/notification"

	local_model "cbe-super-app-cps-action/internal/constants/model"

	notification_constants "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/notification/constants"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
)

func GenerateNotification(notification model.Notification) *model.Notification {
	return &model.Notification{
		ID:                notification.ID,
		NotificationType:  notification.NotificationType,
		NotificationBody:  notification.NotificationBody,
		IsPublic:          notification.IsPublic,
		For:               notification.For,
		CreatedBy:         notification.CreatedBy,
		NotificationParts: notification.NotificationParts,
		Seen:              notification.Seen,
		Enabled:           notification.Enabled,
		IsDeleted:         notification.IsDeleted,
		CreatedAt:         notification.CreatedAt,
		LastModified:      notification.LastModified,
		DeletedAt:         notification.DeletedAt,
	}
}

func BuildCreateNotification(req notify.NotificationRequest) local_model.NotificationDocument {
	return local_model.NotificationDocument{
		Title:            req.Title,
		NotificationBody: req.NotificationBody,
		Category:         notification_constants.BroadcastCategoryOther.String(),
		For:              notification_constants.BroadcastType(req.NotificationType).String(),
	}
}

func BuildUpdateNotification(prev *local_model.NotificationDocument, req notify.NotificationRequest) local_model.NotificationDocument {
	cur := *prev
	if req.Title != "" {
		cur.Title = req.Title
	}
	if req.NotificationBody != "" {
		cur.NotificationBody = req.NotificationBody
	}
	if req.For != "" {
		cur.For = notification_constants.BroadcastType(req.For).String()
	}
	return cur
}
