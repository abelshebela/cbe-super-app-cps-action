package transaction

import (
	"cbe-super-app-cps-action/internal/constants/dto/customer"
	transaction_dto "cbe-super-app-cps-action/internal/constants/dto/transaction"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hugokessem/coreio/core"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/superapp_mapper/checksum"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type TransactionService struct {
	repo          storage.TransactionRepository
	segmentRepo   storage.SuperAppRoleRepository
	limitRepo     storage.TransactionLimitRepository
	customerRepo  storage.CustomerRepository
	coreInterface core.CBECoreAPIInterface
	logger        utils.Logger
}

func NewVaultTransactionService(repo storage.TransactionRepository, segmentRepo storage.SuperAppRoleRepository, limitRepo storage.TransactionLimitRepository, customerRepo storage.CustomerRepository, coreInterface core.CBECoreAPIInterface, logger utils.Logger) service.TransactionService {
	return &TransactionService{
		repo:          repo,
		segmentRepo:   segmentRepo,
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
func (t *TransactionService) FetchTransactionLimitByCustomerNumber(ctx context.Context, filterParams *types.Filter, customerNumber string) (resp types.PaginatedResponse[[]imodel.TransactionLimitChannel], err error) {
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

	// Fetch from core because user may not be found in superapp database, but we still need limits.
	// fetch from core to get the customer group, segment, and subsegment for checksum calculation
	coreCustomerDetail, err := t.coreInterface.CustomerDetail(ctx, core.CustomerDetailParam{
		CustomerNumber: customerNumber,
	})
	if err != nil || coreCustomerDetail == nil || !coreCustomerDetail.Success || coreCustomerDetail.CustomerInfos == nil {
		log.Warnf("[TxnSvc][FetchLimit] core customer lookup failed for customer=%s: %v", customerNumber, err)
		return types.PaginatedResponse[[]imodel.TransactionLimitChannel]{
			Data: []imodel.TransactionLimitChannel{},
		}, localization.ErrorUnexpectedError
	}

	segmentLookupKey := fmt.Sprintf(
		"%s:%s:%s",
		strings.TrimSpace(coreCustomerDetail.CustomerInfos.CustomerGroup),
		coreCustomerDetail.CustomerInfos.CustomerSegment,
		coreCustomerDetail.CustomerInfos.CustomerSubSegment,
	)

	segmentChecksum := checksum.Checksum(segmentLookupKey)

	segmentConfig, err := t.segmentRepo.FindSegmentByCheckSum(ctx, segmentChecksum)
	if err != nil || segmentConfig.SuperappRole == "" {
		log.Warnf("[TxnSvc][FetchLimit] superapp segment lookup failed for key=%s checksum=%s: %+v", segmentLookupKey, segmentChecksum, err)
		return types.PaginatedResponse[[]imodel.TransactionLimitChannel]{
			Data: []imodel.TransactionLimitChannel{},
		}, localization.ErrorUnexpectedError
	}

	// Try to find the customer in superapp database to get self-activation status and expiration.
	customerInfo := customer.CustomerByCIFResponse{}
	if t.customerRepo == nil {
		log.Warnf("[TxnSvc][FetchLimit] customerRepo is nil, proceeding with segment-only defaults for customer=%s", customerNumber)
	} else {
		customerInfo, err = t.customerRepo.GetCustomerByCifNumber(ctx, customerNumber)
		if err != nil {
			log.Warnf("[TxnSvc][FetchLimit] oracle customer lookup failed for customer=%s: %+v", customerNumber, err)
		}
	}

	defaultServiceCode := resolveTransactionLimitServiceCode(segmentConfig.SuperappRole, customerInfo)

	segmentLimitResp, serviceErr := t.coreInterface.CustomerLimitFetchByService(ctx, core.CustomerLimitFetchByServiceParam{
		ServiceCode: defaultServiceCode,
	})
	if serviceErr != nil {
		log.Warnf("[TxnSvc][FetchLimit] CoreIO service lookup failed for service=%s: %v", defaultServiceCode, serviceErr)
	}

	individualLimitResp, err := t.coreInterface.CustomerLimitFetchByCustomerNumber(ctx, core.CustomerLimitFetchByCIFParam{
		CustomerNumber: customerNumber,
	})
	if err != nil {
		log.Warnf("[TxnSvc][FetchLimit] CoreIO customer lookup failed for customer=%s: %v", customerNumber, err)
	}

	baseChannels := make([]imodel.TransactionLimitChannel, 0)
	if segmentLimitResp != nil && segmentLimitResp.Success {
		baseChannels = mapServiceResultToTransactionLimitChannels(segmentLimitResp)
	} else if segmentLimitResp != nil {
		log.Warnf("[TxnSvc][FetchLimit] CoreIO service lookup returned no usable base data for service=%s: %v", defaultServiceCode, segmentLimitResp.Message)
	}

	individualChannels := make([]imodel.TransactionLimitChannel, 0)
	if individualLimitResp != nil && individualLimitResp.Success && individualLimitResp.Detail != nil && individualLimitResp.Detail.GUserChannel != nil {
		individualChannels = mapCIFResultToTransactionLimitChannels(individualLimitResp)
	} else if individualLimitResp != nil {
		log.Warnf("[TxnSvc][FetchLimit] CoreIO customer lookup returned no usable override data for customer=%s: %v", customerNumber, individualLimitResp.Messages)
	}

	mergedChannels := mergeTransactionLimitChannels(baseChannels, individualChannels)
	filteredChannels := mergedChannels
	search = strings.TrimSpace(strings.ToLower(search))
	if search != "" {
		filteredChannels = make([]imodel.TransactionLimitChannel, 0, len(mergedChannels))
		for _, channel := range mergedChannels {
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

	log.Infof("[TxnSvc][FetchLimit] merged %d channels for customer=%s", len(filteredChannels), customerNumber)
	return types.PaginatedResponse[[]imodel.TransactionLimitChannel]{
		Data: filteredChannels[start:end],
		Meta: meta,
	}, nil
}

// FetchAllTransactionLimits implements service.TransactionService.
func (t *TransactionService) FetchAllTransactionLimits(ctx context.Context, filterParams types.Filter) (*types.PaginatedResponse[[]imodel.TransactionLimit], error) {
	return t.limitRepo.FindAllWithPagination(ctx, filterParams)
}

type transactionLimitServiceRow struct {
	Channel              string
	ServiceCode          string
	ServiceName          string
	MaximumAmount        float64
	TransactionFrequency int
}

// mapServiceResultToTransactionLimitChannels converts a CoreIO service limit response into
// channel-level rows used as the default segment baseline.
func mapServiceResultToTransactionLimitChannels(result *core.CustomerLimitFetchByServiceResult) []imodel.TransactionLimitChannel {
	rows := make([]transactionLimitServiceRow, 0)
	for _, limit := range result.Detail {
		rows = append(rows, transactionLimitServiceRow{
			Channel:              strings.TrimSpace(limit.ChannelType),
			ServiceCode:          strings.TrimSpace(limit.ServiceCode),
			ServiceName:          strings.TrimSpace(limit.ServiceName),
			MaximumAmount:        parseFloatOrZero(limit.ServiceMaxAmount),
			TransactionFrequency: parseIntOrZero(limit.ServiceCount),
		})
	}
	return aggregateTransactionLimitRows(rows)
}

// mapCIFResultToTransactionLimitChannels converts a CoreIO CIF limit response into
// channel-level rows used as per-customer overrides.
func mapCIFResultToTransactionLimitChannels(result *core.CustomerLimitFetchByCIFResult) []imodel.TransactionLimitChannel {
	rows := make([]transactionLimitServiceRow, 0)

	for _, ch := range result.Detail.GUserChannel.MUserChannel {
		channelName := strings.TrimSpace(ch.UserChannelType)
		if channelName == "" {
			continue
		}

		if ch.SGServiceType != nil {
			for _, svc := range ch.SGServiceType.Services {
				rows = append(rows, transactionLimitServiceRow{
					Channel:              channelName,
					ServiceCode:          strings.TrimSpace(svc.Name),
					ServiceName:          strings.TrimSpace(svc.Name),
					MaximumAmount:        parseFloatOrZero(svc.ServiceMaxAmt),
					TransactionFrequency: parseIntOrZero(svc.UserMaxCnt),
				})
			}
		}
	}

	return aggregateTransactionLimitRows(rows)
}

func mergeTransactionLimitChannels(baseChannels, overrideChannels []imodel.TransactionLimitChannel) []imodel.TransactionLimitChannel {
	channelIndex := make(map[string]int)
	merged := make([]imodel.TransactionLimitChannel, 0, len(baseChannels))

	addChannel := func(channel imodel.TransactionLimitChannel) int {
		channelKey := strings.ToUpper(strings.TrimSpace(channel.Channel))
		if idx, exists := channelIndex[channelKey]; exists {
			return idx
		}
		merged = append(merged, imodel.TransactionLimitChannel{Channel: channel.Channel, Services: make([]imodel.TransactionLimitEntry, 0, len(channel.Services))})
		channelIndex[channelKey] = len(merged) - 1
		return len(merged) - 1
	}

	serviceIndex := func(channel imodel.TransactionLimitChannel) map[string]int {
		index := make(map[string]int, len(channel.Services))
		for i, service := range channel.Services {
			key := strings.ToUpper(strings.TrimSpace(service.ServiceName))
			if key != "" {
				index[key] = i
			}
		}
		return index
	}

	for _, channel := range baseChannels {
		idx := addChannel(channel)
		merged[idx].Channel = channel.Channel
		merged[idx].Services = append(merged[idx].Services, channel.Services...)
	}

	for _, channel := range overrideChannels {
		idx := addChannel(channel)
		if merged[idx].Channel == "" {
			merged[idx].Channel = channel.Channel
		}
		currentIndex := serviceIndex(merged[idx])
		for _, service := range channel.Services {
			key := strings.ToUpper(strings.TrimSpace(service.ServiceName))
			if key == "" {
				continue
			}
			if existingIdx, exists := currentIndex[key]; exists {
				merged[idx].Services[existingIdx] = service
				continue
			}
			merged[idx].Services = append(merged[idx].Services, service)
			currentIndex[key] = len(merged[idx].Services) - 1
		}
	}

	return merged
}

func aggregateTransactionLimitRows(rows []transactionLimitServiceRow) []imodel.TransactionLimitChannel {
	channelIndex := make(map[string]int)
	serviceIndexByChannel := make(map[string]map[string]int)
	result := make([]imodel.TransactionLimitChannel, 0)

	for _, row := range rows {
		channelKey := strings.ToUpper(strings.TrimSpace(row.Channel))
		serviceKey := strings.ToUpper(strings.TrimSpace(row.ServiceCode))
		if serviceKey == "" {
			serviceKey = strings.ToUpper(strings.TrimSpace(row.ServiceName))
		}
		if channelKey == "" || serviceKey == "" {
			continue
		}

		channelIdx, exists := channelIndex[channelKey]
		if !exists {
			result = append(result, imodel.TransactionLimitChannel{Channel: row.Channel, Services: []imodel.TransactionLimitEntry{}})
			channelIdx = len(result) - 1
			channelIndex[channelKey] = channelIdx
			serviceIndexByChannel[channelKey] = make(map[string]int)
		}

		serviceIdx, exists := serviceIndexByChannel[channelKey][serviceKey]
		entry := imodel.TransactionLimitEntry{
			ServiceName:          row.ServiceName,
			MaximumAmount:        row.MaximumAmount,
			TransactionFrequency: row.TransactionFrequency,
		}
		if exists {
			result[channelIdx].Services[serviceIdx] = entry
			continue
		}

		result[channelIdx].Services = append(result[channelIdx].Services, entry)
		serviceIndexByChannel[channelKey][serviceKey] = len(result[channelIdx].Services) - 1
	}

	return result
}

func parseFloatOrZero(value string) float64 {
	parsed, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
	if err != nil {
		return 0
	}
	return parsed
}

func parseIntOrZero(value string) int {
	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return 0
	}
	return parsed
}

func resolveTransactionLimitServiceCode(serviceCode string, customerInfo customer.CustomerByCIFResponse) string {
	fmt.Printf("[TxnSvc][ResolveServiceCode] Resolving service code for initial serviceCode=%s, customerInfo=%+v\n", serviceCode, customerInfo)
	if serviceCode == "" {
		serviceCode = strings.TrimSpace(customerInfo.CustomerSegmentation)
	}
	serviceCode = strings.ToUpper(serviceCode)

	if customerInfo.IsSelfActivated && customerInfo.SelfActivationRuleExpiration != nil && time.Now().Before(*customerInfo.SelfActivationRuleExpiration) {
		if strings.Contains(strings.ToUpper(serviceCode), "MASS") {
			return "GLOBAL-ONLINE.CONV"
		}
		return "GLOBAL-ONLINE.IFB"
	}

	if serviceCode == "" {
		return "GLOBAL-"
	}
	return "GLOBAL-" + serviceCode
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
