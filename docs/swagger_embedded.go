// Package docs - embed merged swagger.json so doc.json is served from this file.
// Workflow: swag init -g cmd/main.go -o docs && go generate ./docs  (or: go run scripts/merge_swagger_examples.go)
// then build. The binary will serve this embedded spec (with examples).
// To verify: open /api/v1/cbesuperapp/cps_action/swagger/doc.json and check for "x-examples-merged": true.
package docs

import _ "embed"

//go:generate go run ../scripts/merge_swagger_examples.go ..
//go:embed swagger.json
var SwaggerJSONBytes []byte
