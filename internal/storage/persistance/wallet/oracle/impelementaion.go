package wallet_oracle

import (
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type WalletStorage struct {
	db     DBTX
	redis  storage.RedisRepository
	logger utils.Logger
}

func NewWalletOracleRepository(db *sql.DB, redis storage.RedisRepository, log utils.Logger) storage.WalletOracleRepository {
	return &WalletStorage{
		db:     db,
		redis:  redis,
		logger: log,
	}
}

// Create implements [storage.WalletOracleRepository].
func (q *WalletStorage) Create(ctx context.Context, wallet *model.WalletOracle) error {
	q.logger.Infof("[WalletStorage][Create] Creating wallet with UniqueCode: %s", wallet.UniqueCode)

	enabled, deleted := 0, 0
	if wallet.Enabled {
		enabled = 1
	}
	if wallet.IsDeleted {
		deleted = 1
	}
	query := `
			INSERT INTO wallets (
				id, name, unique_code, service_id, avatar, enabled, is_deleted, created_at, last_modified_at, deleted_at
			) VALUES (
				SYS_GUID(), :1, :2, HEXTORAW(:3), :4, :5, :6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, :7
			)`

	_, err := q.db.ExecContext(ctx, query,
		wallet.Name,
		wallet.UniqueCode,
		wallet.ServiceID, // should be hex string
		wallet.Avatar,
		enabled,
		deleted,
		wallet.DeletedAt,
	)
	if err != nil {
		q.logger.Errorf("[WalletStorage][Create] failed to insert wallet: %v", err)
		return err
	}
	storage.BumpRedisCacheKey(ctx, q.redis, constants.RedisCacheKeyWallet)
	q.logger.Infof("[WalletStorage][Create] Successfully created wallet with UniqueCode: %s", wallet.UniqueCode)
	return nil
}

func (q *WalletStorage) Delete(ctx context.Context, id string) error {
	q.logger.Infof("[WalletStorage][Delete] Deleting wallet with ID: %s", id)
	_, err := q.db.ExecContext(ctx, "UPDATE wallets SET is_deleted=1, deleted_at=CURRENT_TIMESTAMP WHERE id=:1 AND is_deleted=0", id)
	if err != nil {
		q.logger.Errorf("[WalletStorage][Delete] failed to delete wallet: %v", err)
		return err
	}
	storage.BumpRedisCacheKey(ctx, q.redis, constants.RedisCacheKeyWallet)
	return nil
}

func (q *WalletStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	q.logger.Infof("[WalletStorage][EnableOrDisable] Setting enabled=%v for wallet ID: %s", enable, id)
	var enabled int
	if enable {
		enabled = 1
	} else {
		enabled = 0
	}

	_, err := q.db.ExecContext(ctx, "UPDATE wallets SET enabled=:1, last_modified_at=CURRENT_TIMESTAMP WHERE id=:2", enabled, id)
	if err != nil {
		q.logger.Errorf("[WalletStorage][EnableOrDisable] failed: %v", err)
		return err
	}
	storage.BumpRedisCacheKey(ctx, q.redis, constants.RedisCacheKeyWallet)
	return nil
}

func (q *WalletStorage) Find(ctx context.Context, code string, name string) (*model.WalletOracle, error) {
	q.logger.Infof("[WalletStorage][Find] Finding wallet with UniqueCode: %s and Name: %s", code, name)

	if code == "" && name == "" {
		q.logger.Warnf("[WalletStorage][Find] unique_code and name are empty")
		return nil, localization.ErrorUnexpectedError
	}

	query := fmt.Sprintf("SELECT RAWTOHEX(id), name, unique_code, RAWTOHEX(service_id), avatar, enabled, is_deleted, created_at, last_modified_at, deleted_at FROM wallets WHERE UPPER(unique_code)=:1 AND UPPER(name)=:2 AND is_deleted=0")
	row := q.db.QueryRowContext(ctx, query, strings.ToUpper(code), strings.ToUpper(name))
	var wallet model.WalletOracle
	err := row.Scan(&wallet.ID,
		&wallet.Name,
		&wallet.UniqueCode,
		&wallet.ServiceID,
		&wallet.Avatar,
		&wallet.Enabled,
		&wallet.IsDeleted,
		&wallet.CreatedAt,
		&wallet.LastModifiedAt,
		&wallet.DeletedAt)
	if err == sql.ErrNoRows {
		q.logger.Infof("[WalletStorage][Find] No wallet found with UniqueCode: %s and Name: %s", code, name)
		return nil, nil
	} else if err != nil {
		q.logger.Errorf("[WalletStorage][Find] failed: %v", err)
		return nil, err
	}
	return &wallet, nil
}

func (q *WalletStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.WalletOracle], error) {
	// TODO: Implement paginated query
	// Build WHERE clause from filterParam
	// rows, err := q.db.QueryContext(ctx, "SELECT ... FROM wallets WHERE ... OFFSET :1 ROWS FETCH NEXT :2 ROWS ONLY", ...)
	// if err != nil {
	// 	q.logger.Errorf("[WalletStorage][FindAllWithPagination] failed: %v", err)
	// 	return nil, err
	// }
	// defer rows.Close()
	// var wallets []model.WalletOracle
	// for rows.Next() {
	// 	var wallet model.WalletOracle
	// 	if err := rows.Scan(...); err != nil {
	// 		q.logger.Errorf("[WalletStorage][FindAllWithPagination] scan failed: %v", err)
	// 		continue
	// 	}
	// 	wallets = append(wallets, wallet)
	// }
	// meta := ... // BuildPaginationMeta equivalent
	// return &types.PaginatedResponse[[]model.WalletOracle]{Data: wallets, Meta: meta}, nil
	return &types.PaginatedResponse[[]model.WalletOracle]{}, nil
}

