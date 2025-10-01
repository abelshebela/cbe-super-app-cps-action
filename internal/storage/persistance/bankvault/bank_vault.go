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

// -----------------------------
// Product methods
// -----------------------------
func (r *bankVaultRepositary) Create(ctx context.Context, product *model.BankVaultProduct) (string, error) {
	// Convert time.Duration to days for Oracle NUMBER(19,0) constraint
	lockPeriodDays := int64(product.LockPeriod.Hours() / 24)

	// Use custom Oracle insert to handle lock_period correctly
	return r.createBankVaultCustom(ctx, product, lockPeriodDays)
}

func (r *bankVaultRepositary) createBankVaultCustom(ctx context.Context, product *model.BankVaultProduct, lockPeriodDays int64) (string, error) {
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
		lockPeriodDays, // Use days directly instead of nanoseconds
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
	params := sqlc.FindBankVaultParams{}

	if v, ok := filterParam.Filters["is_active"].(bool); ok {
		params.IsActive = sql.NullBool{Bool: v, Valid: true}
	}
	if v, ok := filterParam.Filters["name"].(string); ok && v != "" {
		params.NameQuery = sql.NullString{String: v, Valid: true}
	}
	if v, ok := filterParam.Filters["currency"].(string); ok && v != "" {
		params.Currency = sql.NullString{String: v, Valid: true}
	}
	if v, ok := filterParam.Filters["method"].(string); ok && v != "" {
		params.Method = sqlc.NullAccrualMethod{AccrualMethod: sqlc.AccrualMethod(v), Valid: true}
	}
	if v, ok := filterParam.Filters["frequency"].(string); ok && v != "" {
		params.Frequency = sqlc.NullAccrualFrequency{AccrualFrequency: sqlc.AccrualFrequency(v), Valid: true}
	}
	if filterParam.Page > 0 {
		params.Page = sql.NullInt64{Int64: int64(filterParam.Page), Valid: true}
	}
	if filterParam.PerPage > 0 {
		params.Limit = sql.NullInt64{Int64: int64(filterParam.PerPage), Valid: true}
	}

	rows, err := r.queries.FindBankVault(ctx, params)
	if err != nil {
		return nil, err
	}

	products := make([]*model.BankVaultProduct, 0, len(rows))
	var total int64
	for _, row := range rows {
		p := &model.BankVaultProduct{
			ID:                row.ID,
			Name:              row.Name,
			Description:       row.Description,
			Currency:          row.Currency,
			RateBps:           row.RateBps,
			Method:            constants.AccrualMethod(row.Method),
			Frequency:         constants.AccrualFrequency(row.Frequency),
			LockPeriod:        row.LockPeriod,
			MinAmount:         row.MinAmount,
			MaxAmount:         row.MaxAmount,
			EarlyUnlockFeeBps: row.EarlyUnlockFeeBps,
			IsActive:          sqlBoolToBool(row.IsActive),
			CreatedAt:         row.CreatedAt,
			UpdatedAt:         row.UpdatedAt,
		}
		if row.DeletedAt.Valid {
			t := row.DeletedAt.Time
			p.DeletedAt = &t
		}
		products = append(products, p)
		total = row.TotalCount
	}

	resp := types.PaginatedResponse[[]*model.BankVaultProduct]{
		Data: products,
		Meta: types.PaginationMeta{
			TotalDocs:  total,
			Limit:      int(params.Limit.Int64),
			Page:       int(params.Page.Int64),
			TotalPages: 0,
		},
	}
	return &resp, nil
}

func (r *bankVaultRepositary) FindByID(ctx context.Context, id string) (*model.BankVaultProduct, error) {
	row, err := r.queries.FindBankVaultById(ctx, id)
	if err != nil {
		return nil, err
	}

	product := &model.BankVaultProduct{
		ID:                row.ID,
		Name:              row.Name,
		Description:       row.Description,
		Currency:          row.Currency,
		RateBps:           row.RateBps,
		Method:            constants.AccrualMethod(row.Method),
		Frequency:         constants.AccrualFrequency(row.Frequency),
		LockPeriod:        row.LockPeriod,
		MinAmount:         row.MinAmount,
		MaxAmount:         row.MaxAmount,
		EarlyUnlockFeeBps: row.EarlyUnlockFeeBps,
		IsActive:          row.IsActive,
		CreatedAt:         row.CreatedAt,
		UpdatedAt:         row.UpdatedAt,
	}
	if row.DeletedAt.Valid {
		t := row.DeletedAt.Time
		product.DeletedAt = &t
	}

	return product, nil
}

func (r *bankVaultRepositary) Update(ctx context.Context, id string, product *model.BankVaultProduct) error {
	params := sqlc.UpdateBankVaultParams{ID: id}

	// Since product fields are values, set all as provided
	if product.Description != "" {
		params.Description = &product.Description
	}
	if f, ok := product.MinAmount.Float64(); ok {
		params.MinAmount = &sql.NullFloat64{Float64: f, Valid: true}
	}
	if f, ok := product.MaxAmount.Float64(); ok {
		params.MaxAmount = &sql.NullFloat64{Float64: f, Valid: true}
	}

	// IsActive is not updated here handled by EnableOrDisable methods
	_, err := r.queries.UpdateBankVault(ctx, params)
	return err
}

func (r *bankVaultRepositary) Delete(ctx context.Context, id string) (string, error) {
	bankvaultProduct, err := r.queries.FindBankVaultAndLocks(ctx, id)
	if err != nil {
		return "", err
	}

	if bankvaultProduct.IsActive {
		for _, lock := range bankvaultProduct.Locks {
			if lock.Status == constants.VaultStatusActive {
				return "", fmt.Errorf("CANNOT_DELETE_BANK_PRODUCT %s", id)
			}
		}
	}

	return r.queries.DeleteBankVault(ctx, id)
}

func (r *bankVaultRepositary) EnableOrDisable(ctx context.Context, id string, enable bool) {
	if enable {
		_, _ = r.queries.ActivateBankVault(ctx, id)
		return
	}
	_, _ = r.queries.DeactivateBankVault(ctx, id)
}
