package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

type PermissionString string

const (
	UserPayloadKey = contextKey("userPayload")
)

type Payload struct {
	UserID         string   `json:"user_id"`
	FullName       string   `json:"full_name"`
	PhoneNumber    string   `json:"phone_number"`
	UserRealm      string   `json:"user_realm"`
	UserRole       string   `json:"user_role"`
	Permissions    []string `json:"permissions"`
	PrimaryAuth    string   `json:"primary_auth"`
	SessionExpiry  string   `json:"session_expiry"`
	TokenType      string   `json:"token_type"` // device_lookup, verify_otp, login, permanent, forget_pin
	DeviceUUID     string   `json:"device_uuid,omitempty"`
	Phone          string   `json:"phone,omitempty"`
	OTPFor         string   `json:"otp_for,omitempty"`
	ResetSessionID string   `json:"reset_session_id,omitempty"`
	jwt.RegisteredClaims
}

// Token types constants
const (
	TokenTypeDeviceLookup = "device_lookup"
	TokenTypeVerifyOTP    = "verify_otp"
	TokenTypeLogin        = "login"
	TokenTypePermanent    = "permanent"
	TokenTypeForgetPin    = "forget_pin"
)

// Token permissions constants
const (
	PermissionDeviceLookup = "device_lookup"
	PermissionVerifyOTP    = "verify_otp"
	PermissionLogin        = "login"
	PermissionSetPin       = "set_pin"
	PermissionCreatePin    = "create_pin"
	PermissionResetPin     = "reset_pin"
	PermissionAccess       = "access"
)

// Implement the Valid method to satisfy the jwt.Claims interface
func (p *Payload) Valid() error {
	sessionExpiry, err := time.Parse(time.RFC3339, p.SessionExpiry)
	if err != nil {
		return fmt.Errorf("invalid session expiry format")
	}

	if time.Now().After(sessionExpiry) {
		return fmt.Errorf("token has expired")
	}
	return nil
}
