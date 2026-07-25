package transaction

import (
	transaction_dto "cbe-super-app-cps-action/internal/constants/dto/transaction"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/hugokessem/coreio/core"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type TransactionService struct {
	repo          storage.TransactionRepository
	limitRepo     storage.TransactionLimitRepository
	coreInterface core.CBECoreAPIInterface
	logger        utils.Logger
}

func NewVaultTransactionService(repo storage.TransactionRepository, limitRepo storage.TransactionLimitRepository, coreInterface core.CBECoreAPIInterface, logger utils.Logger) service.TransactionService {
	return &TransactionService{
		repo:          repo,
		limitRepo:     limitRepo,
		coreInterface: coreInterface,
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

// FetchTransactionLimitByCustomerNumber fetches real-time customer limits from CoreIO.
func (t *TransactionService) FetchTransactionLimitByCustomerNumber(ctx context.Context, customerNumber string) (*imodel.TransactionLimit, error) {
	log := local_util.LoggerFromCtx(ctx, t.logger)

	result, err := t.coreInterface.CustomerLimitFetchByCustomerNumber(ctx, core.CustomerLimitFetchByCIFParam{
		CustomerNumber: customerNumber,
	})
	if err != nil {
		log.Errorf("[TxnSvc][FetchLimit] CoreIO error for customer=%s: %v", customerNumber, err)
		return nil, err
	}

	if !result.Success {
		log.Warnf("[TxnSvc][FetchLimit] CoreIO unsuccessful for customer=%s: %v", customerNumber, result.Messages)
		return nil, errors.New(localization.ErrorResourceNotFound.Code)
	}

	if result.Detail == nil || result.Detail.GUserChannel == nil {
		log.Warnf("[TxnSvc][FetchLimit] CoreIO returned empty channel data for customer=%s", customerNumber)
		return nil, errors.New(localization.ErrorResourceNotFound.Code)
	}

	log.Infof("[TxnSvc][FetchLimit] CoreIO returned %d channels for customer=%s",
		len(result.Detail.GUserChannel.MUserChannel), customerNumber)
	return mapCIFResultToTransactionLimit(customerNumber, result), nil
}

// FetchAllTransactionLimits implements service.TransactionService.
func (t *TransactionService) FetchAllTransactionLimits(ctx context.Context, filterParams types.Filter) (*types.PaginatedResponse[[]imodel.TransactionLimit], error) {
	return t.limitRepo.FindAllWithPagination(ctx, filterParams)
}

// mapCIFResultToTransactionLimit converts a CoreIO CustomerLimitView response into
// the internal TransactionLimit model, grouping services by channel.
func mapCIFResultToTransactionLimit(customerNumber string, result *core.CustomerLimitFetchByCIFResult) *imodel.TransactionLimit {
	channels := make([]imodel.TransactionLimitChannel, 0)

	for _, ch := range result.Detail.GUserChannel.MUserChannel {
		channelName := strings.TrimSpace(ch.UserChannelType)
		if channelName == "" {
			continue
		}

		var services []imodel.TransactionLimitEntry
		if ch.SGServiceType != nil {
			for _, svc := range ch.SGServiceType.Services {
				maxAmount, _ := strconv.ParseFloat(strings.TrimSpace(svc.ServiceMaxAmt), 64)
				freq, _ := strconv.Atoi(strings.TrimSpace(svc.UserMaxCnt))
				services = append(services, imodel.TransactionLimitEntry{
					ServiceName:          strings.TrimSpace(svc.Name),
					MaximumAmount:        maxAmount,
					TransactionFrequency: freq,
				})
			}
		}

		channels = append(channels, imodel.TransactionLimitChannel{
			Channel:  channelName,
			Services: services,
		})
	}

	return &imodel.TransactionLimit{
		CustomerNumber: customerNumber,
		ChannelLimits:  channels,
		UpdatedAt:      time.Now(),
	}
}
