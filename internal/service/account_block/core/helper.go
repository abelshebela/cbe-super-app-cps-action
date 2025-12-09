package core

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/model"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func GenerateCPSAction(ctx context.Context, codeType string, enabled bool, ids []string, reason string, actionType constants.ActionType, requestActionType constants.RequestAction) model.CPSAction {
	maker := local_util.ExtractUserFromContext(ctx)

	var action model.EnableDisableAction
	action.Codes = ids
	action.Reason = reason

	uniqueID := bson.NewObjectID().Hex()

	cpsAction := lib.CpsModelBuilder(uniqueID, maker, nil, action, string(requestActionType), string(actionType))

	return cpsAction
}
