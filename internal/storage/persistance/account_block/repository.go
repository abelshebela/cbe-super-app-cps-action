package account_block

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"

	account_block_dto "cbe-super-app-cps-action/internal/constants/dto/account_block"

	"github.com/google/uuid"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// ─── SQL constants ───────────────────────────────────────────────────────────

const (
	insertAccountBlock = `INSERT INTO ACCOUNT_BLOCKS (
		id, name, code, address, parent_id, slug, type,
		is_enabled, city_id, district_id, region_id,
		is_deleted, created_at, updated_at
	) VALUES (
		:id, :name, :code, :address, :parent_id, :slug, :type,
		1, :city_id, :district_id, :region_id,
		0, SYSTIMESTAMP, SYSTIMESTAMP
	)`

	softDeleteAccountBlock = `UPDATE ACCOUNT_BLOCKS
		SET is_deleted = 1, updated_at = SYSTIMESTAMP
		WHERE id = :id AND is_deleted = 0`

	selectAccountBlockByID = `SELECT
		id, name, code, address, parent_id, slug, type,
		is_enabled, city_id, district_id, region_id,
		is_deleted, created_at, updated_at
	FROM ACCOUNT_BLOCKS
	WHERE id = :id AND is_deleted = 0`

	selectWithAncestors = `SELECT
		id, name, code, address, parent_id, slug, type,
		is_enabled, city_id, district_id, region_id,
		is_deleted, created_at, updated_at, LEVEL as depth
	FROM ACCOUNT_BLOCKS
	WHERE is_deleted = 0
	START WITH id = :id
	CONNECT BY PRIOR parent_id = id
	ORDER BY LEVEL ASC`

	listAccountBlocksByType = `SELECT
		id, name, code, address, parent_id, slug, type,
		is_enabled, city_id, district_id, region_id,
		is_deleted, created_at, updated_at,
		COUNT(*) OVER() AS total_count
	FROM ACCOUNT_BLOCKS
	WHERE type = :type
	  AND is_deleted = 0
	  AND (:search IS NULL
	       OR LOWER(name) LIKE '%%' || LOWER(:search) || '%%'
	       OR LOWER(code) LIKE '%%' || LOWER(:search) || '%%'
	       OR LOWER(address) LIKE '%%' || LOWER(:search) || '%%')
	  AND (:region_id IS NULL OR region_id = :region_id)
	  AND (:district_id IS NULL OR district_id = :district_id)
	  AND (:city_id IS NULL OR city_id = :city_id)
	  AND (:is_enabled IS NULL OR is_enabled = :is_enabled)
	ORDER BY created_at DESC
	OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`
)

// ─── Repository struct ──────────────────────────────────────────────────────

type AccountBlockStorage struct {
	db     *sql.DB
	redis  storage.RedisRepository
	logger utils.Logger
}

