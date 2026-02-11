package donation_company

import (
	"cbe-super-app-cps-action/internal/constants/dto/donation_company"
	inbound "cbe-super-app-cps-action/internal/constants/interfaces/donation_company"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"context"
	"net/http"

	"cbe-super-app-cps-action/internal/handlers/rest/http/donation_company/core"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"cbe-super-app-cps-action/internal/constants"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
)

type paginatedDonationCompanyListResponse types.PaginatedResponse[[]donation_company.DonationCompanyListResponse]
type donationCompanyAdapter struct {
	logger             utils.Logger
	donationCompanyApp service.DonationCompanyService
}

func InitDonationCompanyAdapter(donationCompanyApp service.DonationCompanyService, logger utils.Logger) inbound.DonationCompanyAdapter {
	return &donationCompanyAdapter{
		logger:             logger,
		donationCompanyApp: donationCompanyApp,
	}
}

// FetchDonationCompany godoc
//
//	@Summary		List donation companies
//	@Description	Retrieve donation companies with pagination and optional search
//	@Tags			Donation Company
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int																			false	"Page number"		default(1)
//	@Param			per_page	query		int																			false	"Items per page"	default(10)
//	@Param			search		query		string																		false	"Search term"
//	@Success		200			{object}	localization.StandardResponse{data=paginatedDonationCompanyListResponse}	"Donation companies retrieved successfully"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}										"Internal server error"
//	@Security		BearerAuth
//	@Router			/donation_company [get]
func (d *donationCompanyAdapter) FetchDonationCompany(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "fetchDonationCompany", "handler", "donationCompany")
	defer span.End()
	filterParams := local_util.ExtractFilterParams(r)

	search := r.URL.Query().Get("search")
	filter := r.URL.Query().Get("filter")

	if err := local_util.NoSpecialChars(search); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := local_util.NoSpecialChars(filter); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if filterParams.Page < 0 || filterParams.PerPage < 0 {
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}

	donationCompanies, err := d.donationCompanyApp.FetchDonationCompany(ctx, filterParams)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.Int("donation_company.count", len(donationCompanies.Data)))
	localization.SendSuccessResponse(w, localization.SuccessDonationCompaniesFetched, donationCompanies)
}

// FetchDonationCompanyByID godoc
//
//	@Summary		Get donation company by ID
//	@Description	Retrieve a donation company's details by ID
//	@Tags			Donation Company
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string																				true	"Donation Company ID"
//	@Success		200	{object}	localization.StandardResponse{data=donation_company.DonationCompanyListResponse}	"Donation company retrieved successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}												"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}												"Not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}												"Internal server error"
//	@Security		BearerAuth
//	@Router			/donation_company/{id} [get]
func (d *donationCompanyAdapter) FetchDonationCompanyByID(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "fetchDonationCompanyById", "handler", "donationCompany")
	defer span.End()
	id := chi.URLParam(r, "id")
	if id == "" {
		d.logger.Errorf("donation company ID is required to fetch one")
		localization.SendErrorResponse(w, localization.ErrorDonationCompanyIdRequired, nil, nil)
		return
	}

	span.SetAttributes(attribute.String("donation_company.id", id))
	donationCompany, err := d.donationCompanyApp.FetchDonationCompanyByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessDonationCompanyFetched, donationCompany)
}

