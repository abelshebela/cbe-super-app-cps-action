package logistics_merchant_oracle

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"database/sql"
	"time"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type LogisticsMerchantOracle struct {
	db     *sql.DB
	cfg    *config.VaultConfig
	logger utils.Logger
}

// Create implements [storage.LogisticsMerchantOracleRepository].
func (l *LogisticsMerchantOracle) Create(ctx context.Context, logisticsMerchant model.LogisticsMerchant) error {

	stmt := `INSERT INTO MERCHANTS (
		ID, MERCHANT_ACCOUNT_NUMBER, MERCHANT_CODE, MERCHANT_NAME, SETTLEMENT_METHOD, MERCHANT_TYPE, CONTACT_EMAIL, CONTACT_PHONE, IS_ENABLED, IS_DELETED, CREATED_AT, LAST_MODIFIED_AT, DELETED_AT
	) VALUES (
		SYS_GUID(), :1, :2, :3, :4, :5, :6, :7, :8, :9, SYSTIMESTAMP, SYSTIMESTAMP, NULL
	)`

	_, err := l.db.ExecContext(ctx, stmt,
		logisticsMerchant.BankAccountNumber,
		logisticsMerchant.MerchantID,
		logisticsMerchant.MerchantName,
		logisticsMerchant.SettlementMethod,
		logisticsMerchant.MerchantType,
		logisticsMerchant.Enabled,
		logisticsMerchant.IsDeleted,
		logisticsMerchant.CreatedAt,
		logisticsMerchant.UpdatedAt,
	)
	if err != nil {
		l.logger.Errorf("Failed to create logistics merchant: %v", err)
		return err
	}
	return nil
}

// Delete implements [storage.LogisticsMerchantOracleRepository].
func (l *LogisticsMerchantOracle) Delete(ctx context.Context, id string) error {
	stmt := `DELETE FROM MERCHANTS WHERE ID = HEXTORAW(:1)`
	_, err := l.db.ExecContext(ctx, stmt, id)
	if err != nil {
		l.logger.Errorf("Failed to delete logistics merchant: %v", err)
		return err
	}
	return nil
}

// EnableOrDisable implements [storage.LogisticsMerchantOracleRepository].
func (l *LogisticsMerchantOracle) EnableOrDisable(ctx context.Context, ids []string, enable bool) error {
	for _, id := range ids {
		stmt := `UPDATE MERCHANTS SET IS_ENABLED = :1, LAST_MODIFIED_AT = SYSTIMESTAMP WHERE ID = HEXTORAW(:2) AND IS_DELETED = 0`
		_, err := l.db.ExecContext(ctx, stmt, boolToInt(enable), id)
		if err != nil {
			l.logger.Errorf("Failed to enable/disable logistics merchant: %v", err)
			return err
		}
	}
	return nil
}

// FindAllWithPagination implements [storage.LogisticsMerchantOracleRepository].
func (l *LogisticsMerchantOracle) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.LogisticsMerchant], error) {
	where := "WHERE IS_DELETED = 0 AND MERCHANT_TYPE = 'LOGISTICS'"
	args := []interface{}{}
	if filterParam.Search != "" {
		where += " AND (LOWER(MERCHANT_NAME) LIKE :1 OR LOWER(MERCHANT_CODE) LIKE :2)"
		args = append(args, "%"+filterParam.Search+"%", "%"+filterParam.Search+"%")
	}
	countQuery := "SELECT COUNT(*) FROM MERCHANTS " + where
	var total int64
	err := l.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		l.logger.Errorf("Failed to count logistics merchants: %v", err)
		return nil, err
	}
	page := filterParam.Page
	perPage := filterParam.PerPage
	offset := (page - 1) * perPage
	query := `SELECT RAWTOHEX(ID), MERCHANT_ACCOUNT_NUMBER, MERCHANT_CODE, MERCHANT_NAME, SETTLEMENT_METHOD, MERCHANT_TYPE, IS_ENABLED, IS_DELETED, CREATED_AT, LAST_MODIFIED_AT, DELETED_AT
		FROM MERCHANTS ` + where + ` ORDER BY CREATED_AT DESC OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`
	args = append(args, offset, perPage)
	rows, err := l.db.QueryContext(ctx, query, args...)
	if err != nil {
		l.logger.Errorf("Failed to fetch paginated logistics merchants: %v", err)
		return nil, err
	}
	defer rows.Close()
	var data []model.LogisticsMerchant
	for rows.Next() {
		var m model.LogisticsMerchant
		var enabled, isDeleted int
		var createdAt, updatedAt, deletedAt sql.NullTime
		if err := rows.Scan(&m.ID, &m.BankAccountNumber, &m.MerchantID, &m.MerchantName, &m.SettlementMethod, &m.MerchantType, &enabled, &isDeleted, &createdAt, &updatedAt, &deletedAt); err != nil {
			l.logger.Errorf("Failed to scan logistics merchant: %v", err)
			return nil, err
		}
		m.Enabled = enabled == 1
		m.IsDeleted = isDeleted == 1
		if createdAt.Valid {
			m.CreatedAt = createdAt.Time
		} else {
			m.CreatedAt = time.Time{}
		}
		if updatedAt.Valid {
			m.UpdatedAt = updatedAt.Time
		} else {
			m.UpdatedAt = time.Time{}
		}
		if deletedAt.Valid {
			m.DeletedAt = deletedAt.Time
		} else {
			m.DeletedAt = time.Time{}
		}
		data = append(data, m)
	}
	meta := local_util.BuildPaginationMeta(total, page, perPage)
	return &types.PaginatedResponse[[]model.LogisticsMerchant]{
		Data: data,
		Meta: meta,
	}, nil
}

