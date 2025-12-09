package transaction

import (
	transaction_dto "cbe-super-app-cps-action/internal/constants/dto/transaction"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type TransactionService struct {
	repo   storage.TransactionRepository
	logger utils.Logger
}

// FetchAllTransactions implements service.TransactionService.
func (t *TransactionService) FetchAllTransactions(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]transaction_dto.FullTransaction], error) {
	return t.repo.FindAllWithPagination(ctx, *filterParams)
}

// FetchTransactionByID implements service.TransactionService.
func (t *TransactionService) FetchTransactionByID(ctx context.Context, id string) (transaction_dto.FullTransaction, error) {
	return t.repo.FindTransactionByID(ctx, id)
}

func NewTransactionService(repo storage.TransactionRepository, logger utils.Logger) service.TransactionService {
	return &TransactionService{
		repo:   repo,
		logger: logger,
	}
}
