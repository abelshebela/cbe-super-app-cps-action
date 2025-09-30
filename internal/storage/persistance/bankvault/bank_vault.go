package bankvault

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/persistance/bankvault/gen/sqlc"
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	utils "cbe-super-app-cps-action/pkgs/utils"

	"github.com/shopspring/decimal"
	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type bankVaultRepositary struct {
	db      *sql.DB
	queries *sqlc.Queries
	logger  shared_utils.Logger
}

func sqlBoolToBool(nb sql.NullBool) bool {
	if nb.Valid {
		return nb.Bool
	}
	return false
}

func NewBankVaultRepository(db *sql.DB, logger shared_utils.Logger) storage.BankVaultRepository {
	return &bankVaultRepositary{
		db:      db,
		queries: sqlc.New(db),
		logger:  logger,
	}

}

func (r *bankVaultRepositary) ExecuteInTransaction(ctx context.Context, fn func(txRepo storage.BankVaultRepository) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	txRepo := &bankVaultRepositary{
		db:      r.db,
		queries: r.queries.WithTx(tx),
		logger:  r.logger,
	}

	if err := fn(txRepo); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *bankVaultRepositary) Create(ctx context.Context, product *model.BankVaultProduct) (string, error) {
	// Convert time.Duration to nanoseconds
	lockPeriodNanoseconds := product.LockPeriod.Nanoseconds()

	return r.createBankVaultCustom(ctx, product, lockPeriodNanoseconds)
}

func (r *bankVaultRepositary) createBankVaultCustom(ctx context.Context, product *model.BankVaultProduct, lockPeriodNanoseconds int64) (string, error) {
	// Check for duplicate name first
	_, err := r.queries.FindBankVaultByName(ctx, product.Name)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return "", err
		}
	} else {
		return "", fmt.Errorf("DUPLICATE_BANK_PRODUCT")
	}

	// Custom Oracle insert query
	query := `
		INSERT INTO bank_vault_products (
			name, description, currency, rate_bps, method, frequency, 
			lock_period, min_amount, max_amount, early_unlock_fee_bps, is_active
		) VALUES (
			:1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11
		)
		RETURNING id INTO :12
	`

	var idBytes []byte
	_, err = r.db.ExecContext(ctx, query,
		product.Name,
		product.Description,
		strings.ToUpper(product.Currency),
		product.RateBps,
		string(product.Method),
		string(product.Frequency),
		lockPeriodNanoseconds, // Use nanoseconds like microservice
		product.MinAmount,
		product.MaxAmount,
		product.EarlyUnlockFeeBps,
		utils.NullBoolToInt(sql.NullBool{Bool: product.IsActive, Valid: true}),
		sql.Out{Dest: &idBytes},
	)

	if err != nil {
		return "", err
	}

	id := strings.ToUpper(hex.EncodeToString(idBytes))
	return id, nil
}

func (r *bankVaultRepositary) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.BankVaultProduct], error) {

	query := `
		SELECT
			id, name, description, currency, rate_bps, method, frequency, 
			lock_period, min_amount, max_amount, early_unlock_fee_bps, 
			is_active, is_deleted, created_at, updated_at, deleted_at,
			COUNT(*) OVER() AS total_count
		FROM bank_vault_products
		WHERE 1=1
	`

	args := []interface{}{}
	argIndex := 1

	if v, ok := filterParam.Filters["is_active"].(bool); ok {
		query += fmt.Sprintf(" AND is_active = :%d", argIndex)
		args = append(args, v)
		argIndex++
	}
	if v, ok := filterParam.Filters["name"].(string); ok && v != "" {
		query += fmt.Sprintf(" AND UPPER(name) LIKE UPPER(:%d)", argIndex)
		args = append(args, "%"+v+"%")
		argIndex++
	}
	if v, ok := filterParam.Filters["currency"].(string); ok && v != "" {
		query += fmt.Sprintf(" AND currency = :%d", argIndex)
		args = append(args, v)
		argIndex++
	}
	if v, ok := filterParam.Filters["method"].(string); ok && v != "" {
		query += fmt.Sprintf(" AND UPPER(method) = UPPER(:%d)", argIndex)
		args = append(args, v)
		argIndex++
	}
	if v, ok := filterParam.Filters["frequency"].(string); ok && v != "" {
		query += fmt.Sprintf(" AND UPPER(frequency) = UPPER(:%d)", argIndex)
		args = append(args, v)
		argIndex++
	}

	query += " ORDER BY created_at DESC"

	// Add pagination
	limit := 50
	offset := 0
	if filterParam.PerPage > 0 {
		limit = filterParam.PerPage
	}
	if filterParam.Page > 0 {
		offset = (filterParam.Page - 1) * limit
	}

	query += fmt.Sprintf(" OFFSET %d ROWS FETCH NEXT %d ROWS ONLY", offset, limit)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]*model.BankVaultProduct, 0)
	var total int64

	for rows.Next() {
		var (
			rateBps_num           string
			minAmount_num         string
			maxAmount_num         string
			earlyUnlockFeeBps_num string
			isActive_num          int
			isDeleted_num         int
			deletedAt             sql.NullTime
		)

		var product model.BankVaultProduct
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Currency,
			&rateBps_num,
			&product.Method,
			&product.Frequency,
			&product.LockPeriod,
			&minAmount_num,
			&maxAmount_num,
			&earlyUnlockFeeBps_num,
			&isActive_num,
			&isDeleted_num,
			&product.CreatedAt,
			&product.UpdatedAt,
			&deletedAt,
			&total,
		)

		if err != nil {
			return nil, err
		}

		if rateBps, err := decimal.NewFromString(rateBps_num); err == nil {
			product.RateBps = rateBps
		}
		if minAmount, err := decimal.NewFromString(minAmount_num); err == nil {
			product.MinAmount = minAmount
		}
		if maxAmount, err := decimal.NewFromString(maxAmount_num); err == nil {
			product.MaxAmount = maxAmount
		}
		if earlyUnlockFeeBps, err := decimal.NewFromString(earlyUnlockFeeBps_num); err == nil {
			product.EarlyUnlockFeeBps = earlyUnlockFeeBps
		}

		// Convert Oracle NUMBER(1) to bool
		product.IsActive = isActive_num == 1
		product.IsDeleted = isDeleted_num == 1

		if deletedAt.Valid {
			product.DeletedAt = &deletedAt.Time
		}

		products = append(products, &product)
	}

	resp := types.PaginatedResponse[[]*model.BankVaultProduct]{
		Data: products,
		Meta: types.PaginationMeta{
			TotalDocs:  total,
			Limit:      limit,
			Page:       filterParam.Page,
			TotalPages: 0,
		},
	}
	return &resp, nil
}

