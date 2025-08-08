package unlink

import (
	"context"
	"fmt"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/common"
	contexts "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	local_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type unlinkCustomer struct {
	logger            utils.Logger
	userDal           dal.MongoDal[member.User, member.User]
	cpsDal            dal.MongoDal[model.CPSAction, model.CPSAction]
	archivedUserDal   dal.MongoDal[model.ArchivedUser, model.ArchivedUser]
	linkedDal         dal.MongoDal[model.LinkedAccount, model.LinkedAccount]
	archivedLinkedDal dal.MongoDal[model.ArchivedLinkedAccount, model.ArchivedLinkedAccount]
	client            *mongo.Client
	dbName            string
	collection        []string
}

type UnlinkAccount interface {
	GetUserByAccount(ctx context.Context, accNumber string) (*any, error)
	UnlinkUserCif(ctx context.Context, userCode string) error
	Authorize(ctx context.Context, cpsAction any) (any, error)
}

func NewUnlinkPersistence(client *mongo.Client, database string, collection []string, logger utils.Logger) UnlinkAccount {

	return &unlinkCustomer{
		logger:            logger,
		userDal:           dal.NewMongoDal[member.User, member.User](client, database, collection[1]),
		cpsDal:            dal.NewMongoDal[model.CPSAction, model.CPSAction](client, database, collection[0]),
		linkedDal:         dal.NewMongoDal[model.LinkedAccount, model.LinkedAccount](client, database, collection[4]),
		archivedUserDal:   dal.NewMongoDal[model.ArchivedUser, model.ArchivedUser](client, database, collection[2]),
		archivedLinkedDal: dal.NewMongoDal[model.ArchivedLinkedAccount, model.ArchivedLinkedAccount](client, database, collection[3]),
		client:            client,
		dbName:            database,
		collection:        collection,
	}
}

func (u *unlinkCustomer) GetUserByAccount(ctx context.Context, accNumber string) (*any, error) {
	if accNumber == "" {
		return nil, fmt.Errorf("ACCOUNT_NUMBER_CANNOT_BE_EMPTY")
	}

	// Step 1: Fetch linked account(s)
	filter := bson.M{
		"account_number": accNumber,
	}
	projection := bson.M{
		"account_number":      1,
		"linked_status":       1,
		"account_holder_name": 1,
		"customer_number":     1,
		"linked_date":         1,
		"account_type":        1,
		"is_account_active":   1,
		"account_branch_code": 1,
		"linked_branch":       1,
		"_id":                 1,
	}

	linkedAccounts, err := u.linkedDal.FindOne(ctx, filter, projection)
	if err != nil {
		u.logger.Errorf("failed to fetch linked accounts: %v", err)
		return nil, fmt.Errorf("FAILED_TO_FETCH_LINKED_ACCOUNTS")
	}

	if linkedAccounts.CustomerNumber == "" {
		u.logger.Warnf("no linked account found for account number: %s", accNumber)
		return nil, fmt.Errorf("NO_DOCUMENT_FOUND")
	}
	// Step 2: Fetch user by customer number
	customerNumber := linkedAccounts.CustomerNumber
	filter = bson.M{
		"customer_number": customerNumber,
	}
	projection = bson.M{}

	userDoc, err := u.userDal.FindOne(ctx, filter, projection)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			u.logger.Warnf("user not found for customer number: %s", customerNumber)
			return nil, fmt.Errorf("NO_DOCUMENT_FOUND")
		}
		u.logger.Errorf("error fetching user: %v", err)
		return nil, fmt.Errorf("UNHANDLED_SERVER_ERROR")
	}

	// Step 3: Unmarshal user and attach linked account
	userData, err := local_util.JsonUnmarshal[bson.M](userDoc)
	if err != nil {
		u.logger.Errorf("failed to unmarshal user document: %v", err)
		return nil, fmt.Errorf("UNHANDLED_SERVER_ERROR")
	}
	(*userData)["linked_account"] = linkedAccounts

	var data any = userData
	return &data, nil
}

