package budget_category

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	action_entity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	budget_entity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/budget_category"
	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/budget_category"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type BudgetCategoryRepo struct {
	client *mongo.Client
	logger utils.Logger
	cpsDal dal.MongoDal[action_entity.CPSAction, action_entity.CPSAction]
	dal    dal.MongoDal[budget_entity.BudgetCategory, budget_entity.BudgetCategory]
}

func NewBudgetCategoryRepo(client *mongo.Client, dbName string, logger utils.Logger) domain.BudgetCategoryRepository {
	dalBudget := dal.NewMongoDal[budget_entity.BudgetCategory, budget_entity.BudgetCategory](client, dbName, "budget_categories")
	dalCPS := dal.NewMongoDal[action_entity.CPSAction, action_entity.CPSAction](client, dbName, "cps_actions")
	return &BudgetCategoryRepo{
		client: client,
		dal:    dalBudget,
		cpsDal: dalCPS,
		logger: logger,
	}
}

type ActionRequest interface {
	dto.CreateBudgetCategoryRequest |
		dto.UpdateBudgetCategoryRequest |
		dto.DeleteBudgetCategoryRequest
	GetActionType() action_entity.ActionType
	GetRequestAction() string
}

func (b *BudgetCategoryRepo) CreateAction(
	ctx context.Context,
	data interface{},
	maker action_entity.User,
) (action_entity.CPSAction, error) {

	var actionType action_entity.ActionType
	var requestAction string

	switch req := data.(type) {
	case dto.CreateBudgetCategoryRequest:
		actionType = action_entity.ActionCreate
		requestAction = req.RequestAction
	case dto.UpdateBudgetCategoryRequest:
		actionType = action_entity.ActionUpdate
		requestAction = req.RequestAction
	case dto.DeleteBudgetCategoryRequest:
		actionType = action_entity.ActionDelete
		requestAction = req.RequestAction
	default:
		return action_entity.CPSAction{}, fmt.Errorf("unsupported action type: %T", data)
	}

	actionId := utils.Random(10, &utils.PreSufix{Prefix: "CPS_"})
	cpsAction := action_entity.CPSAction{
		ActionCode:       actionId,
		MakerID:          maker.UserID,
		MakerName:        maker.FullName,
		MakerPhoneNumber: maker.PhoneNumber,
		ActionType:       actionType,
		RequestAction:    action_entity.RequestAction(requestAction),
		ActionStatus:     action_entity.ActionPending,
		CurrentAction:    data,
		CreatedAt:        time.Now(),
		LastModifiedAt:   time.Now(),
	}

	result, err := b.cpsDal.InsertOne(ctx, cpsAction)
	if err != nil {
		return action_entity.CPSAction{}, fmt.Errorf("failed to create action: %w", err)
	}

	return result, nil
}

func (b *BudgetCategoryRepo) UpdateAction(ctx context.Context, actionId string, checker action_entity.User, status action_entity.ActionStatus) (action_entity.CPSAction, error) {
	update := bson.M{
		"action_status":    status,
		"checker":          checker,
		"last_modified_at": time.Now(),
	}
	objID, err := bson.ObjectIDFromHex(actionId)
	if err != nil {
		b.logger.Errorf("failed to convert action ID to ObjectID: %v", err)
		return action_entity.CPSAction{}, fmt.Errorf("INVALID_ID")
	}
	return b.cpsDal.UpdateOne(ctx, bson.M{"_id": objID}, update)
}

func (b *BudgetCategoryRepo) FindActionById(ctx context.Context, actionId string) (*action_entity.CPSAction, error) {
	oid, err := bson.ObjectIDFromHex(actionId)
	if err != nil {
		b.logger.Errorf("failed to convert action ID to ObjectID: %v", err)
		return nil, err
	}
	return b.cpsDal.FindOne(ctx, bson.M{"_id": oid}, bson.M{})
}

