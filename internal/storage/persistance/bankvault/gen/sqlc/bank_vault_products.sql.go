package sqlc

import (
	constants "cbe-super-app-cps-action/internal/constants"
	utils "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

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
WHERE id = HEXTORAW(REPLACE(UPPER(:id), '-', '')) AND deleted_at IS NULL
RETURNING id INTO :result`

func (q *Queries) ActivateBankVault(ctx context.Context, id string) (string, error) {
	var result string
	res, err := q.db.ExecContext(
		ctx,
		activateBankVault,
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

const deactivateBankVault = `-- name: DeactivateBankVault :one
UPDATE bank_vault_products
SET is_active = 0, updated_at = SYSTIMESTAMP
WHERE id = HEXTORAW(REPLACE(UPPER(:id), '-', '')) AND deleted_at IS NULL
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
  description,
  currency,
  TO_CHAR(rate_bps) as rate_bps,
  method,
  frequency,
  lock_period,
  TO_CHAR(min_amount) as min_amount,
  TO_CHAR(max_amount) as max_amount,
  TO_CHAR(early_unlock_rate_bps) as early_unlock_rate_bps,
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
  AND (:method IS NULL OR UPPER(method) = UPPER(:method))
  AND (:frequency IS NULL OR UPPER(frequency) = UPPER(:frequency))
ORDER BY created_at DESC
OFFSET NVL(:offset, 0) ROWS
FETCH NEXT NVL(:limit, 50) ROWS ONLY
`

type FindBankVaultParams struct {
	IsActive  sql.NullBool         `json:"is_active"`
	IsDeleted sql.NullBool         `json:"is_deleted"`
	NameQuery sql.NullString       `json:"name"`
	Currency  sql.NullString       `json:"currency"`
	Method    NullAccrualMethod    `json:"method"`
	Frequency NullAccrualFrequency `json:"frequency"`
	Page      sql.NullInt64        `json:"page"`
	Limit     sql.NullInt64        `json:"limit"`
}

type FindBankVaultRow struct {
	ID                 string                     `json:"id"`
	Name               string                     `json:"name"`
	Description        string                     `json:"description"`
	Currency           string                     `json:"currency"`
	RateBps            decimal.Decimal            `json:"rate_bps"`
	Method             constants.AccrualMethod    `json:"method"`
	Frequency          constants.AccrualFrequency `json:"frequency"`
	LockPeriod         time.Duration              `json:"lock_period"`
	MinAmount          decimal.Decimal            `json:"min_amount"`
	MaxAmount          decimal.Decimal            `json:"max_amount"`
	EarlyUnlockRateBps sql.NullBool               `json:"early_unlock_rate_bps"`
	IsActive           sql.NullBool               `json:"is_active"`
	IsDeleted          sql.NullBool               `json:"is_deleted "`
	CreatedAt          time.Time                  `json:"created_at"`
	UpdatedAt          time.Time                  `json:"updated_at"`
	DeletedAt          sql.NullTime               `json:"deleted_at"`
	TotalCount         int64                      `json:"total_count"`
}

// --- AccrualMethod ---
func nullAccrualMethodToPtr(am NullAccrualMethod) *string {
	if !am.Valid {
		return (*string)(nil)
	}
	a := string(am.AccrualMethod)
	return &a
}

// --- AccrualFrequency ---
func nullAccrualFrequencyToPtr(af NullAccrualFrequency) *string {
	if !af.Valid {
		return (*string)(nil)
	}
	f := string(af.AccrualFrequency)
	return &f
}

func (q *Queries) FindBankVault(ctx context.Context, arg FindBankVaultParams) ([]FindBankVaultRow, error) {
	namePtr := utils.NullStringToPtrLike(arg.NameQuery)
	currencyPtr := utils.NullStringToPtr(arg.Currency)
	methodPtr := nullAccrualMethodToPtr(arg.Method)
	frequencyPtr := nullAccrualFrequencyToPtr(arg.Frequency)
	limitPtr := utils.NullInt64ToPtr(arg.Limit)
	offsetPtr := (arg.Page.Int64 - 1) * arg.Limit.Int64

	rows, err := q.db.QueryContext(ctx, findBankVault,
		sql.Named("is_active", arg.IsActive),
		sql.Named("is_deleted", arg.IsDeleted),
		sql.Named("name", namePtr),
		sql.Named("currency", currencyPtr),
		sql.Named("method", methodPtr),
		sql.Named("frequency", frequencyPtr),
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
			&i.Description,
			&i.Currency,
			&i.RateBps,
			&i.Method,
			&i.Frequency,
			&i.LockPeriod,
			&i.MinAmount,
			&i.MaxAmount,
			&i.EarlyUnlockRateBps,
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
  description,
  currency,
  rate_bps,
  method,
  frequency,
  lock_period,
  min_amount,
  max_amount,
  early_unlock_rate_bps,
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
		rateBps_num   godror.Number
		minAmount_num godror.Number
		maxAmount_num godror.Number
	)

	row := q.db.QueryRowContext(ctx, findBankVaultById, id)
	var i BankVaultProduct
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Currency,
		&rateBps_num,
		&i.Method,
		&i.Frequency,
		&i.LockPeriod,
		&minAmount_num,
		&maxAmount_num,
		&i.EarlyUnlockRateBps,
		&i.IsActive,
		&i.IsDeleted,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	i.RateBps, _ = decimal.NewFromString(rateBps_num.String())
	i.MinAmount, _ = decimal.NewFromString(minAmount_num.String())
	i.MaxAmount, _ = decimal.NewFromString(maxAmount_num.String())

	return i, err
}

const findBankVaultByName = `--name: FindBankVaultByName :one
SELECT
  id,
  name,
  description,
  currency,
  rate_bps,
  method,
  frequency,
  lock_period,
  min_amount,
  max_amount,
  early_unlock_rate_bps,
  is_active,
  created_at,
  updated_at,
  deleted_at
FROM bank_vault_products
WHERE UPPER(name) = UPPER(:name) AND deleted_at IS NULL
`

func (q *Queries) FindBankVaultByName(ctx context.Context, name string) (BankVaultProduct, error) {
	var (
		rateBps_num           godror.Number
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
		&i.Description,
		&i.Currency,
		&rateBps_num,
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
	i.ID = strings.ToUpper(hex.EncodeToString([]byte(i.ID)))
	i.RateBps, _ = decimal.NewFromString(rateBps_num.String())
	i.MinAmount, _ = decimal.NewFromString(minAmount_num.String())
	i.MaxAmount, _ = decimal.NewFromString(maxAmount_num.String())

	return i, err
}

const saveBankVault = `-- name: SaveBankVault :one
INSERT INTO bank_vault_products (
  id,
  name,
  description,
  currency,
  rate_bps,
  method,
  frequency,
  lock_period,
  min_amount,
  max_amount,
  early_unlock_rate_bps,
  is_active
) VALUES (
  :1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11, :12
)
RETURNING id INTO :13
`

type SaveBankVaultParams struct {
	ID                 string                     `json:"id"`
	Name               string                     `json:"name"`
	Description        string                     `json:"description"`
	Currency           string                     `json:"currency"`
	RateBps            decimal.Decimal            `json:"rate_bps"`
	Method             constants.AccrualMethod    `json:"method"`
	Frequency          constants.AccrualFrequency `json:"frequency"`
	LockPeriod         time.Duration              `json:"lock_period"`
	MinAmount          decimal.Decimal            `json:"min_amount"`
	MaxAmount          decimal.Decimal            `json:"max_amount"`
	EarlyUnlockRateBps sql.NullBool               `json:"early_unlock_rate_bps"`
	IsActive           sql.NullBool               `json:"is_active"`
}

func (q *Queries) SaveBankVault(ctx context.Context, arg SaveBankVaultParams) (string, error) {
	var id string

	_, err := q.db.ExecContext(ctx, saveBankVault,
		generateUUID(),
		arg.Name,
		arg.Description,
		strings.ToUpper(arg.Currency),
		arg.RateBps,
		string(arg.Method),
		string(arg.Frequency),
		int64(arg.LockPeriod),
		arg.MinAmount,
		arg.MaxAmount,
		utils.NullBoolToInt(arg.EarlyUnlockRateBps),
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
	description = COALESCE(:1, description),
	min_amount  = COALESCE(:2, min_amount),
	max_amount  = COALESCE(:3, max_amount),
	is_active   = COALESCE(:4, is_active),
	updated_at  = SYSTIMESTAMP
WHERE id = :5 AND deleted_at IS NULL
RETURNING id INTO :result
`

type UpdateBankVaultParams struct {
	Description *string          `json:"description"`
	MinAmount   *sql.NullFloat64 `json:"min_amount"`
	MaxAmount   *sql.NullFloat64 `json:"max_amount"`
	IsActive    sql.NullBool     `json:"is_active"`
	ID          string           `json:"id"`
}

func (q *Queries) UpdateBankVault(ctx context.Context, arg UpdateBankVaultParams) (string, error) {
	var idBytes []byte

	res, err := q.db.ExecContext(ctx, updateBankVault,
		arg.Description,
		arg.MinAmount,
		arg.MaxAmount,
		utils.NullBoolToIntPtr(arg.IsActive),
		arg.ID,
		sql.Out{Dest: &idBytes},
	)
	if err != nil {
		return "", fmt.Errorf("failed to update bank vault: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return "", fmt.Errorf("bank vault not found or already deleted")
	}
	id := strings.ToUpper(hex.EncodeToString(idBytes))

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
  rate_bps,
  method,
  frequency,
  lock_period,
  created_at,
  updated_at,
  deleted_at
FROM locked_vaults
WHERE product_id = :1 AND deleted_at IS NULL
`

//   early_unlock_rate_bps,

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
		return BankVaultProductWithLocks{}, fmt.Errorf("failed to query locks: %w", err)
	}
	defer rows.Close()

	var (
		principal_num godror.Number
		rateBps_num   godror.Number
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
			&rateBps_num,
			&l.Method,
			&l.Frequency,
			// &l.EarlyUnlockRateBps,
			&l.LockPeriod,
			&l.CreatedAt,
			&l.UpdatedAt,
			&l.DeletedAt,
		); err != nil {
			return BankVaultProductWithLocks{}, err
		}

		l.Principal, _ = decimal.NewFromString(principal_num.String())
		l.RateBps, _ = decimal.NewFromString(rateBps_num.String())

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
