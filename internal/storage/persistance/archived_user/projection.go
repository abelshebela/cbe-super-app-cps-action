package archived_user

import (
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
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
		"main_account":       1,
		"branch_name":        1,
		"district_name":      1,
		"kyc_level":          1,
		"is_verified":        1,
		"is_blocked":         1,
		"last_login":         1,
	}
}

func UserToArchivedUser(user *member.User) *model.ArchivedUser {
	if user == nil {
		return nil
	}
	return &model.ArchivedUser{
		ID:                  user.ID,
		UserCode:            user.UserCode,
		FullName:            user.FullName,
		BranchCode:          user.BranchCode,
		PhoneNumber:         user.PhoneNumber,
		Gender:              user.Gender,
		Avatar:              user.Avatar,
		Email:               user.Email,
		PushToken:           user.PushToken,
		MemberType:          user.MemberType,
		IsBlocked:           user.IsBlocked,
		LoginAttemptCount:   user.LoginAttemptCount,
		LastLoginAttempt:    user.LastLoginAttempt,
		LastLogin:           user.LastLogin,
		LoginPIN:            user.LoginPIN,
		DeviceUUID:          user.DeviceUUID,
		AppVersion:          user.AppVersion,
		Platform:            user.Platform,
		APPInstallationDate: user.APPInstallationDate,
		CustomerNumber:      user.CustomerNumber,
		DeviceStatus:        user.DeviceStatus,
		Enabled:             user.Enabled,
		FirstPinSet:         user.FirstPinSet,
		CreatedAt:           user.CreatedAt,
		LastModifiedAt:      user.LastModifiedAt,
	}
}
