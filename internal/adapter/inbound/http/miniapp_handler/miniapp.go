package miniapphandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
)

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
}
