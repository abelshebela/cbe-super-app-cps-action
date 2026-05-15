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

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"github.com/google/uuid"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// ─── SQL constants ───────────────────────────────────────────────────────────

const (
	insertAccountBlock = `
		INSERT INTO ACCOUNT_BLOCKS (
			id, name, code, address, slug, type,
			is_enabled, district_id, region_id,
			is_deleted, created_at, updated_at
		) VALUES (
			:id, :name, :code, :address, :slug, :type,
			1, :district_id, :region_id,
			0, SYSTIMESTAMP, SYSTIMESTAMP
		)`

	softDeleteAccountBlock = `
		UPDATE ACCOUNT_BLOCKS
		SET is_deleted = 1,
		    updated_at = SYSTIMESTAMP
		WHERE id = HEXTORAW(:id)
		  AND is_deleted = 0`

	selectAccountBlockByID = `
		SELECT
			RAWTOHEX(id) AS id,
			name,
			code,
			address,
			slug,
			type,
			is_enabled,
			RAWTOHEX(district_id) AS district_id,
			RAWTOHEX(region_id) AS region_id,
			is_deleted,
			created_at,
			updated_at
		FROM ACCOUNT_BLOCKS
		WHERE id = HEXTORAW(:id)
		  AND is_deleted = 0`

	selectWithAncestors = `
		SELECT
			RAWTOHEX(id) AS id,
			name,
			code,
			address,
			slug,
			type,
			is_enabled,
			RAWTOHEX(district_id) AS district_id,
			RAWTOHEX(region_id) AS region_id,
			is_deleted,
			created_at,
			updated_at,
			LEVEL AS depth
		FROM ACCOUNT_BLOCKS
		WHERE is_deleted = 0
		START WITH id = HEXTORAW(:id)
		ORDER BY LEVEL ASC`

	listAccountBlocksByType = `
		SELECT
			RAWTOHEX(id) AS id,
			name,
			code,
			address,
			slug,
			type,
			is_enabled,
			RAWTOHEX(district_id) AS district_id,
			RAWTOHEX(region_id) AS region_id,
			is_deleted,
			created_at,
			updated_at,
			COUNT(*) OVER() AS total_count
		FROM ACCOUNT_BLOCKS
		WHERE type = :type
		  AND is_deleted = 0
		  AND (:search IS NULL
		       OR LOWER(name) LIKE '%' || LOWER(:search) || '%'
		       OR LOWER(code) LIKE '%' || LOWER(:search) || '%'
		       OR LOWER(address) LIKE '%' || LOWER(:search) || '%')
		  AND (:region_id IS NULL OR region_id = HEXTORAW(:region_id))
		  AND (:district_id IS NULL OR district_id = HEXTORAW(:district_id))
		  AND (:is_enabled IS NULL OR is_enabled = :is_enabled)
		ORDER BY created_at DESC
		OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`
)

// ─── Repository struct ──────────────────────────────────────────────────────

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

// ─── Scan helpers ───────────────────────────────────────────────────────────
func scanAccountBlock(scanner interface{ Scan(dest ...any) error }) (*imodel.AccountBlock, error) {
	var ab imodel.AccountBlock
	var districtID, regionID sql.NullString
	var isEnabledInt, isDeletedInt int
	var createdAt, updatedAt sql.NullTime

	err := scanner.Scan(
		&ab.ID, &ab.Name, &ab.Code, &ab.Address,
		&ab.Slug, &ab.Type,
		&isEnabledInt, &districtID, &regionID,
		&isDeletedInt, &createdAt, &updatedAt,
	)
	if err != nil {
		return nil, err
	}

	ab.IsEnabled = isEnabledInt == 1
	ab.IsDeleted = isDeletedInt == 1

	if createdAt.Valid {
		ab.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		ab.UpdatedAt = updatedAt.Time
	}
	if districtID.Valid {
		ab.DistrictID = &districtID.String
	}
	if regionID.Valid {
		ab.RegionID = &regionID.String
	}

	return &ab, nil
}

func isEnabledToInt(enabled bool) int {
	if enabled {
		return 1
	}
	return 0
}

func nullStr(s *string) interface{} {
	if s == nil {
		return nil
	}
	return *s
}

// nullIfEmptyFilter treats "", whitespace-only strings as SQL NULL for optional id filters.
func nullIfEmptyFilter(v interface{}) interface{} {
	if v == nil {
		return nil
	}
	if s, ok := v.(string); ok && strings.TrimSpace(s) == "" {
		return nil
	}
	return v
}

const maxParentDepth = 64

func (a *AccountBlockStorage) fetchBlockByID(ctx context.Context, id string) (*imodel.AccountBlock, error) {
	row := a.db.QueryRowContext(ctx, selectAccountBlockByID, sql.Named("id", id))
	return scanAccountBlockFromRow(row)
}

