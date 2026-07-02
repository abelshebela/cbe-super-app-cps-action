package account_block

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"cbe-super-app-cps-action/internal/constants"
	account_block_dto "cbe-super-app-cps-action/internal/constants/dto/account_block"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"

	local_helper "cbe-super-app-cps-action/internal/storage/persistance/account_block/helper"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const companySelectColumns = `
	RAWTOHEX(id) AS id,
	branch_code,
	branch_name,
	district_name,
	region_name,
	federal_region_name,
	dao_code,
	account_type,
	is_enabled,
	created_at,
	last_modified_at`

const (
	insertCompany = `
		INSERT INTO COMPANY (
			id, branch_code, branch_name, district_name, region_name,
			federal_region_name, dao_code, account_type,
			is_enabled, created_at, last_modified_at
		) VALUES (
			HEXTORAW(:id), :branch_code, :branch_name, :district_name, :region_name,
			:federal_region_name, :dao_code, :account_type,
			1, SYSTIMESTAMP, SYSTIMESTAMP
		)`

	deleteCompanyByID = `
		DELETE FROM COMPANY
		WHERE id = HEXTORAW(:id)`

	selectCompanyByID = `
		SELECT ` + companySelectColumns + `
		FROM COMPANY
		WHERE id = HEXTORAW(:id)`
)

type AccountBlockStorage struct {
	db                 *sql.DB
	client             *mongo.Client
	mongoDB            string
	mongoCpsActionColl string
	redis              storage.RedisRepository
	cpsActionRepo      dal.MongoDal[model.CPSAction, model.CPSAction]
	logger             utils.Logger
}

func NewAccountBlockRepository(client *mongo.Client, cfg *config.VaultConfig, cpsCollection string, db *sql.DB, redis storage.RedisRepository, logger utils.Logger) storage.AccountBlockRepository {
	return &AccountBlockStorage{
		db:                 db,
		client:             client,
		mongoDB:            cfg.MongoDBDatabase,
		mongoCpsActionColl: cpsCollection,
		cpsActionRepo:      dal.NewMongoDal[model.CPSAction, model.CPSAction](client, cfg, cfg.MongoDBDatabase, cpsCollection),
		redis:              redis,
		logger:             logger,
	}
}

func (a *AccountBlockStorage) fetchBlockByID(ctx context.Context, id string) (*imodel.AccountBlock, error) {
	row := a.db.QueryRowContext(ctx, selectCompanyByID, sql.Named("id", id))
	return local_helper.ScanCompanyBranch(row)
}

func (a *AccountBlockStorage) populateParents(blocks []*imodel.AccountBlock) error {
	for _, b := range blocks {
		if b == nil {
			continue
		}
		if b.Type == imodel.TypeBranch {
			local_helper.AttachParentChain(b)
		}
	}
	return nil
}

func (a *AccountBlockStorage) FindByFilterKey(ctx context.Context, field, value string) (*imodel.AccountBlock, error) {
	log := local_util.LoggerFromCtx(ctx, a.logger)

	log.Infof("[AccountBlockStorage][FindByFilterKey] field=%s value=%s", field, value)

	allowed := map[string]string{
		"name":                "branch_name",
		"code":                "branch_code",
		"branch_code":         "branch_code",
		"district_name":       "district_name",
		"federal_region_name": "federal_region_name",
	}
	column, ok := allowed[field]
	if !ok {
		return nil, errors.New("invalid filter key")
	}

	query := fmt.Sprintf(`SELECT %s FROM COMPANY WHERE %s = :val`, companySelectColumns, column)
	row := a.db.QueryRowContext(ctx, query, sql.Named("val", value))
	ab, err := local_helper.ScanCompanyBranch(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		log.Errorf("[AccountBlockStorage][FindByFilterKey] failed: %v", err)
		return nil, err
	}
	return ab, nil
}

