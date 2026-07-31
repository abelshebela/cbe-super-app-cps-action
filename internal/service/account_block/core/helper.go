package core

import (
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/lib"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/types"
	local_util "github.com/abelshebela/cbe-super-app-cps-action/pkgs/utils"
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"github.com/google/uuid"
)

func GenerateCPSAction(ctx context.Context, codeType string, enabled bool, prev []types.EnableDisableAction, curr []types.EnableDisableAction, reason string, actionType constants.ActionType, requestActionType constants.RequestAction) model.CPSAction {
	maker := local_util.ExtractUserFromContext(ctx)

	uniqueID := uuid.New().String()
	cpsAction := lib.CpsModelBuilder(uniqueID, maker, prev, curr, string(requestActionType), string(actionType))

	return cpsAction
}
