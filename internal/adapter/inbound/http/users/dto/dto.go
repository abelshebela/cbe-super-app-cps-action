package dto

import (
	"fmt"
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// DeviceLookupRequest represents the request for device lookup
type DeviceLookupRequest struct {
	Platform          string            `json:"platform"`
	AppVersion        string            `json:"app_version"`
	DeviceUUID        string            `json:"device_uuid"`
	SourceApp         string            `json:"source_app"`
	InstallationDate  string            `json:"installation_date"`
	AdditionalHeaders map[string]string `json:"additional_headers,omitempty"`
}

// Validate validates DeviceLookupRequest
func (r DeviceLookupRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.DeviceUUID, validation.Required, validation.Length(10, 100)),
		validation.Field(&r.Platform, validation.Required, validation.In("android", "ios", "web")),
		validation.Field(&r.AppVersion, validation.Required),
	)
}

// DeviceLookupResponse represents the response for device lookup
type DeviceLookupResponse struct {
	Token        string `json:"token"`
	DeviceUUID   string `json:"device_uuid"`
	Platform     string `json:"platform"`
	AppVersion   string `json:"app_version"`
	SourceApp    string `json:"source_app"`
	IsRegistered bool   `json:"is_registered"`
}

// FetchLinkedAccountsRequest represents the request for fetching linked accounts
type FetchLinkedAccountsRequest struct {
	UserID string `json:"user_id"`
}

// Validate validates FetchLinkedAccountsRequest
func (r FetchLinkedAccountsRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.UserID, validation.Required),
	)
}

// FetchLinkedAccountsResponse represents the response for fetching linked accounts
type FetchLinkedAccountsResponse struct {
	UserID         string                `json:"user_id"`
	LinkedAccounts []LinkedAccountDetail `json:"linked_accounts"`
	TotalAccounts  int                   `json:"total_accounts"`
}

// LinkedAccountDetail represents a linked account detail
type LinkedAccountDetail struct {
	AccountNumber     string `json:"account_number"`
	AccountBranchCode string `json:"account_branch_code"`
	LinkedBranch      string `json:"linked_branch"`
	IsAccountActive   bool   `json:"is_account_active"`
	LinkedStatus      bool   `json:"linked_status"`
	CurrencyCode      string `json:"currency_code"`
}

// OTPRequest represents the request for OTP generation
type OTPRequest struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	OTPFor string `json:"otp_for"`
}

// Validate validates OTPRequest
func (r OTPRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.UserID, validation.Required),
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(&r.OTPFor, validation.Required, validation.In("change_email", "pin_set", "login", "registration")),
	)
}

// OTPResponse represents the response for OTP generation
type OTPResponse struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	OTPSent   bool   `json:"otp_sent"`
	ExpiresIn int    `json:"expires_in_minutes"`
}

// OTPVerification represents the request for OTP verification
type OTPVerification struct {
	UserID string `json:"user_id"`
	OTP    string `json:"otp"`
	OTPFor string `json:"otp_for"`
}

// Validate validates OTPVerification
func (r OTPVerification) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.UserID, validation.Required),
		validation.Field(&r.OTP, validation.Required, validation.Length(6, 6), validation.Match(regexp.MustCompile(`^\d{6}$`))),
		validation.Field(&r.OTPFor, validation.Required),
	)
}

// ChangePinRequest represents the request for changing PIN
type ChangePinRequest struct {
	UserID string `json:"user_id"`
	OldPin string `json:"old_pin"`
	NewPin string `json:"new_pin"`
}

// Validate validates ChangePinRequest
func (r ChangePinRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.UserID, validation.Required),
		validation.Field(&r.OldPin, validation.Required, validation.Length(6, 6), validation.Match(regexp.MustCompile(`^\d{6}$`))),
		validation.Field(&r.NewPin, validation.Required, validation.Length(6, 6), validation.Match(regexp.MustCompile(`^\d{6}$`))),
		validation.Field(&r.NewPin, validation.By(func(value interface{}) error {
			pin := value.(string)
			if pin == r.OldPin {
				return fmt.Errorf("new PIN must be different from old PIN")
			}
			return nil
		})),
	)
}

// VerifyOtpRequest represents the request for OTP verification
type VerifyOtpRequest struct {
	Otp    string `json:"otp_code"`
	OtpFor string `json:"otp_for,omitempty"`
}

// Validate validates VerifyOtpRequest
func (r VerifyOtpRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Otp, validation.Required, validation.Length(6, 6), validation.Match(regexp.MustCompile(`^\d{6}$`))),
		validation.Field(&r.OtpFor, validation.In("pin_set", "login", "registration", "change_email", "forget_pin")),
	)
}

// VerifyOtpResponse represents the response for OTP verification
type VerifyOtpResponse struct {
	Token       string `json:"token"`
	UserID      string `json:"user_id"`
	PhoneNumber string `json:"phone_number"`
	Verified    bool   `json:"verified"`
}

