package core

import (
	"bytes"
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/external_call/account_lookup"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func HandleCPSActionForEventMerchant(ctx context.Context, cpsService service.CPSActionService, uniqueID string, requestAction constants.RequestAction, curData, prevData interface{}, actionType constants.ActionType) error {
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
	merchantRepo storage.EventMerchantRepository,
	data *types.CheckMiniAppMerchant,
	opts *types.MiniAppMerchantExistOptions,
) (bool, error) {
	if data == nil {
		return false, nil
	}

	var conditions []bson.M

	if data.BankAccountNumber != "" {
		conditions = append(conditions, bson.M{"bank_account_number": data.BankAccountNumber, "merchant_id": data.MerchantCode})
	}
	if data.MerchantCode != "" {
		conditions = append(conditions, bson.M{"merchant_id": data.MerchantCode})
	}
	if data.Email != "" {
		conditions = append(conditions, bson.M{"email": data.Email})
	}
	if data.PhoneNumber != "" {
		conditions = append(conditions, bson.M{"phone_number": data.PhoneNumber})
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
		if err.Error() == localization.ErrorEventMerchantNotFound.Code {
			return false, nil
		}
		return false, err
	}

	if res == nil {
		return false, nil
	}

	if res.BankAccountNumber == data.BankAccountNumber {
		return false, errors.New(localization.ErrorAccountNumberAlreadyExists.Code)
	}
	if res.MerchantID == data.MerchantCode {
		return false, errors.New(localization.ErrorMerchantIDAlreadyExists.Code)
	}

	return true, nil
}
func MergeEventMerchantData(old, data *model.EventMerchant) *model.EventMerchant {
	now := time.Now()

	return &model.EventMerchant{
		ID:                old.ID,
		MerchantID:        local_util.NonEmptyString(data.MerchantID, old.MerchantID),
		SettlementMethod:  local_util.NonEmptyString(data.SettlementMethod, old.SettlementMethod),
		MerchantName:      local_util.NonEmptyString(data.MerchantName, old.MerchantName),
		MerchantType:      local_util.NonEmptyString(data.MerchantType, old.MerchantType),
		PhoneNumber:       local_util.NonEmptyString(data.PhoneNumber, old.PhoneNumber),
		Email:             local_util.NonEmptyString(data.Email, old.Email),
		BankAccountNumber: local_util.NonEmptyString(data.BankAccountNumber, old.BankAccountNumber),
		Enabled:           old.Enabled,
		IsDeleted:         old.IsDeleted,
		CreatedAt:         old.CreatedAt,
		UpdatedAt:         now,
	}
}

func UpdateERP(ctx context.Context, cfg *config.VaultConfig, bankAccountNumber, merchantID string, logger utils.Logger) error {
	ctx, span := local_util.TraceLogger(ctx, "core", "UpdateERP", "LogisticsMerchant", "UpdateERP")
	defer span.End()

	base := "https://qaapisuperapp.cbe.com.et/api/v1/cbesuperapp/ecommerce"
	if cfg != nil && cfg.OddoEcommerceBaseUrl != "" {
		base = cfg.OddoEcommerceBaseUrl
	}
	base += "/cps/merchant/update/" + merchantID
	apiKey := ""
	if cfg != nil && cfg.ApiKey != "" {
		apiKey = cfg.ApiKey
	}

	reqBody := map[string]string{
		"cps_account_number": bankAccountNumber,
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		logger.Errorf("Failed to marshal ERP update body: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, base, bytes.NewBuffer(jsonBody))
	if err != nil {
		logger.Errorf("Failed to build ERP update request: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-api-key", apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("ERP update request failed: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		logger.Errorf("ERP update failed: %s", string(bodyBytes))
		return errors.New("ERP update failed")
	}

	return nil
}
