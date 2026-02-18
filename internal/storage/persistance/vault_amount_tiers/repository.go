package vaultamounttiers

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/persistance/vault_amount_tiers/gen/sqlc"
	"context"
	"database/sql"
	"errors"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type VaultAmountTierRepository struct {
	db      *sql.DB
	queries *sqlc.Queries
	logger  shared_utils.Logger
}

func NewVaultAmountTierRepository(db *sql.DB, logger shared_utils.Logger) storage.VaultAmountTierRepository {
	return &VaultAmountTierRepository{
		db:      db,
		queries: sqlc.New(db),
		logger:  logger,
	}
}

func (r *VaultAmountTierRepository) Create(ctx context.Context, vaultAmountTier *model.VaultAmountTier) (string, error) {
	arg := sqlc.VaultAmountTier{
		VaultCategoryID: vaultAmountTier.VaultCategoryID,
		MinAmount:       vaultAmountTier.MinAmount,
		MaxAmount:       vaultAmountTier.MaxAmount,
		Interest:        vaultAmountTier.Interest,
		IsActive:        sql.NullBool{Bool: vaultAmountTier.IsActive, Valid: true},
	}
	return r.queries.SaveAmountTier(ctx, arg)
}

func (r *VaultAmountTierRepository) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.VaultAmountTier], error) {
	params := sqlc.VaultAmountTierParams{}

	if v, ok := filterParam.Filters["vault_category_id"].(string); ok {
		params.VaultCategoryID = v
	}
	if v, ok := filterParam.Filters["is_active"].(bool); ok {
		params.IsActive = sql.NullBool{Bool: v, Valid: true}
	}
	if filterParam.Page > 0 {
		params.Page = sql.NullInt64{Int64: int64(filterParam.Page), Valid: true}
	}
	if filterParam.PerPage > 0 {
		params.Limit = sql.NullInt64{Int64: int64(filterParam.PerPage), Valid: true}
	}

	tiers, err := r.queries.FindAmountTier(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(localization.ErrorAmountTierNotFound.Code)
		}
		return nil, errors.New(localization.ErrorAmountTierNotFound.Code)
	}

	var result []*model.VaultAmountTier
	for _, t := range tiers {
		m := &model.VaultAmountTier{
			ID:              t.ID,
			VaultCategoryID: t.VaultCategoryID,
			MinAmount:       t.MinAmount,
			MaxAmount:       t.MaxAmount,
			Interest:        t.Interest,
			IsActive:        t.IsActive.Bool,
			CreatedAt:       t.CreatedAt.GoString(),
			UpdatedAt:       t.UpdatedAt.GoString(),
		}
		result = append(result, m)
	}

	return &types.PaginatedResponse[[]*model.VaultAmountTier]{
		Data: result,
		Meta: types.PaginationMeta{
			TotalDocs: int64(len(result)),
			Page:      filterParam.Page,
			Limit:     filterParam.PerPage,
		},
	}, nil
}

func (r *VaultAmountTierRepository) FindByID(ctx context.Context, id string) (*model.VaultAmountTier, error) {
	t, err := r.queries.FindAmountTierById(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(localization.ErrorAmountTierNotFound.Code)
		}
		return nil, errors.New(localization.ErrorAmountTierNotFound.Code)
	}
	m := &model.VaultAmountTier{
		ID:              t.ID,
		VaultCategoryID: t.VaultCategoryID,
		MinAmount:       t.MinAmount,
		MaxAmount:       t.MaxAmount,
		Interest:        t.Interest,
		IsActive:        t.IsActive.Bool,
		CreatedAt:       t.CreatedAt.GoString(),
		UpdatedAt:       t.UpdatedAt.GoString(),
	}
	return m, nil
}

func (r *VaultAmountTierRepository) Update(ctx context.Context, id string, vaultAmountTier *model.VaultAmountTier) error {
	arg := sqlc.VaultAmountTier{
		ID:        id,
		MinAmount: vaultAmountTier.MinAmount,
		MaxAmount: vaultAmountTier.MaxAmount,
		Interest:  vaultAmountTier.Interest,
	}
	_, err := r.queries.UpdateAmountTier(ctx, arg)
	if err != nil {
		r.logger.Errorf("[Update] failed to update amount tier: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (r *VaultAmountTierRepository) Delete(ctx context.Context, id string) (string, error) {
	return r.queries.DeleteAmountTier(ctx, id)
}

func (r *VaultAmountTierRepository) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	_, err := r.queries.EnableOrDisableAmountTier(ctx, id, enable)
	if err != nil {
		r.logger.Errorf("[EnableOrDisable] failed to enable/disable amount tier: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}
