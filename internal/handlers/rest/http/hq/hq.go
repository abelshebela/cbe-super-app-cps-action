package hqhandler

import (
	"encoding/json"
	"net/http"

	hqDto "cbe-super-app-cps-action/internal/constants/dto/hq"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"go.opentelemetry.io/otel/attribute"

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
//
//	@Summary		Get HQ by ID
//	@Description	Fetch HQ details by ID
//	@Tags			HQ
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string						true	"HQ ID"
//	@Success		200	{object}	localization.ResponseCode	"HQ fetched successfully"
//	@Failure		400	{object}	localization.ResponseCode	"HQ ID required"
//	@Failure		404	{object}	localization.ResponseCode	"HQ not found"
//	@Failure		500	{object}	localization.ResponseCode	"Internal server error"
//	@Security		BearerAuth
//	@Router			/hq/{id} [get]
func (a *hqAdapter) GetHQ(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "", "getHq", "handler", "hq")
	defer span.End()
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorResponse(w, localization.ErrorHQIDRequired, nil, nil)
		return
	}

	span.SetAttributes(attribute.String("hq.id", id))

	hqResp, err := a.hqApp.GetHQ(ctx, id)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessHQFetched, hqResp)
}

// GetAllHQ godoc
//
//	@Summary		Get all HQs
//	@Description	Fetch a list of HQs with pagination and filters
//	@Tags			HQ
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int							false	"Page number (default 1)"
//	@Param			per_page	query		int							false	"Items per page (default 10, max 100)"
//	@Param			sort		query		string						false	"Sort field"
//	@Param			order		query		string						false	"Sort order (asc/desc)"
//	@Success		200			{object}	localization.ResponseCode	"HQs fetched successfully"
//	@Failure		500			{object}	localization.ResponseCode	"Internal server error"
//	@Security		BearerAuth
//	@Router			/hq [get]
func (a *hqAdapter) GetAllHQ(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "", "getAllHq", "handler", "hq")
	defer span.End()
	filter := local_util.ExtractFilterParams(r)
	list, err := a.hqApp.GetHQDetail(ctx, *filter)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.Int("hq.count", len(list.Data)))
	localization.SendSuccessResponse(w, localization.SuccessHQsFetched, list)
}

// GetBlockTime godoc
//
//	@Summary		Get HQ block time
//	@Description	Retrieve current HQ block time configuration
//	@Tags			HQ
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	localization.ResponseCode	"HQ block time fetched successfully"
//	@Failure		500	{object}	localization.ResponseCode	"Internal server error"
//	@Security		BearerAuth
//	@Router			/hq/block_time [get]
func (a *hqAdapter) GetBlockTime(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "", "getHqBlockTime", "handler", "hq")
	defer span.End()

	resp, err := a.hqApp.GetBlockTime(ctx)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessHQBlockTimeFetched, resp)
}

// GetArchiveTime godoc
//
//	@Summary		Get HQ archive time
//	@Description	Retrieve current HQ archive time configuration
//	@Tags			HQ
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	localization.ResponseCode	"HQ archive time fetched successfully"
//	@Failure		500	{object}	localization.ResponseCode	"Internal server error"
//	@Security		BearerAuth
//	@Router			/hq/archive_time [get]
func (a *hqAdapter) GetArchiveTime(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "", "getHqArchiveTime", "handler", "hq")
	defer span.End()

	resp, err := a.hqApp.GetArchiveTime(ctx)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessHQArchiveTimeFetched, resp)
}

// GetPasswordExpiry godoc
//
//	@Summary		Get HQ password expiry
//	@Description	Retrieve current HQ password expiry configuration
//	@Tags			HQ
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	localization.ResponseCode	"HQ password expiry fetched successfully"
//	@Failure		500	{object}	localization.ResponseCode	"Internal server error"
//	@Security		BearerAuth
//	@Router			/hq/password_expiry [get]
func (a *hqAdapter) GetPasswordExpiry(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "", "getHqPasswordExpiry", "handler", "hq")
	defer span.End()

	resp, err := a.hqApp.GetPasswordExpiry(ctx)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessHQPasswordExpiryFetched, resp)
}

// UpdateBlockTimeRequest godoc
//
//	@Summary		Update HQ block time
//	@Description	Submit a request to update HQ block time
//	@Tags			HQ
//	@Accept			json
//	@Produce		json
//	@Param			request	body		hqDto.UpdateBlockTimeRequest	true	"Update block time request payload"
//	@Success		200		{object}	localization.ResponseCode		"HQ block time update request submitted"
//	@Failure		400		{object}	localization.ResponseCode		"Invalid request payload"
//	@Failure		500		{object}	localization.ResponseCode		"Internal server error"
//	@Security		BearerAuth
//	@Router			/hq/block_time [post]
func (a *hqAdapter) UpdateBlockTimeRequest(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "", "updateHqBlockTime", "handler", "hq")
	defer span.End()
	var request hqDto.UpdateBlockTimeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		span.RecordError(err)
		localization.SendBadRequestResponse(w, localization.ErrorInvalidHQRequest.Code)
		return
	}

	if err := request.Validate(); err != nil {
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := a.hqApp.UpdateBlockTime(ctx, request); err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessHQBlockTimeUpdateRequestSubmitted, nil)
}

// UpdateArchiveTimeRequest godoc
//
//	@Summary		Update HQ archive time
//	@Description	Submit a request to update HQ archive time
//	@Tags			HQ
//	@Accept			json
//	@Produce		json
//	@Param			request	body		hqDto.UpdateArchiveTimeRequest	true	"Update archive time request payload"
//	@Success		200		{object}	localization.ResponseCode		"HQ archive time update request submitted"
//	@Failure		400		{object}	localization.ResponseCode		"Invalid request payload"
//	@Failure		500		{object}	localization.ResponseCode		"Internal server error"
//	@Security		BearerAuth
//	@Router			/hq/archive_time [post]
func (a *hqAdapter) UpdateArchiveTimeRequest(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "", "updateHqArchiveTime", "handler", "hq")
	defer span.End()
	var request hqDto.UpdateArchiveTimeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		span.RecordError(err)
		localization.SendBadRequestResponse(w, localization.ErrorInvalidHQRequest.Code)
		return
	}

	if err := request.Validate(); err != nil {
		span.RecordError(err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := a.hqApp.UpdateArchiveTime(ctx, request); err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessHQArchiveTimeUpdateRequestSubmitted, nil)
}

// UpdatePasswordExpiryRequest godoc
//
//	@Summary		Update HQ password expiry
//	@Description	Submit a request to update HQ password expiry
//	@Tags			HQ
//	@Accept			json
//	@Produce		json
//	@Param			request	body		hqDto.UpdatePasswordExpiryRequest	true	"Update password expiry request payload"
//	@Success		200		{object}	localization.ResponseCode			"HQ password expiry update request submitted"
//	@Failure		400		{object}	localization.ResponseCode			"Invalid request payload"
//	@Failure		500		{object}	localization.ResponseCode			"Internal server error"
//	@Security		BearerAuth
//	@Router			/hq/password_expiry [post]
func (a *hqAdapter) UpdatePasswordExpiryRequest(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "", "updateHqPasswordExpiry", "handler", "hq")
	defer span.End()
	var request hqDto.UpdatePasswordExpiryRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		span.RecordError(err)
		localization.SendBadRequestResponse(w, localization.ErrorInvalidHQRequest.Code)
		return
	}

	if err := request.Validate(); err != nil {
		span.RecordError(err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := a.hqApp.UpdatePasswordExpiry(ctx, request); err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessHQPasswordExpiryUpdateRequestSubmitted, nil)
}
