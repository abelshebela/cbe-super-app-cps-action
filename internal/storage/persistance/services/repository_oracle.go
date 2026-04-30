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
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/kafka"

	service_dto "cbe-super-app-cps-action/internal/constants/dto/services"

	"github.com/godror/godror"
	"github.com/hugokessem/coreio/core"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

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
	servicesTable   = "services"
	accessListTable = "access_lists"
	// serviceKeysTable     = "service_keys"
	maxPaginationDefault = 50
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

func (s *ServicesStorage) getServiceCaps(ctx context.Context, serviceID string) ([]imodel.Cap, error) {
	const q = `
SELECT source, currency, single_cap, minimum_transfer_cap
FROM service_cap
WHERE service_id = HEXTORAW(:1)`

	rows, err := s.db.QueryContext(ctx, q, serviceID)
	if err != nil {
		s.logger.Errorf("[ServicesRepo][getServiceCaps] query caps failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer rows.Close()

	var caps []imodel.Cap
	for rows.Next() {
		var sourceStr string
		var currency, singleCap, minTransferCap string
		if err := rows.Scan(&sourceStr, &currency, &singleCap, &minTransferCap); err != nil {
			s.logger.Errorf("[ServicesRepo][getServiceCaps] scan failed: %v", err)
			return nil, local_util.HandleDBError(err)
		}

		caps = append(caps, imodel.Cap{
			Source:             constants.SourceApp(sourceStr),
			Currency:           currency,
			SingleCap:          singleCap,
			MinimumTransferCap: minTransferCap,
		})
	}
	return caps, nil
}

func (s *ServicesStorage) insertService(ctx context.Context, tx *sql.Tx, accountDetail core.AccountLookupResult, service *imodel.Service) (string, error) {
	var serviceID string
	if strings.TrimSpace(service.ServiceKeyId) == "" {
		return "", errors.New(localization.ErrorInvalidID.Code)
	}

	const checkQ = `
		SELECT RAWTOHEX(id)
		FROM accounts
		WHERE account_number = :1
	`

	const bankQ = `
		SELECT RAWTOHEX(id)
		FROM banks
		WHERE is_cbe    = 1
		  AND is_enabled = 1
		  AND is_deleted = 0
		  AND ROWNUM     = 1
	`

	const insertAccountQ = `
		INSERT INTO accounts (
			bank_id,
			account_holder_name,
			account_number,
			account_currency,
			account_type,
			account_branch,
			customer_number,
			created_at,
			last_modified_at
		) VALUES (
			HEXTORAW(:bank_id),
			:customer_name,
			:account_number,
			:currency,
			:account_type,
			:branch,
			:customer_number,
			:created_at,
			:modified_at
		)
		RETURNING RAWTOHEX(id) INTO :id
	`

	var accountID string

	// ── Step 1: check if the GL account already exists ──────────────────────
	s.logger.Debugf("[insertService] checking if GL account %s exists", service.ProductGlAccount)

	err := tx.QueryRowContext(ctx, checkQ, service.ProductGlAccount).Scan(&accountID)
	switch {
	case err == nil:
		// Row found — skip insert, fall through to service insertion.
		s.logger.Infof(
			"[insertService] GL account %s already exists (id=%s), skipping account insert",
			service.ProductGlAccount, accountID,
		)

	case errors.Is(err, sql.ErrNoRows):
		// ── Step 2: resolve the CBE bank ID ─────────────────────────────────
		s.logger.Debugf("[insertService] GL account not found, resolving CBE bank ID")

		var bankID string
		err = tx.QueryRowContext(ctx, bankQ).Scan(&bankID)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			s.logger.Errorf("[insertService] no active CBE bank found")
			return "", errors.New("cbe bank not found")

		case err != nil:
			s.logger.Errorf("[insertService] failed to fetch CBE bank: %v", err)
			return "", local_util.HandleDBError(err)
		}

		// ── Step 3: insert the new GL account ───────────────────────────────
		s.logger.Debugf("[insertService] inserting GL account %s with bank_id %s", service.ProductGlAccount, bankID)

		_, err = tx.ExecContext(ctx, insertAccountQ,
			sql.Named("bank_id", bankID),
			sql.Named("account_number", service.ProductGlAccount),
			sql.Named("customer_name", accountDetail.Detail.CustomerName),
			sql.Named("currency", service.ProductGlAccountCurrency),
			sql.Named("account_type", accountDetail.Detail.AccountType),
			sql.Named("branch", accountDetail.Detail.BranchCode),
			sql.Named("customer_number", accountDetail.Detail.CustomerID),
			sql.Named("created_at", service.CreatedAt),
			sql.Named("modified_at", service.LastModifiedAt),
			sql.Named("id", sql.Out{Dest: &accountID}),
		)
		if err != nil {
			if oraErr, ok := godror.AsOraErr(err); ok {
				switch oraErr.Code() {
				case 1: // ORA-00001: unique constraint — race condition between SELECT and INSERT
					s.logger.Warnf(
						"[insertService] concurrent insert detected for GL account %s, fetching existing id",
						service.ProductGlAccount,
					)
					fetchErr := tx.QueryRowContext(ctx, checkQ, service.ProductGlAccount).Scan(&accountID)
					if fetchErr != nil {
						s.logger.Errorf("[insertService] fallback fetch after race failed: %v", fetchErr)
						return "", local_util.HandleDBError(fetchErr)
					}
					// accountID is now populated — fall through to service insertion

				case 2291: // ORA-02291: FK violation — bank_id not in banks
					s.logger.Errorf("[insertService] bank_id %s not found in banks table", bankID)
					return "", errors.New("bank id not found")

				default:
					s.logger.Errorf("[insertService] account insert failed (ORA-%05d): %v", oraErr.Code(), err)
					return "", local_util.HandleDBError(err)
				}
			} else {
				s.logger.Errorf("[insertService] account insert failed: %v", err)
				return "", local_util.HandleDBError(err)
			}
		}

		s.logger.Infof(
			"[insertService] GL account %s created (id=%s)",
			service.ProductGlAccount, accountID,
		)

	default:
		s.logger.Errorf("[insertService] failed to check existing GL account: %v", err)
		return "", local_util.HandleDBError(err)
	}

	// ── Step 4: insert into services ────────────────────────────────────────
	s.logger.Debugf("[insertService] inserting service with ServiceKeyId %s and GL account %s", service.ServiceKeyId, service.ProductGlAccount)

	const insertServiceQ = `
		INSERT INTO services (
			access_list_id,
			service_code,
			minimum_fraud_amount,
			product_gl_account_number,
			product_gl_account_currency,
			created_at,
			last_modified_at
		)
		VALUES (
			HEXTORAW(:1),:2,:3,:4,:5,:6,:7
		)
		RETURNING RAWTOHEX(id) INTO :8
	`

	if _, err := tx.ExecContext(ctx, insertServiceQ,
		service.ServiceKeyId,
		service.ServiceCode,
		service.MinimumFraudAmount,
		service.ProductGlAccount,
		service.ProductGlAccountCurrency,
		service.CreatedAt,
		service.LastModifiedAt,
		sql.Out{Dest: &serviceID},
	); err != nil {
		s.logger.Errorf("[insertService] insert service failed: %v", err)

		if oraErr, ok := godror.AsOraErr(err); ok {
			if oraErr.Code() == 2291 {
				return "", errors.New(localization.ErrorAccountNumberNotFound.Code)
			}
		}

		return "", local_util.HandleDBError(err)
	}

	// ── Step 5: insert cap rows ──────────────────────────────────────────────
	if len(service.Cap) == 0 {
		return serviceID, nil
	}

	const insertCapQ = `
		INSERT INTO service_cap (
			service_id,
			source,
			currency,
			single_cap,
			minimum_transfer_cap
		)
		VALUES (
			HEXTORAW(:1),:2,:3,:4,:5
		)
	`

	for _, cap := range service.Cap {
		if _, err := tx.ExecContext(ctx, insertCapQ,
			serviceID,
			string(cap.Source),
			cap.Currency,
			cap.SingleCap,
			cap.MinimumTransferCap,
		); err != nil {
			s.logger.Errorf("[insertService] insert cap failed for service %s: %v", serviceID, err)
			return "", local_util.HandleDBError(err)
		}
	}

	return serviceID, nil
}

func (s *ServicesStorage) updateServiceCaps(ctx context.Context, tx *sql.Tx, serviceID string, caps []imodel.Cap) error {
	const deleteCapsQ = `DELETE FROM service_cap WHERE service_id = HEXTORAW(:1)`
	if _, err := tx.ExecContext(ctx, deleteCapsQ, serviceID); err != nil {
		s.logger.Errorf("[ServicesRepo][updateServiceCaps] delete caps failed: %v", err)
		return local_util.HandleDBError(err)
	}

	if len(caps) == 0 {
		return nil
	}

	const capQ = `
INSERT INTO service_cap (
  service_id,
  source,
  currency,
  single_cap,
  minimum_transfer_cap
)
VALUES (
  HEXTORAW(:1),:2,:3,:4,:5
)`

	for _, cap := range caps {
		if _, err := tx.ExecContext(ctx, capQ,
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

func (s *ServicesStorage) Create(ctx context.Context, accountDetail core.AccountLookupResult, service *imodel.Service) error {
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

	s.logger.Debugf("Creating service with ServiceKeyId %s and ProductGlAccount %s", service.ServiceKeyId, service.ProductGlAccount)
	serviceID, err := s.insertService(ctx, tx, accountDetail, service)
	if err != nil {
		s.logger.Errorf("[ServicesRepo][Create] insert failed: %v", err)
		return err
	}
	service.ID = serviceID

	if err := tx.Commit(); err != nil {
		s.logger.Errorf("[ServicesRepo][Create] commit failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

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

func (s *ServicesStorage) Update(ctx context.Context, id string, service *imodel.Service) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Errorf("[ServicesRepo][Update] begin tx failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = tx.Rollback() }()
	serviceKeyID := strings.TrimSpace(service.ServiceKeyId)

	// 	// Ensure service exists before inserting/updating service_keys.
	// 	const svcExistsQ = `
	// SELECT 1
	// FROM services
	// WHERE id = HEXTORAW(:1)
	// FETCH FIRST 1 ROWS ONLY`
	// 	var exists int
	// 	if err := tx.QueryRowContext(ctx, svcExistsQ, id).Scan(&exists); err != nil {
	// 		if errors.Is(err, sql.ErrNoRows) {
	// 			return errors.New(localization.ErrorServiceNotFound.Code)
	// 		}
	// 		s.logger.Errorf("[ServicesRepo][Update] service exists check failed: %v", err)
	// 		return local_util.HandleDBError(err)
	// 	}

	// 	if serviceKeyID == "" {
	// 		return errors.New(localization.ErrorInvalidID.Code)
	// 	}

	// 	// Sync enabled/disabled state on the related service_keys row.
	// 	const qKey = `
	// UPDATE service_keys
	// SET
	//   is_enabled = :1,
	//   is_deleted = :2,
	//   last_modified_at = SYSTIMESTAMP,
	//   deleted_at = CASE WHEN :2 = 1 THEN SYSTIMESTAMP ELSE NULL END
	// WHERE id = HEXTORAW(:3)`

	// 	resKey, err := tx.ExecContext(ctx, qKey,
	// 		boolToOracleNumber(service.Enabled),
	// 		boolToOracleNumber(service.IsDeleted),
	// 		serviceKeyID,
	// 	)
	// 	if err != nil {
	// 		s.logger.Errorf("[ServicesRepo][Update] update service_keys failed: %v", err)
	// 		return local_util.HandleDBError(err)
	// 	}
	// 	rowsKey, _ := resKey.RowsAffected()
	// 	if rowsKey == 0 {
	// 		return errors.New(localization.ErrorServiceNotFound.Code)
	// 	}

	// 1) Update services row fields.
	s.logger.Debugf("Updating service with ID %s and ServiceKeyId %s", id, serviceKeyID)
	const q = `
UPDATE services
SET
  access_list_id = HEXTORAW(:1),
  service_code = :2,
  minimum_fraud_amount = :3,
  product_gl_account_number = :4,
  product_gl_account_currency = :5,
  last_modified_at = SYSTIMESTAMP
WHERE id = HEXTORAW(:6)`

	res, err := tx.ExecContext(ctx, q,
		serviceKeyID,
		service.ServiceCode,
		service.MinimumFraudAmount,
		service.ProductGlAccount,
		service.ProductGlAccountCurrency,
		id,
	)
	if err != nil {
		s.logger.Errorf("[ServicesRepo][Update] update services failed: %v", err)
		return local_util.HandleDBError(err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New(localization.ErrorServiceNotFound.Code)
	}

	// 2) Replace caps.
	if err := s.updateServiceCaps(ctx, tx, id, service.Cap); err != nil {
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

func (s *ServicesStorage) Delete(ctx context.Context, serviceID, accessListID string) error {
	var id string
	if serviceID != "" {
		id = serviceID
	} else {
		srv, err := s.FindServiceByAccessListID(ctx, accessListID)
		if err != nil && err.Error() != localization.ErrorServiceNotFound.Code && err.Error() != localization.ErrorResourceNotFound.Code {
			return err
		}
		if srv != nil && srv.ID != "" {
			id = srv.ID
		}
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Errorf("[ServicesRepo][Delete] begin tx failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = tx.Rollback() }()

	if id != "" {
		const updateCapsQ = `
UPDATE service_cap
SET
  is_deleted = 1,
  deleted_at = SYSTIMESTAMP,
  last_modified_at = SYSTIMESTAMP
WHERE service_id = HEXTORAW(:1)
  AND is_deleted = 0`
		if _, err := tx.ExecContext(ctx, updateCapsQ, id); err != nil {
			s.logger.Errorf("[ServicesRepo][Delete] update service_cap failed: %v", err)
			return local_util.HandleDBError(err)
		}

		const updateServiceQ = `
UPDATE services
SET
  is_deleted = 1,
  deleted_at = SYSTIMESTAMP,
  last_modified_at = SYSTIMESTAMP
WHERE id = HEXTORAW(:1)
  AND is_deleted = 0`
		res, err := tx.ExecContext(ctx, updateServiceQ, id)
		if err != nil {
			s.logger.Errorf("[ServicesRepo][Delete] update services failed: %v", err)
			return local_util.HandleDBError(err)
		}

		rows, _ := res.RowsAffected()
		if rows == 0 {
			return errors.New(localization.ErrorServiceNotFound.Code)
		}

		const q = `
	UPDATE access_lists
	SET
	  is_deleted = 1,
	  deleted_at = SYSTIMESTAMP,
	  last_modified_at = SYSTIMESTAMP
	WHERE id = (
	  SELECT access_list_id
	  FROM services
	  WHERE id = HEXTORAW(:1)
	)
	  AND is_deleted = 0`

		result, err := tx.ExecContext(ctx, q, id)
		if err != nil {
			s.logger.Errorf("[ServicesRepo][Delete] delete service failed: %v", err)
			return local_util.HandleDBError(err)
		}

		rows, _ = result.RowsAffected()
		if rows == 0 {
			return errors.New(localization.ErrorAccessListNotFound.Code)
		}
	} else if accessListID != "" {
		const q = `
	UPDATE access_lists
	SET
	  is_deleted = 1,
	  deleted_at = SYSTIMESTAMP,
	  last_modified_at = SYSTIMESTAMP
	WHERE id = HEXTORAW(:1)
	  AND is_deleted = 0`

		result, err := tx.ExecContext(ctx, q, accessListID)
		if err != nil {
			s.logger.Errorf("[ServicesRepo][Delete] delete service failed: %v", err)
			return local_util.HandleDBError(err)
		}

		rows, _ := result.RowsAffected()
		if rows == 0 {
			return errors.New(localization.ErrorAccessListNotFound.Code)
		}
	}

	if err := tx.Commit(); err != nil {
		s.logger.Errorf("[ServicesRepo][Delete] commit failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (s *ServicesStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	const q = `
UPDATE access_lists
SET
  is_enabled = :1,
  last_modified_at = SYSTIMESTAMP
WHERE id = (
  SELECT access_list_id
  FROM services
  WHERE id = HEXTORAW(:2)
)
  AND is_deleted = 0`

	res, err := s.db.ExecContext(ctx, q, boolToOracleNumber(enable), id)
	if err != nil {
		s.logger.Errorf("[ServicesRepo][EnableOrDisable] update failed: %v", err)
		return local_util.HandleDBError(err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New(localization.ErrorServiceNotFound.Code)
	}

	_, err = s.db.ExecContext(ctx, `UPDATE services SET last_modified_at = SYSTIMESTAMP WHERE id = HEXTORAW(:1)`, id)
	if err != nil {
		s.logger.Errorf("[ServicesRepo][EnableOrDisable] update services last_modified_at failed: %v", err)
		return local_util.HandleDBError(err)
	}

	return nil
}

func (s *ServicesStorage) FindByID(ctx context.Context, id string) (*service_dto.ServiceResponse, error) {
	const q = `
SELECT
  RAWTOHEX(s.id),
  RAWTOHEX(s.access_list_id),
  sk.name,
  sk.service_key,
  s.service_code,
  s.minimum_fraud_amount,
  s.product_gl_account_number,
  s.product_gl_account_currency,
  sk.is_enabled,
  sk.is_deleted,
  s.created_at,
  s.last_modified_at,
  sk.deleted_at
FROM services s
JOIN access_lists sk ON sk.id = s.access_list_id
WHERE s.id = HEXTORAW(:1) AND sk.is_deleted = 0`

	var serviceID string
	var svc service_dto.ServiceResponse
	var deletedAt sql.NullTime

	err := s.db.QueryRowContext(ctx, q, id).Scan(
		&serviceID,
		&svc.ServiceKeyId,
		&svc.ServiceName,
		&svc.ServiceKey,
		&svc.ServiceCode,
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

func (s *ServicesStorage) FindByAccessListID(ctx context.Context, accessListID string) (bool, error) {
	q := `SELECT 1 FROM services WHERE access_list_id = HEXTORAW(:1) AND is_deleted = 0 FETCH FIRST 1 ROWS ONLY`

	var exists int
	err := s.db.QueryRowContext(ctx, q, accessListID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

func (s *ServicesStorage) FindServiceByAccessListID(ctx context.Context, accessListID string) (*service_dto.ServiceResponse, error) {
	const q = `
SELECT
  RAWTOHEX(s.id),
  RAWTOHEX(s.access_list_id),
  sk.name,
  sk.service_key,
  s.service_code,
  s.minimum_fraud_amount,
  s.product_gl_account_number,
  s.product_gl_account_currency,
  sk.is_enabled,
  sk.is_deleted,
  s.created_at,
  s.last_modified_at,
  sk.deleted_at
FROM services s
JOIN access_lists sk ON sk.id = s.access_list_id
WHERE sk.id = HEXTORAW(:1) AND sk.is_deleted = 0`

	var serviceID string
	var svc service_dto.ServiceResponse
	var deletedAt sql.NullTime

	err := s.db.QueryRowContext(ctx, q, accessListID).Scan(
		&serviceID,
		&svc.ServiceKeyId,
		&svc.ServiceName,
		&svc.ServiceKey,
		&svc.ServiceCode,
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

	svc.ID = serviceID
	if deletedAt.Valid {
		svc.DeletedAt = &deletedAt.Time
	}

	caps, err := s.getServiceCaps(ctx, svc.ID)
	if err != nil {
		return nil, err
	}
	svc.Cap = caps

	return &svc, nil
}

func (s *ServicesStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]service_dto.ServiceResponse], error) {
	limit := int64(maxPaginationDefault)
	page := int64(1)
	if filterParam.PerPage > 0 {
		limit = int64(filterParam.PerPage)
	}
	if filterParam.Page > 0 {
		page = int64(filterParam.Page)
	}
	offset := (page - 1) * limit

	fromClause := fmt.Sprintf(`FROM %s s JOIN %s sk ON sk.id = s.access_list_id`, servicesTable, accessListTable)

	clauses := []string{"sk.is_deleted = 0", "s.is_deleted = 0"}
	var args []interface{}

	search := strings.TrimSpace(filterParam.Search)
	if search != "" {
		clauses = append(clauses,
			`(
				LOWER(RAWTOHEX(s.access_list_id)) LIKE '%' || LOWER(:search) || '%'
				OR LOWER(s.service_code) LIKE '%' || LOWER(:search) || '%'
				OR LOWER(sk.name) LIKE '%' || LOWER(:search) || '%'
				OR LOWER(sk.service_key) LIKE '%' || LOWER(:search) || '%'
			)`,
		)
		args = append(args, sql.Named("search", search))
	}

	if filterParam.Filters != nil {
		if v, ok := filterParam.Filters["enabled"]; ok {
			if b, ok2 := parseBoolFilter(v); ok2 {
				clauses = append(clauses, "sk.is_enabled = :enabled")
				args = append(args, sql.Named("enabled", boolToOracleNumber(b)))
			}
		}

		if v, ok := filterParam.Filters["service_name"]; ok {
			if sName, ok2 := v.(string); ok2 && strings.TrimSpace(sName) != "" {
				clauses = append(clauses, "LOWER(sk.name) = LOWER(:service_name)")
				args = append(args, sql.Named("name", sName))
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
				clauses = append(clauses, "LOWER(sk.service_key) = LOWER(:service_key)")
				args = append(args, sql.Named("service_key", sKey))
			}
		}
	}

	where := strings.Join(clauses, " AND ")
	countQ := fmt.Sprintf(`SELECT COUNT(*) %s WHERE %s`, fromClause, where)
	var total int64
	if err := s.db.QueryRowContext(ctx, countQ, args...).Scan(&total); err != nil {
		s.logger.Errorf("[ServicesRepo][FindAllWithPagination] count failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	listQ := fmt.Sprintf(`
SELECT
  RAWTOHEX(s.id),
  RAWTOHEX(s.access_list_id),
  sk.name,
  sk.service_key,
  s.service_code,
  s.minimum_fraud_amount,
  s.product_gl_account_number,
  s.product_gl_account_currency,
  sk.is_enabled,
  sk.is_deleted,
  s.created_at,
  s.last_modified_at,
  sk.deleted_at
%s
WHERE %s
ORDER BY s.created_at DESC
OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`, fromClause, where)

	listArgs := append(args, sql.Named("offset", offset), sql.Named("limit", limit))
	rows, err := s.db.QueryContext(ctx, listQ, listArgs...)
	if err != nil {
		s.logger.Errorf("[ServicesRepo][FindAllWithPagination] list query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer rows.Close()

	var list []service_dto.ServiceResponse
	for rows.Next() {
		var serviceID string
		var svc service_dto.ServiceResponse
		var deletedAt sql.NullTime

		if err := rows.Scan(
			&serviceID,
			&svc.ServiceKeyId,
			&svc.ServiceName,
			&svc.ServiceKey,
			&svc.ServiceCode,
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
	return &types.PaginatedResponse[[]service_dto.ServiceResponse]{Data: list, Meta: meta}, nil
}

func (s *ServicesStorage) CheckServiceExistence(ctx context.Context, serviceCode, serviceKey, serviceName string) (bool, error) {
	clauses := make([]string, 0, 3)
	args := make([]interface{}, 0, 3)

	if strings.TrimSpace(serviceCode) != "" {
		clauses = append(clauses, "s.service_code = :service_code")
		args = append(args, sql.Named("service_code", serviceCode))
	}
	if strings.TrimSpace(serviceKey) != "" {
		clauses = append(clauses, "sk.service_key = :service_key")
		args = append(args, sql.Named("service_key", serviceKey))
	}
	if strings.TrimSpace(serviceName) != "" {
		clauses = append(clauses, "sk.name = :service_name")
		args = append(args, sql.Named("name", serviceName))
	}
	if len(clauses) == 0 {
		return false, nil
	}

	where := strings.Join(clauses, " OR ")
	q := fmt.Sprintf(
		`SELECT COUNT(*) FROM %s s JOIN %s sk ON sk.id = s.access_list_id WHERE sk.is_deleted = 0 AND (%s)`,
		servicesTable,
		accessListTable,
		where,
	)

	var count int64
	if err := s.db.QueryRowContext(ctx, q, args...).Scan(&count); err != nil {
		s.logger.Errorf("[ServicesRepo][CheckServiceExistence] count failed: %v", err)
		return false, local_util.HandleDBError(err)
	}

	return count > 0, nil
}

func (s *ServicesStorage) FindAllServiceListWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.ServiceKey], error) {
	limit := int64(maxPaginationDefault)
	page := int64(1)
	if filterParam.PerPage > 0 {
		limit = int64(filterParam.PerPage)
	}
	if filterParam.Page > 0 {
		page = int64(filterParam.Page)
	}
	offset := (page - 1) * limit

	clauses := []string{"is_enabled = is_enabled", "is_deleted = 0"}
	clauses = []string{}
	args := []interface{}{}

	search := strings.TrimSpace(filterParam.Search)
	if search != "" {
		clauses = append(clauses,
			`(
				LOWER(name) LIKE '%' || LOWER(:search) || '%'
				OR LOWER(service_key) LIKE '%' || LOWER(:search) || '%'
			)`,
		)
		args = append(args, sql.Named("search", search))
	}

	if filterParam.Filters != nil {
		if v, ok := filterParam.Filters["is_enabled"]; ok {
			if b, ok2 := parseBoolFilter(v); ok2 {
				clauses = append(clauses, "is_enabled = :enabled")
				args = append(args, sql.Named("enabled", boolToOracleNumber(b)))
			}
		}
		if v, ok := filterParam.Filters["service_name"]; ok {
			if sName, ok2 := v.(string); ok2 && strings.TrimSpace(sName) != "" {
				clauses = append(clauses, "LOWER(name) = LOWER(:service_name)")
				args = append(args, sql.Named("name", sName))
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

	countQ := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE %s AND is_deleted = 0`, accessListTable, where)
	var total int64
	if err := s.db.QueryRowContext(ctx, countQ, args...).Scan(&total); err != nil {
		s.logger.Errorf("[ServicesRepo][FindAllServiceListWithPagination] count failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	listQ := fmt.Sprintf(`
SELECT
  RAWTOHEX(id),
  name,
  service_key,
  is_enabled,
  is_ussd_enabled,
  created_at,
  last_modified_at
FROM %s
WHERE %s AND is_deleted = 0
ORDER BY created_at DESC
OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`, accessListTable, where)

	listArgs := append(args, sql.Named("offset", offset), sql.Named("limit", limit))
	rows, err := s.db.QueryContext(ctx, listQ, listArgs...)
	if err != nil {
		s.logger.Errorf("[ServicesRepo][FindAllServiceListWithPagination] list query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer rows.Close()

	var list []imodel.ServiceKey
	for rows.Next() {
		var listID string
		var item imodel.ServiceKey
		if err := rows.Scan(
			&listID,
			&item.ServiceName,
			&item.ServiceKey,
			&item.IsEnabled,
			&item.IsUSSDEnabled,
			&item.CreatedAt,
			&item.LastModifiedAt,
		); err != nil {
			s.logger.Errorf("[ServicesRepo][FindAllServiceListWithPagination] scan failed: %v", err)
			return nil, local_util.HandleDBError(err)
		}

		item.ID = listID

		list = append(list, item)
	}

	meta := local_util.BuildPaginationMeta(total, int(page), int(limit))
	return &types.PaginatedResponse[[]imodel.ServiceKey]{Data: list, Meta: meta}, nil
}

func (s *ServicesStorage) FindServiceListByID(ctx context.Context, id string) (*imodel.ServiceKey, error) {
	const q = `
SELECT RAWTOHEX(id), name, service_key, is_enabled, is_ussd_enabled, created_at, last_modified_at
FROM access_lists
WHERE id = HEXTORAW(:1) AND is_deleted = 0`

	var item imodel.ServiceKey
	var listID string
	err := s.db.QueryRowContext(ctx, q, id).Scan(
		&listID,
		&item.ServiceName,
		&item.ServiceKey,
		&item.IsEnabled,
		&item.IsUSSDEnabled,
		&item.CreatedAt,
		&item.LastModifiedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(localization.ErrorAccessListNotFound.Code)
		}
		return nil, local_util.HandleDBError(err)
	}

	item.ID = listID

	return &item, nil
}

func (s *ServicesStorage) FindServiceListByNameOrKey(ctx context.Context, name, key string) (*imodel.ServiceKey, error) {
	conds := make([]string, 0, 2)
	args := make([]interface{}, 0, 2)

	if strings.TrimSpace(name) != "" {
		conds = append(conds, "LOWER(name) LIKE '%' || LOWER(:name) || '%'")
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
SELECT RAWTOHEX(id), name, service_key, is_enabled, created_at, last_modified_at
FROM %s
WHERE (%s)
FETCH FIRST 1 ROWS ONLY`, accessListTable, where)

	var item imodel.ServiceKey
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
			return nil, errors.New(localization.ErrorAccessListNotFound.Code)
		}
		return nil, local_util.HandleDBError(err)
	}

	item.ID = listID
	return &item, nil
}

// FindServiceListByExactNameOrKey returns a row when an active access_list exists with the same
// name or service_key as the given values (full-string equality, case-insensitive). Used before create.
func (s *ServicesStorage) FindServiceListByExactNameOrKey(ctx context.Context, name, key string) (*imodel.ServiceKey, error) {
	conds := make([]string, 0, 2)
	args := make([]interface{}, 0, 2)

	if strings.TrimSpace(name) != "" {
		conds = append(conds, "LOWER(TRIM(name)) = LOWER(TRIM(:svc_name))")
		args = append(args, sql.Named("svc_name", strings.TrimSpace(name)))
	}
	if strings.TrimSpace(key) != "" {
		conds = append(conds, "LOWER(TRIM(service_key)) = LOWER(TRIM(:svc_key))")
		args = append(args, sql.Named("svc_key", strings.TrimSpace(key)))
	}

	if len(conds) == 0 {
		return nil, errors.New(localization.ErrorNoDataProvided.Code)
	}

	where := strings.Join(conds, " OR ")
	query := fmt.Sprintf(`
SELECT RAWTOHEX(id), name, service_key, is_enabled, created_at, last_modified_at
FROM %s
WHERE is_deleted = 0 AND (%s)
FETCH FIRST 1 ROWS ONLY`, accessListTable, where)

	var item imodel.ServiceKey
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
			return nil, errors.New(localization.ErrorAccessListNotFound.Code)
		}
		return nil, local_util.HandleDBError(err)
	}

	item.ID = listID
	return &item, nil
}

func (s *ServicesStorage) CreateServiceKey(ctx context.Context, serviceList *imodel.ServiceKey) error {
	serviceList.IsEnabled = true

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Errorf("[ServicesRepo][CreateServiceKey] begin tx failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = tx.Rollback() }()

	const q = `
INSERT INTO access_lists (
  name,
  service_key,
  account_type,
  is_enabled
)
VALUES (
  :1,:2,:3,:4
)`

	if _, err := tx.ExecContext(ctx, q,
		serviceList.ServiceName,
		serviceList.ServiceKey,
		serviceList.AccountType,
		boolToOracleNumber(serviceList.IsEnabled),
	); err != nil {
		s.logger.Errorf("[ServicesRepo][CreateServiceKey] insert failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	// Synchronize access_list with this service list
	// 	const accessQ = `
	// INSERT INTO access_list (
	//   key,
	//   enabled,
	//   access_list_name,
	//   ussd_enabled,
	//   created_at,
	//   last_modified_at
	// )
	// VALUES (
	//   :1,:2,:3,:4,:5,:6
	// )`
	// 	if _, err := tx.ExecContext(ctx, accessQ,
	// 		serviceList.ServiceKey,
	// 		boolToOracleNumber(serviceList.IsEnabled),
	// 		serviceList.ServiceName,
	// 		0,
	// 		serviceList.CreatedAt,
	// 		serviceList.LastModifiedAt,
	// 	); err != nil {
	// 		s.logger.Errorf("[ServicesRepo][CreateServiceKey] insert access_list failed: %v", err)
	// 		return local_util.HandleDBError(err)
	// 	}

	if err := tx.Commit(); err != nil {
		s.logger.Errorf("[ServicesRepo][CreateServiceKey] commit failed: %v", err)
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

func (s *ServicesStorage) UpdateServiceKey(ctx context.Context, id, serviceKey string, serviceList *imodel.ServiceKey) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Errorf("[ServicesRepo][UpdateServiceKey] begin tx failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = tx.Rollback() }()
	const q = `
UPDATE access_lists
SET
  name = :1,
  service_key = :2,
  account_type = :3,
  last_modified_at = SYSTIMESTAMP
WHERE id = :4 AND service_key = :5`

	res, err := tx.ExecContext(ctx, q,
		serviceList.ServiceName,
		serviceList.ServiceKey,
		serviceList.AccountType,
		id,
		serviceKey,
	)
	if err != nil {
		s.logger.Errorf("[ServicesRepo][UpdateServiceKey] update failed: %v", err)
		return local_util.HandleDBError(err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New(localization.ErrorAccessListNotFound.Code)
	}

	// 	const accessUpdateQ = `
	// UPDATE access_list
	// SET
	//   key = :1,
	//   access_list_name = :2,
	//   last_modified_at = SYSTIMESTAMP
	// WHERE key = :3`
	// 	if _, err := tx.ExecContext(ctx, accessUpdateQ,
	// 		serviceList.ServiceKey,
	// 		serviceList.ServiceName,
	// 		serviceKey,
	// 	); err != nil {
	// 		s.logger.Errorf("[ServicesRepo][UpdateServiceKey] update access_list failed: %v", err)
	// 		return local_util.HandleDBError(err)
	// 	}

	if err := tx.Commit(); err != nil {
		s.logger.Errorf("[ServicesRepo][UpdateServiceKey] commit failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (s *ServicesStorage) EnableOrDisableServiceList(ctx context.Context, id string, enable bool) error {
	const q = `
UPDATE access_lists
SET
  is_enabled = :1,
  last_modified_at = SYSTIMESTAMP
WHERE id = :2`

	res, err := s.db.ExecContext(ctx, q, boolToOracleNumber(enable), id)
	if err != nil {
		s.logger.Errorf("[ServicesRepo][EnableOrDisableServiceList] update failed: %v", err)
		return local_util.HandleDBError(err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New(localization.ErrorAccessListNotFound.Code)
	}

	// result, err := s.FindServiceListByID(ctx, id)
	// if err != nil {
	// 	s.logger.Errorf("[ServicesRepo][EnableOrDisableServiceList] find service list failed: %v", err)
	// 	return err
	// }

	// const accessEnableQ = `
	// UPDATE access_list
	// SET
	//   enabled = :1,
	//   last_modified_at = SYSTIMESTAMP
	// WHERE key = :2`
	// if _, err := s.db.ExecContext(ctx, accessEnableQ, boolToOracleNumber(enable), result.ServiceKey); err != nil {
	// 	s.logger.Errorf("[ServicesRepo][EnableOrDisableServiceList] update access_list failed: %v", err)
	// 	return local_util.HandleDBError(err)
	// }

	return nil
}

func (s *ServicesStorage) DeleteServiceKey(ctx context.Context, id string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Errorf("[ServicesRepo][Create] begin tx failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = tx.Rollback() }()

	// check if services use's this access list
	const checkQ = `SELECT COUNT(*) FROM services WHERE access_list_id = HEXTORAW(:1) AND is_deleted = 0`
	var count int64
	if err := tx.QueryRowContext(ctx, checkQ, id).Scan(&count); err != nil {
		s.logger.Errorf("[ServicesRepo][DeleteServiceList] check services failed: %v", err)
		return local_util.HandleDBError(err)
	}
	if count > 0 {
		return errors.New(localization.ErrorServiceListInUse.Code)
	}

	// Delete the service list from access_lists and then delete the service which has the access_list_id
	const q = `
	UPDATE access_lists 
	SET
		is_deleted = 1, 
		deleted_at = SYSTIMESTAMP, 
		last_modified_at = SYSTIMESTAMP 
	WHERE id = HEXTORAW(:1) AND is_deleted = 0`
	res, err := tx.ExecContext(ctx, q, id)
	if err != nil {
		s.logger.Errorf("[ServicesRepo][DeleteServiceList] delete failed: %v", err)
		return local_util.HandleDBError(err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New(localization.ErrorAccessListNotFound.Code)
	}

	if err := tx.Commit(); err != nil {
		s.logger.Errorf("[ServicesRepo][DeleteServiceList] commit failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}
