package bank

import (
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
	"regexp"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
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

// FindBIC implements [storage.BankRepository].
func (b *BankStorage) FindByBIC(ctx context.Context, bic string) (*model.Bank, error) {
	b.logger.Infof("[FindBIC] fetching bank by BIC: %s", bic)
	bank, err := b.dal.FindOne(ctx, bson.M{"bic_code": bic, "enabled": true, "is_deleted": false}, nil)
	if err != nil {
		b.logger.Errorf("[FindBIC] failed to fetch bank: %v", err)
		return nil, err
	}
	return bank, nil
}

func NewBankRepository(client *mongo.Client, cfg *config.VaultConfig, dbName string, collection string, logger utils.Logger) storage.BankRepository {
	return &BankStorage{
		dal:    dal.NewMongoDal[model.Bank, model.Bank](client, cfg, dbName, collection),
		client: client,
		logger: logger,
	}
}

func (b *BankStorage) Create(ctx context.Context, bank *model.Bank) error {
	b.logger.Infof("[Create] creating bank")
	bank.ID = bson.NewObjectID()
	_, err := b.dal.InsertOne(ctx, *bank)
	if err != nil {
		b.logger.Errorf("[Create] failed to create bank: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	b.logger.Infof("[Create] bank created successfully")
	return nil
}

func (b *BankStorage) Update(ctx context.Context, id string, bank *model.Bank) error {
	b.logger.Infof("[Update] updating bank for id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		b.logger.Errorf("[Update] invalid object id: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	updateData := BankMapper(*bank)

	_, err = b.dal.UpdateOne(ctx, filter, updateData)
	if err != nil {
		b.logger.Errorf("[Update] failed to update bank: %v", err)
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	b.logger.Infof("[Update] bank updated successfully")
	return nil
}

func (b *BankStorage) Delete(ctx context.Context, id string) error {
	b.logger.Infof("[Delete] deleting bank for id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		b.logger.Errorf("[Delete] invalid object id: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	err = b.dal.DeleteOne(ctx, filter)
	if err != nil {
		b.logger.Errorf("[Delete] failed to delete bank: %v", err)
		return err
	}
	b.logger.Infof("[Delete] bank deleted successfully")
	return nil
}

func (b *BankStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	b.logger.Infof("[EnableOrDisable] processing bank enable/disable for id: %s, enabled: %v", id, enable)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		b.logger.Errorf("[EnableOrDisable] invalid object id: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID}
	update := bson.M{"enabled": enable}
	_, err = b.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			b.logger.Errorf("[EnableOrDisable] bank not found")
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		b.logger.Errorf("[EnableOrDisable] failed to enable/disable bank: %v", err)
		return err
	}
	b.logger.Infof("[EnableOrDisable] bank enable/disable completed successfully")
	return nil
}

func (b *BankStorage) FindByID(ctx context.Context, id string) (*model.Bank, error) {
	b.logger.Infof("[FindByID] fetching bank by id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		b.logger.Errorf("[FindByID] invalid object id: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID}

	result, err := b.dal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			b.logger.Errorf("[FindByID] bank not found")
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		b.logger.Errorf("[FindByID] failed to find bank: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	if result.IsDeleted {
		b.logger.Errorf("[FindByID] bank is deleted")
		return nil, mongo.ErrNoDocuments
	}
	b.logger.Infof("[FindByID] bank retrieved successfully")
	return result, nil
}

func (s *BankStorage) FindByNameOrBIC(
	ctx context.Context, bic, name string) (*model.Bank, error) {

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

	filter := bson.M{"$or": conditions}
	bank, err := s.dal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.logger.Infof("[FindByNameOrBICOrCode] no bank found matching the criteria")
			return nil, nil
		}
		s.logger.Errorf("[FindByNameOrBICOrCode] failed to find bank: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return bank, nil
}

func (s *BankStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.Bank], error) {
	filter := bson.M{"is_deleted": false}
	searchKeys := bson.M{}
	allowedKeys := []string{"search", "branch_code", "branch_name", "enabled", "enabled", "is_deleted"}

	if filterParam.Search != "" && filterParam.Search != "enabled" {
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
	if filterParam.Search == "enabled" {
		filter["enabled"] = true
	}
	data, err := s.dal.FindAllWithPaginationE(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			s.logger.Infof("[FindAllWithPagination] no banks found")
			return nil, localization.ErrorResourceNotFound
		}
		s.logger.Errorf("[FindAllWithPagination] failed to fetch banks: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := s.dal.TotalCount(ctx, filter)
	if err != nil {
		s.logger.Errorf("[FindAllWithPagination] failed to count banks: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	s.logger.Infof("[FindAllWithPagination] retrieved %d banks", len(data))

	return &types.PaginatedResponse[[]model.Bank]{
		Data: data,
		Meta: meta,
	}, nil
}
