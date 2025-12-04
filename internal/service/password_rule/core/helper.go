package core

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
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
