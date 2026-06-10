package account_sub_type_dto

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"regexp"
	"strings"
)

var specialCharRegex = regexp.MustCompile(`[^a-zA-Z0-9\s_\-]`)

var validAccountTypes = map[string]bool{
	"IFB": true,
	"CB":  true,
}

var validGenders = map[string]bool{
	"MALE":   true,
	"FEMALE": true,
	"OTHER":  true,
}

func ValidateCreateRequest(req *CreateAccountSubTypeRequest) localization.ResponseCode {
	if strings.TrimSpace(req.AccountSubTypeName) == "" {
		return localization.ErrorAccountSubTypeNameRequired
	}
	if specialCharRegex.MatchString(req.AccountSubTypeName) {
		return localization.ErrorAccountSubTypeNameSpecialChar
	}
	if strings.TrimSpace(req.AccountSubTypeCode) == "" {
		return localization.ErrorAccountSubTypeCodeRequired
	}
	if specialCharRegex.MatchString(req.AccountSubTypeCode) {
		return localization.ErrorAccountSubTypeCodeSpecialChar
	}
	if strings.TrimSpace(req.AccountType) == "" {
		return localization.ErrorAccountSubTypeAccountTypeRequired
	}
	if !validAccountTypes[strings.ToUpper(req.AccountType)] {
		return localization.ErrorAccountSubTypeInvalidAccountType
	}
	if strings.TrimSpace(req.Gender) == "" {
		return localization.ErrorAccountSubTypeGenderRequired
	}
	if !validGenders[strings.ToUpper(req.Gender)] {
		return localization.ErrorAccountSubTypeInvalidGender
	}
	return localization.ResponseCode{}
}

func ValidateUpdateRequest(req *UpdateAccountSubTypeRequest) localization.ResponseCode {
	if req.AccountSubTypeName != "" {
		if specialCharRegex.MatchString(req.AccountSubTypeName) {
			return localization.ErrorAccountSubTypeNameSpecialChar
		}
	}
	if req.AccountSubTypeCode != "" {
		if specialCharRegex.MatchString(req.AccountSubTypeCode) {
			return localization.ErrorAccountSubTypeCodeSpecialChar
		}
	}
	if req.AccountType != "" {
		if !validAccountTypes[strings.ToUpper(req.AccountType)] {
			return localization.ErrorAccountSubTypeInvalidAccountType
		}
	}
	if req.Gender != "" {
		if !validGenders[strings.ToUpper(req.Gender)] {
			return localization.ErrorAccountSubTypeInvalidGender
		}
	}
	return localization.ResponseCode{}
}

func SanitizeCreateRequest(req *CreateAccountSubTypeRequest) {
	req.AccountSubTypeName = strings.TrimSpace(req.AccountSubTypeName)
	req.AccountSubTypeCode = strings.TrimSpace(strings.ToUpper(req.AccountSubTypeCode))
	req.AccountType = strings.TrimSpace(strings.ToUpper(req.AccountType))
	req.Gender = strings.TrimSpace(strings.ToUpper(req.Gender))
}

func SanitizeUpdateRequest(req *UpdateAccountSubTypeRequest) {
	req.AccountSubTypeName = strings.TrimSpace(req.AccountSubTypeName)
	req.AccountSubTypeCode = strings.TrimSpace(strings.ToUpper(req.AccountSubTypeCode))
	req.AccountType = strings.TrimSpace(strings.ToUpper(req.AccountType))
	req.Gender = strings.TrimSpace(strings.ToUpper(req.Gender))
}
