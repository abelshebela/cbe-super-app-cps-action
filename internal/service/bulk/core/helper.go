package core

import (
	local_model "cbe-super-app-cps-action/internal/constants/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	shared_type "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/types"
)

func MapParentChildRelationship(relations []local_model.AccessItemRelation, accessList []model.APPAccessList) []model.APPAccessList {
	var result []model.APPAccessList
	// Build parent to children map and a set of all child keys
	parentToChildren := make(map[string][]string)
	childSet := make(map[string]model.APPAccessList)
	for _, rel := range relations {
		parentToChildren[rel.ParentKey] = append(parentToChildren[rel.ParentKey], rel.ChildKey)
	}

	for _, al := range accessList {
		childSet[al.Key] = al
	}

	for parentKey, childKeys := range parentToChildren {
		parent, ok := childSet[parentKey]
		if !ok {
			continue
		}
		for _, childKey := range childKeys {
			child, ok := childSet[childKey]
			if ok && child.Key != parentKey {
				parent.SubAccessList = append(parent.SubAccessList, modelToSubAccessList(&child))
			}
		}
		result = append(result, parent)
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
