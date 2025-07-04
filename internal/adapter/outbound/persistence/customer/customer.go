package customer

import (
	"context"
	"errors"
	"fmt"
	// "net/http"
	"regexp"
	"strings"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/customer/entity"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CustomerDetailRepo struct {
	client   *mongo.Client
	mongoDal dal.MongoDal[member.User, member.User]
	logger   utils.Logger
}

var _ outbound.CustomerDetailRepository = (*CustomerDetailRepo)(nil)

func normalizePhone(search string) bson.M {
	re := regexp.MustCompile(`^\+?251[79]\d{8}$`)

	trimmedSearch := strings.TrimSpace(search)

	if re.MatchString(trimmedSearch) {
		return bson.M{
			"phone_number": bson.M{
				"$regex":   trimmedSearch,
				"$options": "i",
			},
		}
	}

	return bson.M{
		"phone_number.number": bson.M{
			"$regex":   "^$",
			"$options": "i",
		},
	}
}

func InitCustomerDetail(client *mongo.Client, database string, collection string, logger utils.Logger) *CustomerDetailRepo {
	mongoDal := dal.NewMongoDal[member.User, member.User](client, database, collection)

	return &CustomerDetailRepo{
		client:   client,
		mongoDal: mongoDal,
		logger:   logger,
	}
}

func (c *CustomerDetailRepo) GetCustomersDetail(ctx context.Context, filterParams *constant.Filter) (*entity.CustomerRespose, error) {
	filter := bson.M{}
	projection := bson.M{}

	if filterParams.Search != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"full_name": bson.M{"$regex": filterParams.Search, "$options": "i"}},
				normalizePhone(filterParams.Search),
				{"user_code": bson.M{"$regex": filterParams.Search, "$options": "i"}},
				{"id": bson.M{"$regex": filterParams.Search, "$options": "i"}},
			},
		}
	}

	if filterParams.Filters != "" {
		filter["account_status"] = filterParams.Filters
	}

	skip := (filterParams.Page - 1) * filterParams.PerPage

	customers, err := c.mongoDal.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(filterParams.PerPage))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.logger.Errorf("no customers data found", err)

			return nil, fmt.Errorf("NO_CUSTOMER_DATA_FOUND")
		}
		c.logger.Errorf("failed to get customers data", err)
		return nil, fmt.Errorf("FAILED_TO_GET_CUSTOMER")
	}

	total, err := c.mongoDal.TotalCount(ctx, bson.M{})
	if err != nil {
		c.logger.Errorf("failed to get total counts", err)

		return nil, fmt.Errorf("FAILED_TO_GET_CUSTOMER_COUNT")
	}

	return &entity.CustomerRespose{
		Page:      1,
		Customers: customers,
		Limit:     constant.DefaultPerPage,
		Total:     total,
	}, nil
}

func (c *CustomerDetailRepo) GetCustomerByID(ctx context.Context, id string) (*member.User, error) {
	userID, err := bson.ObjectIDFromHex(id)
	if userID.IsZero() {
		c.logger.Errorf("invalid id provided", err)
		return nil, fmt.Errorf("INVALID_ID")
	}

	if err != nil {
		c.logger.Errorf("failed to convert id to object id", err)

		return nil, fmt.Errorf("FAILED_TO_CONVERT_ID")
	}

	filter := bson.M{"_id": userID}
	projection := bson.M{}
	member, err := c.mongoDal.FindOne(ctx, filter, projection)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.logger.Errorf("customer not found", err)
			return nil, fmt.Errorf("NO_CUSTOMER_DATA_FOUND")
		}
		c.logger.Errorf("failed to get customer", err)

		return nil, fmt.Errorf("FAILED_TO_GET_CUSTOMER")
	}
	return member, nil
}
