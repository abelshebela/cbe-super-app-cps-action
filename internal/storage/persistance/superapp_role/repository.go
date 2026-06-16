package superapprole

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	imodel "cbe-super-app-cps-action/internal/constants/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type superAppRoleStorage struct {
	db     *sql.DB
	logger utils.Logger
}

func NewSuperAppRoleRepository(db *sql.DB, logger utils.Logger) storage.SuperAppRoleRepository {
	return &superAppRoleStorage{db: db, logger: logger}
}

// buildWhere builds the WHERE clause and args for SEGMENTS queries.
func buildWhere(filterParam types.Filter) (string, []interface{}) {
	clauses := []string{"1=1"}
	var args []interface{}

	if s := strings.TrimSpace(filterParam.Search); s != "" {
		clauses = append(clauses,
			`(LOWER(SUPERAPP_ROLE) LIKE '%'||LOWER(:search)||'%'
			  OR LOWER(SUPERAPP_ROLE_LABEL) LIKE '%'||LOWER(:search)||'%'
			  OR LOWER(CUSTOMER_GROUP) LIKE '%'||LOWER(:search)||'%'
			  OR LOWER(CUSTOMER_SEGMENT) LIKE '%'||LOWER(:search)||'%')`)
		args = append(args, sql.Named("search", s))
	}

	if filterParam.Filters != nil {
		if v, ok := filterParam.Filters["is_enabled"]; ok {
			if b, ok2 := parseBoolFilter(v); ok2 {
				clauses = append(clauses, "IS_ENABLED = :is_enabled")
				args = append(args, sql.Named("is_enabled", boolToInt(b)))
			}
		}
	}

	return strings.Join(clauses, " AND "), args
}

func (r *superAppRoleStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.SuperAppRoleGroup], error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	limit := 50
	page := 1
	if filterParam.PerPage > 0 {
		limit = filterParam.PerPage
	}
	if filterParam.Page > 0 {
		page = filterParam.Page
	}
	offset := (page - 1) * limit

	where, args := buildWhere(filterParam)

	total, err := r.countDistinctRoles(ctx, where, args)
	if err != nil {
		log.Errorf("[SuperAppRoleRepo][FindAll] count err: %v", err)
		return nil, err
	}

	roleGroups, err := r.fetchRoleGroups(ctx, where, args, limit, offset)
	if err != nil {
		log.Errorf("[SuperAppRoleRepo][FindAll] role groups err: %v", err)
		return nil, err
	}

	if len(roleGroups) > 0 {
		if err := r.fillSegments(ctx, where, args, roleGroups); err != nil {
			log.Errorf("[SuperAppRoleRepo][FindAll] fill segments err: %v", err)
			return nil, err
		}
	}

	meta := local_util.BuildPaginationMeta(int64(total), page, limit)
	return &types.PaginatedResponse[[]imodel.SuperAppRoleGroup]{
		Data: roleGroups,
		Meta: meta,
	}, nil
}

func (r *superAppRoleStorage) countDistinctRoles(ctx context.Context, where string, args []interface{}) (int, error) {
	q := fmt.Sprintf(`SELECT COUNT(1) FROM (
    SELECT 1
    FROM SEGMENTS
    WHERE %s
    GROUP BY SUPERAPP_ROLE
) grouped_roles`, where)
	var total int
	if err := r.db.QueryRowContext(ctx, q, args...).Scan(&total); err != nil {
		return 0, local_util.HandleDBError(err)
	}
	return total, nil
}

func (r *superAppRoleStorage) fetchRoleGroups(ctx context.Context, where string, args []interface{}, limit, offset int) ([]imodel.SuperAppRoleGroup, error) {
	q := fmt.Sprintf(`
SELECT SUPERAPP_ROLE,
       MAX(SUPERAPP_ROLE_LABEL) KEEP (DENSE_RANK LAST ORDER BY LAST_MODIFIED_AT, CREATED_AT) AS SUPERAPP_ROLE_LABEL,
       MIN(IS_ENABLED) AS GROUP_ENABLED,
       MIN(CREATED_AT) AS FIRST_CREATED,
       MAX(LAST_MODIFIED_AT) AS LAST_TOUCHED
FROM SEGMENTS
WHERE %s
GROUP BY SUPERAPP_ROLE
ORDER BY SUPERAPP_ROLE
OFFSET :pg_offset ROWS FETCH NEXT :pg_limit ROWS ONLY`, where)

	listArgs := append(args,
		sql.Named("pg_offset", offset),
		sql.Named("pg_limit", limit),
	)

	rows, err := r.db.QueryContext(ctx, q, listArgs...)
	if err != nil {
		return nil, local_util.HandleDBError(err)
	}
	defer rows.Close()

	return collectRoleGroups(rows)
}

