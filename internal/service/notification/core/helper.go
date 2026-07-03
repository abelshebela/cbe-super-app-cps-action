package core

import (
	"encoding/json"
	"strings"

	// "cbe-super-app-cps-action/internal/constants"
	notify "cbe-super-app-cps-action/internal/constants/dto/notification"

	local_model "cbe-super-app-cps-action/internal/constants/model"

	shared_notification "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/notification/dto"

	notification_constants "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/notification/constants"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"go.mongodb.org/mongo-driver/v2/bson"
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

func BindNotificationFromAction(current any) (shared_notification.BroadcastInAppNotificationMessage, error) {
	out := shared_notification.BroadcastInAppNotificationMessage{}
	if current == nil {
		return out, nil
	}

	values, err := actionValueMap(current)
	if err != nil {
		return out, err
	}

	out.Title = stringFromAny(values, "title")
	out.Message = stringFromAny(values, "message")
	out.Category = stringFromAny(values, "category")
	out.BroadcastType = stringFromAny(values, "type", "broadcast_type", "broadcasttype", "broadcasetype", "broadcaseType")

	return out, nil
}

func actionValueMap(current any) (map[string]any, error) {
	switch v := current.(type) {
	case shared_notification.BroadcastInAppNotificationMessage:
		return map[string]any{
			"title":         v.Title,
			"message":       v.Message,
			"category":      v.Category,
			"type":          v.BroadcastType,
			"image_url":     v.ImageURL,
			"action_url":    v.ActionURL,
			"data":          v.Data,
			"expires_at":    v.ExpiresAt,
			"broadcasttype": v.BroadcastType,
		}, nil
	case *shared_notification.BroadcastInAppNotificationMessage:
		if v == nil {
			return map[string]any{}, nil
		}
		return actionValueMap(*v)
	case bson.M:
		return normalizeMap(map[string]any(v)), nil
	case bson.D:
		out := make(map[string]any, len(v))
		for _, elem := range v {
			out[elem.Key] = normalizeAny(elem.Value)
		}
		return normalizeMap(out), nil
	case map[string]any:
		return normalizeMap(v), nil
	case string:
		var out map[string]any
		if err := json.Unmarshal([]byte(v), &out); err != nil {
			return nil, err
		}
		return normalizeMap(out), nil
	case []byte:
		var out map[string]any
		if err := json.Unmarshal(v, &out); err != nil {
			return nil, err
		}
		return normalizeMap(out), nil
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		var out map[string]any
		if err := json.Unmarshal(b, &out); err != nil {
			return nil, err
		}
		return normalizeMap(out), nil
	}
}

func normalizeMap(in map[string]any) map[string]any {
	if in == nil {
		return nil
	}
	out := make(map[string]any, len(in))
	for key, value := range in {
		out[canonicalKey(key)] = normalizeAny(value)
	}
	return out
}

func normalizeAny(v any) any {
	switch val := v.(type) {
	case bson.M:
		return normalizeMap(map[string]any(val))
	case bson.D:
		out := make(map[string]any, len(val))
		for _, elem := range val {
			out[canonicalKey(elem.Key)] = normalizeAny(elem.Value)
		}
		return out
	case map[string]any:
		return normalizeMap(val)
	case []any:
		items := make([]any, len(val))
		for i, item := range val {
			items[i] = normalizeAny(item)
		}
		return items
	default:
		return val
	}
}

func canonicalKey(key string) string {
	key = strings.ToLower(strings.TrimSpace(key))
	key = strings.ReplaceAll(key, "_", "")
	key = strings.ReplaceAll(key, "-", "")
	key = strings.ReplaceAll(key, " ", "")
	return key
}

func stringFromAny(values map[string]any, keys ...string) string {
	if v, ok := lookupAny(values, keys...); ok {
		switch val := v.(type) {
		case string:
			return val
		case []byte:
			return string(val)
		default:
			return strings.TrimSpace(strings.Trim(fmtAny(val), "\""))
		}
	}
	return ""
}
func lookupAny(values map[string]any, keys ...string) (any, bool) {
	if len(values) == 0 {
		return nil, false
	}
	for _, key := range keys {
		candidate := canonicalKey(key)
		if v, ok := values[candidate]; ok {
			return v, true
		}
	}
	return nil, false
}

func fmtAny(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}
