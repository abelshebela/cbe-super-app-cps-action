package glue

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Route struct {
	Method      string
	Path        string
	Handler     func(w http.ResponseWriter, r *http.Request)
	Middlewares []func(next http.Handler) http.Handler
}

func RegisterRoutes(router chi.Router, routes []Route) {
	// whitelist := []string{"cps_action", "cps_actions"}
	// actionRouteGuard := middleware.CPSActionRouteGuard(whitelist)

	for _, route := range routes {
		// route.Middlewares = append(route.Middlewares, actionRouteGuard)
		router.With(route.Middlewares...).Method(route.Method, route.Path, http.HandlerFunc(route.Handler))
	}
}
