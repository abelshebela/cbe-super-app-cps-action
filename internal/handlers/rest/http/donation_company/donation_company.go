package donation_company

import (
	"cbe-super-app-cps-action/internal/constants/dto/donation_company"
	inbound "cbe-super-app-cps-action/internal/constants/interfaces/donation_company"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"net/http"

	"cbe-super-app-cps-action/internal/handlers/rest/http/donation_company/core"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
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
// @Summary List donation companies
// @Description Retrieve donation companies with pagination and optional search
// @Tags Donation Company
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param search query string false "Search term"
// @Success 200 {object} localization.StandardResponse{data=paginatedDonationCompanyListResponse} "Donation companies retrieved successfully"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Internal server error"
// @Security BearerAuth
// @Router /donation_company [get]
func (d *donationCompanyAdapter) FetchDonationCompany(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)
	if filterParams.Page < 0 || filterParams.PerPage < 0 {
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}
	ctx := r.Context()

	donationCompanies, err := d.donationCompanyApp.FetchDonationCompany(ctx, filterParams)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessDonationCompaniesFetched, donationCompanies)
}

// FetchDonationCompanyByID godoc
// @Summary Get donation company by ID
// @Description Retrieve a donation company's details by ID
// @Tags Donation Company
// @Accept json
// @Produce json
// @Param id path string true "Donation Company ID"
// @Success 200 {object} localization.StandardResponse{data=donation_company.DonationCompanyListResponse} "Donation company retrieved successfully"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Bad request"
// @Failure 404 {object} localization.StandardResponse{data=nil} "Not found"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Internal server error"
// @Security BearerAuth
// @Router /donation_company/{id} [get]
func (d *donationCompanyAdapter) FetchDonationCompanyByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		d.logger.Errorf("donation company ID is required to fetch one")
		localization.SendErrorResponse(w, localization.ErrorDonationCompanyIdRequired, nil, nil)
		return
	}
	ctx := r.Context()

	donationCompany, err := d.donationCompanyApp.FetchDonationCompanyByID(ctx, id)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessDonationCompanyFetched, donationCompany)
}

// CreateDonationCompany godoc
// @Summary Create a new donation company
// @Description Create a new donation company with the provided information
// @Tags Donation Company
// @Accept multipart/form-data
// @Produce json
// @Param company_name formData string true "Company name"
// @Param company_logo formData file true "Company logo image file"
// @Param account_number formData string true "Account number"
// @Success 200 {object} localization.StandardResponse{data=nil} "Donation company creation request sent successfully"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Bad request"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Internal server error"
// @Security BearerAuth
// @Router /donation_company [post]
func (d *donationCompanyAdapter) CreateDonationCompany(w http.ResponseWriter, r *http.Request) {
	req, err := core.ParseRequestFromMultipartForm(r, true)
	if err != nil {
		d.logger.Errorf("failed to parse donation company request from multipart form: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		d.logger.Errorf("donation company request validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := d.donationCompanyApp.CreateDonationCompany(r.Context(), req); err != nil {
		d.logger.Errorf("failed to create donation company: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	d.logger.Infof("donation company creation request submitted successfully")
	localization.SendSuccessResponse(w, localization.SuccessDonationCompanyCreateRequestSent, nil)
}

// UpdateDonationCompany godoc
// @Summary Update a donation company
// @Description Update a donation company with the provided information
// @Tags Donation Company
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Donation Company ID"
// @Param company_name formData string false "Company name"
// @Param company_logo formData file false "Company logo image file"
// @Param account_number formData string false "Account number"
// @Success 200 {object} localization.StandardResponse{data=nil} "Donation company updated successfully"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Bad request"
// @Failure 404 {object} localization.StandardResponse{data=nil} "Not found"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Internal server error"
// @Security BearerAuth
// @Router /donation_company/{id} [patch]
func (d *donationCompanyAdapter) UpdateDonationCompany(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		d.logger.Errorf("donation company ID is required for update")
		localization.SendErrorResponse(w, localization.ErrorDonationCompanyIdRequired, nil, nil)
		return
	}

	req, err := core.ParseRequestFromMultipartForm(r, false)
	if err != nil {
		d.logger.Errorf("failed to parse donation company update request from multipart form: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := req.ValidateForUpdate(); err != nil {
		d.logger.Errorf("donation company update request validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	updatedDonationCompany, err := d.donationCompanyApp.UpdateDonationCompany(r.Context(), id, req)
	if err != nil {
		d.logger.Errorf("failed to update donation company: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	d.logger.Infof("donation company update request submitted successfully", updatedDonationCompany)
	localization.SendSuccessResponse(w, localization.SuccessDonationCompanyUpdated, nil)
}

func (d *donationCompanyAdapter) AccountLookup(w http.ResponseWriter, r *http.Request){
	accountNumber := chi.URLParam(r, "account_number")
	if accountNumber == "" {
		d.logger.Errorf("account nmumber is required to fetch one")
		localization.SendErrorResponse(w, localization.ErrorAccountNumberRequired, nil, nil)
		return
	}
	ctx := r.Context()

	accountInfo, err := d.donationCompanyApp.AccountLookup(ctx, accountNumber)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessAccountInfoFetched, accountInfo)
}


// EnableDonationCompany godoc
// @Summary Enable a donation company
// @Description Enable a donation company by ID
// @Tags Donation Company
// @Accept json
// @Produce json
// @Param id path string true "Donation company ID"
// @Success 200 {object} localization.StandardResponse{data=nil} "Donation Company enable request sent successfully"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Bad request"
// @Failure 404 {object} localization.StandardResponse{data=nil} "Donation not found"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Internal server error"
// @Security BearerAuth
// @Router /donation_company/enable/{id} [patch]

func (d *donationCompanyAdapter) EnableDonationCompany(w http.ResponseWriter, r *http.Request) {
id := chi.URLParam(r, "id")
	if id == "" {
		d.logger.Errorf("donation company ID is required for diable")
		localization.SendErrorResponse(w, localization.ErrorDonationCompanyIdRequired, nil, nil)
		return
	}

	if err:= d.donationCompanyApp.EnableDonationCompany(r.Context(),id);err!= nil{
			d.logger.Errorf("failed to enable donation company: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
		localization.SendSuccessResponse(w, localization.SuccessDonationCompanyEnableRequestSent, nil)

}

// disableDonationCompany godoc
// @Summary Enable a donation company
// @Description Enable a donation company by ID
// @Tags Donation company
// @Accept json
// @Produce json
// @Param id path string true "Donation company ID"
// @Success 200 {object} localization.StandardResponse{data=nil} "Donation company enable request sent successfully"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Bad request"
// @Failure 404 {object} localization.StandardResponse{data=nil} "Donation not found"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Internal server error"
// @Security BearerAuth
// @Router /donation_company/enable/{id} [patch]

func (d *donationCompanyAdapter) DisableDonationCompany(w http.ResponseWriter, r *http.Request) {
id := chi.URLParam(r, "id")
	if id == "" {
		d.logger.Errorf("donation company ID is required for disable")
		localization.SendErrorResponse(w, localization.ErrorDonationCompanyIdRequired, nil, nil)
		return
	}

	if err:= d.donationCompanyApp.DisableDonationCompany(r.Context(),id);err!= nil{
			d.logger.Errorf("failed to disable donation company: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
		localization.SendSuccessResponse(w, localization.SuccessDonationCompanyDisableRequestSent, nil)

}