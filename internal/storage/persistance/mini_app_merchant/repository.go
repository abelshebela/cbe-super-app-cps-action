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
	createdMerchant, err := m.dal.InsertOne(ctx, *merchant)
	if err != nil {
		m.logger.Errorf("Failed to create mini app merchant: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return &createdMerchant, nil
}

func (m *MiniAppMerchantStorage) Update(ctx context.Context, id string, merchant *model.MiniAppMerchant) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		m.logger.Errorf("Invalid ID format: %s, error: %v", id, err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}
	updateData := MiniAppMerchantMapper(*merchant)

	_, err = m.dal.UpdateOne(ctx, filter, updateData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			m.logger.Warnf("Mini app merchant not found for update, id: %s", id)
			return errors.New(localization.ErrorMiniAppMerchantNotFound.Code)
		}
		m.logger.Errorf("Failed to update mini app merchant, id: %s, error: %v", id, err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (m *MiniAppMerchantStorage) Delete(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		m.logger.Errorf("Invalid ID format: %s, error: %v", id, err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}
	update := bson.M{
		"is_deleted":       true,
		"deleted_at":       time.Now(),
		"last_modified_at": time.Now(),
	}

	_, err = m.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			m.logger.Warnf("Mini app merchant not found for deletion, id: %s", id)
			return errors.New(localization.ErrorMiniAppMerchantNotFound.Code)
		}
		m.logger.Errorf("Failed to delete mini app merchant, id: %s, error: %v", id, err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (m *MiniAppMerchantStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		m.logger.Errorf("Invalid ID format: %s, error: %v", id, err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	update := bson.M{"enabled": enable, "last_modified_at": time.Now()}

	_, err = m.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			m.logger.Warnf("Mini app merchant not found for enable/disable, id: %s", id)
			return errors.New(localization.ErrorMiniAppMerchantNotFound.Code)
		}
		m.logger.Errorf("Failed to enable/disable mini app merchant, id: %s, error: %v", id, err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (m *MiniAppMerchantStorage) FindByID(ctx context.Context, id string) (*model.MiniAppMerchant, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		m.logger.Errorf("Invalid ID format: %s, error: %v", id, err)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	result, err := m.dal.FindOne(context.Background(), filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			m.logger.Warnf("Mini app merchant not found, id: %s", id)
			return nil, errors.New(localization.ErrorMiniAppMerchantNotFound.Code)
		}
		m.logger.Errorf("Failed to find mini app merchant, id: %s, error: %v", id, err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return result, nil
}

func (s *MiniAppMerchantStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.MiniAppMerchant], error) {
	searchKeys := bson.M{}
	allowedKeys := []string{"merchant_type", "merchant_code", "merchant_name", "email", "phone_number", "enabled"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"merchant_id": searchRegex},
			{"merchant_name": searchRegex},
			{"merchant_code": searchRegex},
			{"email": searchRegex},
			{"phone_number": searchRegex},
		}
	}

	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	filter["is_deleted"] = false

	data, err := s.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		s.logger.Errorf("Failed to fetch paginated mini app merchants: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := s.dal.TotalCount(ctx, filter)
	if err != nil {
		s.logger.Errorf("Failed to count total mini app merchants: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]*model.MiniAppMerchant]{
		Data: data,
		Meta: meta,
	}, nil
}

func (m *MiniAppMerchantStorage) Exists(ctx context.Context, data *model.CheckMiniAppMerchant, opts *model.MiniAppMerchantExistOptions) (bool, error) {
	if data == nil {
		return false, nil
	}

	var conditions []bson.M
	if data.BankAccountNumber != "" {
		conditions = append(conditions, bson.M{"bank_account_number": data.BankAccountNumber})
	}
	if data.Email != "" {
		conditions = append(conditions, bson.M{"kyc.representative.email": data.Email})
	}
	if data.PhoneNumber != "" {
		conditions = append(conditions, bson.M{"kyc.representative.phone": data.PhoneNumber})
	}

	if len(conditions) == 0 {
		return false, nil
	}

	filter := bson.M{
		"is_deleted": false,
		"$or":        conditions,
	}

	if opts != nil && opts.ExcludeID != "" {
		objID, err := primitive.ObjectIDFromHex(opts.ExcludeID)
		if err != nil {
			m.logger.Errorf("Invalid exclude ID: %v", err)
			return false, errors.New(localization.ErrorInvalidID.Code)
		}
		filter["_id"] = bson.M{"$ne": objID} // exclude self
	}

	coll := m.client.Database(m.dbName).Collection(m.collection)

	err := coll.FindOne(
		ctx, // ✅ use passed context
		filter,
		options.FindOne().SetProjection(bson.M{"_id": 1}),
	).Err()

	if err == mongo.ErrNoDocuments {
		return false, nil
	}
	if err != nil {
		m.logger.Errorf("Exists check failed: %v", err)
		return false, errors.New(localization.ErrorMiniAppMerchantExistsCheckFailed.Code)
	}

	return true, nil
}

func (p *MiniAppMerchantStorage) AddMiniApp(ctx context.Context, merchantID string, miniApp model.MiniApps) error {
	p.logger.Infof("AddMiniApp: merchant_id=%s, mini_app_id=%s", merchantID, miniApp.ID)

	// Validate input
	if merchantID == "" {
		p.logger.Errorf("Empty merchant ID provided")
		return errors.New(localization.ErrorInvalidID.Code)
	}

	// Use v2.x ObjectID conversion consistently
	objID, err := bson.ObjectIDFromHex(merchantID)
	if err != nil {
		p.logger.Errorf("Invalid merchant ID: %s, error: %v", merchantID, err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	
	// Diagnostic logging: Count documents matching the filter
	count, err := p.client.Database(p.dbName).Collection(p.collection).CountDocuments(ctx, filter)
	if err != nil {
		p.logger.Errorf("Failed to count documents for diagnostic, merchant_id=%s, error=%v", merchantID, err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	
	p.logger.Infof("Diagnostic: Found %d documents matching filter for merchant_id=%s", count, merchantID)
	
	if count == 0 {
		p.logger.Warnf("No merchant found with ID=%s and is_deleted=false", merchantID)
		return errors.New(localization.ErrorMiniAppMerchantNotFound.Code)
	}

	// Log the filter for debugging
	p.logger.Infof("Diagnostic: Filter used: %+v", filter)
	
	// Check if mini_apps field exists and is an array
	var existingDoc bson.M
	err = p.client.Database(p.dbName).Collection(p.collection).FindOne(ctx, filter).Decode(&existingDoc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			p.logger.Warnf("Merchant not found during field check, merchant_id=%s", merchantID)
			return errors.New(localization.ErrorMiniAppMerchantNotFound.Code)
		}
		p.logger.Errorf("Failed to check existing document, merchant_id=%s, error=%v", merchantID, err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	
	// Log the mini_apps field type and value
	if miniAppsField, exists := existingDoc["mini_apps"]; exists {
		p.logger.Infof("Diagnostic: mini_apps field exists, type: %T, value: %+v", miniAppsField, miniAppsField)
	} else {
		p.logger.Infof("Diagnostic: mini_apps field does not exist")
	}

	// Step 1: Ensure mini_apps is an array (only if needed)
	normalizeFilter := bson.M{
		"_id": objID, 
		"is_deleted": false,
		"$or": []bson.M{
			{"mini_apps": bson.M{"$exists": false}},
			{"mini_apps": nil},
			{"mini_apps": bson.M{"$not": bson.M{"$type": "array"}}}, // Handle non-array types
		},
	}
	
	normalizeResult, err := p.client.Database(p.dbName).Collection(p.collection).UpdateOne(
		ctx,
		normalizeFilter,
		bson.M{"$set": bson.M{"mini_apps": bson.A{}}},
	)
	if err != nil {
		p.logger.Errorf("Failed to normalize mini_apps array, merchant_id=%s, error=%v", merchantID, err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	
	if normalizeResult.ModifiedCount > 0 {
		p.logger.Infof("Diagnostic: Normalized mini_apps field for merchant_id=%s", merchantID)
	}

	// Step 2: Push the mini app using FindOneAndUpdate
	var result model.MiniAppMerchant
	err = p.client.Database(p.dbName).Collection(p.collection).FindOneAndUpdate(
		ctx,
		filter,
		bson.M{"$push": bson.M{"mini_apps": miniApp}},
		options.FindOneAndUpdate().
			SetReturnDocument(options.After).
			SetUpsert(false), // Explicitly set upsert to false
	).Decode(&result)
	
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			
			p.logger.Warnf("Merchant not found for adding mini app, merchant_id=%s", merchantID)
			return errors.New(localization.ErrorMiniAppMerchantNotFound.Code)
		}
		p.logger.Errorf("Failed to add mini app, merchant_id=%s, mini_app_id=%s, error=%v", merchantID, miniApp.ID, err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	p.logger.Infof("Successfully added mini app, merchant_id=%s, mini_app_id=%s", merchantID, miniApp.ID)
	return nil
}


func (p *MiniAppMerchantStorage) UpdateMiniAppEnabledState(ctx context.Context, merchantID string, miniAppID string, enabled bool) error {
	p.logger.Infof("UpdateMiniAppEnabledState: merchant_id=%s, mini_app_id=%s, enabled=%v", merchantID, miniAppID, enabled)

	objID, err := bson.ObjectIDFromHex(merchantID)
	if err != nil {
		p.logger.Errorf("Invalid merchant ID: %s, error: %v", merchantID, err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	miniAppObjID, err := bson.ObjectIDFromHex(miniAppID)
	if err != nil {
		p.logger.Errorf("Invalid mini app ID: %s, error: %v", miniAppID, err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{
		"_id":           objID,
		"is_deleted":    false,
		"mini_apps._id": miniAppObjID,
	}

	update := bson.M{
		"$set": bson.M{
			"mini_apps.$.enabled":          enabled,
			"mini_apps.$.last_modified_at": time.Now(),
		},
	}

	result, err := p.client.Database(p.dbName).Collection(p.collection).UpdateOne(ctx, filter, update)
	if err != nil {
		p.logger.Errorf("Failed to update mini app enabled state, merchant_id=%s, mini_app_id=%s, error=%v", merchantID, miniAppID, err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	if result.MatchedCount == 0 {
		p.logger.Warnf("Merchant or mini app not found for update, merchant_id=%s, mini_app_id=%s", merchantID, miniAppID)
		return errors.New(localization.ErrorMiniAppMerchantNotFound.Code)
	}

	p.logger.Infof("Successfully updated mini app enabled state, merchant_id=%s, mini_app_id=%s, enabled=%v", merchantID, miniAppID, enabled)
	return nil
}

func (p *MiniAppMerchantStorage) SoftDeleteMiniApp(ctx context.Context, merchantID string, miniAppID string) error {
	p.logger.Infof("SoftDeleteMiniApp: merchant_id=%s, mini_app_id=%s", merchantID, miniAppID)

	objID, err := bson.ObjectIDFromHex(merchantID)
	if err != nil {
		p.logger.Errorf("Invalid merchant ID: %s, error: %v", merchantID, err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	miniAppObjID, err := bson.ObjectIDFromHex(miniAppID)
	if err != nil {
		p.logger.Errorf("Invalid mini app ID: %s, error: %v", miniAppID, err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{
		"_id":           objID,
		"is_deleted":    false,
		"mini_apps._id": miniAppObjID,
	}

	update := bson.M{
		"$set": bson.M{
			"mini_apps.$.is_deleted":       true,
			"mini_apps.$.deleted_at":       time.Now(),
			"mini_apps.$.last_modified_at": time.Now(),
		},
	}

	result, err := p.client.Database(p.dbName).Collection(p.collection).UpdateOne(ctx, filter, update)
	if err != nil {
		p.logger.Errorf("Failed to soft delete mini app, merchant_id=%s, mini_app_id=%s, error=%v", merchantID, miniAppID, err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	if result.MatchedCount == 0 {
		p.logger.Warnf("Merchant or mini app not found for deletion, merchant_id=%s, mini_app_id=%s", merchantID, miniAppID)
		return errors.New(localization.ErrorMiniAppMerchantNotFound.Code)
	}

	p.logger.Infof("Successfully soft deleted mini app, merchant_id=%s, mini_app_id=%s", merchantID, miniAppID)
	return nil
}