func directParentIDByType(b *imodel.AccountBlock) string {
	if b == nil {
		return ""
	}

	trimmed := func(s *string) string {
		if s == nil {
			return ""
		}
		return strings.TrimSpace(*s)
	}

	switch b.Type {
	case imodel.TypeRegion:
		// Region is the top-most level in current hierarchy.
		return ""
	case imodel.TypeDistrict:
		return trimmed(b.RegionID)
	// case imodel.TypeCity:
	// 	// City usually belongs to a district; fallback to region if needed.
	// 	if id := trimmed(b.DistrictID); id != "" {
	// 		return id
	// 	}
	// 	return trimmed(b.RegionID)
	case imodel.TypeBranch:
		// Branch should resolve to district first, then city, then region.
		if id := trimmed(b.DistrictID); id != "" {
			return id
		}
		// if id := trimmed(b.CityID); id != "" {
		// 	return id
		// }
		return trimmed(b.RegionID)
	default:
		return ""
	}
}

// populateParentChain resolves hierarchy from type-specific foreign keys and sets Parent recursively.
func (a *AccountBlockStorage) populateParentChain(ctx context.Context, b *imodel.AccountBlock, depth int) error {
	if b == nil || depth > maxParentDepth {
		return nil
	}
	pid := directParentIDByType(b)
	if pid == "" {
		return nil
	}
	parent, err := a.fetchBlockByID(ctx, pid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}
	if err := a.populateParentChain(ctx, parent, depth+1); err != nil {
		return err
	}

	// For branches, keep district as direct parent and make sure region is attached
	// as the next parent even when district row is incomplete.
	if b.Type == imodel.TypeBranch &&
		parent.Type == imodel.TypeDistrict &&
		parent.Parent == nil &&
		b.RegionID != nil {
		rid := strings.TrimSpace(*b.RegionID)
		if rid != "" && rid != parent.ID {
			region, rErr := a.fetchBlockByID(ctx, rid)
			if rErr == nil {
				parent.Parent = region
			} else if !errors.Is(rErr, sql.ErrNoRows) {
				return rErr
			}
		}
	}

	b.Parent = parent
	return nil
}

func (a *AccountBlockStorage) populateParentsAndReasons(ctx context.Context, blocks []*imodel.AccountBlock) error {
	a.logger.Debugf("[populateParentsAndReasons] called with %d blocks", len(blocks))
	for i, b := range blocks {
		if b == nil {
			a.logger.Warnf("[populateParentsAndReasons] block at index %d is nil, skipping", i)
			continue
		}
		a.logger.Debugf("[populateParentsAndReasons] populating parent chain for block index %d, id=%v", i, b.ID)
		if err := a.populateParentChain(ctx, b, 0); err != nil {
			a.logger.Errorf("[populateParentsAndReasons] error populating parent chain for block index %d, id=%v: %v", i, b.ID, err)
			return err
		}
	}
	a.logger.Debugf("[populateParentsAndReasons] finished parent chains, attaching reasons...")
	err := a.attachReasons(ctx, blocks)
	if err != nil {
		a.logger.Errorf("[populateParentsAndReasons] error attaching reasons: %v", err)
		return err
	}
	a.logger.Debugf("[populateParentsAndReasons] completed successfully")
	return nil
}

// ─── Create methods ─────────────────────────────────────────────────────────

func (a *AccountBlockStorage) createBlock(ctx context.Context, block *imodel.AccountBlock) error {
	block.ID = uuid.New().String()
	_, err := a.db.ExecContext(ctx, insertAccountBlock,
		sql.Named("id", block.ID),
		sql.Named("name", block.Name),
		sql.Named("code", block.Code),
		sql.Named("address", block.Address),
		sql.Named("slug", block.Slug),
		sql.Named("type", string(block.Type)),
		sql.Named("district_id", nullStr(block.DistrictID)),
		sql.Named("region_id", nullStr(block.RegionID)),
	)
	if err != nil {
		return err
	}
	storage.BumpRedisCacheKey(ctx, a.redis, constants.RedisCacheKeyAccountBlock)
	return nil
}

func (a *AccountBlockStorage) CreateBranch(ctx context.Context, branch *imodel.AccountBlock) error {
	a.logger.Infof("[AccountBlockStorage][CreateBranch] creating branch: %s", branch.Name)
	branch.Type = imodel.TypeBranch
	if err := a.createBlock(ctx, branch); err != nil {
		a.logger.Errorf("[AccountBlockStorage][CreateBranch] failed: %v", err)
		return fmt.Errorf("failed to create branch")
	}
	a.logger.Infof("[AccountBlockStorage][CreateBranch] branch created: %s", branch.ID)
	return nil
}

