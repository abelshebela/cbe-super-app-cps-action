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
	ID            string       `json:"id" bson:"_id,omitempty" gorm:"primaryKey"`
	Name          string       `json:"name" bson:"name"`
	CoverImageURL string       `json:"cover_image_url" bson:"cover_image_url"`
	InterestType  string       `json:"interest_type" bson:"interest_type"`
	Deadlock      bool         `json:"deadlock" bson:"deadlock"`
	Tiers         []VaultTiers `json:"tiers" bson:"tiers"`
	IsActive      bool         `json:"is_active" bson:"is_active"`
	CreatedAt     time.Time    `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at" bson:"updated_at"`
}

// type Vault struct {
// 	ID                         string     `db:"id" json:"id"`
// 	Name                       string     `db:"name" json:"name"`
// 	UserID                     string     `db:"user_id" json:"user_id"`
// 	Balance                    float64    `db:"balance" json:"balance"`
// 	VaultCategoryID            *string    `db:"vault_category_id" json:"vault_category_id,omitempty"`
// 	AssociatedVaultProductID   *string    `db:"associated_vault_product_id" json:"associated_vault_product_id,omitempty"`
// 	VaultGLAccountNo           *string    `db:"vault_gl_account_no" json:"vault_gl_account_no,omitempty"`
// 	InterestTierID             *string    `db:"interest_tier_id" json:"interest_tier_id,omitempty"`
// 	MaturityDisbursementOption *string    `db:"maturity_disbursement_option" json:"maturity_disbursement_option,omitempty"`
// 	VaultType                  string     `db:"vault_type" json:"vault_type"`
// 	LifecycleStatus            string     `db:"lifecycle_status" json:"lifecycle_status"`
// 	EndDate                    *time.Time `db:"end_date" json:"end_date,omitempty"`
// 	TargetAmount               *float64   `db:"target_amount" json:"target_amount,omitempty"`
// 	RemindersEnabled           bool       `db:"reminders_enabled" json:"reminders_enabled"`
// 	ReminderLeadDays           *int       `db:"reminder_lead_days" json:"reminder_lead_days,omitempty"`
// 	DefaultFrequency           *string    `db:"default_frequency" json:"default_frequency,omitempty"`
// 	Purpose                    *string    `db:"purpose" json:"purpose,omitempty"`
// 	TermsAccepted              bool       `db:"terms_accepted" json:"terms_accepted"`
// 	TermsAcceptedAt            *time.Time `db:"terms_accepted_at" json:"terms_accepted_at,omitempty"`
// 	IsDeadlocked               bool       `db:"is_deadlocked" json:"is_deadlocked"`
// 	MaturityStatus             string     `db:"maturity_status" json:"maturity_status"`
// 	PrincipalBalance           float64    `db:"principal_balance" json:"principal_balance"`
// 	InterestPosted             float64    `db:"interest_posted" json:"interest_posted"`
// 	InterestPending            float64    `db:"interest_pending" json:"interest_pending"`
// 	MinBalancePeriod           *float64   `db:"min_balance_period" json:"min_balance_period,omitempty"`
// 	PenaltyAccrued             *float64   `db:"penalty_accrued" json:"penalty_accrued,omitempty"`
// 	MaturityAmount             *float64   `db:"maturity_amount" json:"maturity_amount,omitempty"`
// }

// type DeadlockRequest struct {
// 	ID                  string    `json:"id" bson:"_id,omitempty" gorm:"primaryKey"`
// 	VaultName           string    `json:"vault_name" bson:"vault_name"`
// 	VaultID             string    `json:"vault_id" bson:"vault_id"`
// 	VaultType           string    `json:"vault_type" bson:"vault_type"`
// 	MemberName          string    `json:"member_name" bson:"member_name"`
// 	MemberAccountNumber string    `json:"member_account_number" bson:"member_account_number"`
// 	Status              string    `json:"status" bson:"status"`
// 	CreatedAt           time.Time `json:"created_at" bson:"created_at"`
// 	UpdatedAt           time.Time `json:"updated_at" bson:"updated_at"`
// }

// type WithdrawalStatusHistory struct {
// 	ID           string    `json:"id" bson:"_id,omitempty" gorm:"primaryKey"`
// 	WithdrawalID string    `json:"withdrawal_id" bson:"withdrawal_id"`
// 	Status       string    `json:"status" bson:"status"`
// 	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
// }

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

type DeadlockRequest struct {
	ID                  string    `json:"id" bson:"_id,omitempty" gorm:"primaryKey"`
	VaultName           string    `json:"vault_name" bson:"vault_name"`
	VaultID             string    `json:"vault_id" bson:"vault_id"`
	VaultType           string    `json:"vault_type" bson:"vault_type"`
	MemberName          string    `json:"member_name" bson:"member_name"`
	MemberAccountNumber string    `json:"member_account_number" bson:"member_account_number"`
	Amount              string    `json:"amount" bson:"amount"`
	Category            string    `json:"category" bson:"category"`
	NumberOfMembers     string    `json:"number_of_members" bson:"number_of_members"`
	Duration            string    `json:"duration" bson:"duration"`
	EndDate             string    `json:"end_date" bson:"end_date"`
	Status              string    `json:"status" bson:"status"`
	CreatedAt           time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" bson:"updated_at"`
}
