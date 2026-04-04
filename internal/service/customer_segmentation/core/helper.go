package core

import (
	// "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	imodel "cbe-super-app-cps-action/internal/constants/model"
)

func MapCustomerSegmentationToMap(seg imodel.CustomerSegmentation, removedSegmentIds []string) map[string]interface{} {
	subSegments := make([]map[string]interface{}, 0)
	for _, sub := range seg.CustomerSegments {
		subSegments = append(subSegments, map[string]interface{}{
			"cus_sub_segment": sub.CustomerSubSegment,
			"cust_group":      sub.CustomerGroup,
			"cust_segment":    sub.CustomerSegment,
		})
	}

	return map[string]interface{}{
		"removed_customer_segment_ids": removedSegmentIds,
		"customer_role": map[string]interface{}{
			"id":   seg.CustomerRole.ID,
			"name": seg.CustomerRole.Name,
		},
		"t24_customer_sub_segments": subSegments,
		"created_at":                seg.CreatedAt,
		"is_enabled":                seg.IsEnabled,
		"updated_at":                seg.UpdatedAt,
		"is_deleted":                seg.IsDeleted,
	}
}
