package validation_rule

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// ValidationRuleMapper maps a ValidationRule model to a bson.M for updates
func ValidationRuleMapper(rule model.ValidationRule) bson.M {
	return bson.M{
		"$set": bson.M{
			"entity_type":      rule.EntityType,
			"validation_for":   rule.ValidationFor,
			"identifier":       rule.Identifier,
			"min_length":       rule.MinLength,
			"max_length":       rule.MaxLength,
			"enabled":          rule.Enabled,
			"service_id":       rule.ServiceID,
			"last_modified_at": time.Now(),
		},
	}
} 