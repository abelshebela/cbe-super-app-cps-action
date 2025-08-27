package accountvalidation

import "time"

type ValidationRuleDTO struct {
	ID            string `json:"id"`
	EntityType    string `json:"entity_type"`
	ValidationFor string `json:"validation_for"`
	Identifier    string `json:"identifier"`
	MinLength     uint8  `json:"min_length"`
	MaxLength     uint8  `json:"max_length"`
	Enabled       bool   `json:"enabled"`
	IsDeleted     bool   `json:"is_deleted"`
	ServiceID     string `json:"service_id"`
}

type ValidationRule struct {
	ID             string    `json:"id"`
	EntityType     string    `json:"entity_type"`
	ValidationFor  string    `json:"validation_for"`
	Identifier     string    `json:"identifier"`
	MinLength      uint8     `json:"min_length"`
	MaxLength      uint8     `json:"max_length"`
	Enabled        bool      `json:"enabled"`
	IsDeleted      bool      `json:"is_deleted"`
	CreatedAt      time.Time `json:"created_at"`
	LastModifiedAt time.Time `json:"last_modified_at"`
	ServiceID      string    `json:"service_id"`
}

type User struct {
	ID          string
	FullName    string
	PhoneNumber string
	Department  string
}

type UpdateAccountValidationResponse struct {
	ActionID string `json:"action_id"`
}

type Request struct {
	ActionCode     string `json:"action_code"`
	Decison        bool   `json:"decison"`
	RejectedReason string `json:"rejected_reason"`
}

type ApproveRejectRequest struct {
	ActionCode     string      `json:"action_code"`
	Decison        DecisonEnum `json:"decison"`
	RejectedReason string      `json:"rejected_reason"`
}

type UpdateAccountValidationRequest struct {
	ValidationRuleDTO
}