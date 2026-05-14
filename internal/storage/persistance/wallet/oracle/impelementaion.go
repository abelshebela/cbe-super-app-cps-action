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
	db     *sql.DB
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

// Create implements [storage.WalletOracleRepository].
// Table WALLETS: UNIQUE_CODE, SERVICES_SELF, SERVICES_OTHER, SERVICES_AGENT (see db/migrations).
func (q *WalletStorage) Create(ctx context.Context, wallet *model.WalletOracle) error {
	q.logger.Infof("[WalletStorage][Create] Creating wallet with unique_code: %s", wallet.UniqueCode)

	enabled := 0
	if wallet.Enabled {
		enabled = 1
	}
	deleted := 0
	if wallet.IsDeleted {
		deleted = 1
	}

	// Insert wallet and fetch generated ID using godror's sql.Out
	walletInsert := `
		INSERT INTO wallets (
			id, wallet_name, unique_code, logo, is_enabled, is_deleted, created_at, last_modified_at, deleted_at
		) VALUES (
			SYS_GUID(), :1, :2, :3, :4, :5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL
		) RETURNING RAWTOHEX(id) INTO :6`

	var walletID string
	_, err := q.db.ExecContext(ctx, walletInsert,
		wallet.Name,
		wallet.UniqueCode,
		wallet.Avatar, // logo
		enabled,
		deleted,
		sql.Out{Dest: &walletID},
	)
	if err != nil {
		q.logger.Errorf("[WalletStorage][Create] failed to insert wallet: %v", err)
		return localization.ErrorUnexpectedError
	}

	// Insert wallet_services for each type if enabled and service id exists
	serviceTypes := []struct {
		flag      bool
		stype     string
		serviceID string
	}{
		{wallet.Self, "SELF", wallet.SelfServiceID},
		{wallet.Other, "OTHER", wallet.OtherServiceID},
		{wallet.Agent, "AGENT", wallet.AgentServiceID},
	}

	wsInsert := `INSERT INTO wallet_services (
		id, service_id, wallet_id, service_type, is_enabled, is_deleted, created_at, last_modified_at
	) VALUES (
		SYS_GUID(), HEXTORAW(:1), HEXTORAW(:2), :3, 1, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
	)`

	for _, st := range serviceTypes {
		if st.flag && st.serviceID != "" {
			_, err := q.db.ExecContext(ctx, wsInsert, st.serviceID, walletID, st.stype)
			if err != nil {
				q.logger.Errorf("[WalletStorage][Create] failed to insert wallet_service type %s: %v", st.stype, err)
				return localization.ErrorUnexpectedError
			}
			q.logger.Infof("[WalletStorage][Create] Inserted wallet_service type %s for wallet_id: %s", st.stype, walletID)
		}
	}

	storage.BumpRedisCacheKey(ctx, q.redis, constants.RedisCacheKeyWallet)
	q.logger.Infof("[WalletStorage][Create] Successfully created wallet and wallet_services with unique_code: %s", wallet.UniqueCode)
	return nil
}

func (q *WalletStorage) Delete(ctx context.Context, id string) error {
	q.logger.Infof("[WalletStorage][Delete] Deleting wallet with ID: %s", id)
	_, err := q.db.ExecContext(ctx, "UPDATE WALLETS SET is_deleted=1, deleted_at=SYSTIMESTAMP, last_modified_at=SYSTIMESTAMP WHERE id=HEXTORAW(:1) AND is_deleted=0", id)
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

	_, err := q.db.ExecContext(ctx, "UPDATE wallets SET is_enabled=:1, last_modified_at=CURRENT_TIMESTAMP WHERE id=HEXTORAW(:2)", enabled, id)
	if err != nil {
		q.logger.Errorf("[WalletStorage][EnableOrDisable] failed: %v", err)
		return err
	}
	storage.BumpRedisCacheKey(ctx, q.redis, constants.RedisCacheKeyWallet)
	return nil
}

