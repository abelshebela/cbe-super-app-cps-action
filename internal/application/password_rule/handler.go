package passwordrule

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/password_rule/services"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type PasswordRuleHandler struct {
	services.PasswordRuleService
	logger utils.Logger
}

func InitPasswordRuleHandler(service services.PasswordRuleService, logger utils.Logger) *PasswordRuleHandler {
	return &PasswordRuleHandler{
		PasswordRuleService: service,
		logger:              logger,
	}
}

func (h *PasswordRuleHandler) RequestPasswordRuleUpdate(ctx context.Context, rule *action.PasswordRule, maker action.User, department string) (string, error) {
	
	actionID, err := h.PasswordRuleService.RequestPasswordRuleUpdate(ctx, rule, maker, department)
	if err != nil {
		h.logger.Errorf("Failed to request password rule update: %v", err)
		return "", err
	}
	return actionID, nil
}


func (h *PasswordRuleHandler) GetPasswordRuleUpdateActionByID(ctx context.Context, actionID string) (*action.CPSAction, error) {
	action, err := h.PasswordRuleService.GetPasswordRuleUpdateActionByID(ctx, actionID)
	if err != nil {
		h.logger.Errorf("Failed to get password rule update action by ID: %v", err)
		return nil, err
	}
	return action, nil
}

func (h *PasswordRuleHandler) GetUpdateAction(ctx context.Context, maker action.User) (*action.CPSAction, error) {
	action, err := h.PasswordRuleService.GetUpdateAction(ctx, maker)
	if err != nil {
		h.logger.Errorf("Failed to get update action: %v", err)
		return nil, err
	}
	return action, nil
}
