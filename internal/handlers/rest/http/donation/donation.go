package donation

import (
	dto "cbe-super-app-cps-action/internal/constants/dto/donation"
	donation_interface "cbe-super-app-cps-action/internal/constants/interfaces/donation"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	core "cbe-super-app-cps-action/internal/handlers/rest/http/donation/core"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"cbe-super-app-cps-action/internal/constants"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
)

type paginatedDonationResponse types.PaginatedResponse[[]dto.DonationListResponse]
type donationAdapter struct {
	donationApp service.DonationService
	logger      utils.Logger
}

func NewDonationAdapter(donationApp service.DonationService, logger utils.Logger) donation_interface.DonationHandler {
	return &donationAdapter{
		donationApp: donationApp,
		logger:      logger,
	}
}

// CreateDonation godoc
//
//	@Summary		Create a new donation
//	@Description	Create a new donation with the provided information
//	@Tags			Donation
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			donation_code			formData	string									false	"Donation code"
//	@Param			company_id				formData	string									true	"Company ID"
//	@Param			category_id				formData	string									true	"Category ID"
//	@Param			title					formData	string									true	"Title"
//	@Param			is_featured				formData	bool									false	"Is featured"
//	@Param			target					formData	integer									true	"Target amount"
//	@Param			donation_description	formData	string									true	"Donation description"
//	@Param			donation_images			formData	file									true	"Donation images (allow multiple with the same field name)"
//	@Param			cover_image				formData	file									false	"Cover image"
//	@Param			start_date				formData	string									false	"Start date (YYYY-MM-DD)"
//	@Param			end_date				formData	string									false	"End date (YYYY-MM-DD)"
//	@Param			enabled					formData	bool									false	"Enabled"
//	@Success		200						{object}	localization.StandardResponse{data=nil}	"Donation creation request sent successfully"
//	@Failure		400						{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500						{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/donation [post]
func (d *donationAdapter) CreateDonation(w http.ResponseWriter, r *http.Request) {
	fmt.Println(">>> [HANDLER] CreateDonation - ENTRY")
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "createDonation", "handler", "donation")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, d.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	fmt.Println(">>> [HANDLER] CreateDonation - parsing multipart form")
	req, err := core.ParseRequestFromMultipartForm(r, true)
	if err != nil {
		fmt.Println(">>> [HANDLER] CreateDonation - parse form ERROR:", err)
		span.RecordError(err)
		log.Errorf("[DonationH] parse form err: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	fmt.Printf(">>> [HANDLER] CreateDonation - parsed: title=%s, companyID=%s, categoryID=%s\n", req.Title, req.CompanyID, req.CategoryID)

	fmt.Println(">>> [HANDLER] CreateDonation - validating request")
	if err := req.Validate(); err != nil {
		fmt.Println(">>> [HANDLER] CreateDonation - validate ERROR:", err)
		span.RecordError(err)
		log.Errorf("[DonationH] validate err: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	span.SetAttributes(
		attribute.String("donation.company_id", req.CompanyID),
		attribute.String("donation.category_id", req.CategoryID),
		attribute.String("donation.title", req.Title),
	)

	fmt.Println(">>> [HANDLER] CreateDonation - calling service CreateDonation")
	if err := d.donationApp.CreateDonation(ctx, req); err != nil {
		fmt.Println(">>> [HANDLER] CreateDonation - service ERROR:", err)
		span.RecordError(err)
		log.Errorf("[DonationH][Create] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	fmt.Println(">>> [HANDLER] CreateDonation - service call SUCCESS, isMakerOnly:", md.IsMakerOnly)
	if md.IsMakerOnly {
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessDonationCreatedSP, nil)
	} else {
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessDonationCreateRequestSent, nil)
	}
	fmt.Println(">>> [HANDLER] CreateDonation - EXIT")
}

// UpdateDonation godoc
//
//	@Summary		Update a donation
//	@Description	Update a donation with the provided information
//	@Tags			Donation
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id						path		string									true	"Donation ID"
//	@Param			donation_code			formData	string									false	"Donation code"
//	@Param			company_id				formData	string									false	"Company ID"
//	@Param			category_id				formData	string									false	"Category ID"
//	@Param			title					formData	string									false	"Title"
//	@Param			is_featured				formData	bool									false	"Is featured"
//	@Param			target					formData	integer									false	"Target amount"
//	@Param			donation_description	formData	string									false	"Donation description"
//	@Param			donation_images			formData	file									false	"Donation images (allow multiple with the same field name)"
//	@Param			cover_image				formData	file									false	"Cover image"
//	@Param			start_date				formData	string									false	"Start date (YYYY-MM-DD)"
//	@Param			end_date				formData	string									false	"End date (YYYY-MM-DD)"
//	@Param			enabled					formData	bool									false	"Enabled"
//	@Success		200						{object}	localization.StandardResponse{data=nil}	"Donation update request sent successfully"
//	@Failure		400						{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500						{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/donation/{id} [patch]
func (d *donationAdapter) UpdateDonation(w http.ResponseWriter, r *http.Request) {
	fmt.Println(">>> [HANDLER] UpdateDonation - ENTRY")
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "updateDonation", "handler", "donation")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, d.logger)
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	id := core.ExtractIDFromURL(r)
	fmt.Println(">>> [HANDLER] UpdateDonation - extracted id:", id)
	if id == "" {
		log.Errorf("[DonationH] id required")
		localization.SendBadRequestResponse(w, "donation ID is required")
		return
	}

	fmt.Println(">>> [HANDLER] UpdateDonation - parsing multipart form")
	req, err := core.ParseRequestFromMultipartForm(r, false)
	if err != nil {
		fmt.Println(">>> [HANDLER] UpdateDonation - parse form ERROR:", err)
		span.RecordError(err)
		log.Errorf("[DonationH] parse form err: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	fmt.Println(">>> [HANDLER] UpdateDonation - validating request")
	if err := core.ValidateForUpdate(req); err != nil {
		fmt.Println(">>> [HANDLER] UpdateDonation - validate ERROR:", err)
		span.RecordError(err)
		log.Errorf("[DonationH] validate err: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.String("donation.id", id))

	fmt.Println(">>> [HANDLER] UpdateDonation - calling service UpdateDonation")
	err = d.donationApp.UpdateDonation(ctx, id, req)
	if err != nil {
		fmt.Println(">>> [HANDLER] UpdateDonation - service ERROR:", err)
		span.RecordError(err)
		log.Errorf("[DonationH][Update] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	fmt.Println(">>> [HANDLER] UpdateDonation - service call SUCCESS")
	if md.IsMakerOnly {
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessDonationUpdatedSP, nil)
	} else {
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessDonationUpdateRequestSent, nil)
	}
	fmt.Println(">>> [HANDLER] UpdateDonation - EXIT")
}

// FetchDonation godoc
//
//	@Summary		List donations
//	@Description	Retrieve donations with pagination and optional search
//	@Tags			Donation
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int																false	"Page number"		default(1)
//	@Param			per_page	query		int																false	"Items per page"	default(10)
//	@Param			search		query		string															false	"Search term"
//	@Success		200			{object}	localization.StandardResponse{data=[]paginatedDonationResponse}	"Donations retrieved successfully"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}							"Internal server error"
//	@Security		BearerAuth
//	@Router			/donation [get]
func (d *donationAdapter) FetchDonation(w http.ResponseWriter, r *http.Request) {
	fmt.Println(">>> [HANDLER] FetchDonation - ENTRY")
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "fetchDonations", "handler", "donation")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, d.logger)

	filterParams := local_util.ExtractFilterParams(r)
	fmt.Printf(">>> [HANDLER] FetchDonation - filterParams: page=%d, perPage=%d, search=%s\n", filterParams.Page, filterParams.PerPage, filterParams.Search)

	search := r.URL.Query().Get("search")
	filter := r.URL.Query().Get("filter")

	if err := local_util.NoSpecialChars(search); err != nil {
		fmt.Println(">>> [HANDLER] FetchDonation - search special chars ERROR:", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := local_util.NoSpecialChars(filter); err != nil {
		fmt.Println(">>> [HANDLER] FetchDonation - filter special chars ERROR:", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if filterParams.Page < 0 || filterParams.PerPage < 0 {
		fmt.Println(">>> [HANDLER] FetchDonation - invalid page/perPage")
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}

	fmt.Println(">>> [HANDLER] FetchDonation - calling service FetchDonation")
	donations, err := d.donationApp.FetchDonation(ctx, filterParams)
	if err != nil {
		fmt.Println(">>> [HANDLER] FetchDonation - service ERROR:", err)
		span.RecordError(err)
		log.Errorf("[DonationH][FetchAll] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	fmt.Printf(">>> [HANDLER] FetchDonation - got %d donations\n", len(donations.Data))
	span.SetAttributes(attribute.Int("donation.count", len(donations.Data)))
	localization.SendSuccessResponse(w, localization.SuccessDonationFetched, donations)
	fmt.Println(">>> [HANDLER] FetchDonation - EXIT")
}

// FetchDonationByID godoc
//
//	@Summary		Get donation by ID
//	@Description	Retrieve a donation's details by ID
//	@Tags			Donation
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string															true	"Donation ID"
//	@Success		200	{object}	localization.StandardResponse{data=donation.DonationResponse}	"Donation retrieved successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}							"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}							"Not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}							"Internal server error"
//	@Security		BearerAuth
//	@Router			/donation/{id} [get]
func (d *donationAdapter) FetchDonationByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println(">>> [HANDLER] FetchDonationByID - ENTRY")
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "fetchDonationById", "handler", "donation")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, d.logger)
	id := core.ExtractIDFromURL(r)
	fmt.Println(">>> [HANDLER] FetchDonationByID - id:", id)
	if id == "" {
		log.Errorf("[DonationH] id required")
		localization.SendBadRequestResponse(w, "donation ID is required")
		return
	}

	span.SetAttributes(attribute.String("donation.id", id))

	fmt.Println(">>> [HANDLER] FetchDonationByID - calling service")
	donation, err := d.donationApp.FetchDonationByID(ctx, id)
	if err != nil {
		fmt.Println(">>> [HANDLER] FetchDonationByID - service ERROR:", err)
		span.RecordError(err)
		log.Errorf("[DonationH][FetchByID] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	fmt.Println(">>> [HANDLER] FetchDonationByID - SUCCESS, sending response")
	localization.SendSuccessResponse(w, localization.SuccessDonationFetched, donation)
	fmt.Println(">>> [HANDLER] FetchDonationByID - EXIT")
}

// UpdateDonationImage godoc
//
//	@Summary		Update a donation image
//	@Description	Update a specific donation image by ID
//	@Tags			Donation
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id				path		string									true	"Donation ID"
//	@Param			image_id		formData	string									true	"Image ID"
//	@Param			donation_images	formData	file									true	"New image file"
//	@Success		200				{object}	localization.StandardResponse{data=nil}	"Donation image update request sent successfully"
//	@Failure		400				{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500				{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/donation/image/{id} [patch]
func (d *donationAdapter) UpdateDonationImage(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "updateDonationImage", "handler", "donation")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, d.logger)
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	id := core.ExtractIDFromURL(r)
	if id == "" {
		log.Errorf("[DonationH] id required")
		localization.SendBadRequestResponse(w, "donation ID is required")
		return
	}

	req, err := core.ParseImageUpdateRequestFromMultipartForm(r)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[DonationH][UpdateImage] parse form err: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	// if err := core.ValidateImageUpdateRequest(req); err != nil {
	// 	log.Errorf("image update validation failed: %v", err)
	// 	localization.SendBadRequestResponse(w, err.Error())
	// 	return
	// }

	span.SetAttributes(attribute.String("donation.id", id))
	if err := d.donationApp.UpdateDonationImage(ctx, id, req); err != nil {
		span.RecordError(err)
		log.Errorf("[DonationH][UpdateImage] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessDonationImageUpdatedSP, nil)
	} else {
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessDonationImageUpdateRequestSent, nil)
	}
}

// DeleteDonationImage godoc
//
//	@Summary		Delete a donation image
//	@Description	Delete a specific donation image by ID
//	@Tags			Donation
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string									true	"Donation ID"
//	@Param			request	body		donation.DonationImageDeleteRequest		true	"Image delete request"
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Donation image delete request sent successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/donation/image/{id} [delete]
func (d *donationAdapter) DeleteDonationImage(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "deleteDonationImage", "handler", "donation")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, d.logger)
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	id := core.ExtractIDFromURL(r)
	if id == "" {
		log.Errorf("[DonationH] id required")
		localization.SendBadRequestResponse(w, "donation ID is required")
		return
	}

	var req dto.DonationImageDeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		log.Errorf("[DonationH][DeleteImage] decode err: %v", err)
		localization.SendBadRequestResponse(w, "invalid request body")
		return
	}

	if req.ImageID == "" {
		log.Errorf("[DonationH][DeleteImage] image id required")
		localization.SendBadRequestResponse(w, "image ID is required")
		return
	}

	span.SetAttributes(
		attribute.String("donation.id", id),
		attribute.String("donation.image_id", req.ImageID),
	)
	if err := d.donationApp.DeleteDonationImage(ctx, id, req.ImageID); err != nil {
		span.RecordError(err)
		log.Errorf("[DonationH][DeleteImage] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessDonationImageDeletedSP, nil)
	} else {
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessDonationImageDeleteRequestSent, nil)
	}
}

// AddDonationImage godoc
//
//	@Summary		Add donation image(s)
//	@Description	Add one or more images to a donation
//	@Tags			Donation
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id				path		string									true	"Donation ID"
//	@Param			donation_images	formData	file									true	"Donation images (allow multiple with the same field name)"
//	@Success		200				{object}	localization.StandardResponse{data=nil}	"Donation image add request sent successfully"
//	@Failure		400				{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500				{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/donation/image/{id} [post]
func (d *donationAdapter) AddDonationImage(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "addDonationImage", "handler", "donation")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, d.logger)

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	id := core.ExtractIDFromURL(r)
	if id == "" {
		log.Errorf("[DonationH] id required")
		localization.SendBadRequestResponse(w, "donation ID is required")
		return
	}

	req, err := core.ParseRequestFromMultipartForm(r, false)
	if err != nil {
		span.RecordError(err)
		log.Errorf("[DonationH] parse form err: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := core.ValidateImageAdd(req); err != nil {
		span.RecordError(err)
		log.Errorf("[DonationH][AddImage] validate err: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.String("donation.id", id))
	if err := d.donationApp.AddDonationImage(ctx, id, req); err != nil {
		span.RecordError(err)
		log.Errorf("[DonationH][AddImage] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessDonationImageAddedSP, nil)
	} else {
		localization.SendSuccessResponse(w, localization.SuccessDonationImageAddRequestSent, nil)
	}
}

// EnableDonation godoc
//
//	@Summary		Enable a donation
//	@Description	Enable a donation by ID
//	@Tags			Donation
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Donation ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Donation enable request sent successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}	"Donation not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/donation/enable/{id} [patch]
func (d *donationAdapter) EnableDonation(w http.ResponseWriter, r *http.Request) {
	fmt.Println(">>> [HANDLER] EnableDonation - ENTRY")
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "enableDonation", "handler", "donation")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, d.logger)
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	id := core.ExtractIDFromURL(r)
	fmt.Println(">>> [HANDLER] EnableDonation - id:", id)
	if id == "" {
		log.Errorf("[DonationH] id required")
		localization.SendBadRequestResponse(w, "donation ID is required")
		return
	}

	span.SetAttributes(attribute.String("donation.id", id))
	fmt.Println(">>> [HANDLER] EnableDonation - calling service")
	if err := d.donationApp.EnableDonation(ctx, id); err != nil {
		fmt.Println(">>> [HANDLER] EnableDonation - service ERROR:", err)
		span.RecordError(err)
		log.Errorf("[DonationH][Enable] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	fmt.Println(">>> [HANDLER] EnableDonation - SUCCESS")
	if md.IsMakerOnly {
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessDonationEnabledSP, nil)
	} else {
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessDonationEnableRequestSent, nil)
	}
	fmt.Println(">>> [HANDLER] EnableDonation - EXIT")
}

// DisableDonation godoc
//
//	@Summary		Disable a donation
//	@Description	Disable a donation by ID
//	@Tags			Donation
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Donation ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Donation disable request sent successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}	"Donation not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/donation/disable/{id} [patch]
func (d *donationAdapter) DisableDonation(w http.ResponseWriter, r *http.Request) {
	fmt.Println(">>> [HANDLER] DisableDonation - ENTRY")
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "disableDonation", "handler", "donation")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, d.logger)
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)
	localization.UpdateWriterContext(w, ctx)

	id := core.ExtractIDFromURL(r)
	fmt.Println(">>> [HANDLER] DisableDonation - id:", id)
	if id == "" {
		log.Errorf("[DonationH] id required")
		localization.SendBadRequestResponse(w, "donation ID is required")
		return
	}

	span.SetAttributes(attribute.String("donation.id", id))
	fmt.Println(">>> [HANDLER] DisableDonation - calling service")
	if err := d.donationApp.DisableDonation(ctx, id); err != nil {
		fmt.Println(">>> [HANDLER] DisableDonation - service ERROR:", err)
		span.RecordError(err)
		log.Errorf("[DonationH][Disable] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	fmt.Println(">>> [HANDLER] DisableDonation - SUCCESS")
	if md.IsMakerOnly {
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessDonationDisabledSP, nil)
	} else {
		w = localization.ApplyActionCodeHeaderFromWriter(w, ctx)
		localization.SendSuccessResponse(w, localization.SuccessDonationDisableRequestSent, nil)
	}
	fmt.Println(">>> [HANDLER] DisableDonation - EXIT")
}

