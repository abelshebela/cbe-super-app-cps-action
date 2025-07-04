package entities

import "time"

type ValidationRule struct {
	ID             string    `json:"id" bson:"id"`
	EntityType     string    `json:"entity_type" bson:"entity_type"`
	ValidationFor  string    `json:"validation_for" bson:"validation_for"`
	Identifier     string    `json:"identifier" bson:"identifier"`
	MinLength      uint8     `json:"min_length" bson:"min_length"`
	MaxLength      uint8     `json:"max_length" bson:"max_length"`
	Enabled        bool      `json:"enabled" bson:"enabled"`
	IsDeleted      bool      `json:"is_deleted" bson:"is_deleted"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
	LastModifiedAt time.Time `json:"last_modified_at" bson:"last_modified_at"`
	ServiceID      string    `json:"service_id" bson:"service_id"`
}
