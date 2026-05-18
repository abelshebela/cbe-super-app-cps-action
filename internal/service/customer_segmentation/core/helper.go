package core

import (
	imodel "cbe-super-app-cps-action/internal/constants/model"
)

func MapCustomerSegmentationToMap(seg imodel.CustomerSegmentation, removedSegmentIds []string) map[string]interface{} {
	if len(seg.Customer) == 0 && len(seg.CustomerSegments) > 0 {
		seg.SyncCustomerFromSegments()
	}

	customer := make([]map[string]interface{}, 0, len(seg.Customer))
	for _, entry := range seg.Customer {
		customer = append(customer, map[string]interface{}{
			"group": map[string]interface{}{
				"cust_group":       entry.Group.CustGroup,
				"cust_group_label": entry.Group.CustGroupLabel,
			},
			"segment": map[string]interface{}{
				"cust_segment_name":  entry.Segment.CustSegmentName,
				"cust_segment_label": entry.Segment.CustSegmentLabel,
			},
			"sub_segment": map[string]interface{}{
				"cust_sub_segment_name":  entry.SubSegment.CustSubSegmentName,
				"cust_sub_segment_label": entry.SubSegment.CustSubSegmentLabel,
			},
		})
	}

	return map[string]interface{}{
		"removed_customer_segment_ids": removedSegmentIds,
		"customer_role": map[string]interface{}{
			"id":    seg.CustomerRole.ID,
			"name":  seg.CustomerRole.Name,
			"label": seg.CustomerRole.Label,
		},
		"customer":     customer,
		"created_at":   seg.CreatedAt,
		"is_enabled":   seg.IsEnabled,
		"updated_at":   seg.UpdatedAt,
		"is_deleted":   seg.IsDeleted,
	}
}
