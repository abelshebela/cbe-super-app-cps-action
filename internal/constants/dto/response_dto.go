package dto

import "time"

type DeviceLookupResponse struct {
	DeviceUUID  string    `json:"device_uuid,omitempty"`
	UserID      string    `json:"user_id,omitempty"`
	UserCode    string    `json:"user_code,omitempty"`
	FullName    string    `json:"full_name,omitempty"`
	PhoneNumber string    `json:"phone_number,omitempty"`
	Email       string    `json:"email,omitempty"`
	Platform    string    `json:"platform,omitempty"`
	AppVersion  string    `json:"app_version,omitempty"`
	IsLatest    bool      `json:"is_latest,omitempty"`
	UserFound   bool      `json:"user_found"`
	Token       string    `json:"temp_token,omitempty"`
	TokenType   string    `json:"token_type,omitempty"`
	TokenExpiry time.Time `json:"token_expiry"`
	NextStep    string    `json:"next_step,omitempty"`
	OTPCode     string    `json:"otp_code,omitempty"`
	OTPFor      string    `json:"otp_for,omitempty"`
	Message     string    `json:"-"`
	Status      int       `json:"-"`
}

type SetPinResponse struct {
	UserID      string    `json:"user_id"`
	PinSet      bool      `json:"pin_set"`
	Token       string    `json:"access_token,omitempty"`
	TokenType   string    `json:"token_type,omitempty"`
	TokenExpiry time.Time `json:"token_expiry,omitempty"`
	NextStep    string    `json:"next_step,omitempty"`
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

type VerifyOtpResponse struct {
	UserID      string    `json:"user_id,omitempty"`
	PhoneNumber string    `json:"phone_number"`
	OTPVerified bool      `json:"otp_verified"`
	Token       string    `json:"temp_token,omitempty"`
	TokenType   string    `json:"token_type,omitempty"`
	TokenExpiry time.Time `json:"token_expiry,omitempty"`
	NextStep    string    `json:"next_step"`
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
	Token            string    `json:"access_token,omitempty"`
	TokenType        string    `json:"token_type,omitempty"`
	TokenExpiry      time.Time `json:"token_expiry,omitempty"`
	NextStep         string    `json:"next_step,omitempty"`
}

type LoginResponse struct {
	Token          string    `json:"access_token"`
	UserID         string    `json:"user_id"`
	UserCode       string    `json:"user_code,omitempty"`
	FullName       string    `json:"full_name"`
	PhoneNumber    string    `json:"phone_number"`
	KYCLevel       uint8     `json:"kyc_level"`
	IsVerified     bool      `json:"is_verified"`
	DeviceUUID     string    `json:"device_uuid"`
	Platform       string    `json:"platform,omitempty"`
	AppVersion     string    `json:"app_version,omitempty"`
	LoginTime      time.Time `json:"login_time"`
	SessionExpires time.Time `json:"session_expires"`
	LastLogin      time.Time `json:"last_login,omitzero"`
	LoginAttempts  int       `json:"login_attempts"`
}

type RegisterResponse struct {
	RegistrationID     string    `json:"registration_id"`
	PhoneNumber        string    `json:"phone_number"`
	DeviceUUID         string    `json:"device_uuid"`
	Platform           string    `json:"platform"`
	OTPSent            bool      `json:"otp_sent"`
	Otp                string    `json:"otp,omitempty"`
	OTPExpiryMinutes   int       `json:"otp_expiry_minutes"`
	Token              string    `json:"temp_token,omitempty"`
	TokenType          string    `json:"token_type,omitempty"`
	TokenExpiry        time.Time `json:"token_expiry,omitempty"`
	OTPCode            string    `json:"otp_code,omitempty"`
	NextStep           string    `json:"next_step"`
	RegistrationStatus string    `json:"registration_status"`
}
