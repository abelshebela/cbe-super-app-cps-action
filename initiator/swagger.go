package initiator

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.uber.org/zap"
)

const (
	swaggerUIHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>CBE Super App - API Documentation</title>
    <link rel="stylesheet" type="text/css" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5.10.0/swagger-ui.css">
    <style>
        body {
            margin: 0;
            padding: 0;
        }
        .topbar {
            display: none;
        }
        .swagger-ui .info {
            margin: 20px 0;
        }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5.10.0/swagger-ui-bundle.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5.10.0/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            const ui = SwaggerUIBundle({
                urls: [
                    {url: "/swagger/account_block.yaml", name: "Account Block"},
                    {url: "/swagger/account_validation.yaml", name: "Account Validation"},
                    {url: "/swagger/ad.yaml", name: "Advertisement"},
                    {url: "/swagger/amount_based_auth.yaml", name: "Amount Based Auth"},
                    {url: "/swagger/avatar.yaml", name: "Avatar"},
                    {url: "/swagger/bank.yaml", name: "Bank"},
                    {url: "/swagger/bps_user.yaml", name: "BPS User"},
                    {url: "/swagger/budget.yaml", name: "Budget"},
                    {url: "/swagger/bulk_service.yaml", name: "Bulk Service"},
                    {url: "/swagger/cps_action_handler.yaml", name: "CPS Action"},
                    {url: "/swagger/cps_user.yaml", name: "CPS User"},
                    {url: "/swagger/customer.yaml", name: "Customer"},
                    {url: "/swagger/department.yaml", name: "Department"},
                    {url: "/swagger/donation.yaml", name: "Donation"},
                    {url: "/swagger/donation_category.yaml", name: "Donation Category"},
                    {url: "/swagger/donation_company.yaml", name: "Donation Company"},
                    {url: "/swagger/event.yaml", name: "Event"},
                    {url: "/swagger/fayda.yaml", name: "Fayda"},
                    {url: "/swagger/feedback.yaml", name: "Feedback"},
                    {url: "/swagger/hq.yaml", name: "HQ"},
                    {url: "/swagger/mini_app.yaml", name: "Mini App"},
                    {url: "/swagger/mini_app_merchant.yaml", name: "Mini App Merchant"},
                    {url: "/swagger/notifications.yaml", name: "Notifications"},
                    {url: "/swagger/password_rule.yaml", name: "Password Rule"},
                    {url: "/swagger/permission.yaml", name: "Permission"},
                    {url: "/swagger/portal_card.yaml", name: "Portal Card"},
                    {url: "/swagger/product_code.yaml", name: "Product Code"},
                    {url: "/swagger/service_details.yaml", name: "Service Details"},
                    {url: "/swagger/wallet.yaml", name: "Wallet"}
                ],
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout",
                persistAuthorization: true,
                displayRequestDuration: true,
                filter: true,
                tryItOutEnabled: true
            });
            window.ui = ui;
        };
    </script>
</body>
</html>`
)

// SwaggerHandler handles Swagger UI and API spec file serving
type SwaggerHandler struct {
	logger     utils.Logger
	swaggerDir string
	fileSystem http.FileSystem
}

// NewSwaggerHandler creates a new Swagger handler
func NewSwaggerHandler(logger utils.Logger) *SwaggerHandler {
	swaggerDir := "./docs/swagger"

	// Create swagger directory if it doesn't exist
	if err := os.MkdirAll(swaggerDir, 0755); err != nil {
		logger.Errorf("Failed to create swagger directory", zap.Error(err))
	}

	return &SwaggerHandler{
		logger:     logger,
		swaggerDir: swaggerDir,
		fileSystem: http.Dir(swaggerDir),
	}
}

// ServeSwaggerUI serves the Swagger UI HTML page
func (h *SwaggerHandler) ServeSwaggerUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(swaggerUIHTML)); err != nil {
		h.logger.Errorf("Failed to write Swagger UI response", zap.Error(err))
	}
}

// ServeSwaggerSpec serves the Swagger YAML specification files
func (h *SwaggerHandler) ServeSwaggerSpec(w http.ResponseWriter, r *http.Request) {
	// Extract the file name from the URL path
	path := strings.TrimPrefix(r.URL.Path, "/swagger/")

	// Security check: prevent directory traversal
	if strings.Contains(path, "..") {
		h.logger.Warnf("Directory traversal attempt detected", zap.String("path", path))
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	// Only allow YAML files
	if !strings.HasSuffix(path, ".yaml") && !strings.HasSuffix(path, ".yml") {
		http.Error(w, "Only YAML files are allowed", http.StatusBadRequest)
		return
	}

	// Construct the full file path
	filePath := filepath.Join(h.swaggerDir, path)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		h.logger.Warnf("Swagger spec file not found", zap.String("path", filePath))
		http.Error(w, "Swagger spec not found", http.StatusNotFound)
		return
	}

	// Read and serve the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		h.logger.Errorf("Failed to read swagger spec file", zap.Error(err), zap.String("path", filePath))
		http.Error(w, "Failed to read swagger spec", http.StatusInternalServerError)
		return
	}

	// Set appropriate headers
	w.Header().Set("Content-Type", "application/x-yaml")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(content); err != nil {
		h.logger.Errorf("Failed to write swagger spec response", zap.Error(err))
	}

	h.logger.Infof("Served swagger spec", zap.String("file", path))
}

// ListSwaggerSpecs returns a list of available Swagger spec files
func (h *SwaggerHandler) ListSwaggerSpecs(w http.ResponseWriter, r *http.Request) {
	var specs []string

	err := filepath.Walk(h.swaggerDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (strings.HasSuffix(info.Name(), ".yaml") || strings.HasSuffix(info.Name(), ".yml")) {
			relPath, _ := filepath.Rel(h.swaggerDir, path)
			specs = append(specs, relPath)
		}
		return nil
	})

	if err != nil {
		h.logger.Errorf("Failed to list swagger specs", zap.Error(err))
		http.Error(w, "Failed to list swagger specs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := fmt.Sprintf(`{"specs": ["%s"]}`, strings.Join(specs, `", "`))
	if _, err := w.Write([]byte(response)); err != nil {
		h.logger.Errorf("Failed to write swagger specs list", zap.Error(err))
	}
}

// InitSwaggerRoutes initializes Swagger-related routes
func InitSwaggerRoutes(router http.Handler, logger utils.Logger) http.Handler {
	swaggerHandler := NewSwaggerHandler(logger)

	mux, ok := router.(*http.ServeMux)
	if !ok {
		logger.Warnf("Router is not *http.ServeMux, cannot register Swagger routes")
		return router
	}

	// Serve Swagger UI
	mux.HandleFunc("/swagger/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/swagger/" || r.URL.Path == "/swagger" {
			swaggerHandler.ServeSwaggerUI(w, r)
		} else {
			swaggerHandler.ServeSwaggerSpec(w, r)
		}
	})

	// Redirect /docs to /swagger/
	mux.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
	})

	// List available specs
	mux.HandleFunc("/swagger/specs", swaggerHandler.ListSwaggerSpecs)

	logger.Infof("Swagger routes initialized successfully")
	return router
}
