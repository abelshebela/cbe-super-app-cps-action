package event_merchant_oracle

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"cbe-super-app-cps-action/internal/constants/types"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)




type EventMerchantOracleRepository struct {
	OracleCliant   *sql.DB
	logger         utils.Logger
}

func NewEventMerchantOracleRepository(OracleCliant *sql.DB,logger utils.Logger) storage.EventMerchantRepository{
	return &EventMerchantOracleRepository{
		OracleCliant: OracleCliant,
		logger: logger,
	}
}


func (m *EventMerchantOracleRepository) Create(ctx context.Context,merchant model.EventMerchant) error {
	query :=  `INSERT INTO event_merchant (ID,MERCHANT_ID,MERCHANT_TYPE,SETTLEMENT_METHOD,METCHANT_NAME,BANK_ACCOUNT_NUMBER,EMAIL,PHONE_NUMBER,ENABLED,IS_DELETED,CREATED_AT,UPDATED_AT) VALUES (SYST_GEN(),:1,:2,:3,:4,:5,:6,:7,:8,:9,:10,:11,:12)`

	row,err := m.OracleCliant.ExecContext(ctx,query,merchant.ID,merchant.MerchantType,merchant.SettlementMethod,merchant.MerchantName,merchant.BankAccountNumber,merchant.Email,merchant.PhoneNumber,merchant.Enabled,merchant.IsDeleted,merchant.CreatedAt,merchant.UpdatedAt)

	if err != nil {
		m.logger.Errorf("[persistance oracle create ] got error while crearing event merchant error:",err)
		return  errors.New(localization.ErrorUnexpectedError.Code)
	}

	affected_row,err :=row.RowsAffected()

	if affected_row != 1 {
		return errors.New(localization.ErrorUnexpectedError.Code)

	}
	return nil
	
}

