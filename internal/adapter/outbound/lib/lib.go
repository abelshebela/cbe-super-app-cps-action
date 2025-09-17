package lib

import (
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
)

func MapCPSAction(src model.CPSAction) *model.CpsActionNormalized {
	return &model.CpsActionNormalized{
		ID:                 src.ID.Hex(),
		ActionCode:         src.ActionCode,
		UniqueId:           src.UniqueId,
		MakerID:            src.MakerID,
		MakerName:          src.MakerName,
		MakerPhoneNumber:   src.MakerPhoneNumber,
		CheckerID:          src.CheckerID,
		CheckerName:        src.CheckerName,
		CheckerPhoneNumber: src.CheckerPhoneNumber,
		Department:         src.Department,
		RejectionReason:    src.RejectionReason,
		PreviousAction:     src.PreviousAction,
		CurrentAction:      src.CurrentAction,
		ActionStatus:       src.ActionStatus,
		ActionType:         src.ActionType,
		RequestAction:      src.RequestAction,
		CreatedAt:          src.CreatedAt,
		LastModifiedAt:     src.LastModifiedAt,
		MakerActionTime:    src.MakerActionTime,
		CheckerActionTime:  *src.CheckerActionTime,
	}
}
