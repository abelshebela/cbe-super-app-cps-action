package password_rule

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
)

type PasswordRuleStorage struct {
	dal    dal.MongoDal[model.PasswordRule, model.PasswordRule]
	client *mongo.Client
	logger utils.Logger
}

func NewPasswordRuleRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.PasswordRuleRepository {
	return &PasswordRuleStorage{
		dal:    dal.NewMongoDal[model.PasswordRule, model.PasswordRule](client, dbName, collection),
		client: client,
		logger: logger,
	}
}

func (p *PasswordRuleStorage) Create(ctx context.Context, rule *model.PasswordRule) error {
	_, err := p.dal.InsertOne(ctx, *rule)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (p *PasswordRuleStorage) Update(ctx context.Context, id string, rule *model.PasswordRule) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		p.logger.Errorf("Failed to convert id to ObjectID: %s, error: %v", id, err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	updateData := PasswordRuleMapper(*rule)

	_, err = p.dal.UpdateOne(ctx, filter, updateData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			p.logger.Errorf("Failed to update password rule, id: %s, error: %v", id, err)
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		p.logger.Errorf("UpdateOne failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (p *PasswordRuleStorage) Delete(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	return p.dal.DeleteOne(ctx, filter)
}

func (p *PasswordRuleStorage) FindByID(ctx context.Context, id string) (*model.PasswordRule, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}

	result, err := p.dal.FindOne(ctx, filter, nil)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (p *PasswordRuleStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.PasswordRule], error) {
	filter := bson.M{
		"is_deleted": false,
	}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		filter["name"] = searchRegex
	}

	skip := int64((filterParam.Page - 1) * filterParam.PerPage)
	limit := int64(filterParam.PerPage)

	data, err := p.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		return nil, err
	}

	total, err := p.dal.TotalCount(ctx, filter)
	if err != nil {
		return nil, err
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]*model.PasswordRule]{
		Data: data,
		Meta: meta,
	}, nil
}

func (p *PasswordRuleStorage) FindCurrentRule(ctx context.Context) (*model.PasswordRule, error) {
	ruleModel, err := p.dal.FindOne(ctx, bson.M{}, bson.M{})
	if err != nil || ruleModel == nil {
		return nil, err
	}

	return ruleModel, nil
}
