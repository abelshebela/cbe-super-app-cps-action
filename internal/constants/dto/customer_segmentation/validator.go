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
	capitalizeCustomerEntries(c.Customer)
}

func (r *UpdateCustomerSegmentationRequest) CapitilizeCustomerSegmentationRequest() {
	capitalizeCustomerEntries(r.Customer)
}

func capitalizeCustomerEntries(entries []imodel.CustomerEntry) {
	for i := range entries {
		entries[i].Group.CustGroup = strings.TrimSpace(strings.ToUpper(entries[i].Group.CustGroup))
		entries[i].Group.CustGroupLabel = strings.TrimSpace(entries[i].Group.CustGroupLabel)
		entries[i].Group.Status = imodel.NormalizeCustomerSegmentChangeStatus(entries[i].Group.Status)
		entries[i].Segment.CustSegmentName = strings.TrimSpace(strings.ToUpper(entries[i].Segment.CustSegmentName))
		entries[i].Segment.CustSegmentLabel = strings.TrimSpace(entries[i].Segment.CustSegmentLabel)
		entries[i].Segment.Status = imodel.NormalizeCustomerSegmentChangeStatus(entries[i].Segment.Status)
		entries[i].SubSegment.CustSubSegmentName = strings.TrimSpace(strings.ToUpper(entries[i].SubSegment.CustSubSegmentName))
		entries[i].SubSegment.CustSubSegmentLabel = strings.TrimSpace(entries[i].SubSegment.CustSubSegmentLabel)
		entries[i].SubSegment.Status = imodel.NormalizeCustomerSegmentChangeStatus(entries[i].SubSegment.Status)
	}
}

func validateUniqueCustomerEntries(value interface{}) error {
	entries, ok := value.([]imodel.CustomerEntry)
	if !ok {
		return errors.New("invalid customer entries")
	}
	return validateUniqueCustomerEntriesList(entries, false)
}

func validateUniqueCustomerEntriesForUpdate(value interface{}) error {
	entries, ok := value.([]imodel.CustomerEntry)
	if !ok {
		return errors.New("invalid customer entries")
	}
	return validateUniqueCustomerEntriesList(entries, true)
}

