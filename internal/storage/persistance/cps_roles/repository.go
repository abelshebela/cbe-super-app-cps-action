package cpsroles

import (
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/kafka"
	"context"
	"database/sql"
	"encoding/json"
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

// CheckUserExistence implements [storage.CPSRolesRepository].
func (m *cpsRoleStorage) CheckUserExistence(ctx context.Context, roleID string) error {
	log := local_util.LoggerFromCtx(ctx, m.logger)

	log.Infof("[CPSRolesStorage][CheckUserExistence] checking user existence for roleID: %s", roleID)
	const q = `SELECT COUNT(*) FROM USERS WHERE CUSTOMER_SEGMENTATION = :1`
	var count int
	if err := m.db.QueryRowContext(ctx, q, roleID).Scan(&count); err != nil {
		log.Errorf("[CPSRolesStorage][CheckUserExistence] query failed: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return local_util.HandleDBError(err)
	}
	if count > 0 {
		return fmt.Errorf("there are still %d users under this role", count)
	}
	return nil
}

func NewCPSRolesStorage(cfg *config.VaultConfig, db *sql.DB, kafkaProducer kafka.ClientOrchestrationProducer, logger utils.Logger) storage.CPSRolesRepository {
	return &cpsRoleStorage{
		cfg:           cfg,
		db:            db,
		kafkaProducer: kafkaProducer,
		logger:        logger,
	}
}

const superAppRoleTable = "SUPERAPP_ROLES"
const accessListTable = "ACCESS_LISTS"
const accessListCustomerSegmentationTable = "ACCESS_LIST_BY_SUPERAPP_ROLE"

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
	log := local_util.LoggerFromCtx(ctx, m.logger)

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
	INSERT INTO SUPERAPP_ROLES (
	NAME,
	LABEL,
	ROLE_CODE,
	DESCRIPTION,
	IS_ENABLED,
	IS_DELETED,
	CREATED_AT,
	LAST_MODIFIED_AT
	)
	VALUES (
	UPPER(:1), :2, :3, :4, :5, 0, :6, :7
	)
	RETURNING RAWTOHEX(ID) INTO :8`

	if _, err := m.db.ExecContext(ctx, q,
		req.Name,
		req.Lable,
		req.RoleCode,
		req.Description,
		boolToOracleNumber(enabled),
		req.CreatedAt,
		req.UpdatedAt,
		sql.Out{Dest: &id},
	); err != nil {
		log.Errorf("[CPSRolesStorage][Create] insert failed: %v", err)
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
	log := local_util.LoggerFromCtx(ctx, m.logger)

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
	if strings.TrimSpace(req.Lable) != "" {
		sets = append(sets, "LABEL = :label")
		args = append(args, sql.Named("label", strings.TrimSpace(req.Lable)))
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
		log.Errorf("[CPSRolesStorage][Update] update failed: %v", err)
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
	log := local_util.LoggerFromCtx(ctx, m.logger)

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

	// clauses := []string{"IS_DELETED = 0"}
	// var args []interface{}

	// search := strings.TrimSpace(filterParam.Search)
	// if search != "" {
	// 	clauses = append(clauses, "LOWER(NAME) LIKE '%' || LOWER(:search) || '%'")
	// 	args = append(args, sql.Named("search", search))
	// }

	// if filterParam.Filters != nil {
	// 	if v, ok := filterParam.Filters["enabled"]; ok {
	// 		if b, ok2 := parseBoolFilter(v); ok2 {
	// 			clauses = append(clauses, "IS_ENABLED = :enabled")
	// 			args = append(args, sql.Named("enabled", boolToOracleNumber(b)))
	// 		}
	// 	}
	// }

	// where := strings.Join(clauses, " AND ")

	clauses := []string{"SR.IS_DELETED = 0"}
	var args []interface{}

	search := strings.TrimSpace(filterParam.Search)
	if search != "" {
		clauses = append(clauses, "LOWER(SR.NAME) LIKE '%' || LOWER(:search) || '%'")
		args = append(args, sql.Named("search", search))
	}

	if filterParam.Filters != nil {
		if v, ok := filterParam.Filters["enabled"]; ok {
			if b, ok2 := parseBoolFilter(v); ok2 {
				clauses = append(clauses, "SR.IS_ENABLED = :enabled")
				args = append(args, sql.Named("enabled", boolToOracleNumber(b)))
			}
		}
	}

	where := strings.Join(clauses, " AND ")

	countQ := fmt.Sprintf(`
		SELECT COUNT(1)
		FROM %s SR
		WHERE %s
	`, superAppRoleTable, where)

	var total int64

	err := m.db.QueryRowContext(ctx, countQ, args...).Scan(&total)
	if err != nil {
		log.Errorf("[CPSRolesStorage][countCPSRoles] count failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	// 	listQ := fmt.Sprintf(`
	// SELECT
	//   RAWTOHEX(SR.ID),
	//   SR.NAME,
	//   SR.ROLE_CODE,
	//   SR.DESCRIPTION,
	//   SR.IS_ENABLED,
	//   SR.IS_DELETED,
	//   SR.CREATED_AT,
	//   SR.LAST_MODIFIED_AT,
	//   SR.DELETED_AT
	//   (*AL AS SR.ENABLED_SERVICES),
	// FROM %s AS SR JOIN %s AS AL ON AL.SUPERAPP_ROLE_ID = SR.ID AND AL.IS_DELETED = 0
	// WHERE %s
	// ORDER BY CREATED_AT DESC
	// OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`, superAppRoleTable, accessListTable, where)

	selectQuery := FetchCPSRolesQuery(where)
	listArgs := append(args, sql.Named("offset", offset), sql.Named("limit", limit))
	rows, err := m.db.QueryContext(ctx, selectQuery, listArgs...)
	if err != nil {
		log.Errorf("[CPSRolesStorage][FindAllWithPagination] list query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer rows.Close()

	var out []imodel.CPSRoles

	for rows.Next() {
		var (
			id                          string
			name, roleCode, label, desc sql.NullString
			enabledN, isDeletedN        int
			createdAt, updatedAt        sql.NullTime
			delT                        sql.NullTime

			enabledServicesJSON  sql.NullString
			disabledServicesJSON sql.NullString
		)

		err := rows.Scan(
			&id,
			&name,
			&label,
			&roleCode,
			&desc,
			&enabledN,
			&isDeletedN,
			&createdAt,
			&updatedAt,
			&delT,
		)

		if err != nil {
			log.Errorf("[CPSRolesStorage][FindAllWithPagination] scan failed: %v", err)
			return nil, local_util.HandleDBError(err)
		}

		enabled := enabledN == 1

		role := imodel.CPSRoles{
			ID:          id,
			Name:        name.String,
			RoleCode:    roleCode.String,
			Lable:       label.String,
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

		// -----------------------------
		// Parse ENABLED_SERVICES JSON
		// -----------------------------
		if enabledServicesJSON.Valid && enabledServicesJSON.String != "" {
			_ = json.Unmarshal([]byte(enabledServicesJSON.String), &role.EnabledServices)
		}

		// -----------------------------
		// Parse DISABLED_SERVICES JSON
		// -----------------------------
		if disabledServicesJSON.Valid && disabledServicesJSON.String != "" {
			_ = json.Unmarshal([]byte(disabledServicesJSON.String), &role.DisabledServices)
		}

		out = append(out, role)
	}

	meta := local_util.BuildPaginationMeta(total, int(page), int(limit))
	return &types.PaginatedResponse[[]imodel.CPSRoles]{Data: out, Meta: meta}, nil
}

func (m *cpsRoleStorage) FindById(ctx context.Context, id string) (*imodel.CPSRoles, error) {
	log := local_util.LoggerFromCtx(ctx, m.logger)

	idHex, ok := normalizeRawHex32(id)
	if !ok {
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	const q = `
SELECT
  RAWTOHEX(SR.ID) AS ID,
  NVL(SR.NAME, '') AS NAME,
  NVL(SR.LABEL, '') AS LABEL,
  NVL(SR.ROLE_CODE, '') AS ROLE_CODE,
  NVL(SR.DESCRIPTION, '') AS DESCRIPTION,
  SR.IS_ENABLED,
  SR.IS_DELETED,
  SR.CREATED_AT,
  SR.LAST_MODIFIED_AT,
  SR.DELETED_AT,

  -- ENABLED SERVICES
  NVL(
    (
      SELECT JSON_ARRAYAGG(
        JSON_OBJECT(
          'key' VALUE RAWTOHEX(AL.ID),
          'access_list_name' VALUE AL.NAME,
          'access_list_id' VALUE AL.SERVICE_KEY
        ) RETURNING CLOB
      )
      FROM ACCESS_LISTS AL
      WHERE AL.IS_ENABLED = 1
        AND AL.IS_DELETED = 0
        AND NOT EXISTS (
          SELECT 1
          FROM ACCESS_LIST_BY_SUPERAPP_ROLE ACS
          WHERE ACS.ACCESS_LIST_ID = AL.SERVICE_KEY
            AND ACS.SUPERAPP_ROLE_ID = SR.ID
            AND ACS.IS_ENABLED = 1
        )
    ),
    TO_CLOB('[]')
  ) AS ENABLED_SERVICES,

  -- DISABLED SERVICES
  NVL(
    (
      SELECT JSON_ARRAYAGG(
        JSON_OBJECT(
          'key' VALUE RAWTOHEX(AL.ID),
          'access_list_name' VALUE AL.NAME,
          'access_list_id' VALUE AL.SERVICE_KEY
        ) RETURNING CLOB
      )
      FROM ACCESS_LISTS AL
      WHERE AL.IS_ENABLED = 1
        AND AL.IS_DELETED = 0
        AND EXISTS (
          SELECT 1
          FROM ACCESS_LIST_BY_SUPERAPP_ROLE ACS
          WHERE ACS.ACCESS_LIST_ID = AL.SERVICE_KEY
            AND ACS.SUPERAPP_ROLE_ID = SR.ID
            AND ACS.IS_ENABLED = 1
        )
    ),
    TO_CLOB('[]')
  ) AS DISABLED_SERVICES

FROM SUPERAPP_ROLES SR
WHERE SR.ID = HEXTORAW(:1)
  AND SR.IS_DELETED = 0
`

	var (
		r imodel.CPSRoles

		name, roleCode, label, desc sql.NullString
		enabledN, deletedN          int

		createdAt, updatedAt, deletedAt sql.NullTime

		enabledServicesJSON  sql.NullString
		disabledServicesJSON sql.NullString
	)

	err := m.db.QueryRowContext(ctx, q, idHex).Scan(
		&r.ID,
		&name,
		&label,
		&roleCode,
		&desc,
		&enabledN,
		&deletedN,
		&createdAt,
		&updatedAt,
		&deletedAt,
		&enabledServicesJSON,
		&disabledServicesJSON,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		log.Errorf("[CPSRolesStorage][FindById] query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	// -------------------
	// Basic mapping
	// -------------------
	r.Name = name.String
	r.Lable = label.String
	r.RoleCode = roleCode.String
	r.Description = desc.String
	r.Enabled = PtrBool(enabledN == 1)
	r.IsDeleted = deletedN == 1

	SetTime(&r.CreatedAt, createdAt)
	SetTime(&r.UpdatedAt, updatedAt)
	SetTime(&r.DeletedAt, deletedAt)

	// -------------------
	// Services mapping
	// -------------------
	UnmarshalJSON(enabledServicesJSON, &r.EnabledServices)
	UnmarshalJSON(disabledServicesJSON, &r.DisabledServices)

	// normalize null → empty array
	if r.EnabledServices == nil {
		r.EnabledServices = []imodel.ServiceAccessInfo{}
	}
	if r.DisabledServices == nil {
		r.DisabledServices = []imodel.ServiceAccessInfo{}
	}

	// -------------------
	// Not implemented yet
	// -------------------
	r.MakerActions = nil
	r.CheckerActions = nil
	r.AuditorActions = nil

	return &r, nil
}

func (m *cpsRoleStorage) FindSupperAppRoleByAccessList(ctx context.Context, accessListID string) (bool, error) {
	log := local_util.LoggerFromCtx(ctx, m.logger)

	const q = `SELECT ID FROM ACCESS_LIST_BY_SUPERAPP_ROLE WHERE SUPERAPP_ROLE_ID = :1 AND IS_DELETED = 0`

	var id string
	err := m.db.QueryRowContext(ctx, q, accessListID).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		log.Errorf("[CPSRolesStorage][FindSupperAppRoleByAccessList] query failed: %v", err)
		return false, err
	}
	return true, nil
}

func (m *cpsRoleStorage) FindByNameOrRoleCode(ctx context.Context, name, roleCode string) (*imodel.CPSRoles, error) {
	log := local_util.LoggerFromCtx(ctx, m.logger)

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
  LABEL,
  ROLE_CODE,
  DESCRIPTION,
  IS_ENABLED,
  IS_DELETED,
  CREATED_AT,
  LAST_MODIFIED_AT,
  DELETED_AT
FROM %s
WHERE IS_DELETED = 0 AND (%s)
FETCH FIRST 1 ROWS ONLY`, superAppRoleTable, cond)

	var (
		id                            string
		dbName, dbRole, label, dbDesc sql.NullString
		enabledN, isDeletedN          int
		createdAt, updatedAt, delT    sql.NullTime
	)

	err := m.db.QueryRowContext(ctx, q, args...).Scan(
		&id,
		&dbName,
		&label,
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
		log.Errorf("[CPSRolesStorage][FindByNameOrRoleCode] query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	enabled := enabledN == 1
	out := &imodel.CPSRoles{
		ID:          id,
		Name:        dbName.String,
		Lable:       label.String,
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
	log := local_util.LoggerFromCtx(ctx, m.logger)

	idHex, ok := normalizeRawHex32(id)
	if !ok {
		return errors.New(localization.ErrorInvalidID.Code)
	}

	const q = `
UPDATE SUPERAPP_ROLES
SET
  IS_ENABLED    = :1,
  LAST_MODIFIED_AT = SYSTIMESTAMP
WHERE ID = HEXTORAW(:2) AND IS_DELETED = 0`

	res, err := m.db.ExecContext(ctx, q, boolToOracleNumber(enable), idHex)
	if err != nil {
		log.Errorf("[CPSRolesStorage][EnableOrDisable] update failed: %v", err)
		return local_util.HandleDBError(err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}
	return nil
}

func (m *cpsRoleStorage) FindByCustomerSegmentation(ctx context.Context, customerSegment string) (*imodel.CPSRoles, error) {
	log := local_util.LoggerFromCtx(ctx, m.logger)

	customerSegment = strings.TrimSpace(customerSegment)
	if customerSegment == "" {
		return nil, errors.New(localization.ErrorNoDataProvided.Code)
	}

	// Resolve superapp role for a customer segmentation code (CUSTOMER_SEGMENTATIONS.NAME) via sub-segments.
	const q = `
SELECT
  RAWTOHEX(cr.ID),
  cr.NAME,
  cr.LABEL,
  cr.ROLE_CODE,
  cr.DESCRIPTION,
  cr.IS_ENABLED,
  cr.IS_DELETED,
  cr.CREATED_AT,
  cr.LAST_MODIFIED_AT,
  cr.DELETED_AT
FROM CUSTOMER_SEGMENTATIONS cs
JOIN CUSTOMER_SUB_SEGMENTS css
  ON css.CUSTOMER_SEGMENTATION_ID = cs.ID AND css.IS_DELETED = 0 AND css.IS_ENABLED = 1
JOIN SUPERAPP_ROLES cr
  ON cr.ID = css.SUPERAPP_ROLE_ID AND cr.IS_DELETED = 0 AND cr.IS_ENABLED = 1
WHERE cs.IS_DELETED = 0
  AND cs.IS_ENABLED = 1
  AND UPPER(TRIM(cs.NAME)) = UPPER(TRIM(:1))
FETCH FIRST 1 ROWS ONLY`

	var (
		id                          string
		name, roleCode, label, desc sql.NullString
		enabledN, isDeletedN        int
		createdAt, updatedAt, delT  sql.NullTime
	)

	err := m.db.QueryRowContext(ctx, q, customerSegment).Scan(
		&id,
		&name,
		&label,
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
		log.Errorf("[CPSRolesStorage][FindByCustomerSegmentation] query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	enabled := enabledN == 1
	out := &imodel.CPSRoles{
		ID:          id,
		Name:        name.String,
		Lable:       label.String,
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

func (m *cpsRoleStorage) FindByCustomerSegmentationByID(ctx context.Context, customerSegment string) (*imodel.CPSRoles, error) {
	log := local_util.LoggerFromCtx(ctx, m.logger)

	customerSegment = strings.TrimSpace(customerSegment)
	if customerSegment == "" {
		return nil, errors.New(localization.ErrorNoDataProvided.Code)
	}

	const q = `
SELECT
  RAWTOHEX(ID),
  NAME,
  IS_ENABLED
FROM SUPERAPP_ROLES 
WHERE  ID = HEXTORAW(:1)`

	var (
		id                          string
		name, roleCode, label, desc sql.NullString
		enabledN, isDeletedN        int
		createdAt, updatedAt, delT  sql.NullTime
	)

	err := m.db.QueryRowContext(ctx, q, customerSegment).Scan(
		&id,
		&name,
		&enabledN,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		log.Errorf("[CPSRolesStorage][FindByCustomerSegmentationByID] query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	enabled := enabledN == 1
	out := &imodel.CPSRoles{
		ID:          id,
		Name:        name.String,
		Lable:       label.String,
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
	log := local_util.LoggerFromCtx(ctx, m.logger)

	// Oracle schema for access-list segmentation is not yet available in this service.
	log.Warnf("[CPSRolesStorage][EnableServiceAccess] not implemented for oracle; roleID=%s keys=%v", roleID, accessListKeys)
	return errors.New(localization.ErrorUnexpectedError.Code)
}

func (m *cpsRoleStorage) DisableServiceAccess(ctx context.Context, roleID string, accessListKeys []string) error {
	log := local_util.LoggerFromCtx(ctx, m.logger)

	// Oracle schema for access-list segmentation is not yet available in this service.
	log.Warnf("[CPSRolesStorage][DisableServiceAccess] not implemented for oracle; roleID=%s keys=%v", roleID, accessListKeys)
	return errors.New(localization.ErrorUnexpectedError.Code)
}

func (m *cpsRoleStorage) Delete(ctx context.Context, id string) error {
	log := local_util.LoggerFromCtx(ctx, m.logger)

	idHex, ok := normalizeRawHex32(id)
	if !ok {
		return errors.New(localization.ErrorInvalidID.Code)
	}

	const q = `
UPDATE SUPERAPP_ROLES
SET
  IS_DELETED = 1,
  LAST_MODIFIED_AT = SYSTIMESTAMP,
  DELETED_AT = SYSTIMESTAMP
WHERE ID = HEXTORAW(:1) AND IS_DELETED = 0`

	res, err := m.db.ExecContext(ctx, q, idHex)
	if err != nil {
		log.Errorf("[CPSRolesStorage][Delete] delete failed: %v", err)
		return local_util.HandleDBError(err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}
	return nil
}
