package customerkyc

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/pkgs/utils"
	"mime/multipart"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (r CreateCustomerKYCRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.AccountType, validation.Required),
		validation.Field(&r.CustomerName, validation.Required),
		validation.Field(&r.Address, validation.Required),
		validation.Field(&r.Nationality,
			validation.Required,
			validation.By(utils.NoSpecialChars),
		),
		validation.Field(&r.MaritalStatus, validation.Required),
		validation.Field(&r.EmploymentStatus, validation.Required),
		validation.Field(&r.TermsAndConditions,
			validation.Required,
		),
		validation.Field(&r.LivenessCheck),
		validation.Field(&r.VerificationResult),
	)
}

func (r VerificationResult) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.FaceMatchScore),
		validation.Field(&r.LivenessResult),
		validation.Field(&r.DocumentAuthenticityResult),
	)
}

func (r CustomerInfoRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.FirstName,
			validation.Required,
			validation.By(utils.NoSpecialChars),
		),
		validation.Field(&r.MiddleName,
			validation.By(utils.NoSpecialChars),
		),
		validation.Field(&r.LastName,
			validation.Required,
			validation.By(utils.NoSpecialChars),
		),
		validation.Field(&r.PhoneNumber,
			validation.Required,
		),
		validation.Field(&r.Email),
		validation.Field(&r.DateOfBirth, validation.Required),
		validation.Field(&r.Gender, validation.Required),
		validation.Field(&r.MotherName,
			validation.Required,
			validation.By(utils.NoSpecialChars),
		),
	)
}

func (r AddressRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Country,
			validation.Required,
			validation.By(utils.NoSpecialChars),
		),
		validation.Field(&r.Region,
			validation.Required,
			validation.By(utils.NoSpecialChars),
		),
		validation.Field(&r.City,
			validation.Required,
			validation.By(utils.NoSpecialChars),
		),
		validation.Field(&r.SubCity, validation.By(utils.NoSpecialChars)),
		validation.Field(&r.Wereda, validation.By(utils.NoSpecialChars)),
		validation.Field(&r.Kebele, validation.By(utils.NoSpecialChars)),
		validation.Field(&r.HouseNumber),
	)
}

func (r UpdateKYCStatusRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.KYCStatus, validation.Required, validation.In("APPROVED", "REJECTED").Error("must be either of 'APPROVED' or 'REJECTED'")),
	)
}

// ********************************************************************************************************************** //

func enumOrNil(value interface{}) error {
	if value == nil {
		return nil
	}
	return nil
}

func nestedOrNilCustomer(value interface{}) error {
	if value == nil {
		return nil
	}
	if v, ok := value.(*CustomerInfoUpdateRequest); ok {
		return v.Validate()
	}
	return nil
}

func nestedOrNilAddress(value interface{}) error {
	if value == nil {
		return nil
	}
	if v, ok := value.(*AddressUpdateRequest); ok {
		return v.Validate()
	}
	return nil
}

func nestedOrNilAliveness(value interface{}) error {
	if value == nil {
		return nil
	}
	if v, ok := value.(*LivenessCheckRequest); ok {
		return v.Validate()
	}
	return nil
}

func (r *CustomerInfoUpdateRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.FirstName, validation.By(utils.NoSpecialChars)),
		validation.Field(&r.MiddleName, validation.By(utils.NoSpecialChars)),
		validation.Field(&r.LastName, validation.By(utils.NoSpecialChars)),
		validation.Field(&r.PhoneNumber),
		validation.Field(&r.DateOfBirth),
		validation.Field(&r.Gender),
		validation.Field(&r.MotherName, validation.By(utils.NoSpecialChars)),
	)
}

func (r *AddressUpdateRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Country, validation.By(utils.NoSpecialChars)),
		validation.Field(&r.Region, validation.By(utils.NoSpecialChars)),
		validation.Field(&r.City, validation.By(utils.NoSpecialChars)),
		validation.Field(&r.SubCity, validation.By(utils.NoSpecialChars)),
		validation.Field(&r.Wereda, validation.By(utils.NoSpecialChars)),
		validation.Field(&r.Kebele, validation.By(utils.NoSpecialChars)),
		validation.Field(&r.HouseNumber),
	)
}

func (r *LivenessCheckRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.IDCardFront, validation.Required, validation.By(validateImage)),
		validation.Field(&r.IDCardBack, validation.Required, validation.By(validateImage)),
		validation.Field(&r.LivenessCheckVideo, validation.Required, validation.By(validateVideo)),
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

func validateVideo(value interface{}) error {
	file, ok := value.(*multipart.FileHeader)
	if !ok || file == nil {
		return localization.ErrorMissingOrInvalidVideo
	}
	if !utils.IsValidVideo(file) {
		return localization.ErrorMissingOrInvalidVideo
	}
	return nil
}
