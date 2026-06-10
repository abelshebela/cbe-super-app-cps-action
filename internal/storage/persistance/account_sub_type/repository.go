package account_sub_type_oracle

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type Repository struct {
	db     *sql.DB
	logger utils.Logger
}

var _ storage.AccountSubTypeOracleRepository = (*Repository)(nil)

func NewAccountSubTypeOracleRepository(db *sql.DB, logger utils.Logger) storage.AccountSubTypeOracleRepository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

func (r *Repository) Create(ctx context.Context, ast *imodel.AccountSubType) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	q := `INSERT INTO ACCOUNT_SUB_TYPES (ID, ACCOUNT_TYPE, GENDER, ACCOUNT_SUB_TYPE_NAME, ACCOUNT_SUB_TYPE_CODE, IS_ENABLED, IS_DELETED, CREATED_AT, LAST_MODIFIED_AT)
		VALUES (SYS_GUID(), :1, :2, :3, :4, :5, :6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`
	_, err := r.db.ExecContext(ctx, q,
		ast.AccountType,
		ast.Gender,
		ast.AccountSubTypeName,
		ast.AccountSubTypeCode,
		ast.IsEnabled,
		0,
	)
	if err != nil {
		log.Errorf("[AccountSubTypeOracle][Create] insert failed: %v", err)
		return err
	}
	return nil
}

func (r *Repository) Update(ctx context.Context, id string, ast *imodel.AccountSubType) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	setClauses := []string{}
	args := []interface{}{}
	idx := 1

	if ast.AccountType != "" {
		setClauses = append(setClauses, fmt.Sprintf("ACCOUNT_TYPE = :%d", idx))
		args = append(args, ast.AccountType)
		idx++
	}
	if ast.Gender != "" {
		setClauses = append(setClauses, fmt.Sprintf("GENDER = :%d", idx))
		args = append(args, ast.Gender)
		idx++
	}
	if ast.AccountSubTypeName != "" {
		setClauses = append(setClauses, fmt.Sprintf("ACCOUNT_SUB_TYPE_NAME = :%d", idx))
		args = append(args, ast.AccountSubTypeName)
		idx++
	}
	if ast.AccountSubTypeCode != "" {
		setClauses = append(setClauses, fmt.Sprintf("ACCOUNT_SUB_TYPE_CODE = :%d", idx))
		args = append(args, ast.AccountSubTypeCode)
		idx++
	}

	if len(setClauses) == 0 {
		return nil
	}

	setClauses = append(setClauses, fmt.Sprintf("LAST_MODIFIED_AT = :%d", idx))
	args = append(args, time.Now())
	idx++

	q := fmt.Sprintf("UPDATE ACCOUNT_SUB_TYPES SET %s WHERE ID = HEXTORAW(:%d)", strings.Join(setClauses, ", "), idx)
	args = append(args, id)

	_, err := r.db.ExecContext(ctx, q, args...)
	if err != nil {
		log.Errorf("[AccountSubTypeOracle][Update] failed: %v", err)
		return err
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	q := `UPDATE ACCOUNT_SUB_TYPES SET IS_DELETED = 1, LAST_MODIFIED_AT = CURRENT_TIMESTAMP WHERE ID = HEXTORAW(:1)`
	_, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		log.Errorf("[AccountSubTypeOracle][Delete] failed: %v", err)
		return err
	}
	return nil
}

func (r *Repository) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	var v int
	if enable {
		v = 1
	}
	q := `UPDATE ACCOUNT_SUB_TYPES SET IS_ENABLED = :1, LAST_MODIFIED_AT = CURRENT_TIMESTAMP WHERE ID = HEXTORAW(:2)`
	_, err := r.db.ExecContext(ctx, q, v, id)
	if err != nil {
		log.Errorf("[AccountSubTypeOracle][EnableOrDisable] failed: %v", err)
		return err
	}
	return nil
}

func (r *Repository) FindByID(ctx context.Context, id string) (*imodel.AccountSubType, error) {
	q := `SELECT RAWTOHEX(ID) AS id, ACCOUNT_TYPE, GENDER, ACCOUNT_SUB_TYPE_NAME, ACCOUNT_SUB_TYPE_CODE, IS_ENABLED, IS_DELETED, CREATED_AT, LAST_MODIFIED_AT
		FROM ACCOUNT_SUB_TYPES WHERE ID = HEXTORAW(:1) AND IS_DELETED = 0`
	row := r.db.QueryRowContext(ctx, q, id)
	return r.scanRow(row)
}

