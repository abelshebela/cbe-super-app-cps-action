package logistics_merchant_oracle

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"database/sql"
	"fmt"
	"strings"

	shared_model "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type LogisticsMerchantOracle struct {
	db     *sql.DB
	logger utils.Logger
}

func NewLogisticsMerchantOracleRepository(db *sql.DB, log utils.Logger) storage.LogisticsMerchantOracleRepository {
	return &LogisticsMerchantOracle{
		db:     db,
		logger: log,
	}
}

// CheckAccountOrCreate implements [storage.LogisticsMerchantOracleRepository].
func (q *LogisticsMerchantOracle) CheckAccountOrCreate(ctx context.Context, account_detail *shared_model.AccountDetail) error {
	checkAccountExistsQuery := `SELECT id FROM accounts where account_number = :1`
	var id string
	err := q.db.QueryRowContext(ctx, checkAccountExistsQuery, account_detail.AccountNumber).Scan(&id)
	if err != nil {
		if err != sql.ErrNoRows {
			q.logger.Errorf("[LogisticsMerchantOracle][CheckAccountOrCreate] Failed to check account for account_number=%s: %v", account_detail.AccountNumber, err)
			return localization.ErrorUnexpectedError
		}
		q.logger.Infof("[LogisticsMerchantOracle][CheckAccountOrCreate] Account not found, creating new account for account_number=%s", account_detail.AccountNumber)
	} else {
		return nil // Account already exists, no error
	}

	findBankIDForISCBEQuery := `SELECT id FROM banks where is_cbe=1 AND is_enabled=1 AND is_deleted=0`

	var bankID string
	err = q.db.QueryRowContext(ctx, findBankIDForISCBEQuery).Scan(&bankID)
	if err != nil {
		if err == sql.ErrNoRows {
			q.logger.Errorf("[LogisticsMerchantOracle][CheckAccountOrCreate] No enabled ISCBE bank found")
			return localization.ErrorUnexpectedError
		}
		q.logger.Errorf("[LogisticsMerchantOracle][CheckAccountOrCreate] Failed to find ISCBE bank: %v", err)
		return localization.ErrorUnexpectedError
	}

	createAccountQuery := `INSERT INTO accounts (bank_id,account_holder_name ,account_number, account_currency,account_type,customer_number, created_at) VALUES (:1, :2, :3, :4, :5, :6, CURRENT_TIMESTAMP)`
	_, err = q.db.ExecContext(ctx, createAccountQuery, bankID, account_detail.CustomerName, account_detail.AccountNumber, account_detail.Currency, account_detail.AccountType, account_detail.CustomerID)
	if err != nil {
		q.logger.Errorf("[LogisticsMerchantOracle][CheckAccountOrCreate] Failed to create account for account_number=%s: %v", account_detail.AccountNumber, err)
		return localization.ErrorUnexpectedError
	}
	q.logger.Infof("[LogisticsMerchantOracle][CheckAccountOrCreate] Successfully created account for account_number=%s", account_detail.AccountNumber)
	return nil
}

// Create implements [storage.LogisticsMerchantRepository].
func (q *LogisticsMerchantOracle) Create(ctx context.Context, logisticsMerchant model.LogisticsMerchantOracle) error {
	// Insert into MERCHANTS table
	query := `INSERT INTO merchants (
	       id, merchant_account_number, merchant_code, merchant_name, settlement_method, merchant_type, contact_email, contact_phone, is_enabled, is_deleted, created_at, last_modified_at
       ) VALUES (
	       SYS_GUID(), :1, :2, :3, :4, :5, :6, :7, :8, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
       )`
	isEnabled := 1
	if !logisticsMerchant.IsEnabled {
		isEnabled = 0
	}
	_, err := q.db.ExecContext(ctx, query,
		logisticsMerchant.MerchantAccountNumber,
		logisticsMerchant.MerchantCode,
		logisticsMerchant.MerchantName,
		logisticsMerchant.SettlementMethod,
		logisticsMerchant.MerchantType,
		logisticsMerchant.ContactEmail,
		logisticsMerchant.ContactPhone,
		isEnabled,
	)
	if err != nil {
		q.logger.Errorf("[LogisticsMerchantOracle][Create] Failed to create merchant: %v", err)
		return localization.ErrorUnexpectedError
	}
	return nil
}

// Delete implements [storage.LogisticsMerchantRepository].
func (q *LogisticsMerchantOracle) Delete(ctx context.Context, id string) error {
	// Soft delete: set is_deleted, deleted_at, updated_at
	query := `UPDATE logistics_merchants SET is_deleted = 1, deleted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP WHERE id = :1`
	res, err := q.db.ExecContext(ctx, query, id)
	if err != nil {
		q.logger.Errorf("[LogisticsMerchantOracle][Delete] Failed to delete id=%s: %v", id, err)
		return err // Replace with localization if needed
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		q.logger.Errorf("[LogisticsMerchantOracle][Delete] No rows affected for id=%s", id)
		// Return not found error
		return sql.ErrNoRows
	}
	return nil
}