// Find returns a non-deleted wallet matching any provided criterion (OR).
// Pass only the fields you want to check: e.g. ("", name, "") for name uniqueness, or (code, "", "") for code.
// Previously this used AND across all three, so duplicates on name or code alone were never detected.
func (q *WalletStorage) Find(ctx context.Context, code string, name string, service_id string) (*model.WalletOracle, error) {
	code = strings.TrimSpace(code)
	name = strings.TrimSpace(name)
	service_id = strings.TrimSpace(service_id)

	q.logger.Infof("[WalletStorage][Find] code=%q name=%q service_id=%q", code, name, service_id)

	if code == "" && name == "" && service_id == "" {
		q.logger.Warnf("[WalletStorage][Find] unique_code, name, and service_id are empty")
		return nil, localization.ErrorUnexpectedError
	}

	var parts []string
	var args []interface{}
	idx := 1
	if code != "" {
		parts = append(parts, fmt.Sprintf("UPPER(TRIM(w.unique_code)) = UPPER(TRIM(:%d))", idx))
		args = append(args, code)
		idx++
	}
	if name != "" {
		parts = append(parts, fmt.Sprintf("UPPER(TRIM(w.wallet_name)) = UPPER(TRIM(:%d))", idx))
		args = append(args, name)
		idx++
	}
	if service_id != "" {
		parts = append(parts, fmt.Sprintf("RAWTOHEX(ws.service_id) = UPPER(:%d)", idx))
		args = append(args, strings.ToUpper(service_id))
		idx++
	}

	condition := strings.Join(parts, " OR ")
	if condition == "" {
		condition = "1=1"
	}

	query := fmt.Sprintf(`
SELECT 
  RAWTOHEX(w.id), w.wallet_name, w.unique_code, w.is_enabled, w.logo, w.is_deleted, w.created_at, w.last_modified_at, w.deleted_at
FROM wallets w
WHERE w.is_deleted = 0 AND (%s)
FETCH FIRST 1 ROW ONLY
`, condition)

	row := q.db.QueryRowContext(ctx, query, args...)
	var (
		wallet model.WalletOracle
	)
	err := row.Scan(
		&wallet.ID,
		&wallet.Name,
		&wallet.UniqueCode,
		&wallet.Enabled,
		&wallet.Avatar,
		&wallet.IsDeleted,
		&wallet.CreatedAt,
		&wallet.LastModifiedAt,
		&wallet.DeletedAt,
	)
	if err == sql.ErrNoRows {
		q.logger.Infof("[WalletStorage][Find] no wallet matched for code=%q name=%q service_id=%q", code, name, service_id)
		return nil, nil
	}
	if err != nil {
		q.logger.Errorf("[WalletStorage][Find] failed: %v", err)
		return nil, err
	}

	return &wallet, nil
}

