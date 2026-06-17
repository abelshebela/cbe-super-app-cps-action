package donation_category_oracle

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"cbe-super-app-cps-action/internal/constants/dto/donation_category"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type repository struct {
	db     DBTX
	logger utils.Logger
}

var oracleHexIDRegex = regexp.MustCompile(`^[a-fA-F0-9]{32}$`)

var _ storage.DonationCategoryRepository = (*repository)(nil)

func NewRepository(db DBTX, logger utils.Logger) storage.DonationCategoryRepository {
	return &repository{
		db:     db,
		logger: logger,
	}
}

// WithTx returns a bare repo bound to tx. Kafka publishes are skipped — the
// caller publishes after Commit.
func (r *repository) WithTx(tx *sql.Tx) *repository {
	return &repository{db: tx, logger: r.logger}
}

func (r *repository) Create(ctx context.Context, c *imodel.DonationCategoryOracle) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[DonationCategoryOracle][Create] inserting category=%s", c.CategoryName)

	q := `INSERT INTO DONATION_CATEGORIES (CATEGORY_NAME, ICON) VALUES (:1, :2)`
	if _, err := r.db.ExecContext(ctx, q, c.CategoryName, c.Icon); err != nil {
		log.Errorf("[DonationCategoryOracle][Create] insert failed: %v", err)
		return local_util.HandleDBError(err)
	}
	return nil
}

// Update preserves existing CATEGORY_NAME / ICON when the input is empty
// (COALESCE(NULLIF, ...)). LAST_MODIFIED_AT comes from the caller.
func (r *repository) Update(ctx context.Context, id string, c *imodel.DonationCategoryOracle) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[DonationCategoryOracle][Update] id=%s", id)

	q := `UPDATE DONATION_CATEGORIES
	      SET CATEGORY_NAME    = COALESCE(NULLIF(:1, ''), CATEGORY_NAME),
	          ICON             = COALESCE(NULLIF(:2, ''), ICON),
	          LAST_MODIFIED_AT = :3
	      WHERE ID = HEXTORAW(:4) AND IS_DELETED = 0`

	res, err := r.db.ExecContext(ctx, q, c.CategoryName, c.Icon, c.LastModifiedAt, id)
	if err != nil {
		log.Errorf("[DonationCategoryOracle][Update] failed: %v", err)
		return local_util.HandleDBError(err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}
	return nil
}

func (r *repository) EnableDisable(ctx context.Context, id string, enable bool) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[DonationCategoryOracle][EnableDisable] id=%s enable=%v", id, enable)

	v := 0
	if enable {
		v = 1
	}

	q := `UPDATE DONATION_CATEGORIES
	      SET ENABLED = :1, LAST_MODIFIED_AT = :2
	      WHERE ID = HEXTORAW(:3) AND IS_DELETED = 0`

	res, err := r.db.ExecContext(ctx, q, v, time.Now(), id)
	if err != nil {
		log.Errorf("[DonationCategoryOracle][EnableDisable] failed: %v", err)
		return local_util.HandleDBError(err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}
	return nil
}

