package adapter

import (
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/bps"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/model"
	infra_mongo "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/adapter/outbound/mongo"
	domain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/action"
	outbound "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/outbound/bulk_services"
)

type outboundStore struct {
	MongoDalCPSAction *infra_mongo.MongoDal[model.CPSAction, model.CPSAction]
	MongoDalBPSUser   *infra_mongo.MongoDal[bps.BPSUser, bps.BPSUser]
	MongoDalServices  *infra_mongo.MongoDal[model.Service, model.Service]
	MongoDalMember    *infra_mongo.MongoDal[member.User, member.User]
}

func NewOutBoundStore(client *mongo.Client, dbName string, collectionNames []string) outbound.OutboundInfra {
	mongoDalCPSAction := infra_mongo.NewMongoDal[model.CPSAction, model.CPSAction](client, dbName, collectionNames[0])
	mongoDalBPSUser := infra_mongo.NewMongoDal[bps.BPSUser, bps.BPSUser](client, dbName, collectionNames[1])
	mongoDalService := infra_mongo.NewMongoDal[model.Service, model.Service](client, dbName, collectionNames[2])
	mongoDalMember := infra_mongo.NewMongoDal[member.User, member.User](client, dbName, collectionNames[3])
	return &outboundStore{
		MongoDalCPSAction: mongoDalCPSAction,
		MongoDalBPSUser:   mongoDalBPSUser,
		MongoDalServices:  mongoDalService,
		MongoDalMember:    mongoDalMember,
	}
}

