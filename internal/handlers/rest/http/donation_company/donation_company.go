package donation_company

import (
	inbound "cbe-super-app-cps-action/internal/constants/interfaces/donation_company"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	"net/http"

	"cbe-super-app-cps-action/internal/handlers/rest/http/donation_company/core"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

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
