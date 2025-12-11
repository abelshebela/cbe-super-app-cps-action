package model

import (
	"cbe-super-app-cps-action/internal/constants"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/godror/godror"
	"github.com/shopspring/decimal"
)

type BankVaultProduct struct {
	ID                         string                  `json:"id" bson:"id"`
	Name                       string                  `json:"name" bson:"name"`
	Currency                   string                  `json:"currency" bson:"currency"`
	Interest                   decimal.Decimal         `json:"interest" bson:"interest"`
	Method                     constants.AccrualMethod `json:"method" bson:"method"`
	Frequency                  int64                   `json:"frequency" bson:"frequency"`
	LockPeriod                 float64                 `json:"lock_period" bson:"lock_period"`
	MinAmount                  decimal.Decimal         `json:"min_amount" bson:"min_amount"`
	MaxAmount                  decimal.Decimal         `json:"max_amount" bson:"max_amount"`
	ApplyInterestOnEarlyUnlock bool                    `json:"apply_interest_on_early_unlock" bson:"apply_interest_on_early_unlock"`
	IsActive                   bool                    `json:"is_active" bson:"is_active"`
	CreatedAt                  time.Time               `json:"created_at" bson:"created_at"`
	UpdatedAt                  time.Time               `json:"updated_at" bson:"updated_at"`
	DeletedAt                  *time.Time              `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
	CreatedBy                  string                  `json:"created_by" bson:"created_by"`
	UpdatedBy                  string                  `json:"updated_by" bson:"updated_by"`
	IsDeleted                  bool                    `json:"is_deleted" bson:"is_deleted"`
}

type BankVaultProductPatch struct {
	// Description *string          `json:"description,omitempty" bson:"description,omitempty"`
	MinAmount *decimal.Decimal `json:"min_amount,omitempty" bson:"min_amount,omitempty"`
	MaxAmount *decimal.Decimal `json:"max_amount,omitempty" bson:"max_amount,omitempty"`
	IsActive  *bool            `json:"is_active,omitempty" bson:"is_active,omitempty"`
}

type UpdateBankVault struct {
	// Description *string   `json:"description" bson:"description"`
	MinAmount *float64  `json:"min_amount" bson:"min_amount"`
	MaxAmount *float64  `json:"max_amount" bson:"max_amount"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	UpdatedBy *string   `json:"updated_by" bson:"updated_by"`
}

type LockedVault struct {
	ID                         string                  `json:"id"`
	CustomerID                 string                  `json:"customer_id"`
	LinkedAccount              string                  `json:"linked_account"`
	AccountHolderName          string                  `json:"account_holder_name"`
	TransactionReference       string                  `json:"transaction_reference"`
	ProductID                  string                  `json:"product_id"`
	Principal                  decimal.Decimal         `json:"principal"`
	StartDate                  time.Time               `json:"start_date"`
	MaturityDate               time.Time               `json:"maturity_date"`
	Status                     constants.VaultStatus   `json:"status"`
	TermsVersion               string                  `json:"terms_version"`
	TermsAcceptedAt            time.Time               `json:"terms_accepted_at"`
	MinAmount                  decimal.Decimal         `json:"min_amount"`
	MaxAmount                  decimal.Decimal         `json:"max_amount"`
	Interest                   decimal.Decimal         `json:"interest"`
	Method                     constants.AccrualMethod `json:"method"`
	Frequency                  int64                   `json:"frequency"`
	ApplyInterestOnEarlyUnlock *bool                   `json:"apply_interest_on_early_unlock,omitempty"`
	LockPeriod                 string                  `json:"lock_period"`
	CreatedAt                  time.Time               `json:"created_at"`
	UpdatedAt                  time.Time               `json:"updated_at"`
	ClosedAt                   *time.Time              `json:"closed_at,omitempty"`
	DeletedAt                  *time.Time              `json:"deleted_at,omitempty"`
}

