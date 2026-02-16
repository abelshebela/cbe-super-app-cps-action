package model

import (
	"time"
)

type VaultTiers struct {
	ID           string `json:"id" bson:"_id,omitempty" gorm:"primaryKey"`
	Name         string `json:"name" bson:"name"`
	TierInterest string `json:"tier_interest" bson:"tier_interest"`
	MinAmount    string `json:"min" bson:"min"`
	MaxAmount    string `json:"max" bson:"max"`
}

type VaultCategory struct {
	ID               string       `json:"id" bson:"_id,omitempty" gorm:"primaryKey"`
	Name             string       `json:"name" bson:"name"`
	CoverImageURL    string       `json:"cover_image_url" bson:"cover_image_url"`
	InterestType     string       `json:"interest_type" bson:"interest_type"`
	CategoryInterest string       `json:"category_interest" bson:"category_interest"`
	Deadlock         bool         `json:"deadlock" bson:"deadlock"`
	Tiers            []VaultTiers `json:"tiers" bson:"tiers"`
	IsActive         bool         `json:"is_active" bson:"is_active"`
	CreatedAt        time.Time    `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time    `json:"updated_at" bson:"updated_at"`
}

type Withdrawal struct {
	ID                    string    `json:"id" bson:"_id,omitempty" gorm:"primaryKey"`
	LockedVaultID         string    `json:"locked_vault_id" bson:"locked_vault_id"`
	Amount                string    `json:"amount" bson:"amount"`
	WithdrawerName        string    `json:"withdrawer_name" bson:"withdrawer_name"`
	WithdrawerPhoneNumber string    `json:"withdrawer_phone_number" bson:"withdrawer_phone_number"`
	Status                string    `json:"status" bson:"status"`
	IsActive              bool      `json:"is_active" bson:"is_active"`
	CreatedAt             time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" bson:"updated_at"`
}

type WithdrawalStatusHistory struct {
	ID           string    `json:"id" bson:"_id,omitempty" gorm:"primaryKey"`
	WithdrawalID string    `json:"withdrawal_id" bson:"withdrawal_id"`
	Status       string    `json:"status" bson:"status"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
}

type VaultTransaction struct {
	ID                      int64     `json:"id"`
	Amount                  float64   `json:"amount"`
	BalanceAfter            *float64  `json:"balance_after,omitempty"`
	CreatedAt               time.Time `json:"created_at"`
	CreditAccountHolderName *string   `json:"credit_account_holder_name,omitempty"`
	CreditAccountNumber     *string   `json:"credit_account_number,omitempty"`
	Currency                string    `json:"currency"`
	DebitAccountHolderName  *string   `json:"debit_account_holder_name,omitempty"`
	DebitAccountNumber      *string   `json:"debit_account_number,omitempty"`
	DebitUserID             *string   `json:"debit_user_id,omitempty"`
	ExternalReference       *string   `json:"external_reference,omitempty"`
	FTNumber                *string   `json:"ft_number,omitempty"`
	InterestDelta           *float64  `json:"interest_delta,omitempty"`
	IsIFB                   int64     `json:"is_ifb"`
	PaidAmount              float64   `json:"paid_amount"`
	PrincipalDelta          *float64  `json:"principal_delta,omitempty"`
	ReceiptLink             *string   `json:"receipt_link,omitempty"`
	ReferenceID             *string   `json:"reference_id,omitempty"`
	ReferenceType           *string   `json:"reference_type,omitempty"`
	ServiceFee              float64   `json:"service_fee"`
	TotalAmount             float64   `json:"total_amount"`
	TransactionID           string    `json:"transaction_id"`
	TransactionReason       *string   `json:"transaction_reason,omitempty"`
	TransactionType         string    `json:"transaction_type"`
	VAT                     float64   `json:"vat"`
	VaultID                 *int64    `json:"vault_id,omitempty"`
	VaultTxType             *string   `json:"vault_tx_type,omitempty"`
}
