package miniapps

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// CreateMiniAppProxyHandler initializes and returns the ReverseProxy handler
func CreateMiniAppProxyHandler(logger utils.Logger, cfg *config.VaultConfig) http.Handler {
	miniAppTargetServiceURL := cfg.CBEBaseURL + "/mini-apps/cps_action"

	targetURL, err := url.Parse(miniAppTargetServiceURL)
	if err != nil {
		logger.Errorf("Failed to parse target URL | error: %v", err)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Proxy configuration error", http.StatusInternalServerError)
		})
	}

	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: targetURL.Scheme,
		Host:   targetURL.Host,
	})

	originalDirector := proxy.Director

	proxy.Director = func(r *http.Request) {
		originalDirector(r)
		r.Host = targetURL.Host

		rctx := chi.RouteContext(r.Context())
		wildcardPath := ""
		if rctx != nil {
			wildcardPath = chi.URLParam(r, "*")
		}

		r.Header.Set("x-source-secret", cfg.JwtSecretKey)
		finalPath := targetURL.Path

		if wildcardPath != "" {
			finalPath = strings.TrimSuffix(finalPath, "/")
			finalPath = finalPath + "/" + strings.TrimPrefix(wildcardPath, "/")
		}

		r.URL.Path = finalPath

		logger.Infof("Proxying request to: %s%s", targetURL.Host, r.URL.String())
	}

	return proxy
}

func CreateMiniAppCategoryProxyHandler(logger utils.Logger, cfg *config.VaultConfig) http.Handler {
	miniAppTargetServiceURL := cfg.CBEBaseURL + "/mini-apps/cps_action/categories"

	targetURL, err := url.Parse(miniAppTargetServiceURL)
	if err != nil {
		logger.Errorf("Failed to parse target URL | error: %v", err)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Proxy configuration error", http.StatusInternalServerError)
		})
	}

	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: targetURL.Scheme,
		Host:   targetURL.Host,
	})

	originalDirector := proxy.Director

	proxy.Director = func(r *http.Request) {
		originalDirector(r)
		r.Host = targetURL.Host

		rctx := chi.RouteContext(r.Context())
		wildcardPath := ""
		if rctx != nil {
			wildcardPath = chi.URLParam(r, "*")
		}

		r.Header.Set("x-source-secret", cfg.JwtSecretKey)
		finalPath := targetURL.Path

		if wildcardPath != "" {
			finalPath = strings.TrimSuffix(finalPath, "/")
			finalPath = finalPath + "/" + strings.TrimPrefix(wildcardPath, "/")
		}

		r.URL.Path = finalPath

		logger.Infof("Proxying request to: %s%s", targetURL.Host, r.URL.String())
	}

	return proxy
}
