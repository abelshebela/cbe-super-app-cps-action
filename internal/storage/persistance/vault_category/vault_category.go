package vaultcategory

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/persistance/vault_category/gen/sqlc"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	imodel "cbe-super-app-cps-action/internal/constants/model"

	// "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type VaultCategoryRepository struct {
	db     *sql.DB
	logger shared_utils.Logger
}

func NewVaultCategoryRepository(db *sql.DB, logger shared_utils.Logger) storage.VaultCategoryRepository {
	return &VaultCategoryRepository{
		db:     db,
		logger: logger,
	}
}

// Create creates a new vault category in Oracle and returns its ID
func (r *VaultCategoryRepository) Create(ctx context.Context, entity *imodel.VaultCategory) (string, error) {
	q := sqlc.New(r.db)
	if _, err := q.FindVaultCategoryByName(ctx, entity.Name); err == nil {
		return "", errors.New(localization.ErrorDuplicateGroupVaultCategory.Code)
	} else if !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	qtx := q.WithTx(tx)

	categoryID, err := qtx.SaveVaultCategory(ctx, entity)
	if err != nil {
		return "", err
	}

	err = qtx.SaveVaultTiers(ctx, categoryID, entity)
	if err != nil {
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	return strings.ToUpper(categoryID), nil
}

// FindAllWithPagination lists categories with filters and pagination, including deleted ones
func (r *VaultCategoryRepository) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*imodel.VaultCategory], error) {
	q := sqlc.New(r.db)
	params := sqlc.FindVaultCategoryParams{}
	if v, ok := filterParam.Filters["is_active"].(bool); ok {
		params.IsActive = sql.NullBool{Bool: v, Valid: true}
	}
	if v, ok := filterParam.Filters["name"].(string); ok && v != "" {
		params.NameQuery = sql.NullString{String: v, Valid: true}
	}
	if filterParam.Page > 0 && filterParam.PerPage > 0 {
		params.Page = sql.NullInt64{Int64: int64(filterParam.Page), Valid: true}
		params.Limit = sql.NullInt64{Int64: int64(filterParam.PerPage), Valid: true}
	} else if filterParam.PerPage > 0 {
		params.Page = sql.NullInt64{Int64: 1, Valid: true}
		params.Limit = sql.NullInt64{Int64: int64(filterParam.PerPage), Valid: true}
	}

	rows, err := q.FindVaultCategory(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(localization.ErrorVaultGroupCategoryNotFound.Code)
		}
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	list := make([]*model.VaultCategory, 0, len(rows))
	var total int64
	for _, rrow := range rows {
		e := &model.VaultCategory{
			ID:            rrow.ID,
			Name:          rrow.Name,
			CoverImageURL: rrow.CoverImage,
			IsActive:      rrow.IsActive,
			CreatedAt:     rrow.CreatedAt,
			UpdatedAt:     rrow.UpdatedAt,
		}
		total = rrow.TotalCount
		list = append(list, e)
	}

	limit := int(params.Limit.Int64)
	if limit == 0 {
		limit = 50
	}
	page := int(params.Page.Int64)
	if page == 0 {
		page = 1
	}

	resp := types.PaginatedResponse[[]*model.VaultCategory]{
		Data: list,
		Meta: types.PaginationMeta{
			TotalDocs:  total,
			Limit:      limit,
			Page:       page,
			TotalPages: 0,
		},
	}
	return &resp, nil
}

// FindByID fetches a single category by id, including deleted ones
func (r *VaultCategoryRepository) FindByID(ctx context.Context, id string) (*imodel.VaultCategory, error) {
	q := sqlc.New(r.db)
	rrow, err := q.FindVaultCategoryById(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(localization.ErrorVaultGroupCategoryNotFound.Code)
		}
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	e := &model.VaultCategory{
		ID:            rrow.ID,
		Name:          rrow.Name,
		CoverImageURL: rrow.CoverImage,
		IsActive:      rrow.IsActive,
		CreatedAt:     rrow.CreatedAt,
		UpdatedAt:     rrow.UpdatedAt,
	}
	return e, nil
}

func (r *VaultCategoryRepository) FindByName(ctx context.Context, name string) (*imodel.VaultCategory, error) {
	q := sqlc.New(r.db)
	category, err := q.FindVaultCategoryByName(ctx, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	vc := &model.VaultCategory{
		ID:            category.ID,
		Name:          category.Name,
		CoverImageURL: category.CoverImage,
		IsActive:      category.IsActive,
		CreatedAt:     category.CreatedAt,
		UpdatedAt:     category.UpdatedAt,
	}

	return vc, nil
}

// Update updates name/description; prevents updates on deleted records
func (r *VaultCategoryRepository) Update(ctx context.Context, id string, entity *imodel.VaultCategory) error {
	// Ensure not deleted
	_, err := r.FindByID(ctx, id)
	if err != nil {
		return err
	}

	params := sqlc.UpdateVaultCategoryParams{
		Name:       sql.NullString{String: entity.Name, Valid: entity.Name != ""},
		CoverImage: sql.NullString{String: entity.CoverImageURL, Valid: entity.CoverImageURL != ""},
		ID:         id,
	}
	_, err = sqlc.New(r.db).UpdateVaultCategory(ctx, params)
	return err
}

// Delete performs a soft delete; prevents double delete
func (r *VaultCategoryRepository) Delete(ctx context.Context, id string) (string, error) {
	_, err := r.FindByID(ctx, id)
	if err != nil {
		return "", err
	}
	q := sqlc.New(r.db)
	ret, err := q.DeleteVaultCategory(ctx, id)
	if err != nil {
		return "", err
	}
	if ret == "" {
		return "", fmt.Errorf("VAULT_CATEGORY_ALREADY_DELETED %s", id)
	}
	return strings.ToUpper(ret), nil
}

// EnableOrDisable toggles active state; prevents toggling deleted records
func (r *VaultCategoryRepository) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	_, err := r.FindByID(ctx, id)
	if err != nil {
		return err
	}
	q := sqlc.New(r.db)
	if enable {
		_, err = q.ActivateVaultCategory(ctx, id)
		return err
	}
	_, err = q.DeactivateVaultCategory(ctx, id)
	return err
}