func validateUniqueCustomerEntriesList(entries []imodel.CustomerEntry, skipDeleted bool) error {
	nameMap := make(map[string]int)
	segmentMap := make(map[string]int)

	var duplicateNames []string
	var duplicateSegments []string

	for _, e := range entries {
		if skipDeleted {
			if imodel.NormalizeCustomerSegmentChangeStatus(e.SubSegment.Status) == imodel.CustomerSegmentStatusDeleted {
				continue
			}
			if imodel.NormalizeCustomerSegmentChangeStatus(e.Segment.Status) == imodel.CustomerSegmentStatusDeleted {
				continue
			}
		}
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

func validateChangeStatus(value interface{}) error {
	status, ok := value.(imodel.CustomerSegmentChangeStatus)
	if !ok {
		if s, ok := value.(string); ok {
			status = imodel.CustomerSegmentChangeStatus(s)
		} else {
			return nil
		}
	}
	if !status.IsSet() {
		return nil
	}
	switch imodel.NormalizeCustomerSegmentChangeStatus(status) {
	case imodel.CustomerSegmentStatusUpdated, imodel.CustomerSegmentStatusDeleted, imodel.CustomerSegmentStatusUnchanged:
		return nil
	default:
		return fmt.Errorf("invalid status: %s", status)
	}
}

func validateGroupBlockForUpdate(g imodel.CustGroupBlock) error {
	status := imodel.NormalizeCustomerSegmentChangeStatus(g.Status)
	if status == "" {
		return validation.ValidateStruct(&g,
			validation.Field(&g.CustGroup, validation.Required, validation.By(utils.TrimWhiteSpace), validation.By(utils.NoSpecialChars)),
			validation.Field(&g.CustGroupLabel, validation.By(utils.TrimWhiteSpace), validation.By(utils.NoSpecialChars)),
		)
	}
	if err := validateChangeStatus(g.Status); err != nil {
		return err
	}
	if status == imodel.CustomerSegmentStatusUpdated {
		return validation.ValidateStruct(&g,
			validation.Field(&g.ID, validation.Required, validation.By(utils.TrimWhiteSpace), validation.By(utils.NoSpecialChars)),
			validation.Field(&g.Status, validation.By(validateChangeStatus)),
			validation.Field(&g.CustGroup, validation.Required, validation.By(utils.TrimWhiteSpace), validation.By(utils.NoSpecialChars)),
			validation.Field(&g.CustGroupLabel, validation.By(utils.TrimWhiteSpace), validation.By(utils.NoSpecialChars)),
		)
	}
	return validation.ValidateStruct(&g,
		validation.Field(&g.ID, validation.Required, validation.By(utils.TrimWhiteSpace), validation.By(utils.NoSpecialChars)),
		validation.Field(&g.Status, validation.By(validateChangeStatus)),
		validation.Field(&g.CustGroupLabel, validation.By(utils.TrimWhiteSpace), validation.By(utils.NoSpecialChars)),
	)
}

func validateSegmentBlockForUpdate(s imodel.CustSegmentBlock) error {
	status := imodel.NormalizeCustomerSegmentChangeStatus(s.Status)
	if status == "" {
		return validation.ValidateStruct(&s,
			validation.Field(&s.CustSegmentName, validation.Required, validation.By(utils.TrimWhiteSpace), validation.By(utils.NoSpecialChars)),
			validation.Field(&s.CustSegmentLabel, validation.By(utils.TrimWhiteSpace), validation.By(utils.NoSpecialChars)),
		)
	}
	if err := validateChangeStatus(s.Status); err != nil {
		return err
	}
	if status == imodel.CustomerSegmentStatusUpdated {
		return validation.ValidateStruct(&s,
			validation.Field(&s.ID, validation.Required, validation.By(utils.TrimWhiteSpace), validation.By(utils.NoSpecialChars)),
			validation.Field(&s.Status, validation.By(validateChangeStatus)),
			validation.Field(&s.CustSegmentName, validation.Required, validation.By(utils.TrimWhiteSpace), validation.By(utils.NoSpecialChars)),
			validation.Field(&s.CustSegmentLabel, validation.By(utils.TrimWhiteSpace), validation.By(utils.NoSpecialChars)),
		)
	}
	return validation.ValidateStruct(&s,
		validation.Field(&s.ID, validation.Required, validation.By(utils.TrimWhiteSpace), validation.By(utils.NoSpecialChars)),
		validation.Field(&s.Status, validation.By(validateChangeStatus)),
		validation.Field(&s.CustSegmentLabel, validation.By(utils.TrimWhiteSpace), validation.By(utils.NoSpecialChars)),
	)
}

func validateSubSegmentBlockForUpdate(sub imodel.CustSubSegmentBlock) error {
	status := imodel.NormalizeCustomerSegmentChangeStatus(sub.Status)
	if status == "" {
		return validation.ValidateStruct(&sub,
			validation.Field(&sub.CustSubSegmentName, validation.Required, validation.By(utils.TrimWhiteSpace), validation.By(utils.NoSpecialChars)),
			validation.Field(&sub.CustSubSegmentLabel, validation.By(utils.TrimWhiteSpace), validation.By(utils.NoSpecialChars)),
		)
	}
	if err := validateChangeStatus(sub.Status); err != nil {
		return err
	}
	if status == imodel.CustomerSegmentStatusUpdated {
		return validation.ValidateStruct(&sub,
			validation.Field(&sub.ID, validation.Required, validation.By(utils.TrimWhiteSpace), validation.By(utils.NoSpecialChars)),
			validation.Field(&sub.Status, validation.By(validateChangeStatus)),
			validation.Field(&sub.CustSubSegmentName, validation.Required, validation.By(utils.TrimWhiteSpace), validation.By(utils.NoSpecialChars)),
			validation.Field(&sub.CustSubSegmentLabel, validation.By(utils.TrimWhiteSpace), validation.By(utils.NoSpecialChars)),
		)
	}
	return validation.ValidateStruct(&sub,
		validation.Field(&sub.ID, validation.Required, validation.By(utils.TrimWhiteSpace), validation.By(utils.NoSpecialChars)),
		validation.Field(&sub.Status, validation.By(validateChangeStatus)),
		validation.Field(&sub.CustSubSegmentLabel, validation.By(utils.TrimWhiteSpace), validation.By(utils.NoSpecialChars)),
	)
}

func validateUpdateCustomerEntry(c imodel.CustomerEntry) error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Group, validation.Required, validation.By(func(value interface{}) error {
			g, ok := value.(imodel.CustGroupBlock)
			if !ok {
				return errors.New("invalid customer group")
			}
			return validateGroupBlockForUpdate(g)
		})),
		validation.Field(&c.Segment, validation.Required, validation.By(func(value interface{}) error {
			s, ok := value.(imodel.CustSegmentBlock)
			if !ok {
				return errors.New("invalid customer segment")
			}
			return validateSegmentBlockForUpdate(s)
		})),
		validation.Field(&c.SubSegment, validation.Required, validation.By(func(value interface{}) error {
			sub, ok := value.(imodel.CustSubSegmentBlock)
			if !ok {
				return errors.New("invalid customer sub segment")
			}
			return validateSubSegmentBlockForUpdate(sub)
		})),
	)
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
			validation.By(validateUniqueCustomerEntriesForUpdate),
			validation.Each(validation.By(func(value interface{}) error {
				c, ok := value.(imodel.CustomerEntry)
				if !ok {
					return errors.New("invalid customer sub segment")
				}
				return validateUpdateCustomerEntry(c)
			})),
		),
	)
}
