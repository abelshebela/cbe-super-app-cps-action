package core

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/service"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
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

func CascadeEnableDisableMiniApps(ctx context.Context, miniRepo storage.MiniAppRepository, merchantID string, enabled bool) error {
	filter := types.Filter{Filters: map[string]interface{}{"merchant_id": merchantID, "is_deleted": false}, Page: 1, PerPage: 100}
	for {
		res, err := miniRepo.FindAllWithPagination(ctx, filter)
		if err != nil {
			return err
		}
		if res == nil || len(res.Data) == 0 {
			return nil
		}
		lib.GoRoutinBaker(types.BakerOptions{Sequential: false, UseMutex: false},
			func() {
				for _, m := range res.Data {
					_ = miniRepo.EnableOrDisable(ctx, m.ID.Hex(), enabled)
				}
			},
		)
		if len(res.Data) < filter.PerPage {
			return nil
		}
		filter.Page++
	}
}

func CascadeDeleteMiniApps(ctx context.Context, miniRepo storage.MiniAppRepository, merchantID string) error {
	filter := types.Filter{Filters: map[string]interface{}{"merchant_id": merchantID, "is_deleted": false}, Page: 1, PerPage: 100}
	for {
		res, err := miniRepo.FindAllWithPagination(ctx, filter)
		if err != nil {
			return err
		}
		if res == nil || len(res.Data) == 0 {
			return nil
		}
		lib.GoRoutinBaker(types.BakerOptions{Sequential: false, UseMutex: false},
			func() {
				for _, m := range res.Data {
					_ = miniRepo.Delete(ctx, m.ID.Hex())
				}
			},
		)
		if len(res.Data) < filter.PerPage {
			return nil
		}
		filter.Page++
	}
}

func CheckMerchantExists(ctx context.Context, merchantRepo storage.MiniAppMerchantRepository, data *model.CheckMiniAppMerchant, opts *model.MiniAppMerchantExistOptions) (bool, error) {
	if data == nil {
		return false, nil
	}

	var conditions []map[string]interface{}
	if data.BankAccountNumber != "" {
		conditions = append(conditions, map[string]interface{}{"bank_account_number": data.BankAccountNumber})
	}
	if data.Email != "" {
		conditions = append(conditions, map[string]interface{}{"kyc.representative.email": data.Email})
	}
	if data.PhoneNumber != "" {
		conditions = append(conditions, map[string]interface{}{"kyc.representative.phone": data.PhoneNumber})
	}

	if len(conditions) == 0 {
		return false, nil
	}

	filter := bson.M{
		"is_deleted": false,
		"$or":        conditions,
	}

	if opts != nil && opts.ExcludeID != "" {
		if objID, err := bson.ObjectIDFromHex(opts.ExcludeID); err == nil {
			filter["_id"] = bson.M{"$ne": objID}
		} else {
			return false, err
		}
	}

	res, err := merchantRepo.FindOne(ctx, filter)

	if err != nil {

		if err.Error() == "ERROR_MINI_APP_MERCHANT_NOT_FOUND" {
			return false, nil
		}
		return false, err
	}

	return res != nil, nil
}