func (u *unlinkCustomer) UnlinkUserCif(ctx context.Context, userCode string) error {
	u.logger.Infof("Unlinking the user from the system requsting: %s", userCode)
	actionExist, err := u.checkPendingAction(ctx, string(model.RequestUnlinkUser))

	if err != nil {
		u.logger.Errorf("error checking pending action: %v", err)
		if err.Error() != "mongo: no documents in result" {
			return err
		}
	}
	if actionExist {
		u.logger.Warnf("pending request exists for unlink device, id: %s", userCode)
		return fmt.Errorf("PENDING_REQUEST_EXISTS")
	}

	prev, err := u.userDal.FindOne(ctx, bson.M{"user_code": userCode}, nil)
	if err != nil {
		return err
	}

	err = u.createCpsAction(ctx, userCode, prev, bson.M{"user_code": userCode}, string(model.RequestUnlinkUser))
	if err != nil {
		return nil
	}
	return nil
}
func (u *unlinkCustomer) Authorize(ctx context.Context, cpsAction any) (any, error) {

	cpsActionData, err := local_util.JsonUnmarshal[model.CPSAction](cpsAction)
	if err != nil {
		return nil, err
	}

	userCode := cpsActionData.UniqueId
	if userCode == "" {
		return nil, fmt.Errorf("USERCODE_CANT_BE_EMPTY")
	}
	filter := bson.M{
		"user_code": userCode,
	}
	user, err := u.userDal.FindOne(ctx, filter, bson.M{})
	if err != nil {
		return nil, err
	}
	filter = bson.M{
		"customer_number": user.CustomerNumber,
	}
	linkedAccount, err := u.linkedDal.FindAll(ctx, filter, bson.M{})
	if err != nil {
		return nil, err
	}

	archivedUser, err := local_util.JsonUnmarshal[model.ArchivedUser](user)
	if err != nil {
		return nil, err
	}
	_, err = u.archivedUserDal.InsertOne(ctx, *archivedUser)
	if err != nil {
		return nil, err
	}
	var archivedLinkedAccount []*model.ArchivedLinkedAccount
	for _, account := range linkedAccount {
		acc, err := local_util.JsonUnmarshal[model.ArchivedLinkedAccount](account)
		if err != nil {
			return nil, err
		}
		archivedLinkedAccount = append(archivedLinkedAccount, acc)
	}

	for _, account := range archivedLinkedAccount {
		_, err = u.archivedLinkedDal.InsertOne(ctx, *account)
		if err != nil {
			return nil, err
		}
	}

	err = u.HardDeleteByID(ctx, u.client, u.dbName, u.collection[1], user.ID)
	if err != nil {
		return nil, err
	}

	for _, account := range linkedAccount {
		err = u.HardDeleteByID(ctx, u.client, u.dbName, u.collection[4], account.ID)
		if err != nil {
			return nil, err
		}
	}

	return cpsActionData, nil
}

func (u *unlinkCustomer) createCpsAction(ctx context.Context, id string, previousAction any, currentAction any, requestAction string) error {
	userData := contexts.ExtractContext(ctx)

	cpsAction := &model.CPSAction{
		ID:               bson.NewObjectID(),
		ActionCode:       local_util.GenerateUniqueActionCode(20),
		UniqueId:         id,
		MakerID:          userData.UserID,
		MakerName:        userData.FullName,
		MakerPhoneNumber: userData.PhoneNumber,
		Department:       userData.Department,
		ActionStatus:     "PENDING",
		CurrentAction:    currentAction,
		PreviousAction:   previousAction,
		IsDeleted:        false,
		CreatedAt:        time.Now(),
		LastModifiedAt:   time.Now(),
		RequestAction:    requestAction,
		MakerActionTime:  time.Now(),
	}

	u.logger.Infof("Creating CPS action: %+v", cpsAction)
	_, err := u.cpsDal.InsertOne(ctx, *cpsAction)
	if err != nil {
		u.logger.Errorf("failed to insert CPS action: %v", err)
		return common.DefineError.General["UNHANDLED_SERVER_ERROR"]
	}

	u.logger.Infof("Successfully created CPS action for id: %s", id)
	return nil
}
func (u *unlinkCustomer) checkPendingAction(ctx context.Context, requestAction string) (bool, error) {
	userData := contexts.ExtractContext(ctx)
	filter := bson.M{
		"maker_id":       userData.UserID,
		"department":     userData.Department,
		"action_status":  "PENDING",
		"request_action": requestAction,
	}

	u.logger.Infof("Checking for pending action with filter: %+v", filter)
	dataCpsAction, err := u.cpsDal.FindOne(ctx, filter, nil)
	if err != nil {
		u.logger.Errorf("error finding pending action: %v", err)
		return false, nil
	}

	if dataCpsAction == nil {
		u.logger.Infof("No pending action found for filter: %+v", filter)
		return false, nil
	}

	u.logger.Warnf("Pending action exists for filter: %+v", filter)
	return true, nil
}
func (u *unlinkCustomer) HardDeleteByID(ctx context.Context, client *mongo.Client, databaseName string, collectionName string, id interface{}) error {
	db := client.Database(databaseName)
	collection := db.Collection(collectionName)
	filter := bson.M{"_id": id}

	res, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		u.logger.Errorf("failed to hard delete document from %s: %v", collectionName, err)
		return err
	}
	if res.DeletedCount == 0 {
		u.logger.Warnf("no document found to delete in %s with id: %v", collectionName, id)
		return fmt.Errorf("NO_DOCUMENT_FOUND")
	}
	u.logger.Infof("Successfully hard deleted document from %s with id: %v", collectionName, id)
	return nil
}
