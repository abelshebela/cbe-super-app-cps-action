package service_details

import (
	dto "cbe-super-app-cps-action/internal/constants/dto/service_details"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func ValidateTotalCapAgainstServices(ctx context.Context, serviceRepo storage.ServiceDetailsRepository, logger utils.Logger, newTotalCap uint64) error {
	logger.Infof("Validating new total cap %d against all service caps", newTotalCap)

	// Get all services with their caps
	projection := bson.M{
		"service_name": 1,
		"service_code": 1,
		"cap":          1,
	}

	services, err := serviceRepo.FindAllWithPagination(ctx, projection, types.Filter{Page: 1, PerPage: 1000})
	if err != nil {
		logger.Errorf("error fetching services for total cap validation: %v", err)
		return err
	}

	for _, service := range services.Data {
		if service.Cap.ISingleCap > newTotalCap ||
			service.Cap.IDailyCap > newTotalCap ||
			service.Cap.CorporateSingleCap > newTotalCap ||
			service.Cap.CorporateDailyCap > newTotalCap {

			logger.Warnf(
				"Service %s (ID: %s) has a cap exceeding new total cap %d",
				service.ServiceName, service.ID.Hex(), newTotalCap,
			)
			return errors.New(localization.ErrorTotalMaxTransferCannotBeLessExistTransfers.Code)
		}
	}

	logger.Infof("All service caps are within the new total cap %d", newTotalCap)
	return nil
}

func MapServiceDetailsForUpdate(existingService *model.ServiceDetails, newData interface{}) (*model.ServiceDetails, error) {
	updatedService := *existingService

	// Try to directly cast to the expected DTO type first
	if dto, ok := newData.(dto.ServiceFeeDetailDTO); ok {

		if dto.ServiceType != "" {
			updatedService.ServiceType = dto.ServiceType
		}
		if dto.PaymentType != "" {
			updatedService.PaymentType = string(dto.PaymentType)
		}
		if dto.AboveServiceFee > 0 {
			updatedService.AboveServiceFee = dto.AboveServiceFee
		}

		// Map GL entries from DTO
		if dto.CBglEntry.CBglProductAccount != "" {
			updatedService.CBEGLEntry.ProductAccount = dto.CBglEntry.CBglProductAccount
		}
		if dto.CBglEntry.CBglProductBranchcode != "" {
			updatedService.CBEGLEntry.ProductBranchCode = dto.CBglEntry.CBglProductBranchcode
		}
		if dto.CBglEntry.CBglServiceAccount != "" {
			updatedService.CBEGLEntry.ServiceAccount = dto.CBglEntry.CBglServiceAccount
		}
		if dto.CBglEntry.CBglServiceBranchcode != "" {
			updatedService.CBEGLEntry.ServiceBranchCode = dto.CBglEntry.CBglServiceBranchcode
		}

		if dto.IFBglEntry.IFBglProductAccount != "" {
			updatedService.CBEIFBGLEntry.ProductAccount = dto.IFBglEntry.IFBglProductAccount
		}
		if dto.IFBglEntry.IFBglProductBranchcode != "" {
			updatedService.CBEIFBGLEntry.ProductBranchCode = dto.IFBglEntry.IFBglProductBranchcode
		}
		if dto.IFBglEntry.IFBglServiceAccount != "" {
			updatedService.CBEIFBGLEntry.ServiceAccount = dto.IFBglEntry.IFBglServiceAccount
		}
		if dto.IFBglEntry.IFBglServiceBranchcode != "" {
			updatedService.CBEIFBGLEntry.ServiceBranchCode = dto.IFBglEntry.IFBglServiceBranchcode
		}

		// Map tiers from DTO
		if len(dto.Tiers) > 0 {
			tiers := make([]types.Tier, len(dto.Tiers))
			for i, tier := range dto.Tiers {
				tiers[i] = types.Tier{
					Min:       tier.Min,
					Max:       tier.Max,
					FeeAmount: tier.FeeAmount,
				}
			}
			updatedService.Tiers = tiers
			updatedService.AboveAmount = tiers[len(tiers)-1].Max
		}

		// Ensure the ID is preserved
		updatedService.ID = existingService.ID
		return &updatedService, nil
	}

	// Handle SingleMaxTransferRequest DTO
	if singleCapReq, ok := newData.(dto.SingleMaxTransferRequest); ok {
		// Map cap values from SingleMaxTransferRequest
		if singleCapReq.ISingleCap > 0 {
			updatedService.Cap.ISingleCap = singleCapReq.ISingleCap
		}
		if singleCapReq.IDailyCap > 0 {
			updatedService.Cap.IDailyCap = singleCapReq.IDailyCap
		}
		if singleCapReq.CSingleCap > 0 {
			updatedService.Cap.CorporateSingleCap = singleCapReq.CSingleCap
		}
		if singleCapReq.CDailyCap > 0 {
			updatedService.Cap.CorporateDailyCap = singleCapReq.CDailyCap
		}

		// Ensure the ID is preserved
		updatedService.ID = existingService.ID
		return &updatedService, nil
	}

	// Handle MinimumTransferUpdateRequest DTO
	if minCapReq, ok := newData.(dto.MinimumTransferUpdateRequest); ok {
		// Map minimum amount from MinimumTransferUpdateRequest
		if minCapReq.Minimum > 0 {
			updatedService.Cap.MinAmount = minCapReq.Minimum
		}

		// Ensure the ID is preserved
		updatedService.ID = existingService.ID
		return &updatedService, nil
	}

	// Fallback to generic map handling if direct cast fails
	var incoming map[string]interface{}
	if err := BindAction(newData, &incoming); err != nil {
		return nil, err
	}

	// Map simple scalar fields if present
	if v, ok := incoming["service_type"].(string); ok && v != "" {
		updatedService.ServiceType = v
	}
	if v, ok := incoming["payment_type"].(string); ok && v != "" {
		updatedService.PaymentType = v
	}
	if v, ok := incoming["above_amount"]; ok {
		switch n := v.(type) {
		case float64:
			updatedService.AboveAmount = uint64(n)
		case int:
			updatedService.AboveAmount = uint64(n)
		case int64:
			updatedService.AboveAmount = uint64(n)
		case uint64:
			updatedService.AboveAmount = n
		}
	}
	if v, ok := incoming["above_service_fee"]; ok {

		switch n := v.(type) {
		case float64:
			fmt.Println("float")
			updatedService.AboveServiceFee = uint64(n)
		case int:
			updatedService.AboveServiceFee = uint64(n)
		case int64:
			updatedService.AboveServiceFee = uint64(n)
		case uint64:
			updatedService.AboveServiceFee = n
		}
	}

	// Map cap fields if present
	if v, ok := incoming["individual_single_cap"]; ok {
		switch n := v.(type) {
		case float64:
			updatedService.Cap.ISingleCap = uint64(n)
		case int:
			updatedService.Cap.ISingleCap = uint64(n)
		case int64:
			updatedService.Cap.ISingleCap = uint64(n)
		case uint64:
			updatedService.Cap.ISingleCap = n
		}
	}
	if v, ok := incoming["individual_daily_cap"]; ok {
		switch n := v.(type) {
		case float64:
			updatedService.Cap.IDailyCap = uint64(n)
		case int:
			updatedService.Cap.IDailyCap = uint64(n)
		case int64:
			updatedService.Cap.IDailyCap = uint64(n)
		case uint64:
			updatedService.Cap.IDailyCap = n
		}
	}
	if v, ok := incoming["corporate_single_cap"]; ok {
		switch n := v.(type) {
		case float64:
			updatedService.Cap.CorporateSingleCap = uint64(n)
		case int:
			updatedService.Cap.CorporateSingleCap = uint64(n)
		case int64:
			updatedService.Cap.CorporateSingleCap = uint64(n)
		case uint64:
			updatedService.Cap.CorporateSingleCap = n
		}
	}
	if v, ok := incoming["corporate_daily_cap"]; ok {
		switch n := v.(type) {
		case float64:
			updatedService.Cap.CorporateDailyCap = uint64(n)
		case int:
			updatedService.Cap.CorporateDailyCap = uint64(n)
		case int64:
			updatedService.Cap.CorporateDailyCap = uint64(n)
		case uint64:
			updatedService.Cap.CorporateDailyCap = n
		}
	}
	if v, ok := incoming["min_amount"]; ok {
		switch n := v.(type) {
		case float64:
			updatedService.Cap.MinAmount = uint64(n)
		case int:
			updatedService.Cap.MinAmount = uint64(n)
		case int64:
			updatedService.Cap.MinAmount = uint64(n)
		case uint64:
			updatedService.Cap.MinAmount = n
		}
	}

	// Map GL entries - accept either canonical keys or alternate cbgl_/ifbgl_ prefixed keys
	if v, ok := incoming["cbe_gl_entry"].(map[string]interface{}); ok {
		if acc, ok := v["product_account"].(string); ok && acc != "" {
			updatedService.CBEGLEntry.ProductAccount = acc
		}
		if bc, ok := v["product_branch_code"].(string); ok && bc != "" {
			updatedService.CBEGLEntry.ProductBranchCode = bc
		}
		if acc, ok := v["service_account"].(string); ok && acc != "" {
			updatedService.CBEGLEntry.ServiceAccount = acc
		}
		if bc, ok := v["service_branch_code"].(string); ok && bc != "" {
			updatedService.CBEGLEntry.ServiceBranchCode = bc
		}
		if acc, ok := v["vat_account"].(string); ok && acc != "" {
			updatedService.CBEGLEntry.VatAccount = acc
		}
		if bc, ok := v["vat_branch_code"].(string); ok && bc != "" {
			updatedService.CBEGLEntry.VatBranchCode = bc
		}
	}
	if v, ok := incoming["cbgl_entry"].(map[string]interface{}); ok {
		if acc, ok := v["cbgl_product_account"].(string); ok && acc != "" {
			updatedService.CBEGLEntry.ProductAccount = acc
		}
		if bc, ok := v["cbgl_product_branchcode"].(string); ok && bc != "" {
			updatedService.CBEGLEntry.ProductBranchCode = bc
		}
		if acc, ok := v["cbgl_service_account"].(string); ok && acc != "" {
			updatedService.CBEGLEntry.ServiceAccount = acc
		}
		if bc, ok := v["cbgl_service_branchcode"].(string); ok && bc != "" {
			updatedService.CBEGLEntry.ServiceBranchCode = bc
		}
	}

	if v, ok := incoming["cbe_ifb_gl_entry"].(map[string]interface{}); ok {
		if acc, ok := v["product_account"].(string); ok && acc != "" {
			updatedService.CBEIFBGLEntry.ProductAccount = acc
		}
		if bc, ok := v["product_branch_code"].(string); ok && bc != "" {
			updatedService.CBEIFBGLEntry.ProductBranchCode = bc
		}
		if acc, ok := v["service_account"].(string); ok && acc != "" {
			updatedService.CBEIFBGLEntry.ServiceAccount = acc
		}
		if bc, ok := v["service_branch_code"].(string); ok && bc != "" {
			updatedService.CBEIFBGLEntry.ServiceBranchCode = bc
		}
		if acc, ok := v["vat_account"].(string); ok && acc != "" {
			updatedService.CBEIFBGLEntry.VatAccount = acc
		}
		if bc, ok := v["vat_branch_code"].(string); ok && bc != "" {
			updatedService.CBEIFBGLEntry.VatBranchCode = bc
		}
	}
	if v, ok := incoming["ifbgl_entry"].(map[string]interface{}); ok {
		if acc, ok := v["ifbgl_product_account"].(string); ok && acc != "" {
			updatedService.CBEIFBGLEntry.ProductAccount = acc
		}
		if bc, ok := v["ifbgl_product_branchcode"].(string); ok && bc != "" {
			updatedService.CBEIFBGLEntry.ProductBranchCode = bc
		}
		if acc, ok := v["ifbgl_service_account"].(string); ok && acc != "" {
			updatedService.CBEIFBGLEntry.ServiceAccount = acc
		}
		if bc, ok := v["ifbgl_service_branchcode"].(string); ok && bc != "" {
			updatedService.CBEIFBGLEntry.ServiceBranchCode = bc
		}
	}

	// Map tiers with ID preservation
	if v, ok := incoming["tiers"].([]interface{}); ok {
		incomingTiers := make([]types.Tier, 0, len(v))
		for _, it := range v {
			if tierMap, ok := it.(map[string]interface{}); ok {
				var t types.Tier
				if min, ok := tierMap["min"].(float64); ok {
					t.Min = uint64(min)
				}
				if max, ok := tierMap["max"].(float64); ok {
					t.Max = uint64(max)
				}
				if fee, ok := tierMap["fee_amount"].(float64); ok {
					t.FeeAmount = uint64(fee)
				}
				incomingTiers = append(incomingTiers, t)
			}
		}

		// Assign tiers as-is; persistence layer will omit any id fields during update
		updatedService.Tiers = incomingTiers
	}

	// Ensure the ID is preserved
	updatedService.ID = existingService.ID

	return &updatedService, nil
}

func MapHQForTotalCapUpdate(existingHQ *model.HQ, newTotalCap uint64) *model.HQ {
	updatedHQ := *existingHQ
	updatedHQ.TotalCap = newTotalCap
	return &updatedHQ
}

func BindAction(source any, target any) error {
	bytes, err := json.Marshal(source)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, target)
}

