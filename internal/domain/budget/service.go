package budget

import (
	"context"
	"fmt"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/budget/entities"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type BudgetService struct {
	repo   Repository
	logger utils.Logger
}

func InitBudgetDomain(repo Repository, logger utils.Logger) *BudgetService {
	return &BudgetService{
		repo:   repo,
		logger: logger,
	}
}

func (s *BudgetService) CreateIcon(ctx context.Context, cpsAction entities.CPSAction) (*entities.CPSAction, error) {
	action, err := s.repo.CreateIconAction(ctx, cpsAction)
	if err != nil {
		return nil, err
	}
	return action, nil
}

func (s *BudgetService) FetchIcons(ctx context.Context, filterParams *constant.Filter) (*entities.FetchIconResponse, error) {
	icons, err := s.repo.FetchIcons(ctx, filterParams)
	if err != nil {
		return nil, err
	}

	return icons, nil
}

func (s *BudgetService) UpdateIcon(ctx context.Context, id string, cpsAction entities.CPSAction) (*entities.CPSAction, error) {
	action, err := s.repo.UpdateIcon(ctx, id, cpsAction)
	if err != nil {
		return nil, err
	}

	return action, nil
}

func (s *BudgetService) CreateColor(ctx context.Context, hexCode string, cpsAction entities.CPSAction) (*entities.CPSAction, error) {
	exist, err := s.repo.CheckColorExist(ctx, hexCode)
	if err != nil {
		return nil, err
	}

	if exist {
		s.logger.Errorf("Color %s already exists", hexCode)
		return nil, fmt.Errorf("CONFLICT_KEY")
	}
	action, err := s.repo.CreateColor(ctx, hexCode, cpsAction)
	if err != nil {
		return nil, err
	}

	return action, nil
}

func (s *BudgetService) FetchColors(ctx context.Context, filterParams *constant.Filter) (*entities.FetchColorsResponse, error) {
	colors, err := s.repo.ListAllColor(ctx, filterParams)
	if err != nil {
		return nil, err
	}

	return colors, nil
}

func (s *BudgetService) UpdateColor(ctx context.Context, id, hexCode string, cpsAction entities.CPSAction) (*entities.CPSAction, error) {
	existing, err := s.repo.GetByIDColor(ctx, id)
	if err != nil || existing == nil {
		return nil, err
	}

	colors, err := s.repo.CreateAction(ctx, cpsAction)
	if err != nil {
		return nil, err
	}

	return colors, nil
}

func (s *BudgetService) ApproveAction(ctx context.Context, cpsAction entities.CPSAction) (*entities.CPSAction, error) {
	action, err := s.repo.ApproveAction(ctx, cpsAction)
	if err != nil {
		return nil, err
	}
	return action, nil
}
