package vault

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/persistance/vault/gen/sqlc"
	"context"
	"database/sql"
	"errors"
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
		return "", errors.New(localization.ErrorDuplicateVaultCategory.Code)
	} else if !errors.Is(err, sql.ErrNoRows) {
		r.logger.Errorf("[Create] failed to check vault category existence: %v", err)
		return "", errors.New(localization.ErrorUnexpectedError.Code)
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Errorf("[Create] failed to begin transaction: %v", err)
		return "", errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer tx.Rollback()

	qtx := q.WithTx(tx)

	categoryID, err := qtx.SaveVaultCategory(ctx, entity)
	if err != nil {
		r.logger.Errorf("[Create] failed to save vault category: %v", err)
		return "", errors.New(localization.ErrorUnexpectedError.Code)
	}

	err = qtx.SaveVaultTiers(ctx, categoryID, entity)
	if err != nil {
		r.logger.Errorf("[Create] failed to save vault tiers: %v", err)
		return "", errors.New(localization.ErrorUnexpectedError.Code)
	}

	if err := tx.Commit(); err != nil {
		r.logger.Errorf("[Create] failed to commit transaction: %v", err)
		return "", errors.New(localization.ErrorUnexpectedError.Code)
	}

	return strings.ToUpper(categoryID), nil
}

func (r *VaultCategoryRepository) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*imodel.VaultCategory], error) {
	q := sqlc.New(r.db)

	params := sqlc.FindVaultCategoryParams{}

	if filterParam.Filters != nil {
		if v, ok := filterParam.Filters["is_active"].(bool); ok {
			params.IsActive = sql.NullBool{Bool: v, Valid: true}
		}
		if v, ok := filterParam.Filters["name"].(string); ok && v != "" {
			params.NameQuery = sql.NullString{String: v, Valid: true}
		}
	}

	if filterParam.PerPage > 0 {
		params.Limit = sql.NullInt64{Int64: int64(filterParam.PerPage), Valid: true}
	} else {
		params.Limit = sql.NullInt64{Int64: 50, Valid: true}
	}

	if filterParam.Page > 0 {
		params.Page = sql.NullInt64{Int64: int64(filterParam.Page), Valid: true}
	} else {
		params.Page = sql.NullInt64{Int64: 1, Valid: true}
	}

	rows, err := q.FindVaultCategories(ctx, params)
	if err != nil {
		r.logger.Errorf("failed to find vault categories: %v", err)
		if err == sql.ErrNoRows {
			return nil, errors.New(localization.ErrorVaultCategoryNotFound.Code)
		}
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	categories := make([]*imodel.VaultCategory, 0, len(rows))
	var total int64

	for _, row := range rows {
		c := &imodel.VaultCategory{
			ID:               row.ID,
			Name:             row.Name,
			CoverImageURL:    row.CoverImageURL,
			InterestType:     row.InterestType,
			CategoryInterest: row.CategoryInterest,
			Deadlock:         row.Deadlock,
			Tiers:            row.Tiers,
			IsActive:         row.IsActive,
			CreatedAt:        row.CreatedAt,
			UpdatedAt:        row.UpdatedAt,
		}

		total = row.TotalCount
		categories = append(categories, c)
	}

	limit := int(params.Limit.Int64)
	if limit <= 0 {
		limit = 50
	}
	page := int(params.Page.Int64)
	if page <= 0 {
		page = 1
	}

	totalPages := 0
	if limit > 0 && total > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}

	pagingCounter := 0
	if total > 0 {
		pagingCounter = (page-1)*limit + 1
	}

	hasPrevPage := page > 1
	hasNextPage := int64(page*limit) < total

	var prevPage *int
	var nextPage *int
	if hasPrevPage {
		p := page - 1
		prevPage = &p
	}
	if hasNextPage {
		n := page + 1
		nextPage = &n
	}

	resp := types.PaginatedResponse[[]*imodel.VaultCategory]{
		Data: categories,
		Meta: types.PaginationMeta{
			TotalDocs:     total,
			Limit:         limit,
			TotalPages:    totalPages,
			Page:          page,
			PagingCounter: pagingCounter,
			HasPrevPage:   hasPrevPage,
			HasNextPage:   hasNextPage,
			PrevPage:      prevPage,
			NextPage:      nextPage,
		},
	}

	return &resp, nil
}

// FindByID fetches a single category by id, including deleted ones
func (r *VaultCategoryRepository) FindByID(ctx context.Context, id string) (*imodel.VaultCategory, error) {
	q := sqlc.New(r.db)
	category, tiers, err := q.FindVaultCategoryWithTiers(ctx, id)
	if err != nil {
		r.logger.Errorf("failed to find vault category with it's tier: %v", err)
		if err == sql.ErrNoRows {
			return nil, errors.New(localization.ErrorVaultCategoryNotFound.Code)
		}
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	e := &imodel.VaultCategory{
		ID:               category.ID,
		Name:             category.Name,
		CoverImageURL:    category.CoverImageURL,
		InterestType:     category.InterestType,
		CategoryInterest: category.CategoryInterest,
		Deadlock:         category.Deadlock,
		Tiers:            tiers,
		IsActive:         category.IsActive,
		CreatedAt:        category.CreatedAt,
		UpdatedAt:        category.UpdatedAt,
	}
	return e, nil
}

func (r *VaultCategoryRepository) FindByName(ctx context.Context, name string) (*imodel.VaultCategory, error) {
	q := sqlc.New(r.db)
	category, err := q.FindVaultCategoryByName(ctx, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(localization.ErrorVaultCategoryNotFound.Code)
		}
		r.logger.Errorf("[FindByName] failed to find vault category by name: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	vc := &model.VaultCategory{
		ID:            category.ID,
		Name:          category.Name,
		CoverImageURL: category.CoverImageURL,
		IsActive:      category.IsActive,
		CreatedAt:     category.CreatedAt,
		UpdatedAt:     category.UpdatedAt,
	}

	return vc, nil
}

func (r *VaultCategoryRepository) Update(ctx context.Context, id string, entity *imodel.VaultCategory) error {
	q := sqlc.New(r.db)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Errorf("[Update] failed to begin transaction: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer tx.Rollback()

	qtx := q.WithTx(tx)

	_, err = qtx.UpdateVaultCategory(ctx, id, entity)
	if err != nil {
		r.logger.Errorf("[Update] failed to update vault category on transaction: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	err = qtx.UpdateVaultTiers(ctx, id, entity)
	if err != nil {
		r.logger.Errorf("[Update] failed to update vault tiers on transaction: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	if err := tx.Commit(); err != nil {
		r.logger.Errorf("[Update] failed to commit transaction: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
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
		r.logger.Errorf("[Delete] failed to delete vault category: %v", err)
		return "", errors.New(localization.ErrorUnexpectedError.Code)
	}
	if ret == "" {
		return "", errors.New(localization.ErrorVaultCategoryNotFound.Code)
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
		if err != nil {
			r.logger.Errorf("[EnableOrDisable] failed to activate vault category: %v", err)
			return errors.New(localization.ErrorUnexpectedError.Code)
		}
		return nil
	}
	_, err = q.DeactivateVaultCategory(ctx, id)
	if err != nil {
		r.logger.Errorf("[EnableOrDisable] failed to deactivate vault category: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}
