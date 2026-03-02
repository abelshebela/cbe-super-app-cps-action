package vault

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	local_util "cbe-super-app-cps-action/pkgs/utils"
)

const (
	insertWithdrawal = `INSERT INTO withdrawals (
		id, locked_vault_id, amount, withdrawer_name, withdrawer_phone_number,
		status, is_active, created_at, updated_at
	) VALUES (
		:1, :2, :3, :4, :5, :6, :7, SYSTIMESTAMP, SYSTIMESTAMP
	)`
	updateWithdrawalStatus = `UPDATE withdrawals SET status = :1, updated_at = SYSTIMESTAMP WHERE id = :2`
	selectWithdrawalByID   = `SELECT id, locked_vault_id, amount, withdrawer_name, withdrawer_phone_number, status, is_active, created_at, updated_at FROM withdrawals WHERE id = :1`
	countWithdrawals       = `SELECT COUNT(*) FROM withdrawals WHERE is_active = 1`
	listWithdrawals        = `SELECT id, locked_vault_id, amount, withdrawer_name, withdrawer_phone_number, status, is_active, created_at, updated_at
		FROM withdrawals WHERE is_active = 1 ORDER BY created_at DESC OFFSET :offset ROWS FETCH NEXT :limit ROWS ONLY`
)

func (r *VaultCategoryRepository) CreateWithdrawalRequest(ctx context.Context, withdrawal *imodel.Withdrawal) error {
	id := withdrawal.ID
	if id == "" {
		id = uuid.New().String()
	}
	isActive := 0
	if withdrawal.IsActive {
		isActive = 1
	}
	_, err := r.db.ExecContext(ctx, insertWithdrawal,
		id,
		withdrawal.LockedVaultID,
		withdrawal.Amount,
		withdrawal.WithdrawerName,
		withdrawal.WithdrawerPhoneNumber,
		withdrawal.Status,
		isActive,
	)
	if err != nil {
		r.logger.Errorf("failed to create withdrawal request: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (r *VaultCategoryRepository) UpdateWithdrawalRequest(ctx context.Context, id string, status string) error {
	res, err := r.db.ExecContext(ctx, updateWithdrawalStatus, status, id)
	if err != nil {
		r.logger.Errorf("failed to update withdrawal request: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}

	// updatedRequest, err := r.GetWithdrawalRequest(ctx, id)
	// if err != nil {
	// 	r.logger.Errorf("failed to fetch updated withdrawal request: %v", err)
	// 	return errors.New(localization.ErrorUnexpectedError.Code)
	// }

	// r.kafkaProducer.PublishMessage(
	// 	ctx,
	// 	updatedRequest,
	// 	string(constants.ClientOrchestrationServicesTopic),
	// 	r.cfg.CPSServiceUpdate,
	// 	"vault withdrawal request authorized",
	// )

	return nil
}

func (r *VaultCategoryRepository) GetWithdrawalRequest(ctx context.Context, id string) (*imodel.Withdrawal, error) {
	var w imodel.Withdrawal
	var isActive int
	err := r.db.QueryRowContext(ctx, selectWithdrawalByID, id).Scan(
		&w.ID,
		&w.LockedVaultID,
		&w.Amount,
		&w.WithdrawerName,
		&w.WithdrawerPhoneNumber,
		&w.Status,
		&isActive,
		&w.CreatedAt,
		&w.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(localization.ErrorVaultCategoryNotFound.Code)
		}
		r.logger.Errorf("failed to get withdrawal request: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return &w, nil
}

func (r *VaultCategoryRepository) GetAllWithdrawalRequests(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.Withdrawal], error) {
	limit := int64(50)
	page := int64(1)
	if filterParam.PerPage > 0 {
		limit = int64(filterParam.PerPage)
	}
	if filterParam.Page > 0 {
		page = int64(filterParam.Page)
	}
	offset := (page - 1) * limit

	var total int64
	if err := r.db.QueryRowContext(ctx, countWithdrawals).Scan(&total); err != nil {
		r.logger.Errorf("failed to count withdrawal requests: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	rows, err := r.db.QueryContext(ctx, listWithdrawals, sql.Named("offset", offset), sql.Named("limit", limit))
	if err != nil {
		r.logger.Errorf("failed to list withdrawal requests: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer rows.Close()

	var list []imodel.Withdrawal
	for rows.Next() {
		var w imodel.Withdrawal
		var isActive int
		if err := rows.Scan(
			&w.ID,
			&w.LockedVaultID,
			&w.Amount,
			&w.WithdrawerName,
			&w.WithdrawerPhoneNumber,
			&w.Status,
			&isActive,
			&w.CreatedAt,
			&w.UpdatedAt,
		); err != nil {
			r.logger.Errorf("failed to scan withdrawal request: %v", err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		list = append(list, w)
	}

	meta := local_util.BuildPaginationMeta(total, int(page), int(limit))
	return &types.PaginatedResponse[[]imodel.Withdrawal]{
		Data: list,
		Meta: meta,
	}, nil
}
