package customer_group

import (
	"errors"
	"strings"
)

func (r *CreateSegmentRequest) Validate() error {
	if strings.TrimSpace(r.CustomerGroup) == "" {
		return errors.New("customer_group is required")
	}
	if strings.TrimSpace(r.CustomerGroupLabel) == "" {
		return errors.New("customer_group_label is required")
	}
	if strings.TrimSpace(r.CustomerSegment) == "" {
		return errors.New("customer_segment is required")
	}
	if strings.TrimSpace(r.CustomerSegmentLabel) == "" {
		return errors.New("customer_segment_label is required")
	}
	if strings.TrimSpace(r.CustomerSubsegment) == "" {
		return errors.New("customer_subsegment is required")
	}
	if strings.TrimSpace(r.CustomerSubsegmentLabel) == "" {
		return errors.New("customer_subsegment_label is required")
	}
	if strings.TrimSpace(r.SuperappRole) == "" {
		return errors.New("superapp_role is required")
	}
	if strings.TrimSpace(r.SuperappRoleLabel) == "" {
		return errors.New("superapp_role_label is required")
	}
	return nil
}

func (r *UpdateSegmentRequest) Validate() error {
	if strings.TrimSpace(r.CustomerGroup) == "" {
		return errors.New("customer_group is required")
	}
	if strings.TrimSpace(r.CustomerGroupLabel) == "" {
		return errors.New("customer_group_label is required")
	}
	if strings.TrimSpace(r.CustomerSegment) == "" {
		return errors.New("customer_segment is required")
	}
	if strings.TrimSpace(r.CustomerSegmentLabel) == "" {
		return errors.New("customer_segment_label is required")
	}
	if strings.TrimSpace(r.CustomerSubsegment) == "" {
		return errors.New("customer_subsegment is required")
	}
	if strings.TrimSpace(r.CustomerSubsegmentLabel) == "" {
		return errors.New("customer_subsegment_label is required")
	}
	if strings.TrimSpace(r.SuperappRole) == "" {
		return errors.New("superapp_role is required")
	}
	if strings.TrimSpace(r.SuperappRoleLabel) == "" {
		return errors.New("superapp_role_label is required")
	}
	return nil
}

func (r *CreateSegmentRequest) Normalize() {
	r.CustomerGroup = strings.ToUpper(strings.TrimSpace(r.CustomerGroup))
	r.CustomerGroupLabel = strings.TrimSpace(r.CustomerGroupLabel)
	r.CustomerSegment = strings.ToUpper(strings.TrimSpace(r.CustomerSegment))
	r.CustomerSegmentLabel = strings.TrimSpace(r.CustomerSegmentLabel)
	r.CustomerSubsegment = strings.ToUpper(strings.TrimSpace(r.CustomerSubsegment))
	r.CustomerSubsegmentLabel = strings.TrimSpace(r.CustomerSubsegmentLabel)
	r.SuperappRole = strings.ToUpper(strings.TrimSpace(r.SuperappRole))
	r.SuperappRoleLabel = strings.TrimSpace(r.SuperappRoleLabel)
}

func (r *UpdateSegmentRequest) Normalize() {
	r.CustomerGroup = strings.ToUpper(strings.TrimSpace(r.CustomerGroup))
	r.CustomerGroupLabel = strings.TrimSpace(r.CustomerGroupLabel)
	r.CustomerSegment = strings.ToUpper(strings.TrimSpace(r.CustomerSegment))
	r.CustomerSegmentLabel = strings.TrimSpace(r.CustomerSegmentLabel)
	r.CustomerSubsegment = strings.ToUpper(strings.TrimSpace(r.CustomerSubsegment))
	r.CustomerSubsegmentLabel = strings.TrimSpace(r.CustomerSubsegmentLabel)
	r.SuperappRole = strings.ToUpper(strings.TrimSpace(r.SuperappRole))
	r.SuperappRoleLabel = strings.TrimSpace(r.SuperappRoleLabel)
}
