package customersegmentation

import (
	"cbe-super-app-cps-action/pkgs/utils"
	"errors"
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
)

func (c *CreateCustomerSegmentationRequest) CapitilizeCustomerSegmentationRequest() {
	for i := range c.CustomerSubSegments {
		c.CustomerSubSegments[i].Name = strings.TrimSpace(strings.ToUpper(c.CustomerSubSegments[i].Name))
		c.CustomerSubSegments[i].CustomerSegment = strings.TrimSpace(strings.ToUpper(c.CustomerSubSegments[i].CustomerSegment))
		c.CustomerSubSegments[i].CustomerGroup = strings.TrimSpace(strings.ToUpper(c.CustomerSubSegments[i].CustomerGroup))
	}
}

func validateUniqueSubSegments(value interface{}) error {
	subSegments, ok := value.([]CustomerSubSegments)
	if !ok {
		return errors.New("invalid customer sub segments")
	}

	nameMap := make(map[string]int)
	segmentMap := make(map[string]int)

	var duplicateNames []string
	var duplicateSegments []string

	for _, s := range subSegments {
		nameMap[strings.TrimSpace(strings.ToUpper(s.Name))]++
		segmentMap[strings.TrimSpace(strings.ToUpper(s.CustomerSegment))]++
	}

	for k, v := range nameMap {
		if v > 1 {
			duplicateNames = append(duplicateNames, k)
		}
	}

	for k, v := range segmentMap {
		if v > 1 {
			duplicateSegments = append(duplicateSegments, k)
		}
	}

	var errMessages []string

	if len(duplicateNames) > 0 {
		errMessages = append(errMessages,
			fmt.Sprintf("duplicate name(s): %s", strings.Join(duplicateNames, ", ")))
	}

	if len(duplicateSegments) > 0 {
		errMessages = append(errMessages,
			fmt.Sprintf("duplicate customer segment(s): %s", strings.Join(duplicateSegments, ", ")))
	}

	if len(errMessages) > 0 {
		return errors.New(strings.Join(errMessages, " | "))
	}

	return nil
}

func (r CreateCustomerSegmentationRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.CustomerRole,
			validation.Required,
			validation.By(utils.TrimWhiteSpace),
			validation.By(utils.NoSpecialChars),
		),
		validation.Field(&r.CustomerSubSegments,
			validation.Required,
			validation.By(validateUniqueSubSegments),
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
			validation.By(validateUniqueSubSegments),
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
