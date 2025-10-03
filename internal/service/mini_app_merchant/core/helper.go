package core

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/service"
	"context"
	"errors"
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

func MergeMiniAppMerchantData(old, data *model.MiniAppMerchant) *model.MiniAppMerchant {
	now := time.Now()
	return &model.MiniAppMerchant{
		ID:                old.ID,
		Code:              old.Code,
		MerchantName:      local_util.NonEmptyString(data.MerchantName, old.MerchantName),
		MerchantType:      local_util.NonEmptyString(data.MerchantType, old.MerchantType),
		PhoneNumber:       local_util.NonEmptyString(data.PhoneNumber, old.PhoneNumber),
		Email:             local_util.NonEmptyString(data.Email, old.Email),
		BankAccountNumber: local_util.NonEmptyString(data.BankAccountNumber, old.BankAccountNumber),
		Enabled:           old.Enabled,
		IsDeleted:         old.IsDeleted,
		CreatedAt:         old.CreatedAt,
		LastModifiedAt:    now,
		KYC: model.KYC{
			Status: old.KYC.Status,
			Representative: model.KYCInformation{
				Name:  local_util.NonEmptyString(data.KYC.Representative.Name, old.KYC.Representative.Name),
				Email: local_util.NonEmptyString(data.KYC.Representative.Email, old.KYC.Representative.Email),
				Phone: local_util.NonEmptyString(data.KYC.Representative.Phone, old.KYC.Representative.Phone),
			},
		},
		Branches: old.Branches,
		MiniApps: old.MiniApps,
	}
}

func HandleCPSActionForMiniAppMerchant(ctx context.Context, cpsService service.CPSActionService, uniqueID string, requestAction constants.RequestAction, curData, prevData interface{}, actionType constants.ActionType) error {
	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	cpsAction := lib.CpsModelBuilder(uniqueID, maker, prevData, curData, string(requestAction), string(actionType))

	err := cpsService.CreateCPSAction(ctx, &cpsAction)
	if err != nil {
		return err
	}
	return nil
}
