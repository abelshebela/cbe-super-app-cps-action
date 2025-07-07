package budget_category

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	action_entity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	budget_entity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/budget_category"
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

type BudgetCategoryRepoInterface interface {
	CreateBudgetCategoryAction(ctx context.Context, budgetCategory dto.CreateBudgetCategoryRequest, maker action_entity.User) (string, error)
	UpdateBudgetCategoryAction(ctx context.Context, budgetCategory dto.UpdateBudgetCategoryRequest, maker action_entity.User) (string, error)
	DeleteBudgetCategoryAction(ctx context.Context, budgetCategory dto.DeleteBudgetCategoryRequest, maker action_entity.User) (string, error)
	ApproveBudgetCategoryAction(ctx context.Context, actionId string, approve bool, checker action_entity.User) error

	CreateBudgetCategory(ctx context.Context, req dto.CreateBudgetCategoryRequest) (budget_entity.BudgetCategory, error)
	UpdateBudgetCategory(ctx context.Context, req dto.UpdateBudgetCategoryRequest) (budget_entity.BudgetCategory, error)
	DeleteBudgetCategory(ctx context.Context, req dto.DeleteBudgetCategoryRequest) error

	GetBudgetCategory(ctx context.Context, req dto.GetBudgetCategoryRequest) (*budget_entity.BudgetCategory, error)
	GetAllBudgetCategory(ctx context.Context, req dto.GetAllBudgetCategoryRequest) ([]*budget_entity.BudgetCategory, error)
}

func NewBudgetCategoryRepo(client *mongo.Client, dbName string) BudgetCategoryRepoInterface {
	dalBudget := dal.NewMongoDal[budget_entity.BudgetCategory, budget_entity.BudgetCategory](client, dbName, "budget_categories")
	dalCPS := dal.NewMongoDal[action_entity.CPSAction, action_entity.CPSAction](client, dbName, "cps_actions")
	return &BudgetCategoryRepo{
		client: client,
		dal:    dalBudget,
		cpsDal: dalCPS,
	}
}

func (b *BudgetCategoryRepo) CreateBudgetCategoryAction(ctx context.Context, req dto.CreateBudgetCategoryRequest, maker action_entity.User) (string, error) {
	actionId := utils.Random(10, &utils.PreSufix{Prefix: "CPS_"})
	cpsAction := action_entity.CPSAction{
		ActionCode:     actionId,
		Maker:          maker,
		ActionType:     action_entity.ActionCreate,
		RequestAction:  action_entity.RequestAction("CREATE_BUDGET_CATEGORY"),
		ActionStatus:   action_entity.ActionPending,
		CurrentAction:  req,
		CreatedAt:      time.Now(),
		LastModifiedAt: time.Now(),
		IsDeleted:      false,
	}
	_, err := b.cpsDal.InsertOne(ctx, cpsAction)
	if err != nil {
		return "", err
	}
	return actionId, nil
}

func (b *BudgetCategoryRepo) UpdateBudgetCategoryAction(ctx context.Context, req dto.UpdateBudgetCategoryRequest, maker action_entity.User) (string, error) {

	update := bson.M{
		"action_status":    action_entity.ActionPending,
		"maker":            maker,
		"action_type":      action_entity.ActionUpdate,
		"request_action":   action_entity.RequestAction("UPDATE_BUDGET_CATEGORY"),
		"current_action":   req,
		"last_modified_at": time.Now(),
	}
	_, err := b.cpsDal.UpdateOne(ctx, bson.M{"action_code": req.ActionCode}, update)
	if err != nil {
		return "", err
	}
	return req.ActionCode, nil
}

func (b *BudgetCategoryRepo) ApproveBudgetCategoryAction(ctx context.Context, actionId string, approve bool, checker action_entity.User) error {
	update := bson.M{
		"action_status":    action_entity.ActionApproved,
		"checker":          checker,
		"last_modified_at": time.Now(),
	}
	_, err := b.cpsDal.UpdateOne(ctx, bson.M{"action_code": actionId}, update)
	if err != nil {
		return err
	}
	return nil
}

func (b *BudgetCategoryRepo) DeleteBudgetCategoryAction(ctx context.Context, req dto.DeleteBudgetCategoryRequest, maker action_entity.User) (string, error) {
	actionId := utils.Random(10, &utils.PreSufix{Prefix: "CPS_"})
	return actionId, nil
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

func (b *BudgetCategoryRepo) UpdateBudgetCategory(ctx context.Context, req dto.UpdateBudgetCategoryRequest) (budget_entity.BudgetCategory, error) {
	oid, err := bson.ObjectIDFromHex(req.ID)
	if err != nil {
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
		return budget_entity.BudgetCategory{}, fmt.Errorf("no fields to update")
	}

	filter := bson.M{
		"_id":        oid,
		"is_deleted": false,
	}

	updated, err := b.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		return budget_entity.BudgetCategory{}, fmt.Errorf("failed to update category: %w", err)
	}

	return updated, nil
}

func (b *BudgetCategoryRepo) DeleteBudgetCategory(ctx context.Context, req dto.DeleteBudgetCategoryRequest) error {
	oid, err := bson.ObjectIDFromHex(req.ID)
	if err != nil {
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
