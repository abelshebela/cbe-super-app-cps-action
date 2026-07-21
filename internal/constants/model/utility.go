package model

import (
	"encoding/json"
	"time"
)

type UtilityKafkaMessage struct {
	Action       string           `json:"action"`        // "create"|"update"|"enable"|"disable"|"delete"
	UniqueTokens []string         `json:"unique_tokens"` // CREATE only — UniqueId = UniqueTokens[0]
	ID           string           `json:"id"`            // non-CREATE: MongoDB ObjectId of existing CPS action
	Maker        UtilityMakerInfo `json:"maker"`
	Payload      json.RawMessage  `json:"payload"` // arbitrary JSON — stored as-is
	CreatedAt    time.Time        `json:"created_at"`
}

type UtilityMakerInfo struct {
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
}