func (a *AccountBlockStorage) getBranchesByIDs(ctx context.Context, ids []string) ([]*imodel.AccountBlock, error) {
	log := local_util.LoggerFromCtx(ctx, a.logger)
	if len(ids) == 0 {
		return nil, nil
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, 0, len(ids))
	for i, id := range ids {
		paramName := fmt.Sprintf("id_%d", i)
		placeholders[i] = "HEXTORAW(:" + paramName + ")"
		args = append(args, sql.Named(paramName, id))
	}
	query := fmt.Sprintf(`SELECT %s FROM COMPANY WHERE id IN (%s)`, companySelectColumns, strings.Join(placeholders, ","))

	rows, err := a.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Errorf("[getBranchesByIDs] query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var results []*imodel.AccountBlock
	for rows.Next() {
		ab, err := local_helper.ScanCompanyBranch(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, ab)
	}
	return results, rows.Err()
}

func (a *AccountBlockStorage) getRegionsByNames(ctx context.Context, names []string) ([]*imodel.AccountBlock, error) {
	if len(names) == 0 {
		return nil, nil
	}

	clause, args := local_helper.BuildInClause("federal_region_name", names, "region")
	query := fmt.Sprintf(`
		SELECT federal_region_name,
		       MIN(is_enabled) AS is_enabled,
		       MIN(created_at) AS created_at,
		       MAX(last_modified_at) AS last_modified_at
		FROM COMPANY
		WHERE %s
		GROUP BY federal_region_name`, clause)

	rows, err := a.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*imodel.AccountBlock
	for rows.Next() {
		ab, err := local_helper.ScanRegionAggregate(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, ab)
	}
	return results, rows.Err()
}

func (a *AccountBlockStorage) getDistrictsByNames(ctx context.Context, names []string) ([]*imodel.AccountBlock, error) {
	if len(names) == 0 {
		return nil, nil
	}

	clause, args := local_helper.BuildInClause("district_name", names, "district")
	query := fmt.Sprintf(`
		SELECT federal_region_name, district_name,
		       MIN(is_enabled) AS is_enabled,
		       MIN(created_at) AS created_at,
		       MAX(last_modified_at) AS last_modified_at
		FROM COMPANY
		WHERE %s
		GROUP BY federal_region_name, district_name`, clause)

	rows, err := a.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*imodel.AccountBlock
	for rows.Next() {
		ab, err := local_helper.ScanDistrictAggregate(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, ab)
	}
	return results, rows.Err()
}

func (a *AccountBlockStorage) GetBranchesByIds(ctx context.Context, ids []string) ([]*imodel.AccountBlock, error) {
	log := local_util.LoggerFromCtx(ctx, a.logger)
	log.Infof("[AccountBlockStorage][GetBranchesByIds] fetching %d branches", len(ids))
	results, err := a.getBranchesByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	if err := a.populateParents(results); err != nil {
		return nil, err
	}
	return results, nil
}

func (a *AccountBlockStorage) GetRegionsByIds(ctx context.Context, ids []string) ([]*imodel.AccountBlock, error) {
	log := local_util.LoggerFromCtx(ctx, a.logger)
	log.Infof("[AccountBlockStorage][GetRegionsByIds] fetching %d regions", len(ids))
	return a.getRegionsByNames(ctx, ids)
}

func (a *AccountBlockStorage) GetDistrictsByIds(ctx context.Context, ids []string) ([]*imodel.AccountBlock, error) {
	log := local_util.LoggerFromCtx(ctx, a.logger)
	log.Infof("[AccountBlockStorage][GetDistrictsByIds] fetching %d districts", len(ids))
	return a.getDistrictsByNames(ctx, ids)
}

func (a *AccountBlockStorage) findAllBranchesWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*imodel.AccountBlock], error) {
	offset := (filterParam.Page - 1) * filterParam.PerPage
	if offset < 0 {
		offset = 0
	}
	limit := filterParam.PerPage
	if limit <= 0 {
		limit = 50
	}

	var search interface{}
	var isEnabledFilter interface{}
	var regionNames, districtNames []string
	var isEnableCheck *bool

	if filterParam.Search != "" {
		search = filterParam.Search
	}
	if filterParam.Filters != nil {
		if v, ok := filterParam.Filters["region_id"]; ok {
			regionNames = local_helper.NamesFromFilterValue(v)
		}
		if v, ok := filterParam.Filters["district_id"]; ok {
			districtNames = local_helper.NamesFromFilterValue(v)
		}
		if v, ok := filterParam.Filters["is_enabled"]; ok {
			if enabled, isBool := v.(bool); isBool {
				isEnabledFilter = local_util.BoolToInt(enabled)
			}
		}
		if v, ok := filterParam.Filters["status_check"]; ok {
			if boolVal, ok := v.(bool); ok {
				isEnableCheck = &boolVal
			}
		}
	}

	var filterClauses []string
	var args []interface{}
	if clause, clauseArgs := local_helper.BuildInClause("federal_region_name", regionNames, "region"); clause != "" {
		filterClauses = append(filterClauses, clause)
		args = append(args, clauseArgs...)
	}
	if clause, clauseArgs := local_helper.BuildInClause("district_name", districtNames, "district"); clause != "" {
		filterClauses = append(filterClauses, clause)
		args = append(args, clauseArgs...)
	}

	filterSQL := ""
	if len(filterClauses) > 0 {
		filterSQL = " AND " + strings.Join(filterClauses, " AND ")
	}

	nameSearchCondition := "LOWER(TRIM(branch_name)) LIKE '%%' || LOWER(TRIM(:search)) || '%%'"
	if isEnableCheck != nil {
		nameSearchCondition = "LOWER(TRIM(branch_name)) = LOWER(TRIM(:search))"
	}

	query := fmt.Sprintf(`
		SELECT %s, COUNT(*) OVER() AS total_count
		FROM COMPANY
		WHERE 1=1
		  %s
		  AND (:search IS NULL
		       OR %s
		       OR LOWER(branch_code) LIKE '%%' || LOWER(:search) || '%%'
		       OR LOWER(region_name) LIKE '%%' || LOWER(:search) || '%%'
		       OR LOWER(district_name) LIKE '%%' || LOWER(:search) || '%%'
		       OR LOWER(federal_region_name) LIKE '%%' || LOWER(:search) || '%%')
		  AND (:is_enabled IS NULL OR is_enabled = :is_enabled)
		ORDER BY created_at DESC
		OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`, companySelectColumns, filterSQL, nameSearchCondition)

	args = append(args,
		sql.Named("search", search),
		sql.Named("is_enabled", isEnabledFilter),
		sql.Named("offset", offset),
		sql.Named("limit", limit),
	)

	return a.scanPaginatedBranches(ctx, query, args, filterParam, offset)
}

