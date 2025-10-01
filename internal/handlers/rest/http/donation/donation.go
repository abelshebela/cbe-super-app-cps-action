package donation

import (
	dto "cbe-super-app-cps-action/internal/constants/dto/donation"
	donation_interface "cbe-super-app-cps-action/internal/constants/interfaces/donation"
	"cbe-super-app-cps-action/internal/constants/localization"
	core "cbe-super-app-cps-action/internal/handlers/rest/http/donation/core"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"encoding/json"
	"net/http"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

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

func (d *donationAdapter) CreateDonation(w http.ResponseWriter, r *http.Request) {
	req, err := core.ParseRequestFromMultipartForm(r, true)
	if err != nil {
		d.logger.Errorf("failed to parse request from multipart form: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		d.logger.Errorf("request validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := d.donationApp.CreateDonation(r.Context(), req); err != nil {
		d.logger.Errorf("failed to create donation: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessDonationCreateRequestSent, nil)
}

func (d *donationAdapter) UpdateDonation(w http.ResponseWriter, r *http.Request) {
	id := core.ExtractIDFromURL(r)
	if id == "" {
		d.logger.Errorf("donation ID is required")
		localization.SendBadRequestResponse(w, "donation ID is required")
		return
	}

	req, err := core.ParseRequestFromMultipartForm(r, false)
	if err != nil {
		d.logger.Errorf("failed to parse request from multipart form: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := core.ValidateForUpdate(req); err != nil {
		d.logger.Errorf("request validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	err = d.donationApp.UpdateDonation(r.Context(), id, req)
	if err != nil {
		d.logger.Errorf("failed to update donation: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessDonationUpdateRequestSent, nil)
}

func (d *donationAdapter) FetchDonation(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)
	if filterParams.Page < 0 || filterParams.PerPage < 0 {
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}

	donations, err := d.donationApp.FetchDonation(r.Context(), filterParams)
	if err != nil {
		d.logger.Errorf("failed to fetch donations: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessDonationFetched, donations)
}

func (d *donationAdapter) FetchDonationByID(w http.ResponseWriter, r *http.Request) {
	id := core.ExtractIDFromURL(r)
	if id == "" {
		d.logger.Errorf("donation ID is required")
		localization.SendBadRequestResponse(w, "donation ID is required")
		return
	}

	donation, err := d.donationApp.FetchDonationByID(r.Context(), id)
	if err != nil {
		d.logger.Errorf("failed to fetch donation by ID: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessDonationFetched, donation)
}

func (d *donationAdapter) UpdateDonationImage(w http.ResponseWriter, r *http.Request) {
	id := core.ExtractIDFromURL(r)
	if id == "" {
		d.logger.Errorf("donation ID is required")
		localization.SendBadRequestResponse(w, "donation ID is required")
		return
	}

	req, err := core.ParseImageUpdateRequestFromMultipartForm(r)
	if err != nil {
		d.logger.Errorf("failed to parse image update request from multipart form: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := core.ValidateImageUpdateRequest(req); err != nil {
		d.logger.Errorf("image update validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := d.donationApp.UpdateDonationImage(r.Context(), id, req); err != nil {
		d.logger.Errorf("failed to update donation image: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessDonationImageUpdateRequestSent, nil)
}

func (d *donationAdapter) DeleteDonationImage(w http.ResponseWriter, r *http.Request) {
	id := core.ExtractIDFromURL(r)
	if id == "" {
		d.logger.Errorf("donation ID is required")
		localization.SendBadRequestResponse(w, "donation ID is required")
		return
	}

	var req dto.DonationImageDeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		d.logger.Errorf("Failed to decode JSON request: %v", err)
		localization.SendBadRequestResponse(w, "invalid request body")
		return
	}

	if req.ImageID == "" {
		d.logger.Errorf("image ID is required")
		localization.SendBadRequestResponse(w, "image ID is required")
		return
	}

	if err := d.donationApp.DeleteDonationImage(r.Context(), id, req.ImageID); err != nil {
		d.logger.Errorf("failed to delete donation image: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessDonationImageDeleteRequestSent, nil)
}

func (d *donationAdapter) AddDonationImage(w http.ResponseWriter, r *http.Request) {
	id := core.ExtractIDFromURL(r)
	if id == "" {
		d.logger.Errorf("donation ID is required")
		localization.SendBadRequestResponse(w, "donation ID is required")
		return
	}

	req, err := core.ParseRequestFromMultipartForm(r, false)
	if err != nil {
		d.logger.Errorf("failed to parse request from multipart form: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := core.ValidateImageAdd(req); err != nil {
		d.logger.Errorf("image validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := d.donationApp.AddDonationImage(r.Context(), id, req); err != nil {
		d.logger.Errorf("failed to add donation image: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessDonationImageAddRequestSent, nil)
}

func (d *donationAdapter) EnableDonation(w http.ResponseWriter, r *http.Request) {
	id := core.ExtractIDFromURL(r)
	if id == "" {
		d.logger.Errorf("donation ID is required")
		localization.SendBadRequestResponse(w, "donation ID is required")
		return
	}

	if err := d.donationApp.EnableDonation(r.Context(), id); err != nil {
		d.logger.Errorf("failed to enable donation: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessDonationEnableRequestSent, nil)
}

func (d *donationAdapter) DisableDonation(w http.ResponseWriter, r *http.Request) {
	id := core.ExtractIDFromURL(r)
	if id == "" {
		d.logger.Errorf("donation ID is required")
		localization.SendBadRequestResponse(w, "donation ID is required")
		return
	}

	if err := d.donationApp.DisableDonation(r.Context(), id); err != nil {
		d.logger.Errorf("failed to disable donation: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessDonationDisableRequestSent, nil)
}
