# 📚 Swagger API Documentation

## 🚀 Quick Start

### Access Swagger UI
```
http://localhost:8080/swagger/index.html
```

### Authenticate
1. Click **"Authorize"** button
2. Enter: `Bearer <your-jwt-token>`
3. Click **"Authorize"**

### Test Endpoints
- Browse available endpoints
- Click "Try it out"
- Execute requests
- View responses

---

## 📖 Documentation Files

| File | Purpose |
|------|---------|
| **SWAGGER_QUICK_START.md** | Quick reference guide - start here! |
| **SWAGGER_ANNOTATION_GUIDE.md** | Comprehensive documentation with examples |
| **SWAGGER_IMPLEMENTATION_SUMMARY.md** | Technical implementation details |

---

## 🔄 Regenerate Documentation

After modifying API endpoints or annotations:

```bash
swag init -g cmd/main.go --output docs
```

Then restart the application.

---

## 📊 Currently Documented

### ✅ Account Block Module (16 endpoints)
- **Branches**: List, Get, Enable, Disable
- **Regions**: List, Get, Enable, Disable
- **Districts**: List, Get, Enable, Disable
- **Cities**: List, Get, Enable, Disable

### ✅ Account Validation Module (3 endpoints)
- List validation rules
- Get validation rule by ID
- Update validation rule

---

## 🛠️ Adding New Endpoints

### 1. Add Annotations
```go
// HandlerName godoc
// @Summary Short description
// @Description Detailed description
// @Tags Module Name
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} ResponseType
// @Security BearerAuth
// @Router /path/{id} [get]
func (h *Handler) HandlerName(w http.ResponseWriter, r *http.Request) {
    // Implementation
}
```

### 2. Regenerate
```bash
swag init -g cmd/main.go --output docs
```

### 3. Restart & Test
```bash
go run cmd/main.go
```

Visit: `http://localhost:8080/swagger/index.html`

---

## ⚠️ Important Rules

### Use Full Package Names
❌ **Wrong**: `@Param request body dto.Type`  
✅ **Correct**: `@Param request body packagename.Type`

### Regenerate After Changes
- Adding new endpoints
- Modifying annotations
- Changing DTOs
- Updating API metadata

---

## 🔗 Useful URLs

| Resource | URL |
|----------|-----|
| Swagger UI | http://localhost:8080/swagger/index.html |
| Docs Redirect | http://localhost:8080/docs |
| Swagger JSON | http://localhost:8080/swagger/doc.json |
| API Base | http://localhost:8080/api/v1/cbesuperapp/cps_action |

---

## 📚 Learn More

- **Quick Start**: Read `SWAGGER_QUICK_START.md`
- **Detailed Guide**: Read `SWAGGER_ANNOTATION_GUIDE.md`
- **Implementation**: Read `SWAGGER_IMPLEMENTATION_SUMMARY.md`

---

## 🎯 Benefits

✅ Documentation lives with code  
✅ Always in sync  
✅ Type-safe  
✅ Auto-generated  
✅ Interactive testing  
✅ Easy to maintain  

---

**Status**: ✅ Production Ready  
**Version**: OpenAPI 3.0  
**Generator**: swaggo/swag
