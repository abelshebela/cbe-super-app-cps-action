package notification

import (
	shared_notification "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/notification/dto"
)

func MapNotificationToResponse(entity *shared_notification.BroadcastInAppNotificationMessage) *NotificationResponse {
	return &NotificationResponse{
		Title:            entity.Title,
		NotificationType: entity.BroadcastType,
		NotificationBody: entity.Message,
		For:              string(entity.BroadcastType),
	}
}

func MapNotificationRequestToDomain(req *NotificationRequest) *NotificationRequest {
	if req == nil {
		return nil
	}

	return &NotificationRequest{
		NotificationType: req.NotificationType,
		NotificationBody: req.NotificationBody,
		IsPublic:         req.IsPublic,
		For:              req.For,
		CreatedBy:        req.CreatedBy,
		Title:            req.Title,
	}
}
