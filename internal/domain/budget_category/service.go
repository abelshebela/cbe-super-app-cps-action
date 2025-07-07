package budget_category

import (
	"context"
	"encoding/json"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	action "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type BudgetCategoryService struct {
	repository Repository
	actionRepo action.IActionRepository
	logger     utils.Logger
}

func NewBudgetCategoryService(repository Repository, actionRepo action.IActionRepository, logger utils.Logger) *BudgetCategoryService {
	return &BudgetCategoryService{
		repository: repository,
		actionRepo: actionRepo,
		logger:     logger,
	}
}

func (s *BudgetCategoryService) CreateBudgetCategoryAction(ctx context.Context, req dto.CreateBudgetCategoryRequest, maker action.User) (string, error) {
	actionId := utils.Random(10, &utils.PreSufix{Prefix: "CPS_"})
	cpsAction := action.CPSAction{
		ActionCode:     actionId,
		Maker:          maker,
		ActionType:     action.ActionCreate,
		RequestAction:  action.RequestAction("CREATE_BUDGET_CATEGORY"),
		ActionStatus:   action.ActionPending,
		CurrentAction:  req,
		CreatedAt:      time.Now(),
		LastModifiedAt: time.Now(),
	}
	_, err := s.actionRepo.CreateCpsAction(ctx, cpsAction)
	if err != nil {
		return "", err
	}
	return actionId, nil
}

func (s *BudgetCategoryService) UpdateBudgetCategoryAction(ctx context.Context, req dto.UpdateBudgetCategoryRequest, maker action.User) (string, error) {
	actionId := utils.Random(10, &utils.PreSufix{Prefix: "CPS_"})
	cpsAction := action.CPSAction{
		ActionCode:     actionId,
		Maker:          maker,
		ActionType:     action.ActionUpdate,
		RequestAction:  action.RequestAction("UPDATE_BUDGET_CATEGORY"),
		ActionStatus:   action.ActionPending,
		CurrentAction:  req,
		CreatedAt:      time.Now(),
		LastModifiedAt: time.Now(),
	}
	_, err := s.actionRepo.CreateCpsAction(ctx, cpsAction)
	if err != nil {
		return "", err
	}
	return actionId, nil
}

func (s *BudgetCategoryService) DeleteBudgetCategoryAction(ctx context.Context, req dto.DeleteBudgetCategoryRequest, maker action.User) (string, error) {
	actionId := utils.Random(10, &utils.PreSufix{Prefix: "CPS_"})
	cpsAction := action.CPSAction{
		ActionCode:     actionId,
		Maker:          maker,
		ActionType:     action.ActionDelete,
		RequestAction:  action.RequestAction("DELETE_BUDGET_CATEGORY"),
		ActionStatus:   action.ActionPending,
		CurrentAction:  req,
		CreatedAt:      time.Now(),
		LastModifiedAt: time.Now(),
	}
	_, err := s.actionRepo.CreateCpsAction(ctx, cpsAction)
	if err != nil {
		return "", err
	}
	return actionId, nil
}
func (s *BudgetCategoryService) ApproveBudgetCategoryAction(ctx context.Context, actionId string, approve bool, checker action.User) error {
	cpsAction, err := s.actionRepo.FetchCpsActionById(ctx, actionId)
	if err != nil {
		return err
	}

	if cpsAction.ActionStatus == action.ActionApproved {
		return dto.ErrActionApproved
	}

	if cpsAction.ActionStatus == action.ActionRejected {
		return dto.ErrActionRejected
	}

	var operationErr error

	if approve {
		cpsAction.ActionStatus = action.ActionApproved

		switch cpsAction.ActionType {
		case action.ActionCreate:
			var req dto.CreateBudgetCategoryRequest
			if err := bindAction(cpsAction.CurrentAction, &req); err != nil {
				return err
			}
			_, operationErr = s.CreateBudgetCategory(ctx, req)

		case action.ActionUpdate:
			var req dto.UpdateBudgetCategoryRequest
			if err := bindAction(cpsAction.CurrentAction, &req); err != nil {
				return err
			}
			_, operationErr = s.UpdateBudgetCategory(ctx, req)
			if operationErr == nil {
				s.logger.Infof("Budget category updated successfully")
			}

		case action.ActionDelete:
			var req dto.DeleteBudgetCategoryRequest
			if err := bindAction(cpsAction.CurrentAction, &req); err != nil {
				return err
			}
			operationErr = s.DeleteBudgetCategory(ctx, req)
		}
	} else {
		cpsAction.ActionStatus = action.ActionRejected

		if cpsAction.ActionType == action.ActionCreate {
			// var req dto.CreateBudgetCategoryRequest
			// if err := bindAction(cpsAction.CurrentAction, &req); err != nil {
			// 	return err
			// }

			// deleteObjectBody := config.DeleteObjectBody{
			// 	BucketName: req.BucketName,
			// 	ObjectName: req.ObjectName,
			// }

			// if _, err := s.minioClient.DeleteObject(ctx, deleteObjectBody); err != nil {
			// 	s.logger.Errorf("failed to delete object from minio: %v", err)
			// 	return err
			// }
		}
	}

	cpsAction.Checker = checker
	cpsAction.LastModifiedAt = time.Now()

	if _, updateErr := s.actionRepo.UpdateCpsAction(ctx, *cpsAction); updateErr != nil {
		if operationErr != nil {
			s.logger.Errorf("operationErr: %v", operationErr)
		}
		return updateErr
	}

	return operationErr
}

func (s *BudgetCategoryService) CreateBudgetCategory(ctx context.Context, budgetCategory dto.CreateBudgetCategoryRequest) (BudgetCategory, error) {
	return s.repository.CreateBudgetCategory(ctx, budgetCategory)
}

func (s *BudgetCategoryService) UpdateBudgetCategory(ctx context.Context, budgetCategory dto.UpdateBudgetCategoryRequest) (BudgetCategory, error) {
	return s.repository.UpdateBudgetCategory(ctx, budgetCategory)
}

func (s *BudgetCategoryService) DeleteBudgetCategory(ctx context.Context, budgetCategory dto.DeleteBudgetCategoryRequest) error {
	return s.repository.DeleteBudgetCategory(ctx, budgetCategory)
}

func (s *BudgetCategoryService) GetBudgetCategory(ctx context.Context, budgetCategory dto.GetBudgetCategoryRequest) (*BudgetCategory, error) {
	return s.repository.GetBudgetCategory(ctx, budgetCategory)
}

func (s *BudgetCategoryService) GetAllBudgetCategory(ctx context.Context, getAllBudgetCategory dto.GetAllBudgetCategoryRequest) ([]*BudgetCategory, error) {
	return s.repository.GetAllBudgetCategory(ctx, getAllBudgetCategory)
}

func bindAction(source interface{}, target interface{}) error {
	bytes, err := json.Marshal(source)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, target)
}