func (a *AccountBlockStorage) findAllRegionsWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*imodel.AccountBlock], error) {
	offset := (filterParam.Page - 1) * filterParam.PerPage
	if offset < 0 {
		offset = 0
	}
	limit := filterParam.PerPage
	if limit <= 0 {
		limit = 50
	}

	var search interface{}
	var isEnabledFilter interface{}

	if filterParam.Search != "" {
		search = filterParam.Search
	}
	if filterParam.Filters != nil {
		if v, ok := filterParam.Filters["is_enabled"]; ok {
			if enabled, isBool := v.(bool); isBool {
				isEnabledFilter = local_util.BoolToInt(enabled)
			}
		}
	}

	query := `
		WITH regions AS (
			SELECT federal_region_name,
			       MIN(is_enabled) AS is_enabled,
			       MIN(created_at) AS created_at,
			       MAX(last_modified_at) AS last_modified_at
			FROM COMPANY
			WHERE (:search IS NULL
			       OR LOWER(federal_region_name) LIKE '%' || LOWER(:search) || '%'
			       OR LOWER(region_name) LIKE '%' || LOWER(:search) || '%')
			  AND federal_region_name IS NOT NULL
			GROUP BY federal_region_name
		)
		SELECT federal_region_name, is_enabled, created_at, last_modified_at,
		       COUNT(*) OVER() AS total_count
		FROM regions
		WHERE (:is_enabled IS NULL OR is_enabled = :is_enabled)
		ORDER BY federal_region_name
		OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`

	args := []interface{}{
		sql.Named("search", search),
		sql.Named("is_enabled", isEnabledFilter),
		sql.Named("offset", offset),
		sql.Named("limit", limit),
	}

	return a.scanPaginatedRegions(ctx, query, args, filterParam, offset)
}