// ExportDonationList godoc
//
//	@Summary		Export donation list
//	@Description	Export donations within a date range as a CSV file and return the download link
//	@Tags			Donation
//	@Accept			json
//	@Produce		json
//	@Param			file_type	query		string									true	"Export file type (e.g. csv)"
//	@Param			From		query		string									true	"Start date (YYYY-MM-DD)"
//	@Param			To			query		string									true	"End date (YYYY-MM-DD)"
//	@Success		200			{object}	localization.StandardResponse{data=string}	"File link returned successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/donation/export [get]
func (a *donationAdapter) ExportDonationList(w http.ResponseWriter, r *http.Request) {
	fmt.Println(">>> [HANDLER] ExportDonationList - ENTRY")
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "exportDonationList", "handler", "donation")
	defer span.End()
	log := local_util.LoggerFromCtx(ctx, a.logger)

	fileType := r.URL.Query().Get("file_type")
	from := r.URL.Query().Get("From")
	to := r.URL.Query().Get("To")
	fmt.Printf(">>> [HANDLER] ExportDonationList - params: file_type=%q, From=%q, To=%q\n", fileType, from, to)

	if fileType == "" || from == "" || to == "" {
		fmt.Println(">>> [HANDLER] ExportDonationList - MISSING required params")
		log.Warnf("[DonationH][Export] missing required params: file_type=%q, From=%q, To=%q", fileType, from, to)
		localization.SendBadRequestResponse(w, localization.ErrorRequiredFieldMissing.Message)
		return
	}

	// Normalize date strings (accepts YYYY-MM-DD or RFC3339)
	fmt.Println(">>> [HANDLER] ExportDonationList - normalizing date range")
	fromNorm, toNorm, err := local_util.FormatDateRangeToUTCStrings(from, to)
	if err != nil {
		fmt.Println(">>> [HANDLER] ExportDonationList - date normalization ERROR:", err)
		log.Warnf("[DonationH][Export] invalid date format: From=%s, To=%s, err=%v", from, to, err)
		localization.SendBadRequestResponse(w, localization.ErrorInvalidFormat.Message)
		return
	}
	fmt.Printf(">>> [HANDLER] ExportDonationList - normalized: from=%s, to=%s\n", fromNorm, toNorm)

	fmt.Println(">>> [HANDLER] ExportDonationList - parsing start date")
	startDate, err := local_util.ValidateTimeAndParse(fromNorm)
	if err != nil {
		fmt.Println(">>> [HANDLER] ExportDonationList - start date parse ERROR:", err)
		log.Warnf("[DonationH][Export] invalid start date: %s", from)
		localization.SendBadRequestResponse(w, localization.ErrorInvalidFormat.Message)
		return
	}
	fmt.Println(">>> [HANDLER] ExportDonationList - startDate:", startDate)

	fmt.Println(">>> [HANDLER] ExportDonationList - parsing end date")
	endDate, err := local_util.ValidateTimeAndParse(toNorm)
	if err != nil {
		fmt.Println(">>> [HANDLER] ExportDonationList - end date parse ERROR:", err)
		log.Warnf("[DonationH][Export] invalid end date: %s", to)
		localization.SendBadRequestResponse(w, localization.ErrorInvalidFormat.Message)
		return
	}
	fmt.Println(">>> [HANDLER] ExportDonationList - endDate:", endDate)

	if endDate.Before(startDate) {
		fmt.Println(">>> [HANDLER] ExportDonationList - endDate before startDate")
		log.Warnf("[DonationH][Export] invalid date range: start=%v, end=%v", startDate, endDate)
		localization.SendBadRequestResponse(w, localization.ErrorInvalidFormat.Message)
		return
	}

	span.SetAttributes(
		attribute.String("export.file_type", fileType),
		attribute.String("export.from", from),
		attribute.String("export.to", to),
	)

	fmt.Println(">>> [HANDLER] ExportDonationList - calling service ExportDonationData")
	fileLink, err := a.donationApp.ExportDonationData(ctx, startDate, endDate, fileType)
	if err != nil {
		fmt.Println(">>> [HANDLER] ExportDonationList - service ERROR:", err)
		span.RecordError(err)
		log.Errorf("[DonationH][Export] svc err: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	fmt.Println(">>> [HANDLER] ExportDonationList - SUCCESS, fileLink:", fileLink)
	localization.SendSuccessResponse(w, localization.DonationDataExportedSuccess, fileLink)
	fmt.Println(">>> [HANDLER] ExportDonationList - EXIT")
}
