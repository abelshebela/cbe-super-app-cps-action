package otp

import (
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func OtpProjection() bson.M {
	return bson.M{
		"_id":          1,
		"otp_code":     1,
		"otp_for":      1,
		"status":       1,
		"user_code":    1,
		"device_uuid":  1,
		"phone_number": 1,
		"expires_at":   1,
	}
}

func OtpDataBuilder(o model.OTP, data bson.M) {
	if o.PhoneNumber != "" {
		data["phone_number"] = o.PhoneNumber
	}
	if o.UserCode != "" {
		data["user_code"] = o.UserCode
	}
	if o.OTPCode != "" {
		data["otp_code"] = o.OTPCode
	}
	if o.UserRealm != "" {
		data["user_realm"] = o.UserRealm
	}
	if o.OTPFor != "" {
		data["otp_for"] = o.OTPFor
	}
	if o.Status != "" {
		data["status"] = o.Status
	}
	if o.DeviceUUID != nil && *o.DeviceUUID != "" {
		data["device_uuid"] = *o.DeviceUUID
	}
	if !o.ExpiresAt.IsZero() {
		data["expires_at"] = o.ExpiresAt
	}

}
