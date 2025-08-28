package department_core

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"net/http"
	"strings"
)

func ValidateDepartmentRequest(r *http.Request, data interface{}) localization.ResponseCode {
	var department string
	var portal_cards, permission_groups []string

	if department != "" && isInvalidFormat(department) {
		return localization.ErrorInvalidFormatForDepartmentName
	}

	for _, portal_card := range portal_cards {
		if portal_card != "" && isInvalidFormat(portal_card) {
			return localization.ErrorInvalidFormatForDepartmentPortalCards
		}

	}
	for _, permission_group := range permission_groups {
		if permission_group != "" && isInvalidFormat(permission_group) {
			return localization.ErrorInvalidFormatForDepartmentPermissionGroups
		}

	}

	return localization.ResponseCode{}
}

// Helper function to check for invalid formats
func isInvalidFormat(s string) bool {
	if s == "" {
		return false
	}
	str := strings.TrimSpace(s)
	if str == "" {
		return true
	}
	for _, r := range str {
		if !(r == ' ' || (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')) {
			return true
		}
	}
	return false
}
