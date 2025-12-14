package cpsactionhandler

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
)

func MapCPSActionToApproval(existingAction *model.CPSAction, userData *types.UserContext, prevChecker []types.Checker) *model.CPSAction {
	usersData := make([]types.Checker, 0)

	CheckerUser := types.Checker{
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
		Department:       existingAction.Department,
		CreatedAt:        existingAction.CreatedAt,
		LastModifiedAt:   existingAction.LastModifiedAt,
		MakerActionTime:  existingAction.MakerActionTime,

		// Add approval information
		ActionStatus: constants.Approved,
		CheckerUsers: usersData,
	}
}
