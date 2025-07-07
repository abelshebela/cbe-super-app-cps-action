package permission_handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/permission"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/permission/entities"
	inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/permission"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
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
			Data:           map[string]string{"message": "Invalid JSON payload"},
		}
		resp.SendJSON()
		return
	}

	if err := request.Validate(); err != nil {
		h.logger.Warnf("[CreatePermissionGroup] validation failed: %v", err)
		resp := common.Response[any]{
			ResponseWriter: w,
			Status:         http.StatusBadRequest,
			Data:           map[string]string{"message": "Invalid input provided"},
		}
		resp.SendJSON()
		return
	}

	ctx := ctx_util.ExtractUserContext(r)

	cpsAction := entities.CPSAction{
		MakerID:          ctx.UserID,
		MakerName:        ctx.FullName,
		MakerPhoneNumber: ctx.PhoneNumber,
		Department:       ctx.Department,
		ActionStatus:     entities.ActionPending,
		ActionType:       entities.ActionCreate,
		RequestAction:    entities.RequestPermissionGroup,
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

	h.logger.Infof("[CreatePermissionGroup] request sent successfully by user: %s", ctx.UserID)
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
	ctx := ctx_util.ExtractUserContext(r)

	h.logger.Infof("[ApprovePermissionGroup] approving action %s", actionCode)

	checker := entities.CPSAction{
		CheckerID:          ctx.UserID,
		CheckerName:        ctx.FullName,
		CheckerPhoneNumber: ctx.PhoneNumber,
		Department:         ctx.Department,
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
