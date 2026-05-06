package ecommercemerchant

// import (
// 	"cbe-super-app-cps-action/internal/storage"
// 	"context"
// 	"database/sql"
// 	"errors"
// 	"fmt"
// 	"strings"
// 	"time"

// 	"cbe-super-app-cps-action/internal/constants/localization"
// 	"cbe-super-app-cps-action/internal/constants/types"
// 	local_util "cbe-super-app-cps-action/pkgs/utils"

// 	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
// 	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
// 	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
// 	"go.mongodb.org/mongo-driver/v2/bson"
// )

// type EcommerceMerchantStorage struct {
// 	cfg    *config.VaultConfig
// 	db     *sql.DB
// 	logger utils.Logger
// }

// func NewEcommerceMerchantRepository(cfg *config.VaultConfig, db *sql.DB, logger utils.Logger) storage.EcommerceMerchantRepository {
// 	return &EcommerceMerchantStorage{
// 		cfg:    cfg,
// 		db:     db,
// 		logger: logger,
// 	}
// }

// const merchantsTable = "MERCHANTS"

// func boolToOracleNumber(v bool) int {
// 	if v {
// 		return 1
// 	}
// 	return 0
// }

// func parseBoolFilter(v interface{}) (bool, bool) {
// 	switch t := v.(type) {
// 	case bool:
// 		return t, true
// 	case int:
// 		return t != 0, true
// 	case int32:
// 		return t != 0, true
// 	case int64:
// 		return t != 0, true
// 	case float64:
// 		return t != 0, true
// 	case string:
// 		switch strings.ToLower(strings.TrimSpace(t)) {
// 		case "true", "1", "yes", "y":
// 			return true, true
// 		case "false", "0", "no", "n":
// 			return false, true
// 		}
// 	}
// 	return false, false
// }

// func normalizeRawHex(id string) (string, bool) {
// 	s := strings.TrimSpace(id)
// 	s = strings.TrimPrefix(s, "0x")
// 	s = strings.TrimPrefix(s, "0X")
// 	s = strings.ToLower(s)
// 	if len(s) != 24 && len(s) != 32 {
// 		return "", false
// 	}
// 	for _, r := range s {
// 		if (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') {
// 			continue
// 		}
// 		return "", false
// 	}
// 	return s, true
// }

// func toObjectIDCompatibleHex(id string) string {
// 	if len(id) >= 24 {
// 		return id[:24]
// 	}
// 	return id
// }

// func (m *EcommerceMerchantStorage) resolveAccountID(ctx context.Context, accountNumber string) (*string, error) {
// 	acc := strings.TrimSpace(accountNumber)
// 	if acc == "" {
// 		return nil, nil
// 	}

// 	const q = `
// SELECT RAWTOHEX(ID)
// FROM ACCOUNTS
// WHERE ACCOUNT_NUMBER = :1
//   AND ROWNUM = 1`

// 	var accountID string
// 	if err := m.db.QueryRowContext(ctx, q, acc).Scan(&accountID); err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return nil, nil
// 		}
// 		return nil, local_util.HandleDBError(err)
// 	}

// 	return &accountID, nil
// }

// func (m *EcommerceMerchantStorage) Create(ctx context.Context, merchant *model.EcommerceMerchant) (*model.EcommerceMerchant, error) {
// 	accountID, err := m.resolveAccountID(ctx, merchant.BankAccountNumber)
// 	if err != nil {
// 		m.logger.Errorf("[EcommerceMerchantStorage][Create] resolve account failed: %v", err)
// 		return nil, err
// 	}

// 	now := time.Now()
// 	if merchant.CreatedAt.IsZero() {
// 		merchant.CreatedAt = now
// 	}
// 	if merchant.UpdatedAt.IsZero() {
// 		merchant.UpdatedAt = now
// 	}

// 	var insertedID string
// 	const q = `
// INSERT INTO MERCHANTS (
//   ACCOUNT_ID,
//   MERCHANT_CODE,
//   MERCHANT_NAME,
//   SETTLEMENT_METHOD,
//   MERCHANT_TYPE,
//   CONTACT_EMAIL,
//   CONTACT_PHONE,
//   IS_ENABLED,
//   IS_DELETED,
//   CREATED_AT,
//   LAST_MODIFIED_AT
// ) VALUES (
//   CASE WHEN :1 IS NULL THEN NULL ELSE HEXTORAW(:1) END,
//   :2,
//   :3,
//   :4,
//   'ECOMMERCE',
//   :5,
//   :6,
//   :7,
//   0,
//   :8,
//   :9
// )
// RETURNING RAWTOHEX(ID) INTO :10`

