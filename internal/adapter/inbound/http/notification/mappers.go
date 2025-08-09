package notification

import (
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/notification"
)

func MapNotificationToResponse(entity *entities.Notification) *NotificationResponse {
	return &NotificationResponse{
		ID:                entity.ID,
		Title:             entity.Title,
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

func MapNotificationRequestToDomain(req *NotificationRequest) *entities.NotificationRequest {
	if req == nil {
		return nil
	}

	return &entities.NotificationRequest{
		NotificationType: req.NotificationType,
		NotificationBody: req.NotificationBody,
		IsPublic:         req.IsPublic,
		For:              req.For,
		CreatedBy:        req.CreatedBy,
		Title:            req.Title,
	}
}
