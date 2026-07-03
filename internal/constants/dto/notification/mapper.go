package notification

import (
	local_model "cbe-super-app-cps-action/internal/constants/model"
)

func MapNotificationToResponse(entity *local_model.NotificationDocument) *NotificationResponse {
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
		NotificationType: req.For,
		NotificationBody: req.NotificationBody,
		IsPublic:         req.IsPublic,
		For:              req.For,
		CreatedBy:        req.CreatedBy,
		Title:            req.Title,
	}
}
