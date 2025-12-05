package core

import (
	"cbe-super-app-cps-action/internal/constants/dto/donation"
	"cbe-super-app-cps-action/internal/constants/localization"

	"cbe-super-app-cps-action/pkgs/utils"
	"errors"

	"mime/multipart"
	"net/http"

	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func ParseImageUpdateRequestFromMultipartForm(r *http.Request) (donation.DonationImageUpdateRequest, error) {
	var req donation.DonationImageUpdateRequest

	if err := r.ParseMultipartForm(315 << 20); err != nil {
		return req, errors.New("failed to parse multipart form")
	}

	req.ImageID = r.FormValue("image_id")
	if req.ImageID == "" {
		return req, errors.New("image_id is required")
	}

	_, imageHeader, err := utils.ParseMultipartFormFile(r, "donation_images", 10<<20)
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return req, errors.New("image file is required")
		}
		return req, err
	}
	req.Image = imageHeader

	return req, nil
}

func ParseRequestFromMultipartForm(r *http.Request, isCreate bool) (donation.DonationRequest, error) {
	var req donation.DonationRequest

	if err := r.ParseMultipartForm(315 << 20); err != nil {
		return req, errors.New("failed to parse multipart form")
	}

	req.DonationCode = r.FormValue("donation_code")
	req.CompanyID = r.FormValue("company_id")
	req.CategoryID = r.FormValue("category_id")
	req.Title = r.FormValue("title")
	req.RemovedImages = r.Form["removed_images"]
	req.DonationDescription = r.FormValue("donation_description")
	startDateStr := r.FormValue("start_date")
	endDateStr := r.FormValue("end_date")
	if isFeatured := r.FormValue("is_featured"); isFeatured != "" {
		req.IsFeatured = isFeatured == "true"
	}

	if targetStr := r.FormValue("target"); targetStr != "" {
		converter := utils.NewNumericConverter()
		if target, err := converter.ToInt32(targetStr, "target"); err != nil {
			return req, err
		} else {
			req.Target = target
		}
	}

	if enabledStr := r.FormValue("enabled"); enabledStr != "" {
		req.Enabled = enabledStr == "true"
	}

	var endDate, startDate time.Time
	var err error

	if endDateStr != "" {
		endDate, err = utils.ParseDateString(endDateStr)
		if err != nil {
			return req, err
		}
		req.EndDate = endDate
	}

	if startDateStr != "" {
		startDate, err = utils.ParseDateString(startDateStr)
		if err != nil {
			return req, err
		}
		req.StartDate = startDate
	} else {
		req.StartDate = time.Now()
	}

	_, coverImageHeader, err := utils.ParseMultipartFormFile(r, "cover_image", 10<<20)
	if err != nil {
		if !errors.Is(err, http.ErrMissingFile) {
			return req, err
		}
	} else {
		req.CoverImage = coverImageHeader
	}

	form := r.MultipartForm
	if form != nil && form.File != nil {
		if files, exists := form.File["donation_images"]; exists {
			req.DonationImages = files
		}
	}
	return req, nil
}

func ValidateForUpdate(req donation.DonationRequest) error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.Title,
			validation.When(req.Title != "", validation.Length(5, 200).Error("title must be between 5 and 200 characters")),
			validation.When(req.Title != "", validation.By(utils.NoSpecialChars)),
		),
		validation.Field(&req.Target,
			validation.When(req.Target != 0,
				validation.Min(0).Error("donation amount must be greater than or equal to 0"),
				validation.By(validateDonationAmount)),
		),
		validation.Field(&req.DonationImages,
			validation.When(req.DonationImages != nil, validation.By(validateImages)),
		),
		validation.Field(&req.CoverImage,
			validation.When(req.CoverImage != nil, validation.By(validateImage)),
		),
		validation.Field(&req.StartDate,
			validation.When(!req.StartDate.IsZero(), validation.By(validateStartDate)),
		),
		validation.Field(&req.EndDate,
			validation.When(!req.EndDate.IsZero(), validation.By(validateEndDate(req.StartDate))),
		),
	)
}

func ValidateImageUpdate(req donation.DonationRequest) error {
	if req.DonationImages == nil {
		return errors.New("cover image is required for update")
	}
	return validateImage(req.DonationImages)
}