// 	_, err = m.db.ExecContext(
// 		ctx,
// 		q,
// 		accountID,
// 		strings.TrimSpace(merchant.Code),
// 		strings.TrimSpace(merchant.MerchantName),
// 		strings.TrimSpace(merchant.SettlementMethod),
// 		strings.TrimSpace(merchant.Email),
// 		strings.TrimSpace(merchant.PhoneNumber),
// 		boolToOracleNumber(merchant.Enabled),
// 		merchant.CreatedAt,
// 		merchant.UpdatedAt,
// 		sql.Out{Dest: &insertedID},
// 	)
// 	if err != nil {
// 		m.logger.Errorf("[EcommerceMerchantStorage][Create] insert failed: %v", err)
// 		return nil, local_util.HandleDBError(err)
// 	}

// 	if objID, convErr := bson.ObjectIDFromHex(toObjectIDCompatibleHex(strings.ToLower(insertedID))); convErr == nil {
// 		merchant.ID = objID
// 	}
// 	types.SetId(ctx, insertedID)
// 	return merchant, nil
// }

// func (m *EcommerceMerchantStorage) Update(ctx context.Context, id string, merchant *model.EcommerceMerchant) error {
// 	idHex, ok := normalizeRawHex(id)
// 	if !ok {
// 		return errors.New(localization.ErrorInvalidID.Code)
// 	}

// 	sets := []string{"LAST_MODIFIED_AT = SYSTIMESTAMP"}
// 	var args []interface{}

// 	if v := strings.TrimSpace(merchant.Code); v != "" {
// 		sets = append(sets, "MERCHANT_CODE = :merchant_code")
// 		args = append(args, sql.Named("merchant_code", v))
// 	}
// 	if v := strings.TrimSpace(merchant.MerchantName); v != "" {
// 		sets = append(sets, "MERCHANT_NAME = :merchant_name")
// 		args = append(args, sql.Named("merchant_name", v))
// 	}
// 	if v := strings.TrimSpace(merchant.SettlementMethod); v != "" {
// 		sets = append(sets, "SETTLEMENT_METHOD = :settlement_method")
// 		args = append(args, sql.Named("settlement_method", v))
// 	}
// 	if v := strings.TrimSpace(merchant.Email); v != "" {
// 		sets = append(sets, "CONTACT_EMAIL = :contact_email")
// 		args = append(args, sql.Named("contact_email", v))
// 	}
// 	if v := strings.TrimSpace(merchant.PhoneNumber); v != "" {
// 		sets = append(sets, "CONTACT_PHONE = :contact_phone")
// 		args = append(args, sql.Named("contact_phone", v))
// 	}
// 	if strings.TrimSpace(merchant.BankAccountNumber) != "" {
// 		accountID, err := m.resolveAccountID(ctx, merchant.BankAccountNumber)
// 		if err != nil {
// 			return err
// 		}
// 		sets = append(sets, "ACCOUNT_ID = CASE WHEN :account_id IS NULL THEN NULL ELSE HEXTORAW(:account_id) END")
// 		args = append(args, sql.Named("account_id", accountID))
// 	}

// 	if len(sets) == 1 {
// 		return errors.New(localization.ErrorNoDataProvided.Code)
// 	}

// 	q := fmt.Sprintf(`
// UPDATE %s
// SET %s
// WHERE IS_DELETED = 0
//   AND (
//     RAWTOHEX(ID) = UPPER(:id)
//     OR SUBSTR(RAWTOHEX(ID), 1, 24) = UPPER(:id)
//   )`, merchantsTable, strings.Join(sets, ", "))
// 	args = append(args, sql.Named("id", idHex))

// 	res, err := m.db.ExecContext(ctx, q, args...)
// 	if err != nil {
// 		m.logger.Errorf("[EcommerceMerchantStorage][Update] update failed: %v", err)
// 		return local_util.HandleDBError(err)
// 	}
// 	rows, _ := res.RowsAffected()
// 	if rows == 0 {
// 		return errors.New(localization.ErrorResourceNotFound.Code)
// 	}

// 	return nil
// }

// func (m *EcommerceMerchantStorage) Delete(ctx context.Context, id string) error {
// 	idHex, ok := normalizeRawHex(id)
// 	if !ok {
// 		return errors.New(localization.ErrorInvalidID.Code)
// 	}

