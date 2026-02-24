package core

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	"context"
	"errors"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	merchantDto "cbe-super-app-cps-action/internal/constants/dto/ecommerce-merchant"

	"cbe-super-app-cps-action/internal/storage/external_call/account_lookup"

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

func nonEmpty(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

func MergeBranches(oldBranches, newBranches []model.BranchInformation) []model.BranchInformation {
	oldMap := make(map[string]model.BranchInformation, len(oldBranches))
	for _, ob := range oldBranches {
		oldMap[ob.BranchCode] = ob
	}

	merged := make([]model.BranchInformation, 0, len(newBranches))

	for _, nb := range newBranches {
		// If the branch code exists in old branches, merge the information
		if ob, ok := oldMap[nb.BranchCode]; ok {
			merged = append(merged, model.BranchInformation{
				BranchCode:          ob.BranchCode,
				BranchName:          nonEmpty(nb.BranchName, ob.BranchName),
				BranchAddress:       nonEmpty(nb.BranchAddress, ob.BranchAddress),
				BranchOwner:         nonEmpty(nb.BranchOwner, ob.BranchOwner),
				BranchAccountNumber: nonEmpty(nb.BranchAccountNumber, ob.BranchAccountNumber),
			})
			delete(oldMap, nb.BranchCode)
		} else {
			merged = append(merged, model.BranchInformation{
				BranchCode:          nb.BranchCode,
				BranchName:          nb.BranchName,
				BranchAddress:       nb.BranchAddress,
				BranchOwner:         nb.BranchOwner,
				BranchAccountNumber: nb.BranchAccountNumber,
			})
		}
	}

	for _, leftover := range oldMap {
		merged = append(merged, leftover)
	}

	return merged
}

func MergeMiniAppMerchantData(old, data *model.EcommerceMerchant) *model.EcommerceMerchant {
	now := time.Now()

	updatedBranches := MergeBranches(old.Branches, data.Branches)

	return &model.EcommerceMerchant{
		ID:                old.ID,
		Code:              old.Code,
		MerchantName:      local_util.NonEmptyString(data.MerchantName, old.MerchantName),
		PhoneNumber:       local_util.NonEmptyString(data.PhoneNumber, old.PhoneNumber),
		Email:             local_util.NonEmptyString(data.Email, old.Email),
		BankAccountNumber: local_util.NonEmptyString(data.BankAccountNumber, old.BankAccountNumber),
		Enabled:           old.Enabled,
		IsDeleted:         old.IsDeleted,
		CreatedAt:         old.CreatedAt,
		UpdatedAt:         now,
		Branches:          updatedBranches,
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

func ValidateAccountNumberWithExternalAPI(ctx context.Context, accountNumber string, accountLookupService account_lookup.Account) (*model.AccountDetail, error) {
	accountRequest := model.AccountLookUpRequest{
		AccountNumber: accountNumber,
	}

	accountDetail, err := accountLookupService.LookupAccountByAccountNumber(ctx, accountRequest)
	if err != nil {
		return nil, errors.New(localization.ErrorAccountNumberNotFound.Code)
	}

	if accountDetail == nil {
		return nil, errors.New(localization.ErrorAccountNumberNotFound.Code)
	}

	return accountDetail, nil
}

func CheckMerchantExists(
	ctx context.Context,
	merchantRepo storage.EcommerceMerchantRepository,
	data *types.CheckMiniAppMerchant,
	opts *types.MiniAppMerchantExistOptions,
) (bool, error) {
	if data == nil {
		return false, nil
	}

	var conditions []bson.M
	if data.MerchantCode != "" {
		conditions = append(conditions, bson.M{"merchant_code": data.MerchantCode})
	}

	if len(conditions) == 0 {
		return false, nil
	}

	filter := bson.M{
		"is_deleted": false,
		"$or":        conditions,
	}

	if opts != nil && opts.ExcludeID != "" {
		objID, err := bson.ObjectIDFromHex(opts.ExcludeID)
		if err != nil {
			return false, err
		}
		filter["_id"] = bson.M{"$ne": objID}
	}

	res, err := merchantRepo.FindOne(ctx, filter)
	if err != nil {
		if err.Error() == localization.ErrorMiniAppMerchantNotFound.Code {
			return false, nil
		}
		return false, err
	}

	if res == nil {
		return false, nil
	}

	// if res.BankAccountNumber == data.BankAccountNumber {
	// 	return false, errors.New(localization.ErrorAccountNumberAlreadyExists.Code)
	// }
	// if res.Email == data.Email {
	// 	return false, errors.New(localization.ErrorEmailAlreadyExist.Code)
	// }
	// if res.PhoneNumber == data.PhoneNumber {
	// 	return false, errors.New(localization.ErrorPhonenumberAlreadyExist.Code)
	// }
	if res.Code == data.MerchantCode {
		return false, errors.New(localization.ErrorCodeAlreadyExist.Code)
	}

	return true, nil
}

func ToMiniAppMerchantResponseDTO(domain *model.EcommerceMerchant) *merchantDto.MiniAppMerchantResponseDTO {
	return &merchantDto.MiniAppMerchantResponseDTO{
		ID:            domain.ID.Hex(),
		Code:          domain.Code,
		MerchantName:  domain.MerchantName,
		AccountNumber: domain.BankAccountNumber,
		Enabled:       domain.Enabled,
		IsDeleted:     domain.IsDeleted,
		CreatedAt:     domain.CreatedAt,
		LastModified:  domain.UpdatedAt,
	}
}

// Convert DTO to Domain model for service layer
func ToMiniAppMerchantDomainFromUpdateDTO(d *merchantDto.EcommerceMerchant) *model.EcommerceMerchant {
	return &model.EcommerceMerchant{
		ID:                bson.NewObjectID(),
		Code:              d.MerchantCode,
		MerchantName:      d.MerchantName,
		BankAccountNumber: d.AccountNumber,
		// Email:             d.Email,
		// PhoneNumber:       d.PhoneNumber,
		SettlementMethod: d.SettlementMethod,
		Branches:         d.Branches,
		Enabled:          true,
		IsDeleted:        false,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}
