package core

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"log"
	// "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
)

func HandleCPSAction(ctx context.Context, cpsService service.CPSActionService, uniqueID string, requestAction constants.RequestAction, curData, prevData interface{}, actionType constants.ActionType) error {
	userData := local_util.ExtractUserFromContext(ctx)
	// if incomplet := local_util.IsIncomplete(userData); incomplet {
	// 	log.Println("User data incomplete for CPS action", "userCode", userData.UserCode)
	// 	return errors.New(localization.ErrorIncompleteUserInfo.Code)
	// }

	cpsAction := lib.CpsModelBuilder(uniqueID, userData, prevData, curData, string(requestAction), string(actionType))

	if err := cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		log.Println("Failed to create CPS action", "error", err)
		return err
	}
	log.Println("Successfully created CPS action", "userCode", userData.UserCode, "uniqueID", uniqueID)
	return nil
}

func ValidateCreate(req imodel.Services, service storage.ServicesRepository) error {
	if req.ServiceCode == "" || req.ServiceName == "" {
		return localization.ErrorRequiredFieldMissing
	}

	serviceDoc, err := service.FindAllWithPagination(context.Background(), types.Filter{Filters: map[string]interface{}{
		"service_code": req.ServiceCode,
		"service_name": req.ServiceName,
	}})
	if err != nil {
		return err
	}
	if len(serviceDoc.Data) > 0 {
		return localization.ErrorServiceExists
	}

	return nil
}
