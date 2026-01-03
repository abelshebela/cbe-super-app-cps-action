package model

import (
	"time"

	"cbe-super-app-cps-action/internal/constants"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type OTP struct {
	ID          bson.ObjectID       `json:"id,omitempty" bson:"_id,omitempty"`
	PhoneNumber string              `json:"phone_number" bson:"phone_number"`
	FullName    string              `json:"full_name" bson:"full_name"`
	UserRealm   constants.Realm     `json:"user_realm" bson:"user_realm"`
	Email       string              `json:"email" bson:"email"`
	UserCode    string              `json:"user_code" bson:"user_code"`
	OTPCode     string              `json:"otp_code" bson:"otp_code"`
	BillNo      *string             `json:"bill_no,omitempty" bson:"bill_no,omitempty"`
	DeviceUUID  *string             `json:"device_uuid,omitempty" bson:"device_uuid,omitempty"`
	OTPFor      constants.OTPFor    `json:"otp_for" bson:"otp_for"`
	Status      constants.OTPStatus `json:"status" bson:"status"`
	ExpiresAt   time.Time           `json:"expires_at" bson:"expires_at"`
	CreatedAt   time.Time           `json:"created_at" bson:"created_at"`
	IsDeleted   bool                `json:"is_deleted" bson:"is_deleted"`
}
