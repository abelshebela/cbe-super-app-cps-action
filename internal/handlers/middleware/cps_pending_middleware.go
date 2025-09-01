package middleware

import (
	"context"
	"net/http"

	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	contexts "cbe-super-app-cps-action/pkgs/context"
)

type CPSActionChecker interface {
	CPSActionExists(ctx context.Context, user model.CheckCPSAction) (bool, error)
}

type CPSActionMiddlewareFactory struct {
	checker CPSActionChecker
}

func NewCPSActionMiddlewareFactory(checker CPSActionChecker) *CPSActionMiddlewareFactory {
	return &CPSActionMiddlewareFactory{checker}
}

func (f *CPSActionMiddlewareFactory) RequireNoPendingCPSActionGuard(requestAction string) func(http.Handler) http.Handler {
	return RequireNoPendingCPSActionGuard(requestAction, f.checker)
}

func RequireNoPendingCPSActionGuard(requestAction string, checker CPSActionChecker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := contexts.ExtractUserContext(r)

			check := model.CheckCPSAction{
				UserCode:      user.UserCode,
				FullName:      user.FullName,
				Department:    user.Department,
				PhoneNumber:   user.PhoneNumber,
				RequestAction: requestAction,
			}

			_, err := checker.CPSActionExists(r.Context(), check)
			if err != nil {
				fieldError := localization.CreateFieldError("", err.Error(), "", "")
				localization.SendErrorResponse(w, localization.ResponseCode{}, localization.CreateFieldErrors(fieldError), nil)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
