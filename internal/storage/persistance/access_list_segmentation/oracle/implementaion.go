package access_list_segmentation_oracle

import (
	access_list_segmentation_dto "cbe-super-app-cps-action/internal/constants/dto/access_list_segmentation"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	local_model "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/kafka"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"strings"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type accessListSegmentationOracle struct {
	db            DBTX
	cfg           config.VaultConfig
	kafkaProducer *kafka.AccessListSegmentationProducer
	accBlock      storage.AccountBlockRepository
	redisClient   storage.RedisRepository
	logger        utils.Logger
}

func NewAccessListSegmentationOracle(db DBTX, cfg config.VaultConfig, kafkaProducer *kafka.AccessListSegmentationProducer, redisRepository storage.RedisRepository, logger utils.Logger) storage.AccessListSegmentationRepositoryOracle {
	return &accessListSegmentationOracle{
		db:            db,
		cfg:           cfg,
		kafkaProducer: kafkaProducer,
		redisClient:   redisRepository,
		logger:        logger,
	}
}

// BulkDisable implements [storage.AccessListSegmentationRepository].
func (q *accessListSegmentationOracle) BulkDisable(ctx context.Context, req access_list_segmentation_dto.BulkDisableAccessListSegmentationRequest) error {
	q.logger.Infof("[AccessListSegmentation][BulkDisable] bulk disable request: %+v", req)
	var filter string
	var args []interface{}
	if len(req.Keys) == 0 {
		return nil
	}
	// Build IN clause with positional parameters and HEXTORAW for each key
	inClause := make([]string, len(req.Keys))
	args = append(args, req.ID)
	for i := range req.Keys {
		inClause[i] = fmt.Sprintf("HEXTORAW(:%d)", i+2)
		args = append(args, req.Keys[i])
	}

	var stmt string
	switch req.SegmentationType {
	case "account":
		filter = fmt.Sprintf("segmented_id = HEXTORAW(:1) AND access_list_key IN (%s)", strings.Join(inClause, ","))
		stmt = "DELETE FROM access_list_customer_seg WHERE " + filter
	case "block":
		filter = fmt.Sprintf("segmented_id = :1 AND access_list_key IN (%s)", strings.Join(inClause, ","))
		stmt = "DELETE FROM access_list_geo_seg WHERE " + filter
	default:
		return localization.ErrorUnexpectedError
	}

	res, err := q.db.ExecContext(ctx, stmt, args...)
	if err != nil {
		q.logger.Errorf("[AccessListSegmentation][BulkDisable] failed to bulk disable access list segmentation: %v", err)
		return localization.ErrorUnexpectedError
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		q.logger.Warnf("[AccessListSegmentation][BulkDisable] no rows deleted for filter: %s", filter)
		// Optionally, return a specific error or nil if that's not an error in your logic
		return localization.ErrorUnexpectedError
	}

	// q.redisClient.Delete(ctx, q.cfg.RedirectURI)

	return nil
}

func (q *accessListSegmentationOracle) CreateAccountSegment(ctx context.Context, accessListSegmentation access_list_segmentation_dto.CreateAccessListSegmentationRequest) error {
	n := len(accessListSegmentation.AccessListKeys)
	if n == 0 {
		return nil
	}

	valueStrings := make([]string, 0, n)
	valueArgs := make([]interface{}, 0, n*6)

	docs := make([]local_model.AccessListSegmentation, 0, n)

	paramIdx := 1
	for _, key := range accessListSegmentation.AccessListKeys {
		valueStrings = append(valueStrings, fmt.Sprintf("(SYS_GUID(), HEXTORAW(:%d), :%d, :%d, :%d, :%d, :%d)", paramIdx, paramIdx+1, paramIdx+2, paramIdx+3, paramIdx+4, paramIdx+5))
		valueArgs = append(valueArgs,
			key,                                   // access_list_key
			accessListSegmentation.SegmentationID, // segmented_id
			1,                                     // enabled
			time.Now(),                            // created_at
			time.Now(),                            // updated_at
			nil,                                   // deleted_at
		)
		docs = append(docs, local_model.AccessListSegmentation{
			AccessListKey: key,
			SegmentedID:   accessListSegmentation.SegmentationID,
			Enabled:       true,
		})
		paramIdx += 6
	}

	stmt := `
	       INSERT INTO ACCESS_LIST_BY_SUPERAPP_ROLE (
		   id, access_list_key, segmented_id, enabled, created_at, updated_at, deleted_at
	       ) VALUES ` + strings.Join(valueStrings, ",")

	_, err := q.db.ExecContext(ctx, stmt, valueArgs...)
	if err != nil {
		q.logger.Errorf("[AccessListSegmentation][CreateAccountSegment] failed to insert rows: %v", err)
		return err
	}

	res := map[string]any{
		"docs": docs,
		"type": "account-segment",
	}
	q.kafkaProducer.PublishMessage(ctx, res, "create", q.cfg.KafkaCustomerSegmentaionTopic, "create account-segment")

	return nil
}

func (q *accessListSegmentationOracle) CreateBlockSegment(ctx context.Context, accessListSegmentation access_list_segmentation_dto.CreateAccessListSegmentationRequest) error {
	if len(accessListSegmentation.SegmentationID) == 0 {
		return fmt.Errorf("SegmentationID is required")
	}

	n := len(accessListSegmentation.AccessListKeys)
	if n == 0 {
		return nil
	}

	now := time.Now()
	valueStrings := make([]string, 0, n)
	valueArgs := make([]interface{}, 0, n*8)

	docs := make([]local_model.AccessListSegmentation, 0, n)

	// Use positional parameters for each value
	paramIdx := 1
	for _, key := range accessListSegmentation.AccessListKeys {
		valueStrings = append(valueStrings, fmt.Sprintf("(SYS_GUID(), HEXTORAW(:%d), :%d, :%d, :%d, :%d, :%d, :%d)", paramIdx, paramIdx+1, paramIdx+2, paramIdx+3, paramIdx+4, paramIdx+5, paramIdx+6))
		valueArgs = append(valueArgs,
			key,                                   // access_list_key
			accessListSegmentation.SegmentationID, // segmented_id
			accessListSegmentation.Type,           // type
			1,                                     // enabled
			now,                                   // created_at
			now,                                   // updated_at
			nil,                                   // deleted_at
		)
		docs = append(docs, local_model.AccessListSegmentation{
			AccessListKey: key,
			SegmentedID:   accessListSegmentation.SegmentationID,
			Enabled:       true,
		})
		paramIdx += 7
	}

	stmt := `
		   INSERT INTO ACCESS_LIST_GEO_SEG (
		   id, access_list_key, segmented_id, type, enabled, created_at, updated_at, deleted_at
		   ) VALUES ` + strings.Join(valueStrings, ",")

	_, err := q.db.ExecContext(ctx, stmt, valueArgs...)
	if err != nil {
		q.logger.Errorf("[AccessListSegmentation][CreateBlockSegment] failed to insert rows: %v", err)
		return localization.ErrorUnexpectedError
	}

	res := map[string]any{
		"docs": docs,
		"type": "block-segment",
	}
	q.kafkaProducer.PublishMessage(ctx, res, "create", q.cfg.KafkaCustomerSegmentaionTopic, "create block-segment")

	return nil
}

// FindAllForBlock implements [storage.AccessListSegmentationRepository].
func (q *accessListSegmentationOracle) FindAllForBlock(ctx context.Context, geographicalID string) ([]local_model.APPAccessList, error) {

	query := `SELECT RAWTOHEX(g.access_list_key), a.name,a.service_key, RAWTOHEX(a.id), g.enabled
		FROM ACCESS_LIST_GEO_SEG g
		JOIN ACCESS_LISTS a ON g.access_list_key = a.id
		WHERE g.segmented_id = :1`
	rows, err := q.db.QueryContext(ctx, query, geographicalID)
	if err != nil {
		q.logger.Errorf("[AccessListSegmentation][FindAllForBlock] query failed: %v", err)
		return nil, err
	}
	defer rows.Close()
	var result []local_model.APPAccessList
	for rows.Next() {
		var key, accessListName, accessListServiceKey, accessListID string
		var enabled bool
		err := rows.Scan(&key, &accessListName, &accessListServiceKey, &accessListID, &enabled)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				q.logger.Infof("[AccessListSegmentation][FindAllForBlock] no access list segmentation found for geographicalID: %s", geographicalID)
				return nil, localization.ErrorResourceNotFound
			}
			q.logger.Errorf("[AccessListSegmentation][FindAllForBlock] query failed: %v", err)
			return nil, err
		}
		q.logger.Infof("[AccessListSegmentation][FindAllForBlock] found access list segmentation: key=%s, name=%s, service_key=%s, id=%s, enabled=%v", key, accessListName, accessListServiceKey, accessListID, enabled)
		result = append(result, local_model.APPAccessList{
			ID:             accessListID,
			Key:            accessListServiceKey,
			AccessListName: accessListName,
			Enabled:        enabled,
		})
	}
	return result, nil
}

