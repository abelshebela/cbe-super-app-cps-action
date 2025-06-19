package services

import (
	"context"
	"fmt"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/action"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/password_rule/repository"
)

type PasswordRuleService interface {
	RequestPasswordRuleUpdate(ctx context.Context, rule *action.PasswordRule, maker action.User) (string, error)
	ApproveOrRejectPasswordRuleAction(ctx context.Context, actionID string, approve bool, checker action.User, rejectionReason *string) error
	GetPasswordRuleUpdateActionByID(ctx context.Context, actionID string) (*action.CPSAction, error)
	GetUpdateAction(ctx context.Context, maker action.User) (*action.CPSAction, error) 

}

type passwordRuleService struct {
	repo repository.PasswordRuleRepository
}

func NewPasswordRuleService(repo repository.PasswordRuleRepository) PasswordRuleService {
	return &passwordRuleService{repo: repo}
}

func (s *passwordRuleService) RequestPasswordRuleUpdate(ctx context.Context, rule *action.PasswordRule, maker action.User) (string, error) {
    existing, err := s.repo.GetUpdateAction(ctx, maker)
    if err != nil {
        return "", err
    }
    if existing != nil && existing.ActionStatus == action.ActionPending && existing.ActionType == action.ActionUpdate {
        return "", fmt.Errorf("pending update action already exists for this maker")
    }
    return s.repo.CreatePasswordRuleUpdateAction(ctx, rule, maker)
}

func (s *passwordRuleService) ApproveOrRejectPasswordRuleAction(ctx context.Context, actionID string, approve bool, checker action.User, rejectionReason *string) error {
	return s.repo.ApproveOrRejectPasswordRuleAction(ctx, actionID, approve, checker, rejectionReason)
}

func (s *passwordRuleService) GetPasswordRuleUpdateActionByID(ctx context.Context, actionID string) (*action.CPSAction, error) {
	return s.repo.GetPasswordRuleUpdateActionByID(ctx, actionID)
}
func (s *passwordRuleService) GetUpdateAction(ctx context.Context, maker action.User) (*action.CPSAction, error) {
	return s.repo.GetUpdateAction(ctx, maker)
}
