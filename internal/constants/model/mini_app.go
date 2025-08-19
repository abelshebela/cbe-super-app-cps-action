package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// MiniApp represents a mini application within the super app ecosystem
type MiniApp struct {
	ID            bson.ObjectID `bson:"_id,omitempty"`
	Code          string        `bson:"code"`
	Name          string        `bson:"name"`
	Description   string        `bson:"description"`
	Version       string        `bson:"version"`
	IconURL       string        `bson:"icon_url"`
	BannerURL     string        `bson:"banner_url"`
	Category      string        `bson:"category"`
	DeveloperInfo string        `bson:"developer_info"`
	APIEndpoint   string        `bson:"api_endpoint"`
	WebhookURL    string        `bson:"webhook_url"`
	Enabled       bool          `bson:"enabled"`
	IsDeleted     bool          `bson:"is_deleted"`
	CreatedAt     time.Time     `bson:"created_at"`
	UpdatedAt     time.Time     `bson:"updated_at"`
	DeletedAt     time.Time     `bson:"deleted_at,omitempty"`
}