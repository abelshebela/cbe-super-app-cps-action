package transaction

import (
	transaction_dto "cbe-super-app-cps-action/internal/constants/dto/transaction"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type TransactionService struct {
	repo      storage.TransactionRepository
	limitRepo storage.TransactionLimitRepository
	logger    utils.Logger
}

func NewVaultTransactionService(repo storage.TransactionRepository, limitRepo storage.TransactionLimitRepository, logger utils.Logger) service.TransactionService {
	return &TransactionService{
		repo:      repo,
		limitRepo: limitRepo,
		logger:    logger,
	}
}

// FindTransactionByCifOrAccountNumberOrFT implements service.TransactionService.
func (t *TransactionService) FindTransactionByCifOrAccountNumberOrFT(ctx context.Context, identifier string) (transaction_dto.VaultTransaction, error) {
	return t.repo.FindTransactionByCifOrAccountNumberOrFT(ctx, identifier)
}

// FetchAllTransactions implements service.TransactionService.
func (t *TransactionService) FetchAllTransactions(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]transaction_dto.VaultTransaction], error) {
	return t.repo.FindAllWithPagination(ctx, *filterParams)
}

// FetchTransactionByID implements service.TransactionService.
func (t *TransactionService) FetchTransactionByID(ctx context.Context, id string) (transaction_dto.VaultTransaction, error) {
	return t.repo.FindTransactionByID(ctx, id)
}

// FetchTransactionLimitByCustomerNumber implements service.TransactionService.
func (t *TransactionService) FetchTransactionLimitByCustomerNumber(ctx context.Context, customerNumber string) (*imodel.TransactionLimit, error) {
	return t.limitRepo.FindByCustomerNumber(ctx, customerNumber)
}

// FetchAllTransactionLimits implements service.TransactionService.
func (t *TransactionService) FetchAllTransactionLimits(ctx context.Context, filterParams types.Filter) (*types.PaginatedResponse[[]imodel.TransactionLimit], error) {
	return t.limitRepo.FindAllWithPagination(ctx, filterParams)
}
