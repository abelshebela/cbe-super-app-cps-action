package customer

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	// "net/http"
	"regexp"
	"strings"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/customer/entity"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"

	// "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	member "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
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

func (c *CustomerDetailRepo) GetCustomersDetail(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*member.User], error) {
	filter := bson.M{
		"is_deleted": false,
	}
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

	// if filterParams.Filters != "" {
	// 	filter["account_status"] = filterParams.Filters
	// }
	if filterParams.Filters != "" {
		blocked, err := strconv.ParseBool(filterParams.Filters)
		if err == nil {
			filter["is_blocked"] = blocked
		}
	}
	skip := (filterParams.Page - 1) * filterParams.PerPage
	limit := filterParams.PerPage

	customers, err := c.mongoDal.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(limit))
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

	meta := common_util.BuildPaginationMeta(total, limit, filterParams.Page)
	fmt.Println("============================================")
	// customerList := make([]member.User, len(customers))

	for index, user := range customers {
		data := *user
		fmt.Printf("index:%v data: %v", index, data)
	}
	// fmt.Printf("list of users : %v", *customers)
	fmt.Println("============================================")
	return &common_util.PaginatedResponse[[]*member.User]{
		Data: customers,
		Meta: meta,
	}, nil

}

func (c *CustomerDetailRepo) GetBlockedCustomer(ctx context.Context, filterParams *constant.Filter) (*entity.CustomerRespose, error) {
	filter := bson.M{
		"is_deleted": false,
		"is_blocked": true,
	}
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

	fmt.Println("filter-------------", filter)
	customers, err := c.mongoDal.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(filterParams.PerPage))
	fmt.Println("filtered==========", err)
	fmt.Println("filtered==========", customers)
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

func (c *CustomerDetailRepo) GetFaydaCustomersDeatil(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*member.User], error) {
	filter := bson.M{
		"is_deleted": false,
	}
	filter["kyc.level"] = int32(1)

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
		filter["is_blocked"] = filterParams.Filters
	}

	// skip := (filterParams.Page - 1) * filterParams.PerPage

	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	fmt.Printf("filtered by kyc 1: %v", filter)
	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")

	skip := (filterParams.Page - 1) * filterParams.PerPage
	limit := filterParams.PerPage
	customers, err := c.mongoDal.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(filterParams.PerPage))
	fmt.Println("---------------------------------------------------")
	fmt.Printf("fayda customer: %v ", customers)
	fmt.Println("----------------------------------------------------")
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

	meta := common_util.BuildPaginationMeta(total, limit, filterParams.Page)

	return &common_util.PaginatedResponse[[]*member.User]{
		Data: customers,
		Meta: meta,
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

func (c *CustomerDetailRepo) GetFaydaCustomerByID(ctx context.Context, id string) (*member.User, error) {
	userID, err := bson.ObjectIDFromHex(id)
	if userID.IsZero() {
		c.logger.Errorf("invalid id provided", err)
		return nil, fmt.Errorf("INVALID_ID")
	}

	if err != nil {
		c.logger.Errorf("failed to convert id to object id", err)

		return nil, fmt.Errorf("FAILED_TO_CONVERT_ID")
	}

	filter := bson.M{"_id": userID, "kyc.level": 1}
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