func NewAccountBlockRepository(

	db *sql.DB,
	redis storage.RedisRepository,
	logger utils.Logger,
) storage.AccountBlockRepository {
	// One-time table/constraint creation with existence checks
	// Helper: check if table exists
	tableExists := func(name string) bool {
		var cnt int
		err := db.QueryRow("SELECT COUNT(*) FROM USER_TABLES WHERE TABLE_NAME = :1", strings.ToUpper(name)).Scan(&cnt)
		return err == nil && cnt > 0
	}
	// Helper: check if constraint exists
	constraintExists := func(name string) bool {
		var cnt int
		err := db.QueryRow("SELECT COUNT(*) FROM USER_CONSTRAINTS WHERE CONSTRAINT_NAME = :1", strings.ToUpper(name)).Scan(&cnt)
		return err == nil && cnt > 0
	}
	// ACCESS_LIST_CUSTOMER_SEG
	if !tableExists("ACCESS_LIST_CUSTOMER_SEG") {
		stmt := `CREATE TABLE ACCESS_LIST_CUSTOMER_SEG (
			   ID RAW(16) DEFAULT SYS_GUID() PRIMARY KEY,
			   ACCESS_LIST_KEY RAW(16),
			   SEGMENTED_ID RAW(16),
			   ENABLED NUMBER(1) DEFAULT 1,
			   CREATED_AT TIMESTAMP,
			   UPDATED_AT TIMESTAMP,
			   DELETED_AT TIMESTAMP,
			   CONSTRAINT FK_AL_CUSTOMER_ACCESS_LIST FOREIGN KEY (ACCESS_LIST_KEY) REFERENCES ACCESS_LIST (id) ON DELETE CASCADE,
			   CONSTRAINT FK_AL_CUSTOMER_SEGMENT FOREIGN KEY (SEGMENTED_ID) REFERENCES CUSTOMER_SEGMENTATIONS (ID) ON DELETE CASCADE
		   )`
		if _, err := db.Exec(stmt); err != nil {
			logger.Warnf("ACCESS_LIST_CUSTOMER_SEG table creation: %v", err)
		} else {
			logger.Infof("ACCESS_LIST_CUSTOMER_SEG table created successfully")
		}
	}

	// ACCESS_LIST_GEO_SEG
	if !tableExists("ACCESS_LIST_GEO_SEG") {
		stmt := `CREATE TABLE ACCESS_LIST_GEO_SEG (
			   ID RAW(16) DEFAULT SYS_GUID() PRIMARY KEY,
			   ACCESS_LIST_KEY RAW(16),
			   SEGMENTED_ID VARCHAR2(32),
			   TYPE VARCHAR2(1),
			   ENABLED NUMBER(1) DEFAULT 1,
			   CREATED_AT TIMESTAMP,
			   UPDATED_AT TIMESTAMP,
			   DELETED_AT TIMESTAMP,
			   CONSTRAINT FK_AL_GEO_ACCESS_LIST FOREIGN KEY (ACCESS_LIST_KEY) REFERENCES ACCESS_LIST (ID) ON DELETE CASCADE,
			   CONSTRAINT FK_AL_GEO_BLOCK FOREIGN KEY (SEGMENTED_ID) REFERENCES ACCOUNT_BLOCKS (ID) ON DELETE CASCADE
		   )`
		if _, err := db.Exec(stmt); err != nil {
			logger.Warnf("ACCESS_LIST_GEO_SEG table creation: %v", err)
		} else {
			logger.Infof("ACCESS_LIST_GEO_SEG table created successfully")
		}
	}

	// ACCOUNT_BLOCKS table
	if !tableExists("ACCOUNT_BLOCKS") {
		stmt := `CREATE TABLE ACCOUNT_BLOCKS (
			   ID RAW(24) PRIMARY KEY,
			   NAME VARCHAR2(255),
			   CODE VARCHAR2(100) UNIQUE,
			   ADDRESS VARCHAR2(500),
			   PARENT_ID RAW(24),
			   SLUG VARCHAR2(255),
			   TYPE VARCHAR2(1),
			   IS_ENABLED NUMBER(1) DEFAULT 1,
			   CITY_ID RAW(24),
			   REGION_ID RAW(24),
			   DISTRICT_ID RAW(24),
			   IS_DELETED NUMBER(1) DEFAULT 0,
			   CREATED_AT TIMESTAMP,
			   UPDATED_AT TIMESTAMP,
			   CONSTRAINT FK_ACCOUNT_BLOCK_PARENT FOREIGN KEY (CITY_ID) REFERENCES ACCOUNT_BLOCKS (ID) ON DELETE CASCADE,
			   CONSTRAINT FK_ACCOUNT_BLOCK_PARENT2 FOREIGN KEY (REGION_ID) REFERENCES ACCOUNT_BLOCKS (ID) ON DELETE CASCADE,
			   CONSTRAINT FK_ACCOUNT_BLOCK_PARENT3 FOREIGN KEY (DISTRICT_ID) REFERENCES ACCOUNT_BLOCKS (ID) ON DELETE CASCADE
		   )`
		if _, err := db.Exec(stmt); err != nil {
			logger.Errorf("ACCOUNT_BLOCKS table creation failed (API will return ORA-00942 until migration or DDL succeeds): %v", err)
		} else {
			logger.Infof("ACCOUNT_BLOCKS table created successfully")
		}
	}
	// Add type check constraint if not exists
	if !constraintExists("CHK_AB_TYPE") {
		stmt := `ALTER TABLE ACCOUNT_BLOCKS ADD CONSTRAINT CHK_AB_TYPE CHECK (TYPE IN ('R', 'D', 'C', 'B'))`
		if _, err := db.Exec(stmt); err != nil {
			logger.Errorf("CHK_AB_TYPE constraint on ACCOUNT_BLOCKS failed: %v", err)
		} else {
			logger.Infof("CHK_AB_TYPE constraint added successfully")
		}
	}

	// CUSTOMER_GROUPS
	if !tableExists("CUSTOMER_GROUPS") {
		stmt := `CREATE TABLE CUSTOMER_GROUPS (
			   ID RAW(16) DEFAULT SYS_GUID() PRIMARY KEY,
			   NAME VARCHAR2(32) NOT NULL,
			   IS_ENABLED NUMBER(1) DEFAULT 1,
			   IS_DELETED NUMBER(1) DEFAULT 0,
			   CREATED_AT TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			   LAST_MODIFIED_AT TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			   DELETED_AT TIMESTAMP
		   )`
		if _, err := db.Exec(stmt); err != nil {
			logger.Warnf("CUSTOMER_GROUPS table creation: %v", err)
		} else {
			logger.Infof("CUSTOMER_GROUPS table created successfully")
		}
	}
	// SUPERAPP_ROLE
	if !tableExists("SUPERAPP_ROLE") {
		stmt := `CREATE TABLE SUPERAPP_ROLE (
			   ID RAW(16) DEFAULT SYS_GUID() PRIMARY KEY,
			   NAME VARCHAR2(32) NOT NULL,
			   ROLE_CODE VARCHAR2(32) NOT NULL UNIQUE,
			   DESCRIPTION VARCHAR2(128),
			   ENABLED NUMBER(1) DEFAULT 1,
			   IS_DELETED NUMBER(1) DEFAULT 0,
			   CREATED_AT TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			   LAST_MODIFIED_AT TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			   DELETED_AT TIMESTAMP
		   )`
		if _, err := db.Exec(stmt); err != nil {
			logger.Warnf("SUPERAPP_ROLE table creation: %v", err)
		} else {
			logger.Infof("SUPERAPP_ROLE table created successfully")
		}
	}
	// CUSTOMER_SEGMENTATIONS
	if !tableExists("CUSTOMER_SEGMENTATIONS") {
		stmt := `CREATE TABLE CUSTOMER_SEGMENTATIONS (
			   ID RAW(16) DEFAULT SYS_GUID() PRIMARY KEY,
			   NAME VARCHAR2(32) NOT NULL,
			   CUSTOMER_GROUPS_ID RAW(16) NOT NULL,
			   IS_ENABLED NUMBER(1) DEFAULT 1,
			   IS_DELETED NUMBER(1) DEFAULT 0,
			   CREATED_AT TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			   LAST_MODIFIED_AT TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			   DELETED_AT TIMESTAMP,
			   CONSTRAINT FK_CUSTOMER_GROUPS FOREIGN KEY (CUSTOMER_GROUPS_ID) REFERENCES CUSTOMER_GROUPS (ID) ON DELETE CASCADE
		   )`
		if _, err := db.Exec(stmt); err != nil {
			logger.Warnf("CUSTOMER_SEGMENTATIONS table creation: %v", err)
		} else {
			logger.Infof("CUSTOMER_SEGMENTATIONS table created successfully")
		}
	}
	// CUSTOMER_SUB_SEGMENTS
	if !tableExists("CUSTOMER_SUB_SEGMENTS") {
		stmt := `CREATE TABLE CUSTOMER_SUB_SEGMENTS (
			   ID RAW(16) DEFAULT SYS_GUID() PRIMARY KEY,
			   NAME VARCHAR2(32) NOT NULL,
			   CUSTOMER_SEGMENTATIONS_ID RAW(16) NOT NULL,
			   SUPERAPP_ROLE_ID RAW(16) NOT NULL,
			   IS_ENABLED NUMBER(1) DEFAULT 1,
			   IS_DELETED NUMBER(1) DEFAULT 0,
			   CREATED_AT TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			   LAST_MODIFIED_AT TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			   DELETED_AT TIMESTAMP,
			   CONSTRAINT FK_CUSTOMER_SEGMENTATIONS FOREIGN KEY (CUSTOMER_SEGMENTATIONS_ID) REFERENCES CUSTOMER_SEGMENTATIONS (ID) ON DELETE CASCADE,
			   CONSTRAINT FK_SUB_SEG_SUPERAPP_ROLE FOREIGN KEY (SUPERAPP_ROLE_ID) REFERENCES SUPERAPP_ROLE (ID)
		   )`
		if _, err := db.Exec(stmt); err != nil {
			logger.Warnf("CUSTOMER_SUB_SEGMENTS table creation: %v", err)
		} else {
			logger.Infof("CUSTOMER_SUB_SEGMENTS table created successfully")
		}
	}
	return &AccountBlockStorage{
		db:     db,
		redis:  redis,
		logger: logger,
	}
}

