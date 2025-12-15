package core

import (
	"cbe-super-app-cps-action/internal/constants"
	walletDto "cbe-super-app-cps-action/internal/constants/dto/wallet"
	"cbe-super-app-cps-action/internal/constants/lib"

	// "cbe-super-app-cps-action/internal/constants/types"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	shared_type "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/types"

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

func ToCreateWalletDoc(name, code, URL string, self, other, agent bool, walletType string) *model.Wallet {
	return &model.Wallet{
		Name:   name,
		Code:   code,
		Avatar: URL,
		Services: shared_type.Services{
			Self:  self,
			Other: other,
			Agent: agent,
		},
		Type: walletType,
	}
}

// note: this comparision might not be needed if the existing data is first in the request form and the user update those values
func ToUpdateWalletDoc(existing model.Wallet, req walletDto.WalletRequest, fieldsProvided map[string]bool) (*model.Wallet, int) {
	wallet := existing
	changeCount := 0

	if fieldsProvided["self"] && req.Self != existing.Services.Self {
		changeCount++
		wallet.Services.Self = req.Self
	}

	if fieldsProvided["other"] && req.Other != existing.Services.Other {
		changeCount++
		wallet.Services.Other = req.Other
	}

	if fieldsProvided["agent"] && req.Agent != existing.Services.Agent {
		changeCount++
		wallet.Services.Agent = req.Agent
	}

	if req.Name != "" && req.Name != existing.Name {
		changeCount++
		wallet.Name = req.Name
	}

	if req.Code != "" && req.Code != existing.Code {
		changeCount++
		wallet.Code = req.Code
	}
	if req.Type != "" && req.Type != existing.Type {
		changeCount++
		wallet.Type = req.Type
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
