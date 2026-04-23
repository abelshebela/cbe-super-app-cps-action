package amount_based_auth_oracle

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	local_model "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type Repository struct {
	db     *sql.DB
	logger shared_utils.Logger
}

var _ storage.AmountBasedAuthOracleRepository = (*Repository)(nil)

func NewAmountBasedAuthOracleRepository(db *sql.DB, logger shared_utils.Logger) storage.AmountBasedAuthOracleRepository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

func formatNullTime(t sql.NullTime) time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}

func (r *Repository) Create(ctx context.Context, tier *local_model.AuthTierOracle) error {
	if tier == nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	q := `INSERT INTO AMOUNT_BASED_AUTH_TIERS
		(id, currency, min_amount, max_amount, method, enabled, is_deleted, created_at, last_modified)
		VALUES (SYS_GUID(), :1, :2, :3, :4, :5, :6, SYSTIMESTAMP, SYSTIMESTAMP)`

	_, err := r.db.ExecContext(ctx, q,
		string(tier.Currency),
		tier.MinAmount,
		tier.MaxAmount,
		string(tier.Method),
		tier.Enabled,
		tier.IsDeleted,
	)
	if err != nil {
		r.logger.Errorf("[AmountBasedAuthOracle][Create] failed: %v", err)
		return local_util.HandleDBError(err)
	}
	return nil
}

