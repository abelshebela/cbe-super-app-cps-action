package sqlc

import (
	utils "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/constants"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"github.com/godror/godror"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func generateUUID() string {
	return uuid.New().String()
}

const activateBankVault = `-- name: ActivateBankVault :one
UPDATE bank_vault_products
SET is_active = 1, updated_at = SYSTIMESTAMP
WHERE id = :1 AND deleted_at IS NULL
RETURNING id INTO :2`

func (q *Queries) ActivateBankVault(ctx context.Context, id string) (string, error) {
	var result string

	res, err := q.db.ExecContext(
		ctx,
		activateBankVault,
		id,
		sql.Out{Dest: &result},
	)
	if err != nil {
		return "", fmt.Errorf("failed to activate bank product: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return "", fmt.Errorf("no bank product found with id %s", id)
	}

	return result, nil
}

const deactivateBankVault = `-- name: DeactivateBankVault :one
UPDATE bank_vault_products
SET is_active = 0, updated_at = SYSTIMESTAMP
WHERE id = :1 AND deleted_at IS NULL
RETURNING id INTO :result`

func (q *Queries) DeactivateBankVault(ctx context.Context, id string) (string, error) {
	var result string
	res, err := q.db.ExecContext(
		ctx,
		deactivateBankVault,
		id,
		sql.Out{Dest: &result},
	)

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return "", fmt.Errorf("no bank product found with id %s", id)
	}

	if err != nil {
		return "", fmt.Errorf("failed to activate bank product: %w", err)
	}

	return result, nil
}

const deleteBankVault = `-- name: DeleteBankVault :one
UPDATE bank_vault_products
SET deleted_at = SYSTIMESTAMP, is_deleted = 1, updated_at = SYSTIMESTAMP
WHERE id = :1 AND deleted_at IS NULL
`

func (q *Queries) DeleteBankVault(ctx context.Context, id string) (string, error) {
	res, err := q.db.ExecContext(ctx, deleteBankVault, id)
	if err != nil {
		return "", fmt.Errorf("failed to delete bank product: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return "", fmt.Errorf("NO_BANK_PRODUCT_FOUND")
	}

	return id, nil
}

const findBankVault = `-- name: FindBankVault :many
SELECT
  id,
  name,
  currency,
  TO_CHAR(interest) as interest,
  method,
  frequency,
  lock_period,
  TO_CHAR(min_amount) as min_amount,
  TO_CHAR(max_amount) as max_amount,
  TO_CHAR(apply_interest_on_early_unlock) as apply_interest_on_early_unlock,
  is_active,
  is_deleted,
  created_at,
  updated_at,
  deleted_at,
  COUNT(*) OVER() AS total_count
FROM bank_vault_products
WHERE (:is_deleted IS NULL AND deleted_at IS NULL AND is_deleted = 0
  OR :is_deleted IS NOT NULL AND is_deleted = :is_deleted)
  AND (:is_active IS NULL OR is_active = :is_active)
  AND (:name IS NULL OR UPPER(name) LIKE :name)
  AND (:currency IS NULL OR currency = :currency)
ORDER BY created_at DESC
OFFSET NVL(:offset, 0) ROWS
FETCH NEXT NVL(:limit, 50) ROWS ONLY
`

type FindBankVaultParams struct {
	IsActive  sql.NullBool   `json:"is_active"`
	IsDeleted sql.NullBool   `json:"is_deleted"`
	NameQuery sql.NullString `json:"name"`
	Currency  sql.NullString `json:"currency"`
	// Method    NullAccrualMethod `json:"method"`
	// Frequency int64             `json:"frequency"`
	Page  sql.NullInt64 `json:"page"`
	Limit sql.NullInt64 `json:"limit"`
}

type FindBankVaultRow struct {
	ID                         string                  `json:"id"`
	Name                       string                  `json:"name"`
	Currency                   string                  `json:"currency"`
	Interest                   decimal.Decimal         `json:"interest"`
	Method                     constants.AccrualMethod `json:"method"`
	Frequency                  int64                   `json:"frequency"`
	LockPeriod                 float64                 `json:"lock_period"`
	MinAmount                  decimal.Decimal         `json:"min_amount"`
	MaxAmount                  decimal.Decimal         `json:"max_amount"`
	ApplyInterestOnEarlyUnlock sql.NullBool            `json:"apply_interest_on_early_unlock"`
	IsActive                   sql.NullBool            `json:"is_active"`
	IsDeleted                  sql.NullBool            `json:"is_deleted "`
	CreatedAt                  time.Time               `json:"created_at"`
	UpdatedAt                  time.Time               `json:"updated_at"`
	DeletedAt                  sql.NullTime            `json:"deleted_at"`
	TotalCount                 int64                   `json:"total_count"`
}

// --- AccrualMethod ---
func nullAccrualMethodToPtr(am NullAccrualMethod) *string {
	if !am.Valid {
		return (*string)(nil)
	}
	a := string(am.AccrualMethod)
	return &a
}

func (q *Queries) FindBankVault(ctx context.Context, arg FindBankVaultParams) ([]FindBankVaultRow, error) {
	namePtr := utils.NullStringToPtrLike(arg.NameQuery)
	currencyPtr := utils.NullStringToPtr(arg.Currency)
	// methodPtr := nullAccrualMethodToPtr(arg.Method)
	limitPtr := utils.NullInt64ToPtr(arg.Limit)
	offsetPtr := (arg.Page.Int64 - 1) * arg.Limit.Int64

	rows, err := q.db.QueryContext(ctx, findBankVault,
		sql.Named("is_active", arg.IsActive),
		sql.Named("is_deleted", arg.IsDeleted),
		sql.Named("name", namePtr),
		sql.Named("currency", currencyPtr),
		sql.Named("offset", offsetPtr),
		sql.Named("limit", limitPtr),
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []FindBankVaultRow
	for rows.Next() {
		var i FindBankVaultRow

		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Currency,
			&i.Interest,
			&i.Method,
			&i.Frequency,
			&i.LockPeriod,
			&i.MinAmount,
			&i.MaxAmount,
			&i.ApplyInterestOnEarlyUnlock,
			&i.IsActive,
			&i.IsDeleted,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.TotalCount,
		); err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

const findBankVaultById = `-- name: FindBankVaultById :one
SELECT
  id,
  name,
  currency,
  interest,
  method,
  frequency,
  lock_period,
  min_amount,
  max_amount,
  apply_interest_on_early_unlock,
  is_active,
  is_deleted,
  created_at,
  updated_at,
  deleted_at
FROM bank_vault_products
WHERE id = :id AND deleted_at IS NULL
`

func (q *Queries) FindBankVaultById(ctx context.Context, id string) (BankVaultProduct, error) {
	var (
		Interest_num  godror.Number
		minAmount_num godror.Number
		maxAmount_num godror.Number
	)

	row := q.db.QueryRowContext(ctx, findBankVaultById, id)
	var i BankVaultProduct
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Currency,
		&Interest_num,
		&i.Method,
		&i.Frequency,
		&i.LockPeriod,
		&minAmount_num,
		&maxAmount_num,
		&i.ApplyInterestOnEarlyUnlock,
		&i.IsActive,
		&i.IsDeleted,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	i.Interest, _ = decimal.NewFromString(Interest_num.String())
	i.MinAmount, _ = decimal.NewFromString(minAmount_num.String())
	i.MaxAmount, _ = decimal.NewFromString(maxAmount_num.String())

	return i, err
}

const findBankVaultByName = `--name: FindBankVaultByName :one
SELECT
  id,
  name,
  currency,
  interest,
  method,
  frequency,
  lock_period,
  min_amount,
  max_amount,
  apply_interest_on_early_unlock,
  is_active,
  created_at,
  updated_at,
  deleted_at
FROM bank_vault_products
WHERE UPPER(name) = UPPER(:name) AND deleted_at IS NULL
`

func (q *Queries) FindBankVaultByName(ctx context.Context, name string) (BankVaultProduct, error) {
	var (
		Interest_num          godror.Number
		minAmount_num         godror.Number
		maxAmount_num         godror.Number
		earlyUnlockFeeBps_num godror.Number
	)

	row := q.db.QueryRowContext(ctx,
		findBankVaultByName,
		sql.Named("name", name),
	)
	var i BankVaultProduct
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Currency,
		&Interest_num,
		&i.Method,
		&i.Frequency,
		&i.LockPeriod,
		&minAmount_num,
		&maxAmount_num,
		&earlyUnlockFeeBps_num,
		&i.IsActive,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	i.Interest, _ = decimal.NewFromString(Interest_num.String())
	i.MinAmount, _ = decimal.NewFromString(minAmount_num.String())
	i.MaxAmount, _ = decimal.NewFromString(maxAmount_num.String())

	return i, err
}

const saveBankVault = `-- name: SaveBankVault :one
INSERT INTO bank_vault_products (
  id,
  name,
  currency,
  interest,
  method,
  frequency,
  lock_period,
  min_amount,
  max_amount,
  apply_interest_on_early_unlock,
  is_active
) VALUES (
  :1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11
)
RETURNING id INTO :13
`

type SaveBankVaultParams struct {
	ID                         string          `json:"id"`
	Name                       string          `json:"name"`
	Currency                   string          `json:"currency"`
	Interest                   decimal.Decimal `json:"interest"`
	Method                     string          `json:"method"`
	Frequency                  int64           `json:"frequency"`
	LockPeriod                 float64         `json:"lock_period"`
	MinAmount                  decimal.Decimal `json:"min_amount"`
	MaxAmount                  decimal.Decimal `json:"max_amount"`
	ApplyInterestOnEarlyUnlock sql.NullBool    `json:"apply_interest_on_early_unlock"`
	IsActive                   sql.NullBool    `json:"is_active"`
}

func (q *Queries) SaveBankVault(ctx context.Context, arg SaveBankVaultParams) (string, error) {
	var id string

	_, err := q.db.ExecContext(ctx, saveBankVault,
		generateUUID(),
		arg.Name,
		strings.ToUpper(arg.Currency),
		arg.Interest,
		string(arg.Method),
		arg.Frequency,
		int64(arg.LockPeriod),
		arg.MinAmount,
		arg.MaxAmount,
		utils.NullBoolToInt(arg.ApplyInterestOnEarlyUnlock),
		utils.NullBoolToInt(arg.IsActive),
		sql.Out{Dest: &id},
	)
	if err != nil {
		return "", err
	}

	return id, nil
}

const updateBankVault = `-- name: UpdateBankVault :one
UPDATE bank_vault_products
SET
	min_amount  = COALESCE(:1, min_amount),
	max_amount  = COALESCE(:2, max_amount),
	is_active   = COALESCE(:3, is_active),
	updated_at  = SYSTIMESTAMP
WHERE id = :4 AND deleted_at IS NULL
RETURNING id INTO :result
`

type UpdateBankVaultParams struct {
	MinAmount *sql.NullFloat64 `json:"min_amount"`
	MaxAmount *sql.NullFloat64 `json:"max_amount"`
	IsActive  sql.NullBool     `json:"is_active"`
	ID        string           `json:"id"`
}

func (q *Queries) UpdateBankVault(ctx context.Context, arg UpdateBankVaultParams) (string, error) {
	var id string

	res, err := q.db.ExecContext(ctx, updateBankVault,
		arg.MinAmount,
		arg.MaxAmount,
		utils.NullBoolToIntPtr(arg.IsActive),
		arg.ID,
		sql.Out{Dest: &id},
	)
	if err != nil {
		return "", fmt.Errorf("failed to update bank vault: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return "", fmt.Errorf("bank vault not found or already deleted")
	}

	return id, nil
}

const findLocksByProductID = `-- name: FindLocksByProductID :many
SELECT
  id,
  customer_id,
  linked_account,
  product_id,
  principal,
  start_date,
  maturity_date,
  status,
  terms_version,
  terms_accepted_at,
  closed_at,
  interest,
  method,
  frequency,
  lock_period,
  created_at,
  updated_at,
  deleted_at
FROM vaults
WHERE product_id = :1 AND deleted_at IS NULL
`

//   apply_interest_on_early_unlock,

type BankVaultProductWithLocks struct {
	BankVaultProduct
	Locks []LockedVault `json:"locks"`
}

func (q *Queries) FindBankVaultAndLocks(ctx context.Context, id string) (BankVaultProductWithLocks, error) {
	product, err := q.FindBankVaultById(ctx, id)
	if err != nil {
		return BankVaultProductWithLocks{}, err
	}

	rows, err := q.db.QueryContext(ctx, findLocksByProductID, id)
	if err != nil {
		fmt.Printf("failed to query locked vaults from vaults table: %v", err)
		return BankVaultProductWithLocks{}, fmt.Errorf("failed to query locks: %w", err)
	}
	defer rows.Close()

	var (
		principal_num godror.Number
		Interest_num  godror.Number
	)

	locks := []LockedVault{}
	for rows.Next() {
		var l LockedVault
		if err := rows.Scan(
			&l.ID,
			&l.CustomerID,
			&l.LinkedAccount,
			&l.ProductID,
			&principal_num,
			&l.StartDate,
			&l.MaturityDate,
			&l.Status,
			&l.TermsVersion,
			&l.TermsAcceptedAt,
			&l.ClosedAt,
			&Interest_num,
			&l.Method,
			&l.Frequency,
			// &l.ApplyInterestOnEarlyUnlock,
			&l.LockPeriod,
			&l.CreatedAt,
			&l.UpdatedAt,
			&l.DeletedAt,
		); err != nil {
			return BankVaultProductWithLocks{}, err
		}

		l.Principal, _ = decimal.NewFromString(principal_num.String())
		l.Interest, _ = decimal.NewFromString(Interest_num.String())

		locks = append(locks, l)
	}

	if err := rows.Err(); err != nil {
		return BankVaultProductWithLocks{}, err
	}

	return BankVaultProductWithLocks{
		BankVaultProduct: product,
		Locks:            locks,
	}, nil
}

const getAllLockedVaults = `-- name: GetAllLockedVaults  :many
SELECT
  id,
  customer_id,
  linked_account,
  account_holder_name,
  product_id,
  principal,
  start_date,
  maturity_date,
  status,
  terms_version,
  terms_accepted_at,
  closed_at,
  min_amount,
  max_amount,
  interest,
  method,
  frequency,
  apply_interest_on_early_unlock,
  lock_period,
  created_at,
  updated_at,
  deleted_at,
  COUNT(*) OVER() AS total_count
FROM locked_vaults
WHERE deleted_at IS NULL
ORDER BY created_at DESC
OFFSET NVL(:offset, 0) ROWS
FETCH NEXT NVL(:limit, 50) ROWS ONLY
`

type ListLockedVaultParams struct {
	Status     NullLockedVaultStatus `json:"status"`
	ProductID  sql.NullString        `json:"product_id"`
	CustomerID sql.NullString        `json:"customer_id"`
	StartFrom  sql.NullTime          `json:"start_from"`
	MaturityTo sql.NullTime          `json:"maturity_to"`
	Page       sql.NullInt64         `json:"offset_count"`
	Limit      sql.NullInt64         `json:"limit_count"`
}

type ListLocksRow struct {
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
	ClosedAt                   sql.NullTime            `json:"closed_at"`
	MinAmount                  decimal.Decimal         `json:"min_amount"`
	MaxAmount                  decimal.Decimal         `json:"max_amount"`
	Interest                   decimal.Decimal         `json:"interest"`
	Method                     constants.AccrualMethod `json:"method"`
	Frequency                  int64                   `json:"frequency"`
	ApplyInterestOnEarlyUnlock NullBoolNumber          `json:"apply_interest_on_early_unlock"`
	LockPeriod                 int64                   `json:"lock_period"`
	CreatedAt                  time.Time               `json:"created_at"`
	UpdatedAt                  time.Time               `json:"updated_at"`
	DeletedAt                  sql.NullTime            `json:"deleted_at"`
	TotalCount                 int64                   `json:"total_count"`
}

func (q *Queries) GetAllLockedVaults(ctx context.Context, arg ListLockedVaultParams) ([]ListLocksRow, error) {
	offset := (arg.Page.Int64 - 1) * arg.Limit.Int64

	rows, err := q.db.QueryContext(ctx, getAllLockedVaults,
		sql.Named("offset", offset),
		sql.Named("limit", arg.Limit),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		principal_num godror.Number
		Interest_num  godror.Number
		minAmount_num godror.Number
		maxAmount_num godror.Number
	)

	items := []ListLocksRow{}
	for rows.Next() {
		var i ListLocksRow
		if err := rows.Scan(
			&i.ID,
			&i.CustomerID,
			&i.LinkedAccount,
			&i.AccountHolderName,
			&i.TransactionReference,
			&i.ProductID,
			&principal_num,
			&i.StartDate,
			&i.MaturityDate,
			&i.Status,
			&i.TermsVersion,
			&i.TermsAcceptedAt,
			&i.ClosedAt,
			&minAmount_num,
			&maxAmount_num,
			&Interest_num,
			&i.Method,
			&i.Frequency,
			&i.ApplyInterestOnEarlyUnlock,
			&i.LockPeriod,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.TotalCount,
		); err != nil {
			return nil, err
		}

		i.Principal, _ = decimal.NewFromString(principal_num.String())
		i.MinAmount, _ = decimal.NewFromString(minAmount_num.String())
		i.MaxAmount, _ = decimal.NewFromString(maxAmount_num.String())
		i.Interest, _ = decimal.NewFromString(Interest_num.String())

		items = append(items, i)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

const getAllGroupVaults = `--- name: GetAllGroupVaults :many
SELECT
  id,
  vault_name,
  vault_category,
  target_amount,
  vault_type,
  status,
  created_at,
  updated_at,
  deleted_at,
  COUNT(*) OVER() AS total_count
FROM vaults
WHERE vault_type = 'GROUP' AND deleted_at IS NULL
ORDER BY created_at DESC
OFFSET NVL(:offset, 0) ROWS
FETCH NEXT NVL(:limit, 50) ROWS ONLY
`

type GetAllGroupVaultsParams struct {
	Page  sql.NullInt64 `json:"offset_count"`
	Limit sql.NullInt64 `json:"limit_count"`
}

func (q *Queries) GetAllGroupVaults(ctx context.Context, arg GetAllGroupVaultsParams) ([]Vault, error) {
	offset := (arg.Page.Int64 - 1) * arg.Limit.Int64

	rows, err := q.db.QueryContext(ctx, getAllGroupVaults,
		sql.Named("offset", offset),
		sql.Named("limit", arg.Limit),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []Vault{}
	for rows.Next() {
		var i Vault

		var targetAmountNum godror.Number

		if err := rows.Scan(
			&i.ID,
			&i.VaultName,
			&i.VaultCategory,
			&targetAmountNum,
			// &i.IsDeadlock,
			&i.VaultType,
			// &i.InterestRate,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.TotalCount,
		); err != nil {
			return nil, err
		}

		i.TargetAmount, _ = decimal.NewFromString(targetAmountNum.String())

		items = append(items, i)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

const getAllTransactions = `--- name: GetAllTransactions :many
SELECT
  id,
  transaction_id,
  ft_number,
  debit_branch_code,
  debit_district_code,
  debit_user_id,
  debit_account_number,
  debit_account_holder_name,
  credit_user_id,
  credit_account_number,
  credit_account_holder_name,
  institution_code,
  institution_name,
  currency,
  service_fee,
  tip_amount,
  paid_amount,
  vat,
  amount,
  total_amount,
  external_reference,
  transaction_reason,
  transaction_type,
  transaction_status,
  is_ifb,
  is_reversed,
  paid_at,
  reversed_at,
  metadata,
  created_at,
  last_modified_at,
  COUNT(*) OVER() AS total_count
FROM transaction
WHERE ft_number = :1 AND deleted_at IS NULL
ORDER BY created_at DESC
OFFSET NVL(:offset, 0) ROWS
FETCH NEXT NVL(:limit, 50) ROWS ONLY
`

type GetAllTransactionsParams struct {
	Page  sql.NullInt64 `json:"offset_count"`
	Limit sql.NullInt64 `json:"limit_count"`
}

func (q *Queries) GetAllTransactions(ctx context.Context, arg GetAllTransactionsParams) ([]model.Transaction, error) {
	offset := (arg.Page.Int64 - 1) * arg.Limit.Int64

	rows, err := q.db.QueryContext(ctx, getAllTransactions,
		sql.Named("offset", offset),
		sql.Named("limit", arg.Limit),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []model.Transaction{}
	for rows.Next() {
		var i model.Transaction

		var serviceFee godror.Number
		var tipAmount godror.Number
		var paidAmount godror.Number
		var vat godror.Number
		var amount godror.Number
		var totalAmount godror.Number

		var metadata sql.NullString

		if err := rows.Scan(
			&i.ID,
			&i.TransactionID,
			&i.FTNumber,
			&i.DebitBranchCode,
			&i.DebitDistrictCode,
			&i.DebitUserID,
			&i.DebitAccountNumber,
			&i.DebitAccountHolderName,
			&i.CreditUserID,
			&i.CreditAccountNumber,
			&i.CreditAccountHolderName,
			&i.InstitutionCode,
			&i.InstitutionName,
			&i.Currency,
			&serviceFee,
			&tipAmount,
			&paidAmount,
			&vat,
			&amount,
			&totalAmount,
			&i.ExternalReference,
			&i.TransactionReason,
			&i.TransactionType,
			&i.TransactionStatus,
			&i.IsIFB,
			&i.IsReversed,
			&i.PaidAt,
			&i.ReversedAt,
			&metadata,
			&i.CreatedAt,
			&i.LastModifiedAt,
			&i.TotalCount,
		); err != nil {
			return nil, err
		}

		// convert decimal fields
		i.ServiceFee, _ = decimal.NewFromString(serviceFee.String())
		i.TipAmount, _ = decimal.NewFromString(tipAmount.String())
		i.PaidAmount, _ = decimal.NewFromString(paidAmount.String())
		i.VAT, _ = decimal.NewFromString(vat.String())
		i.Amount, _ = decimal.NewFromString(amount.String())
		i.TotalAmount, _ = decimal.NewFromString(totalAmount.String())

		// metadata
		if metadata.Valid {
			i.Metadata = json.RawMessage(metadata.String)
		}

		items = append(items, i)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