func (q *WalletStorage) FindAllWithPaginationForGRPC(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.WalletOracle], error) {
	q.logger.Infof("[WalletStorage][FindAllWithPaginationForGRPC] called with filter: %+v", filterParam)

	// Filtering
	var filters []string
	var args []interface{}
	idx := 1
	filters = append(filters, "w.is_deleted = 0")
	if filterParam.Filters["name"] != nil {
		filters = append(filters, fmt.Sprintf("LOWER(w.name) LIKE :%d", idx))
		args = append(args, "%"+strings.ToLower(fmt.Sprint(filterParam.Filters["name"]))+"%")
		idx++
	}
	if filterParam.Filters["code"] != nil {
		filters = append(filters, fmt.Sprintf("LOWER(w.unique_code) LIKE :%d", idx))
		args = append(args, "%"+strings.ToLower(fmt.Sprint(filterParam.Filters["code"]))+"%")
		idx++
	}
	if filterParam.Filters["enabled"] != nil {
		filters = append(filters, fmt.Sprintf("w.enabled = :%d", idx))
		args = append(args, filterParam.Filters["enabled"])
		idx++
	}
	if filterParam.Search != "" && filterParam.Search != "enabled" {
		filters = append(filters, fmt.Sprintf("(LOWER(w.name) LIKE :%d OR LOWER(w.unique_code) LIKE :%d)", idx, idx+1))
		args = append(args, "%"+strings.ToLower(filterParam.Search)+"%", "%"+strings.ToLower(filterParam.Search)+"%")
		idx += 2
	}
	if filterParam.Search == "enabled" {
		filters = append(filters, "w.enabled = 1")
	}
	whereClause := ""
	if len(filters) > 0 {
		whereClause = "WHERE " + strings.Join(filters, " AND ")
	}

	// Sorting
	sortClause := "ORDER BY w.created_at DESC"
	if val, ok := filterParam.Filters["sort_name"]; ok {
		if s, ok := val.(string); ok {
			switch strings.ToUpper(s) {
			case "ASC":
				sortClause = "ORDER BY w.name ASC"
			case "DESC":
				sortClause = "ORDER BY w.name DESC"
			}
		}
	}

	// Pagination
	page := filterParam.Page
	perPage := filterParam.PerPage
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}
	offset := (page - 1) * perPage
	// Total count
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM wallets w %s`, whereClause)
	var total int
	err := q.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		q.logger.Errorf("[WalletStorage][FindAllWithPaginationForGRPC] count failed: %v", err)
		return nil, err
	}

	// Out-of-bounds validation
	if offset+perPage >= total && total > 0 {
		// meta := local_util.BuildPaginationMeta(int64(total), filterParam.Page, filterParam.PerPage)
		// return &types.PaginatedResponse[[]model.WalletOracle]{
		// 	Data: []model.WalletOracle{},
		// 	Meta: meta,
		// }, nil
		offset = 0
		perPage = 10
	}

	// Query
	// query := fmt.Sprintf(`
	// 		SELECT RAWTOHEX(w.id), w.name, w.unique_code, RAWTOHEX(w.service_id), w.avatar, w.enabled, w.is_deleted, w.created_at, w.last_modified_at, w.deleted_at,
	// 			   s.service_key, s.service_code
	// 		FROM wallets w
	// 		LEFT JOIN services s ON w.service_id = s.id
	// 		%s
	// 		%s
	// 		OFFSET %d ROWS FETCH NEXT %d ROWS ONLY
	// 	`, whereClause, sortClause, offset, perPage)

	query := fmt.Sprintf(`
			SELECT RAWTOHEX(w.id), w.name, w.unique_code, RAWTOHEX(w.service_id), w.avatar, w.enabled, w.is_deleted, w.created_at, w.last_modified_at, w.deleted_at
			FROM wallets w
			LEFT JOIN services s ON w.service_id = s.id
			%s
			%s
			OFFSET %d ROWS FETCH NEXT %d ROWS ONLY
		`, whereClause, sortClause, offset, perPage)

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		q.logger.Errorf("[WalletStorage][FindAllWithPaginationForGRPC] query failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	var wallets []model.WalletOracle
	for rows.Next() {
		var wallet model.WalletOracle
		var serviceKey, serviceCode sql.NullString
		err := rows.Scan(
			&wallet.ID,
			&wallet.Name,
			&wallet.UniqueCode,
			&wallet.ServiceID,
			&wallet.Avatar,
			&wallet.Enabled,
			&wallet.IsDeleted,
			&wallet.CreatedAt,
			&wallet.LastModifiedAt,
			&wallet.DeletedAt,
			&serviceKey,
			&serviceCode,
		)
		if err != nil {
			q.logger.Warnf("[WalletStorage][FindAllWithPaginationForGRPC] scan failed: %v", err)
			continue
		}
		if serviceKey.Valid {
			wallet.ServiceKey = serviceKey.String
		}
		if serviceCode.Valid {
			wallet.ServiceCode = serviceCode.String
		}
		wallets = append(wallets, wallet)
	}

	meta := local_util.BuildPaginationMeta(int64(total), filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]model.WalletOracle]{
		Data: wallets,
		Meta: meta,
	}, nil
}

