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
}

type UnlinkAccount interface {
	GetUserByAccount(ctx context.Context, accNumber string) (*local_util.PaginatedResponse[*any], error)
	UnlinkUserCif(ctx context.Context, userCode string) error
	Authorize(ctx context.Context, cpsAction any) error
}

func NewUnlinkPersistence(client *mongo.Client, database string, collection []string, logger utils.Logger) UnlinkAccount {

	return &unlinkCustomer{
		logger:            logger,
		userDal:           dal.NewMongoDal[member.User, member.User](client, database, collection[1]),
		cpsDal:            dal.NewMongoDal[model.CPSAction, model.CPSAction](client, database, collection[0]),
		archivedUserDal:   dal.NewMongoDal[model.ArchivedUser, model.ArchivedUser](client, database, collection[2]),
		archivedLinkedDal: dal.NewMongoDal[model.ArchivedLinkedAccount, model.ArchivedLinkedAccount](client, database, collection[3]),
	}
}

func (u *unlinkCustomer) GetUserByAccount(ctx context.Context, accNumber string) (*local_util.PaginatedResponse[*any], error) {

	if accNumber != "" {
		return nil, fmt.Errorf("ACCOUNT_NUMBER_CAN_NOT_BE_EMPTY")
	}
	filter := bson.M{
		"account_number": accNumber,
	}

	projection := bson.M{
		"account_number":      1,
		"linked_status":       1,
		"account_holder_name": 1,
		"linked_date":         1,
		"account_type":        1,
		"is_account_active":   1,
		"account_branch_code": 1,
		"linked_branch":       1,
		"_id":                 1,
	}

	linkedAccount, err := u.linkedDal.FindAllWithPagination(ctx, filter, projection, 1, 1)
	if err != nil {
		return nil, err
	}

	if len(linkedAccount) == 0 {
		u.logger.Warnf("no data found for minimum transfer cap with filter: %+v", filter)
		return nil, fmt.Errorf("NO_DOC_FOUND")
	}
	total, err := u.linkedDal.TotalCount(ctx, filter)
	meta := local_util.BuildPaginationMeta(total, 1, 1)
	if err != nil {
		return nil, fmt.Errorf("NO_DOC_FOUND")
	}

	filter = bson.M{
		"customer_code": linkedAccount[0].CustomerNumber,
	}
	projection = bson.M{
		"full_name": 1,
		"user_code": 1,
		"_id":       1,
	}
	user, err := u.userDal.FindOne(ctx, filter, projection)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, fmt.Errorf("NO_DOC_FOUND")
		}
		return nil, fmt.Errorf("")
	}

	userData, err := local_util.JsonUnmarshal[bson.M](user)
	if err != nil {
		return nil, fmt.Errorf("UNHANDLED_SERVER_ERROR")
	}

	(*userData)["linked_account"] = linkedAccount

	var data any = userData
	return &local_util.PaginatedResponse[*any]{
		Data: &data,
		Meta: meta,
	}, nil
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

	prev, err := u.userDal.FindOne(ctx, bson.M{"user_code": userCode, "is_deleted": false}, nil)
	if err != nil {
		return err
	}

	err = u.createCpsAction(ctx, userCode, prev, bson.M{"user_code": userCode}, string(model.RequestUnlinkUser))
	if err != nil {
		return nil
	}
	return nil
}
func (u *unlinkCustomer) Authorize(ctx context.Context, cpsAction any) error {

	current, err := local_util.JsonUnmarshal[model.CPSAction](cpsAction)
	if err != nil {
		return err
	}

	userCode := current.UniqueId
	if userCode == "" {
		return fmt.Errorf("USERCODE_CANT_BE_EMPTY")
	}
	filter := bson.M{
		"user_code": userCode,
	}
	user, err := u.userDal.FindOne(ctx, filter, bson.M{})
	if err != nil {
		return err
	}
	filter = bson.M{
		"customer_number": user.CustomerNumber,
	}
	linkedAccount, err := u.linkedDal.FindAll(ctx, filter, bson.M{})
	if err != nil {
		return err
	}

	archivedUser, err := local_util.JsonUnmarshal[model.ArchivedUser](user)
	if err != nil {
		return err
	}
	_, err = u.archivedUserDal.InsertOne(ctx, *archivedUser)
	if err != nil {
		return err
	}
	archivedLinkedAccount, err := local_util.JsonUnmarshal[model.ArchivedLinkedAccount](linkedAccount)
	if err != nil {
		return err
	}

	_, err = u.archivedLinkedDal.InsertOne(ctx, *archivedLinkedAccount)
	if err != nil {
		return err
	}

	err = u.userDal.DeleteOne(ctx, bson.M{"user_code": user.UserCode})
	if err != nil {
		return err
	}

	err = u.linkedDal.DeleteOne(ctx, bson.M{"customer_number": user.CustomerNumber})
	if err != nil {
		return err
	}
	return nil
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
