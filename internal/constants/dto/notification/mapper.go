package notification

import (
	"cbe-super-app-cps-action/internal/constants/model"
)

func MapNotificationToResponse(entity *model.Notification) *NotificationResponse {
	return &NotificationResponse{
		ID:                entity.ID.Hex(),
		Title:             entity.Title,
		NotificationCode:  entity.NotificationCode,
		NotificationType:  entity.NotificationType,
		NotificationBody:  entity.NotificationBody,
		IsPublic:          entity.IsPublic,
		For:               string(entity.For),
		CreatedBy:         entity.CreatedBy,
		NotificationParts: entity.NotificationParts,
		Seen:              entity.Seen,
		Status:            string(entity.Status),
		Enabled:           entity.Enabled,
		IsDeleted:         entity.IsDeleted,
		CreatedAt:         entity.CreatedAt,
		LastModified:      entity.LastModified,
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
