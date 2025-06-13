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
	MongoDalAccounts  *infra_mongo.MongoDal[model.LinkedAccount, model.LinkedAccount]
}

func NewOutBoundStore(client *mongo.Client, dbName string, collectionNames []string) outbound.OutboundInfra {
	mongoDalCPSAction := infra_mongo.NewMongoDal[model.CPSAction, model.CPSAction](client, dbName, collectionNames[0])
	mongoDalBPSUser := infra_mongo.NewMongoDal[bps.BPSUser, bps.BPSUser](client, dbName, collectionNames[1])
	mongoDalService := infra_mongo.NewMongoDal[model.Service, model.Service](client, dbName, collectionNames[2])
	mongoDalMember := infra_mongo.NewMongoDal[member.User, member.User](client, dbName, collectionNames[3])
	mongoDalAccounts := infra_mongo.NewMongoDal[model.LinkedAccount, model.LinkedAccount](client, dbName, collectionNames[4])
	return &outboundStore{
		MongoDalCPSAction: mongoDalCPSAction,
		MongoDalBPSUser:   mongoDalBPSUser,
		MongoDalServices:  mongoDalService,
		MongoDalMember:    mongoDalMember,
		MongoDalAccounts:  mongoDalAccounts,
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

func (o *outboundStore) FetchAccountsByCif(ctx context.Context, cif string) ([]domain.LinkedAccount, error) {
	filter := map[string]interface{}{"customer_number": cif}
	data, err := o.MongoDalAccounts.FindAll(ctx, filter, nil)
	if err != nil {
		return nil, err
	}
	var accounts []domain.LinkedAccount
	for _, v := range data {
		account := domain.LinkedAccount{
			ID:                stringToPointer(v.ID.Hex()),
			UserID:            stringToPointer(v.UserID.Hex()),
			CustomerNumber:    v.CustomerNumber,
			AccountNumber:     v.AccountNumber,
			AccountHolderName: v.AccountHolderName,
			AccountType:       v.AccountType,
			BranchCode:        v.BranchCode,
			LinkedStatus:      v.LinkedStatus,
			LastLinkedStatus:  v.LastLinkedStatus,
			LinkedAt:          v.LinkedAt,
			LinkedBranch:      v.LinkedBranch,
			RegistrationType:  domain.RegistrationType(v.RegistrationType),
			IsAccountActive:   v.IsAccountActive,
			AndOrStatus:       v.AndOrStatus,
			AccountBranchCode: v.AccountBranchCode,
			CurrencyCode:      v.CurrencyCode,
			IsMain:            v.IsMain,
			MakerAndChecker: struct {
				Linkers struct {
					Maker   string
					Checker string
				}
				Unlinkers struct {
					Maker   string
					Checker string
				}
			}{
				Linkers: struct {
					Maker   string
					Checker string
				}{
					Maker:   v.MakerAndChecker.Linkers.Maker,
					Checker: v.MakerAndChecker.Linkers.Checker,
				},
				Unlinkers: struct {
					Maker   string
					Checker string
				}{
					Maker:   v.MakerAndChecker.Unlinkers.Maker,
					Checker: v.MakerAndChecker.Unlinkers.Checker,
				},
			},
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (o *outboundStore) UpdateAccounts(ctx context.Context, linkedAccounts []domain.LinkedAccount) ([]domain.LinkedAccount, error) {
	var updatedAccounts []domain.LinkedAccount
	for _, acc := range linkedAccounts {
		filter := map[string]interface{}{"_id": acc.ID}
		update := map[string]interface{}{
			"$set": map[string]interface{}{
				"user_id":             acc.UserID,
				"customer_number":     acc.CustomerNumber,
				"account_number":      acc.AccountNumber,
				"account_holder_name": acc.AccountHolderName,
				"account_type":        acc.AccountType,
				"branch_code":         acc.BranchCode,
				"linked_status":       acc.LinkedStatus,
				"last_linked_status":  acc.LastLinkedStatus,
				"linked_at":           acc.LinkedAt,
				"linker_branch":       acc.LinkedBranch,
				"registration_type":   acc.RegistrationType,
				"is_account_active":   acc.IsAccountActive,
				"and_or_status":       acc.AndOrStatus,
				"account_branch_code": acc.AccountBranchCode,
				"currency":            acc.CurrencyCode,
				"is_main":             acc.IsMain,
				"maker_and_checker": map[string]interface{}{
					"linkers": map[string]interface{}{
						"maker":   acc.MakerAndChecker.Linkers.Maker,
						"checker": acc.MakerAndChecker.Linkers.Checker,
					},
					"unlinkers": map[string]interface{}{
						"maker":   acc.MakerAndChecker.Unlinkers.Maker,
						"checker": acc.MakerAndChecker.Unlinkers.Checker,
					},
				},
				"created_at": acc.CreatedAt,
				"updated_at": acc.UpdatedAt,
			},
		}
		_, err := o.MongoDalAccounts.UpdateOne(ctx, filter, update)
		if err != nil {
			return nil, err
		}
		updatedAccounts = append(updatedAccounts, acc)
	}
	return updatedAccounts, nil
}

func (o *outboundStore) UpdateAccount(ctx context.Context, linkedAccount domain.LinkedAccount) (domain.LinkedAccount, error) {
	filter := map[string]interface{}{"_id": linkedAccount.ID}
	update := map[string]interface{}{
		"$set": map[string]interface{}{
			"user_id":             linkedAccount.UserID,
			"customer_number":     linkedAccount.CustomerNumber,
			"account_number":      linkedAccount.AccountNumber,
			"account_holder_name": linkedAccount.AccountHolderName,
			"account_type":        linkedAccount.AccountType,
			"branch_code":         linkedAccount.BranchCode,
			"linked_status":       linkedAccount.LinkedStatus,
			"last_linked_status":  linkedAccount.LastLinkedStatus,
			"linked_at":           linkedAccount.LinkedAt,
			"linker_branch":       linkedAccount.LinkedBranch,
			"registration_type":   linkedAccount.RegistrationType,
			"is_account_active":   linkedAccount.IsAccountActive,
			"and_or_status":       linkedAccount.AndOrStatus,
			"account_branch_code": linkedAccount.AccountBranchCode,
			"currency":            linkedAccount.CurrencyCode,
			"is_main":             linkedAccount.IsMain,
			"maker_and_checker": map[string]interface{}{
				"linkers": map[string]interface{}{
					"maker":   linkedAccount.MakerAndChecker.Linkers.Maker,
					"checker": linkedAccount.MakerAndChecker.Linkers.Checker,
				},
				"unlinkers": map[string]interface{}{
					"maker":   linkedAccount.MakerAndChecker.Unlinkers.Maker,
					"checker": linkedAccount.MakerAndChecker.Unlinkers.Checker,
				},
			},
			"created_at": linkedAccount.CreatedAt,
			"updated_at": linkedAccount.UpdatedAt,
		},
	}
	_, err := o.MongoDalAccounts.UpdateOne(ctx, filter, update)
	if err != nil {
		return domain.LinkedAccount{}, err
	}
	return linkedAccount, nil
}

func (o *outboundStore) FetchLastCpsActionByMakerID(ctx context.Context, makerId string) (domain.CPSAction, error) {
	filter := map[string]interface{}{"maker_id": makerId}

	d, err := o.MongoDalCPSAction.FindAll(ctx, filter, nil)
	if err != nil {
		return domain.CPSAction{}, err
	}
	var data model.CPSAction
	data = *d[len(d)-1] // Get the last action for the maker
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
