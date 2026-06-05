package core

import (
	"strings"
	"time"

	role_delegation_dto "cbe-super-app-cps-action/internal/constants/dto/role_delegation"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	local_util "cbe-super-app-cps-action/pkgs/utils"
)

// BuildRoleDelegationRequest validates and converts the request payload into a service model.
func BuildRoleDelegationRequest(body role_delegation_dto.RoleDelegationRequest) (imodel.RoleDelegation, error) {
	userID := strings.TrimSpace(body.UserID)
	jobTitleID := strings.TrimSpace(body.JobTitleID)
	startDate := strings.TrimSpace(body.StartDate)
	endDate := strings.TrimSpace(body.EndDate)

	if userID == "" || jobTitleID == "" || startDate == "" || endDate == "" {
		return imodel.RoleDelegation{}, localization.ErrorRequiredFieldMissing
	}

	if err := local_util.NoSpecialChars(userID); err != nil {
		return imodel.RoleDelegation{}, localization.ErrorInvalidCPSUserID
	}

	if err := local_util.NoSpecialChars(jobTitleID); err != nil {
		return imodel.RoleDelegation{}, localization.ErrorInvalidJobTitleID
	}

	parsedStart, err := time.Parse(time.RFC3339, startDate)
	if err != nil {
		return imodel.RoleDelegation{}, localization.ErrorInvalidDate
	}

	parsedEnd, err := time.Parse(time.RFC3339, endDate)
	if err != nil {
		return imodel.RoleDelegation{}, localization.ErrorInvalidDate
	}

	now := time.Now()
	if parsedStart.Before(now) || parsedEnd.Before(now) {
		return imodel.RoleDelegation{}, localization.ErrorInvalidDate
	}

	if !parsedStart.Before(parsedEnd) {
		return imodel.RoleDelegation{}, localization.ErrorInvalidDate
	}

	return imodel.RoleDelegation{
		UserID:   userID,
		StartAt:  parsedStart,
		EndAt:    parsedEnd,
		JobTitle: jobTitleID,
	}, nil
}
