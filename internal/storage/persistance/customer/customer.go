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
	client   *mongo.Client
	mongoDal dal.MongoDal[model.User, model.User]
	logger   utils.Logger
}

func InitCustomerDetail(client *mongo.Client, database string, collection string, logger utils.Logger) storage.CustomerRepository {
	mongoDal := dal.NewMongoDal[model.User, model.User](client, database, collection)

	return &CustomerRepository{
		client:   client,
		mongoDal: mongoDal,
		logger:   logger,
	}
}

func (p *CustomerRepository) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.User], error) {
	// 1. Base filter (only active records)
	// filter := bson.M{"is_deleted": false}
	// Add additional filters from filterParam.Filters if provided
	// for k, v := range filterParam.Filters {
	// 	filter[k] = v
	// }
	searchKeys := bson.M{}

	// 2. Allowed filterable/searchable fields
	allowedKeys := []string{}

	// 3. Add search (if provided)
	// if filterParam.Search != "" {
	// 	searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
	// 	searchKeys["field1"] = searchRegex // choose your searchable field(s)
	// }

	// 4. Build filter, skip, limit
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	// 5. Fetch data
	data, err := p.mongoDal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)

	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	fmt.Printf("data: %v", data)

	// 6. Count total
	total, err := p.mongoDal.TotalCount(ctx, filter)
	if err != nil {
		return nil, errors.New(localization.ErrorNoDataProvided.Code)
	}

	// 7. Build pagination metadata
	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	// 8. Return standard paginated response
	return &types.PaginatedResponse[[]*model.User]{
		Data: data,
		Meta: meta,
	}, nil
}

func (p *CustomerRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	fmt.Println("persistance 1")
	objID, err := bson.ObjectIDFromHex(id)
	filter := map[string]interface{}{"_id": objID}
	user, err := p.mongoDal.FindOne(ctx, filter, nil)
	fmt.Println("persistance 2")

	if err != nil {
		p.logger.Errorf("Failed to fetch user by ID: %v", err)
		return nil, err
	}
	if user != nil {
		p.logger.Infof("no user found for given id: %v", id)
		return nil, fmt.Errorf("no user found for given id")
	}
	return user, nil
}
