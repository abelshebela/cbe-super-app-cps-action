package model

import (
	"cbe-super-app-cps-action/internal/constants"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Notification struct {
	ID                bson.ObjectID                `json:"id,omitempty" bson:"_id"`
	Title             string                       `json:"title" bson:"title"`
	NotificationType  string                       `json:"notification_type" bson:"notification_type"`
	NotificationBody  string                       `json:"notification_body" bson:"notification_body"`
	IsPublic          bool                         `json:"is_public" bson:"is_public"`
	For               constants.NotificationFor    `json:"for" bson:"for"`
	CreatedBy         string                       `json:"created_by" bson:"created_by"`
	NotificationParts any                          `json:"notification_parts" bson:"notification_parts"`
	Seen              bool                         `json:"seen" bson:"seen"`
	Enabled           bool                         `json:"enabled" bson:"enabled"`
	Status            constants.NotificationStatus `json:"status" bson:"status"`
	IsDeleted         bool                         `json:"is_deleted" bson:"is_deleted"`
	CreatedAt         time.Time                    `json:"created_at" bson:"created_at"`
	LastModified      time.Time                    `json:"last_modified" bson:"last_modified"`
	DeletedAt         time.Time                    `json:"deleted_at" bson:"deleted_at"`
}