func collectRoleGroups(rows *sql.Rows) ([]imodel.SuperAppRoleGroup, error) {
	var result []imodel.SuperAppRoleGroup
	for rows.Next() {
		var role, label string
		var groupEnabled int
		var firstCreated, lastTouched sql.NullTime
		if err := rows.Scan(&role, &label, &groupEnabled, &firstCreated, &lastTouched); err != nil {
			return nil, local_util.HandleDBError(err)
		}
		g := imodel.SuperAppRoleGroup{
			SuperappRole:      role,
			SuperappRoleLabel: label,
			IsEnabled:         groupEnabled == 1,
			Segments:          []imodel.Segment{},
		}
		if firstCreated.Valid {
			g.CreatedAt = firstCreated.Time
		}
		if lastTouched.Valid {
			g.LastModifiedAt = lastTouched.Time
		}
		result = append(result, g)
	}
	return result, nil
}

func (r *superAppRoleStorage) fillSegments(ctx context.Context, where string, args []interface{}, groups []imodel.SuperAppRoleGroup) error {
	roleNames := make([]string, len(groups))
	for i, g := range groups {
		roleNames[i] = g.SuperappRole
	}

	placeholders := make([]string, len(roleNames))
	roleArgs := make([]interface{}, len(roleNames))
	for i, rn := range roleNames {
		key := fmt.Sprintf("rl%d", i)
		placeholders[i] = ":" + key
		roleArgs[i] = sql.Named(key, rn)
	}

	q := fmt.Sprintf(`
SELECT RAWTOHEX(ID), CUSTOMER_GROUP, CUSTOMER_GROUP_LABEL,
       CUSTOMER_SEGMENT, CUSTOMER_SEGMENT_LABEL,
       CUSTOMER_SUBSEGMENT, CUSTOMER_SUBSEGMENT_LABEL,
       SUPERAPP_ROLE, SUPERAPP_ROLE_LABEL,
       IS_ENABLED, CREATED_AT, LAST_MODIFIED_AT
FROM SEGMENTS
WHERE SUPERAPP_ROLE IN (%s)
  AND %s
ORDER BY SUPERAPP_ROLE, CREATED_AT`, strings.Join(placeholders, ", "), where)

	listArgs := append(roleArgs, args...)
	rows, err := r.db.QueryContext(ctx, q, listArgs...)
	if err != nil {
		return local_util.HandleDBError(err)
	}
	defer rows.Close()

	indexMap := make(map[string]int, len(groups))
	for i, g := range groups {
		indexMap[g.SuperappRole] = i
	}

	for rows.Next() {
		seg, err := scanSegment(rows)
		if err != nil {
			return err
		}
		if idx, ok := indexMap[seg.SuperappRole]; ok {
			groups[idx].Segments = append(groups[idx].Segments, seg)
		}
	}
	return nil
}

func scanSegment(rows *sql.Rows) (imodel.Segment, error) {
	var seg imodel.Segment
	var enabledN int
	var createdAt, lastModAt sql.NullTime
	err := rows.Scan(
		&seg.ID,
		&seg.CustomerGroup, &seg.CustomerGroupLabel,
		&seg.CustomerSegment, &seg.CustomerSegmentLabel,
		&seg.CustomerSubsegment, &seg.CustomerSubsegmentLabel,
		&seg.SuperappRole, &seg.SuperappRoleLabel,
		&enabledN, &createdAt, &lastModAt,
	)
	if err != nil {
		return seg, local_util.HandleDBError(err)
	}
	seg.ID = strings.ToLower(seg.ID)
	seg.IsEnabled = enabledN == 1
	if createdAt.Valid {
		seg.CreatedAt = createdAt.Time
	}
	if lastModAt.Valid {
		seg.LastModifiedAt = lastModAt.Time
	}
	return seg, nil
}

func (r *superAppRoleStorage) EnableByRole(ctx context.Context, superappRole string) error {
	return r.setEnabled(ctx, superappRole, 1)
}

