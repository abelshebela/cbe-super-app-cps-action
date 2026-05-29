package customergroup

import (
	"context"
	"crypto/md5"
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

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/superapp_mapper/checksum"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type customerGroupStorage struct {
	cfg    *config.VaultConfig
	db     *sql.DB
	logger utils.Logger
}

func NewCustomerGroupRepository(cfg *config.VaultConfig, db *sql.DB, logger utils.Logger) storage.CustomerGroupRepository {
	return &customerGroupStorage{
		cfg:    cfg,
		db:     db,
		logger: logger,
	}
}

func computeCheckSum(group, segment, subsegment string) string {
	raw := fmt.Sprintf("%s:%s:%s", group, segment, subsegment)
	sum := md5.Sum([]byte(raw))
	return fmt.Sprintf("%x", sum)
}

func normalizeSegmentHex(id string) (string, bool) {
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

func boolToNumber(v bool) int {
	if v {
		return 1
	}
	return 0
}

func (r *customerGroupStorage) DuplicateCheck(ctx context.Context, action, id, group, segment, subsegment string) (bool, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)
	var dupQ string
	var segments, subsegments string
	if segment == "" {
		segments = "*"
	} else {
		segment = strings.TrimSpace(segment)
	}
	if subsegment == "" {
		subsegments = "*"
	} else {
		subsegment = strings.TrimSpace(subsegment)
	}
	data := fmt.Sprintf("%s:%s:%s", strings.TrimSpace(group), segments, subsegments)
	checkSum := checksum.Checksum(data)

	var count int

	// dupQ = `SELECT COUNT(*) FROM SEGMENTS WHERE CHECK_SUM = :1`

	if action == "update" {
		dupQ = `SELECT COUNT(*) FROM SEGMENTS WHERE CHECK_SUM = :1 AND ID != HEXTORAW(:2)`
		idHex, ok := normalizeSegmentHex(id)
		if !ok {
			return false, errors.New(localization.ErrorInvalidID.Code)
		}
		if err := r.db.QueryRowContext(ctx, dupQ, checkSum, idHex).Scan(&count); err != nil {
			log.Errorf("[CustomerGroup][DuplicateCheck] duplicate check failed: %v", err)
			return false, local_util.HandleDBError(err)
		}
	} else {
		dupQ = `SELECT COUNT(*) FROM SEGMENTS WHERE CHECK_SUM = :1`
		if err := r.db.QueryRowContext(ctx, dupQ, checkSum).Scan(&count); err != nil {
			log.Errorf("[CustomerGroup][DuplicateCheck] duplicate check failed: %v", err)
			return false, local_util.HandleDBError(err)
		}
	}

	return count > 0, nil
}

func (r *customerGroupStorage) Create(ctx context.Context, seg *imodel.Segment) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)
	log.Infof("[CustomerGroup][Create] creating segment: %+v", seg)

	if seg == nil {
		return errors.New(localization.ErrorNoDataProvided.Code)
	}
	var segment, subsegment string
	if seg.CustomerSegment == "" {
		segment = "*"
	} else {
		segment = strings.TrimSpace(seg.CustomerSegment)
	}
	if seg.CustomerSubsegment == "" {
		subsegment = "*"
	} else {
		subsegment = strings.TrimSpace(seg.CustomerSubsegment)
	}
	data := fmt.Sprintf("%s:%s:%s", strings.TrimSpace(seg.CustomerGroup), segment, subsegment)
	checkSum := checksum.Checksum(data)

	var count int
	const dupQ = `SELECT COUNT(*) FROM SEGMENTS WHERE CHECK_SUM = :1`
	if err := r.db.QueryRowContext(ctx, dupQ, checkSum).Scan(&count); err != nil {
		log.Errorf("[CustomerGroup][Create] duplicate check failed: %v", err)
		return local_util.HandleDBError(err)
	}

	if count > 0 {
		log.Warnf("[CustomerGroup][Create] duplicate combination detected checksum=%s", checkSum)
		return errors.New(localization.ErrorCustomerGroupAlreadyExists.Code)
	}

	now := time.Now()
	var newID string
	const insertQ = `
INSERT INTO SEGMENTS (
	CUSTOMER_GROUP,
	CUSTOMER_GROUP_LABEL,
	CUSTOMER_SEGMENT,
	CUSTOMER_SEGMENT_LABEL,
	CUSTOMER_SUBSEGMENT,
	CUSTOMER_SUBSEGMENT_LABEL,
	SUPERAPP_ROLE,
	SUPERAPP_ROLE_LABEL,
	CHECK_SUM,
	IS_ENABLED,
	CREATED_AT,
	LAST_MODIFIED_AT
) VALUES (
	:1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11, :12
) RETURNING RAWTOHEX(ID) INTO :13`

	if _, err := r.db.ExecContext(ctx, insertQ,
		seg.CustomerGroup,
		seg.CustomerGroupLabel,
		seg.CustomerSegment,
		seg.CustomerSegmentLabel,
		seg.CustomerSubsegment,
		seg.CustomerSubsegmentLabel,
		seg.SuperappRole,
		seg.SuperappRoleLabel,
		checkSum,
		boolToNumber(seg.IsEnabled),
		now,
		now,
		sql.Out{Dest: &newID},
	); err != nil {
		log.Errorf("[CustomerGroup][Create] insert failed: %v", err)
		return local_util.HandleDBError(err)
	}

	seg.ID = strings.ToLower(newID)
	seg.CheckSum = checkSum
	seg.CreatedAt = now
	seg.LastModifiedAt = now

	log.Infof("[CustomerGroup][Create] created with id=%s", seg.ID)
	return nil
}