func (o *outboundStore) GetAllHqServices(ctx context.Context) ([]domain.Service, error) {
	data, err := o.MongoDalServices.FindAll(ctx, nil, nil)
	if err != nil {
		return nil, err
	}
	var d []domain.Service
	for _, v := range data {
		x := domain.Service{
			ID:             stringToPointer(v.ID.Hex()),
			Key:            v.Key,
			ServiceName:    v.ServiceName,
			SingleCap:      v.SingleCap,
			MinAmount:      v.MinAmount,
			DailyCap:       v.DailyCap,
			Flag:           v.Flag,
			CreatedAt:      v.CreatedAt,
			LastModifiedAt: v.LastModifiedAt,
		}
		d = append(d, x)
	}
	return d, nil
}
func (o *outboundStore) GetAllHqServicesPaginated(ctx context.Context, offset, limit int) ([]domain.Service, error) {
	skip := int64(offset)      // Convert offset to int64 for MongoDB compatibility
	limitInt64 := int64(limit) // Convert limit to int64 for MongoDB compatibility

	data, err := o.MongoDalServices.FindAllWithPagination(ctx, nil, nil, skip, limitInt64)
	if err != nil {
		return nil, err
	}

	var d []domain.Service
	for _, v := range data {
		x := domain.Service{
			ID:             stringToPointer(v.ID.Hex()),
			Key:            v.Key,
			ServiceName:    v.ServiceName,
			SingleCap:      v.SingleCap,
			MinAmount:      v.MinAmount,
			DailyCap:       v.DailyCap,
			Flag:           v.Flag,
			CreatedAt:      v.CreatedAt,
			LastModifiedAt: v.LastModifiedAt,
		}
		d = append(d, x)
	}
	return d, nil
}
func (o *outboundStore) GetHqServiceById(ctx context.Context, id string) (domain.Service, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return domain.Service{}, err
	}
	filter := map[string]interface{}{"_id": objID}
	data, err := o.MongoDalServices.FindOne(ctx, filter, nil)
	if err != nil {
		return domain.Service{}, err
	}
	return domain.Service{
		ID:             stringToPointer(data.ID.Hex()),
		Key:            data.Key,
		ServiceName:    data.ServiceName,
		SingleCap:      data.SingleCap,
		MinAmount:      data.MinAmount,
		DailyCap:       data.DailyCap,
		Flag:           data.Flag,
		CreatedAt:      data.CreatedAt,
		LastModifiedAt: data.LastModifiedAt,
	}, nil
}
func (o *outboundStore) UpdateHqService(ctx context.Context, service domain.Service) error {
	filter := map[string]interface{}{
		"_id": service.ID,
	}
	update := map[string]interface{}{
		"$set": map[string]interface{}{
			"key":            service.Key,
			"serviceName":    service.ServiceName,
			"singleCap":      service.SingleCap,
			"minAmount":      service.MinAmount,
			"dailyCap":       service.DailyCap,
			"flag":           service.Flag,
			"lastModifiedAt": service.LastModifiedAt,
		},
	}
	_, err := o.MongoDalServices.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (o *outboundStore) CreateCpsAction(ctx context.Context, Action domain.CPSAction) (domain.CPSAction, error) {
	CpsAction := model.CPSAction{
		ActionCode:         Action.ActionCode,
		MakerID:            Action.MakerAndChecker.Maker.UserID,
		MakerName:          Action.MakerAndChecker.Maker.FullName,
		MakerPhoneNumber:   Action.MakerAndChecker.Maker.PhoneNumber,
		CheckerID:          &Action.MakerAndChecker.Checker.UserID,
		CheckerName:        &Action.MakerAndChecker.Checker.FullName,
		CheckerPhoneNumber: &Action.MakerAndChecker.Checker.PhoneNumber,
		Department:         Action.Department,
		RejectionReason:    Action.RejectionReason,
		PreviosAction:      struct{}{}, // Assuming this is empty for now
		CurrentAction: model.CurrentAction{
			Id:     Action.CurrentAction.Id,
			Action: Action.CurrentAction.Action,
		},
		ActionStatus:   model.ActionStatus(Action.ActionStatus),
		ActionType:     model.ActionType(Action.ActionType),
		RequestAction:  model.RequestAction(Action.RequestAction),
		CreatedAt:      Action.CreatedAt,
		LastModifiedAt: Action.LastModifiedAt,
	}
	data, err := o.MongoDalCPSAction.InsertOne(ctx, CpsAction)
	if err != nil {
		return domain.CPSAction{}, nil
	}
	result := domain.CPSAction{
		ID:              data.ID.Hex(),
		ActionCode:      data.ActionCode,
		Department:      data.Department,
		RejectionReason: data.RejectionReason,
		PreviosAction:   data.PreviosAction, // Assuming this is directly mapped
		CurrentAction: domain.CurrentAction{
			Id:     data.CurrentAction.Id,
			Action: data.CurrentAction.Action,
		},
		ActionStatus:   domain.ActionStatus(data.ActionStatus),
		ActionType:     domain.ActionType(data.ActionType),
		RequestAction:  domain.RequestAction(data.RequestAction),
		CreatedAt:      data.CreatedAt,
		LastModifiedAt: data.LastModifiedAt,
	}
	return result, nil
}
func (o *outboundStore) UpdateCpsAction(ctx context.Context, Action domain.CPSAction) error {
	CpsAction := model.CPSAction{
		ActionCode:         Action.ActionCode,
		MakerID:            Action.MakerAndChecker.Maker.UserID,
		MakerName:          Action.MakerAndChecker.Maker.FullName,
		MakerPhoneNumber:   Action.MakerAndChecker.Maker.PhoneNumber,
		CheckerID:          &Action.MakerAndChecker.Checker.UserID,
		CheckerName:        &Action.MakerAndChecker.Checker.FullName,
		CheckerPhoneNumber: &Action.MakerAndChecker.Checker.PhoneNumber,
		Department:         Action.Department,
		RejectionReason:    Action.RejectionReason,
		PreviosAction:      Action.PreviosAction, // Assuming this is directly mapped
		CurrentAction: model.CurrentAction{
			Id:     Action.CurrentAction.Id,
			Action: Action.CurrentAction.Action,
		},
		ActionStatus:   model.ActionStatus(Action.ActionStatus),
		ActionType:     model.ActionType(Action.ActionType),
		RequestAction:  model.RequestAction(Action.RequestAction),
		CreatedAt:      Action.CreatedAt,
		LastModifiedAt: Action.LastModifiedAt,
	}
	filter := map[string]interface{}{
		"action_code": CpsAction.ActionCode,
	}
	update := map[string]interface{}{
		"$set": map[string]interface{}{
			"makerID":            CpsAction.MakerID,
			"makerName":          CpsAction.MakerName,
			"makerPhoneNumber":   CpsAction.MakerPhoneNumber,
			"checkerID":          CpsAction.CheckerID,
			"checkerName":        CpsAction.CheckerName,
			"checkerPhoneNumber": CpsAction.CheckerPhoneNumber,
			"department":         CpsAction.Department,
			"rejectionReason":    CpsAction.RejectionReason,
			"previosAction":      CpsAction.PreviosAction,
			"currentAction": map[string]interface{}{
				"id":     CpsAction.CurrentAction.Id,
				"action": CpsAction.CurrentAction.Action,
			},
			"actionStatus":   CpsAction.ActionStatus,
			"actionType":     CpsAction.ActionType,
			"requestAction":  CpsAction.RequestAction,
			"createdAt":      CpsAction.CreatedAt,
			"lastModifiedAt": CpsAction.LastModifiedAt,
		},
	}
	_, err := o.MongoDalCPSAction.UpdateOne(ctx, filter, update)
	return err
}
func (o *outboundStore) FetchCpsActionById(ctx context.Context, Action_Id string) (domain.CPSAction, error) {
	filter := map[string]interface{}{"action_code": Action_Id}
	data, err := o.MongoDalCPSAction.FindOne(ctx, filter, nil)
	if err != nil {
		return domain.CPSAction{}, err
	}
	result := domain.CPSAction{
		ID:              data.ID.Hex(),
		ActionCode:      data.ActionCode,
		Department:      data.Department,
		RejectionReason: data.RejectionReason,
		PreviosAction:   data.PreviosAction, // Assuming this is directly mapped
		CurrentAction: domain.CurrentAction{
			Id:     data.CurrentAction.Id,
			Action: data.CurrentAction.Action,
		},
		ActionStatus:   domain.ActionStatus(data.ActionStatus),
		ActionType:     domain.ActionType(data.ActionType),
		RequestAction:  domain.RequestAction(data.RequestAction),
		CreatedAt:      data.CreatedAt,
		LastModifiedAt: data.LastModifiedAt,
	}
	return result, nil
}

func stringToPointer(s string) *string {
	return &s
}
