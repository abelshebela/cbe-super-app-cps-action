package password_rule

import (
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// PasswordRuleMapper maps a PasswordRule model to a bson.M for updates
func PasswordRuleMapper(rule model.PasswordRule) bson.M {
	return bson.M{
		"name":            rule.Name,
		"min_length":      rule.MinLength,
		"max_length":      rule.MaxLength,
		"numbers":         rule.Numbers,
		"capital_letters": rule.CapitalLetters,
		"small_letters":   rule.SmallLetters,
		"characters":      rule.Characters,
		"updated_at":      time.Now(),
	}
}
