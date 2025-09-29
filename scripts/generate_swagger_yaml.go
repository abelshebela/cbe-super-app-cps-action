package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type OpenAPISpec struct {
	OpenAPI    string     `yaml:"openapi"`
	Info       APIInfo    `yaml:"info"`
	Servers    []Server   `yaml:"servers"`
	Paths      Paths      `yaml:"paths"`
	Components Components `yaml:"components,omitempty"`
}

type APIInfo struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
}

type Server struct {
	URL         string `yaml:"url"`
	Description string `yaml:"description"`
}

type Paths map[string]map[string]interface{}
type Components struct {
	Schemas map[string]interface{} `yaml:"schemas,omitempty"`
}

func main() {
	handlersDir := "internal/handlers/rest/http"
	outputDir := "docs/swagger"

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Get all module directories
	modules, err := ioutil.ReadDir(handlersDir)
	if err != nil {
		log.Fatalf("Failed to read handlers directory: %v", err)
	}

	for _, module := range modules {
		if !module.IsDir() || shouldSkipModule(module.Name()) {
			continue
		}

		moduleName := module.Name()
		modulePath := filepath.Join(handlersDir, moduleName)

		// Create a basic OpenAPI spec for the module
		spec := createBaseSpec(moduleName)

		// Add paths from the module
		spec.Paths = extractPathsFromModule(modulePath, moduleName)

		// Generate YAML
		yamlData, err := yaml.Marshal(spec)
		if err != nil {
			log.Printf("Failed to marshal YAML for %s: %v", moduleName, err)
			continue
		}

		// Write to file
		outputFile := filepath.Join(outputDir, fmt.Sprintf("%s.yaml", moduleName))
		if err := ioutil.WriteFile(outputFile, yamlData, 0644); err != nil {
			log.Printf("Failed to write YAML file for %s: %v", moduleName, err)
			continue
		}

		log.Printf("Generated OpenAPI YAML for %s at %s", moduleName, outputFile)
	}
}

func shouldSkipModule(name string) bool {
	// Skip hidden directories and common non-module directories
	return strings.HasPrefix(name, ".") || name == "swagger" || name == "middleware"
}

func createBaseSpec(moduleName string) OpenAPISpec {
	return OpenAPISpec{
		OpenAPI: "3.0.0",
		Info: APIInfo{
			Title:       fmt.Sprintf("%s API", strings.Title(moduleName)),
			Description: fmt.Sprintf("API documentation for the %s module", moduleName),
			Version:     "1.0.0",
		},
		Servers: []Server{
			{
				URL:         "http://localhost:8080/api/v1/cbesuperapp/cps_action",
				Description: "Local Development Server",
			},
		},
		Paths:      make(Paths),
		Components: Components{Schemas: make(map[string]interface{})},
	}
}

func extractPathsFromModule(modulePath, moduleName string) Paths {
	paths := make(Paths)

	// This is a simplified example - in a real implementation, you would:
	// 1. Parse the Go files in the module directory
	// 2. Extract route information from the HTTP handlers
	// 3. Convert to OpenAPI paths
	// For now, we'll just add a placeholder path
	paths["/api/"+moduleName] = map[string]interface{}{
		"get": map[string]interface{}{
			"summary":     fmt.Sprintf("Get %s data", moduleName),
			"description": fmt.Sprintf("Retrieves %s data from the server", moduleName),
			"responses": map[string]interface{}{
				"200": map[string]interface{}{
					"description": "Successful response",
				},
			},
		},
	}

	return paths
}
