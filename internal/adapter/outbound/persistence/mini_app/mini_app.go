package miniapp

import (
	"context"
	"errors"
	"fmt"
	"time"

	dal "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/infra"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/mappers"
	model "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
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
		return nil, fmt.Errorf(common_util.GeneralDBQueryFailed)
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

	objectID, err := common_util.ParsePrimitiveObjectID(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objectID, "is_deleted": false}
	miniApp, err := o.MongoDalMiniApp.FindOne(ctx, filter, nil)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf(common_util.NotFound)
		}

		return nil, fmt.Errorf(common_util.GeneralDBQueryFailed)
	}
	res := mappers.ToDomainMiniApp(*miniApp)
	return &res, nil
}

func (o *MiniAppPersistence) DeleteMiniAppAction(ctx context.Context, action *miniApp_domain.MiniApp) (*miniApp_domain.MiniApp, error) {

	objID, err := bson.ObjectIDFromHex(action.ID)
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

func (o *MiniAppPersistence) UpdateMinApp(ctx context.Context, miniApp *miniApp_domain.MiniApp) (*miniApp_domain.MiniApp, error) {
	objID, err := common_util.ParsePrimitiveObjectID(miniApp.ID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id":        objID,
		"is_deleted": false,
	}

	update := bson.M{}
	if miniApp.AppName != "" {
		update["app_name"] = miniApp.AppName
	}
	if miniApp.AppIcon != "" {
		update["app_icon"] = miniApp.AppIcon
	}
	if miniApp.CommissionGLAccount != "" {
		update["commison_gl_account"] = miniApp.CommissionGLAccount
	}
	if miniApp.AppType != "" {
		update["app_type"] = miniApp.AppType
	}
	if len(miniApp.ProductCode) > 0 {
		update["product_code"] = miniApp.ProductCode
	}
	if len(miniApp.Credential) > 0 {
		update["credential"] = miniApp.Credential
	}

	update["is_event_mini_app"] = miniApp.IsEventMiniApp
	update["is_three_click"] = miniApp.IsThreeClick
	update["last_modified_at"] = time.Now()

	if len(update) == 1 {
		return nil, fmt.Errorf(common_util.NoDataProvidedForUpdate)
	}

	mini, err := o.MongoDalMiniApp.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf(common_util.NotFound)
		}
		o.logger.Warnf(err.Error(), "while updating")
		return nil, fmt.Errorf(common_util.GeneralDBUpdateFailed)
	}

	res := mappers.ToDomainMiniApp(mini)
	return &res, nil
}

func (o *MiniAppPersistence) CreateMiniApp(ctx context.Context, action *miniApp_domain.MiniApp) (*miniApp_domain.MiniApp, error) {

	miniAppDoc, err := mappers.ToModelMiniApp(action)
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

func (o *MiniAppPersistence) EnableDisableMiniApp(ctx context.Context, id string, enabled bool) (*miniApp_domain.MiniApp, error) {
	objID, err := common_util.ParsePrimitiveObjectID(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id":        objID,
		"is_deleted": false,
	}

	fmt.Println(enabled, "Enabled")
	update := bson.M{
		"enabled":          enabled,
		"last_modified_at": time.Now(),
	}

	miniApp, err := o.MongoDalMiniApp.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			o.logger.Errorf("MiniApp with ID %s not found for enable/disable", id)
			return nil, fmt.Errorf(common_util.NotFound)
		}

		o.logger.Errorf("Failed to update enabled state for MiniApp ID %s: %v", id, err)
		return nil, fmt.Errorf(common_util.GeneralDBUpdateFailed)
	}

	res := mappers.ToDomainMiniApp(miniApp)
	return &res, nil
}
