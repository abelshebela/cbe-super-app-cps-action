package core

import (
	imodel "cbe-super-app-cps-action/internal/constants/model"
)

func MapSegmentToMap(seg imodel.Segment) map[string]interface{} {
	return map[string]interface{}{
		"id":                        seg.ID,
		"customer_group":            seg.CustomerGroup,
		"customer_group_label":      seg.CustomerGroupLabel,
		"customer_segment":          seg.CustomerSegment,
		"customer_segment_label":    seg.CustomerSegmentLabel,
		"customer_subsegment":       seg.CustomerSubsegment,
		"customer_subsegment_label": seg.CustomerSubsegmentLabel,
		"superapp_role":             seg.SuperappRole,
		"superapp_role_label":       seg.SuperappRoleLabel,
		"is_enabled":                seg.IsEnabled,
		"created_at":                seg.CreatedAt,
		"last_modified_at":          seg.LastModifiedAt,
	}
}
