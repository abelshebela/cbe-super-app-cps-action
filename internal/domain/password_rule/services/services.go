package services

import (
    "context"

    "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/action"
    "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/password_rule/repository"
)

type PasswordRuleService interface {
    RequestPasswordRuleUpdate(ctx context.Context, rule *action.PasswordRule, maker action.User) (string, error)
    ApproveOrRejectPasswordRuleAction(ctx context.Context, actionID string, approve bool, checker action.User, rejectionReason *string) error
    GetPasswordRuleUpdateActionByID(ctx context.Context, actionID string) (*action.CPSAction, error)
}

type passwordRuleService struct {
    repo repository.PasswordRuleRepository
}

func NewPasswordRuleService(repo repository.PasswordRuleRepository) PasswordRuleService {
    return &passwordRuleService{repo: repo}
}

func (s *passwordRuleService) RequestPasswordRuleUpdate(ctx context.Context, rule *action.PasswordRule, maker action.User) (string, error) {
    return s.repo.CreatePasswordRuleUpdateAction(ctx, rule, maker)
}

func (s *passwordRuleService) ApproveOrRejectPasswordRuleAction(ctx context.Context, actionID string, approve bool, checker action.User, rejectionReason *string) error {
    return s.repo.ApproveOrRejectPasswordRuleAction(ctx, actionID, approve, checker, rejectionReason)
}

func (s *passwordRuleService) GetPasswordRuleUpdateActionByID(ctx context.Context, actionID string) (*action.CPSAction, error) {
    return s.repo.GetPasswordRuleUpdateActionByID(ctx, actionID)
}