func (a *AccountBlockStorage) findAllDistrictsWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*imodel.AccountBlock], error) {
	offset := (filterParam.Page - 1) * filterParam.PerPage
	if offset < 0 {
		offset = 0
	}
	limit := filterParam.PerPage
	if limit <= 0 {
		limit = 50
	}

	var search interface{}
	var isEnabledFilter interface{}
	var regionNames []string

	if filterParam.Search != "" {
		search = filterParam.Search
	}
	if filterParam.Filters != nil {
		if v, ok := filterParam.Filters["region_id"]; ok {
			regionNames = local_helper.NamesFromFilterValue(v)
		}
		if v, ok := filterParam.Filters["is_enabled"]; ok {
			if enabled, isBool := v.(bool); isBool {
				isEnabledFilter = local_util.BoolToInt(enabled)
			}
		}
	}

	var filterClauses []string
	var args []interface{}
	if clause, clauseArgs := local_helper.BuildInClause("federal_region_name", regionNames, "region"); clause != "" {
		filterClauses = append(filterClauses, clause)
		args = append(args, clauseArgs...)
	}

	filterSQL := ""
	if len(filterClauses) > 0 {
		filterSQL = " AND " + strings.Join(filterClauses, " AND ")
	}

	query := fmt.Sprintf(`
		WITH districts AS (
			SELECT federal_region_name, district_name,
			       MIN(is_enabled) AS is_enabled,
			       MIN(created_at) AS created_at,
			       MAX(last_modified_at) AS last_modified_at
			FROM COMPANY
			WHERE (:search IS NULL
			       OR LOWER(district_name) LIKE '%%' || LOWER(:search) || '%%')
			  AND district_name IS NOT NULL
			  %s
			GROUP BY federal_region_name, district_name
		)
		SELECT federal_region_name, district_name, is_enabled, created_at, last_modified_at,
		       COUNT(*) OVER() AS total_count
		FROM districts
		WHERE (:is_enabled IS NULL OR is_enabled = :is_enabled)
		ORDER BY district_name
		OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`, filterSQL)

	args = append(args,
		sql.Named("search", search),
		sql.Named("is_enabled", isEnabledFilter),
		sql.Named("offset", offset),
		sql.Named("limit", limit),
	)

	return a.scanPaginatedDistricts(ctx, query, args, filterParam, offset)
}

func buildPaginationMeta(filterParam types.Filter, offset int, totalCount int64) types.PaginationMeta {
	totalPages := int((totalCount + int64(filterParam.PerPage) - 1) / int64(filterParam.PerPage))
	var prevPage, nextPage *int
	if filterParam.Page > 1 {
		p := filterParam.Page - 1
		prevPage = &p
	}
	if filterParam.Page < totalPages {
		n := filterParam.Page + 1
		nextPage = &n
	}
	return types.PaginationMeta{
		TotalDocs:     totalCount,
		Limit:         filterParam.PerPage,
		TotalPages:    totalPages,
		Page:          filterParam.Page,
		PagingCounter: offset + 1,
		HasPrevPage:   filterParam.Page > 1,
		HasNextPage:   filterParam.Page < totalPages,
		PrevPage:      prevPage,
		NextPage:      nextPage,
	}
}