func (r *Repository) FindByCodeOrName(ctx context.Context, code, name string) (*imodel.AccountSubType, error) {
	if code == "" && name == "" {
		return nil, nil
	}
	var q string
	var arg interface{}
	if code != "" {
		q = `SELECT RAWTOHEX(ID) AS id, ACCOUNT_TYPE, GENDER, ACCOUNT_SUB_TYPE_NAME, ACCOUNT_SUB_TYPE_CODE, IS_ENABLED, IS_DELETED, CREATED_AT, LAST_MODIFIED_AT
			FROM ACCOUNT_SUB_TYPES WHERE UPPER(ACCOUNT_SUB_TYPE_CODE) = UPPER(:1) AND IS_DELETED = 0`
		arg = code
	} else {
		q = `SELECT RAWTOHEX(ID) AS id, ACCOUNT_TYPE, GENDER, ACCOUNT_SUB_TYPE_NAME, ACCOUNT_SUB_TYPE_CODE, IS_ENABLED, IS_DELETED, CREATED_AT, LAST_MODIFIED_AT
			FROM ACCOUNT_SUB_TYPES WHERE UPPER(ACCOUNT_SUB_TYPE_NAME) = UPPER(:1) AND IS_DELETED = 0`
		arg = name
	}
	row := r.db.QueryRowContext(ctx, q, arg)
	ast, err := r.scanRow(row)
	if err != nil && err == localization.ErrorResourceNotFound {
		return nil, nil
	}
	return ast, err
}

func (r *Repository) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.AccountSubType], error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	var filters []string
	var args []interface{}
	idx := 1

	if filterParam.Search != "" {
		search := "%" + strings.ToUpper(filterParam.Search) + "%"
		filters = append(filters, fmt.Sprintf(`(UPPER(ACCOUNT_SUB_TYPE_NAME) LIKE :%d OR UPPER(ACCOUNT_SUB_TYPE_CODE) LIKE :%d)`, idx, idx+1))
		args = append(args, search, search)
		idx += 2
	}

	whereClause := "IS_DELETED = 0"
	if len(filters) > 0 {
		whereClause += " AND " + strings.Join(filters, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM ACCOUNT_SUB_TYPES WHERE %s", whereClause)
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		log.Errorf("[AccountSubTypeOracle][FindAllWithPagination] count failed: %v", err)
		return nil, err
	}

	page := filterParam.Page
	if page < 1 {
		page = 1
	}
	perPage := filterParam.PerPage
	if perPage < 1 {
		perPage = 10
	}
	offset := (page - 1) * perPage

	if total == 0 {
		meta := local_util.BuildPaginationMeta(0, page, perPage)
		return &types.PaginatedResponse[[]imodel.AccountSubType]{Data: []imodel.AccountSubType{}, Meta: meta}, nil
	}

	selectQuery := fmt.Sprintf(
		`SELECT RAWTOHEX(ID) AS id, ACCOUNT_TYPE, GENDER, ACCOUNT_SUB_TYPE_NAME, ACCOUNT_SUB_TYPE_CODE, IS_ENABLED, IS_DELETED, CREATED_AT, LAST_MODIFIED_AT
		 FROM ACCOUNT_SUB_TYPES WHERE %s ORDER BY CREATED_AT DESC OFFSET :%d ROWS FETCH NEXT :%d ROWS ONLY`,
		whereClause, idx, idx+1,
	)
	args = append(args, offset, perPage)

	rows, err := r.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		log.Errorf("[AccountSubTypeOracle][FindAllWithPagination] query failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	var list []imodel.AccountSubType
	for rows.Next() {
		var ast imodel.AccountSubType
		var createdAt, lastModifiedAt sql.NullTime
		if err := rows.Scan(
			&ast.ID,
			&ast.AccountType,
			&ast.Gender,
			&ast.AccountSubTypeName,
			&ast.AccountSubTypeCode,
			&ast.IsEnabled,
			&ast.IsDeleted,
			&createdAt,
			&lastModifiedAt,
		); err != nil {
			return nil, err
		}
		ast.CreatedAt = formatNullTime(createdAt)
		ast.LastModifiedAt = formatNullTime(lastModifiedAt)
		list = append(list, ast)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	meta := local_util.BuildPaginationMeta(total, page, perPage)
	return &types.PaginatedResponse[[]imodel.AccountSubType]{Data: list, Meta: meta}, nil
}

func (r *Repository) scanRow(row *sql.Row) (*imodel.AccountSubType, error) {
	var ast imodel.AccountSubType
	var createdAt, lastModifiedAt sql.NullTime
	err := row.Scan(
		&ast.ID,
		&ast.AccountType,
		&ast.Gender,
		&ast.AccountSubTypeName,
		&ast.AccountSubTypeCode,
		&ast.IsEnabled,
		&ast.IsDeleted,
		&createdAt,
		&lastModifiedAt,
	)
	if err == sql.ErrNoRows {
		return nil, localization.ErrorResourceNotFound
	}
	if err != nil {
		return nil, err
	}
	ast.CreatedAt = formatNullTime(createdAt)
	ast.LastModifiedAt = formatNullTime(lastModifiedAt)
	return &ast, nil
}

func formatNullTime(t sql.NullTime) string {
	if !t.Valid {
		return ""
	}
	return t.Time.UTC().Format(time.RFC3339)
}
