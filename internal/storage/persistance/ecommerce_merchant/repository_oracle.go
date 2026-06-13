package ecommercemerchant

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

const (
	merchantsTable        = "MERCHANTS"
	merchantBranchesTable = "MERCHANT_BRANCHES"
	defaultPageSize       = 50
)

type EcommerceMerchantStorage struct {
	db     *sql.DB
	logger utils.Logger
}

func NewEcommerceMerchantRepository(db *sql.DB, logger utils.Logger) storage.EcommerceMerchantRepository {
	return &EcommerceMerchantStorage{
		db:     db,
		logger: logger,
	}
}

func boolToOracleNumber(v bool) int {
	if v {
		return 1
	}
	return 0
}

func nullString(v string) sql.NullString {
	trimmed := strings.TrimSpace(v)
	if trimmed == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: trimmed, Valid: true}
}

func parseOracleBool(v interface{}) (bool, bool) {
	switch b := v.(type) {
	case bool:
		return b, true
	case int:
		return b != 0, true
	case int32:
		return b != 0, true
	case int64:
		return b != 0, true
	case float64:
		return b != 0, true
	case string:
		switch strings.ToLower(strings.TrimSpace(b)) {
		case "true", "1", "yes", "y":
			return true, true
		case "false", "0", "no", "n":
			return false, true
		}
	}
	return false, false
}

func (m *EcommerceMerchantStorage) insertBranches(ctx context.Context, tx *sql.Tx, merchantID string, branches []model.BranchInformation) error {
	log := local_util.LoggerFromCtx(ctx, m.logger)

	if len(branches) == 0 {
		return nil
	}

	const insertBranchQ = `
INSERT INTO MERCHANT_BRANCHES (
	MERCHANT_ID,
	BRANCH_CODE,
	BRANCH_NAME,
	BRANCH_ADDRESS,
	BRANCH_OWNER,
	BRANCH_ACCOUNT_NUMBER,
	IS_ENABLED,
	IS_DELETED,
	CREATED_AT,
	LAST_MODIFIED_AT
) VALUES (
	HEXTORAW(:1), :2, :3, :4, :5, :6, 1, 0, SYSTIMESTAMP, SYSTIMESTAMP
)`

	for _, branch := range branches {
		if _, err := tx.ExecContext(
			ctx,
			insertBranchQ,
			merchantID,
			branch.BranchCode,
			branch.BranchName,
			nullString(branch.BranchAddress),
			nullString(branch.BranchOwner),
			branch.BranchAccountNumber,
		); err != nil {
			log.Errorf("[EcommerceMerchantRepo][insertBranches] insert failed: %v", err)
			return local_util.HandleDBError(err)
		}
	}

	return nil
}

func (m *EcommerceMerchantStorage) upsertBranches(ctx context.Context, tx *sql.Tx, merchantID string, branches []model.BranchInformation) error {
	log := local_util.LoggerFromCtx(ctx, m.logger)

	const markDeletedQ = `
UPDATE MERCHANT_BRANCHES
SET IS_DELETED = 1, DELETED_AT = SYSTIMESTAMP, LAST_MODIFIED_AT = SYSTIMESTAMP
WHERE MERCHANT_ID = HEXTORAW(:1) AND IS_DELETED = 0`

	if _, err := tx.ExecContext(ctx, markDeletedQ, merchantID); err != nil {
		log.Errorf("[EcommerceMerchantRepo][upsertBranches] mark delete failed: %v", err)
		return local_util.HandleDBError(err)
	}

	return m.insertBranches(ctx, tx, merchantID, branches)
}

