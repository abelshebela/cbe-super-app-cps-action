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
	UserID string `json:"userid"`
	// UserCode             string   `json:"usercode"`
	FullName string `json:"fullname"`
	// OrganizationID       string   `json:"organizationid"`
	PhoneNumber string `json:"phonenumber"`
	// UserEmail            string   `json:"useremail"`
	UserRealm string `json:"userrealm"`
	UserRole  string `json:"userrole"`
	// IFBMember            bool     `json:"ifbmember"`
	// DeviceUUID           string   `json:"deviceuuid"`
	// UserDeviceLinkedDate string   `json:"userDeviceLinkedDate"`
	Permissions   []string `json:"permissions"`
	PrimaryAuth   string   `json:"primaryauth"`
	SessionExpiry string   `json:"sessionexpiry"`
	// PublicKey     string   `json:"publicKey"`
	jwt.RegisteredClaims
}

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
