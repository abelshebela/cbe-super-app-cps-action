package bankvault

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/persistance/bankvault/gen/sqlc"
	"cbe-super-app-cps-action/pkgs/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
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
		return err
	}

	if err := tx.Commit(); err != nil {
		r.logger.Errorf("[ExecuteInTransaction] failed to commit transaction: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (r *bankVaultRepositary) Create(ctx context.Context, product *model.BankVaultProduct) (string, error) {
	params := sqlc.SaveBankVaultParams{
		Name:       product.Name,
		Currency:   product.Currency,
		Interest:   product.Interest,
		Method:     string(product.Method),
		Frequency:  product.Frequency,
		LockPeriod: product.LockPeriod,
		MinAmount:  product.MinAmount,
		MaxAmount:  product.MaxAmount,
		ApplyInterestOnEarlyUnlock: sql.NullBool{
			Bool:  product.ApplyInterestOnEarlyUnlock,
			Valid: true,
		},
		IsActive: sql.NullBool{
			Bool:  true,
			Valid: true,
		},
	}

	id, err := r.queries.SaveBankVault(ctx, params)
	if err != nil {
		r.logger.Errorf("[Create] failed to save bank vault: %v", err)
		return "", errors.New(localization.ErrorUnexpectedError.Code)
	}
	return id, nil
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
	if filterParam.Page > 0 {
		params.Page = sql.NullInt64{Int64: int64(filterParam.Page), Valid: true}
	}
	if filterParam.PerPage > 0 {
		params.Limit = sql.NullInt64{Int64: int64(filterParam.PerPage), Valid: true}
	}

	rows, err := r.queries.FindBankVault(ctx, params)
	if err != nil {
		return nil, errors.New(localization.ErrorResourceNotFound.Code)
	}

	products := make([]*model.BankVaultProduct, 0, len(rows))
	var total int64
	for _, row := range rows {
		p := &model.BankVaultProduct{
			ID:                         row.ID,
			Name:                       row.Name,
			Currency:                   row.Currency,
			Interest:                   row.Interest,
			Method:                     string(row.Method),
			Frequency:                  row.Frequency,
			LockPeriod:                 row.LockPeriod,
			MinAmount:                  row.MinAmount,
			MaxAmount:                  row.MaxAmount,
			ApplyInterestOnEarlyUnlock: utils.NullBoolToBool(row.ApplyInterestOnEarlyUnlock),
			IsActive:                   utils.NullBoolToBool(row.IsActive),
			IsDeleted:                  utils.NullBoolToBool(row.IsDeleted),
			CreatedAt:                  row.CreatedAt,
			UpdatedAt:                  row.UpdatedAt,
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

func (r *bankVaultRepositary) FindAllBankLockedVaultsWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.LockedVault], error) {
	params := sqlc.ListLockedVaultParams{}

	rows, err := r.queries.GetAllLockedVaults(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {

			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		r.logger.Errorf("failed to get locked vaults: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	lockedVaults := make([]*model.LockedVault, 0, len(rows))
	var total int64

	for _, row := range rows {
		lv := &model.LockedVault{
			ID:                         row.ID,
			CustomerID:                 row.CustomerID,
			LinkedAccount:              row.LinkedAccount,
			AccountHolderName:          row.AccountHolderName,
			TransactionReference:       row.TransactionReference,
			ProductID:                  row.ProductID,
			Principal:                  row.Principal,
			StartDate:                  row.StartDate,
			MaturityDate:               row.MaturityDate,
			Status:                     row.Status,
			TermsVersion:               row.TermsVersion,
			TermsAcceptedAt:            row.TermsAcceptedAt,
			MinAmount:                  row.MinAmount,
			MaxAmount:                  row.MaxAmount,
			Interest:                   row.Interest.Div(decimal.NewFromInt(100)),
			Method:                     row.Method,
			Frequency:                  row.Frequency,
			ApplyInterestOnEarlyUnlock: nil,
			LockPeriod:                 fmt.Sprintf("%d months", utils.DurationToMonths(row.LockPeriod)),
			CreatedAt:                  row.CreatedAt,
			UpdatedAt:                  row.UpdatedAt,
		}

		if row.ApplyInterestOnEarlyUnlock.Valid {
			val := row.ApplyInterestOnEarlyUnlock.Bool
			lv.ApplyInterestOnEarlyUnlock = &val
		}

		if row.ClosedAt.Valid {
			t := row.ClosedAt.Time
			lv.ClosedAt = &t
		}
		if row.DeletedAt.Valid {
			t := row.DeletedAt.Time
			lv.DeletedAt = &t
		}

		lockedVaults = append(lockedVaults, lv)
		total = row.TotalCount
	}

	resp := types.PaginatedResponse[[]*model.LockedVault]{
		Data: lockedVaults,
		Meta: types.PaginationMeta{
			TotalDocs:  total,
			Limit:      int(params.Limit.Int64),
			Page:       filterParam.Page,
			TotalPages: int(total),
		},
	}

	return &resp, nil
}

func (r *bankVaultRepositary) FindAllGroupVaultWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.GroupVault], error) {
	params := sqlc.GetAllGroupVaultsParams{}

	rows, err := r.queries.GetAllGroupVaults(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(localization.ErrorVaultCategoryNotFound.Code)
		}
		r.logger.Errorf("failed to get group vaults: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	groupVaults := make([]*model.GroupVault, 0, len(rows))
	var total int64

	for _, row := range rows {
		gv := &model.GroupVault{
			ID:            row.ID,
			VaultName:     row.VaultName,
			VaultCategory: row.VaultCategory,
			TargetAmount:  row.TargetAmount,
			// IsDeadlock: row.IsDeadlock,
			VaultType: row.VaultType,
			// InterestRate: row.InterestRate,
			Status: row.Status,
			// CreatedAt:  row.CreatedAt,
			// UpdatedAt:  row.UpdatedAt,
			// DeletedAt:  row.DeletedAt,
		}

		if row.CreatedAt.Valid {
			t, err := time.Parse(time.RFC3339, row.CreatedAt.String)
			if err == nil {
				gv.CreatedAt = t
			} else {
				// Try another common format if RFC3339 fails, or log it
				// For now assuming RFC3339 based on codebase convention, but catching error
				r.logger.Errorf("failed to parse CreatedAt: %v", err)
			}
		}

		if row.UpdatedAt.Valid {
			t, err := time.Parse(time.RFC3339, row.UpdatedAt.String)
			if err == nil {
				gv.UpdatedAt = t
			} else {
				r.logger.Errorf("failed to parse UpdatedAt: %v", err)
			}
		}

		if row.DeletedAt.Valid {
			t, err := time.Parse(time.RFC3339, row.DeletedAt.String)
			if err == nil {
				gv.DeletedAt = &t
			} else {
				r.logger.Errorf("failed to parse DeletedAt: %v", err)
			}
		}

		groupVaults = append(groupVaults, gv)
		total = row.TotalCount
	}

	resp := &types.PaginatedResponse[[]*model.GroupVault]{
		Data: groupVaults,
		Meta: types.PaginationMeta{
			TotalDocs:  total,
			Limit:      int(params.Limit.Int64),
			Page:       filterParam.Page,
			TotalPages: int(total),
		},
	}

	return resp, nil
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
		ID:                         row.ID,
		Name:                       row.Name,
		Currency:                   row.Currency,
		Interest:                   row.Interest,
		Method:                     string(row.Method),
		Frequency:                  row.Frequency,
		LockPeriod:                 row.LockPeriod,
		MinAmount:                  row.MinAmount,
		MaxAmount:                  row.MaxAmount,
		ApplyInterestOnEarlyUnlock: utils.NullBoolToBool(row.ApplyInterestOnEarlyUnlock),
		IsActive:                   utils.NullBoolToBool(row.IsActive),
		IsDeleted:                  utils.NullBoolToBool(row.IsDeleted),
		CreatedAt:                  row.CreatedAt,
		UpdatedAt:                  row.UpdatedAt,
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
		if err == sql.ErrNoRows {
			return errors.New(localization.ErrorNoBankProductFound.Code)
		}
		r.logger.Errorf("[FindBankVaultByName] failed to find bank vault by name: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (r *bankVaultRepositary) Update(ctx context.Context, id string, product *model.BankVaultProduct) error {
	existingProduct, err := r.FindByID(ctx, id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	if existingProduct.IsDeleted {
		return errors.New(localization.ErrorCannotDeletedBankProduct.Code)
	}

	params := sqlc.UpdateBankVaultParams{ID: id}

	if f, ok := product.MinAmount.Float64(); ok {
		params.MinAmount = &sql.NullFloat64{Float64: f, Valid: true}
	}
	if f, ok := product.MaxAmount.Float64(); ok {
		params.MaxAmount = &sql.NullFloat64{Float64: f, Valid: true}
	}

	_, err = r.queries.UpdateBankVault(ctx, params)
	if err != nil {
		r.logger.Errorf("[Update] failed to update bank vault: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
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
		if err != nil {
			return errors.New(localization.ErrorUnexpectedError.Code)
		}
		return nil
	}

	_, err = r.queries.DeactivateBankVault(ctx, id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}