func (a *AccountBlockStorage) scanPaginatedBranches(ctx context.Context, query string, args []interface{}, filterParam types.Filter, offset int) (*types.PaginatedResponse[[]*imodel.AccountBlock], error) {
	rows, err := a.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*imodel.AccountBlock
	var totalCount int64
	for rows.Next() {
		var ab imodel.AccountBlock
		var districtName, regionName, federalRegionName sql.NullString
		var isEnabledInt int
		var createdAt, updatedAt sql.NullTime

		if err := rows.Scan(
			&ab.ID, &ab.Code, &ab.Name,
			&districtName, &regionName, &federalRegionName,
			&ab.DaoCode, &ab.AccountType, &isEnabledInt,
			&createdAt, &updatedAt,
			&totalCount,
		); err != nil {
			return nil, err
		}

		ab.Type = imodel.TypeBranch
		ab.IsEnabled = isEnabledInt == 1
		if districtName.Valid {
			ab.DistrictName = districtName.String
		}
		if regionName.Valid {
			ab.RegionName = regionName.String
		}
		if federalRegionName.Valid {
			ab.FederalRegionName = federalRegionName.String
		}
		if createdAt.Valid {
			ab.CreatedAt = createdAt.Time
		}
		if updatedAt.Valid {
			ab.UpdatedAt = updatedAt.Time
		}
		results = append(results, &ab)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if err := a.populateParents(results); err != nil {
		return nil, err
	}

	return &types.PaginatedResponse[[]*imodel.AccountBlock]{
		Data: results,
		Meta: buildPaginationMeta(filterParam, offset, totalCount),
	}, nil
}

func (a *AccountBlockStorage) scanPaginatedRegions(ctx context.Context, query string, args []interface{}, filterParam types.Filter, offset int) (*types.PaginatedResponse[[]*imodel.AccountBlock], error) {
	rows, err := a.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*imodel.AccountBlock
	var totalCount int64
	for rows.Next() {
		var ab imodel.AccountBlock
		var isEnabledInt int
		var createdAt, updatedAt sql.NullTime

		if err := rows.Scan(
			&ab.FederalRegionName,
			&isEnabledInt, &createdAt, &updatedAt,
			&totalCount,
		); err != nil {
			return nil, err
		}

		ab.ID = ab.FederalRegionName
		ab.Name = ab.FederalRegionName
		ab.Type = imodel.TypeRegion
		ab.IsEnabled = isEnabledInt == 1
		frn := ab.FederalRegionName
		ab.RegionID = &frn
		if createdAt.Valid {
			ab.CreatedAt = createdAt.Time
		}
		if updatedAt.Valid {
			ab.UpdatedAt = updatedAt.Time
		}
		results = append(results, &ab)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &types.PaginatedResponse[[]*imodel.AccountBlock]{
		Data: results,
		Meta: buildPaginationMeta(filterParam, offset, totalCount),
	}, nil
}

func (a *AccountBlockStorage) scanPaginatedDistricts(ctx context.Context, query string, args []interface{}, filterParam types.Filter, offset int) (*types.PaginatedResponse[[]*imodel.AccountBlock], error) {
	rows, err := a.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*imodel.AccountBlock
	var totalCount int64
	for rows.Next() {
		var ab imodel.AccountBlock
		var isEnabledInt int
		var createdAt, updatedAt sql.NullTime

		var federalRegionName sql.NullString

		if err := rows.Scan(
			&federalRegionName, &ab.DistrictName,
			&isEnabledInt, &createdAt, &updatedAt,
			&totalCount,
		); err != nil {
			return nil, err
		}

		ab.ID = ab.DistrictName
		ab.Name = ab.DistrictName
		ab.Type = imodel.TypeDistrict
		ab.IsEnabled = isEnabledInt == 1
		if federalRegionName.Valid {
			ab.FederalRegionName = federalRegionName.String
			frn := federalRegionName.String
			ab.RegionID = &frn
		}
		dn := ab.DistrictName
		ab.DistrictID = &dn
		if createdAt.Valid {
			ab.CreatedAt = createdAt.Time
		}
		if updatedAt.Valid {
			ab.UpdatedAt = updatedAt.Time
		}
		if ab.FederalRegionName != "" {
			ab.Parent = &imodel.AccountBlock{
				ID:                ab.FederalRegionName,
				Name:              ab.FederalRegionName,
				FederalRegionName: ab.FederalRegionName,
				Type:              imodel.TypeRegion,
			}
		}
		results = append(results, &ab)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &types.PaginatedResponse[[]*imodel.AccountBlock]{
		Data: results,
		Meta: buildPaginationMeta(filterParam, offset, totalCount),
	}, nil
}

func (a *AccountBlockStorage) FindAllBranchesWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*imodel.AccountBlock], error) {
	log := local_util.LoggerFromCtx(ctx, a.logger)
	log.Infof("[AccountBlockStorage][FindAllBranchesWithPagination] fetching branches")
	return a.findAllBranchesWithPagination(ctx, filterParam)
}

func (a *AccountBlockStorage) FindAllRegionsWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*imodel.AccountBlock], error) {
	log := local_util.LoggerFromCtx(ctx, a.logger)
	log.Infof("[AccountBlockStorage][FindAllRegionsWithPagination] fetching regions")
	return a.findAllRegionsWithPagination(ctx, filterParam)
}

