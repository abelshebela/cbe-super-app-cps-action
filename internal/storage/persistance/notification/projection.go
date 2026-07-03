package notification

import (
	local_model "cbe-super-app-cps-action/internal/constants/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// NotificationMapper maps a Notification model to a bson.M for updates
func NotificationMapper(notification local_model.NotificationDocument) bson.M {
	// 	 Title             string                      `json:"title" bson:"title"`
	// Message           string                      `json:"message" bson:"message"`
	// BroadcastCategory constants.BroadcastCategory `json:"category" bson:"category"`
	// BroadcastType     constants.BroadcastType     `json:"type" bson:"type"`
	// ImageURL          string                      `json:"image_url,omitempty" bson:"image_url,omitempty"`
	// ActionURL         string                      `json:"action_url,omitempty" bson:"action_url,omitempty"`
	// Data              map[string]interface{}      `json:"data,omitempty" bson:"data,omitempty"`
	// CreatedAt         time.Time                   `json:"created_at" bson:"created_at"`
	// UpdatedAt         time.Time                   `json:"updated_at" bson:"updated_at"`
	// ExpiresAt         *time.Time                  `json:"expires_at,omitempty" bson:"expires_at,omitempty"`
	// KafkaMessageID    string                      `json:"kafka_message_id,omitempty" bson:"kafka_message_id,omitempty"`

	// return bson.M{
	// 	"title":              notification.Title,
	// 	"notification_code":  notification.NotificationCode,
	// 	"notification_type":  notification.NotificationType,
	// 	"notification_body":  notification.NotificationBody,
	// 	"is_public":          notification.IsPublic,
	// 	"for":                notification.For,
	// 	"created_by":         notification.CreatedBy,
	// 	"notification_parts": notification.NotificationParts,
	// 	"seen":               notification.Seen,
	// 	"enabled":            notification.Enabled,
	// 	"status":             notification.Status,
	// 	"last_modified":      time.Now(),
	// }

	return bson.M{
		"title":    notification.Title,
		"message":  notification.Message,
		"category": notification.Category,
		"type":     notification.BroadcastType,
	}
}

// func ToNotificationDocument(domain *model.Notification) (*model.NotificationDocument, error) {
// 	var ID bson.ObjectID

// 	if !domain.ID.IsZero() {
// 		ID = domain.ID
// 	} else {
// 		ID = bson.NewObjectID()
// 	}

// 	return &model.NotificationDocument{
// 		ID:                ID,
// 		NotificationCode:  domain.NotificationCode,
// 		NotificationType:  domain.NotificationType,
// 		NotificationBody:  domain.NotificationBody,
// 		IsPublic:          domain.IsPublic,
// 		For:               string(domain.For),
// 		CreatedBy:         domain.CreatedBy,
// 		NotificationParts: domain.NotificationParts,
// 		Seen:              domain.Seen,
// 		Title:             domain.Title,
// 		Status:            string(domain.Status),
// 		Enabled:           domain.Enabled,
// 		IsDeleted:         domain.IsDeleted,
// 		CreatedAt:         domain.CreatedAt,
// 		LastModified:      domain.LastModified,
// 	}, nil
// }
