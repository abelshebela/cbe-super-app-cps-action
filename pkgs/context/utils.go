// Package context contains utility methods for context related functionalities.
package context

import (
	"net/http"

	constant "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/utils"
)

type UserContext struct {
	UserCode    string
	UserID      string
	FullName    string
	PhoneNumber string
	Department  string
	BranchCode  []string
}

func ExtractUserContext(r *http.Request) UserContext {
	// This method extracts users data from the middleware context
	get := func(key string) string {
		val, _ := r.Context().Value(constant.ContextKey(key)).(string)
		return val
	}

	branchCode, _ := r.Context().Value(constant.ContextKey("branch_code")).([]string)

	return UserContext{
		UserCode:    get("user_code"),
		UserID:      get("user_id"),
		FullName:    get("full_name"),
		PhoneNumber: get("phone_number"),
		Department:  get("department"),
		BranchCode:  branchCode,
	}
}

func (u UserContext) IsIncomplete() bool {
	return u.UserID == "" || u.FullName == "" || u.PhoneNumber == "" || u.Department == ""
}