// EnableOrDisable implements [storage.LogisticsMerchantRepository].
func (q *LogisticsMerchantOracle) EnableOrDisable(ctx context.Context, ids []string, enable bool) error {
	q.logger.Debugf("[LogisticsMerchantOracle][EnableOrDisable] Setting enabled=%t for ids=%v", enable, ids)
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		q.logger.Errorf("[LogisticsMerchantOracle][EnableOrDisable] Failed to begin transaction: %v", err)
		return localization.ErrorUnexpectedError
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	query := `UPDATE merchants SET is_enabled = :1, last_modified_at = CURRENT_TIMESTAMP WHERE id = HEXTORAW(:2) AND is_deleted = 0`
	var enableInt int
	if enable {
		enableInt = 1
	} else {
		enableInt = 0
	}
	for _, id := range ids {
		res, execErr := tx.ExecContext(ctx, query, enableInt, id)
		if execErr != nil {
			q.logger.Errorf("[LogisticsMerchantOracle][EnableOrDisable] Failed for id=%s: %v", id, execErr)
			err = localization.ErrorUnexpectedError
			return err
		}
		rows, _ := res.RowsAffected()
		if rows == 0 {
			q.logger.Warnf("[LogisticsMerchantOracle][EnableOrDisable] No rows affected for id=%s", id)
		}
	}
	return nil
}

// FindAllWithPagination implements [storage.LogisticsMerchantRepository].
func (q *LogisticsMerchantOracle) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.LogisticsMerchantOracle], error) {
	// Build WHERE clause
	where := "WHERE is_deleted = 0"
	args := []interface{}{}
	argIdx := 1
	if filterParam.Search != "" {
		where += " AND (LOWER(merchant_code) LIKE :" + string(rune(argIdx)) + " OR LOWER(merchant_name) LIKE :" + string(rune(argIdx)) + " OR LOWER(contact_email) LIKE :" + string(rune(argIdx)) + " OR LOWER(contact_phone) LIKE :" + string(rune(argIdx)) + ")"
		search := "%" + filterParam.Search + "%"
		args = append(args, search, search, search, search)
		argIdx += 4
	}
	// Add more filters as needed (merchant_type, is_enabled, etc.)
	if filterParam.Filters != nil {
		if v, ok := filterParam.Filters["merchant_type"].(string); ok && v != "" {
			where += " AND merchant_type = :" + string(rune(argIdx))
			args = append(args, v)
			argIdx++
		}
		if v, ok := filterParam.Filters["is_enabled"].(bool); ok {
			where += " AND is_enabled = :" + string(rune(argIdx))
			if v {
				args = append(args, 1)
			} else {
				args = append(args, 0)
			}
			argIdx++
		}
	}
	// Pagination
	page := filterParam.Page
	perPage := filterParam.PerPage
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}
	offset := (page - 1) * perPage
	// Count total
	countQuery := "SELECT COUNT(*) FROM merchants " + where
	var total int64
	err := q.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		q.logger.Errorf("[LogisticsMerchantOracle][FindAllWithPagination] Count error: %v", err)
		return nil, localization.ErrorUnexpectedError
	}
	// Data query
	dataQuery := `SELECT id, merchant_account_number, merchant_code, merchant_name, settlement_method, merchant_type, contact_email, contact_phone, is_enabled, is_deleted, created_at, last_modified_at, deleted_at FROM merchants ` + where + ` ORDER BY created_at DESC OFFSET :` + string(rune(argIdx)) + ` ROWS FETCH NEXT :` + string(rune(argIdx+1)) + ` ROWS ONLY`
	args = append(args, offset, perPage)
	rows, err := q.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		q.logger.Errorf("[LogisticsMerchantOracle][FindAllWithPagination] Query error: %v", err)
		return nil, localization.ErrorUnexpectedError
	}
	defer rows.Close()
	var merchants []model.LogisticsMerchantOracle
	for rows.Next() {
		var m model.LogisticsMerchantOracle
		err := rows.Scan(
			&m.ID,
			&m.MerchantAccountNumber,
			&m.MerchantCode,
			&m.MerchantName,
			&m.SettlementMethod,
			&m.MerchantType,
			&m.ContactEmail,
			&m.ContactPhone,
			&m.IsEnabled,
			&m.IsDeleted,
			&m.CreatedAt,
			&m.LastModifiedAt,
			&m.DeletedAt,
		)
		if err != nil {
			q.logger.Errorf("[LogisticsMerchantOracle][FindAllWithPagination] Row scan error: %v", err)
			continue
		}
		merchants = append(merchants, m)
	}
	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]model.LogisticsMerchantOracle]{
		Data: merchants,
		Meta: meta,
	}, nil
}

