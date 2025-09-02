package mini_app_merchant

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
	"fmt"
	"time"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (m *MiniAppMerchantStorage) Create(ctx context.Context, merchant *model.MiniAppMerchant) (*model.MiniAppMerchant, error) {
	if merchant.ID.IsZero() {
		merchant.ID = bson.ObjectID(primitive.NewObjectID())
	}
	createdMerchant, err := m.dal.InsertOne(ctx, *merchant) // returns struct
	if err != nil {
		return nil, err
	}
	return &createdMerchant, nil
}

func (w *MiniAppMerchantStorage) Update(ctx context.Context, id string, miniAppMerchant *model.MiniAppMerchant) error {
	var update bson.M
	objID, err := bson.ObjectIDFromHex(id)

	if err != nil {
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	update = MiniAppMerchantMapper(*miniAppMerchant)

	if len(update) == 2 {
		return errors.New(localization.ErrorNoDataProvided.Code)
	}

	_, err = w.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New(localization.ErrorWalletNotFound.Code)
		}
		w.logger.Errorf("Failed to update miniapp merchant: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (w *MiniAppMerchantStorage) Delete(ctx context.Context, id string) error {
	fmt.Println("We are in in there in repo")
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	update := bson.M{"is_deleted": true, "deleted_at": time.Now()}

	_, err = w.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return err
		}
		w.logger.Errorf("Failed to delete miniapp merchant: %v", err)
		return err
	}
	return nil
}

// EnableOrDisable toggles the Mini App Merchant's active status.
func (m *MiniAppMerchantStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	m.logger.Debugf("EnableOrDisable called with id=%s, enable=%v", id, enable)

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		m.logger.Debugf("Invalid ObjectID: %s, error: %v", id, err)
		return err
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	update := bson.M{"enabled": enable, "last_modified_at": time.Now()}

	m.logger.Debugf("EnableOrDisable filter: %+v, update: %+v", filter, update)

	_, err = m.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			m.logger.Warnf("Mini App Merchant ID %s not found for enable/disable", id)
			return err
		}
		m.logger.Errorf("Failed to enable/disable Mini App Merchant ID %s: %v", id, err)
		return err
	}
	m.logger.Debugf("EnableOrDisable succeeded for id=%s", id)
	return nil
}

func (s *MiniAppMerchantStorage) FindByID(ctx context.Context, id string) (*model.MiniAppMerchant, error) {
	idObj, ok := local_util.StringToObjectID(id)
	if !ok {
		s.logger.Errorf("Invalid ObjectID for fetch by id: %s", id)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": idObj, "is_deleted": false}
	projection := bson.M{}
	result, err := s.dal.FindOne(ctx, filter, projection)
	if err != nil {
		s.logger.Errorf("Error finding CPSAction: %v", err)
		code, _ := local_util.HandleMongoError(err)
		return nil, errors.New(code)
	}
	s.logger.Infof("Successfully found ServiceDetails: %+v", result)
	return result, nil

}

func (m *MiniAppMerchantStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.MiniAppMerchant], error) {
	filter := bson.M{"is_deleted": nil}

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

func (m *MiniAppMerchantStorage) Exists(ctx context.Context, data *model.CheckMiniAppMerchant, opts *model.MiniAppMerchantExistOptions) (*model.MiniAppMerchant, error) {
	if data.BankAccountNumber == "" {
		return nil, errors.New("bank account number is required")
	}

	filter := bson.M{
		"is_deleted":          nil,
		"bank_account_number": data.BankAccountNumber,
	}

	// Exclude a specific ID (useful for updates)
	if opts != nil && opts.ExcludeID != "" {
		objID, err := primitive.ObjectIDFromHex(opts.ExcludeID)
		if err != nil {
			m.logger.Errorf("invalid exclude ID: %v", err)
			return nil, errors.New(localization.ErrorInvalidID.Code)
		}
		filter["_id"] = bson.M{"$ne": objID}
	}

	result, err := m.dal.FindOne(ctx, filter, nil)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // not found is fine
		}
		m.logger.Errorf("failed to check merchant by bank account number: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return result, nil
}

func (p *MiniAppMerchantStorage) AddMiniApp(ctx context.Context, merchantID string, miniApp model.MiniApps) error {
	p.logger.Infof("AddMiniApp: merchant_id=%s, mini_app_id=%s", merchantID, miniApp.ID)

	objID, err := bson.ObjectIDFromHex(merchantID)
	if err != nil {
		p.logger.Warnf("AddMiniApp: invalid merchant_id=%s, error=%v", merchantID, err)
		return errors.New(localization.ErrorMerchantNotFound.Code)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}

	// First make sure mini_apps is initialized as an array if it’s null
	_, initErr := p.client.Database(p.dbName).Collection(p.collection).UpdateOne(
		ctx,
		bson.M{"_id": objID, "mini_apps": bson.M{"$type": "null"}},
		bson.M{"$set": bson.M{"mini_apps": bson.A{}}},
	)
	if initErr != nil {
		p.logger.Debugf("AddMiniApp: no initialization needed for merchant_id=%s (mini_apps not null)", merchantID)
	}

	// Now safely push the new mini app
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
