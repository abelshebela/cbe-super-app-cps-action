package miniapp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/mappers"
	dal "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/infra"
	model "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	miniApp_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp"
	miniApp "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/mini_app"

	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type MiniAppPersistence struct {
	MongoDalMiniApp dal.MongoDal[model.MiniApp, model.MiniApp]
	logger          utils.Logger
}

func InitMiniAppPersistence(client *mongo.Client, DB_name string, collections []string, logger utils.Logger) miniApp.MiniRepository {
	return &MiniAppPersistence{
		MongoDalMiniApp: dal.NewMongoDal[model.MiniApp, model.MiniApp](client, DB_name, collections[0]),
		logger:          logger,
	}
}

var _ miniApp.MiniRepository = (*MiniAppPersistence)(nil)

func (o *MiniAppPersistence) ListMiniApp(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*miniApp_domain.MiniApp], error) {
	filter := bson.M{"is_deleted": false}

	if filterParams.Search != "" {
		searchRegex := bson.M{"$regex": filterParams.Search, "$options": "i"}

		filter["$or"] = []bson.M{
			{"app_name": searchRegex},
			{"commison_gl_account": searchRegex},
			{"app_type.uat": searchRegex},
			{"app_type.production": searchRegex},
			{"app_type.test": searchRegex},
			{"app_type.dev": searchRegex},
			{"merchant_id": searchRegex},
			{"product_code.product_code": searchRegex},
		}
	}

	skip := (filterParams.Page - 1) * filterParams.PerPage
	limit := filterParams.PerPage

	miniAppsDocs, err := o.MongoDalMiniApp.FindAllWithPagination(ctx, filter, bson.M{}, int64(skip), int64(limit))
	if err != nil {
		return nil, err
	}

	var miniApps []*miniApp_domain.MiniApp
	for _, doc := range miniAppsDocs {
		converted := mappers.ToDomainMiniApp(*doc)
		miniApps = append(miniApps, &converted)
	}

	total, err := o.MongoDalMiniApp.TotalCount(ctx, filter)
	if err != nil {
		return nil, err
	}

	meta := common_util.BuildPaginationMeta(total, limit, filterParams.Page)

	return &common_util.PaginatedResponse[[]*miniApp_domain.MiniApp]{
		Data: miniApps,
		Meta: meta,
	}, nil
}

func (o *MiniAppPersistence) DetailMiniAppByID(ctx context.Context, id string) (*miniApp_domain.MiniApp, error) {

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("INVALID_ID")
	}
	filter := bson.M{"_id": objectID, "is_deleted": false}
	miniApp, err := o.MongoDalMiniApp.FindOne(ctx, filter, nil)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("NOT_FOUND")
		}

		return nil, fmt.Errorf("GENERAL_DB_QUERY_FAILED")
	}
	res := mappers.ToDomainMiniApp(*miniApp)
	return &res, nil
}

func (o *MiniAppPersistence) DeleteMiniAppAction(ctx context.Context, action *entities.CPSAction) (*miniApp_domain.MiniApp, error) {
	var actionData miniApp_domain.MiniApp
	bytes, err := json.Marshal(action.PreviousAction)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal CurrentAction: %w", err)
	}

	if err := json.Unmarshal(bytes, &actionData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal CurrentAction into Advert: %w", err)
	}

	objID, err := bson.ObjectIDFromHex(actionData.ID)
	if err != nil {
		return nil, fmt.Errorf(common_util.InvalidID)
	}
	update := bson.M{"is_deleted": true, "deleted_at": time.Now()}
	filter := bson.M{"_id": objID, "is_deleted": false}
	res, err := o.MongoDalMiniApp.UpdateOne(ctx, filter, update)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf(common_util.NotFound)
		}

		return nil, fmt.Errorf(common_util.GeneralDBUpdateFailed)
	}
	mappedApp := mappers.ToDomainMiniApp(res)
	return &mappedApp, nil
}

func (o *MiniAppPersistence) UpdateMinApp(ctx context.Context, action *entities.CPSAction) (*miniApp_domain.MiniApp, error) {

	var actionData miniApp_domain.MiniApp
	bytes, err := json.Marshal(action.CurrentAction)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal CurrentAction: %w", err)
	}

	if err := json.Unmarshal(bytes, &actionData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal CurrentAction into Advert: %w", err)
	}

	objID, err := bson.ObjectIDFromHex(actionData.ID)
	if err != nil {
		return nil, fmt.Errorf(common_util.InvalidID)
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
		return nil, fmt.Errorf(common_util.GeneralDBUpdateFailed)
	}

	res := mappers.ToDomainMiniApp(mini)

	return &res, nil
}

func (o *MiniAppPersistence) CreateMiniApp(ctx context.Context, action *entities.CPSAction) (*miniApp_domain.MiniApp, error) {
	var actionData miniApp_domain.MiniApp
	bytes, err := json.Marshal(action.CurrentAction)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal CurrentAction: %w", err)
	}

	if err := json.Unmarshal(bytes, &actionData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal CurrentAction into Advert: %w", err)
	}

	miniAppDoc, err := mappers.ToModelMiniApp(&actionData)
	if err != nil {
		return nil, err
	}

	miniApp, err := o.MongoDalMiniApp.InsertOne(ctx, *miniAppDoc)
	if err != nil {
		return nil, err
	}

	res := mappers.ToDomainMiniApp(miniApp)
	return &res, nil
}
