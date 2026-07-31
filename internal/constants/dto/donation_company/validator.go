package donation_company

import (
	"mime/multipart"

	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/localization"
	"github.com/abelshebela/cbe-super-app-cps-action/pkgs/utils"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (d DonationCompanyRequest) ValidateForUpdate() error {
	// First check if at least one field is provided
	if d.CompanyName == "" && d.CompanyLogo == nil && d.PhoneNumber == "" && d.Email == "" && d.Address == "" && d.CompanyDescription == "" {
		return validation.NewError("validation_at_least_one_field", "at least one field must be provided for update")
	}

	return validation.ValidateStruct(&d,
		validation.Field(&d.CompanyName,
			validation.When(d.CompanyName != "", validation.Required.Error("company name is required"),
				validation.Length(3, 100).Error("company name must be between 3 and 100 characters"),
				validation.By(utils.NoSpecialChars)),
		),
		validation.Field(&d.CompanyLogo,
			validation.When(d.CompanyLogo != nil, validation.By(validateImage)),
		),
		validation.Field(&d.PhoneNumber,
			validation.When(d.PhoneNumber != "", validation.Required.Error("phone number is required"),
				validation.Length(9, 13).Error("phone number must be between 9 and 13 characters"),
			)),
		validation.Field(&d.Email,
			validation.When(d.Email != "", validation.Required.Error("email is required"))),
		validation.Field(&d.Address,
			validation.When(d.Address != "", validation.Required.Error("address is required"))),
	)
}

func (d DonationCompanyRequest) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.CompanyName,
			validation.Required.Error("company name is required"),
			validation.Length(3, 100).Error("company name must be between 3 and 100 characters"),
			validation.By(utils.NoSpecialChars),
		),
		validation.Field(&d.CompanyDescription,
			validation.Length(10, 500).Error("description must be between 10 and 500 characters"),
		),
		validation.Field(&d.CompanyLogo,
			validation.Required.Error("company logo is required"),
			validation.By(validateImage),
		),
		validation.Field(&d.PhoneNumber,
			validation.When(d.PhoneNumber != "", validation.Required.Error("phone number is required"),
				validation.Length(9, 13).Error("phone number must be between 9 and 13 characters"),
			)),
		validation.Field(&d.Email,
			validation.When(d.Email != "", validation.Required.Error("email is required"))),
		validation.Field(&d.Address,
			validation.When(d.Address != "", validation.Required.Error("address is required"))),
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

	if file.Size > (10 << 20) {
		return validation.NewError("logo", "file size exceeds 10MB limit")
	}

	return nil
}
