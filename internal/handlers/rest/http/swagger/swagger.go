package swagger

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
)

// go:embed swagger-ui/*
var swaggerUIFS embed.FS

// ModuleSpecs maps module names to their OpenAPI specification files
var ModuleSpecs = map[string]string{
	"bank": "/docs/swagger/bank.yaml",
	// Add other modules here as they are documented
}

// RegisterSwaggerRoutes sets up the Swagger UI routes
func RegisterSwaggerRoutes(r chi.Router) {
	// Serve the main Swagger UI
	r.Get("/swagger/*", func(w http.ResponseWriter, r *http.Request) {
		// Get the requested file path
		filePath := strings.TrimPrefix(r.URL.Path, "/swagger/")
		if filePath == "" {
			filePath = "index.html"
		}

		// Handle API spec files
		if strings.HasSuffix(filePath, ".yaml") || strings.HasSuffix(filePath, ".json") {
			serveAPISpec(w, r, filePath)
			return
		}

		// Special handling for the root path
		if filePath == "swagger" || filePath == "swagger/" {
			http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
			return
		}

		// Try to serve the file from the embedded filesystem
		data, err := fs.ReadFile(swaggerUIFS, path.Join("swagger-ui", filePath))
		if err != nil {
			// If file not found, serve index.html for SPA routing
			if filePath != "index.html" {
				http.Redirect(w, r, "/swagger/", http.StatusFound)
				return
			}
			http.NotFound(w, r)
			return
		}

		// Set content type based on file extension
		contentType := "text/plain"
		switch {
		case strings.HasSuffix(filePath, ".html"):
			contentType = "text/html"
		case strings.HasSuffix(filePath, ".css"):
			contentType = "text/css"
		case strings.HasSuffix(filePath, ".js"):
			contentType = "application/javascript"
		case strings.HasSuffix(filePath, ".json"):
			contentType = "application/json"
		case strings.HasSuffix(filePath, ".yaml"), strings.HasSuffix(filePath, ".yml"):
			contentType = "application/x-yaml"
		case strings.HasSuffix(filePath, ".png"):
			contentType = "image/png"
		case strings.HasSuffix(filePath, ".svg"):
			contentType = "image/svg+xml"
		}

		w.Header().Set("Content-Type", contentType)
		w.Write(data)
	})
}

// serveAPISpec serves the OpenAPI specification files
func serveAPISpec(w http.ResponseWriter, r *http.Request, filePath string) {
	// Read the YAML file
	data, err := os.ReadFile(filepath.Join(".", filePath))
	if err != nil {
		http.Error(w, "Specification not found", http.StatusNotFound)
		return
	}

	// Set the appropriate content type
	contentType := "application/x-yaml"
	if strings.HasSuffix(filePath, ".json") {
		contentType = "application/json"
	}

	w.Header().Set("Content-Type", contentType)
	w.Write(data)
}
