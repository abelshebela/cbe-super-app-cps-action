package entity

import (
	"encoding/json"
	"time"
)

type KafkaMessage struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Payload   json.RawMessage        `json:"payload"`
	CreatedAt time.Time              `json:"created_at"`
	Headers   map[string]interface{} `json:"headers,omitempty"`
}

type FeedbackKafkaMessage struct {
	UserID     string                 `json:"user_id"`
	FeedbackID string                 `json:"feedback_id"`
	Responses  map[string]Response    `json:"responses"`
	CreatedAt  time.Time              `json:"created_at"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
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
