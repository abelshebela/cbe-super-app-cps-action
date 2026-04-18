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
	r.logger.Infof("[CustomerSegmentation][assertSuperAppRoleExistsTx] called with roleHex='%s'", roleHex)
	const q = `
	SELECT 1
	FROM SUPERAPP_ROLE
	WHERE ID = HEXTORAW(:1) AND IS_DELETED = 0 AND IS_ENABLED = 1`
	var one int
	r.logger.Debugf("[CustomerSegmentation][assertSuperAppRoleExistsTx] Executing query for roleHex: '%s'", roleHex)
	err := tx.QueryRowContext(ctx, q, roleHex).Scan(&one)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Warnf("[CustomerSegmentation][assertSuperAppRoleExistsTx] No SUPERAPP_ROLE found for roleHex: '%s'", roleHex)
			return errors.New(localization.ErrorResourceNotFound.Code)
		}
		r.logger.Errorf("[CustomerSegmentation][assertSuperAppRoleExistsTx] query failed: %v", err)
		return local_util.HandleDBError(err)
	}
	r.logger.Infof("[CustomerSegmentation][assertSuperAppRoleExistsTx] SUPERAPP_ROLE exists for roleHex: '%s'", roleHex)
	return nil
}

func (r *customerStorage) findOrCreateCustomerGroupTx(ctx context.Context, tx *sql.Tx, name string, now time.Time, isEnabled bool) (string, error) {
	r.logger.Infof("[CustomerSegmentation][findOrCreateCustomerGroupTx] called with name='%s', isEnabled=%v, now=%v", name, isEnabled, now)

	const findQ = `
	SELECT RAWTOHEX(ID)
	FROM CUSTOMER_GROUPS
	WHERE UPPER(TRIM(NAME)) = UPPER(TRIM(:1)) AND IS_DELETED = 0`
	var existing string
	r.logger.Debugf("[CustomerSegmentation][findOrCreateCustomerGroupTx] Executing find query for group name: '%s'", name)
	err := tx.QueryRowContext(ctx, findQ, name).Scan(&existing)
	if err == nil {
		r.logger.Infof("[CustomerSegmentation][findOrCreateCustomerGroupTx] Found existing group: %s", existing)
		return strings.ToLower(existing), nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		r.logger.Errorf("[CustomerSegmentation][findOrCreateCustomerGroupTx] find failed: %v", err)
		return "", local_util.HandleDBError(err)
	}

	r.logger.Infof("[CustomerSegmentation][findOrCreateCustomerGroupTx] No existing group found, creating new group for name: '%s'", name)
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
	r.logger.Debugf("[CustomerSegmentation][findOrCreateCustomerGroupTx] Executing insert query for group name: '%s', isEnabled: %v", name, isEnabled)
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
	r.logger.Infof("[CustomerSegmentation][findOrCreateCustomerGroupTx] Created new group with ID: %s", newID)
	return strings.ToLower(newID), nil
}

