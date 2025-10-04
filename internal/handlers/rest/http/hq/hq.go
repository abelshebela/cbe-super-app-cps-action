package hqhandler

import (
	"encoding/json"
	"net/http"

	hqDto "cbe-super-app-cps-action/internal/constants/dto/hq"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type hqAdapter struct {
	hqApp  service.HQService
	logger utils.Logger
}

func InitHQAdapter(hqApp service.HQService, logger utils.Logger) *hqAdapter {
	return &hqAdapter{
		hqApp:  hqApp,
		logger: logger,
	}
}

// GetHQ godoc
// @Summary Get HQ by ID
// @Description Fetch HQ details by ID
// @Tags HQ
// @Accept json
// @Produce json
// @Param id path string true "HQ ID"
// @Success 200 {object} localization.ResponseCode "HQ fetched successfully"
// @Failure 400 {object} localization.ResponseCode "HQ ID required"
// @Failure 404 {object} localization.ResponseCode "HQ not found"
// @Failure 500 {object} localization.ResponseCode "Internal server error"
// @Security BearerAuth
// @Router /hq/{id} [get]
func (a *hqAdapter) GetHQ(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorResponse(w, localization.ErrorHQIDRequired, nil, nil)
		return
	}

	hqResp, err := a.hqApp.GetHQ(r.Context(), id)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessHQFetched, hqResp)
}

// GetAllHQ godoc
// @Summary Get all HQs
// @Description Fetch a list of HQs with pagination and filters
// @Tags HQ
// @Accept json
// @Produce json
// @Param page query int false "Page number (default 1)"
// @Param per_page query int false "Items per page (default 10, max 100)"
// @Param sort query string false "Sort field"
// @Param order query string false "Sort order (asc/desc)"
// @Success 200 {object} localization.ResponseCode "HQs fetched successfully"
// @Failure 500 {object} localization.ResponseCode "Internal server error"
// @Security BearerAuth
// @Router /hq [get]
func (a *hqAdapter) GetAllHQ(w http.ResponseWriter, r *http.Request) {
	filter := local_util.ExtractFilterParams(r)
	list, err := a.hqApp.GetHQDetail(r.Context(), *filter)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessHQsFetched, list)
}

// GetBlockTime godoc
// @Summary Get HQ block time
// @Description Retrieve current HQ block time configuration
// @Tags HQ
// @Accept json
// @Produce json
// @Success 200 {object} localization.ResponseCode "HQ block time fetched successfully"
// @Failure 500 {object} localization.ResponseCode "Internal server error"
// @Security BearerAuth
// @Router /hq/block_time [get]
func (a *hqAdapter) GetBlockTime(w http.ResponseWriter, r *http.Request) {
	resp, err := a.hqApp.GetBlockTime(r.Context())
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessHQBlockTimeFetched, resp)
}

// GetArchiveTime godoc
// @Summary Get HQ archive time
// @Description Retrieve current HQ archive time configuration
// @Tags HQ
// @Accept json
// @Produce json
// @Success 200 {object} localization.ResponseCode "HQ archive time fetched successfully"
// @Failure 500 {object} localization.ResponseCode "Internal server error"
// @Security BearerAuth
// @Router /hq/archive_time [get]
func (a *hqAdapter) GetArchiveTime(w http.ResponseWriter, r *http.Request) {
	resp, err := a.hqApp.GetArchiveTime(r.Context())
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessHQArchiveTimeFetched, resp)
}

// GetPasswordExpiry godoc
// @Summary Get HQ password expiry
// @Description Retrieve current HQ password expiry configuration
// @Tags HQ
// @Accept json
// @Produce json
// @Success 200 {object} localization.ResponseCode "HQ password expiry fetched successfully"
// @Failure 500 {object} localization.ResponseCode "Internal server error"
// @Security BearerAuth
// @Router /hq/password_expiry [get]
func (a *hqAdapter) GetPasswordExpiry(w http.ResponseWriter, r *http.Request) {
	resp, err := a.hqApp.GetPasswordExpiry(r.Context())
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessHQPasswordExpiryFetched, resp)
}

// UpdateBlockTimeRequest godoc
// @Summary Update HQ block time
// @Description Submit a request to update HQ block time
// @Tags HQ
// @Accept json
// @Produce json
// @Param request body hqDto.UpdateBlockTimeRequest true "Update block time request payload"
// @Success 200 {object} localization.ResponseCode "HQ block time update request submitted"
// @Failure 400 {object} localization.ResponseCode "Invalid request payload"
// @Failure 500 {object} localization.ResponseCode "Internal server error"
// @Security BearerAuth
// @Router /hq/block_time [post]
func (a *hqAdapter) UpdateBlockTimeRequest(w http.ResponseWriter, r *http.Request) {
	var request hqDto.UpdateBlockTimeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidHQRequest.Code)
		return
	}

	if err := request.Validate(); err != nil {
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	err := a.hqApp.UpdateBlockTime(r.Context(), request)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessHQBlockTimeUpdateRequestSubmitted, nil)
}

// UpdateArchiveTimeRequest godoc
// @Summary Update HQ archive time
// @Description Submit a request to update HQ archive time
// @Tags HQ
// @Accept json
// @Produce json
// @Param request body hqDto.UpdateArchiveTimeRequest true "Update archive time request payload"
// @Success 200 {object} localization.ResponseCode "HQ archive time update request submitted"
// @Failure 400 {object} localization.ResponseCode "Invalid request payload"
// @Failure 500 {object} localization.ResponseCode "Internal server error"
// @Security BearerAuth
// @Router /hq/archive_time [post]
func (a *hqAdapter) UpdateArchiveTimeRequest(w http.ResponseWriter, r *http.Request) {
	var request hqDto.UpdateArchiveTimeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidHQRequest.Code)
		return
	}

	if err := request.Validate(); err != nil {
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	err := a.hqApp.UpdateArchiveTime(r.Context(), request)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessHQArchiveTimeUpdateRequestSubmitted, nil)
}

// UpdatePasswordExpiryRequest godoc
// @Summary Update HQ password expiry
// @Description Submit a request to update HQ password expiry
// @Tags HQ
// @Accept json
// @Produce json
// @Param request body hqDto.UpdatePasswordExpiryRequest true "Update password expiry request payload"
// @Success 200 {object} localization.ResponseCode "HQ password expiry update request submitted"
// @Failure 400 {object} localization.ResponseCode "Invalid request payload"
// @Failure 500 {object} localization.ResponseCode "Internal server error"
// @Security BearerAuth
// @Router /hq/password_expiry [post]
func (a *hqAdapter) UpdatePasswordExpiryRequest(w http.ResponseWriter, r *http.Request) {
	var request hqDto.UpdatePasswordExpiryRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		localization.SendBadRequestResponse(w, localization.ErrorInvalidHQRequest.Code)
		return
	}

	if err := request.Validate(); err != nil {
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	err := a.hqApp.UpdatePasswordExpiry(r.Context(), request)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessHQPasswordExpiryUpdateRequestSubmitted, nil)
}
