package core

import (
	"bytes"
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	local_model "cbe-super-app-cps-action/internal/constants/model"
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

)

func HandleCPSActionForLogisticsMerchant(ctx context.Context, cpsService service.CPSActionService, uniqueID string, requestAction constants.RequestAction, curData, prevData interface{}, actionType constants.ActionType) error {
	maker := local_util.ExtractUserFromContext(ctx)
	isErp, _ := ctx.Value(constants.ContextKey("is_erp")).(bool)

	if !isErp && local_util.IsIncomplete(maker) {
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
	merchantRepo storage.LogisticsMerchantOracleRepository,
	data *types.CheckMerchant,
	opts *types.MiniAppMerchantExistOptions,
) (bool, error) {
	if data == nil || (data.BankAccountNumber == "" && data.MerchantCode == "") {
		return false, nil
	}

	excludeID := ""
	if opts != nil {
		excludeID = opts.ExcludeID
	}

	res, err := merchantRepo.FindByAccountOrMerchantCode(ctx, data.BankAccountNumber, data.MerchantCode, excludeID)
	if err != nil {
		return false, err
	}
	if res == nil {
		return false, nil
	}

	if data.BankAccountNumber != "" && res.BankAccountNumber == data.BankAccountNumber {
		return false, errors.New(localization.ErrorAccountNumberAlreadyExists.Code)
	}
	if data.MerchantCode != "" && res.MerchantID == data.MerchantCode {
		return false, errors.New(localization.ErrorMerchantCodeAlreadyExists.Code)
	}

	return true, nil
}
func MergeLogisticsMerchantData(old, data *local_model.LogisticsMerchant) *local_model.LogisticsMerchant {
	now := time.Now()

	return &local_model.LogisticsMerchant{
		ID:                old.ID,
		MerchantID:        local_util.NonEmptyString(data.MerchantID, old.MerchantID),
		SettlementMethod:  local_util.NonEmptyString(data.SettlementMethod, old.SettlementMethod),
		MerchantName:      local_util.NonEmptyString(data.MerchantName, old.MerchantName),
		MerchantType:      local_util.NonEmptyString(data.MerchantType, old.MerchantType),
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
		logger.Errorf("[LogMerchCore][ERPUpdate] marshal err: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, base, bytes.NewBuffer(jsonBody))
	if err != nil {
		logger.Errorf("[LogMerchCore][ERPUpdate] build req err: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-api-key", apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("[LogMerchCore][ERPUpdate] request err: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		logger.Errorf("[LogMerchCore][ERPUpdate] failed: %s", string(bodyBytes))
		return errors.New("ERP update failed")
	}

	return nil
}
