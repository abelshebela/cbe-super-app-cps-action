package customersegmentation

import (
	"cbe-super-app-cps-action/pkgs/utils"
	"errors"

	validation "github.com/go-ozzo/ozzo-validation"
)

func (r CreateCustomerSegmentationRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.CustomerRole,
			validation.Required,
			validation.By(utils.TrimWhiteSpace),
			validation.By(utils.NoSpecialChars),
		),
		validation.Field(&r.CustomerSubSegments,
			validation.Required,
			validation.Each(validation.By(func(value interface{}) error {
				c, ok := value.(CustomerSubSegments)
				if !ok {
					return errors.New("invalid customer classification")
				}

				return validation.ValidateStruct(&c,
					validation.Field(&c.Name,
						validation.Required,
						validation.By(utils.TrimWhiteSpace),
						validation.By(utils.NoSpecialChars),
					),
					validation.Field(&c.CustomerSegment,
						validation.Required,
						validation.By(utils.TrimWhiteSpace),
						validation.By(utils.NoSpecialChars),
					),
					validation.Field(&c.CustomerGroup,
						validation.Required,
						validation.By(utils.TrimWhiteSpace),
						validation.By(utils.NoSpecialChars),
					),
				)
			})),
		),
	)
}

func (r UpdateCustomerSegmentationRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.CustomerSubSegments,
			validation.Required,
			validation.Each(
				validation.By(func(value interface{}) error {
					c, ok := value.(CustomerSubSegments)
					if !ok {
						return errors.New("invalid customer sub segment")
					}

					return validation.ValidateStruct(&c,
						validation.Field(&c.Name,
							validation.Required,
							validation.By(utils.TrimWhiteSpace),
							validation.By(utils.NoSpecialChars),
						),
						validation.Field(&c.CustomerSegment,
							validation.Required,
							validation.By(utils.TrimWhiteSpace),
							validation.By(utils.NoSpecialChars),
						),
						validation.Field(&c.CustomerGroup,
							validation.Required,
							validation.By(utils.TrimWhiteSpace),
							validation.By(utils.NoSpecialChars),
						),
					)
				}),
			),
		),
	)
}
