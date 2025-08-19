package dto

import (
	"mime/multipart"

	"cbe-super-app-cps-action/internal/constants"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// DeviceLookupRequest represents the request for device lookup
type DeviceLookupRequest struct {
	Platform                    constants.Platform `json:"platform"`
	AppVersion                  string             `json:"app_version"`
	DeviceUUID                  string             `json:"device_uuid"`
	SourceApp                   string             `json:"source_app"`
	ApplicationInstallationDate string             `json:"installation_date"`
}

type ChangePinRequest struct {
	UserID string `json:"user_id"`
	OldPin string `json:"old_pin"`
	NewPin string `json:"new_pin"`
}

type SetPinRequest struct {
	NewPin     string          `json:"new_pin" validate:"required,min=6,max=6"`
	UserID     string          `json:"user_id"`
	DeviceUUID string          `json:"device_uuid"`
	Realm      constants.Realm `json:"realm"`
}

type SetProfileThemeRequest struct {
	ThemeType string `json:"theme_type"`
}

type VerifyOTPRequest struct {
	UserID      string `json:"user_id"`
	PhoneNumber string `json:"phone_number"`
	OTP         string `json:"otp_code"`
	DeviceUUID  string `json:"device_uuid"`
	UserRealm   string `json:"user_realm"`
	OtpFor      string `json:"otp_for"`
	Action      string `json:"action"`
}

type RegisterRequest struct {
	Phone               string             `json:"phone" validate:"required"`
	DeviceUUID          string             `json:"device_uuid"`
	Platform            constants.Platform `json:"platform"`
	FullName            string             `json:"full_name"`
	Email               string             `json:"email,omitempty"`
	AppVersion          string             `json:"app_version"`
	SourceApp           string             `json:"source_app"`
	APPInstallationDate string             `json:"application_installation_date"`
}

type ResetPinRequest struct {
	ResetSessionID string `json:"reset_session_id" validate:"required"`
	Phone          string `json:"phone" validate:"required"`
	DeviceUUID     string `json:"device_uuid" validate:"required"`
	OTP            string `json:"otp" validate:"required,min=6,max=6"`
	NewPin         string `json:"new_pin" validate:"required,min=6,max=6"`
	UserID         string `json:"user_id"`
}

type UpdateProfilePicture struct {
	UserId         string
	ProfilePicture *multipart.FileHeader `form:"profile_picture"`
	File           multipart.File
}

type LoginRequest struct {
	DeviceUUID string
	Phone      string `json:"phone"`
	Pin        string `json:"pin"`
}

func (r LoginRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Pin, validation.Required.Error("pin is required"), validation.Length(6, 6), is.Digit))
}

type VerifyForgetPinOtpRequest struct {
	UserId         string `json:"user_id"`
	FullName       string `json:"full_name"`
	ResetSessionID string `json:"reset_session_id" validate:"required"`
	Phone          string `json:"phone" validate:"required"`
	DeviceUUID     string `json:"device_uuid" validate:"required"`
	OTP            string `json:"otp" validate:"required,min=6,max=6"`
}

type PhoneLoginRequest struct {
	Phone string `json:"phone,omitempty" validate:"required"`
}

type PinStrengthRequest struct {
	NewPin int `json:"new_pin" bson:"new_pin"`
}

type CompleteRegistrationRequest struct {
	RegistrationID string `json:"registration_id" validate:"required"`
	Phone          string `json:"phone" validate:"required"`
	DeviceUUID     string `json:"device_uuid" validate:"required"`
	Platform       string `json:"platform" validate:"required,oneof=android ios web"`
	FullName       string `json:"full_name" validate:"required"`
	OTP            string `json:"otp" validate:"required,min=6,max=6"`
}

type ForgetPinSendOtpRequest struct {
	Phone string `json:"phone"`
}

type OTPRequest struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

type ResetPinWithTokenRequest struct {
	ResetSessionID string `json:"reset_session_id" validate:"required"`
	Phone          string `json:"phone" validate:"required"`
	DeviceUUID     string `json:"device_uuid" validate:"required"`
	NewPin         string `json:"new_pin" validate:"required,min=6,max=6"`
}

type OTPVerification struct {
	UserID string `json:"user_id" bson:"user_id"`
	Email  string `json:"email" bson:"email"`
	OTP    string `json:"otp" bson:"otp"`
}

// SMSRequest represents the request for sending SMS
type SMSRequest struct {
	Recipient   string `json:"recipient" validate:"required"`
	MessageBody string `json:"message_body" validate:"required"`
}
