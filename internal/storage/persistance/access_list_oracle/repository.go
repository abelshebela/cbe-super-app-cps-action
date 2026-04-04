package access_list_oracle

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type Repository struct {
	db     *sql.DB
	redis  storage.RedisRepository
	logger utils.Logger
}

func NewAccessListOracleRepository(db *sql.DB, redis storage.RedisRepository, logger utils.Logger) *Repository {
	return &Repository{
		db:     db,
		redis:  redis,
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

// accessListSelectCols matches ACCESS_LIST (Oracle): NAME, SERVICE_KEY, flags, timestamps.
const accessListSelectCols = `ID, NAME, SERVICE_KEY, IS_ENABLED, IS_DELETED, CREATED_AT, LAST_MODIFIED_AT, DELETED_AT`

func scanRowToAPPAccessList(scanner interface {
	Scan(dest ...any) error
}) (model.APPAccessList, error) {
	var idRaw []byte
	var name, serviceKey string
	var isEn, isDel int
	var createdAt, lastMod, deletedAt sql.NullTime

	if err := scanner.Scan(&idRaw, &name, &serviceKey, &isEn, &isDel, &createdAt, &lastMod, &deletedAt); err != nil {
		return model.APPAccessList{}, err
	}

	_ = idRaw
	_ = deletedAt
	_ = createdAt
	_ = lastMod
	_ = isDel

	return model.APPAccessList{
		Key:            serviceKey,
		AccessListName: name,
		Enabled:        isEn == 1,
		USSDEnabled:    false,
	}, nil
}

func (r *Repository) FindAllWithPagination(ctx context.Context, filterParams types.Filter) (*types.PaginatedResponse[[]model.APPAccessList], error) {
	page, perPage := normalizePagination(filterParams)
	offset := (page - 1) * perPage

	whereParts := []string{"IS_DELETED = 0"}
	args := []interface{}{}
	argIdx := 1

	if strings.TrimSpace(filterParams.Search) != "" {
		search := "%" + strings.ToUpper(strings.TrimSpace(filterParams.Search)) + "%"
		whereParts = append(whereParts, fmt.Sprintf("(UPPER(NAME) LIKE :%d OR UPPER(SERVICE_KEY) LIKE :%d)", argIdx, argIdx+1))
		args = append(args, search, search)
		argIdx += 2
	}

	if filterParams.Filters != nil {
		if v, ok := filterParams.Filters["enabled"]; ok {
			if b, valid := parseBoolFilter(v); valid {
				n := 0
				if b {
					n = 1
				}
				whereParts = append(whereParts, fmt.Sprintf("IS_ENABLED = :%d", argIdx))
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

	query := "SELECT " + accessListSelectCols + " FROM ACCESS_LIST WHERE " + whereClause +
		fmt.Sprintf(" ORDER BY NAME OFFSET :%d ROWS FETCH NEXT :%d ROWS ONLY", argIdx, argIdx+1)
	args = append(args, offset, perPage)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.Errorf("[AccessListOracle][FindAllWithPagination] query failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer rows.Close()

	result := []model.APPAccessList{}
	for rows.Next() {
		item, err := scanRowToAPPAccessList(rows)
		if err != nil {
			r.logger.Errorf("[AccessListOracle][FindAllWithPagination] scan failed: %v", err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		result = append(result, item)
	}
	if err := rows.Err(); err != nil {
		r.logger.Errorf("[AccessListOracle][FindAllWithPagination] rows: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return &types.PaginatedResponse[[]model.APPAccessList]{
		Data: result,
		Meta: meta,
	}, nil
}

func (r *Repository) FindAll(ctx context.Context) ([]model.APPAccessList, error) {
	query := `SELECT ` + accessListSelectCols + `
		FROM ACCESS_LIST
		WHERE IS_DELETED = 0
		ORDER BY NAME`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		r.logger.Errorf("[AccessListOracle][FindAll] query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer rows.Close()

	result := []model.APPAccessList{}
	for rows.Next() {
		item, err := scanRowToAPPAccessList(rows)
		if err != nil {
			r.logger.Errorf("[AccessListOracle][FindAll] scan failed: %v", err)
			return nil, local_util.HandleDBError(err)
		}
		result = append(result, item)
	}
	if err := rows.Err(); err != nil {
		return nil, local_util.HandleDBError(err)
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

	query := `UPDATE ACCESS_LIST
		SET IS_ENABLED = :1,
		    LAST_MODIFIED_AT = SYSTIMESTAMP
		WHERE UPPER(SERVICE_KEY) = UPPER(:2)` //nolint:goconst // Oracle positional binds

	for _, key := range keys {
		if _, err := r.db.ExecContext(ctx, query, next, strings.TrimSpace(key)); err != nil {
			r.logger.Errorf("[AccessListOracle][Update] update failed for key=%s err=%v", key, err)
			return errors.New(localization.ErrorFailToUpdateBulkService.Code)
		}
	}

	storage.BumpRedisCacheKey(ctx, r.redis, constants.RedisCacheKeyAccessList)
	return nil
}
