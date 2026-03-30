package services

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

	"github.com/google/uuid"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	shared_constants "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/constants"
	model "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"

	local_util "cbe-super-app-cps-action/pkgs/utils"
)

type ServicesStorage struct {
	cfg           *config.VaultConfig
	db            *sql.DB
	kafkaProducer kafka.ClientOrchestrationProducer
	redis         storage.RedisRepository
	logger        utils.Logger
}

func NewServicesRepository(db *sql.DB, cfg *config.VaultConfig, kafkaProducer kafka.ClientOrchestrationProducer, redis storage.RedisRepository, logger utils.Logger) storage.ServicesRepository {
	return &ServicesStorage{
		cfg:           cfg,
		db:            db,
		kafkaProducer: kafkaProducer,
		redis:         redis,
		logger:        logger,
	}
}

const (
	servicesTable        = "services"
	appAccessListTable   = "app_access_list"
	serviceKeysTable     = "service_keys"
	maxPaginationDefault = 50
)

func generateUUID() string {
	return uuid.New().String()
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

// func parseObjectIDHex(id string) (bson.ObjectID, error) {
// 	oid, err := bson.ObjectIDFromHex(id)
// 	if err != nil {
// 		return bson.NilObjectID, errors.New(localization.ErrorInvalidID.Code)
// 	}
// 	return oid, nil
// }

func (s *ServicesStorage) getServiceCaps(ctx context.Context, serviceID string) ([]model.Cap, error) {
	const q = `
SELECT source, currency, single_cap, minimum_transfer_cap
FROM service_cap
WHERE service_id = :1`

	rows, err := s.db.QueryContext(ctx, q, serviceID)
	if err != nil {
		s.logger.Errorf("[ServicesRepo][getServiceCaps] query caps failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer rows.Close()

	var caps []model.Cap
	for rows.Next() {
		var sourceStr string
		var currency, singleCap, minTransferCap string
		if err := rows.Scan(&sourceStr, &currency, &singleCap, &minTransferCap); err != nil {
			s.logger.Errorf("[ServicesRepo][getServiceCaps] scan failed: %v", err)
			return nil, local_util.HandleDBError(err)
		}

		caps = append(caps, model.Cap{
			Source:             shared_constants.SourceApp(sourceStr),
			Currency:           currency,
			SingleCap:          singleCap,
			MinimumTransferCap: minTransferCap,
		})
	}
	return caps, nil
}

func (s *ServicesStorage) insertService(ctx context.Context, tx *sql.Tx, service *model.Service) error {
	serviceID := generateUUID()
	deletedAtVal := time.Now()
	if service.DeletedAt != nil {
		deletedAtVal = *service.DeletedAt
	}

	// Insert the main service row (caps are stored in service_cap).
	const q = `
INSERT INTO services (
  id,
  service_name,
  service_code,
  service_key,
  minimum_fraud_amount,
  product_gl_account,
  product_gl_account_currency,
  enabled,
  is_deleted,
  created_at,
  last_modified_at,
  deleted_at
) 
VALUES (
  :1,:2,:3,:4,:5,:6,:7,
  :8,:9,:10,:11,:12
)`

	_, err := tx.ExecContext(ctx, q,
		serviceID,
		service.ServiceName,
		service.ServiceCode,
		service.ServiceKey,
		service.MinimumFraudAmount,
		service.ProductGlAccount,
		service.ProductGlAccountCurrency,
		boolToOracleNumber(service.Enabled),
		boolToOracleNumber(service.IsDeleted),
		service.CreatedAt,
		service.LastModifiedAt,
		deletedAtVal,
	)
	if err != nil {
		s.logger.Errorf("[ServicesRepo][insertService] insert failed: %v", err)
		return local_util.HandleDBError(err)
	}

	// Insert cap rows (one row per cap entry).
	if len(service.Cap) == 0 {
		return nil
	}

	const capQ = `
INSERT INTO service_cap (
  id,
  service_id,
  source,
  currency,
  single_cap,
  minimum_transfer_cap
)
VALUES (
  :1,:2,:3,:4,:5,:6
)`

	for _, cap := range service.Cap {
		_, err := tx.ExecContext(ctx, capQ,
			generateUUID(),
			serviceID,
			string(cap.Source),
			cap.Currency,
			cap.SingleCap,
			cap.MinimumTransferCap,
		)
		if err != nil {
			s.logger.Errorf("[ServicesRepo][insertService] insert cap failed: %v", err)
			return local_util.HandleDBError(err)
		}
	}

	return nil
}

func (s *ServicesStorage) updateServiceCaps(ctx context.Context, tx *sql.Tx, serviceID string, caps []model.Cap) error {
	const deleteCapsQ = `DELETE FROM service_cap WHERE service_id = :1`
	if _, err := tx.ExecContext(ctx, deleteCapsQ, serviceID); err != nil {
		s.logger.Errorf("[ServicesRepo][updateServiceCaps] delete caps failed: %v", err)
		return local_util.HandleDBError(err)
	}

	if len(caps) == 0 {
		return nil
	}

	const capQ = `
INSERT INTO service_cap (
  id,
  service_id,
  source,
  currency,
  single_cap,
  minimum_transfer_cap
)
VALUES (
  :1,:2,:3,:4,:5,:6
)`

	for _, cap := range caps {
		if _, err := tx.ExecContext(ctx, capQ,
			generateUUID(),
			serviceID,
			string(cap.Source),
			cap.Currency,
			cap.SingleCap,
			cap.MinimumTransferCap,
		); err != nil {
			s.logger.Errorf("[ServicesRepo][updateServiceCaps] insert cap failed: %v", err)
			return local_util.HandleDBError(err)
		}
	}

	return nil
}

func (s *ServicesStorage) Create(ctx context.Context, service *model.Service) error {
	service.ID = generateUUID()
	if service.CreatedAt.IsZero() {
		service.CreatedAt = time.Now()
	}
	if service.LastModifiedAt.IsZero() {
		service.LastModifiedAt = time.Now()
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Errorf("[ServicesRepo][Create] begin tx failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = tx.Rollback() }()

	if err := s.insertService(ctx, tx, service); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		s.logger.Errorf("[ServicesRepo][Create] commit failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	// Keep existing behavior from the old Mongo implementation.
	_ = s.kafkaProducer.PublishMessage(
		ctx,
		service,
		string(constants.ClientOrchestrationServicesTopic),
		string(constants.ClientOrchestrationServicesTopic),
		"new service created",
	)
	_ = s.redis.Set(ctx, s.cfg.CPSServiceUpdate, service, -1)

	return nil
}

func (s *ServicesStorage) Update(ctx context.Context, id string, service *model.Service) error {
	// if _, err := parseObjectIDHex(id); err != nil {
	// 	return err
	// }

	updatedServiceID := id
	// if service != nil && service.ID != bson.NilObjectID {
	// 	updatedServiceID = service.ID.Hex()
	// }

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Errorf("[ServicesRepo][Update] begin tx failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = tx.Rollback() }()

	// 1) Update main services row.
	const q = `
UPDATE services SET
  service_name = :1,
  service_code = :2,
  service_key = :3,
  minimum_fraud_amount = :4,
  product_gl_account = :5,
  product_gl_account_currency = :6,
  enabled = :7,
  is_deleted = :8,
  last_modified_at = SYSTIMESTAMP
WHERE id = :9 AND is_deleted = 0`

	res, err := tx.ExecContext(ctx, q,
		service.ServiceName,
		service.ServiceCode,
		service.ServiceKey,
		service.MinimumFraudAmount,
		service.ProductGlAccount,
		service.ProductGlAccountCurrency,
		boolToOracleNumber(service.Enabled),
		boolToOracleNumber(service.IsDeleted),
		updatedServiceID,
	)
	if err != nil {
		s.logger.Errorf("[ServicesRepo][Update] update service failed: %v", err)
		return local_util.HandleDBError(err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New(localization.ErrorServiceNotFound.Code)
	}

	// 2) Replace caps.
	if err := s.updateServiceCaps(ctx, tx, updatedServiceID, service.Cap); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		s.logger.Errorf("[ServicesRepo][Update] commit failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	_ = s.kafkaProducer.PublishMessage(
		ctx,
		service,
		string(constants.ClientOrchestrationServicesTopic),
		string(constants.ClientOrchestrationServicesTopic),
		"service authorized and updated",
	)
	_ = s.redis.Set(ctx, s.cfg.CPSServiceUpdate, service, -1)

	return nil
}

func (s *ServicesStorage) Delete(ctx context.Context, id string) error {
	// if _, err := parseObjectIDHex(id); err != nil {
	// 	return err
	// }

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Errorf("[ServicesRepo][Delete] begin tx failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = tx.Rollback() }()

	const deleteCapsQ = `DELETE FROM service_cap WHERE service_id = :1`
	if _, err := tx.ExecContext(ctx, deleteCapsQ, id); err != nil {
		s.logger.Errorf("[ServicesRepo][Delete] delete caps failed: %v", err)
		return local_util.HandleDBError(err)
	}

	const q = `
UPDATE services SET
  is_deleted = 1,
  deleted_at = SYSTIMESTAMP,
  last_modified_at = SYSTIMESTAMP
WHERE id = :1`

	res, err := tx.ExecContext(ctx, q, id)
	if err != nil {
		s.logger.Errorf("[ServicesRepo][Delete] delete service failed: %v", err)
		return local_util.HandleDBError(err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New(localization.ErrorServiceNotFound.Code)
	}

	if err := tx.Commit(); err != nil {
		s.logger.Errorf("[ServicesRepo][Delete] commit failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (s *ServicesStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	// if _, err := parseObjectIDHex(id); err != nil {
	// 	return err
	// }

	const q = `
UPDATE services
SET
  enabled = :1,
  last_modified_at = SYSTIMESTAMP
WHERE id = :2 AND is_deleted = 0`

	res, err := s.db.ExecContext(ctx, q, boolToOracleNumber(enable), id)
	if err != nil {
		s.logger.Errorf("[ServicesRepo][EnableOrDisable] update failed: %v", err)
		return local_util.HandleDBError(err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New(localization.ErrorServiceNotFound.Code)
	}

	return nil
}

func (s *ServicesStorage) FindByID(ctx context.Context, id string) (*model.Service, error) {
	// if _, err := parseObjectIDHex(id); err != nil {
	// 	return nil, err
	// }

	const q = `
SELECT
  id,
  service_name,
  service_code,
  service_key,
  minimum_fraud_amount,
  product_gl_account,
  product_gl_account_currency,
  enabled,
  is_deleted,
  created_at,
  last_modified_at,
  deleted_at
FROM services
WHERE id = :1 AND is_deleted = 0`

	var serviceID string
	var svc model.Service
	var deletedAt sql.NullTime

	err := s.db.QueryRowContext(ctx, q, id).Scan(
		&serviceID,
		&svc.ServiceName,
		&svc.ServiceCode,
		&svc.ServiceKey,
		&svc.MinimumFraudAmount,
		&svc.ProductGlAccount,
		&svc.ProductGlAccountCurrency,
		&svc.Enabled,
		&svc.IsDeleted,
		&svc.CreatedAt,
		&svc.LastModifiedAt,
		&deletedAt,
	)
	if err != nil {
		s.logger.Errorf("[ServicesRepo][FindByID] query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	// oid, err := bson.ObjectIDFromHex(serviceID)
	// if err != nil {
	// 	return nil, errors.New(localization.ErrorInvalidID.Code)
	// }
	svc.ID = serviceID
	if deletedAt.Valid {
		svc.DeletedAt = &deletedAt.Time
	}

	caps, err := s.getServiceCaps(ctx, id)
	if err != nil {
		return nil, err
	}
	svc.Cap = caps

	return &svc, nil
}

func (s *ServicesStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.Service], error) {
	limit := int64(maxPaginationDefault)
	page := int64(1)
	if filterParam.PerPage > 0 {
		limit = int64(filterParam.PerPage)
	}
	if filterParam.Page > 0 {
		page = int64(filterParam.Page)
	}
	offset := (page - 1) * limit

	clauses := []string{"s.is_deleted = 0"}
	var args []interface{}

	search := strings.TrimSpace(filterParam.Search)
	if search != "" {
		clauses = append(clauses,
			`(
				LOWER(s.service_name) LIKE '%' || LOWER(:search) || '%'
				OR LOWER(s.service_code) LIKE '%' || LOWER(:search) || '%'
				OR LOWER(s.service_key) LIKE '%' || LOWER(:search) || '%'
			)`,
		)
		args = append(args, sql.Named("search", search))
	}

	if filterParam.Filters != nil {
		if v, ok := filterParam.Filters["enabled"]; ok {
			if b, ok2 := parseBoolFilter(v); ok2 {
				clauses = append(clauses, "s.enabled = :enabled")
				args = append(args, sql.Named("enabled", boolToOracleNumber(b)))
			}
		}

		if v, ok := filterParam.Filters["service_name"]; ok {
			if sName, ok2 := v.(string); ok2 && strings.TrimSpace(sName) != "" {
				clauses = append(clauses, "LOWER(s.service_name) = LOWER(:service_name)")
				args = append(args, sql.Named("service_name", sName))
			}
		}
		if v, ok := filterParam.Filters["service_code"]; ok {
			if sCode, ok2 := v.(string); ok2 && strings.TrimSpace(sCode) != "" {
				clauses = append(clauses, "LOWER(s.service_code) = LOWER(:service_code)")
				args = append(args, sql.Named("service_code", sCode))
			}
		}
		if v, ok := filterParam.Filters["service_key"]; ok {
			if sKey, ok2 := v.(string); ok2 && strings.TrimSpace(sKey) != "" {
				clauses = append(clauses, "LOWER(s.service_key) = LOWER(:service_key)")
				args = append(args, sql.Named("service_key", sKey))
			}
		}
	}

	where := strings.Join(clauses, " AND ")

	countQ := fmt.Sprintf(`SELECT COUNT(*) FROM %s s WHERE %s`, servicesTable, where)
	var total int64
	if err := s.db.QueryRowContext(ctx, countQ, args...).Scan(&total); err != nil {
		s.logger.Errorf("[ServicesRepo][FindAllWithPagination] count failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	listQ := fmt.Sprintf(`
SELECT
  id,
  service_name,
  service_code,
  service_key,
  minimum_fraud_amount,
  product_gl_account,
  product_gl_account_currency,
  enabled,
  is_deleted,
  created_at,
  last_modified_at,
  deleted_at
FROM %s s
WHERE %s
ORDER BY created_at DESC
OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`, servicesTable, where)

	listArgs := append(args, sql.Named("offset", offset), sql.Named("limit", limit))
	rows, err := s.db.QueryContext(ctx, listQ, listArgs...)
	if err != nil {
		s.logger.Errorf("[ServicesRepo][FindAllWithPagination] list query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer rows.Close()

	var list []model.Service
	for rows.Next() {
		var serviceID string
		var svc model.Service
		var deletedAt sql.NullTime

		if err := rows.Scan(
			&serviceID,
			&svc.ServiceName,
			&svc.ServiceCode,
			&svc.ServiceKey,
			&svc.MinimumFraudAmount,
			&svc.ProductGlAccount,
			&svc.ProductGlAccountCurrency,
			&svc.Enabled,
			&svc.IsDeleted,
			&svc.CreatedAt,
			&svc.LastModifiedAt,
			&deletedAt,
		); err != nil {
			s.logger.Errorf("[ServicesRepo][FindAllWithPagination] scan failed: %v", err)
			return nil, local_util.HandleDBError(err)
		}

		svc.ID = serviceID
		if deletedAt.Valid {
			svc.DeletedAt = &deletedAt.Time
		}

		caps, err := s.getServiceCaps(ctx, serviceID)
		if err != nil {
			return nil, err
		}
		svc.Cap = caps

		list = append(list, svc)
	}

	meta := local_util.BuildPaginationMeta(total, int(page), int(limit))
	return &types.PaginatedResponse[[]model.Service]{Data: list, Meta: meta}, nil
}

func (s *ServicesStorage) CheckServiceExistence(ctx context.Context, serviceCode, serviceKey, serviceName string) (bool, error) {
	clauses := make([]string, 0, 3)
	args := make([]interface{}, 0, 3)

	if strings.TrimSpace(serviceCode) != "" {
		clauses = append(clauses, "service_code = :service_code")
		args = append(args, sql.Named("service_code", serviceCode))
	}
	if strings.TrimSpace(serviceKey) != "" {
		clauses = append(clauses, "service_key = :service_key")
		args = append(args, sql.Named("service_key", serviceKey))
	}
	if strings.TrimSpace(serviceName) != "" {
		clauses = append(clauses, "service_name = :service_name")
		args = append(args, sql.Named("service_name", serviceName))
	}
	if len(clauses) == 0 {
		return false, nil
	}

	where := strings.Join(clauses, " OR ")
	q := fmt.Sprintf(`SELECT COUNT(*) FROM services WHERE is_deleted = 0 AND (%s)`, where)

	var count int64
	if err := s.db.QueryRowContext(ctx, q, args...).Scan(&count); err != nil {
		s.logger.Errorf("[ServicesRepo][CheckServiceExistence] count failed: %v", err)
		return false, local_util.HandleDBError(err)
	}

	return count > 0, nil
}

func (s *ServicesStorage) FindAllServiceListWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.ServiceList], error) {
	limit := int64(maxPaginationDefault)
	page := int64(1)
	if filterParam.PerPage > 0 {
		limit = int64(filterParam.PerPage)
	}
	if filterParam.Page > 0 {
		page = int64(filterParam.Page)
	}
	offset := (page - 1) * limit

	clauses := []string{"is_enabled = is_enabled"} // placeholder to keep join logic simple
	clauses = []string{}
	args := []interface{}{}

	search := strings.TrimSpace(filterParam.Search)
	if search != "" {
		clauses = append(clauses,
			`(
				LOWER(service_name) LIKE '%' || LOWER(:search) || '%'
				OR LOWER(service_key) LIKE '%' || LOWER(:search) || '%'
			)`,
		)
		args = append(args, sql.Named("search", search))
	}

	if filterParam.Filters != nil {
		if v, ok := filterParam.Filters["enabled"]; ok {
			if b, ok2 := parseBoolFilter(v); ok2 {
				clauses = append(clauses, "is_enabled = :enabled")
				args = append(args, sql.Named("enabled", boolToOracleNumber(b)))
			}
		}
		if v, ok := filterParam.Filters["service_name"]; ok {
			if sName, ok2 := v.(string); ok2 && strings.TrimSpace(sName) != "" {
				clauses = append(clauses, "LOWER(service_name) = LOWER(:service_name)")
				args = append(args, sql.Named("service_name", sName))
			}
		}
		if v, ok := filterParam.Filters["service_key"]; ok {
			if sKey, ok2 := v.(string); ok2 && strings.TrimSpace(sKey) != "" {
				clauses = append(clauses, "LOWER(service_key) = LOWER(:service_key)")
				args = append(args, sql.Named("service_key", sKey))
			}
		}
	}

	where := "1=1"
	if len(clauses) > 0 {
		where = strings.Join(clauses, " AND ")
	}

	countQ := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE %s`, serviceKeysTable, where)
	var total int64
	if err := s.db.QueryRowContext(ctx, countQ, args...).Scan(&total); err != nil {
		s.logger.Errorf("[ServicesRepo][FindAllServiceListWithPagination] count failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	listQ := fmt.Sprintf(`
SELECT
  id,
  service_name,
  service_key,
  is_enabled,
  created_at,
  last_modified_at
FROM %s
WHERE %s
ORDER BY created_at DESC
OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`, serviceKeysTable, where)

	listArgs := append(args, sql.Named("offset", offset), sql.Named("limit", limit))
	rows, err := s.db.QueryContext(ctx, listQ, listArgs...)
	if err != nil {
		s.logger.Errorf("[ServicesRepo][FindAllServiceListWithPagination] list query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer rows.Close()

	var list []model.ServiceList
	for rows.Next() {
		var listID string
		var item model.ServiceList
		if err := rows.Scan(
			&listID,
			&item.ServiceName,
			&item.ServiceKey,
			&item.IsEnabled,
			&item.CreatedAt,
			&item.LastModifiedAt,
		); err != nil {
			s.logger.Errorf("[ServicesRepo][FindAllServiceListWithPagination] scan failed: %v", err)
			return nil, local_util.HandleDBError(err)
		}

		// oid, err := bson.ObjectIDFromHex(listID)
		// if err != nil {
		// 	return nil, errors.New(localization.ErrorInvalidID.Code)
		// }
		item.ID = listID

		list = append(list, item)
	}

	meta := local_util.BuildPaginationMeta(total, int(page), int(limit))
	return &types.PaginatedResponse[[]model.ServiceList]{Data: list, Meta: meta}, nil
}

func (s *ServicesStorage) FindServiceListByID(ctx context.Context, id string) (*model.ServiceList, error) {
	// if _, err := parseObjectIDHex(id); err != nil {
	// 	return nil, err
	// }

	const q = `
SELECT id, service_name, service_key, is_enabled, created_at, last_modified_at
FROM service_keys
WHERE id = :1`

	var item model.ServiceList
	var listID string
	err := s.db.QueryRowContext(ctx, q, id).Scan(
		&listID,
		&item.ServiceName,
		&item.ServiceKey,
		&item.IsEnabled,
		&item.CreatedAt,
		&item.LastModifiedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(localization.ErrorServiceListNotFound.Code)
		}
		return nil, local_util.HandleDBError(err)
	}

	// oid, err := bson.ObjectIDFromHex(listID)
	// if err != nil {
	// 	return nil, errors.New(localization.ErrorInvalidID.Code)
	// }
	item.ID = listID

	return &item, nil
}

func (s *ServicesStorage) FindServiceListByNameOrKey(ctx context.Context, name, key string) (*model.ServiceList, error) {
	conds := make([]string, 0, 2)
	args := make([]interface{}, 0, 2)

	if strings.TrimSpace(name) != "" {
		conds = append(conds, "LOWER(service_name) LIKE '%' || LOWER(:name) || '%'")
		args = append(args, sql.Named("name", name))
	}
	if strings.TrimSpace(key) != "" {
		conds = append(conds, "service_key = :key")
		args = append(args, sql.Named("key", key))
	}

	if len(conds) == 0 {
		return nil, errors.New(localization.ErrorNoDataProvided.Code)
	}

	where := strings.Join(conds, " OR ")
	query := fmt.Sprintf(`
SELECT id, service_name, service_key, is_enabled, created_at, last_modified_at
FROM %s
WHERE (%s)
FETCH FIRST 1 ROWS ONLY`, serviceKeysTable, where)

	var item model.ServiceList
	var listID string
	err := s.db.QueryRowContext(ctx, query, args...).Scan(
		&listID,
		&item.ServiceName,
		&item.ServiceKey,
		&item.IsEnabled,
		&item.CreatedAt,
		&item.LastModifiedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(localization.ErrorServiceListNotFound.Code)
		}
		return nil, local_util.HandleDBError(err)
	}

	// oid, err := bson.ObjectIDFromHex(listID)
	// if err != nil {
	// 	return nil, errors.New(localization.ErrorInvalidID.Code)
	// }
	item.ID = listID
	return &item, nil
}

func (s *ServicesStorage) CreateServiceList(ctx context.Context, serviceList *model.ServiceList) error {
	serviceList.ID = generateUUID()
	now := time.Now()
	if serviceList.CreatedAt.IsZero() {
		serviceList.CreatedAt = now
	}
	if serviceList.LastModifiedAt.IsZero() {
		serviceList.LastModifiedAt = now
	}

	serviceList.IsEnabled = true

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Errorf("[ServicesRepo][CreateServiceList] begin tx failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = tx.Rollback() }()

	const q = `
INSERT INTO service_keys (
  id,
  service_name,
  service_key,
  is_enabled,
  created_at,
  last_modified_at
)
VALUES (
  :1,:2,:3,:4,:5,:6
)`

	if _, err := tx.ExecContext(ctx, q,
		serviceList.ID,
		serviceList.ServiceName,
		serviceList.ServiceKey,
		boolToOracleNumber(serviceList.IsEnabled),
		serviceList.CreatedAt,
		serviceList.LastModifiedAt,
	); err != nil {
		s.logger.Errorf("[ServicesRepo][CreateServiceList] insert failed: %v", err)
		return local_util.HandleDBError(err)
	}

	// Synchronize app_access_list with this service list
	const accessQ = `
INSERT INTO app_access_list (
  id,
  key,
  enabled,
  access_list_name,
  ussd_enabled,
  created_at,
  last_modified_at
)
VALUES (
  :1,:2,:3,:4,:5,:6,:7
)`
	accessID := bson.NewObjectID().Hex()
	if _, err := tx.ExecContext(ctx, accessQ,
		accessID,
		serviceList.ServiceKey,
		boolToOracleNumber(serviceList.IsEnabled),
		serviceList.ServiceName,
		0,
		serviceList.CreatedAt,
		serviceList.LastModifiedAt,
	); err != nil {
		s.logger.Errorf("[ServicesRepo][CreateServiceList] insert access_list failed: %v", err)
		return local_util.HandleDBError(err)
	}

	if err := tx.Commit(); err != nil {
		s.logger.Errorf("[ServicesRepo][CreateServiceList] commit failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	_ = s.kafkaProducer.PublishMessage(
		ctx,
		serviceList,
		string(constants.ClientOrchestrationServicesTopic),
		string(constants.ClientOrchestrationServicesTopic),
		"new service list created",
	)
	_ = s.redis.Set(ctx, s.cfg.CPSServiceUpdate, serviceList, -1)

	return nil
}

func (s *ServicesStorage) UpdateServiceList(ctx context.Context, id, serviceKey string, serviceList *model.ServiceList) error {
	// if _, err := parseObjectIDHex(id); err != nil {
	// 	return err
	// }

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Errorf("[ServicesRepo][UpdateServiceList] begin tx failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = tx.Rollback() }()
	const q = `
UPDATE service_keys
SET
  service_name = :1,
  service_key = :2,
  last_modified_at = SYSTIMESTAMP
WHERE id = :3 AND service_key = :4`

	res, err := tx.ExecContext(ctx, q,
		serviceList.ServiceName,
		serviceList.ServiceKey,
		id,
		serviceKey,
	)
	if err != nil {
		s.logger.Errorf("[ServicesRepo][UpdateServiceList] update failed: %v", err)
		return local_util.HandleDBError(err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New(localization.ErrorServiceListNotFound.Code)
	}

	// Synchronize Oracle app_access_list with this update.
	// Old Mongo logic updated:
	// - access_list_name -> serviceList.ServiceName
	// - key -> serviceList.ServiceKey
	const accessUpdateQ = `
UPDATE app_access_list
SET
  key = :1,
  access_list_name = :2,
  last_modified_at = SYSTIMESTAMP
WHERE key = :3`
	if _, err := tx.ExecContext(ctx, accessUpdateQ,
		serviceList.ServiceKey,
		serviceList.ServiceName,
		serviceKey,
	); err != nil {
		s.logger.Errorf("[ServicesRepo][UpdateServiceList] update access_list failed: %v", err)
		return local_util.HandleDBError(err)
	}

	if err := tx.Commit(); err != nil {
		s.logger.Errorf("[ServicesRepo][UpdateServiceList] commit failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (s *ServicesStorage) EnableOrDisableServiceList(ctx context.Context, id, serviceKey string, enable bool) error {
	// if _, err := parseObjectIDHex(id); err != nil {
	// 	return err
	// }

	const q = `
UPDATE service_keys
SET
  is_enabled = :1,
  last_modified_at = SYSTIMESTAMP
WHERE id = :2 AND service_key = :3`

	res, err := s.db.ExecContext(ctx, q, boolToOracleNumber(enable), id, serviceKey)
	if err != nil {
		s.logger.Errorf("[ServicesRepo][EnableOrDisableServiceList] update failed: %v", err)
		return local_util.HandleDBError(err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New(localization.ErrorServiceListNotFound.Code)
	}

	// Synchronize Oracle app_access_list enable/disable with this toggle.
	const accessEnableQ = `
UPDATE app_access_list
SET
  enabled = :1,
  last_modified_at = SYSTIMESTAMP
WHERE key = :2`
	if _, err := s.db.ExecContext(ctx, accessEnableQ, boolToOracleNumber(enable), serviceKey); err != nil {
		s.logger.Errorf("[ServicesRepo][EnableOrDisableServiceList] update access_list failed: %v", err)
		return local_util.HandleDBError(err)
	}

	return nil
}
