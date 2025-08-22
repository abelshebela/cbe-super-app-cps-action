package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// InAppNotificationDocument represents the in-app notification document in MongoDB
type InAppNotificationDocument struct {
	ID                bson.ObjectID `bson:"_id,omitempty"`
	NotificationType  string        `bson:"notificationType"`
	NotificationBody  string        `bson:"notificationBody"`
	IsPublic          bool          `bson:"isPublic"`
	For               string        `bson:"for"`
	CreatedBy         bson.ObjectID `bson:"createdBy"`
	NotificationParts any           `bson:"notificationParts"`
	Seen              bool          `bson:"seen"`
	Enabled           bool          `bson:"enabled"`
	IsDeleted         bool          `bson:"isDeleted"`
	CreatedAt         time.Time     `bson:"createdAt"`
	LastModified      time.Time     `bson:"lastModified"`
}
