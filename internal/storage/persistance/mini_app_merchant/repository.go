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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MiniAppMerchantStorage struct {
	dal    dal.MongoDal[model.MiniAppMerchant, model.MiniAppMerchant]
	client *mongo.Client
	logger utils.Logger
}

func NewMiniAppMerchantRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.MiniAppMerchantRepository {
	return &MiniAppMerchantStorage{
		dal:    dal.NewMongoDal[model.MiniAppMerchant, model.MiniAppMerchant](client, dbName, collection),
		client: client,
		logger: logger,
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

func (m *MiniAppMerchantStorage) Update(ctx context.Context, id string, merchant *model.MiniAppMerchant) (*model.MiniAppMerchant, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	updateData := MiniAppMerchantMapper(*merchant)

	_, err = m.dal.UpdateOne(ctx, filter, bson.M{"$set": updateData})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(localization.ErrorFileNotFound.Code)
		}
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// Fetch the updated document and return it
	updatedMerchant, err := m.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return updatedMerchant, nil
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
func (m *MiniAppMerchantStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"enabled": enable}}
	_, err = m.dal.UpdateOne(ctx, filter, update)
	return err
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

func (m *MiniAppMerchantStorage) Exists(
	ctx context.Context,
	data *model.CheckMiniAppMerchant,
	opts *model.MiniAppMerchantExistOptions,
) (bool, error) {
	filter := bson.M{
		"is_deleted": false,
		"$or":        []bson.M{},
	}

	if data.BankAccountNumber != "" {
		filter["$or"] = append(filter["$or"].([]bson.M), bson.M{"bank_account_number": data.BankAccountNumber})
	}
	if data.Email != "" {
		filter["$or"] = append(filter["$or"].([]bson.M), bson.M{"email": data.Email})
	}
	if data.PhoneNumber != "" {
		filter["$or"] = append(filter["$or"].([]bson.M), bson.M{"phone_number": data.PhoneNumber})
	}

	if len(filter["$or"].([]bson.M)) == 0 {
		return false, nil
	}

	if opts != nil && opts.ExcludeID != "" {
		objID, err := primitive.ObjectIDFromHex(opts.ExcludeID)
		if err != nil {
			m.logger.Errorf("invalid exclude ID: %v", err)
			return false, errors.New(localization.ErrorInvalidID.Code)
		}
		filter["_id"] = bson.M{"$ne": objID}
	}

	count, err := m.dal.TotalCount(ctx, filter)
	if err != nil {
		m.logger.Errorf("failed to check merchant info exists: %v", err)
		return false, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return count > 0, nil
}
