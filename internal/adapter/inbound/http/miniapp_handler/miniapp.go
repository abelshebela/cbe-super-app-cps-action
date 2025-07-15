package miniapphandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	// miniAppApplication "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/mini_app"
	// miniapp_application "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/mini_app"
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
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	maker, UserErr := h.createUser(w, r)
	if UserErr != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Process the request using application logic
	response, err := h.Application.MakerCreateMiniApp(r.Context(), req, *maker)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusNoContent, "Failed to create mini app")
		return
	}
	utils.WriteSuccessResponse(w, response, "successful")
}

func (h *HttpStore) CheckerMiniApp(w http.ResponseWriter, r *http.Request) {
	var req dto.MiniAppCheckerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	checker, UserErr := h.createUser(w, r)
	if UserErr != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	err := h.Application.CheckerCreateMiniApp(r.Context(), req.Action_id, req.Action, *checker)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusNoContent, "")
		return
	}
	utils.WriteSuccessResponse(w, nil, "successful")
	return
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
	// var req dto.MiniAppDeleteRequest
	// if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	// 	utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
	// 	return
	// }
	// defer r.Body.Close()

	maker, UserErr := h.createUser(w, r)
	if UserErr != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	ActionCode, err := h.Application.MakerDeleteMiniApp(r.Context(), maker)
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
	id := r.URL.Query().Get("id")
	if id == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Missing mini app ID")
		return
	}
	detail, err := h.Application.DetailMiniAppByID(r.Context(), id)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusNoContent, "Failed to get mini app detail")
		return
	}
	utils.WriteSuccessResponse(w, detail, "successful")
}
