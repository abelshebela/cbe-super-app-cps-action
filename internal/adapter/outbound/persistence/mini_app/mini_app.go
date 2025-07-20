package miniapp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	dal "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/infra"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	model "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"

	utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	miniApp "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/mini_app"
)

type MiniAppPersistence struct {
	MongoDalMiniApp   dal.MongoDal[model.MiniApp, model.MiniApp]
	MongoDalCPSAction dal.MongoDal[model.CPSAction, model.CPSAction]
	logger            utils.Logger
}

func InitMiniAppPersistence(client *mongo.Client, DB_name string, collections []string, logger utils.Logger) miniApp.Outbound {
	return &MiniAppPersistence{
		MongoDalMiniApp:   dal.NewMongoDal[model.MiniApp, model.MiniApp](client, DB_name, collections[0]),
		MongoDalCPSAction: dal.NewMongoDal[model.CPSAction, model.CPSAction](client, DB_name, collections[1]),
	}
}

var _ miniApp.Outbound = (*MiniAppPersistence)(nil)

func (o *MiniAppPersistence) CreateMiniAppAction(ctx context.Context, Action model.CPSAction) (model.CPSAction, error) {

	fmt.Println("perstance===========================")
	fmt.Printf("current action to store in DB:%v", Action.CurrentAction)
	fmt.Println("perstance===========================")

	CpsAction := model.CPSAction{
		ActionCode:         Action.ActionCode,
		MakerID:            Action.MakerID,
		MakerName:          Action.MakerName,
		MakerPhoneNumber:   Action.MakerPhoneNumber,
		CheckerID:          Action.CheckerID,
		CheckerName:        Action.CheckerName,
		CheckerPhoneNumber: Action.CheckerPhoneNumber,
		Department:         Action.Department,
		RejectionReason:    Action.RejectionReason,
		MakerActionTime:    Action.MakerActionTime,
		PreviousAction:     Action.PreviousAction,
		CurrentAction:      Action.CurrentAction,
		ActionStatus:       string(model.ActionPending),
		ActionType:         Action.ActionType,
		RequestAction:      Action.RequestAction,
		CreatedAt:          time.Now(),
		LastModifiedAt:     time.Now(),
	}

	data, err := o.MongoDalCPSAction.InsertOne(ctx, CpsAction)
	if err != nil {
		return model.CPSAction{}, err
	}

	result := model.CPSAction{
		ID:                 data.ID,
		ActionCode:         data.ActionCode,
		UniqueId:           data.UniqueId,
		MakerID:            data.MakerID,
		MakerName:          data.MakerName,
		MakerPhoneNumber:   data.MakerPhoneNumber,
		CheckerID:          data.CheckerID,
		CheckerName:        data.CheckerName,
		CheckerPhoneNumber: data.CheckerPhoneNumber,
		Department:         data.Department,
		RejectionReason:    data.RejectionReason,
		PreviousAction: func() interface{} {
			var v interface{}
			if b, ok := data.PreviousAction.(json.RawMessage); ok {
				_ = json.Unmarshal(b, &v)
			}
			return v
		}(),
		CurrentAction: func() interface{} {
			var v interface{}
			if b, ok := data.CurrentAction.(json.RawMessage); ok {
				_ = json.Unmarshal(b, &v)
			}
			return v
		}(),
		// PreviousAction:     Action.PreviousAction,
		// CurrentAction:     Action.CurrentAction,
		ActionStatus:      data.ActionStatus,
		ActionType:        data.ActionType,
		RequestAction:     data.RequestAction,
		CreatedAt:         data.CreatedAt,
		LastModifiedAt:    data.LastModifiedAt,
		MakerActionTime:   data.MakerActionTime,
		CheckerActionTime: data.CheckerActionTime,
	}
	return result, nil
}

func (o *MiniAppPersistence) DeleteMiniAppAction(ctx context.Context, action model.CPSAction, id string) (model.CPSAction, error) {
	return model.CPSAction{}, nil
}

