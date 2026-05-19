package core

import (
	"log"
	"sort"

	local_model "cbe-super-app-cps-action/internal/constants/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	shared_type "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/types"
)

// MapParentChildRelationship builds parent rows with SubAccessList from ACCESS_ITEMS_RELATION edges.
// Self-edges (PARENT_KEY = CHILD_KEY) define a standalone parent and are not added as children.
// func MapParentChildRelationship(relations []local_model.AccessItemRelation, accessList []model.APPAccessList) []model.APPAccessList {
// 	nodesByKey := make(map[string]model.APPAccessList, len(accessList))
// 	for _, al := range accessList {
// 		nodesByKey[al.ID] = al
// 	}

// 	// Deduplicate (PARENT_KEY, CHILD_KEY) pairs
// 	parentChildren := make(map[string]map[string]struct{})
// 	for _, rel := range relations {
// 		if parentChildren[rel.ParentKey] == nil {
// 			parentChildren[rel.ParentKey] = make(map[string]struct{})
// 		}
// 		parentChildren[rel.ParentKey][rel.ChildKey] = struct{}{}
// 	}

// 	accountedFor := make(map[string]bool)
// 	var result []model.APPAccessList

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

// 	for _, al := range accessList {
// 		if !accountedFor[al.Key] {
// 			result = append(result, al)
// 		}
// 	}
// 	return result
// }

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

// Helper to convert APPAccessList to SubAccessList
func modelToSubAccessList(al *model.APPAccessList) shared_type.SubAccessList {
	return shared_type.SubAccessList{
		Key:            al.Key,
		Enabled:        al.Enabled,
		AccessListName: al.AccessListName,
	}
}

func SplitEnabledDisabledTree(accessList []model.APPAccessList) (
	enabled []model.APPAccessList,
	disabled []model.APPAccessList,
) {
	for _, node := range accessList {

		if node.Enabled {
			enabled = append(enabled, node)
		} else {
			disabled = append(disabled, node)
			continue
		}
	}
	return
}
