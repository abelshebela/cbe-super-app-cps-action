package core

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/model"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"encoding/json"
	"fmt"
	"reflect"

	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func HandleCPSAction(ctx context.Context, cpsService service.CPSActionService, uniqueID string, requestAction constants.RequestAction, curData, prevData interface{}, actionType constants.ActionType) error {
	userData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(userData); incomplet {
		log.Println("User data incomplete for CPS action", "userCode", userData.UserCode)
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}

	cpsAction := lib.CpsModelBuilder(uniqueID, userData, prevData, curData, string(requestAction), string(actionType))

	log.Println("Creating CPS action", "userCode", userData.UserCode, "uniqueID", uniqueID, "actionType", actionType)
	if err := cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		log.Println("Failed to create CPS action", "error", err)
		return err
	}
	log.Println("Successfully created CPS action", "userCode", userData.UserCode, "uniqueID", uniqueID)
	return nil
}

func ParseHQ(raw interface{}) (*model.HQ, error) {
	switch v := raw.(type) {
	case *model.HQ:
		return v, nil
	case bson.D:
		var hq model.HQ
		// Convert bson.D to bytes using bson.Marshal
		b, err := bson.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal bson.D: %v", err)
		}
		// Unmarshal directly into HQ struct
		if err := bson.Unmarshal(b, &hq); err != nil {
			return nil, fmt.Errorf("failed to unmarshal bson.D to HQ: %v", err)
		}
		return &hq, nil
	case map[string]interface{}:
		var hq model.HQ
		b, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal map: %v", err)
		}
		if err := json.Unmarshal(b, &hq); err != nil {
			return nil, fmt.Errorf("failed to unmarshal map to HQ: %v", err)
		}
		return &hq, nil
	default:
		return nil, fmt.Errorf("unsupported type: %T", raw)
	}
}

func ToJSONBytes(val interface{}) ([]byte, error) {
	switch v := val.(type) {
	case nil:
		return nil, fmt.Errorf("value is nil")
	case json.RawMessage:
		return v, nil
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	case map[string]interface{}, []interface{}:
		return json.Marshal(v)
	default:
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.Struct {
			return json.Marshal(v)
		}
		return nil, fmt.Errorf("unsupported type: %T", v)
	}

}

func ConvertToHQ(currentAction interface{}) (*model.HQ, error) {
	bsonBytes, err := bson.Marshal(currentAction)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal CurrentAction: %w", err)
	}

	var hq model.HQ
	if err := bson.Unmarshal(bsonBytes, &hq); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to HQ: %w", err)
	}

	return &hq, nil
}
