package account_block

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
	local_util "cbe-super-app-cps-action/pkgs/utils"
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

func (a *AccountBlockStorage) insertDisableReasonForBlocks(ctx context.Context, blockIDs []string, reason *types.Reason) error {
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
			INSERT INTO ACCOUNT_BLOCKS_DISABLED_REASONS (account_block_id, reason_text, created_by, created_at)
			VALUES (HEXTORAW(:account_block_id), :reason_text, :created_by, :created_at)`,
			sql.Named("account_block_id", bid),
			sql.Named("reason_text", rText),
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

func (a *AccountBlockStorage) GetPreviousReasons(ctx context.Context, accountBlockID string) ([]imodel.AccountBlockReason, error) {
	log := local_util.LoggerFromCtx(ctx, a.logger)

	log.Infof("[AccountBlockStorage][GetPreviousReasons] account_block_id=%s", accountBlockID)
	if _, err := a.fetchBlockByID(ctx, accountBlockID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		return nil, err
	}
	m, err := a.fetchReasonsMap(ctx, []string{accountBlockID})
	if err != nil {
		log.Errorf("[AccountBlockStorage][GetPreviousReasons] fetch failed: %v", err)
		return nil, err
	}
	reasons := m[normalizeHexID(accountBlockID)]
	if reasons == nil {
		reasons = []imodel.AccountBlockReason{}
	}
	return reasons, nil
}
