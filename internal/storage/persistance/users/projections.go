package users

import (
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/errors"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func UserProjection() bson.M {
	return bson.M{
		"_id":                           1,
		"username":                      1,
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
		"device_status":                 1,
		"enabled":                       1,
		"first_pin_set":                 1,
		"device_uuid":                   1,
	}
}

func UserIdFilterAttachMent(id string, filter bson.M) error {
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.ErrUnexpected
	}
	filter = bson.M{
		"_id": objId,
	}
	return nil
}

func UserPhoneFilterAttachment(phone string, filter bson.M) {
	filter = bson.M{
		"phone_number": phone,
	}
}

func UserDeviceUUIDAttachment(deviceUUID string, filter bson.M) {
	filter = bson.M{
		"device_uuid": deviceUUID,
	}
}

func UserBuilder(update model.User, data bson.M) {

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
	if update.Fayda.FaydaID != "" {
		data["fayda.id_number"] = update.Fayda.FaydaID
	}
	if update.Fayda.EmploymentStatus != "" {
		data["fayda.employement_status"] = update.Fayda.EmploymentStatus
	}
	if update.Avatar != "" {
		data["photo"] = update.Avatar
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

}
