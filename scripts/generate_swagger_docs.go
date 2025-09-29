package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type SwaggerDoc struct {
	OpenAPI    string                    `json:"openapi"`
	Info       Info                      `json:"info"`
	Servers    []Server                  `json:"servers"`
	Paths      map[string]map[string]any `json:"paths"`
	Components Components                `json:"components"`
}

type Info struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Contact     struct {
		Name  string `json:"name,omitempty"`
		Email string `json:
```go
"email,omitempty"`
	} `json:"contact,omitempty"`
}

type Server struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

type Components struct {
	SecuritySchemes map[string]SecurityScheme `json:"securitySchemes,omitempty"`
	Schemas         map[string]any            `json:"schemas,omitempty"`
}

type SecurityScheme struct {
	Type         string `json:"type"`
	Scheme       string `json:"scheme,omitempty"`
	BearerFormat string `json:"bearerFormat,omitempty"`
}

func main() {
	handlersDir := flag.String("handlers", "internal/handlers/rest/http", "Path to the handlers directory")
	outputDir := flag.String("output", "docs/swagger", "Output directory for Swagger files")
	flag.Parse()

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Process each module in the handlers directory
	modules, err := os.ReadDir(*handlersDir)
	if err != nil {
		log.Fatalf("Failed to read handlers directory: %v", err)
	}

	for _, module := range modules {
		if !module.IsDir() {
			continue
		}

		modulePath := filepath.Join(*handlersDir, module.Name())
		doc := createBaseSwaggerDoc(module.Name())
		processModule(modulePath, doc, module.Name())

		// Write the Swagger JSON file
		outputFile := filepath.Join(*outputDir, module.Name()+".json")
		jsonData, err := json.MarshalIndent(doc, "", "  ")
		if err != nil {
			log.Printf("Failed to marshal Swagger JSON for %s: %v", module.Name(), err)
			continue
		}

		if err := os.WriteFile(outputFile, jsonData, 0644); err != nil {
			log.Printf("Failed to write Swagger file for %s: %v", module.Name(), err)
			continue
		}

		log.Printf("Generated Swagger documentation for %s at %s", module.Name(), outputFile)
	}
}

func createBaseSwaggerDoc(moduleName string) *SwaggerDoc {
	doc := &SwaggerDoc{
		OpenAPI: "3.0.0",
		Info: Info{
			Title:       fmt.Sprintf("%s API", strings.Title(moduleName)),
			Description: fmt.Sprintf("API documentation for the %s module", moduleName),
			Version:     "1.0.0",
			Contact: struct {
				Name  string `json:"name,omitempty"`
				Email string `json:"email,omitempty"`
			}{
				Name:  "API Support",
				Email: "contact@eaglelionsystems.com",
			},
		},
		Servers: []Server{
			{
				URL:         "http://localhost:8080/api/v1/cbesuperapp/cps_action",
				Description: "Local Development Server",
			},
		},
		Paths: make(map[string]map[string]any),
		Components: Components{
			SecuritySchemes: map[string]SecurityScheme{
				"BearerAuth": {
					Type:         "http",
					Scheme:       "bearer",
					BearerFormat: "JWT",
				},
			},
			Schemas: make(map[string]any),
		},
	}

	return doc
}

func processModule(modulePath string, doc *SwaggerDoc, moduleName string) {
	// Walk through all Go files in the module directory
	err := filepath.Walk(modulePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip non-Go files and test files
		if info.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		// Parse the Go file
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			log.Printf("Failed to parse file %s: %v", path, err)
			return nil
		}

		// Process each function in the file
		ast.Inspect(node, func(n ast.Node) bool {
			// Look for function declarations
			fn, ok := n.(*ast.FuncDecl)
			if !ok || fn.Doc == nil {
				return true
			}

			// Extract API documentation from comments
			for _, comment := range fn.Doc.List {
				text := comment.Text
				if strings.HasPrefix(text, "// @") {
					// Process Swagger annotations
					processSwaggerAnnotation(text, doc, moduleName)
				}
			}

			return true
		})

		return nil
	})

	if err != nil {
		log.Printf("Error processing module %s: %v", moduleName, err)
	}
}

func processSwaggerAnnotation(annotation string, doc *SwaggerDoc, moduleName string) {
	// This is a simplified example - in a real implementation, you would parse
	// the Swagger annotations and update the Swagger document accordingly
	// For now, we'll just log the annotation
	log.Printf("Found Swagger annotation in %s: %s", moduleName, annotation)
}
