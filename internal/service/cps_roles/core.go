package cpsroles

import (
	"context"
	"strings"

	cps_roles_dto "cbe-super-app-cps-action/internal/constants/dto/cps_roles"
	"cbe-super-app-cps-action/internal/constants/types"

	local_util "cbe-super-app-cps-action/pkgs/utils"
	"github.com/hugokessem/coreio/core"
)

func (r *cpsRoleService) GlobalLimits(ctx context.Context) {}

func (r *cpsRoleService) ServiceLevelLimit(ctx context.Context, serviceCode string) (*core.CustomerLimitFetchByServiceResult, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	response, err := r.core.CustomerLimitFetchByService(core.CustomerLimitFetchByServiceParam{
		ServiceCode: serviceCode,
	})
	if err != nil {
		log.Errorf("failed to make request to get service-level limit", err)
		return nil, err
	}

	if !response.Success {
		var message string
		for _, msg := range response.Message {
			message += msg
		}

		log.Warnf("(core) failed to get customer-level limit", message)

		return nil, err
	}

	if response.Detail == nil {
		log.Errorf("service-level limit successfully fetched")
		return nil, err
	}

	return response, nil
}

// ***************************************************** //
// Implementations
// ***************************************************** //

func (r *cpsRoleService) GetGlobalLimits(ctx context.Context) (*cps_roles_dto.GlobalLimitResponse, error) {
	return nil, nil
}

func (r *cpsRoleService) GetServiceLevelLimits(ctx context.Context, roleCode string, filterParam *types.Filter) (*types.PaginatedResponse[[]cps_roles_dto.ServiceLevelLimitResponse], error) {
	serviceCode := "GLOBAL" + "-" + strings.ToUpper(roleCode)
	limits, err := r.ServiceLevelLimit(ctx, serviceCode)
	if err != nil {
		return nil, err
	}

	serviceMap := make(map[string]*cps_roles_dto.ServiceLevelLimitResponse)

	if limits != nil && limits.Detail != nil && limits.Detail.GChannelType != nil {
		for _, ch := range limits.Detail.GChannelType.MChannelType {

			for _, svc := range ch.SGServiceTypes.GServiceType {

				if _, exists := serviceMap[svc.Name]; !exists {
					serviceMap[svc.Name] = &cps_roles_dto.ServiceLevelLimitResponse{
						Name: svc.Name,
					}
				}

				entry := serviceMap[svc.Name]

				switch ch.ChannelType {

				case "APP":
					entry.SuperAppMaxLimit = svc.CHANNELMAXLIMIT
					entry.SuperAppTranFreq = svc.CHANNELCOUNT

				case "USSD":
					entry.USSDMaxLimit = svc.CHANNELMAXLIMIT
					entry.USSDTranFreq = svc.CHANNELCOUNT
				}
			}
		}
	}

	results := make([]cps_roles_dto.ServiceLevelLimitResponse, 0, len(serviceMap))
	for _, v := range serviceMap {
		results = append(results, *v)
	}

	totalDocs := int64(len(results))
	page := filterParam.Page
	if page < 1 {
		page = 1
	}
	limit := filterParam.PerPage
	if limit < 1 {
		limit = 10
	}

	totalPages := int(totalDocs) / limit
	if int(totalDocs)%limit > 0 {
		totalPages++
	}

	start := (page - 1) * limit
	if start > int(totalDocs) {
		start = int(totalDocs)
	}
	end := start + limit
	if end > int(totalDocs) {
		end = int(totalDocs)
	}

	paginatedResults := results[start:end]

	meta := types.PaginationMeta{
		TotalDocs:     totalDocs,
		Limit:         limit,
		TotalPages:    totalPages,
		Page:          page,
		PagingCounter: start + 1,
		HasPrevPage:   page > 1,
		HasNextPage:   page < totalPages,
	}
	if meta.HasPrevPage {
		prev := page - 1
		meta.PrevPage = &prev
	}
	if meta.HasNextPage {
		next := page + 1
		meta.NextPage = &next
	}

	return &types.PaginatedResponse[[]cps_roles_dto.ServiceLevelLimitResponse]{
		Data: paginatedResults,
		Meta: meta,
	}, nil
}
