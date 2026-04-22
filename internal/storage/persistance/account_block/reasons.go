package account_block

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
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

func (a *AccountBlockStorage) deleteReasonsForBlockIDs(ctx context.Context, blockIDs []string) error {
	if len(blockIDs) == 0 {
		return nil
	}
	ph := make([]string, len(blockIDs))
	args := make([]interface{}, 0, len(blockIDs))
	for i, id := range blockIDs {
		n := fmt.Sprintf("bid_%d", i)
		ph[i] = ":" + n
		args = append(args, sql.Named(n, id))
	}
	q := fmt.Sprintf(`DELETE FROM account_blocks_disabled_reasons WHERE account_block_id IN (%s)`, strings.Join(ph, ","))
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
			INSERT INTO account_blocks_disabled_reasons (account_block_id, reason_text, created_by, created_at)
			VALUES (:account_block_id, :reason_text, :created_by, :created_at)`,
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
	out := make(map[string][]imodel.AccountBlockReason)
	if len(blockIDs) == 0 {
		return out, nil
	}
	ph := make([]string, len(blockIDs))
	args := make([]interface{}, 0, len(blockIDs))
	for i, id := range blockIDs {
		n := fmt.Sprintf("rid_%d", i)
		ph[i] = ":" + n
		args = append(args, sql.Named(n, id))
	}
	q := fmt.Sprintf(`
		SELECT id, account_block_id, reason_text, created_by, created_at
		FROM account_blocks_disabled_reasons
		WHERE account_block_id IN (%s)
		ORDER BY created_at ASC`, strings.Join(ph, ","))
	rows, err := a.db.QueryContext(ctx, q, args...)
	if err != nil {
		// Keep account-block APIs available even if disable-reason table
		// is not present in a given environment yet.
		if isOracleTableMissingErr(err) {
			a.logger.Warnf("[AccountBlockStorage][fetchReasonsMap] reason table missing, continuing without reasons: %v", err)
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
		bidKey := strings.TrimSpace(blockID.String)
		rid := ""
		if id.Valid {
			rid = strings.TrimSpace(id.String)
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

// attachReasons loads disable history for each block and its parent chain (same block ids as keys).
func (a *AccountBlockStorage) attachReasons(ctx context.Context, roots []*imodel.AccountBlock) error {
	if len(roots) == 0 {
		return nil
	}
	ids := make([]string, 0)
	seen := make(map[string]struct{})
	var collect func(*imodel.AccountBlock)
	collect = func(b *imodel.AccountBlock) {
		if b == nil {
			return
		}
		if _, ok := seen[b.ID]; ok {
			return
		}
		seen[b.ID] = struct{}{}
		ids = append(ids, b.ID)
		collect(b.Parent)
	}
	for _, b := range roots {
		collect(b)
	}
	m, err := a.fetchReasonsMap(ctx, ids)
	if err != nil {
		return err
	}
	var apply func(*imodel.AccountBlock)
	apply = func(b *imodel.AccountBlock) {
		if b == nil {
			return
		}
		rr := m[b.ID]
		if rr == nil {
			rr = []imodel.AccountBlockReason{}
		}
		b.DisableReason = rr
		apply(b.Parent)
	}
	for _, b := range roots {
		apply(b)
	}
	return nil
}
