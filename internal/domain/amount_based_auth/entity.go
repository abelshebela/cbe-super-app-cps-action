package amount_based_auth_domain

import "time"

type Method string

const (
	OPEN      Method = "OPEN"
	PIN       Method = "PIN"
	OTPANDPIN Method = "OTP_PIN"
)

type AuthTier struct {
	ID           string    `json:"id" bson:"id"`
	MinAmount    uint64    `json:"min_amount" bson:"min_amount"`
	MaxAmount    uint64    `json:"max_amount" bson:"max_amount"`
	Method       Method    `json:"method" bson:"method"`
	Enabled      bool      `json:"enabled" bson:"enabled"`
	IsDeleted    bool      `json:"is_deleted" bson:"is_deleted"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	LastModified time.Time `json:"last_modified" bson:"last_modified"`
}
