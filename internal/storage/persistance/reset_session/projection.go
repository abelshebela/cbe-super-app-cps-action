package reset_session

import (
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func ResetSessionProjection() bson.M {
	return bson.M{
		"_id":               1,
		"user_id":           1,
		"phone_number":      1,
		"device_uuid":       1,
		"expires_at":        1,
		"created_at":        1,
		"attempts":          1,
		"max_attempts":      1,
		"otp":               1,
		"otp_for":           1,
		"verified_at":       1,
		"completed_at":      1,
		"enabled":           1,
		"access_restricted": 1,
		"status":            1,
		"restrictions":      1,
	}
}

func ResetSessionIdFilterAttachment(id string, filter bson.M) error {
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.ErrUnexpected
	}
	filter["_id"] = objId
	return nil
}

func ResetSessionPhoneFilterAttachment(phoneNumber string) bson.M {
	return bson.M{"phone_number": phoneNumber}
}

func ResetSessionDeviceUUIDFilterAttachment(deviceUUID string, filter bson.M) {
	filter["device_uuid"] = deviceUUID
}
