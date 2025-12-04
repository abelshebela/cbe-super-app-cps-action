package notification

import "time"

type NotificationResponse struct {
	ID                string    `json:"id"`
	Title             string    `json:"title"`
	IdNum             string    `json:"id_num"`
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