func (r *bankVaultRepositary) FindByID(ctx context.Context, id string) (*model.BankVaultProduct, error) {

	query := `
		SELECT
			id, name, description, currency, rate_bps, method, frequency, 
			lock_period, min_amount, max_amount, early_unlock_fee_bps, 
			is_active, is_deleted, created_at, updated_at, deleted_at
		FROM bank_vault_products
		WHERE id = :1
	`

	var (
		rateBps_num           string
		minAmount_num         string
		maxAmount_num         string
		earlyUnlockFeeBps_num string
		isActive_num          int
		isDeleted_num         int
	)

	row := r.db.QueryRowContext(ctx, query, id)
	var product model.BankVaultProduct
	var deletedAt sql.NullTime

	err := row.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Currency,
		&rateBps_num,
		&product.Method,
		&product.Frequency,
		&product.LockPeriod,
		&minAmount_num,
		&maxAmount_num,
		&earlyUnlockFeeBps_num,
		&isActive_num,
		&isDeleted_num,
		&product.CreatedAt,
		&product.UpdatedAt,
		&deletedAt,
	)

	if err != nil {
		return nil, err
	}

	// Convert Oracle numbers to decimal
	if rateBps, err := decimal.NewFromString(rateBps_num); err == nil {
		product.RateBps = rateBps
	}
	if minAmount, err := decimal.NewFromString(minAmount_num); err == nil {
		product.MinAmount = minAmount
	}
	if maxAmount, err := decimal.NewFromString(maxAmount_num); err == nil {
		product.MaxAmount = maxAmount
	}
	if earlyUnlockFeeBps, err := decimal.NewFromString(earlyUnlockFeeBps_num); err == nil {
		product.EarlyUnlockFeeBps = earlyUnlockFeeBps
	}

	// Convert Oracle NUMBER(1) to bool
	product.IsActive = isActive_num == 1
	product.IsDeleted = isDeleted_num == 1

	if deletedAt.Valid {
		product.DeletedAt = &deletedAt.Time
	}

	return &product, nil
}

func (r *bankVaultRepositary) Update(ctx context.Context, id string, product *model.BankVaultProduct) error {
	// First check if the product exists and is not deleted
	existingProduct, err := r.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if existingProduct.IsDeleted {
		return fmt.Errorf("CANNOT_UPDATE_DELETED_BANK_PRODUCT %s", id)
	}

	params := sqlc.UpdateBankVaultParams{ID: id}

	if product.Description != "" {
		params.Description = &product.Description
	}
	if f, ok := product.MinAmount.Float64(); ok {
		params.MinAmount = &sql.NullFloat64{Float64: f, Valid: true}
	}
	if f, ok := product.MaxAmount.Float64(); ok {
		params.MaxAmount = &sql.NullFloat64{Float64: f, Valid: true}
	}

	_, err = r.queries.UpdateBankVault(ctx, params)
	return err
}

func (r *bankVaultRepositary) Delete(ctx context.Context, id string) (string, error) {

	product, err := r.FindByID(ctx, id)
	if err != nil {
		return "", err
	}

	if product.IsDeleted {
		return "", fmt.Errorf("BANK_PRODUCT_ALREADY_DELETED %s", id)
	}

	if product.IsActive {
		bankvaultProduct, err := r.queries.FindBankVaultAndLocks(ctx, id)
		if err != nil {
			return "", err
		}

		for _, lock := range bankvaultProduct.Locks {
			if lock.Status == constants.VaultStatusActive {
				return "", fmt.Errorf("CANNOT_DELETE_BANK_PRODUCT %s", id)
			}
		}
	}

	// Use custom delete query to ensure we only delete non-deleted products
	query := `
		UPDATE bank_vault_products
		SET is_deleted = 1, deleted_at = SYSTIMESTAMP, updated_at = SYSTIMESTAMP
		WHERE id = :1 AND is_deleted = 0
		RETURNING id INTO :2
	`

	var result string
	_, err = r.db.ExecContext(ctx, query, id, sql.Out{Dest: &result})
	if err != nil {
		return "", err
	}

	if result == "" {
		return "", fmt.Errorf("BANK_PRODUCT_ALREADY_DELETED %s", id)
	}

	return result, nil
}

func (r *bankVaultRepositary) EnableOrDisable(ctx context.Context, id string, enable bool) error {

	product, err := r.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if product.IsDeleted {
		return fmt.Errorf("CANNOT_ENABLE_DISABLE_DELETED_BANK_PRODUCT %s", id)
	}

	if enable {
		_, err = r.queries.ActivateBankVault(ctx, id)
		return err
	}
	_, err = r.queries.DeactivateBankVault(ctx, id)
	return err
}
