package mappers

import (
	"encoding/json"
	"fmt"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	constants "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// DomainToModelCPSAction converts a domain CPSAction to a model CPSAction for MongoDB storage.
func DomainToModelCPSAction(domainAction domain.CPSAction) (*model.CPSAction, error) {
	// Validate required fields
	if domainAction.ActionCode == "" {
		return nil, fmt.Errorf("ActionCode is required")
	}

	// Validate enum fields
	if !constants.IsValidActionStatus(string(domainAction.ActionStatus)) {
		return nil, fmt.Errorf("invalid ActionStatus: %s", domainAction.ActionStatus)
	}
	if !constants.IsValidActionType(string(domainAction.ActionType)) {
		return nil, fmt.Errorf("invalid ActionType: %s", domainAction.ActionType)
	}
	if !constants.IsValidRequestAction(string(domainAction.RequestAction)) {
		return nil, fmt.Errorf("invalid RequestAction: %s", domainAction.RequestAction)
	}

	// Convert domain ID to BSON ObjectID
	var objID bson.ObjectID
	if domainAction.ID != "" {
		var err error
		objID, err = bson.ObjectIDFromHex(domainAction.ID)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", common_util.InvalidID, err)
		}
	} else {
		objID = bson.NewObjectID()
	}

	// Convert CurrentAction to map[string]interface{} if needed
	var currentAction any
	switch v := domainAction.CurrentAction.(type) {
	case map[string]any:
		currentAction = v
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal CurrentAction: %w", err)
		}
		if err := json.Unmarshal(b, &currentAction); err != nil {
			return nil, fmt.Errorf("failed to unmarshal CurrentAction: %w", err)
		}
	}

	return &model.CPSAction{
		ID:                 objID,
		ActionCode:         domainAction.ActionCode,
		UniqueId:           domainAction.UniqueID,
		MakerID:            domainAction.MakerID,
		MakerName:          domainAction.MakerName,
		MakerPhoneNumber:   domainAction.MakerPhoneNumber,
		CheckerID:          domainAction.CheckerID,
		CheckerName:        domainAction.CheckerName,
		CheckerPhoneNumber: domainAction.CheckerPhoneNumber,
		Department:         domainAction.Department,
		RejectionReason:    domainAction.RejectionReason,
		PreviousAction:     domainAction.PreviousAction,
		CurrentAction:      currentAction,
		ActionStatus:       string(domainAction.ActionStatus),
		ActionType:         string(domainAction.ActionType),
		RequestAction:      string(domainAction.RequestAction),
		CreatedAt:          domainAction.CreatedAt,
		LastModifiedAt:     domainAction.LastModifiedAt,
		MakerActionTime:    domainAction.MakerActionTime,
		CheckerActionTime:  domainAction.CheckerActionTime,
	}, nil
}

// BSONToMap recursively converts BSON types to map[string]interface{} or slices.
func BSONToMap(i any) any {
	switch v := i.(type) {
	case bson.D:
		m := make(map[string]interface{})
		for _, e := range v {
			m[e.Key] = BSONToMap(e.Value)
		}
		return m
	case bson.A:
		a := make([]interface{}, len(v))
		for i, e := range v {
			a[i] = BSONToMap(e)
		}
		return a
	case []interface{}:
		a := make([]interface{}, len(v))
		for i, e := range v {
			a[i] = BSONToMap(e)
		}
		return a
	case primitive.DateTime:
		return v.Time()
	default:
		return v
	}
}

// ModelToDomainCPSAction converts a model CPSAction to a domain CPSAction.
func ModelToDomainCPSAction(modelAction model.CPSAction) (*domain.CPSAction, error) {

	return &domain.CPSAction{
		ID:                 modelAction.ID.Hex(),
		ActionCode:         modelAction.ActionCode,
		UniqueID:           modelAction.UniqueId,
		MakerID:            modelAction.MakerID,
		MakerName:          modelAction.MakerName,
		MakerPhoneNumber:   modelAction.MakerPhoneNumber,
		CheckerID:          modelAction.CheckerID,
		CheckerName:        modelAction.CheckerName,
		CheckerPhoneNumber: modelAction.CheckerPhoneNumber,
		Department:         modelAction.Department,
		RejectionReason:    modelAction.RejectionReason,
		PreviousAction:     modelAction.PreviousAction,
		CurrentAction:      BSONToMap(modelAction.CurrentAction),
		ActionStatus:       constants.ActionStatus(modelAction.ActionStatus),
		ActionType:         constants.ActionType(modelAction.ActionType),
		RequestAction:      constants.RequestAction(modelAction.RequestAction),
		CreatedAt:          modelAction.CreatedAt,
		LastModifiedAt:     modelAction.LastModifiedAt,
		MakerActionTime:    modelAction.MakerActionTime,
		CheckerActionTime:  modelAction.CheckerActionTime,
	}, nil
}