// 	const q = `
// UPDATE MERCHANTS
// SET
//   IS_DELETED = 1,
//   LAST_MODIFIED_AT = SYSTIMESTAMP,
//   DELETED_AT = SYSTIMESTAMP
// WHERE IS_DELETED = 0
//   AND (
//     RAWTOHEX(ID) = UPPER(:1)
//     OR SUBSTR(RAWTOHEX(ID), 1, 24) = UPPER(:1)
//   )`

// 	res, err := m.db.ExecContext(ctx, q, idHex)
// 	if err != nil {
// 		m.logger.Errorf("[EcommerceMerchantStorage][Delete] delete failed: %v", err)
// 		return local_util.HandleDBError(err)
// 	}
// 	rows, _ := res.RowsAffected()
// 	if rows == 0 {
// 		return errors.New(localization.ErrorResourceNotFound.Code)
// 	}
// 	return nil
// }

// func (m *EcommerceMerchantStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
// 	idHex, ok := normalizeRawHex(id)
// 	if !ok {
// 		return errors.New(localization.ErrorInvalidID.Code)
// 	}

// 	const q = `
// UPDATE MERCHANTS
// SET
//   IS_ENABLED = :1,
//   LAST_MODIFIED_AT = SYSTIMESTAMP
// WHERE IS_DELETED = 0
//   AND (
//     RAWTOHEX(ID) = UPPER(:2)
//     OR SUBSTR(RAWTOHEX(ID), 1, 24) = UPPER(:2)
//   )`

// 	res, err := m.db.ExecContext(ctx, q, boolToOracleNumber(enable), idHex)
// 	if err != nil {
// 		m.logger.Errorf("[EcommerceMerchantStorage][EnableOrDisable] update failed: %v", err)
// 		return local_util.HandleDBError(err)
// 	}
// 	rows, _ := res.RowsAffected()
// 	if rows == 0 {
// 		return errors.New(localization.ErrorResourceNotFound.Code)
// 	}
// 	return nil
// }

// func (m *EcommerceMerchantStorage) FindByID(ctx context.Context, id string) (*model.EcommerceMerchant, error) {
// 	idHex, ok := normalizeRawHex(id)
// 	if !ok {
// 		return nil, errors.New(localization.ErrorInvalidID.Code)
// 	}

// 	const q = `
// SELECT
//   RAWTOHEX(M.ID),
//   NVL(M.MERCHANT_CODE, ''),
//   NVL(M.MERCHANT_NAME, ''),
//   NVL(M.SETTLEMENT_METHOD, ''),
//   NVL(M.CONTACT_EMAIL, ''),
//   NVL(M.CONTACT_PHONE, ''),
//   NVL(A.ACCOUNT_NUMBER, ''),
//   M.IS_ENABLED,
//   M.IS_DELETED,
//   M.CREATED_AT,
//   M.LAST_MODIFIED_AT,
//   M.DELETED_AT
// FROM MERCHANTS M
// LEFT JOIN ACCOUNTS A ON A.ID = M.ACCOUNT_ID
// WHERE M.IS_DELETED = 0
//   AND (
//     RAWTOHEX(M.ID) = UPPER(:1)
//     OR SUBSTR(RAWTOHEX(M.ID), 1, 24) = UPPER(:1)
//   )
// FETCH FIRST 1 ROWS ONLY`

// 	var (
// 		rawID                                                   string
// 		code, merchantName, settlementMethod, email, phone, acc sql.NullString
// 		enabledN, isDeletedN                                    int
// 		createdAt, updatedAt, deletedAt                         sql.NullTime
// 	)

// 	err := m.db.QueryRowContext(ctx, q, idHex).Scan(
// 		&rawID,
// 		&code,
// 		&merchantName,
// 		&settlementMethod,
// 		&email,
// 		&phone,
// 		&acc,
// 		&enabledN,
// 		&isDeletedN,
// 		&createdAt,
// 		&updatedAt,
// 		&deletedAt,
// 	)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return nil, errors.New(localization.ErrorResourceNotFound.Code)
// 		}
// 		m.logger.Errorf("[EcommerceMerchantStorage][FindByID] query failed: %v", err)
// 		return nil, local_util.HandleDBError(err)
// 	}