// FindByID implements [storage.LogisticsMerchantOracleRepository].
func (l *LogisticsMerchantOracle) FindByID(ctx context.Context, id string) (*model.LogisticsMerchant, error) {
	stmt := `SELECT RAWTOHEX(ID), MERCHANT_ACCOUNT_NUMBER, MERCHANT_CODE, MERCHANT_NAME, SETTLEMENT_METHOD, MERCHANT_TYPE, IS_ENABLED, IS_DELETED, CREATED_AT, LAST_MODIFIED_AT, DELETED_AT FROM MERCHANTS WHERE ID = HEXTORAW(:1) AND IS_DELETED = 0 AND MERCHANT_TYPE = 'LOGISTICS'`
	var m model.LogisticsMerchant
	var enabled, isDeleted int
	var createdAt, updatedAt, deletedAt sql.NullTime
	err := l.db.QueryRowContext(ctx, stmt, id).Scan(&m.ID, &m.BankAccountNumber, &m.MerchantID, &m.MerchantName, &m.SettlementMethod, &m.MerchantType, &enabled, &isDeleted, &createdAt, &updatedAt, &deletedAt)
	if err != nil {
		l.logger.Errorf("Failed to find logistics merchant by id: %v", err)
		return nil, err
	}
	m.Enabled = enabled == 1
	m.IsDeleted = isDeleted == 1
	if createdAt.Valid {
		m.CreatedAt = createdAt.Time
	} else {
		m.CreatedAt = time.Time{}
	}
	if updatedAt.Valid {
		m.UpdatedAt = updatedAt.Time
	} else {
		m.UpdatedAt = time.Time{}
	}
	if deletedAt.Valid {
		m.DeletedAt = deletedAt.Time
	} else {
		m.DeletedAt = time.Time{}
	}
	return &m, nil
}