// FindByID implements [storage.LogisticsMerchantRepository].
func (q *LogisticsMerchantOracle) FindByID(ctx context.Context, id string) (*model.LogisticsMerchantOracle, error) {
	q.logger.Debugf("[LogisticsMerchantOracle][FindByID] Fetching logistics merchant with id=%s", id)
	query := `SELECT id, merchant_account_number, merchant_code, merchant_name, settlement_method, merchant_type, contact_email, contact_phone, is_enabled, is_deleted, created_at, last_modified_at, deleted_at FROM merchants WHERE id = HEXTORAW(:1) AND is_deleted = 0`
	row := q.db.QueryRowContext(ctx, query, id)
	var m model.LogisticsMerchantOracle
	err := row.Scan(
		&m.ID,
		&m.MerchantAccountNumber,
		&m.MerchantCode,
		&m.MerchantName,
		&m.SettlementMethod,
		&m.MerchantType,
		&m.ContactEmail,
		&m.ContactPhone,
		&m.IsEnabled,
		&m.IsDeleted,
		&m.CreatedAt,
		&m.LastModifiedAt,
		&m.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			q.logger.Warnf("[LogisticsMerchantOracle][FindByID] Not found: id=%s", id)
			return nil, localization.ErrorResourceNotFound
		}
		q.logger.Errorf("[LogisticsMerchantOracle][FindByID] Query error: %v", err)
		return nil, localization.ErrorUnexpectedError
	}
	return &m, nil
}

// FindOne implements [storage.LogisticsMerchantRepository].
func (q *LogisticsMerchantOracle) FindOne(ctx context.Context, filter bson.M) (*model.LogisticsMerchantOracle, error) {
	var conditions []string
	var args []interface{}
	paramIdx := 1

	// BankAccountNumber + MerchantCode
	if bankAcc, ok := filter["bank_account_number"].(string); ok && bankAcc != "" {
		if merchantCode, ok := filter["merchant_id"].(string); ok && merchantCode != "" {
			cond := fmt.Sprintf("(merchant_account_number = :%d AND merchant_code = :%d)", paramIdx, paramIdx+1)
			conditions = append(conditions, cond)
			args = append(args, bankAcc, merchantCode)
			paramIdx += 2
		}
	}

	// MerchantCode only (if not already used above)
	if merchantCode, ok := filter["merchant_id"].(string); ok && merchantCode != "" {
		// Only add merchant_code only condition if bank_account_number was not present
		if _, ok := filter["bank_account_number"].(string); !ok || filter["bank_account_number"].(string) == "" {
			cond := fmt.Sprintf("merchant_code = :%d", paramIdx)
			conditions = append(conditions, cond)
			args = append(args, merchantCode)
			paramIdx++
		}
	}

	if len(conditions) == 0 {
		return nil, nil
	}

	whereClause := "is_deleted = 0 AND (" + strings.Join(conditions, " OR ") + ")"

	// ExcludeID (optional)
	if excludeID, ok := filter["_id"].(string); ok && excludeID != "" {
		whereClause += fmt.Sprintf(" AND id != HEXTORAW(:%d)", paramIdx)
		args = append(args, excludeID)
		paramIdx++
	}

	query := `SELECT id, merchant_account_number, merchant_code, merchant_name, settlement_method, merchant_type, contact_email, contact_phone, is_enabled, is_deleted, created_at, last_modified_at, deleted_at FROM merchants WHERE ` + whereClause + ` FETCH NEXT 1 ROWS ONLY`

	row := q.db.QueryRowContext(ctx, query, args...)
	var m model.LogisticsMerchantOracle
	err := row.Scan(
		&m.ID,
		&m.MerchantAccountNumber,
		&m.MerchantCode,
		&m.MerchantName,
		&m.SettlementMethod,
		&m.MerchantType,
		&m.ContactEmail,
		&m.ContactPhone,
		&m.IsEnabled,
		&m.IsDeleted,
		&m.CreatedAt,
		&m.LastModifiedAt,
		&m.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		q.logger.Errorf("[LogisticsMerchantOracle][FindOne] Query error: %v", err)
		return nil, err
	}
	return &m, nil
}

// Update implements [storage.LogisticsMerchantRepository].
func (q *LogisticsMerchantOracle) Update(ctx context.Context, id string, logisticsMerchant model.LogisticsMerchantOracle) error {
	// Update merchant fields except id, created_at
	query := `UPDATE merchants SET 
	       merchant_account_number = :1,
	       merchant_code = :2,
	       merchant_name = :3,
	       settlement_method = :4,
	       merchant_type = :5,
	       contact_email = :6,
	       contact_phone = :7,
	       is_enabled = :8,
	       last_modified_at = CURRENT_TIMESTAMP
	       WHERE id = HEXTORAW(:9) AND is_deleted = 0`
	isEnabled := 1
	if !logisticsMerchant.IsEnabled {
		isEnabled = 0
	}
	res, err := q.db.ExecContext(ctx, query,
		logisticsMerchant.MerchantAccountNumber,
		logisticsMerchant.MerchantCode,
		logisticsMerchant.MerchantName,
		logisticsMerchant.SettlementMethod,
		logisticsMerchant.MerchantType,
		logisticsMerchant.ContactEmail,
		logisticsMerchant.ContactPhone,
		isEnabled,
		id,
	)
	if err != nil {
		q.logger.Errorf("[LogisticsMerchantOracle][Update] Failed to update merchant id=%s: %v", id, err)
		return localization.ErrorUnexpectedError
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		q.logger.Warnf("[LogisticsMerchantOracle][Update] No rows affected for id=%s", id)
		return localization.ErrorResourceNotFound
	}
	return nil
}
