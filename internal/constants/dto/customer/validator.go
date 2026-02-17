package customer

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"errors"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (f *FaydaApproveRequest) Validate() error {
	return validation.ValidateStruct(&f,
		validation.Field(&f.RiskLevel, validation.Required, validation.Empty.Error("risk level can not be empty")),
	)
}

func (f *SearchCustomerByCIRequest) Validate() error {
	if f.CifOrAccountNumber == "" {
		return errors.New(localization.ErrorInvalidRequest.Code)
	}

	if !regexp.MustCompile(`^[0-9]{10,18}$`).MatchString(f.CifOrAccountNumber) && f.CifOrAccountNumber != "" {
		return errors.New(localization.ErrorCustomerAccountNumberMustContainOnlyNumbers.Code)
	}
	return nil
}
