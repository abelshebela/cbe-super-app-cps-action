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
		ID:             model.ID.Hex(),
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

// func Maptolocal(shared []model.APPAccessList) []local_model.APPAccessList {
// 	var res []local_model.APPAccessList
// 	for _, s := range shared {
// 		res = append(res, local_model.APPAccessList{
// 			ID:             s.ID,
// 			Key:            s.Key,
// 			Enabled:        s.Enabled,
// 			AccessListName: s.AccessListName,
// 			USSDEnabled:    s.USSDEnabled,
// 		})
// 	}
// 	return res
// }

// func MapParentChildRelationship(accessListSegmentation []local_model.AccessListSegmentation, relations []local_model.AccessItemRelation, accessList *[]model.APPAccessList) ([]model.APPAccessList, []model.APPAccessList) {
// 	// Build parent-child map from relations
// 	childToParent := make(map[string]string)
// 	parentToChildren := make(map[string][]string)
// 	for _, rel := range relations {
// 		childToParent[rel.ChildKey] = rel.ParentKey
// 		parentToChildren[rel.ParentKey] = append(parentToChildren[rel.ParentKey], rel.ChildKey)
// 	}

// 	// Helper to build a map of key to APPAccessList pointer
// 	accessListMap := make(map[string]*model.APPAccessList)
// 	for i := range *accessList {
// 		al := &(*accessList)[i]
// 		accessListMap[al.Key] = al
// 	}

// 	// Prepare all as top-level initially
// 	topLevel := make(map[string]*model.APPAccessList)
// 	for k, v := range accessListMap {
// 		topLevel[k] = v
// 	}

// 	// Assign children to parents for accessList
// 	for childKey, parentKey := range childToParent {
// 		child, okChild := accessListMap[childKey]
// 		parent, okParent := accessListMap[parentKey]
// 		if okChild && okParent {
// 			// Add child as SubAccessList to parent
// 			parent.SubAccessList = append(parent.SubAccessList,
// 				modelToSubAccessList(child))
// 			// Remove child from top-level
// 			delete(topLevel, childKey)
// 		}
// 	}

// 	// Now do the same for accessListSegmentation, but create new APPAccessList for each
// 	segMap := make(map[string]*model.APPAccessList)
// 	for _, seg := range accessListSegmentation {
// 		segMap[seg.AccessListKey] = &model.APPAccessList{
// 			Key:            seg.AccessListKey,
// 			Enabled:        seg.Enabled,
// 			AccessListName: seg.AccessListName,
// 			SubAccessList:  nil,
// 			USSDEnabled:    false, // Not available in segmentation, set default
// 		}
// 	}
// 	segTopLevel := make(map[string]*model.APPAccessList)
// 	for k, v := range segMap {
// 		segTopLevel[k] = v
// 	}
// 	for childKey, parentKey := range childToParent {
// 		child, okChild := segMap[childKey]
// 		parent, okParent := segMap[parentKey]
// 		if okChild && okParent {
// 			parent.SubAccessList = append(parent.SubAccessList, modelToSubAccessList(child))
// 			delete(segTopLevel, childKey)
// 		}
// 	}

// 	// Collect results
// 	var result []model.APPAccessList
// 	for _, v := range topLevel {
// 		result = append(result, *v)
// 	}
// 	var segResult []model.APPAccessList
// 	for _, v := range segTopLevel {
// 		segResult = append(segResult, *v)
// 	}
// 	return result, segResult
// }

// Helper to convert APPAccessList to SubAccessList
func modelToSubAccessList(al *local_model.APPAccessList) types.SubAccessList {
	return types.SubAccessList{
		Key:            al.Key,
		Enabled:        al.Enabled,
		AccessListName: al.AccessListName,
	}
}

