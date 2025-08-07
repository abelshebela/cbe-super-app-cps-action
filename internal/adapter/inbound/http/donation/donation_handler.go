package donation

import (
	"net/http"

	"mime/multipart"

	donation_application "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/donation"
	dto "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	cps_entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	c "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)

type DonationHttpStore struct {
	Application donation_application.DonationAbstract
	logger      c.Logger
}

func NewDonationHTTPHandler(app donation_application.DonationAbstract, logger c.Logger) *DonationHttpStore {
	return &DonationHttpStore{
		Application: app,
		logger:      logger,
	}
}

func (h *DonationHttpStore) extractUserAndMaker(w http.ResponseWriter, r *http.Request) (cps_entities.User, bool) {
	userContext := ctx_util.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		utils.SendErrorResponse(w, utils.IncompleteUserInfo, 0, nil)
		return cps_entities.User{}, false
	}
	maker := cps_entities.User{
		UserCode:    userContext.UserID,
		FullName:    userContext.FullName,
		PhoneNumber: userContext.PhoneNumber,
		Department:  userContext.Department,
	}
	return maker, true
}

func (h *DonationHttpStore) extractID(w http.ResponseWriter, r *http.Request) (string, bool) {
	id, ok := utils.GetParam(r, "id")
	if !ok {
		h.logger.Errorf("missing or invalid parameter 'id'")
		utils.SendErrorResponse(w, utils.InvalidInputParameters, 0, nil)
		return "", false
	}
	return id, true
}

func (h *DonationHttpStore) CreateDonationCategory(w http.ResponseWriter, r *http.Request) {
	maker, ok := h.extractUserAndMaker(w, r)
	if !ok {
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.SendErrorResponse(w, "INVALID_MULTIPART_FORM", http.StatusBadRequest, nil)
		return
	}

	categoryName := r.FormValue("category_name")
	fileHeaders := r.MultipartForm.File["donation_icon"]
	if len(fileHeaders) == 0 {
		utils.SendErrorResponse(w, "donation_icon is required", http.StatusBadRequest, nil)
		return
	}
	fileHeader := fileHeaders[0]

	donationRequest := dto.DonationCategoryRequest{
		CategoryName: categoryName,
		Icon:         fileHeader,
	}

	if err := donationRequest.Validate(); err != nil {
		h.logger.Errorf("validation failed: %v", err)
		utils.SendErrorResponse(w, "INVALID_INPUT", http.StatusBadRequest, map[string]interface{}{"errors": err})
		return
	}

	if err := h.Application.CreateDonationCategory(r.Context(), donationRequest, maker); err != nil {
		h.logger.Errorf("failed to create donation category: %v", err)
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	utils.WriteSuccessResponse(w, nil, "Donation category created successfully")
}

func (h *DonationHttpStore) FetchDonationCategory(w http.ResponseWriter, r *http.Request) {
	filterParams := common_util.ExtractMongoFilterParams(r)

	paginatedResponse, err := h.Application.FetchDonationCategory(r.Context(), filterParams)
	if err != nil {
		h.logger.Errorf("failed to fetch donation categories: %v", err)
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	utils.WriteSuccessResponse(w, paginatedResponse, "Donation categories fetched successfully")
}

func (h *DonationHttpStore) FetchDonationCategoryByID(w http.ResponseWriter, r *http.Request) {
	id, ok := h.extractID(w, r)
	if !ok {
		return
	}

	donationCategory, err := h.Application.FetchDonationCategoryByID(r.Context(), id)
	if err != nil {
		h.logger.Errorf("failed to fetch donation category by ID: %v", err)
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	utils.WriteSuccessResponse(w, donationCategory, "Donation category fetched successfully")
}

func (h *DonationHttpStore) UpdateDonationCategory(w http.ResponseWriter, r *http.Request) {
	id, ok := h.extractID(w, r)
	if !ok {
		return
	}

	maker, ok := h.extractUserAndMaker(w, r)
	if !ok {
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.SendErrorResponse(w, "INVALID_MULTIPART_FORM", http.StatusBadRequest, nil)
		return
	}

	categoryName := r.FormValue("category_name")
	fileHeaders := r.MultipartForm.File["donation_icon"]

	var fileHeader *multipart.FileHeader
	if len(fileHeaders) > 0 {
		fileHeader = fileHeaders[0]
	}

	donationRequest := dto.DonationCategoryRequest{
		CategoryName: categoryName,
		Icon:         fileHeader,
	}

	if err := donationRequest.ValidateForUpdate(); err != nil {
		h.logger.Errorf("validation failed: %v", err)
		utils.SendErrorResponse(w, "INVALID_INPUT", http.StatusBadRequest, map[string]interface{}{"errors": err})
		return
	}

	if err := h.Application.UpdateDonationCategory(r.Context(), id, donationRequest, maker); err != nil {
		h.logger.Errorf("failed to update donation category: %v", err)
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	utils.WriteSuccessResponse(w, nil, "Donation category updated successfully")
}
