package donation_oracle

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	donation_dto "cbe-super-app-cps-action/internal/constants/dto/donation"
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

var _ storage.DonationRepository = (*repository)(nil)

func NewRepository(db DBTX, logger utils.Logger) storage.DonationRepository {
	return &repository{db: db, logger: logger}
}

func (r *repository) WithTx(tx *sql.Tx) *repository {
	return &repository{db: tx, logger: r.logger}
}

// Create wraps INSERT + DONATION_IMAGES inserts in a single transaction when
// r.db is *sql.DB; if it's already *sql.Tx (via WithTx), the caller commits.
func (r *repository) Create(ctx context.Context, d *imodel.DonationOracle) error {
	r.logger.Infof("[DonationOracle][Create] code=%s title=%s", d.DonationCode, d.Title)

	if dbHandle, ok := r.db.(*sql.DB); ok {
		tx, err := dbHandle.BeginTx(ctx, nil)
		if err != nil {
			r.logger.Errorf("[DonationOracle][Create] BeginTx failed: %v", err)
			return local_util.HandleDBError(err)
		}
		defer func() { _ = tx.Rollback() }()
		if err := r.WithTx(tx).createInner(ctx, d); err != nil {
			return err
		}
		if err := tx.Commit(); err != nil {
			r.logger.Errorf("[DonationOracle][Create] Commit failed: %v", err)
			return local_util.HandleDBError(err)
		}
		return nil
	}
	return r.createInner(ctx, d)
}

func (r *repository) createInner(ctx context.Context, d *imodel.DonationOracle) error {
	insertParent := `INSERT INTO DONATIONS
		(DONATION_CODE, SERVICE_ID, COMPANY_ID, CATEGORY_ID, TITLE, IS_FEATURED, TARGET, CURRENT_AMOUNT,
		 DONATION_DESCRIPTION, COVER_IMAGE, START_DATE, END_DATE, IS_DELETED, ENABLED)
		VALUES (:1, HEXTORAW(:2), HEXTORAW(:3), HEXTORAW(:4), :5, :6, :7, :8, :9, :10, :11, :12, :13, :14)`

	if _, err := r.db.ExecContext(ctx, insertParent,
		d.DonationCode,
		d.ServiceID,
		d.CompanyID,
		d.CategoryID,
		d.Title,
		boolToInt(d.IsFeatured),
		nullStringFromString(d.Target),
		nullStringFromString(d.CurrentAmount),
		nullStringFromString(d.DonationDescription),
		nullStringFromString(d.CoverImage),
		d.StartDate,
		nullTimeFromTime(d.EndDate),
		boolToInt(d.IsDeleted),
		boolToInt(d.Enabled),
	); err != nil {
		r.logger.Errorf("[DonationOracle][Create] parent insert failed: %v", err)
		return local_util.HandleDBError(err)
	}

	var newID string
	if err := r.db.QueryRowContext(ctx,
		`SELECT RAWTOHEX(ID) FROM DONATIONS WHERE DONATION_CODE = :1`, d.DonationCode,
	).Scan(&newID); err != nil {
		r.logger.Errorf("[DonationOracle][Create] resolve new id failed: %v", err)
		return local_util.HandleDBError(err)
	}

	if err := r.replaceImages(ctx, newID, d.DonationImages); err != nil {
		return err
	}
	return nil
}

func (r *repository) Update(ctx context.Context, id string, d *imodel.DonationOracle) error {
	r.logger.Infof("[DonationOracle][Update] id=%s", id)

	if dbHandle, ok := r.db.(*sql.DB); ok {
		tx, err := dbHandle.BeginTx(ctx, nil)
		if err != nil {
			r.logger.Errorf("[DonationOracle][Update] BeginTx failed: %v", err)
			return local_util.HandleDBError(err)
		}
		defer func() { _ = tx.Rollback() }()
		if err := r.WithTx(tx).updateInner(ctx, id, d); err != nil {
			return err
		}
		if err := tx.Commit(); err != nil {
			r.logger.Errorf("[DonationOracle][Update] Commit failed: %v", err)
			return local_util.HandleDBError(err)
		}
		return nil
	}
	return r.updateInner(ctx, id, d)
}

