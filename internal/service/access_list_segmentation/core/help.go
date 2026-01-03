package access_list_segmentation_core

import (
	access_list_segmentation_dto "cbe-super-app-cps-action/internal/constants/dto/access_list_segmentation"
	local_model "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"

	t "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/types"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
)

// ConvertPaginatedModelToDTO converts a paginated response of []AccessListSegmentation to a paginated response of []AccessListSegmentationResponse
func ConvertPaginatedModelToDTO(paginated types.PaginatedResponse[[]local_model.AccessListSegmentation]) types.PaginatedResponse[[]access_list_segmentation_dto.AccessListSegmentationResponse] {
	var dtoList []access_list_segmentation_dto.AccessListSegmentationResponse
	for _, m := range paginated.Data {
		dtoList = append(dtoList, MapModelToDTO(m))
	}
	return types.PaginatedResponse[[]access_list_segmentation_dto.AccessListSegmentationResponse]{
		Data: dtoList,
		Meta: paginated.Meta,
	}
}

func MapModelToDTO(model local_model.AccessListSegmentation) access_list_segmentation_dto.AccessListSegmentationResponse {
	return access_list_segmentation_dto.AccessListSegmentationResponse{
		ID:             model.ID.Hex(),
		Type:           model.Type,
		SegmentType:    model.SegmentationType,
		SegmentCode:    model.SegmentationCode,
		SegmentName:    model.SegmentationName,
		SegmentedID:    model.SegmentedID.Hex(),
		AccessListKey:  model.AccessListKey,
		AccessListName: model.AccessListName,
		Enabled:        model.Enabled,
		CreatedAt:      model.CreatedAt,
		UpdatedAt:      model.UpdatedAt,
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
func FindNoneSegmentedAccessList(ctx context.Context, accessListServiceRepo storage.AppAccessListRepository, accessListSegmentation []local_model.AccessListSegmentation) []model.APPAccessList {
	var ids []string
	for _, seg := range accessListSegmentation {
		ids = append(ids, seg.AccessListKey)
	}
	idSet := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		idSet[id] = struct{}{}
	}

	var res []model.APPAccessList
	als, _ := accessListServiceRepo.FindAllByKeys(ctx, ids)
	for _, al := range als {
		// Skip parent if its key is in ids
		if _, found := idSet[al.Key]; found {
			continue
		}
		// Filter sub_access_list children whose key is in ids
		var filteredSubs []t.SubAccessList
		for _, sub := range al.SubAccessList {
			if _, found := idSet[sub.Key]; !found {
				filteredSubs = append(filteredSubs, sub)
			}
		}
		al.SubAccessList = filteredSubs
		res = append(res, al)
	}
	return res
}
