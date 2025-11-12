package model

import (
	"time"
)

type SitotaTransaction struct {
	ID                     string     `bson:"id,omitempty" json:"id"`
	SenderName             string     `bson:"sender_name" json:"sender_name"`
	SenderAccountNumber    string     `bson:"sender_account_number" json:"sender_account_number"`
	SenderPhoneNumber      string     `bson:"sender_phone_number" json:"sender_phone_number"`
	RecipientName          string     `bson:"recipient_name" json:"recipient_name"`
	RecipientAccountNumber string     `bson:"recipient_account_number" json:"recipient_account_number"`
	RecipientPhoneNumber   string     `bson:"recipient_phone_number" json:"recipient_phone_number"`
	SitotaAmount           float64    `bson:"sitota_amount" json:"sitota_amount"`
	GLAccountNumber        string     `bson:"gl_account_number" json:"gl_account_number"`
	Status                 string     `bson:"status" json:"status"`
	ClaimedAt              *time.Time `bson:"claimed_at,omitempty" json:"claimed_at"`
	CreatedAt              time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt              time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt              *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}
