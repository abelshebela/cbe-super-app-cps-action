package core

import (
	imodel "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/model"
)

func mapCustomerBlock(id, name, label string, status imodel.CustomerSegmentChangeStatus) map[string]interface{} {
	block := map[string]interface{}{
		"name":  name,
		"label": label,
	}
	if id != "" {
		block["id"] = id
	}
	if status.IsSet() {
		block["status"] = imodel.NormalizeCustomerSegmentChangeStatus(status)
	}
	return block
}

func MapCustomerSegmentationToMap(seg imodel.CustomerSegmentation) map[string]interface{} {
	customer := make([]map[string]interface{}, 0, len(seg.Customer))
	for _, entry := range seg.Customer {
		customer = append(customer, map[string]interface{}{
			"group": mapCustomerBlock(
				entry.Group.ID,
				entry.Group.CustGroup,
				entry.Group.CustGroupLabel,
				entry.Group.Status,
			),
			"segment": mapCustomerBlock(
				entry.Segment.ID,
				entry.Segment.CustSegmentName,
				entry.Segment.CustSegmentLabel,
				entry.Segment.Status,
			),
			"sub_segment": mapCustomerBlock(
				entry.SubSegment.ID,
				entry.SubSegment.CustSubSegmentName,
				entry.SubSegment.CustSubSegmentLabel,
				entry.SubSegment.Status,
			),
		})
	}

	return map[string]interface{}{
		"customer_role": map[string]interface{}{
			"id":    seg.CustomerRole.ID,
			"name":  seg.CustomerRole.Name,
			"label": seg.CustomerRole.Label,
		},
		"customer":   customer,
		"created_at": seg.CreatedAt,
		"is_enabled": seg.IsEnabled,
		"updated_at": seg.UpdatedAt,
		"is_deleted": seg.IsDeleted,
	}
}