func (m *EcommerceMerchantStorage) updateBranchesByID(ctx context.Context, tx *sql.Tx, merchantID string, branches []model.BranchInformation) error {
	log := local_util.LoggerFromCtx(ctx, m.logger)

	if len(branches) == 0 {
		return nil
	}

	const updateBranchQ = `
UPDATE MERCHANT_BRANCHES
SET
	BRANCH_CODE = :1,
	BRANCH_NAME = :2,
	BRANCH_ADDRESS = :3,
	BRANCH_OWNER = :4,
	BRANCH_ACCOUNT_NUMBER = :5,
	LAST_MODIFIED_AT = SYSTIMESTAMP
WHERE ID = HEXTORAW(:6) AND MERCHANT_ID = HEXTORAW(:7) AND IS_DELETED = 0`

	const branchExistsQ = `
SELECT BRANCH_CODE
FROM MERCHANT_BRANCHES
WHERE ID = HEXTORAW(:1) AND MERCHANT_ID = HEXTORAW(:2) AND IS_DELETED = 0`

	const branchCodeConflictQ = `
SELECT 1
FROM MERCHANT_BRANCHES
WHERE MERCHANT_ID = HEXTORAW(:1)
  AND UPPER(TRIM(BRANCH_CODE)) = UPPER(TRIM(:2))
  AND ID <> HEXTORAW(:3)
  AND IS_DELETED = 0
FETCH FIRST 1 ROWS ONLY`

	for _, branch := range branches {
		branchID := strings.TrimSpace(fmt.Sprint(branch.ID))
		if branchID == "" {
			// If branch id is not provided on update, treat as a new branch for this merchant.
			if err := m.insertBranches(ctx, tx, merchantID, []model.BranchInformation{branch}); err != nil {
				return err
			}
			continue
		}

		var existingCode string
		if err := tx.QueryRowContext(ctx, branchExistsQ, branchID, merchantID).Scan(&existingCode); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errors.New(localization.ErrorEcommerceMerchantNotFound.Code)
			}
			log.Errorf("[EcommerceMerchantRepo][updateBranchesByID] fetch existing branch failed: %v", err)
			return local_util.HandleDBError(err)
		}

		finalBranchCode := strings.TrimSpace(branch.BranchCode)
		if finalBranchCode == "" {
			finalBranchCode = existingCode
		}

		if !strings.EqualFold(strings.TrimSpace(existingCode), finalBranchCode) {
			var conflict int
			err := tx.QueryRowContext(ctx, branchCodeConflictQ, merchantID, finalBranchCode, branchID).Scan(&conflict)
			if err == nil {
				return errors.New(localization.ErrorCodeAlreadyExist.Code)
			}
			if !errors.Is(err, sql.ErrNoRows) {
				log.Errorf("[EcommerceMerchantRepo][updateBranchesByID] check branch code conflict failed: %v", err)
				return local_util.HandleDBError(err)
			}
		}

		res, err := tx.ExecContext(
			ctx,
			updateBranchQ,
			finalBranchCode,
			branch.BranchName,
			nullString(branch.BranchAddress),
			nullString(branch.BranchOwner),
			branch.BranchAccountNumber,
			branchID,
			merchantID,
		)
		if err != nil {
			log.Errorf("[EcommerceMerchantRepo][updateBranchesByID] update failed: %v", err)
			return local_util.HandleDBError(err)
		}

		rows, _ := res.RowsAffected()
		if rows == 0 {
			return errors.New(localization.ErrorEcommerceMerchantNotFound.Code)
		}
	}

	return nil
}

