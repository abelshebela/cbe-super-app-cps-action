package miniapphandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	"github.com/go-chi/chi/v5"

	// miniAppApplication "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/mini_app"
	// miniapp_application "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/mini_app"
	contexts "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	// util "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	// inboundMiniApp "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/miniapp"
	// miniApp_application "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/mini_app"
	// util "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// type miniAppAdapter struct {
// 	miniAppHandler miniApp_application.ApplicationAbstracts
// 	logger         util.Logger
// }

// func InitMiniAppHandler(miniAppHandler miniAppApplication.ApplicationAbstracts, logger util.Logger) inboundMiniApp.MiniAppInbound {
// 	return miniAppAdapter{
// 		miniAppHandler: miniAppHandler,
// 		logger:         logger,
// 	}

// }

func (h *HttpStore) createUser(w http.ResponseWriter, r *http.Request) (*model.User, error) {
	ctx_extract := ctx_util.ExtractUserContext(r)
	if ctx_extract.IsIncomplete() {
		return nil, fmt.Errorf("Unauthrozed")
	}
	maker := &model.User{
		UserCode:    ctx_extract.UserID,
		FullName:    ctx_extract.FullName,
		PhoneNumber: ctx_extract.PhoneNumber,
	}

	return maker, nil
}

func (h *HttpStore) MakerCreateMiniApp(w http.ResponseWriter, r *http.Request) {
	var req dto.MiniAppCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println("handelr decode")
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	makerUser := contexts.ExtractUserContext(r)

	var maker model.User

	maker.FullName = makerUser.FullName
	maker.UserCode = makerUser.UserCode
	maker.PhoneNumber = makerUser.PhoneNumber
	Department := makerUser.Department

	response, err := h.Application.MakerCreateMiniApp(r.Context(), req, maker, Department)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusNoContent, "Failed to create mini app")
		return
	}
	utils.WriteSuccessResponse(w, response, "mini App request successfully created")
}

func (h *HttpStore) CheckerMiniApp(w http.ResponseWriter, r *http.Request) {
	var req dto.MiniAppCheckerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	CheckerUser := contexts.ExtractUserContext(r)

	var Checker model.User

	Checker.FullName = CheckerUser.FullName
	Checker.UserCode = CheckerUser.UserCode
	Checker.PhoneNumber = CheckerUser.PhoneNumber
	Department := CheckerUser.Department

	ApprovedAction, err := h.Application.CheckerCreateMiniApp(r.Context(), req.Action_id, req.Action, Checker, Department)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "")
		return
	}
	utils.WriteSuccessResponse(w, ApprovedAction, "successful")

}

func (h *HttpStore) MakerUpdateMiniApp(w http.ResponseWriter, r *http.Request) {
	var req dto.MiniAppCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	maker, UserErr := h.createUser(w, r)
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
	maker, UserErr := h.createUser(w, r)
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
	// id := r.URL.Query().Get("id")
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
