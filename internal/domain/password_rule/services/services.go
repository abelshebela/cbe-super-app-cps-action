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
	CheckPasswordRule(ctx context.Context, password string) (bool, string)
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

func (s *passwordRuleService) CheckPasswordRule(ctx context.Context, password string) (bool, string) {

	rule, err := s.repo.GetCurrentPasswordRule(ctx)
	if err != nil || rule == nil {
		return false, "could not retrieve password rule"
	}

	if len(password) < rule.MinLength {
		return false, fmt.Sprintf("password must be at least %d characters", rule.MinLength)
	}

	if rule.MaxLength > 0 && len(password) > rule.MaxLength {
		return false, fmt.Sprintf("password must be at most %d characters", rule.MaxLength)
	}

	if rule.Numbers {
		hasNumber := false
		for _, c := range password {
			if c >= '0' && c <= '9' {
				hasNumber = true
				break
			}
		}
		if !hasNumber {
			return false, "password must contain at least one number"
		}
	}

	if rule.CapitalLetters {
		hasUpper := false
		for _, c := range password {
			if c >= 'A' && c <= 'Z' {
				hasUpper = true
				break
			}
		}
		if !hasUpper {
			return false, "password must contain at least one uppercase letter"
		}
	}
	
	if rule.SmallLetters {
		hasLower := false
		for _, c := range password {
			if c >= 'a' && c <= 'z' {
				hasLower = true
				break
			}
		}
		if !hasLower {
			return false, "password must contain at least one lowercase letter"
		}
	}

	if rule.Characters {
		hasSpecial := false
		for _, c := range password {
			if (c >= 33 && c <= 47) || (c >= 58 && c <= 64) || (c >= 91 && c <= 96) || (c >= 123 && c <= 126) {
				hasSpecial = true
				break
			}
		}
		if !hasSpecial {
			return false, "password must contain at least one special character"
		}
	}
	return true, ""
}
