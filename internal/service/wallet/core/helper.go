package core

import (
	"cbe-super-app-cps-action/internal/constants"
	walletDto "cbe-super-app-cps-action/internal/constants/dto/wallet"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/types"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/service"
	"context"
	cRand "crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
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

	digits := "0123456789"
	max := big.NewInt(int64(len(digits)))
	code := make([]byte, 7)

	for i := range code {
		n, err := cRand.Int(cRand.Reader, max)
		if err != nil {
			logger.Errorf("Failed to generate random digit", "error", err)
			return "", fmt.Errorf("failed to generate random digit: %v", err)
		}
		code[i] = digits[n.Int64()]
	}

	value = strings.ReplaceAll(value, " ", "")
	prefix = strings.ReplaceAll(prefix, " ", "")
	result := strings.Join([]string{prefix, value, string(code)}, "-")
	logger.Infof("Successfully generated prefixed name", "result", result)
	return result, nil

}

func ToCreateWalletDoc(name, code, URL string, self, other, agent bool) *model.Wallet {
	return &model.Wallet{
		Name:   name,
		Code:   code,
		Avatar: URL,
		Services: types.Services{
			Self:  self,
			Other: other,
			Agent: agent,
		},
	}
}

// note: this comparision might not be needed if the existing data is first in the request form and the user update those values
func ToUpdateWalletDoc(existing model.Wallet, req walletDto.WalletRequest) (*model.Wallet, int) {
	var wallet model.Wallet
	change_count := 0
	if req.Agent == existing.Services.Agent {
		wallet.Services.Agent = existing.Services.Agent
	} else {
		change_count++
		wallet.Services.Agent = req.Agent
	}
	if req.Other == existing.Services.Other {
		wallet.Services.Other = existing.Services.Other
	} else {
		change_count++
		wallet.Services.Other = req.Other
	}
	if req.Self == existing.Services.Self {
		wallet.Services.Self = existing.Services.Self
	} else {
		change_count++
		wallet.Services.Self = req.Self
	}
	if req.Name == existing.Name {
		wallet.Name = existing.Name
	} else {
		change_count++
		wallet.Name = req.Name
	}
	if req.Code == existing.Code {
		wallet.Code = existing.Code
	} else {
		change_count++
		wallet.Code = req.Code
	}
	wallet.Avatar = existing.Avatar
	return &wallet, change_count
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
