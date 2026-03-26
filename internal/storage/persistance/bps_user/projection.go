package bps_user

import (
	// "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	// local_model "cbe-super-app-cps-action/internal/constants/model"
	bpsUserDto "cbe-super-app-cps-action/internal/constants/dto/bps_user"

	bps_model "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/bps"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// BPSUserMapper maps BPSUser model to BSON for database operations
func BPSUserMapper(data bps_model.BPSUser) bson.M {
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
	// if data.Role != "" {
	// 	result["role"] = data.Role
	// }
	if len(data.BranchCode) > 0 {
		result["branch_code"] = data.BranchCode
	}

	if data.Email != "" {
		result["email"] = data.Email
	}

	if data.Username != "" {
		result["username"] = data.Username
	}
	if data.JobTitle != "" {
		result["job_title"] = data.JobTitle
	}

	// result["role"] = data.Role
	result["enabled"] = data.Enabled
	return result
}

func BPSUserResponseMapper(data bps_model.BPSUser) *bpsUserDto.BPSUserResposenDTO {
	return &bpsUserDto.BPSUserResposenDTO{
		ID:                data.ID,
		UserCode:          data.UserCode,
		FullName:          data.FullName,
		Username:          data.Username,
		Email:             data.Email,
		PhoneNumber:       data.PhoneNumber,
		BranchCode:        data.BranchCode,
		BranchName:        data.BranchName,
		HomeBranch:        data.HomeBranch,
		Role:              data.Role,
		Realm:             data.Realm,
		LoginAttemptCount: data.LoginAttemptCount,
		FirstPasswordSet:  data.FirstPasswordSet,
		Enabled:           data.Enabled,
		IsDeleted:         data.IsDeleted,
		OTPVerifyCount:    data.OTPVerifyCount,
		OTPLastTriedAt:    data.OTPLastTriedAt,
		JobTitle:          data.JobTitle,
		OTPLastVerifiedAt: data.OTPLastVerifiedAt,
		IsFirstTimeLogin:  data.IsFirstTimeLogin,
		LastLoginAttempt:  data.LastLoginAttempt,
		NextLoginAttempt:  data.NextLoginAttempt,
		LastLogin:         data.LastLogin,
		CreatedAt:         data.CreatedAt,
		LastModifiedAt:    data.LastModifiedAt,
	}
}
