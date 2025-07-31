package cpsaction

import (
	"context"
	"errors"
	"fmt"
	"time"

	infra_mongo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/infra"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/mappers"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	entity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/cps_actions"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type cpsActionStore struct {
	MongoCPSAction infra_mongo.MongoDal[model.CPSAction, model.CPSAction]
	client         *mongo.Client
	logger         utils.Logger
}

func NewOutBoundStore(client *mongo.Client, dbName string, collectionName string, logger utils.Logger) repo.CPSActionRepository {
	MongoCPSAction := infra_mongo.NewMongoDal[model.CPSAction, model.CPSAction](client, dbName, collectionName)
	return &cpsActionStore{
		MongoCPSAction: MongoCPSAction,
		logger:         logger,
		client:         client,
	}
}

func (p *cpsActionStore) RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	p.logger.Debugf("Starting MongoDB session for transaction")

	session, err := p.client.StartSession()
	if err != nil {
		p.logger.Errorf("failed to start MongoDB session: %v", err)
		return fmt.Errorf(common_util.UnhandledServerError)
	}
	defer session.EndSession(ctx)

	return mongo.WithSession(ctx, session, func(txCtx context.Context) error {
		p.logger.Debugf("Starting MongoDB transaction")

		if err := session.StartTransaction(); err != nil {
			p.logger.Errorf("failed to start transaction: %v", err)
			return fmt.Errorf(common_util.UnhandledServerError)
		}

		err := fn(txCtx)
		if err != nil {
			p.logger.Errorf("transaction logic failed: %v", err)
			if abortErr := session.AbortTransaction(txCtx); abortErr != nil {
				p.logger.Errorf("failed to abort transaction: %v", abortErr)
			} else {
				p.logger.Debugf("Transaction aborted successfully")
			}
			return fmt.Errorf(common_util.UnhandledServerError)
		}

		if err := session.CommitTransaction(txCtx); err != nil {
			p.logger.Errorf("failed to commit transaction: %v", err)
			return fmt.Errorf(common_util.UnhandledServerError)
		}

		p.logger.Debugf("Transaction committed successfully")
		return nil
	})
}

func (a *cpsActionStore) CPSActionExists(ctx context.Context, user entity.CheckCPSAction) (bool, error) {
	filter := bson.M{
		"maker_phone_number": user.PhoneNumber,
		"action_status":      model.ActionPending,
		"department":         user.Department,
		"request_action":     user.RequestAction,
	}

	projection := bson.M{
		"action_code": 1,
		"_id":         1,
	}

	existingAction, err := a.MongoCPSAction.FindOne(ctx, filter, projection)
	if err != nil && err != mongo.ErrNoDocuments {
		a.logger.Errorf("failed to get ad", err)
		return false, fmt.Errorf("GENERAL_DB_QUERY_FAILED")
	}
	if existingAction != nil {
		a.logger.Infof("pending cps action present", user.FullName, user.UserCode, user.Department)
		return true, fmt.Errorf("PENDING_CPS_ACTION_PRESENT")
	}

	return false, nil
}

func (o *cpsActionStore) CreateCPSAction(ctx context.Context, action *entity.CPSAction) (*entity.CPSAction, error) {
	modelAction, err := mappers.DomainToModelCPSAction(*action)
	if err != nil {
		o.logger.Errorf("Domain to Model conversion failed: %v", err)
		return nil, err
	}
	data, err := o.MongoCPSAction.InsertOne(ctx, *modelAction)
	if err != nil {
		o.logger.Errorf("Insert CPSAction failed: %v", err)
		return nil, fmt.Errorf(common_util.GeneralDBQueryFailed)
	}
	o.logger.Infof("CPSAction created successfully with ActionCode: %s", data.ActionCode)
	return mappers.ModelToDomainCPSAction(data), nil
}

func (o *cpsActionStore) UpdateCPSAction(ctx context.Context, action *entity.CPSAction) (*entity.CPSAction, error) {
	modelAction, err := mappers.DomainToModelCPSAction(*action)
	if err != nil {
		o.logger.Errorf("Domain to Model conversion failed: %v", err)
		return nil, err
	}

	filter := bson.M{"action_code": modelAction.ActionCode}
	update := mappers.BuildCPSActionUpdate(modelAction)
	res, err := o.MongoCPSAction.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			o.logger.Warnf("No document found with ActionCode: %s", modelAction.ActionCode)
			return nil, fmt.Errorf(common_util.NotFound)
		}
		o.logger.Errorf("Update CPSAction failed: %v", err)
		return nil, fmt.Errorf(common_util.GeneralDBUpdateFailed)
	}
	o.logger.Infof("CPSAction updated successfully: %s", modelAction.ActionCode)
	return mappers.ModelToDomainCPSAction(res), nil
}

func (o *cpsActionStore) ApproveCPSAction(ctx context.Context, action *entity.AuthorizeCPSAction) (*entity.CPSAction, error) {
	o.logger.Infof("Approving CPSAction with ActionCode: %s", action.ActionCode)
	return o.updateCPSActionStatus(ctx, action, string(model.ActionApproved))
}

