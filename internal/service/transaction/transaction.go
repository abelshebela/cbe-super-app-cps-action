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
	customerRepo  storage.CustomerRepository
	coreInterface core.CBECoreAPIInterface
	logger        utils.Logger
}

func NewVaultTransactionService(repo storage.TransactionRepository, limitRepo storage.TransactionLimitRepository, customerRepo storage.CustomerRepository, coreInterface core.CBECoreAPIInterface, logger utils.Logger) service.TransactionService {
	return &TransactionService{
		repo:          repo,
		limitRepo:     limitRepo,
		customerRepo:  customerRepo,
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
// Results are normalized into channel rows and paginated in the same style as customer service lookups.
func (t *TransactionService) FetchTransactionLimitByCustomerNumber(ctx context.Context, filterParams *types.Filter, customerNumber string) (types.PaginatedResponse[[]imodel.TransactionLimitChannel], error) {
	log := local_util.LoggerFromCtx(ctx, t.logger)
	log.Infof("[TxnSvc][FetchLimit] Fetching transaction limits for customer=%s", customerNumber)
	page := 1
	perPage := 10
	search := ""
	if filterParams != nil {
		page = filterParams.Page
		perPage = filterParams.PerPage
		search = filterParams.Search
	}
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}
	customerInfo, err := t.customerRepo.GetCustomerByCifNumber(ctx, customerNumber)
	if err != nil {
		log.Warnf("[TxnSvc][FetchLimit] oracle customer lookup failed for customer=%s: %v", customerNumber, err)
		return types.PaginatedResponse[[]imodel.TransactionLimitChannel]{}, err
	}
	selfActiveTime := *customerInfo.SelfActivationRuleExpiration
	serviceCode := "GLOBAL-"
	if customerInfo.IsSelfActivated && time.Now().Before(selfActiveTime) {
		if strings.Contains(customerInfo.SuperappRole, "MASS") {
			serviceCode += "ONLINE.CONV"
		} else {
			serviceCode += "ONLINE.IFB"
		}
	} else {
		serviceCode += customerInfo.SuperappRole
	}

	serviceRes, err := t.coreInterface.CustomerLimitFetchByService(ctx, core.CustomerLimitFetchByServiceParam{
		ServiceCode: customerInfo.SuperappRole,
	})
	if err != nil {
		log.Errorf("[TxnSvc][FetchLimit] CoreIO error for service=%s: %v", serviceCode, err)
		return types.PaginatedResponse[[]imodel.TransactionLimitChannel]{}, localization.ErrorUnexpectedError
	}

	result, err := t.coreInterface.CustomerLimitFetchByCustomerNumber(ctx, core.CustomerLimitFetchByCIFParam{
		CustomerNumber: customerNumber,
	})
	if err != nil {
		log.Errorf("[TxnSvc][FetchLimit] CoreIO error for customer=%s: %v", customerNumber, err)
		return types.PaginatedResponse[[]imodel.TransactionLimitChannel]{}, err
	}
	if result == nil {
		log.Warnf("[TxnSvc][FetchLimit] CoreIO returned nil response for customer=%s", customerNumber)
		return types.PaginatedResponse[[]imodel.TransactionLimitChannel]{
			Data: []imodel.TransactionLimitChannel{},
			Meta: local_util.BuildPaginationMeta(0, page, perPage),
		}, nil
	}

	if !result.Success {
		log.Warnf("[TxnSvc][FetchLimit] CoreIO unsuccessful for customer=%s: %v", customerNumber, result.Messages)
		return types.PaginatedResponse[[]imodel.TransactionLimitChannel]{}, errors.New(localization.ErrorResourceNotFound.Code)
	}

	if result.Detail == nil || result.Detail.GUserChannel == nil {
		log.Warnf("[TxnSvc][FetchLimit] CoreIO returned empty channel data for customer=%s", customerNumber)
		return types.PaginatedResponse[[]imodel.TransactionLimitChannel]{
			Data: []imodel.TransactionLimitChannel{},
			Meta: local_util.BuildPaginationMeta(0, page, perPage),
		}, nil
	}

	channels := mapCIFResultToTransactionLimitChannels(result)
	filteredChannels := channels
	search = strings.TrimSpace(strings.ToLower(search))
	if search != "" {
		filteredChannels = make([]imodel.TransactionLimitChannel, 0, len(channels))
		for _, channel := range channels {
			if transactionLimitChannelMatchesSearch(channel, search) {
				filteredChannels = append(filteredChannels, channel)
			}
		}
	}

	totalDocs := int64(len(filteredChannels))
	meta := local_util.BuildPaginationMeta(totalDocs, page, perPage)

	start := (page - 1) * perPage
	if start >= len(filteredChannels) {
		return types.PaginatedResponse[[]imodel.TransactionLimitChannel]{
			Data: []imodel.TransactionLimitChannel{},
			Meta: meta,
		}, nil
	}

	end := start + perPage
	if end > len(filteredChannels) {
		end = len(filteredChannels)
	}

	log.Infof("[TxnSvc][FetchLimit] CoreIO returned %d channels for customer=%s", len(filteredChannels), customerNumber)
	return types.PaginatedResponse[[]imodel.TransactionLimitChannel]{
		Data: filteredChannels[start:end],
		Meta: meta,
	}, nil
}

// FetchAllTransactionLimits implements service.TransactionService.
func (t *TransactionService) FetchAllTransactionLimits(ctx context.Context, filterParams types.Filter) (*types.PaginatedResponse[[]imodel.TransactionLimit], error) {
	return t.limitRepo.FindAllWithPagination(ctx, filterParams)
}

// mapCIFResultToTransactionLimitChannels converts a CoreIO CustomerLimitView response into
// channel-level transaction limit rows.
func mapCIFResultToTransactionLimitChannels(result *core.CustomerLimitFetchByCIFResult) []imodel.TransactionLimitChannel {
	channels := make([]imodel.TransactionLimitChannel, 0)

	for _, ch := range result.Detail.GUserChannel.MUserChannel {
		channelName := strings.TrimSpace(ch.UserChannelType)
		if channelName == "" {
			continue
		}

		services := make([]imodel.TransactionLimitEntry, 0)
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

	return channels
}

func transactionLimitChannelMatchesSearch(channel imodel.TransactionLimitChannel, search string) bool {
	if strings.Contains(strings.ToLower(channel.Channel), search) {
		return true
	}

	for _, service := range channel.Services {
		if strings.Contains(strings.ToLower(service.ServiceName), search) {
			return true
		}
	}

	return false
}