func (q *WalletStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.WalletOracle], error) {
	// Use the same logic as FindAllWithPaginationForGRPC
	return q.FindAllWithPaginationForGRPC(context.Background(), filterParam)
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
		filters = append(filters, fmt.Sprintf("LOWER(w.wallet_name) LIKE :%d", idx))
		args = append(args, "%"+strings.ToLower(fmt.Sprint(val))+"%")
		idx++
	}

	// Code filter
	if val, ok := filterParam.Filters["code"]; ok && val != nil {
		filters = append(filters, fmt.Sprintf("LOWER(w.unique_code) LIKE :%d", idx))
		args = append(args, "%"+strings.ToLower(fmt.Sprint(val))+"%")
		idx++
	}

	// Optional enabled filter — Oracle column is NUMBER; godror rejects Go bool binds (ORA-00932).
	if val, ok := filterParam.Filters["enabled"]; ok && val != nil {
		var on bool
		var apply bool
		switch b := val.(type) {
		case bool:
			on = b
			apply = true
		case string:
			s := strings.TrimSpace(b)
			if s == "" {
				break
			}
			on = strings.EqualFold(s, "true") || s == "1"
			apply = true
		case float64:
			on = b != 0
			apply = true
		}
		if apply {
			n := 0
			if on {
				n = 1
			}
			filters = append(filters, fmt.Sprintf("w.is_enabled = :%d", idx))
			args = append(args, n)
			idx++
		}
	}

	// Search filter
	if filterParam.Search != "" {
		search := "%" + strings.ToLower(filterParam.Search) + "%"
		filters = append(filters,
			fmt.Sprintf("(LOWER(w.wallet_name) LIKE :%d OR LOWER(w.unique_code) LIKE :%d)", idx, idx+1))
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
				sortClause = "ORDER BY w.wallet_name ASC"
			case "DESC":
				sortClause = "ORDER BY w.wallet_name DESC"
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

	// Main query: fetch wallet, and for each wallet, fetch service IDs and codes for SELF, OTHER, AGENT
	query := fmt.Sprintf(`
SELECT 
  RAWTOHEX(w.id), w.wallet_name, w.unique_code, w.is_enabled, w.logo, w.is_deleted, w.created_at, w.last_modified_at, w.deleted_at
FROM wallets w %s %s
OFFSET :%d ROWS FETCH NEXT :%d ROWS ONLY`, whereClause, sortClause, idx, idx+1)

	// Append pagination args ONCE
	args = append(args, offset, perPage)

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
		var (
			wallet model.WalletOracle
		)
		err := rows.Scan(
			&wallet.ID,
			&wallet.Name,
			&wallet.UniqueCode,
			&wallet.Enabled,
			&wallet.Avatar,
			&wallet.IsDeleted,
			&wallet.CreatedAt,
			&wallet.LastModifiedAt,
			&wallet.DeletedAt,
		)
		if err != nil {
			q.logger.Warnf("[WalletStorage][FindAllWithPaginationForGRPC] scan failed: %v", err)
			continue
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
	query := `
SELECT 
  RAWTOHEX(w.id), w.wallet_name, w.unique_code, w.is_enabled, w.logo, w.is_deleted, w.created_at, w.last_modified_at, w.deleted_at,
  MAX(CASE WHEN ws.service_type = 'SELF' THEN RAWTOHEX(ws.service_id) END) AS self_service_id,
  MAX(CASE WHEN ws.service_type = 'OTHER' THEN RAWTOHEX(ws.service_id) END) AS other_service_id,
  MAX(CASE WHEN ws.service_type = 'AGENT' THEN RAWTOHEX(ws.service_id) END) AS agent_service_id,

  MAX(CASE WHEN ws.service_type = 'SELF' THEN s.service_code END) AS self_service_code,
  MAX(CASE WHEN ws.service_type = 'OTHER' THEN s.service_code END) AS other_service_code,
  MAX(CASE WHEN ws.service_type = 'AGENT' THEN s.service_code END) AS agent_service_code
FROM wallets w

LEFT JOIN wallet_services ws ON w.id = ws.wallet_id AND ws.is_deleted = 0 AND ws.is_enabled = 1
LEFT JOIN services s ON ws.service_id = s.id
WHERE w.id = HEXTORAW(:1)
GROUP BY w.id, w.wallet_name, w.unique_code, w.is_enabled, w.logo, w.is_deleted, w.created_at, w.last_modified_at, w.deleted_at`

	row := q.db.QueryRowContext(ctx, query, id)
	var (
		wallet                                              model.WalletOracle
		selfServiceID, otherServiceID, agentServiceID       sql.NullString
		selfServiceCode, otherServiceCode, agentServiceCode sql.NullString
	)
	err := row.Scan(
		&wallet.ID,
		&wallet.Name,
		&wallet.UniqueCode,
		&wallet.Enabled,
		&wallet.Avatar,
		&wallet.IsDeleted,
		&wallet.CreatedAt,
		&wallet.LastModifiedAt,
		&wallet.DeletedAt,
		&selfServiceID,
		&otherServiceID,
		&agentServiceID,
		&selfServiceCode,
		&otherServiceCode,
		&agentServiceCode,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		q.logger.Errorf("[WalletStorage][FindByID] failed: %v", err)
		return nil, err
	}

	wallet.SelfServiceID = ""
	if selfServiceID.Valid {
		wallet.SelfServiceID = selfServiceID.String
		if wallet.SelfServiceID != "" {
			wallet.Self = true
		}
	}
	wallet.OtherServiceID = ""
	if otherServiceID.Valid {
		wallet.OtherServiceID = otherServiceID.String
		if wallet.OtherServiceID != "" {
			wallet.Other = true
		}
	}
	wallet.AgentServiceID = ""
	if agentServiceID.Valid {
		wallet.AgentServiceID = agentServiceID.String
		if wallet.AgentServiceID != "" {
			wallet.Agent = true
		}
	}
	wallet.SelfServiceCode = ""
	if selfServiceCode.Valid {
		wallet.SelfServiceCode = selfServiceCode.String
	}
	wallet.OtherServiceCode = ""
	if otherServiceCode.Valid {
		wallet.OtherServiceCode = otherServiceCode.String
	}
	wallet.AgentServiceCode = ""
	if agentServiceCode.Valid {
		wallet.AgentServiceCode = agentServiceCode.String
	}

	return &wallet, nil
}
func (q *WalletStorage) FindWalletServiceByID(ctx context.Context, id string) (*model.WalletService, error) {
	q.logger.Infof("[WalletStorage][FindWalletServiceByID] Finding wallet service with ID: %s", id)
	query := `
SELECT 
  RAWTOHEX(ws.id),ws.service_id, RAWTOHEX(ws.wallet_id), ws.service_type, ws.is_enabled, ws.is_deleted, ws.created_at, ws.last_modified_at, ws.deleted_at
FROM WALLET_SERVICES ws
WHERE ws.id = HEXTORAW(:1)`

	row := q.db.QueryRowContext(ctx, query, id)
	var (
		wallet model.WalletService
	)
	err := row.Scan(
		&wallet.ID,
		&wallet.ServiceID,
		&wallet.WalletID,
		&wallet.ServiceType,
		&wallet.IsEnabled,
		&wallet.IsDeleted,
		&wallet.CreatedAt,
		&wallet.LastModifiedAt,
		&wallet.DeletedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		q.logger.Errorf("[WalletStorage][FindByID] failed: %v", err)
		return nil, err
	}

	return &wallet, nil
}

func (q *WalletStorage) FindByIDForGRPC(ctx context.Context, id string) (*model.WalletOracle, error) {
	q.logger.Infof("[WalletStorage][FindByIDForGRPC] Finding wallet with ID: %s", id)
	query := `
SELECT 
  RAWTOHEX(w.id), w.wallet_name, w.unique_code, w.is_enabled, w.logo, w.is_deleted, w.created_at, w.last_modified_at, w.deleted_at,
  MAX(CASE WHEN ws.service_type = 'SELF' THEN RAWTOHEX(ws.service_id) END) AS self_service_id,
  MAX(CASE WHEN ws.service_type = 'OTHER' THEN RAWTOHEX(ws.service_id) END) AS other_service_id,
  MAX(CASE WHEN ws.service_type = 'AGENT' THEN RAWTOHEX(ws.service_id) END) AS agent_service_id,

  MAX(CASE WHEN ws.service_type = 'SELF' THEN s.service_code END) AS self_service_code,
  MAX(CASE WHEN ws.service_type = 'OTHER' THEN s.service_code END) AS other_service_code,
  MAX(CASE WHEN ws.service_type = 'AGENT' THEN s.service_code END) AS agent_service_code,

  MAX(CASE WHEN ws.service_type = 'SELF' THEN ws.is_enabled END) AS self_service_enabled,
  MAX(CASE WHEN ws.service_type = 'OTHER' THEN ws.is_enabled END) AS other_service_enabled,
  MAX(CASE WHEN ws.service_type = 'AGENT' THEN ws.is_enabled END) AS agent_service_enabled
FROM wallets w
LEFT JOIN wallet_services ws ON w.id = ws.wallet_id AND ws.is_deleted = 0 AND ws.is_enabled = 1
LEFT JOIN services s ON ws.service_id = s.id
WHERE w.id = HEXTORAW(:1) AND w.is_deleted = 0
GROUP BY w.id, w.wallet_name, w.unique_code, w.is_enabled, w.logo, w.is_deleted, w.created_at, w.last_modified_at, w.deleted_at`

	row := q.db.QueryRowContext(ctx, query, id)
	var (
		wallet                                                       model.WalletOracle
		selfServiceID, otherServiceID, agentServiceID                sql.NullString
		selfServiceCode, otherServiceCode, agentServiceCode          sql.NullString
		selfServiceEnabled, otherServiceEnabled, agentServiceEnabled sql.NullInt16
	)
	err := row.Scan(
		&wallet.ID,
		&wallet.Name,
		&wallet.UniqueCode,
		&wallet.Enabled,
		&wallet.Avatar,
		&wallet.IsDeleted,
		&wallet.CreatedAt,
		&wallet.LastModifiedAt,
		&wallet.DeletedAt,
		&selfServiceID,
		&otherServiceID,
		&agentServiceID,
		&selfServiceCode,
		&otherServiceCode,
		&agentServiceCode,
		&selfServiceEnabled,
		&otherServiceEnabled,
		&agentServiceEnabled,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		q.logger.Errorf("[WalletStorage][FindByIDForGRPC] failed: %v", err)
		return nil, err
	}

	wallet.SelfServiceID = ""
	if selfServiceID.Valid {
		wallet.SelfServiceID = selfServiceID.String
		if wallet.SelfServiceID != "" {
			wallet.Self = true
		}
	}
	wallet.OtherServiceID = ""
	if otherServiceID.Valid {
		wallet.OtherServiceID = otherServiceID.String
		if wallet.OtherServiceID != "" {
			wallet.Other = true
		}
	}
	wallet.AgentServiceID = ""
	if agentServiceID.Valid {
		wallet.AgentServiceID = agentServiceID.String
		if wallet.AgentServiceID != "" {
			wallet.Agent = true
		}
	}
	wallet.SelfServiceCode = ""
	if selfServiceCode.Valid {
		wallet.SelfServiceCode = selfServiceCode.String
	}
	wallet.OtherServiceCode = ""
	if otherServiceCode.Valid {
		wallet.OtherServiceCode = otherServiceCode.String
	}
	wallet.AgentServiceCode = ""
	if agentServiceCode.Valid {
		wallet.AgentServiceCode = agentServiceCode.String
	}
	wallet.SelfServiceEnabled = 0
	if selfServiceEnabled.Valid {
		wallet.SelfServiceEnabled = int(selfServiceEnabled.Int16)
	}
	wallet.OtherServiceEnabled = 0
	if otherServiceEnabled.Valid {
		wallet.OtherServiceEnabled = int(otherServiceEnabled.Int16)
	}
	wallet.AgentServiceEnabled = 0
	if agentServiceEnabled.Valid {
		wallet.AgentServiceEnabled = int(agentServiceEnabled.Int16)
	}

	return &wallet, nil
}

func (q *WalletStorage) Update(ctx context.Context, id string, wallet *model.WalletOracle) error {
	q.logger.Infof("[WalletStorage][Update] Updating wallet with ID: %s", id)
	wallet.LastModifiedAt = time.Now()

	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		q.logger.Errorf("[WalletStorage][Update] failed to begin transaction: %v", err)
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	query := `
		UPDATE wallets SET
			wallet_name = :1,
			unique_code = :2,
			logo = :3,
			last_modified_at = CURRENT_TIMESTAMP
		WHERE id = HEXTORAW(:7)`
	_, err = tx.ExecContext(ctx, query,
		wallet.Name,
		wallet.UniqueCode,
		wallet.Avatar,
		id,
	)
	if err != nil {
		q.logger.Errorf("[WalletStorage][Update] failed to update wallet: %v", err)
		return err
	}

	// Helper for wallet_services update/insert/delete
	type serviceUpdate struct {
		flag      bool
		stype     string
		serviceID string
	}
	serviceTypes := []serviceUpdate{
		{wallet.Self, "SELF", wallet.SelfServiceID},
		{wallet.Other, "OTHER", wallet.OtherServiceID},
		{wallet.Agent, "AGENT", wallet.AgentServiceID},
	}

	// Always soft-delete existing wallet_services of this type for this wallet
	_, err = tx.ExecContext(ctx, `DELETE FROM wallet_services ws WHERE ws.wallet_id=HEXTORAW(:1)`, id)
	if err != nil {
		q.logger.Errorf("[WalletStorage][Update] failed to delete wallet_service for wallet %s: %v", id, err)
		return err
	}
	for _, st := range serviceTypes {
		if st.flag && st.serviceID != "" {
			// Insert new wallet_service for this type
			_, err = tx.ExecContext(ctx, `INSERT INTO wallet_services (
				id, service_id, wallet_id, service_type, is_enabled, is_deleted, created_at, last_modified_at, deleted_at
			) VALUES (
				SYS_GUID(), HEXTORAW(:1), HEXTORAW(:2), :3, 1, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL
			)`, st.serviceID, id, st.stype)
			if err != nil {
				q.logger.Errorf("[WalletStorage][Update] failed to insert wallet_service type %s: %v", st.stype, err)
				return err
			}
			q.logger.Infof("[WalletStorage][Update] Inserted wallet_service type %s for wallet_id: %s", st.stype, id)
		}
	}

	storage.BumpRedisCacheKey(ctx, q.redis, constants.RedisCacheKeyWallet)
	return nil
}

// CheckServiceIDInWalletService implements [storage.WalletOracleRepository].
func (q *WalletStorage) CheckServiceIDInWalletService(ctx context.Context, selfServiceID string, otherServiceID string, agentServiceID string) ([]string, error) {
	// Collect non-empty IDs
	ids := make([]string, 0, 3)
	if selfServiceID != "" {
		ids = append(ids, selfServiceID)
	}
	if otherServiceID != "" {
		ids = append(ids, otherServiceID)
	}
	if agentServiceID != "" {
		ids = append(ids, agentServiceID)
	}
	if len(ids) == 0 {
		return []string{}, nil
	}

	// Build placeholders for query
	placeholders := make([]string, len(ids))
	for i := range ids {
		placeholders[i] = fmt.Sprintf("HEXTORAW(:%d)", i+1)
	}
	query := fmt.Sprintf("SELECT RAWTOHEX(SERVICE_ID) FROM WALLET_SERVICES WHERE SERVICE_ID IN (%s) AND is_deleted = 0", strings.Join(placeholders, ", "))

	rows, err := q.db.QueryContext(ctx, query, toInterfaceSlice(ids)...)
	if err != nil {
		q.logger.Errorf("[WalletRepo][CheckIfIDsExist] query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer rows.Close()

	existingIDs := make([]string, 0, len(ids))
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			q.logger.Errorf("[WalletRepo][CheckIfIDsExist] scan failed: %v", err)
			return nil, local_util.HandleDBError(err)
		}
		existingIDs = append(existingIDs, id)
	}
	return existingIDs, nil
}

// toInterfaceSlice converts a string slice to an interface{} slice for variadic SQL args
func toInterfaceSlice(strs []string) []interface{} {
	res := make([]interface{}, len(strs))
	for i, v := range strs {
		res[i] = v
	}
	return res
}
