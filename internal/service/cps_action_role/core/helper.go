package core

import (
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Dif struct {
	MakerDiff   []bson.ObjectID
	CheckerDiff [][]bson.ObjectID
	AuditorDiff []bson.ObjectID
}

func DiffFinder(cur *model.CPSActionRole, new *model.CPSActionRole) (add Dif, remove Dif) {

	// Helper to diff two slices of ObjectID, preserving order
	diff := func(cur, next []bson.ObjectID) (toAdd, toRemove []bson.ObjectID) {
		curMap := make(map[bson.ObjectID]struct{}, len(cur))
		nextMap := make(map[bson.ObjectID]struct{}, len(next))
		for _, id := range cur {
			curMap[id] = struct{}{}
		}
		for _, id := range next {
			nextMap[id] = struct{}{}
		}
		// To add: in next but not in cur, preserve order from next
		for _, id := range next {
			if _, found := curMap[id]; !found {
				toAdd = append(toAdd, id)
			}
		}
		// To remove: in cur but not in next, preserve order from cur
		for _, id := range cur {
			if _, found := nextMap[id]; !found {
				toRemove = append(toRemove, id)
			}
		}
		return
	}

	// Maker diff (single slice)
	add.MakerDiff, remove.MakerDiff = diff(cur.AssignedMakersRoles, new.AssignedMakersRoles)

	// Checker diff (slice of slices)
	add.CheckerDiff = make([][]bson.ObjectID, len(new.AssignedCheckersRoles))
	remove.CheckerDiff = make([][]bson.ObjectID, len(cur.AssignedCheckersRoles))
	maxLen := len(cur.AssignedCheckersRoles)
	if len(new.AssignedCheckersRoles) > maxLen {
		maxLen = len(new.AssignedCheckersRoles)
	}
	for i := 0; i < maxLen; i++ {
		var curIDs, newIDs []bson.ObjectID
		if i < len(cur.AssignedCheckersRoles) {
			curIDs = cur.AssignedCheckersRoles[i]
		}
		if i < len(new.AssignedCheckersRoles) {
			newIDs = new.AssignedCheckersRoles[i]
		}
		toAdd, toRemove := diff(curIDs, newIDs)
		if i < len(new.AssignedCheckersRoles) {
			add.CheckerDiff[i] = toAdd
		}
		if i < len(cur.AssignedCheckersRoles) {
			remove.CheckerDiff[i] = toRemove
		}
	}

	// Auditor diff (single slice)
	add.AuditorDiff, remove.AuditorDiff = diff(cur.AssignedAuditorRoles, new.AssignedAuditorRoles)

	return
}

func MapResponseToModel(res *model.CPSActionRoleResposne) *model.CPSActionRole {
	m := &model.CPSActionRole{
		ID:            res.ID,
		ActionCode:    res.ActionCode,
		ActionName:    res.ActionName,
		ApproverCount: res.ApproverCount,
		IsMakerOnly:   res.IsMakerOnly,
		Enabled:       res.Enabled,
		UpdatedAt:     res.UpdatedAt,
		CreatedAt:     res.CreatedAt,
	}

	// Extract IDs from Makers
	for _, item := range res.AssignedMakersRoles {
		if id, ok := extractID(item); ok {
			m.AssignedMakersRoles = append(m.AssignedMakersRoles, id)
		}
	}

	// Extract IDs from Auditors
	for _, item := range res.AssignedAuditorRoles {
		if id, ok := extractID(item); ok {
			m.AssignedAuditorRoles = append(m.AssignedAuditorRoles, id)
		}
	}

	// Extract IDs from Checkers (nested arrays)
	for _, checkerArr := range res.AssignedCheckersRoles {
		var ids []bson.ObjectID
		for _, item := range checkerArr {
			if id, ok := extractID(item); ok {
				ids = append(ids, id)
			}
		}
		m.AssignedCheckersRoles = append(m.AssignedCheckersRoles, ids)
	}

	return m
}

// Helper to get ObjectID from either a string, a map, or an ObjectID type
func extractID(data interface{}) (bson.ObjectID, bool) {
	// If it's already an ObjectID
	if id, ok := data.(bson.ObjectID); ok {
		return id, true
	}
	// If it's a map (from the lookup document)
	if m, ok := data.(map[string]interface{}); ok {
		if idStr, ok := m["_id"].(string); ok {
			id, _ := bson.ObjectIDFromHex(idStr)
			return id, true
		}
	}
	return bson.NilObjectID, false
}
