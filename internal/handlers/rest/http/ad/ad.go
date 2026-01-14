package ad

import (
	"cbe-super-app-cps-action/internal/handlers/rest/http/ad/core"
	"cbe-super-app-cps-action/internal/service"
	"net/http"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"go.opentelemetry.io/otel/attribute"

	ad_dto "cbe-super-app-cps-action/internal/constants/dto/ad"
	advertInbound "cbe-super-app-cps-action/internal/constants/interfaces/ad"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"

	"cbe-super-app-cps-action/internal/constants"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type paginated_advert_response types.PaginatedResponse[[]*ad_dto.AdvertResponse]
type advertAdapter struct {
	advertApplication service.AdvertService
	logger            utils.Logger
}

func InitAdvertAdapter(advertApplication service.AdvertService, logger utils.Logger) advertInbound.ADAdapter {
	return &advertAdapter{
		logger:            logger,
		advertApplication: advertApplication,
	}
}

// CreateAdvert godoc
//
//	@Summary		Create a new advert (maker)
//	@Description	Submit an advert create request. Requires multipart/form-data with optional banner image.
//	@Tags			Adverts
//	@Accept			mpfd
//	@Produce		json
//	@Param			title			formData	string									true	"Title"				minLength(3)		maxLength(20)	example("New Promo")
//	@Param			description		formData	string									true	"Description"		minLength(30)		maxLength(100)	example("Enjoy our new promotion valid this weekend only.")
//	@Param			advert_for		formData	string									true	"Advert audience"	Enums(IFB,CB,ALL)	example(ALL)
//	@Param			banner_image	formData	file									false	"Banner image (<=2MB; jpeg/png/gif/webp)"
//	@Success		200				{object}	localization.StandardResponse{data=nil}	"Advert create request sent"
//	@Failure		400				{object}	localization.StandardResponse{data=nil}	"Invalid request"
//	@Failure		409				{object}	localization.StandardResponse{data=nil}	"Duplicate or conflict"
//	@Failure		500				{object}	localization.StandardResponse{data=nil}	"Server error"
//	@Security		BearerAuth
//	@Router			/adverts [post]
func (a *advertAdapter) CreateAdvert(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "createAdvert", "handler", "advert")
	defer span.End()

	req, err := core.ParseBannerImage(r, true)
	if err != nil {
		span.RecordError(err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	if err := req.Validate(false); err != nil {
		span.RecordError(err)
		a.logger.Errorf("advert create  update request validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	domainReq, err := core.ToAdvert(req)
	if err != nil {
		span.RecordError(err)
		a.logger.Errorf("[event.CreateAdvert] failed to convert to domain advert, error: %v", err)
		localization.SendBadRequestResponse(w, localization.ErrorInvalidRequest.Message)
		return
	}

	span.SetAttributes(
		attribute.String("advert.title", req.Title),
		attribute.String("advert.for", req.AdvertFor),
	)

	err = a.advertApplication.CreateAdvert(ctx, &domainReq, req.BannerImage)
	if err != nil {
		span.RecordError(err)
		a.logger.Errorf("[CreateAdvert] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	IsMakerOnly, ok := r.Context().Value(constants.ContextKey("is_maker_only")).(bool)
	if IsMakerOnly && ok {
		localization.SendSuccessResponse(w, localization.SuccessAdvertCreatedSP, nil)
	}
	a.logger.Infof("[CreateAdvert] request sent successfully for title: %s", req.Title)
	localization.SendSuccessResponse(w, localization.SuccessAdvertCreateRequestSent, nil)

}

// FetchAdverts godoc
//
//	@Summary		List adverts
//	@Description	Fetch adverts with pagination, optional text search, and field filtering.
//	@Tags			Adverts
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int																false	"Page number"															minimum(1)	default(1)		example(1)
//	@Param			per_page	query		int																false	"Items per page"														minimum(1)	maximum(100)	default(10)	example(10)
//	@Param			search		query		string															false	"Text search on title & description (case-insensitive partial match)"	example("promo")
//	@Param			title		query		string															false	"Filter by exact title match"											example("Summer Promo")
//	@Param			description	query		string															false	"Filter by exact description match"										example("Discount event")
//	@Param			enabled		query		bool															false	"Filter by enabled flag"												example(true)
//	@Success		200			{object}	localization.StandardResponse{data=paginated_advert_response}	"Adverts fetched successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}							"Invalid pagination params"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}							"Server error"
//	@Security		BearerAuth
//	@Router			/adverts [get]
func (a *advertAdapter) FetchAdverts(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "fetchAdverts", "handler", "advert")
	defer span.End()
	filterParams := local_util.ExtractFilterParams(r)
	list, err := a.advertApplication.FetchAdverts(ctx, *filterParams)
	if err != nil {
		span.RecordError(err)
		a.logger.Errorf("[FetchAdverts] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	docs := core.ToAdvertResponses(list.Data)
	res := paginated_advert_response{
		Data: docs,
		Meta: list.Meta,
	}
	a.logger.Infof("[FetchAdverts] retrieved %d adverts", len(docs))
	localization.SendSuccessResponse(w, localization.SuccessAdvertsFetched, res)
}

// FetchAdvertByID godoc
//
//	@Summary		Get advert by ID
//	@Description	Retrieve a single advert by its identifier.
//	@Tags			Adverts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string													true	"Advert ID"
//	@Success		200	{object}	localization.StandardResponse{data=ad.AdvertResponse}	"Advert fetched successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}					"Invalid ID"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}					"Advert not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}					"Server error"
//	@Security		BearerAuth
//	@Router			/adverts/{id} [get]
func (a *advertAdapter) FetchAdvertByID(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "fetchAdvertById", "handler", "advert")
	defer span.End()
	id, err := core.ExtractID(r, a.logger)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.String("advert.id", id))

	data, err := a.advertApplication.FetchAdvertByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		a.logger.Errorf("[FetchAdvertByID] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	res := core.ToAdvertResponse(*data)
	a.logger.Infof("[FetchAdvertByID] advert retrieved successfully for id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessAdvertFetched, res)
}

// UpdateAdvert godoc
//
//	@Summary		Update existing advert (maker)
//	@Description	Submit an advert update request. Provide only fields to change. Multipart/form-data supported for banner_image.
//	@Tags			Adverts
//	@Accept			mpfd
//	@Produce		json
//	@Param			id				path		string									true	"Advert ID"
//	@Param			title			formData	string									false	"Title"				minLength(3)	maxLength(20)	example("Weekend Promo")
//	@Param			description		formData	string									false	"Description"		minLength(30)	maxLength(100)
//	@Param			advert_for		formData	string									false	"Advert audience"	Enums(IFB,CB,ALL)
//	@Param			banner_image	formData	file									false	"Banner image (<=2MB; jpeg/png/gif/webp)"
//	@Success		200				{object}	localization.StandardResponse{data=nil}	"Advert update request sent"
//	@Failure		400				{object}	localization.StandardResponse{data=nil}	"No data provided for update / invalid payload"
//	@Failure		404				{object}	localization.StandardResponse{data=nil}	"Advert not found"
//	@Failure		500				{object}	localization.StandardResponse{data=nil}	"Server error"
//	@Security		BearerAuth
//	@Router			/adverts/{id} [patch]
func (a *advertAdapter) UpdateAdvert(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "updateAdvert", "handler", "advert")
	defer span.End()
	id, err := core.ExtractID(r, a.logger)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	req, err := core.ParseBannerImage(r, false)
	if err != nil {
		span.RecordError(err)
		a.logger.Errorf("[event.UpdateAdvert] failed to parse and validate advert request, id: %s, error: %v", id, err.Error())
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if err := req.Validate(true); err != nil {
		span.RecordError(err)
		a.logger.Errorf("advert update request validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	domainReq, _ := core.ToAdvert(req)

	if domainReq.Title == "" && domainReq.Description == "" && domainReq.AdvertFor == "" && req.BannerImage == nil {
		a.logger.Errorf("[event.UpdateAdvert] no data provided for update, id: %s", id)
		localization.SendBadRequestResponse(w, localization.ErrorNoDataProvidedForUpdate.Message)
		return
	}

	span.SetAttributes(
		attribute.String("advert.id", id),
		attribute.String("advert.for", string(domainReq.AdvertFor)),
	)

	err = a.advertApplication.UpdateAdvert(ctx, id, &domainReq, req.BannerImage)
	if err != nil {
		span.RecordError(err)
		a.logger.Errorf("[UpdateAdvert] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	IsMakerOnly, ok := r.Context().Value(constants.ContextKey("is_maker_only")).(bool)
	if IsMakerOnly && ok {
		localization.SendSuccessResponse(w, localization.SuccessAdvertUpdatedSP, nil)
	}
	a.logger.Infof("[UpdateAdvert] request sent successfully for id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessAdvertUpdateRequestSent, nil)
}

// DeleteAdvert godoc
//
//	@Summary		Delete advert (maker)
//	@Description	Submit an advert delete request.
//	@Tags			Adverts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Advert ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Delete request sent"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Invalid ID"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}	"Advert not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Server error"
//	@Security		BearerAuth
//	@Router			/advert/{id} [delete]
func (a *advertAdapter) DeleteAdvert(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "deleteAdvert", "handler", "advert")
	defer span.End()
	id, err := core.ExtractID(r, a.logger)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.String("advert.id", id))

	err = a.advertApplication.DeleteAdvert(ctx, id)
	if err != nil {
		span.RecordError(err)
		a.logger.Errorf("[DeleteAdvert] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	IsMakerOnly, ok := r.Context().Value(constants.ContextKey("is_maker_only")).(bool)
	if IsMakerOnly && ok {
		localization.SendSuccessResponse(w, localization.SuccessAdvertDeletedSP, nil)
	}
	a.logger.Infof("[DeleteAdvert] request sent successfully for id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessAdvertDeleteRequestSent, nil)
}

// EnableAdvert godoc
//
//	@Summary		Enable advert (checker)
//	@Description	Approve enable request for an advert.
//	@Tags			Adverts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Advert ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Enable request sent"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Invalid ID"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}	"Advert not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Server error"
//	@Security		BearerAuth
//	@Router			/adverts/{id}/enable [patch]
func (a *advertAdapter) EnableAdvert(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "enableAdvert", "handler", "advert")
	defer span.End()
	id, err := core.ExtractID(r, a.logger)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.String("advert.id", id))

	err = a.advertApplication.EnableDisableAdvert(ctx, id, true)
	if err != nil {
		span.RecordError(err)
		a.logger.Errorf("[EnableAdvert] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	IsMakerOnly, ok := r.Context().Value(constants.ContextKey("is_maker_only")).(bool)
	if IsMakerOnly && ok {
		localization.SendSuccessResponse(w, localization.SuccessAdvertEnabledSP, nil)
	}
	a.logger.Infof("[EnableAdvert] request sent successfully for id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessAdvertEnableRequestSent, nil)

}

// DisableAdvert godoc
//
//	@Summary		Disable advert (checker)
//	@Description	Approve disable request for an advert.
//	@Tags			Adverts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Advert ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Disable request sent"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Invalid ID"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}	"Advert not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Server error"
//	@Security		BearerAuth
//	@Router			/adverts/{id}/disable [patch]
func (a *advertAdapter) DisableAdvert(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "disableAdvert", "handler", "advert")
	defer span.End()
	id, err := core.ExtractID(r, a.logger)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.String("advert.id", id))

	err = a.advertApplication.EnableDisableAdvert(ctx, id, false)
	if err != nil {
		span.RecordError(err)
		a.logger.Errorf("[DisableAdvert] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	IsMakerOnly, ok := r.Context().Value(constants.ContextKey("is_maker_only")).(bool)
	if IsMakerOnly && ok {
		localization.SendSuccessResponse(w, localization.SuccessAdvertDisabledSP, nil)
	}
	a.logger.Infof("[DisableAdvert] request sent successfully for id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessAdvertDisableRequestSent, nil)

}
