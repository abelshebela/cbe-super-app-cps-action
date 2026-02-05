package testinit

import (
	"context"
	"net/http"

	"cbe-super-app-cps-action/internal/constants"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func NewLogger() utils.Logger {
	return utils.NewLogger()
}

func WithUserContext(r *http.Request, userID, fullName, phoneNumber, department, userRole, userCode string) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, constants.ContextKey("user_id"), userID)
	ctx = context.WithValue(ctx, constants.ContextKey("full_name"), fullName)
	ctx = context.WithValue(ctx, constants.ContextKey("phone_number"), phoneNumber)
	ctx = context.WithValue(ctx, constants.ContextKey("department"), department)
	ctx = context.WithValue(ctx, constants.ContextKey("user_role"), userRole)
	ctx = context.WithValue(ctx, constants.ContextKey("user_code"), userCode)
	return r.WithContext(ctx)
}

func WithRoleCode(r *http.Request, roleCode string) *http.Request {
	ctx := context.WithValue(r.Context(), constants.ContextKey("role_code"), roleCode)
	return r.WithContext(ctx)
}
