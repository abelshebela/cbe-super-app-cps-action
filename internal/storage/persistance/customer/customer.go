package customer

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"context"
	"fmt"

	"cbe-super-app-cps-action/internal/storage"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"cbe-super-app-cps-action/internal/constants/types"
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
	// Build filter from filterParam
	fmt.Println("perstistance 1")

	filter := make(map[string]interface{})
	// if filterParam.Search != "" {
	// 	filter["$text"] = map[string]interface{}{"$search": filterParam.Search}
	// }
	for k, v := range filterParam.Filters {
		filter[k] = v
	}

	// Pagination options
	page := filterParam.Page
	if page < 1 {
		page = 1
	}
	limit := filterParam.PerPage
	if limit < 1 {
		limit = 10
	}
	skip := (page - 1) * limit

	// Projection (if needed, otherwise pass nil or default)
	var projection map[string]interface{}
	// If you have projection in Filter, handle here

	// Convert skip and limit to int64
	skip64 := int64(skip)
	limit64 := int64(limit)

	// Query database
	users, err := p.mongoDal.FindAllWithPagination(ctx, filter, projection, skip64, limit64)
	if err != nil {
		p.logger.Errorf("Failed to fetch users: ", err)
		return nil, err
	}
	fmt.Println("perstistance 2")

	// Prepare paginated response
	total := int64(len(users)) // You should ideally get the total count from DB, not just the current page length

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	pagingCounter := skip + 1
	hasPrevPage := page > 1
	hasNextPage := page < totalPages

	var prevPage *int
	var nextPage *int
	if hasPrevPage {
		p := page - 1
		prevPage = &p
	}
	if hasNextPage {
		n := page + 1
		nextPage = &n
	}

	resp := &types.PaginatedResponse[[]*model.User]{
		Data: users,
		Meta: types.PaginationMeta{
			TotalDocs:     total,
			Limit:         limit,
			TotalPages:    totalPages,
			Page:          page,
			PagingCounter: pagingCounter,
			HasPrevPage:   hasPrevPage,
			HasNextPage:   hasNextPage,
			PrevPage:      prevPage,
			NextPage:      nextPage,
		},
	}

	return resp, nil
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
