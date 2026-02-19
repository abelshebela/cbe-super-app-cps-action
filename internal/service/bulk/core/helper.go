package core

import (
	local_model "cbe-super-app-cps-action/internal/constants/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	shared_type "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/types"
)

func MapParentChildRelationship(accessListSegmentation []local_model.AccessListSegmentation, relations []local_model.AccessItemRelation, accessList *[]model.APPAccessList) ([]model.APPAccessList, []model.APPAccessList) {
	// Build parent-child map from relations
	childToParent := make(map[string]string)
	parentToChildren := make(map[string][]string)
	for _, rel := range relations {
		childToParent[rel.ChildKey] = rel.ParentKey
		parentToChildren[rel.ParentKey] = append(parentToChildren[rel.ParentKey], rel.ChildKey)
	}

	// Helper to build a map of key to APPAccessList pointer
	accessListMap := make(map[string]*model.APPAccessList)
	for i := range *accessList {
		al := &(*accessList)[i]
		accessListMap[al.Key] = al
	}

	// Prepare all as top-level initially
	topLevel := make(map[string]*model.APPAccessList)
	for k, v := range accessListMap {
		topLevel[k] = v
	}

	// Assign children to parents for accessList
	for childKey, parentKey := range childToParent {
		child, okChild := accessListMap[childKey]
		parent, okParent := accessListMap[parentKey]
		if okChild && okParent {
			// Add child as SubAccessList to parent
			parent.SubAccessList = append(parent.SubAccessList,
				modelToSubAccessList(child))
			// Remove child from top-level
			delete(topLevel, childKey)
		}
	}

	// Now do the same for accessListSegmentation, but create new APPAccessList for each
	segMap := make(map[string]*model.APPAccessList)
	for _, seg := range accessListSegmentation {
		segMap[seg.AccessListKey] = &model.APPAccessList{
			Key:            seg.AccessListKey,
			Enabled:        seg.Enabled,
			AccessListName: seg.AccessListName,
			SubAccessList:  nil,
			USSDEnabled:    false, // Not available in segmentation, set default
		}
	}
	segTopLevel := make(map[string]*model.APPAccessList)
	for k, v := range segMap {
		segTopLevel[k] = v
	}
	for childKey, parentKey := range childToParent {
		child, okChild := segMap[childKey]
		parent, okParent := segMap[parentKey]
		if okChild && okParent {
			parent.SubAccessList = append(parent.SubAccessList, modelToSubAccessList(child))
			delete(segTopLevel, childKey)
		}
	}

	// Collect results
	var result []model.APPAccessList
	for _, v := range topLevel {
		result = append(result, *v)
	}
	var segResult []model.APPAccessList
	for _, v := range segTopLevel {
		segResult = append(segResult, *v)
	}
	return result, segResult
}

// Helper to convert APPAccessList to SubAccessList
func modelToSubAccessList(al *model.APPAccessList) shared_type.SubAccessList {
	return shared_type.SubAccessList{
		Key:            al.Key,
		Enabled:        al.Enabled,
		AccessListName: al.AccessListName,
	}
}
