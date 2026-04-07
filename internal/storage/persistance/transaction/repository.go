package transaction_repo

import (
	transaction_dto "cbe-super-app-cps-action/internal/constants/dto/transaction"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	transaction_core "cbe-super-app-cps-action/internal/storage/persistance/transaction/core"
	sqlc "cbe-super-app-cps-action/internal/storage/persistance/transaction/sql"
	"context"
	"database/sql"
	"errors"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type TransactionRepository struct {
	db     *sql.DB
	logger utils.Logger
}

// FindTransactionByCifOrAccountNumberOrFT implements storage.TransactionRepository.
func (t *TransactionRepository) FindTransactionByCifOrAccountNumberOrFT(ctx context.Context, identifier string) (transaction_dto.FullTransaction, error) {
	q := sqlc.New(t.db)
	transaction, err := q.FindTransactionByCifOrAccountNumberOrFT(ctx, identifier)
	if err != nil {
		t.logger.Errorf("failed to find transaction by CIF/AccountNumber/FT: %v", err)
		return transaction_dto.FullTransaction{}, errors.New(localization.ErrorResourceNotFound.Code)
	}
	return transaction_core.MapFullTransactionToNative(transaction), nil
}

// FindAllWithPagination implements storage.TransactionRepository.
func (t *TransactionRepository) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]transaction_dto.FullTransaction], error) {
	q := sqlc.New(t.db)
	transactionsData, err := q.FindTransactionWithParam(ctx, &filterParam)
	if err != nil {
		t.logger.Errorf("failed to find transactions with param: %v", err)
		return nil, errors.New(localization.ErrorResourceNotFound.Code)
	}
	return transactionsData, nil
}

// FindTransactionByID implements storage.TransactionRepository.
func (t *TransactionRepository) FindTransactionByID(ctx context.Context, id string) (transaction_dto.FullTransaction, error) {
	q := sqlc.New(t.db)
	transaction, err := q.FindTransactionByID(ctx, id)
	if err != nil {
		t.logger.Errorf("failed to find transaction by id: %v", err)
		return transaction_dto.FullTransaction{}, errors.New(localization.ErrorResourceNotFound.Code)
	}
	return transaction_core.MapFullTransactionToNative(transaction), nil
}

func NewTransactionRepository(db *sql.DB, logger utils.Logger) storage.TransactionRepository {
	return &TransactionRepository{
		db:     db,
		logger: logger,
	}
}