func (o *MiniAppPersistence) CreateMiniApp(ctx context.Context, miniapp model.MiniApp) (model.CPSAction, error) {
	// var credential []domain.CredentialInformation
	// var product_code []domain.ProductCode
	// for _, creds := range miniapp.Credential {
	// 	credential = append(credential, domain.CredentialInformation{
	// 		Environment:   domain.EnvironmentType(creds.Environment),
	// 		MerchantAppID: creds.MerchantAppID,
	// 		FabricAppID:   creds.FabricAppID,
	// 		ShortCode:     creds.ShortCode,
	// 		AppSecret:     creds.AppSecret,
	// 		PrivateKey:    creds.PrivateKey,
	// 		PublicKey:     creds.PublicKey,
	// 	})
	// }
	// for _, pc := range miniapp.ProductCode {
	// 	product_code = append(product_code, domain.ProductCode{
	// 		BranchType:  domain.BranchType(pc.BranchType),
	// 		ProductCode: pc.ProductCode,
	// 	})
	// // }

	_, err := o.MongoDalMiniApp.InsertOne(ctx, miniapp)

	// 	model.MiniApp{
	// 	AppName:           miniapp.AppName,
	// 	AppIcon:           miniapp.AppIcon,
	// 	CommisonGLAccount: miniapp.CommisonGLAccount,
	// 	AppType: domain.AppType{
	// 		UAT:        miniapp.AppType.UAT,
	// 		Production: miniapp.AppType.Production,
	// 		Test:       miniapp.AppType.Test,
	// 		Dev:        miniapp.AppType.Dev,
	// 	},
	// 	MerchantID:     miniapp.MerchantID,
	// 	ProductCode:    product_code,
	// 	Credential:     credential,
	// 	IsEventMiniApp: miniapp.IsEventMiniApp,
	// 	IsThreeClick:   miniapp.IsThreeClick,
	// 	Enabled:        miniapp.Enabled,
	// 	IsDeleted:      miniapp.IsDeleted,
	// 	CreatedAt:      time.Now(),
	// 	LastModifiedAt: time.Now(),
	// 	DeletedAt:      time.Time{},
	// })
	if err != nil {
		return model.CPSAction{}, err
	}
	return model.CPSAction{}, nil
}

func (o *MiniAppPersistence) GetMiniAppActionId(ctx context.Context, action_id string) (model.CPSAction, error) {
	filter := map[string]interface{}{"action_code": action_id}
	data, err := o.MongoDalCPSAction.FindOne(ctx, filter, nil)
	if err != nil {
		return model.CPSAction{}, err
	}
	result := model.CPSAction{
		ID:              data.ID,
		ActionCode:      data.ActionCode,
		Department:      data.Department,
		RejectionReason: data.RejectionReason,
		PreviousAction:  data.PreviousAction,
		CurrentAction:   data.CurrentAction,
		ActionStatus:    data.ActionStatus,
		ActionType:      data.ActionType,
		RequestAction:   data.RequestAction,
		CreatedAt:       data.CreatedAt,
		LastModifiedAt:  data.LastModifiedAt,
	}
	return result, nil
}

func (o *MiniAppPersistence) UpdateCpsAction(ctx context.Context, action model.CPSAction) (model.CPSAction, error) {
	filter := map[string]interface{}{"action_code": action.ActionCode}
	update := map[string]interface{}{
		"checker_id":           action.CheckerID,
		"checker_name":         action.CheckerName,
		"checker_phone_number": action.CheckerPhoneNumber,
		"department":           action.Department,
		"rejection_reason":     action.RejectionReason,
		// "previos_action": func() json.RawMessage {
		// 	b, _ := json.Marshal(action.PreviousAction)
		// 	return b
		// }(),
		// "current_action":      action.CurrentAction,
		"action_status":       action.ActionStatus,
		"action_type":         action.ActionType,
		"request_action":      action.RequestAction,
		"last_modified_at":    time.Now(),
		"checker_action_time": time.Now(),
	}
	updatedAction, err := o.MongoDalCPSAction.UpdateOne(ctx, filter, update)
	return updatedAction, err
}

func (o *MiniAppPersistence) ListMiniApp(ctx context.Context) ([]*model.MiniApp, error) {
	filter := map[string]interface{}{"is_deleted": false}
	miniApps, err := o.MongoDalMiniApp.FindAll(ctx, filter, nil)
	if err != nil {
		return nil, err
	}
	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++")
	fmt.Printf("List of mini Apps %v", miniApps)
	fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++")

	return miniApps, nil
}

func (o *MiniAppPersistence) DetailMiniAppByID(ctx context.Context, id string) (model.MiniApp, error) {

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return model.MiniApp{}, err
	}
	filter := bson.M{"_id": objectID}
	miniApp, err := o.MongoDalMiniApp.FindOne(ctx, filter, nil)
	if err != nil {
		return model.MiniApp{}, err
	}
	return *miniApp, nil
}
