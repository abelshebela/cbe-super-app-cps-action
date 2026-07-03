package core

import (
	"errors"
	"net/mail"
	"strings"
	"time"

	role_delegation_dto "cbe-super-app-cps-action/internal/constants/dto/role_delegation"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	local_util "cbe-super-app-cps-action/pkgs/utils"
)

func BuildRoleDelegationRequestForUpdate(
	body role_delegation_dto.RoleDelegationRequest,
) (imodel.RoleDelegation, error) {

	now := time.Now().UTC()

	var startAt time.Time
	var endAt time.Time

	// At least one field must be provided
	if body.StartAt.IsZero() && body.EndAt.IsZero() {
		return imodel.RoleDelegation{}, localization.ErrorInvalidDate
	}

	// Validate start date if provided
	if !body.StartAt.IsZero() {
		startAt = body.StartAt.UTC()

		if startAt.Before(now) {
			return imodel.RoleDelegation{}, localization.ErrorInvalidDate
		}
	}

	// Validate end date if provided
	if !body.EndAt.IsZero() {
		endAt = body.EndAt.UTC()

		if endAt.Before(now) {
			return imodel.RoleDelegation{}, localization.ErrorInvalidDate
		}
	}

	// If both are provided, validate the range
	if !body.StartAt.IsZero() && !body.EndAt.IsZero() {
		if !startAt.Before(endAt) {
			return imodel.RoleDelegation{}, localization.ErrorInvalidDate
		}
	}

	return imodel.RoleDelegation{
		StartAt: startAt,
		EndAt:   endAt,
	}, nil
}

