package sqlc

import (
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
		VALUES (gen_random_uuid(), :1, :2, :3, :4, :5, :6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`
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
	query := `DELETE FROM banks WHERE id = :1`
	_, err := q.db.ExecContext(ctx, query, id)
	return err
}

func (q *Queries) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	query := `UPDATE banks SET is_enabled = :1, update_at = CURRENT_TIMESTAMP WHERE id = :2`
	_, err := q.db.ExecContext(ctx, query, enable, id)
	return err
}

func (q *Queries) FindByID(ctx context.Context, id string) (*imodel.BankOracle, error) {
	query := `SELECT id, bank_name, logo, bic_code, is_enabled, account_length, has_alha_numeric, create_at, update_at FROM banks WHERE id = :1`
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
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &bank, nil
}

func (q *Queries) FindByBIC(ctx context.Context, bic string) (*imodel.BankOracle, error) {
	query := `SELECT id, bank_name, logo, bic_code, is_enabled, account_length, has_alpha_numeric, create_at, update_at FROM banks WHERE bic_code = :1 AND is_enabled = 1`
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
	var conditions []string
	var args []interface{}
	idx := 1
	if name != "" {
		conditions = append(conditions, fmt.Sprintf("bank_name = :%d", idx))
		args = append(args, name)
		idx++
	}
	if bic != "" {
		conditions = append(conditions, fmt.Sprintf("bic_code = :%d", idx))
		args = append(args, bic)
		idx++
	}
	if len(conditions) == 0 {
		return nil, nil
	}
	query := fmt.Sprintf(`SELECT id, bank_name, logo, bic_code, is_enabled, account_length, has_alpha_numeric, create_at, update_at FROM banks WHERE %s`, strings.Join(conditions, " OR "))
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

func (q *Queries) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.BankOracle], error) {
	// Build filtering logic
	var filters []string
	var args []interface{}
	idx := 1

	filters = append(filters, "1=1") // always true, simplifies appending ANDs

	if filterParam.Search != "" && filterParam.Search != "enabled" {
		search := "%" + filterParam.Search + "%"
		filters = append(filters, "(bank_name LIKE :"+fmt.Sprint(idx)+" OR bic_code LIKE :"+fmt.Sprint(idx)+" OR type LIKE :"+fmt.Sprint(idx)+")")
		args = append(args, search)
		idx++
	}
	if filterParam.Search == "enabled" {
		filters = append(filters, "is_enabled = 1")
	}

	whereClause := strings.Join(filters, " AND ")

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM banks WHERE %s", whereClause)
	var total int64
	err := q.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, err
	}

	// Pagination
	offset := (filterParam.Page - 1) * filterParam.PerPage
	limit := filterParam.PerPage
	if int64(offset) >= total {
		limit = 0
	} else if int64(offset)+int64(limit) > total {
		limit = int(total) - offset
	}

	// Fetch paginated results
	selectQuery := fmt.Sprintf(`SELECT id, bank_name, logo, bic_code, is_enabled, account_length, has_alpha_numeric, create_at, update_at FROM banks WHERE %s ORDER BY create_at DESC OFFSET :%d ROWS FETCH NEXT :%d ROWS ONLY`, whereClause, idx, idx+1)
	args = append(args, offset, limit)
	rows, err := q.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	banks := []imodel.BankOracle{}
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
			return nil, err
		}
		banks = append(banks, bank)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	resp := &types.PaginatedResponse[[]imodel.BankOracle]{
		Data: banks,
		Meta: meta,
	}
	return resp, nil
}