// FindAllForBlockParents implements [storage.AccessListSegmentationRepositoryOracle].
func (q *accessListSegmentationOracle) FindAllForBlockParents(ctx context.Context, geographicalID string) ([]local_model.APPAccessList, error) {
	// Step 1: Fetch city_id, region_id, district_id for the given account block id
	var regionID, districtID sql.NullString
	err := q.db.QueryRowContext(ctx, `
    SELECT rawtohex(region_id), rawtohex(district_id)
    FROM ACCOUNT_BLOCKS
    WHERE id = HEXTORAW(:1)
`, geographicalID).Scan(&regionID, &districtID)
	if err != nil {
		q.logger.Errorf("[AccessListSegmentation][FindAllForBlockParents] failed to fetch block: %v", err)
		return nil, err
	}

	// Step 2: Build a slice of non-null IDs
	var ids []string
	if regionID.Valid {
		ids = append(ids, regionID.String)
	}
	if districtID.Valid {
		ids = append(ids, districtID.String)
	}
	if len(ids) == 0 {
		return nil, nil // No city/region/district to look up
	}

	// Step 3: Query ACCESS_LIST_GEO_SEG for these IDs
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf(":%d", i+1)
		args[i] = id
	}
	query := fmt.Sprintf(`
    SELECT RAWTOHEX(g.access_list_key), a.name, a.service_key, RAWTOHEX(a.id), g.enabled
    FROM ACCESS_LIST_GEO_SEG g
    JOIN ACCESS_LISTS a ON g.access_list_key = a.id
    WHERE g.segmented_id IN (%s)
`, strings.Join(placeholders, ","))

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		q.logger.Errorf("[AccessListSegmentation][FindAllForBlock] query failed: %v", err)
		return nil, err
	}
	defer rows.Close()
	var result []local_model.APPAccessList
	for rows.Next() {
		var key, accessListName, accessListServiceKey, accessListID string
		var enabled bool
		err := rows.Scan(&key, &accessListName, &accessListServiceKey, &accessListID, &enabled)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				q.logger.Infof("[AccessListSegmentation][FindAllForBlock] no access list segmentation found for geographicalID: %s", geographicalID)
				return nil, localization.ErrorResourceNotFound
			}
			q.logger.Errorf("[AccessListSegmentation][FindAllForBlock] query failed: %v", err)
			return nil, err
		}
		q.logger.Infof("[AccessListSegmentation][FindAllForBlock] found access list segmentation: key=%s, name=%s, service_key=%s, id=%s, enabled=%v", key, accessListName, accessListServiceKey, accessListID, enabled)
		result = append(result, local_model.APPAccessList{
			ID:             accessListID,
			Key:            accessListServiceKey,
			AccessListName: accessListName,
			Enabled:        enabled,
		})
	}
	return result, nil
}

