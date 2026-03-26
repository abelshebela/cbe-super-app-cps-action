package vault

import (
	"context"
	"database/sql"
	"errors"

	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	local_util "cbe-super-app-cps-action/pkgs/utils"
)

const (
	selectDeadlockRequestByID = `SELECT
		id,
		vault_name,
		vault_id,
		vault_type,
		member_name,
		member_account_number,
		status,
		created_at,
		updated_at
	FROM deadlock_requests
	WHERE id = :1`

	// countDeadlockRequests = `SELECT COUNT(*) FROM deadlock_requests`

	listDeadlockRequests = `SELECT
		id,
		vault_name,
		vault_id,
		vault_type,
		member_name,
		member_account_number,
		status,
		created_at,
		updated_at,
		COUNT(*) OVER() AS total_count
	FROM deadlock_requests
	ORDER BY created_at DESC
	OFFSET NVL(:1, 0) ROWS
	FETCH NEXT NVL(:2, 50) ROWS ONLY`

	updateDeadlockRequestStatus = `UPDATE deadlock_requests SET status = :1, updated_at = SYSTIMESTAMP WHERE id = :2`
	selectVaultIDByDeadlockReq  = `SELECT vault_id FROM deadlock_requests WHERE id = :1`
	updateVaultDeadlockFalse    = `UPDATE vault SET is_deadlocked = 0 WHERE id = :1`
)

// NOTE: Withdrawal request flow is currently disabled in this service.
// The `imodel.Withdrawal` type is not present (commented out in the model),
// so keeping this method would break compilation.
//
// func (r *VaultCategoryRepository) CreateWithdrawalRequest(ctx context.Context, withdrawal *imodel.Withdrawal) error {
// 	id := withdrawal.ID
// 	if id == "" {
// 		id = uuid.New().String()
// 	}
// 	isActive := 0
// 	if withdrawal.IsActive {
// 		isActive = 1
// 	}
// 	_, err := r.db.ExecContext(ctx, insertWithdrawal,
// 		id,
// 		withdrawal.LockedVaultID,
// 		withdrawal.Amount,
// 		withdrawal.WithdrawerName,
// 		withdrawal.WithdrawerPhoneNumber,
// 		withdrawal.Status,
// 		isActive,
// 	)
// 	if err != nil {
// 		r.logger.Errorf("failed to create withdrawal request: %v", err)
// 		return errors.New(localization.ErrorUnexpectedError.Code)
// 	}
// 	return nil
// }

func (r *VaultCategoryRepository) UpdateVaultDeadlock(ctx context.Context, id string, status string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.logger.Errorf("failed to begin tx for deadlock update: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var vaultID string
	if err := tx.QueryRowContext(ctx, selectVaultIDByDeadlockReq, id).Scan(&vaultID); err != nil {
		if err == sql.ErrNoRows {
			return errors.New(localization.ErrorResourceNotFound.Code)
		}
		r.logger.Errorf("failed to fetch vault id by deadlock request id: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	res, err := tx.ExecContext(ctx, updateDeadlockRequestStatus, status, id)
	if err != nil {
		r.logger.Errorf("failed to update deadlock request status: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}

	vaultRes, err := tx.ExecContext(ctx, updateVaultDeadlockFalse, vaultID)
	if err != nil {
		r.logger.Errorf("failed to update vault deadlock flag: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	vaultRows, _ := vaultRes.RowsAffected()
	if vaultRows == 0 {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}

	if err := tx.Commit(); err != nil {
		r.logger.Errorf("failed to commit deadlock unlock tx: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (r *VaultCategoryRepository) GetDeadlockRequest(ctx context.Context, id string) (*imodel.DeadlockRequest, error) {
	var request imodel.DeadlockRequest
	err := r.db.QueryRowContext(ctx, selectDeadlockRequestByID, id).Scan(
		&request.ID,
		&request.VaultName,
		&request.VaultID,
		&request.VaultType,
		&request.MemberName,
		&request.MemberAccountNumber,
		&request.Status,
		&request.CreatedAt,
		&request.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(localization.ErrorVaultCategoryNotFound.Code)
		}
		r.logger.Errorf("failed to get deadlock request: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return &request, nil
}

func (r *VaultCategoryRepository) GetAllDeadlockedRequests(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.DeadlockRequest], error) {
	limit := int64(50)
	page := int64(1)

	if filterParam.PerPage > 0 {
		limit = int64(filterParam.PerPage)
	}
	if filterParam.Page > 0 {
		page = int64(filterParam.Page)
	}

	offset := (page - 1) * limit

	rows, err := r.db.QueryContext(ctx, listDeadlockRequests, offset, limit)
	if err != nil {
		r.logger.Errorf("failed to list deadlock requests: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer rows.Close()

	var list []imodel.DeadlockRequest
	var total int64

	for rows.Next() {
		var request imodel.DeadlockRequest
		var count int64

		if err := rows.Scan(
			&request.ID,
			&request.VaultName,
			&request.VaultID,
			&request.VaultType,
			&request.MemberName,
			&request.MemberAccountNumber,
			&request.Status,
			&request.CreatedAt,
			&request.UpdatedAt,
			&count,
		); err != nil {
			r.logger.Errorf("failed to scan deadlock request: %v", err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}

		total = count
		list = append(list, request)
	}

	meta := local_util.BuildPaginationMeta(total, int(page), int(limit))

	return &types.PaginatedResponse[[]imodel.DeadlockRequest]{
		Data: list,
		Meta: meta,
	}, nil
}
