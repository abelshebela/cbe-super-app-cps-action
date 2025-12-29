package core

import (
	"cbe-super-app-cps-action/internal/constants"
	walletDto "cbe-super-app-cps-action/internal/constants/dto/wallet"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/types"

	// "cbe-super-app-cps-action/internal/constants/types"
	local_model "cbe-super-app-cps-action/internal/constants/model"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func NonEmptyString(s, fallback string) string {
	if s != "" {
		return s
	}
	return fallback
}

func GeneratePrefixedName(prefix, value string, logger shared_utils.Logger) (string, error) {
	logger.Infof("Generating prefixed name", "prefix", prefix, "value", value)

	if prefix == "" || value == "" {
		logger.Errorf("Invalid input for GeneratePrefixedName", "prefix", prefix, "value", value)
		return "", fmt.Errorf("prefix and value must not be empty")
	}

	value = strings.ToUpper(strings.ReplaceAll(value, " ", "_"))
	prefix = strings.ToUpper(strings.ReplaceAll(prefix, " ", "_"))
	result := strings.Join([]string{prefix, value}, "-")
	logger.Infof("Successfully generated prefixed name", "result", result)
	return result, nil

}

func ToCreateWalletDoc(name, code, URL, serviceCode string, self, other, agent *bool) *local_model.Wallet {
	return &local_model.Wallet{
		Name:        name,
		UniqueCode:  code,
		Avatar:      URL,
		ServiceCode: serviceCode,
		Services: types.Services{
			Self:  *self,
			Other: *other,
			Agent: *agent,
		},
	}
}

// note: this comparision might not be needed if the existing data is first in the request form and the user update those values
func ToUpdateWalletDoc(existing local_model.Wallet, req walletDto.WalletRequest, serviceCode string) (*local_model.Wallet, int) {
	wallet := existing
	changeCount := 0

	if req.Self != nil && *req.Self != existing.Services.Self {
		changeCount++
		wallet.Services.Self = *req.Self
	}

	if req.Other != nil && existing.Services.Other {
		changeCount++
		wallet.Services.Other = *req.Other
	}

	if req.Agent != nil && existing.Services.Agent {
		changeCount++
		wallet.Services.Agent = *req.Agent
	}

	if req.Name != "" && req.Name != existing.Name {
		changeCount++
		wallet.Name = req.Name
	}

	if req.UniqueCode != "" && req.UniqueCode != existing.UniqueCode {
		changeCount++
		wallet.UniqueCode = req.UniqueCode
	}
	if serviceCode != "" && serviceCode != existing.ServiceCode {
		changeCount++
		wallet.ServiceCode = serviceCode
	}

	return &wallet, changeCount
}

func HandleCPSAction(ctx context.Context, cpsService service.CPSActionService, uniqueID string, requestAction constants.RequestAction, curData, prevData interface{}, actionType constants.ActionType) error {
	userData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(userData); incomplet {
		log.Println("User data incomplete for CPS action", "userCode", userData.UserCode)
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}

	cpsAction := lib.CpsModelBuilder(uniqueID, userData, prevData, curData, string(requestAction), string(actionType))

	if err := cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		log.Println("Failed to create CPS action", "error", err)
		return err
	}
	log.Println("Successfully created CPS action", "userCode", userData.UserCode, "uniqueID", uniqueID)
	return nil
}
