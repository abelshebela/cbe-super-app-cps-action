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

func oracleNumToBool(n sql.NullInt64) bool {
	return n.Valid && n.Int64 != 0
}

func scanWalletOracleCore(
	row interface {
		Scan(dest ...interface{}) error
	},
	wallet *model.WalletOracle,
	withService bool,
	serviceKey, serviceCode *sql.NullString,
) error {
	var en, del, sself, soth, sag sql.NullInt64
	var deletedAt sql.NullTime
	dest := []interface{}{
		&wallet.ID,
		&wallet.Name,
		&wallet.UniqueCode,
		&wallet.ServiceID,
		&en,
		&wallet.Avatar,
		&sself,
		&soth,
		&sag,
		&del,
		&wallet.CreatedAt,
		&wallet.LastModifiedAt,
		&deletedAt,
	}
	if withService {
		dest = append(dest, serviceKey, serviceCode)
	}
	if err := row.Scan(dest...); err != nil {
		return err
	}
	wallet.Enabled = oracleNumToBool(en)
	wallet.IsDeleted = oracleNumToBool(del)
	wallet.Self = oracleNumToBool(sself)
	wallet.Other = oracleNumToBool(soth)
	wallet.Agent = oracleNumToBool(sag)
	if deletedAt.Valid {
		t := deletedAt.Time
		wallet.DeletedAt = &t
	} else {
		wallet.DeletedAt = nil
	}
	return nil
}

// Create implements [storage.WalletOracleRepository].
// Table WALLETS: UNIQUE_CODE, SERVICES_SELF, SERVICES_OTHER, SERVICES_AGENT (see db/migrations).
func (q *WalletStorage) Create(ctx context.Context, wallet *model.WalletOracle) error {
	q.logger.Infof("[WalletStorage][Create] Creating wallet with unique_code: %s", wallet.UniqueCode)

	enabled, deleted := 0, 0
	if wallet.Enabled {
		enabled = 1
	}
	if wallet.IsDeleted {
		deleted = 1
	}
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

	query := `
			INSERT INTO wallets (
				id, name, unique_code, service_id, enabled, avatar,
				services_self, services_other, services_agent, is_deleted,
				created_at, last_modified_at, deleted_at
			) VALUES (
				SYS_GUID(), :1, :2, HEXTORAW(:3), :4, :5, :6, :7, :8, :9, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL
			)`

	_, err := q.db.ExecContext(ctx, query,
		wallet.Name,
		wallet.UniqueCode,
		wallet.ServiceID,
		enabled,
		wallet.Avatar,
		self,
		other,
		agent,
		deleted,
	)
	if err != nil {
		q.logger.Errorf("[WalletStorage][Create] failed to insert wallet: %v", err)
		return err
	}
	storage.BumpRedisCacheKey(ctx, q.redis, constants.RedisCacheKeyWallet)
	q.logger.Infof("[WalletStorage][Create] Successfully created wallet with unique_code: %s", wallet.UniqueCode)
	return nil
}

func (q *WalletStorage) Delete(ctx context.Context, id string) error {
	q.logger.Infof("[WalletStorage][Delete] Deleting wallet with ID: %s", id)
	_, err := q.db.ExecContext(ctx, "UPDATE wallets SET is_deleted=1, deleted_at=CURRENT_TIMESTAMP, last_modified_at=CURRENT_TIMESTAMP WHERE id=HEXTORAW(:1) AND is_deleted=0", id)
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

	_, err := q.db.ExecContext(ctx, "UPDATE wallets SET enabled=:1, last_modified_at=CURRENT_TIMESTAMP WHERE id=HEXTORAW(:2)", enabled, id)
	if err != nil {
		q.logger.Errorf("[WalletStorage][EnableOrDisable] failed: %v", err)
		return err
	}
	storage.BumpRedisCacheKey(ctx, q.redis, constants.RedisCacheKeyWallet)
	return nil
}

