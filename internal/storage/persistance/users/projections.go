package users

import (
	"cbe-super-app-cps-action/internal/constants/localization"

	"errors"

	member "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func UserProjection() bson.M {
	return bson.M{
		"_id":                           1,
		"username":                      1,
		"user_code":                     1,
		"full_name":                     1,
		"login_pin":                     1,
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
		"customer_number":               1,
	}
}

func UserIdFilterAttachMent(id string) (bson.M, error) {
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return bson.M{
		"_id":        objId,
		"is_deleted": false,
	}, nil
}

func UserPhoneFilterAttachment(phone string) bson.M {
	return bson.M{
		"phone_number": phone,
		"is_deleted":   false,
	}
}

func UserDeviceUUIDAttachment(deviceUUID string) bson.M {
	return bson.M{
		"device_uuid": deviceUUID,
		"is_deleted":  false,
	}
}

func UserBuilder(update member.User) bson.M {
	data := bson.M{}
	if update.FullName != "" {
		data["full_name"] = update.FullName
	}
	if update.PhoneNumber != "" {
		data["phone_number"] = update.PhoneNumber
	}
	if update.Gender != "" {
		data["gender"] = update.Gender
	}
	if update.Avatar != "" {
		data["photo"] = update.Avatar
		data["avatar"] = update.Avatar
	}
	if update.Email != "" {
		data["email"] = update.Email
	}
	if update.DeviceUUID != "" {
		data["device_uuid"] = update.DeviceUUID
	}

	if update.IsBlocked {
		data["is_blocked"] = update.IsBlocked
	}

	if update.BranchCode != "" {
		data["branch_code"] = update.BranchCode
	}

	return data
}
