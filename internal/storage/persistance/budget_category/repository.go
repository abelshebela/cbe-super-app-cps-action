package budget_category

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/kafka"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type BudgetCategoryStorage struct {
	budgetCategoryDal dal.MongoDal[model.BudgetCategory, model.BudgetCategory]
	client            *mongo.Client
	kafkaProducer     kafka.ClientOrchestrationProducer
	logger            utils.Logger
}

var _ storage.BudgetCategoryRepository = (*BudgetCategoryStorage)(nil)

func NewBudgetCategoryRepository(client *mongo.Client, cfg *config.VaultConfig, dbName string, collection string, kafkaProducer kafka.ClientOrchestrationProducer, logger utils.Logger) storage.BudgetCategoryRepository {
	return &BudgetCategoryStorage{
		budgetCategoryDal: dal.NewMongoDal[model.BudgetCategory, model.BudgetCategory](client, cfg, dbName, collection),
		client:            client,
		kafkaProducer:     kafkaProducer,
		logger:            logger,
	}
}

func (b *BudgetCategoryStorage) CreateBudgetCategory(ctx context.Context, budgetCategory *model.BudgetCategory) error {
	newBudgetCategory, err := b.budgetCategoryDal.InsertOne(ctx, *budgetCategory)
	if err != nil {
		b.logger.Errorf("failed to create budget category: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	b.kafkaProducer.PublishMessage(ctx, newBudgetCategory, string(constants.ClientOrchestrationBudgetCategoryTopic), string(constants.ClientOrchestrationBudgetCategoryTopic), "new budget category created")
	return nil
}

func (b *BudgetCategoryStorage) UpdateBudgetCategory(ctx context.Context, id string, budgetCategory *model.BudgetCategory) error {
	b.logger.Infof("[UpdateBudgetCategory] updating budget category for id: %s", id)
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		b.logger.Errorf("[UpdateBudgetCategory] invalid object id: %v", err)
		return errors.New(localization.ErrorPINInvalid.Code)
	}

	filter := bson.M{"_id": objectID, "is_deleted": false}
	updateData := BudgetCategoryMapper(*budgetCategory)
	updatedBudgetCategory, err := b.budgetCategoryDal.UpdateOne(ctx, filter, updateData)
	if err != nil {
		b.logger.Errorf("[UpdateBudgetCategory] failed to update budget category: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	b.kafkaProducer.PublishMessage(ctx, updatedBudgetCategory, string(constants.ClientOrchestrationBudgetCategoryTopic), string(constants.ClientOrchestrationBudgetCategoryTopic), "budget category updated")

	b.logger.Infof("[UpdateBudgetCategory] budget category updated successfully")
	return nil
}

func (b *BudgetCategoryStorage) FindBudgetCategoryByID(ctx context.Context, id string) (*model.BudgetCategory, error) {
	b.logger.Infof("[FindBudgetCategoryByID] fetching budget category by id: %s", id)
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		b.logger.Errorf("[FindBudgetCategoryByID] invalid object id: %v", err)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	budgetCategory, err := b.budgetCategoryDal.FindOne(ctx, bson.M{"_id": objectID, "is_deleted": false}, bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			b.logger.Errorf("[FindBudgetCategoryByID] budget category not found")
			return nil, errors.New(localization.ErrorFileNotFound.Code)
		}
		b.logger.Errorf("[FindBudgetCategoryByID] failed to fetch budget category: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	b.logger.Infof("[FindBudgetCategoryByID] budget category retrieved successfully")
	return budgetCategory, nil
}

func (b *BudgetCategoryStorage) FindAllBudgetCategories(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]model.BudgetCategory], error) {
	filter := bson.M{"is_deleted": false}
	allowedKeys := []string{"enabled", "name"}

	filter, skip, limit := lib.FilterBuilder(*filterParams, bson.M{}, allowedKeys)

	if filterParams.Search != "" {
		searchRegex := bson.M{"$regex": filterParams.Search, "$options": "i"}
		filter["$or"] = []bson.M{
			{"name": searchRegex},
		}
	}

	data, err := b.budgetCategoryDal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		b.logger.Errorf("[FindAllBudgetCategories] failed to fetch budget categories: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := b.budgetCategoryDal.TotalCount(ctx, filter)
	if err != nil {
		b.logger.Errorf("[FindAllBudgetCategories] failed to count budget categories: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParams.Page, filterParams.PerPage)
	b.logger.Infof("[FindAllBudgetCategories] retrieved %d budget categories", len(data))
	return &types.PaginatedResponse[[]model.BudgetCategory]{Data: data, Meta: meta}, nil
}

func (b *BudgetCategoryStorage) DeleteBudgetCategory(ctx context.Context, id string) error {
	b.logger.Infof("[DeleteBudgetCategory] deleting budget category for id: %s", id)
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		b.logger.Errorf("[DeleteBudgetCategory] invalid object id: %v", err)
		return errors.New(localization.ErrorPINInvalid.Code)
	}

	filter := bson.M{"_id": objectID, "is_deleted": false}
	update := bson.M{"is_deleted": true, "updated_at": time.Now()}
	_, err = b.budgetCategoryDal.UpdateOne(ctx, filter, update)
	if err != nil {
		b.logger.Errorf("[DeleteBudgetCategory] failed to delete budget category: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	b.logger.Infof("[DeleteBudgetCategory] budget category deleted successfully")
	return nil
}

func (b *BudgetCategoryStorage) EnableOrDisableBudgetCategory(ctx context.Context, id string, enable bool) error {
	b.logger.Infof("[EnableOrDisableBudgetCategory] processing budget category enable/disable for id: %s, enabled: %v", id, enable)
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		b.logger.Errorf("[EnableOrDisableBudgetCategory] invalid object id: %v", err)
		return errors.New(localization.ErrorPINInvalid.Code)
	}

	filter := bson.M{"_id": objectID, "is_deleted": false}
	update := bson.M{"enabled": enable, "updated_at": time.Now()}
	updatedBudgetCategory, err := b.budgetCategoryDal.UpdateOne(ctx, filter, update)
	if err != nil {
		b.logger.Errorf("[EnableOrDisableBudgetCategory] failed to enable/disable budget category: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	b.kafkaProducer.PublishMessage(ctx, updatedBudgetCategory, string(constants.ClientOrchestrationBudgetCategoryTopic), string(constants.ClientOrchestrationBudgetCategoryTopic), "budget category enable/disable updated")

	b.logger.Infof("[EnableOrDisableBudgetCategory] budget category enable/disable completed successfully")
	return nil
}

func (b *BudgetCategoryStorage) FindByName(ctx context.Context, name string) (*model.BudgetCategory, error) {
	b.logger.Infof("[FindByName] searching for budget category by name")

	filter := bson.M{
		"name": bson.M{
			"$regex":   "^" + regexp.QuoteMeta(name) + "$", // exact match, case-insensitive
			"$options": "i",
		},
		"is_deleted": false,
	}

	result, err := b.budgetCategoryDal.FindOne(ctx, filter, nil)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			b.logger.Infof("[FindByName] budget category not found")
			return nil, nil
		}
		b.logger.Errorf("[FindByName] failed to find budget category: %v", err)
		return nil, fmt.Errorf("%s", "ERROR_PRODUCT_CODE_DATABASE_QUERY_FAILED")
	}
	b.logger.Infof("[FindByName] budget category retrieved successfully")
	return result, nil
}
