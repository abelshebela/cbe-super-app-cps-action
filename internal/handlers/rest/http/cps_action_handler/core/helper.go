package cpsactionhandler

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
)

func MapCPSActionToApproval(existingAction *model.CPSAction, userData *types.UserContext) *model.CPSAction {
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
		ActionStatus:       constants.Approved,
		CheckerID:          userData.UserID,
		CheckerName:        userData.FullName,
		CheckerPhoneNumber: userData.PhoneNumber,
	}
}
