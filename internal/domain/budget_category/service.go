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
	repository BudgetCategoryRepository
	logger     utils.Logger
}

func NewBudgetCategoryService(repository BudgetCategoryRepository, logger utils.Logger) *BudgetCategoryService {
	return &BudgetCategoryService{
		repository: repository,
		logger:     logger,
	}
}

func (s *BudgetCategoryService) CreateAction(ctx context.Context, req any, maker action.User) (action.CPSAction, error) {
	return s.repository.CreateAction(ctx, req, maker)
}

func (s *BudgetCategoryService) ApproveAction(ctx context.Context, approveRequest dto.ApproveBudgetCategoryRequest, checker action.User) (action.CPSAction, error) {
	cpsAction, err := s.repository.FindActionById(ctx, approveRequest.ActionID)
	if err != nil {
		s.logger.Errorf("find action by id error: %v", err)
		return action.CPSAction{}, err
	}

	if cpsAction.ActionStatus == action.ActionApproved {
		s.logger.Errorf("action already approved")
		return action.CPSAction{}, dto.ErrActionApproved
	}

	if cpsAction.ActionStatus == action.ActionRejected {
		s.logger.Errorf("action already rejected")
		return action.CPSAction{}, dto.ErrActionRejected
	}

	var operationErr error

	if approveRequest.Status == "APPROVED" {
		cpsAction.ActionStatus = action.ActionApproved
		operationErr = s.executeApprovedAction(ctx, cpsAction)
	} else {
		cpsAction.ActionStatus = action.ActionRejected
		cpsAction.RejectionReason = &approveRequest.Reason
	}

	cpsAction.CheckerID = checker.UserID
	cpsAction.CheckerName = checker.FullName
	cpsAction.CheckerPhoneNumber = checker.PhoneNumber
	cpsAction.LastModifiedAt = time.Now()

	if _, updateErr := s.repository.UpdateAction(ctx, cpsAction.ID, checker, cpsAction.ActionStatus); updateErr != nil {
		if operationErr != nil {
			s.logger.Errorf("operation error: %v", operationErr)
		}
		return *cpsAction, updateErr
	}

	return *cpsAction, operationErr
}

func (s *BudgetCategoryService) executeApprovedAction(ctx context.Context, cpsAction *action.CPSAction) error {
	switch cpsAction.ActionType {
	case action.ActionCreate:
		var req dto.CreateBudgetCategoryRequest
		if err := bindAction(cpsAction.CurrentAction, &req); err != nil {
			return err
		}
		_, err := s.CreateBudgetCategory(ctx, req)
		return err

	case action.ActionUpdate:
		var req dto.UpdateBudgetCategoryRequest
		if err := bindAction(cpsAction.CurrentAction, &req); err != nil {
			return err
		}
		_, err := s.UpdateBudgetCategory(ctx, req)
		if err == nil {
			s.logger.Infof("Budget category updated successfully")
		}
		return err

	case action.ActionDelete:
		var req dto.DeleteBudgetCategoryRequest
		if err := bindAction(cpsAction.CurrentAction, &req); err != nil {
			return err
		}
		return s.DeleteBudgetCategory(ctx, req)

	default:
		return nil
	}
}

func (s *BudgetCategoryService) CreateBudgetCategory(ctx context.Context, budgetCategory dto.CreateBudgetCategoryRequest) (BudgetCategory, error) {
	return s.repository.CreateBudgetCategory(ctx, budgetCategory)
}
func (s *BudgetCategoryService) FindBudgetCategoryById(ctx context.Context, id string) (*BudgetCategory, error) {
	return s.repository.FindBudgetCategoryById(ctx, id)
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

func (s *BudgetCategoryService) FindActionById(ctx context.Context, actionId string) (*action.CPSAction, error) {
	return s.repository.FindActionById(ctx, actionId)
}

func (s *BudgetCategoryService) UpdateAction(ctx context.Context, actionId string, checker action.User, status action.ActionStatus) (action.CPSAction, error) {
	return s.repository.UpdateAction(ctx, actionId, checker, status)
}

func bindAction(source interface{}, target interface{}) error {
	bytes, err := json.Marshal(source)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, target)
}
