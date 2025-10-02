-- name: SaveVaultGroupCategory :one
INSERT INTO group_vault_categories (
  name,
  description,
  is_active
) VALUES (
  $1,
  $2,
  $3
)
RETURNING id;

-- name: UpdateVaultGroupCategory :one
UPDATE group_vault_categories
SET
  name        = COALESCE($1,        name),
  description = COALESCE($2, description),
  updated_at  = NOW()
WHERE id = $3 AND is_deleted = FALSE
RETURNING id;

-- name: ActivateVaultGroupCategory :one
UPDATE group_vault_categories
SET is_active = TRUE, updated_at = NOW()
WHERE id = $1 AND is_deleted = FALSE
RETURNING id;

-- name: DeactivateVaultGroupCategory :one
UPDATE group_vault_categories
SET is_active = FALSE, updated_at = NOW()
WHERE id = $1 AND is_deleted = FALSE
RETURNING id;

-- name: DeleteVaultGroupCategory :one
UPDATE group_vault_categories
SET is_deleted = TRUE, deleted_at = NOW(), updated_at = NOW()
WHERE id = $1 AND is_deleted = FALSE
RETURNING id;

-- name: FindVaultGroupCategory :many
SELECT
  id,
  name,
  description,
  is_active,
  is_deleted,
  created_at,
  updated_at,
  deleted_at,
  COUNT(*) OVER() AS total_count
FROM group_vault_categories
WHERE ($1::boolean IS NULL OR is_active = $1::boolean)
  AND ($2::text IS NULL OR name ILIKE ('%' || $2::text || '%'))
ORDER BY created_at DESC
LIMIT  COALESCE($4::int, 50)
OFFSET COALESCE($3::int, 0);

-- name: FindVaultGroupCategoryById :one
SELECT
  id,
  name,
  description,
  is_active,
  is_deleted,
  created_at,
  updated_at,
  deleted_at
FROM group_vault_categories
WHERE id = $1;

-- name: FindVaultGroupCategoryByName :one
SELECT id, name, description, is_active, is_deleted, created_at, updated_at, deleted_at
FROM group_vault_categories
WHERE name = $1;


