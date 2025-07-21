package miniapphandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	miniapp_application "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/mini_app"
	Inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/miniapp"
	contexts "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	util "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"github.com/go-chi/chi/v5"
)

type HttpStore struct {
	Application miniapp_application.ApplicationAbstracts
	logger      util.Logger
}

func NewMiniAppAdapter(app miniapp_application.ApplicationAbstracts, logger util.Logger) Inbound.MiniAppInbound {
	return &HttpStore{
		Application: app,
		logger:      logger,
	}
}

func (h *HttpStore) createUser(r *http.Request) (*model.User, error) {
	ctx_extract := ctx_util.ExtractUserContext(r)
	if ctx_extract.IsIncomplete() {
		return nil, fmt.Errorf("Unauthrozed")
	}
	maker := &model.User{
		UserCode:    ctx_extract.UserID,
		FullName:    ctx_extract.FullName,
		PhoneNumber: ctx_extract.PhoneNumber,
		Department:  ctx_extract.Department,
	}

	return maker, nil
}

func (h *HttpStore) MakerCreateMiniApp(w http.ResponseWriter, r *http.Request) {
	var req dto.MiniAppCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	makerUser := contexts.ExtractUserContext(r)

	if makerUser.IsIncomplete() {
		utils.SendErrorResponse(w, utils.IncompleteUserInfo, 0, nil)
		return
	}

	var maker model.User

	maker.FullName = makerUser.FullName
	maker.UserCode = makerUser.UserCode
	maker.PhoneNumber = makerUser.PhoneNumber
	maker.Department = makerUser.Department

	response, err := h.Application.MakerCreateMiniApp(r.Context(), req, maker)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusNoContent, "Failed to create mini app")
		return
	}

	utils.WriteSuccessResponse(w, response, "mini App request successfully created")
}

func (h *HttpStore) MakerUpdateMiniApp(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		h.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}
	var req dto.MiniAppCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	req.ID = id

	maker, UserErr := h.createUser(r)
	if UserErr != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	response, err := h.Application.MakerUpdateMiniApp(r.Context(), req, *maker)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusNoContent, "Failed to update mini app")
		return
	}
	utils.WriteSuccessResponse(w, response, "successful")
}

func (h *HttpStore) MakerDeleteMiniApp(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	maker, UserErr := h.createUser(r)
	if UserErr != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	ActionCode, err := h.Application.MakerDeleteMiniApp(r.Context(), maker, id)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusNoContent, "Failed to delete mini app")
		return
	}
	utils.WriteSuccessResponse(w, ActionCode, "successful")
}

func (h *HttpStore) ListMiniApp(w http.ResponseWriter, r *http.Request) {
	list, err := h.Application.ListMiniApp(r.Context())
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusNoContent, "Failed to list mini apps")
		return
	}
	utils.WriteSuccessResponse(w, list, "successful")
}

func (h *HttpStore) DetailMiniAppByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Missing mini app ID")
		return
	}
	detail, err := h.Application.DetailMiniAppByID(r.Context(), id)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to get mini app detail")
		return
	}
	utils.WriteSuccessResponse(w, detail, "successful feached miniapp by ID")
}