// FindAllBySegmentIDorSegmentCode implements [storage.AccessListSegmentationRepository].
func (q *accessListSegmentationOracle) FindAllForAccount(ctx context.Context, customer_seg_id string) ([]local_model.APPAccessList, error) {
	q.logger.Infof("[AccessListSegmentationOracle][FindAllForAccount] called with customer_seg_id: %s", customer_seg_id)

	query := `SELECT RAWTOHEX(g.access_list_key), a.name,a.service_key, RAWTOHEX(a.id), g.enabled
		FROM ACCESS_LIST_BY_SUPERAPP_ROLE g
		JOIN ACCESS_LISTS a ON g.access_list_key = a.id
		WHERE g.segmented_id = HEXTORAW(:1)`
	rows, err := q.db.QueryContext(ctx, query, customer_seg_id)
	if err != nil {
		q.logger.Errorf("[AccessListSegmentation][FindAllForAccount] query failed: %v", err)
		return nil, err
	}
	defer rows.Close()
	var result []local_model.APPAccessList
	for rows.Next() {
		var key, accessListName, accessListServiceKey, accessListID string
		var enabled bool
		err := rows.Scan(&key, &accessListName, &accessListServiceKey, &accessListID, &enabled)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				q.logger.Infof("[AccessListSegmentation][FindAllForAccount] no access list segmentation found for customer_seg_id: %s", customer_seg_id)
				return nil, localization.ErrorResourceNotFound
			}
			q.logger.Errorf("[AccessListSegmentation][FindAllForAccount] query failed: %v", err)
			return nil, err
		}
		q.logger.Infof("[AccessListSegmentation][FindAllForAccount] found access list segmentation: key=%s, name=%s, service_key=%s, id=%s, enabled=%v", key, accessListName, accessListServiceKey, accessListID, enabled)
		result = append(result, local_model.APPAccessList{
			ID:             accessListID,
			Key:            accessListServiceKey,
			AccessListName: accessListName,
			Enabled:        enabled,
		})
	}
	q.logger.Infof("[AccessListSegmentation][FindAllForAccount] total access list segmentations found for customer_seg_id %s: %d", customer_seg_id, len(result))
	return result, nil
}

