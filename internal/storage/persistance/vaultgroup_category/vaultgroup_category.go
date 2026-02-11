package vaultgroupcategory

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/persistance/vaultgroup_category/gen/sqlc"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type VaultGroupCategoryRepository struct {
	db     *sql.DB
	logger shared_utils.Logger
}

func NewVaultGroupCategoryRepository(db *sql.DB, logger shared_utils.Logger) storage.VaultGroupCategoryRepository {
	return &VaultGroupCategoryRepository{
		db:     db,
		logger: logger,
	}
}

// Create creates a new vault group category in Oracle and returns its ID
func (r *VaultGroupCategoryRepository) Create(ctx context.Context, entity *model.VaultCategory) (string, error) {
	q := sqlc.New(r.db)
	if _, err := q.FindVaultGroupCategoryByName(ctx, entity.Name); err == nil {
		return "", errors.New(localization.ErrorDuplicateGroupVaultCategory.Code)
	} else if !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}

	params := sqlc.SaveVaultGroupCategoryParams{
		Name:         strings.ToUpper(entity.Name),
		CategoryType: entity.CategoryType,
		CoverImage:   entity.CoverImage,
		IsActive:     sql.NullBool{Bool: entity.IsActive, Valid: true},
	}
	id, err := q.SaveVaultGroupCategory(ctx, params)
	if err != nil {
		return "", err
	}
	return strings.ToUpper(id), nil
}

// FindAllWithPagination lists categories with filters and pagination, including deleted ones
func (r *VaultGroupCategoryRepository) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.VaultCategory], error) {
	q := sqlc.New(r.db)
	params := sqlc.FindVaultGroupCategoryParams{}
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

	rows, err := q.FindVaultGroupCategory(ctx, params)
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
			ID:         rrow.ID,
			Name:       rrow.Name,
			CoverImage: rrow.CoverImage,
			IsActive:   rrow.IsActive,
			IsDeleted:  rrow.IsDeleted,
		}
		if rrow.CreatedAt.Valid {
			e.CreatedAt = rrow.CreatedAt.Time
		}
		if rrow.UpdatedAt.Valid {
			e.UpdatedAt = rrow.UpdatedAt.Time
		}
		if rrow.DeletedAt.Valid {
			e.DeletedAt = &rrow.DeletedAt.Time
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
func (r *VaultGroupCategoryRepository) FindByID(ctx context.Context, id string) (*model.VaultCategory, error) {
	q := sqlc.New(r.db)
	rrow, err := q.FindVaultGroupCategoryById(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(localization.ErrorVaultGroupCategoryNotFound.Code)
		}
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	e := &model.VaultCategory{
		ID:         rrow.ID,
		Name:       rrow.Name,
		CoverImage: rrow.CoverImage,
		IsActive:   rrow.IsActive,
		IsDeleted:  rrow.IsDeleted,
	}
	if rrow.CreatedAt.Valid {
		e.CreatedAt = rrow.CreatedAt.Time
	}
	if rrow.UpdatedAt.Valid {
		e.UpdatedAt = rrow.UpdatedAt.Time
	}
	if rrow.DeletedAt.Valid {
		e.DeletedAt = &rrow.DeletedAt.Time
	}
	return e, nil
}

func (r *VaultGroupCategoryRepository) GetGroupcategoryByName(ctx context.Context, groupName string) (*model.VaultCategory, error) {
	q := sqlc.New(r.db)
	category, err := q.FindVaultGroupCategoryByName(ctx, groupName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(localization.ErrorVaultGroupCategoryNotFound.Code)
		}
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	vc := &model.VaultCategory{
		ID:         category.ID,
		Name:       category.Name,
		CoverImage: category.CoverImage,
		IsActive:   category.IsActive,
		IsDeleted:  category.IsDeleted,
	}
	if category.CreatedAt.Valid {
		vc.CreatedAt = category.CreatedAt.Time
	}
	if category.UpdatedAt.Valid {
		vc.UpdatedAt = category.UpdatedAt.Time
	}
	if category.DeletedAt.Valid {
		vc.DeletedAt = &category.DeletedAt.Time
	}

	return vc, nil
}

// Update updates name/description; prevents updates on deleted records
func (r *VaultGroupCategoryRepository) Update(ctx context.Context, id string, entity *model.VaultCategory) error {
	// Ensure not deleted
	current, err := r.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if current.IsDeleted {
		return fmt.Errorf("CANNOT_UPDATE_DELETED_VAULT_GROUP_CATEGORY %s", id)
	}

	name := sql.NullString{String: entity.Name, Valid: entity.Name != ""}
	coverImage := sql.NullString{String: entity.CoverImage, Valid: entity.CoverImage != ""}
	params := sqlc.UpdateVaultGroupCategoryParams{Name: name, CoverImage: coverImage, ID: id}
	_, err = sqlc.New(r.db).UpdateVaultGroupCategory(ctx, params)
	return err
}

// Delete performs a soft delete; prevents double delete
func (r *VaultGroupCategoryRepository) Delete(ctx context.Context, id string) (string, error) {
	current, err := r.FindByID(ctx, id)
	if err != nil {
		return "", err
	}
	if current.IsDeleted {
		return "", fmt.Errorf("VAULT_GROUP_CATEGORY_ALREADY_DELETED %s", id)
	}
	q := sqlc.New(r.db)
	ret, err := q.DeleteVaultGroupCategory(ctx, id)
	if err != nil {
		return "", err
	}
	if ret == "" {
		return "", fmt.Errorf("VAULT_GROUP_CATEGORY_ALREADY_DELETED %s", id)
	}
	return strings.ToUpper(ret), nil
}

// EnableOrDisable toggles active state; prevents toggling deleted records
func (r *VaultGroupCategoryRepository) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	current, err := r.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if current.IsDeleted {
		return fmt.Errorf("CANNOT_ENABLE_DISABLE_DELETED_VAULT_GROUP_CATEGORY %s", id)
	}
	q := sqlc.New(r.db)
	if enable {
		_, err = q.ActivateVaultGroupCategory(ctx, id)
		return err
	}
	_, err = q.DeactivateVaultGroupCategory(ctx, id)
	return err
}
