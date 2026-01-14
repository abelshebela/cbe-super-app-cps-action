package donation

import (
	dto "cbe-super-app-cps-action/internal/constants/dto/donation"
	donation_interface "cbe-super-app-cps-action/internal/constants/interfaces/donation"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	core "cbe-super-app-cps-action/internal/handlers/rest/http/donation/core"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"encoding/json"
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
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "createDonation", "handler", "donation")
	defer span.End()
	req, err := core.ParseRequestFromMultipartForm(r, true)
	if err != nil {
		span.RecordError(err)
		d.logger.Errorf("failed to parse request from multipart form: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		span.RecordError(err)
		d.logger.Errorf("request validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	span.SetAttributes(
		attribute.String("donation.company_id", req.CompanyID),
		attribute.String("donation.category_id", req.CategoryID),
		attribute.String("donation.title", req.Title),
	)

	if err := d.donationApp.CreateDonation(ctx, req); err != nil {
		span.RecordError(err)
		d.logger.Errorf("failed to create donation: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	IsMakerOnly, ok := r.Context().Value(constants.ContextKey("is_maker_only")).(bool)
	if IsMakerOnly && ok {
		localization.SendSuccessResponse(w, localization.SuccessDonationCreatedSP, nil)
	}
	localization.SendSuccessResponse(w, localization.SuccessDonationCreateRequestSent, nil)
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
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "updateDonation", "handler", "donation")
	defer span.End()
	id := core.ExtractIDFromURL(r)
	if id == "" {
		d.logger.Errorf("donation ID is required")
		localization.SendBadRequestResponse(w, "donation ID is required")
		return
	}

	req, err := core.ParseRequestFromMultipartForm(r, false)
	if err != nil {
		span.RecordError(err)
		d.logger.Errorf("failed to parse request from multipart form: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := core.ValidateForUpdate(req); err != nil {
		span.RecordError(err)
		d.logger.Errorf("request validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.String("donation.id", id))

	err = d.donationApp.UpdateDonation(ctx, id, req)
	if err != nil {
		span.RecordError(err)
		d.logger.Errorf("failed to update donation: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	IsMakerOnly, ok := r.Context().Value(constants.ContextKey("is_maker_only")).(bool)
	if IsMakerOnly && ok {
		localization.SendSuccessResponse(w, localization.SuccessDonationUpdatedSP, nil)
	}
	localization.SendSuccessResponse(w, localization.SuccessDonationUpdateRequestSent, nil)
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
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "fetchDonations", "handler", "donation")
	defer span.End()
	filterParams := local_util.ExtractFilterParams(r)
	if filterParams.Page < 0 || filterParams.PerPage < 0 {
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}

	donations, err := d.donationApp.FetchDonation(ctx, filterParams)
	if err != nil {
		span.RecordError(err)
		d.logger.Errorf("failed to fetch donations: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.Int("donation.count", len(donations.Data)))
	localization.SendSuccessResponse(w, localization.SuccessDonationFetched, donations)
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
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "fetchDonationById", "handler", "donation")
	defer span.End()
	id := core.ExtractIDFromURL(r)
	if id == "" {
		d.logger.Errorf("donation ID is required")
		localization.SendBadRequestResponse(w, "donation ID is required")
		return
	}

	span.SetAttributes(attribute.String("donation.id", id))

	donation, err := d.donationApp.FetchDonationByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		d.logger.Errorf("failed to fetch donation by ID: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessDonationFetched, donation)
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
	id := core.ExtractIDFromURL(r)
	if id == "" {
		d.logger.Errorf("donation ID is required")
		localization.SendBadRequestResponse(w, "donation ID is required")
		return
	}

	req, err := core.ParseImageUpdateRequestFromMultipartForm(r)
	if err != nil {
		span.RecordError(err)
		d.logger.Errorf("failed to parse image update request from multipart form: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	// if err := core.ValidateImageUpdateRequest(req); err != nil {
	// 	d.logger.Errorf("image update validation failed: %v", err)
	// 	localization.SendBadRequestResponse(w, err.Error())
	// 	return
	// }

	span.SetAttributes(attribute.String("donation.id", id))
	if err := d.donationApp.UpdateDonationImage(ctx, id, req); err != nil {
		span.RecordError(err)
		d.logger.Errorf("failed to update donation image: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	IsMakerOnly, ok := r.Context().Value(constants.ContextKey("is_maker_only")).(bool)
	if IsMakerOnly && ok {
		localization.SendSuccessResponse(w, localization.SuccessDonationImageUpdatedSP, nil)
	}
	localization.SendSuccessResponse(w, localization.SuccessDonationImageUpdateRequestSent, nil)
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
	id := core.ExtractIDFromURL(r)
	if id == "" {
		d.logger.Errorf("donation ID is required")
		localization.SendBadRequestResponse(w, "donation ID is required")
		return
	}

	var req dto.DonationImageDeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		d.logger.Errorf("Failed to decode JSON request: %v", err)
		localization.SendBadRequestResponse(w, "invalid request body")
		return
	}

	if req.ImageID == "" {
		d.logger.Errorf("image ID is required")
		localization.SendBadRequestResponse(w, "image ID is required")
		return
	}

	span.SetAttributes(
		attribute.String("donation.id", id),
		attribute.String("donation.image_id", req.ImageID),
	)
	if err := d.donationApp.DeleteDonationImage(ctx, id, req.ImageID); err != nil {
		span.RecordError(err)
		d.logger.Errorf("failed to delete donation image: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	IsMakerOnly, ok := r.Context().Value(constants.ContextKey("is_maker_only")).(bool)
	if IsMakerOnly && ok {
		localization.SendSuccessResponse(w, localization.SuccessDonationImageDeletedSP, nil)
	}
	localization.SendSuccessResponse(w, localization.SuccessDonationImageDeleteRequestSent, nil)
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
	id := core.ExtractIDFromURL(r)
	if id == "" {
		d.logger.Errorf("donation ID is required")
		localization.SendBadRequestResponse(w, "donation ID is required")
		return
	}

	req, err := core.ParseRequestFromMultipartForm(r, false)
	if err != nil {
		span.RecordError(err)
		d.logger.Errorf("failed to parse request from multipart form: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := core.ValidateImageAdd(req); err != nil {
		span.RecordError(err)
		d.logger.Errorf("image validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.String("donation.id", id))
	if err := d.donationApp.AddDonationImage(ctx, id, req); err != nil {
		span.RecordError(err)
		d.logger.Errorf("failed to add donation image: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	IsMakerOnly, ok := r.Context().Value(constants.ContextKey("is_maker_only")).(bool)
	if IsMakerOnly && ok {
		localization.SendSuccessResponse(w, localization.SuccessDonationImageAddedSP, nil)
	}
	localization.SendSuccessResponse(w, localization.SuccessDonationImageAddRequestSent, nil)
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
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "enableDonation", "handler", "donation")
	defer span.End()
	id := core.ExtractIDFromURL(r)
	if id == "" {
		d.logger.Errorf("donation ID is required")
		localization.SendBadRequestResponse(w, "donation ID is required")
		return
	}

	span.SetAttributes(attribute.String("donation.id", id))
	if err := d.donationApp.EnableDonation(ctx, id); err != nil {
		span.RecordError(err)
		d.logger.Errorf("failed to enable donation: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	IsMakerOnly, ok := r.Context().Value(constants.ContextKey("is_maker_only")).(bool)
	if IsMakerOnly && ok {
		localization.SendSuccessResponse(w, localization.SuccessDonationEnabledSP, nil)
	}
	localization.SendSuccessResponse(w, localization.SuccessDonationEnableRequestSent, nil)
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
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "disableDonation", "handler", "donation")
	defer span.End()
	id := core.ExtractIDFromURL(r)
	if id == "" {
		d.logger.Errorf("donation ID is required")
		localization.SendBadRequestResponse(w, "donation ID is required")
		return
	}

	span.SetAttributes(attribute.String("donation.id", id))
	if err := d.donationApp.DisableDonation(ctx, id); err != nil {
		span.RecordError(err)
		d.logger.Errorf("failed to disable donation: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	IsMakerOnly, ok := r.Context().Value(constants.ContextKey("is_maker_only")).(bool)
	if IsMakerOnly && ok {
		localization.SendSuccessResponse(w, localization.SuccessDonationDisabledSP, nil)
	}
	localization.SendSuccessResponse(w, localization.SuccessDonationDisableRequestSent, nil)
}