// FindAllBySegmentIDorSegmentCodeAndKeys implements [storage.AccessListSegmentationRepository].
func (q *accessListSegmentationOracle) FindAllByBlockAndKeys(ctx context.Context, segmentIDorCode string, keys []string) ([]model.AccessListSegmentation, error) {
	q.logger.Infof("[AccessListSegmentationOracle][FindAllByBlockAndKeys] called with segmentIDorCode: %s, keys: %v", segmentIDorCode, keys)
	if len(keys) == 0 {
		return nil, fmt.Errorf("no keys provided")
	}
	query := `SELECT RAWTOHEX(id), RAWTOHEX(access_list_key), RAWTOHEX(segmented_id), type, enabled FROM ACCESS_LIST_GEO_SEG WHERE segmented_id = HEXTORAW(:1) AND access_list_key IN (`
	placeholders := make([]string, len(keys))
	args := make([]interface{}, 0, len(keys)+1)
	args = append(args, segmentIDorCode)
	for i, k := range keys {
		placeholders[i] = fmt.Sprintf("HEXTORAW(:key%d)", i)
		args = append(args, k)
	}
	query = query + strings.Join(placeholders, ",") + `)`
	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		q.logger.Errorf("[AccessListSegmentationOracle][FindAllByBlockAndKeys] query failed: %v", err)
		return nil, err
	}
	defer rows.Close()
	var result []model.AccessListSegmentation
	for rows.Next() {
		var seg model.AccessListSegmentation
		err := rows.Scan(&seg.ID, &seg.AccessListKey, &seg.SegmentedID, &seg.Type, &seg.Enabled)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				q.logger.Infof("[AccessListSegmentationOracle][FindAllByBlockAndKeys] no access list segmentation found for segmentIDorCode: %s and keys: %v", segmentIDorCode, keys)
				return nil, localization.ErrorResourceNotFound
			}
			q.logger.Errorf("[AccessListSegmentationOracle][FindAllByBlockAndKeys] query failed: %v", err)
			return nil, err
		}
		result = append(result, seg)
	}
	return result, nil
}

