package sqlc

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"database/sql"
	"fmt"
	"strings"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func NewBankRepository(db *sql.DB, log utils.Logger) storage.BankOracleRepository {
	return &Queries{
		db:     db,
		logger: log,
	}
}
func (q *Queries) Create(ctx context.Context, bank *imodel.BankOracle) error {
	query := `INSERT INTO banks (id, bank_name, logo, bic_code, is_enabled, account_length, has_alpha_numeric, create_at, update_at)
		VALUES (SYS_GUID(), :1, :2, :3, :4, :5, :6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`
	_, err := q.db.ExecContext(ctx, query,
		bank.BankName,
		bank.Logo,
		bank.BICCode,
		bank.IsEnabled,
		bank.AccountLength,
		bank.HasAlphaNumeric,
	)
	return err
}

func (q *Queries) Update(ctx context.Context, id string, bank *imodel.BankOracle) error {
	query := `UPDATE banks SET bank_name = :1, logo = :2, bic_code = :3, is_enabled = :4, account_length = :5, has_alpha_numeric = :6, update_at = CURRENT_TIMESTAMP WHERE id = :7`
	_, err := q.db.ExecContext(ctx, query,
		bank.BankName,
		bank.Logo,
		bank.BICCode,
		bank.IsEnabled,
		bank.AccountLength,
		bank.HasAlphaNumeric,
		id,
	)
	return err
}

