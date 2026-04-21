package access_list_segmentation_dto

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"errors"
	"strings"
)

func (c *EnableDisableAccessListSegmentationRequest) Validate() error {
	if c.AccessListKeys == nil {
		return errors.New(localization.ErrorAccessListKeysRequired.Code)
	}
	if c.SegmentationType == "" {

		return errors.New(localization.ErrorSegmentationTypeRequired.Code)
	}
	return nil
}
func (c *CreateAccessListSegmentationRequest) Validate() error {
	if c.AccessListKeys == nil {
		return errors.New(localization.ErrorAccessListKeysRequired.Code)
	}
	st := strings.ToLower(strings.TrimSpace(c.SegmentType))

	if st == "" {
		return errors.New("segment type is required")
	}
	if st != "block" && st != "account" {
		return errors.New("segment type must be either 'Block' or 'Account'")
	}
	if strings.TrimSpace(c.SegmentationID) == "" {
		return errors.New(localization.ErrorAccessListSegmentationIDSRequired.Code)
	}
	c.SegmentType = st
	return nil
}

func (u *UpdateAccessListSegmentationRequest) Validate() error {
	if strings.TrimSpace(u.ID) == "" {
		return errors.New("id is required")
	}
	st := strings.ToLower(strings.TrimSpace(u.SegmentType))

	if st != "" && st != "block" && st != "account" {
		return errors.New("segment type must be either 'Block' or 'Account'")
	}
	u.SegmentType = st
	return nil
}