func (a *AccountBlockStorage) CreateRegion(ctx context.Context, region *imodel.AccountBlock) error {
	a.logger.Infof("[AccountBlockStorage][CreateRegion] creating region: %s", region.Name)
	region.Type = imodel.TypeRegion
	if err := a.createBlock(ctx, region); err != nil {
		a.logger.Errorf("[AccountBlockStorage][CreateRegion] failed: %v", err)
		return fmt.Errorf("failed to create region")
	}
	a.logger.Infof("[AccountBlockStorage][CreateRegion] region created: %s", region.ID)
	return nil
}

func (a *AccountBlockStorage) CreateDistrict(ctx context.Context, district *imodel.AccountBlock) error {
	a.logger.Infof("[AccountBlockStorage][CreateDistrict] creating district: %s", district.Name)
	district.Type = imodel.TypeDistrict
	if err := a.createBlock(ctx, district); err != nil {
		a.logger.Errorf("[AccountBlockStorage][CreateDistrict] failed: %v", err)
		return fmt.Errorf("failed to create district")
	}
	a.logger.Infof("[AccountBlockStorage][CreateDistrict] district created: %s", district.ID)
	return nil
}

// func (a *AccountBlockStorage) CreateCity(ctx context.Context, city *imodel.AccountBlock) error {
// 	a.logger.Infof("[AccountBlockStorage][CreateCity] creating city: %s", city.Name)
// 	city.Type = imodel.TypeCity
// 	if err := a.createBlock(ctx, city); err != nil {
// 		a.logger.Errorf("[AccountBlockStorage][CreateCity] failed: %v", err)
// 		return fmt.Errorf("failed to create city")
// 	}
// 	a.logger.Infof("[AccountBlockStorage][CreateCity] city created: %s", city.ID)
// 	return nil
// }

// ─── Delete methods ─────────────────────────────────────────────────────────

func (a *AccountBlockStorage) deleteBlock(ctx context.Context, id string, entity string) error {
	a.logger.Infof("[AccountBlockStorage][Delete%s] deleting id: %s", entity, id)
	res, err := a.db.ExecContext(ctx, softDeleteAccountBlock,
		sql.Named("id", id),
	)
	if err != nil {
		a.logger.Errorf("[AccountBlockStorage][Delete%s] failed: %v", entity, err)
		return fmt.Errorf("failed to delete %s", entity)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("%s not found", entity)
	}
	storage.BumpRedisCacheKey(ctx, a.redis, constants.RedisCacheKeyAccountBlock)
	a.logger.Infof("[AccountBlockStorage][Delete%s] deleted: %s", entity, id)
	return nil
}

func (a *AccountBlockStorage) DeleteBranch(ctx context.Context, id string) error {
	return a.deleteBlock(ctx, id, "Branch")
}

func (a *AccountBlockStorage) DeleteRegion(ctx context.Context, id string) error {
	return a.deleteBlock(ctx, id, "Region")
}

func (a *AccountBlockStorage) DeleteDistrict(ctx context.Context, id string) error {
	return a.deleteBlock(ctx, id, "District")
}

// func (a *AccountBlockStorage) DeleteCity(ctx context.Context, id string) error {
// 	return a.deleteBlock(ctx, id, "City")
// }

// ─── Find by ID (with parent hierarchy) ─────────────────────────────────────

func (a *AccountBlockStorage) findByIDWithParents(ctx context.Context, id string) (*imodel.AccountBlock, error) {
	rows, err := a.db.QueryContext(ctx, selectWithAncestors,
		sql.Named("id", id),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blocks []*imodel.AccountBlock
	for rows.Next() {
		var ab imodel.AccountBlock
		var districtID, regionID sql.NullString
		var isEnabledInt, isDeletedInt, depth int
		var createdAt, updatedAt sql.NullTime

		err := rows.Scan(
			&ab.ID, &ab.Name, &ab.Code, &ab.Address,
			&ab.Slug, &ab.Type,
			&isEnabledInt, &districtID, &regionID,
			&isDeletedInt, &createdAt, &updatedAt,
			&depth,
		)
		if err != nil {
			return nil, err
		}

		ab.IsEnabled = isEnabledInt == 1
		ab.IsDeleted = isDeletedInt == 1
		if createdAt.Valid {
			ab.CreatedAt = createdAt.Time
		}
		if updatedAt.Valid {
			ab.UpdatedAt = updatedAt.Time
		}
		if districtID.Valid {
			ab.DistrictID = &districtID.String
		}
		if regionID.Valid {
			ab.RegionID = &regionID.String
		}

		blocks = append(blocks, &ab)
	}

	if len(blocks) == 0 {
		return nil, sql.ErrNoRows
	}

	// Build the parent chain: blocks[0] is the target, blocks[1] is its parent, etc.
	for i := len(blocks) - 1; i > 0; i-- {
		blocks[i-1].Parent = blocks[i]
	}

	if err := a.attachReasons(ctx, []*imodel.AccountBlock{blocks[0]}); err != nil {
		return nil, err
	}

	return blocks[0], nil
}

// ─── Find by filter key ─────────────────────────────────────────────────────

func (a *AccountBlockStorage) FindByFilterKey(ctx context.Context, field, value string) (*imodel.AccountBlock, error) {
	a.logger.Infof("[AccountBlockStorage][FindByFilterKey] field=%s value=%s", field, value)

	// Whitelist allowed filter keys to prevent SQL injection
	allowed := map[string]bool{"name": true, "code": true, "slug": true, "type": true}
	if !allowed[field] {
		return nil, errors.New("invalid filter key")
	}

	query := fmt.Sprintf(selectByFilterKey, field)
	row := a.db.QueryRowContext(ctx, query, sql.Named("val", value))
	ab, err := scanAccountBlockFromRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		a.logger.Errorf("[AccountBlockStorage][FindByFilterKey] failed: %v", err)
		return nil, err
	}
	if err := a.populateParentChain(ctx, ab, 0); err != nil {
		return nil, err
	}
	if err := a.attachReasons(ctx, []*imodel.AccountBlock{ab}); err != nil {
		return nil, err
	}
	return ab, nil
}

