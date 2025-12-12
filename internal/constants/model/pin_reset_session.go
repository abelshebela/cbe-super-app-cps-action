package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type PinResetSession struct {
	ID               bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID           string        `bson:"user_id" json:"user_id"`
	PhoneNumber      string        `bson:"phone_number" json:"phone_number"`
	DeviceUUID       string        `bson:"device_uuid" json:"device_uuid"`
	ExpiresAt        time.Time     `bson:"expires_at" json:"expires_at"`
	CreatedAt        time.Time     `bson:"created_at" json:"created_at"`
	Attempts         int           `bson:"attempts" json:"attempts"`
	MaxAttempts      int           `bson:"max_attempts" json:"max_attempts"`
	OTP              string        `bson:"otp" json:"otp"`
	OTPFor           string        `bson:"otp_for" json:"otp_for"`
	VerifiedAt       *time.Time    `bson:"verified_at,omitempty" json:"verified_at,omitempty"`
	CompletedAt      *time.Time    `bson:"completed_at,omitempty" json:"completed_at,omitempty"`
	Enabled          bool          `bson:"enabled" json:"enabled"`
	AccessRestricted bool          `bson:"access_restricted" json:"access_restricted"`
	Status           string        `bson:"status" json:"status"`
	Restrictions     []string      `bson:"restrictions" json:"restrictions"`
	IsDeleted        bool          `bson:"is_deleted" json:"is_deleted"`
}
