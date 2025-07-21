package miniapp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	dal "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/infra"
	model "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	miniApp "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/mini_app"


	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type MiniAppPersistence struct {
	MongoDalMiniApp   dal.MongoDal[model.MiniApp, model.MiniApp]
	MongoDalCPSAction dal.MongoDal[model.CPSAction, model.CPSAction]
	logger            utils.Logger
}

func InitMiniAppPersistence(client *mongo.Client, DB_name string, collections []string, logger utils.Logger) miniApp.MiniRepository {
	return &MiniAppPersistence{
		MongoDalMiniApp:   dal.NewMongoDal[model.MiniApp, model.MiniApp](client, DB_name, collections[0]),
		MongoDalCPSAction: dal.NewMongoDal[model.CPSAction, model.CPSAction](client, DB_name, collections[1]),
		logger:            logger,
	}
}

var _ miniApp.MiniRepository = (*MiniAppPersistence)(nil)

func (o *MiniAppPersistence) CreateMiniAppAction(ctx context.Context, Action *entities.CPSAction) (*entities.CPSAction, error) {
	CpsAction := model.CPSAction{
		ID:                 bson.NewObjectID(),
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
		ActionType:         string(Action.ActionType),
		RequestAction:      string(Action.RequestAction),
		CreatedAt:          time.Now(),
		LastModifiedAt:     time.Now(),
	}

	data, err := o.MongoDalCPSAction.InsertOne(ctx, CpsAction)
	if err != nil {
		return nil, err
	}

	res := entities.ToDomainCPSAction(&data)
	return res, nil
}

func (o *MiniAppPersistence) DeleteMiniAppAction(ctx context.Context, id string) (*model.MiniApp, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf(common_util.InvalidID)
	}
	update := bson.M{"is_deleted": true, "deleted_at": time.Now()}
	filter := bson.M{"_id": objID, "is_deleted": false}
	res, err := o.MongoDalMiniApp.UpdateOne(ctx, filter, update)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("NOT_FOUND")
		}

		return nil, fmt.Errorf("GENERAL_DB_UPDATE_FAILED")
	}
	return &res, nil
}

func (o *MiniAppPersistence) UpdateMinApp(ctx context.Context, action *entities.CPSAction, id string) (*model.MiniApp, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf(common_util.InvalidID)
	}

	actionData, err := o.ExtractActionData(action.CurrentAction)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id":        objID,
		"is_deleted": false,
	}

	update := bson.M{
		"app_name":            actionData.AppName,
		"app_icon":            actionData.AppIcon,
		"commison_gl_account": actionData.CommisonGLAccount,
		"app_type":            actionData.AppType,
		"product_code":        actionData.ProductCode,
		"credential":          actionData.Credential,
		"is_event_mini_app":   actionData.IsEventMiniApp,
		"is_three_click":      actionData.IsThreeClick,
		"last_modified_at":    time.Now(),
	}

	mini, err := o.MongoDalMiniApp.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("NOT_FOUND")
		}

		o.logger.Warnf(err.Error(), "while updating")
		return nil, fmt.Errorf("GENERAL_DB_UPDATE_FAILED")
	}

	return &mini, nil
}

func (o *MiniAppPersistence) CreateMiniApp(ctx context.Context, action *entities.CPSAction) (*model.MiniApp, error) {
	actionData, err := o.ExtractActionData(action.CurrentAction)
	if err != nil {
		return nil, err
	}

	req := model.MiniApp{
		ID:                bson.NewObjectID(),
		AppName:           actionData.AppName,
		AppIcon:           actionData.AppIcon,
		CommisonGLAccount: actionData.CommisonGLAccount,
		AppType: model.AppType{
			UAT:        actionData.AppType.UAT,
			Production: actionData.AppType.Production,
			Test:       actionData.AppType.Test,
			Dev:        actionData.AppType.Dev,
		},
		MerchantID: actionData.MerchantID,
		ProductCode: func() []model.ProductCode {
			var result []model.ProductCode
			for _, p := range actionData.ProductCode {
				result = append(result, model.ProductCode{
					ID:          p.ID,
					BranchType:  p.BranchType,
					ProductCode: p.ProductCode,
				})
			}
			return result
		}(),
		Credential: func() []model.CredentialInformation {
			var result []model.CredentialInformation
			for _, c := range actionData.Credential {
				result = append(result, model.CredentialInformation{
					ID:            c.ID,
					Environment:   c.Environment,
					MerchantAppID: c.MerchantAppID,
					FabricAppID:   c.FabricAppID,
					ShortCode:     c.ShortCode,
					AppSecret:     c.AppSecret,
					PrivateKey:    c.PrivateKey,
					PublicKey:     c.PublicKey,
				})
			}
			return result
		}(),
		IsEventMiniApp: actionData.IsEventMiniApp,
		IsThreeClick:   actionData.IsThreeClick,
		Enabled:        actionData.Enabled,
		IsDeleted:      actionData.IsDeleted,
		CreatedAt:      time.Now(),
		LastModifiedAt: time.Now(),
		DeletedAt:      time.Time{},
	}

	miniApp, err := o.MongoDalMiniApp.InsertOne(ctx, req)

	if err != nil {
		return nil, err
	}
	return &miniApp, nil
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
		"action_status":        action.ActionStatus,
		"action_type":          action.ActionType,
		"request_action":       action.RequestAction,
		"last_modified_at":     time.Now(),
		"checker_action_time":  time.Now(),
	}
	updatedAction, err := o.MongoDalCPSAction.UpdateOne(ctx, filter, update)
	return updatedAction, err
}

func (o *MiniAppPersistence) ListMiniApp(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*model.MiniApp], error) {
	filter := map[string]interface{}{"is_deleted": false}

	if filterParams.Filters != "" {
		blocked, err := strconv.ParseBool(filterParams.Filters)
		if err == nil {
			filter["is_blocked"] = blocked
		}
	}
	skip := (filterParams.Page - 1) * filterParams.PerPage
	limit := filterParams.PerPage
	miniApps, err := o.MongoDalMiniApp.FindAllWithPagination(ctx, filter, bson.M{}, int64(skip), int64(limit))
	if err != nil {
		return nil, err
	}

	total, err := o.MongoDalMiniApp.TotalCount(ctx, bson.M{})
	meta := common_util.BuildPaginationMeta(total, limit, filterParams.Page)

	return &common_util.PaginatedResponse[[]*model.MiniApp]{
		Data: miniApps,
		Meta: meta,
	}, nil

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

func (o *MiniAppPersistence) ExtractActionData(currentAction interface{}) (*model.MiniApp, error) {

	var actionData model.MiniApp

	// First try BSON Marshal + Unmarshal
	data, err := bson.Marshal(currentAction)
	if err == nil {
		if err := bson.Unmarshal(data, &actionData); err == nil {
			return &actionData, nil
		}
		o.logger.Warnf("BSON Unmarshal failed, falling back to JSON: %v", err)
	} else {
		o.logger.Warnf("BSON Marshal failed, falling back to JSON: %v", err)
	}

	// Fallback to JSON if BSON fails
	jsonData, err := json.Marshal(currentAction)
	if err != nil {
		o.logger.Errorf("JSON marshal failed: %v", err)
		return nil, fmt.Errorf("INVALID_ACTION_DATA")
	}

	if err := json.Unmarshal(jsonData, &actionData); err != nil {
		o.logger.Errorf("JSON unmarshal to MiniApp failed: %v", err)
		return nil, fmt.Errorf("INVALID_ACTION_DATA")
	}

	return &actionData, nil
}
