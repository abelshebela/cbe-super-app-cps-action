package customersegmentation

import (
	"cbe-super-app-cps-action/pkgs/utils"

	validation "github.com/go-ozzo/ozzo-validation"
)

func (r CreateCustomerSegmentationRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.CustomerRole,
			validation.Required,
			validation.By(utils.NoSpecialChars),
			validation.By(utils.TrimWhiteSpace),
		),
		validation.Field(&r.CustomerSegment,
			validation.Required,
			validation.By(utils.NoSpecialChars),
			validation.By(utils.TrimWhiteSpace),
		),
		validation.Field(&r.CustomerSubSegment,
			validation.Required,
			validation.By(utils.NoSpecialChars),
			validation.By(utils.TrimWhiteSpace),
		),
		validation.Field(&r.CustomerGroup,
			validation.Required,
			validation.By(utils.NoSpecialChars),
			validation.By(utils.TrimWhiteSpace),
		),
	)
}

func (r UpdateCustomerSegmentationRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.CustomerRole),
		validation.Field(&r.CustomerSegment),
		validation.Field(&r.CustomerSubSegment),
		validation.Field(&r.CustomerGroup),
	)
}