// SetPinRequest represents the request for setting PIN
type SetPinRequest struct {
	NewPin     string  `json:"new_pin"`
	Otp        string  `json:"otp"`
	DeviceUUID *string `json:"device_uuid,omitempty"`
	UserRealm  string  `json:"user_realm,omitempty"`
	OtpFor     string  `json:"otp_for,omitempty"`
}

// Validate validates SetPinRequest
func (r SetPinRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.NewPin, validation.Required, validation.Length(6, 6), validation.Match(regexp.MustCompile(`^\d{6}$`))),
		validation.Field(&r.Otp, validation.Required, validation.Length(6, 6), validation.Match(regexp.MustCompile(`^\d{6}$`))),
		validation.Field(&r.UserRealm, validation.In("member", "admin", "merchant")),
		validation.Field(&r.OtpFor, validation.In("pin_set", "login", "registration")),
	)
}

// SetPinResponse represents the response for setting PIN
type SetPinResponse struct {
	Token    string `json:"token"`
	UserID   string `json:"user_id"`
	PinSet   bool   `json:"pin_set"`
	SetTime  string `json:"set_time"`
	NextStep string `json:"next_step"`
}

// RegisterRequest represents the request for user registration
type RegisterRequest struct {
	Phone      string `json:"phone"`
	DeviceUUID string `json:"device_uuid"`
	Platform   string `json:"platform"`
	FullName   string `json:"full_name,omitempty"`
	Email      string `json:"email,omitempty"`
}

// Validate validates RegisterRequest
func (r RegisterRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Phone, validation.Required, validation.By(validatePhoneNumber)),
		validation.Field(&r.DeviceUUID, validation.Required, validation.Length(10, 100)),
		validation.Field(&r.Platform, validation.Required, validation.In("android", "ios", "web")),
		validation.Field(&r.FullName, validation.Length(0, 100)),
		validation.Field(&r.Email, validation.When(r.Email != "", is.Email)),
	)
}

// RegisterResponse represents the response for user registration
type RegisterResponse struct {
	RegistrationID string `json:"registration_id"`
	Phone          string `json:"phone"`
	DeviceUUID     string `json:"device_uuid"`
	Platform       string `json:"platform"`
	OTPSent        bool   `json:"otp_sent"`
	ExpiresIn      int    `json:"expires_in_minutes"`
	NextStep       string `json:"next_step"`
}

// LoginRequest represents the request for user login
type LoginRequest struct {
	Phone      string `json:"phone"`
	DeviceUUID string `json:"device_uuid"`
	Pin        string `json:"pin"`
}

// Validate validates LoginRequest
func (r LoginRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Phone, validation.Required, validation.By(validatePhoneNumber)),
		validation.Field(&r.DeviceUUID, validation.Required, validation.Length(10, 100)),
		validation.Field(&r.Pin, validation.Required, validation.Length(6, 6), validation.Match(regexp.MustCompile(`^\d{6}$`))),
	)
}

// LoginResponse represents the response for user login
type LoginResponse struct {
	Token          string `json:"token"`
	UserID         string `json:"user_id"`
	UserCode       string `json:"user_code"`
	FullName       string `json:"full_name"`
	PhoneNumber    string `json:"phone_number"`
	KYCLevel       uint8  `json:"kyc_level"`
	IsVerified     bool   `json:"is_verified"`
	DeviceUUID     string `json:"device_uuid"`
	Platform       string `json:"platform"`
	AppVersion     string `json:"app_version"`
	LoginTime      string `json:"login_time"`
	SessionExpires string `json:"session_expires"`
}

// ForgetPinSendOtpRequest represents the request for sending OTP for PIN reset
type ForgetPinSendOtpRequest struct {
	Phone      string `json:"phone"`
	DeviceUUID string `json:"device_uuid"`
}

// Validate validates ForgetPinSendOtpRequest
func (r ForgetPinSendOtpRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Phone, validation.Required, validation.By(validatePhoneNumber)),
		validation.Field(&r.DeviceUUID, validation.Required, validation.Length(10, 100)),
	)
}

// ForgetPinSendOtpResponse represents the response for sending OTP for PIN reset
type ForgetPinSendOtpResponse struct {
	PhoneNumber      string `json:"phone_number"`
	DeviceUUID       string `json:"device_uuid"`
	OTPSent          bool   `json:"otp_sent"`
	OTPExpiryMinutes int    `json:"otp_expiry_minutes"`
	ResetSessionID   string `json:"reset_session_id"`
	NextStep         string `json:"next_step"`
}

// ResetPinRequest represents the request for resetting PIN
type ResetPinRequest struct {
	ResetSessionID string `json:"reset_session_id"`
	Phone          string `json:"phone"`
	DeviceUUID     string `json:"device_uuid"`
	OTP            string `json:"otp"`
	NewPin         string `json:"new_pin"`
}