const selectByFilterKey = `SELECT
	RAWTOHEX(id), name, code, address, slug, type,
	is_enabled, RAWTOHEX(district_id), RAWTOHEX(region_id),
	is_deleted, created_at, updated_at
FROM ACCOUNT_BLOCKS
WHERE %s = :val AND is_deleted = 0`

func scanAccountBlockFromRow(row *sql.Row) (*imodel.AccountBlock, error) {
	var ab imodel.AccountBlock
	var districtID, regionID sql.NullString
	var isEnabledInt, isDeletedInt int
	var createdAt, updatedAt sql.NullTime

	err := row.Scan(
		&ab.ID, &ab.Name, &ab.Code, &ab.Address,
		&ab.Slug, &ab.Type,
		&isEnabledInt, &districtID, &regionID,
		&isDeletedInt, &createdAt, &updatedAt,
	)
	if err != nil {
		return nil, err
	}

	ab.IsEnabled = isEnabledInt == 1
	ab.IsDeleted = isDeletedInt == 1
	if createdAt.Valid {
		ab.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		ab.UpdatedAt = updatedAt.Time
	}
	if districtID.Valid {
		ab.DistrictID = &districtID.String
	}
	if regionID.Valid {
		ab.RegionID = &regionID.String
	}
	return &ab, nil
}

// ─── Get by IDs ─────────────────────────────────────────────────────────────