// func MapParentChildRelationship(relations []local_model.AccessItemRelation, accessList []local_model.APPAccessList) []local_model.APPAccessList {
// 	// Kept in sync with service/bulk/core.MapParentChildRelationship
// 	nodesByKey := make(map[string]local_model.APPAccessList, len(accessList))
// 	for _, al := range accessList {
// 		nodesByKey[al.Key] = al
// 	}

// 	parentChildren := make(map[string]map[string]struct{})
// 	for _, rel := range relations {
// 		if parentChildren[rel.ParentKey] == nil {
// 			parentChildren[rel.ParentKey] = make(map[string]struct{})
// 		}
// 		parentChildren[rel.ParentKey][rel.ChildKey] = struct{}{}
// 	}

// 	accountedFor := make(map[string]bool)
// 	var result []local_model.APPAccessList

// 	for parentKey, childSet := range parentChildren {
// 		parent, ok := nodesByKey[parentKey]
// 		if !ok {
// 			continue
// 		}
// 		childKeys := make([]string, 0, len(childSet))
// 		for ck := range childSet {
// 			if ck != parentKey {
// 				childKeys = append(childKeys, ck)
// 			}
// 		}
// 		sort.Strings(childKeys)
// 		for _, childKey := range childKeys {
// 			child, ok := nodesByKey[childKey]
// 			if !ok {
// 				continue
// 			}
// 			parent.SubAccessList = append(parent.SubAccessList, modelToSubAccessList(&child))
// 			accountedFor[childKey] = true
// 		}
// 		result = append(result, parent)
// 		accountedFor[parent.Key] = true
// 	}

//		for _, al := range accessList {
//			if !accountedFor[al.Key] {
//				result = append(result, al)
//			}
//		}
//		return result
//	}

// MapParentChildRelationship builds parent rows with SubAccessList from ACCESS_ITEMS_RELATION edges.
// Self-edges (PARENT_KEY = CHILD_KEY) define a standalone parent and are not added as children.
func MapParentChildRelationship(relations []local_model.AccessItemRelation, accessList []model.APPAccessList) []model.APPAccessList {
	log.Printf("[DEBUG] MapParentChildRelationship called: %d relations, %d accessList", len(relations), len(accessList))
	nodesByKey := make(map[string]model.APPAccessList, len(accessList))
	for _, al := range accessList {
		log.Printf("[DEBUG] Adding accessList node: key=%s", al.Key)
		nodesByKey[al.Key] = al
	}

	// Deduplicate (PARENT_KEY, CHILD_KEY) pairs
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

	for parentKey, childSet := range parentChildren {
		log.Printf("[DEBUG] Processing parentKey=%s with %d children", parentKey, len(childSet))
		parent, ok := nodesByKey[parentKey]
		if !ok {
			log.Printf("[WARN] Parent key not found in nodesByKey: %s", parentKey)
			continue
		}
		childKeys := make([]string, 0, len(childSet))
		for ck := range childSet {
			if ck != parentKey {
				childKeys = append(childKeys, ck)
			}
		}
		sort.Strings(childKeys)
		for _, childKey := range childKeys {
			log.Printf("[DEBUG] Processing childKey=%s for parentKey=%s", childKey, parentKey)
			child, ok := nodesByKey[childKey]
			if !ok {
				log.Printf("[WARN] Child key not found in nodesByKey: %s", childKey)
				continue
			}
			parent.SubAccessList = append(parent.SubAccessList, modelToSubAccessList(&child))
			accountedFor[childKey] = true
		}
		log.Printf("[DEBUG] Appending parent to result: key=%s, subAccessListCount=%d", parent.Key, len(parent.SubAccessList))
		result = append(result, parent)
		accountedFor[parent.Key] = true
	}

	for _, al := range accessList {
		if !accountedFor[al.Key] {
			log.Printf("[DEBUG] Appending unaccounted accessList: key=%s", al.Key)
			result = append(result, al)
		}
	}
	log.Printf("[DEBUG] MapParentChildRelationship returning %d results", len(result))
	return result
}
