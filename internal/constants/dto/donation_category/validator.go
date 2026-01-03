package donation_category

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/pkgs/utils"
	"mime/multipart"

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
