# Action Registry Usage Example

## GetActionNameFromPath Function

The `GetActionNameFromPath` function in `internal/handlers/middleware/action_registry.go` allows you to resolve action names from HTTP methods and paths using the CPS action registry.

### Usage Examples

```go
package main

import (
    "fmt"
    "cbe-super-app-cps-action/internal/handlers/middleware"
)

func main() {
    // Example 1: GET banks
    actionName := middleware.GetActionNameFromPath("GET", "/api/v1/cbesuperapp/cps_action/banks")
    fmt.Printf("Action: %s\n", actionName) // Output: BANK

    // Example 2: POST notifications
    actionName = middleware.GetActionNameFromPath("POST", "/api/v1/cbesuperapp/cps_action/notifications")
    fmt.Printf("Action: %s\n", actionName) // Output: NOTIFICATIONS

    // Example 3: PATCH banks/{id}
    actionName = middleware.GetActionNameFromPath("PATCH", "/api/v1/cbesuperapp/cps_action/banks/123")
    fmt.Printf("Action: %s\n", actionName) // Output: BANK

    // Example 4: DELETE news/category/{id}
    actionName = middleware.GetActionNameFromPath("DELETE", "/api/v1/cbesuperapp/cps_action/news/category/123")
    fmt.Printf("Action: %s\n", actionName) // Output: NEWSCATEGORY

    // Example 5: POST fayda_account/disable/{user_code}
    actionName = middleware.GetActionNameFromPath("POST", "/api/v1/cbesuperapp/cps_action/fayda_account/disable/USER123")
    fmt.Printf("Action: %s\n", actionName) // Output: FAYDA

    // Example 6: Unknown path
    actionName = middleware.GetActionNameFromPath("GET", "/api/v1/cbesuperapp/cps_action/unknown")
    fmt.Printf("Action: %s\n", actionName) // Output: (empty string)
}
```

### How It Works

1. **Path Normalization**: Removes the base API prefix `/api/v1/cbesuperapp/cps_action`
2. **Special Cases**: Handles special endpoints like news/category, news/tag, fayda_account
3. **Exact Match**: First tries to find an exact match in the registry
4. **Pattern Matching**: Uses `strictAvatarMatch` to handle dynamic segments like `{id}`
5. **Resource Matching**: Falls back to resource-based matching if no exact match found

### Registry Mapping

The function uses the `cpsActionRegistry` map which contains mappings like:

```go
"POST notifications":   "NOTIFICATIONS",
"PATCH notifications":  "NOTIFICATIONS", 
"DELETE notifications": "NOTIFICATIONS",
"POST banks":           "BANK",
"PATCH banks":          "BANK",
"DELETE banks":         "BANK",
"DELETE news/category/{id}": "NEWSCATEGORY",
"POST /fayda_account/disable/{user_code}": "FAYDA",
```

### Integration with Role Validation

The `ValidateRequiredRoles` middleware now uses this function to dynamically resolve action names from request paths, providing more accurate role-based access control.

```go
// In the middleware, it's used like this:
actionName := GetActionNameFromPath(r.Method, r.URL.Path)
if actionName != "" {
    result, err := repo.FindByRoleAndAction(ctx, roleCode, actionName, 1)
    // ... validate role has viewer/maker/checker/auditor indices
}
```
