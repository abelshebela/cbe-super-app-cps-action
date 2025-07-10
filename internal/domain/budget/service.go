package budget

import (
	"context"
	"errors"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/budget/entities"
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

func (s *BudgetService) FetchIcons(ctx context.Context) ([]*entities.Icon, error) {
	icons, err := s.repo.FetchIcons(ctx)
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
	action, err := s.repo.CreateColor(ctx, hexCode, cpsAction)
	if err != nil {
		return nil, err
	}

	return action, nil
}

func (s *BudgetService) FetchColors(ctx context.Context) ([]*entities.Color, error) {
	colors, err := s.repo.ListAllColor(ctx)
	if err != nil {
		return nil, err
	}

	return colors, nil
}

func (s *BudgetService) UpdateColor(ctx context.Context, id, hexCode string, cpsAction entities.CPSAction) (*entities.CPSAction, error) {
	existing, err := s.repo.GetByIDColor(ctx, id)
	if err != nil || existing == nil {
		return nil, errors.New("color not found")
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