func (a *AccountBlockStorage) getByIds(ctx context.Context, ids []string, entityType imodel.AccountBlockType) ([]*imodel.AccountBlock, error) {
	a.logger.Infof("[getByIds] called with %d ids, entityType=%s", len(ids), entityType)
	if len(ids) == 0 {
		a.logger.Warnf("[getByIds] empty ids slice, returning nil")
		return nil, nil
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, 0, len(ids)+1)
	for i, id := range ids {
		paramName := fmt.Sprintf("id_%d", i)
		placeholders[i] = "HEXTORAW(:" + paramName + ")"
		args = append(args, sql.Named(paramName, id))
		a.logger.Debugf("[getByIds] param: %s = %s", paramName, id)
	}
	args = append(args, sql.Named("type", string(entityType)))

	query := fmt.Sprintf(`SELECT
			RAWTOHEX(id), name, code, address, slug, type,
			is_enabled, RAWTOHEX(district_id), RAWTOHEX(region_id),
			is_deleted, created_at, updated_at
		FROM ACCOUNT_BLOCKS
		WHERE id IN (%s) AND type = :type AND is_deleted = 0`, strings.Join(placeholders, ","))

	a.logger.Debugf("[getByIds] query: %s", query)
	a.logger.Debugf("[getByIds] args: %+v", args)

	rows, err := a.db.QueryContext(ctx, query, args...)
	if err != nil {
		a.logger.Errorf("[getByIds] QueryContext error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var results []*imodel.AccountBlock
	rowNum := 0
	for rows.Next() {
		ab, err := scanAccountBlock(rows)
		if err != nil {
			a.logger.Errorf("[getByIds] scanAccountBlock error at row %d: %v", rowNum, err)
			return nil, err
		}
		a.logger.Debugf("[getByIds] scanned row %d: %+v", rowNum, ab)
		results = append(results, ab)
		rowNum++
	}
	if err := rows.Err(); err != nil {
		a.logger.Errorf("[getByIds] rows.Err: %v", err)
		return nil, err
	}
	a.logger.Infof("[getByIds] fetched %d rows", len(results))
	if len(results) == 0 {
		a.logger.Warnf("[getByIds] no results found for ids: %v", ids)
		return results, nil
	}
	if err := a.populateParentsAndReasons(ctx, results); err != nil {
		a.logger.Errorf("[getByIds] populateParentsAndReasons error: %v", err)
		return nil, err
	}
	a.logger.Infof("[getByIds] returning %d results", len(results))
	return results, nil
}

func (a *AccountBlockStorage) GetBranchesByIds(ctx context.Context, ids []string) ([]*imodel.AccountBlock, error) {
	a.logger.Infof("[AccountBlockStorage][GetBranchesByIds] fetching %d branches", len(ids))
	return a.getByIds(ctx, ids, imodel.TypeBranch)
}

func (a *AccountBlockStorage) GetRegionsByIds(ctx context.Context, ids []string) ([]*imodel.AccountBlock, error) {
	a.logger.Infof("[AccountBlockStorage][GetRegionsByIds] fetching %d regions", len(ids))
	return a.getByIds(ctx, ids, imodel.TypeRegion)
}

func (a *AccountBlockStorage) GetDistrictsByIds(ctx context.Context, ids []string) ([]*imodel.AccountBlock, error) {
	a.logger.Infof("[AccountBlockStorage][GetDistrictsByIds] fetching %d districts", len(ids))
	return a.getByIds(ctx, ids, imodel.TypeDistrict)
}

// ─── Paginated list ─────────────────────────────────────────────────────────

// hexIDsFromFilterValue normalizes region_id / district_id values from ExtractFilterParams
// (comma-separated string, []interface{} from repeated keys or parseValue splits, []string, etc.).
func hexIDsFromFilterValue(v interface{}) []string {
	if v == nil {
		return nil
	}
	switch t := v.(type) {
	case string:
		var ids []string
		for _, part := range strings.Split(t, ",") {
			if s := strings.TrimSpace(part); s != "" {
				ids = append(ids, s)
			}
		}
		return ids
	case []string:
		ids := make([]string, 0, len(t))
		for _, s := range t {
			if s = strings.TrimSpace(s); s != "" {
				ids = append(ids, s)
			}
		}
		return ids
	case []interface{}:
		var ids []string
		for _, x := range t {
			ids = append(ids, hexIDsFromFilterValue(x)...)
		}
		return ids
	default:
		if s := strings.TrimSpace(fmt.Sprint(t)); s != "" {
			return []string{s}
		}
		return nil
	}
}

func buildHexIDMatchClause(column string, ids []string, paramPrefix string) (string, []interface{}) {
	if len(ids) == 0 {
		return "", nil
	}
	conditions := make([]string, len(ids))
	named := make([]interface{}, 0, len(ids))
	for i, id := range ids {
		name := fmt.Sprintf("%s_%d", paramPrefix, i)
		conditions[i] = fmt.Sprintf("%s = HEXTORAW(:%s)", column, name)
		named = append(named, sql.Named(name, id))
	}
	return "(" + strings.Join(conditions, " OR ") + ")", named
}

func (a *AccountBlockStorage) findAllWithPagination(ctx context.Context, filterParam types.Filter, entityType imodel.AccountBlockType) (*types.PaginatedResponse[[]*imodel.AccountBlock], error) {
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
	var regionIDs []string
	var districtIDs []string

	if filterParam.Search != "" {
		search = filterParam.Search
	}

	if filterParam.Filters != nil {
		if v, ok := filterParam.Filters["region_id"]; ok {
			regionIDs = hexIDsFromFilterValue(v)
		}
		if v, ok := filterParam.Filters["district_id"]; ok {
			districtIDs = hexIDsFromFilterValue(v)
		}
		if v, ok := filterParam.Filters["is_enabled"]; ok {
			if enabled, isBool := v.(bool); isBool {
				isEnabledFilter = isEnabledToInt(enabled)
			}
		}
	}

	var idClauses []string
	var args []interface{}
	if clause, clauseArgs := buildHexIDMatchClause("region_id", regionIDs, "region"); clause != "" {
		idClauses = append(idClauses, clause)
		args = append(args, clauseArgs...)
	}
	if clause, clauseArgs := buildHexIDMatchClause("district_id", districtIDs, "district"); clause != "" {
		idClauses = append(idClauses, clause)
		args = append(args, clauseArgs...)
	}

	idFilterSQL := ""
	if len(idClauses) > 0 {
		idFilterSQL = " AND " + strings.Join(idClauses, " AND ")
	}

	query := fmt.Sprintf(`SELECT
			RAWTOHEX(id) AS id,
			name,
			code,
			address,
			slug,
			type,
			is_enabled,
			RAWTOHEX(district_id) AS district_id,
			RAWTOHEX(region_id) AS region_id,
			is_deleted,
			created_at,
			updated_at,
			COUNT(*) OVER() AS total_count
		FROM ACCOUNT_BLOCKS
		WHERE type = :type
		  AND is_deleted = 0
		  %s
		  AND (:search IS NULL
		       OR LOWER(name) LIKE '%%' || LOWER(:search) || '%%'
		       OR LOWER(code) LIKE '%%' || LOWER(:search) || '%%'
		       OR LOWER(address) LIKE '%%' || LOWER(:search) || '%%')
		  AND (:is_enabled IS NULL OR is_enabled = :is_enabled)
		ORDER BY created_at DESC
		OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`, idFilterSQL)

	args = append(args,
		sql.Named("type", string(entityType)),
		sql.Named("search", search),
		sql.Named("is_enabled", isEnabledFilter),
		sql.Named("offset", offset),
		sql.Named("limit", limit),
	)

	rows, err := a.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*imodel.AccountBlock
	var totalCount int64

	for rows.Next() {
		var ab imodel.AccountBlock
		var districtID, regionID sql.NullString
		var isEnabledInt, isDeletedInt int
		var createdAt, updatedAt sql.NullTime

		err := rows.Scan(
			&ab.ID, &ab.Name, &ab.Code, &ab.Address,
			&ab.Slug, &ab.Type,
			&isEnabledInt, &districtID, &regionID,
			&isDeletedInt, &createdAt, &updatedAt,
			&totalCount,
		)
		if err != nil {
			return nil, err
		}

		ab.IsEnabled = isEnabledInt == 1
		ab.IsDeleted = isDeletedInt == 1
		if createdAt.Valid {
			ab.CreatedAt = createdAt.Time
		}
		if updatedAt.Valid {
			ab.UpdatedAt = updatedAt.Time
		}
		if districtID.Valid {
			ab.DistrictID = &districtID.String
		}
		if regionID.Valid {
			ab.RegionID = &regionID.String
		}
		results = append(results, &ab)
	}

	a.logger.Debugf("[findAllWithPagination] results count: %d", len(results))
	for i, ab := range results {
		if ab == nil {
			a.logger.Warnf("[findAllWithPagination] results[%d] is nil", i)
		} else {
			a.logger.Debugf("[findAllWithPagination] results[%d]: id=%v, name=%v, type=%v", i, ab.ID, ab.Name, ab.Type)
		}
	}

	if err := a.populateParentsAndReasons(ctx, results); err != nil {
		a.logger.Errorf("[AccountBlock][findAllWithPagination] error in populateParentsAndReasons: %v", err)
		return nil, err
	}

	a.logger.Debugf("[findAllWithPagination] populateParentsAndReasons completed")
	a.logger.Debugf("[findAllWithPagination] totalCount: %d, perPage: %d", totalCount, filterParam.PerPage)

	totalPages := int((totalCount + int64(filterParam.PerPage) - 1) / int64(filterParam.PerPage))
	a.logger.Debugf("[findAllWithPagination] totalPages: %d, currentPage: %d", totalPages, filterParam.Page)
	var prevPage, nextPage *int
	if filterParam.Page > 1 {
		p := filterParam.Page - 1
		prevPage = &p
		a.logger.Debugf("[findAllWithPagination] prevPage: %d", p)
	}
	if filterParam.Page < totalPages {
		n := filterParam.Page + 1
		nextPage = &n
		a.logger.Debugf("[findAllWithPagination] nextPage: %d", n)
	}

	resp := &types.PaginatedResponse[[]*imodel.AccountBlock]{
		Data: results,
		Meta: types.PaginationMeta{
			TotalDocs:     totalCount,
			Limit:         filterParam.PerPage,
			TotalPages:    totalPages,
			Page:          filterParam.Page,
			PagingCounter: offset + 1,
			HasPrevPage:   filterParam.Page > 1,
			HasNextPage:   filterParam.Page < totalPages,
			PrevPage:      prevPage,
			NextPage:      nextPage,
		},
	}
	a.logger.Debugf("[findAllWithPagination] returning paginated response: %+v", resp.Meta)
	return resp, nil
}

func (a *AccountBlockStorage) FindAllBranchesWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*imodel.AccountBlock], error) {
	a.logger.Infof("[AccountBlockStorage][FindAllBranchesWithPagination] fetching branches")
	return a.findAllWithPagination(ctx, filterParam, imodel.TypeBranch)
}

