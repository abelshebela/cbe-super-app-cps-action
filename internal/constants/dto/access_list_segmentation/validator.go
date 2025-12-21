package access_list_segmentation_dto

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"errors"
	"strings"
)

func (c CreateAccessListSegmentationRequest) Validate() error {
	if c.ServiceID == "" {
		return errors.New(localization.ErrorServiceIdRequired.Code)
	}
	validTypes := map[string]struct{}{"R": {}, "D": {}, "C": {}, "U": {}, "B": {}}
	t := strings.TrimSpace(c.Type)
	st := strings.TrimSpace(c.SegmentType)
	if t == "" {
		return errors.New("type is required")
	}

	if st == "" {
		return errors.New("segment type is required")
	}
	if st != "Block" && st != "Account" {
		return errors.New("segment type must be either 'Block' or 'Account'")
	}
	if st == "Block" {
		if _, ok := validTypes[t]; !ok {
			return errors.New("type must be one of: R, D, C, U, B")
		}
		if len(c.SegmentedID) == 0 {
			return errors.New("segmented_id is required for segment type 'Block'")
		}

	}
	if st == "Account" {
		if c.SegmentCode == "" {
			return errors.New("segment code is required for segment type 'Account'")
		}
		if c.SegmentName == "" {
			return errors.New("segment name is required for segment type 'Account'")
		}
	}

	return nil
}

func (u UpdateAccessListSegmentationRequest) Validate() error {
	if strings.TrimSpace(u.ID) == "" {
		return errors.New("id is required")
	}
	validTypes := map[string]struct{}{"R": {}, "D": {}, "C": {}, "U": {}, "B": {}}
	t := strings.TrimSpace(u.Type)
	st := strings.TrimSpace(u.SegmentType)
	if t == "" {
		return errors.New("type is required")
	}
	if st == "" {
		return errors.New("segment type is required")
	}
	if st != "Block" && st != "Account" {
		return errors.New("segment type must be either 'Block' or 'Account'")
	}
	if st == "Block" {
		if _, ok := validTypes[t]; !ok {
			return errors.New("type must be one of: R, D, C, U, B")
		}
	}
	if st == "Account" {
		if u.SegmentCode == "" {
			return errors.New("segment code is required for segment type 'Account'")
		}
		if u.SegmentName == "" {
			return errors.New("segment name is required for segment type 'Account'")
		}
	}
	return nil
}
