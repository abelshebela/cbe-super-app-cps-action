package mini_app_merchant

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MiniAppMerchantStorage struct {
	dal        dal.MongoDal[model.MiniAppMerchant, model.MiniAppMerchant]
	client     *mongo.Client
	logger     utils.Logger
	dbName     string
	collection string
}

func NewMiniAppMerchantRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.MiniAppMerchantRepository {
	return &MiniAppMerchantStorage{
		dal:        dal.NewMongoDal[model.MiniAppMerchant, model.MiniAppMerchant](client, dbName, collection),
		client:     client,
		logger:     logger,
		dbName:     dbName,
		collection: collection,
	}
}

func (m *MiniAppMerchantStorage) Create(ctx context.Context, merchant *model.MiniAppMerchant) error {
	_, err := m.dal.InsertOne(ctx, *merchant)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (m *MiniAppMerchantStorage) Update(ctx context.Context, id string, merchant *model.MiniAppMerchant) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	updateData := MiniAppMerchantMapper(*merchant)

	_, err = m.dal.UpdateOne(ctx, filter, updateData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (m *MiniAppMerchantStorage) Delete(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	return m.dal.DeleteOne(ctx, filter)
}

func (m *MiniAppMerchantStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"enabled": enable}}
	_, err = m.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
	}
	return nil
}

func (m *MiniAppMerchantStorage) FindByID(ctx context.Context, id string) (*model.MiniAppMerchant, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID}

	result, err := m.dal.FindOne(ctx, filter, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *MiniAppMerchantStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.MiniAppMerchant], error) {
	filter := bson.M{"is_deleted": false}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		filter["$or"] = []bson.M{
			{"name": searchRegex},
			{"merchant_id": searchRegex},
		}
	}

	skip := int64((filterParam.Page - 1) * filterParam.PerPage)
	limit := int64(filterParam.PerPage)

	data, err := m.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		return nil, err
	}

	total, err := m.dal.TotalCount(ctx, filter)
	if err != nil {
		return nil, err
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]*model.MiniAppMerchant]{
		Data: data,
		Meta: meta,
	}, nil
}

func (p *MiniAppMerchantStorage) DetailMiniAppByID(ctx context.Context, merchantID string) (*model.MiniAppMerchant, error) {
	objID, err := bson.ObjectIDFromHex(merchantID)
	if err != nil {
		p.logger.Errorf("invalid merchant ObjectID %s: %v", merchantID, err)
		return nil, errors.New(localization.ErrorMerchantNotFound.Code)
	}

	doc, err := p.dal.FindOne(ctx, bson.M{"_id": objID}, bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			p.logger.Warnf("no Mini App merchant found for ID: %s, %+v :", objID, err)
			return nil, errors.New(localization.ErrorMerchantNotFound.Code)
		}
		p.logger.Errorf("failed to get miniapp merchant by ID: %v", err)
		return nil, errors.New(localization.ErrorMerchantNotFound.Code)
	}

	return ToMiniAppMerchantDomain(doc), nil
}

func (p *MiniAppMerchantStorage) AddMiniApp(ctx context.Context, merchantID string, miniApp model.MiniApps) error {
	p.logger.Infof("AddMiniApp: merchant_id=%s, mini_app_id=%s", merchantID, miniApp.ID)

	objID, err := bson.ObjectIDFromHex(merchantID)
	if err != nil {
		p.logger.Warnf("AddMiniApp: invalid merchant_id=%s, error=%v", merchantID, err)
		return errors.New(localization.ErrorMerchantNotFound.Code)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}

	update := bson.M{"$push": bson.M{"mini_apps": miniApp}}
	var result model.MiniAppMerchant
	err = p.client.Database(p.dbName).Collection(p.collection).FindOneAndUpdate(
		ctx,
		filter,
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			p.logger.Warnf("AddMiniApp: merchant not found, merchant_id=%s", merchantID)
			return errors.New(localization.ErrorMerchantNotFound.Code)
		}
		p.logger.Warnf("AddMiniApp: failed to add mini app, merchant_id=%s, mini_app_id=%s, error=%v", merchantID, miniApp.ID, err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	p.logger.Infof("AddMiniApp: successfully added mini app, merchant_id=%s, mini_app_id=%s", merchantID, miniApp.ID)
	return nil
}

func (p *MiniAppMerchantStorage) UpdateMiniAppEnabledState(ctx context.Context, merchantID string, miniAppID string, enabled bool) error {
	p.logger.Infof("UpdateMiniAppEnabledState: merchant_id=%s, mini_app_id=%s, enabled=%v", merchantID, miniAppID, enabled)

	objID, err := bson.ObjectIDFromHex(merchantID)
	if err != nil {
		p.logger.Errorf("invalid merchant ObjectID %s: %v", merchantID, err)
		return errors.New(localization.ErrorMerchantNotFound.Code)
	}

	filter := bson.M{
		"_id":                  objID,
		"mini_apps.id":         miniAppID,
		"mini_apps.is_deleted": false,
		"is_deleted":           false,
	}
	update := bson.M{"mini_apps.$.enabled": enabled}

	_, err = p.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			p.logger.Warnf("UpdateMiniAppEnabledState: mini app not found, merchant_id=%s, mini_app_id=%s", merchantID, miniAppID)
			return errors.New(localization.ErrorMiniAppNotFound.Code)
		}
		p.logger.Warnf("UpdateMiniAppEnabledState: failed to update enabled state, merchant_id=%s, mini_app_id=%s, error=%v", merchantID, miniAppID, err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	p.logger.Infof("UpdateMiniAppEnabledState: successfully updated enabled=%v, merchant_id=%s, mini_app_id=%s", enabled, merchantID, miniAppID)
	return nil
}

func (p *MiniAppMerchantStorage) SoftDeleteMiniApp(ctx context.Context, merchantID string, miniAppID string) error {
	p.logger.Infof("SoftDeleteMiniApp: merchant_id=%s, mini_app_id=%s", merchantID, miniAppID)

	objID, err := bson.ObjectIDFromHex(merchantID)
	if err != nil {
		p.logger.Errorf("invalid merchant ObjectID %s: %v", merchantID, err)
		return errors.New(localization.ErrorMerchantNotFound.Code)
	}

	filter := bson.M{
		"_id":                  objID,
		"mini_apps.id":         miniAppID,
		"mini_apps.is_deleted": false,
		"is_deleted":           false,
	}
	update := bson.M{"mini_apps.$.is_deleted": true}

	_, err = p.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			p.logger.Warnf("SoftDeleteMiniApp: mini app not found, merchant_id=%s, mini_app_id=%s", merchantID, miniAppID)
			return errors.New(localization.ErrorMiniAppNotFound.Code)
		}
		p.logger.Warnf("SoftDeleteMiniApp: failed to soft delete mini app, merchant_id=%s, mini_app_id=%s, error=%v", merchantID, miniAppID, err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	p.logger.Infof("SoftDeleteMiniApp: successfully soft deleted mini app, merchant_id=%s, mini_app_id=%s", merchantID, miniAppID)
	return nil
}
