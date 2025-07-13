package permission_handler

import (
	"encoding/json"
	// "fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/permission"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/permission/entities"
	inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/permission"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"

	util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	// "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	local_commen "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/common"
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
		util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	if err := request.Validate(); err != nil {
		h.logger.Warnf("[CreatePermissionGroup] validation failed: %v", err)
		util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	maker_ID, ok1 := r.Context().Value(constant.ContextKey("user_code")).(string)
	full_name, ok2 := r.Context().Value(constant.ContextKey("full_name")).(string)
	phone_number, ok3 := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department, ok4 := r.Context().Value(constant.ContextKey("department")).(string)

	if !ok1 || maker_ID == "" {
		util.SendErrorResponse(w, "missing or invalid user_code in context", http.StatusBadRequest, nil)
		return
	}
	if !ok2 || full_name == "" {
		util.SendErrorResponse(w, "missing or invalid full_name in context", http.StatusBadRequest, nil)
		return
	}
	if !ok3 || phone_number == "" {
		util.SendErrorResponse(w, "missing or invalid phone_number in context", http.StatusBadRequest, nil)
		return
	}
	if !ok4 || department == "" {
		util.SendErrorResponse(w, "missing or invalid department in context", http.StatusBadRequest, nil)
		return
	}

	cpsAction := model.CPSAction{

		MakerID:          maker_ID,
		MakerName:        full_name,
		MakerPhoneNumber: phone_number,
		Department:       department,
		ActionStatus:     string(model.ActionPending),
		ActionType:       string(entities.ActionCreate),
		RequestAction:    string(entities.RequestPermissionGroup),
	}

	cpsAction, err := h.permissionService.CreatePermissionGroup(request.GroupName, request.Role, request.PermissionCategoryLists, cpsAction)
	if err != nil {
		h.logger.Errorf("[CreatePermissionGroup] service error: %v", err)
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	h.logger.Infof("[CreatePermissionGroup] request sent successfully by user: %s", maker_ID)
	def, _ := local_commen.GetSuccessResponseByCode("SUCCESS ")
	data, _ := util.StructToMap(cpsAction)
	util.BaseResponseMaker(data, w, def.Message, http.StatusCreated)
}

func (h *PermissionHandler) ApprovePermissionGroup(w http.ResponseWriter, r *http.Request) {
	actionCode := chi.URLParam(r, "action_code")
	ctx := ctx_util.ExtractUserContext(r)

	h.logger.Infof("[ApprovePermissionGroup] approving action %s", actionCode)

	checker := model.CPSAction{
		CheckerID:          ctx.UserID,
		CheckerName:        ctx.FullName,
		CheckerPhoneNumber: ctx.PhoneNumber,
		Department:         ctx.Department,
	}

	err := h.permissionService.ApprovePermissionGroup(actionCode, checker)
	if err != nil {
		h.logger.Errorf("[ApprovePermissionGroup] service error: %v", err)
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	def, _ := local_commen.GetSuccessResponseByCode("SUCCESS")
	util.BaseResponseMaker(nil, w, def.Message, http.StatusCreated)
}