// Delete is a soft-delete (IS_DELETED = 1). Returns nil if the row is missing
// or already deleted.
func (r *repository) Delete(ctx context.Context, id string) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[DonationCategoryOracle][Delete] id=%s", id)

	q := `UPDATE DONATION_CATEGORIES SET IS_DELETED = 1
	      WHERE ID = HEXTORAW(:1) AND IS_DELETED = 0`

	if _, err := r.db.ExecContext(ctx, q, id); err != nil {
		log.Errorf("[DonationCategoryOracle][Delete] failed: %v", err)
		return local_util.HandleDBError(err)
	}
	return nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*donation_category.DonationCategoryListResponse, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[DonationCategoryOracle][FindByID] id=%s", id)

	if !oracleHexIDRegex.MatchString(strings.TrimSpace(id)) {
		log.Errorf("[DonationCategoryOracle][FindByID] invalid id format: %s", id)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	q := `SELECT ` + donationCategorySelectCols + `
	      FROM DONATION_CATEGORIES
	      WHERE ID = HEXTORAW(:1) AND IS_DELETED = 0`

	row := r.db.QueryRowContext(ctx, q, id)
	out, err := scanDonationCategoryRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		log.Errorf("[DonationCategoryOracle][FindByID] scan failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return out, nil
}

// FindByName: case-insensitive lookup among non-deleted rows. Returns
// ErrorResourceNotFound on miss so core.DonationNameExists string-matches.
func (r *repository) FindByName(ctx context.Context, name string) (*donation_category.DonationCategoryListResponse, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[DonationCategoryOracle][FindByName] name=%s", name)

	q := `SELECT ` + donationCategorySelectCols + `
	      FROM DONATION_CATEGORIES
	      WHERE UPPER(CATEGORY_NAME) = UPPER(:1) AND IS_DELETED = 0
	      FETCH FIRST 1 ROW ONLY`

	row := r.db.QueryRowContext(ctx, q, name)
	out, err := scanDonationCategoryRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		log.Errorf("[DonationCategoryOracle][FindByName] scan failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return out, nil
}

func (r *repository) FindAllWithPagination(
	ctx context.Context,
	filterParam types.Filter,
) (*types.PaginatedResponse[[]donation_category.DonationCategoryListResponse], error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	page := filterParam.Page
	if page < 1 {
		page = 1
	}
	perPage := filterParam.PerPage
	if perPage < 1 {
		perPage = 10
	}
	offset := (page - 1) * perPage

	conds := []string{"IS_DELETED = 0"}
	args := []interface{}{}
	idx := 1

	if s := strings.TrimSpace(filterParam.Search); s != "" {
		conds = append(conds, fmt.Sprintf("UPPER(CATEGORY_NAME) LIKE UPPER(:%d)", idx))
		args = append(args, "%"+s+"%")
		idx++
	}

	if filterParam.Filters != nil {
		if v, ok := filterParam.Filters["enabled"]; ok {
			on, apply := lib.BoolFromInterface(v)
			if apply {
				n := 0
				if on {
					n = 1
				}
				conds = append(conds, fmt.Sprintf("ENABLED = :%d", idx))
				args = append(args, n)
				idx++
			}
		}
		if v, ok := filterParam.Filters["category_name"]; ok {
			if s, _ := v.(string); strings.TrimSpace(s) != "" {
				conds = append(conds, fmt.Sprintf("UPPER(CATEGORY_NAME) = UPPER(:%d)", idx))
				args = append(args, s)
				idx++
			}
		}
	}

	whereClause := strings.Join(conds, " AND ")

	countQuery := `SELECT COUNT(*) FROM DONATION_CATEGORIES WHERE ` + whereClause
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		log.Errorf("[DonationCategoryOracle][FindAllWithPagination] count failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	meta := local_util.BuildPaginationMeta(total, page, perPage)
	if total == 0 {
		return &types.PaginatedResponse[[]donation_category.DonationCategoryListResponse]{
			Data: []donation_category.DonationCategoryListResponse{},
			Meta: meta,
		}, nil
	}

	selectQuery := `SELECT ` + donationCategorySelectCols + `
	                FROM DONATION_CATEGORIES
	                WHERE ` + whereClause + `
	                ORDER BY CREATED_AT DESC
	                OFFSET :` + fmt.Sprintf("%d", idx) + ` ROWS FETCH NEXT :` + fmt.Sprintf("%d", idx+1) + ` ROWS ONLY`
	args = append(args, offset, perPage)

	rows, err := r.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		log.Errorf("[DonationCategoryOracle][FindAllWithPagination] query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer rows.Close()

	out := []donation_category.DonationCategoryListResponse{}
	for rows.Next() {
		item, err := scanDonationCategoryRow(rows)
		if err != nil {
			log.Errorf("[DonationCategoryOracle][FindAllWithPagination] scan failed: %v", err)
			return nil, local_util.HandleDBError(err)
		}
		out = append(out, *item)
	}
	if err := rows.Err(); err != nil {
		log.Errorf("[DonationCategoryOracle][FindAllWithPagination] rows: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	return &types.PaginatedResponse[[]donation_category.DonationCategoryListResponse]{
		Data: out,
		Meta: meta,
	}, nil
}