// FindOne implements [storage.LogisticsMerchantOracleRepository].
func (l *LogisticsMerchantOracle) FindOne(ctx context.Context, filter bson.M) (*model.LogisticsMerchant, error) {
	// For Oracle, only support lookup by merchant_code or merchant_id for now
	var stmt string
	var arg string
	if v, ok := filter["merchant_code"]; ok {
		stmt = `SELECT RAWTOHEX(ID), MERCHANT_ACCOUNT_NUMBER, MERCHANT_CODE, MERCHANT_NAME, SETTLEMENT_METHOD, MERCHANT_TYPE, IS_ENABLED, IS_DELETED, CREATED_AT, LAST_MODIFIED_AT, DELETED_AT FROM MERCHANTS WHERE MERCHANT_CODE = :1 AND IS_DELETED = 0 AND MERCHANT_TYPE = 'LOGISTICS'`
		arg = v.(string)
	} else if v, ok := filter["merchant_id"]; ok {
		stmt = `SELECT RAWTOHEX(ID), MERCHANT_ACCOUNT_NUMBER, MERCHANT_CODE, MERCHANT_NAME, SETTLEMENT_METHOD, MERCHANT_TYPE, IS_ENABLED, IS_DELETED, CREATED_AT, LAST_MODIFIED_AT, DELETED_AT FROM MERCHANTS WHERE MERCHANT_CODE = :1 AND IS_DELETED = 0 AND MERCHANT_TYPE = 'LOGISTICS'`
		arg = v.(string)
	}

	var m model.LogisticsMerchant
	var enabled, isDeleted int
	var createdAt, updatedAt, deletedAt sql.NullTime
	err := l.db.QueryRowContext(ctx, stmt, arg).Scan(&m.ID, &m.BankAccountNumber, &m.MerchantID, &m.MerchantName, &m.SettlementMethod, &m.MerchantType, &enabled, &isDeleted, &createdAt, &updatedAt, &deletedAt)
	if err != nil {
		l.logger.Errorf("Failed to find logistics merchant: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	m.Enabled = enabled == 1
	m.IsDeleted = isDeleted == 1
	if createdAt.Valid {
		m.CreatedAt = createdAt.Time
	} else {
		m.CreatedAt = time.Time{}
	}
	if updatedAt.Valid {
		m.UpdatedAt = updatedAt.Time
	} else {
		m.UpdatedAt = time.Time{}
	}
	if deletedAt.Valid {
		m.DeletedAt = deletedAt.Time
	} else {
		m.DeletedAt = time.Time{}
	}
	return &m, nil
}

// FindByAccountOrMerchantCode implements [storage.LogisticsMerchantOracleRepository].
// Returns the first active logistics merchant that matches either accountNumber or
// merchantCode. If excludeID is non-empty, that row is skipped (used during updates).
func (l *LogisticsMerchantOracle) FindByAccountOrMerchantCode(ctx context.Context, accountNumber, merchantCode, excludeID string) (*model.LogisticsMerchant, error) {
	baseQ := `SELECT RAWTOHEX(ID), MERCHANT_ACCOUNT_NUMBER, MERCHANT_CODE, MERCHANT_NAME,
	           SETTLEMENT_METHOD, MERCHANT_TYPE, IS_ENABLED, IS_DELETED,
	           CREATED_AT, LAST_MODIFIED_AT, DELETED_AT
	          FROM MERCHANTS
	          WHERE IS_DELETED = 0 AND MERCHANT_TYPE = 'LOGISTICS'
	            AND (MERCHANT_ACCOUNT_NUMBER = :1 OR MERCHANT_CODE = :2)`

	args := []interface{}{accountNumber, merchantCode}
	if excludeID != "" {
		baseQ += ` AND ID != HEXTORAW(:3)`
		args = append(args, excludeID)
	}
	baseQ += ` FETCH FIRST 1 ROWS ONLY`

	var m model.LogisticsMerchant
	var enabled, isDeleted int
	var createdAt, updatedAt, deletedAt sql.NullTime

	err := l.db.QueryRowContext(ctx, baseQ, args...).Scan(
		&m.ID, &m.BankAccountNumber, &m.MerchantID, &m.MerchantName,
		&m.SettlementMethod, &m.MerchantType, &enabled, &isDeleted,
		&createdAt, &updatedAt, &deletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		l.logger.Errorf("[LogisticsMerchantOracle][FindByAccountOrMerchantCode] err: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	m.Enabled = enabled == 1
	m.IsDeleted = isDeleted == 1
	if createdAt.Valid {
		m.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		m.UpdatedAt = updatedAt.Time
	}
	if deletedAt.Valid {
		m.DeletedAt = deletedAt.Time
	}
	return &m, nil
}

// Update implements [storage.LogisticsMerchantOracleRepository].
func (l *LogisticsMerchantOracle) Update(ctx context.Context, id string, logisticsMerchant model.LogisticsMerchant) error {
	stmt := `UPDATE MERCHANTS SET MERCHANT_ACCOUNT_NUMBER = :1, MERCHANT_CODE = :2, MERCHANT_NAME = :3, SETTLEMENT_METHOD = :4, MERCHANT_TYPE = :5, IS_ENABLED = :6, LAST_MODIFIED_AT = SYSTIMESTAMP WHERE ID = HEXTORAW(:7)`
	_, err := l.db.ExecContext(ctx, stmt,
		logisticsMerchant.BankAccountNumber,
		logisticsMerchant.MerchantID,
		logisticsMerchant.MerchantName,
		logisticsMerchant.SettlementMethod,
		logisticsMerchant.MerchantType,
		boolToInt(logisticsMerchant.Enabled),
		id,
	)
	if err != nil {
		l.logger.Errorf("Failed to update logistics merchant: %v", err)
		return err
	}
	return nil
}

// Helper to convert bool to int for Oracle NUMBER(1)
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func NewLogisticsMerchantOracle(db *sql.DB, cfg *config.VaultConfig, logger utils.Logger) storage.LogisticsMerchantOracleRepository {
	return &LogisticsMerchantOracle{
		db:     db,
		cfg:    cfg,
		logger: logger,
	}
}
