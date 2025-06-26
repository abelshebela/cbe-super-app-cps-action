package users

import "time"

type FullName struct {
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
}

type LinkedAccountDetail struct {
	AccountNumber     string `json:"account_number"`
	AccountBranchCode string `json:"account_branch_code"`
	LinkedBranch      string `json:"linked_branch"`
	IsAccountActive   bool   `json:"is_account_active"`
	LinkedStatus      bool   `json:"linked_status"`
	CurrencyCode      string `json:"currency_code"`
}

type DeviceInfo struct {
	DeviceUUID string `json:"device_uuid"`
	AppVersion string `json:"app_version"`
}

type LoginPIN struct {
	PIN              string    `json:"pin" bson:"pin"`
	PINHistory       [4]string `json:"pin_history" bson:"pin_histroy"`
	LastPINCreatedAt time.Time `json:"last_pin_created_at" bson:"last_pin_created_at"`
}

type User struct {
	ID        string
	Email     string
	FullName  string
	IsDeleted bool
	Device    *DeviceInfo `json:"device,omitempty" bson:"device,omitempty"`
	LoginPIN  LoginPIN    `json:"login_pin,omitempty" bson:"login_pin,omitempty"`
}

type UserEmail struct {
	ID    string
	Email string `json:"email" bson:"email"`
}

type LinkedAccountResponse struct {
	UserID         string                `json:"user_id"`
	FullName       string                `json:"full_name"`
	LinkedAccounts []LinkedAccountDetail `json:"linked_accounts"`
}

type GenerateOTPRequest struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

type OTPRequest struct {
	UserID string
	Email  string
}

type OTPVerification struct {
	UserID string
	Email  string
	OTP    string
}

type OTPRecord struct {
	ID        string `json:"id" bson:"_id"`
	UserCode  string
	UserID    string
	Email     string
	OTP       string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type ChangePinRequest struct {
	UserID string `json:"user_id"`
	OldPin string `json:"old_pin"`
	NewPin string `json:"new_pin"`
}
