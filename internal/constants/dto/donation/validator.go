package donation

import (
	"mime/multipart"
	"strings"
	"time"

	"cbe-super-app-cps-action/pkgs/utils"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (d DonationRequest) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.CompanyID,
			validation.Required.Error("company ID is required"),
		),
		validation.Field(&d.CategoryID,
			validation.Required.Error("category ID is required"),
		),
		validation.Field(&d.Title,
			validation.Required.Error("title is required"),
			validation.Length(5, 200).Error("title must be between 5 and 200 characters"),
			validation.By(utils.NoSpecialChars),
		),
		validation.Field(&d.Target,
			validation.Required.Error("target is required"),
			validation.Min(1).Error("donation amount must be greater than 0"),
			validation.By(validateDonationAmount),
		),
		validation.Field(&d.DonationDescription,
			validation.Required.Error("donation description is required"),
		),
		validation.Field(&d.DonationImages,
			validation.Required.Error("cover image is required"),
			validation.By(validateDonationImages),
			
		),
		validation.Field(&d.CoverImage,
			validation.Required.Error("cover image is required"),
			validation.By(validateImage),
		),

		validation.Field(&d.StartDate,
			validation.When(!d.StartDate.IsZero(), validation.By(validateStartDate)),
		),
		validation.Field(&d.EndDate,
			validation.Required.Error("end date is required"),
			validation.By(validateEndDate(d.StartDate)),
		),
	)
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
func validateDonationImages(value interface{}) error {
	files, ok := value.([]*multipart.FileHeader)
	if !ok {
		return validation.NewError("validation_images_invalid", "invalid image files")
	}

	for _, file := range files {
		// Check file size
		if file.Size > 10*1024*1024 {
			return validation.NewError("validation_image_size",
				"image exceeds the 10MB size limit")
		}

		// Check file extension
		if !hasAllowedExtension(file.Filename, []string{".jpg", ".jpeg", ".png", ".gif"}) {
			return validation.NewError("validation_image_format", "image must be JPG, JPEG, PNG, or GIF")
		}
	}

	return nil
}
func hasAllowedExtension(filename string, allowed []string) bool {
	if filename == "" {
		return false
	}

	filename = strings.ToLower(strings.TrimSpace(filename))
	for _, ext := range allowed {
		if strings.HasSuffix(filename, strings.ToLower(ext)) {
			return true
		}
	}
	return false
}

func validateImage(value interface{}) error {
	file, ok := value.(*multipart.FileHeader)
	if !ok {
		return validation.NewError("validation_image_invalid", "invalid image file")
	}

	// Check file size
	if file.Size > 10*1024*1024 {
		return validation.NewError("validation_image_size", "image file size must not exceed 10MB")
	}

	// Check file extension
	if !hasAllowedExtension(file.Filename, []string{".jpg", ".jpeg", ".png", ".gif"}) {
		return validation.NewError("validation_image_format", "image must be JPG, JPEG, PNG, or GIF")
	}

	return nil
}
