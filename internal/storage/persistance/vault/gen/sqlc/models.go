package sqlc

import (
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"time"
)

type VaultTiers struct {
	ID           string `json:"id" bson:"_id,omitempty" gorm:"primaryKey"`
	Name         string `json:"name" bson:"name"`
	CategoryID   string `json:"category_id" bson:"category_id"`
	TierInterest string `json:"tier_interest" bson:"tier_interest"`
	MinAmount    string `json:"min_amount" bson:"min_amount"`
	MaxAmount    string `json:"max_amount" bson:"max_amount"`
}

type VaultCategory struct {
	ID               string              `json:"id" bson:"_id,omitempty" gorm:"primaryKey"`
	Name             string              `json:"name" bson:"name"`
	CoverImageURL    string              `json:"cover_image_url" bson:"cover_image_url"`
	InterestType     string              `json:"interest_type" bson:"interest_type"`
	CategoryInterest string              `json:"category_interest" bson:"category_interest"`
	Deadlock         bool                `json:"deadlock" bson:"deadlock"`
	IsActive         bool                `json:"is_active" bson:"is_active"`
	CreatedAt        time.Time           `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time           `json:"updated_at" bson:"updated_at"`
	TotalCount       int64               `json:"total_count"`
	Tiers            []imodel.VaultTiers `json:"tiers" bson:"tiers"`
}

type VaultTransaction struct {
	ID                      int64     `json:"id"`
	Amount                  float64   `json:"amount"`
	FTNumber                *string   `json:"ft_number,omitempty"`
	BalanceAfter            *float64  `json:"balance_after,omitempty"`
	CreatedAt               time.Time `json:"created_at"`
	CreditAccountHolderName *string   `json:"credit_account_holder_name,omitempty"`
	CreditAccountNumber     *string   `json:"credit_account_number,omitempty"`
	Currency                string    `json:"currency"`
	DebitAccountHolderName  *string   `json:"debit_account_holder_name,omitempty"`
	DebitAccountNumber      *string   `json:"debit_account_number,omitempty"`
	DebitUserID             *string   `json:"debit_user_id,omitempty"`
	ExternalReference       *string   `json:"external_reference,omitempty"`
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
	TotalCount              int64     `json:"total_count"`
}
