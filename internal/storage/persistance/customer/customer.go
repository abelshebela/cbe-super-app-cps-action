package customer

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"context"
	"errors"
	"fmt"
	"strings"

	"cbe-super-app-cps-action/internal/storage"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/types"
	local_util "cbe-super-app-cps-action/pkgs/utils"
)

type CustomerRepository struct {
	client           *mongo.Client
	mongoDal         dal.MongoDal[model.User, model.User]
	linkedAccountDal dal.MongoDal[model.LinkedAccount, model.LinkedAccount]
	logger           utils.Logger
}

func InitCustomerDetail(client *mongo.Client, database string, collection []string, logger utils.Logger) storage.CustomerRepository {
	mongoDal := dal.NewMongoDal[model.User, model.User](client, database, collection[0])
	linkedAccountDal := dal.NewMongoDal[model.LinkedAccount, model.LinkedAccount](client, database, collection[1])
	return &CustomerRepository{
		client:           client,
		mongoDal:         mongoDal,
		logger:           logger,
		linkedAccountDal: linkedAccountDal,
	}
}

func (p *CustomerRepository) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.User], error) {

	searchKeys := bson.M{}

	allowedKeys := []string{"gender", "branch_code", "kyc_level", "is_blocked", "enabled", "bps_reject_status"}

	if filterParam.Search != "" {
		// Build regex for Ethiopian phone numbers: match any prefix (+2519, +2517, +25109, +25107, 09, 07, 9, 7) and main number
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		mainNumber := filterParam.Search
		mainNumber = strings.TrimPrefix(mainNumber, " ")
		mainNumber = strings.TrimPrefix(mainNumber, "+")
		mainNumber = strings.TrimPrefix(mainNumber, "251")
		mainNumber = strings.TrimPrefix(mainNumber, "0")
		phonePattern := fmt.Sprintf("(?:\\+?2510?|0)?%s$", mainNumber)
		phoneRegex := bson.M{"$regex": phonePattern, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"full_name": searchRegex},
			{"phone_number": phoneRegex},
			{"gender": searchRegex},
			{"user_name": searchRegex},
			{"user_code": searchRegex},
			{"is_blocked": searchRegex},
			{"kyc_level": searchRegex},
		}
	}

	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	data, err := p.mongoDal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		p.logger.Errorf("[FindAllWithPagination] failed to fetch customers: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := p.mongoDal.TotalCount(ctx, filter)
	if err != nil {
		p.logger.Errorf("[FindAllWithPagination] failed to count customers: %v", err)
		return nil, errors.New(localization.ErrorNoDataProvided.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	p.logger.Infof("[FindAllWithPagination] retrieved %d customers", len(data))

	return &types.PaginatedResponse[[]*model.User]{
		Data: data,
		Meta: meta,
	}, nil
}

func (p *CustomerRepository) Update(ctx context.Context, id string, data model.User) error {
	p.logger.Infof("[Update] updating customer for id: %s", id)
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		p.logger.Errorf("[Update] invalid object id: %v", err)
		return localization.ErrorUnexpectedError
	}

	filter, update := FaydaEnable(objId, data)
	_, err = p.mongoDal.UpdateOne(ctx, filter, update)
	if err != nil {
		p.logger.Errorf("[Update] failed to update customer: %v", err)
		return err
	}
	p.logger.Infof("[Update] customer updated successfully")
	return nil
}
func (p *CustomerRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	p.logger.Infof("[FindByID] fetching customer by id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		p.logger.Errorf("[FindByID] invalid object id: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID}
	// user, err := p.mongoDal.FindOne(ctx, filter, UserProjection())
	user, err := p.mongoDal.FindOne(ctx, filter, nil)
	if err != nil {
		code, _ := local_util.HandleMongoError(err)
		if code == localization.ErrorResourceNotFound.Code {
			p.logger.Errorf("[FindByID] customer not found")
			return nil, fmt.Errorf("%s", code)
		}
		p.logger.Errorf("[FindByID] failed to fetch customer: %v", err)
		return nil, err
	}

	if user == nil {
		p.logger.Errorf("[FindByID] customer not found for id: %s", id)
		return nil, fmt.Errorf("no user found for given id")
	}
	p.logger.Infof("[FindByID] customer retrieved successfully")
	return user, nil
}

// EnableOrDisable implements storage.CustomerRepository.
func (b *CustomerRepository) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	b.logger.Infof("[EnableOrDisable] processing customer enable/disable for id: %s, enabled: %v", id, enable)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		b.logger.Errorf("[EnableOrDisable] invalid object id: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID}
	update := bson.M{"enabled": enable}
	_, err = b.mongoDal.UpdateOne(ctx, filter, update)
	if err != nil {
		b.logger.Errorf("[EnableOrDisable] failed to enable/disable customer: %v", err)
		return err
	}
	b.logger.Infof("[EnableOrDisable] customer enable/disable completed successfully")
	return nil
}

func (c *CustomerRepository) FetchLinkedAccount(ctx context.Context, customerNumber string) ([]*model.LinkedAccount, error) {
	c.logger.Infof("[FetchLinkedAccount] fetching linked accounts for customer number")
	filter := bson.M{
		"customer_number": customerNumber,
	}

	linkedAccount, err := c.linkedAccountDal.FindAll(ctx, filter, bson.M{})
	if err != nil {
		c.logger.Errorf("[FetchLinkedAccount] failed to fetch linked accounts: %v", err)
		return []*model.LinkedAccount{}, err
	}
	c.logger.Infof("[FetchLinkedAccount] retrieved %d linked accounts", len(linkedAccount))

	return linkedAccount, nil
}
