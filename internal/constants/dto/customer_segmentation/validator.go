package customersegmentation

import (
	"cbe-super-app-cps-action/pkgs/utils"
	"errors"
	"fmt"
	"strings"

	imodel "cbe-super-app-cps-action/internal/constants/model"

	validation "github.com/go-ozzo/ozzo-validation"
)

func (c *CreateCustomerSegmentationRequest) CapitilizeCustomerSegmentationRequest() {
	for i := range c.Customer {
		c.Customer[i].Group.CustGroup = strings.TrimSpace(strings.ToUpper(c.Customer[i].Group.CustGroup))
		c.Customer[i].Group.CustGroupLabel = strings.TrimSpace(c.Customer[i].Group.CustGroupLabel)
		c.Customer[i].Segment.CustSegmentName = strings.TrimSpace(strings.ToUpper(c.Customer[i].Segment.CustSegmentName))
		c.Customer[i].Segment.CustSegmentLabel = strings.TrimSpace(c.Customer[i].Segment.CustSegmentLabel)
		c.Customer[i].SubSegment.CustSubSegmentName = strings.TrimSpace(strings.ToUpper(c.Customer[i].SubSegment.CustSubSegmentName))
		c.Customer[i].SubSegment.CustSubSegmentLabel = strings.TrimSpace(c.Customer[i].SubSegment.CustSubSegmentLabel)
	}
}

func validateUniqueCustomerEntries(value interface{}) error {
	entries, ok := value.([]imodel.CustomerEntry)
	if !ok {
		return errors.New("invalid customer entries")
	}

	nameMap := make(map[string]int)
	segmentMap := make(map[string]int)

	var duplicateNames []string
	var duplicateSegments []string

	for _, e := range entries {
		nameMap[strings.TrimSpace(strings.ToUpper(e.SubSegment.CustSubSegmentName))]++
		segmentMap[strings.TrimSpace(strings.ToUpper(e.Segment.CustSegmentName))]++
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

func validateCustomerEntry(c imodel.CustomerEntry) error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Group, validation.Required, validation.By(func(value interface{}) error {
			g, ok := value.(imodel.CustGroupBlock)
			if !ok {
				return errors.New("invalid customer group")
			}
			return validation.ValidateStruct(&g,
				validation.Field(&g.CustGroup,
					validation.Required,
					validation.By(utils.TrimWhiteSpace),
					validation.By(utils.NoSpecialChars),
				),
				validation.Field(&g.CustGroupLabel,
					validation.By(utils.TrimWhiteSpace),
					validation.By(utils.NoSpecialChars),
				),
			)
		})),
		validation.Field(&c.Segment, validation.Required, validation.By(func(value interface{}) error {
			s, ok := value.(imodel.CustSegmentBlock)
			if !ok {
				return errors.New("invalid customer segment")
			}
			return validation.ValidateStruct(&s,
				validation.Field(&s.CustSegmentName,
					validation.Required,
					validation.By(utils.TrimWhiteSpace),
					validation.By(utils.NoSpecialChars),
				),
				validation.Field(&s.CustSegmentLabel,
					validation.By(utils.TrimWhiteSpace),
					validation.By(utils.NoSpecialChars),
				),
			)
		})),
		validation.Field(&c.SubSegment, validation.Required, validation.By(func(value interface{}) error {
			sub, ok := value.(imodel.CustSubSegmentBlock)
			if !ok {
				return errors.New("invalid customer sub segment")
			}
			return validation.ValidateStruct(&sub,
				validation.Field(&sub.CustSubSegmentName,
					validation.Required,
					validation.By(utils.TrimWhiteSpace),
					validation.By(utils.NoSpecialChars),
				),
				validation.Field(&sub.CustSubSegmentLabel,
					validation.By(utils.TrimWhiteSpace),
					validation.By(utils.NoSpecialChars),
				),
			)
		})),
	)
}

func (r CreateCustomerSegmentationRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.CustomerRole,
			validation.Required,
			validation.By(utils.TrimWhiteSpace),
			validation.By(utils.NoSpecialChars),
		),
		validation.Field(&r.Customer,
			validation.Required,
			validation.By(validateUniqueCustomerEntries),
			validation.Each(validation.By(func(value interface{}) error {
				c, ok := value.(imodel.CustomerEntry)
				if !ok {
					return errors.New("invalid customer classification")
				}
				return validateCustomerEntry(c)
			})),
		),
	)
}

func (r UpdateCustomerSegmentationRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Customer,
			validation.Required,
			validation.By(validateUniqueCustomerEntries),
			validation.Each(validation.By(func(value interface{}) error {
				c, ok := value.(imodel.CustomerEntry)
				if !ok {
					return errors.New("invalid customer sub segment")
				}
				return validateCustomerEntry(c)
			})),
		),
	)
}
