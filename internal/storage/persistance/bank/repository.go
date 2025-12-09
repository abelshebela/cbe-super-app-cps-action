package bank

import (
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
	"fmt"
	"regexp"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type BankStorage struct {
	dal    dal.MongoDal[model.Bank, model.Bank]
	client *mongo.Client
	logger utils.Logger
}

func NewBankRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.BankRepository {
	return &BankStorage{
		dal:    dal.NewMongoDal[model.Bank, model.Bank](client, dbName, collection),
		client: client,
		logger: logger,
	}
}

func (b *BankStorage) Create(ctx context.Context, bank *model.Bank) error {
	bank.ID = bson.NewObjectID()
	_, err := b.dal.InsertOne(ctx, *bank)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (b *BankStorage) Update(ctx context.Context, id string, bank *model.Bank) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	updateData := BankMapper(*bank)

	_, err = b.dal.UpdateOne(ctx, filter, updateData)
	if err != nil {
		b.logger.Errorf("failed to update bank error: %v", err)
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (b *BankStorage) Delete(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	return b.dal.DeleteOne(ctx, filter)
}

func (b *BankStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID}
	update := bson.M{"enabled": enable}
	_, err = b.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
	}
	return nil
}

func (b *BankStorage) FindByID(ctx context.Context, id string) (*model.Bank, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID}

	result, err := b.dal.FindOne(ctx, filter, nil)
	if err != nil {
		return nil, err
	}
	if result.IsDeleted {
		return nil, mongo.ErrNoDocuments
	}
	return result, nil
}

func (s *BankStorage) FindByNameOrBICOrCode(
	ctx context.Context,
	bic, code, name string, account_length *int) (*model.Bank, error) {

	// Build conditions dynamically, only for non-empty parameters
	conditions := []bson.M{}

	if name != "" {
		conditions = append(conditions, bson.M{
			"name": bson.M{"$regex": "^" + regexp.QuoteMeta(name) + "$", "$options": "i"},
		})
	}

	if bic != "" {
		conditions = append(conditions, bson.M{
			"bic": bson.M{"$regex": "^" + regexp.QuoteMeta(bic) + "$", "$options": "i"},
		})
	}

	if code != "" {
		conditions = append(conditions, bson.M{
			"code": bson.M{"$regex": "^" + regexp.QuoteMeta(code) + "$", "$options": "i"},
		})
	}

	// If no conditions provided, return error or handle appropriately
	if len(conditions) == 0 {
		if account_length != nil {
			return nil, nil
		}
		return nil, errors.New("at least one search parameter must be provided")
	}

	filter := bson.M{"$or": conditions}
	return s.dal.FindOne(ctx, filter, nil)
}

func (s *BankStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Bank], error) {
	filter := bson.M{"is_deleted": false}
	searchKeys := bson.M{}
	allowedKeys := []string{"search", "branch_code", "branch_name", "enabled", "enabled", "is_deleted"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"name": searchRegex},
			{"code": searchRegex},
			{"bic": searchRegex},
			{"created_at": searchRegex},
			{"last_modified_at": searchRegex},
		}

	}
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	data, err := s.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	total, err := s.dal.TotalCount(ctx, filter)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]*model.Bank]{
		Data: data,
		Meta: meta,
	}, nil
}
