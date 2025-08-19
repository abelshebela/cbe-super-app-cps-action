package archived_user

import (
	"go.mongodb.org/mongo-driver/bson"
)

func Projection() bson.M {
	return bson.M{
		"_id":                1,
		"user_code":          1,
		"full_name":          1,
		"mother_name":        1,
		"nationality":        1,
		"birth_date":         1,
		"residential_status": 1,
		"phone_number":       1,
		"gender":             1,
		"marital_status":     1,
		"email":              1,
		"created_at":         1,
		"account_status":     1,
		"kyc_level":          1,
		"is_verified":        1,
		"is_blocked":         1,
		"last_login":         1,
	}
}
