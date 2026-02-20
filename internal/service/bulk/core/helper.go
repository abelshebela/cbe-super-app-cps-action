package core

import (
	local_model "cbe-super-app-cps-action/internal/constants/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	shared_type "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/types"
)

func MapParentChildRelationship(relations []local_model.AccessItemRelation, accessList []model.APPAccessList) []model.APPAccessList {
	// Build a map from key to APPAccessList for quick lookup
	accessListMap := make(map[string]*model.APPAccessList)
	for i := range accessList {
		al := &accessList[i]
		accessListMap[al.Key] = al
	}

	// Build parent to children map
	parentToChildren := make(map[string][]string)
	for _, rel := range relations {
		parentToChildren[rel.ParentKey] = append(parentToChildren[rel.ParentKey], rel.ChildKey)
	}

	// For each parent, add its children to SubAccessList
	for parentKey, childKeys := range parentToChildren {
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

	return accessList
}

// Helper to convert APPAccessList to SubAccessList
func modelToSubAccessList(al *model.APPAccessList) shared_type.SubAccessList {
	return shared_type.SubAccessList{
		Key:            al.Key,
		Enabled:        al.Enabled,
		AccessListName: al.AccessListName,
	}
}