// Validate validates ResetPinRequest
func (r ResetPinRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ResetSessionID, validation.Required),
		validation.Field(&r.Phone, validation.Required, validation.By(validatePhoneNumber)),
		validation.Field(&r.DeviceUUID, validation.Required, validation.Length(10, 100)),
		validation.Field(&r.OTP, validation.Required, validation.Length(6, 6), validation.Match(regexp.MustCompile(`^\d{6}$`))),
		validation.Field(&r.NewPin, validation.Required, validation.Length(6, 6), validation.Match(regexp.MustCompile(`^\d{6}$`))),
	)
}

// ResetPinResponse represents the response for resetting PIN
type ResetPinResponse struct {
	UserID           string   `json:"user_id"`
	UserCode         string   `json:"user_code"`
	FullName         string   `json:"full_name"`
	PhoneNumber      string   `json:"phone_number"`
	PinReset         bool     `json:"pin_reset"`
	ResetTime        string   `json:"reset_time"`
	AccessRestricted bool     `json:"access_restricted"`
	Restrictions     []string `json:"restrictions"`
	NextStep         string   `json:"next_step"`
}

// CompleteRegistrationRequest represents the request for completing registration
type CompleteRegistrationRequest struct {
	RegistrationID string `json:"registration_id"`
	Phone          string `json:"phone"`
	DeviceUUID     string `json:"device_uuid"`
	Platform       string `json:"platform"`
	FullName       string `json:"full_name"`
	OTP            string `json:"otp"`
}

// Validate validates CompleteRegistrationRequest
func (r CompleteRegistrationRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.RegistrationID, validation.Required),
		validation.Field(&r.Phone, validation.Required, validation.By(validatePhoneNumber)),
		validation.Field(&r.DeviceUUID, validation.Required, validation.Length(10, 100)),
		validation.Field(&r.Platform, validation.Required, validation.In("android", "ios", "web")),
		validation.Field(&r.FullName, validation.Required, validation.Length(1, 100)),
		validation.Field(&r.OTP, validation.Required, validation.Length(6, 6), validation.Match(regexp.MustCompile(`^\d{6}$`))),
	)
}

// CompleteRegistrationResponse represents the response for completing registration
type CompleteRegistrationResponse struct {
	UserID      string `json:"user_id"`
	UserCode    string `json:"user_code"`
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email,omitempty"`
	Token       string `json:"token"`
	Registered  bool   `json:"registered"`
	NextStep    string `json:"next_step"`
}

// ProfilePictureResponse represents the response for profile picture upload
type ProfilePictureResponse struct {
	Message string `json:"message"`
	URL     string `json:"url"`
}

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Error   string      `json:"error"`
	Details interface{} `json:"details,omitempty"`
}

// SuccessResponse represents a standard success response
type SuccessResponse struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

// Helper validation functions

// validatePhoneNumber validates phone number format
func validatePhoneNumber(value interface{}) error {
	phone, ok := value.(string)
	if !ok {
		return fmt.Errorf("phone number must be a string")
	}

	// Remove any non-digit characters
	phone = regexp.MustCompile(`[^\d]`).ReplaceAllString(phone, "")

	// Check if it's a valid Ethiopian phone number
	if len(phone) < 10 || len(phone) > 12 {
		return fmt.Errorf("phone number must be between 10 and 12 digits")
	}

	// Ethiopian phone numbers typically start with 7 or 9
	if !strings.HasPrefix(phone, "7") && !strings.HasPrefix(phone, "9") {
		return fmt.Errorf("phone number must start with 7 or 9")
	}

	return nil
}

// validatePIN validates PIN format and strength
func validatePIN(value interface{}) error {
	pin, ok := value.(string)
	if !ok {
		return fmt.Errorf("PIN must be a string")
	}

	if len(pin) != 6 {
		return fmt.Errorf("PIN must be exactly 6 digits")
	}

	// Check if PIN contains only digits
	if !regexp.MustCompile(`^\d{6}$`).MatchString(pin) {
		return fmt.Errorf("PIN must contain only digits")
	}

	// Check for common weak PINs
	weakPINs := []string{"000000", "111111", "123456", "654321", "999999"}
	for _, weakPIN := range weakPINs {
		if pin == weakPIN {
			return fmt.Errorf("PIN is too weak, please choose a different PIN")
		}
	}

	// Check for sequential digits
	if isSequential(pin) {
		return fmt.Errorf("PIN cannot be sequential digits")
	}

	// Check for repeated digits
	if isRepeated(pin) {
		return fmt.Errorf("PIN cannot be repeated digits")
	}

	return nil
}

// isSequential checks if a string contains sequential digits
func isSequential(s string) bool {
	if len(s) < 3 {
		return false
	}

	// Check ascending sequence
	ascending := true
	descending := true

	for i := 1; i < len(s); i++ {
		if s[i] != s[i-1]+1 {
			ascending = false
		}
		if s[i] != s[i-1]-1 {
			descending = false
		}
	}

	return ascending || descending
}

// isRepeated checks if a string contains repeated digits
func isRepeated(s string) bool {
	if len(s) < 2 {
		return false
	}

	first := s[0]
	for i := 1; i < len(s); i++ {
		if s[i] != first {
			return false
		}
	}

	return true
}
