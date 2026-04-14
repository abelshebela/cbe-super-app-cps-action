package cpsactioncore

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"errors"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
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

func DataFormatter(from string, to string, log utils.Logger) (time.Time, time.Time, error) {

	ValidStartDate, ValidEndDate, err := local_util.FormatDateRangeToUTCStrings(from, to)
	if err != nil {
		log.Warnf("Invalid Start date is given ", from)
		return time.Time{}, time.Time{}, errors.New(localization.ErrorInvalidFormat.Message)
	}

	// ValidStartDate, err := local_util.ValidateTimeAndParse(FormatedFrom)
	// if err != nil {
	// 	log.Warnf("Invalid Start date is given ", from)
	// 	return time.Time{}, time.Time{}, errors.New(localization.ErrorInvalidFormat.Message)
	// }

	// ValidEndDate, err := local_util.ValidateTimeAndParse(formatedTo)
	// if err != nil {
	// 	log.Warnf("Invalid End date is given ", from)
	// 	return time.Time{}, time.Time{}, errors.New(localization.ErrorInvalidFormat.Message)
	// }

	isValidOrder, err := local_util.ValidateTimeRangeOrder(ValidStartDate, ValidEndDate)
	if err != nil {
		log.Warnf("get error while validating start and end date order error:", err)
		return time.Time{}, time.Time{}, errors.New(localization.ErrorInvalidFormat.Message)
	}

	if !isValidOrder {
		log.Warnf("end date can not be before Start Date: %v, End Date:%v", ValidStartDate, ValidEndDate)
		return time.Time{}, time.Time{}, errors.New(localization.ErrorInvalidFormat.Message)
	}

	return ValidStartDate, ValidEndDate, nil
}
