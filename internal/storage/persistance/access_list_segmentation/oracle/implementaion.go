package access_list_segmentation_oracle

import (
	access_list_segmentation_dto "cbe-super-app-cps-action/internal/constants/dto/access_list_segmentation"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	local_model "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/kafka"
	"cbe-super-app-cps-action/internal/storage/persistance/access_list_segmentation/core"
	"context"
	"fmt"

	"strings"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	shared_model "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type accessListSegmentationOracle struct {
	db            DBTX
	cfg           config.VaultConfig
	kafkaProducer *kafka.AccessListSegmentationProducer
	accBlock      storage.AccountBlockRepository
	logger        utils.Logger
}

func NewAccessListSegmentationOracle(db DBTX, cfg config.VaultConfig, kafkaProducer *kafka.AccessListSegmentationProducer, logger utils.Logger) storage.AccessListSegmentationRepositoryOracle {
	return &accessListSegmentationOracle{
		db:            db,
		cfg:           cfg,
		kafkaProducer: kafkaProducer,
		logger:        logger,
	}
}

// BulkDisable implements [storage.AccessListSegmentationRepository].
func (q *accessListSegmentationOracle) BulkDisable(ctx context.Context, req access_list_segmentation_dto.BulkDisableAccessListSegmentationRequest) error {
	q.logger.Infof("[AccessListSegmentation][BulkDisable] bulk disable request: %+v", req)
	var filter string
	var args []interface{}
	var segType string
	if req.SegmentationType == "Account" {
		filter = "segmentation_code = :1 AND access_list_key IN ("
		args = append(args, req.ID)
		for i, key := range req.Keys {
			if i > 0 {
				filter += ","
			}
			filter += fmt.Sprintf(":key%d", i+2)
			args = append(args, key)
		}
		filter += ")"
		segType = "account-segment"
	} else {
		filter = "segmentation_id = :1 AND access_list_key IN ("
		args = append(args, req.ID)
		for i, key := range req.Keys {
			if i > 0 {
				filter += ","
			}
			filter += fmt.Sprintf(":key%d", i+2)
			args = append(args, key)
		}
		filter += ")"
		segType = "block-segment"
	}
	var stmt string
	if req.SegmentationType == "Account" {
		stmt = "DELETE FROM access_list_customer_seg WHERE " + filter

	} else {
		stmt = "DELETE FROM access_list_geo_seg WHERE " + filter

	}
	_, err := q.db.ExecContext(ctx, stmt, args...)
	if err != nil {
		q.logger.Errorf("[AccessListSegmentation][BulkDisable] failed to bulk disable access list segmentation: %v", err)
		return err
	}

	branches := core.GetAllBranches(ctx, req.ID, q.accBlock)

	// Prepare Kafka message
	als := make([]local_model.AccessListSegmentation, len(req.Keys))
	for i, key := range req.Keys {
		doc := local_model.AccessListSegmentation{
			AccessListKey: key,
			Enabled:       req.Enabled,
		}
		doc.SegmentedID = req.ID
		als[i] = doc
	}
	res := map[string]any{
		"docs":     als,
		"type":     segType,
		"branches": branches,
	}
	q.kafkaProducer.PublishMessage(ctx, res, "delete", q.cfg.KafkaCustomerSegmentaionTopic, "bulk delete access-list-segmentation")
	return nil
}

