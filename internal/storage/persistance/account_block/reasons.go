package account_block

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/localization"
	imodel "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/model"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/types"
	local_util "github.com/abelshebela/cbe-super-app-cps-action/pkgs/utils"
)

func isOracleTableMissingErr(err error) bool {
	if err == nil {
		return false
	}
	// Oracle: ORA-00942 table or view does not exist
	return strings.Contains(strings.ToUpper(err.Error()), "ORA-00942")
}

func bindOracleTime(t time.Time) time.Time {
	if t.IsZero() {
		return time.Now().UTC()
	}
	u := t.UTC()
	if y := u.Year(); y < 1970 || y > 9999 {
		return time.Now().UTC()
	}
	return u
}

func normalizeHexID(id string) string {
	return strings.ToUpper(strings.TrimSpace(id))
}

func (a *AccountBlockStorage) deleteReasonsForBlockIDs(ctx context.Context, blockIDs []string) error {
	if len(blockIDs) == 0 {
		return nil
	}
	ph := make([]string, len(blockIDs))
	args := make([]interface{}, 0, len(blockIDs))
	for i, id := range blockIDs {
		n := fmt.Sprintf("bid_%d", i)
		ph[i] = "HEXTORAW(:" + n + ")"
		args = append(args, sql.Named(n, id))
	}
	q := fmt.Sprintf(`DELETE FROM ACCOUNT_BLOCKS_DISABLED_REASONS WHERE ACCOUNT_BLOCK_ID IN (%s)`, strings.Join(ph, ","))
	_, err := a.db.ExecContext(ctx, q, args...)
	return err
}

func (a *AccountBlockStorage) insertDisableReasonForBlocks(ctx context.Context, blockType string, blockIDs []string, reason *types.Reason) error {
	if len(blockIDs) == 0 {
		return nil
	}

	rt := bindOracleTime(time.Now())
	rText, rBy := "", ""
	if reason != nil {
		rText, rBy = reason.Reason, reason.CreatedBy
		if !reason.CreatedAt.IsZero() {
			rt = bindOracleTime(reason.CreatedAt)
		}
	}

	for _, bid := range blockIDs {
		_, err := a.db.ExecContext(ctx, `
			INSERT INTO COMPANY_BLOCK_REASON (entity_type, entity_id, reason, created_by, created_at)
			VALUES (:entity_type, :entity_id, :reason, :created_by, :created_at)`,
			sql.Named("entity_type", blockType),
			sql.Named("entity_id", bid),
			sql.Named("reason", rText),
			sql.Named("created_by", rBy),
			sql.Named("created_at", rt),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *AccountBlockStorage) fetchReasonsMap(ctx context.Context, blockIDs []string) (map[string][]imodel.AccountBlockReason, error) {
	log := local_util.LoggerFromCtx(ctx, a.logger)

	out := make(map[string][]imodel.AccountBlockReason)
	if len(blockIDs) == 0 {
		return out, nil
	}
	ph := make([]string, len(blockIDs))
	args := make([]interface{}, 0, len(blockIDs))
	for i, id := range blockIDs {
		n := fmt.Sprintf("rid_%d", i)
		ph[i] = "HEXTORAW(:" + n + ")"
		args = append(args, sql.Named(n, id))
	}
	q := fmt.Sprintf(`
		SELECT RAWTOHEX(id), RAWTOHEX(account_block_id), reason_text, created_by, created_at
		FROM ACCOUNT_BLOCKS_DISABLED_REASONS
		WHERE account_block_id IN (%s)
		ORDER BY created_at DESC`, strings.Join(ph, ","))
	rows, err := a.db.QueryContext(ctx, q, args...)
	if err != nil {
		// Keep account-block APIs available even if disable-reason table
		// is not present in a given environment yet.
		if isOracleTableMissingErr(err) {
			log.Warnf("[AccountBlockStorage][fetchReasonsMap] reason table missing, continuing without reasons: %v", err)
			return out, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id, blockID, rText, rBy sql.NullString
		var rAt sql.NullTime
		if err := rows.Scan(&id, &blockID, &rText, &rBy, &rAt); err != nil {
			return nil, err
		}
		if !blockID.Valid {
			continue
		}
		at := time.Time{}
		if rAt.Valid {
			at = rAt.Time
		}
		bidKey := normalizeHexID(blockID.String)
		rid := ""
		if id.Valid {
			rid = normalizeHexID(id.String)
		}
		out[bidKey] = append(out[bidKey], imodel.AccountBlockReason{
			ID:        rid,
			Reason:    rText.String,
			CreatedBy: rBy.String,
			CreatedAt: at,
		})
	}
	return out, rows.Err()
}

func (a *AccountBlockStorage) GetPreviousReasons(ctx context.Context, entityType string, identifier string) ([]imodel.AccountBlockReason, error) {
	log := local_util.LoggerFromCtx(ctx, a.logger)

	entityType = strings.ToUpper(strings.TrimSpace(entityType))
	identifier = strings.TrimSpace(identifier)

	log.Infof("[AccountBlockStorage][GetPreviousReasons] entity_type=%s entity_id=%s", entityType, identifier)

	switch entityType {
	case string(imodel.TypeBranch):
		ok := local_util.IsOracleHexID(identifier)
		if identifier == "" || !ok {
			return nil, errors.New(localization.ErrorInvalidID.Code)
		}

		if _, err := a.fetchBlockByID(ctx, identifier); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, errors.New(localization.ErrorBranchNotFound.Code)
			}
			return nil, err
		}
	case string(imodel.TypeRegion):
		regions, err := a.getRegionsByNames(ctx, []string{identifier})
		if err != nil {
			return nil, err
		}
		if len(regions) == 0 {
			return nil, errors.New(localization.ErrorRegionNotFound.Code)
		}
	case string(imodel.TypeDistrict):
		districts, err := a.getDistrictsByNames(ctx, []string{identifier})
		if err != nil {
			return nil, err
		}
		if len(districts) == 0 {
			return nil, errors.New(localization.ErrorDistrictNotFound.Code)
		}
	default:
		return nil, errors.New(localization.ErrorResourceNotFound.Code)
	}

	rows, err := a.db.QueryContext(ctx, `
		SELECT RAWTOHEX(id), reason, created_by, created_at
		FROM COMPANY_BLOCK_REASON
		WHERE entity_type = :entity_type AND entity_id = :entity_id
		ORDER BY created_at DESC`,
		sql.Named("entity_type", entityType),
		sql.Named("entity_id", identifier),
	)
	if err != nil {
		if isOracleTableMissingErr(err) {
			log.Warnf("[AccountBlockStorage][GetPreviousReasons] reason table missing, returning empty list: %v", err)
			return []imodel.AccountBlockReason{}, nil
		}
		log.Errorf("[AccountBlockStorage][GetPreviousReasons] fetch failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	var reasons []imodel.AccountBlockReason
	for rows.Next() {
		var id, rText, rBy sql.NullString
		var rAt sql.NullTime
		if err := rows.Scan(&id, &rText, &rBy, &rAt); err != nil {
			return nil, err
		}
		at := time.Time{}
		if rAt.Valid {
			at = rAt.Time
		}
		rid := ""
		if id.Valid {
			rid = normalizeHexID(id.String)
		}
		reasons = append(reasons, imodel.AccountBlockReason{
			ID:        rid,
			Reason:    rText.String,
			CreatedBy: rBy.String,
			CreatedAt: at,
		})
	}
	if reasons == nil {
		reasons = []imodel.AccountBlockReason{}
	}
	return reasons, rows.Err()
}
