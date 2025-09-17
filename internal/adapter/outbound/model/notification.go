package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type NotificationDocument struct {
	ID                bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Title             string        `bson:"title,omitempty" json:"title"`
	NotificationType  string        `bson:"notification_type" json:"notification_type"`
	NotificationBody  string        `bson:"notification_body" json:"notification_body"`
	IsPublic          bool          `bson:"is_public" json:"is_public"`
	For               string        `bson:"for" json:"for"`
	CreatedBy         string        `bson:"created_by" json:"created_by"`
	NotificationParts any           `bson:"notification_parts" json:"notification_parts"`
	Seen              bool          `bson:"seen" json:"seen"`
	Status            string        `json:"status" bson:"status"`
	Enabled           bool          `bson:"enabled" json:"enabled"`
	IsDeleted         bool          `bson:"is_deleted" json:"is_deleted"`
	CreatedAt         time.Time     `bson:"created_at" json:"created_at"`
	LastModified      time.Time     `bson:"last_modified" json:"last_modified"`
	DeletedAt         time.Time     `bson:"deleted_at" json:"deleted_at"`
}
