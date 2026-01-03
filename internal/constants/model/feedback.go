package model

import (
	"cbe-super-app-cps-action/internal/constants/types"
	"time"

	"encoding/json"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Feedback struct {
	ID        bson.ObjectID             `json:"id" bson:"_id,omitempty"`
	UserID    string                    `json:"user_id" bson:"user_id"`
	Responses map[string]types.Response `json:"responses" bson:"responses"`
	CreatedAt time.Time                 `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time                 `json:"updated_at" bson:"updated_at"`
}

type FeedbackResponse struct {
	Page      int         `json:"page"`
	Feedbacks []*Feedback `json:"feedbacks"`
	Limit     int         `json:"limit"`
	Total     int64       `json:"total"`
}

// kafka messages
type KafkaMessage struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Payload   json.RawMessage        `json:"payload"`
	CreatedAt time.Time              `json:"created_at"`
	Headers   map[string]interface{} `json:"headers,omitempty"`
}

type FeedbackKafkaMessage struct {
	UserID     string                    `json:"user_id"`
	FeedbackID string                    `json:"feedback_id"`
	Responses  map[string]types.Response `json:"responses"`
	CreatedAt  time.Time                 `json:"created_at"`
	Metadata   map[string]interface{}    `json:"metadata,omitempty"`
}

type Response struct {
	Question string      `json:"question" bson:"question"`
	Answer   interface{} `json:"answer" bson:"answer"`
	Type     string      `json:"type" bson:"type"`
	Options  interface{} `json:"options" bson:"options"`
}

func NewKafkaMessage(id, msgType string, payload interface{}) (*KafkaMessage, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return &KafkaMessage{
		ID:        id,
		Type:      msgType,
		Payload:   payloadBytes,
		CreatedAt: time.Now(),
		Headers:   make(map[string]interface{}),
	}, nil
}

func (k *KafkaMessage) UnmarshalPayload(v interface{}) error {
	return json.Unmarshal(k.Payload, v)
}

type CustomerFeedback struct {
	ID            bson.ObjectID `json:"id" bson:"_id,omitempty"`
	CustomerName  string        `json:"customer_name" bson:"customer_name"`
	CustomerCode  string        `json:"customer_code" bson:"customer_code"`
	Email         string        `json:"email" bson:"email"`
	PhoneNumber   string        `json:"phone_number" bson:"phone_number"`
	AccountNumber string        `json:"account_number" bson:"account_number"`
	Rating        int           `json:"rating" bson:"rating"`
	DeviceModel   string        `json:"device_model" bson:"device_model"`
	Message       string        `json:"message" bson:"message"`
	SentAt        time.Time     `json:"sent_at" bson:"sent_at"`
	CreatedAt     time.Time     `json:"created_at" bson:"created_at"`
}