func (q *Queries) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM BANKS WHERE id = :1`
	_, err := q.db.ExecContext(ctx, query, id)
	return err
}

func (q *Queries) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	query := `UPDATE banks SET is_enabled = :1, update_at = CURRENT_TIMESTAMP WHERE id = :2`
	var err error
	if enable {
		_, err = q.db.ExecContext(ctx, query, 1, id)
	} else {
		_, err = q.db.ExecContext(ctx, query, 0, id)

	}
	return err
}

func (q *Queries) FindByID(ctx context.Context, id string) (*imodel.BankOracle, error) {
	query := `SELECT RAWTOHEX(ID) AS id, bank_name, logo, bic_code, is_enabled, account_length, has_alpha_numeric, create_at, update_at FROM BANKS WHERE id = :1`
	row := q.db.QueryRowContext(ctx, query, id)
	var bank imodel.BankOracle
	err := row.Scan(
		&bank.ID,
		&bank.BankName,
		&bank.Logo,
		&bank.BICCode,
		&bank.IsEnabled,
		&bank.AccountLength,
		&bank.HasAlphaNumeric,
		&bank.CreateAt,
		&bank.UpdateAt,
	)
	if err == sql.ErrNoRows {
		return nil, localization.ErrorResourceNotFound
	}
	if err != nil {
		return nil, err
	}
	return &bank, nil
}

func (q *Queries) FindByBIC(ctx context.Context, bic string) (*imodel.BankOracle, error) {
	query := `SELECT RAWTOHEX(ID) AS id, bank_name, logo, bic_code, is_enabled, account_length, has_alpha_numeric, create_at, update_at FROM BANKS WHERE bic_code = :1 AND is_enabled = 1`
	row := q.db.QueryRowContext(ctx, query, bic)
	var bank imodel.BankOracle
	err := row.Scan(
		&bank.ID,
		&bank.BankName,
		&bank.Logo,
		&bank.BICCode,
		&bank.IsEnabled,
		&bank.AccountLength,
		&bank.HasAlphaNumeric,
		&bank.CreateAt,
		&bank.UpdateAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &bank, nil
}

func (q *Queries) FindByNameOrBIC(ctx context.Context, bic, name string) (*imodel.BankOracle, error) {
	q.logger.Infof("[BankOracleRepository][FindByNameOrBIC] called with name: %s, bic: %s", name, bic)
	var conditions []string
	var args []interface{}
	idx := 1
	if name != "" {
		conditions = append(conditions, fmt.Sprintf("UPPER(bank_name) = :%d", idx))
		args = append(args, strings.ToUpper(name))
		idx++
	}
	if bic != "" {
		conditions = append(conditions, fmt.Sprintf("UPPER(bic_code) = :%d", idx))
		args = append(args, strings.ToUpper(bic))
		idx++
	}
	if len(conditions) == 0 {
		return nil, nil
	}
	query := fmt.Sprintf(`SELECT RAWTOHEX(ID) AS id, bank_name, logo, bic_code, is_enabled, account_length, has_alpha_numeric, create_at, update_at FROM BANKS WHERE %s`, strings.Join(conditions, " OR "))
	row := q.db.QueryRowContext(ctx, query, args...)
	var bank imodel.BankOracle
	err := row.Scan(
		&bank.ID,
		&bank.BankName,
		&bank.Logo,
		&bank.BICCode,
		&bank.IsEnabled,
		&bank.AccountLength,
		&bank.HasAlphaNumeric,
		&bank.CreateAt,
		&bank.UpdateAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &bank, nil
}

// func (q *Queries) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.BankOracle], error) {
// 	q.logger.Infof("[BankOracleRepository][FindAllWithPagination] called with filter: %+v", filterParam)
// 	// Build filtering logic
// 	var filters []string
// 	var args []interface{}
// 	idx := 1

// 	filters = append(filters, "1=1") // always true, simplifies appending ANDs

// 	if filterParam.Search != "" && filterParam.Search != "enabled" {
// 		search := "%" + strings.ToUpper(filterParam.Search) + "%"
// 		filters = append(filters, fmt.Sprintf("(UPPER(bank_name) LIKE :%d OR UPPER(bic_code) LIKE :%d)", idx, idx+1))
// 		args = append(args, search, search)
// 		idx += 2
// 	}
// 	if filterParam.Search == "enabled" {
// 		filters = append(filters, "is_enabled = 1")
// 	}

// 	oracleQuery := lib.BuildOracleFilter(filterParam, map[string]string{"bank_name": "UPPER(bank_name)", "bic_code": "UPPER(bic_code)"}, []string{"bank_name", "bic_code", "is_enabled"})

// 	whereClause := oracleQuery.WhereClause
// 	args = append(args, oracleQuery.Args)
// 	offset := oracleQuery.Offset
// 	limit := oracleQuery.Limit

// 	// whereClause := strings.Join(filters, " AND ")
// 	q.logger.Debugf("[BankOracleRepository][FindAllWithPagination] whereClause: %s, args: %+v", whereClause, args)

// 	// Count total
// 	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM BANKS WHERE %s", whereClause)
// 	q.logger.Debugf("[BankOracleRepository][FindAllWithPagination] countQuery: %s", countQuery)
// 	var total int64

// 	err := q.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
// 	if err != nil {
// 		q.logger.Errorf("[BankOracleRepository][FindAllWithPagination] failed to count: %v", err)
// 		return nil, err
// 	}
// 	q.logger.Infof("[BankOracleRepository][FindAllWithPagination] total records: %d", total)

// 	if total == 0 {
// 		meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
// 		resp := &types.PaginatedResponse[[]imodel.BankOracle]{
// 			Data: []imodel.BankOracle{},
// 			Meta: meta,
// 		}
// 		q.logger.Infof("[BankOracleRepository][FindAllWithPagination] no records found")
// 		return resp, nil
// 	}

// 	// offset := (filterParam.Page - 1) * filterParam.PerPage
// 	// limit := filterParam.PerPage
// 	// If requested offset is beyond total, return all data (no pagination)
// 	if int64(offset) >= total {
// 		offset = 0
// 		limit = total
// 	} else if int64(offset)+int64(limit) > total {
// 		limit = total - offset
// 	}
// 	q.logger.Debugf("[BankOracleRepository][FindAllWithPagination] offset: %d, limit: %d", offset, limit)

// 	// Fetch paginated results
// 	selectQuery := fmt.Sprintf(`SELECT RAWTOHEX(ID) AS id, bank_name, logo, bic_code, is_enabled, account_length, has_alpha_numeric, create_at, update_at FROM BANKS WHERE %s ORDER BY create_at DESC OFFSET :%d ROWS FETCH NEXT :%d ROWS ONLY`, whereClause, idx, idx+1)
// 	args = append(args, offset, limit)
// 	q.logger.Debugf("[BankOracleRepository][FindAllWithPagination] selectQuery: %s", selectQuery)
// 	rows, err := q.db.QueryContext(ctx, selectQuery, args...)
// 	if err != nil {
// 		q.logger.Errorf("[BankOracleRepository][FindAllWithPagination] failed to fetch rows: %v", err)
// 		return nil, err
// 	}
// 	defer rows.Close()
// 	banks := []imodel.BankOracle{}
// 	q.logger.Infof("[BankOracleRepository][FindAllWithPagination] scanning rows...")
// 	for rows.Next() {
// 		var bank imodel.BankOracle
// 		err := rows.Scan(
// 			&bank.ID,
// 			&bank.BankName,
// 			&bank.Logo,
// 			&bank.BICCode,
// 			&bank.IsEnabled,
// 			&bank.AccountLength,
// 			&bank.HasAlphaNumeric,
// 			&bank.CreateAt,
// 			&bank.UpdateAt,
// 		)
// 		if err != nil {
// 			q.logger.Errorf("[BankOracleRepository][FindAllWithPagination] failed to scan row: %v", err)
// 			return nil, err
// 		}
// 		q.logger.Debugf("[BankOracleRepository][FindAllWithPagination] scanned bank: %+v", bank)
// 		banks = append(banks, bank)
// 	}
// 	if err := rows.Err(); err != nil {
// 		q.logger.Errorf("[BankOracleRepository][FindAllWithPagination] rows iteration failed: %v", err)
// 		return nil, err
// 	}

// 	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
// 	q.logger.Debugf("[BankOracleRepository][FindAllWithPagination] meta: %+v", meta)
// 	resp := &types.PaginatedResponse[[]imodel.BankOracle]{
// 		Data: banks,
// 		Meta: meta,
// 	}
// 	q.logger.Infof("[BankOracleRepository][FindAllWithPagination] success, returning response")
// 	return resp, nil
// }

func (q *Queries) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.BankOracle], error) {
	q.logger.Infof("[BankOracleRepository][FindAllWithPagination] called with filter: %+v", filterParam)
	// Build filtering logic
	var filters []string
	var args []interface{}
	idx := 1

	filters = append(filters, "1=1") // always true, simplifies appending ANDs

	if filterParam.Search != "" && filterParam.Search != "enabled" {
		search := "%" + strings.ToUpper(filterParam.Search) + "%"
		filters = append(filters, fmt.Sprintf("(UPPER(bank_name) LIKE :%d OR UPPER(bic_code) LIKE :%d)", idx, idx+1))
		args = append(args, search, search)
		idx += 2
	}
	if val, ok := filterParam.Filters["enabled"]; ok {
		if b, ok := val.(bool); ok {
			if b {
				filters = append(filters, "is_enabled = 1")
			} else {
				filters = append(filters, "is_enabled = 0")
			}
		}
	}

	whereClause := strings.Join(filters, " AND ")
	q.logger.Debugf("[BankOracleRepository][FindAllWithPagination] whereClause: %s, args: %+v", whereClause, args)

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM BANKS WHERE %s", whereClause)
	q.logger.Debugf("[BankOracleRepository][FindAllWithPagination] countQuery: %s", countQuery)
	var total int64

	err := q.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		q.logger.Errorf("[BankOracleRepository][FindAllWithPagination] failed to count: %v", err)
		return nil, err
	}
	q.logger.Infof("[BankOracleRepository][FindAllWithPagination] total records: %d", total)

	if total == 0 {
		meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
		resp := &types.PaginatedResponse[[]imodel.BankOracle]{
			Data: []imodel.BankOracle{},
			Meta: meta,
		}
		q.logger.Infof("[BankOracleRepository][FindAllWithPagination] no records found")
		return resp, nil
	}

	offset := (filterParam.Page - 1) * filterParam.PerPage
	limit := filterParam.PerPage
	// If requested offset is beyond total, return all data (no pagination)
	if int64(offset) >= total {
		offset = 0
		limit = int(total)
	} else if int64(offset)+int64(limit) > total {
		limit = int(total) - offset
	}
	q.logger.Debugf("[BankOracleRepository][FindAllWithPagination] offset: %d, limit: %d", offset, limit)

	// Determine sort order for name
	orderBy := "create_at DESC"
	if val, ok := filterParam.Filters["sort_name"]; ok {
		q.logger.Debugf("[BankOracleRepository][FindAllWithPagination] sort filter: %v", val)
		if s, ok := val.(string); ok && (strings.ToUpper(s) == "ASC" || strings.ToUpper(s) == "DESC") {
			orderBy = fmt.Sprintf("UPPER(bank_name) %s", strings.ToUpper(s))
		}
	}

	// Fetch paginated results
	selectQuery := fmt.Sprintf(`SELECT RAWTOHEX(ID) AS id, bank_name, logo, bic_code, is_enabled, account_length, has_alpha_numeric, create_at, update_at FROM BANKS WHERE %s ORDER BY %s OFFSET :%d ROWS FETCH NEXT :%d ROWS ONLY`, whereClause, orderBy, idx, idx+1)
	args = append(args, offset, limit)
	q.logger.Debugf("[BankOracleRepository][FindAllWithPagination] selectQuery: %s", selectQuery)
	rows, err := q.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		q.logger.Errorf("[BankOracleRepository][FindAllWithPagination] failed to fetch rows: %v", err)
		return nil, err
	}
	defer rows.Close()
	banks := []imodel.BankOracle{}
	q.logger.Infof("[BankOracleRepository][FindAllWithPagination] scanning rows...")
	for rows.Next() {
		var bank imodel.BankOracle
		err := rows.Scan(
			&bank.ID,
			&bank.BankName,
			&bank.Logo,
			&bank.BICCode,
			&bank.IsEnabled,
			&bank.AccountLength,
			&bank.HasAlphaNumeric,
			&bank.CreateAt,
			&bank.UpdateAt,
		)
		if err != nil {
			q.logger.Errorf("[BankOracleRepository][FindAllWithPagination] failed to scan row: %v", err)
			return nil, err
		}
		q.logger.Debugf("[BankOracleRepository][FindAllWithPagination] scanned bank: %+v", bank)
		banks = append(banks, bank)
	}
	if err := rows.Err(); err != nil {
		q.logger.Errorf("[BankOracleRepository][FindAllWithPagination] rows iteration failed: %v", err)
		return nil, err
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	q.logger.Debugf("[BankOracleRepository][FindAllWithPagination] meta: %+v", meta)
	resp := &types.PaginatedResponse[[]imodel.BankOracle]{
		Data: banks,
		Meta: meta,
	}
	q.logger.Infof("[BankOracleRepository][FindAllWithPagination] success, returning response")
	return resp, nil
}