func (a *AccountBlockStorage) FindAllRegionsWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*imodel.AccountBlock], error) {
	a.logger.Infof("[AccountBlockStorage][FindAllRegionsWithPagination] fetching regions")
	return a.findAllWithPagination(ctx, filterParam, imodel.TypeRegion)
}

func (a *AccountBlockStorage) FindAllDistrictsWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*imodel.AccountBlock], error) {
	a.logger.Infof("[AccountBlockStorage][FindAllDistrictsWithPagination] fetching districts")
	return a.findAllWithPagination(ctx, filterParam, imodel.TypeDistrict)
}

// ─── Enable / Disable ───────────────────────────────────────────────────────

func (a *AccountBlockStorage) enableOrDisable(ctx context.Context, ids []string, reason *types.Reason, enabled bool, entityType imodel.AccountBlockType) error {
	if len(ids) == 0 {
		return nil
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, 0, len(ids)+5)
	for i, id := range ids {
		paramName := fmt.Sprintf("id_%d", i)
		placeholders[i] = "HEXTORAW(:" + paramName + ")"
		args = append(args, sql.Named(paramName, id))
	}

	args = append(args,
		sql.Named("is_enabled", isEnabledToInt(enabled)),
		sql.Named("type", string(entityType)),
	)

	if enabled {
		if err := a.deleteReasonsForBlockIDs(ctx, ids); err != nil {
			return err
		}
	} else {
		if err := a.deleteReasonsForBlockIDs(ctx, ids); err != nil {
			return err
		}
		if err := a.insertDisableReasonForBlocks(ctx, ids, reason); err != nil {
			return err
		}
	}

	query := fmt.Sprintf(`UPDATE ACCOUNT_BLOCKS
		SET is_enabled = :is_enabled,
		    updated_at = SYSTIMESTAMP
		WHERE id IN (%s) AND type = :type`, strings.Join(placeholders, ","))

	_, err := a.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	storage.BumpRedisCacheKey(ctx, a.redis, constants.RedisCacheKeyAccountBlock)
	return nil
}

