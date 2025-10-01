# Swagger Annotation-Based Documentation Guide

## Overview

This project uses **swaggo/swag** for automatic Swagger/OpenAPI documentation generation from Go annotations. The documentation is generated from code comments and served via an interactive Swagger UI.

## Access URLs

Once the application is running, access the Swagger documentation at:

- **Swagger UI**: `http://localhost:8080/swagger/index.html`
- **API Documentation**: `http://localhost:8080/docs` (redirects to Swagger UI)
- **Swagger JSON**: `http://localhost:8080/swagger/doc.json`

## Architecture

### Components

1. **Annotations in Handlers**: Each HTTP handler function has Swagger annotations describing the endpoint
2. **Generated Documentation**: The `swag` CLI tool generates documentation files in `docs/`
3. **Swagger UI**: Served via `github.com/swaggo/http-swagger` package

### File Structure

```
cbe-supper-app-cps-action/
├── cmd/
│   └── main.go                          # Main API metadata annotations
├── docs/                                 # Generated Swagger documentation
│   ├── docs.go                          # Generated Go code
│   ├── swagger.json                     # OpenAPI JSON spec
│   └── swagger.yaml                     # OpenAPI YAML spec
├── internal/
│   ├── constants/dto/                   # Data Transfer Objects with comments
│   │   ├── account_block/
│   │   │   ├── request_dto.go
│   │   │   └── response_dto.go
│   │   └── account_validation/
│   │       ├── request.go
│   │       └── response.go
│   └── handlers/rest/http/              # HTTP handlers with Swagger annotations
│       ├── account_block/
│       │   └── account_block.go         # Account Block handlers
│       └── account_validation/
│           └── account_validation.go    # Account Validation handlers
└── initiator/
    └── route.go                         # Routes including Swagger endpoints
```

## Documented Modules

### Account Block Module

**16 Endpoints** organized into 4 categories:

#### Branches
- `GET /account_block/branches` - Get all branches (paginated)
- `GET /account_block/branches/{branch_code}` - Get branch by code
- `POST /account_block/branches/enable` - Enable branches
- `POST /account_block/branches/disable` - Disable branches

#### Regions
- `GET /account_block/regions` - Get all regions (paginated)
- `GET /account_block/regions/{region_code}` - Get region by code
- `POST /account_block/regions/enable` - Enable regions
- `POST /account_block/regions/disable` - Disable regions

#### Districts
- `GET /account_block/districts` - Get all districts (paginated)
- `GET /account_block/districts/{district_code}` - Get district by code
- `POST /account_block/districts/enable` - Enable districts
- `POST /account_block/districts/disable` - Disable districts

#### Cities
- `GET /account_block/city` - Get all cities (paginated)
- `GET /account_block/city/{city_code}` - Get city by code
- `POST /account_block/city/enable` - Enable cities
- `POST /account_block/city/disable` - Disable cities

### Account Validation Module

**3 Endpoints**:

- `GET /account_validation` - Get all validation rules (paginated)
- `GET /account_validation/{id}` - Get validation rule by ID
- `PATCH /account_validation/update/{id}` - Update validation rule

## How to Add Swagger Annotations

### 1. Main API Metadata (cmd/main.go)

```go
// @title CPS Action API
// @version 1.0.0
// @description API documentation for CPS Action service

// @contact.name API Support
// @contact.email contact@eaglelionsystems.com

// @host localhost:8080
// @BasePath /api/v1/cbesuperapp/cps_action

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
    initiator.Init(context.Background())
}
```

### 2. Handler Annotations

Add annotations above each handler function:

```go
// GetBranchByCode godoc
// @Summary Get branch by code
// @Description Retrieve a specific branch by its branch code
// @Tags Account Block - Branches
// @Accept json
// @Produce json
// @Param branch_code path string true "Branch Code"
// @Success 200 {object} map[string]interface{} "Branch retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /account_block/branches/{branch_code} [get]
func (a *accountBlockAdapter) GetBranchByCode(w http.ResponseWriter, r *http.Request) {
    // Handler implementation
}
```

### 3. DTO Comments

Add comments to all DTOs for better documentation:

```go
// EnableOrDisableBranches represents the request to enable or disable branches
type EnableOrDisableBranches struct {
    BranchCodes []string `json:"branches_code" example:"BR001,BR002"`
}
```

## Annotation Reference

### Common Annotations

