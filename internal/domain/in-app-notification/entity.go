package inappnotification

import "time"

type InAppNotification struct {
	ID                string
	UserID            string
	NotificationType  string
	NotificationBody  string
	IsPublic          bool
	For               string
	CreatedBy         string
	NotificationParts any
	Seen              bool
	Enabled           bool
	IsDeleted         bool
	CreatedAt         time.Time
	LastModified      time.Time
}
