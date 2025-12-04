package cpsaction

import (
	"cbe-super-app-cps-action/internal/constants/localization"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (cp *ActionRequest) Validate() error {
	return validation.ValidateStruct(
		cp,
		validation.Field(
			&cp.RejectionReason,
			validation.Required.Error(localization.MsgCPSActionRejectionReason),
			validation.Length(10, 300).Error(localization.MsgCPSActionRejectionReason),
		),
	)
}
