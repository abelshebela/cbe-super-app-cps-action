package access_list_segmentation_core

import (
	access_list_segmentation_dto "cbe-super-app-cps-action/internal/constants/dto/access_list_segmentation"
	local_model "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
)

// ConvertPaginatedModelToDTO converts a paginated response of []*AccessListSegmentation to a paginated response of []AccessListSegmentationResponse
func ConvertPaginatedModelToDTO(paginated types.PaginatedResponse[[]*local_model.AccessListSegmentation]) types.PaginatedResponse[[]access_list_segmentation_dto.AccessListSegmentationResponse] {
	var dtoList []access_list_segmentation_dto.AccessListSegmentationResponse
	for _, m := range paginated.Data {
		if m != nil {
			dtoList = append(dtoList, MapModelToDTO(*m))
		}
	}
	return types.PaginatedResponse[[]access_list_segmentation_dto.AccessListSegmentationResponse]{
		Data: dtoList,
		Meta: paginated.Meta,
	}
}

func MapModelToDTO(model local_model.AccessListSegmentation) access_list_segmentation_dto.AccessListSegmentationResponse {
	return access_list_segmentation_dto.AccessListSegmentationResponse{
		ID:          model.ID.Hex(),
		Type:        model.Type,
		SegmentType: model.SegmentationType,
		SegmentCode: model.SegmentationCode,
		SegmentName: model.SegmentationName,
		SegmentedID: model.SegmentedID.Hex(),
		ServiceID:   model.ServiceID.Hex(),
		ServiceName: model.ServiceName,
		Enabled:     model.Enabled,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}

func GetMissingIds(request []string, models []map[string]interface{}) []string {
	var missingIds []string
	idSet := make(map[string]struct{})
	for _, m := range models {
		if idVal, ok := m["_id"]; ok {
			if idStr, ok := idVal.(string); ok {
				idSet[idStr] = struct{}{}
			}
		}
	}
	for _, reqID := range request {
		if _, found := idSet[reqID]; !found {
			missingIds = append(missingIds, reqID)
		}
	}
	return missingIds
}
