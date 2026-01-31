package cpsactionhandler

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/types"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
)

func MapCPSActionToApproval(existingAction *model.CPSAction, userData *types.UserContext, prevChecker []model.Checker) *model.CPSAction {
	usersData := make([]model.Checker, 0)

	CheckerUser := model.Checker{
		CheckerID:          userData.UserID,
		CheckerName:        userData.FullName,
		CheckerPhoneNumber: userData.PhoneNumber,
	}

	usersData = append(prevChecker, CheckerUser)

	return &model.CPSAction{
		// Preserve all original data
		ID:               existingAction.ID,
		ActionCode:       existingAction.ActionCode,
		UniqueId:         existingAction.UniqueId,
		MakerID:          existingAction.MakerID,
		MakerName:        existingAction.MakerName,
		MakerPhoneNumber: existingAction.MakerPhoneNumber,
		RequestAction:    existingAction.RequestAction,
		PreviousAction:   existingAction.PreviousAction,
		CurrentAction:    existingAction.CurrentAction,
		ActionType:       existingAction.ActionType,
		CreatedAt:        existingAction.CreatedAt,
		LastModifiedAt:   existingAction.LastModifiedAt,
		MakerActionTime:  existingAction.MakerActionTime,

		// Add approval information
		ActionStatus: constants.Approved,
		CheckerUsers: usersData,
	}
}
