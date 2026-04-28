package core

import (
	"cbe-super-app-cps-action/internal/constants"
	walletDto "cbe-super-app-cps-action/internal/constants/dto/wallet"
	"cbe-super-app-cps-action/internal/constants/lib"

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
	logger.Infof("[WalletCore][GenPrefix] prefix: %s value: %s", prefix, value)

	if prefix == "" || value == "" {
		logger.Errorf("[WalletCore][GenPrefix] invalid input prefix: %s value: %s", prefix, value)
		return "", fmt.Errorf("prefix and value must not be empty")
	}

	value = strings.ToUpper(strings.ReplaceAll(value, " ", "_"))
	prefix = strings.ToUpper(strings.ReplaceAll(prefix, " ", "_"))
	result := strings.Join([]string{prefix, value}, "-")
	logger.Infof("[WalletCore][GenPrefix] result: %s", result)
	return result, nil

}

func ToCreateWalletDoc(name, code, URL string, self, other, agent *bool, selfServiceID, otherServiceID, agentServiceID string) *local_model.WalletOracle {
	return &local_model.WalletOracle{
		Name:           name,
		UniqueCode:     code,
		Avatar:         URL,
		Self:           *self,
		Other:          *other,
		Agent:          *agent,
		SelfServiceID:  selfServiceID,
		OtherServiceID: otherServiceID,
		AgentServiceID: agentServiceID,
	}
}

// note: this comparision might not be needed if the existing data is first in the request form and the user update those values
func ToUpdateWalletDoc(existing local_model.WalletOracle, req walletDto.WalletRequest) (*local_model.WalletOracle, int) {
	wallet := existing
	changeCount := 0

	if req.Self != nil && *req.Self != existing.Self {
		changeCount++
		wallet.Self = *req.Self
	}

	if req.Other != nil && *req.Other != existing.Other {
		changeCount++
		wallet.Other = *req.Other
	}

	if req.Agent != nil && *req.Agent != existing.Agent {
		changeCount++
		wallet.Agent = *req.Agent
	}

	if req.Name != "" && req.Name != existing.Name {
		changeCount++
		wallet.Name = req.Name
	}

	if req.UniqueCode != "" && req.UniqueCode != existing.UniqueCode {
		changeCount++
		wallet.UniqueCode = req.UniqueCode
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
