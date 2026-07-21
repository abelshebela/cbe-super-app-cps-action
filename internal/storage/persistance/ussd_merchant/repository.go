package ussd_merchant

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	ussd_merchant_dto "cbe-super-app-cps-action/internal/constants/dto/ussd_merchant"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	merchantsTable     = "MERCHANTS"
	ussdMerchantsTable = "USSD_MERCHANT_SERVICE"
)

const selectJoin = `
SELECT
	RAWTOHEX(M.ID),
	M.MERCHANT_CODE,
	U.CREDENTIAL,
	M.MERCHANT_NAME,
	M.SETTLEMENT_METHOD,
	M.CONTACT_PHONE,
	M.CONTACT_EMAIL,
	RAWTOHEX(U.SERVICE_ID),
	M.MERCHANT_ACCOUNT_NUMBER,
	U.LOGO,
	M.IS_ENABLED,
	M.IS_DELETED,
	M.CREATED_AT,
	M.LAST_MODIFIED_AT,
	M.DELETED_AT
FROM MERCHANTS M
JOIN USSD_MERCHANT_SERVICE U ON U.MERCHANT_ID = M.ID AND U.IS_DELETED = 0`

type UssdMerchantRepository struct {
	db     *sql.DB
	logger utils.Logger
}

func NewUssdMerchant(db *sql.DB, logger utils.Logger) storage.UssdMerchantRepository {
	return &UssdMerchantRepository{
		db:     db,
		logger: logger,
	}
}