func (m *EventMerchantOracleRepository) Update (ctx context.Context,id string ,merchant model.EventMerchant) error {
	query,args :=  buildUpdateQuery(id,merchant)

	rows,err := m.OracleCliant.ExecContext(ctx,query,args...)
	if err != nil {
		m.logger.Errorf("[EVENT MERCHANT persistance UPDATE] got error while updated event merchant error : ",err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	row_affected ,_ := rows.RowsAffected()
	if row_affected != 1 {
		m.logger.Errorf("[event_merchant persistance update]no affected row found on update")
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}
func (m *EventMerchantOracleRepository) Delete(ctx context.Context,id string) error{
	query := `UPDATE  EVENT_MERCAHNT SET IS_DELETED = 1,DELETED_AT = SYSDATA WHERE ID = :1`

	rows,err := m.OracleCliant.ExecContext(ctx,query,id)
	if err != nil {
		m.logger.Errorf("[EVENT MERCHANT persistance Delet] got error while delete event merchant error : ",err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	row_affected ,_ := rows.RowsAffected()
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
        placeholders[i] = fmt.Sprintf(":%d", i+1)
        args = append(args, ids[i])
    }
    inClause := strings.Join(placeholders, ",")

    // 1. Check all IDs exist and are not deleted
    checkQuery := fmt.Sprintf("SELECT COUNT(1) FROM EVENT_MERCHANT WHERE ID IN (%s) AND IS_DELETED = 0", inClause)
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

    // 2. Update ENABLED for all valid IDs
    updateArgs := []interface{}{enable}
    updateArgs = append(updateArgs, args...)
    updateQuery := fmt.Sprintf("UPDATE EVENT_MERCHANT SET ENABLED = :1 WHERE ID IN (%s) AND IS_DELETED = 0", inClause)
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

func (m *EventMerchantOracleRepository) FindByID(ctx context.Context,id string)(*model.EventMerchant,error) {

query := `SELECT ID, MERCHANT_ID, MERCHANT_TYPE, SETTLEMENT_METHOD, MERCHANT_NAME, BANK_ACCOUNT_NUMBER, EMAIL, PHONE_NUMBER, ENABLED, IS_DELETED, CREATED_AT, UPDATED_AT, DELETED_AT FROM EVENT_MERCHANT WHERE ID = :1`

row := m.OracleCliant.QueryRowContext(ctx, query, id)

var idStr string // for Oracle ID column
MerchantData := model.EventMerchant{}

err := row.Scan(
    &idStr,
    &MerchantData.MerchantID,
    &MerchantData.MerchantType,
    &MerchantData.SettlementMethod,
    &MerchantData.MerchantName,
    &MerchantData.BankAccountNumber,
    &MerchantData.Email,
    &MerchantData.PhoneNumber,
    &MerchantData.Enabled,
    &MerchantData.IsDeleted,
    &MerchantData.CreatedAt,
    &MerchantData.UpdatedAt,
    &MerchantData.DeletedAt,
)


if err != nil {
	if errors.Is(err,sql.ErrNoRows){
		m.logger.Warnf("[event merchant persistance findbyID] No row find with given ID")
		return nil,errors.New(localization.UserNotFoundWithGivenID.Code)
		}
	}
	return &MerchantData,nil
}


func (m *EventMerchantOracleRepository) FindAllWithPagination(ctx context.Context,filterParam types.Filter)(*types.PaginatedResponse[[]model.EventMerchant],error){
	query ,args := buildQueryFromFilter(filterParam)
    TotalCountQuery := query
    CountArgs := args
    if filterParam.Page <= 0{
        filterParam.Page =1
    }
    if filterParam.PerPage <= 0 || filterParam.PerPage >=100{
        filterParam.PerPage = 10
    }
    // -------------------------
	// PAGINATION
	// -------------------------
	if filterParam.Page <= 0 {
		filterParam.Page = 1
	}

	if filterParam.PerPage <= 0 {
		filterParam.PerPage = 10
	}

	offset := (filterParam.Page - 1) * filterParam.PerPage

	query += fmt.Sprintf(`
		ORDER BY ID DESC
		OFFSET :%d ROWS
		FETCH NEXT :%d ROWS ONLY
	`)

	args = append(args, offset, filterParam.PerPage)
   var merchants []model.EventMerchant
    rows, err := m.OracleCliant.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

for rows.Next() {
    var merchant model.EventMerchant
    err := rows.Scan(
        &merchant.ID,
        &merchant.MerchantID,
        &merchant.MerchantType,
        &merchant.SettlementMethod,
        &merchant.MerchantName,
        &merchant.BankAccountNumber,
        &merchant.Email,
        &merchant.PhoneNumber,
        &merchant.Enabled,
        &merchant.IsDeleted,
        &merchant.CreatedAt,
        &merchant.UpdatedAt,
        &merchant.DeletedAt,
    )
    if err != nil {
        return nil, err
    }
    merchants = append(merchants, merchant)
    }
    var total int
    err = m.OracleCliant.QueryRowContext(ctx,TotalCountQuery,CountArgs...).Scan(&total)
    if err != nil {
        m.logger.Errorf("[event merchant persistance FindAllWithPagination] find error while counting doucment err:",err)
        return nil,errors.New(localization.ErrorUnexpectedError.Code)
    }
    meta := local_util.BuildPaginationMeta(int64(total),filterParam.Page,filterParam.PerPage)


    return &types.PaginatedResponse[[]model.EventMerchant]{
        Data: merchants,
        Meta: meta,
    },nil 

}


func buildQueryFromFilter(filter types.Filter) (string, []interface{}) {

	query := `SELECT * FROM EVENT_MERCANT WHERE 1=1`

	args := []interface{}{}
	argPos := 1

	allowedKeys := []string{
		"merchant_type",
		"merchant_id",
		"merchant_name",
		"email",
		"phone_number",
		"enabled",
		"bank_account_number",
	}

	// -------------------------
	// SEARCH (OR across fields)
	// -------------------------
	if filter.Search != "" {

		searchFields := []string{
			"MERCHANT_ID",
			"MERCHANT_NAME",
			"EMAIL",
			"PHONE_NUMBER",
			"BANK_ACCOUNT_NUMBER",
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

		if !contains(allowedKeys, key) {
			continue
		}

		query += fmt.Sprintf(`
			AND %s = :%d
		`, strings.ToUpper(key), argPos)

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
    return  nil,nil
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
    var isEnabled int
    var isDeleted int

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
        &merchant.CreatedAt,
        &merchant.UpdatedAt,
        &merchant.DeletedAt,
    )
    if err != nil {
        if errors.Is(err,sql.ErrNoRows){
            return nil,nil
        }
        return nil, err
    }

    merchant.Enabled = isEnabled == 1
    merchant.IsDeleted = isDeleted == 1

    return &merchant, nil
}









































func buildUpdateQuery(id string, merchant model.EventMerchant) (string, []interface{}) {
    query := "UPDATE EVENT_MERCHANT SET "
    args := []interface{}{}
    setParts := []string{}
    index := 1

    if merchant.MerchantID != "" {
        setParts = append(setParts, fmt.Sprintf("MERCHANT_ID = :%d", index))
        args = append(args, merchant.MerchantID)
        index++
    }
    if merchant.MerchantType != "" {
        setParts = append(setParts, fmt.Sprintf("MERCHANT_TYPE = :%d", index))
        args = append(args, merchant.MerchantType)
        index++
    }
    if merchant.SettlementMethod != "" {
        setParts = append(setParts, fmt.Sprintf("SETTLEMENT_METHOD = :%d", index))
        args = append(args, merchant.SettlementMethod)
        index++
    }
    if merchant.MerchantName != "" {
        setParts = append(setParts, fmt.Sprintf("MERCHANT_NAME = :%d", index))
        args = append(args, merchant.MerchantName)
        index++
    }
    if merchant.BankAccountNumber != "" {
        setParts = append(setParts, fmt.Sprintf("BANK_ACCOUNT_NUMBER = :%d", index))
        args = append(args, merchant.BankAccountNumber)
        index++
    }
    if merchant.Email != "" {
        setParts = append(setParts, fmt.Sprintf("EMAIL = :%d", index))
        args = append(args, merchant.Email)
        index++
    }
    if merchant.PhoneNumber != "" {
        setParts = append(setParts, fmt.Sprintf("PHONE_NUMBER = :%d", index))
        args = append(args, merchant.PhoneNumber)
        index++
    }
    // For booleans, you may want to always update, or only if changed. Here, always update:
    setParts = append(setParts, fmt.Sprintf("ENABLED = :%d", index))
    args = append(args, merchant.Enabled)
    index++

    setParts = append(setParts, fmt.Sprintf("IS_DELETED = :%d", index))
    args = append(args, merchant.IsDeleted)
    index++

    // UpdatedAt should always be set to now
    setParts = append(setParts, fmt.Sprintf("UPDATED_AT = :%d", index))
    args = append(args, merchant.UpdatedAt)
    index++

    // Optionally handle DeletedAt if needed
    if !merchant.DeletedAt.IsZero() {
        setParts = append(setParts, fmt.Sprintf("DELETED_AT = :%d", index))
        args = append(args, merchant.DeletedAt)
        index++
    }

    query += strings.Join(setParts, ", ")
    query += fmt.Sprintf(" WHERE ID = :%d", index)
    args = append(args, id)

    return query, args
}