// FindAllBySegmentIDorSegmentCodeAndKeys implements [storage.AccessListSegmentationRepository].
func (q *accessListSegmentationOracle) FindAllByAccountAndKeys(ctx context.Context, segmentIDorCode string, keys []string) ([]model.AccessListSegmentation, error) {
	q.logger.Infof("[AccessListSegmentationOracle][FindAllByAccountAndKeys] called with segmentIDorCode: %s, keys: %v", segmentIDorCode, keys)
	if len(keys) == 0 {
		return nil, fmt.Errorf("no keys provided")
	}
	query := `SELECT RAWTOHEX(ID), RAWTOHEX(access_list_key), RAWTOHEX(segmented_id), enabled FROM ACCESS_LIST_BY_SUPERAPP_ROLE WHERE segmented_id = HEXTORAW(:1) AND access_list_key IN (`
	placeholders := make([]string, len(keys))
	args := make([]interface{}, 0, len(keys)+1)
	args = append(args, segmentIDorCode)
	for i, k := range keys {
		placeholders[i] = fmt.Sprintf("HEXTORAW(:key%d)", i)
		args = append(args, k)
	}
	query = query + strings.Join(placeholders, ",") + `)`
	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		q.logger.Errorf("[AccessListSegmentationOracle][FindAllByAccountAndKeys] query failed: %v", err)
		return nil, err
	}
	defer rows.Close()
	var result []model.AccessListSegmentation
	for rows.Next() {
		var segID string
		var seg model.AccessListSegmentation
		err := rows.Scan(&segID, &seg.AccessListKey, &seg.SegmentedID, &seg.Enabled)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				q.logger.Infof("[AccessListSegmentationOracle][FindAllByAccountAndKeys] no access list segmentation found for segmentIDorCode: %s and keys: %v", segmentIDorCode, keys)
				return nil, localization.ErrorResourceNotFound
			}
			q.logger.Errorf("[AccessListSegmentationOracle][FindAllByAccountAndKeys] query failed: %v", err)
			return nil, err
		}
		result = append(result, seg)
	}
	return result, nil
}