func (m *EcommerceMerchantStorage) getBranchesByMerchantID(ctx context.Context, merchantID string) ([]model.BranchInformation, error) {
	log := local_util.LoggerFromCtx(ctx, m.logger)

	const q = `
SELECT RAWTOHEX(ID), BRANCH_CODE, BRANCH_NAME, BRANCH_ADDRESS, BRANCH_OWNER, BRANCH_ACCOUNT_NUMBER, IS_ENABLED
FROM MERCHANT_BRANCHES
WHERE MERCHANT_ID = HEXTORAW(:1) AND IS_DELETED = 0
ORDER BY CREATED_AT ASC`

	rows, err := m.db.QueryContext(ctx, q, merchantID)
	if err != nil {
		log.Errorf("[EcommerceMerchantRepo][getBranchesByMerchantID] query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer rows.Close()

	branches := make([]model.BranchInformation, 0)
	for rows.Next() {
		var b model.BranchInformation
		var address, owner sql.NullString
		if err := rows.Scan(&b.ID, &b.BranchCode, &b.BranchName, &address, &owner, &b.BranchAccountNumber, &b.IsEnabled); err != nil {
			log.Errorf("[EcommerceMerchantRepo][getBranchesByMerchantID] scan failed: %v", err)
			return nil, local_util.HandleDBError(err)
		}
		if address.Valid {
			b.BranchAddress = address.String
		}
		if owner.Valid {
			b.BranchOwner = owner.String
		}
		branches = append(branches, b)
	}

	return branches, nil
}

func (m *EcommerceMerchantStorage) scanMerchant(row interface {
	Scan(dest ...interface{}) error
}) (*model.EcommerceMerchant, error) {
	merchant := &model.EcommerceMerchant{}
	var enabled, isDeleted int
	var deletedAt sql.NullTime

	err := row.Scan(
		&merchant.ID,
		&merchant.Code,
		&merchant.MerchantName,
		&merchant.SettlementMethod,
		&merchant.BankAccountNumber,
		&merchant.Email,
		&merchant.PhoneNumber,
		&enabled,
		&isDeleted,
		&merchant.CreatedAt,
		&merchant.UpdatedAt,
		&deletedAt,
	)
	if err != nil {
		return nil, err
	}

	merchant.Enabled = enabled == 1
	merchant.IsDeleted = isDeleted == 1
	if deletedAt.Valid {
		merchant.DeletedAt = &deletedAt.Time
	}

	return merchant, nil
}

func (m *EcommerceMerchantStorage) Create(ctx context.Context, merchant *model.EcommerceMerchant) (*model.EcommerceMerchant, error) {
	log := local_util.LoggerFromCtx(ctx, m.logger)

	const insertMerchantQ = `
INSERT INTO MERCHANTS (
	MERCHANT_ACCOUNT_NUMBER,
	MERCHANT_CODE,
	MERCHANT_NAME,
	SETTLEMENT_METHOD,
	MERCHANT_TYPE,
	CONTACT_EMAIL,
	CONTACT_PHONE,
	IS_ENABLED,
	IS_DELETED,
	CREATED_AT,
	LAST_MODIFIED_AT
) VALUES (
	:1, :2, :3, :4, 'ECOMMERCE', :5, :6, :7, :8, SYSTIMESTAMP, SYSTIMESTAMP
)
RETURNING RAWTOHEX(ID) INTO :9`

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		log.Errorf("[EcommerceMerchantRepo][Create] begin tx failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = tx.Rollback() }()

	created := *merchant
	var merchantID string

	if _, err := tx.ExecContext(
		ctx,
		insertMerchantQ,
		merchant.BankAccountNumber,
		merchant.Code,
		merchant.MerchantName,
		merchant.SettlementMethod,
		nullString(merchant.Email),
		nullString(merchant.PhoneNumber),
		boolToOracleNumber(merchant.Enabled),
		boolToOracleNumber(merchant.IsDeleted),
		sql.Out{Dest: &merchantID},
	); err != nil {
		log.Errorf("[EcommerceMerchantRepo][Create] insert merchant failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	if err := m.insertBranches(ctx, tx, merchantID, merchant.Branches); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("[EcommerceMerchantRepo][Create] commit failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// Get created merchant with branches to return in response and set in context metadata for potential use in handlers.
	createdM, err := m.FindByID(ctx, merchantID)
	if err != nil {
		log.Errorf("[EcommerceMerchantRepo][Create] find created merchant failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	branches := make([]types.B, 0, len(createdM.Branches))
	for _, b := range createdM.Branches {
		branches = append(branches, types.B{ID: b.ID, Code: b.BranchCode})
	}

	types.SetMerchant(ctx, &types.Merchant{ID: createdM.ID, Branches: branches})

	return &created, nil
}

func (m *EcommerceMerchantStorage) Update(ctx context.Context, id string, merchant *model.EcommerceMerchant) error {
	log := local_util.LoggerFromCtx(ctx, m.logger)

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		log.Errorf("[EcommerceMerchantRepo][Update] begin tx failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = tx.Rollback() }()

	updateParts := make([]string, 0, 6)
	args := make([]interface{}, 0, 8)
	idx := 1

	if strings.TrimSpace(merchant.MerchantName) != "" {
		updateParts = append(updateParts, fmt.Sprintf("MERCHANT_NAME = :%d", idx))
		args = append(args, merchant.MerchantName)
		idx++
	}
	if strings.TrimSpace(merchant.SettlementMethod) != "" {
		updateParts = append(updateParts, fmt.Sprintf("SETTLEMENT_METHOD = :%d", idx))
		args = append(args, merchant.SettlementMethod)
		idx++
	}
	if strings.TrimSpace(merchant.BankAccountNumber) != "" {
		updateParts = append(updateParts, fmt.Sprintf("MERCHANT_ACCOUNT_NUMBER = :%d", idx))
		args = append(args, merchant.BankAccountNumber)
		idx++
	}
	if strings.TrimSpace(merchant.Email) != "" {
		updateParts = append(updateParts, fmt.Sprintf("CONTACT_EMAIL = :%d", idx))
		args = append(args, merchant.Email)
		idx++
	}
	if strings.TrimSpace(merchant.PhoneNumber) != "" {
		updateParts = append(updateParts, fmt.Sprintf("CONTACT_PHONE = :%d", idx))
		args = append(args, merchant.PhoneNumber)
		idx++
	}

	merchantFound := false
	if len(updateParts) > 0 {
		updateParts = append(updateParts, "LAST_MODIFIED_AT = SYSTIMESTAMP")
		query := fmt.Sprintf(
			"UPDATE MERCHANTS SET %s WHERE ID = HEXTORAW(:%d) AND IS_DELETED = 0",
			strings.Join(updateParts, ", "),
			idx,
		)
		args = append(args, id)

		res, err := tx.ExecContext(ctx, query, args...)
		if err != nil {
			log.Errorf("[EcommerceMerchantRepo][Update] update merchant failed: %v", err)
			return local_util.HandleDBError(err)
		}

		rows, _ := res.RowsAffected()
		if rows == 0 {
			return errors.New(localization.ErrorEcommerceMerchantNotFound.Code)
		}
		merchantFound = true
	}

	// If no merchant fields were supplied, still ensure the merchant exists before branch updates.
	if !merchantFound {
		const existsQ = `SELECT 1 FROM MERCHANTS WHERE ID = HEXTORAW(:1) AND IS_DELETED = 0`
		var exists int
		if err := tx.QueryRowContext(ctx, existsQ, id).Scan(&exists); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errors.New(localization.ErrorEcommerceMerchantNotFound.Code)
			}
			return local_util.HandleDBError(err)
		}
	}

	if len(merchant.Branches) > 0 {
		if err := m.updateBranchesByID(ctx, tx, id, merchant.Branches); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("[EcommerceMerchantRepo][Update] commit failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (m *EcommerceMerchantStorage) Delete(ctx context.Context, id string) error {
	log := local_util.LoggerFromCtx(ctx, m.logger)

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		log.Errorf("[EcommerceMerchantRepo][Delete] begin tx failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = tx.Rollback() }()

	const deleteMerchantQ = `
UPDATE MERCHANTS
SET
	IS_DELETED = 1,
	DELETED_AT = SYSTIMESTAMP,
	LAST_MODIFIED_AT = SYSTIMESTAMP
WHERE ID = HEXTORAW(:1) AND IS_DELETED = 0`

	res, err := tx.ExecContext(ctx, deleteMerchantQ, id)
	if err != nil {
		log.Errorf("[EcommerceMerchantRepo][Delete] merchant delete failed: %v", err)
		return local_util.HandleDBError(err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New(localization.ErrorEcommerceMerchantNotFound.Code)
	}

	const deleteBranchesQ = `
UPDATE MERCHANT_BRANCHES
SET
	IS_DELETED = 1,
	DELETED_AT = SYSTIMESTAMP,
	LAST_MODIFIED_AT = SYSTIMESTAMP
WHERE MERCHANT_ID = HEXTORAW(:1) AND IS_DELETED = 0`
	if _, err := tx.ExecContext(ctx, deleteBranchesQ, id); err != nil {
		log.Errorf("[EcommerceMerchantRepo][Delete] branch delete failed: %v", err)
		return local_util.HandleDBError(err)
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("[EcommerceMerchantRepo][Delete] commit failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (m *EcommerceMerchantStorage) DeleteBranch(ctx context.Context, id string) error {
	log := local_util.LoggerFromCtx(ctx, m.logger)

	const q = `
UPDATE MERCHANT_BRANCHES
SET
	IS_DELETED = 1,
	DELETED_AT = SYSTIMESTAMP,
	LAST_MODIFIED_AT = SYSTIMESTAMP
WHERE ID = HEXTORAW(:1) AND IS_DELETED = 0`

	res, err := m.db.ExecContext(ctx, q, id)
	if err != nil {
		log.Errorf("[EcommerceMerchantRepo][DeleteBranch] failed: %v", err)
		return local_util.HandleDBError(err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New(localization.ErrorEcommerceMerchantNotFound.Code)
	}

	return nil
}

func (m *EcommerceMerchantStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	log := local_util.LoggerFromCtx(ctx, m.logger)

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		log.Errorf("[EcommerceMerchantRepo][EnableOrDisable] begin tx failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = tx.Rollback() }()

	const merchantQ = `
UPDATE MERCHANTS
SET
	IS_ENABLED = :1,
	LAST_MODIFIED_AT = SYSTIMESTAMP
WHERE ID = HEXTORAW(:2) AND IS_DELETED = 0`

	res, err := tx.ExecContext(ctx, merchantQ, boolToOracleNumber(enable), id)
	if err != nil {
		log.Errorf("[EcommerceMerchantRepo][EnableOrDisable] merchant update failed: %v", err)
		return local_util.HandleDBError(err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New(localization.ErrorEcommerceMerchantNotFound.Code)
	}

	const branchQ = `
UPDATE MERCHANT_BRANCHES
SET
	IS_ENABLED = :1,
	LAST_MODIFIED_AT = SYSTIMESTAMP
WHERE MERCHANT_ID = HEXTORAW(:2) AND IS_DELETED = 0`
	if _, err := tx.ExecContext(ctx, branchQ, boolToOracleNumber(enable), id); err != nil {
		log.Errorf("[EcommerceMerchantRepo][EnableOrDisable] branch update failed: %v", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("[EcommerceMerchantRepo][EnableOrDisable] commit failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (m *EcommerceMerchantStorage) EnableOrDisableBranch(ctx context.Context, id string, enable bool) error {
	log := local_util.LoggerFromCtx(ctx, m.logger)

	const q = `
UPDATE MERCHANT_BRANCHES
SET
	IS_ENABLED = :1,
	LAST_MODIFIED_AT = SYSTIMESTAMP
WHERE ID = HEXTORAW(:2) AND IS_DELETED = 0`

	res, err := m.db.ExecContext(ctx, q, boolToOracleNumber(enable), id)
	if err != nil {
		log.Errorf("[EcommerceMerchantRepo][EnableOrDisableBranch] failed: %v", err)
		return local_util.HandleDBError(err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New(localization.ErrorEcommerceMerchantNotFound.Code)
	}

	return nil
}

func (m *EcommerceMerchantStorage) FindByID(ctx context.Context, id string) (*model.EcommerceMerchant, error) {
	log := local_util.LoggerFromCtx(ctx, m.logger)

	const q = `
SELECT
	RAWTOHEX(ID),
	MERCHANT_CODE,
	MERCHANT_NAME,
	SETTLEMENT_METHOD,
	MERCHANT_ACCOUNT_NUMBER,
	CONTACT_EMAIL,
	CONTACT_PHONE,
	IS_ENABLED,
	IS_DELETED,
	CREATED_AT,
	LAST_MODIFIED_AT,
	DELETED_AT
FROM MERCHANTS
WHERE ID = HEXTORAW(:1) AND IS_DELETED = 0`

	merchant, err := m.scanMerchant(m.db.QueryRowContext(ctx, q, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(localization.ErrorEcommerceMerchantNotFound.Code)
		}
		log.Errorf("[EcommerceMerchantRepo][FindByID] query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	branches, err := m.getBranchesByMerchantID(ctx, merchant.ID)
	if err != nil {
		return nil, err
	}
	merchant.Branches = branches

	return merchant, nil
}

func (m *EcommerceMerchantStorage) FindBranchByID(ctx context.Context, id string) (*model.BranchInformation, error) {
	log := local_util.LoggerFromCtx(ctx, m.logger)

	const q = `
SELECT RAWTOHEX(ID), BRANCH_CODE, BRANCH_NAME, BRANCH_ADDRESS, BRANCH_OWNER, BRANCH_ACCOUNT_NUMBER, IS_ENABLED
FROM MERCHANT_BRANCHES
WHERE ID = HEXTORAW(:1) AND IS_DELETED = 0`

	var branch model.BranchInformation
	var address, owner sql.NullString
	if err := m.db.QueryRowContext(ctx, q, id).Scan(
		&branch.ID,
		&branch.BranchCode,
		&branch.BranchName,
		&address,
		&owner,
		&branch.BranchAccountNumber,
		&branch.IsEnabled,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(localization.ErrorEcommerceMerchantBranchNotFound.Code)
		}
		log.Errorf("[EcommerceMerchantRepo][FindBranchByID] query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	if address.Valid {
		branch.BranchAddress = address.String
	}
	if owner.Valid {
		branch.BranchOwner = owner.String
	}

	return &branch, nil
}

func (m *EcommerceMerchantStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.EcommerceMerchant], error) {
	log := local_util.LoggerFromCtx(ctx, m.logger)

	limit := int64(defaultPageSize)
	page := int64(1)
	if filterParam.PerPage > 0 {
		limit = int64(filterParam.PerPage)
	}
	if filterParam.Page > 0 {
		page = int64(filterParam.Page)
	}
	offset := (page - 1) * limit

	clauses := []string{"IS_DELETED = 0", "MERCHANT_TYPE = 'ECOMMERCE'"}
	args := make([]interface{}, 0)

	search := strings.TrimSpace(filterParam.Search)
	if search != "" {
		clauses = append(clauses, `(LOWER(MERCHANT_CODE) LIKE '%' || LOWER(:search) || '%' OR LOWER(MERCHANT_NAME) LIKE '%' || LOWER(:search) || '%' OR LOWER(MERCHANT_ACCOUNT_NUMBER) LIKE '%' || LOWER(:search) || '%')`)
		args = append(args, sql.Named("search", search))
	}

	if filterParam.Filters != nil {
		if v, ok := filterParam.Filters["enabled"]; ok {
			if b, ok2 := parseOracleBool(v); ok2 {
				clauses = append(clauses, "IS_ENABLED = :enabled")
				args = append(args, sql.Named("enabled", boolToOracleNumber(b)))
			}
		}
		if v, ok := filterParam.Filters["merchant_code"]; ok {
			if val, ok2 := v.(string); ok2 && strings.TrimSpace(val) != "" {
				clauses = append(clauses, "LOWER(MERCHANT_CODE) = LOWER(:merchant_code)")
				args = append(args, sql.Named("merchant_code", val))
			}
		}
		if v, ok := filterParam.Filters["merchant_name"]; ok {
			if val, ok2 := v.(string); ok2 && strings.TrimSpace(val) != "" {
				clauses = append(clauses, "LOWER(MERCHANT_NAME) = LOWER(:merchant_name)")
				args = append(args, sql.Named("merchant_name", val))
			}
		}
		if v, ok := filterParam.Filters["bank_account_number"]; ok {
			if val, ok2 := v.(string); ok2 && strings.TrimSpace(val) != "" {
				clauses = append(clauses, "MERCHANT_ACCOUNT_NUMBER = :merchant_account_number")
				args = append(args, sql.Named("merchant_account_number", val))
			}
		}
	}

	where := strings.Join(clauses, " AND ")
	countQ := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", merchantsTable, where)
	var total int64
	if err := m.db.QueryRowContext(ctx, countQ, args...).Scan(&total); err != nil {
		log.Errorf("[EcommerceMerchantRepo][FindAllWithPagination] count failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	listQ := fmt.Sprintf(`
SELECT
	RAWTOHEX(ID),
	MERCHANT_CODE,
	MERCHANT_NAME,
	SETTLEMENT_METHOD,
	MERCHANT_ACCOUNT_NUMBER,
	CONTACT_EMAIL,
	CONTACT_PHONE,
	IS_ENABLED,
	IS_DELETED,
	CREATED_AT,
	LAST_MODIFIED_AT,
	DELETED_AT
FROM %s
WHERE %s
ORDER BY CREATED_AT DESC
OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`, merchantsTable, where)

	listArgs := append(args, sql.Named("offset", offset), sql.Named("limit", limit))
	rows, err := m.db.QueryContext(ctx, listQ, listArgs...)
	if err != nil {
		log.Errorf("[EcommerceMerchantRepo][FindAllWithPagination] list query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer rows.Close()

	list := make([]model.EcommerceMerchant, 0)
	for rows.Next() {
		merchant, err := m.scanMerchant(rows)
		if err != nil {
			log.Errorf("[EcommerceMerchantRepo][FindAllWithPagination] scan failed: %v", err)
			return nil, local_util.HandleDBError(err)
		}

		branches, err := m.getBranchesByMerchantID(ctx, merchant.ID)
		if err != nil {
			return nil, err
		}
		merchant.Branches = branches

		list = append(list, *merchant)
	}

	meta := local_util.BuildPaginationMeta(total, int(page), int(limit))
	return &types.PaginatedResponse[[]model.EcommerceMerchant]{
		Data: list,
		Meta: meta,
	}, nil
}

func (m *EcommerceMerchantStorage) FindOneO(ctx context.Context, data *types.CheckMerchant, opts *types.MiniAppMerchantExistOptions) (*model.EcommerceMerchant, error) {
	log := local_util.LoggerFromCtx(ctx, m.logger)

	clauses := []string{"IS_DELETED = 0", "MERCHANT_TYPE = 'ECOMMERCE'"}
	args := make([]interface{}, 0)
	paramIdx := 1

	orClauses := make([]string, 0, 2)
	if strings.TrimSpace(data.MerchantCode) != "" {
		orClauses = append(orClauses, fmt.Sprintf("LOWER(MERCHANT_CODE) = LOWER(:%d)", paramIdx))
		args = append(args, data.MerchantCode)
		paramIdx++
	}
	if strings.TrimSpace(data.BankAccountNumber) != "" {
		orClauses = append(orClauses, fmt.Sprintf("MERCHANT_ACCOUNT_NUMBER = :%d", paramIdx))
		args = append(args, data.BankAccountNumber)
		paramIdx++
	}

	if len(orClauses) == 0 {
		return nil, errors.New(localization.ErrorEcommerceMerchantNotFound.Code)
	}
	clauses = append(clauses, "("+strings.Join(orClauses, " OR ")+")")

	if opts != nil && strings.TrimSpace(opts.ExcludeID) != "" {
		clauses = append(clauses, fmt.Sprintf("RAWTOHEX(ID) != UPPER(:%d)", paramIdx))
		args = append(args, opts.ExcludeID)
	}

	const selectCols = `
SELECT
	RAWTOHEX(ID),
	MERCHANT_CODE,
	MERCHANT_NAME,
	SETTLEMENT_METHOD,
	MERCHANT_ACCOUNT_NUMBER,
	CONTACT_EMAIL,
	CONTACT_PHONE,
	IS_ENABLED,
	IS_DELETED,
	CREATED_AT,
	LAST_MODIFIED_AT,
	DELETED_AT
FROM MERCHANTS`

	q := selectCols + " WHERE " + strings.Join(clauses, " AND ") + " FETCH FIRST 1 ROWS ONLY"

	merchant, err := m.scanMerchant(m.db.QueryRowContext(ctx, q, args...))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(localization.ErrorEcommerceMerchantNotFound.Code)
		}
		log.Errorf("[EcommerceMerchantRepo][FindOneO] query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	branches, err := m.getBranchesByMerchantID(ctx, merchant.ID)
	if err != nil {
		return nil, err
	}
	merchant.Branches = branches

	return merchant, nil
}