func (r *customerStorage) findOrCreateSegmentationTx(ctx context.Context, tx *sql.Tx, groupHex, segName string, now time.Time, isEnabled bool) (string, error) {
	r.logger.Infof("[CustomerSegmentation][findOrCreateSegmentationTx] called with groupHex='%s', segName='%s', isEnabled=%v, now=%v", groupHex, segName, isEnabled, now)

	const findQ = `
	SELECT RAWTOHEX(ID)
	FROM CUSTOMER_SEGMENTATIONS
	WHERE CUSTOMER_GROUP_ID = HEXTORAW(:1)
		AND UPPER(TRIM(NAME)) = UPPER(TRIM(:2))
		AND IS_DELETED = 0`
	var existing string
	r.logger.Debugf("[CustomerSegmentation][findOrCreateSegmentationTx] Executing find query for groupHex: '%s', segName: '%s'", groupHex, segName)
	err := tx.QueryRowContext(ctx, findQ, groupHex, segName).Scan(&existing)
	if err == nil {
		r.logger.Infof("[CustomerSegmentation][findOrCreateSegmentationTx] Found existing segmentation: %s", existing)
		return strings.ToLower(existing), nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		r.logger.Errorf("[CustomerSegmentation][findOrCreateSegmentationTx] find failed: %v", err)
		return "", local_util.HandleDBError(err)
	}

	r.logger.Infof("[CustomerSegmentation][findOrCreateSegmentationTx] No existing segmentation found, creating new segmentation for groupHex: '%s', segName: '%s'", groupHex, segName)
	var newID string
	const insertQ = `
	INSERT INTO CUSTOMER_SEGMENTATIONS (
		NAME,
		CUSTOMER_GROUP_ID,
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
	r.logger.Debugf("[CustomerSegmentation][findOrCreateSegmentationTx] Executing insert query for segName: '%s', groupHex: '%s', isEnabled: %v", segName, groupHex, isEnabled)
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
	r.logger.Infof("[CustomerSegmentation][findOrCreateSegmentationTx] Created new segmentation with ID: %s", newID)
	return strings.ToLower(newID), nil
}

func (r *customerStorage) resolveSuperAppRoleIDTx(ctx context.Context, tx *sql.Tx, idHex string) (string, error) {
	const roleFromSegQ = `
SELECT RAWTOHEX(css.SUPERAPP_ROLE_ID)
FROM CUSTOMER_SUB_SEGMENTS css
WHERE css.CUSTOMER_SEGMENTATION_ID = HEXTORAW(:1)
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
	r.logger.Infof("[CustomerSegmentation][Create] called with segmentation: %+v", seg)
	if seg == nil {
		r.logger.Warnf("[CustomerSegmentation][Create] segmentation is nil")
		return errors.New(localization.ErrorNoDataProvided.Code)
	}
	if len(seg.CustomerSegments) == 0 {
		r.logger.Warnf("[CustomerSegmentation][Create] segmentation has no CustomerSegments")
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

	r.logger.Debugf("[CustomerSegmentation][Create] Normalizing CustomerRole.ID: %s", seg.CustomerRole.ID)
	roleIDHex, ok := normalizeRawHex32(seg.CustomerRole.ID)
	if !ok {
		r.logger.Warnf("[CustomerSegmentation][Create] Invalid CustomerRole.ID: %s", seg.CustomerRole.ID)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	r.logger.Debugf("[CustomerSegmentation][Create] Beginning transaction")
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Errorf("[CustomerSegmentation][Create] begin tx failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = tx.Rollback() }()

	r.logger.Debugf("[CustomerSegmentation][Create] Asserting SUPERAPP_ROLE exists for roleIDHex: %s", roleIDHex)
	if err := r.assertSuperAppRoleExistsTx(ctx, tx, roleIDHex); err != nil {
		r.logger.Warnf("[CustomerSegmentation][Create] SUPERAPP_ROLE does not exist for roleIDHex: %s", roleIDHex)
		return err
	}

	const insertSubQ = `
INSERT INTO CUSTOMER_SUB_SEGMENTS (
  NAME,
  CUSTOMER_SEGMENTATION_ID,
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

	for i, sub := range seg.CustomerSegments {
		r.logger.Debugf("[CustomerSegmentation][Create] Processing CustomerSegment[%d]: %+v", i, sub)
		groupHex, err := r.findOrCreateCustomerGroupTx(ctx, tx, sub.CustomerGroup, now, seg.IsEnabled)
		if err != nil {
			r.logger.Errorf("[CustomerSegmentation][Create] findOrCreateCustomerGroupTx failed for CustomerGroup='%s': %v", sub.CustomerGroup, err)
			return err
		}
		segHex, err := r.findOrCreateSegmentationTx(ctx, tx, groupHex, sub.CustomerSegment, now, seg.IsEnabled)
		if err != nil {
			r.logger.Errorf("[CustomerSegmentation][Create] findOrCreateSegmentationTx failed for CustomerSegment='%s': %v", sub.CustomerSegment, err)
			return err
		}
		r.logger.Debugf("[CustomerSegmentation][Create] Inserting sub segment: CustomerSubSegment='%s', segHex='%s', roleIDHex='%s'", sub.CustomerSubSegment, segHex, roleIDHex)
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

	r.logger.Debugf("[CustomerSegmentation][Create] Committing transaction")
	if err := tx.Commit(); err != nil {
		r.logger.Errorf("[CustomerSegmentation][Create] commit failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	r.logger.Infof("[CustomerSegmentation][Create] Successfully created segmentation and committed transaction")
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
WHERE ID = HEXTORAW(:1) AND IS_DELETED = 0 AND IS_ENABLED = 1`
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
  CUSTOMER_SEGMENTATION_ID,
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
WHERE CUSTOMER_SEGMENTATION_ID = HEXTORAW(:1) AND IS_DELETED = 0`
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
  sar.IS_ENABLED,
  sar.IS_DELETED,
  MIN(cs.IS_ENABLED),
  MIN(css.IS_ENABLED),
  MIN(css.CREATED_AT),
  MAX(GREATEST(cs.LAST_MODIFIED_AT, css.LAST_MODIFIED_AT))
FROM CUSTOMER_SUB_SEGMENTS css
JOIN SUPERAPP_ROLE sar ON sar.ID = css.SUPERAPP_ROLE_ID AND sar.IS_DELETED = 0 AND sar.IS_ENABLED = 1
JOIN CUSTOMER_SEGMENTATIONS cs ON cs.ID = css.CUSTOMER_SEGMENTATION_ID AND cs.IS_DELETED = 0
JOIN CUSTOMER_GROUPS cg ON cg.ID = cs.CUSTOMER_GROUP_ID AND cg.IS_DELETED = 0 AND cg.IS_ENABLED = 1
WHERE css.SUPERAPP_ROLE_ID = HEXTORAW(:1)
  AND css.IS_DELETED = 0
  AND css.IS_ENABLED = 1
GROUP BY sar.ID, sar.NAME, sar.IS_ENABLED, sar.IS_DELETED`

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
JOIN CUSTOMER_SEGMENTATIONS cs ON cs.ID = css.CUSTOMER_SEGMENTATION_ID AND cs.IS_DELETED = 0
JOIN CUSTOMER_GROUPS cg ON cg.ID = cs.CUSTOMER_GROUP_ID AND cg.IS_DELETED = 0 AND cg.IS_ENABLED = 1
JOIN SUPERAPP_ROLE sar ON sar.ID = css.SUPERAPP_ROLE_ID AND sar.IS_DELETED = 0 AND sar.IS_ENABLED = 1
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
JOIN SUPERAPP_ROLE sar ON sar.ID = css.SUPERAPP_ROLE_ID AND sar.IS_DELETED = 0 AND sar.IS_ENABLED = 1
WHERE css.CUSTOMER_SEGMENTATION_ID = HEXTORAW(:1)
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
WHERE ID = HEXTORAW(:1) AND IS_DELETED = 0 AND IS_ENABLED = 1`
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
JOIN CUSTOMER_SEGMENTATIONS cs ON cs.ID = css.CUSTOMER_SEGMENTATION_ID AND cs.IS_DELETED = 0
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

func (r *customerStorage) FindAllWithPagination(
	ctx context.Context,
	filterParam types.Filter,
) (*types.PaginatedResponse[[]imodel.CustomerSegmentation], error) {

	// ----------------------------
	// Pagination
	// ----------------------------
	limit := int64(50)
	page := int64(1)

	if filterParam.PerPage > 0 {
		limit = int64(filterParam.PerPage)
	}
	if filterParam.Page > 0 {
		page = int64(filterParam.Page)
	}

	offset := (page - 1) * limit

	// ----------------------------
	// Filters
	// ----------------------------
	clauses := []string{"1=1"}
	var args []interface{}

	if search := strings.TrimSpace(filterParam.Search); search != "" {
		clauses = append(clauses, `
			(
				LOWER(cg.NAME) LIKE '%' || LOWER(:search) || '%'
				OR LOWER(cs.NAME) LIKE '%' || LOWER(:search) || '%'
				OR LOWER(css.NAME) LIKE '%' || LOWER(:search) || '%'
				OR LOWER(sar.NAME) LIKE '%' || LOWER(:search) || '%'
				OR LOWER(sar.ROLE_CODE) LIKE '%' || LOWER(:search) || '%'
			)
		`)
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

	// ----------------------------
	// JOIN BLOCK
	// ----------------------------
	joinFrom := `
FROM CUSTOMER_SUB_SEGMENTS css

JOIN CUSTOMER_SEGMENTATIONS cs
  ON cs.ID = css.CUSTOMER_SEGMENTATION_ID
 AND cs.IS_DELETED = 0

JOIN CUSTOMER_GROUPS cg
  ON cg.ID = cs.CUSTOMER_GROUP_ID
 AND cg.IS_DELETED = 0
 AND cg.IS_ENABLED = 1

JOIN SUPERAPP_ROLE sar
  ON sar.ID = css.SUPERAPP_ROLE_ID
 AND sar.IS_DELETED = 0
 AND sar.IS_ENABLED = 1
`

	// ----------------------------
	// COUNT QUERY
	// ----------------------------
	countQ := fmt.Sprintf(`
SELECT COUNT(DISTINCT css.SUPERAPP_ROLE_ID)
%s
WHERE css.IS_DELETED = 0 AND %s
`, joinFrom, where)

	var total int64
	if err := r.db.QueryRowContext(ctx, countQ, args...).Scan(&total); err != nil {
		r.logger.Errorf("[CustomerSegmentation][FindAll] count failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	// ----------------------------
	// LIST QUERY
	// ----------------------------
	listQ := fmt.Sprintf(`
SELECT
  RAWTOHEX(css.SUPERAPP_ROLE_ID) AS ROLE_ID,

  sar.NAME AS ROLE_NAME,
  sar.ROLE_CODE,
  sar.DESCRIPTION,

  sar.IS_ENABLED,
  sar.IS_DELETED,

  sar.CREATED_AT,
  sar.LAST_MODIFIED_AT,
  sar.DELETED_AT,

  RAWTOHEX(cs.ID) AS SEGMENT_ID,
  cs.NAME AS SEGMENT_NAME,

  cg.NAME AS GROUP_NAME,

  RAWTOHEX(css.ID) AS SUB_SEGMENT_ID,
  css.NAME AS SUB_SEGMENT_NAME

%s
WHERE css.IS_DELETED = 0 AND %s
ORDER BY sar.CREATED_AT DESC
OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY
`, joinFrom, where)

	listArgs := append(args,
		sql.Named("offset", offset),
		sql.Named("limit", limit),
	)

	rows, err := r.db.QueryContext(ctx, listQ, listArgs...)
	if err != nil {
		r.logger.Errorf("[CustomerSegmentation][FindAll] list query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer rows.Close()

	// ----------------------------
	// GROUPING RESULT
	// ----------------------------
	roleMap := make(map[string]*imodel.CustomerSegmentation)
	order := make([]string, 0)

	for rows.Next() {

		var (
			roleID string

			roleName, roleCode, roleDesc    string
			enabledN, deletedN              int
			createdAt, updatedAt, deletedAt sql.NullTime

			segID, segName string
			groupName      string
			subSegID       string
			subSegName     string
		)

		if err := rows.Scan(
			&roleID,
			&roleName,
			&roleCode,
			&roleDesc,
			&enabledN,
			&deletedN,
			&createdAt,
			&updatedAt,
			&deletedAt,
			&segID,
			&segName,
			&groupName,
			&subSegID,
			&subSegName,
		); err != nil {
			r.logger.Errorf("[CustomerSegmentation][FindAll] scan failed: %v", err)
			return nil, local_util.HandleDBError(err)
		}

		role, exists := roleMap[roleID]
		if !exists {

			role = &imodel.CustomerSegmentation{
				ID: roleID,
				CustomerRole: imodel.CustomerRoleInfo{
					ID:   roleID,
					Name: roleName,
				},
				CustomerSegments: make([]imodel.CustSegment, 0),
				IsEnabled:        enabledN == 1,
				IsDeleted:        deletedN == 1,
			}

			if createdAt.Valid {
				role.CreatedAt = createdAt.Time
			}
			if updatedAt.Valid {
				role.UpdatedAt = updatedAt.Time
			}

			roleMap[roleID] = role
			order = append(order, roleID)
		}

		// ----------------------------
		// APPEND MULTIPLE SEGMENTS
		// ----------------------------
		role.CustomerSegments = append(role.CustomerSegments, imodel.CustSegment{
			Id:                 subSegID,
			CustomerGroup:      groupName,
			CustomerSegment:    segName,
			CustomerSubSegment: subSegName,
		})
	}

	// ----------------------------
	// FINAL RESPONSE BUILD
	// ----------------------------
	result := make([]imodel.CustomerSegmentation, 0, len(order))
	for _, id := range order {
		result = append(result, *roleMap[id])
	}

	meta := local_util.BuildPaginationMeta(total, int(page), int(limit))

	return &types.PaginatedResponse[[]imodel.CustomerSegmentation]{
		Data: result,
		Meta: meta,
	}, nil
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
JOIN CUSTOMER_SUB_SEGMENTS css ON css.CUSTOMER_SEGMENTATION_ID = cs.ID
JOIN SUPERAPP_ROLE sar ON sar.ID = css.SUPERAPP_ROLE_ID AND sar.IS_DELETED = 0 AND sar.IS_ENABLED = 1
JOIN CUSTOMER_GROUPS cg ON cg.ID = cs.CUSTOMER_GROUP_ID AND cg.IS_DELETED = 0 AND cg.IS_ENABLED = 1
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
