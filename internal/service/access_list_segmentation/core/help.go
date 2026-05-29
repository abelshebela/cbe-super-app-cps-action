package access_list_segmentation_core

import (
	access_list_segmentation_dto "cbe-super-app-cps-action/internal/constants/dto/access_list_segmentation"
	"cbe-super-app-cps-action/internal/constants/model"
	local_model "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"log"
	"sort"
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
		ID:             model.ID,
		Type:           model.Type,
		SegmentType:    model.SegmentationType,
		SegmentCode:    model.SegmentationCode,
		SegmentName:    model.SegmentationName,
		SegmentedID:    model.SegmentedID,
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
func FindNoneSegmentedAccessList(ctx context.Context, accessListServiceRepo storage.BulkServiceRepository, accessListSegmentation, accessListSegmentationFromParents []local_model.APPAccessList) []local_model.APPAccessList {
	var ids []string
	for _, seg := range accessListSegmentation {
		ids = append(ids, seg.Key)
	}
	for _, seg := range accessListSegmentationFromParents {
		ids = append(ids, seg.Key)
	}
	idSet := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		idSet[id] = struct{}{}
	}

	var res []local_model.APPAccessList
	als, _ := accessListServiceRepo.FindAllForSegmentation(ctx)
	for _, al := range als {
		// Skip parent if its key is in ids
		if _, found := idSet[al.Key]; found {
			continue
		}
		// Filter sub_access_list children whose key is in ids
		// var filteredSubs []t.SubAccessList
		// for _, sub := range al.SubAccessList {
		// 	if _, found := idSet[sub.Key]; !found {
		// 		filteredSubs = append(filteredSubs, sub)
		// 	}
		// }
		// al.SubAccessList = filteredSubs
		res = append(res, al)
	}

	return res
}

// Helper to convert APPAccessList to SubAccessList
func modelToSubAccessList(al *local_model.APPAccessList) types.SubAccessList {
	return types.SubAccessList{
		Key:            al.Key,
		Enabled:        al.Enabled,
		AccessListName: al.AccessListName,
	}
}

// MapParentChildRelationship builds parent rows with SubAccessList from ACCESS_ITEMS_RELATION edges.
// Self-edges (PARENT_KEY = CHILD_KEY) define a standalone parent and are not added as children.
func MapParentChildRelationship(relations []local_model.AccessItemRelation, accessList []model.APPAccessList) []model.APPAccessList {
	log.Printf("[DEBUG] MapParentChildRelationship called: %d relations, %d accessList", len(relations), len(accessList))
	nodesByID := make(map[string]model.APPAccessList, len(accessList))
	for _, al := range accessList {
		log.Printf("[DEBUG] Adding accessList node: id=%s", al.ID)
		nodesByID[al.ID] = al
	}

	// Deduplicate (PARENT_ID, CHILD_ID) pairs
	parentChildren := make(map[string]map[string]struct{})
	for _, rel := range relations {
		log.Printf("[DEBUG] Relation: parent=%s child=%s", rel.ParentKey, rel.ChildKey)
		if parentChildren[rel.ParentKey] == nil {
			parentChildren[rel.ParentKey] = make(map[string]struct{})
		}
		parentChildren[rel.ParentKey][rel.ChildKey] = struct{}{}
	}

	accountedFor := make(map[string]bool)
	var result []model.APPAccessList

	for parentID, childSet := range parentChildren {
		log.Printf("[DEBUG] Processing parentID=%s with %d children", parentID, len(childSet))
		parent, ok := nodesByID[parentID]
		if !ok {
			log.Printf("[WARN] Parent id not found in nodesByID: %s", parentID)
			continue
		}
		childIDs := make([]string, 0, len(childSet))
		for cid := range childSet {
			if cid != parentID {
				childIDs = append(childIDs, cid)
			}
		}
		sort.Strings(childIDs)
		for _, childID := range childIDs {
			log.Printf("[DEBUG] Processing childID=%s for parentID=%s", childID, parentID)
			child, ok := nodesByID[childID]
			if !ok {
				log.Printf("[WARN] Child id not found in nodesByID: %s", childID)
				continue
			}
			parent.SubAccessList = append(parent.SubAccessList, modelToSubAccessList(&child))
			accountedFor[childID] = true
		}
		log.Printf("[DEBUG] Appending parent to result: id=%s, subAccessListCount=%d", parent.ID, len(parent.SubAccessList))
		result = append(result, parent)
		accountedFor[parent.ID] = true
	}

	for _, al := range accessList {
		if !accountedFor[al.ID] {
			log.Printf("[DEBUG] Appending unaccounted accessList: id=%s", al.ID)
			result = append(result, al)
		}
	}
	log.Printf("[DEBUG] MapParentChildRelationship returning %d results", len(result))
	return result
}
