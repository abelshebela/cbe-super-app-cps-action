package core

import (
	"cbe-super-app-cps-action/internal/constants"
	passwordrule "cbe-super-app-cps-action/internal/constants/dto/password_rule"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	local_model "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"log"
)

func HandleCPSAction(ctx context.Context, cpsService service.CPSActionService, uniqueID string, requestAction constants.RequestAction, curData, prevData interface{}, actionType constants.ActionType) error {
	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		log.Printf("incomplete user context for fayda account | context = %v", maker)
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	cpsAction := lib.CpsModelBuilder(uniqueID, maker, prevData, curData, string(requestAction), string(actionType))

	err := cpsService.CreateCPSAction(ctx, &cpsAction)
	if err != nil {
		log.Printf("[passwordRule.HandleCPSAction] failed to create CPS action | action: %s, error: %v", requestAction, err)
		return err
	}
	return nil
}

func PasswordRuleDtoToModel(existing local_model.PasswordRule, dto passwordrule.PasswordRuleUpdate) local_model.PasswordRule {
	if dto.Rule.Name != "" {
		existing.Name = dto.Rule.Name
	}
	if dto.Rule.MinLength != 0 {
		existing.MinLength = dto.Rule.MinLength
	}
	if dto.Rule.MaxLength != 0 {
		existing.MaxLength = dto.Rule.MaxLength
	}

	if dto.Rule.Numbers != nil {
		existing.Numbers = dto.Rule.Numbers
	}
	if dto.Rule.CapitalLetters != nil {
		existing.CapitalLetters = dto.Rule.CapitalLetters
	}
	if dto.Rule.SmallLetters != nil {
		existing.SmallLetters = dto.Rule.SmallLetters
	}
	if dto.Rule.Characters != nil {
		existing.Characters = dto.Rule.Characters
	}

	return existing
}
