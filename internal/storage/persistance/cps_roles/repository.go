package cpsroles

import (
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/kafka"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	imodel "cbe-super-app-cps-action/internal/constants/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type cpsRoleStorage struct {
	cfg           *config.VaultConfig
	db            *sql.DB
	kafkaProducer kafka.ClientOrchestrationProducer
	logger        utils.Logger
}

func NewCPSRolesStorage(cfg *config.VaultConfig, db *sql.DB, kafkaProducer kafka.ClientOrchestrationProducer, logger utils.Logger) storage.CPSRolesRepository {
	return &cpsRoleStorage{
		cfg:           cfg,
		db:            db,
		kafkaProducer: kafkaProducer,
		logger:        logger,
	}
}

const superAppRoleTable = "SUPERAPP_ROLE"

func boolToOracleNumber(v bool) int {
	if v {
		return 1
	}
	return 0
}

func parseBoolFilter(v interface{}) (bool, bool) {
	switch t := v.(type) {
	case bool:
		return t, true
	case int:
		return t != 0, true
	case int32:
		return t != 0, true
	case int64:
		return t != 0, true
	case float64:
		return t != 0, true
	case string:
		switch strings.ToLower(strings.TrimSpace(t)) {
		case "true", "1", "yes", "y":
			return true, true
		case "false", "0", "no", "n":
			return false, true
		}
	}
	return false, false
}

func normalizeRawHex32(id string) (string, bool) {
	s := strings.TrimSpace(id)
	s = strings.TrimPrefix(s, "0x")
	s = strings.TrimPrefix(s, "0X")
	s = strings.ToLower(s)
	if len(s) != 32 {
		return "", false
	}
	for _, r := range s {
		if (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') {
			continue
		}
		return "", false
	}
	return s, true
}

func (m *cpsRoleStorage) Create(ctx context.Context, req imodel.CPSRoles) error {
	req.Name = strings.ToUpper(strings.TrimSpace(req.Name))
	req.RoleCode = strings.TrimSpace(req.RoleCode)

	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	now := time.Now()
	if req.CreatedAt.IsZero() {
		req.CreatedAt = now
	}
	if req.UpdatedAt.IsZero() {
		req.UpdatedAt = now
	}

	var id string
	const q = `
INSERT INTO SUPERAPP_ROLE (
  NAME,
  ROLE_CODE,
  DESCRIPTION,
  ENABLED,
  IS_DELETED,
  CREATED_AT,
  LAST_MODIFIED_AT
)
VALUES (
  UPPER(:1),:2,:3,:4,0,:5,:6
)
RETURNING RAWTOHEX(ID) INTO :7`

	if _, err := m.db.ExecContext(ctx, q,
		req.Name,
		req.RoleCode,
		req.Description,
		boolToOracleNumber(enabled),
		req.CreatedAt,
		req.UpdatedAt,
		sql.Out{Dest: &id},
	); err != nil {
		m.logger.Errorf("[CPSRolesStorage][Create] insert failed: %v", err)
		return local_util.HandleDBError(err)
	}

	req.ID = id
	req.IsDeleted = false
	req.Enabled = &enabled

	_ = m.kafkaProducer.PublishMessage(
		ctx,
		req,
		string(constants.ClientOrchestrationServicesTopic),
		string(constants.CustomerRoleUpdatedTopic),
		"cps customer role created",
	)

	return nil
}

func (m *cpsRoleStorage) Update(ctx context.Context, id string, req imodel.CPSRoles) error {
	idHex, ok := normalizeRawHex32(id)
	if !ok {
		return errors.New(localization.ErrorInvalidID.Code)
	}

	sets := []string{"LAST_MODIFIED_AT = SYSTIMESTAMP"}
	var args []interface{}

	if strings.TrimSpace(req.Name) != "" {
		sets = append(sets, "NAME = :name")
		args = append(args, sql.Named("name", strings.ToUpper(strings.TrimSpace(req.Name))))
	}
	if strings.TrimSpace(req.RoleCode) != "" {
		sets = append(sets, "ROLE_CODE = :role_code")
		args = append(args, sql.Named("role_code", strings.TrimSpace(req.RoleCode)))
	}
	if strings.TrimSpace(req.Description) != "" {
		sets = append(sets, "DESCRIPTION = :description")
		args = append(args, sql.Named("description", strings.TrimSpace(req.Description)))
	}

	if len(sets) == 1 {
		return errors.New(localization.ErrorNoDataProvided.Code)
	}

	q := fmt.Sprintf(`
UPDATE %s
SET %s
WHERE ID = HEXTORAW(:id) AND IS_DELETED = 0`, superAppRoleTable, strings.Join(sets, ", "))
	args = append(args, sql.Named("id", idHex))

	res, err := m.db.ExecContext(ctx, q, args...)
	if err != nil {
		m.logger.Errorf("[CPSRolesStorage][Update] update failed: %v", err)
		return local_util.HandleDBError(err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}

	updated, _ := m.FindById(ctx, idHex)
	if updated != nil {
		_ = m.kafkaProducer.PublishMessage(
			ctx,
			updated,
			string(constants.ClientOrchestrationServicesTopic),
			string(constants.CustomerRoleUpdatedTopic),
			"cps customer role updated",
		)
	}

	return nil
}

func (m *cpsRoleStorage) FindAllWithPagination(ctx context.Context, filterParam *types.Filter) (*types.PaginatedResponse[[]imodel.CPSRoles], error) {
	if filterParam == nil {
		filterParam = &types.Filter{}
	}

	limit := int64(50)
	page := int64(1)
	if filterParam.PerPage > 0 {
		limit = int64(filterParam.PerPage)
	}
	if filterParam.Page > 0 {
		page = int64(filterParam.Page)
	}
	offset := (page - 1) * limit

	clauses := []string{"IS_DELETED = 0"}
	var args []interface{}

	search := strings.TrimSpace(filterParam.Search)
	if search != "" {
		clauses = append(clauses, "LOWER(NAME) LIKE '%' || LOWER(:search) || '%'")
		args = append(args, sql.Named("search", search))
	}

	if filterParam.Filters != nil {
		if v, ok := filterParam.Filters["enabled"]; ok {
			if b, ok2 := parseBoolFilter(v); ok2 {
				clauses = append(clauses, "ENABLED = :enabled")
				args = append(args, sql.Named("enabled", boolToOracleNumber(b)))
			}
		}
	}

	where := strings.Join(clauses, " AND ")

	countQ := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE %s`, superAppRoleTable, where)
	var total int64
	if err := m.db.QueryRowContext(ctx, countQ, args...).Scan(&total); err != nil {
		m.logger.Errorf("[CPSRolesStorage][FindAllWithPagination] count failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	listQ := fmt.Sprintf(`
SELECT
  RAWTOHEX(ID),
  NAME,
  ROLE_CODE,
  DESCRIPTION,
  ENABLED,
  IS_DELETED,
  CREATED_AT,
  LAST_MODIFIED_AT,
  DELETED_AT
FROM %s
WHERE %s
ORDER BY CREATED_AT DESC
OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`, superAppRoleTable, where)

	listArgs := append(args, sql.Named("offset", offset), sql.Named("limit", limit))
	rows, err := m.db.QueryContext(ctx, listQ, listArgs...)
	if err != nil {
		m.logger.Errorf("[CPSRolesStorage][FindAllWithPagination] list query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer rows.Close()

	var out []imodel.CPSRoles
	for rows.Next() {
		var (
			id                         string
			name, roleCode, desc       sql.NullString
			enabledN, isDeletedN       int
			createdAt, updatedAt, delT sql.NullTime
		)
		if err := rows.Scan(&id, &name, &roleCode, &desc, &enabledN, &isDeletedN, &createdAt, &updatedAt, &delT); err != nil {
			m.logger.Errorf("[CPSRolesStorage][FindAllWithPagination] scan failed: %v", err)
			return nil, local_util.HandleDBError(err)
		}

		enabled := enabledN == 1
		role := imodel.CPSRoles{
			ID:          id,
			Name:        name.String,
			RoleCode:    roleCode.String,
			Description: desc.String,
			Enabled:     &enabled,
			IsDeleted:   isDeletedN == 1,
		}
		if createdAt.Valid {
			role.CreatedAt = createdAt.Time
		}
		if updatedAt.Valid {
			role.UpdatedAt = updatedAt.Time
		}
		if delT.Valid {
			role.DeletedAt = delT.Time
		}

		out = append(out, role)
	}

	meta := local_util.BuildPaginationMeta(total, int(page), int(limit))
	return &types.PaginatedResponse[[]imodel.CPSRoles]{Data: out, Meta: meta}, nil
}

func (m *cpsRoleStorage) FindById(ctx context.Context, id string) (*imodel.CPSRoles, error) {
	idHex, ok := normalizeRawHex32(id)
	if !ok {
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	const q = `
SELECT
  RAWTOHEX(ID),
  NAME,
  ROLE_CODE,
  DESCRIPTION,
  ENABLED,
  IS_DELETED,
  CREATED_AT,
  LAST_MODIFIED_AT,
  DELETED_AT
FROM SUPERAPP_ROLE
WHERE ID = HEXTORAW(:1) AND IS_DELETED = 0`

	var (
		roleID                    string
		name, roleCode, desc      sql.NullString
		enabledN, isDeletedN      int
		createdAt, updatedAt, del sql.NullTime
	)

	err := m.db.QueryRowContext(ctx, q, idHex).Scan(
		&roleID,
		&name,
		&roleCode,
		&desc,
		&enabledN,
		&isDeletedN,
		&createdAt,
		&updatedAt,
		&del,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		m.logger.Errorf("[CPSRolesStorage][FindById] query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	enabled := enabledN == 1
	role := &imodel.CPSRoles{
		ID:          roleID,
		Name:        name.String,
		RoleCode:    roleCode.String,
		Description: desc.String,
		Enabled:     &enabled,
		IsDeleted:   isDeletedN == 1,
	}
	if createdAt.Valid {
		role.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		role.UpdatedAt = updatedAt.Time
	}
	if del.Valid {
		role.DeletedAt = del.Time
	}

	// Oracle implementation currently does not populate Maker/Checker/Auditor or service access lists.
	role.MakerActions = nil
	role.CheckerActions = nil
	role.AuditorActions = nil
	role.EnabledServices = nil
	role.DisabledServices = nil

	return role, nil
}

func (m *cpsRoleStorage) FindByNameOrRoleCode(ctx context.Context, name, roleCode string) (*imodel.CPSRoles, error) {
	name = strings.TrimSpace(name)
	roleCode = strings.TrimSpace(roleCode)
	if name == "" && roleCode == "" {
		return nil, errors.New(localization.ErrorNoDataProvided.Code)
	}

	clauses := []string{}
	var args []interface{}
	if name != "" {
		clauses = append(clauses, "LOWER(NAME) = LOWER(:name)")
		args = append(args, sql.Named("name", name))
	}
	if roleCode != "" {
		clauses = append(clauses, "ROLE_CODE = :role_code")
		args = append(args, sql.Named("role_code", roleCode))
	}

	cond := strings.Join(clauses, " OR ")
	q := fmt.Sprintf(`
SELECT
  RAWTOHEX(ID),
  NAME,
  ROLE_CODE,
  DESCRIPTION,
  ENABLED,
  IS_DELETED,
  CREATED_AT,
  LAST_MODIFIED_AT,
  DELETED_AT
FROM %s
WHERE IS_DELETED = 0 AND (%s)
FETCH FIRST 1 ROWS ONLY`, superAppRoleTable, cond)

	var (
		id                         string
		dbName, dbRole, dbDesc     sql.NullString
		enabledN, isDeletedN       int
		createdAt, updatedAt, delT sql.NullTime
	)

	err := m.db.QueryRowContext(ctx, q, args...).Scan(
		&id,
		&dbName,
		&dbRole,
		&dbDesc,
		&enabledN,
		&isDeletedN,
		&createdAt,
		&updatedAt,
		&delT,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		m.logger.Errorf("[CPSRolesStorage][FindByNameOrRoleCode] query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	enabled := enabledN == 1
	out := &imodel.CPSRoles{
		ID:          id,
		Name:        dbName.String,
		RoleCode:    dbRole.String,
		Description: dbDesc.String,
		Enabled:     &enabled,
		IsDeleted:   isDeletedN == 1,
	}
	if createdAt.Valid {
		out.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		out.UpdatedAt = updatedAt.Time
	}
	if delT.Valid {
		out.DeletedAt = delT.Time
	}
	return out, nil
}

func (m *cpsRoleStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	idHex, ok := normalizeRawHex32(id)
	if !ok {
		return errors.New(localization.ErrorInvalidID.Code)
	}

	const q = `
UPDATE SUPERAPP_ROLE
SET
  ENABLED    = :1,
  LAST_MODIFIED_AT = SYSTIMESTAMP
WHERE ID = HEXTORAW(:2) AND IS_DELETED = 0`

	res, err := m.db.ExecContext(ctx, q, boolToOracleNumber(enable), idHex)
	if err != nil {
		m.logger.Errorf("[CPSRolesStorage][EnableOrDisable] update failed: %v", err)
		return local_util.HandleDBError(err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}
	return nil
}

func (m *cpsRoleStorage) FindByCustomerSegmentation(ctx context.Context, customerSegment string) (*imodel.CPSRoles, error) {
	customerSegment = strings.TrimSpace(customerSegment)
	if customerSegment == "" {
		return nil, errors.New(localization.ErrorNoDataProvided.Code)
	}

	// Resolve superapp role for a customer segmentation code (CUSTOMER_SEGMENTATIONS.NAME) via sub-segments.
	const q = `
SELECT
  RAWTOHEX(cr.ID),
  cr.NAME,
  cr.ROLE_CODE,
  cr.DESCRIPTION,
  cr.ENABLED,
  cr.IS_DELETED,
  cr.CREATED_AT,
  cr.LAST_MODIFIED_AT,
  cr.DELETED_AT
FROM CUSTOMER_SEGMENTATIONS cs
JOIN CUSTOMER_SUB_SEGMENTS css
  ON css.CUSTOMER_SEGMENTATIONS_ID = cs.ID AND css.IS_DELETED = 0 AND css.IS_ENABLED = 1
JOIN SUPERAPP_ROLE cr
  ON cr.ID = css.SUPERAPP_ROLE_ID AND cr.IS_DELETED = 0 AND cr.ENABLED = 1
WHERE cs.IS_DELETED = 0
  AND cs.IS_ENABLED = 1
  AND UPPER(TRIM(cs.NAME)) = UPPER(TRIM(:1))
FETCH FIRST 1 ROWS ONLY`

	var (
		id                         string
		name, roleCode, desc       sql.NullString
		enabledN, isDeletedN       int
		createdAt, updatedAt, delT sql.NullTime
	)

	err := m.db.QueryRowContext(ctx, q, customerSegment).Scan(
		&id,
		&name,
		&roleCode,
		&desc,
		&enabledN,
		&isDeletedN,
		&createdAt,
		&updatedAt,
		&delT,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		m.logger.Errorf("[CPSRolesStorage][FindByCustomerSegmentation] query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	enabled := enabledN == 1
	out := &imodel.CPSRoles{
		ID:          id,
		Name:        name.String,
		RoleCode:    roleCode.String,
		Description: desc.String,
		Enabled:     &enabled,
		IsDeleted:   isDeletedN == 1,
	}
	if createdAt.Valid {
		out.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		out.UpdatedAt = updatedAt.Time
	}
	if delT.Valid {
		out.DeletedAt = delT.Time
	}
	return out, nil
}

func (m *cpsRoleStorage) EnableServiceAccess(ctx context.Context, roleID string, accessListKeys []string) error {
	// Oracle schema for access-list segmentation is not yet available in this service.
	m.logger.Warnf("[CPSRolesStorage][EnableServiceAccess] not implemented for oracle; roleID=%s keys=%v", roleID, accessListKeys)
	return errors.New(localization.ErrorUnexpectedError.Code)
}

func (m *cpsRoleStorage) DisableServiceAccess(ctx context.Context, roleID string, accessListKeys []string) error {
	// Oracle schema for access-list segmentation is not yet available in this service.
	m.logger.Warnf("[CPSRolesStorage][DisableServiceAccess] not implemented for oracle; roleID=%s keys=%v", roleID, accessListKeys)
	return errors.New(localization.ErrorUnexpectedError.Code)
}

func (m *cpsRoleStorage) Delete(ctx context.Context, id string) error {
	idHex, ok := normalizeRawHex32(id)
	if !ok {
		return errors.New(localization.ErrorInvalidID.Code)
	}

	const q = `
UPDATE SUPERAPP_ROLE
SET
  IS_DELETED = 1,
  LAST_MODIFIED_AT = SYSTIMESTAMP,
  DELETED_AT = SYSTIMESTAMP
WHERE ID = HEXTORAW(:1) AND IS_DELETED = 0`

	res, err := m.db.ExecContext(ctx, q, idHex)
	if err != nil {
		m.logger.Errorf("[CPSRolesStorage][Delete] delete failed: %v", err)
		return local_util.HandleDBError(err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}
	return nil
}
