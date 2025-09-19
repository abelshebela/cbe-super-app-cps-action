package archived_user

import (
	"cbe-super-app-cps-action/internal/constants/model"

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
		"mainaccount":        1,
		"branch":             1,
		"district":           1,
		"kyc_level":          1,
		"is_verified":        1,
		"is_blocked":         1,
		"last_login":         1,
	}
}

func UserToArchivedUser(user *model.User) *model.ArchivedUser {
	if user == nil {
		return nil
	}
	return &model.ArchivedUser{
		ID:                user.ID,
		UserCode:          user.UserCode,
		FullName:          user.FullName,
		MotherName:        user.MotherName,
		Nationality:       user.Nationality,
		BirthDate:         user.BirthDate,
		ResidentialStatus: user.ResidentialStatus,
		PhoneNumber:       user.PhoneNumber,
		Gender:            user.Gender,
		Email:             user.Email,
		CreatedAt:         user.CreatedAt,
		AccountStatus:     user.AccountStatus,
		KYCLevel:          user.KYCLevel,
		IsVerified:        user.IsVerified,
		IsBlocked:         user.IsBlocked,
		LastLogin:         user.LastLogin,
		// Add other fields as needed if present in both structs
	}
}