func ValidateImageUpdateRequest(req donation.DonationImageUpdateRequest) error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.ImageID, validation.Required.Error("image_id is required")),
		validation.Field(&req.Image, validation.Required.Error("image file is required"), validation.By(validateImage)),
	)
}

func ValidateImageAdd(req donation.DonationRequest) error {
	if len(req.DonationImages) == 0 {
		return errors.New("at least one donation image is required")
	}
	return validateImages(req.DonationImages)
}

func ExtractIDFromURL(r *http.Request) string {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}

func validateDonationAmount(value interface{}) error {
	amount, ok := value.(int32)
	if !ok {
		return validation.NewError("validation_amount_invalid", "invalid donation amount")
	}

	if amount <= 0 {
		return validation.NewError("validation_amount_zero", "donation amount must be greater than zero")
	}

	if amount > 100000000 {
		return validation.NewError("validation_amount_too_large", "donation amount must not exceed 100,000,000")
	}

	return nil
}

func validateEndDate(startDate time.Time) validation.RuleFunc {
	return func(value interface{}) error {
		endDate, ok := value.(time.Time)
		if !ok {
			return validation.NewError("validation_end_date_invalid", "invalid end date")
		}
		if !startDate.IsZero() && endDate.Before(startDate) {
			return validation.NewError("validation_end_date_before_start", "end date must be after start date")
		}
		if endDate.Before(time.Now()) {
			return validation.NewError("validation_end_date_past", "end date cannot be in the past")
		}
		return nil
	}
}

func validateStartDate(value interface{}) error {
	startDate, ok := value.(time.Time)
	if !ok {
		return validation.NewError("validation_start_date_invalid", "invalid start date")
	}
	if startDate.Before(time.Now().AddDate(0, 0, -1)) {
		return validation.NewError("validation_start_date_past", "start date cannot be in the past")
	}
	return nil
}

// func validateDonationImages(value interface{}) error {
// 	files, ok := value.([]*multipart.FileHeader)
// 	if !ok {
// 		return validation.NewError("validation_images_invalid", "invalid image files")
// 	}

// 	for _, file := range files {
// 		if file.Size > 10*1024*1024 {
// 			return validation.NewError("validation_image_size",
// 				"image exceeds the 10MB size limit")
// 		}

// 		if !hasAllowedExtension(file.Filename, []string{".jpg", ".jpeg", ".png", ".gif"}) {
// 			return validation.NewError("validation_image_format", "image must be JPG, JPEG, PNG, or GIF")
// 		}
// 	}

// 	return nil
// }

// func hasAllowedExtension(filename string, allowed []string) bool {
// 	if filename == "" {
// 		return false
// 	}

// 	filename = strings.ToLower(strings.TrimSpace(filename))
// 	for _, ext := range allowed {
// 		if strings.HasSuffix(filename, strings.ToLower(ext)) {
// 			return true
// 		}
// 	}
// 	return false
// }

// func validateImage(value interface{}) error {
// 	file, ok := value.(*multipart.FileHeader)
// 	if !ok {
// 		return validation.NewError("validation_image_invalid", "invalid image file")
// 	}

// 	if file.Size > 10*1024*1024 {
// 		return validation.NewError("validation_image_size", "image file size must not exceed 10MB")
// 	}

// 	if !hasAllowedExtension(file.Filename, []string{".jpg", ".jpeg", ".png", ".gif"}) {
// 		return validation.NewError("validation_image_format", "image must be JPG, JPEG, PNG, or GIF")
// 	}

// 	return nil
// }

func validateImage(value interface{}) error {
	file, ok := value.(*multipart.FileHeader)
	if !ok || file == nil {
		return localization.ErrorMissingOrInvalidImage
	}
	if !utils.IsValidImage(file) {
		return localization.ErrorMissingOrInvalidImage
	}

	if file.Size > (15 << 20) {
		return validation.NewError("logo", localization.MsgFileTooLarge)
	}

	return nil
}

func validateImages(value interface{}) error {
	files, ok := value.([]*multipart.FileHeader)
	if !ok || len(files) == 0 {
		return localization.ErrorMissingOrInvalidImage
	}

	for _, f := range files {
		if err := validateImage(f); err != nil {
			return err
		}
	}

	return nil
}
