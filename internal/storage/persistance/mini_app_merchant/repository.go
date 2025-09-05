package mini_app_merchant

import (
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
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

func (m *MiniAppMerchantStorage) Update(ctx context.Context, id string, merchant *model.MiniAppMerchant)  error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{"_id": objID, "is_deleted": nil}
	updateData := MiniAppMerchantMapper(*merchant)

	_, err = m.dal.UpdateOne(ctx, filter, bson.M{"$set": updateData})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return  errors.New(localization.ErrorFileNotFound.Code)
		}
		return  errors.New(localization.ErrorUnexpectedError.Code)
	}

	return  nil
}

func (m *MiniAppMerchantStorage) Delete(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objID}
	err = m.dal.DeleteOne(ctx, filter)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return err
}

// EnableOrDisable toggles the merchant's active status.
// EnableOrDisable toggles the Mini App Merchant's active status.
func (m *MiniAppMerchantStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID, "is_deleted": nil}
	update := bson.M{"enabled": enable, "last_modified_at": time.Now()}

	_, err = m.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			m.logger.Warnf("Mini App Merchant ID %s not found for enable/disable", id)
			return errors.New(localization.ErrorMiniAppMerchantNotFound.Code)
		}
		m.logger.Errorf("Failed to enable/disable Mini App Merchant ID %s: %v", id, err)
		return errors.New(localization.ErrorUnexpectedError.Code)
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

func (s *MiniAppMerchantStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.MiniAppMerchant], error) {
	// 1. Base filter (only active records)
	filter := bson.M{"is_deleted": false}
	searchKeys := bson.M{}

	// 2. Allowed filterable/searchable fields
	allowedKeys := []string{"merchant_type"}

	// 3. Add search (if provided)
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["merchant_id"] = searchRegex // choose your searchable field(s)
	}

	// 4. Build filter, skip, limit
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	// 5. Fetch data
	data, err := s.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	// 6. Count total
	total, err := s.dal.TotalCount(ctx, filter)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	// 7. Build pagination metadata
	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	// 8. Return standard paginated response
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