func (r *Repository) CreateMany(ctx context.Context, tiers []local_model.AuthTierOracle) error {
	for i := range tiers {
		t := tiers[i]
		if err := r.Create(ctx, &t); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) Update(ctx context.Context, id string, update *local_model.AuthTierOracle) error {
	if update == nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	q := `UPDATE AMOUNT_BASED_AUTH_TIERS
		SET currency = :1,
			min_amount = :2,
			max_amount = :3,
			method = :4,
			enabled = :5,
			is_deleted = :6,
			last_modified = SYSTIMESTAMP
		WHERE id = HEXTORAW(:7) AND is_deleted = 0`

	_, err := r.db.ExecContext(ctx, q,
		string(update.Currency),
		update.MinAmount,
		update.MaxAmount,
		string(update.Method),
		update.Enabled,
		update.IsDeleted,
		id,
	)
	if err != nil {
		r.logger.Errorf("[AmountBasedAuthOracle][Update] failed: %v", err)
		return local_util.HandleDBError(err)
	}
	return nil
}

func (r *Repository) DeleteByCurrency(ctx context.Context, currency string) error {
	q := `UPDATE AMOUNT_BASED_AUTH_TIERS
		SET is_deleted = 1,
			last_modified = CURRENT_TIMESTAMP
		WHERE currency = :1 AND is_deleted = 0`

	_, err := r.db.ExecContext(ctx, q, currency)
	if err != nil {
		r.logger.Errorf("[AmountBasedAuthOracle][DeleteByCurrency] failed: %v", err)
		return local_util.HandleDBError(err)
	}
	return nil
}

func (r *Repository) DeleteByID(ctx context.Context, id string) error {
	q := `UPDATE AMOUNT_BASED_AUTH_TIERS
		SET is_deleted = 1,
			last_modified = SYSTIMESTAMP
		WHERE id = HEXTORAW(:1) AND is_deleted = 0`

	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		r.logger.Errorf("[AmountBasedAuthOracle][DeleteByID] failed: %v", err)
		return local_util.HandleDBError(err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return local_util.HandleDBError(err)
	}
	if n == 0 {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}
	return nil
}

func (r *Repository) CurrencyExists(ctx context.Context, currency string) (bool, error) {
	q := `SELECT COUNT(*) FROM AMOUNT_BASED_AUTH_TIERS WHERE currency = :1 AND is_deleted = 0`

	var total int64
	if err := r.db.QueryRowContext(ctx, q, currency).Scan(&total); err != nil {
		r.logger.Errorf("[AmountBasedAuthOracle][CurrencyExists] failed: %v", err)
		return false, local_util.HandleDBError(err)
	}
	return total > 0, nil
}

func (r *Repository) FindByID(ctx context.Context, id string) (*local_model.AuthTierOracle, error) {
	q := `SELECT RAWTOHEX(id) AS id,
			currency, min_amount, max_amount, method,
			enabled, is_deleted, created_at, last_modified
		FROM AMOUNT_BASED_AUTH_TIERS
		WHERE id = HEXTORAW(:1) AND is_deleted = 0`

	row := r.db.QueryRowContext(ctx, q, id)
	var tier local_model.AuthTierOracle
	var createdAt, lastModified sql.NullTime

	if err := row.Scan(
		&tier.ID,
		&tier.Currency,
		&tier.MinAmount,
		&tier.MaxAmount,
		&tier.Method,
		&tier.Enabled,
		&tier.IsDeleted,
		&createdAt,
		&lastModified,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		return nil, err
	}

	tier.CreatedAt = formatNullTime(createdAt)
	tier.LastModified = formatNullTime(lastModified)
	return &tier, nil
}

func (r *Repository) FindActiveByCurrencyAndMethod(ctx context.Context, currency string, method constants.Method) ([]local_model.AuthTierOracle, error) {
	q := `SELECT RAWTOHEX(id) AS id,
			currency, min_amount, max_amount, method,
			enabled, is_deleted, created_at, last_modified
		FROM AMOUNT_BASED_AUTH_TIERS
		WHERE currency = :1 AND method = :2 AND is_deleted = 0
		ORDER BY last_modified DESC`

	rows, err := r.db.QueryContext(ctx, q, currency, string(method))
	if err != nil {
		return nil, local_util.HandleDBError(err)
	}
	defer rows.Close()

	var result []local_model.AuthTierOracle
	for rows.Next() {
		var tier local_model.AuthTierOracle
		var createdAt, lastModified sql.NullTime
		if err := rows.Scan(
			&tier.ID,
			&tier.Currency,
			&tier.MinAmount,
			&tier.MaxAmount,
			&tier.Method,
			&tier.Enabled,
			&tier.IsDeleted,
			&createdAt,
			&lastModified,
		); err != nil {
			return nil, err
		}
		tier.CreatedAt = formatNullTime(createdAt)
		tier.LastModified = formatNullTime(lastModified)
		result = append(result, tier)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Repository) FindActiveByCurrency(ctx context.Context, currency string) ([]local_model.AuthTierOracle, error) {
	q := `SELECT RAWTOHEX(id) AS id,
			currency, min_amount, max_amount, method,
			enabled, is_deleted, created_at, last_modified
		FROM AMOUNT_BASED_AUTH_TIERS
		WHERE currency = :1 AND is_deleted = 0
		ORDER BY last_modified DESC`

	rows, err := r.db.QueryContext(ctx, q, currency)
	if err != nil {
		return nil, local_util.HandleDBError(err)
	}
	defer rows.Close()

	var result []local_model.AuthTierOracle
	for rows.Next() {
		var tier local_model.AuthTierOracle
		var createdAt, lastModified sql.NullTime
		if err := rows.Scan(
			&tier.ID,
			&tier.Currency,
			&tier.MinAmount,
			&tier.MaxAmount,
			&tier.Method,
			&tier.Enabled,
			&tier.IsDeleted,
			&createdAt,
			&lastModified,
		); err != nil {
			return nil, err
		}
		tier.CreatedAt = formatNullTime(createdAt)
		tier.LastModified = formatNullTime(lastModified)
		result = append(result, tier)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *Repository) FindAllActiveForSearch(ctx context.Context, search string) ([]local_model.AuthTierOracle, error) {
	search = strings.TrimSpace(search)
	var q string
	var args []interface{}

	q = `SELECT RAWTOHEX(id) AS id,
			currency, min_amount, max_amount, method,
			enabled, is_deleted, created_at, last_modified
		FROM AMOUNT_BASED_AUTH_TIERS
		WHERE is_deleted = 0`

	if search != "" {
		q += ` AND (UPPER(method) LIKE :1 OR UPPER(currency) LIKE :1)`
		args = append(args, "%"+strings.ToUpper(search)+"%")
	}

	q += ` ORDER BY currency, method`

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		r.logger.Errorf("[GetAllAmountBasedAuth] persistance query context error: %v",err)
		return nil, local_util.HandleDBError(err)
	}
	defer rows.Close()

	var result []local_model.AuthTierOracle
	for rows.Next() {
		var tier local_model.AuthTierOracle
		var createdAt, lastModified sql.NullTime
		if err := rows.Scan(
			&tier.ID,
			&tier.Currency,
			&tier.MinAmount,
			&tier.MaxAmount,
			&tier.Method,
			&tier.Enabled,
			&tier.IsDeleted,
			&createdAt,
			&lastModified,
		); err != nil {
		r.logger.Errorf("[GetAllAmountBasedAuth] scan  error:%v",err)
			return nil, err
		}
		tier.CreatedAt = formatNullTime(createdAt)
		tier.LastModified = formatNullTime(lastModified)
		result = append(result, tier)
	}
	if err := rows.Err(); err != nil {
		r.logger.Errorf("[GetAllAmountBasedAuth] row error : %v",err)

		return nil, err
	}
	return result, nil
}
