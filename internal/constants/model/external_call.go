package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ExternalCall struct {
	ID         bson.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
	URL        string            `json:"url" bson:"url"`
	Method     string            `json:"method" bson:"method"`
	Payload    string            `json:"payload" bson:"payload"`
	Headers    map[string]string `json:"headers" bson:"headers"`
	StatusCode int               `json:"status_code" bson:"status_code"`
	Response   string            `json:"response" bson:"response"`
	Error      string            `json:"error,omitempty" bson:"error,omitempty"`
	Duration   time.Duration     `json:"duration" bson:"duration"`
	CreatedAt  time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at" bson:"updated_at"`
}
