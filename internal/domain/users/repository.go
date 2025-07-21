package users

import (
	"cbe-super-app-member-users/internal/adapter/outbound/model"
	"cbe-super-app-member-users/pkgs/entities"

	"cbe-super-app-member-users/pkgs/entities/type_definition"
	"context"
)

type UserRepository interface {
	FindByID(ctx context.Context, id string) (*User, error)
	FindActiveLinkedAccounts(ctx context.Context, userID string) ([]LinkedAccountDetail, error)
	FindByEmail(ctx context.Context, email string) (*UserEmail, error)
	StoreOTP(ctx context.Context, otp *OTPRecord) error
	FindOTP(ctx context.Context, userID, email string) (*OTPRecord, error)
	UpdateProfileTheme(ctx context.Context, id string, themeType string) (*model.User, error)
	UpdateUserEmail(ctx context.Context, userID, email string) error
	UpdateProfileImageURL(ctx context.Context, id string, imageURL string) error
	UnlinkDevice(ctx context.Context, userID string, deviceID string) error
	ChangePin(ctx context.Context, userID string, loginPIN type_definition.LoginPIN) error
	GetOneHQ(ctx context.Context, req map[string]interface{}) (*HQ, error)
	GetOneUser(ctx context.Context, req map[string]interface{}) (*entities.User, error)
	// OTP CRUD
	CreateOtp(ctx context.Context, otp *OTPRecord) error
	GetOtpByID(ctx context.Context, id string) (*OTPRecord, error)
	UpdateOtp(ctx context.Context, otp *OTPRecord) error
	UpdateOneUser(ctx context.Context, filter map[string]interface{}, req map[string]interface{}) error
	DeleteOtp(ctx context.Context, phone, otpCode, otpFor string) error
	DeleteOtpHard(ctx context.Context, phone, otpCode, otpFor string) error
	// Registration methods
	FindUserByPhone(ctx context.Context, phone string) (*User, error)
	FindUserByPhoneForLogin(ctx context.Context, phone string, pin string) (*User, error)
	FindUserByDevice(ctx context.Context, deviceUUID string) (*User, error)
	FindPendingRegistration(ctx context.Context, userID, deviceUUID string) (*RegistrationRecord, error)
	FindPendingRegistrationByID(ctx context.Context, registrationID string) (*RegistrationRecord, error)
	CreatePendingRegistration(ctx context.Context, registration *RegistrationRecord) error
	DeletePendingRegistration(ctx context.Context, registrationID string) error
	UpdatePendingRegistration(ctx context.Context, registration *RegistrationRecord) error
	CreateUser(ctx context.Context, user *User) error
	// Login methods
	IncrementLoginAttempts(ctx context.Context, userID string) error
	ResetLoginAttempts(ctx context.Context, userID string) error
	UpdateLastLogin(ctx context.Context, userID string) error
	// PIN Reset methods
	CreatePinResetSession(ctx context.Context, session *PinResetSession) error
	FindPinResetSession(ctx context.Context, sessionID string) (*PinResetSession, error)
	FindPinResetSessionByPhone(ctx context.Context, phone, deviceUUID string) (*PinResetSession, error)
	UpdatePinResetSession(ctx context.Context, session *PinResetSession) error
	DeletePinResetSession(ctx context.Context, sessionID string) error
	IncrementPinResetAttempts(ctx context.Context, sessionID string) error
	ResetPinResetAttempts(ctx context.Context, sessionID string) error
}
