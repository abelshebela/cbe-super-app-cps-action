package core

import (
	local_model "cbe-super-app-cps-action/internal/constants/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	shared_type "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/types"
)

func MapParentChildRelationship(relations []local_model.AccessItemRelation, accessList []model.APPAccessList) []model.APPAccessList {
	var result []model.APPAccessList

	// Build a map from key to APPAccessList for quick lookup
	accessListMap := make(map[string]*model.APPAccessList)
	for i := range accessList {
		al := &accessList[i]
		accessListMap[al.Key] = al
	}

	// Build parent to children map and a set of all child keys
	parentToChildren := make(map[string][]string)
	childSet := make(map[string]struct{})
	for _, rel := range relations {
		parentToChildren[rel.ParentKey] = append(parentToChildren[rel.ParentKey], rel.ChildKey)
		childSet[rel.ChildKey] = struct{}{}
	}

	// For each parent, add its children to SubAccessList
	for parentKey, childKeys := range parentToChildren {
		if parentKey == "" {
			// Treat children with no parent as standalone nodes
			for _, childKey := range childKeys {
				child, ok := accessListMap[childKey]
				if ok {
					result = append(result, *child)
				}
			}
			continue
		}

		// Handle case where a node's parent is itself
		if parent, ok := accessListMap[parentKey]; ok && parentKey == parent.Key {
			result = append(result, *parent)
			continue
		}

		parent, ok := accessListMap[parentKey]
		if !ok {
			continue
		}
		parent.SubAccessList = []shared_type.SubAccessList{}
		for _, childKey := range childKeys {
			child, ok := accessListMap[childKey]
			if ok {
				parent.SubAccessList = append(parent.SubAccessList, modelToSubAccessList(child))
			}
		}
	}

	// Only return access lists that are not children (i.e., parents or standalone)
	for i := range accessList {
		al := accessList[i]
		if _, isChild := childSet[al.Key]; !isChild {
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
