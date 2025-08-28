package core

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/model"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func GenerateCPSAction(ctx context.Context, codeType string, enabled bool, codes []string, actionType constants.ActionType, requestActionType constants.RequestAction) model.CPSAction {
	maker := local_util.ExtractUserFromContext(ctx)

	var actions []any
	for _, code := range codes {
		switch codeType {
		case "BRANCH":
			actions = append(actions, model.Branch{
				BranchCode: code,
				Enabled:    enabled,
			})
		case "REGION":
			actions = append(actions, model.Region{
				RegionCode: code,
				Enabled:    enabled,
			})
		case "DISTRICT":
			actions = append(actions, model.District{
				DistrictCode: code,
				Enabled:      enabled,
			})
		case "CITY":
			actions = append(actions, model.City{
				CityCode: code,
				Enabled:  enabled,
			})
		}
	}

	uniqueID := bson.NewObjectID().Hex()

	cpsAction := lib.CpsModelBuilder(uniqueID, maker, nil, actions, string(requestActionType), string(actionType) )

	return cpsAction
}