package donation

import (
	"fmt"
	"net/http"
	"strings"

	"mime/multipart"

	"time"

	"strconv"

	donation_application "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/donation"
	dto "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	cps_entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	c "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// parseDateString tries to parse a date string in multiple formats
func parseDateString(dateStr string) (time.Time, error) {
	formats := []string{
		"02/01/2006",                // DD/MM/YYYY
		"01/02/2006",                // MM/DD/YYYY
		"2006-01-02",                // YYYY-MM-DD
		"2006-01-02T15:04:05Z07:00", // ISO format
		"2006-01-02T15:04:05",       // ISO format without timezone
		"02-01-2006",                // DD-MM-YYYY
		"01-02-2006",                // MM-DD-YYYY
	}

	for _, format := range formats {
		if parsed, err := time.Parse(format, dateStr); err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

func parseDonationAmount(amountStr string) (int, error) {
	if amountStr == "" {
		return 0, fmt.Errorf("donation amount is required")
	}

	// Clean the input string to handle various formats
	cleanedAmount := cleanAmountString(amountStr)

	// Parse as float first, then convert to int
	amount, err := strconv.ParseFloat(cleanedAmount, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid donation amount format")
	}

	// Convert to int (this will truncate decimals)
	return int(amount), nil
}

// cleanAmountString cleans and normalizes amount strings to be compatible with decimal.NewFromString
func cleanAmountString(amountStr string) string {
	// Remove leading/trailing whitespace
	amountStr = strings.TrimSpace(amountStr)

	// Check if this looks like European format (comma as decimal separator)
	// European format: "1.000,50" -> 1000.50
	// US format: "1,000.50" -> 1000.50
	isEuropeanFormat := false
	if strings.Contains(amountStr, ",") && strings.Contains(amountStr, ".") {
		// If both comma and dot exist, check the pattern
		// European: "1.000,50" (dot before comma)
		// US: "1,000.50" (comma before dot)
		commaIndex := strings.Index(amountStr, ",")
		dotIndex := strings.Index(amountStr, ".")
		if dotIndex < commaIndex {
			isEuropeanFormat = true
		}
	} else if strings.Contains(amountStr, ",") && !strings.Contains(amountStr, ".") {
		// If only comma exists, it might be European format
		// Check if there are digits after the comma
		parts := strings.Split(amountStr, ",")
		if len(parts) == 2 && len(parts[1]) <= 2 {
			// Likely European format (e.g., "1000,50")
			isEuropeanFormat = true
		}
	}

	// Remove currency symbols and other non-numeric characters except digits, dots, and minus
	// This handles cases like "$1,000.50", "€1.000,50", "1,000 USD", etc.
	var cleaned strings.Builder
	decimalFound := false
	minusFound := false

	for _, char := range amountStr {
		switch {
		case char >= '0' && char <= '9':
			cleaned.WriteRune(char)
		case char == '.' && !decimalFound:
			if isEuropeanFormat {
				// In European format, dot is thousands separator, skip it
				continue
			}
			cleaned.WriteRune(char)
			decimalFound = true
		case char == ',' && !decimalFound:
			if isEuropeanFormat {
				// In European format, comma is decimal separator
				cleaned.WriteRune('.')
				decimalFound = true
			} else {
				// Skip thousands separators before decimal point
				continue
			}
		case char == ',' && decimalFound:
			// Skip thousands separators after decimal point
			continue
		case char == '-' && !minusFound && cleaned.Len() == 0:
			// Allow minus sign only at the beginning
			cleaned.WriteRune(char)
			minusFound = true
		}
	}

	result := cleaned.String()

	// Handle edge cases
	if result == "" || result == "-" {
		return "0"
	}

	// Remove trailing decimal point
	if strings.HasSuffix(result, ".") {
		result = strings.TrimSuffix(result, ".")
	}

	return result
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

	utils.WriteSuccessResponse(w, nil, "Donation company created successfully")
}

func (h *DonationHttpStore) FetchDonationCompany(w http.ResponseWriter, r *http.Request) {
	filterParams := common_util.ExtractMongoFilterParams(r)

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

	companyName := r.FormValue("company_name")
	accountNumber := r.FormValue("account_number")
	fileHeaders := r.MultipartForm.File["company_logo"]

	var fileHeader *multipart.FileHeader
	if len(fileHeaders) > 0 {
		fileHeader = fileHeaders[0]
	}

	companyRequest := dto.DonationCompanyRequest{
		CompanyName:   companyName,
		CompanyLogo:   fileHeader,
		AccountNumber: accountNumber,
	}

	if err := companyRequest.ValidateForUpdate(); err != nil {
		h.logger.Errorf("validation failed: %v", err)
		utils.SendErrorResponse(w, "INVALID_INPUT", http.StatusBadRequest, map[string]interface{}{"errors": err})
		return
	}

	if err := h.Application.UpdateDonationCompany(r.Context(), id, companyRequest, maker); err != nil {
		h.logger.Errorf("failed to update donation company: %v", err)
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	utils.WriteSuccessResponse(w, nil, "Donation company updated successfully")
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
	if len(fileHeaders) == 0 {
		utils.SendErrorResponse(w, "DONATION_IMAGES_REQUIRED", http.StatusBadRequest, nil)
		return
	}

	target, err := parseDonationAmount(targetStr)
	if err != nil {
		utils.SendErrorResponse(w, "INVALID_DONATION_AMOUNT", http.StatusBadRequest, nil)
		return
	}

	var endDate, startDate time.Time
	if endDateStr != "" {
		endDate, err = parseDateString(endDateStr)
		if err != nil {
			utils.SendErrorResponse(w, "INVALID_END_DATE_FORMAT", http.StatusBadRequest, nil)
			return
		}
	}

	if startDateStr != "" {
		startDate, err = parseDateString(startDateStr)
		if err != nil {
			utils.SendErrorResponse(w, "INVALID_START_DATE_FORMAT", http.StatusBadRequest, nil)
			return
		}
	}

	donationRequest := dto.DonationRequest{
		CompanyID:           companyID,
		CategoryID:          categoryID,
		Title:               title,
		IsFeatured:          isFeatured,
		Target:              target,
		DonationDescription: donationDescription,
		DonationImages:      fileHeaders,
		EndDate:             endDate,
		StartDate:           startDate,
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

	utils.WriteSuccessResponse(w, nil, "Donation created successfully")
}

func (h *DonationHttpStore) FetchDonation(w http.ResponseWriter, r *http.Request) {
	filterParams := common_util.ExtractMongoFilterParams(r)

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

	fileHeaders := r.MultipartForm.File["donation_images"]

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
		endDate, err = parseDateString(endDateStr)
		if err != nil {
			utils.SendErrorResponse(w, "INVALID_END_DATE_FORMAT", http.StatusBadRequest, nil)
			return
		}
	}

	if startDateStr != "" {
		startDate, err = parseDateString(startDateStr)
		if err != nil {
			utils.SendErrorResponse(w, "INVALID_START_DATE_FORMAT", http.StatusBadRequest, nil)
			return
		}
	}

	donationRequest := dto.DonationRequest{
		CompanyID:           companyID,
		CategoryID:          categoryID,
		Title:               title,
		IsFeatured:          isFeatured,
		Target:              target,
		DonationDescription: donationDescription,
		DonationImages:      fileHeaders,
		EndDate:             endDate,
		StartDate:           startDate,
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

	utils.WriteSuccessResponse(w, nil, "Donation updated successfully")
}