func (a *AccountBlockStorage) FindAllDistrictsWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*imodel.AccountBlock], error) {
	log := local_util.LoggerFromCtx(ctx, a.logger)
	log.Infof("[AccountBlockStorage][FindAllDistrictsWithPagination] fetching districts")
	return a.findAllDistrictsWithPagination(ctx, filterParam)
}

func (a *AccountBlockStorage) enableOrDisableByColumn(ctx context.Context, blockType, column string, values []string, reason *types.Reason, enabled bool) error {
	if len(values) == 0 {
		return nil
	}

	clause, args := local_helper.BuildInClause(column, values, "val")
	args = append(args, sql.Named("is_enabled", local_util.BoolToInt(enabled)))

	if !enabled {
		if err := a.insertDisableReasonForBlocks(ctx, blockType, values, reason); err != nil {
			return err
		}
	}

	query := fmt.Sprintf(`UPDATE COMPANY
		SET is_enabled = :is_enabled,
		    last_modified_at = SYSTIMESTAMP
		WHERE %s`, clause)

	_, err := a.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	storage.BumpRedisCacheKey(ctx, a.redis, constants.RedisCacheKeyAccountBlock)
	return nil
}

func (a *AccountBlockStorage) branchIDsForColumnValues(ctx context.Context, column string, values []string) ([]string, error) {
	clause, args := local_helper.BuildInClause(column, values, "val")
	query := fmt.Sprintf(`SELECT RAWTOHEX(id) FROM COMPANY WHERE %s`, clause)
	rows, err := a.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (a *AccountBlockStorage) enableOrDisableByIDs(ctx context.Context, blockType string, ids []string, reason *types.Reason, enabled bool) error {
	if len(ids) == 0 {
		return nil
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, 0, len(ids)+1)
	for i, id := range ids {
		paramName := fmt.Sprintf("id_%d", i)
		placeholders[i] = "HEXTORAW(:" + paramName + ")"
		args = append(args, sql.Named(paramName, id))
	}
	args = append(args, sql.Named("is_enabled", local_util.BoolToInt(enabled)))

	if !enabled {
		if err := a.insertDisableReasonForBlocks(ctx, blockType, ids, reason); err != nil {
			return err
		}
	}

	query := fmt.Sprintf(`UPDATE COMPANY
		SET is_enabled = :is_enabled,
		    last_modified_at = SYSTIMESTAMP
		WHERE id IN (%s)`, strings.Join(placeholders, ","))

	_, err := a.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	storage.BumpRedisCacheKey(ctx, a.redis, constants.RedisCacheKeyAccountBlock)
	return nil
}

func (a *AccountBlockStorage) EnableOrDisableBranches(ctx context.Context, ids []string, reason *types.Reason, enabled bool) error {
	log := local_util.LoggerFromCtx(ctx, a.logger)
	log.Infof("[AccountBlockStorage][EnableOrDisableBranches] ids=%v enabled=%v", ids, enabled)
	err := a.enableOrDisableByIDs(ctx, "B", ids, reason, enabled)
	if err != nil {
		log.Errorf("[AccountBlockStorage][EnableOrDisableBranches] failed: %v", err)
	}
	return err
}

func (a *AccountBlockStorage) EnableOrDisableRegions(ctx context.Context, ids []string, reason *types.Reason, enabled bool) error {
	log := local_util.LoggerFromCtx(ctx, a.logger)
	log.Infof("[AccountBlockStorage][EnableOrDisableRegions] ids=%v enabled=%v", ids, enabled)
	err := a.enableOrDisableByColumn(ctx, "R", "federal_region_name", ids, reason, enabled)
	if err != nil {
		log.Errorf("[AccountBlockStorage][EnableOrDisableRegions] failed: %v", err)
	}
	return err
}

func (a *AccountBlockStorage) EnableOrDisableDistricts(ctx context.Context, ids []string, reason *types.Reason, enabled bool) error {
	log := local_util.LoggerFromCtx(ctx, a.logger)
	log.Infof("[AccountBlockStorage][EnableOrDisableDistricts] ids=%v enabled=%v", ids, enabled)
	err := a.enableOrDisableByColumn(ctx, "D", "district_name", ids, reason, enabled)
	if err != nil {
		log.Errorf("[AccountBlockStorage][EnableOrDisableDistricts] failed: %v", err)
	}
	return err
}

func (a *AccountBlockStorage) GetAccountBlockDetails(ctx context.Context, id string, filterParam types.Filter) (*types.PaginatedResponse[[]account_block_dto.AccountBlockActionResponse], error) {
	log := local_util.LoggerFromCtx(ctx, a.logger)

	log.Infof("[AccountBlockStorage][GetAccountBlockDetails] fetching CPS actions for account block id: %s", id)

	cpsCollection := a.client.Database(a.mongoDB).Collection(a.mongoCpsActionColl)

	accountBlockRequestActions := []string{
		string(constants.RequestEnableBranches),
		string(constants.RequestDisableBranches),
		string(constants.RequestEnableCities),
		string(constants.RequestDisableCities),
		string(constants.RequestEnableDistricts),
		string(constants.RequestDisableDistricts),
		string(constants.RequestEnableRegions),
		string(constants.RequestDisableRegions),
	}

	matchFilter := bson.M{
		"is_deleted":         false,
		"previous_action.id": id,
		"request_action":     bson.M{"$in": accountBlockRequestActions},
	}

	skip := int64((filterParam.Page - 1) * filterParam.PerPage)
	limit := int64(filterParam.PerPage)

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: matchFilter}},
		{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
		{{Key: "$skip", Value: skip}},
		{{Key: "$limit", Value: limit}},
	}

	cur, err := cpsCollection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Errorf("[AccountBlockStorage][GetAccountBlockDetails] aggregation failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = cur.Close(ctx) }()

	var results []model.CPSAction
	if err := cur.All(ctx, &results); err != nil {
		log.Errorf("[AccountBlockStorage][GetAccountBlockDetails] failed to decode CPS actions: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	if len(results) == 0 {
		log.Errorf("[AccountBlockStorage][GetAccountBlockDetails] no CPS actions found for id: %s", id)
		return nil, errors.New(localization.ErrorActionNotFound.Code)
	}

	total, err := cpsCollection.CountDocuments(ctx, matchFilter)
	if err != nil {
		log.Errorf("[GetAccountBlockDetails] failed to count CPS actions: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	response := make([]account_block_dto.AccountBlockActionResponse, 0, len(results))
	for _, cpsAction := range results {
		actionResponse := account_block_dto.AccountBlockActionResponse{
			ID:                  cpsAction.ID.Hex(),
			ActionCode:          cpsAction.ActionCode,
			UniqueId:            cpsAction.UniqueId,
			MakerID:             cpsAction.MakerID,
			MakerName:           cpsAction.MakerName,
			MakerPhoneNumber:    cpsAction.MakerPhoneNumber,
			CheckerUsers:        local_helper.ConvertCheckers(cpsAction.CheckerUsers),
			AuditorUsers:        local_helper.ConvertAuditors(cpsAction.AuditorUsers),
			AuditorCount:        cpsAction.AuditorCount,
			AuditorStatus:       account_block_dto.AuditorStatus(cpsAction.AuditorStatus),
			CurrentAuditorIndex: cpsAction.CurrentAuditorIndex,
			CheckerCount:        cpsAction.CheckerCount,
			CurrentCheckerIndex: cpsAction.CurrentCheckerIndex,
			RoleCode:            cpsAction.RoleCode,
			RejectionReason:     cpsAction.RejectionReason,
			CanceledReason:      cpsAction.CanceledReason,
			ActionStatus:        cpsAction.ActionStatus,
			ActionType:          cpsAction.ActionType,
			IsDeleted:           cpsAction.IsDeleted,
			RequestAction:       cpsAction.RequestAction,
			Version:             cpsAction.Version,
			ReversedByRoleID:    cpsAction.ReversedByRoleID,
			ReversedByID:        cpsAction.ReversedByID,
			ReversedByName:      cpsAction.ReversedByName,
			ReversedAt:          cpsAction.ReversedAt,
			CreatedAt:           cpsAction.CreatedAt,
			LastModifiedAt:      cpsAction.LastModifiedAt,
			MakerActionTime:     cpsAction.MakerActionTime,
		}

		var previousAction interface{}
		if cpsAction.ActionStatus == string(constants.ActionApproved) {
			previousAction = getMatchingAction(cpsAction.CurrentAction, id, a.logger)
		} else if cpsAction.ActionStatus == string(constants.ActionPending) || cpsAction.ActionStatus == string(constants.ActionRejected) {
			previousAction = getMatchingAction(cpsAction.PreviousAction, id, a.logger)
		} else {
			previousAction = getMatchingAction(cpsAction.PreviousAction, id, a.logger)
		}

		actionResponse.PreviousAction = previousAction
		response = append(response, actionResponse)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	log.Infof("[AccountBlockStorage][GetAccountBlockDetails] successfully mapped %d CPS actions", len(response))

	return &types.PaginatedResponse[[]account_block_dto.AccountBlockActionResponse]{
		Data: response,
		Meta: meta,
	}, nil
}

func getMatchingAction(actionData interface{}, accountBlockID string, logger utils.Logger) interface{} {
	if actionData == nil {
		return nil
	}

	actions, err := local_util.JsonUnmarshal[[]types.EnableDisableAction](actionData)
	if err == nil && actions != nil {
		for _, action := range *actions {
			if action.ID == accountBlockID {
				return action
			}
		}
		return nil
	}

	singleAction, err := local_util.JsonUnmarshal[types.EnableDisableAction](actionData)
	if err == nil && singleAction != nil {
		if singleAction.ID == accountBlockID {
			return *singleAction
		}
		return nil
	}

	if actionMap, ok := actionData.(map[string]interface{}); ok {
		if id, exists := actionMap["id"]; exists {
			if idStr, ok := id.(string); ok && idStr == accountBlockID {
				return actionMap
			}
		}
	}

	logger.Warnf("[AccountBlockStorage][getMatchingAction] failed to extract matching action for id: %s", accountBlockID)
	return nil
}

func (a *AccountBlockStorage) GetAllBranchesByDistrictOrRegion(ctx context.Context, id string) ([]imodel.AccountBlock, error) {
	query := `SELECT ` + companySelectColumns + `
		FROM COMPANY
		WHERE federal_region_name = :id OR district_name = :id`

	rows, err := a.db.QueryContext(ctx, query, sql.Named("id", id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []imodel.AccountBlock
	for rows.Next() {
		ab, err := local_helper.ScanCompanyBranch(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, *ab)
	}
	ptrs := make([]*imodel.AccountBlock, len(results))
	for i := range results {
		ptrs[i] = &results[i]
	}
	if err := a.populateParents(ptrs); err != nil {
		return nil, err
	}
	return results, rows.Err()
}

func (a *AccountBlockStorage) GetBranchByIds(ctx context.Context, id string) (*imodel.AccountBlock, error) {
	log := local_util.LoggerFromCtx(ctx, a.logger)

	log.Infof("[AccountBlockStorage][GetBranchByIds] fetching branch by id: %s", id)
	ab, err := a.fetchBlockByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Errorf("[AccountBlockStorage][GetBranchByIds] branch not found")
			return nil, errors.New(localization.ErrorBranchNotFound.Code)
		}
		log.Errorf("[AccountBlockStorage][GetBranchByIds] failed to fetch branch: %v", err)
		return nil, err
	}
	log.Infof("[AccountBlockStorage][GetBranchByIds] branch retrieved successfully")
	return ab, nil
}

func (a *AccountBlockStorage) GetBranchByCode(ctx context.Context, code string) (*imodel.AccountBlock, error) {
	log := local_util.LoggerFromCtx(ctx, a.logger)
	log.Infof("[AccountBlockStorage][GetBranchByCode] fetching branch by Code: %s", code)

	query := `SELECT ` + companySelectColumns + `
		FROM COMPANY
		WHERE branch_code = :code
		  AND is_enabled = 1`

	row := a.db.QueryRowContext(ctx, query, sql.Named("code", code))
	branch, err := local_helper.ScanCompanyBranch(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Infof("[AccountBlockStorage][GetBranchByCode] branch not found: %s", code)
			return nil, nil
		}
		log.Errorf("[AccountBlockStorage][GetBranchByCode] failed to fetch branch: %v", err)
		return nil, err
	}

	return branch, nil
}