// FindAllWithPagination implements [storage.AccessListSegmentationRepository].
func (q *accessListSegmentationOracle) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.AccessListSegmentation], error) {
	q.logger.Infof("[AccessListSegmentation][FindAllWithPagination] filter params: %+v", filterParam)

	// Validate pagination params
	if filterParam.Page < 1 || filterParam.PerPage < 1 {
		return nil, fmt.Errorf("invalid pagination: page and perPage must be >= 1")
	}

	// Validate filter keys
	allowedKeys := map[string]struct{}{"type": {}, "segmented_id": {}, "created_at": {}, "updated_at": {}, "enabled": {}}
	for key := range filterParam.Filters {
		if _, ok := allowedKeys[key]; !ok {
			return nil, fmt.Errorf("invalid filter key: %s", key)
		}
	}

	// Build WHERE clause
	var whereClauses []string
	var args []interface{}
	if filterParam.Search != "" {
		whereClauses = append(whereClauses, "(LOWER(type) LIKE :search OR LOWER(enabled) LIKE :search)")
		args = append(args, "%"+strings.ToLower(filterParam.Search)+"%")
	}
	for key := range allowedKeys {
		if val, ok := filterParam.Filters[key]; ok && val != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("%s = :%s", key, key))
			args = append(args, val)
		}
	}
	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM ACCESS_LIST_GEO_SEG %s`, whereSQL)
	countRow := q.db.QueryRowContext(ctx, countQuery, args...)
	var total int64
	if err := countRow.Scan(&total); err != nil {
		q.logger.Errorf("[AccessListSegmentation][FindAllWithPagination] count error: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	page := filterParam.Page
	perPage := filterParam.PerPage
	offset := (page - 1) * perPage

	if int64(offset+perPage) >= total {
		q.logger.Infof("[AccessListSegmentation][FindAllWithPagination] offset %d exceeds total %d, returning empty result", offset, total)
		return &types.PaginatedResponse[[]model.AccessListSegmentation]{
			Data: []model.AccessListSegmentation{},
			Meta: local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage),
		}, nil
	}

	// Main query
	query := fmt.Sprintf(`SELECT RAWTOHEX(id), RAWTOHEX(access_list_key), RAWTOHEX(segmented_id), type, enabled, created_at, updated_at, deleted_at FROM ACCESS_LIST_GEO_SEG %s ORDER BY created_at DESC OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`, whereSQL)
	argsWithPag := append(args, offset, perPage)
	rows, err := q.db.QueryContext(ctx, query, argsWithPag...)
	if err != nil {
		q.logger.Errorf("[AccessListSegmentation][FindAllWithPagination] query failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer rows.Close()
	var data []model.AccessListSegmentation
	for rows.Next() {
		var seg model.AccessListSegmentation
		err := rows.Scan(&seg.ID, &seg.AccessListKey, &seg.SegmentedID, &seg.Type, &seg.Enabled, &seg.CreatedAt, &seg.UpdatedAt, &seg.DeletedAt)
		if err != nil {
			q.logger.Errorf("[AccessListSegmentation][FindAllWithPagination] scan error: %v", err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		data = append(data, seg)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	q.logger.Infof("[AccessListSegmentation][FindAllWithPagination] retrieved %d segmentations", len(data))

	return &types.PaginatedResponse[[]model.AccessListSegmentation]{
		Data: data,
		Meta: meta,
	}, nil
}

// FindBlockSegmentByID implements [storage.AccessListSegmentationRepository].
func (q *accessListSegmentationOracle) FindBlockSegmentByID(ctx context.Context, id string) (*model.AccessListSegmentation, error) {
	q.logger.Infof("[AccessListSegmentationOracle][FindBlockSegmentByID] called with id: %s", id)
	query := `SELECT RAWTOHEX(id), RAWTOHEX(access_list_key), RAWTOHEX(segmented_id), type, enabled, created_at, updated_at, deleted_at FROM ACCESS_LIST_GEO_SEG WHERE id = HEXTORAW(:1)`
	row := q.db.QueryRowContext(ctx, query, id)
	var seg model.AccessListSegmentation
	err := row.Scan(&seg.ID, &seg.AccessListKey, &seg.SegmentedID, &seg.Type, &seg.Enabled, &seg.CreatedAt, &seg.UpdatedAt, &seg.DeletedAt)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			q.logger.Infof("[AccessListSegmentationOracle][FindBlockSegmentByID] no access list segmentation found for id: %s", id)
			return nil, localization.ErrorResourceNotFound
		}
		q.logger.Errorf("[AccessListSegmentationOracle][FindBlockSegmentByID] query failed: %v", err)
		return nil, err
	}
	return &seg, nil
}

// FindAccountSegmentByID implements [storage.AccessListSegmentationRepository].
func (q *accessListSegmentationOracle) FindAccountSegmentByID(ctx context.Context, id string) (*model.AccessListSegmentation, error) {
	q.logger.Infof("[AccessListSegmentationOracle][FindAccountSegmentByID] called with id: %s", id)
	query := `SELECT RAWTOHEX(id), RAWTOHEX(access_list_key), RAWTOHEX(segmented_id), type, enabled, created_at, updated_at, deleted_at FROM ACCESS_LIST_BY_SUPERAPP_ROLE WHERE id = HEXTORAW(:1)`
	row := q.db.QueryRowContext(ctx, query, id)
	var seg model.AccessListSegmentation
	err := row.Scan(&seg.ID, &seg.AccessListKey, &seg.SegmentedID, &seg.Type, &seg.Enabled, &seg.CreatedAt, &seg.UpdatedAt, &seg.DeletedAt)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			q.logger.Infof("[AccessListSegmentationOracle][FindAccountSegmentByID] no access list segmentation found for id: %s", id)
			return nil, localization.ErrorResourceNotFound
		}
		q.logger.Errorf("[AccessListSegmentationOracle][FindAccountSegmentByID] query failed: %v", err)
		return nil, err
	}
	return &seg, nil
}

// FindByIDAndType implements [storage.AccessListSegmentationRepository].
func (q *accessListSegmentationOracle) FindByIDAndType(ctx context.Context, ids string, t string) (*model.AccessListSegmentation, error) {
	q.logger.Infof("[AccessListSegmentationOracle][FindByIDAndType] called with ids: %s, type: %s", ids, t)
	query := `SELECT RAWTOHEX(id), RAWTOHEX(access_list_key), RAWTOHEX(segmented_id), type, enabled, created_at, updated_at, deleted_at FROM ACCESS_LIST_GEO_SEG WHERE id = HEXTORAW(:1) AND type = :2`
	row := q.db.QueryRowContext(ctx, query, ids, t)
	var seg model.AccessListSegmentation
	err := row.Scan(&seg.ID, &seg.AccessListKey, &seg.SegmentedID, &seg.Type, &seg.Enabled, &seg.CreatedAt, &seg.UpdatedAt, &seg.DeletedAt)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			q.logger.Infof("[AccessListSegmentationOracle][FindByIDAndType] no access list segmentation found for id: %s and type: %s", ids, t)
			return nil, localization.ErrorResourceNotFound
		}
		q.logger.Errorf("[AccessListSegmentationOracle][FindByIDAndType] query failed: %v", err)
		return nil, err
	}
	return &seg, nil

}

// FindByIDS implements [storage.AccessListSegmentationRepository].
func (q *accessListSegmentationOracle) FindByIDS(ctx context.Context, ids []string, t string) (*model.AccessListSegmentation, error) {
	panic("unimplemented")
}

// FindBySegmentIDAndAccessListKeys implements [storage.AccessListSegmentationRepository].
func (q *accessListSegmentationOracle) FindBySegmentIDAndAccessListKeys(ctx context.Context, id string, keys []string) (*model.AccessListSegmentation, error) {
	q.logger.Infof("[AccessListSegmentationOracle][FindBySegmentIDAndAccessListKeys] called with id: %s, keys: %v", id, keys)
	var cityID, regionID, districtID sql.NullString
	err := q.db.QueryRowContext(ctx, `
		SELECT city_id, region_id, district_id
		FROM ACCOUNT_BLOCKS
		WHERE id = :1
	`, id).Scan(&cityID, &regionID, &districtID)
	if err != nil {
		q.logger.Errorf("[AccessListSegmentation][FindBySegmentIDAndAccessListKeys] failed to fetch block: %v", err)
		return nil, err
	}
	if len(keys) == 0 {
		return nil, fmt.Errorf("no keys provided")
	}

	// Build list of non-null ids (city, region, district)
	var ids []string
	if cityID.Valid && strings.TrimSpace(cityID.String) != "" {
		ids = append(ids, cityID.String)
	}
	if regionID.Valid && strings.TrimSpace(regionID.String) != "" {
		ids = append(ids, regionID.String)
	}
	if districtID.Valid && strings.TrimSpace(districtID.String) != "" {
		ids = append(ids, districtID.String)
	}
	if len(ids) == 0 {
		return nil, nil // No city/region/district to look up
	}
	// Build IN clause for segmented_id
	segPlaceholders := make([]string, len(ids))
	segArgs := make([]interface{}, len(ids))
	for i, v := range ids {
		segPlaceholders[i] = fmt.Sprintf(":seg%d", i)
		segArgs[i] = v
	}

	// Build IN clause for access_list_key
	keyPlaceholders := make([]string, len(keys))
	keyArgs := make([]interface{}, len(keys))
	for i, k := range keys {
		keyPlaceholders[i] = fmt.Sprintf("HEXTORAW(:key%d)", i)
		keyArgs[i] = k
	}

	// Combine all args: first seg ids, then keys
	args := append(segArgs, keyArgs...)
	query := fmt.Sprintf(`SELECT RAWTOHEX(id), RAWTOHEX(access_list_key), RAWTOHEX(segmented_id), enabled, created_at, updated_at, deleted_at FROM ACCESS_LIST_GEO_SEG WHERE segmented_id IN (%s) AND access_list_key IN (%s)`, strings.Join(segPlaceholders, ","), strings.Join(keyPlaceholders, ","))
	row := q.db.QueryRowContext(ctx, query, args...)
	var seg model.AccessListSegmentation
	err = row.Scan(&seg.ID, &seg.AccessListKey, &seg.SegmentedID, &seg.Type, &seg.Enabled, &seg.CreatedAt, &seg.UpdatedAt, &seg.DeletedAt)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			q.logger.Errorf("[AccessListSegmentationOracle][FindBySegmentIDAndAccessListKeys] query failed: %v", err)
			return nil, nil
		}
		return nil, err
	}
	return &seg, nil
}

// FindByAccountSegmentationAndAccessListKeys implements [storage.AccessListSegmentationRepository].
func (q *accessListSegmentationOracle) FindByAccountSegmentationAndAccessListKeys(ctx context.Context, customerSegments string, segmentKeys []string) (*model.AccessListSegmentation, error) {

	if len(segmentKeys) == 0 {
		return nil, fmt.Errorf("no keys provided")
	}
	placeholders := make([]string, len(segmentKeys))
	for i := range segmentKeys {
		placeholders[i] = fmt.Sprintf("HEXTORAW(:key%d)", i)
	}
	args := make([]interface{}, 0, len(segmentKeys)+1)
	args = append(args, customerSegments)
	for _, k := range segmentKeys {
		args = append(args, k)
	}
	query := fmt.Sprintf(`SELECT RAWTOHEX(id), RAWTOHEX(access_list_key), RAWTOHEX(segmented_id), enabled, created_at, updated_at, deleted_at FROM ACCESS_LIST_BY_SUPERAPP_ROLE WHERE SEGMENTED_ID = HEXTORAW(:1) AND access_list_key IN (%s)`, strings.Join(placeholders, ","))
	row := q.db.QueryRowContext(ctx, query, args...)
	var seg model.AccessListSegmentation
	err := row.Scan(&seg.ID, &seg.AccessListKey, &seg.SegmentedID, &seg.Enabled, &seg.CreatedAt, &seg.UpdatedAt, &seg.DeletedAt)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			q.logger.Errorf("[AccessListSegmentationOracle][FindByAccountSegmentationAndAccessListKeys] query failed: %v", err)
			return nil, nil
		}
		return nil, err
	}
	return &seg, nil
}

// FindBySegmentationAndServiceID implements [storage.AccessListSegmentationRepository].
func (q *accessListSegmentationOracle) FindBySegmentationAndServiceID(ctx context.Context, segmentationID string, serviceID string) (*model.AccessListSegmentation, error) {
	panic("unimplemented")
}

// FindParentChildRelationship implements [storage.AccessListSegmentationRepository].
func (q *accessListSegmentationOracle) FindParentChildRelationship(ctx context.Context) ([]model.AccessItemRelation, error) {

	query := `SELECT parent_key, child_key FROM ACCESS_ITEMS_RELATION`
	rows, err := q.db.QueryContext(ctx, query)
	if err != nil {
		q.logger.Errorf("[AccessListSegmentation][FindParentChildRelationship] query failed: %v", err)
		return nil, err
	}
	var relations []local_model.AccessItemRelation
	for rows.Next() {
		var rel local_model.AccessItemRelation
		if err := rows.Scan(&rel.ParentKey, &rel.ChildKey); err != nil {
			q.logger.Errorf("[AccessListSegmentation][FindParentChildRelationship] failed to scan row: %v", err)
			return nil, localization.ErrorUnexpectedError
		}
		relations = append(relations, rel)
	}

	return relations, nil
}

// Update implements [storage.AccessListSegmentationRepository].
func (q *accessListSegmentationOracle) Update(ctx context.Context, id string, accessListSegmentation model.AccessListSegmentation) error {

	// update := []string{}
	// if accessListSegmentation.Type != "" {
	// 	update["type"] = accessListSegmentation.Type
	// }
	// if accessListSegmentation.SegmentedID != "" {
	// 	update["segmented_id"] = accessListSegmentation.SegmentedID
	// }
	// if accessListSegmentation.AccessListKey != "" {
	// 	update["access_list_key"] = accessListSegmentation.AccessListKey
	// }
	// update["updated_at"] = time.Now()

	return nil
}
