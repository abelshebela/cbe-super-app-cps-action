package transaction

import (
	transaction_dto "cbe-super-app-cps-action/internal/constants/dto/transaction"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/constants/localization"
	account_lookup "cbe-super-app-cps-action/internal/storage/external_call/account_lookup"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/hugokessem/coreio/core"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type TransactionService struct {
	repo          storage.TransactionRepository
	limitRepo     storage.TransactionLimitRepository
	accountLookup account_lookup.Account
	logger        utils.Logger
}

func NewVaultTransactionService(repo storage.TransactionRepository, limitRepo storage.TransactionLimitRepository, accountLookup account_lookup.Account, logger utils.Logger) service.TransactionService {
	return &TransactionService{
		repo:          repo,
		limitRepo:     limitRepo,
		accountLookup: accountLookup,
		logger:        logger,
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

// FetchTransactionLimitByCustomerNumber fetches real-time limits from CoreIO.
func (t *TransactionService) FetchTransactionLimitByCustomerNumber(ctx context.Context, customerNumber string) (*imodel.TransactionLimit, error) {
	log := local_util.LoggerFromCtx(ctx, t.logger)

	result, err := t.accountLookup.SearchCustomerServiceLimitByCIF(ctx, customerNumber)
	if err != nil {
		log.Errorf("[TxnSvc][FetchLimit] CoreIO error for customer=%s: %v", customerNumber, err)
		return nil, err
	}
	if !result.Success || len(result.Detail) == 0 {
		log.Warnf("[TxnSvc][FetchLimit] CoreIO returned no records for customer=%s: %v", customerNumber, result.Messages)
		return nil, errors.New(localization.ErrorResourceNotFound.Code)
	}

	log.Infof("[TxnSvc][FetchLimit] CoreIO returned %d entries for customer=%s", len(result.Detail), customerNumber)
	return mapCoreIOToTransactionLimit(customerNumber, result), nil
}

// FetchAllTransactionLimits implements service.TransactionService.
func (t *TransactionService) FetchAllTransactionLimits(ctx context.Context, filterParams types.Filter) (*types.PaginatedResponse[[]imodel.TransactionLimit], error) {
	return t.limitRepo.FindAllWithPagination(ctx, filterParams)
}

// mapCoreIOToTransactionLimit groups CoreIO detail rows into the internal TransactionLimit model.
func mapCoreIOToTransactionLimit(customerNumber string, result *core.CustomerLimitFetchByCIFReturnServiceResult) *imodel.TransactionLimit {
	channelMap := make(map[string][]imodel.TransactionLimitEntry)
	for _, d := range result.Detail {
		maxAmount, _ := strconv.ParseFloat(d.Limit, 64)
		freq, _ := strconv.Atoi(d.Count)
		channelMap[d.Channel] = append(channelMap[d.Channel], imodel.TransactionLimitEntry{
			ServiceName:          d.ServiceName,
			MaximumAmount:        maxAmount,
			TransactionFrequency: freq,
		})
	}

	channels := make([]imodel.TransactionLimitChannel, 0, len(channelMap))
	for ch, services := range channelMap {
		channels = append(channels, imodel.TransactionLimitChannel{
			Channel:  ch,
			Services: services,
		})
	}

	return &imodel.TransactionLimit{
		CustomerNumber: customerNumber,
		ChannelLimits:  channels,
		UpdatedAt:      time.Now(),
	}
}
