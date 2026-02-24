// merge_swagger_examples reads docs/swagger.json and docs/swagger.examples.json,
// injects response examples from the examples file into the swagger spec, and
// writes the result back to docs/swagger.json. Run after: swag init -g cmd/main.go -o docs
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	swaggerPath     = "docs/swagger.json"
	examplesPath    = "docs/swagger.examples.json"
	pathsKey        = "paths"
	responsesKey    = "responses"
	schemaKey       = "schema"
	exampleKey      = "example"
	applicationJSON = "application/json"
)

func main1() {
	baseDir := "."
	if len(os.Args) > 1 {
		baseDir = os.Args[1]
	}

	swaggerFile := filepath.Join(baseDir, swaggerPath)
	examplesFile := filepath.Join(baseDir, examplesPath)

	swaggerBytes, err := os.ReadFile(swaggerFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read swagger: %v\n", err)
		os.Exit(1)
	}

	var swagger map[string]interface{}
	if err := json.Unmarshal(swaggerBytes, &swagger); err != nil {
		fmt.Fprintf(os.Stderr, "parse swagger: %v\n", err)
		os.Exit(1)
	}

	examplesBytes, err := os.ReadFile(examplesFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read examples: %v (file may not exist yet)\n", err)
		os.Exit(1)
	}

	var examplesDoc struct {
		Paths map[string]map[string]struct {
			Responses map[string]map[string]interface{} `json:"responses"`
		} `json:"paths"`
	}
	if err := json.Unmarshal(examplesBytes, &examplesDoc); err != nil {
		fmt.Fprintf(os.Stderr, "parse examples: %v\n", err)
		os.Exit(1)
	}

	paths, _ := swagger[pathsKey].(map[string]interface{})
	if paths == nil {
		fmt.Fprintf(os.Stderr, "swagger has no paths\n")
		os.Exit(1)
	}

	// Mark spec so we can verify the merged spec is being served (e.g. open doc.json and check for this).
	swagger["x-examples-merged"] = true

	injected := 0
	for path, pathExamples := range examplesDoc.Paths {
		pathSpec, _ := paths[path].(map[string]interface{})
		if pathSpec == nil {
			continue
		}

		for method, methodExamples := range pathExamples {
			methodSpec, _ := pathSpec[method].(map[string]interface{})
			if methodSpec == nil {
				continue
			}
			responses, _ := methodSpec[responsesKey].(map[string]interface{})
			if responses == nil {
				continue
			}
			for statusCode, responseExamples := range methodExamples.Responses {
				responseSpec, _ := responses[statusCode].(map[string]interface{})
				if responseSpec == nil {
					continue
				}
				exampleBody, ok := responseExamples[applicationJSON]
				if !ok {
					continue
				}
				// Replace schema with a simple type:object + example so Swagger UI displays it.
				// allOf + schema.example is often ignored by Swagger UI; a plain schema with example is shown.
				responseSpec[schemaKey] = map[string]interface{}{
					"type":     "object",
					exampleKey: exampleBody,
				}
				injected++
			}
		}
	}

	out, err := json.MarshalIndent(swagger, "", "    ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "marshal swagger: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(swaggerFile, out, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "write swagger: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Injected %d response examples into %s\n", injected, swaggerFile)
}