func (q *WalletStorage) FindByID(ctx context.Context, id string) (*model.WalletOracle, error) {
	q.logger.Infof("[WalletStorage][FindByID] Finding wallet with ID: %s", id)
	row := q.db.QueryRowContext(ctx, "SELECT RAWTOHEX(id), name, unique_code, RAWTOHEX(service_id), avatar, enabled, is_deleted, created_at, last_modified_at, deleted_at FROM wallets WHERE id=HEXTORAW(:1)", id)
	var wallet model.WalletOracle
	err := row.Scan(&wallet.ID,
		&wallet.Name,
		&wallet.UniqueCode,
		&wallet.ServiceID,
		&wallet.Avatar,
		&wallet.Enabled,
		&wallet.IsDeleted,
		&wallet.CreatedAt,
		&wallet.LastModifiedAt,
		&wallet.DeletedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		q.logger.Errorf("[WalletStorage][FindByID] failed: %v", err)
		return nil, err
	}
	return &wallet, nil
}

func (q *WalletStorage) FindByIDForGRPC(ctx context.Context, id string) (*model.WalletOracle, error) {
	q.logger.Infof("[WalletStorage][FindByIDForGRPC] Finding wallet with ID: %s", id)
	row := q.db.QueryRowContext(ctx, `
			SELECT RAWTOHEX(w.id), w.name, w.unique_code, RAWTOHEX(w.service_id), w.avatar, w.enabled, w.is_deleted, w.created_at, w.last_modified_at, w.deleted_at,
				   s.service_key, s.service_code
			FROM wallets w
			LEFT JOIN services s ON w.service_id = s.id
			WHERE w.id = HEXTORAW(:1)`, id)
	var wallet model.WalletOracle
	err := row.Scan(
		&wallet.ID,
		&wallet.Name,
		&wallet.UniqueCode,
		&wallet.ServiceID,
		&wallet.Avatar,
		&wallet.Enabled,
		&wallet.IsDeleted,
		&wallet.CreatedAt,
		&wallet.LastModifiedAt,
		&wallet.DeletedAt,
		&wallet.ServiceKey,
		&wallet.ServiceCode,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		q.logger.Errorf("[WalletStorage][FindByIDForGRPC] failed: %v", err)
		return nil, err
	}
	return &wallet, nil
}

func (q *WalletStorage) Update(ctx context.Context, id string, wallet *model.WalletOracle) error {
	self, other, agent := 0, 0, 0
	if wallet.Self {
		self = 1
	}
	if wallet.Other {
		other = 1
	}
	if wallet.Agent {
		agent = 1
	}
	q.logger.Infof("[WalletStorage][Update] Updating wallet with ID: %s", id)
	wallet.LastModifiedAt = time.Now()
	query := `
			UPDATE wallets SET
				name = :1,
				unique_code = :2,
				service_id = HEXTORAW(:3),
				avatar = :4,
				self = :5,
				other = :6,
				agent = :9,
				last_modified_at = :10
			WHERE id = HEXTORAW(:11) `
	_, err := q.db.ExecContext(ctx, query,
		wallet.Name,
		wallet.UniqueCode,
		wallet.ServiceID, // hex string
		wallet.Avatar,
		self,
		other,
		agent,
		wallet.LastModifiedAt,
		id, // hex string
	)
	if err != nil {
		q.logger.Errorf("[WalletStorage][Update] failed: %v", err)
		return err
	}
	storage.BumpRedisCacheKey(ctx, q.redis, constants.RedisCacheKeyWallet)
	return nil
}
