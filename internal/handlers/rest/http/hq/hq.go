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

func (a *hqAdapter) GetAllHQ(w http.ResponseWriter, r *http.Request) {
	filter := local_util.ExtractFilterParams(r)
	list, err := a.hqApp.GetHQDetail(r.Context(), *filter)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessHQsFetched, list)
}

func (a *hqAdapter) GetBlockTime(w http.ResponseWriter, r *http.Request) {
	resp, err := a.hqApp.GetBlockTime(r.Context())
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessHQBlockTimeFetched, resp)
}

func (a *hqAdapter) GetArchiveTime(w http.ResponseWriter, r *http.Request) {
	resp, err := a.hqApp.GetArchiveTime(r.Context())
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessHQArchiveTimeFetched, resp)
}

func (a *hqAdapter) GetPasswordExpiry(w http.ResponseWriter, r *http.Request) {
	resp, err := a.hqApp.GetPasswordExpiry(r.Context())
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessHQPasswordExpiryFetched, resp)
}

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
