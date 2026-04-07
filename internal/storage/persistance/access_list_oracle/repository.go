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

// accessListSelectCols matches ACCESS_LISTS (Oracle): NAME, SERVICE_KEY, flags, timestamps.
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
	countQuery := "SELECT COUNT(*) FROM ACCESS_LISTS WHERE " + whereClause

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

	query := "SELECT " + accessListSelectCols + " FROM ACCESS_LISTS WHERE " + whereClause +
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
		FROM ACCESS_LISTS
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

func (r *Repository) FindAllForSegmentation(ctx context.Context) ([]model.APPAccessList, error) {
	query := `SELECT RAWTOHEX(ID), NAME, SERVICE_KEY, IS_ENABLED, IS_DELETED, CREATED_AT, LAST_MODIFIED_AT, DELETED_AT
		FROM ACCESS_LISTS
		WHERE IS_DELETED = 0
		ORDER BY NAME`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		r.logger.Errorf("[AccessListOracle][FindAllForSegmentation] query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer rows.Close()

	result := []model.APPAccessList{}
	for rows.Next() {
		var idRaw string
		var name, serviceKey string
		var isEn, isDel int
		var createdAt, lastMod, deletedAt sql.NullTime

		if err := rows.Scan(&idRaw, &name, &serviceKey, &isEn, &isDel, &createdAt, &lastMod, &deletedAt); err != nil {
			return nil, err
		}
		if err != nil {
			r.logger.Errorf("[AccessListOracle][FindAllForSegmentation] scan failed: %v", err)
			return nil, local_util.HandleDBError(err)
		}
		result = append(result, model.APPAccessList{
			Key:            idRaw,
			AccessListName: name,
			Enabled:        isEn == 1,
		})
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

	query := `UPDATE ACCESS_LISTS
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

func (r *Repository) FindAllByKeys(ctx context.Context, keys []string) ([]model.APPAccessList, error) {
	r.logger.Infof("[AccessListOracle][FindAllByKeys] checking access list for keys")
	if len(keys) == 0 {
		return nil, nil
	}

	// Build the IN clause with the correct number of bind variables
	inClause := make([]string, len(keys))
	args := make([]interface{}, len(keys))
	for i, key := range keys {
		inClause[i] = fmt.Sprintf("HEXTORAW(:%d)", i+1)
		args[i] = strings.TrimSpace(key)
	}
	query := `SELECT RAWTOHEX(ID), NAME, SERVICE_KEY, IS_ENABLED FROM ACCESS_LISTS WHERE IS_ENABLED = 1 AND IS_DELETED = 0 AND ID IN (` + strings.Join(inClause, ",") + ")"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.Errorf("[AccessListOracle][FindAllKeys] query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer rows.Close()

	// Map from hex key to APPAccessList
	found := make(map[string]model.APPAccessList)
	for rows.Next() {
		var idHex, name, serviceKey string
		var isEn int
		if err := rows.Scan(&idHex, &name, &serviceKey, &isEn); err != nil {
			r.logger.Errorf("[AccessListOracle][FindAllKeys] scan failed: %v", err)
			return nil, local_util.HandleDBError(err)
		}
		found[strings.ToUpper(strings.TrimSpace(idHex))] = model.APPAccessList{
			Key:            serviceKey,
			AccessListName: name,
			Enabled:        isEn == 1,
			USSDEnabled:    false,
		}
	}
	if err := rows.Err(); err != nil {
		r.logger.Errorf("[AccessListOracle][FindAllKeys] rows: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	// Build result, returning the hex key if not found
	result := make([]model.APPAccessList, len(keys))
	for i, key := range keys {
		hexKey := strings.ToUpper(strings.TrimSpace(key))
		if item, ok := found[hexKey]; ok {
			result[i] = item
		} else {
			return nil, errors.New(fmt.Sprintf("Access list with id %s not found", key))
		}
	}
	return result, nil
}

func (r *Repository) FindByKeys(ctx context.Context, keys []string) (map[string]string, error) {
	r.logger.Infof("[AccessListOracle][FindByKeys] checking access list for keys")
	if len(keys) == 0 {
		return map[string]string{}, nil
	}

	// Build the IN clause for SERVICE_KEY
	inClause := make([]string, len(keys))
	args := make([]interface{}, len(keys))
	for i, key := range keys {
		inClause[i] = fmt.Sprintf(":%d", i+1)
		args[i] = strings.TrimSpace(key)
	}
	query := `SELECT SERVICE_KEY, NAME FROM ACCESS_LISTS WHERE IS_DELETED = 0 AND SERVICE_KEY IN (` + strings.Join(inClause, ",") + ")"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.Errorf("[AccessListOracle][FindByKeys] query failed: %v", err)
		return map[string]string{}, local_util.HandleDBError(err)
	}
	defer rows.Close()

	found := make(map[string]string)
	for rows.Next() {
		var serviceKey, name string
		if err := rows.Scan(&serviceKey, &name); err != nil {
			r.logger.Errorf("[AccessListOracle][FindByKeys] scan failed: %v", err)
			return map[string]string{}, local_util.HandleDBError(err)
		}
		found[strings.TrimSpace(serviceKey)] = name
	}
	if err := rows.Err(); err != nil {
		r.logger.Errorf("[AccessListOracle][FindByKeys] rows: %v", err)
		return map[string]string{}, local_util.HandleDBError(err)
	}

	// Check for missing keys and build result
	var missing []string
	for _, key := range keys {
		k := strings.TrimSpace(key)
		if _, ok := found[k]; !ok {
			missing = append(missing, k)
		}
	}
	if len(missing) > 0 {
		msg := "keys: " + strings.Join(missing, ", ") + " not found in access list"
		if len(missing) == 1 {
			msg = "key: " + strings.Join(missing, ", ") + " not found in access list"
		}
		r.logger.Errorf("[AccessListOracle][FindByKeys] %s", msg)
		return map[string]string{}, errors.New(msg)
	}
	return found, nil
}
