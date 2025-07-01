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
