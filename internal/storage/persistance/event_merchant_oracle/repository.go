package event_merchant_oracle

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"cbe-super-app-cps-action/internal/constants/types"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"cbe-super-app-cps-action/internal/constants/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type EventMerchantOracleRepository struct {
	OracleCliant *sql.DB
	logger       utils.Logger
}

func NewEventMerchantOracleRepository(OracleCliant *sql.DB, logger utils.Logger) storage.EventMerchantRepository {
	return &EventMerchantOracleRepository{
		OracleCliant: OracleCliant,
		logger:       logger,
	}
}
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func (m *EventMerchantOracleRepository) Create(ctx context.Context, merchant model.EventMerchant) error {
	// query := `INSERT INTO MERCHANTS (ID,MERCHANT_ID,MERCHANT_TYPE,SETTLEMENT_METHOD,METCHANT_NAME,BANK_ACCOUNT_NUMBER,EMAIL,PHONE_NUMBER,ENABLED,IS_DELETED,CREATED_AT,UPDATED_AT) VALUES (SYST_GEN(),:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12)`

	query := `INSERT INTO MERCHANTS (
    ID,
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
    LAST_MODIFIED_AT,
    DELETED_AT
) VALUES (
    SYS_GUID(), :1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11, :12
)`

	row, err := m.OracleCliant.ExecContext(
		ctx, query,
		merchant.BankAccountNumber,    // :1
		merchant.MerchantID,           // :2
		merchant.MerchantName,         // :3
		merchant.SettlementMethod,     // :4
		merchant.MerchantType,         // :5
		merchant.Email,                // :6
		merchant.PhoneNumber,          // :7
		boolToInt(merchant.Enabled),   // :8 (should be int: 1/0)
		boolToInt(merchant.IsDeleted), // :9 (should be int: 1/0)
		merchant.CreatedAt,            // :10
		merchant.UpdatedAt,            // :11
		merchant.DeletedAt,            // :12
	)
	if err != nil {
		m.logger.Errorf("[persistance oracle create ] got error while crearing event merchant error:", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	affected_row, err := row.RowsAffected()

	if affected_row != 1 {
		return errors.New(localization.ErrorUnexpectedError.Code)

	}
	return nil

}

func (m *EventMerchantOracleRepository) Update(ctx context.Context, id string, merchant model.EventMerchant) error {
	query, args := buildUpdateQuery(id, merchant)

	rows, err := m.OracleCliant.ExecContext(ctx, query, args...)
	if err != nil {
		m.logger.Errorf("[EVENT MERCHANT persistance UPDATE] got error while updated event merchant error : ", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	row_affected, _ := rows.RowsAffected()
	if row_affected != 1 {
		m.logger.Errorf("[event_merchant persistance update]no affected row found on update")
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}
func (m *EventMerchantOracleRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE  MERCHANTS SET IS_DELETED = 1,DELETED_AT = SYSDATE WHERE ID = HEXTORAW(:1) AND MERCHANT_TYPE = 'EVENT'`

	rows, err := m.OracleCliant.ExecContext(ctx, query, id)
	if err != nil {
		m.logger.Errorf("[EVENT MERCHANT persistance Delet] got error while delete event merchant error : ", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	row_affected, _ := rows.RowsAffected()
	if row_affected != 1 {
		m.logger.Errorf("[event_merchant persistance delete]no affected row found on update")
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil

}

func (m *EventMerchantOracleRepository) EnableOrDisable(ctx context.Context, ids []string, enable bool) error {
	if len(ids) == 0 {
		return nil // nothing to do
	}

	// Build placeholders for the IN clause
	placeholders := make([]string, len(ids))
	args := []interface{}{}
	for i := range ids {
		placeholders[i] = fmt.Sprintf("HEXTORAW(:%d)", i+2) // +2 because :1 is for IS_ENABLED
		args = append(args, ids[i])
	}
	inClause := strings.Join(placeholders, ",")

	// 1. Check all IDs exist and are not deleted
	checkQuery := fmt.Sprintf("SELECT COUNT(1) FROM MERCHANTS WHERE ID IN (%s) AND IS_DELETED = 0 AND MERCHANT_TYPE = 'EVENT'", inClause)
	var count int
	err := m.OracleCliant.QueryRowContext(ctx, checkQuery, args...).Scan(&count)
	if err != nil {
		m.logger.Errorf("[EVENT MERCHANT EnableOrDisable] error checking IDs: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	if count != len(ids) {
		m.logger.Warnf("[EVENT MERCHANT EnableOrDisable] not all IDs exist or are not deleted, requested: %d, found: %d", len(ids), count)
		return errors.New("one or more IDs do not exist or are deleted")
	}

	// 2. Update IS_ENABLED for all valid IDs
	isEnabled := 0
	if enable {
		isEnabled = 1
	}
	updateArgs := []interface{}{isEnabled}
	updateArgs = append(updateArgs, args...)
	updateQuery := fmt.Sprintf("UPDATE MERCHANTS SET IS_ENABLED = :1 WHERE ID IN (%s) AND IS_DELETED = 0 AND MERCHANT_TYPE = 'EVENT'", inClause)
	result, err := m.OracleCliant.ExecContext(ctx, updateQuery, updateArgs...)
	if err != nil {
		m.logger.Errorf("[EVENT MERCHANT EnableOrDisable] error updating: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	rowsAffected, _ := result.RowsAffected()
	if int(rowsAffected) != len(ids) {
		m.logger.Warnf("[EVENT MERCHANT EnableOrDisable] not all IDs updated, requested: %d, updated: %d", len(ids), rowsAffected)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

//   //MINE
// func (m *EventMerchantOracleRepository) EnableOrDisable(ctx context.Context, ids []string, enable bool) error {
// 	if len(ids) == 0 {
// 		return nil // nothing to do
// 	}

// 	// Build placeholders for the IN clause
// 	placeholders := make([]string, len(ids))
// 	args := []interface{}{}
// 	for i := range ids {
// 		placeholders[i] = fmt.Sprintf(":%d", i+1)
// 		args = append(args, ids[i])
// 	}
// 	inClause := strings.Join(placeholders, ",")

// 	// 1. Check all IDs exist and are not deleted
// 	checkQuery := fmt.Sprintf("SELECT COUNT(1) FROM MERCHANTS WHERE ID IN (%s) AND IS_DELETED = 0 AND MERCHANT_TYPE = 'EVENT'", inClause)
// 	var count int
// 	err := m.OracleCliant.QueryRowContext(ctx, checkQuery, args...).Scan(&count)
// 	if err != nil {
// 		m.logger.Errorf("[EVENT MERCHANT EnableOrDisable] error checking IDs: %v", err)
// 		return errors.New(localization.ErrorUnexpectedError.Code)
// 	}
// 	if count != len(ids) {
// 		m.logger.Warnf("[EVENT MERCHANT EnableOrDisable] not all IDs exist or are not deleted, requested: %d, found: %d", len(ids), count)
// 		return errors.New("one or more IDs do not exist or are deleted")
// 	}

// 	// 2. Update ENABLED for all valid IDs
// 	updateArgs := []interface{}{enable}
// 	updateArgs = append(updateArgs, args...)
// 	updateQuery := fmt.Sprintf("UPDATE MERCHANTS SET ENABLED = :1 WHERE ID IN (%s) AND IS_DELETED = 0 AND MERCHANT_TYPE = 'EVENT'", inClause)
// 	result, err := m.OracleCliant.ExecContext(ctx, updateQuery, updateArgs...)
// 	if err != nil {
// 		m.logger.Errorf("[EVENT MERCHANT EnableOrDisable] error updating: %v", err)
// 		return errors.New(localization.ErrorUnexpectedError.Code)
// 	}
// 	rowsAffected, _ := result.RowsAffected()
// 	if int(rowsAffected) != len(ids) {
// 		m.logger.Warnf("[EVENT MERCHANT EnableOrDisable] not all IDs updated, requested: %d, updated: %d", len(ids), rowsAffected)
// 		return errors.New(localization.ErrorUnexpectedError.Code)
// 	}

// 	return nil
// }

func (m *EventMerchantOracleRepository) FindByID(ctx context.Context, id string) (*model.EventMerchant, error) {
	query := `SELECT 
        ID, 
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
        LAST_MODIFIED_AT, 
        DELETED_AT 
    FROM MERCHANTS 
    WHERE ID = HEXTORAW(:1) AND IS_DELETED = 0 AND MERCHANT_TYPE = 'EVENT'`

	row := m.OracleCliant.QueryRowContext(ctx, query, id)

	MerchantData := model.EventMerchant{}
	var (
		isEnabled, isDeleted            int
		byteData                        []byte
		createdAt, updatedAt, deletedAt sql.NullTime
	)

	err := row.Scan(
		&byteData,
		&MerchantData.BankAccountNumber,
		&MerchantData.MerchantID,
		&MerchantData.MerchantName,
		&MerchantData.SettlementMethod,
		&MerchantData.MerchantType,
		&MerchantData.Email,
		&MerchantData.PhoneNumber,
		&isEnabled,
		&isDeleted,
		&createdAt,
		&updatedAt,
		&deletedAt,
	)

	MerchantData.Enabled = isEnabled == 1
	MerchantData.IsDeleted = isDeleted == 1
	if createdAt.Valid {
		MerchantData.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		MerchantData.UpdatedAt = updatedAt.Time
	}
	if deletedAt.Valid {
		MerchantData.DeletedAt = deletedAt.Time
	}
	MerchantData.ID = hex.EncodeToString(byteData)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			m.logger.Warnf("[event merchant persistance findbyID] No row find with given ID")
			return nil, errors.New(localization.UserNotFoundWithGivenID.Code)
		}
		return nil, err
	}
	return &MerchantData, nil
}

func (m *EventMerchantOracleRepository) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.EventMerchant], error) {
	// Set default pagination values
	if filterParam.Page <= 0 {
		filterParam.Page = 1
	}
	if filterParam.PerPage <= 0 || filterParam.PerPage > 100 {
		filterParam.PerPage = 10
	}

	// Build base query and args
	baseQuery, args := buildQueryFromFilter(filterParam)

	// Build paginated query
	paginatedQuery := baseQuery + `
        ORDER BY ID DESC
        OFFSET :%d ROWS FETCH NEXT :%d ROWS ONLY
    `
	// Add offset and limit as parameters
	offset := (filterParam.Page - 1) * filterParam.PerPage
	paginatedQuery = fmt.Sprintf(paginatedQuery, len(args)+1, len(args)+2)
	args = append(args, offset, filterParam.PerPage)

	// Query for data
	rows, err := m.OracleCliant.QueryContext(ctx, paginatedQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var merchants []model.EventMerchant
	for rows.Next() {
		var merchant model.EventMerchant

		var (
			idBytes                         []byte
			createdAt, updatedAt, deletedAt sql.NullTime
		)
		err := rows.Scan(
			&idBytes,                    // ID (RAW(16))
			&merchant.BankAccountNumber, // MERCHANT_ACCOUNT_NUMBER
			&merchant.MerchantID,        // MERCHANT_CODE
			&merchant.MerchantName,      // MERCHANT_NAME
			&merchant.SettlementMethod,  // SETTLEMENT_METHOD
			&merchant.MerchantType,      // MERCHANT_TYPE
			&merchant.Email,             // CONTACT_EMAIL
			&merchant.PhoneNumber,       // CONTACT_PHONE
			&merchant.Enabled,           // IS_ENABLED
			&merchant.IsDeleted,         // IS_DELETED
			&createdAt,                  // CREATED_AT
			&updatedAt,                  // LAST_MODIFIED_AT
			&deletedAt,                  // DELETED_AT
		)
		if createdAt.Valid {
			merchant.CreatedAt = createdAt.Time
		}
		if updatedAt.Valid {
			merchant.UpdatedAt = updatedAt.Time
		}
		if deletedAt.Valid {
			merchant.DeletedAt = deletedAt.Time
		}
		merchant.ID = hex.EncodeToString(idBytes)

		if err != nil {
			return nil, err
		}
		merchants = append(merchants, merchant)
	}

	// Build count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM (%s)", baseQuery)
	var total int
	err = m.OracleCliant.QueryRowContext(ctx, countQuery, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		m.logger.Errorf("[event merchant persistance FindAllWithPagination] error counting documents: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(int64(total), filterParam.Page, filterParam.PerPage)
	return &types.PaginatedResponse[[]model.EventMerchant]{
		Data: merchants,
		Meta: meta,
	}, nil
}

func buildQueryFromFilter(filter types.Filter) (string, []interface{}) {

	// Use explicit columns for compatibility with model.EventMerchant
	query := `SELECT 
    ID, 
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
    LAST_MODIFIED_AT, 
    DELETED_AT 
FROM MERCHANTS 
WHERE 1=1 AND MERCHANT_TYPE = 'EVENT' AND IS_DELETED = 0`

	args := []interface{}{}
	argPos := 1

	allowedKeys := []string{
		"MERCHANT_CODE",
		"MERCHANT_NAME",
		"CONTACT_EMAIL",
		"CONTACT_PHONE",
		"IS_ENABLED",
		"MERCHANT_ACCOUNT_NUMBER",
	}

	goStructToOracleFileds := map[string]string{
		"merchant_id":"MERCHANT_CODE",
		"merchant_name":"MERCHANT_NAME",
		"email":"CONTACT_EMAIL",
		"phone_number":"CONTACT_PHONE",
		"enabled":"IS_ENABLED",
		"bank_account_number":"MERCHANT_ACCOUNT_NUMBER",

	}

	// -------------------------
	// SEARCH (OR across fields)
	// -------------------------
	if filter.Search != "" {

		searchFields := []string{
			"MERCHANT_CODE",
			"MERCHANT_NAME",
			"CONTACT_EMAIL",
			"CONTACT_PHONE",
			"MERCHANT_ACCOUNT_NUMBER",
		}

		searchParts := []string{}
		searchValue := "%" + filter.Search + "%"

		for _, field := range searchFields {

			searchParts = append(
				searchParts,
				fmt.Sprintf("LOWER(%s) LIKE LOWER(:%d)", field, argPos),
			)

			args = append(args, searchValue)
			argPos++
		}

		query += fmt.Sprintf(`
			AND (%s)
		`, strings.Join(searchParts, " OR "))
	}

	// -------------------------
	// FILTERS (AND)
	// -------------------------
	for key, value := range filter.Filters {

		if !contains(allowedKeys, strings.ToUpper(goStructToOracleFileds[key])) {
			continue
		}

		query += fmt.Sprintf(`
			AND %s = :%d
		`, strings.ToUpper(goStructToOracleFileds[key]), argPos)

		args = append(args, value)
		argPos++
	}

	return query, args
}

func contains(arr []string, key string) bool {
	for _, v := range arr {
		if v == key {
			return true
		}
	}
	return false
}

func (m *EventMerchantOracleRepository) FindOne(ctx context.Context, query bson.M) (*model.EventMerchant, error) {

	// ExistedMerchant := model.EventMerchant{}
	// err := m.OracleCliant.QueryRowContext(ctx, query).Scan(
	// 	&ExistedMerchant.ID,
	// 	&ExistedMerchant.MerchantID,
	// 	&ExistedMerchant.MerchantType,
	// 	&ExistedMerchant.SettlementMethod,
	// 	&ExistedMerchant.MerchantName,
	// 	&ExistedMerchant.BankAccountNumber,
	// 	&ExistedMerchant.Email,
	// 	&ExistedMerchant.PhoneNumber,
	// 	&ExistedMerchant.Enabled,
	// 	&ExistedMerchant.IsDeleted,
	// 	&ExistedMerchant.CreatedAt,
	// 	&ExistedMerchant.UpdatedAt,
	// 	&ExistedMerchant.DeletedAt,
	// )
	// if err != nil {
	// 	return nil, err
	// }
	// return &ExistedMerchant, nil

	//unimplimented
	return nil, nil
}

func (m *EventMerchantOracleRepository) FindOneO(ctx context.Context, data *types.CheckMerchant) (*model.EventMerchant, error) {
	query := `
        SELECT
            ID,
            MERCHANT_CODE,
            MERCHANT_TYPE,
            SETTLEMENT_METHOD,
            MERCHANT_NAME,
            MERCHANT_ACCOUNT_NUMBER,
            CONTACT_EMAIL,
            CONTACT_PHONE,
            IS_ENABLED,
            IS_DELETED,
            CREATED_AT,
            LAST_MODIFIED_AT,
            DELETED_AT
        FROM MERCHANTS
        WHERE MERCHANT_TYPE = 'EVENT'
            AND (
                MERCHANT_ACCOUNT_NUMBER = :1 OR
                CONTACT_EMAIL = :2 OR
                CONTACT_PHONE = :3 OR
                MERCHANT_CODE = :4
            )
            AND IS_DELETED = 0
    `

	merchant := model.EventMerchant{}
	var (
		isEnabled int
		isDeleted int

		createdAt, updatedAt, deletedAt sql.NullTime
	)

	err := m.OracleCliant.QueryRowContext(
		ctx,
		query,
		data.BankAccountNumber,
		data.Email,
		data.PhoneNumber,
		data.MerchantCode,
	).Scan(
		&merchant.ID,
		&merchant.MerchantID,
		&merchant.MerchantType,
		&merchant.SettlementMethod,
		&merchant.MerchantName,
		&merchant.BankAccountNumber,
		&merchant.Email,
		&merchant.PhoneNumber,
		&isEnabled,
		&isDeleted,
		&createdAt,
		&updatedAt,
		&deletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	merchant.Enabled = isEnabled == 1
	merchant.IsDeleted = isDeleted == 1

	return &merchant, nil
}

func buildUpdateQuery(id string, merchant model.EventMerchant) (string, []interface{}) {
	query := "UPDATE MERCHANTS SET "
	args := []interface{}{}
	setParts := []string{}
	index := 1

	if merchant.BankAccountNumber != "" {
		setParts = append(setParts, fmt.Sprintf("MERCHANT_ACCOUNT_NUMBER = :%d", index))
		args = append(args, merchant.BankAccountNumber)
		index++
	}
	if merchant.MerchantID != "" {
		setParts = append(setParts, fmt.Sprintf("MERCHANT_CODE = :%d", index))
		args = append(args, merchant.MerchantID)
		index++
	}
	if merchant.MerchantName != "" {
		setParts = append(setParts, fmt.Sprintf("MERCHANT_NAME = :%d", index))
		args = append(args, merchant.MerchantName)
		index++
	}
	if merchant.SettlementMethod != "" {
		setParts = append(setParts, fmt.Sprintf("SETTLEMENT_METHOD = :%d", index))
		args = append(args, merchant.SettlementMethod)
		index++
	}
	if merchant.MerchantType != "" {
		setParts = append(setParts, fmt.Sprintf("MERCHANT_TYPE = :%d", index))
		args = append(args, merchant.MerchantType)
		index++
	}
	if merchant.Email != "" {
		setParts = append(setParts, fmt.Sprintf("CONTACT_EMAIL = :%d", index))
		args = append(args, merchant.Email)
		index++
	}
	if merchant.PhoneNumber != "" {
		setParts = append(setParts, fmt.Sprintf("CONTACT_PHONE = :%d", index))
		args = append(args, merchant.PhoneNumber)
		index++
	}
	// // Always update IS_ENABLED and IS_DELETED
	// setParts = append(setParts, fmt.Sprintf("IS_ENABLED = :%d", index))
	// args = append(args, merchant.Enabled)
	// index++

	// setParts = append(setParts, fmt.Sprintf("IS_DELETED = :%d", index))
	// args = append(args, merchant.IsDeleted)
	// index++

	// Always update LAST_MODIFIED_AT
	setParts = append(setParts, fmt.Sprintf("LAST_MODIFIED_AT = :%d", index))
	args = append(args, merchant.UpdatedAt)
	index++

	// Optionally handle DELETED_AT if needed
	// if !merchant.DeletedAt.IsZero() {
	// 	setParts = append(setParts, fmt.Sprintf("DELETED_AT = :%d", index))
	// 	args = append(args, merchant.DeletedAt)
	// 	index++
	// }

	query += strings.Join(setParts, ", ")
	query += fmt.Sprintf(" WHERE MERCHANT_TYPE = 'EVENT' AND ID = :%d", index)
	args = append(args, id)

	return query, args
}

// func buildUpdateQuery(id string, merchant model.EventMerchant) (string, []interface{}) {
// 	query := "UPDATE MERCHANTS SET "
// 	args := []interface{}{}
// 	setParts := []string{}
// 	index := 1

// 	if merchant.MerchantID != "" {
// 		setParts = append(setParts, fmt.Sprintf("MERCHANT_ID = :%d", index))
// 		args = append(args, merchant.MerchantID)
// 		index++
// 	}
// 	if merchant.MerchantType != "" {
// 		setParts = append(setParts, fmt.Sprintf("MERCHANT_TYPE = :%d", index))
// 		args = append(args, merchant.MerchantType)
// 		index++
// 	}
// 	if merchant.SettlementMethod != "" {
// 		setParts = append(setParts, fmt.Sprintf("SETTLEMENT_METHOD = :%d", index))
// 		args = append(args, merchant.SettlementMethod)
// 		index++
// 	}
// 	if merchant.MerchantName != "" {
// 		setParts = append(setParts, fmt.Sprintf("MERCHANT_NAME = :%d", index))
// 		args = append(args, merchant.MerchantName)
// 		index++
// 	}
// 	if merchant.BankAccountNumber != "" {
// 		setParts = append(setParts, fmt.Sprintf("BANK_ACCOUNT_NUMBER = :%d", index))
// 		args = append(args, merchant.BankAccountNumber)
// 		index++
// 	}
// 	if merchant.Email != "" {
// 		setParts = append(setParts, fmt.Sprintf("EMAIL = :%d", index))
// 		args = append(args, merchant.Email)
// 		index++
// 	}
// 	if merchant.PhoneNumber != "" {
// 		setParts = append(setParts, fmt.Sprintf("PHONE_NUMBER = :%d", index))
// 		args = append(args, merchant.PhoneNumber)
// 		index++
// 	}
// 	// For booleans, you may want to always update, or only if changed. Here, always update:
// 	setParts = append(setParts, fmt.Sprintf("ENABLED = :%d", index))
// 	args = append(args, merchant.Enabled)
// 	index++

// 	setParts = append(setParts, fmt.Sprintf("IS_DELETED = :%d", index))
// 	args = append(args, merchant.IsDeleted)
// 	index++

// 	// UpdatedAt should always be set to now
// 	setParts = append(setParts, fmt.Sprintf("UPDATED_AT = :%d", index))
// 	args = append(args, merchant.UpdatedAt)
// 	index++

// 	// Optionally handle DeletedAt if needed
// 	if !merchant.DeletedAt.IsZero() {
// 		setParts = append(setParts, fmt.Sprintf("DELETED_AT = :%d", index))
// 		args = append(args, merchant.DeletedAt)
// 		index++
// 	}

// 	query += strings.Join(setParts, ", ")
// 	query += fmt.Sprintf(" WHERE MERCHANT_TYPE = 'EVENT' AND ID = :%d", index)
// 	args = append(args, id)

// 	return query, args
// }