func (b *BudgetCategoryRepo) ApproveAction(ctx context.Context, approveRequest dto.ApproveBudgetCategoryRequest, checker action_entity.User) (action_entity.CPSAction, error) {
	oid, err := bson.ObjectIDFromHex(approveRequest.ActionID)
	if err != nil {
		b.logger.Errorf("failed to convert action ID to ObjectID: %v", err)
		return action_entity.CPSAction{}, err
	}
	update := bson.M{
		"action_status":    action_entity.ActionApproved,
		"checker":          checker,
		"last_modified_at": time.Now(),
	}
	return b.cpsDal.UpdateOne(ctx, bson.M{"_id": oid}, update)
}

func (b *BudgetCategoryRepo) CreateBudgetCategory(ctx context.Context, req dto.CreateBudgetCategoryRequest) (budget_entity.BudgetCategory, error) {
	budgetCategory := budget_entity.BudgetCategory{
		Name:          req.Name,
		Icon:          req.Icon,
		Description:   req.Description,
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
		IsDeleted:     false,
	}
	return b.dal.InsertOne(ctx, budgetCategory)
}

func (b *BudgetCategoryRepo) FindBudgetCategoryById(ctx context.Context, id string) (*budget_entity.BudgetCategory, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		b.logger.Errorf("failed to convert action ID to ObjectID: %v", err)
		return nil, err
	}
	filter := bson.M{
		"_id":        oid,
		"is_deleted": false,
	}
	return b.dal.FindOne(ctx, filter, bson.M{})
}

func (b *BudgetCategoryRepo) UpdateBudgetCategory(ctx context.Context, req dto.UpdateBudgetCategoryRequest) (budget_entity.BudgetCategory, error) {
	oid, err := bson.ObjectIDFromHex(req.ID)
	if err != nil {
		b.logger.Errorf("failed to convert action ID to ObjectID: %v", err)
		return budget_entity.BudgetCategory{}, fmt.Errorf("invalid ID: %w", err)
	}

	update := bson.M{
		"last_updated_at": time.Now(),
	}

	for key, value := range map[string]string{
		"name":        req.Name,
		"icon":        req.Icon,
		"description": req.Description,
	} {
		if strings.TrimSpace(value) != "" {
			update[key] = value
		}
	}

	if len(update) == 1 {
		b.logger.Errorf("no fields to update")
		return budget_entity.BudgetCategory{}, fmt.Errorf("no fields to update")
	}

	filter := bson.M{
		"_id":        oid,
		"is_deleted": false,
	}

	return b.dal.UpdateOne(ctx, filter, update)
}

func (b *BudgetCategoryRepo) DeleteBudgetCategory(ctx context.Context, req dto.DeleteBudgetCategoryRequest) error {
	oid, err := bson.ObjectIDFromHex(req.ID)
	if err != nil {
		b.logger.Errorf("failed to convert action ID to ObjectID: %v", err)
		return err
	}
	return b.dal.DeleteOne(ctx, bson.M{
		"_id":        oid,
		"is_deleted": false,
	})
}

func (b *BudgetCategoryRepo) GetBudgetCategory(ctx context.Context, req dto.GetBudgetCategoryRequest) (*budget_entity.BudgetCategory, error) {
	oid, err := bson.ObjectIDFromHex(req.ID)
	if err != nil {
		b.logger.Errorf("failed to convert action ID to ObjectID: %v", err)
		return nil, err
	}
	filter := bson.M{
		"_id":        oid,
		"is_deleted": false,
	}
	return b.dal.FindOne(ctx, filter, bson.M{})
}

func (b *BudgetCategoryRepo) GetAllBudgetCategory(ctx context.Context, req dto.GetAllBudgetCategoryRequest) ([]*budget_entity.BudgetCategory, error) {
	filter := bson.M{
		"is_deleted": false,
	}
	return b.dal.FindAllWithPagination(ctx, filter, bson.M{}, req.Page, req.Limit)
}
