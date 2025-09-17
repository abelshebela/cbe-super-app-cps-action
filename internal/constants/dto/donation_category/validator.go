package donation_category

import (
	"mime/multipart"
	"strings"

	"cbe-super-app-cps-action/pkgs/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)


func (d DonationCategoryRequest) ValidateForUpdate() error {
	// First check if at least one field is provided
	if d.CategoryName == "" && d.Icon == nil {
		return validation.NewError("validation_at_least_one_field", "at least one field must be provided for update")
	}

	return validation.ValidateStruct(&d,
		validation.Field(&d.CategoryName,
			validation.When(d.CategoryName != "", validation.Required.Error("category name is required"),
				validation.Length(3, 100).Error("category name must be between 3 and 100 characters"),
				validation.By(utils.NoSpecialChars)),
		),
		validation.Field(&d.Icon,
			validation.When(d.Icon != nil, validation.By(validateImage)),
		),
	)
}
func (d DonationCategoryRequest) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.CategoryName,
			validation.Required.Error("category name is required"),
			validation.Length(3, 100).Error("category name must be between 3 and 100 characters"),
			validation.By(utils.NoSpecialChars),
		),
		validation.Field(&d.Icon,
			validation.Required.Error("category icon is required"),
			validation.By(validateImage),
		),
	)
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