package sqlc

import (
	"cbe-super-app-cps-action/pkgs/utils"
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func generateUUID() string {
	return uuid.New().String()
}

const activateVaultGroupCategory = `-- name: ActivateVaultGroupCategory :one
UPDATE vault_categories
SET is_active = 1, updated_at = SYSTIMESTAMP
WHERE id = :1 AND is_deleted = 0
RETURNING id INTO :result`

func (q *Queries) ActivateVaultGroupCategory(ctx context.Context, id string) (string, error) {
	var result string
	res, err := q.db.ExecContext(ctx, activateVaultGroupCategory, id, sql.Out{Dest: &result})
	if err != nil {
		return "", fmt.Errorf("failed to activate vault group category: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return "", fmt.Errorf("no vault group category found with id %s", id)
	}
	return result, nil
}

const deactivateVaultGroupCategory = `-- name: DeactivateVaultGroupCategory :one
UPDATE vault_categories
SET is_active = 0, updated_at = SYSTIMESTAMP
WHERE id = :1 AND is_deleted = 0
RETURNING id INTO :result`

func (q *Queries) DeactivateVaultGroupCategory(ctx context.Context, id string) (string, error) {
	var result string
	res, err := q.db.ExecContext(ctx, deactivateVaultGroupCategory, id, sql.Out{Dest: &result})
	if err != nil {
		return "", fmt.Errorf("failed to deactivate vault group category: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return "", fmt.Errorf("no vault group category found with id %s", id)
	}
	return result, nil
}

const deleteVaultGroupCategory = `-- name: DeleteVaultGroupCategory :one
UPDATE vault_categories
SET is_deleted = 1, deleted_at = SYSTIMESTAMP, updated_at = SYSTIMESTAMP
WHERE id = :1
RETURNING id INTO :result`

func (q *Queries) DeleteVaultGroupCategory(ctx context.Context, id string) (string, error) {
	var result string
	fmt.Println("===============", id)
	res, err := q.db.ExecContext(ctx, deleteVaultGroupCategory, id, sql.Out{Dest: &result})
	if err != nil {
		return "", fmt.Errorf("failed to delete vault group category: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return "", fmt.Errorf("no vault group category found with id %s", id)
	}
	return result, nil
}

const findVaultGroupCategory = `-- name: FindVaultGroupCategory :many
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

type FindVaultGroupCategoryParams struct {
	IsActive  sql.NullBool   `json:"is_active"`
	NameQuery sql.NullString `json:"name"`
	Page      sql.NullInt64  `json:"page"`
	Limit     sql.NullInt64  `json:"limit"`
}

// type FindVaultGroupCategoryRow struct {
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

func (q *Queries) FindVaultGroupCategory(ctx context.Context, arg FindVaultGroupCategoryParams) ([]VaultCategory, error) {
	namePtr := utils.NullStringToPtrLike(arg.NameQuery)
	limitPtr := utils.NullInt64ToPtr(arg.Limit)
	offsetPtr := (arg.Page.Int64 - 1) * arg.Limit.Int64

	rows, err := q.db.QueryContext(ctx, findVaultGroupCategory,
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

const findVaultGroupCategoryById = `-- name: FindVaultGroupCategoryById :one
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

func (q *Queries) FindVaultGroupCategoryById(ctx context.Context, id string) (VaultCategory, error) {
	row := q.db.QueryRowContext(ctx, findVaultGroupCategoryById, id)
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

const findVaultGroupCategoryByName = `-- name: FindVaultGroupCategoryByName :one
SELECT id, name, is_active, is_deleted, created_at, updated_at, deleted_at
FROM vault_categories
WHERE UPPER(name) = UPPER(:name)`

func (q *Queries) FindVaultGroupCategoryByName(ctx context.Context, name string) (VaultCategory, error) {
	row := q.db.QueryRowContext(ctx, findVaultGroupCategoryByName, sql.Named("name", name))
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

const saveVaultGroupCategory = `-- name: SaveVaultGroupCategory :one
INSERT INTO vault_categories (
  id,
  name,
  category_type,
  cover_image,
  is_active
) VALUES (
  UPPER(:1),
  :2,
  :3,
  :4,
  :5
)
RETURNING id INTO :result`

type SaveVaultGroupCategoryParams struct {
	Name         string       `json:"name"`
	CategoryType string       `json:"category_type"`
	CoverImage   string       `json:"cover_image"`
	IsActive     sql.NullBool `json:"is_active"`
}

func (q *Queries) SaveVaultGroupCategory(ctx context.Context, arg SaveVaultGroupCategoryParams) (string, error) {
	var id string
	_, err := q.db.ExecContext(ctx, saveVaultGroupCategory,
		generateUUID(),
		strings.ToUpper(arg.Name),
		arg.CategoryType,
		arg.CoverImage,
		utils.NullBoolToInt(arg.IsActive),
		sql.Out{Dest: &id},
	)
	if err != nil {
		return "", fmt.Errorf("failed to save vault group category: %w", err)
	}
	return id, nil
}

const updateVaultGroupCategory = `-- name: UpdateVaultGroupCategory :one
UPDATE vault_categories
SET
  name 				  = COALESCE(UPPER(:1), name),
  cover_image         = COALESCE(:2, cover_image),
  updated_at  = SYSTIMESTAMP
WHERE id = :3 AND is_deleted = 0
RETURNING id INTO :result`

type UpdateVaultGroupCategoryParams struct {
	Name       sql.NullString `json:"name"`
	CoverImage sql.NullString `json:"cover_image"`
	ID         string         `json:"id"`
}

func (q *Queries) UpdateVaultGroupCategory(ctx context.Context, arg UpdateVaultGroupCategoryParams) (string, error) {
	var id string
	res, err := q.db.ExecContext(ctx, updateVaultGroupCategory,
		arg.Name,
		arg.CoverImage,
		arg.ID,
		sql.Out{Dest: &id},
	)
	if err != nil {
		return "", fmt.Errorf("failed to update vault group category: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return "", fmt.Errorf("vault group category not found or already deleted")
	}
	return id, nil
}
