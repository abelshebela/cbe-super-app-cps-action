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
	"go.mongodb.org/mongo-driver/v2/bson"
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

const (
	customerSegmentationTable = "customer_segmentations"
	customerSubSegmentsTable  = "customer_sub_segments"
)

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

func convertRawHexToObjectID(rawHex string) bson.ObjectID {
	rawHex = strings.TrimSpace(rawHex)
	if len(rawHex) >= 24 {
		if oid, err := bson.ObjectIDFromHex(rawHex[:24]); err == nil {
			return oid
		}
	}
	return bson.NewObjectID()
}

func (r *customerStorage) Create(ctx context.Context, seg *imodel.CustomerSegmentation) error {
	if seg == nil {
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

	var id string

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Errorf("[CustomerSegmentation][Create] begin tx failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = tx.Rollback() }()

	const insertSegQ = `
INSERT INTO CUSTOMER_SEGMENTATIONS (
  ROLE_NAME,
  IS_ENABLED,
  IS_DELETED,
  CREATED_AT,
  UPDATED_AT
)
VALUES (
  :1,:2,:3,:4,:5
)
RETURNING RAWTOHEX(ID) INTO :6`

	if _, err := tx.ExecContext(ctx, insertSegQ,
		seg.CustomerRole.Name,
		boolToOracleNumber(seg.IsEnabled),
		boolToOracleNumber(seg.IsDeleted),
		seg.CreatedAt,
		seg.UpdatedAt,
		sql.Out{Dest: &id},
	); err != nil {
		r.logger.Errorf("[CustomerSegmentation][Create] insert segmentation failed: %v", err)
		return local_util.HandleDBError(err)
	}

	if len(seg.CustomerSubSegments) > 0 {
		const insertSubQ = `
INSERT INTO CUSTOMER_SUB_SEGMENTS (
  CUSTOMER_SEG_ID,
  NAME,
  CUSTOMER_GROUP,
  CUSTOMER_SEGMENT
)
VALUES (
  HEXTORAW(:1),:2,:3,:4
)`

		for _, sub := range seg.CustomerSubSegments {
			if _, err := tx.ExecContext(ctx, insertSubQ,
				id,
				sub.Name,
				sub.CustomerGroup,
				sub.CustomerSegment,
			); err != nil {
				r.logger.Errorf("[CustomerSegmentation][Create] insert sub segment failed: %v", err)
				return local_util.HandleDBError(err)
			}
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
	if strings.TrimSpace(id) == "" {
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

	const updateSegQ = `
UPDATE CUSTOMER_SEGMENTATIONS
SET
  ROLE_NAME   = :1,
  UPDATED_AT  = SYSTIMESTAMP
WHERE ID = HEXTORAW(:2) AND IS_DELETED = 0`

	res, err := tx.ExecContext(ctx, updateSegQ,
		seg.CustomerRole.Name,
		id,
	)
	if err != nil {
		r.logger.Errorf("[CustomerSegmentation][Update] update segmentation failed: %v", err)
		return local_util.HandleDBError(err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}

	const deleteSubQ = `DELETE FROM CUSTOMER_SUB_SEGMENTS WHERE CUSTOMER_SEG_ID = HEXTORAW(:1)`
	if _, err := tx.ExecContext(ctx, deleteSubQ, id); err != nil {
		r.logger.Errorf("[CustomerSegmentation][Update] delete sub segments failed: %v", err)
		return local_util.HandleDBError(err)
	}

	if len(seg.CustomerSubSegments) > 0 {
		const insertSubQ = `
INSERT INTO CUSTOMER_SUB_SEGMENTS (
  CUSTOMER_SEG_ID,
  NAME,
  CUSTOMER_GROUP,
  CUSTOMER_SEGMENT
)
VALUES (
  HEXTORAW(:1),:2,:3,:4
)`
		for _, sub := range seg.CustomerSubSegments {
			if _, err := tx.ExecContext(ctx, insertSubQ,
				id,
				sub.Name,
				sub.CustomerGroup,
				sub.CustomerSegment,
			); err != nil {
				r.logger.Errorf("[CustomerSegmentation][Update] insert sub segment failed: %v", err)
				return local_util.HandleDBError(err)
			}
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
	if strings.TrimSpace(id) == "" {
		return errors.New(localization.ErrorInvalidID.Code)
	}

	const q = `
UPDATE CUSTOMER_SEGMENTATIONS
SET
  IS_ENABLED  = :1,
  UPDATED_AT  = SYSTIMESTAMP
WHERE ID = HEXTORAW(:2) AND IS_DELETED = 0`

	res, err := r.db.ExecContext(ctx, q, boolToOracleNumber(enable), id)
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
	if strings.TrimSpace(id) == "" {
		return errors.New(localization.ErrorInvalidID.Code)
	}

	const q = `
UPDATE CUSTOMER_SEGMENTATIONS
SET
  IS_DELETED = 1,
  UPDATED_AT = SYSTIMESTAMP
WHERE ID = HEXTORAW(:1)`

	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		r.logger.Errorf("[CustomerSegmentation][Delete] delete failed: %v", err)
		return local_util.HandleDBError(err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}
	return nil
}

func (r *customerStorage) FindByID(ctx context.Context, id string) (*imodel.CustomerSegmentation, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	const segQ = `
SELECT
  RAWTOHEX(ID),
  ROLE_NAME,
  IS_ENABLED,
  IS_DELETED,
  CREATED_AT,
  UPDATED_AT
FROM CUSTOMER_SEGMENTATIONS
WHERE ID = HEXTORAW(:1) AND IS_DELETED = 0`

	var (
		segID     string
		roleName  sql.NullString
		isEnabled int
		isDeleted int
		createdAt sql.NullTime
		updatedAt sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, segQ, id).Scan(
		&segID,
		&roleName,
		&isEnabled,
		&isDeleted,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		r.logger.Errorf("[CustomerSegmentation][FindByID] query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	seg := imodel.CustomerSegmentation{
		ID: convertRawHexToObjectID(segID),
		CustomerRole: imodel.CustomerRoleInfo{
			// Oracle schema stores only role name, so use segmentation-derived ID
			// to avoid returning zero ObjectID in responses.
			ID:   convertRawHexToObjectID(segID),
			Name: roleName.String,
		},
		IsEnabled: isEnabled == 1,
		IsDeleted: isDeleted == 1,
	}
	if createdAt.Valid {
		seg.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		seg.UpdatedAt = updatedAt.Time
	}

	const subQ = `
SELECT
  NAME,
  CUSTOMER_GROUP,
  CUSTOMER_SEGMENT
FROM CUSTOMER_SUB_SEGMENTS
WHERE CUSTOMER_SEG_ID = HEXTORAW(:1)
ORDER BY ID`

	rows, err := r.db.QueryContext(ctx, subQ, segID)
	if err != nil {
		r.logger.Errorf("[CustomerSegmentation][FindByID] sub segments query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			name, group, segment sql.NullString
		)
		if err := rows.Scan(&name, &group, &segment); err != nil {
			r.logger.Errorf("[CustomerSegmentation][FindByID] sub segments scan failed: %v", err)
			return nil, local_util.HandleDBError(err)
		}
		seg.CustomerSubSegments = append(seg.CustomerSubSegments, imodel.CustomerSubSegments{
			Name:            name.String,
			CustomerGroup:   group.String,
			CustomerSegment: segment.String,
		})
	}

	return &seg, nil
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

	clauses := []string{"cs.IS_DELETED = 0"}
	var args []interface{}

	search := strings.TrimSpace(filterParam.Search)
	if search != "" {
		clauses = append(clauses,
			`(
				LOWER(cs.ROLE_NAME) LIKE '%' || LOWER(:search) || '%'
				OR EXISTS (
					SELECT 1 FROM CUSTOMER_SUB_SEGMENTS css
					WHERE css.CUSTOMER_SEG_ID = cs.ID
					  AND (
					    LOWER(css.NAME) LIKE '%' || LOWER(:search) || '%'
					    OR LOWER(css.CUSTOMER_GROUP) LIKE '%' || LOWER(:search) || '%'
					    OR LOWER(css.CUSTOMER_SEGMENT) LIKE '%' || LOWER(:search) || '%'
					  )
				)
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

	countQ := fmt.Sprintf(`SELECT COUNT(*) FROM %s cs WHERE %s`, customerSegmentationTable, where)
	var total int64
	if err := r.db.QueryRowContext(ctx, countQ, args...).Scan(&total); err != nil {
		r.logger.Errorf("[CustomerSegmentation][FindAllWithPagination] count failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	listQ := fmt.Sprintf(`
SELECT
  RAWTOHEX(cs.ID),
  cs.ROLE_NAME,
  cs.IS_ENABLED,
  cs.IS_DELETED,
  cs.CREATED_AT,
  cs.UPDATED_AT
FROM %s cs
WHERE %s
ORDER BY cs.CREATED_AT DESC
OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`, customerSegmentationTable, where)

	listArgs := append(args, sql.Named("offset", offset), sql.Named("limit", limit))
	rows, err := r.db.QueryContext(ctx, listQ, listArgs...)
	if err != nil {
		r.logger.Errorf("[CustomerSegmentation][FindAllWithPagination] list query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer rows.Close()

	var list []imodel.CustomerSegmentation
	for rows.Next() {
		var (
			segID                    string
			roleName                 sql.NullString
			isEnabled, isDeleted     int
			createdAt, updatedAtTime sql.NullTime
		)
		if err := rows.Scan(
			&segID,
			&roleName,
			&isEnabled,
			&isDeleted,
			&createdAt,
			&updatedAtTime,
		); err != nil {
			r.logger.Errorf("[CustomerSegmentation][FindAllWithPagination] scan failed: %v", err)
			return nil, local_util.HandleDBError(err)
		}

		seg := imodel.CustomerSegmentation{
			ID: convertRawHexToObjectID(segID),
			CustomerRole: imodel.CustomerRoleInfo{
				ID:   convertRawHexToObjectID(segID),
				Name: roleName.String,
			},
			IsEnabled: isEnabled == 1,
			IsDeleted: isDeleted == 1,
		}
		if createdAt.Valid {
			seg.CreatedAt = createdAt.Time
		}
		if updatedAtTime.Valid {
			seg.UpdatedAt = updatedAtTime.Time
		}

		const subQ = `
SELECT
  NAME,
  CUSTOMER_GROUP,
  CUSTOMER_SEGMENT
FROM CUSTOMER_SUB_SEGMENTS
WHERE CUSTOMER_SEG_ID = HEXTORAW(:1)
ORDER BY ID`

		subRows, err := r.db.QueryContext(ctx, subQ, segID)
		if err != nil {
			r.logger.Errorf("[CustomerSegmentation][FindAllWithPagination] sub segments query failed: %v", err)
			return nil, local_util.HandleDBError(err)
		}

		for subRows.Next() {
			var (
				name, group, segment sql.NullString
			)
			if err := subRows.Scan(&name, &group, &segment); err != nil {
				_ = subRows.Close()
				r.logger.Errorf("[CustomerSegmentation][FindAllWithPagination] sub segments scan failed: %v", err)
				return nil, local_util.HandleDBError(err)
			}
			seg.CustomerSubSegments = append(seg.CustomerSubSegments, imodel.CustomerSubSegments{
				Name:            name.String,
				CustomerGroup:   group.String,
				CustomerSegment: segment.String,
			})
		}
		_ = subRows.Close()

		list = append(list, seg)
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
  RAWTOHEX(cs.ID),
  cs.ROLE_NAME,
  cs.IS_ENABLED,
  cs.IS_DELETED,
  cs.CREATED_AT,
  cs.UPDATED_AT
FROM CUSTOMER_SEGMENTATIONS cs
WHERE cs.IS_DELETED = 0
  AND EXISTS (
    SELECT 1 FROM CUSTOMER_SUB_SEGMENTS css
    WHERE css.CUSTOMER_SEG_ID = cs.ID
      AND css.CUSTOMER_SEGMENT = :1
  )
FETCH FIRST 1 ROWS ONLY`

	var (
		segID                    string
		roleName                 sql.NullString
		isEnabled, isDeleted     int
		createdAt, updatedAtTime sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, q, customerSegment).Scan(
		&segID,
		&roleName,
		&isEnabled,
		&isDeleted,
		&createdAt,
		&updatedAtTime,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		r.logger.Errorf("[CustomerSegmentation][FindByCustomerSegmentation] query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	seg := imodel.CustomerSegmentation{
		ID: convertRawHexToObjectID(segID),
		CustomerRole: imodel.CustomerRoleInfo{
			ID:   convertRawHexToObjectID(segID),
			Name: roleName.String,
		},
		IsEnabled: isEnabled == 1,
		IsDeleted: isDeleted == 1,
	}

	if createdAt.Valid {
		seg.CreatedAt = createdAt.Time
	}
	if updatedAtTime.Valid {
		seg.UpdatedAt = updatedAtTime.Time
	}

	const subQ = `
SELECT
  NAME,
  CUSTOMER_GROUP,
  CUSTOMER_SEGMENT
FROM CUSTOMER_SUB_SEGMENTS
WHERE CUSTOMER_SEG_ID = HEXTORAW(:1)
ORDER BY ID`

	rows, err := r.db.QueryContext(ctx, subQ, segID)
	if err != nil {
		r.logger.Errorf("[CustomerSegmentation][FindByCustomerSegmentation] sub segments query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			name, group, segment sql.NullString
		)
		if err := rows.Scan(&name, &group, &segment); err != nil {
			r.logger.Errorf("[CustomerSegmentation][FindByCustomerSegmentation] sub segments scan failed: %v", err)
			return nil, local_util.HandleDBError(err)
		}
		seg.CustomerSubSegments = append(seg.CustomerSubSegments, imodel.CustomerSubSegments{
			Name:            name.String,
			CustomerGroup:   group.String,
			CustomerSegment: segment.String,
		})
	}

	return &seg, nil
}
