package customer

import "go.mongodb.org/mongo-driver/v2/bson"

func UserProjection() bson.M {
	return bson.M{
		"_id":                           1,
		"username":                      1,
		"full_name":                     1,
		"login_pin":                     1,
		"region_name":                   1,
		"branch_name":                   1,
		"is_blocked":                    1,
		"is_account_blocked":            1,
		"phone_number":                  1,
		"avatar":                        1,
		"blocked_on_cps":                1,
		"profile_theme_type":            1,
		"kyc_level":                     1,
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
