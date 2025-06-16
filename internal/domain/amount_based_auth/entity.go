package amount_based_auth_domain

import "time"

type Method string

const (
	OTP        Method = "otp"
	OPEN       Method = "open"
	OTPANDPINT Method = "otp+pin"
)

type AuthTier struct {
	MinAmount    uint64    `json:"minAmount" bson:"minAmount"`
	MaxAmount    uint64    `json:"maxAmount" bson:"maxAmount"`
	Method       string    `json:"method" bson:"method"`
	Enabled      bool      `json:"enabled" bson:"enabled"`
	IsDeleted    bool      `json:"isDeleted" bson:"isDeleted"`
	CreatedAt    time.Time `json:"createdAt" bson:"createdAt"`
	LastModified time.Time `json:"lastModified" bson:"lastModified"`
}
