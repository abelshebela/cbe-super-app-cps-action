package notification

import "time"

type NotificationStatus string

const (
	StatusPending NotificationStatus = "PENDING"
	StatusSent    NotificationStatus = "SENT"
	StatusSeen    NotificationStatus = "SEEN"
)

type NotificationFor string

const (
	ForIFB NotificationFor = "IFB"
	ForCB  NotificationFor = "CB"
	ForAll NotificationFor = "ALL"
)

type Notification struct {
	ID                string             `json:"id,omitempty" bson:"_id"`
	Title             string             `json:"title" bson:"title"`
	NotificationType  string             `json:"notification_type" bson:"notification_type"`
	NotificationBody  string             `json:"notification_body" bson:"notification_body"`
	IsPublic          bool               `json:"is_public" bson:"is_public"`
	For               NotificationFor    `json:"for" bson:"for"`
	CreatedBy         string             `json:"created_by" bson:"created_by"`
	NotificationParts any                `json:"notification_parts" bson:"notification_parts"`
	Seen              bool               `json:"seen" bson:"seen"`
	Enabled           bool               `json:"enabled" bson:"enabled"`
	Status            NotificationStatus `json:"status" bson:"status"`
	IsDeleted         bool               `json:"is_deleted" bson:"is_deleted"`
	CreatedAt         time.Time          `json:"created_at" bson:"created_at"`
	LastModified      time.Time          `json:"last_modified" bson:"last_modified"`
	DeletedAt         time.Time          `json:"deleted_at" bson:"deleted_at"`
}
