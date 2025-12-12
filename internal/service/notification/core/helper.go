package core

import (
	"encoding/json"
	"time"

	// "cbe-super-app-cps-action/internal/constants"
	notify "cbe-super-app-cps-action/internal/constants/dto/notification"
	local_utils "cbe-super-app-cps-action/pkgs/utils"

	shared_constants "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/constants"
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

func BuildCreateNotification(req notify.NotificationRequest) model.Notification {
	return model.Notification{
		Title:            req.Title,
		NotificationCode: local_utils.NewNotificationID(),
		NotificationType: req.NotificationType,
		NotificationBody: req.NotificationBody,
		IsPublic:         req.IsPublic,
		For:              shared_constants.NotificationFor(req.For),
		CreatedBy:        req.CreatedBy,
		Status:           shared_constants.StatusPending,
		Seen:             false,
		Enabled:          true,
		IsDeleted:        false,
		CreatedAt:        time.Now(),
		LastModified:     time.Now(),
	}
}

func BuildUpdateNotification(prev *model.Notification, req notify.NotificationRequest) model.Notification {
	cur := *prev
	if req.Title != "" {
		cur.Title = req.Title
	}
	if req.NotificationType != "" {
		cur.NotificationType = req.NotificationType
	}
	if req.NotificationBody != "" {
		cur.NotificationBody = req.NotificationBody
	}
	if req.For != "" {
		cur.For = shared_constants.NotificationFor(req.For)
	}
	if req.CreatedBy != "" {
		cur.CreatedBy = req.CreatedBy
	}
	// IsPublic only set to true explicitly; do not force false on partial updates
	if req.IsPublic {
		cur.IsPublic = true
	}
	cur.LastModified = time.Now()
	return cur
}

func BindNotificationFromAction(current any) (model.Notification, error) {
	var out model.Notification
	if current == nil {
		return out, nil
	}

	switch v := current.(type) {
	case model.Notification:
		return v, nil
	case *model.Notification:
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