func (r *superAppRoleStorage) DisableByRole(ctx context.Context, superappRole string) error {
	return r.setEnabled(ctx, superappRole, 0)
}

func (r *superAppRoleStorage) setEnabled(ctx context.Context, superappRole string, val int) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	role := strings.TrimSpace(superappRole)
	if role == "" {
		return errors.New(localization.ErrorNoDataProvided.Code)
	}
	const q = `UPDATE SEGMENTS SET IS_ENABLED = :1, LAST_MODIFIED_AT = :2 WHERE SUPERAPP_ROLE = :3`
	res, err := r.db.ExecContext(ctx, q, val, time.Now(), role)
	if err != nil {
		log.Errorf("[SuperAppRoleRepo][setEnabled=%d] err for role %s: %v", val, role, err)
		return local_util.HandleDBError(err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}
	return nil
}

func (r *superAppRoleStorage) DeleteByRole(ctx context.Context, superappRole string) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	role := strings.TrimSpace(superappRole)
	if role == "" {
		return errors.New(localization.ErrorNoDataProvided.Code)
	}
	const q = `DELETE FROM SEGMENTS WHERE SUPERAPP_ROLE = :1`
	res, err := r.db.ExecContext(ctx, q, role)
	if err != nil {
		log.Errorf("[SuperAppRoleRepo][DeleteByRole] err for role %s: %v", role, err)
		return local_util.HandleDBError(err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}
	return nil
}

func (r *superAppRoleStorage) RoleExists(ctx context.Context, superappRole string) (bool, error) {
	const q = `SELECT COUNT(1) FROM SEGMENTS WHERE SUPERAPP_ROLE = :1`
	var count int
	if err := r.db.QueryRowContext(ctx, q, superappRole).Scan(&count); err != nil {
		return false, local_util.HandleDBError(err)
	}
	return count > 0, nil
}