func (q *accessListSegmentationOracle) CreateAccountSegment(ctx context.Context, accessListSegmentation access_list_segmentation_dto.CreateAccessListSegmentationRequest) error {
	n := len(accessListSegmentation.AccessListKeys)
	if n == 0 {
		return nil
	}

	valueStrings := make([]string, 0, n)
	valueArgs := make([]interface{}, 0, n*7)

	docs := make([]local_model.AccessListSegmentation, 0, n)

	for i, key := range accessListSegmentation.AccessListKeys {
		valueStrings = append(valueStrings, fmt.Sprintf("(SYS_GUID(), :access_list_key%d, :segmented_id, :enabled, :created_at, :updated_at, :deleted_at)", i))
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
	}

	stmt := `
	       INSERT INTO ACCESS_LIST_CUSTOMER_SEG (
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

	for i, key := range accessListSegmentation.AccessListKeys {
		valueStrings = append(valueStrings, fmt.Sprintf("(SYS_GUID(), :access_list_key%d, :segmented_id, :type, :enabled, :created, :updated, :deleted)", i))
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
	}

	stmt := `
	       INSERT INTO ACCESS_LIST_GEO_SEG (
		   id, access_list_key, segmented_id, type, enabled, created_at, updated_at, deleted_at
	       ) VALUES ` + strings.Join(valueStrings, ",")

	_, err := q.db.ExecContext(ctx, stmt, valueArgs...)
	if err != nil {
		q.logger.Errorf("[AccessListSegmentation][CreateBlockSegment] failed to insert rows: %v", err)
		return err
	}

	res := map[string]any{
		"docs": docs,
		"type": "block-segment",
	}
	q.kafkaProducer.PublishMessage(ctx, res, "create", q.cfg.KafkaCustomerSegmentaionTopic, "create block-segment")

	return nil
}

// FindAllForBlock implements [storage.AccessListSegmentationRepository].
func (q *accessListSegmentationOracle) FindAllForBlock(ctx context.Context, geographicalID string) ([]shared_model.APPAccessList, error) {

	query := `SELECT RAWTOHEX(g.access_list_key), a.access_list_name, g.enabled
		FROM ACCESS_LIST_GEO_SEG g
		JOIN ACCESS_LIST a ON g.access_list_key = a.id
		WHERE g.segmented_id = HEXTORAW(:1)`
	rows, err := q.db.QueryContext(ctx, query, geographicalID)
	if err != nil {
		q.logger.Errorf("[AccessListSegmentation][FindAllForBlock] query failed: %v", err)
		return nil, err
	}
	defer rows.Close()
	var result []shared_model.APPAccessList
	for rows.Next() {
		var key, accessListName string
		var enabled bool
		err := rows.Scan(&key, &accessListName, &enabled)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				q.logger.Infof("[AccessListSegmentation][FindAllForBlock] no access list segmentation found for geographicalID: %s", geographicalID)
				return nil, localization.ErrorResourceNotFound
			}
			q.logger.Errorf("[AccessListSegmentation][FindAllForBlock] query failed: %v", err)
			return nil, err
		}
		q.logger.Infof("[AccessListSegmentation][FindAllForBlock] found access list segmentation: key=%s, name=%s, enabled=%v", key, accessListName, enabled)
		result = append(result, shared_model.APPAccessList{
			Key:            key,
			AccessListName: accessListName,
			Enabled:        enabled,
		})
	}
	return result, nil
}

// FindAllBySegmentIDorSegmentCode implements [storage.AccessListSegmentationRepository].
func (q *accessListSegmentationOracle) FindAllForAccount(ctx context.Context, customer_seg_id string) ([]shared_model.APPAccessList, error) {

	query := `SELECT RAWTOHEX(g.access_list_key), a.access_list_name, g.enabled
		FROM ACCESS_LIST_CUSTOMER_SEG g
		JOIN ACCESS_LIST a ON g.access_list_key = a.id
		WHERE g.segmented_id = HEXTORAW(:1)`
	rows, err := q.db.QueryContext(ctx, query, customer_seg_id)
	if err != nil {
		q.logger.Errorf("[AccessListSegmentation][FindAllForAccount] query failed: %v", err)
		return nil, err
	}
	defer rows.Close()
	var result []shared_model.APPAccessList
	for rows.Next() {
		var key, accessListName string
		var enabled bool
		err := rows.Scan(&key, &accessListName, &enabled)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				q.logger.Infof("[AccessListSegmentation][FindAllForAccount] no access list segmentation found for customer_seg_id: %s", customer_seg_id)
				return nil, localization.ErrorResourceNotFound
			}
			q.logger.Errorf("[AccessListSegmentation][FindAllForAccount] query failed: %v", err)
			return nil, err
		}
		q.logger.Infof("[AccessListSegmentation][FindAllForAccount] found access list segmentation: key=%s, name=%s, enabled=%v", key, accessListName, enabled)
		result = append(result, shared_model.APPAccessList{
			Key:            key,
			AccessListName: accessListName,
			Enabled:        enabled,
		})
	}
	return result, nil
}

// FindAllBySegmentIDorSegmentCodeAndKeys implements [storage.AccessListSegmentationRepository].
func (q *accessListSegmentationOracle) FindAllByBlockAndKeys(ctx context.Context, segmentIDorCode string, keys []string) ([]model.AccessListSegmentation, error) {
	q.logger.Infof("[AccessListSegmentationOracle][FindAllByBlockAndKeys] called with segmentIDorCode: %s, keys: %v", segmentIDorCode, keys)
	if len(keys) == 0 {
		return nil, fmt.Errorf("no keys provided")
	}
	query := `SELECT RAWTOHEX(id), RAWTOHEX(access_list_key), RAWTOHEX(segmented_id), type, enabled, created_at, updated_at, deleted_at FROM ACCESS_LIST_GEO_SEG WHERE segmented_id = HEXTORAW(:1) AND access_list_key IN (`
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
		err := rows.Scan(&seg.ID, &seg.AccessListKey, &seg.SegmentedID, &seg.Type, &seg.Enabled, &seg.CreatedAt, &seg.UpdatedAt, &seg.DeletedAt)
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
	query := `SELECT RAWTOHEX(id), RAWTOHEX(access_list_key), RAWTOHEX(segmented_id), type, enabled, created_at, updated_at, deleted_at FROM ACCESS_LIST_CUSTOMER_SEG WHERE segmented_id = HEXTORAW(:1) AND access_list_key IN (`
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
		var seg model.AccessListSegmentation
		err := rows.Scan(&seg.ID, &seg.AccessListKey, &seg.SegmentedID, &seg.Type, &seg.Enabled, &seg.CreatedAt, &seg.UpdatedAt, &seg.DeletedAt)
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
	panic("unimplemented")
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
	query := `SELECT RAWTOHEX(id), RAWTOHEX(access_list_key), RAWTOHEX(segmented_id), type, enabled, created_at, updated_at, deleted_at FROM ACCESS_LIST_CUSTOMER_SEG WHERE id = HEXTORAW(:1)`
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
	// Check if a record exists in ACCESS_LIST_GEO_SEG with the given segmented_id and any of the keys

	// Inputs are hex strings, convert to RAW for query, output as hex
	if len(keys) == 0 {
		return nil, fmt.Errorf("no keys provided")
	}
	placeholders := make([]string, len(keys))
	for i := range keys {
		placeholders[i] = fmt.Sprintf("HEXTORAW(:key%d)", i)
	}
	args := make([]interface{}, 0, len(keys)+1)
	args = append(args, id)
	for _, k := range keys {
		args = append(args, k)
	}
	query := fmt.Sprintf(`SELECT RAWTOHEX(id), RAWTOHEX(access_list_key), RAWTOHEX(segmented_id), enabled, created_at, updated_at, deleted_at FROM ACCESS_LIST_GEO_SEG WHERE segmented_id = HEXTORAW(:1) AND access_list_key IN (%s)`, strings.Join(placeholders, ","))
	row := q.db.QueryRowContext(ctx, query, args...)
	var seg model.AccessListSegmentation
	err := row.Scan(&seg.ID, &seg.AccessListKey, &seg.SegmentedID, &seg.Type, &seg.Enabled, &seg.CreatedAt, &seg.UpdatedAt, &seg.DeletedAt)
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
	query := fmt.Sprintf(`SELECT RAWTOHEX(id), RAWTOHEX(access_list_key), RAWTOHEX(segmented_id), enabled, created_at, updated_at, deleted_at FROM ACCESS_LIST_CUSTOMER_SEG WHERE SEGMENTED_ID = HEXTORAW(:1) AND access_list_key IN (%s)`, strings.Join(placeholders, ","))
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

	query := `SELECT RAWTOHEX(parent_key), RAWTOHEX(child_key) FROM ACCESS_ITEM_RELATIONS`
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
