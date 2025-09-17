package service

import (
	"context"
	"fmt"

	local_utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type serviceApplication struct {
	serviceDomain ApplicationService
	logger        utils.Logger
}

func NewServiceApplication(domain ApplicationService, logger utils.Logger) ApplicationService {
	return &serviceApplication{
		serviceDomain: domain,
		logger:        logger,
	}
}

func (sa *serviceApplication) GetAllService(ctx context.Context, filterParams local_utils.Filter) (*local_utils.PaginatedResponse[*any], error) {
	data, err := sa.serviceDomain.GetAllService(ctx, filterParams)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (sa *serviceApplication) GetAllMinimumTransferCap(ctx context.Context, filterParams local_utils.Filter) (*local_utils.PaginatedResponse[*any], error) {
	data, err := sa.serviceDomain.GetAllMinimumTransferCap(ctx, filterParams)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (sa *serviceApplication) GetAllMaximumTransferCap(ctx context.Context, filterParams local_utils.Filter) (*local_utils.PaginatedResponse[*any], error) {
	data, err := sa.serviceDomain.GetAllMaximumTransferCap(ctx, filterParams)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (sa *serviceApplication) GetAllServiceFee(ctx context.Context, filterParams local_utils.Filter) (*local_utils.PaginatedResponse[*any], error) {
	data, err := sa.serviceDomain.GetAllServiceFee(ctx, filterParams)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (sa *serviceApplication) GetAllTotalTransferCap(ctx context.Context) (*any, error) {
	data, err := sa.serviceDomain.GetAllTotalTransferCap(ctx)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (sa *serviceApplication) GetServiceFeeDetail(ctx context.Context, id string) (*any, error) {
	data, err := sa.serviceDomain.GetServiceFeeDetail(ctx, id)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (sa *serviceApplication) UpdateServiceFee(ctx context.Context, id string, req any) error {
	err := sa.serviceDomain.UpdateServiceFee(ctx, id, req)
	if err != nil {
		return err
	}
	return nil
}

func (sa *serviceApplication) UpdateSingleMaxTransfer(ctx context.Context, id string, req any) error {
	if err := sa.validateSingleMaxTransfer(ctx, req); err != nil {
		return err
	}

	err := sa.serviceDomain.UpdateSingleMaxTransfer(ctx, id, req)
	if err != nil {
		return err
	}
	return nil
}

func (sa *serviceApplication) UpdateTotalMaxTransferCap(ctx context.Context, id string, newTotalCap uint64) error {
	if err := sa.validateTotalMaxTransferCap(ctx, newTotalCap); err != nil {
		return err
	}

	err := sa.serviceDomain.UpdateTotalMaxTransferCap(ctx, id, newTotalCap)
	if err != nil {
		return err
	}
	return nil
}

func (sa *serviceApplication) UpdateMinimumTransferCap(ctx context.Context, id string, req any) error {
	if err := sa.validateMinimumTransferCap(ctx, req); err != nil {
		return err
	}

	err := sa.serviceDomain.UpdateMinimumTransferCap(ctx, id, req)
	if err != nil {
		return err
	}
	return nil
}

func (sa *serviceApplication) DeleteServiceFeeTire(ctx context.Context, id string) error {
	err := sa.serviceDomain.DeleteServiceFeeTire(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (sa *serviceApplication) validateSingleMaxTransfer(ctx context.Context, req any) error {
	totalCapData, err := sa.serviceDomain.GetAllTotalTransferCap(ctx)
	if err != nil {
		return err
	}

	totalCap, err := sa.extractTotalCap(totalCapData)
	if err != nil {
		return fmt.Errorf("failed to extract total cap: %w", err)
	}

	capData, err := local_utils.JsonUnmarshal[map[string]interface{}](req)
	if err != nil {
		return fmt.Errorf("failed to parse request: %w", err)
	}

	if err := sa.validateCapsAgainstTotal(*capData, totalCap); err != nil {
		return err
	}

	return nil
}

func (sa *serviceApplication) validateTotalMaxTransferCap(ctx context.Context, newTotalCap uint64) error {
	maxCapData, err := sa.serviceDomain.GetAllMaximumTransferCap(ctx, local_utils.Filter{Page: 1, PerPage: 1000})
	if err != nil {
		return fmt.Errorf("failed to get max caps: %w", err)
	}

	minCapData, err := sa.serviceDomain.GetAllMinimumTransferCap(ctx, local_utils.Filter{Page: 1, PerPage: 1000})
	if err != nil {
		return fmt.Errorf("failed to get min caps: %w", err)
	}

	if err := sa.validateTotalCapAgainstExisting(uint(newTotalCap), maxCapData, minCapData); err != nil {
		return err
	}

	return nil
}

func (sa *serviceApplication) validateMinimumTransferCap(ctx context.Context, req any) error {
	totalCapData, err := sa.serviceDomain.GetAllTotalTransferCap(ctx)
	if err != nil {
		return fmt.Errorf("failed to get total cap: %w", err)
	}

	totalCap, err := sa.extractTotalCap(totalCapData)
	if err != nil {
		return fmt.Errorf("failed to extract total cap: %w", err)
	}

	minReq, err := local_utils.JsonUnmarshal[map[string]interface{}](req)
	if err != nil {
		return fmt.Errorf("failed to parse request: %w", err)
	}

	minAmount, ok := (*minReq)["min_amount"].(float64)
	if !ok {
		return fmt.Errorf("invalid min_amount value")
	}

	if uint64(minAmount) > totalCap {
		return fmt.Errorf("min_amount (%d) cannot exceed total cap (%d)", uint64(minAmount), totalCap)
	}

	return nil
}

func (sa *serviceApplication) extractTotalCap(data *any) (uint64, error) {
	if data == nil {
		return 0, fmt.Errorf("NO_TOTAL_CAP_FOUND")
	}

	dataMap, ok := (*data).(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("invalid total cap data format")
	}

	totalCap, exists := dataMap["total_cap"]
	if !exists {
		return 0, fmt.Errorf("total_cap field not found")
	}

	switch v := totalCap.(type) {
	case uint64:
		return v, nil
	case uint:
		return uint64(v), nil
	case int:
		if v < 0 {
			return 0, fmt.Errorf("total cap cannot be negative")
		}
		return uint64(v), nil
	case float64:
		if v < 0 {
			return 0, fmt.Errorf("total cap cannot be negative")
		}
		return uint64(v), nil
	default:
		return 0, fmt.Errorf("unsupported total cap data type: %T", totalCap)
	}
}

func (sa *serviceApplication) validateCapsAgainstTotal(capData map[string]interface{}, totalCap uint64) error {
	if iSingleCap, exists := capData["individual_single_cap"]; exists {
		if val, ok := iSingleCap.(float64); ok && uint64(val) > totalCap {
			return fmt.Errorf("individual_single_cap (%d) cannot exceed total cap (%d)", uint64(val), totalCap)
		}
	}

	if iDailyCap, exists := capData["individual_daily_cap"]; exists {
		if val, ok := iDailyCap.(float64); ok && uint64(val) > totalCap {
			return fmt.Errorf("individual_daily_cap (%d) cannot exceed total cap (%d)", uint64(val), totalCap)
		}
	}

	if cSingleCap, exists := capData["corporate_single_cap"]; exists {
		if val, ok := cSingleCap.(float64); ok && uint64(val) > totalCap {
			return fmt.Errorf("corporate_single_cap (%d) cannot exceed total cap (%d)", uint64(val), totalCap)
		}
	}

	if cDailyCap, exists := capData["corporate_daily_cap"]; exists {
		if val, ok := cDailyCap.(float64); ok && uint64(val) > totalCap {
			return fmt.Errorf("corporate_daily_cap (%d) cannot exceed total cap (%d)", uint64(val), totalCap)
		}
	}

	return nil
}

func (sa *serviceApplication) validateTotalCapAgainstExisting(newTotalCap uint, maxCapData, minCapData *local_utils.PaginatedResponse[*any]) error {
	maxCaps := sa.extractMaxIndividualCaps(maxCapData)
	maxMinAmount := sa.extractMaxMinimumAmount(minCapData)

	if uint64(newTotalCap) <= maxCaps.ISingleCap {
		return fmt.Errorf("total_cap (%d) must be greater than max individual_single_cap (%d)", newTotalCap, maxCaps.ISingleCap)
	}

	if uint64(newTotalCap) <= maxCaps.IDailyCap {
		return fmt.Errorf("total_cap (%d) must be greater than max individual_daily_cap (%d)", newTotalCap, maxCaps.IDailyCap)
	}

	if uint64(newTotalCap) <= maxCaps.CSingleCap {
		return fmt.Errorf("total_cap (%d) must be greater than max corporate_single_cap (%d)", newTotalCap, maxCaps.CSingleCap)
	}

	if uint64(newTotalCap) <= maxCaps.CDailyCap {
		return fmt.Errorf("total_cap (%d) must be greater than max corporate_daily_cap (%d)", newTotalCap, maxCaps.CDailyCap)
	}

	if uint64(newTotalCap) <= maxMinAmount {
		return fmt.Errorf("total_cap (%d) must be greater than max min_amount (%d)", newTotalCap, maxMinAmount)
	}

	return nil
}

type MaxCaps struct {
	ISingleCap uint64
	IDailyCap  uint64
	CSingleCap uint64
	CDailyCap  uint64
}

func (sa *serviceApplication) extractMaxIndividualCaps(data *local_utils.PaginatedResponse[*any]) MaxCaps {
	var maxCaps MaxCaps

	if data == nil || data.Data == nil {
		return maxCaps
	}

	services, ok := (*data.Data).([]interface{})
	if !ok {
		return maxCaps
	}

	for _, service := range services {
		serviceMap, ok := service.(map[string]interface{})
		if !ok {
			continue
		}

		capData, exists := serviceMap["cap"]
		if !exists {
			continue
		}

		capMap, ok := capData.(map[string]interface{})
		if !ok {
			continue
		}

		if iSingleCap, exists := capMap["individual_single_cap"]; exists {
			if val, ok := iSingleCap.(float64); ok && uint64(val) > maxCaps.ISingleCap {
				maxCaps.ISingleCap = uint64(val)
			}
		}

		if iDailyCap, exists := capMap["individual_daily_cap"]; exists {
			if val, ok := iDailyCap.(float64); ok && uint64(val) > maxCaps.IDailyCap {
				maxCaps.IDailyCap = uint64(val)
			}
		}

		if cSingleCap, exists := capMap["corporate_single_cap"]; exists {
			if val, ok := cSingleCap.(float64); ok && uint64(val) > maxCaps.CSingleCap {
				maxCaps.CSingleCap = uint64(val)
			}
		}

		if cDailyCap, exists := capMap["corporate_daily_cap"]; exists {
			if val, ok := cDailyCap.(float64); ok && uint64(val) > maxCaps.CDailyCap {
				maxCaps.CDailyCap = uint64(val)
			}
		}
	}

	return maxCaps
}

func (sa *serviceApplication) extractMaxMinimumAmount(data *local_utils.PaginatedResponse[*any]) uint64 {
	var maxMinAmount uint64

	if data == nil || data.Data == nil {
		return maxMinAmount
	}

	services, ok := (*data.Data).([]interface{})
	if !ok {
		return maxMinAmount
	}

	for _, service := range services {
		serviceMap, ok := service.(map[string]interface{})
		if !ok {
			continue
		}

		capData, exists := serviceMap["cap"]
		if !exists {
			continue
		}

		capMap, ok := capData.(map[string]interface{})
		if !ok {
			continue
		}

		if minAmount, exists := capMap["min_amount"]; exists {
			if val, ok := minAmount.(float64); ok && uint64(val) > maxMinAmount {
				maxMinAmount = uint64(val)
			}
		}
	}

	return maxMinAmount
}
