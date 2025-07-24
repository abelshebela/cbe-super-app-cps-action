package glue

import (
	"cbe-super-app-budget/internal/constants"
	"cbe-super-app-budget/platform/logger"
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Router struct {
	Method      string
	Path        string
	Handler     http.HandlerFunc
	Middlewares []func(http.Handler) http.Handler
	Realm       []constants.Realm
}

func RegisterRoute(
	r chi.Router,
	routes []Router,
	zapLogger logger.Logger,
) {
	for _, route := range routes {
		for _, realm := range route.Realm {
			var handler http.Handler = route.Handler
			var endpoint string

			switch realm {
			case constants.User:
				endpoint = route.Path
			default:
				zapLogger.Fatal(context.Background(), fmt.Sprintf("Invalid Realm %s Registered on Route", realm))
			}

			for i := len(route.Middlewares) - 1; i >= 0; i-- {
				handler = route.Middlewares[i](handler)
			}

			r.Method(route.Method, endpoint, handler)
		}
	}
}
