package access_list_segmentation_oracle

import (
	access_list_segmentation_dto "cbe-super-app-cps-action/internal/constants/dto/access_list_segmentation"
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

// FindAllBySegmentIDorSegmentCode implements [storage.AccessListSegmentationRepository].
func (q *accessListSegmentationOracle) FindAllBySegmentIDorSegmentCode(ctx context.Context, segmentIDorCode string) ([]shared_model.APPAccessList, error) {
	panic("unimplemented")
}

// FindAllBySegmentIDorSegmentCodeAndKeys implements [storage.AccessListSegmentationRepository].
func (q *accessListSegmentationOracle) FindAllBySegmentIDorSegmentCodeAndKeys(ctx context.Context, segmentIDorCode string, keys []string) ([]model.AccessListSegmentation, error) {
	panic("unimplemented")
}

// FindAllWithPagination implements [storage.AccessListSegmentationRepository].
func (q *accessListSegmentationOracle) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.AccessListSegmentation], error) {
	panic("unimplemented")
}

// FindByID implements [storage.AccessListSegmentationRepository].
func (q *accessListSegmentationOracle) FindByID(ctx context.Context, id string) (*model.AccessListSegmentation, error) {
	panic("unimplemented")
}

// FindByIDAndType implements [storage.AccessListSegmentationRepository].
func (q *accessListSegmentationOracle) FindByIDAndType(ctx context.Context, ids string, t string) (*model.AccessListSegmentation, error) {
	panic("unimplemented")
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
	panic("unimplemented")
}

// Update implements [storage.AccessListSegmentationRepository].
func (q *accessListSegmentationOracle) Update(ctx context.Context, id string, accessListSegmentation model.AccessListSegmentation) error {
	panic("unimplemented")
}
