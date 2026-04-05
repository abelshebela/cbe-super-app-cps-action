package customersegmentaion

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/kafka"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	imodel "cbe-super-app-cps-action/internal/constants/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type customerStorage struct {
	cfg           *config.VaultConfig
	db            *sql.DB
	kafkaProducer kafka.ClientOrchestrationProducer
	logger        utils.Logger
}

func NewCustomerSegmentationRepository(cfg *config.VaultConfig, db *sql.DB, kafkaProducer kafka.ClientOrchestrationProducer, logger utils.Logger) storage.CustomerSegmentationRepository {
	return &customerStorage{
		cfg:           cfg,
		db:            db,
		kafkaProducer: kafkaProducer,
		logger:        logger,
	}
}

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

func (r *customerStorage) assertSuperAppRoleExistsTx(ctx context.Context, tx *sql.Tx, roleHex string) error {
	const q = `
SELECT 1
FROM SUPERAPP_ROLE
WHERE ID = HEXTORAW(:1) AND IS_DELETED = 0 AND ENABLED = 1`
	var one int
	err := tx.QueryRowContext(ctx, q, roleHex).Scan(&one)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New(localization.ErrorResourceNotFound.Code)
		}
		r.logger.Errorf("[CustomerSegmentation][assertSuperAppRoleExistsTx] query failed: %v", err)
		return local_util.HandleDBError(err)
	}
	return nil
}

func (r *customerStorage) findOrCreateCustomerGroupTx(ctx context.Context, tx *sql.Tx, name string, now time.Time, isEnabled bool) (string, error) {
	const findQ = `
SELECT RAWTOHEX(ID)
FROM CUSTOMER_GROUPS
WHERE UPPER(TRIM(NAME)) = UPPER(TRIM(:1)) AND IS_DELETED = 0`
	var existing string
	err := tx.QueryRowContext(ctx, findQ, name).Scan(&existing)
	if err == nil {
		return strings.ToLower(existing), nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		r.logger.Errorf("[CustomerSegmentation][findOrCreateCustomerGroupTx] find failed: %v", err)
		return "", local_util.HandleDBError(err)
	}

	var newID string
	const insertQ = `
INSERT INTO CUSTOMER_GROUPS (
  NAME,
  IS_ENABLED,
  IS_DELETED,
  CREATED_AT,
  LAST_MODIFIED_AT
)
VALUES (
  UPPER(TRIM(:1)),
  :2,
  0,
  :3,
  :4
)
RETURNING RAWTOHEX(ID) INTO :5`
	if _, err := tx.ExecContext(ctx, insertQ,
		name,
		boolToOracleNumber(isEnabled),
		now,
		now,
		sql.Out{Dest: &newID},
	); err != nil {
		r.logger.Errorf("[CustomerSegmentation][findOrCreateCustomerGroupTx] insert failed: %v", err)
		return "", local_util.HandleDBError(err)
	}
	return strings.ToLower(newID), nil
}

func (r *customerStorage) findOrCreateSegmentationTx(ctx context.Context, tx *sql.Tx, groupHex, segName string, now time.Time, isEnabled bool) (string, error) {
	const findQ = `
SELECT RAWTOHEX(ID)
FROM CUSTOMER_SEGMENTATIONS
WHERE CUSTOMER_GROUPS_ID = HEXTORAW(:1)
  AND UPPER(TRIM(NAME)) = UPPER(TRIM(:2))
  AND IS_DELETED = 0`
	var existing string
	err := tx.QueryRowContext(ctx, findQ, groupHex, segName).Scan(&existing)
	if err == nil {
		return strings.ToLower(existing), nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		r.logger.Errorf("[CustomerSegmentation][findOrCreateSegmentationTx] find failed: %v", err)
		return "", local_util.HandleDBError(err)
	}

	var newID string
	const insertQ = `
INSERT INTO CUSTOMER_SEGMENTATIONS (
  NAME,
  CUSTOMER_GROUPS_ID,
  IS_ENABLED,
  IS_DELETED,
  CREATED_AT,
  LAST_MODIFIED_AT
)
VALUES (
  UPPER(TRIM(:1)),
  HEXTORAW(:2),
  :3,
  0,
  :4,
  :5
)
RETURNING RAWTOHEX(ID) INTO :6`
	if _, err := tx.ExecContext(ctx, insertQ,
		segName,
		groupHex,
		boolToOracleNumber(isEnabled),
		now,
		now,
		sql.Out{Dest: &newID},
	); err != nil {
		r.logger.Errorf("[CustomerSegmentation][findOrCreateSegmentationTx] insert failed: %v", err)
		return "", local_util.HandleDBError(err)
	}
	return strings.ToLower(newID), nil
}