func (a *AccountBlockStorage) EnableOrDisableBranches(ctx context.Context, ids []string, reason *types.Reason, enabled bool) error {
	a.logger.Infof("[AccountBlockStorage][EnableOrDisableBranches] ids=%v enabled=%v", ids, enabled)
	err := a.enableOrDisable(ctx, ids, reason, enabled, imodel.TypeBranch)
	if err != nil {
		a.logger.Errorf("[AccountBlockStorage][EnableOrDisableBranches] failed: %v", err)
	}
	return err
}

func (a *AccountBlockStorage) EnableOrDisableRegions(ctx context.Context, ids []string, reason *types.Reason, enabled bool) error {
	a.logger.Infof("[AccountBlockStorage][EnableOrDisableRegions] ids=%v enabled=%v", ids, enabled)
	err := a.enableOrDisable(ctx, ids, reason, enabled, imodel.TypeRegion)
	if err != nil {
		a.logger.Errorf("[AccountBlockStorage][EnableOrDisableRegions] failed: %v", err)
	}
	return err
}

func (a *AccountBlockStorage) EnableOrDisableDistricts(ctx context.Context, ids []string, reason *types.Reason, enabled bool) error {
	a.logger.Infof("[AccountBlockStorage][EnableOrDisableDistricts] ids=%v enabled=%v", ids, enabled)
	err := a.enableOrDisable(ctx, ids, reason, enabled, imodel.TypeDistrict)
	if err != nil {
		a.logger.Errorf("[AccountBlockStorage][EnableOrDisableDistricts] failed: %v", err)
	}
	return err
}

// func (a *AccountBlockStorage) EnableOrDisableCities(ctx context.Context, ids []string, reason *types.Reason, enabled bool) error {
// 	a.logger.Infof("[AccountBlockStorage][EnableOrDisableCities] ids=%v enabled=%v", ids, enabled)
// 	err := a.enableOrDisable(ctx, ids, reason, enabled, imodel.TypeCity)
// 	if err != nil {
// 		a.logger.Errorf("[AccountBlockStorage][EnableOrDisableCities] failed: %v", err)
// 	}
// 	return err
// }

// ─── GetAccountBlockDetails ─────────────────────────────────────────────────