func (r *repository) updateInner(ctx context.Context, id string, d *imodel.DonationOracle) error {
	q := `UPDATE DONATIONS
	      SET DONATION_CODE        = :1,
	          SERVICE_ID           = HEXTORAW(:2),
	          COMPANY_ID           = HEXTORAW(:3),
	          CATEGORY_ID          = HEXTORAW(:4),
	          TITLE                = :5,
	          IS_FEATURED          = :6,
	          TARGET               = :7,
	          CURRENT_AMOUNT       = :8,
	          DONATION_DESCRIPTION = :9,
	          COVER_IMAGE          = :10,
	          START_DATE           = :11,
	          END_DATE             = :12,
	          ENABLED              = :13,
	          LAST_MODIFIED_AT     = :14
	      WHERE ID = HEXTORAW(:15) AND IS_DELETED = 0`

	res, err := r.db.ExecContext(ctx, q,
		d.DonationCode,
		d.ServiceID,
		d.CompanyID,
		d.CategoryID,
		d.Title,
		boolToInt(d.IsFeatured),
		nullStringFromString(d.Target),
		nullStringFromString(d.CurrentAmount),
		nullStringFromString(d.DonationDescription),
		nullStringFromString(d.CoverImage),
		d.StartDate,
		nullTimeFromTime(d.EndDate),
		boolToInt(d.Enabled),
		d.LastModifiedAt,
		id,
	)
	if err != nil {
		r.logger.Errorf("[DonationOracle][Update] failed: %v", err)
		return local_util.HandleDBError(err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}

	if err := r.replaceImages(ctx, id, d.DonationImages); err != nil {
		return err
	}
	return nil
}

// Delete is a soft-delete (IS_DELETED = 1). Returns nil if the row is missing
// or already deleted.
func (r *repository) Delete(ctx context.Context, id string) error {
	r.logger.Infof("[DonationOracle][Delete] id=%s", id)
	q := `UPDATE DONATIONS SET IS_DELETED = 1
	      WHERE ID = HEXTORAW(:1) AND IS_DELETED = 0`
	if _, err := r.db.ExecContext(ctx, q, id); err != nil {
		r.logger.Errorf("[DonationOracle][Delete] failed: %v", err)
		return local_util.HandleDBError(err)
	}
	return nil
}

func (r *repository) FindByServiceID(ctx context.Context, serviceID string) (*donation_dto.DonationListResponse, error) {
	r.logger.Infof("[DonationOracle][FindByServiceID] serviceID=%s", serviceID)
	q := `SELECT ` + donationListSelectCols + donationListFromJoin + `
	      WHERE d.SERVICE_ID = HEXTORAW(:1) AND d.IS_DELETED = 0`
	row := r.db.QueryRowContext(ctx, q, serviceID)
	out, err := scanDonationListRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		r.logger.Errorf("[DonationOracle][FindByServiceID] scan failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	return out, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*donation_dto.DonationListResponse, error) {
	r.logger.Infof("[DonationOracle][FindByID] id=%s", id)

	q := `SELECT ` + donationListSelectCols + donationListFromJoin + `
	      WHERE d.ID = HEXTORAW(:1) AND d.IS_DELETED = 0`

	row := r.db.QueryRowContext(ctx, q, id)
	out, err := scanDonationListRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		r.logger.Errorf("[DonationOracle][FindByID] scan failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	images, err := r.fetchImagesFor(ctx, []string{id})
	if err != nil {
		return nil, err
	}

	if imgs, ok := images[id]; ok {
		out.DonationImages = imgs
	}
	return out, nil
}

func (r *repository) FindAllWithPagination(
	ctx context.Context,
	filterParam types.Filter,
) (*types.PaginatedResponse[[]donation_dto.DonationListResponse], error) {

	page := filterParam.Page
	if page < 1 {
		page = 1
	}
	perPage := filterParam.PerPage
	if perPage < 1 {
		perPage = 10
	}
	offset := (page - 1) * perPage

	conds := []string{"d.IS_DELETED = 0"}
	args := []interface{}{}
	idx := 1

	if s := strings.TrimSpace(filterParam.Search); s != "" {
		conds = append(conds, fmt.Sprintf(
			"(UPPER(d.TITLE) LIKE UPPER(:%d) OR UPPER(d.DONATION_CODE) LIKE UPPER(:%d) OR UPPER(d.DONATION_DESCRIPTION) LIKE UPPER(:%d))",
			idx, idx+1, idx+2,
		))
		args = append(args, "%"+s+"%", "%"+s+"%", "%"+s+"%")
		idx += 3
	}

	if filterParam.Filters != nil {
		if v, ok := filterParam.Filters["title"]; ok {
			if s, _ := v.(string); strings.TrimSpace(s) != "" {
				conds = append(conds, fmt.Sprintf("UPPER(d.TITLE) = UPPER(:%d)", idx))
				args = append(args, s)
				idx++
			}
		}
		if v, ok := filterParam.Filters["donation_code"]; ok {
			if s, _ := v.(string); strings.TrimSpace(s) != "" {
				conds = append(conds, fmt.Sprintf("d.DONATION_CODE = :%d", idx))
				args = append(args, s)
				idx++
			}
		}
		if v, ok := filterParam.Filters["target"]; ok {
			if s, _ := v.(string); strings.TrimSpace(s) != "" {
				conds = append(conds, fmt.Sprintf("d.TARGET = :%d", idx))
				args = append(args, s)
				idx++
			}
		}
		if v, ok := filterParam.Filters["is_featured"]; ok {
			if b, applied := lib.BoolFromInterface(v); applied {
				conds = append(conds, fmt.Sprintf("d.IS_FEATURED = :%d", idx))
				args = append(args, boolToInt(b))
				idx++
			}
		}
		if v, ok := filterParam.Filters["enabled"]; ok {
			if b, applied := lib.BoolFromInterface(v); applied {
				// Both enabled=true and enabled=false exclude expired campaigns:
				// strict END_DATE > now. NULL end_dates are excluded by Oracle's
				// 3-valued logic. Use is_expired=true to see expired records.
				conds = append(conds, fmt.Sprintf("d.ENABLED = :%d", idx))
				args = append(args, boolToInt(b))
				idx++
				conds = append(conds, fmt.Sprintf("d.END_DATE > :%d", idx))
				args = append(args, time.Now())
				idx++
			}
		}
		if v, ok := filterParam.Filters["is_expired"]; ok {
			if b, applied := lib.BoolFromInterface(v); applied && b {
				conds = append(conds, fmt.Sprintf("d.END_DATE < :%d", idx))
				args = append(args, time.Now())
				idx++
			}
		}

		// Equality on YYYY-MM-DD expands to a 24-hour BETWEEN range; full
		// RFC3339 inputs do an exact match. `_from`/`_to` provide ranges.
		for _, dc := range []struct{ key, column string }{
			{"start_date", "d.START_DATE"},
			{"end_date", "d.END_DATE"},
			{"created_at", "d.CREATED_AT"},
		} {
			if v, ok := filterParam.Filters[dc.key]; ok {
				if s, _ := v.(string); strings.TrimSpace(s) != "" {
					if t, err := local_util.ParseDateInput(s); err == nil {
						if !strings.Contains(s, "T") {
							conds = append(conds, fmt.Sprintf("%s BETWEEN :%d AND :%d", dc.column, idx, idx+1))
							args = append(args, t, t.Add(24*time.Hour-time.Millisecond))
							idx += 2
						} else {
							conds = append(conds, fmt.Sprintf("%s = :%d", dc.column, idx))
							args = append(args, t)
							idx++
						}
					}
				}
			}
			if v, ok := filterParam.Filters[dc.key+"_from"]; ok {
				if s, _ := v.(string); strings.TrimSpace(s) != "" {
					if t, err := local_util.ParseDateInput(s); err == nil {
						conds = append(conds, fmt.Sprintf("%s >= :%d", dc.column, idx))
						args = append(args, t)
						idx++
					}
				}
			}
			if v, ok := filterParam.Filters[dc.key+"_to"]; ok {
				if s, _ := v.(string); strings.TrimSpace(s) != "" {
					if t, err := local_util.ParseDateInput(s); err == nil {
						if !strings.Contains(s, "T") {
							t = t.Add(24*time.Hour - time.Millisecond)
						}
						conds = append(conds, fmt.Sprintf("%s <= :%d", dc.column, idx))
						args = append(args, t)
						idx++
					}
				}
			}
		}
	}

	whereClause := strings.Join(conds, " AND ")

	countQuery := `SELECT COUNT(*) FROM DONATIONS d WHERE ` + whereClause
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		r.logger.Errorf("[DonationOracle][FindAllWithPagination] count failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	meta := local_util.BuildPaginationMeta(total, page, perPage)
	if total == 0 {
		return &types.PaginatedResponse[[]donation_dto.DonationListResponse]{
			Data: []donation_dto.DonationListResponse{},
			Meta: meta,
		}, nil
	}

	selectQuery := `SELECT ` + donationListSelectCols + donationListFromJoin + `
	                WHERE ` + whereClause + `
	                ORDER BY d.CREATED_AT DESC
	                OFFSET :` + fmt.Sprintf("%d", idx) + ` ROWS FETCH NEXT :` + fmt.Sprintf("%d", idx+1) + ` ROWS ONLY`
	args = append(args, offset, perPage)

	rows, err := r.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		r.logger.Errorf("[DonationOracle][FindAllWithPagination] query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer rows.Close()

	out := []donation_dto.DonationListResponse{}
	ids := []string{}
	for rows.Next() {
		item, err := scanDonationListRow(rows)
		if err != nil {
			r.logger.Errorf("[DonationOracle][FindAllWithPagination] scan failed: %v", err)
			return nil, local_util.HandleDBError(err)
		}
		out = append(out, *item)
		ids = append(ids, item.ID)
	}
	if err := rows.Err(); err != nil {
		r.logger.Errorf("[DonationOracle][FindAllWithPagination] rows: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	if len(ids) > 0 {
		images, err := r.fetchImagesFor(ctx, ids)
		if err != nil {
			return nil, err
		}
		for i := range out {
			if imgs, ok := images[out[i].ID]; ok {
				out[i].DonationImages = imgs
			}
		}
	}

	return &types.PaginatedResponse[[]donation_dto.DonationListResponse]{
		Data: out,
		Meta: meta,
	}, nil
}

func (r *repository) StreamByDateRange(
	ctx context.Context,
	startDate, endDate time.Time,
	handler func(*imodel.DonationOracle) error,
) error {
	q := `SELECT
	        RAWTOHEX(d.ID), d.DONATION_CODE,
	        RAWTOHEX(d.SERVICE_ID), RAWTOHEX(d.COMPANY_ID), RAWTOHEX(d.CATEGORY_ID),
	        d.TITLE, d.IS_FEATURED,
	        d.TARGET, d.CURRENT_AMOUNT, d.DONATION_DESCRIPTION, d.COVER_IMAGE,
	        d.START_DATE, d.END_DATE, d.IS_DELETED, d.ENABLED,
	        d.CREATED_AT, d.LAST_MODIFIED_AT
	      FROM DONATIONS d
	      WHERE d.START_DATE BETWEEN :1 AND :2
	        AND d.IS_DELETED = 0
	      ORDER BY d.START_DATE`

	rows, err := r.db.QueryContext(ctx, q, startDate, endDate)
	if err != nil {
		r.logger.Errorf("[DonationOracle][StreamByDateRange] query failed: %v", err)
		return local_util.HandleDBError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id, code, serviceID, companyID, categoryID, title string
			isFeatured, isDeleted, enabled                    int
			target, currentAmount, descr, cov                 sql.NullString
			startDt                                           time.Time
			endDt                                             sql.NullTime
			createdAt, lastModifiedAt                         time.Time
		)
		if err := rows.Scan(
			&id, &code, &serviceID, &companyID, &categoryID,
			&title, &isFeatured,
			&target, &currentAmount, &descr, &cov,
			&startDt, &endDt, &isDeleted, &enabled,
			&createdAt, &lastModifiedAt,
		); err != nil {
			r.logger.Errorf("[DonationOracle][StreamByDateRange] scan failed: %v", err)
			return local_util.HandleDBError(err)
		}

		d := &imodel.DonationOracle{
			ID:                  id,
			DonationCode:        code,
			ServiceID:           serviceID,
			CompanyID:           companyID,
			CategoryID:          categoryID,
			Title:               title,
			IsFeatured:          isFeatured == 1,
			Target:              target.String,
			CurrentAmount:       currentAmount.String,
			DonationDescription: descr.String,
			CoverImage:          cov.String,
			StartDate:           startDt,
			IsDeleted:           isDeleted == 1,
			Enabled:             enabled == 1,
			CreatedAt:           createdAt,
			LastModifiedAt:      lastModifiedAt,
		}
		if endDt.Valid {
			d.EndDate = endDt.Time
		}

		if err := handler(d); err != nil {
			return err
		}
	}
	if err := rows.Err(); err != nil {
		return local_util.HandleDBError(err)
	}
	return nil
}

func (r *repository) HasActiveDonationsByCategory(ctx context.Context, categoryID string) (bool, error) {
	r.logger.Infof("[DonationOracle][HasActiveDonationsByCategory] category=%s", categoryID)
	return r.hasActiveBy(ctx, "CATEGORY_ID", categoryID)
}

func (r *repository) HasActiveDonationsByCompany(ctx context.Context, companyID string) (bool, error) {
	r.logger.Infof("[DonationOracle][HasActiveDonationsByCompany] company=%s", companyID)
	return r.hasActiveBy(ctx, "COMPANY_ID", companyID)
}

func (r *repository) hasActiveBy(ctx context.Context, column, id string) (bool, error) {
	q := fmt.Sprintf(
		`SELECT COUNT(*) FROM DONATIONS WHERE %s = HEXTORAW(:1) AND IS_DELETED = 0 AND ENABLED = 1`,
		column,
	)
	var n int64
	if err := r.db.QueryRowContext(ctx, q, id).Scan(&n); err != nil {
		r.logger.Errorf("[DonationOracle][hasActiveBy %s] failed: %v", column, err)
		return false, local_util.HandleDBError(err)
	}
	return n > 0, nil
}

func (r *repository) DisableAllByCompany(ctx context.Context, companyID string) error {
	r.logger.Infof("[DonationOracle][DisableAllByCompany] company=%s", companyID)
	q := `UPDATE DONATIONS
	      SET ENABLED = 0, LAST_MODIFIED_AT = :1
	      WHERE COMPANY_ID = HEXTORAW(:2) AND IS_DELETED = 0 AND ENABLED = 1`
	if _, err := r.db.ExecContext(ctx, q, time.Now(), companyID); err != nil {
		r.logger.Errorf("[DonationOracle][DisableAllByCompany] failed: %v", err)
		return local_util.HandleDBError(err)
	}
	return nil
}