func (q *WalletStorage) Find(ctx context.Context, code string, name string, service_id string) (*model.WalletOracle, error) {
	q.logger.Infof("[WalletStorage][Find] Finding wallet with unique_code: %s and Name: %s", code, name)

	if code == "" && name == "" && service_id == "" {
		q.logger.Warnf("[WalletStorage][Find] unique_code, name, and service_id are empty")
		return nil, localization.ErrorUnexpectedError
	}

	query := `SELECT RAWTOHEX(id), name, unique_code, RAWTOHEX(service_id), enabled, avatar, services_self, services_other, services_agent, is_deleted, created_at, last_modified_at, deleted_at FROM wallets WHERE UPPER(unique_code)=:1 AND UPPER(name)=:2 AND RAWTOHEX(service_id)=:3 AND is_deleted=0`
	row := q.db.QueryRowContext(ctx, query, strings.ToUpper(code), strings.ToUpper(name), service_id)
	var wallet model.WalletOracle
	if err := scanWalletOracleCore(row, &wallet, false, nil, nil); err != nil {
		if err == sql.ErrNoRows {
			q.logger.Infof("[WalletStorage][Find] No wallet found with unique_code: %s and Name: %s", code, name)
			return nil, nil
		}
		q.logger.Errorf("[WalletStorage][Find] failed: %v", err)
		return nil, err
	}
	return &wallet, nil
}

func (q *WalletStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.WalletOracle], error) {
	return &types.PaginatedResponse[[]model.WalletOracle]{}, nil
}