// ─── Scan helpers ───────────────────────────────────────────────────────────

func scanAccountBlock(scanner interface{ Scan(dest ...any) error }) (*imodel.AccountBlock, error) {
	var ab imodel.AccountBlock
	var parentID, cityID, districtID, regionID sql.NullString
	var isEnabledInt, isDeletedInt int
	var createdAt, updatedAt sql.NullTime

	err := scanner.Scan(
		&ab.ID, &ab.Name, &ab.Code, &ab.Address,
		&parentID, &ab.Slug, &ab.Type,
		&isEnabledInt, &cityID, &districtID, &regionID,
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
	if parentID.Valid {
		ab.ParentID = &parentID.String
	}
	if cityID.Valid {
		ab.CityID = &cityID.String
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

// fetchBlockByID loads one row by primary key (any type), for parent_id resolution.
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
	case imodel.TypeCity:
		// City usually belongs to a district; fallback to region if needed.
		if id := trimmed(b.DistrictID); id != "" {
			return id
		}
		return trimmed(b.RegionID)
	case imodel.TypeBranch:
		// Branch should resolve to district first, then city, then region.
		if id := trimmed(b.DistrictID); id != "" {
			return id
		}
		if id := trimmed(b.CityID); id != "" {
			return id
		}
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
		sql.Named("parent_id", nullStr(block.ParentID)),
		sql.Named("slug", block.Slug),
		sql.Named("type", string(block.Type)),
		sql.Named("city_id", nullStr(block.CityID)),
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

func (a *AccountBlockStorage) CreateCity(ctx context.Context, city *imodel.AccountBlock) error {
	a.logger.Infof("[AccountBlockStorage][CreateCity] creating city: %s", city.Name)
	city.Type = imodel.TypeCity
	if err := a.createBlock(ctx, city); err != nil {
		a.logger.Errorf("[AccountBlockStorage][CreateCity] failed: %v", err)
		return fmt.Errorf("failed to create city")
	}
	a.logger.Infof("[AccountBlockStorage][CreateCity] city created: %s", city.ID)
	return nil
}

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

func (a *AccountBlockStorage) DeleteCity(ctx context.Context, id string) error {
	return a.deleteBlock(ctx, id, "City")
}

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
		var parentID, cityID, districtID, regionID sql.NullString
		var isEnabledInt, isDeletedInt, depth int
		var createdAt, updatedAt sql.NullTime

		err := rows.Scan(
			&ab.ID, &ab.Name, &ab.Code, &ab.Address,
			&parentID, &ab.Slug, &ab.Type,
			&isEnabledInt, &cityID, &districtID, &regionID,
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
		if parentID.Valid {
			ab.ParentID = &parentID.String
		}
		if cityID.Valid {
			ab.CityID = &cityID.String
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
			return nil, nil
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
	id, name, code, address, parent_id, slug, type,
	is_enabled, city_id, district_id, region_id,
	is_deleted, created_at, updated_at
FROM ACCOUNT_BLOCKS
WHERE %s = :val AND is_deleted = 0`

func scanAccountBlockFromRow(row *sql.Row) (*imodel.AccountBlock, error) {
	var ab imodel.AccountBlock
	var parentID, cityID, districtID, regionID sql.NullString
	var isEnabledInt, isDeletedInt int
	var createdAt, updatedAt sql.NullTime

	err := row.Scan(
		&ab.ID, &ab.Name, &ab.Code, &ab.Address,
		&parentID, &ab.Slug, &ab.Type,
		&isEnabledInt, &cityID, &districtID, &regionID,
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
	if parentID.Valid {
		ab.ParentID = &parentID.String
	}
	if cityID.Valid {
		ab.CityID = &cityID.String
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
		placeholders[i] = ":" + paramName
		args = append(args, sql.Named(paramName, id))
		a.logger.Debugf("[getByIds] param: %s = %s", paramName, id)
	}
	args = append(args, sql.Named("type", string(entityType)))

	query := fmt.Sprintf(`SELECT
			id, name, code, address, parent_id, slug, type,
			is_enabled, city_id, district_id, region_id,
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

func (a *AccountBlockStorage) GetCitiesByIds(ctx context.Context, ids []string) ([]*imodel.AccountBlock, error) {
	a.logger.Infof("[AccountBlockStorage][GetCitiesByIds] fetching %d cities", len(ids))
	return a.getByIds(ctx, ids, imodel.TypeCity)
}

// ─── Paginated list ─────────────────────────────────────────────────────────

func (a *AccountBlockStorage) findAllWithPagination(ctx context.Context, filterParam types.Filter, entityType imodel.AccountBlockType) (*types.PaginatedResponse[[]*imodel.AccountBlock], error) {
	offset := (filterParam.Page - 1) * filterParam.PerPage
	if offset < 0 {
		offset = 0
	}
	limit := filterParam.PerPage
	if limit <= 0 {
		limit = 50
	}

	var search, regionIDFilter, districtIDFilter, cityIDFilter interface{}
	var isEnabledFilter interface{}

	if filterParam.Search != "" {
		search = filterParam.Search
	}

	// Extract filters from map
	if filterParam.Filters != nil {
		if v, ok := filterParam.Filters["region_id"]; ok {
			regionIDFilter = nullIfEmptyFilter(v)
		}
		if v, ok := filterParam.Filters["district_id"]; ok {
			districtIDFilter = nullIfEmptyFilter(v)
		}
		if v, ok := filterParam.Filters["city_id"]; ok {
			cityIDFilter = nullIfEmptyFilter(v)
		}
		if v, ok := filterParam.Filters["is_enabled"]; ok {
			if enabled, isBool := v.(bool); isBool {
				isEnabledFilter = isEnabledToInt(enabled)
			}
		}
	}

	rows, err := a.db.QueryContext(ctx, listAccountBlocksByType,
		sql.Named("type", string(entityType)),
		sql.Named("search", search),
		sql.Named("region_id", regionIDFilter),
		sql.Named("district_id", districtIDFilter),
		sql.Named("city_id", cityIDFilter),
		sql.Named("is_enabled", isEnabledFilter),
		sql.Named("offset", offset),
		sql.Named("limit", limit),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*imodel.AccountBlock
	var totalCount int64

	for rows.Next() {
		var ab imodel.AccountBlock
		var parentID, cityID, districtID, regionID sql.NullString
		var isEnabledInt, isDeletedInt int
		var createdAt, updatedAt sql.NullTime

		err := rows.Scan(
			&ab.ID, &ab.Name, &ab.Code, &ab.Address,
			&parentID, &ab.Slug, &ab.Type,
			&isEnabledInt, &cityID, &districtID, &regionID,
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
		if parentID.Valid {
			ab.ParentID = &parentID.String
		}
		if cityID.Valid {
			ab.CityID = &cityID.String
		}
		if districtID.Valid {
			ab.DistrictID = &districtID.String
		}
		if regionID.Valid {
			ab.RegionID = &regionID.String
		}
		results = append(results, &ab)
	}

	if err := a.populateParentsAndReasons(ctx, results); err != nil {
		a.logger.Errorf("[AccountBlock][findAllWithPagination] error-----: %v", err)
		return nil, err
	}

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

	return &types.PaginatedResponse[[]*imodel.AccountBlock]{
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
	}, nil
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

func (a *AccountBlockStorage) FindAllCitiesWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*imodel.AccountBlock], error) {
	a.logger.Infof("[AccountBlockStorage][FindAllCitiesWithPagination] fetching cities")
	return a.findAllWithPagination(ctx, filterParam, imodel.TypeCity)
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
		placeholders[i] = ":" + paramName
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

func (a *AccountBlockStorage) EnableOrDisableCities(ctx context.Context, ids []string, reason *types.Reason, enabled bool) error {
	a.logger.Infof("[AccountBlockStorage][EnableOrDisableCities] ids=%v enabled=%v", ids, enabled)
	err := a.enableOrDisable(ctx, ids, reason, enabled, imodel.TypeCity)
	if err != nil {
		a.logger.Errorf("[AccountBlockStorage][EnableOrDisableCities] failed: %v", err)
	}
	return err
}

// ─── GetAccountBlockDetails ─────────────────────────────────────────────────

func (a *AccountBlockStorage) GetAccountBlockDetails(ctx context.Context, id string, filterParam types.Filter) (*types.PaginatedResponse[[]account_block_dto.AccountBlockActionResponse], error) {
	a.logger.Infof("[AccountBlockStorage][GetAccountBlockDetails] id=%s", id)

	// Fetch the account block by ID
	row := a.db.QueryRowContext(ctx, selectAccountBlockByID,
		sql.Named("id", id),
	)
	_, err := scanAccountBlockFromRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("account block not found")
		}
		return nil, err
	}

	// For now, return empty paginated response as CPS action details
	// are typically fetched via separate CPS action queries
	return &types.PaginatedResponse[[]account_block_dto.AccountBlockActionResponse]{
		Data: []account_block_dto.AccountBlockActionResponse{},
		Meta: types.PaginationMeta{
			TotalDocs:  0,
			Limit:      filterParam.PerPage,
			TotalPages: 0,
			Page:       filterParam.Page,
		},
	}, nil
}

// ─── GetAllBranches (by parent region/district/city id) ─────────────────────

func (a *AccountBlockStorage) GetAllBranches(ctx context.Context, id string) ([]imodel.AccountBlock, error) {
	a.logger.Infof("[AccountBlockStorage][GetAllBranches] parent_id=%s", id)

	query := `SELECT
		id, name, code, address, parent_id, slug, type,
		is_enabled, city_id, district_id, region_id,
		is_deleted, created_at, updated_at
	FROM ACCOUNT_BLOCKS
	WHERE type = 'B' AND is_deleted = 0 AND (
		city_id = :id OR district_id = :id OR region_id = :id
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
