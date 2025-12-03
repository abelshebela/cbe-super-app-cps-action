package budget_category

import (
	"context"
	"errors"
	"time"

	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type BudgetCategoryStorage struct {
	budgetCategoryDal dal.MongoDal[model.BudgetCategory, model.BudgetCategory]
	client            *mongo.Client
	logger            utils.Logger
}

var _ storage.BudgetCategoryRepository = (*BudgetCategoryStorage)(nil)

func NewBudgetCategoryRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.BudgetCategoryRepository {
	return &BudgetCategoryStorage{
		budgetCategoryDal: dal.NewMongoDal[model.BudgetCategory, model.BudgetCategory](client, dbName, collection),
		client:            client,
		logger:            logger,
	}
}

func (b *BudgetCategoryStorage) CreateBudgetCategory(ctx context.Context, budgetCategory *model.BudgetCategory) error {
	_, err := b.budgetCategoryDal.InsertOne(ctx, *budgetCategory)
	if err != nil {
		b.logger.Errorf("failed to create budget category: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (b *BudgetCategoryStorage) UpdateBudgetCategory(ctx context.Context, id string, budgetCategory *model.BudgetCategory) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorPINInvalid.Code)
	}

	filter := bson.M{"_id": objectID, "is_deleted": false}
	updateData := BudgetCategoryMapper(*budgetCategory)
	_, err = b.budgetCategoryDal.UpdateOne(ctx, filter, updateData)
	if err != nil {
		b.logger.Errorf("failed to update budget category: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (b *BudgetCategoryStorage) FindBudgetCategoryByID(ctx context.Context, id string) (*model.BudgetCategory, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	budgetCategory, err := b.budgetCategoryDal.FindOne(ctx, bson.M{"_id": objectID, "is_deleted": false}, bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New(localization.ErrorFileNotFound.Code)
		}
		b.logger.Errorf("failed to fetch budget category: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return budgetCategory, nil
}

func (b *BudgetCategoryStorage) FindAllBudgetCategories(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.BudgetCategory], error) {
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
		b.logger.Errorf("failed to fetch budget categories: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := b.budgetCategoryDal.TotalCount(ctx, filter)
	if err != nil {
		b.logger.Errorf("failed counting budget categories: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParams.Page, filterParams.PerPage)
	return &types.PaginatedResponse[[]*model.BudgetCategory]{Data: data, Meta: meta}, nil
}

func (b *BudgetCategoryStorage) DeleteBudgetCategory(ctx context.Context, id string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorPINInvalid.Code)
	}

	filter := bson.M{"_id": objectID, "is_deleted": false}
	update := bson.M{"is_deleted": true, "updated_at": time.Now()}
	_, err = b.budgetCategoryDal.UpdateOne(ctx, filter, update)
	if err != nil {
		b.logger.Errorf("failed to delete budget category: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (b *BudgetCategoryStorage) EnableOrDisableBudgetCategory(ctx context.Context, id string, enable bool) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorPINInvalid.Code)
	}

	filter := bson.M{"_id": objectID, "is_deleted": false}
	update := bson.M{"enabled": enable, "updated_at": time.Now()}
	_, err = b.budgetCategoryDal.UpdateOne(ctx, filter, update)
	if err != nil {
		b.logger.Errorf("failed to enable/disable budget category: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}
