package donation_company_oracle

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/dto/donation_company"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/lib"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/localization"
	imodel "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/model"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/types"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage"
	local_util "github.com/abelshebela/cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type repository struct {
	db     DBTX
	logger utils.Logger
}

var _ storage.DonationCompanyRepository = (*repository)(nil)

func NewRepository(db DBTX, logger utils.Logger) storage.DonationCompanyRepository {
	return &repository{
		db:     db,
		logger: logger,
	}
}

func (r *repository) WithTx(tx *sql.Tx) *repository {
	return &repository{db: tx, logger: r.logger}
}

func (r *repository) Create(ctx context.Context, c *imodel.DonationCompanyOracle) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[DonationCompanyOracle][Create] code=%s name=%s", c.CompanyCode, c.CompanyName)

	q := `INSERT INTO DONATION_COMPANIES
	      (COMPANY_NAME, COMPANY_CODE, COMPANY_LOGO, COMPANY_DESCRIPTION,
	       ADDRESS, PHONE_NUMBER, EMAIL,
	       IS_DELETED, ENABLED)
	      VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9)`

	if _, err := r.db.ExecContext(ctx, q,
		c.CompanyName,
		c.CompanyCode,
		c.CompanyLogo,
		c.CompanyDescription,
		c.Address,
		c.PhoneNumber,
		c.Email,
		boolToInt(c.IsDeleted),
		boolToInt(c.Enabled),
	); err != nil {
		log.Errorf("[DonationCompanyOracle][Create] insert failed: %v", err)
		return local_util.HandleDBError(err)
	}
	return nil
}

// Update writes all fields (including IS_DELETED) and matches by ID only —
// soft-deleted rows are reachable, mirroring the legacy Mongo path.
func (r *repository) Update(ctx context.Context, id string, c *imodel.DonationCompanyOracle) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[DonationCompanyOracle][Update] id=%s", id)

	q := `UPDATE DONATION_COMPANIES
	      SET COMPANY_NAME        = :1,
	          COMPANY_CODE        = :2,
	          COMPANY_LOGO        = :3,
	          COMPANY_DESCRIPTION = :4,
	          ADDRESS             = :5,
	          PHONE_NUMBER        = :6,
	          EMAIL               = :7,
	          ENABLED             = :8,
	          IS_DELETED          = :9,
	          LAST_MODIFIED_AT    = :10
	      WHERE ID = HEXTORAW(:11)`

	res, err := r.db.ExecContext(ctx, q,
		c.CompanyName,
		c.CompanyCode,
		c.CompanyLogo,
		c.CompanyDescription,
		c.Address,
		c.PhoneNumber,
		c.Email,
		boolToInt(c.Enabled),
		boolToInt(c.IsDeleted),
		c.LastModifiedAt,
		id,
	)
	if err != nil {
		log.Errorf("[DonationCompanyOracle][Update] failed: %v", err)
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

	log.Infof("[DonationCompanyOracle][Delete] id=%s", id)

	q := `UPDATE DONATION_COMPANIES SET IS_DELETED = 1
	      WHERE ID = HEXTORAW(:1) AND IS_DELETED = 0`

	if _, err := r.db.ExecContext(ctx, q, id); err != nil {
		log.Errorf("[DonationCompanyOracle][Delete] failed: %v", err)
		return local_util.HandleDBError(err)
	}
	return nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*donation_company.DonationCompanyListResponse, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[DonationCompanyOracle][FindByID] id=%s", id)

	q := `SELECT ` + donationCompanySelectCols + `
	      FROM DONATION_COMPANIES
	      WHERE ID = HEXTORAW(:1) AND IS_DELETED = 0`

	row := r.db.QueryRowContext(ctx, q, id)
	out, err := scanDonationCompanyRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		log.Errorf("[DonationCompanyOracle][FindByID] scan failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return out, nil
}

func (r *repository) FindAllWithPagination(
	ctx context.Context,
	filterParam types.Filter,
) (*types.PaginatedResponse[[]donation_company.DonationCompanyListResponse], error) {
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
		conds = append(conds, fmt.Sprintf(
			"(UPPER(COMPANY_NAME) LIKE UPPER(:%d) OR UPPER(COMPANY_CODE) LIKE UPPER(:%d))",
			idx, idx+1,
		))
		args = append(args, "%"+s+"%", "%"+s+"%")
		idx += 2
	}

	if filterParam.Filters != nil {
		if v, ok := filterParam.Filters["enabled"]; ok {
			on, apply := lib.BoolFromInterface(v)
			if apply {
				conds = append(conds, fmt.Sprintf("ENABLED = :%d", idx))
				args = append(args, boolToInt(on))
				idx++
			}
		}
		if v, ok := filterParam.Filters["company_name"]; ok {
			if s, _ := v.(string); strings.TrimSpace(s) != "" {
				conds = append(conds, fmt.Sprintf("UPPER(COMPANY_NAME) = UPPER(:%d)", idx))
				args = append(args, s)
				idx++
			}
		}
		if v, ok := filterParam.Filters["company_code"]; ok {
			if s, _ := v.(string); strings.TrimSpace(s) != "" {
				conds = append(conds, fmt.Sprintf("COMPANY_CODE = :%d", idx))
				args = append(args, s)
				idx++
			}
		}
	}

	whereClause := strings.Join(conds, " AND ")

	countQuery := `SELECT COUNT(*) FROM DONATION_COMPANIES WHERE ` + whereClause
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		log.Errorf("[DonationCompanyOracle][FindAllWithPagination] count failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	meta := local_util.BuildPaginationMeta(total, page, perPage)
	if total == 0 {
		return &types.PaginatedResponse[[]donation_company.DonationCompanyListResponse]{
			Data: []donation_company.DonationCompanyListResponse{},
			Meta: meta,
		}, nil
	}

	selectQuery := `SELECT ` + donationCompanySelectCols + `
	                FROM DONATION_COMPANIES
	                WHERE ` + whereClause + `
	                ORDER BY CREATED_AT DESC
	                OFFSET :` + fmt.Sprintf("%d", idx) + ` ROWS FETCH NEXT :` + fmt.Sprintf("%d", idx+1) + ` ROWS ONLY`
	args = append(args, offset, perPage)

	rows, err := r.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		log.Errorf("[DonationCompanyOracle][FindAllWithPagination] query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer rows.Close()

	out := []donation_company.DonationCompanyListResponse{}
	for rows.Next() {
		item, err := scanDonationCompanyRow(rows)
		if err != nil {
			log.Errorf("[DonationCompanyOracle][FindAllWithPagination] scan failed: %v", err)
			return nil, local_util.HandleDBError(err)
		}
		out = append(out, *item)
	}
	if err := rows.Err(); err != nil {
		log.Errorf("[DonationCompanyOracle][FindAllWithPagination] rows: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	return &types.PaginatedResponse[[]donation_company.DonationCompanyListResponse]{
		Data: out,
		Meta: meta,
	}, nil
}
