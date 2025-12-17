package customer

import (
	member "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func UserProjection() bson.M {
	return bson.M{
		"_id":                           1,
		"username":                      1,
		"full_name":                     1,
		"login_pin":                     1,
		"user_code":                     1,
		"region_name":                   1,
		"branch_name":                   1,
		"is_blocked":                    1,
		"phone_number":                  1,
		"main_account":                  1,
		"avatar":                        1,
		"gender":                        1,
		"created_at":                    1,
		"blocked_on_cps":                1,
		"profile_theme_type":            1,
		"member_type":                   1,
		"is_verified":                   1,
		"application_installation_date": 1,
		"login_attempt_count":           1,
		"last_login_attempt":            1,
		"device_status":                 1,
		"enabled":                       1,
		"first_pin_set":                 1,
		"device_uuid":                   1,
	}
}

func FaydaEnable(objID bson.ObjectID, data member.User) (bson.M, bson.M) {
	filter := bson.M{
		"_id": objID,
	}
	// update := bson.M{
	// 	"fayda_risk_level": data.FaydaRiskLevel,
	// }

	return filter, nil
}
