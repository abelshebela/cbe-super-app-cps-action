package bankvault

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/persistance/bankvault/gen/sqlc"
	"cbe-super-app-cps-action/pkgs/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"

	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type bankVaultRepositary struct {
	db      *sql.DB
	queries *sqlc.Queries
	logger  shared_utils.Logger
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
		return errors.New(localization.ErrorTransactionFailed.Code)
	}
	defer tx.Rollback()

	txRepo := &bankVaultRepositary{
		db:      r.db,
		queries: r.queries.WithTx(tx),
		logger:  r.logger,
	}

	if err := fn(txRepo); err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return tx.Commit()
}

func (r *bankVaultRepositary) Create(ctx context.Context, product *model.BankVaultProduct) (string, error) {
	// Convert time.Duration to nanoseconds
	lockPeriodNanoseconds := product.LockPeriod.Nanoseconds()

	return r.createBankVaultCustom(ctx, product, lockPeriodNanoseconds)
}

func (r *bankVaultRepositary) createBankVaultCustom(ctx context.Context, product *model.BankVaultProduct, lockPeriodNanoseconds int64) (string, error) {
	params := sqlc.SaveBankVaultParams{
		Name:        product.Name,
		Description: product.Description,
		Currency:    product.Currency,
		RateBps:     product.RateBps,
		Method:      product.Method,
		Frequency:   product.Frequency,
		LockPeriod:  product.LockPeriod,
		MinAmount:   product.MinAmount,
		MaxAmount:   product.MaxAmount,
		EarlyUnlockRateBps: sql.NullBool{
			Bool:  product.EarlyUnlockRateBps,
			Valid: true,
		},
		IsActive: sql.NullBool{
			Bool:  product.IsActive,
			Valid: true,
		},
	}

	_, err := r.queries.FindBankVaultByName(ctx, params.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return r.queries.SaveBankVault(ctx, params)
		}
		return "", errors.New(localization.ErrorUnexpectedError.Code)
	} else {
		return "", errors.New(localization.ErrorDuplicateBankProduct.Code)
	}
}

func (r *bankVaultRepositary) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.BankVaultProduct], error) {
	params := sqlc.FindBankVaultParams{}

	if v, ok := filterParam.Filters["is_active"].(bool); ok {
		params.IsActive = sql.NullBool{Bool: v, Valid: true}
	}
	if v, ok := filterParam.Filters["is_deleted"].(bool); ok {
		params.IsDeleted = sql.NullBool{Bool: v, Valid: true}
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
	} else if v, ok := filterParam.Filters["frequency"]; ok {
		if freqStr := fmt.Sprintf("%v", v); freqStr != "" {
			params.Frequency = sqlc.NullAccrualFrequency{AccrualFrequency: sqlc.AccrualFrequency(freqStr), Valid: true}
		}
	}
	if filterParam.Page > 0 {
		params.Page = sql.NullInt64{Int64: int64(filterParam.Page), Valid: true}
	}
	if filterParam.PerPage > 0 {
		params.Limit = sql.NullInt64{Int64: int64(filterParam.PerPage), Valid: true}
	}

	rows, err := r.queries.FindBankVault(ctx, params)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	products := make([]*model.BankVaultProduct, 0, len(rows))
	var total int64
	for _, row := range rows {
		p := &model.BankVaultProduct{
			ID:                 row.ID,
			Name:               row.Name,
			Description:        row.Description,
			Currency:           row.Currency,
			RateBps:            row.RateBps,
			Method:             constants.AccrualMethod(row.Method),
			Frequency:          constants.AccrualFrequency(row.Frequency),
			LockPeriod:         row.LockPeriod,
			MinAmount:          row.MinAmount,
			MaxAmount:          row.MaxAmount,
			EarlyUnlockRateBps: utils.NullBoolToBool(row.EarlyUnlockRateBps),
			IsActive:           utils.NullBoolToBool(row.IsActive),
			IsDeleted:          utils.NullBoolToBool(row.IsDeleted),
			CreatedAt:          row.CreatedAt,
			UpdatedAt:          row.UpdatedAt,
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
			Page:       filterParam.Page,
			TotalPages: 0,
		},
	}
	return &resp, nil
}

func (r *bankVaultRepositary) FindByID(ctx context.Context, id string) (*model.BankVaultProduct, error) {
	row, err := r.queries.FindBankVaultById(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(localization.ErrorNoBankProductFound.Code)
		}
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	product := &model.BankVaultProduct{
		ID:                 row.ID,
		Name:               row.Name,
		Description:        row.Description,
		Currency:           row.Currency,
		RateBps:            row.RateBps,
		Method:             constants.AccrualMethod(row.Method),
		Frequency:          constants.AccrualFrequency(row.Frequency),
		LockPeriod:         row.LockPeriod,
		MinAmount:          row.MinAmount,
		MaxAmount:          row.MaxAmount,
		EarlyUnlockRateBps: utils.NullBoolToBool(row.EarlyUnlockRateBps),
		IsActive:           utils.NullBoolToBool(row.IsActive),
		IsDeleted:          utils.NullBoolToBool(row.IsDeleted),
		CreatedAt:          row.CreatedAt,
		UpdatedAt:          row.UpdatedAt,
	}
	if row.DeletedAt.Valid {
		t := row.DeletedAt.Time
		product.DeletedAt = &t
	}

	return product, nil

}

func (r *bankVaultRepositary) FindBankVaultByName(ctx context.Context, name string) error {
	_, err := r.queries.FindBankVaultByName(ctx, name)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (r *bankVaultRepositary) Update(ctx context.Context, id string, product *model.BankVaultProduct) error {
	// First check if the product exists and is not deleted
	existingProduct, err := r.FindByID(ctx, id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
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
	bankvaultProduct, err := r.queries.FindBankVaultAndLocks(ctx, id)
	if err != nil && err != sql.ErrNoRows {
		return "", errors.New(localization.ErrorUnexpectedError.Code)
	}

	if bankvaultProduct.IsActive.Bool {
		for _, lock := range bankvaultProduct.Locks {
			if lock.Status == constants.VaultStatusActive {
				return "", errors.New(localization.ErrorCannotDeletedBankProduct.Code)
			}
		}
	}

	_, err = r.queries.DeleteBankVault(ctx, id)
	if err != nil {
		return "", errors.New(localization.ErrorUnexpectedError.Code)
	}

	return "", nil
}

func (r *bankVaultRepositary) EnableOrDisable(ctx context.Context, id string, enable bool) error {

	product, err := r.FindByID(ctx, id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	if product.IsDeleted {
		return errors.New(localization.ErrorCannotEnableOrDisable.Code)
	}

	if enable {
		_, err = r.queries.ActivateBankVault(ctx, id)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	_, err = r.queries.DeactivateBankVault(ctx, id)
	return errors.New(localization.ErrorUnexpectedError.Code)
}
