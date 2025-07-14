package account_validation

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

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
