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
