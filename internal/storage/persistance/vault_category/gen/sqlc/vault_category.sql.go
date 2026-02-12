package sqlc

import (
	"cbe-super-app-cps-action/pkgs/utils"
	"context"
	"database/sql"
	"fmt"
	"strings"

	imodel "cbe-super-app-cps-action/internal/constants/model"

	"github.com/google/uuid"
)

func generateUUID() string {
	return uuid.New().String()
}

const activateVaultCategory = `-- name: ActivateVaultCategory :one
UPDATE vault_categories
SET is_active = 1, updated_at = SYSTIMESTAMP
WHERE id = :1 AND is_deleted = 0
RETURNING id INTO :result`

func (q *Queries) ActivateVaultCategory(ctx context.Context, id string) (string, error) {
	var result string
	res, err := q.db.ExecContext(ctx, activateVaultCategory, id, sql.Out{Dest: &result})
	if err != nil {
		return "", fmt.Errorf("failed to activate vault  category: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return "", fmt.Errorf("no vault  category found with id %s", id)
	}
	return result, nil
}

const deactivateVaultCategory = `-- name: DeactivateVaultCategory :one
UPDATE vault_categories
SET is_active = 0, updated_at = SYSTIMESTAMP
WHERE id = :1 AND is_deleted = 0
RETURNING id INTO :result`

func (q *Queries) DeactivateVaultCategory(ctx context.Context, id string) (string, error) {
	var result string
	res, err := q.db.ExecContext(ctx, deactivateVaultCategory, id, sql.Out{Dest: &result})
	if err != nil {
		return "", fmt.Errorf("failed to deactivate vault  category: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return "", fmt.Errorf("no vault  category found with id %s", id)
	}
	return result, nil
}

const deleteVaultCategory = `-- name: DeleteVaultCategory :one
UPDATE vault_categories
SET is_deleted = 1, deleted_at = SYSTIMESTAMP, updated_at = SYSTIMESTAMP
WHERE id = :1
RETURNING id INTO :result`

func (q *Queries) DeleteVaultCategory(ctx context.Context, id string) (string, error) {
	var result string
	fmt.Println("===============", id)
	res, err := q.db.ExecContext(ctx, deleteVaultCategory, id, sql.Out{Dest: &result})
	if err != nil {
		return "", fmt.Errorf("failed to delete vault  category: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return "", fmt.Errorf("no vault  category found with id %s", id)
	}
	return result, nil
}

const findVaultCategory = `-- name: FindVaultCategory :many
SELECT
  id,
  name,
  category_type,
  cover_image,
  is_active,
  created_at,
  updated_at,
  COUNT(*) OVER() AS total_count
FROM vault_categories
WHERE (:is_active IS NULL OR is_active = :is_active)
  AND (:name IS NULL OR UPPER(name) LIKE :name)
  AND deleted_at IS NULL
ORDER BY created_at DESC
OFFSET NVL(:offset, 0) ROWS
FETCH NEXT NVL(:limit, 50) ROWS ONLY`

type FindVaultCategoryParams struct {
	IsActive  sql.NullBool   `json:"is_active"`
	NameQuery sql.NullString `json:"name"`
	Page      sql.NullInt64  `json:"page"`
	Limit     sql.NullInt64  `json:"limit"`
}

// type FindVaultCategoryRow struct {
// 	ID         string       `json:"id"`
// 	Name       string       `json:"name"`
// 	CoverImage string       `json:"cover_image"`
// 	IsActive   bool         `json:"is_active"`
// 	IsDeleted  bool         `json:"is_deleted"`
// 	CreatedAt  sql.NullTime `json:"created_at"`
// 	UpdatedAt  sql.NullTime `json:"updated_at"`
// 	DeletedAt  sql.NullTime `json:"deleted_at"`
// 	TotalCount int64        `json:"total_count"`
// }

func (q *Queries) FindVaultCategory(ctx context.Context, arg FindVaultCategoryParams) ([]VaultCategory, error) {
	namePtr := utils.NullStringToPtrLike(arg.NameQuery)
	limitPtr := utils.NullInt64ToPtr(arg.Limit)
	offsetPtr := (arg.Page.Int64 - 1) * arg.Limit.Int64

	rows, err := q.db.QueryContext(ctx, findVaultCategory,
		sql.Named("is_active", arg.IsActive),
		sql.Named("name", namePtr),
		sql.Named("offset", offsetPtr),
		sql.Named("limit", limitPtr),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []VaultCategory{}
	for rows.Next() {
		var i VaultCategory
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CategoryType,
			&i.CoverImage,
			&i.IsActive,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.TotalCount,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findVaultCategoryById = `-- name: FindVaultCategoryById :one
SELECT
  id,
  name,
  category_type,
  cover_image,
  is_active,
  created_at,
  updated_at
FROM vault_categories
WHERE id = :1 AND deleted_at IS NULL`

func (q *Queries) FindVaultCategoryById(ctx context.Context, id string) (VaultCategory, error) {
	row := q.db.QueryRowContext(ctx, findVaultCategoryById, id)
	var i VaultCategory
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CategoryType,
		&i.CoverImage,
		&i.IsActive,
		&i.CreatedAt,
		&i.UpdatedAt,
	)

	return i, err
}

const findVaultCategoryByName = `-- name: FindVaultCategoryByName :one
SELECT
  id,
  name,
  category_type,
  cover_image,
  is_active,
  created_at,
  updated_at
FROM vault_categories
WHERE UPPER(name) = UPPER(:name)`

func (q *Queries) FindVaultCategoryByName(ctx context.Context, name string) (VaultCategory, error) {
	row := q.db.QueryRowContext(ctx, findVaultCategoryByName, sql.Named("name", name))
	var i VaultCategory
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.CategoryType,
		&i.CoverImage,
		&i.IsActive,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const saveVaultCategory = `-- name: SaveVaultCategory :one
INSERT INTO vault_categories (
    id,
    name,
    cover_image_url,
    interest_type,
    category_interest,
    deadlock,
    is_active,
    created_at,
    updated_at
) VALUES (
    :1, :2, :3, :4, :5, :6, :7, SYSTIMESTAMP, SYSTIMESTAMP
)
RETURNING id INTO :8`

const saveVaultTier = `-- name: SaveVaultTier :exec
INSERT INTO vault_tiers (
    id,
    category_id,
    name,
    tier_interest,
    min_amount,
    max_amount
) VALUES (
    :1, :2, :3, :4, :5, :6
)`

func (q *Queries) SaveVaultCategory(ctx context.Context, arg *imodel.VaultCategory) (string, error) {
	categoryID := generateUUID()

	// Insert category
	_, err := q.db.ExecContext(
		ctx,
		saveVaultCategory,
		categoryID,
		strings.ToUpper(arg.Name),
		arg.CoverImageURL,
		arg.InterestType,
		arg.CategoryInterest,
		arg.Deadlock,
		arg.IsActive,
		sql.Out{Dest: &categoryID},
	)
	if err != nil {
		return "", fmt.Errorf("insert category failed: %w", err)
	}

	return categoryID, nil
}

func (q *Queries) SaveVaultTiers(ctx context.Context, categoryID string, arg *imodel.VaultCategory) error {
	for _, tier := range arg.Tiers {
		tierID := generateUUID()

		_, err := q.db.ExecContext(
			ctx,
			saveVaultTier,
			tierID,
			categoryID,
			tier.Name,
			tier.TierInterest,
			tier.MinAmount,
			tier.MaxAmount,
		)

		if err != nil {
			return fmt.Errorf("insert tier failed: %w", err)
		}
	}
	return nil
}

const updateVaultCategory = `-- name: UpdateVaultCategory :one
UPDATE vault_categories
SET
  name 				  = COALESCE(UPPER(:1), name),
  cover_image         = COALESCE(:2, cover_image),
  updated_at  = SYSTIMESTAMP
WHERE id = :3 AND is_deleted = 0
RETURNING id INTO :result`

type UpdateVaultCategoryParams struct {
	Name       sql.NullString `json:"name"`
	CoverImage sql.NullString `json:"cover_image"`
	ID         string         `json:"id"`
}

func (q *Queries) UpdateVaultCategory(ctx context.Context, arg UpdateVaultCategoryParams) (string, error) {
	var id string
	res, err := q.db.ExecContext(ctx, updateVaultCategory,
		arg.Name,
		arg.CoverImage,
		arg.ID,
		sql.Out{Dest: &id},
	)
	if err != nil {
		return "", fmt.Errorf("failed to update vault  category: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return "", fmt.Errorf("vault  category not found or already deleted")
	}
	return id, nil
}
