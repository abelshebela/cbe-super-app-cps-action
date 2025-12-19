package access_list_segmentation_dto

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"errors"
	"strings"
)

func (c CreateAccessListSegmentationRequest) Validate() error {
	if len(c.SegmentedID) == 0 {
		return errors.New("segmented_id is required and cannot be empty")
	}
	if c.ServiceID == "" {
		return errors.New(localization.ErrorServiceIdRequired.Code)
	}
	validTypes := map[string]struct{}{"R": {}, "D": {}, "C": {}, "U": {}, "B": {}}
	t := strings.TrimSpace(c.Type)
	if t == "" {
		return errors.New("type is required")
	}
	if _, ok := validTypes[t]; !ok {
		return errors.New("type must be one of: R, D, C, U, B")
	}
	return nil
}

func (u UpdateAccessListSegmentationRequest) Validate() error {
	if strings.TrimSpace(u.ID) == "" {
		return errors.New("id is required")
	}
	// if strings.TrimSpace(u.NewSegmentedID) == "" {
	// 	return errors.New("new_segmented_id is required")
	// }
	// if u.NewServiceID == "" {
	// 	return errors.New(localization.ErrorServiceIdRequired.Code)
	// }
	validTypes := map[string]struct{}{"R": {}, "D": {}, "C": {}, "U": {}, "B": {}}
	t := strings.TrimSpace(u.Type)
	if t != "" {
		if _, ok := validTypes[t]; !ok {
			return errors.New("type must be one of: R, D, C, U, B")
		}
	}
	return nil
}
