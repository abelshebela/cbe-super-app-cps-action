package term_and_condition_oracle

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	tac_core "cbe-super-app-cps-action/internal/storage/persistance/term_and_condition/oracle/core"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type repository struct {
	db     DBTX
	logger utils.Logger
}

var _ storage.AccountOpeningTermsRepository = (*repository)(nil)

func NewRepository(db DBTX, logger utils.Logger) storage.AccountOpeningTermsRepository {
	return &repository{db: db, logger: logger}
}

func (r *repository) Create(ctx context.Context, t *imodel.AccountOpeningTerms) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)
	log.Infof("[TACOracle][Create] product=%s version=%s", t.AccountProductID, t.VersionLabel)

	activationTime, err := tac_core.ParseActivationTime(t.ActivationTime)
	if err != nil {
		log.Errorf("[TACOracle][Create] parse activation_time: %v", err)
		return err
	}

	q := `INSERT INTO ACCOUNT_OPENING_TERMS
		(ACCOUNT_PRODUCT_ID, ACTIVATION_TIME, VERSION_LABEL, TERMS_AND_CONDITIONS_PATH, IS_ENABLED, IS_DELETED)
		VALUES (HEXTORAW(:1), :2, :3, :4, :5, :6)`

	if _, err := r.db.ExecContext(ctx, q,
		t.AccountProductID,
		activationTime,
		t.VersionLabel,
		t.TermsAndConditionsPath,
		tac_core.BoolToInt(t.IsEnabled),
		tac_core.BoolToInt(t.IsDeleted),
	); err != nil {
		log.Errorf("[TACOracle][Create] insert: %v", err)
		return local_util.HandleDBError(err)
	}
	return nil
}

func (r *repository) Update(ctx context.Context, id string, t *imodel.AccountOpeningTerms) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)
	log.Infof("[TACOracle][Update] id=%s", id)

	activationTime, err := tac_core.ParseActivationTime(t.ActivationTime)
	if err != nil {
		log.Errorf("[TACOracle][Update] parse activation_time: %v", err)
		return err
	}

	q := `UPDATE ACCOUNT_OPENING_TERMS
		SET ACTIVATION_TIME = :1, VERSION_LABEL = :2,
		    TERMS_AND_CONDITIONS_PATH = :3, LAST_MODIFIED_AT = :4
		WHERE ID = HEXTORAW(:5)`

	res, err := r.db.ExecContext(ctx, q,
		activationTime,
		t.VersionLabel,
		t.TermsAndConditionsPath,
		time.Now(),
		id,
	)
	if err != nil {
		log.Errorf("[TACOracle][Update] exec: %v", err)
		return local_util.HandleDBError(err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)
	log.Infof("[TACOracle][Delete] id=%s", id)

	res, err := r.db.ExecContext(ctx,
		`DELETE FROM ACCOUNT_OPENING_TERMS WHERE ID = HEXTORAW(:1)`,
		id,
	)
	if err != nil {
		log.Errorf("[TACOracle][Delete] exec: %v", err)
		return local_util.HandleDBError(err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}
	return nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*imodel.AccountOpeningTerms, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)
	log.Infof("[TACOracle][FindByID] id=%s", id)

	q := `SELECT ` + tacSelectCols + tacFromTable + `WHERE t.ID = HEXTORAW(:1) AND t.IS_DELETED = 0`
	row := r.db.QueryRowContext(ctx, q, id)
	out, err := scanTACRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		log.Errorf("[TACOracle][FindByID] scan: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return out, nil
}

func (r *repository) FindByProductAndVersion(ctx context.Context, productID, versionLabel string) (*imodel.AccountOpeningTerms, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)
	log.Infof("[TACOracle][FindByProductAndVersion] product=%s version=%s", productID, versionLabel)

	q := `SELECT ` + tacSelectCols + tacFromTable +
		`WHERE t.ACCOUNT_PRODUCT_ID = HEXTORAW(:1) AND t.VERSION_LABEL = :2 AND t.IS_DELETED = 0`
	row := r.db.QueryRowContext(ctx, q, productID, versionLabel)
	out, err := scanTACRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		log.Errorf("[TACOracle][FindByProductAndVersion] scan: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return out, nil
}

func (r *repository) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.AccountOpeningTerms], error) {
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

	conds := []string{"t.IS_DELETED = 0"}
	args := []interface{}{}
	idx := 1

	if s := strings.TrimSpace(filterParam.Search); s != "" {
		conds = append(conds, fmt.Sprintf(
			"(UPPER(t.VERSION_LABEL) LIKE UPPER(:%d) OR TO_CHAR(t.ACTIVATION_TIME, 'YYYY-MM-DD HH24:MI:SS') LIKE :%d)",
			idx, idx+1,
		))
		args = append(args, "%"+s+"%", "%"+s+"%")
		idx += 2
	}

	if filterParam.Filters != nil {
		if v, ok := filterParam.Filters["account_product_id"]; ok {
			if s, _ := v.(string); strings.TrimSpace(s) != "" {
				conds = append(conds, fmt.Sprintf("t.ACCOUNT_PRODUCT_ID = HEXTORAW(:%d)", idx))
				args = append(args, s)
				idx++
			}
		}
		if v, ok := filterParam.Filters["version_label"]; ok {
			if s, _ := v.(string); strings.TrimSpace(s) != "" {
				conds = append(conds, fmt.Sprintf("UPPER(t.VERSION_LABEL) = UPPER(:%d)", idx))
				args = append(args, s)
				idx++
			}
		}
	}

	where := strings.Join(conds, " AND ")

	var total int64
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM ACCOUNT_OPENING_TERMS t WHERE `+where, args...).Scan(&total); err != nil {
		log.Errorf("[TACOracle][FindAll] count: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	meta := local_util.BuildPaginationMeta(total, page, perPage)
	if total == 0 {
		return &types.PaginatedResponse[[]imodel.AccountOpeningTerms]{Data: []imodel.AccountOpeningTerms{}, Meta: meta}, nil
	}

	selectQ := `SELECT ` + tacSelectCols + tacFromTable + `WHERE ` + where +
		` ORDER BY t.CREATED_AT DESC OFFSET :` + fmt.Sprintf("%d", idx) +
		` ROWS FETCH NEXT :` + fmt.Sprintf("%d", idx+1) + ` ROWS ONLY`
	args = append(args, offset, perPage)

	rows, err := r.db.QueryContext(ctx, selectQ, args...)
	if err != nil {
		log.Errorf("[TACOracle][FindAll] query: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer rows.Close()

	out := []imodel.AccountOpeningTerms{}
	for rows.Next() {
		item, err := scanTACRow(rows)
		if err != nil {
			log.Errorf("[TACOracle][FindAll] scan: %v", err)
			return nil, local_util.HandleDBError(err)
		}
		out = append(out, *item)
	}
	if err := rows.Err(); err != nil {
		return nil, local_util.HandleDBError(err)
	}

	return &types.PaginatedResponse[[]imodel.AccountOpeningTerms]{Data: out, Meta: meta}, nil
}
