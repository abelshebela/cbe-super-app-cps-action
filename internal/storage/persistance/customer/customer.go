package customer

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/storage/kafka"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"cbe-super-app-cps-action/internal/storage"

	member "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	customer_dto "cbe-super-app-cps-action/internal/constants/dto/customer"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/types"
	local_util "cbe-super-app-cps-action/pkgs/utils"
)

type CustomerRepository struct {
	client           *mongo.Client
	mongoDal         dal.MongoDal[member.User, member.User]
	linkedAccountDal dal.MongoDal[model.LinkedAccount, model.LinkedAccount]
	logger           utils.Logger
	coll             *mongo.Collection
	kafkaProducer    kafka.ClientOrchestrationProducer
}

func InitCustomerDetail(client *mongo.Client, database string, collection []string, clientOrchestrationProducer kafka.ClientOrchestrationProducer, logger utils.Logger) storage.CustomerRepository {
	mongoDal := dal.NewMongoDal[member.User, member.User](client, database, collection[0])
	linkedAccountDal := dal.NewMongoDal[model.LinkedAccount, model.LinkedAccount](client, database, collection[1])
	return &CustomerRepository{
		client:           client,
		mongoDal:         mongoDal,
		logger:           logger,
		linkedAccountDal: linkedAccountDal,
		coll:             client.Database(database).Collection(collection[0]),
		kafkaProducer:    clientOrchestrationProducer,
	}
}

