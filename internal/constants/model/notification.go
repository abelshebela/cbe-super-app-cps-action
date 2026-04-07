package model

import (
	"cbe-super-app-cps-action/internal/constants"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Notification struct {
	ID                bson.ObjectID                `json:"id,omitempty" bson:"_id"`
	NotificationCode  string                       `json:"notification_code" bson:"notification_code"`
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
	IsFromCPS         bool                         `json:"is_from_cps" bson:"is_from_cps"`
	IsDeleted         bool                         `json:"is_deleted" bson:"is_deleted"`
	CreatedAt         time.Time                    `json:"created_at" bson:"created_at"`
	LastModified      time.Time                    `json:"last_modified" bson:"last_modified"`
	DeletedAt         time.Time                    `json:"deleted_at" bson:"deleted_at"`
}

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