// 	out := &model.EcommerceMerchant{
// 		Code:              code.String,
// 		MerchantName:      merchantName.String,
// 		SettlementMethod:  settlementMethod.String,
// 		Email:             email.String,
// 		PhoneNumber:       phone.String,
// 		BankAccountNumber: acc.String,
// 		Enabled:           enabledN == 1,
// 		IsDeleted:         isDeletedN == 1,
// 		Branches:          []model.BranchInformation{},
// 	}
// 	if createdAt.Valid {
// 		out.CreatedAt = createdAt.Time
// 	}
// 	if updatedAt.Valid {
// 		out.UpdatedAt = updatedAt.Time
// 	}
// 	if deletedAt.Valid {
// 		t := deletedAt.Time
// 		out.DeletedAt = &t
// 	}

// 	if objID, convErr := bson.ObjectIDFromHex(toObjectIDCompatibleHex(strings.ToLower(rawID))); convErr == nil {
// 		out.ID = objID
// 	}
// 	return out, nil
// }

// func (m *EcommerceMerchantStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.EcommerceMerchant], error) {
// 	limit := int64(10)
// 	page := int64(1)
// 	if filterParam.PerPage > 0 {
// 		limit = int64(filterParam.PerPage)
// 	}
// 	if filterParam.Page > 0 {
// 		page = int64(filterParam.Page)
// 	}
// 	offset := (page - 1) * limit

// 	clauses := []string{"M.IS_DELETED = 0", "M.MERCHANT_TYPE = 'ECOMMERCE'"}
// 	var args []interface{}

// 	search := strings.TrimSpace(filterParam.Search)
// 	if search != "" {
// 		clauses = append(clauses, `(LOWER(M.MERCHANT_CODE) LIKE '%' || LOWER(:search) || '%' OR LOWER(M.MERCHANT_NAME) LIKE '%' || LOWER(:search) || '%' OR LOWER(NVL(M.CONTACT_EMAIL, '')) LIKE '%' || LOWER(:search) || '%' OR LOWER(NVL(M.CONTACT_PHONE, '')) LIKE '%' || LOWER(:search) || '%' OR LOWER(NVL(A.ACCOUNT_NUMBER, '')) LIKE '%' || LOWER(:search) || '%')`)
// 		args = append(args, sql.Named("search", search))
// 	}
// 	if filterParam.Filters != nil {
// 		if v, ok := filterParam.Filters["merchant_code"]; ok {
// 			if s, ok2 := v.(string); ok2 && strings.TrimSpace(s) != "" {
// 				clauses = append(clauses, "LOWER(M.MERCHANT_CODE) LIKE '%' || LOWER(:merchant_code) || '%'")
// 				args = append(args, sql.Named("merchant_code", strings.TrimSpace(s)))
// 			}
// 		}
// 		if v, ok := filterParam.Filters["merchant_name"]; ok {
// 			if s, ok2 := v.(string); ok2 && strings.TrimSpace(s) != "" {
// 				clauses = append(clauses, "LOWER(M.MERCHANT_NAME) LIKE '%' || LOWER(:merchant_name) || '%'")
// 				args = append(args, sql.Named("merchant_name", strings.TrimSpace(s)))
// 			}
// 		}
// 		if v, ok := filterParam.Filters["enabled"]; ok {
// 			if b, ok2 := parseBoolFilter(v); ok2 {
// 				clauses = append(clauses, "M.IS_ENABLED = :enabled")
// 				args = append(args, sql.Named("enabled", boolToOracleNumber(b)))
// 			}
// 		}
// 	}

// 	where := strings.Join(clauses, " AND ")
// 	countQ := fmt.Sprintf(`
// SELECT COUNT(1)
// FROM %s M
// LEFT JOIN ACCOUNTS A ON A.ID = M.ACCOUNT_ID
// WHERE %s`, merchantsTable, where)

// 	var total int64
// 	if err := m.db.QueryRowContext(ctx, countQ, args...).Scan(&total); err != nil {
// 		m.logger.Errorf("[EcommerceMerchantStorage][FindAllWithPagination] count failed: %v", err)
// 		return nil, local_util.HandleDBError(err)
// 	}

