package password_rule

import (
	local_model "cbe-super-app-cps-action/internal/constants/model"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// PasswordRuleMapper maps a PasswordRule model to a bson.M for updates
func PasswordRuleMapper(rule local_model.PasswordRule) bson.M {
	result := bson.M{}
	if rule.Name != "" {
		result["name"] = rule.Name
	}
	if rule.MinLength != 0 {
		result["min_length"] = rule.MinLength
	}
	if rule.MaxLength != 0 {
		result["max_length"] = rule.MaxLength
	}
	result["numbers"] = rule.Numbers
	result["capital_letters"] = rule.CapitalLetters
	result["small_letters"] = rule.SmallLetters
	result["characters"] = rule.Characters
	result["updated_at"] = time.Now()
	return result
}
