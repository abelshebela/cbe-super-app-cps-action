package dto

import (
	"time"
)

type FetchLinkedAccountsRequest struct {
	UserID string `json:"user_id"`
}

type GenerateOTPRequest struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

type VerifyOTPRequest struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	OTP    string `json:"otp"`
}

type GenerateOTPResponse struct {
	OTP string `json:"otp"`
}

type VerifyOTPResponse struct {
	Success bool `json:"success"`
}

type ChangePinRequest struct {
	UserID string `json:"user_id"`
	OldPin string `json:"old_pin"`
	NewPin string `json:"new_pin"`
}

type SetProfileThemeRequest struct {
	ThemeType string `json:"theme_type"`
}

type RegisterRequest struct {
	Phone      string `json:"phone" validate:"required"`
	DeviceUUID string `json:"device_uuid" validate:"required"`
	Platform   string `json:"platform" validate:"required,oneof=android ios web"`
	FullName   string `json:"full_name,omitempty"`
	Email      string `json:"email,omitempty"`
}

type RegisterResponse struct {
	RegistrationID     string    `json:"registration_id"`
	PhoneNumber        string    `json:"phone_number"`
	DeviceUUID         string    `json:"device_uuid"`
	Platform           string    `json:"platform"`
	OTPSent            bool      `json:"otp_sent"`
	Otp                string    `json:"otp"`
	OTPExpiryMinutes   int       `json:"otp_expiry_minutes"`
	Token              string    `json:"temp_token,omitempty"`
	TokenType          string    `json:"token_type,omitempty"`
	TokenExpiry        time.Time `json:"token_expiry,omitempty"`
	OTPCode            string    `json:"otp_code,omitempty"`
	NextStep           string    `json:"next_step"`
	RegistrationStatus string    `json:"registration_status"`
}

type LoginRequest struct {
	Phone      string `json:"phone" validate:"required"`
	DeviceUUID string `json:"device_uuid" validate:"required"`
	Pin        string `json:"pin" validate:"required,min=6,max=6"`
}

type PhoneLoginRequest struct {
	Phone string `json:"phone,omitempty" validate:"required"`
}

type LoginResponse struct {
	Token          string    `json:"access_token"`
	UserID         string    `json:"user_id"`
	UserCode       string    `json:"user_code"`
	FullName       string    `json:"full_name"`
	PhoneNumber    string    `json:"phone_number"`
	KYCLevel       uint8     `json:"kyc_level"`
	IsVerified     bool      `json:"is_verified"`
	DeviceUUID     string    `json:"device_uuid"`
	Platform       string    `json:"platform"`
	AppVersion     string    `json:"app_version"`
	LoginTime      time.Time `json:"login_time"`
	SessionExpires time.Time `json:"session_expires"`
	LastLogin      time.Time `json:"last_login"`
	LoginAttempts  int       `json:"login_attempts"`
}

// Forget PIN DTOs
type ForgetPinSendOtpRequest struct {
	Phone      string `json:"phone" validate:"required"`
	DeviceUUID string `json:"device_uuid" validate:"required"`
}

type ForgetPinSendOtpResponse struct {
	PhoneNumber      string    `json:"phone_number"`
	DeviceUUID       string    `json:"device_uuid"`
	OTPSent          bool      `json:"otp_sent"`
	OTPExpiryMinutes int       `json:"otp_expiry_minutes"`
	ResetSessionID   string    `json:"reset_session_id"`
	Token            string    `json:"temp_token,omitempty"`
	TokenType        string    `json:"token_type,omitempty"`
	TokenExpiry      time.Time `json:"token_expiry,omitempty"`
	NextStep         string    `json:"next_step"`
	OTP              string    `json:"otp,omitempty"`
}

type ResetPinRequest struct {
	ResetSessionID string `json:"reset_session_id" validate:"required"`
	Phone          string `json:"phone" validate:"required"`
	DeviceUUID     string `json:"device_uuid" validate:"required"`
	OTP            string `json:"otp" validate:"required,min=6,max=6"`
	NewPin         string `json:"new_pin" validate:"required,min=6,max=6"`
}

type ResetPinResponse struct {
	UserID           string    `json:"user_id,omitempty"`
	UserCode         string    `json:"user_code,omitempty"`
	FullName         string    `json:"full_name,omitempty"`
	PhoneNumber      string    `json:"phone_number,omitempty"`
	PinReset         bool      `json:"pin_reset,omitempty"`
	ResetTime        time.Time `json:"reset_time,omitempty"`
	AccessRestricted bool      `json:"access_restricted,omitempty"`
	Restrictions     []string  `json:"restrictions,omitempty"`
	Token            string    `json:"token,omitempty"`
	TokenType        string    `json:"token_type,omitempty"`
	TokenExpiry      time.Time `json:"token_expiry,omitempty"`
	NextStep         string    `json:"next_step,omitempty"`
}

// Device Lookup DTOs
type DeviceLookupResponse struct {
	DeviceUUID  string    `json:"device_uuid,omitempty"`
	UserID      string    `json:"user_id,omitempty"`
	UserCode    string    `json:"user_code,omitempty"`
	FullName    string    `json:"full_name,omitempty"`
	PhoneNumber string    `json:"phone_number,omitempty"`
	Email       string    `json:"email,omitempty"`
	Platform    string    `json:"platform,omitempty"`
	AppVersion  string    `json:"app_version,omitempty"`
	IsLatest    bool      `json:"is_latest"`
	UserFound   bool      `json:"user_found,omitempty"`
	Token       string    `json:"temp_token,omitempty"`
	TokenType   string    `json:"token_type,omitempty"`
	TokenExpiry time.Time `json:"token_expiry,omitempty"`
	NextStep    string    `json:"next_step,omitempty"`
	OTPCode     string    `json:"otp_code,omitempty"`
	OTPFor      string    `json:"otp_for,omitempty"`
	Message     string    `json:"-"`
	Status      int       `json:"-"`
}

type PinStrengthRequest struct {
	NewPin int `json:"new_pin" bson:"new_pin"`
}

// Verify OTP DTOs
type VerifyOtpResponse struct {
	UserID      string    `json:"user_id,omitempty"`
	PhoneNumber string    `json:"phone_number"`
	OTPVerified bool      `json:"otp_verified"`
	Token       string    `json:"temp_token,omitempty"`
	TokenType   string    `json:"token_type,omitempty"`
	TokenExpiry time.Time `json:"token_expiry,omitempty"`
	NextStep    string    `json:"next_step"`
}

// Set PIN DTOs
type SetPinResponse struct {
	UserID      string    `json:"user_id"`
	PinSet      bool      `json:"pin_set"`
	Token       string    `json:"access_token,omitempty"`
	TokenType   string    `json:"token_type,omitempty"`
	TokenExpiry time.Time `json:"token_expiry,omitempty"`
	NextStep    string    `json:"next_step"`
}

// Complete Registration DTOs
type CompleteRegistrationResponse struct {
	UserID      string    `json:"user_id"`
	UserCode    string    `json:"user_code"`
	FullName    string    `json:"full_name"`
	PhoneNumber string    `json:"phone_number"`
	KYCLevel    uint8     `json:"kyc_level"`
	Status      string    `json:"status"`
	Token       string    `json:"token,omitempty"`
	TokenType   string    `json:"token_type,omitempty"`
	TokenExpiry time.Time `json:"token_expiry,omitempty"`
	NextStep    string    `json:"next_step"`
}