func (p *CustomerRepository) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*customer_dto.CustomerListResponse], error) {
	searchKeys := bson.M{}
	allowedKeys := []string{"gender", "branch_code", "kyc_level", "is_blocked", "enabled", "bps_reject_status"}

	if filterParam.Search != "" {
		// Build regex for Ethiopian phone numbers
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
			{"member_type": searchRegex},
			{"kyc_level": searchRegex},
		}
	}

	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		// Lookup Branch Info from account_block using branch_code
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "account_block"},
			{Key: "localField", Value: "branch_code"},
			{Key: "foreignField", Value: "code"},
			{Key: "as", Value: "branch_info"},
		}}},
		// Unwind branch_info (preserve if null)
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$branch_info"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		// Lookup District Info from account_block using branch_info.district_id
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "account_block"},
			{Key: "localField", Value: "branch_info.district_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "district_info"},
		}}},
		// Add status field based on is_blocked
		{{Key: "$addFields", Value: bson.D{
			{Key: "status", Value: bson.D{
				{Key: "$cond", Value: bson.D{
					{Key: "if", Value: bson.D{{Key: "$eq", Value: bson.A{"$is_blocked", false}}}},
					{Key: "then", Value: "Active"},
					{Key: "else", Value: "Inactive"},
				}},
			}},
		}}},
		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 1},
			{Key: "user_code", Value: 1},
			{Key: "full_name", Value: 1},
			{Key: "phone_number", Value: 1},
			{Key: "branch_code", Value: 1},
			{Key: "gender", Value: 1},
			{Key: "created_at", Value: 1},
			{Key: "is_blocked", Value: 1},
			{Key: "branch_name", Value: "$branch_info.name"},
			{Key: "district_name", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$district_info.name", 0}}}},
			{Key: "status", Value: 1},
		}}},
		{{Key: "$facet", Value: bson.D{
			{Key: "data", Value: bson.A{
				bson.D{{Key: "$skip", Value: skip}},
				bson.D{{Key: "$limit", Value: limit}},
			}},
			{Key: "total", Value: bson.A{
				bson.D{{Key: "$count", Value: "count"}},
			}},
		}}},
	}

	cursor, err := p.coll.Aggregate(ctx, pipeline)
	if err != nil {
		p.logger.Errorf("[FindAllWithPagination] aggregation failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer cursor.Close(ctx)

	var result []struct {
		Data []struct {
			ID           bson.ObjectID `bson:"_id"`
			UserCode     string        `bson:"user_code"`
			FullName     string        `bson:"full_name"`
			PhoneNumber  string        `bson:"phone_number"`
			BranchCode   string        `bson:"branch_code"`
			Gender       string        `bson:"gender"`
			CreatedAt    time.Time     `bson:"created_at"`
			IsBlocked    bool          `bson:"is_blocked"`
			BranchName   string        `bson:"branch_name"`
			DistrictName string        `bson:"district_name"`
			Status       string        `bson:"status"`
		} `bson:"data"`
		Total []struct {
			Count int64 `bson:"count"`
		} `bson:"total"`
	}

	if err := cursor.All(ctx, &result); err != nil {
		p.logger.Errorf("[FindAllWithPagination] cursor decode failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	var data []*customer_dto.CustomerListResponse
	var total int64

	if len(result) > 0 {
		data = make([]*customer_dto.CustomerListResponse, len(result[0].Data))
		for i, d := range result[0].Data {
			data[i] = &customer_dto.CustomerListResponse{
				ID:          d.ID.Hex(),
				UserCode:    d.UserCode,
				FullName:    d.FullName,
				PhoneNumber: d.PhoneNumber,
				BranchCode:  d.BranchCode,
				Gender:      d.Gender,
				CreatedAt:   d.CreatedAt.Format(time.RFC3339),
				IsBlocked:   d.IsBlocked,
			}
		}
		if len(result[0].Total) > 0 {
			total = result[0].Total[0].Count
		}
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	p.logger.Infof("[FindAllWithPagination] retrieved %d customers", len(data))

	return &types.PaginatedResponse[[]*customer_dto.CustomerListResponse]{
		Data: data,
		Meta: meta,
	}, nil
}

func (p *CustomerRepository) Update(ctx context.Context, id string, data member.User) error {
	p.logger.Infof("[Update] updating customer for id: %s", id)
	objId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		p.logger.Errorf("[Update] invalid object id: %v", err)
		return localization.ErrorUnexpectedError
	}

	filter, update := FaydaEnable(objId, data)
	updatedCustomer, err := p.mongoDal.UpdateOne(ctx, filter, update)
	if err != nil {
		p.logger.Errorf("[Update] failed to update customer: %v", err)
		return err
	}
	p.kafkaProducer.PublishMessage(ctx, updatedCustomer, string(constants.ClientOrchestrationMemberTopic), string(constants.ClientOrchestrationMemberTopic), "Customer Update")

	p.logger.Infof("[Update] customer updated successfully")
	return nil
}
func (p *CustomerRepository) FindByID(ctx context.Context, id string) (*member.User, error) {
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
	updatedCustomer, err := b.mongoDal.UpdateOne(ctx, filter, update)
	if err != nil {
		b.logger.Errorf("[EnableOrDisable] failed to enable/disable customer: %v", err)
		return err
	}

	b.kafkaProducer.PublishMessage(ctx, updatedCustomer, string(constants.ClientOrchestrationMemberTopic), string(constants.ClientOrchestrationMemberTopic), "enable/disable Customer")

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

func (p *CustomerRepository) FindCustomerDetailByID(ctx context.Context, id string) (*customer_dto.CustomerDetailResponse, error) {
	p.logger.Infof("[FindCustomerDetailByID] fetching customer detail by id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		p.logger.Errorf("[FindCustomerDetailByID] invalid object id: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "_id", Value: objID}}}},
		// Lookup Linked Accounts using user_id (ObjectId)
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "linked_account"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "user_id"},
			{Key: "as", Value: "linked_accounts"},
		}}},
		// Lookup KYC Data
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "customer_kyc"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "user_id"},
			{Key: "as", Value: "kyc_data"},
		}}},
		// Lookup Branch Info from account_block using branch_code
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "account_block"},
			{Key: "localField", Value: "branch_code"},
			{Key: "foreignField", Value: "code"},
			{Key: "as", Value: "branch_info"},
		}}},
		// Unwind branch_info to use its district_id for next lookup (preserve if null)
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$branch_info"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		// Lookup District Info from account_block using branch_info.district_id
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "account_block"},
			{Key: "localField", Value: "branch_info.district_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "district_info"},
		}}},
		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 1},
			{Key: "user_code", Value: 1},
			{Key: "full_name", Value: 1},
			{Key: "phone_number", Value: 1},
			{Key: "gender", Value: 1},
			{Key: "account_numbers", Value: "$linked_accounts.account_number"},
			{Key: "branch_name", Value: "$branch_info.name"},
			{Key: "district_name", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$district_info.name", 0}}}},
			{Key: "kyc_data", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$kyc_data", 0}}}},
		}}},
	}

	cursor, err := p.coll.Aggregate(ctx, pipeline)
	if err != nil {
		p.logger.Errorf("[FindCustomerDetailByID] aggregation failed: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		ID             bson.ObjectID `bson:"_id"`
		CustomerCode   string        `bson:"user_code"`
		FullName       string        `bson:"full_name"`
		PhoneNumber    string        `bson:"phone_number"`
		Gender         string        `bson:"gender"`
		AccountNumbers []string      `bson:"account_numbers"`
		BranchName     string        `bson:"branch_name"`
		DistrictName   string        `bson:"district_name"`
		KYCData        struct {
			MothersName   string               `bson:"mothers_name"`
			Nationality   string               `bson:"nationality"`
			BirthDate     string               `bson:"birth_date"`
			Address       customer_dto.Address `bson:"address"`
			MonthlyIncome string               `bson:"monthly_income"`
		} `bson:"kyc_data"`
	}

	if err = cursor.All(ctx, &results); err != nil {
		p.logger.Errorf("[FindCustomerDetailByID] cursor decode failed: %v", err)
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("customer not found")
	}

	res := results[0]

	// Split Full Name
	parts := strings.Fields(res.FullName)
	firstName := ""
	middleName := ""
	lastName := ""
	if len(parts) > 0 {
		firstName = parts[0]
	}
	if len(parts) > 1 {
		middleName = parts[1]
	}
	if len(parts) > 2 {
		lastName = strings.Join(parts[2:], " ")
	}

	response := &customer_dto.CustomerDetailResponse{
		ID:             res.ID.Hex(),
		CustomerCode:   res.CustomerCode,
		FirstName:      firstName,
		MiddleName:     middleName,
		LastName:       lastName,
		FullName:       res.FullName,
		PhoneNumber:    res.PhoneNumber,
		Gender:         res.Gender,
		AccountNumbers: res.AccountNumbers,
		BranchName:     res.BranchName,
		DistrictName:   res.DistrictName,
		MothersName:    res.KYCData.MothersName,
		Nationality:    res.KYCData.Nationality,
		BirthDate:      res.KYCData.BirthDate,
		Address:        res.KYCData.Address,
		MonthlyIncome:  res.KYCData.MonthlyIncome,
	}

	return response, nil
}

