package permission_handler

import (
	"encoding/json"
	"net/http"
	"github.com/go-chi/chi/v5"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/application/permission"
	inbound "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/port/inbound/permission"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/application/middleware"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/domain/permission/entities"

)

type PermissionHandler struct {
	permissionService permission.PermissionService
	logger            utils.Logger
}

func NewPermissionHTTPHandler(service permission.PermissionService, logger utils.Logger) inbound.PermissionPortHandler {
	return &PermissionHandler{
		permissionService: service,
		logger:            logger,
	}
}

func (h *PermissionHandler) CreatePermissionGroup(w http.ResponseWriter, r *http.Request) {
	var request permission.CreatePermissionGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Errorf("[CreatePermissionGroup] failed to decode request: %v", err)
		resp := common.Response[any]{
			ResponseWriter: w,
			Status:         http.StatusBadRequest,
			Data: map[string]string{"message":  "Invalid JSON payload"},
		}
		resp.SendJSON()
		return
	}

	if err := request.Validate(); err != nil {
		h.logger.Warnf("[CreatePermissionGroup] validation failed: %v", err)
		resp := common.Response[any]{
			ResponseWriter: w,
			Status:         http.StatusBadRequest,
			Data: map[string]string{"message": "Invalid input provided"},
		}
		resp.SendJSON()
		return
	}

	userCtx := r.Context().Value(middleware.ContextKey("user_payload"))
	user, ok := userCtx.(middleware.UserPayload)
	if !ok {
		h.logger.Errorf("[CreatePermissionGroup] failed to extract user from context")
		resp := common.Response[any]{
			ResponseWriter: w,
			Status:         http.StatusUnauthorized,
			Data: map[string]string{"message": " unauthorized user"},
		}
		resp.SendJSON()
		return
	}

	cpsAction := entities.CPSAction{
		MakerID:            user.UserID,
		MakerName:          user.FullName,
		MakerPhoneNumber:   user.PhoneNumber,
		Department:         user.Department,
		ActionStatus:       entities.ActionPending,
		ActionType:         entities.ActionCreate,
		RequestAction:      entities.RequestPermissionGroup,
	}

	if err := h.permissionService.CreatePermissionGroup(request.GroupName, request.Role, request.PermissionCategoryLists, cpsAction); err != nil {
		h.logger.Errorf("[CreatePermissionGroup] service error: %v", err)
		resp := common.Response[any]{
			ResponseWriter: w,
			Status:         http.StatusBadRequest,
			Data: map[string]string{
			"message": err.Error(),
		},
		}
		resp.SendJSON()
		return
	}

	h.logger.Infof("[CreatePermissionGroup] request sent successfully by user: %s", user.UserID)
	resp := common.Response[any]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data: map[string]string{
			"message": "request sent successfully",
		},
	}
	resp.SendJSON()
}

func (h *PermissionHandler) ApprovePermissionGroup(w http.ResponseWriter, r *http.Request) {
	actionCode := chi.URLParam(r, "action_code")
	userCtx := r.Context().Value(middleware.ContextKey("user_payload"))
	user, ok := userCtx.(middleware.UserPayload)
	if !ok {
		h.logger.Errorf("[ApprovePermissionGroup] failed to extract user from context")
		resp := common.Response[any]{
			ResponseWriter: w,
			Status:         http.StatusUnauthorized,
			Data:           map[string]string{"message": "Unauthorized"},
		}
		resp.SendJSON()
		return
	}

	h.logger.Infof("[ApprovePermissionGroup] approving action %s", actionCode)

	checker := entities.CPSAction{
		CheckerID:   user.UserID,
		CheckerName: user.FullName,
		CheckerPhoneNumber: user.PhoneNumber,
		Department:  user.Department,
	}

	err := h.permissionService.ApprovePermissionGroup(actionCode, checker)
	if err != nil {
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
		Status:         http.StatusOK,
		Data:           map[string]string{"message": "Action approved"},
	}
	resp.SendJSON()
}
