package sqlc

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"time"

	"cbe-super-app-cps-action/internal/constants"

	"github.com/godror/godror"
	"github.com/shopspring/decimal"
)

type AccrualMethod string

func (e *AccrualMethod) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = AccrualMethod(s)
	case string:
		*e = AccrualMethod(s)
	default:
		return fmt.Errorf("unsupported scan type for AccrualMethod: %T", src)
	}
	return nil
}

type NullAccrualMethod struct {
	AccrualMethod AccrualMethod `json:"accrual_method"`
	Valid         bool          `json:"valid"` // Valid is true if AccrualMethod is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullAccrualMethod) Scan(value interface{}) error {
	if value == nil {
		ns.AccrualMethod, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.AccrualMethod.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullAccrualMethod) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.AccrualMethod), nil
}

type LockedVaultStatus string

const (
	LockedVaultStatusACTIVE        LockedVaultStatus = "ACTIVE"
	LockedVaultStatusMATURED       LockedVaultStatus = "MATURED"
	LockedVaultStatusUNLOCKEDEARLY LockedVaultStatus = "UNLOCKED_EARLY"
	LockedVaultStatusPAIDOUT       LockedVaultStatus = "PAID_OUT"
)

func (e *LockedVaultStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = LockedVaultStatus(s)
	case string:
		*e = LockedVaultStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for LockedVaultStatus: %T", src)
	}
	return nil
}

type NullLockedVaultStatus struct {
	LockedVaultStatus LockedVaultStatus `json:"locked_vault_status"`
	Valid             bool              `json:"valid"` // Valid is true if LockedVaultStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullLockedVaultStatus) Scan(value interface{}) error {
	if value == nil {
		ns.LockedVaultStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.LockedVaultStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullLockedVaultStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.LockedVaultStatus), nil
}

type ReceiptKind string

const (
	ReceiptKindLOCK   ReceiptKind = "LOCK"
	ReceiptKindUNLOCK ReceiptKind = "UNLOCK"
	ReceiptKindPAYOUT ReceiptKind = "PAYOUT"
)

func (e *ReceiptKind) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = ReceiptKind(s)
	case string:
		*e = ReceiptKind(s)
	default:
		return fmt.Errorf("unsupported scan type for ReceiptKind: %T", src)
	}
	return nil
}

type NullReceiptKind struct {
	ReceiptKind ReceiptKind `json:"receipt_kind"`
	Valid       bool        `json:"valid"` // Valid is true if ReceiptKind is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullReceiptKind) Scan(value interface{}) error {
	if value == nil {
		ns.ReceiptKind, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.ReceiptKind.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullReceiptKind) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.ReceiptKind), nil
}

type TransactionType string

const (
	TransactionTypePAYOUT TransactionType = "PAYOUT"
	TransactionTypeUNLOCK TransactionType = "UNLOCK"
)

func (e *TransactionType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = TransactionType(s)
	case string:
		*e = TransactionType(s)
	default:
		return fmt.Errorf("unsupported scan type for TransactionType: %T", src)
	}
	return nil
}

type NullTransactionType struct {
	TransactionType TransactionType `json:"transaction_type"`
	Valid           bool            `json:"valid"` // Valid is true if TransactionType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullTransactionType) Scan(value interface{}) error {
	if value == nil {
		ns.TransactionType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.TransactionType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullTransactionType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.TransactionType), nil
}

type BankVaultProduct struct {
	ID                         string          `json:"id"`
	Name                       string          `json:"name"`
	Description                string          `json:"description"`
	Currency                   string          `json:"currency"`
	Interest                   decimal.Decimal `json:"interest"`
	Method                     string          `json:"method"`
	Frequency                  int64           `json:"frequency"`
	LockPeriod                 float64         `json:"lock_period"`
	MinAmount                  decimal.Decimal `json:"min_amount"`
	MaxAmount                  decimal.Decimal `json:"max_amount"`
	ApplyInterestOnEarlyUnlock sql.NullBool    `json:"apply_interest_on_early_unlock"`
	IsActive                   sql.NullBool    `json:"is_active"`
	IsDeleted                  sql.NullBool    `json:"is_deleted"`
	CreatedAt                  time.Time       `json:"created_at"`
	UpdatedAt                  time.Time       `json:"updated_at"`
	DeletedAt                  sql.NullTime    `json:"deleted_at"`
}