func (o *cpsActionStore) RejectCPSAction(ctx context.Context, action *entity.AuthorizeCPSAction) (*entity.CPSAction, error) {
	o.logger.Infof("Rejecting CPSAction with ActionCode: %s", action.ActionCode)
	if action.RejectionReason == "" {
		o.logger.Warnf("Rejection reason is empty for ActionCode: %s", action.ActionCode)
		return nil, fmt.Errorf(common_util.InvalidInput)
	}
	return o.updateCPSActionStatus(ctx, action, string(model.ActionRejected))
}

func (o *cpsActionStore) updateCPSActionStatus(ctx context.Context, action *entity.AuthorizeCPSAction, newStatus string) (*entity.CPSAction, error) {

	filter := bson.M{
		"action_code":   action.ActionCode,
		"department":    action.Department,
		"action_status": model.ActionPending,
	}

	now := time.Now()
	update := bson.M{
		"checker_id":           action.CheckerUser.UserCode,
		"checker_name":         action.CheckerUser.FullName,
		"checker_phone_number": action.CheckerUser.PhoneNumber,
		"action_status":        newStatus,
		"checker_action_time":  now,
		"rejection_reason":     action.RejectionReason,
		"last_modified_at":     now,
	}

	cpsAction, err := o.MongoCPSAction.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			o.logger.Errorf("cps action not found for update, action_code: %s", action.ActionCode)
			return nil, fmt.Errorf(common_util.ActionNotFound)
		}
		o.logger.Errorf("failed to update cps action, action_code: %s, error: %v", action.ActionCode, err)
		return nil, fmt.Errorf(common_util.UnhandledServerError)
	}

	o.logger.Infof("CPSAction status updated successfully: %s", action.ActionCode)
	return mappers.ModelToDomainCPSAction(cpsAction), nil
}

func (o *cpsActionStore) GetCPSActionsByDepartment(ctx context.Context, department string, status string, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entity.CPSAction], error) {
	o.logger.Infof("Getting CPSActions for department: %s", department)
	filter := bson.M{
		"department":    department,
		"action_status": status,
	}

	if filterParams.Search != "" {
		search := filterParams.Search
		filter["$or"] = []bson.M{
			{"action_code": bson.M{"$regex": search, "$options": "i"}},
			{"action_type": bson.M{"$regex": search, "$options": "i"}},
			{"request_action": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	fmt.Println(filter, "filter")
	page := filterParams.Page
	limit := filterParams.PerPage
	skip := (page - 1) * limit

	cpsActionsDocs, err := o.MongoCPSAction.FindAllWithPagination(ctx, filter, bson.M{}, int64(skip), int64(limit))
	if err != nil {
		o.logger.Errorf("failed to fetch cps actions: %v", err)
		return nil, fmt.Errorf(common_util.UnhandledServerError)
	}

	var actions []*entity.CPSAction
	for _, doc := range cpsActionsDocs {
		actions = append(actions, mappers.ModelToDomainCPSAction(*doc))
	}

	total, err := o.MongoCPSAction.TotalCount(ctx, filter)
	if err != nil {
		o.logger.Errorf("failed to get cps actions total counts, error: %v", err)
		return nil, fmt.Errorf(common_util.UnhandledServerError)
	}

	o.logger.Infof("retrieved cps actions, page: %d, count: %d, total: %d", page, len(actions), total)
	meta := common_util.BuildPaginationMeta(total, page, limit)

	return &common_util.PaginatedResponse[[]*entity.CPSAction]{
		Data: actions,
		Meta: meta,
	}, nil
}

func (o *cpsActionStore) GetCPSActionByID(ctx context.Context, id string) (*entity.CPSAction, error) {
	o.logger.Infof("Getting CPSAction by ID: %s", id)
	objID, err := mappers.ObjectIDFromHex(id)
	if err != nil {
		o.logger.Errorf("Invalid ObjectID: %s, error: %v", id, err)
		return nil, err
	}
	filter := bson.M{"_id": objID}
	return o.getCPSActionByFilter(ctx, filter)
}

func (o *cpsActionStore) GetCPSActionByActionCode(ctx context.Context, uniqueID string) (*entity.CPSAction, error) {
	o.logger.Infof("Getting CPSAction by UniqueID: %s", uniqueID)
	filter := bson.M{"action_code": uniqueID}
	return o.getCPSActionByFilter(ctx, filter)
}

func (o *cpsActionStore) getCPSActionByFilter(ctx context.Context, filter bson.M) (*entity.CPSAction, error) {
	o.logger.Infof("Fetching CPSAction with filter: %+v", filter)
	action, err := o.MongoCPSAction.FindOne(ctx, filter, bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			o.logger.Warnf("CPSAction not found with filter: %+v", filter)
			return nil, fmt.Errorf(common_util.AccountNotFound)
		}
		o.logger.Errorf("Failed to fetch CPSAction, error: %v", err)
		return nil, err
	}
	o.logger.Infof("Successfully fetched CPSAction")
	return mappers.ModelToDomainCPSAction(*action), nil
}