func (r *customerStorage) resolveSuperAppRoleIDTx(ctx context.Context, tx *sql.Tx, idHex string) (string, error) {
	const roleFromSegQ = `
SELECT RAWTOHEX(css.SUPERAPP_ROLE_ID)
FROM CUSTOMER_SUB_SEGMENTS css
WHERE css.CUSTOMER_SEGMENTATIONS_ID = HEXTORAW(:1)
  AND css.IS_DELETED = 0
  AND css.IS_ENABLED = 1
FETCH FIRST 1 ROWS ONLY`
	var roleHex string
	err := tx.QueryRowContext(ctx, roleFromSegQ, idHex).Scan(&roleHex)
	if err == nil {
		return strings.ToLower(roleHex), nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		r.logger.Errorf("[CustomerSegmentation][resolveSuperAppRoleIDTx] role-from-seg query failed: %v", err)
		return "", local_util.HandleDBError(err)
	}

	if err := r.assertSuperAppRoleExistsTx(ctx, tx, idHex); err != nil {
		return "", err
	}
	return idHex, nil
}

func (r *customerStorage) Create(ctx context.Context, seg *imodel.CustomerSegmentation) error {
	if seg == nil {
		return errors.New(localization.ErrorNoDataProvided.Code)
	}
	if len(seg.CustomerSegments) == 0 {
		return errors.New(localization.ErrorNoDataProvided.Code)
	}

	now := time.Now()
	if seg.CreatedAt.IsZero() {
		seg.CreatedAt = now
	}
	if seg.UpdatedAt.IsZero() {
		seg.UpdatedAt = now
	}
	seg.IsDeleted = false

	roleIDHex, ok := normalizeRawHex32(seg.CustomerRole.ID)
	if !ok {
		return errors.New(localization.ErrorInvalidID.Code)
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Errorf("[CustomerSegmentation][Create] begin tx failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = tx.Rollback() }()

	if err := r.assertSuperAppRoleExistsTx(ctx, tx, roleIDHex); err != nil {
		return err
	}

	const insertSubQ = `
INSERT INTO CUSTOMER_SUB_SEGMENTS (
  NAME,
  CUSTOMER_SEGMENTATIONS_ID,
  SUPERAPP_ROLE_ID,
  IS_ENABLED,
  IS_DELETED,
  CREATED_AT,
  LAST_MODIFIED_AT
)
VALUES (
  UPPER(TRIM(:1)),
  HEXTORAW(:2),
  HEXTORAW(:3),
  :4,
  0,
  :5,
  :6
)`

	for _, sub := range seg.CustomerSegments {
		groupHex, err := r.findOrCreateCustomerGroupTx(ctx, tx, sub.CustomerGroup, now, seg.IsEnabled)
		if err != nil {
			return err
		}
		segHex, err := r.findOrCreateSegmentationTx(ctx, tx, groupHex, sub.CustomerSegment, now, seg.IsEnabled)
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, insertSubQ,
			sub.CustomerSubSegment,
			segHex,
			roleIDHex,
			boolToOracleNumber(seg.IsEnabled),
			now,
			now,
		); err != nil {
			r.logger.Errorf("[CustomerSegmentation][Create] insert sub segment failed: %v", err)
			return local_util.HandleDBError(err)
		}
	}

	if err := tx.Commit(); err != nil {
		r.logger.Errorf("[CustomerSegmentation][Create] commit failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	_ = r.kafkaProducer.PublishMessage(
		ctx,
		seg,
		string(constants.ClientOrchestrationServicesTopic),
		string(constants.CustomerSegmentationUpdatedTopic),
		"customer segmentation created",
	)

	return nil
}

func (r *customerStorage) Update(ctx context.Context, id string, seg *imodel.CustomerSegmentation) error {
	idHex, ok := normalizeRawHex32(id)
	if !ok {
		return errors.New(localization.ErrorInvalidID.Code)
	}
	if seg == nil {
		return errors.New(localization.ErrorNoDataProvided.Code)
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Errorf("[CustomerSegmentation][Update] begin tx failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = tx.Rollback() }()

	roleIDHex, err := r.resolveSuperAppRoleIDTx(ctx, tx, idHex)
	if err != nil {
		return err
	}

	const touchRoleQ = `
UPDATE SUPERAPP_ROLE
SET LAST_MODIFIED_AT = SYSTIMESTAMP
WHERE ID = HEXTORAW(:1) AND IS_DELETED = 0 AND ENABLED = 1`
	res, err := tx.ExecContext(ctx, touchRoleQ, roleIDHex)
	if err != nil {
		r.logger.Errorf("[CustomerSegmentation][Update] touch role failed: %v", err)
		return local_util.HandleDBError(err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}

	const deleteSubQ = `DELETE FROM CUSTOMER_SUB_SEGMENTS WHERE SUPERAPP_ROLE_ID = HEXTORAW(:1)`
	if _, err := tx.ExecContext(ctx, deleteSubQ, roleIDHex); err != nil {
		r.logger.Errorf("[CustomerSegmentation][Update] delete sub segments failed: %v", err)
		return local_util.HandleDBError(err)
	}

	if len(seg.CustomerSegments) == 0 {
		if err := tx.Commit(); err != nil {
			r.logger.Errorf("[CustomerSegmentation][Update] commit failed: %v", err)
			return errors.New(localization.ErrorUnexpectedError.Code)
		}
		_ = r.kafkaProducer.PublishMessage(
			ctx,
			seg,
			string(constants.ClientOrchestrationServicesTopic),
			string(constants.CustomerSegmentationUpdatedTopic),
			"customer segmentation updated",
		)
		return nil
	}

	now := time.Now()
	const insertSubQ = `
INSERT INTO CUSTOMER_SUB_SEGMENTS (
  NAME,
  CUSTOMER_SEGMENTATIONS_ID,
  SUPERAPP_ROLE_ID,
  IS_ENABLED,
  IS_DELETED,
  CREATED_AT,
  LAST_MODIFIED_AT
)
VALUES (
  UPPER(TRIM(:1)),
  HEXTORAW(:2),
  HEXTORAW(:3),
  :4,
  0,
  :5,
  :6
)`

	for _, sub := range seg.CustomerSegments {
		groupHex, err := r.findOrCreateCustomerGroupTx(ctx, tx, sub.CustomerGroup, now, seg.IsEnabled)
		if err != nil {
			return err
		}
		segHex, err := r.findOrCreateSegmentationTx(ctx, tx, groupHex, sub.CustomerSegment, now, seg.IsEnabled)
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, insertSubQ,
			sub.CustomerSubSegment,
			segHex,
			roleIDHex,
			boolToOracleNumber(seg.IsEnabled),
			now,
			now,
		); err != nil {
			r.logger.Errorf("[CustomerSegmentation][Update] insert sub segment failed: %v", err)
			return local_util.HandleDBError(err)
		}
	}

	if err := tx.Commit(); err != nil {
		r.logger.Errorf("[CustomerSegmentation][Update] commit failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	_ = r.kafkaProducer.PublishMessage(
		ctx,
		seg,
		string(constants.ClientOrchestrationServicesTopic),
		string(constants.CustomerSegmentationUpdatedTopic),
		"customer segmentation updated",
	)

	return nil
}

func (r *customerStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	idHex, ok := normalizeRawHex32(id)
	if !ok {
		return errors.New(localization.ErrorInvalidID.Code)
	}

	const q = `
UPDATE CUSTOMER_SEGMENTATIONS
SET
  IS_ENABLED  = :1,
  LAST_MODIFIED_AT  = SYSTIMESTAMP
WHERE ID = HEXTORAW(:2) AND IS_DELETED = 0`

	res, err := r.db.ExecContext(ctx, q, boolToOracleNumber(enable), idHex)
	if err != nil {
		r.logger.Errorf("[CustomerSegmentation][EnableOrDisable] update failed: %v", err)
		return local_util.HandleDBError(err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}
	return nil
}

func (r *customerStorage) Delete(ctx context.Context, id string) error {
	idHex, ok := normalizeRawHex32(id)
	if !ok {
		return errors.New(localization.ErrorInvalidID.Code)
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Errorf("[CustomerSegmentation][Delete] begin tx failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = tx.Rollback() }()

	const softSubQ = `
UPDATE CUSTOMER_SUB_SEGMENTS
SET IS_DELETED = 1, LAST_MODIFIED_AT = SYSTIMESTAMP
WHERE CUSTOMER_SEGMENTATIONS_ID = HEXTORAW(:1) AND IS_DELETED = 0`
	if _, err := tx.ExecContext(ctx, softSubQ, idHex); err != nil {
		r.logger.Errorf("[CustomerSegmentation][Delete] soft-delete sub segments failed: %v", err)
		return local_util.HandleDBError(err)
	}

	const q = `
UPDATE CUSTOMER_SEGMENTATIONS
SET
  IS_DELETED = 1,
  LAST_MODIFIED_AT = SYSTIMESTAMP
WHERE ID = HEXTORAW(:1) AND IS_DELETED = 0`

	res, err := tx.ExecContext(ctx, q, idHex)
	if err != nil {
		r.logger.Errorf("[CustomerSegmentation][Delete] delete failed: %v", err)
		return local_util.HandleDBError(err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}

	if err := tx.Commit(); err != nil {
		r.logger.Errorf("[CustomerSegmentation][Delete] commit failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (r *customerStorage) fillAggregateBySuperAppRoleID(ctx context.Context, roleHex, topLevelID string) (*imodel.CustomerSegmentation, error) {
	const metaQ = `
SELECT
  RAWTOHEX(sar.ID),
  sar.NAME,
  sar.ENABLED,
  sar.IS_DELETED,
  MIN(cs.IS_ENABLED),
  MIN(css.IS_ENABLED),
  MIN(css.CREATED_AT),
  MAX(GREATEST(cs.LAST_MODIFIED_AT, css.LAST_MODIFIED_AT))
FROM CUSTOMER_SUB_SEGMENTS css
JOIN SUPERAPP_ROLE sar ON sar.ID = css.SUPERAPP_ROLE_ID AND sar.IS_DELETED = 0 AND sar.ENABLED = 1
JOIN CUSTOMER_SEGMENTATIONS cs ON cs.ID = css.CUSTOMER_SEGMENTATIONS_ID AND cs.IS_DELETED = 0
JOIN CUSTOMER_GROUPS cg ON cg.ID = cs.CUSTOMER_GROUPS_ID AND cg.IS_DELETED = 0 AND cg.IS_ENABLED = 1
WHERE css.SUPERAPP_ROLE_ID = HEXTORAW(:1)
  AND css.IS_DELETED = 0
  AND css.IS_ENABLED = 1
GROUP BY sar.ID, sar.NAME, sar.ENABLED, sar.IS_DELETED`

	var (
		roleID, roleName        string
		sarEn, sarDel           int
		minCsEn, minCssEn       int
		firstCreated, lastTouch sql.NullTime
	)
	err := r.db.QueryRowContext(ctx, metaQ, roleHex).Scan(
		&roleID,
		&roleName,
		&sarEn,
		&sarDel,
		&minCsEn,
		&minCssEn,
		&firstCreated,
		&lastTouch,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		r.logger.Errorf("[CustomerSegmentation][fillAggregateBySuperAppRoleID] meta query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	seg := imodel.CustomerSegmentation{
		ID: strings.ToUpper(strings.TrimSpace(topLevelID)),
		CustomerRole: imodel.CustomerRoleInfo{
			ID:   roleID,
			Name: roleName,
		},
		IsEnabled: sarEn == 1 && minCsEn == 1 && minCssEn == 1,
		IsDeleted: sarDel == 1,
	}
	if firstCreated.Valid {
		seg.CreatedAt = firstCreated.Time
	}
	if lastTouch.Valid {
		seg.UpdatedAt = lastTouch.Time
	}

	const rowsQ = `
SELECT
  RAWTOHEX(css.ID),
  cg.NAME,
  cs.NAME,
  css.NAME
FROM CUSTOMER_SUB_SEGMENTS css
JOIN CUSTOMER_SEGMENTATIONS cs ON cs.ID = css.CUSTOMER_SEGMENTATIONS_ID AND cs.IS_DELETED = 0
JOIN CUSTOMER_GROUPS cg ON cg.ID = cs.CUSTOMER_GROUPS_ID AND cg.IS_DELETED = 0 AND cg.IS_ENABLED = 1
JOIN SUPERAPP_ROLE sar ON sar.ID = css.SUPERAPP_ROLE_ID AND sar.IS_DELETED = 0 AND sar.ENABLED = 1
WHERE css.SUPERAPP_ROLE_ID = HEXTORAW(:1)
  AND css.IS_DELETED = 0
  AND css.IS_ENABLED = 1
ORDER BY css.CREATED_AT, css.ID`

	detailRows, err := r.db.QueryContext(ctx, rowsQ, roleHex)
	if err != nil {
		r.logger.Errorf("[CustomerSegmentation][fillAggregateBySuperAppRoleID] detail query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer detailRows.Close()

	for detailRows.Next() {
		var cssId, gName, sName, subName sql.NullString
		if err := detailRows.Scan(&cssId, &gName, &sName, &subName); err != nil {
			r.logger.Errorf("[CustomerSegmentation][fillAggregateBySuperAppRoleID] detail scan failed: %v", err)
			return nil, local_util.HandleDBError(err)
		}
		seg.CustomerSegments = append(seg.CustomerSegments, imodel.CustSegment{
			Id:                 cssId.String,
			CustomerGroup:      gName.String,
			CustomerSegment:    sName.String,
			CustomerSubSegment: subName.String,
		})
	}

	if len(seg.CustomerSegments) == 0 {
		return nil, errors.New(localization.ErrorResourceNotFound.Code)
	}

	return &seg, nil
}

func (r *customerStorage) FindByID(ctx context.Context, id string) (*imodel.CustomerSegmentation, error) {
	idHex, ok := normalizeRawHex32(id)
	if !ok {
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	const roleFromSegQ = `
SELECT RAWTOHEX(css.SUPERAPP_ROLE_ID)
FROM CUSTOMER_SUB_SEGMENTS css
JOIN SUPERAPP_ROLE sar ON sar.ID = css.SUPERAPP_ROLE_ID AND sar.IS_DELETED = 0 AND sar.ENABLED = 1
WHERE css.CUSTOMER_SEGMENTATIONS_ID = HEXTORAW(:1)
  AND css.IS_DELETED = 0
  AND css.IS_ENABLED = 1
FETCH FIRST 1 ROWS ONLY`

	var roleHex string
	err := r.db.QueryRowContext(ctx, roleFromSegQ, idHex).Scan(&roleHex)
	if err == nil {
		rh, ok := normalizeRawHex32(roleHex)
		if !ok {
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		return r.fillAggregateBySuperAppRoleID(ctx, rh, idHex)
	}
	if !errors.Is(err, sql.ErrNoRows) {
		r.logger.Errorf("[CustomerSegmentation][FindByID] role-from-seg query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	const roleOnlyQ = `
SELECT 1
FROM SUPERAPP_ROLE
WHERE ID = HEXTORAW(:1) AND IS_DELETED = 0 AND ENABLED = 1`
	var one int
	if err := r.db.QueryRowContext(ctx, roleOnlyQ, idHex).Scan(&one); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		r.logger.Errorf("[CustomerSegmentation][FindByID] role lookup failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	const docIDQ = `
SELECT RAWTOHEX(MIN(cs.ID))
FROM CUSTOMER_SUB_SEGMENTS css
JOIN CUSTOMER_SEGMENTATIONS cs ON cs.ID = css.CUSTOMER_SEGMENTATIONS_ID AND cs.IS_DELETED = 0
WHERE css.SUPERAPP_ROLE_ID = HEXTORAW(:1)
  AND css.IS_DELETED = 0
  AND css.IS_ENABLED = 1`
	var docID string
	if err := r.db.QueryRowContext(ctx, docIDQ, idHex).Scan(&docID); err != nil {
		r.logger.Errorf("[CustomerSegmentation][FindByID] doc id query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	dh, ok := normalizeRawHex32(docID)
	if !ok {
		return nil, errors.New(localization.ErrorResourceNotFound.Code)
	}
	return r.fillAggregateBySuperAppRoleID(ctx, idHex, dh)
}

func (r *customerStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.CustomerSegmentation], error) {
	limit := int64(50)
	page := int64(1)
	if filterParam.PerPage > 0 {
		limit = int64(filterParam.PerPage)
	}
	if filterParam.Page > 0 {
		page = int64(filterParam.Page)
	}
	offset := (page - 1) * limit

	clauses := []string{
		"css.IS_DELETED = 0",
		"css.IS_ENABLED = 1",
		"cs.IS_DELETED = 0",
		"cg.IS_DELETED = 0",
		"cg.IS_ENABLED = 1",
		"sar.IS_DELETED = 0",
		"sar.ENABLED = 1",
	}
	var args []interface{}

	search := strings.TrimSpace(filterParam.Search)
	if search != "" {
		clauses = append(clauses,
			`(
				LOWER(cg.NAME) LIKE '%' || LOWER(:search) || '%'
				OR LOWER(cs.NAME) LIKE '%' || LOWER(:search) || '%'
				OR LOWER(css.NAME) LIKE '%' || LOWER(:search) || '%'
				OR LOWER(sar.NAME) LIKE '%' || LOWER(:search) || '%'
				OR LOWER(sar.ROLE_CODE) LIKE '%' || LOWER(:search) || '%'
			)`,
		)
		args = append(args, sql.Named("search", search))
	}

	if filterParam.Filters != nil {
		if v, ok := filterParam.Filters["is_enabled"]; ok {
			if b, ok2 := parseBoolFilter(v); ok2 {
				clauses = append(clauses, "cs.IS_ENABLED = :is_enabled")
				args = append(args, sql.Named("is_enabled", boolToOracleNumber(b)))
			}
		}
	}

	where := strings.Join(clauses, " AND ")
	joinFrom := `
FROM CUSTOMER_SUB_SEGMENTS css
JOIN CUSTOMER_SEGMENTATIONS cs ON cs.ID = css.CUSTOMER_SEGMENTATIONS_ID
JOIN CUSTOMER_GROUPS cg ON cg.ID = cs.CUSTOMER_GROUPS_ID
JOIN SUPERAPP_ROLE sar ON sar.ID = css.SUPERAPP_ROLE_ID`

	countQ := fmt.Sprintf(`
SELECT COUNT(*) FROM (
  SELECT css.SUPERAPP_ROLE_ID
  %s
  WHERE %s
  GROUP BY css.SUPERAPP_ROLE_ID
) role_buckets`, joinFrom, where)

	var total int64
	if err := r.db.QueryRowContext(ctx, countQ, args...).Scan(&total); err != nil {
		r.logger.Errorf("[CustomerSegmentation][FindAllWithPagination] count failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	listQ := fmt.Sprintf(`
SELECT
  RAWTOHEX(css.SUPERAPP_ROLE_ID),
  RAWTOHEX(MIN(cs.ID)),
  MAX(cs.CREATED_AT) AS ord_ts
%s
WHERE %s
GROUP BY css.SUPERAPP_ROLE_ID
ORDER BY ord_ts DESC
OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`, joinFrom, where)

	listArgs := append(args, sql.Named("offset", offset), sql.Named("limit", limit))
	rows, err := r.db.QueryContext(ctx, listQ, listArgs...)
	if err != nil {
		r.logger.Errorf("[CustomerSegmentation][FindAllWithPagination] list query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer rows.Close()

	var list []imodel.CustomerSegmentation
	for rows.Next() {
		var roleHexStr, docSegHex string
		var ordTs sql.NullTime
		if err := rows.Scan(&roleHexStr, &docSegHex, &ordTs); err != nil {
			r.logger.Errorf("[CustomerSegmentation][FindAllWithPagination] scan failed: %v", err)
			return nil, local_util.HandleDBError(err)
		}
		rh, ok := normalizeRawHex32(roleHexStr)
		if !ok {
			r.logger.Errorf("[CustomerSegmentation][FindAllWithPagination] invalid role hex: %s", roleHexStr)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		dh, ok := normalizeRawHex32(docSegHex)
		if !ok {
			r.logger.Errorf("[CustomerSegmentation][FindAllWithPagination] invalid doc seg hex: %s", docSegHex)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		seg, err := r.fillAggregateBySuperAppRoleID(ctx, rh, dh)
		if err != nil {
			return nil, err
		}
		list = append(list, *seg)
	}

	meta := local_util.BuildPaginationMeta(total, int(page), int(limit))
	return &types.PaginatedResponse[[]imodel.CustomerSegmentation]{Data: list, Meta: meta}, nil
}

func (r *customerStorage) FindByCustomerSegmentation(ctx context.Context, customerSegment string) (*imodel.CustomerSegmentation, error) {
	customerSegment = strings.TrimSpace(customerSegment)
	if customerSegment == "" {
		return nil, errors.New(localization.ErrorNoDataProvided.Code)
	}

	const q = `
SELECT
  RAWTOHEX(css.SUPERAPP_ROLE_ID),
  RAWTOHEX(MIN(cs.ID))
FROM CUSTOMER_SEGMENTATIONS cs
JOIN CUSTOMER_SUB_SEGMENTS css ON css.CUSTOMER_SEGMENTATIONS_ID = cs.ID
JOIN SUPERAPP_ROLE sar ON sar.ID = css.SUPERAPP_ROLE_ID AND sar.IS_DELETED = 0 AND sar.ENABLED = 1
JOIN CUSTOMER_GROUPS cg ON cg.ID = cs.CUSTOMER_GROUPS_ID AND cg.IS_DELETED = 0 AND cg.IS_ENABLED = 1
WHERE cs.IS_DELETED = 0
  AND css.IS_DELETED = 0
  AND css.IS_ENABLED = 1
  AND UPPER(TRIM(cs.NAME)) = UPPER(TRIM(:1))
GROUP BY css.SUPERAPP_ROLE_ID
ORDER BY MIN(css.CREATED_AT)
FETCH FIRST 1 ROWS ONLY`

	var roleHexStr, docSegHex string
	err := r.db.QueryRowContext(ctx, q, customerSegment).Scan(&roleHexStr, &docSegHex)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		r.logger.Errorf("[CustomerSegmentation][FindByCustomerSegmentation] query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	rh, ok := normalizeRawHex32(roleHexStr)
	if !ok {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	dh, ok := normalizeRawHex32(docSegHex)
	if !ok {
		return nil, errors.New(localization.ErrorResourceNotFound.Code)
	}
	return r.fillAggregateBySuperAppRoleID(ctx, rh, dh)
}

func (s *customerStorage) CheckIfCustomerSubSegmentExists(ctx context.Context, id string) (bool, error) {
	q := fmt.Sprintf(`SELECT COUNT(*) FROM CUSTOMER_SUB_SEGMENTS WHERE SUPERAPP_ROLE_ID = '%s'`, id)
	var count int
	if err := s.db.QueryRowContext(ctx, q).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}
