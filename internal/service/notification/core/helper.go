package core

import (
	"encoding/json"

	// "cbe-super-app-cps-action/internal/constants"
	notify "cbe-super-app-cps-action/internal/constants/dto/notification"

	shared_notification "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/notification/dto"

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

func BuildCreateNotification(req notify.NotificationRequest) shared_notification.BroadcastInAppNotificationMessage {
	return shared_notification.BroadcastInAppNotificationMessage{
		Title:         req.Title,
		Message:       req.NotificationBody,
		Category:      notification_constants.BroadcastCategoryOther.String(),
		BroadcastType: notification_constants.BroadcastType(req.NotificationType).String(),
	}
}

func BuildUpdateNotification(prev *shared_notification.BroadcastInAppNotificationMessage, req notify.NotificationRequest) shared_notification.BroadcastInAppNotificationMessage {
	cur := *prev
	if req.Title != "" {
		cur.Title = req.Title
	}
	if req.NotificationType != "" {
		cur.BroadcastType = req.NotificationType
	}
	if req.NotificationBody != "" {
		cur.Message = req.NotificationBody
	}
	if req.For != "" {
		cur.BroadcastType = notification_constants.BroadcastType(req.For).String()
	}
	return cur
}

func BindNotificationFromAction(current any) (shared_notification.BroadcastInAppNotificationMessage, error) {
	var out shared_notification.BroadcastInAppNotificationMessage
	if current == nil {
		return out, nil
	}

	switch v := current.(type) {
	case shared_notification.BroadcastInAppNotificationMessage:
		return v, nil
	case *shared_notification.BroadcastInAppNotificationMessage:
		if v == nil {
			return out, nil
		}
		return *v, nil
	case string:
		if err := json.Unmarshal([]byte(v), &out); err != nil {
			return out, err
		}
		return out, nil
	case []byte:
		if err := json.Unmarshal(v, &out); err != nil {
			return out, err
		}
		return out, nil
	case map[string]any:
		b, err := json.Marshal(v)
		if err != nil {
			return out, err
		}
		if err := json.Unmarshal(b, &out); err != nil {
			return out, err
		}
		return out, nil
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return out, err
		}
		if err := json.Unmarshal(b, &out); err != nil {
			return out, err
		}
		return out, nil
	}
}