// CreateDonationCompany godoc
//
//	@Summary		Create a new donation company
//	@Description	Create a new donation company with the provided information
//	@Tags			Donation Company
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			company_name	formData	string									true	"Company name"
//	@Param			company_logo	formData	file									true	"Company logo image file"
//	@Param			account_number	formData	string									true	"Account number"
//	@Success		200				{object}	localization.StandardResponse{data=nil}	"Donation company creation request sent successfully"
//	@Failure		400				{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500				{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/donation_company [post]
func (d *donationCompanyAdapter) CreateDonationCompany(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "createDonationCompany", "handler", "donationCompany")
	defer span.End()
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	req, err := core.ParseRequestFromMultipartForm(r, true)
	if err != nil {
		span.RecordError(err)
		d.logger.Errorf("failed to parse donation company request from multipart form: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if err := req.Validate(); err != nil {
		span.RecordError(err)
		d.logger.Errorf("donation company request validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	formattedPhone := local_util.FormatPhoneNumber(req.PhoneNumber)
	req.PhoneNumber = formattedPhone
	span.SetAttributes(
		attribute.String("donation_company.name", req.CompanyName),
		attribute.String("donation_company.account_number", req.AccountNumber),
	)
	if err := d.donationCompanyApp.CreateDonationCompany(ctx, req); err != nil {
		span.RecordError(err)
		d.logger.Errorf("failed to create donation company: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessDonationCompanyCreatedSP, nil)
	} else {
		d.logger.Infof("donation company creation request submitted successfully")
		localization.SendSuccessResponse(w, localization.SuccessDonationCompanyCreateRequestSent, nil)
	}
}

// UpdateDonationCompany godoc
//
//	@Summary		Update a donation company
//	@Description	Update a donation company with the provided information
//	@Tags			Donation Company
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id				path		string									true	"Donation Company ID"
//	@Param			company_name	formData	string									false	"Company name"
//	@Param			company_logo	formData	file									false	"Company logo image file"
//	@Param			account_number	formData	string									false	"Account number"
//	@Success		200				{object}	localization.StandardResponse{data=nil}	"Donation company updated successfully"
//	@Failure		400				{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		404				{object}	localization.StandardResponse{data=nil}	"Not found"
//	@Failure		500				{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/donation_company/{id} [patch]
func (d *donationCompanyAdapter) UpdateDonationCompany(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "updateDonationCompany", "handler", "donationCompany")
	defer span.End()
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	id := chi.URLParam(r, "id")
	if id == "" {
		d.logger.Errorf("donation company ID is required for update")
		localization.SendErrorResponse(w, localization.ErrorDonationCompanyIdRequired, nil, nil)
		return
	}

	req, err := core.ParseRequestFromMultipartForm(r, false)
	if err != nil {
		span.RecordError(err)
		d.logger.Errorf("failed to parse donation company update request from multipart form: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := req.ValidateForUpdate(); err != nil {
		span.RecordError(err)
		d.logger.Errorf("donation company update request validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	span.SetAttributes(
		attribute.String("donation_company.id", id),
		attribute.String("donation_company.name", req.CompanyName),
	)
	updatedDonationCompany, err := d.donationCompanyApp.UpdateDonationCompany(ctx, id, req)
	if err != nil {
		span.RecordError(err)
		d.logger.Errorf("failed to update donation company: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessDonationCompanyUpdatedSP, nil)
	} else {
		d.logger.Infof("donation company update request submitted successfully", updatedDonationCompany)
		localization.SendSuccessResponse(w, localization.SuccessDonationCompanyUpdatedRequestSent, nil)
	}
}

// AccountLookup godoc
//
//	@Summary		Account lookup for donation company
//	@Description	Lookup account information by account number for donation company
//	@Tags			Donation Company
//	@Accept			json
//	@Produce		json
//	@Param			account_number	path		string									true	"Account number"
//	@Success		200				{object}	localization.StandardResponse{data=object}	"Account information retrieved successfully"
//	@Failure		400				{object}	localization.StandardResponse{data=nil}		"Bad request"
//	@Failure		404				{object}	localization.StandardResponse{data=nil}		"Account not found"
//	@Failure		500				{object}	localization.StandardResponse{data=nil}		"Internal server error"
//	@Security		BearerAuth
//	@Router			/account_lookup/{account_number} [get]
func (d *donationCompanyAdapter) AccountLookup(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "donationCompanyAccountLookup", "handler", "donationCompany")
	defer span.End()
	accountNumber := chi.URLParam(r, "account_number")
	if accountNumber == "" {
		d.logger.Errorf("account nmumber is required to fetch one")
		localization.SendErrorResponse(w, localization.ErrorAccountNumberRequired, nil, nil)
		return
	}

	if !local_util.IsValidCBEAccountNumber(accountNumber) {
		d.logger.Errorf("account nmumber is not valid")
		localization.SendBadRequestResponse(w, localization.ErrorAccountNumberNotValid.Message)
		return
	}
	span.SetAttributes(attribute.String("donation_company.account_number", accountNumber))

	accountInfo, err := d.donationCompanyApp.AccountLookup(ctx, accountNumber)
	if err != nil {
		span.RecordError(err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessAccountInfoFetched, accountInfo)
}

// EnableDonationCompany godoc
//	@Summary		Enable a donation company
//	@Description	Enable a donation company by ID
//	@Tags			Donation Company
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Donation company ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Donation Company enable request sent successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}	"Donation not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/donation_company/enable/{id} [patch]

func (d *donationCompanyAdapter) EnableDonationCompany(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "enableDonationCompany", "handler", "donationCompany")
	defer span.End()
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	id := chi.URLParam(r, "id")
	if id == "" {
		d.logger.Errorf("donation company ID is required for diable")
		localization.SendErrorResponse(w, localization.ErrorDonationCompanyIdRequired, nil, nil)
		return
	}

	span.SetAttributes(attribute.String("donation_company.id", id))
	if err := d.donationCompanyApp.EnableDonationCompany(ctx, id); err != nil {
		span.RecordError(err)
		d.logger.Errorf("failed to enable donation company: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessDonationCompanyEnabledSP, nil)
	} else {
		localization.SendSuccessResponse(w, localization.SuccessDonationCompanyEnableRequestSent, nil)

	}
}

// DisableDonationCompany godoc
//
//	@Summary		Disable a donation company
//	@Description	Disable a donation company by ID
//	@Tags			Donation Company
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Donation company ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Donation company disable request sent successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}	"Donation company not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/donation_company/disable/{id} [patch]
func (d *donationCompanyAdapter) DisableDonationCompany(w http.ResponseWriter, r *http.Request) {
	ctx, span := local_util.TraceLogger(r.Context(), "handler", "disableDonationCompany", "handler", "donationCompany")
	defer span.End()
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	id := chi.URLParam(r, "id")
	if id == "" {
		d.logger.Errorf("donation company ID is required for disable")
		localization.SendErrorResponse(w, localization.ErrorDonationCompanyIdRequired, nil, nil)
		return
	}

	span.SetAttributes(attribute.String("donation_company.id", id))
	if err := d.donationCompanyApp.DisableDonationCompany(ctx, id); err != nil {
		span.RecordError(err)
		d.logger.Errorf("failed to disable donation company: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {
		localization.SendSuccessResponse(w, localization.SuccessDonationCompanyDisabledSP, nil)
	} else {
		localization.SendSuccessResponse(w, localization.SuccessDonationCompanyDisableRequestSent, nil)
	}
}