func (r *customerGroupStorage) Update(ctx context.Context, id string, seg *imodel.Segment) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	idHex, ok := normalizeSegmentHex(id)
	if !ok {
		return errors.New(localization.ErrorInvalidID.Code)
	}
	if seg == nil {
		return errors.New(localization.ErrorNoDataProvided.Code)
	}

	var segment, subsegment string
	if seg.CustomerSegment != "" {
		segment = seg.CustomerSegment
	} else {
		segment = "*"
	}
	if seg.CustomerSubsegment != "" {
		subsegment = seg.CustomerSubsegment
	} else {
		subsegment = "*"
	}
	data := fmt.Sprintf("%s:%s:%s", seg.CustomerGroup, segment, subsegment)
	checkSum := checksum.Checksum(data)

	var count int
	const dupQ = `SELECT COUNT(*) FROM SEGMENTS WHERE CHECK_SUM = :1 AND ID != HEXTORAW(:2)`
	if err := r.db.QueryRowContext(ctx, dupQ, checkSum, idHex).Scan(&count); err != nil {
		log.Errorf("[CustomerGroup][Update] duplicate check failed: %v", err)
		return local_util.HandleDBError(err)
	}
	if count > 0 {
		log.Warnf("[CustomerGroup][Update] duplicate combination detected checksum=%s", checkSum)
		return errors.New(localization.ErrorCustomerGroupAlreadyExists.Code)
	}

	now := time.Now()
	const updateQ = `
UPDATE SEGMENTS SET
	CUSTOMER_GROUP           = :1,
	CUSTOMER_GROUP_LABEL     = :2,
	CUSTOMER_SEGMENT         = :3,
	CUSTOMER_SEGMENT_LABEL   = :4,
	CUSTOMER_SUBSEGMENT      = :5,
	CUSTOMER_SUBSEGMENT_LABEL = :6,
	SUPERAPP_ROLE            = :7,
	SUPERAPP_ROLE_LABEL      = :8,
	CHECK_SUM                = :9,
	LAST_MODIFIED_AT         = :10
WHERE ID = HEXTORAW(:11)`

	res, err := r.db.ExecContext(ctx, updateQ,
		seg.CustomerGroup,
		seg.CustomerGroupLabel,
		seg.CustomerSegment,
		seg.CustomerSegmentLabel,
		seg.CustomerSubsegment,
		seg.CustomerSubsegmentLabel,
		seg.SuperappRole,
		seg.SuperappRoleLabel,
		checkSum,
		now,
		idHex,
	)
	if err != nil {
		log.Errorf("[CustomerGroup][Update] update failed: %v", err)
		return local_util.HandleDBError(err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}

	seg.CheckSum = checkSum
	seg.LastModifiedAt = now
	return nil
}

func (r *customerGroupStorage) Delete(ctx context.Context, id string) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	idHex, ok := normalizeSegmentHex(id)
	if !ok {
		return errors.New(localization.ErrorInvalidID.Code)
	}

	const q = `DELETE FROM SEGMENTS WHERE ID = HEXTORAW(:1)`
	res, err := r.db.ExecContext(ctx, q, idHex)
	if err != nil {
		log.Errorf("[CustomerGroup][Delete] delete failed: %v", err)
		return local_util.HandleDBError(err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}
	return nil
}

