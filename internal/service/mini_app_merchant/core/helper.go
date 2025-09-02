package core

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/model"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"encoding/json"

	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	"context"
	"errors"
	"log"
	"time"
)

func NonEmptyString(s, fallback string) string {
	if s != "" {
		return s
	}
	return fallback
}

func NonZeroTime(t, fallback time.Time) time.Time {
	if !t.IsZero() {
		return t
	}
	return fallback
}

func NonZeroUint64(n, fallback uint64) uint64 {
	if n != 0 {
		return n
	}
	return fallback
}

func MergeMiniAppMerchant(data model.MiniAppMerchant) *model.MiniAppMerchant {
	now := time.Now()
	return &model.MiniAppMerchant{
		ID:                data.ID,
		Code:              data.Code,
		MerchantName:      local_util.NonEmptyString(data.MerchantName, data.MerchantName),
		MerchantType:      local_util.NonEmptyString(data.MerchantType, data.MerchantType),
		PhoneNumber:       local_util.NonEmptyString(data.PhoneNumber, data.PhoneNumber),
		Email:             local_util.NonEmptyString(data.Email, data.Email),
		BankAccountNumber: local_util.NonEmptyString(data.BankAccountNumber, data.BankAccountNumber),
		Enabled:           data.Enabled,
		IsDeleted:         data.IsDeleted,
		CreatedAt:         data.CreatedAt,
		LastModifiedAt:    now,
		KYC: model.KYC{
			Status: data.KYC.Status,
			Representative: model.KYCInformation{
				Name:  local_util.NonEmptyString(data.KYC.Representative.Name, data.KYC.Representative.Name),
				Email: local_util.NonEmptyString(data.KYC.Representative.Email, data.KYC.Representative.Email),
				Phone: local_util.NonEmptyString(data.KYC.Representative.Phone, data.KYC.Representative.Phone),
			},
		},
		Branches: data.Branches,
		MiniApps: data.MiniApps,
	}
}

// HandleCPSActionForMiniAppMerchant creates a CPS action related to Mini App Merchant
func HandleCPSActionForMiniAppMerchant(ctx context.Context, cpsService service.CPSActionService, uniqueID string, requestAction constants.RequestAction, curData, prevData interface{}, actionType constants.ActionType) error {
	userData := local_util.ExtractUserFromContext(ctx)

	if incomplet := local_util.IsIncomplete(userData); incomplet {
		log.Printf("User data incomplete for CPS action: %+v", userData)
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}

	cpsAction := lib.CpsModelBuilder(uniqueID, userData, prevData, curData, string(requestAction), string(actionType))

	if err := cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		log.Printf("Failed to create CPS action for Mini App Merchant: %v", err)
		return err
	}
	return nil
}

func BindAction(source any, target any) error {
	bytes, err := json.Marshal(source)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, target)
}
