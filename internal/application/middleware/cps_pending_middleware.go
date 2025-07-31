package middleware

import (
	"context"
	"net/http"

	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	contexts "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
)

type CPSActionChecker interface {
	CPSActionExists(ctx context.Context, user entities.CheckCPSAction) (bool, error)
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

			check := entities.CheckCPSAction{
				UserCode:      user.UserCode,
				FullName:      user.FullName,
				Department:    user.Department,
				PhoneNumber:   user.PhoneNumber,
				RequestAction: requestAction,
			}

			_, err := checker.CPSActionExists(r.Context(), check)
			if err != nil {
				utils.SendErrorResponse(w, err.Error(), 0, nil)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
