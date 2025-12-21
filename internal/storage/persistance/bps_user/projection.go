package bps_user

import (
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// BPSUserMapper maps BPSUser model to BSON for database operations
func BPSUserMapper(data model.BPSUser) bson.M {
	result := bson.M{}
	if data.UserCode != "" {
		result["user_code"] = data.UserCode
	}
	if data.FullName != "" {
		result["full_name"] = data.FullName
	}
	if data.HomeBranch != "" {
		result["home_branch"] = data.HomeBranch
	}
	if data.PhoneNumber != "" {
		result["phone_number"] = data.PhoneNumber
	}
	if data.Role != "" {
		result["role"] = data.Role
	}

	if data.UserCode != "" {
		result["user_code"] = data.UserCode
	}

	if data.UserName != "" {
		result["username"] = data.UserName
	}

	result["enabled"] = data.Enabled
	return result
}
