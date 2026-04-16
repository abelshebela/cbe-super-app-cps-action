package core

import (
	"sort"

	local_model "cbe-super-app-cps-action/internal/constants/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	shared_type "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/types"
)

// MapParentChildRelationship builds parent rows with SubAccessList from ACCESS_ITEMS_RELATION edges.
// Self-edges (PARENT_KEY = CHILD_KEY) define a standalone parent and are not added as children.
func MapParentChildRelationship(relations []local_model.AccessItemRelation, accessList []model.APPAccessList) []model.APPAccessList {
	nodesByKey := make(map[string]model.APPAccessList, len(accessList))
	for _, al := range accessList {
		nodesByKey[al.Key] = al
	}

	// Deduplicate (PARENT_KEY, CHILD_KEY) pairs
	parentChildren := make(map[string]map[string]struct{})
	for _, rel := range relations {
		if parentChildren[rel.ParentKey] == nil {
			parentChildren[rel.ParentKey] = make(map[string]struct{})
		}
		parentChildren[rel.ParentKey][rel.ChildKey] = struct{}{}
	}

	accountedFor := make(map[string]bool)
	var result []model.APPAccessList

	for parentKey, childSet := range parentChildren {
		parent, ok := nodesByKey[parentKey]
		if !ok {
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
			child, ok := nodesByKey[childKey]
			if !ok {
				continue
			}
			parent.SubAccessList = append(parent.SubAccessList, modelToSubAccessList(&child))
			accountedFor[childKey] = true
		}
		result = append(result, parent)
		accountedFor[parent.Key] = true
	}

	for _, al := range accessList {
		if !accountedFor[al.Key] {
			result = append(result, al)
		}
	}
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
