package topup

import (
	topupDto "cbe-super-app-cps-action/internal/constants/dto/topup"
	topupInbound "cbe-super-app-cps-action/internal/constants/interfaces/topup"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"net/http"

	"cbe-super-app-cps-action/internal/constants/localization"
	topupcore "cbe-super-app-cps-action/internal/handlers/rest/http/topup/core"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type PaginatedTopupResponse types.PaginatedResponse[[]*model.Topup]

type topupAdapter struct {
	topupApp service.TopupService
	logger   utils.Logger
}

func InitTopupAdapter(topupApp service.TopupService, logger utils.Logger) topupInbound.TopupAdapter {
	return &topupAdapter{
		topupApp: topupApp,
		logger:   logger,
	}
}

// CreateTopup godoc
// @Summary Create a new topup
// @Description Create a new topup with the provided information
// @Tags Topup
// @Accept multipart/form-data
// @Produce json
// @Param data formData topupDto.TopupRequest false "Topup update data"
// @Param avatar formData file fale "Avatar image file"
// @Success 200 {object} localization.StandardResponse{data=nil} "Topup creation request sent successfully"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Bad request"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Internal server error"
// @Security BearerAuth
// @Router /topups [post]
func (a *topupAdapter) CreateTopup(w http.ResponseWriter, r *http.Request) {
	a.logger.Infof("called the create topup handler: %s", "create_method")

	var req topupDto.TopupRequest
	req, err := topupcore.ParseTopupRequestFromMultipartForm(r, true)
	if err != nil {
		a.logger.Errorf("error fetching topup create request data")
	}

	if err := req.Validate(true); err != nil {
		a.logger.Errorf("topup request validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := a.topupApp.CreateTopup(r.Context(), req); err != nil {
		a.logger.Errorf("failed to create topup: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	a.logger.Infof("topup creation request submitted successfully")
	localization.SendSuccessResponse(w, localization.SuccessTopupCreationRequestSent, nil)
}

// UpdateTopup godoc
// @Summary Update a topup
// @Description Update a topup with the provided information
// @Tags Topup
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Topup ID"
// @Param data formData topupDto.TopupRequest false "Topup update data"
// @Param avatar formData file fale "Avatar image file"
// @Success 200 {object} localization.StandardResponse{data=nil} "Topup update request sent successfully"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Bad request"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Internal server error"
// @Security BearerAuth
// @Router /topups/{id} [patch]
func (a *topupAdapter) UpdateTopup(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	a.logger.Infof("called the topup handler update for topup with id: %s", id)

	if id == "" {
		a.logger.Errorf("topup ID is required for update")
		localization.SendErrorResponse(w, localization.ErrorTopupIDRequired, nil, nil)
		return
	}

	req, err := topupcore.ParseTopupRequestFromMultipartForm(r, false)
	if err != nil {
		a.logger.Errorf("failed to parse topup update request: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	if err := req.Validate(false); err != nil {
		a.logger.Errorf("topup request validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if req.IsEmpty() {
		a.logger.Warnf("no data provided for topup update, topup ID: %s", id)
		localization.SendErrorResponse(w, localization.ErrorTopupUpdateEmptyPayload, nil, nil)
		return
	}

	if err := a.topupApp.UpdateTopup(r.Context(), id, req); err != nil {
		a.logger.Errorf("failed to update topup (ID: %s): %v", id, err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	a.logger.Infof("topup update request submitted successfully, topup ID: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessTopupUpdateRequestSent, nil)
}

// DeleteTopup godoc
// @Summary Delete a topup
// @Description Permanently delete a topup by ID
// @Tags Topup
// @Accept json
// @Produce json
// @Param id path string true "Topup ID"
// @Success 200 {object} localization.StandardResponse{data=nil} "Topup deleted successfully"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Bad request"
// @Failure 404 {object} localization.StandardResponse{data=nil} "Topup not found"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Internal server error"
// @Security BearerAuth
// @Router /topups/{id} [delete]
func (a *topupAdapter) DeleteTopup(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	a.logger.Infof("called the topup handler delete for topup with id: %s", id)

	if id == "" {
		localization.SendErrorResponse(w, localization.ErrorTopupIDRequired, nil, nil)
		return
	}

	if err := a.topupApp.DeleteTopup(r.Context(), id); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessTopupDeleted, nil)
}

// EnableTopup godoc
// @Summary Enable a topup
// @Description Enable a topup by ID
// @Tags Topup
// @Accept json
// @Produce json
// @Param id path string true "Topup ID"
// @Success 200 {object} localization.StandardResponse{data=nil} "Topup enable request submitted"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Bad request"
// @Failure 404 {object} localization.StandardResponse{data=nil} "Topup not found"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Internal server error"
// @Security BearerAuth
// @Router /topups/{id}/enable [patch]
func (a *topupAdapter) Enable(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	a.logger.Infof("called the topup handler enable for topup with id: %s", id)

	if id == "" {
		localization.SendErrorResponse(w, localization.ErrorTopupIDRequired, nil, nil)
		return
	}

	if err := a.topupApp.EnableOrDisableTopup(r.Context(), id, true); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessTopupEnableRequestSubmitted, nil)
}

// DisableTopup godoc
// @Summary Disable a topup
// @Description Disable a topup by ID
// @Tags Topup
// @Accept json
// @Produce json
// @Param id path string true "Topup ID"
// @Success 200 {object} localization.StandardResponse{data=nil} "Topup disable request submitted"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Bad request"
// @Failure 404 {object} localization.StandardResponse{data=nil} "Topup not found"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Internal server error"
// @Security BearerAuth
// @Router /topups/{id}/disable [patch]
func (a *topupAdapter) Disable(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	a.logger.Infof("called the topup handler disable for topup with id: %s", id)

	if id == "" {
		localization.SendErrorResponse(w, localization.ErrorTopupIDRequired, nil, nil)
		return
	}

	if err := a.topupApp.EnableOrDisableTopup(r.Context(), id, false); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessTopupDisableRequestSubmitted, nil)
}

// GetTopup godoc
// @Summary Get topup by ID
// @Description Retrieve a topup's details by ID
// @Tags Topup
// @Accept json
// @Produce json
// @Param id path string true "Topup ID"
// @Success 200 {object} localization.StandardResponse{data=model.Topup} "Topup retrieved successfully"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Bad request"
// @Failure 404 {object} localization.StandardResponse{data=nil} "Topup not found"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Internal server error"
// @Security BearerAuth
// @Router /topups/{id} [get]
func (a *topupAdapter) GetTopup(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	a.logger.Infof("called the fetch topup handler  for topup with id: %s", id)

	if id == "" {
		localization.SendErrorResponse(w, localization.ErrorTopupIDRequired, nil, nil)
		return
	}

	topup, err := a.topupApp.GetTopup(r.Context(), id)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessTopupRetrieved, topup)
}

// GetTopups godoc
// @Summary List topups
// @Description Retrieve topups with pagination and optional search
// @Tags Topup
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param search query string false "Search term"
// @Success 200 {object} localization.StandardResponse{data=PaginatedTopupResponse} "Topups retrieved successfully"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Internal server error"
// @Security BearerAuth
// @Router /topups [get]
func (a *topupAdapter) GetAllTopup(w http.ResponseWriter, r *http.Request) {
	filter := local_util.ExtractFilterParams(r)
	a.logger.Infof("fetching topups with filter: %+v", filter)

	list, err := a.topupApp.GetAllTopup(r.Context(), *filter)
	if err != nil {
		a.logger.Errorf("failed to fetch topups: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	a.logger.Infof("topups fetched successfully")
	localization.SendSuccessResponse(w, localization.SuccessTopupsRetrieved, list)
}
