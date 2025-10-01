package donation_category

import (
	inbound "cbe-super-app-cps-action/internal/constants/interfaces/donation_category"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	"net/http"

	"cbe-super-app-cps-action/internal/handlers/rest/http/donation_category/core"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type donationCategoryAdapter struct {
	logger              utils.Logger
	donationCategoryApp service.DonationCategoryService
}

func InitDonationCategoryAdapter(donationCategoryApp service.DonationCategoryService, logger utils.Logger) inbound.DonationCategoryAdapter {
	return &donationCategoryAdapter{
		logger:              logger,
		donationCategoryApp: donationCategoryApp,
	}
}

func (d *donationCategoryAdapter) FetchDonationCategory(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)
	if filterParams.Page < 0 || filterParams.PerPage < 0 {
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}
	ctx := r.Context()

	donationCategories, err := d.donationCategoryApp.FetchDonationCategory(ctx, filterParams)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessDonationCategoriesFetched, donationCategories)

}

func (d *donationCategoryAdapter) FetchDonationCategoryByID(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	if id == "" {
		d.logger.Errorf("donation category ID is required to fetch one")
		localization.SendErrorResponse(w, localization.ErrorDonationCategoryIDRequired, nil, nil)
		return
	}
	ctx := r.Context()

	serviceFeeDetails, err := d.donationCategoryApp.FetchDonationCategoryByID(ctx, id)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessDonationCategoryFetched, serviceFeeDetails)

}

func (d *donationCategoryAdapter) CreateDonationCategory(w http.ResponseWriter, r *http.Request) {
	req, err := core.ParseRequestFromMultipartForm(r, true)
	if err != nil {
		d.logger.Errorf("failed to parse event request from multipart form: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		d.logger.Errorf("event request validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := d.donationCategoryApp.CreateDonationCategory(r.Context(), req); err != nil {
		d.logger.Errorf("failed to create event: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	d.logger.Infof("event creation request submitted successfully")
	localization.SendSuccessResponse(w, localization.SuccessDonationCategoryCreateRequestSent, nil)
}

func (d *donationCategoryAdapter) UpdateDonationCategory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		d.logger.Errorf("donation category ID is required for update")
		localization.SendErrorResponse(w, localization.ErrorDonationCategoryIDRequired, nil, nil)
		return
	}

	req, err := core.ParseRequestFromMultipartForm(r, false)
	if err != nil {
		d.logger.Errorf("failed to parse donation category update request from multipart form: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := req.ValidateForUpdate(); err != nil {
		d.logger.Errorf("donation category update request validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	updatedDonationCategory, err := d.donationCategoryApp.UpdateDonationCategory(r.Context(), id, req)
	if err != nil {
		d.logger.Errorf("failed to update donation category: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	d.logger.Infof("donation category update request submitted successfully", updatedDonationCategory)
	localization.SendSuccessResponse(w, localization.SuccessDonationCategoryUpdated, nil)
}
