package unlink

import (
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/constants"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type AccountBlockResponse struct {
	ID      string `json:"id" bson:"id"`
	Name    string `json:"name" bson:"name"`
	Enabled bool   `json:"enabled,omitempty" bson:"enabled,omitempty"`
	Reason  string `json:"reason,omitempty" bson:"reason,omitempty"`
}

type ArchivedUserResponse struct {
	ID               bson.ObjectID              `json:"id,omitempty" bson:"_id,omitempty"`
	CustomerCode     string                     `json:"customer_code" bson:"customer_code"`
	CustomerName     string                     `json:"customer_name" bson:"customer_name"`
	BranchCode       string                     `json:"branch_code" bson:"branch_code"`
	Branch           AccountBlockResponse       `json:"branch" bson:"branch"`
	District         AccountBlockResponse       `json:"district" bson:"district"`
	PhoneNumber      string                     `json:"phone_number" bson:"phone_number"`
	AccountNumber    string                     `json:"account_number" bson:"account_number"`
	Language         string                     `json:"language" bson:"language"`
	Avatar           string                     `json:"avatar" bson:"avatar"`
	Email            string                     `json:"email" bson:"email"`
	PushToken        string                     `json:"push_token" bson:"push_token"`
	CustomerNumber   string                     `json:"customer_number" bson:"customer_number"`
	UserCategory     string                     `json:"user_category" bson:"user_category"`
	Industry         string                     `json:"industry" bson:"industry"`
	Sector           string                     `json:"sector" bson:"sector"`
	Ownership        string                     `json:"ownership" bson:"ownership"`
	CustomerSegment  string                     `json:"customer_segment" bson:"customer_segment"`
	BlockedReason    string                     `json:"blocked_reason" bson:"blocked_reason,omitempty"` // -- optional
	DeviceUUID       string                     `json:"device_uuid" bson:"device_uuid"`
	AppVersion       string                     `json:"app_version" bson:"app_version"`
	Gender           constants.Gender           `json:"gender" bson:"gender"`
	MemberType       constants.MemberType       `json:"account_branch_type" bson:"account_branch_type"`
	Platform         constants.Platform         `json:"platform" bson:"platform"`
	DeviceStatus     constants.DeviceStatus     `json:"device_status" bson:"device_status"`
	OnboardingMethod constants.OnboardingMethod `json:"onboarding_method" bson:"onboarding_method"`
	// LoginPIN            types.LoginPIN              `json:"login_pin" bson:"login_pin"`
	BlockedOn           constants.BlockedOn         `json:"blocked_on" bson:"blocked_on,omitempty"`   // -- optional
	EnabledChannels     []constants.EnabledChannels `json:"enabled_channels" bson:"enabled_channels"` // ussd, superapp, internet banking
	LoginAttemptCount   uint8                       `json:"login_attempt_count" bson:"login_attempt_count"`
	LastLoginAttempt    time.Time                   `json:"last_login_attempt" bson:"last_login_attempt"`
	LastLogin           time.Time                   `json:"last_login" bson:"last_login"`
	APPInstallationDate time.Time                   `json:"application_installation_date" bson:"application_installation_date"`
	CreatedAt           time.Time                   `json:"created_at" bson:"created_at"`
	ExpiryAt            time.Time                   `json:"expiry_at" bson:"expiry_at"`
	LastModifiedAt      time.Time                   `json:"last_modified_at" bson:"last_modified_at"`
	IsBlocked           bool                        `json:"is_blocked" bson:"is_blocked"` // default: false
	Enabled             bool                        `json:"enabled" bson:"enabled"`
	FirstPinSet         bool                        `json:"first_pin_set" bson:"first_pin_set"`
	IsActivated         bool                        `json:"is_activated" bson:"is_activated"`
}