// BuildRoleDelegationRequest validates and converts the request payload into a service model.
func BuildRoleDelegationRequestWithExistingUser(body role_delegation_dto.RoleDelegationRequest) (imodel.RoleDelegation, error) {
	delegatedUserID := strings.TrimSpace(body.DelegatedUserID)
	delegatedUserUserCode := strings.TrimSpace(body.DelegatedUserUserCode)
	delegatedUserUserType := strings.TrimSpace(body.DelegatedUserUserType)
	delegationType := strings.TrimSpace(body.DelegationType)
	delegatedUserDepartmentOrBranch := strings.TrimSpace(body.DelegatedUserDepartmentOrBranch)
	delegatedUserExistingRole := strings.TrimSpace(body.DelegatedUserExistingRole)

	delegatorUserID := strings.TrimSpace(body.DelegatorUserID)
	delegatorUserFullName := strings.TrimSpace(body.DelegatorUserFullName)
	delegatorUserJobTitle := strings.TrimSpace(body.DelegatorUserJobTitle)
	delegatorUserRole := strings.TrimSpace(body.DelegatorUserRole)

	newRoleID := strings.TrimSpace(body.NewRoleID)
	newDepartmentOrBranch := strings.TrimSpace(body.NewDepartmentOrBranch)

	reason := strings.TrimSpace(body.Reason)
	revoke := body.RevokeExistingDelegation

	if delegatedUserID == "" {
		return imodel.RoleDelegation{}, errors.New("Delegated user ID is required")
	}

	if delegatedUserUserCode == "" {
		return imodel.RoleDelegation{}, errors.New("Delegated user user code is required")
	}
	if delegatedUserUserType == "" {
		return imodel.RoleDelegation{}, errors.New("Delegated user type is required")
	}
	delegatedUserDepartmentOrBranch = strings.Clone(delegatedUserDepartmentOrBranch)
	if delegatedUserDepartmentOrBranch == "" {
		return imodel.RoleDelegation{}, errors.New("Delegated user department or branch is required")
	}
	// delegatedUserExistingRole = strings.Clone(delegatedUserExistingRole)
	// if delegatedUserExistingRole == "" {
	// 	return imodel.RoleDelegation{}, errors.New("Delegated user existing role is required")
	// }
	if delegationType == "" {
		return imodel.RoleDelegation{}, errors.New("Delegation type is required")
	}
	if delegatorUserID == "" {
		return imodel.RoleDelegation{}, errors.New("Delegator user ID is required")
	}
	if delegatorUserFullName == "" {
		return imodel.RoleDelegation{}, errors.New("Delegator user full name is required")
	}
	if delegatorUserJobTitle == "" {
		return imodel.RoleDelegation{}, errors.New("Delegator user job title is required")
	}
	if delegatorUserRole == "" {
		return imodel.RoleDelegation{}, errors.New("Delegator user role is required")
	}
	if newRoleID == "" {
		return imodel.RoleDelegation{}, errors.New("New role ID is required")
	}
	if newDepartmentOrBranch == "" {
		return imodel.RoleDelegation{}, errors.New("New department or branch is required")
	}
	if reason == "" {
		return imodel.RoleDelegation{}, errors.New("Reason is required")
	}

	body.DelegatedUserUserType = strings.Clone(body.DelegatedUserUserType)
	if body.DelegatedUserUserType != "CPS" && body.DelegatedUserUserType != "BPS" {
		return imodel.RoleDelegation{}, localization.ErrorInvalidDelegationUserType
	}

	if err := local_util.NoSpecialChars(newRoleID); err != nil {
		return imodel.RoleDelegation{}, localization.ErrorInvalidInputParameter
	}

	if body.StartAt.IsZero() || body.EndAt.IsZero() {
		return imodel.RoleDelegation{}, localization.ErrorInvalidDate
	}

	now := time.Now().UTC()
	parsedStart := body.StartAt.UTC()
	parsedEnd := body.EndAt.UTC()

	if parsedStart.Before(now) || parsedEnd.Before(now) {
		return imodel.RoleDelegation{}, localization.ErrorInvalidDate
	}

	if !parsedStart.Before(parsedEnd) {
		return imodel.RoleDelegation{}, localization.ErrorInvalidDate
	}

	return imodel.RoleDelegation{
		DelegatedUserID:                 delegatedUserID,
		DelegatedUserUserCode:           delegatedUserUserCode,
		DelegatedUserUserType:           delegatedUserUserType,
		DelegatedUserDepartmentOrBranch: delegatedUserDepartmentOrBranch,
		DelegatedUserExistingRole:       delegatedUserExistingRole,
		DelegationType:                  delegationType,
		DelegatorUserID:                 delegatorUserID,
		DelegatorUserFullName:           delegatorUserFullName,
		DelegatorUserJobTitle:           delegatorUserJobTitle,
		DelegatorUserRole:               delegatorUserRole,
		NewRoleID:                       newRoleID,
		NewDepartmentOrBranch:           newDepartmentOrBranch,
		StartAt:                         parsedStart,
		EndAt:                           parsedEnd,
		Reason:                          reason,
		RevokeExistingDelegation:        revoke,
	}, nil
}
func BuildRoleDelegationRequestWithNewUser(body role_delegation_dto.RoleDelegationRequest) (imodel.RoleDelegation, error) {
	delegatedUserID := strings.TrimSpace(body.DelegatedUserID)
	delegatedUserFullName := strings.TrimSpace(body.DelegatedUserFullName)
	delegatedUserUserType := strings.TrimSpace(body.DelegatedUserUserType)
	delegatedUserDepartmentOrBranch := strings.TrimSpace(body.DelegatedUserDepartmentOrBranch)
	delegatedUserJobTitle := strings.TrimSpace(body.DelegatedUserJobTitle)
	delegationType := strings.TrimSpace(body.DelegationType)
	delegatedUserPhoneNumber := strings.TrimSpace(body.DelegatedUserPhoneNumber)
	delegatedUserEmail := strings.ToLower(strings.TrimSpace(body.DelegatedUserEmail))
	// delegatedUserExistingRole := strings.ToLower(strings.TrimSpace(body.DelegatedUserExistingRole))

	delegatorUserID := strings.TrimSpace(body.DelegatorUserID)
	delegatorUserFullName := strings.TrimSpace(body.DelegatorUserFullName)
	delegatorUserJobTitle := strings.TrimSpace(body.DelegatorUserJobTitle)
	delegatorUserRole := strings.TrimSpace(body.DelegatorUserRole)
	delegatorUserUserType := strings.TrimSpace(body.DelegatorUserUserType)

	newRoleID := strings.TrimSpace(body.NewRoleID)
	newDepartmentOrBranch := strings.TrimSpace(body.NewDepartmentOrBranch)
	reason := strings.TrimSpace(body.Reason)

	if delegatedUserID == "" {
		return imodel.RoleDelegation{}, errors.New("Delegated user ID is required")
	}
	if delegatedUserFullName == "" {
		return imodel.RoleDelegation{}, errors.New("Delegated user full name is required")
	}
	if delegatedUserUserType == "" {
		return imodel.RoleDelegation{}, errors.New("Delegated user type is required")
	}
	if delegatedUserDepartmentOrBranch == "" {
		return imodel.RoleDelegation{}, errors.New("Delegated user department or branch is required")
	}
	if delegatedUserJobTitle == "" {
		return imodel.RoleDelegation{}, errors.New("Delegated user job title is required")
	}
	if delegationType == "" {
		return imodel.RoleDelegation{}, errors.New("Delegation type is required")
	}
	if delegatedUserPhoneNumber == "" {
		return imodel.RoleDelegation{}, errors.New("Delegated user phone number is required")
	}
	if delegatedUserEmail == "" {
		return imodel.RoleDelegation{}, errors.New("Delegated user email is required")
	}
	if delegatorUserID == "" {
		return imodel.RoleDelegation{}, errors.New("Delegator user ID is required")
	}
	if delegatorUserFullName == "" {
		return imodel.RoleDelegation{}, errors.New("Delegator user full name is required")
	}
	if delegatorUserJobTitle == "" {
		return imodel.RoleDelegation{}, errors.New("Delegator user job title is required")
	}
	if delegatorUserRole == "" {
		return imodel.RoleDelegation{}, errors.New("Delegator user role is required")
	}
	if newRoleID == "" {
		return imodel.RoleDelegation{}, errors.New("New role ID is required")
	}
	if newDepartmentOrBranch == "" {
		return imodel.RoleDelegation{}, errors.New("New department or branch is required")
	}
	if reason == "" {
		return imodel.RoleDelegation{}, errors.New("Reason is required")
	}
	if delegatorUserUserType == "" {
		return imodel.RoleDelegation{}, errors.New("Delegator user type is required")
	}

	body.DelegatedUserUserType = strings.Clone(body.DelegatedUserUserType)
	if body.DelegatedUserUserType != "CPS" && body.DelegatedUserUserType != "BPS" {
		return imodel.RoleDelegation{}, localization.ErrorInvalidDelegationUserType
	}

	if err := local_util.NoSpecialChars(delegatedUserJobTitle); err != nil {
		return imodel.RoleDelegation{}, localization.ErrorInvalidJobTitleID
	}

	if err := local_util.NoSpecialChars(newRoleID); err != nil {
		return imodel.RoleDelegation{}, localization.ErrorInvalidInputParameter
	}

	if _, err := mail.ParseAddress(delegatedUserEmail); err != nil {
		return imodel.RoleDelegation{}, localization.ErrorInvalidEmail
	}

	formattedPhone := local_util.FormatPhoneNumber(delegatedUserPhoneNumber)
	if formattedPhone == "" {
		return imodel.RoleDelegation{}, localization.ErrorInvalidPhoneNumber
	}

	if body.StartAt.IsZero() || body.EndAt.IsZero() {
		return imodel.RoleDelegation{}, localization.ErrorInvalidDate
	}

	now := time.Now().UTC()
	parsedStart := body.StartAt.UTC()
	parsedEnd := body.EndAt.UTC()

	if parsedStart.Before(now) || parsedEnd.Before(now) {
		return imodel.RoleDelegation{}, localization.ErrorInvalidDate
	}

	if !parsedStart.Before(parsedEnd) {
		return imodel.RoleDelegation{}, localization.ErrorInvalidDate
	}

	return imodel.RoleDelegation{
		DelegatedUserID:                 delegatedUserID,
		DelegatedUserFullName:           delegatedUserFullName,
		DelegatedUserUserType:           delegatedUserUserType,
		DelegatedUserUserCode:           local_util.GenerateCPSUserCode(),
		DelegatedUserDepartmentOrBranch: delegatedUserDepartmentOrBranch,
		DelegatedUserJobTitle:           delegatedUserJobTitle,
		DelegationType:                  delegationType,
		DelegatedUserPhoneNumber:        formattedPhone,
		DelegatedUserEmail:              delegatedUserEmail,
		// DelegatedUserExistingRole:       delegatedUserExistingRole,

		DelegatorUserID:       delegatorUserID,
		DelegatorUserFullName: delegatorUserFullName,
		DelegatorUserJobTitle: delegatorUserJobTitle,
		DelegatorUserRole:     delegatorUserRole,
		DelegatorUserUserType: delegatorUserUserType,
		NewRoleID:             newRoleID,
		NewDepartmentOrBranch: newDepartmentOrBranch,
		StartAt:               parsedStart,
		EndAt:                 parsedEnd,
		Reason:                reason,
	}, nil
}
