package customer

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"regexp"
	"strings"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	outbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/customer"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	member "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
)

type CustomerDetailRepo struct {
	client   *mongo.Client
	mongoDal dal.MongoDal[member.User, member.User]
	logger   utils.Logger
}

func normalizePhone(search string) bson.M {
	re := regexp.MustCompile(`^\+?251[79]\d{8}$`)
	trimmed := strings.TrimSpace(search)

	if re.MatchString(trimmed) {
		return bson.M{"phone_number": bson.M{"$regex": trimmed, "$options": "i"}}
	}

	// no match => skip phone condition
	return nil
}

func InitCustomerDetail(client *mongo.Client, database string, collection string, logger utils.Logger) outbound.CustomerRepository {
	mongoDal := dal.NewMongoDal[member.User, member.User](client, database, collection)

	return &CustomerDetailRepo{
		client:   client,
		mongoDal: mongoDal,
		logger:   logger,
	}
}

func (c *CustomerDetailRepo) GetCustomersDetail(ctx context.Context, kycLevel int, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*member.User], error) {
	filter := buildCustomerFilter(filterParams)
	filter["is_deleted"] = false
	filter["kyc_level"] = kycLevel

	return c.paginateFind(ctx, filter, filterParams)
}

func (c *CustomerDetailRepo) GetBlockedCustomer(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*member.User], error) {
	filter := buildCustomerFilter(filterParams)
	filter["is_deleted"] = false
	filter["is_blocked"] = true

	if filterParams.Filters != "" {
		filter["account_status"] = filterParams.Filters
	}

	return c.paginateFind(ctx, filter, filterParams)
}

func (c *CustomerDetailRepo) paginateFind(ctx context.Context, filter bson.M, params *constant.Filter) (*common_util.PaginatedResponse[[]*member.User], error) {
	skip := int64((params.Page - 1) * params.PerPage)
	limit := int64(params.PerPage)

	users, err := c.mongoDal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		return &common_util.PaginatedResponse[[]*member.User]{}, nil
	}

	total, err := c.mongoDal.TotalCount(ctx, filter)
	if err != nil {
		c.logger.Errorf("failed to count customers: %v", err)
		return nil, fmt.Errorf("FAILED_TO_GET_CUSTOMER_COUNT")
	}

	meta := common_util.BuildPaginationMeta(total, params.Page, params.PerPage)

	return &common_util.PaginatedResponse[[]*member.User]{
		Data: users,
		Meta: meta,
	}, nil
}

func buildCustomerFilter(params *constant.Filter) bson.M {
	filter := bson.M{}

	if params.Search != "" {
		search := params.Search
		filter["$or"] = []bson.M{
			{"full_name": bson.M{"$regex": search, "$options": "i"}},
			normalizePhone(search),
			{"user_code": bson.M{"$regex": search, "$options": "i"}},
			{"id": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	if params.Filters != "" {
		if blocked, err := strconv.ParseBool(params.Filters); err == nil {
			filter["is_blocked"] = blocked
		}
	}

	return filter
}

func (c *CustomerDetailRepo) GetCustomerByID(ctx context.Context, id string) (*member.User, error) {
	userID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		c.logger.Errorf("failed to convert id to object id", err)
		return nil, fmt.Errorf("FAILED_TO_CONVERT_ID")
	}
	if userID.IsZero() {
		c.logger.Errorf("invalid id provided", err)
		return nil, fmt.Errorf("INVALID_ID")
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
