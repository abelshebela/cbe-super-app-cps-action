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
	dal        dal.MongoDal[model.PasswordRule, model.PasswordRule]
	client     *mongo.Client
	collection *mongo.Collection
	logger     utils.Logger
}

func NewPasswordRuleRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.PasswordRuleRepository {
	return &PasswordRuleStorage{
		dal:        dal.NewMongoDal[model.PasswordRule, model.PasswordRule](client, dbName, collection),
		client:     client,
		logger:     logger,
		collection: client.Database(dbName).Collection(collection),
	}
}

func (p *PasswordRuleStorage) Create(ctx context.Context, rule *model.PasswordRule) error {
	p.logger.Infof("[Create] creating password rule")
	_, err := p.dal.InsertOne(ctx, *rule)
	if err != nil {
		p.logger.Errorf("[Create] failed to create password rule: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	p.logger.Infof("[Create] password rule created successfully")
	return nil
}

func (p *PasswordRuleStorage) Update(ctx context.Context, id string, rule *model.PasswordRule) error {
	p.logger.Infof("[Update] updating password rule for id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		p.logger.Errorf("[Update] invalid object id: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	updateData := PasswordRuleMapper(*rule)

	_, err = p.dal.UpdateOne(ctx, filter, updateData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			p.logger.Errorf("[Update] password rule not found")
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		p.logger.Errorf("[Update] failed to update password rule: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	p.logger.Infof("[Update] password rule updated successfully")
	return nil
}

func (p *PasswordRuleStorage) Delete(ctx context.Context, id string) error {
	p.logger.Infof("[Delete] deleting password rule for id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		p.logger.Errorf("[Delete] invalid object id: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	err = p.dal.DeleteOne(ctx, filter)
	if err != nil {
		p.logger.Errorf("[Delete] failed to delete password rule: %v", err)
		return err
	}
	p.logger.Infof("[Delete] password rule deleted successfully")
	return nil
}

func (p *PasswordRuleStorage) FindByID(ctx context.Context, id string) (*model.PasswordRule, error) {
	p.logger.Infof("[FindByID] fetching password rule by id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		p.logger.Errorf("[FindByID] invalid object id: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}

	result, err := p.dal.FindOne(ctx, filter, nil)

	if err != nil {
		p.logger.Errorf("[FindByID] failed to find password rule: %v", err)
		return nil, err
	}
	p.logger.Infof("[FindByID] password rule retrieved successfully")
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
		p.logger.Errorf("[FindAllWithPagination] failed to fetch password rules: %v", err)
		return nil, err
	}

	total, err := p.dal.TotalCount(ctx, filter)
	if err != nil {
		p.logger.Errorf("[FindAllWithPagination] failed to count password rules: %v", err)
		return nil, err
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	p.logger.Infof("[FindAllWithPagination] retrieved %d password rules", len(data))

	return &types.PaginatedResponse[[]*model.PasswordRule]{
		Data: data,
		Meta: meta,
	}, nil
}

func (p *PasswordRuleStorage) FindCurrentRule(ctx context.Context) (*model.PasswordRule, error) {
	p.logger.Infof("[FindCurrentRule] fetching current password rule")
	ruleModel, err := p.dal.FindOne(ctx, bson.M{}, bson.M{})
	if err != nil || ruleModel == nil {
		p.logger.Errorf("[FindCurrentRule] failed to find current password rule: %v", err)
		return nil, err
	}
	p.logger.Infof("[FindCurrentRule] current password rule retrieved successfully")
	return ruleModel, nil
}
