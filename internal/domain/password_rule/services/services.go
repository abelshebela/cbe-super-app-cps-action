package services

import (
	"context"
	"fmt"
	"time"

	"encoding/json"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/password_rule/repository"
)

type PasswordRuleService interface {
	RequestPasswordRuleUpdate(ctx context.Context, rule *action.PasswordRule, maker action.User, department string) (string, error)
	ApproveOrRejectPasswordRuleAction(ctx context.Context, actionID string, decision string, checker action.User, rejectionReason *string, department string) error
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

func (s *passwordRuleService) RequestPasswordRuleUpdate(ctx context.Context, rule *action.PasswordRule, maker action.User, department string) (string, error) {
	fmt.Println("rule", rule)
	if rule == nil || rule.ID == "" {
		return "", fmt.Errorf("INVALID_ID")
	}

	dbRule, err := s.repo.GetCurrentPasswordRule(ctx)
	if err != nil || dbRule == nil || dbRule.ID != rule.ID {
		return "", fmt.Errorf("NOT_FOUND")
	}

	pendingActions, err := s.repo.FetchPendingActionsByUniqueID(ctx, rule.ID)
	if err != nil {
		return "", fmt.Errorf("FAILED_TO_FETCH_PENDING_ACTIONS")
	}
	if len(pendingActions) > 0 {
		return "", fmt.Errorf("PENDING_ACTION_EXISTS")
	}

	prevAction := dbRule
	if dbRule == nil {
		prevAction = nil
	}

	cpsAction := action.CPSAction{
		ActionCode:       generateActionCode(),
		MakerID:          maker.UserID,
		MakerName:        maker.FullName,
		MakerPhoneNumber: maker.PhoneNumber,
		Department:       department,
		UniqueId:         rule.ID,
		ActionStatus:     action.ActionPending,
		ActionType:       action.ActionUpdate,
		RequestAction:    action.RequestUpdatePasswordRule,
		PreviosAction:    prevAction,
		CurrentAction:    rule,
		CreatedAt:        time.Now(),
		LastModifiedAt:   time.Now(),
	}

	created, err := s.repo.CreateCpsAction(ctx, cpsAction)
	if err != nil {
		return "", err
	}
	return created.ActionCode, nil
}

func generateActionCode() string {
	return fmt.Sprintf("CPS_%d", time.Now().UnixNano())
}

func (s *passwordRuleService) ApproveOrRejectPasswordRuleAction(ctx context.Context, actionID string, decision string, checker action.User, rejectionReason *string, department string) error {
	cpsAction, err := s.repo.GetPasswordRuleUpdateActionByID(ctx, actionID)
	if err != nil {
		return fmt.Errorf("ACTION_NOT_FOUND")
	}
	if cpsAction.ActionStatus != action.ActionPending {
		return fmt.Errorf("ACTION_NOT_PENDING")
	}
	cpsAction.CheckerID = checker.UserID
	cpsAction.CheckerName = checker.FullName
	cpsAction.CheckerPhoneNumber = checker.PhoneNumber
	cpsAction.Department = department
	cpsAction.CheckerActionTime = time.Now()
	cpsAction.LastModifiedAt = time.Now()

	if decision == "APPROVED" {
		// Unmarshal the domain PasswordRule from CurrentAction
		var rule action.PasswordRule
		currentActionBytes, _ := json.Marshal(cpsAction.CurrentAction)
		if err := json.Unmarshal(currentActionBytes, &rule); err != nil {
			return fmt.Errorf("INVALID_CURRENT_ACTION: %w", err)
		}
		if rule.ID == "" {
			return fmt.Errorf("ID_REQUIRED")
		}
		if err := s.repo.UpdatePasswordRule(ctx, rule); err != nil {
			return fmt.Errorf("FAILED_TO_UPDATE_PASSWORD_RULE: %w", err)
		}
		cpsAction.ActionStatus = action.ActionApproved
	} else if decision == "REJECTED" {
		cpsAction.ActionStatus = action.ActionRejected
		if rejectionReason != nil {
			cpsAction.RejectionReason = rejectionReason
		}
	} else {
		return fmt.Errorf("INVALID_DECISION")
	}
	return s.repo.UpdateCpsAction(ctx, *cpsAction)
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
