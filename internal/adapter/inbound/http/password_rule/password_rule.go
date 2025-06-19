package passwordrule

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/middleware"
	passwordrule "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/password_rule"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/action"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/password_rule/services"
	constant "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

)

type PasswordRuleHTTPHandler struct {
	service services.PasswordRuleService
	logger  utils.Logger
}

func NewPasswordRuleHTTPHandler(service services.PasswordRuleService, logger utils.Logger) *PasswordRuleHTTPHandler {
	return &PasswordRuleHTTPHandler{
		service: service,
		logger:  logger,
	}
}

func (h *PasswordRuleHTTPHandler) RequestPasswordRuleUpdate(w http.ResponseWriter, r *http.Request) {
    var req passwordrule.RequestPasswordRuleUpdateDTO
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.logger.Errorf("failed to bind password rule update data: %v", err)
        resp := common.Response[any]{
            ResponseWriter: w,
            Status:         http.StatusBadRequest,
            Data:           map[string]string{"message": "invalid request"},
        }
        resp.SendJSON()
        return
    }

    userID, _ := r.Context().Value(constant.ContextKey("user_id")).(string)
    fullName, _ := r.Context().Value(constant.ContextKey("full_name")).(string)
    phoneNumber, _ := r.Context().Value(constant.ContextKey("phone_number")).(string)

    if strings.TrimSpace(userID) == "" || strings.TrimSpace(fullName) == "" || strings.TrimSpace(phoneNumber) == "" {
        resp := common.Response[any]{
            ResponseWriter: w,
            Status:         http.StatusUnauthorized,
            Data:           map[string]string{"message": "user info missing in context"},
        }
        resp.SendJSON()
        return
    }

    maker := action.User{
        UserID:      userID,
        FullName:    fullName,
        PhoneNumber: phoneNumber,
    }

    ctx := r.Context()
    _, err := h.service.RequestPasswordRuleUpdate(ctx, &req.Rule, maker)
    if err != nil {
        h.logger.Errorf("RequestPasswordRuleUpdate failed: %v", err)
        resp := common.Response[any]{
            ResponseWriter: w,
            Status:         http.StatusInternalServerError,
            Data:           map[string]string{"message": err.Error()},
        }
        resp.SendJSON()
        return
    }

    resp := common.Response[any]{
        ResponseWriter: w,
        Status:         http.StatusCreated,
        Data:           map[string]string{"message": "Action created successfully"},
    }
    resp.SendJSON()
}
func (h *PasswordRuleHTTPHandler) ApproveOrRejectPasswordRuleAction(w http.ResponseWriter, r *http.Request) {
	var req passwordrule.ApproveOrRejectPasswordRuleActionDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.ActionID) == "" {
		err := fmt.Errorf("action_id is required %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "action_id is required",
		})
		middleware.ErrorHandler(w, err)
		return
	}

	userID, _ := r.Context().Value(constant.ContextKey("user_id")).(string)
	fullName, _ := r.Context().Value(constant.ContextKey("full_name")).(string)
	phoneNumber, _ := r.Context().Value(constant.ContextKey("phone_number")).(string)

	if strings.TrimSpace(userID) == "" || strings.TrimSpace(fullName) == "" || strings.TrimSpace(phoneNumber) == "" {
		err := fmt.Errorf("user info missing in context %w", constant.ErrorDefinition{
			Code:    http.StatusUnauthorized,
			Message: "user info missing in context",
		})
		middleware.ErrorHandler(w, err)
		return
	}

	checker := action.User{
		UserID:      userID,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}

	ctx := r.Context()
	err := h.service.ApproveOrRejectPasswordRuleAction(ctx, req.ActionID, req.Approve, checker, req.RejectionReason)
	if err != nil {
		h.logger.Errorf("ApproveOrRejectPasswordRuleAction failed: %v", err)
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[string]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           "Action processed successfully",
	}
	res.SendJSON()
}

func (h *PasswordRuleHTTPHandler) GetPasswordRuleUpdateActionByID(w http.ResponseWriter, r *http.Request) {
    actionID := r.URL.Query().Get("action_id")
    if actionID == "" {
        err := fmt.Errorf("action_id is required %w", constant.ErrorDefinition{
            Code:    http.StatusBadRequest,
            Message: "action_id is required",
        })
        middleware.ErrorHandler(w, err)
        return
    }

    userRole, _ := r.Context().Value(constant.ContextKey("user_role")).(string)
    if userRole != "CHECKER" && userRole != "IFB-CHECKER" {
        err := fmt.Errorf("forbidden: only checker can access this resource %w", constant.ErrorDefinition{
            Code:    http.StatusForbidden,
            Message: "forbidden: only checker can access this resource",
        })
        middleware.ErrorHandler(w, err)
        return
    }

    ctx := r.Context()
    result, err := h.service.GetPasswordRuleUpdateActionByID(ctx, actionID)
    if err != nil {
        h.logger.Errorf("GetPasswordRuleUpdateActionByID failed: %v", err)
        middleware.ErrorHandler(w, err)
        return
    }
    res := common.Response[*action.CPSAction]{
        ResponseWriter: w,
        Status:         http.StatusOK,
        Data:           result,
    }
    res.SendJSON()
}