func (r *customerGroupStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	idHex, ok := normalizeSegmentHex(id)
	if !ok {
		return errors.New(localization.ErrorInvalidID.Code)
	}

	const q = `UPDATE SEGMENTS SET IS_ENABLED = :1, LAST_MODIFIED_AT = :2 WHERE ID = HEXTORAW(:3)`
	res, err := r.db.ExecContext(ctx, q, boolToNumber(enable), time.Now(), idHex)
	if err != nil {
		log.Errorf("[CustomerGroup][EnableOrDisable] update failed: %v", err)
		return local_util.HandleDBError(err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}
	return nil
}

func (r *customerGroupStorage) FindByID(ctx context.Context, id string) (*imodel.Segment, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	idHex, ok := normalizeSegmentHex(id)
	if !ok {
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	const q = `
SELECT
	RAWTOHEX(ID),
	CUSTOMER_GROUP,
	CUSTOMER_GROUP_LABEL,
	CUSTOMER_SEGMENT,
	CUSTOMER_SEGMENT_LABEL,
	CUSTOMER_SUBSEGMENT,
	CUSTOMER_SUBSEGMENT_LABEL,
	SUPERAPP_ROLE,
	SUPERAPP_ROLE_LABEL,
	IS_ENABLED,
	CREATED_AT,
	LAST_MODIFIED_AT
FROM SEGMENTS
WHERE ID = HEXTORAW(:1)`

	row := r.db.QueryRowContext(ctx, q, idHex)
	seg, err := scanSegment(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		log.Errorf("[CustomerGroup][FindByID] query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return seg, nil
}

func buildSegmentWhere(filterParam types.Filter) (string, []interface{}) {
	clauses := []string{"1=1"}
	var args []interface{}

	if search := strings.TrimSpace(filterParam.Search); search != "" {
		clauses = append(clauses, `(
			LOWER(CUSTOMER_GROUP) LIKE '%'||LOWER(:search)||'%'
			OR LOWER(CUSTOMER_GROUP_LABEL) LIKE '%'||LOWER(:search)||'%'
			OR LOWER(CUSTOMER_SEGMENT) LIKE '%'||LOWER(:search)||'%'
			OR LOWER(CUSTOMER_SEGMENT_LABEL) LIKE '%'||LOWER(:search)||'%'
			OR LOWER(CUSTOMER_SUBSEGMENT) LIKE '%'||LOWER(:search)||'%'
			OR LOWER(CUSTOMER_SUBSEGMENT_LABEL) LIKE '%'||LOWER(:search)||'%'
			OR LOWER(SUPERAPP_ROLE) LIKE '%'||LOWER(:search)||'%'
			OR LOWER(SUPERAPP_ROLE_LABEL) LIKE '%'||LOWER(:search)||'%'
		)`)
		args = append(args, sql.Named("search", search))
	}

	clauses, args = applySegmentFilters(filterParam.Filters, clauses, args)
	return strings.Join(clauses, " AND "), args
}

func applyStringEqFilter(filters map[string]interface{}, key, clause, namedParam string, clauses []string, args []interface{}) ([]string, []interface{}) {
	v, ok := filters[key]
	if !ok {
		return clauses, args
	}
	s, ok := v.(string)
	if !ok || s == "" {
		return clauses, args
	}
	return append(clauses, clause), append(args, sql.Named(namedParam, s))
}

func applySegmentFilters(filters map[string]interface{}, clauses []string, args []interface{}) ([]string, []interface{}) {
	if filters == nil {
		return clauses, args
	}
	if v, ok := filters["is_enabled"]; ok {
		if b, ok2 := parseBoolFilter(v); ok2 {
			clauses = append(clauses, "IS_ENABLED = :is_enabled")
			args = append(args, sql.Named("is_enabled", boolToNumber(b)))
		}
	}
	clauses, args = applyStringEqFilter(filters, "customer_group", "UPPER(CUSTOMER_GROUP) = UPPER(:cg)", "cg", clauses, args)
	clauses, args = applyStringEqFilter(filters, "customer_segment", "UPPER(CUSTOMER_SEGMENT) = UPPER(:cs)", "cs", clauses, args)
	clauses, args = applyStringEqFilter(filters, "superapp_role", "UPPER(SUPERAPP_ROLE) = UPPER(:sar)", "sar", clauses, args)
	return clauses, args
}

func collectSegmentRows(rows *sql.Rows) ([]imodel.Segment, error) {
	result := make([]imodel.Segment, 0)
	for rows.Next() {
		seg, err := scanSegmentRow(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, *seg)
	}
	return result, nil
}

func (r *customerGroupStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.Segment], error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	limit := int64(50)
	page := int64(1)
	if filterParam.PerPage > 0 {
		limit = int64(filterParam.PerPage)
	}
	if filterParam.Page > 0 {
		page = int64(filterParam.Page)
	}
	offset := (page - 1) * limit

	where, args := buildSegmentWhere(filterParam)

	var total int64
	countQ := fmt.Sprintf(`SELECT COUNT(*) FROM SEGMENTS WHERE %s`, where)
	if err := r.db.QueryRowContext(ctx, countQ, args...).Scan(&total); err != nil {
		log.Errorf("[CustomerGroup][FindAll] count failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	listQ := fmt.Sprintf(`
SELECT
	RAWTOHEX(ID),
	CUSTOMER_GROUP,
	CUSTOMER_GROUP_LABEL,
	CUSTOMER_SEGMENT,
	CUSTOMER_SEGMENT_LABEL,
	CUSTOMER_SUBSEGMENT,
	CUSTOMER_SUBSEGMENT_LABEL,
	SUPERAPP_ROLE,
	SUPERAPP_ROLE_LABEL,
	IS_ENABLED,
	CREATED_AT,
	LAST_MODIFIED_AT
FROM SEGMENTS
WHERE %s
ORDER BY CREATED_AT DESC
OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`, where)

	listArgs := append(args, sql.Named("offset", offset), sql.Named("limit", limit))

	rows, err := r.db.QueryContext(ctx, listQ, listArgs...)
	if err != nil {
		log.Errorf("[CustomerGroup][FindAll] list query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer rows.Close()

	result, err := collectSegmentRows(rows)
	if err != nil {
		log.Errorf("[CustomerGroup][FindAll] scan failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	meta := local_util.BuildPaginationMeta(total, int(page), int(limit))
	return &types.PaginatedResponse[[]imodel.Segment]{
		Data: result,
		Meta: meta,
	}, nil
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func scanSegment(row rowScanner) (*imodel.Segment, error) {
	var (
		id, cg, cgl, cs, csl, csu, csul, sar, sarl string
		enabledN                                   int
		createdAt, lastModifiedAt                  sql.NullTime
	)
	if err := row.Scan(&id, &cg, &cgl, &cs, &csl, &csu, &csul, &sar, &sarl, &enabledN, &createdAt, &lastModifiedAt); err != nil {
		return nil, err
	}
	seg := &imodel.Segment{
		ID:                      strings.ToLower(id),
		CustomerGroup:           cg,
		CustomerGroupLabel:      cgl,
		CustomerSegment:         cs,
		CustomerSegmentLabel:    csl,
		CustomerSubsegment:      csu,
		CustomerSubsegmentLabel: csul,
		SuperappRole:            sar,
		SuperappRoleLabel:       sarl,
		IsEnabled:               enabledN == 1,
	}
	if createdAt.Valid {
		seg.CreatedAt = createdAt.Time
	}
	if lastModifiedAt.Valid {
		seg.LastModifiedAt = lastModifiedAt.Time
	}
	return seg, nil
}

func scanSegmentRow(rows *sql.Rows) (*imodel.Segment, error) {
	var (
		id, cg, cgl, cs, csl, csu, csul, sar, sarl string
		enabledN                                   int
		createdAt, lastModifiedAt                  sql.NullTime
	)
	if err := rows.Scan(&id, &cg, &cgl, &cs, &csl, &csu, &csul, &sar, &sarl, &enabledN, &createdAt, &lastModifiedAt); err != nil {
		return nil, err
	}
	seg := &imodel.Segment{
		ID:                      strings.ToLower(id),
		CustomerGroup:           cg,
		CustomerGroupLabel:      cgl,
		CustomerSegment:         cs,
		CustomerSegmentLabel:    csl,
		CustomerSubsegment:      csu,
		CustomerSubsegmentLabel: csul,
		SuperappRole:            sar,
		SuperappRoleLabel:       sarl,
		IsEnabled:               enabledN == 1,
	}
	if createdAt.Valid {
		seg.CreatedAt = createdAt.Time
	}
	if lastModifiedAt.Valid {
		seg.LastModifiedAt = lastModifiedAt.Time
	}
	return seg, nil
}

func parseBoolFilter(v interface{}) (bool, bool) {
	switch t := v.(type) {
	case bool:
		return t, true
	case int:
		return t != 0, true
	case float64:
		return t != 0, true
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
