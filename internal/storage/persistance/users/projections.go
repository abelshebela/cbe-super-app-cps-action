package users

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"errors"

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

func UserBuilder(update model.User) bson.M {
	data := bson.M{}
	if update.FullName != "" {
		data["full_name"] = update.FullName
	}
	if update.MotherName != "" {
		data["mother_name"] = update.MotherName
	}
	if !update.BirthDate.IsZero() {
		data["birth_date"] = update.BirthDate
	}
	if update.Nationality != "" {
		data["nationality"] = update.Nationality
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
	if update.Username != "" {
		data["username"] = update.Username
	}
	if update.IsSelfRegister {
		data["is_self_register"] = true
	}
	if update.IsVerified {
		data["is_verified"] = update.IsVerified
	}
	if update.DeviceUUID != "" {
		data["device_uuid"] = update.DeviceUUID
	}
	if update.LoginPIN != (types.LoginPIN{}) {
		data["login_pin"] = update.LoginPIN
	}
	if update.ProfileThemeType != "" {
		data["profile_theme_type"] = update.ProfileThemeType
	}
	if update.IsBlocked {
		data["is_blocked"] = update.IsBlocked
	}
	if update.Address.Zone != "" {
		data["address.zone"] = update.Address.Zone
	}
	if update.Address.Kebele != "" {
		data["address.kebele"] = update.Address.Kebele
	}
	if update.Address.Woreda != "" {
		data["address.woreda"] = update.Address.Woreda
	}
	if update.Address.Region != "" {
		data["address.region"] = update.Address.Region
	}
	if update.BranchName != "" {
		data["branch_name"] = update.BranchName
	}
	if update.DistrictName != "" {
		data["district_name"] = update.DistrictName
	}
	if update.BranchCode != "" {
		data["branch_code"] = update.BranchCode
	}
	if update.DistrictCode != "" {
		data["district_code"] = update.DistrictCode
	}
	if update.ResidentialStatus != "" {
		data["residential_status"] = update.ResidentialStatus
	}
	if !update.IssuedDate.IsZero() {
		data["issued_date"] = update.IssuedDate
	}

	return data
}
