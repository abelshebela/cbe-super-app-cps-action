package customer

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"context"
	"errors"
	"fmt"

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
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"full_name": searchRegex},
			{"phone_number": searchRegex},
			{"gender": searchRegex},
			{"user_name": searchRegex},
			{"user_code": searchRegex},
		}
	}

	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	if filter["is_blocked"] == nil {
		filter["is_blocked"] = false
	}
	filter["enabled"] = true

	data, err := p.mongoDal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		p.logger.Infof("error while fetching customer data: %v", err.Error())
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := p.mongoDal.TotalCount(ctx, filter)
	if err != nil {
		return nil, errors.New(localization.ErrorNoDataProvided.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]*model.User]{
		Data: data,
		Meta: meta,
	}, nil
}

func (p *CustomerRepository) Update(ctx context.Context, id string, data model.User) error {

	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		p.logger.Errorf("Error while parsing id from string to object")
		return localization.ErrorUnexpectedError
	}

	filter, update := FaydaEnable(objId, data)
	_, err = p.mongoDal.UpdateOne(ctx, filter, update)
	if err != nil {
		p.logger.Errorf("Error while updating the customer data")
		return err
	}

	return nil
}
func (p *CustomerRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID}
	p.logger.Infof("Fetching user with filter: %v", filter)
	// user, err := p.mongoDal.FindOne(ctx, filter, UserProjection())
	user, err := p.mongoDal.FindOne(ctx, filter, nil)
	if err != nil {
		code, _ := local_util.HandleMongoError(err)
		if code == localization.ErrorResourceNotFound.Code {
			return nil, fmt.Errorf("%s", code)
		} else if err != nil {
			return nil, err
		}

		p.logger.Errorf("Failed to fetch user by ID: %v", err)
		return nil, err
	}

	if user == nil {
		p.logger.Infof("no user found for given id: %v", id)
		return nil, fmt.Errorf("no user found for given id")
	}
	return user, nil
}

// EnableOrDisable implements storage.CustomerRepository.
func (b *CustomerRepository) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	b.logger.Infof("Enabling/Disabling customer with ID: %s to %v", id, enable)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		b.logger.Errorf("invalid object id: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID}
	update := bson.M{"enabled": enable}
	_, err = b.mongoDal.UpdateOne(ctx, filter, update)
	if err != nil {
		b.logger.Errorf("error while enabling/disabling customer: %v", err)
		return err
	}
	return nil
}

func (c *CustomerRepository) FetchLinkedAccount(ctx context.Context, customerNumber string) ([]*model.LinkedAccount, error) {

	filter := bson.M{
		"customer_number": customerNumber,
	}

	linkedAccount, err := c.linkedAccountDal.FindAll(ctx, filter, bson.M{})
	if err != nil {
		return []*model.LinkedAccount{}, err
	}

	return linkedAccount, nil
}
