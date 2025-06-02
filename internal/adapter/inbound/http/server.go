package http

import (
	"cbe-super-app-member-users/internal/port/inbound"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	// "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/middleware"
)

func NewHTTPServer(handlers ...inbound.Handler) *http.Server {
    r := chi.NewRouter()

    // Middleware setup
    // r.Use(middleware.ZeroLogMiddleware)
    r.Use(middleware.Recoverer)
    r.Use(middleware.Timeout(60 * time.Second))
    r.Mount("/debug", middleware.Profiler())

    // Register all handlers
    r.Route("/api", func(r chi.Router) {
        for _, handler := range handlers {
            handler.RegisterRoutes(r)
        }
    })

    return &http.Server{
        Addr:         ":8080",
        Handler:      r,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }
}