// 	listQ := fmt.Sprintf(`
// SELECT
//   RAWTOHEX(M.ID),
//   NVL(M.MERCHANT_CODE, ''),
//   NVL(M.MERCHANT_NAME, ''),
//   NVL(M.SETTLEMENT_METHOD, ''),
//   NVL(M.CONTACT_EMAIL, ''),
//   NVL(M.CONTACT_PHONE, ''),
//   NVL(A.ACCOUNT_NUMBER, ''),
//   M.IS_ENABLED,
//   M.IS_DELETED,
//   M.CREATED_AT,
//   M.LAST_MODIFIED_AT,
//   M.DELETED_AT
// FROM %s M
// LEFT JOIN ACCOUNTS A ON A.ID = M.ACCOUNT_ID
// WHERE %s
// ORDER BY M.CREATED_AT DESC
// OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`, merchantsTable, where)

// 	listArgs := append(args, sql.Named("offset", offset), sql.Named("limit", limit))
// 	rows, err := m.db.QueryContext(ctx, listQ, listArgs...)
// 	if err != nil {
// 		m.logger.Errorf("[EcommerceMerchantStorage][FindAllWithPagination] list query failed: %v", err)
// 		return nil, local_util.HandleDBError(err)
// 	}
// 	defer rows.Close()

// 	out := make([]model.EcommerceMerchant, 0)
// 	for rows.Next() {
// 		var (
// 			rawID                                                   string
// 			code, merchantName, settlementMethod, email, phone, acc sql.NullString
// 			enabledN, isDeletedN                                    int
// 			createdAt, updatedAt, deletedAt                         sql.NullTime
// 		)

// 		if err := rows.Scan(
// 			&rawID,
// 			&code,
// 			&merchantName,
// 			&settlementMethod,
// 			&email,
// 			&phone,
// 			&acc,
// 			&enabledN,
// 			&isDeletedN,
// 			&createdAt,
// 			&updatedAt,
// 			&deletedAt,
// 		); err != nil {
// 			m.logger.Errorf("[EcommerceMerchantStorage][FindAllWithPagination] scan failed: %v", err)
// 			return nil, local_util.HandleDBError(err)
// 		}

// 		item := model.EcommerceMerchant{
// 			Code:              code.String,
// 			MerchantName:      merchantName.String,
// 			SettlementMethod:  settlementMethod.String,
// 			Email:             email.String,
// 			PhoneNumber:       phone.String,
// 			BankAccountNumber: acc.String,
// 			Enabled:           enabledN == 1,
// 			IsDeleted:         isDeletedN == 1,
// 			Branches:          []model.BranchInformation{},
// 		}
// 		if createdAt.Valid {
// 			item.CreatedAt = createdAt.Time
// 		}
// 		if updatedAt.Valid {
// 			item.UpdatedAt = updatedAt.Time
// 		}
// 		if deletedAt.Valid {
// 			t := deletedAt.Time
// 			item.DeletedAt = &t
// 		}
// 		if objID, convErr := bson.ObjectIDFromHex(toObjectIDCompatibleHex(strings.ToLower(rawID))); convErr == nil {
// 			item.ID = objID
// 		}

// 		out = append(out, item)
// 	}

// 	meta := local_util.BuildPaginationMeta(total, int(page), int(limit))
// 	return &types.PaginatedResponse[[]model.EcommerceMerchant]{
// 		Data: out,
// 		Meta: meta,
// 	}, nil
// }

// func (m *EcommerceMerchantStorage) FindOne(ctx context.Context, filter bson.M) (*model.EcommerceMerchant, error) {
// 	if v, ok := filter["_id"]; ok {
// 		switch t := v.(type) {
// 		case string:
// 			return m.FindByID(ctx, t)
// 		case bson.ObjectID:
// 			return m.FindByID(ctx, t.Hex())
// 		}
// 	}
// 	if v, ok := filter["merchant_code"]; ok {
// 		code, _ := v.(string)
// 		if strings.TrimSpace(code) == "" {
// 			return nil, errors.New(localization.ErrorNoDataProvided.Code)
// 		}

// 		const q = `
// SELECT RAWTOHEX(ID)
// FROM MERCHANTS
// WHERE IS_DELETED = 0
//   AND LOWER(MERCHANT_CODE) = LOWER(:1)
// FETCH FIRST 1 ROWS ONLY`
// 		var idHex string
// 		if err := m.db.QueryRowContext(ctx, q, strings.TrimSpace(code)).Scan(&idHex); err != nil {
// 			if errors.Is(err, sql.ErrNoRows) {
// 				return nil, errors.New(localization.ErrorResourceNotFound.Code)
// 			}
// 			return nil, local_util.HandleDBError(err)
// 		}
// 		return m.FindByID(ctx, idHex)
// 	}

// 	return nil, errors.New(localization.ErrorNoDataProvided.Code)
// }