type GroupVault struct {
	ID              string          `json:"id"`
	VaultName       string          `json:"vault_name"`
	VaultCategory   string          `json:"vault_category"`
	Purpose         string          `json:"purpose"`
	TargetAmount    decimal.Decimal `json:"target_amount"`
	CollectedAmount decimal.Decimal `json:"collected_amount"`
	MemberCount     int32           `json:"member_count"`
	EndDate         time.Time       `json:"end_date"`
	VaultType       string          `json:"vault_type"`
	Status          string          `json:"status"`
	AdminUserID     string          `json:"admin_user_id"`
	Recurrence      string          `json:"recurrence"`
	NextRun         time.Time       `json:"next_run"`
	Reminder        *bool           `json:"reminder,omitempty"`
	TCVersion       string          `json:"tc_version"`
	OccVersion      int32           `json:"occ_version"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	DeletedAt       *time.Time      `json:"deleted_at,omitempty"`
	TotalCount      int64           `json:"total_count,omitempty"`
}

type (
	TransactionStatus string
	Currency          string
	TransactionType   string
)

const (
	PAID     TransactionStatus = "paid"
	Complete TransactionStatus = "complete"
	PENDING  TransactionStatus = "pending"
	FAILED   TransactionStatus = "failed"
	EXPIRED  TransactionStatus = "expired"
	TIMEOUT  TransactionStatus = "timeout"
	REVERSED TransactionStatus = "reversed"
)

const (
	CBE TransactionType = "cbe"
	OWN
	AirtimeTopUP TransactionType = "Airtime Top-up"
	TopUp        TransactionType = "topup"
)

type Transaction struct {
	ID                      string            `json:"id"`
	TransactionID           string            `json:"transaction_id"`
	FTNumber                string            `json:"ft_number"`
	DebitBranchCode         string            `json:"debit_branch_code"`
	DebitDistrictCode       string            `json:"debit_district_code"`
	DebitUserID             string            `json:"debit_user_id"`
	DebitAccountNumber      string            `json:"debit_account_number"`
	DebitAccountHolderName  string            `json:"debit_account_holder_name"`
	CreditUserID            string            `json:"credit_user_id"`
	CreditAccountNumber     string            `json:"credit_account_number"`
	CreditAccountHolderName string            `json:"credit_account_holder_name"`
	InstitutionCode         string            `json:"institution_code"`
	InstitutionName         string            `json:"institution_name"`
	Currency                Currency          `json:"currency"`
	ServiceFee              decimal.Decimal   `json:"service_fee"`
	TipAmount               decimal.Decimal   `json:"tip_amount"`
	PaidAmount              decimal.Decimal   `json:"paid_amount"`
	VAT                     decimal.Decimal   `json:"vat"`
	Amount                  decimal.Decimal   `json:"amount"`
	TotalAmount             decimal.Decimal   `json:"total_amount"`
	ExternalReference       string            `json:"external_reference,omitempty"`
	TransactionReason       string            `json:"transaction_reason,omitempty"`
	TransactionType         TransactionType   `json:"transaction_type"`
	TransactionStatus       TransactionStatus `json:"transaction_status"`
	IsIFB                   bool              `json:"is_ifb"`
	IsReversed              bool              `json:"is_reversed"`
	PaidAt                  time.Time         `json:"paid_at,omitzero"`
	ReversedAt              time.Time         `json:"reversed_at,omitzero"`
	Metadata                json.RawMessage   `json:"metadata,omitempty"`
	CreatedAt               time.Time         `json:"created_at"`
	LastModifiedAt          time.Time         `json:"last_modified_at,omitzero"`
	TotalCount              int64             `json:"total_count,omitempty"`
}

type TransactionModel struct {
	ID                      string            `json:"id"`
	TransactionID           string            `json:"transaction_id"`
	FTNumber                string            `json:"ft_number"`
	DebitBranchCode         string            `json:"debit_branch_code"`
	DebitDistrictCode       string            `json:"debit_district_code"`
	DebitUserID             string            `json:"debit_user_id"`
	DebitAccountNumber      string            `json:"debit_account_number"`
	DebitAccountHolderName  string            `json:"debit_account_holder_name"`
	CreditUserID            string            `json:"credit_user_id"`
	CreditAccountNumber     string            `json:"credit_account_number"`
	CreditAccountHolderName string            `json:"credit_account_holder_name"`
	InstitutionCode         string            `json:"institution_code"`
	InstitutionName         string            `json:"institution_name"`
	Currency                Currency          `json:"currency"`
	ServiceFee              godror.Number     `json:"service_fee"`
	TipAmount               godror.Number     `json:"tip_amount"`
	PaidAmount              godror.Number     `json:"paid_amount"`
	VAT                     godror.Number     `json:"vat"`
	Amount                  godror.Number     `json:"amount"`
	TotalAmount             godror.Number     `json:"total_amount"`
	ExternalReference       string            `json:"external_reference,omitempty"`
	TransactionReason       string            `json:"transaction_reason,omitempty"`
	TransactionType         TransactionType   `json:"transaction_type"`
	TransactionStatus       TransactionStatus `json:"transaction_status"`
	IsIFB                   string            `json:"is_ifb"`
	IsReversed              sql.NullBool      `json:"is_reversed"`
	PaidAt                  sql.NullTime      `json:"paid_at,omitzero"`
	ReversedAt              sql.NullTime      `json:"reversed_at,omitzero"`
	Metadata                *godror.JSON      `json:"metadata,omitempty"`
	CreatedAt               sql.NullTime      `json:"created_at"`
	LastModifiedAt          sql.NullTime      `json:"last_modified_at,omitzero"`
}
