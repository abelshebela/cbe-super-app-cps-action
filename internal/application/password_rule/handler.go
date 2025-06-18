package passwordrule

import (
    "context"

    "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/action"
    "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/password_rule/services"
    "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type PasswordRuleHandler struct {
    service services.PasswordRuleService
    logger  utils.Logger
}

func InitPasswordRuleHandler(service services.PasswordRuleService, logger utils.Logger) *PasswordRuleHandler {
    return &PasswordRuleHandler{
        service: service,
        logger:  logger,
    }
}

func (h *PasswordRuleHandler) RequestPasswordRuleUpdate(ctx context.Context, rule *action.PasswordRule, maker action.User) (string, error) {
    actionID, err := h.service.RequestPasswordRuleUpdate(ctx, rule, maker)
    if err != nil {
        h.logger.Errorf("Failed to request password rule update: %v", err)
        return "", err
    }
    return actionID, nil
}

func (h *PasswordRuleHandler) ApproveOrRejectPasswordRuleAction(ctx context.Context, actionID string, approve bool, checker action.User, rejectionReason *string) error {
    err := h.service.ApproveOrRejectPasswordRuleAction(ctx, actionID, approve, checker, rejectionReason)
    if err != nil {
        h.logger.Errorf("Failed to approve/reject password rule action: %v", err)
        return err
    }
    return nil
}

func (h *PasswordRuleHandler) GetPasswordRuleUpdateActionByID(ctx context.Context, actionID string) (*action.CPSAction, error) {
    action, err := h.service.GetPasswordRuleUpdateActionByID(ctx, actionID)
    if err != nil {
        h.logger.Errorf("Failed to get password rule update action by ID: %v", err)
        return nil, err
    }
    return action, nil
}