// func (p *CustomerRepository) SearchCustomerByCIForAccountNumber(ctx context.Context, number string) (*customer_dto.CustomerListResponse, error) {
// 	p.logger.Infof("[SearchCustomerByCIForAccountNumber] searching customer by value: %s", number)

// 	// Step 1: Find user_id from linked_account where customer_id and account_number match
// 	linkedAccountColl := p.client.Database(p.coll.Database().Name()).Collection("linked_account")
// 	var linkedResult struct {
// 		UserID interface{} `bson:"user_id"`
// 	}
// 	f := bson.M{"$or": bson.A{
// 		bson.M{"customer_number": number},
// 		bson.M{"account_number": number},
// 	}}

// 	err := linkedAccountColl.FindOne(ctx, f).Decode(&linkedResult)
// 	if err != nil {
// 		code, _ := local_util.HandleMongoError(err)
// 		if code == localization.ErrorResourceNotFound.Code {
// 			p.logger.Errorf("[searchCustomerByCIForAccountNumber] customer not found")
// 			return nil, fmt.Errorf("%s", code)
// 		}
// 		p.logger.Errorf("[searchCustomerByCIForAccountNumber] failed to fetch customer: %v", err)
// 		return nil, err
// 	}

// 	res, err := p.mongoDal.FindOne(ctx, bson.M{"_id": linkedResult.UserID}, nil)
// 	if err != nil {
// 		p.logger.Errorf("[searchCustomerByCIForAccountNumber] failed to find customer: %v", err)
// 		return nil, err
// 	}

// 	response := &customer_dto.CustomerListResponse{
// 		ID:          res.ID.Hex(),
// 		UserCode:    res.UserCode,
// 		FullName:    res.FullName,
// 		PhoneNumber: res.PhoneNumber,
// 		BranchCode:  res.BranchCode,
// 		Gender:      string(res.Gender),
// 		CreatedAt:   res.CreatedAt.Format(time.RFC3339),
// 		IsBlocked:   res.IsBlocked,
// 	}