type LockedVault struct {
	ID                         string                `json:"id"`
	CustomerID                 string                `json:"customer_id"`
	LinkedAccount              string                `json:"linked_account"`
	ProductID                  string                `json:"product_id"`
	Principal                  decimal.Decimal       `json:"principal"`
	StartDate                  time.Time             `json:"start_date"`
	MaturityDate               time.Time             `json:"maturity_date"`
	Status                     constants.VaultStatus `json:"status"`
	TermsVersion               string                `json:"terms_version"`
	TermsAcceptedAt            time.Time             `json:"terms_accepted_at"`
	Interest                   decimal.Decimal       `json:"interest"`
	Method                     string                `json:"method"`
	Frequency                  int64                 `json:"frequency"`
	ApplyInterestOnEarlyUnlock decimal.Decimal       `json:"apply_interest_on_early_unlock"`
	LockPeriod                 time.Duration         `json:"lock_period"`
	CreatedAt                  time.Time             `json:"created_at"`
	UpdatedAt                  time.Time             `json:"updated_at"`
	ClosedAt                   sql.NullTime          `json:"closed_at"`
	DeletedAt                  sql.NullTime          `json:"deleted_at"`
}

type NullBoolNumber struct {
	sql.NullBool
}

func (n *NullBoolNumber) Scan(src any) error {
	if src == nil {
		n.Valid = false
		return nil
	}

	switch v := src.(type) {
	case bool:
		n.Bool = v
		n.Valid = true
		return nil
	case int64:
		n.Bool = v != 0
		n.Valid = true
		return nil
	case float64:
		n.Bool = v != 0
		n.Valid = true
		return nil
	case []byte:
		s := strings.TrimSpace(string(v))
		if s == "" {
			n.Valid = false
			return nil
		}
		if i, err := strconv.Atoi(s); err == nil {
			n.Bool = i != 0
		} else {
			n.Bool = strings.EqualFold(s, "true") || strings.HasPrefix(s, "1")
		}
		n.Valid = true
		return nil
	case string:
		s := strings.TrimSpace(v)
		if s == "" {
			n.Valid = false
			return nil
		}
		if i, err := strconv.Atoi(s); err == nil {
			n.Bool = i != 0
		} else {
			n.Bool = strings.EqualFold(s, "true") || strings.HasPrefix(s, "1")
		}
		n.Valid = true
		return nil
	case godror.Number:
		s := strings.TrimSpace(v.String())
		if s == "" {
			n.Valid = false
			return nil
		}
		if i, err := strconv.Atoi(s); err == nil {
			n.Bool = i != 0
		} else {
			n.Bool = strings.HasPrefix(s, "1")
		}
		n.Valid = true
		return nil
	default:
		// Fallback: try to use fmt.Sprintf then parse
		s := strings.TrimSpace(strings.ToLower(strings.TrimSuffix(strings.TrimPrefix(strings.TrimSpace(strings.Trim(strings.Trim(fmt.Sprintf("%v", src), "\n"), "\r")), " "), " ")))
		if s == "" {
			n.Valid = false
			return nil
		}
		if i, err := strconv.Atoi(s); err == nil {
			n.Bool = i != 0
		} else {
			n.Bool = strings.EqualFold(s, "true") || strings.HasPrefix(s, "1")
		}
		n.Valid = true
		return nil
	}
}

type Vault struct {
	ID            string          `json:"id"`
	VaultName     string          `json:"vault_name"`
	VaultCategory string          `json:"vault_category,omitempty"`
	CategoryID    string          `json:"category_id,omitempty"`
	Purpose       string          `json:"purpose,omitempty"`
	TargetAmount  decimal.Decimal `json:"target_amount,omitempty"`
	EndDate       *time.Time      `json:"end_date,omitempty"`
	IsDeadlock    *bool           `json:"is_deadlock,omitempty"`
	VaultType     string          `json:"vault_type,omitempty"`
	InterestRate  string          `json:"interest_rate"`
	Status        string          `json:"status,omitempty"`
	AdminUserID   string          `json:"admin_user_id"`
	Frequency     string          `json:"frequency,omitempty"`
	NextRun       sql.NullString  `json:"next_run,omitempty"`
	CreatedAt     sql.NullString  `json:"created_at"`
	UpdatedAt     sql.NullString  `json:"updated_at,omitempty"`
	DeletedAt     sql.NullString  `json:"deleted_at,omitempty"`
	Reminder      string          `json:"reminder,omitempty"`
	TotalCount    int64           `json:"total_count,omitempty"`
}