func (a *AccountBlockStorage) GetAccountBlockDetails(ctx context.Context, id string, filterParam types.Filter) (*types.PaginatedResponse[[]account_block_dto.AccountBlockActionResponse], error) {
	a.logger.Infof("[AccountBlockStorage][GetAccountBlockDetails] fetching CPS actions for account block id: %s", id)

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
		a.logger.Errorf("[AccountBlockStorage][GetAccountBlockDetails] aggregation failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() { _ = cur.Close(ctx) }()

	var results []model.CPSAction
	if err := cur.All(ctx, &results); err != nil {
		a.logger.Errorf("[AccountBlockStorage][GetAccountBlockDetails] failed to decode CPS actions: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	if len(results) == 0 {
		a.logger.Errorf("[AccountBlockStorage][GetAccountBlockDetails] no CPS actions found for id: %s", id)
		return nil, errors.New(localization.ErrorActionNotFound.Code)
	}

	total, err := cpsCollection.CountDocuments(ctx, matchFilter)
	if err != nil {
		a.logger.Errorf("[GetAccountBlockDetails] failed to count CPS actions: %v", err)
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
			CheckerUsers:        convertCheckers(cpsAction.CheckerUsers),
			AuditorUsers:        convertAuditors(cpsAction.AuditorUsers),
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
	a.logger.Infof("[AccountBlockStorage][GetAccountBlockDetails] successfully mapped %d CPS actions", len(response))

	return &types.PaginatedResponse[[]account_block_dto.AccountBlockActionResponse]{
		Data: response,
		Meta: meta,
	}, nil
}

func getMatchingAction(actionData interface{}, accountBlockID string, logger utils.Logger) interface{} {
	if actionData == nil {
		return nil
	}

	// Try to unmarshal as array of EnableDisableAction
	actions, err := local_util.JsonUnmarshal[[]types.EnableDisableAction](actionData)
	if err == nil && actions != nil {
		// Find the action matching the account block id
		for _, action := range *actions {
			if action.ID == accountBlockID {
				return action
			}
		}
		// If no match found, return nil
		return nil
	}

	// Try to unmarshal as single EnableDisableAction
	singleAction, err := local_util.JsonUnmarshal[types.EnableDisableAction](actionData)
	if err == nil && singleAction != nil {
		if singleAction.ID == accountBlockID {
			return *singleAction
		}
		return nil
	}

	// If unmarshaling fails, check if it's already a map/object with id field
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

// convertCheckers converts model.Checker to account_block_dto.Checker
func convertCheckers(checkers []model.Checker) []account_block_dto.Checker {
	result := make([]account_block_dto.Checker, 0, len(checkers))
	for _, c := range checkers {
		result = append(result, account_block_dto.Checker{
			CheckerID:          c.CheckerID,
			RoleID:             c.RoleID,
			CheckerIndex:       c.CheckerIndex,
			CheckerName:        c.CheckerName,
			CheckerPhoneNumber: c.CheckerPhoneNumber,
			ApprovedAt:         c.ApprovedAt,
		})
	}
	return result
}

// // convertAuditors converts model.Auditor to account_block_dto.Auditor
func convertAuditors(auditors []model.Auditor) []account_block_dto.Auditor {
	result := make([]account_block_dto.Auditor, 0, len(auditors))
	for _, a := range auditors {
		result = append(result, account_block_dto.Auditor{
			AuditorID:          a.AuditorID,
			RoleID:             a.RoleID,
			AuditorIndex:       a.AuditorIndex,
			AuditorName:        a.AuditorName,
			AuditorPhoneNumber: a.AuditorPhoneNumber,
			AuditorReason:      a.AuditorReason,
			AuditorMark:        account_block_dto.AuditorMark(a.AuditorMark),
			ApprovedAt:         a.ApprovedAt,
		})
	}
	return result
}

// ─── GetAllBranches (by parent region/district/city id) ─────────────────────

func (a *AccountBlockStorage) GetAllBranches(ctx context.Context, id string) ([]imodel.AccountBlock, error) {

	query := `SELECT
		RAWTOHEX(id), name, code, address, slug, type,
		is_enabled, RAWTOHEX(district_id), RAWTOHEX(region_id),
		is_deleted, created_at, updated_at
	FROM ACCOUNT_BLOCKS
	WHERE type = 'B' AND is_deleted = 0 AND (
		district_id = HEXTORAW(:id) OR region_id = HEXTORAW(:id)
	)`

	rows, err := a.db.QueryContext(ctx, query, sql.Named("id", id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []imodel.AccountBlock
	for rows.Next() {
		ab, err := scanAccountBlock(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, *ab)
	}
	ptrs := make([]*imodel.AccountBlock, len(results))
	for i := range results {
		ptrs[i] = &results[i]
	}
	if err := a.populateParentsAndReasons(ctx, ptrs); err != nil {
		return nil, err
	}
	return results, nil
}

// ─── GetBranchByIds (single branch by ID, used by service internally) ───────

func (a *AccountBlockStorage) GetBranchByIds(ctx context.Context, id string) (*imodel.AccountBlock, error) {
	a.logger.Infof("[AccountBlockStorage][GetBranchByIds] fetching branch by id: %s", id)
	ab, err := a.findByIDWithParents(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			a.logger.Errorf("[AccountBlockStorage][GetBranchByIds] branch not found")
			return nil, errors.New(localization.ErrorBranchNotFound.Code)
		}
		a.logger.Errorf("[AccountBlockStorage][GetBranchByIds] failed to fetch branch: %v", err)
		return nil, err
	}
	a.logger.Infof("[AccountBlockStorage][GetBranchByIds] branch retrieved successfully")
	return ab, nil
}
