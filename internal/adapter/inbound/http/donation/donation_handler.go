package donation

import (
	"encoding/json"
	"fmt"
	"net/http"

	"mime/multipart"

	"time"

	"strconv"

	donation_application "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/donation"
	dto "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	cps_entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"

	// common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	c "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func parseDonationAmount(amountStr string) (int, error) {
	if amountStr == "" {
		return 0, fmt.Errorf("donation amount is required")
	}

	cleanedAmount := utils.CleanAmountString(amountStr)

	amount, err := strconv.ParseFloat(cleanedAmount, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid donation amount format")
	}

	return int(amount), nil
}

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

	utils.WriteSuccessResponse(w, nil, "Donation category create request sent  successfully")
}

func (h *DonationHttpStore) FetchDonationCategory(w http.ResponseWriter, r *http.Request) {
	filterParams := utils.ExtractMongoFilterParams(r)

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

	// Get form values - these can be empty strings for partial updates
	categoryName := r.FormValue("category_name")
	fileHeaders := r.MultipartForm.File["donation_icon"]

	var fileHeader *multipart.FileHeader
	if len(fileHeaders) > 0 {
		fileHeader = fileHeaders[0]
	}

	// Debug logging to understand what's being received
	h.logger.Infof("UpdateDonationCategory - ID: %s, CategoryName: '%s', HasIcon: %v", id, categoryName, fileHeader != nil)

	// Check if at least one field is provided for the update
	if categoryName == "" && fileHeader == nil {
		utils.SendErrorResponse(w, "at least one field must be provided for update", http.StatusBadRequest, nil)
		return
	}

	// Create request with provided values (empty strings are allowed for partial updates)
	donationRequest := dto.DonationCategoryRequest{
		CategoryName: categoryName, // Can be empty string for partial updates
		Icon:         fileHeader,   // Can be nil for partial updates
	}

	// Validate the request (validation now allows partial updates)
	if err := donationRequest.ValidateForUpdate(); err != nil {
		h.logger.Errorf("validation failed: %v", err)
		h.logger.Errorf("validation details - CategoryName: '%s', Icon: %v", donationRequest.CategoryName, donationRequest.Icon)
		utils.SendErrorResponse(w, "INVALID_INPUT", http.StatusBadRequest, map[string]interface{}{"errors": err})
		return
	}

	if err := h.Application.UpdateDonationCategory(r.Context(), id, donationRequest, maker); err != nil {
		h.logger.Errorf("failed to update donation category: %v", err)
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	utils.WriteSuccessResponse(w, nil, "Donation category update request sent successfully")
}

func (h *DonationHttpStore) CreateDonationCompany(w http.ResponseWriter, r *http.Request) {
	maker, ok := h.extractUserAndMaker(w, r)
	if !ok {
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.SendErrorResponse(w, "INVALID_MULTIPART_FORM", http.StatusBadRequest, nil)
		return
	}

	companyName := r.FormValue("company_name")
	accountNumber := r.FormValue("account_number")
	fileHeaders := r.MultipartForm.File["company_logo"]
	if len(fileHeaders) == 0 {
		utils.SendErrorResponse(w, "company_logo is required", http.StatusBadRequest, nil)
		return
	}
	fileHeader := fileHeaders[0]

	companyRequest := dto.DonationCompanyRequest{
		CompanyName:   companyName,
		CompanyLogo:   fileHeader,
		AccountNumber: accountNumber,
	}

	if err := companyRequest.Validate(); err != nil {
		h.logger.Errorf("validation failed: %v", err)
		utils.SendErrorResponse(w, "INVALID_INPUT", http.StatusBadRequest, map[string]interface{}{"errors": err})
		return
	}

	if err := h.Application.CreateDonationCompany(r.Context(), companyRequest, maker); err != nil {
		h.logger.Errorf("failed to create donation company: %v", err)
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	utils.WriteSuccessResponse(w, nil, "Donation company create request sent  successfully")
}

func (h *DonationHttpStore) FetchDonationCompany(w http.ResponseWriter, r *http.Request) {
	filterParams := utils.ExtractMongoFilterParams(r)

	paginatedResponse, err := h.Application.FetchDonationCompany(r.Context(), filterParams)
	if err != nil {
		h.logger.Errorf("failed to fetch donation companies: %v", err)
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	utils.WriteSuccessResponse(w, paginatedResponse, "Donation companies fetched successfully")
}

func (h *DonationHttpStore) FetchDonationCompanyByID(w http.ResponseWriter, r *http.Request) {
	id, ok := h.extractID(w, r)
	if !ok {
		return
	}

	donationCompany, err := h.Application.FetchDonationCompanyByID(r.Context(), id)
	if err != nil {
		h.logger.Errorf("failed to fetch donation company by ID: %v", err)
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	utils.WriteSuccessResponse(w, donationCompany, "Donation company fetched successfully")
}

func (h *DonationHttpStore) UpdateDonationCompany(w http.ResponseWriter, r *http.Request) {
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

	// Get form values - these can be empty strings for partial updates
	companyName := r.FormValue("company_name")
	accountNumber := r.FormValue("account_number")
	fileHeaders := r.MultipartForm.File["company_logo"]

	var fileHeader *multipart.FileHeader
	if len(fileHeaders) > 0 {
		fileHeader = fileHeaders[0]
	}

	// Debug logging to understand what's being received
	h.logger.Infof("UpdateDonationCompany - ID: %s, CompanyName: '%s', AccountNumber: '%s', HasLogo: %v", id, companyName, accountNumber, fileHeader != nil)

	// Check if at least one field is provided for the update
	if companyName == "" && accountNumber == "" && fileHeader == nil {
		utils.SendErrorResponse(w, "at least one field must be provided for update", http.StatusBadRequest, nil)
		return
	}

	// Create request with provided values (empty strings are allowed for partial updates)
	companyRequest := dto.DonationCompanyRequest{
		CompanyName:   companyName,   // Can be empty string for partial updates
		CompanyLogo:   fileHeader,    // Can be nil for partial updates
		AccountNumber: accountNumber, // Can be empty string for partial updates
	}

	// Validate the request (validation now allows partial updates)
	if err := companyRequest.ValidateForUpdate(); err != nil {
		h.logger.Errorf("validation failed: %v", err)
		h.logger.Errorf("validation details - CompanyName: '%s', AccountNumber: '%s', HasLogo: %v", companyRequest.CompanyName, companyRequest.AccountNumber, companyRequest.CompanyLogo != nil)
		utils.SendErrorResponse(w, "INVALID_INPUT", http.StatusBadRequest, map[string]interface{}{"errors": err})
		return
	}

	if err := h.Application.UpdateDonationCompany(r.Context(), id, companyRequest, maker); err != nil {
		h.logger.Errorf("failed to update donation company: %v", err)
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	utils.WriteSuccessResponse(w, nil, "Donation company update request sent successfully")
}

func (h *DonationHttpStore) CreateDonation(w http.ResponseWriter, r *http.Request) {
	maker, ok := h.extractUserAndMaker(w, r)
	if !ok {
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.SendErrorResponse(w, "INVALID_MULTIPART_FORM", http.StatusBadRequest, nil)
		return
	}

	companyID := r.FormValue("company_id")
	categoryID := r.FormValue("category_id")
	title := r.FormValue("title")
	isFeatured := r.FormValue("is_featured") == "true"
	targetStr := r.FormValue("target")
	donationDescription := r.FormValue("donation_description")
	endDateStr := r.FormValue("end_date")
	startDateStr := r.FormValue("start_date")

	fileHeaders := r.MultipartForm.File["donation_images"]
	coverImageFile := r.MultipartForm.File["cover_image"]

	if len(fileHeaders) == 0 {
		utils.SendErrorResponse(w, "DONATION_IMAGES_REQUIRED", http.StatusBadRequest, nil)
		return
	}
		if len(coverImageFile) == 0 {
		utils.SendErrorResponse(w, "COVER_IMAGES_REQUIRED", http.StatusBadRequest, nil)
		return
	}

	target, err := parseDonationAmount(targetStr)
	if err != nil {
		utils.SendErrorResponse(w, "INVALID_DONATION_AMOUNT", http.StatusBadRequest, nil)
		return
	}

	var endDate, startDate time.Time
	if endDateStr != "" {
		endDate, err = utils.ParseDateString(endDateStr)
		if err != nil {
			utils.SendErrorResponse(w, "INVALID_END_DATE_FORMAT", http.StatusBadRequest, nil)
			return
		}
	}

	if startDateStr != "" {
		startDate, err = utils.ParseDateString(startDateStr)
		if err != nil {
			utils.SendErrorResponse(w, "INVALID_START_DATE_FORMAT", http.StatusBadRequest, nil)
			return
		}
	}

	var coverImage *multipart.FileHeader
	if len(coverImageFile) > 0 {
		coverImage = coverImageFile[0]
	}

	donationRequest := dto.DonationRequest{
		CompanyID:           companyID,
		CategoryID:          categoryID,
		Title:               title,
		IsFeatured:          isFeatured,
		Target:              target,
		DonationDescription: donationDescription,
		EndDate:             endDate,
		StartDate:           startDate,
		DonationImages:      fileHeaders,
		CoverImage:          coverImage,
	}

	if err := donationRequest.Validate(); err != nil {
		h.logger.Errorf("validation failed: %v", err)
		utils.SendErrorResponse(w, "INVALID_INPUT", http.StatusBadRequest, map[string]interface{}{"errors": err})
		return
	}

	if err := h.Application.CreateDonation(r.Context(), donationRequest, maker); err != nil {
		h.logger.Errorf("failed to create donation: %v", err)
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	utils.WriteSuccessResponse(w, nil, "Donation create request sent successfully")
}

func (h *DonationHttpStore) FetchDonation(w http.ResponseWriter, r *http.Request) {
	filterParams := utils.ExtractMongoFilterParams(r)

	paginatedResponse, err := h.Application.FetchDonation(r.Context(), filterParams)
	if err != nil {
		h.logger.Errorf("failed to fetch donations: %v", err)
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	utils.WriteSuccessResponse(w, paginatedResponse, "Donations fetched successfully")
}

func (h *DonationHttpStore) FetchDonationByID(w http.ResponseWriter, r *http.Request) {
	id, ok := h.extractID(w, r)
	if !ok {
		return
	}

	donation, err := h.Application.FetchDonationByID(r.Context(), id)
	if err != nil {
		h.logger.Errorf("failed to fetch donation by ID: %v", err)
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	utils.WriteSuccessResponse(w, donation, "Donation fetched successfully")
}

func (h *DonationHttpStore) UpdateDonation(w http.ResponseWriter, r *http.Request) {
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

	companyID := r.FormValue("company_id")
	categoryID := r.FormValue("category_id")
	title := r.FormValue("title")
	isFeatured := r.FormValue("is_featured") == "true"
	targetStr := r.FormValue("target")
	donationDescription := r.FormValue("donation_description")
	endDateStr := r.FormValue("end_date")
	startDateStr := r.FormValue("start_date")

	coverImageFile := r.MultipartForm.File["cover_image"]

	var target int
	var err error
	if targetStr != "" {
		target, err = parseDonationAmount(targetStr)
		if err != nil {
			utils.SendErrorResponse(w, "INVALID_DONATION_AMOUNT", http.StatusBadRequest, nil)
			return
		}
	}

	var endDate, startDate time.Time
	if endDateStr != "" {
		endDate, err = utils.ParseDateString(endDateStr)
		if err != nil {
			utils.SendErrorResponse(w, "INVALID_END_DATE_FORMAT", http.StatusBadRequest, nil)
			return
		}
	}

	if startDateStr != "" {
		startDate, err = utils.ParseDateString(startDateStr)
		if err != nil {
			utils.SendErrorResponse(w, "INVALID_START_DATE_FORMAT", http.StatusBadRequest, nil)
			return
		}
	}

	var coverImage *multipart.FileHeader
	if len(coverImageFile) > 0 {
		coverImage = coverImageFile[0]
	}

	donationRequest := dto.DonationRequest{
		CompanyID:           companyID,
		CategoryID:          categoryID,
		Title:               title,
		IsFeatured:          isFeatured,
		Target:              target,
		DonationDescription: donationDescription,
		EndDate:             endDate,
		StartDate:           startDate,
		CoverImage:          coverImage,
	}

	if err := donationRequest.ValidateForUpdate(); err != nil {
		h.logger.Errorf("validation failed: %v", err)
		utils.SendErrorResponse(w, "INVALID_INPUT", http.StatusBadRequest, map[string]interface{}{"errors": err})
		return
	}

	if err := h.Application.UpdateDonation(r.Context(), id, donationRequest, maker); err != nil {
		h.logger.Errorf("failed to update donation: %v", err)
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	utils.WriteSuccessResponse(w, nil, "Donation update request sent successfully")
}

func (h *DonationHttpStore) UpdateDonationImage(w http.ResponseWriter, r *http.Request) {
	donationID, ok := h.extractID(w, r)
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

	imageID := r.FormValue("image_id")
	if imageID == "" {
		utils.SendErrorResponse(w, "IMAGE_ID_REQUIRED", http.StatusBadRequest, nil)
		return
	}

	fileHeaders := r.MultipartForm.File["donation_images"]
	if len(fileHeaders) == 0 {
		utils.SendErrorResponse(w, "IMAGE_REQUIRED", http.StatusBadRequest, nil)
		return
	}

	if err := h.Application.UpdateDonationImage(r.Context(), donationID, imageID, fileHeaders[0], maker); err != nil {
		h.logger.Errorf("failed to update donation image: %v", err)
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	utils.WriteSuccessResponse(w, nil, "Donation image update request sent successfully")
}

func (h *DonationHttpStore) DeleteDonationImage(w http.ResponseWriter, r *http.Request) {
	donationID, ok := h.extractID(w, r)
	if !ok {
		return
	}

	maker, ok := h.extractUserAndMaker(w, r)
	if !ok {
		return
	}

	// Parse JSON body instead of form values
	var requestBody struct {
		ImageID string `json:"image_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		utils.SendErrorResponse(w, "INVALID_JSON_BODY", http.StatusBadRequest, nil)
		return
	}

	if requestBody.ImageID == "" {
		utils.SendErrorResponse(w, "IMAGE_ID_REQUIRED", http.StatusBadRequest, nil)
		return
	}

	if err := h.Application.DeleteDonationImage(r.Context(), donationID, requestBody.ImageID, maker); err != nil {
		h.logger.Errorf("failed to delete donation image: %v", err)
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	utils.WriteSuccessResponse(w, nil, "Donation image delete request sent successfully")
}

func (h *DonationHttpStore) AddDonationImage(w http.ResponseWriter, r *http.Request) {
	donationID, ok := h.extractID(w, r)
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

	fileHeaders := r.MultipartForm.File["donation_images"]
	if len(fileHeaders) == 0 {
		utils.SendErrorResponse(w, "IMAGE_REQUIRED", http.StatusBadRequest, nil)
		return
	}

	if err := h.Application.AddDonationImage(r.Context(), donationID, fileHeaders[0], maker); err != nil {
		h.logger.Errorf("failed to add donation image: %v", err)
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	utils.WriteSuccessResponse(w, nil, "Donation image add request sent successfully")
}
