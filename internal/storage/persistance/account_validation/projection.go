package accountvalidation

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// AccountMapper maps ValidationRule model to BSON for database operations
func AccountValidationMapper(data model.ValidationRule) bson.M {
	result := bson.M{}

	if data.EntityType != "" {
		result["entity_type"] = data.EntityType
	}
	if data.ValidationFor != "" {
		result["validation_for"] = data.ValidationFor
	}
	if data.Identifier != "" {
		result["identifier"] = data.Identifier
	}
	if data.MinLength != 0 {
		result["min_length"] = data.MinLength
	}
	if data.MaxLength != 0 {
		result["max_length"] = data.MaxLength
	}

	if data.Enabled != nil {
		result["enabled"] = data.Enabled

	}
	if data.IsDeleted != nil {
		result["is_deleted"] = data.IsDeleted
	}
	// timestamps (always included)
	if !data.CreatedAt.IsZero() {
		result["created_at"] = data.CreatedAt
	}
	result["last_modified_at"] = time.Now() // or data.LastModifiedAt if you want to keep it

	if data.ServiceID != "" {
		result["service_id"] = data.ServiceID
	}

	return result
}
