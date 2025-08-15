package notification

import (
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type NotificationRequest struct {
	NotificationType string `json:"notification_type"`
	NotificationBody string `json:"notification_body"`
	IsPublic         bool   `json:"is_public"`
	For              string `json:"for"`
	CreatedBy        string `json:"created_by"`
	Title            string `json:"title"`
}

type NotificationResponse struct {
	ID                string    `json:"id"`
	Title             string    `json:"title"`
	NotificationType  string    `json:"notification_type"`
	NotificationBody  string    `json:"notification_body"`
	IsPublic          bool      `json:"is_public"`
	For               string    `json:"for"`
	CreatedBy         string    `json:"created_by"`
	NotificationParts any       `json:"notification_parts"`
	Seen              bool      `json:"seen"`
	Enabled           bool      `json:"enabled"`
	Status            string    `json:"status"`
	IsDeleted         bool      `json:"is_deleted"`
	CreatedAt         time.Time `json:"created_at"`
	LastModified      time.Time `json:"last_modified"`
	DeletedAt         time.Time `json:"deleted_at"`
}

var allowedNotificationFor = []string{"IFB", "CB", "ALL"}

func (r NotificationRequest) Validate(isCreate bool) error {
	enumValues := toInterfaceSlice(allowedNotificationFor)

	var fieldRules []*validation.FieldRules

	if isCreate {
		fieldRules = []*validation.FieldRules{
			validation.Field(&r.NotificationType, validation.Required.Error("notification_type is required")),
			validation.Field(&r.NotificationBody, validation.Required.Error("notification_body is required")),
			validation.Field(&r.Title, validation.Required.Error("title is required")),
			validation.Field(&r.For,
				validation.Required.Error("for is required"),
				validation.In(enumValues...).Error("invalid value for 'for'"),
			),
		}
	} else {
		fieldRules = []*validation.FieldRules{
			validation.Field(&r.NotificationType),
			validation.Field(&r.NotificationBody),
			validation.Field(&r.Title),
			validation.Field(&r.For,
				validation.In(enumValues...).Error("invalid value for 'for'"),
			),
		}
	}

	return validation.ValidateStruct(&r, fieldRules...)
}

func toInterfaceSlice(strs []string) []interface{} {
	res := make([]interface{}, len(strs))
	for i, s := range strs {
		res[i] = s
	}
	return res
}

func (r NotificationRequest) IsEmpty() bool {
	return strings.TrimSpace(r.NotificationType) == "" &&
		strings.TrimSpace(r.NotificationBody) == "" &&
		strings.TrimSpace(r.For) == "" &&
		strings.TrimSpace(r.Title) == ""
}
