package mappers

import (
	"encoding/json"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func DomainToModelCPSAction(domainAction domain.CPSAction) *model.CPSAction {

	var objID bson.ObjectID
	if domainAction.ID != "" {
		var err error
		objID, err = bson.ObjectIDFromHex(domainAction.ID)
		if err != nil {
			objID = bson.NilObjectID
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
		b, _ := json.Marshal(v)
		if err := json.Unmarshal(b, &currentAction); err != nil {
			currentAction = nil
		}
	}

	return &model.CPSAction{
		ID:                 objID,
		ActionCode:         domainAction.ActionCode,
		UniqueId:           domainAction.UniqueId,
		MakerID:            domainAction.MakerID,
		MakerName:          domainAction.MakerName,
		MakerPhoneNumber:   domainAction.MakerPhoneNumber,
		CheckerID:          domainAction.CheckerID,
		CheckerName:        domainAction.CheckerName,
		CheckerPhoneNumber: domainAction.CheckerPhoneNumber,
		Department:         domainAction.Department,
		RejectionReason: func() string {
			if domainAction.RejectionReason != nil {
				return *domainAction.RejectionReason
			}
			return ""
		}(),
		PreviosAction:     domainAction.PreviosAction,
		CurrentAction:     currentAction,
		ActionStatus:      string(domainAction.ActionStatus),
		ActionType:        string(domainAction.ActionType),
		RequestAction:     string(domainAction.RequestAction),
		CreatedAt:         domainAction.CreatedAt,
		LastModifiedAt:    domainAction.LastModifiedAt,
		MakerActionTime:   domainAction.MakerActionTime,
		CheckerActionTime: domainAction.CheckerActionTime,
	}
}

// Helper to recursively convert bson.D to map[string]interface{}
func BsonDToMap(i any) any {
	switch v := i.(type) {
	case bson.D:
		m := make(map[string]interface{})
		for _, e := range v {
			m[e.Key] = BsonDToMap(e.Value)
		}
		return m
	case []any:
		for i, e := range v {
			v[i] = BsonDToMap(e)
		}
		return v
	default:
		return v
	}
}

func ModelToDomainCPSAction(modelAction model.CPSAction) *domain.CPSAction {
	var rejectionReason *string
	if modelAction.RejectionReason != "" {
		rejectionReason = &modelAction.RejectionReason
	}

	return &domain.CPSAction{
		ID:                 modelAction.ID.Hex(),
		ActionCode:         modelAction.ActionCode,
		UniqueId:           modelAction.UniqueId,
		MakerID:            modelAction.MakerID,
		MakerName:          modelAction.MakerName,
		MakerPhoneNumber:   modelAction.MakerPhoneNumber,
		CheckerID:          modelAction.CheckerID,
		CheckerName:        modelAction.CheckerName,
		CheckerPhoneNumber: modelAction.CheckerPhoneNumber,
		Department:         modelAction.Department,
		RejectionReason:    rejectionReason,
		PreviosAction:      modelAction.PreviosAction,
		CurrentAction:      BsonDToMap(modelAction.CurrentAction),
		ActionStatus:       domain.ActionStatus(modelAction.ActionStatus),
		ActionType:         domain.ActionType(modelAction.ActionType),
		RequestAction:      domain.RequestAction(modelAction.RequestAction),
		CreatedAt:          modelAction.CreatedAt,
		LastModifiedAt:     modelAction.LastModifiedAt,
		MakerActionTime:    modelAction.MakerActionTime,
		CheckerActionTime:  modelAction.CheckerActionTime,
	}
}
