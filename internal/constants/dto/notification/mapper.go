package notification

import (
	local_model "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/model"
)

func MapNotificationToResponse(entity *local_model.NotificationDocument) *NotificationResponse {
	return &NotificationResponse{
		Title:            entity.Title,
		NotificationType: entity.For,
		NotificationBody: entity.NotificationBody,
		For:              string(entity.For),
	}
}

func MapNotificationRequestToDomain(req *NotificationRequest) *NotificationRequest {
	if req == nil {
		return nil
	}

	return &NotificationRequest{
		NotificationType: req.For,
		NotificationBody: req.NotificationBody,
		IsPublic:         req.IsPublic,
		For:              req.For,
		CreatedBy:        req.CreatedBy,
		Title:            req.Title,
	}
}