func (r *superAppRoleStorage) FindRoleBlockedAccessLists(ctx context.Context, superappRole string) ([]imodel.APPAccessList, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	const q = `SELECT RAWTOHEX(g.ACCESS_LIST_ID), a.NAME, a.SERVICE_KEY, g.IS_ENABLED
		FROM ACCESS_LIST_BY_SUPERAPP_ROLE g
		JOIN ACCESS_LISTS a ON g.ACCESS_LIST_ID = a.ID
		WHERE g.SUPERAPP_ROLE_ID = :1
		  AND g.IS_ENABLED = 1
		  AND g.IS_DELETED = 0`

	rows, err := r.db.QueryContext(ctx, q, superappRole)
	if err != nil {
		log.Errorf("[SuperAppRoleRepo][FindRoleBlocked] query err: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer rows.Close()
	return scanAccessLists(rows)
}

func (r *superAppRoleStorage) FindGloballyEnabledAccessLists(ctx context.Context) ([]imodel.APPAccessList, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	const q = `SELECT RAWTOHEX(ID), NAME, SERVICE_KEY, IS_ENABLED
		FROM ACCESS_LISTS
		WHERE IS_ENABLED = 1 AND IS_DELETED = 0
		ORDER BY NAME`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		log.Errorf("[SuperAppRoleRepo][FindGloballyEnabled] query err: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer rows.Close()
	return scanAccessLists(rows)
}

func (r *superAppRoleStorage) FindGloballyDisabledAccessLists(ctx context.Context) ([]imodel.APPAccessList, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	const q = `SELECT RAWTOHEX(ID), NAME, SERVICE_KEY, IS_ENABLED
		FROM ACCESS_LISTS
		WHERE IS_ENABLED = 0 AND IS_DELETED = 0
		ORDER BY NAME`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		log.Errorf("[SuperAppRoleRepo][FindGloballyDisabled] query err: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer rows.Close()
	return scanAccessLists(rows)
}

func (r *superAppRoleStorage) FindBlockedAccessListsByIDs(ctx context.Context, superappRole string, accessListIDs []string) ([]imodel.APPAccessList, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	if len(accessListIDs) == 0 {
		return nil, nil
	}

	placeholders := make([]string, len(accessListIDs))
	args := make([]interface{}, len(accessListIDs)+1)
	args[0] = superappRole
	for i, id := range accessListIDs {
		placeholders[i] = fmt.Sprintf("HEXTORAW(:%d)", i+2)
		args[i+1] = id
	}

	q := fmt.Sprintf(`SELECT RAWTOHEX(g.ACCESS_LIST_ID) AS ID, a.NAME , a.SERVICE_KEY, g.IS_ENABLED
		FROM ACCESS_LIST_BY_SUPERAPP_ROLE g
		JOIN ACCESS_LISTS a ON g.ACCESS_LIST_ID = a.ID
		WHERE g.SUPERAPP_ROLE_ID = :1
		  AND g.ACCESS_LIST_ID IN (%s)
		  AND g.IS_ENABLED = 1
		  AND g.IS_DELETED = 0`, strings.Join(placeholders, ", "))

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		log.Errorf("[SuperAppRoleRepo][FindBlockedByIDs] query err: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer rows.Close()
	return scanAccessLists(rows)
}
func (r *superAppRoleStorage) FindAccessListsByIDs(ctx context.Context, ids []string) ([]imodel.APPAccessList, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	if len(ids) == 0 {
		return nil, nil
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))

	for i, id := range ids {
		placeholders[i] = fmt.Sprintf(":%d", i+1)
		args[i] = strings.ToUpper(id)
	}

	q := fmt.Sprintf(`
	SELECT RAWTOHEX(ID), NAME, SERVICE_KEY, IS_ENABLED
	FROM ACCESS_LISTS
	WHERE RAWTOHEX(ID) IN (%s)
	AND IS_DELETED = 0
`, strings.Join(placeholders, ", "))

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		log.Errorf("[SuperAppRoleRepo][FindAccessListsByIDs] query err: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer rows.Close()

	return scanAccessLists(rows)
}

func scanAccessLists(rows *sql.Rows) ([]imodel.APPAccessList, error) {
	var result []imodel.APPAccessList
	for rows.Next() {
		var id, name, serviceKey string
		var isEnabled int
		if err := rows.Scan(&id, &name, &serviceKey, &isEnabled); err != nil {
			return nil, local_util.HandleDBError(err)
		}
		result = append(result, imodel.APPAccessList{
			ID:             id,
			Key:            serviceKey,
			AccessListName: name,
			Enabled:        isEnabled == 1,
		})
	}
	return result, nil
}

func (r *superAppRoleStorage) BulkDisableAccessLists(ctx context.Context, superappRole string, accessListIDs []string) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	if len(accessListIDs) == 0 {
		return nil
	}

	for _, alID := range accessListIDs {
		const q = `
		INSERT INTO ACCESS_LIST_BY_SUPERAPP_ROLE (
			ID,
			SUPERAPP_ROLE_ID,
			ACCESS_LIST_ID,
			IS_ENABLED,
			IS_DELETED,
			CREATED_AT,
			LAST_MODIFIED_AT,
			DELETED_AT
		) Values (
			SYS_GUID(), :1, HEXTORAW(:2), 1, 0, SYSDATE, SYSDATE, NULL
		)
	`

		if _, err := r.db.ExecContext(ctx, q, superappRole, alID); err != nil {
			log.Errorf(
				"[SuperAppRoleRepo][Insert] insert err for alID=%s: %v",
				alID,
				err,
			)
			return local_util.HandleDBError(err)
		}
	}
	return nil
}

func (r *superAppRoleStorage) BulkEnableAccessLists(ctx context.Context, superappRole string, accessListIDs []string) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	if len(accessListIDs) == 0 {
		return nil
	}

	placeholders := make([]string, len(accessListIDs))
	args := make([]interface{}, len(accessListIDs)+1)
	args[0] = superappRole
	for i, id := range accessListIDs {
		placeholders[i] = fmt.Sprintf("HEXTORAW(:%d)", i+2)
		args[i+1] = id
	}

	q := fmt.Sprintf(`DELETE FROM ACCESS_LIST_BY_SUPERAPP_ROLE
		WHERE SUPERAPP_ROLE_ID = :1
		  AND ACCESS_LIST_ID IN (%s)`, strings.Join(placeholders, ", "))

	if _, err := r.db.ExecContext(ctx, q, args...); err != nil {
		log.Errorf("[SuperAppRoleRepo][BulkEnable] delete err: %v", err)
		return local_util.HandleDBError(err)
	}
	return nil
}

func parseBoolFilter(v interface{}) (bool, bool) {
	switch t := v.(type) {
	case bool:
		return t, true
	case string:
		switch strings.ToLower(strings.TrimSpace(t)) {
		case "true", "1", "yes":
			return true, true
		case "false", "0", "no":
			return false, true
		}
	}
	return false, false
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
