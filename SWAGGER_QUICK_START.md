# Swagger Quick Start Guide

## 🚀 Quick Access

After starting the application, access Swagger at:

**Swagger UI**: `http://localhost:8080/swagger/index.html`

**Docs Redirect**: `http://localhost:8080/docs`

## 📋 Prerequisites

1. **Install Swag CLI** (one-time setup):
   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

2. **Verify Installation**:
   ```bash
   swag --version
   ```

## 🔄 Regenerate Documentation

Whenever you modify API endpoints or annotations, regenerate the docs:

```bash
swag init -g cmd/main.go --output docs
```

## 🏃 Run the Application

```bash
go run cmd/main.go
```

Or with Air for hot reload:

```bash
air
```

## 🔐 Authentication

1. Open Swagger UI: `http://localhost:8080/swagger/index.html`
2. Click the **"Authorize"** button (top right)
3. Enter your JWT token in the format:
   ```
   Bearer <your-jwt-token>
   ```
4. Click **"Authorize"**
5. Click **"Close"**

Now all requests will include the authentication token!

## 📚 Available Endpoints

### Account Block Module (16 endpoints)

#### Branches
- **GET** `/api/v1/cbesuperapp/cps_action/account_block/branches` - List all branches
- **GET** `/api/v1/cbesuperapp/cps_action/account_block/branches/{branch_code}` - Get branch
- **POST** `/api/v1/cbesuperapp/cps_action/account_block/branches/enable` - Enable branches
- **POST** `/api/v1/cbesuperapp/cps_action/account_block/branches/disable` - Disable branches

#### Regions
- **GET** `/api/v1/cbesuperapp/cps_action/account_block/regions` - List all regions
- **GET** `/api/v1/cbesuperapp/cps_action/account_block/regions/{region_code}` - Get region
- **POST** `/api/v1/cbesuperapp/cps_action/account_block/regions/enable` - Enable regions
- **POST** `/api/v1/cbesuperapp/cps_action/account_block/regions/disable` - Disable regions

#### Districts
- **GET** `/api/v1/cbesuperapp/cps_action/account_block/districts` - List all districts
- **GET** `/api/v1/cbesuperapp/cps_action/account_block/districts/{district_code}` - Get district
- **POST** `/api/v1/cbesuperapp/cps_action/account_block/districts/enable` - Enable districts
- **POST** `/api/v1/cbesuperapp/cps_action/account_block/districts/disable` - Disable districts

#### Cities
- **GET** `/api/v1/cbesuperapp/cps_action/account_block/city` - List all cities
- **GET** `/api/v1/cbesuperapp/cps_action/account_block/city/{city_code}` - Get city
- **POST** `/api/v1/cbesuperapp/cps_action/account_block/city/enable` - Enable cities
- **POST** `/api/v1/cbesuperapp/cps_action/account_block/city/disable` - Disable cities

### Account Validation Module (3 endpoints)

- **GET** `/api/v1/cbesuperapp/cps_action/account_validation` - List all validation rules
- **GET** `/api/v1/cbesuperapp/cps_action/account_validation/{id}` - Get validation rule
- **PATCH** `/api/v1/cbesuperapp/cps_action/account_validation/update/{id}` - Update validation rule

## 🧪 Testing an Endpoint

### Example: Get All Branches

1. Navigate to **Account Block - Branches** section
2. Click on **GET /account_block/branches**
3. Click **"Try it out"**
4. (Optional) Set query parameters:
   - `page`: 1
   - `per_page`: 10
   - `search`: (leave empty or add search term)
5. Click **"Execute"**
6. View the response below

### Example: Enable Branches

1. Navigate to **Account Block - Branches** section
2. Click on **POST /account_block/branches/enable**
3. Click **"Try it out"**
4. Edit the request body:
   ```json
   {
     "branches_code": ["BR001", "BR002"]
   }
   ```
5. Click **"Execute"**
6. View the response

## 🛠️ Adding New Endpoints

### Step 1: Add Annotations to Handler

```go
// GetExample godoc
// @Summary Get example data
// @Description Retrieve example data by ID
// @Tags Examples
// @Accept json
// @Produce json
// @Param id path string true "Example ID"
// @Success 200 {object} map[string]interface{} "Success"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Security BearerAuth
// @Router /examples/{id} [get]
func (h *Handler) GetExample(w http.ResponseWriter, r *http.Request) {
    // Implementation
}
```

### Step 2: Regenerate Documentation

```bash
swag init -g cmd/main.go --output docs
```

### Step 3: Restart Application

```bash
go run cmd/main.go
```

### Step 4: Verify in Swagger UI

Refresh `http://localhost:8080/swagger/index.html` and check your new endpoint!

## 📝 Annotation Cheat Sheet

```go
// Handler godoc
// @Summary Short description
// @Description Detailed description
// @Tags Category Name
// @Accept json
// @Produce json
// @Param name location type required "description" default(value)
// @Success 200 {object} ResponseType "description"
// @Failure 400 {object} ErrorType "description"
// @Security BearerAuth
// @Router /path [method]
```

### Parameter Locations
- `path` - URL path parameter
- `query` - URL query parameter
- `body` - Request body
- `header` - HTTP header

### HTTP Methods
- `[get]`, `[post]`, `[put]`, `[patch]`, `[delete]`

## ⚠️ Common Issues

### Issue: Swagger UI shows "Failed to load API definition"

**Solution**: Make sure docs package is imported in `initiator/route.go`:
```go
import (
    _ "cbe-super-app-cps-action/docs"
)
```

### Issue: Changes not appearing

**Solution**: 
1. Regenerate docs: `swag init -g cmd/main.go --output docs`
2. Restart the application
3. Hard refresh browser (Ctrl+F5)

### Issue: Parse errors when running swag init

**Solution**: Use full package names in annotations, not import aliases:
```go
// ❌ Wrong
// @Param request body dto.Request true "Request"

// ✅ Correct
// @Param request body packagename.Request true "Request"
```

## 📦 Project Structure

```
cbe-supper-app-cps-action/
├── cmd/main.go                    # API metadata annotations
├── docs/                          # Generated Swagger files
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── initiator/route.go             # Swagger UI routes
└── internal/handlers/rest/http/   # Handler annotations
    ├── account_block/
    └── account_validation/
```

## 🔗 Useful Links

- **Swagger UI**: http://localhost:8080/swagger/index.html
- **Swagger JSON**: http://localhost:8080/swagger/doc.json
- **API Base Path**: http://localhost:8080/api/v1/cbesuperapp/cps_action

## 💡 Tips

1. **Use Tags Wisely**: Group related endpoints with consistent tag names
2. **Document Everything**: Include all parameters, responses, and error cases
3. **Test in Swagger**: Use the interactive UI to test endpoints before integration
4. **Keep Docs Updated**: Regenerate after every API change
5. **Add Examples**: Use `example` tags in struct fields for better documentation

## 🎯 Next Steps

1. ✅ Access Swagger UI
2. ✅ Authenticate with Bearer token
3. ✅ Test existing endpoints
4. ✅ Add annotations to new handlers
5. ✅ Regenerate documentation
6. ✅ Share with team!

---

**Need Help?** Check `SWAGGER_ANNOTATION_GUIDE.md` for detailed documentation.
