package bankvault

import (
	"time"

	"github.com/shopspring/decimal"
)

type BankVaultProductResponse struct {
	ID                         string          `json:"id"`
	Name                       string          `json:"name"`
	Currency                   string          `json:"currency"`
	Interest                   decimal.Decimal `json:"interest"`
	Method                     string          `json:"method"`
	Frequency                  int64           `json:"frequency"`
	LockPeriod                 float64         `json:"lock_period"`
	MinAmount                  decimal.Decimal `json:"min_amount"`
	MaxAmount                  decimal.Decimal `json:"max_amount"`
	ApplyInterestOnEarlyUnlock bool            `json:"apply_interest_on_early_unlock"`
	IsActive                   bool            `json:"is_active"`
	IsDeleted                  bool            `json:"is_deleted"`
	CreatedAt                  time.Time       `json:"created_at"`
	UpdatedAt                  time.Time       `json:"updated_at"`
	DeletedAt                  *time.Time      `json:"deleted_at,omitempty"`
}

//	type LockedVault struct {
//		ID                         string                  `json:"id"`
//		CustomerID                 string                  `json:"customer_id"`
//		LinkedAccount              string                  `json:"linked_account"`
//		AccountHolderName          string                  `json:"account_holder_name"`
//		TransactionReference       string                  `json:"transaction_reference"`
//		ProductID                  string                  `json:"product_id"`
//		Principal                  decimal.Decimal         `json:"principal"`
//		StartDate                  time.Time               `json:"start_date"`
//		MaturityDate               time.Time               `json:"maturity_date"`
//		Status                     constants.VaultStatus   `json:"status"`
//		TermsVersion               string                  `json:"terms_version"`
//		TermsAcceptedAt            time.Time               `json:"terms_accepted_at"`
//		MinAmount                  decimal.Decimal         `json:"min_amount"`
//		MaxAmount                  decimal.Decimal         `json:"max_amount"`
//		Interest                   decimal.Decimal         `json:"interest"`
//		Method                     constants.AccrualMethod `json:"method"`
//		Frequency                  string                  `json:"frequency"`
//		ApplyInterestOnEarlyUnlock *bool                   `json:"apply_interest_on_early_unlock,omitempty"`
//		LockPeriod                 int64                   `json:"lock_period"`
//		CreatedAt                  time.Time               `json:"created_at"`
//		UpdatedAt                  time.Time               `json:"updated_at"`
//		ClosedAt                   *time.Time              `json:"closed_at,omitempty"`
//		DeletedAt                  *time.Time              `json:"deleted_at,omitempty"`
//	}
// type GroupVaultResponse struct {
// 	ID              string     `json:"id"`
// 	VaultName       string     `json:"vault_name"`
// 	VaultCategory   string     `json:"vault_category"`
// 	Purpose         string     `json:"purpose"`
// 	TargetAmount    string     `json:"target_amount"`
// 	CollectedAmount string     `json:"collected_amount"`
// 	MemberCount     int32      `json:"member_count"`
// 	EndDate         time.Time  `json:"end_date"`
// 	VaultType       string     `json:"vault_type"`
// 	Status          string     `json:"status"`
// 	AdminUserID     string     `json:"admin_user_id"`
// 	Recurrence      string     `json:"recurrence"`
// 	NextRun         time.Time  `json:"next_run"`
// 	Reminder        *bool      `json:"reminder,omitempty"`
// 	TCVersion       string     `json:"tc_version"`
// 	OccVersion      int32      `json:"occ_version"`
// 	CreatedAt       time.Time  `json:"created_at"`
// 	UpdatedAt       time.Time  `json:"updated_at"`
// 	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
// }
