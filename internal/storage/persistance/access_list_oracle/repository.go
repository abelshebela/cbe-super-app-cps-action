package access_list_oracle

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type Repository struct {
	db     *sql.DB
	logger utils.Logger
}

func NewAccessListOracleRepository(db *sql.DB, logger utils.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

func parseBoolFilter(v interface{}) (bool, bool) {
	switch t := v.(type) {
	case bool:
		return t, true
	case string:
		p, err := strconv.ParseBool(strings.TrimSpace(t))
		if err != nil {
			return false, false
		}
		return p, true
	default:
		return false, false
	}
}

func normalizePagination(filterParams types.Filter) (int, int) {
	page := filterParams.Page
	if page <= 0 {
		page = 1
	}
	perPage := filterParams.PerPage
	if perPage <= 0 {
		perPage = 10
	}
	return page, perPage
}

func (r *Repository) FindAllWithPagination(ctx context.Context, filterParams types.Filter) (*types.PaginatedResponse[[]model.APPAccessList], error) {
	page, perPage := normalizePagination(filterParams)
	offset := (page - 1) * perPage

	// ACCESS_LIST may exist without IS_DELETED (legacy / minimal DDL); do not filter on it.
	whereParts := []string{"1=1"}
	args := []interface{}{}
	argIdx := 1

	if strings.TrimSpace(filterParams.Search) != "" {
		search := "%" + strings.ToUpper(strings.TrimSpace(filterParams.Search)) + "%"
		whereParts = append(whereParts, fmt.Sprintf("(UPPER(access_list_name) LIKE :%d OR UPPER(key) LIKE :%d)", argIdx, argIdx))
		args = append(args, search)
		argIdx++
	}

	if filterParams.Filters != nil {
		if v, ok := filterParams.Filters["enabled"]; ok {
			if b, valid := parseBoolFilter(v); valid {
				n := 0
				if b {
					n = 1
				}
				whereParts = append(whereParts, fmt.Sprintf("enabled = :%d", argIdx))
				args = append(args, n)
				argIdx++
			}
		}
		if v, ok := filterParams.Filters["ussd_enabled"]; ok {
			if b, valid := parseBoolFilter(v); valid {
				n := 0
				if b {
					n = 1
				}
				whereParts = append(whereParts, fmt.Sprintf("ussd_enabled = :%d", argIdx))
				args = append(args, n)
				argIdx++
			}
		}
	}

	whereClause := strings.Join(whereParts, " AND ")
	countQuery := "SELECT COUNT(*) FROM ACCESS_LIST WHERE " + whereClause

	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		r.logger.Errorf("[AccessListOracle][FindAllWithPagination] count failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, page, perPage)
	if total == 0 {
		return &types.PaginatedResponse[[]model.APPAccessList]{
			Data: []model.APPAccessList{},
			Meta: meta,
		}, nil
	}

	query := "SELECT key, enabled, access_list_name, ussd_enabled FROM ACCESS_LIST WHERE " + whereClause + fmt.Sprintf(" ORDER BY access_list_name OFFSET :%d ROWS FETCH NEXT :%d ROWS ONLY", argIdx, argIdx+1)
	args = append(args, offset, perPage)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.Errorf("[AccessListOracle][FindAllWithPagination] query failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer rows.Close()

	result := []model.APPAccessList{}
	for rows.Next() {
		var item model.APPAccessList
		var enabled, ussdEnabled int
		if err := rows.Scan(&item.Key, &enabled, &item.AccessListName, &ussdEnabled); err != nil {
			r.logger.Errorf("[AccessListOracle][FindAllWithPagination] scan failed: %v", err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		item.Enabled = enabled == 1
		item.USSDEnabled = ussdEnabled == 1
		result = append(result, item)
	}

	return &types.PaginatedResponse[[]model.APPAccessList]{
		Data: result,
		Meta: meta,
	}, nil
}

func (r *Repository) FindAll(ctx context.Context) ([]model.APPAccessList, error) {
	query := `SELECT key, enabled, access_list_name, ussd_enabled
		FROM ACCESS_LIST
		ORDER BY access_list_name`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		r.logger.Errorf("[AccessListOracle][FindAll] query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer rows.Close()

	result := []model.APPAccessList{}
	for rows.Next() {
		var item model.APPAccessList
		var enabled, ussdEnabled int
		if err := rows.Scan(&item.Key, &enabled, &item.AccessListName, &ussdEnabled); err != nil {
			r.logger.Errorf("[AccessListOracle][FindAll] scan failed: %v", err)
			return nil, local_util.HandleDBError(err)
		}
		item.Enabled = enabled == 1
		item.USSDEnabled = ussdEnabled == 1
		result = append(result, item)
	}

	return result, nil
}

func (r *Repository) Update(ctx context.Context, keys []string, state bool) error {
	if len(keys) == 0 {
		return nil
	}

	next := 0
	if state {
		next = 1
	}

	// Minimal ACCESS_LIST DDL may omit UPDATE_AT; only toggle enabled.
	query := `UPDATE ACCESS_LIST
		SET enabled = :1
		WHERE UPPER(key) = UPPER(:2)`

	for _, key := range keys {
		if _, err := r.db.ExecContext(ctx, query, next, strings.TrimSpace(key)); err != nil {
			r.logger.Errorf("[AccessListOracle][Update] update failed for key=%s err=%v", key, err)
			return errors.New(localization.ErrorFailToUpdateBulkService.Code)
		}
	}

	return nil
}

