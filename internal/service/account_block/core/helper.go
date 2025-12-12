package core

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func GenerateCPSAction(ctx context.Context, codeType string, enabled bool, prev []model.EnableDisableAction, curr []model.EnableDisableAction, reason string, actionType constants.ActionType, requestActionType constants.RequestAction) model.CPSAction {
	maker := local_util.ExtractUserFromContext(ctx)

	uniqueID := bson.NewObjectID().Hex()
	cpsAction := lib.CpsModelBuilder(uniqueID, maker, prev, curr, string(requestActionType), string(actionType))

	return cpsAction
}