| Annotation | Description | Example |
|------------|-------------|---------|
| `@Summary` | Short description | `@Summary Get all branches` |
| `@Description` | Detailed description | `@Description Retrieve all branches with pagination` |
| `@Tags` | Group endpoints | `@Tags Account Block - Branches` |
| `@Accept` | Request content type | `@Accept json` |
| `@Produce` | Response content type | `@Produce json` |
| `@Param` | Parameter definition | `@Param id path string true "ID"` |
| `@Success` | Success response | `@Success 200 {object} ResponseType` |
| `@Failure` | Error response | `@Failure 400 {object} ErrorType` |
| `@Security` | Security requirement | `@Security BearerAuth` |
| `@Router` | Route path and method | `@Router /path [get]` |

### Parameter Types

```go
// Path parameter
// @Param id path string true "User ID"

// Query parameter
// @Param page query int false "Page number" default(1)

// Body parameter (use full package path, not alias)
// @Param request body accountblock.EnableOrDisableBranches true "Request body"

// Header parameter
// @Param Authorization header string true "Bearer token"
```

### Response Types

```go
// Using struct type (full package path)
// @Success 200 {object} accountblock.BranchResponse

// Using generic map
// @Success 200 {object} map[string]interface{}

// Using array
// @Success 200 {array} accountblock.BranchResponse
```

## Generating Documentation

### Prerequisites

Install the swag CLI tool:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### Generate Documentation

Run this command from the project root:

```bash
swag init -g cmd/main.go --output docs
```

This will:
1. Parse all annotations in your code
2. Generate `docs/docs.go`, `docs/swagger.json`, and `docs/swagger.yaml`
3. Make the documentation available at runtime

### When to Regenerate

Regenerate documentation whenever you:
- Add new endpoints
- Modify existing endpoint annotations
- Change request/response structures
- Update API metadata

## Important Notes

### ⚠️ Package Path Requirements

**CRITICAL**: When referencing DTOs in annotations, always use the **full package name**, not import aliases:

❌ **Wrong** (using alias):
```go
import ab_dto "cbe-super-app-cps-action/internal/constants/dto/account_block"

// @Param request body ab_dto.EnableOrDisableBranches true "Request"
```

✅ **Correct** (using package name):
```go
import ab_dto "cbe-super-app-cps-action/internal/constants/dto/account_block"

// @Param request body accountblock.EnableOrDisableBranches true "Request"
```

### Best Practices

1. **Consistent Tagging**: Use consistent tag names to group related endpoints
2. **Descriptive Summaries**: Keep summaries concise but informative
3. **Detailed Descriptions**: Provide context and usage information
4. **Example Values**: Add `example` tags to struct fields for better documentation
5. **Error Responses**: Document all possible error responses
6. **Security**: Always specify `@Security BearerAuth` for protected endpoints

### Common Issues

**Issue**: `swag init` fails with parse errors
- **Solution**: Check that all DTO references use package names, not aliases

**Issue**: Swagger UI shows "Failed to load API definition"
- **Solution**: Ensure `docs` package is imported in `route.go`: `_ "cbe-super-app-cps-action/docs"`

**Issue**: Changes not reflected in Swagger UI
- **Solution**: Regenerate docs with `swag init` and restart the application

## Testing the API

1. **Start the application**:
   ```bash
   go run cmd/main.go
   ```

2. **Open Swagger UI**:
   Navigate to `http://localhost:8080/swagger/index.html`

3. **Authenticate**:
   - Click the "Authorize" button
   - Enter: `Bearer <your-jwt-token>`
   - Click "Authorize"

4. **Test endpoints**:
   - Select an endpoint
   - Click "Try it out"
   - Fill in parameters
   - Click "Execute"

## Integration with CI/CD

Add to your build pipeline:

```bash
# Install swag
go install github.com/swaggo/swag/cmd/swag@latest

# Generate docs
swag init -g cmd/main.go --output docs

# Build application
go build -o app cmd/main.go
```

## Additional Resources

- [Swaggo Documentation](https://github.com/swaggo/swag)
- [OpenAPI Specification](https://swagger.io/specification/)
- [Swagger UI](https://swagger.io/tools/swagger-ui/)

## Support

For issues or questions about the Swagger documentation:
- Check the [Swaggo GitHub Issues](https://github.com/swaggo/swag/issues)
- Review the annotation examples in existing handlers
- Contact the development team

---

**Last Updated**: 2025-10-01
**Swagger Version**: OpenAPI 3.0
**Generated By**: swaggo/swag