//		return response, nil
//	}
func (p *CustomerRepository) SearchCustomerByCIForAccountNumber(ctx context.Context, number string) (*customer_dto.CustomerListResponse, error) {
	p.logger.Infof("[SearchCustomerByCIForAccountNumber] searching members by value: %s", number)

	pipeline := mongo.Pipeline{
		// 1. JOIN Linked Accounts
		// We start with the 'members' collection and look into 'linked_account'
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "linked_account"},
			{Key: "localField", Value: "_id"},       // member's ObjectID
			{Key: "foreignField", Value: "user_id"}, // linked_account's owner ObjectID
			{Key: "as", Value: "accounts"},
		}}},

		// 2. MULTI-FIELD MATCH
		// This stage checks the current member document AND the joined 'accounts' array simultaneously.
		{{Key: "$match", Value: bson.D{
			{Key: "$or", Value: bson.A{
				bson.M{"phone_number": number},             // Search in 'members'
				bson.M{"accounts.customer_number": number}, // Search in joined 'linked_account' array
				bson.M{"accounts.account_number": number},  // Search in joined 'linked_account' array
			}},
		}}},

		// 3. LIMIT
		// We only need the first member that matches any of the criteria.
		{{Key: "$limit", Value: 1}},

		// 4. PROJECT
		// Formatting the output and ensuring 'user_id' is returned as a string.
		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 1},
			{Key: "user_id", Value: bson.M{"$toString": "$_id"}},
			{Key: "user_code", Value: 1},
			{Key: "full_name", Value: 1},
			{Key: "phone_number", Value: 1},
			{Key: "branch_code", Value: 1},
			{Key: "gender", Value: 1},
			{Key: "created_at", Value: 1},
			{Key: "is_blocked", Value: 1},
			// Extract the first account_number from linked_account array if exists
			{Key: "account_number", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$accounts.account_number", 0}}}},
		}}},
	}

	cursor, err := p.coll.Aggregate(ctx, pipeline) // p.coll must point to 'members'
	if err != nil {
		p.logger.Errorf("[SearchCustomerByCIForAccountNumber] aggregation failed: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		ID            bson.ObjectID `bson:"_id"`
		UserID        string        `bson:"user_id"`
		UserCode      string        `bson:"user_code"`
		AccountNumber string        `bson:"account_number"`
		FullName      string        `bson:"full_name"`
		PhoneNumber   string        `bson:"phone_number"`
		BranchCode    string        `bson:"branch_code"`
		Gender        string        `bson:"gender"`
		CreatedAt     time.Time     `bson:"created_at"`
		IsBlocked     bool          `bson:"is_blocked"`
	}

	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no customer found matching: %s", number)
	}

	res := results[0]
	return &customer_dto.CustomerListResponse{
		ID:            res.ID.Hex(),
		UserID:        res.UserID,
		UserCode:      res.UserCode,
		FullName:      res.FullName,
		PhoneNumber:   res.PhoneNumber,
		BranchCode:    res.BranchCode,
		Gender:        res.Gender,
		CreatedAt:     res.CreatedAt.Format(time.RFC3339),
		IsBlocked:     res.IsBlocked,
		AccountNumber: res.AccountNumber,
	}, nil
}

func (p *CustomerRepository) FindCustomerByIDs(ctx context.Context, ids []string) ([]*member.User, error) {
	p.logger.Infof("[FindCustomerByIDs] fetching customers by ids")
	var objIDs []bson.ObjectID
	for _, idStr := range ids {
		objID, err := bson.ObjectIDFromHex(idStr)
		if err != nil {
			p.logger.Errorf("[FindCustomerByIDs] invalid ObjectID: %s", idStr)
			return nil, errors.New(localization.ErrorInvalidID.Code)
		}
		objIDs = append(objIDs, objID)
	}
	filter := bson.M{"_id": bson.M{"$in": objIDs}}

	customers, err := p.mongoDal.FindAll(ctx, filter, nil)
	if err != nil {
		p.logger.Errorf("[FindCustomerByIDs] failed to find customers by ids: %v", err)
		return nil, err
	}
	return customers, nil
}