func (q *WalletStorage) FindAllWithPaginationForGRPC(
	ctx context.Context,
	filterParam types.Filter,
) (*types.PaginatedResponse[[]model.WalletOracle], error) {

	q.logger.Infof("[WalletStorage][FindAllWithPaginationForGRPC] called with filter: %+v", filterParam)

	var filters []string
	var args []interface{}
	idx := 1

	// Base filter
	filters = append(filters, "w.is_deleted = 0")

	// Name filter
	if val, ok := filterParam.Filters["name"]; ok && val != nil {
		filters = append(filters, fmt.Sprintf("LOWER(w.name) LIKE :%d", idx))
		args = append(args, "%"+strings.ToLower(fmt.Sprint(val))+"%")
		idx++
	}

	// Code filter
	if val, ok := filterParam.Filters["code"]; ok && val != nil {
		filters = append(filters, fmt.Sprintf("LOWER(w.unique_code) LIKE :%d", idx))
		args = append(args, "%"+strings.ToLower(fmt.Sprint(val))+"%")
		idx++
	}

	// ✅ Optional enabled filter (FIXED)
	if val, ok := filterParam.Filters["enabled"]; ok && val != nil {
		_, ok := val.(bool)
		if ok {
			filters = append(filters, fmt.Sprintf("w.enabled = :%d", idx))
			args = append(args, val)
			idx++
		}

	}

	// Search filter
	if filterParam.Search != "" {
		search := "%" + strings.ToLower(filterParam.Search) + "%"
		filters = append(filters,
			fmt.Sprintf("(LOWER(w.name) LIKE :%d OR LOWER(w.unique_code) LIKE :%d)", idx, idx+1))
		args = append(args, search, search)
		idx += 2
	}

	// WHERE clause
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

	// Count query
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM wallets w %s`, whereClause)

	var total int
	err := q.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		q.logger.Errorf("[WalletStorage][FindAllWithPaginationForGRPC] count failed: %v", err)
		return nil, err
	}

	if total > 0 && offset >= total {
		meta := local_util.BuildPaginationMeta(int64(total), page, perPage)
		return &types.PaginatedResponse[[]model.WalletOracle]{
			Data: []model.WalletOracle{},
			Meta: meta,
		}, nil
	}

	// Main query (FIXED placeholder indexing)
	query := fmt.Sprintf(`
SELECT RAWTOHEX(w.id), w.name, w.unique_code, RAWTOHEX(w.service_id),
       w.enabled, w.avatar, w.services_self, w.services_other,
       w.services_agent, w.is_deleted, w.created_at,
       w.last_modified_at, w.deleted_at,
       sk.service_key, s.service_code
FROM WALLETS w
LEFT JOIN services s ON w.service_id = s.id
LEFT JOIN access_lists sk ON sk.id = s.access_list_id
%s
%s
OFFSET :%d ROWS FETCH NEXT :%d ROWS ONLY
`, whereClause, sortClause, idx, idx+1)

	// Append pagination args ONCE
	args = append(args, offset, perPage)

	// Debug logs
	q.logger.Infof("Final Query: %s", query)
	q.logger.Infof("Args: %+v", args)

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

		if err := scanWalletOracleCore(rows, &wallet, true, &serviceKey, &serviceCode); err != nil {
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

	if err := rows.Err(); err != nil {
		return nil, err
	}

	meta := local_util.BuildPaginationMeta(int64(total), page, perPage)

	return &types.PaginatedResponse[[]model.WalletOracle]{
		Data: wallets,
		Meta: meta,
	}, nil
}

func (q *WalletStorage) FindByID(ctx context.Context, id string) (*model.WalletOracle, error) {
	q.logger.Infof("[WalletStorage][FindByID] Finding wallet with ID: %s", id)
	row := q.db.QueryRowContext(ctx, `SELECT RAWTOHEX(id), name, unique_code, RAWTOHEX(service_id), enabled, avatar, services_self, services_other, services_agent, is_deleted, created_at, last_modified_at, deleted_at FROM wallets WHERE id=HEXTORAW(:1)`, id)
	var wallet model.WalletOracle
	if err := scanWalletOracleCore(row, &wallet, false, nil, nil); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		q.logger.Errorf("[WalletStorage][FindByID] failed: %v", err)
		return nil, err
	}
	return &wallet, nil
}

func (q *WalletStorage) FindByIDForGRPC(ctx context.Context, id string) (*model.WalletOracle, error) {
	q.logger.Infof("[WalletStorage][FindByIDForGRPC] Finding wallet with ID: %s", id)
	row := q.db.QueryRowContext(ctx, `
			SELECT RAWTOHEX(w.id), w.name, w.unique_code, RAWTOHEX(w.service_id), w.enabled, w.avatar, w.services_self, w.services_other, w.services_agent, w.is_deleted, w.created_at, w.last_modified_at, w.deleted_at,
				   sk.service_key, s.service_code
			FROM wallets w
			LEFT JOIN services s ON w.service_id = s.id
			LEFT JOIN access_lists sk ON sk.id = s.access_list_id
			WHERE w.id = HEXTORAW(:1)`, id)
	var wallet model.WalletOracle
	var serviceKey, serviceCode sql.NullString
	err := scanWalletOracleCore(row, &wallet, true, &serviceKey, &serviceCode)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		q.logger.Errorf("[WalletStorage][FindByIDForGRPC] failed: %v", err)
		return nil, err
	}
	if serviceKey.Valid {
		wallet.ServiceKey = serviceKey.String
	}
	if serviceCode.Valid {
		wallet.ServiceCode = serviceCode.String
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
				services_self = :5,
				services_other = :6,
				services_agent = :7,
				last_modified_at = CURRENT_TIMESTAMP
			WHERE id = HEXTORAW(:8) `
	_, err := q.db.ExecContext(ctx, query,
		wallet.Name,
		wallet.UniqueCode,
		wallet.ServiceID,
		wallet.Avatar,
		self,
		other,
		agent,
		id,
	)
	if err != nil {
		q.logger.Errorf("[WalletStorage][Update] failed: %v", err)
		return err
	}
	storage.BumpRedisCacheKey(ctx, q.redis, constants.RedisCacheKeyWallet)
	return nil
}
