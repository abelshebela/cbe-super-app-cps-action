package customer

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (f *FaydaApproveRequest) Validate() error {
	return validation.ValidateStruct(&f,
		validation.Field(&f.RiskLevel, validation.Required, validation.Empty.Error("risk level can not be empty")),
	)
}
