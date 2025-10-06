package donation_category

import (
	"cbe-super-app-cps-action/internal/constants/dto/donation_category"
	inbound "cbe-super-app-cps-action/internal/constants/interfaces/donation_category"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"net/http"

	"cbe-super-app-cps-action/internal/handlers/rest/http/donation_category/core"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type paginatedDonationCategoryListResponse types.PaginatedResponse[[]donation_category.DonationCategoryListResponse]
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

// FetchDonationCategory godoc
// @Summary List donation categories
// @Description Retrieve donation categories with pagination and optional search
// @Tags Donation Category
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param search query string false "Search term"
// @Success 200 {object} localization.StandardResponse{data=paginatedDonationCategoryListResponse} "Donation categories fetched successfully"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Internal server error"
// @Security BearerAuth
// @Router /donation_category [get]
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

// FetchDonationCategoryByID godoc
// @Summary Get donation category by ID
// @Description Retrieve a donation category's details by ID
// @Tags Donation Category
// @Accept json
// @Produce json
// @Param id path string true "Donation Category ID"
// @Success 200 {object} localization.StandardResponse{data=donation_category.DonationCategoryListResponse} "Donation category retrieved successfully"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Bad request"
// @Failure 404 {object} localization.StandardResponse{data=nil} "Not found"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Internal server error"
// @Security BearerAuth
// @Router /donation_category/{id} [get]
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

// CreateDonationCategory godoc
// @Summary Create a new donation category
// @Description Create a new donation category with the provided information
// @Tags Donation Category
// @Accept multipart/form-data
// @Produce json
// @Param category_name formData string true "Category name"
// @Param donation_icon formData file true "Donation icon image file"
// @Success 200 {object} localization.StandardResponse{data=nil} "Donation category creation request sent successfully"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Bad request"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Internal server error"
// @Security BearerAuth
// @Router /donation_category [post]
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

// UpdateDonationCategory godoc
// @Summary Update a donation category
// @Description Update a donation category with the provided information
// @Tags Donation Category
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Donation Category ID"
// @Param category_name formData string false "Category name"
// @Param donation_icon formData file false "Donation icon image file"
// @Success 200 {object} localization.StandardResponse{data=nil} "Donation category updated successfully"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Bad request"
// @Failure 404 {object} localization.StandardResponse{data=nil} "Not found"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Internal server error"
// @Security BearerAuth
// @Router /donation_category/{id} [patch]
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