func ValidateMinimumTransferCap(newMinAmount uint64, existingCaps types.Cap) error {

	if newMinAmount >= existingCaps.ISingleCap {

		return fmt.Errorf("min_amount (%d) must be less than individual_single_cap (%d)", newMinAmount, existingCaps.ISingleCap)
	}
	if newMinAmount >= existingCaps.IDailyCap {
		return fmt.Errorf("min_amount (%d) must be less than individual_daily_cap (%d)", newMinAmount, existingCaps.IDailyCap)
	}
	if newMinAmount >= existingCaps.CorporateSingleCap {
		return fmt.Errorf("min_amount (%d) must be less than corporate_single_cap (%d)", newMinAmount, existingCaps.CorporateSingleCap)
	}
	if newMinAmount >= existingCaps.CorporateDailyCap {
		return fmt.Errorf("min_amount (%d) must be less than corporate_daily_cap (%d)", newMinAmount, existingCaps.CorporateDailyCap)
	}
	return nil
}
func TotalCapMapper(existingService *model.HQ, newData interface{}) (*model.HQ, error) {
	updatedService := *existingService
	var incoming map[string]interface{}
	if err := BindAction(newData, &incoming); err != nil {
		return nil, err
	}

	if v, ok := incoming["total_cap"]; ok {
		switch n := v.(type) {
		case float64:
			updatedService.TotalCap = uint64(n)
		case int:
			updatedService.TotalCap = uint64(n)
		case int64:
			updatedService.TotalCap = uint64(n)
		case uint64:
			updatedService.TotalCap = n
		}
	}
	return &updatedService, nil

}