func (u *UssdMerchantRepository) Create(ctx context.Context, data imodel.UssdMerchant) error {
	log := local_util.LoggerFromCtx(ctx, u.logger)

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		log.Errorf("[UssdMerchantRepo][Create] begin tx failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = tx.Rollback() }()

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
	:1, :2, :3, :4, 'USSD_PUSH', :5, :6, :7, 0, SYSTIMESTAMP, SYSTIMESTAMP
)
RETURNING RAWTOHEX(ID) INTO :8`

	enabled := 0
	if data.Enabled {
		enabled = 1
	}

	var merchantID string
	if _, err := tx.ExecContext(
		ctx,
		insertMerchantQ,
		nullString(data.AccountNumber),
		data.MerchantCode,
		data.Name,
		string(data.SettlementMethod),
		nullString(data.Email),
		nullString(data.PhoneNumber),
		enabled,
		sql.Out{Dest: &merchantID},
	); err != nil {
		log.Errorf("[UssdMerchantRepo][Create] insert merchant failed: %v", err)
		return local_util.HandleDBError(err)
	}

	const insertUssdQ = `
INSERT INTO USSD_MERCHANT_SERVICE (
	MERCHANT_ID,
	SERVICE_ID,
	LOGO,
	CREDENTIAL,
	IS_DELETED,
	CREATED_AT,
	LAST_MODIFIED_AT
) VALUES (
	HEXTORAW(:1), HEXTORAW(:2), :3, :4, 0, SYSTIMESTAMP, SYSTIMESTAMP
)`

	if _, err := tx.ExecContext(
		ctx,
		insertUssdQ,
		merchantID,
		nullString(data.Service),
		nullString(data.Logo),
		nullString(data.Credential),
	); err != nil {
		log.Errorf("[UssdMerchantRepo][Create] insert ussd_merchant failed: %v", err)
		return local_util.HandleDBError(err)
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("[UssdMerchantRepo][Create] commit failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (u *UssdMerchantRepository) Update(ctx context.Context, id string, update bson.M) error {
	log := local_util.LoggerFromCtx(ctx, u.logger)

	if len(update) == 0 {
		return nil
	}

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		log.Errorf("[UssdMerchantRepo][Update] begin tx failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = tx.Rollback() }()

	// Split bson.M keys into merchant vs ussd_merchant fields.
	merchantUpdates := make([]string, 0)
	ussdUpdates := make([]string, 0)
	merchantArgs := make([]interface{}, 0)
	ussdArgs := make([]interface{}, 0)
	mIdx := 1
	uIdx := 1

	for key, val := range update {
		switch key {
		case "name":
			merchantUpdates = append(merchantUpdates, fmt.Sprintf("MERCHANT_NAME = :%d", mIdx))
			merchantArgs = append(merchantArgs, val)
			mIdx++
		case "settlement_method":
			merchantUpdates = append(merchantUpdates, fmt.Sprintf("SETTLEMENT_METHOD = :%d", mIdx))
			merchantArgs = append(merchantArgs, val)
			mIdx++
		case "phone_number":
			merchantUpdates = append(merchantUpdates, fmt.Sprintf("CONTACT_PHONE = :%d", mIdx))
			merchantArgs = append(merchantArgs, val)
			mIdx++
		case "email":
			merchantUpdates = append(merchantUpdates, fmt.Sprintf("CONTACT_EMAIL = :%d", mIdx))
			merchantArgs = append(merchantArgs, val)
			mIdx++
		case "account_number":
			merchantUpdates = append(merchantUpdates, fmt.Sprintf("MERCHANT_ACCOUNT_NUMBER = :%d", mIdx))
			merchantArgs = append(merchantArgs, val)
			mIdx++
		case "enabled":
			v, ok := val.(bool)
			if !ok {
				log.Errorf("[UssdMerchantRepo][Update] enabled is not bool: %v", val)
				return errors.New(localization.ErrorInvalidID.Code)
			}
			n := 0
			if v {
				n = 1
			}
			merchantUpdates = append(merchantUpdates, fmt.Sprintf("IS_ENABLED = :%d", mIdx))
			merchantArgs = append(merchantArgs, n)
			mIdx++
		case "service":
			ussdUpdates = append(ussdUpdates, fmt.Sprintf("SERVICE_ID = HEXTORAW(:%d)", uIdx))
			ussdArgs = append(ussdArgs, val)
			uIdx++
		case "logo":
			ussdUpdates = append(ussdUpdates, fmt.Sprintf("LOGO = :%d", uIdx))
			ussdArgs = append(ussdArgs, val)
			uIdx++
		case "credential":
			ussdUpdates = append(ussdUpdates, fmt.Sprintf("CREDENTIAL = :%d", uIdx))
			ussdArgs = append(ussdArgs, val)
			uIdx++
		default:
			log.Errorf("[UssdMerchantRepo][Update] unknown bson key: %s", key)
			return errors.New(localization.ErrorInvalidID.Code)
		}
	}

	merchantFound := false
	if len(merchantUpdates) > 0 {
		merchantUpdates = append(merchantUpdates, "LAST_MODIFIED_AT = SYSTIMESTAMP")
		query := fmt.Sprintf(
			"UPDATE %s SET %s WHERE ID = HEXTORAW(:%d) AND IS_DELETED = 0",
			merchantsTable,
			strings.Join(merchantUpdates, ", "),
			mIdx,
		)
		merchantArgs = append(merchantArgs, id)

		res, err := tx.ExecContext(ctx, query, merchantArgs...)
		if err != nil {
			log.Errorf("[UssdMerchantRepo][Update] update merchant failed: %v", err)
			return local_util.HandleDBError(err)
		}
		rows, _ := res.RowsAffected()
		if rows == 0 {
			return errors.New(localization.ErrorMerchantNotFound.Code)
		}
		merchantFound = true
	}

	if len(ussdUpdates) > 0 {
		ussdUpdates = append(ussdUpdates, "LAST_MODIFIED_AT = SYSTIMESTAMP")
		query := fmt.Sprintf(
			"UPDATE %s SET %s WHERE MERCHANT_ID = HEXTORAW(:%d) AND IS_DELETED = 0",
			ussdMerchantsTable,
			strings.Join(ussdUpdates, ", "),
			uIdx,
		)
		ussdArgs = append(ussdArgs, id)

		res, err := tx.ExecContext(ctx, query, ussdArgs...)
		if err != nil {
			log.Errorf("[UssdMerchantRepo][Update] update ussd_merchant failed: %v", err)
			return local_util.HandleDBError(err)
		}
		rows, _ := res.RowsAffected()
		if rows == 0 && !merchantFound {
			return errors.New(localization.ErrorMerchantNotFound.Code)
		}
	}

	if !merchantFound && len(ussdUpdates) == 0 {
		// Verify merchant exists.
		const existsQ = `SELECT 1 FROM MERCHANTS WHERE ID = HEXTORAW(:1) AND IS_DELETED = 0`
		var exists int
		if err := tx.QueryRowContext(ctx, existsQ, id).Scan(&exists); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errors.New(localization.ErrorMerchantNotFound.Code)
			}
			return local_util.HandleDBError(err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("[UssdMerchantRepo][Update] commit failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (u *UssdMerchantRepository) Delete(ctx context.Context, id string) error {
	log := local_util.LoggerFromCtx(ctx, u.logger)

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		log.Errorf("[UssdMerchantRepo][Delete] begin tx failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = tx.Rollback() }()

	const deleteUssdQ = `
UPDATE USSD_MERCHANT_SERVICE
SET IS_DELETED = 1, DELETED_AT = SYSTIMESTAMP, LAST_MODIFIED_AT = SYSTIMESTAMP
WHERE MERCHANT_ID = HEXTORAW(:1) AND IS_DELETED = 0`

	if _, err := tx.ExecContext(ctx, deleteUssdQ, id); err != nil {
		log.Errorf("[UssdMerchantRepo][Delete] ussd_merchant delete failed: %v", err)
		return local_util.HandleDBError(err)
	}

	const deleteMerchantQ = `
UPDATE MERCHANTS
SET IS_DELETED = 1, DELETED_AT = SYSTIMESTAMP, LAST_MODIFIED_AT = SYSTIMESTAMP
WHERE ID = HEXTORAW(:1) AND IS_DELETED = 0`

	res, err := tx.ExecContext(ctx, deleteMerchantQ, id)
	if err != nil {
		log.Errorf("[UssdMerchantRepo][Delete] merchant delete failed: %v", err)
		return local_util.HandleDBError(err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New(localization.ErrorMerchantNotFound.Code)
	}

	if err := tx.Commit(); err != nil {
		log.Errorf("[UssdMerchantRepo][Delete] commit failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (u *UssdMerchantRepository) FindById(ctx context.Context, id string) (ussd_merchant_dto.UssdMerchantResponse, error) {
	log := local_util.LoggerFromCtx(ctx, u.logger)

	query := selectJoin + "\nWHERE M.ID = HEXTORAW(:1) AND M.IS_DELETED = 0"

	resp, err := scanResponse(u.db.QueryRowContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ussd_merchant_dto.UssdMerchantResponse{}, errors.New(localization.ErrorMerchantNotFound.Code)
		}
		log.Errorf("[UssdMerchantRepo][FindById] query failed: %v", err)
		return ussd_merchant_dto.UssdMerchantResponse{}, local_util.HandleDBError(err)
	}

	return resp, nil
}

func (u *UssdMerchantRepository) Find(ctx context.Context, filter bson.M) (ussd_merchant_dto.UssdMerchantResponse, error) {
	log := local_util.LoggerFromCtx(ctx, u.logger)

	clauses := []string{"M.IS_DELETED = 0"}
	args := make([]interface{}, 0)
	idx := 1

	for key, val := range filter {
		switch key {
		case "service":
			clauses = append(clauses, fmt.Sprintf("U.SERVICE_ID = HEXTORAW(:%d)", idx))
		case "enabled":
			clauses = append(clauses, fmt.Sprintf("M.IS_ENABLED = :%d", idx))
		case "name":
			clauses = append(clauses, fmt.Sprintf("M.MERCHANT_NAME = :%d", idx))
		case "settlement_method":
			clauses = append(clauses, fmt.Sprintf("M.SETTLEMENT_METHOD = :%d", idx))
		case "phone_number":
			clauses = append(clauses, fmt.Sprintf("M.CONTACT_PHONE = :%d", idx))
		case "email":
			clauses = append(clauses, fmt.Sprintf("M.CONTACT_EMAIL = :%d", idx))
		case "account_number":
			clauses = append(clauses, fmt.Sprintf("M.MERCHANT_ACCOUNT_NUMBER = :%d", idx))
		case "merchant_code":
			clauses = append(clauses, fmt.Sprintf("M.MERCHANT_CODE = :%d", idx))
		default:
			log.Errorf("[UssdMerchantRepo][Find] unknown bson key: %s", key)
			return ussd_merchant_dto.UssdMerchantResponse{}, errors.New(localization.ErrorMerchantNotFound.Code)
		}
		args = append(args, val)
		idx++
	}

	query := selectJoin + "\nWHERE " + strings.Join(clauses, " AND ") + "\nFETCH FIRST 1 ROWS ONLY"

	resp, err := scanResponse(u.db.QueryRowContext(ctx, query, args...))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ussd_merchant_dto.UssdMerchantResponse{}, errors.New(localization.ErrorMerchantNotFound.Code)
		}
		log.Errorf("[UssdMerchantRepo][Find] query failed: %v", err)
		return ussd_merchant_dto.UssdMerchantResponse{}, local_util.HandleDBError(err)
	}

	return resp, nil
}

func (u *UssdMerchantRepository) FindByOr(ctx context.Context, phone, email, accountNumber string) (imodel.UssdMerchant, error) {
	log := local_util.LoggerFromCtx(ctx, u.logger)

	orClauses := make([]string, 0, 3)
	args := make([]interface{}, 0, 3)
	idx := 1

	if phone != "" {
		orClauses = append(orClauses, fmt.Sprintf("UPPER(TRIM(M.CONTACT_PHONE)) = UPPER(TRIM(:%d))", idx))
		args = append(args, phone)
		idx++
	}
	if accountNumber != "" {
		orClauses = append(orClauses, fmt.Sprintf("M.MERCHANT_ACCOUNT_NUMBER = :%d", idx))
		args = append(args, accountNumber)
		idx++
	}
	if email != "" {
		orClauses = append(orClauses, fmt.Sprintf("UPPER(TRIM(M.CONTACT_EMAIL)) = UPPER(TRIM(:%d))", idx))
		args = append(args, email)
		idx++
	}

	if len(orClauses) == 0 {
		return imodel.UssdMerchant{}, errors.New(localization.ErrorMerchantNotFound.Code)
	}

	query := selectJoin + "\nWHERE (" + strings.Join(orClauses, " OR ") + ") AND M.IS_DELETED = 0\nFETCH FIRST 1 ROWS ONLY"

	merchant, err := scanMerchant(u.db.QueryRowContext(ctx, query, args...))
	if err != nil {
		log.Errorf("[UssdMerchantRepo][FindByOr] query failed: %v", err)
		return imodel.UssdMerchant{}, local_util.HandleDBError(err)
	}

	return merchant, nil
}

func (u *UssdMerchantRepository) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (types.PaginatedResponse[[]ussd_merchant_dto.UssdMerchantResponse], error) {
	log := local_util.LoggerFromCtx(ctx, u.logger)

	limit := int64(10)
	page := int64(1)
	if filterParam.PerPage > 0 {
		limit = int64(filterParam.PerPage)
	}
	if filterParam.Page > 0 {
		page = int64(filterParam.Page)
	}
	offset := (page - 1) * limit

	clauses := []string{"M.IS_DELETED = 0"}
	args := make([]interface{}, 0)

	search := strings.TrimSpace(filterParam.Search)
	if search != "" {
		clauses = append(clauses, `(LOWER(M.MERCHANT_CODE) LIKE '%' || LOWER(:search) || '%' OR LOWER(M.MERCHANT_NAME) LIKE '%' || LOWER(:search) || '%' OR LOWER(M.SETTLEMENT_METHOD) LIKE '%' || LOWER(:search) || '%' OR LOWER(M.CONTACT_PHONE) LIKE '%' || LOWER(:search) || '%' OR LOWER(M.MERCHANT_ACCOUNT_NUMBER) LIKE '%' || LOWER(:search) || '%')`)
		args = append(args, sql.Named("search", search))
	}

	allowedFilters := map[string]string{
		"enabled":           "M.IS_ENABLED",
		"merchant_code":     "M.MERCHANT_CODE",
		"name":              "M.MERCHANT_NAME",
		"settlement_method": "M.SETTLEMENT_METHOD",
		"phone_number":      "M.CONTACT_PHONE",
		"accountNumber":     "M.MERCHANT_ACCOUNT_NUMBER",
		"account_number":    "M.MERCHANT_ACCOUNT_NUMBER",
	}

	if filterParam.Filters != nil {
		for key, col := range allowedFilters {
			v, ok := filterParam.Filters[key]
			if !ok {
				continue
			}

			if key == "enabled" {
				if b, ok2 := parseOracleBool(v); ok2 {
					clauses = append(clauses, fmt.Sprintf("%s = :%s", col, key))
					args = append(args, sql.Named(key, boolToOracleNumber(b)))
				}
				continue
			}

			if val, ok2 := v.(string); ok2 && strings.TrimSpace(val) != "" {
				clauses = append(clauses, fmt.Sprintf("LOWER(%s) = LOWER(:%s)", col, key))
				args = append(args, sql.Named(key, val))
			}
		}
	}

	where := strings.Join(clauses, " AND ")

	countQ := fmt.Sprintf("SELECT COUNT(*) FROM %s M JOIN %s U ON U.MERCHANT_ID = M.ID AND U.IS_DELETED = 0 WHERE %s", merchantsTable, ussdMerchantsTable, where)
	var total int64
	if err := u.db.QueryRowContext(ctx, countQ, args...).Scan(&total); err != nil {
		log.Errorf("[UssdMerchantRepo][FindAllWithPagination] count failed: %v", err)
		return types.PaginatedResponse[[]ussd_merchant_dto.UssdMerchantResponse]{}, local_util.HandleDBError(err)
	}

	listQ := selectJoin + "\nWHERE " + where + "\nORDER BY M.CREATED_AT DESC\nOFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY"

	listArgs := append(args, sql.Named("offset", offset), sql.Named("limit", limit))
	rows, err := u.db.QueryContext(ctx, listQ, listArgs...)
	if err != nil {
		log.Errorf("[UssdMerchantRepo][FindAllWithPagination] list query failed: %v", err)
		return types.PaginatedResponse[[]ussd_merchant_dto.UssdMerchantResponse]{}, local_util.HandleDBError(err)
	}
	defer rows.Close()

	list := make([]ussd_merchant_dto.UssdMerchantResponse, 0)
	for rows.Next() {
		resp, err := scanResponse(rows)
		if err != nil {
			log.Errorf("[UssdMerchantRepo][FindAllWithPagination] scan failed: %v", err)
			return types.PaginatedResponse[[]ussd_merchant_dto.UssdMerchantResponse]{}, local_util.HandleDBError(err)
		}
		list = append(list, resp)
	}

	meta := local_util.BuildPaginationMeta(total, int(page), int(limit))
	return types.PaginatedResponse[[]ussd_merchant_dto.UssdMerchantResponse]{
		Data: list,
		Meta: meta,
	}, nil
}

// scanResponse scans a row into a UssdMerchantResponse DTO.
func scanResponse(row interface {
	Scan(dest ...interface{}) error
}) (ussd_merchant_dto.UssdMerchantResponse, error) {
	var resp ussd_merchant_dto.UssdMerchantResponse

	var enabled, isDeleted int
	var deletedAt sql.NullTime
	var logo, credential sql.NullString

	err := row.Scan(
		&resp.ID,
		&resp.MerchantCode,
		&credential,
		&resp.Name,
		&resp.SettlementMethod,
		&resp.PhoneNumber,
		&resp.Email,
		&resp.Service,
		&resp.AccountNumber,
		&logo,
		&enabled,
		&isDeleted,
		&resp.CreatedAt,
		&resp.UpdatedAt,
		&deletedAt,
	)
	if err != nil {
		return ussd_merchant_dto.UssdMerchantResponse{}, err
	}

	if logo.Valid {
		resp.Logo = logo.String
	}
	if credential.Valid {
		resp.Credential = credential.String
	}
	resp.Enabled = enabled == 1
	resp.IsDeleted = isDeleted == 1
	if deletedAt.Valid {
		resp.DeletedAt = deletedAt.Time
	}

	return resp, nil
}

// scanMerchant scans a row into an UssdMerchant model.
func scanMerchant(row interface {
	Scan(dest ...interface{}) error
}) (imodel.UssdMerchant, error) {
	var m imodel.UssdMerchant
	var enabled, isDeleted int
	var deletedAt sql.NullTime
	var logo, credential sql.NullString

	err := row.Scan(
		&m.ID,
		&m.MerchantCode,
		&credential,
		&m.Name,
		&m.SettlementMethod,
		&m.PhoneNumber,
		&m.Email,
		&m.Service,
		&m.AccountNumber,
		&logo,
		&enabled,
		&isDeleted,
		&m.CreatedAt,
		&m.UpdatedAt,
		&deletedAt,
	)
	if err != nil {
		return imodel.UssdMerchant{}, err
	}

	if logo.Valid {
		m.Logo = logo.String
	}
	if credential.Valid {
		m.Credential = credential.String
	}
	m.Enabled = enabled == 1
	m.IsDeleted = isDeleted == 1
	if deletedAt.Valid {
		m.DeletedAt = deletedAt.Time
	}

	return m, nil
}

func nullString(v string) sql.NullString {
	trimmed := strings.TrimSpace(v)
	if trimmed == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: trimmed, Valid: true}
}

func boolToOracleNumber(v bool) int {
	if v {
		return 1
	}
	return 0
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
