# Donation Category Feature Pattern Guide

This document outlines the complete implementation pattern for the Donation Category feature, serving as a template for building similar features in the CBE Super App CPS Action system.

## Table of Contents
1. [Architecture Overview](#architecture-overview)
2. [Layer-by-Layer Implementation](#layer-by-layer-implementation)
3. [Code Patterns and Rules](#code-patterns-and-rules)
4. [File Structure](#file-structure)
5. [Implementation Checklist](#implementation-checklist)
6. [Best Practices](#best-practices)

## Architecture Overview

The Donation Category feature follows a clean architecture pattern with clear separation of concerns:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP Layer    │───▶│  Service Layer  │───▶│ Persistence     │
│  (Handlers)     │    │   (Business)    │    │   (Storage)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Routing       │    │   Interfaces    │    │   Repository    │
│  (Endpoints)    │    │  (Contracts)    │    │  (Data Access)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Layer-by-Layer Implementation

### 1. Model Layer (`internal/constants/model/`)

**File:** `donation_category.go`

```go
type DonationCategory struct {
    ID             bson.ObjectID `bson:"_id,omitempty" json:"id"`
    CategoryName   string        `json:"category_name" bson:"category_name"`
    Icon           string        `json:"donation_icon" bson:"donation_icon"`
    IsDeleted      bool          `json:"is_deleted" bson:"is_deleted"`
    CreatedAt      time.Time     `json:"created_at" bson:"created_at"`
    LastModifiedAt time.Time     `json:"last_modified_at" bson:"last_modified_at"`
}
```

**Rules:**
- Use `bson.ObjectID` for MongoDB IDs
- Include standard audit fields: `IsDeleted`, `CreatedAt`, `LastModifiedAt`
- Use consistent naming: snake_case for database, camelCase for JSON
- Add appropriate BSON and JSON tags

### 2. DTO Layer (`internal/constants/dto/donation_category/`)

#### Request DTO (`request_dto.go`)
```go
type DonationCategoryRequest struct {
    CategoryName string                `json:"category_name" bson:"category_name"`
    Icon         *multipart.FileHeader `json:"donation_icon,omitempty" bson:"donation_icon,omitempty"`
}

type DonationCategoryCPSRequest struct {
    ID           string `json:"id,omitempty" bson:"id,omitempty"`
    CategoryName string `json:"category_name" bson:"category_name"`
    Icon         string `json:"donation_icon,omitempty" bson:"donation_icon,omitempty"`
}
```

#### Response DTO (`response_dto.go`)
```go
type DonationCategoryListResponse struct {
    ID             string `json:"id" bson:"id"`
    CategoryName   string `json:"category_name" bson:"category_name"`
    Icon           string `json:"donation_icon" bson:"donation_icon"`
    IsDeleted      bool   `json:"is_deleted" bson:"is_deleted"`
    CreatedAt      string `json:"created_at" bson:"created_at"`
    LastModifiedAt string `json:"last_modified_at" bson:"last_modified_at"`
}

type DonationCategoryResponse struct {
    CategoryName string `json:"category_name" bson:"category_name"`
    Icon         string `json:"donation_icon" bson:"donation_icon"`
}
```

#### Validator (`validator.go`)
```go
func (d DonationCategoryRequest) Validate() error {
    return validation.ValidateStruct(&d,
        validation.Field(&d.CategoryName,
            validation.Required.Error("category name is required"),
            validation.Length(3, 100).Error("category name must be between 3 and 100 characters"),
            validation.By(utils.NoSpecialChars),
        ),
        validation.Field(&d.Icon,
            validation.Required.Error("category icon is required"),
            validation.By(validateImage),
        ),
    )
}

func (d DonationCategoryRequest) ValidateForUpdate() error {
    if d.CategoryName == "" && d.Icon == nil {
        return validation.NewError("validation_at_least_one_field", "at least one field must be provided for update")
    }
    // ... validation logic
}
```

**Rules:**
- Create separate DTOs for different contexts (HTTP, CPS, Internal)
- Implement validation methods for create and update operations
- Use `multipart.FileHeader` for file uploads in HTTP requests
- Use string URLs for CPS actions
- Include proper validation rules with meaningful error messages

### 3. Interface Layer (`internal/constants/interfaces/donation_category/`)

**File:** `donation_category.go`
```go
type DonationCategoryAdapter interface {
    CreateDonationCategory(w http.ResponseWriter, r *http.Request)
    UpdateDonationCategory(w http.ResponseWriter, r *http.Request)
    FetchDonationCategory(w http.ResponseWriter, r *http.Request)
    FetchDonationCategoryByID(w http.ResponseWriter, r *http.Request)
}
```

**Rules:**
- Define interfaces for each layer (Handler, Service, Repository)
- Use descriptive method names
- Follow Go interface conventions
- Keep interfaces focused and cohesive

### 4. Routing Layer (`internal/glue/routing/donation_category/`)

**File:** `donation_category.go`
```go
func Init(router chi.Router, handler donation_category.DonationCategoryAdapter, authMiddleware middleware.AuthMiddleware) {
    routes := []glue.Route{
        {
            Method:  http.MethodPost,
            Path:    "/donation_category",
            Handler: handler.CreateDonationCategory,
            Middlewares: []func(next http.Handler) http.Handler{
                authMiddleware.AuthenticateToken,
                authMiddleware.AccessControl([]string{role.Maker, role.IFBMaker}),
            },
        },
        // ... other routes
    }
    glue.RegisterRoutes(router, routes)
}
```

**Rules:**
- Use RESTful URL patterns
- Apply appropriate middleware (auth, role-based access)
- Group related routes in the same file
- Use the `glue.RegisterRoutes` pattern for consistency


### 5. Handler Layer (`internal/handlers/rest/http/donation_category/`)

#### Main Handler (`donation_category.go`)
```go
type donationCategoryAdapter struct {
    logger     utils.Logger
    donationCategoryApp service.DonationCategoryService
}

func (d *donationCategoryAdapter) CreateDonationCategory(w http.ResponseWriter, r *http.Request) {
    req, err := core.ParseRequestFromMultipartForm(r, true)
    if err != nil {
        d.logger.Errorf("failed to parse event request from multipart form: %v", err)
        localization.SendBadRequestResponse(w, err.Error())
        return
    }

    if err := req.Validate(); err != nil {
        d.logger.Errorf("event request validation failed: %v", err)
        localization.SendBadRequestResponse(w, err.Error())
        return
    }

    if err := d.donationCategoryApp.CreateDonationCategory(r.Context(), req); err != nil {
        d.logger.Errorf("failed to create event: %v", err)
        localization.SendErrorByCodeResponse(w, err.Error())
        return
    }

    d.logger.Infof("event creation request submitted successfully")
    localization.SendSuccessResponse(w, localization.SuccessEventCreationRequestSubmitted, nil)
}
```

#### Helper (`core/helper.go`)
```go
func ParseRequestFromMultipartForm(r *http.Request, isCreate bool) (dto.DonationCategoryRequest, error) {
    var req dto.DonationCategoryRequest
    _, fileHeader, err := utils.ParseMultipartFormFile(r, "donation_icon", 2<<20)
    if err != nil {
        if err.Error() != localization.ErrorMissingFile.Code || isCreate {
            return req, errors.New(localization.ErrorInvalidFileUpload.Code)
        }
    } else {
        req.Icon = fileHeader
    }
    req.CategoryName = r.FormValue("category_name")
    return req, nil
}
```

**Rules:**
- Extract request parsing logic to helper functions
- Use consistent error handling and logging
- Validate requests before passing to service layer
- Use localization for all responses
- Handle file uploads properly with size limits

### 6. Service Layer (`internal/service/donation_category/`)

#### Main Service (`donation_category.go`)
```go
type DonationCategory struct {
    DonationCategoryRepo storage.DonationCategoryRepository
    cpsService           service.CPSActionService
    logger               utils.Logger
    minio                config.MinioClientInterface
    bucketName           string
    cfg                  *config.VaultConfig
    minioEndPoint        string
}

func (d *DonationCategory) CreateDonationCategory(ctx context.Context, donationCategory dto.DonationCategoryRequest) error {
    makerData := local_util.ExtractUserFromContext(ctx)
    if incomplet := local_util.IsIncomplete(makerData); incomplet {
        return errors.New(localization.ErrorAccountNumberRequired.Code)
    }
    
    ok, err := core.DonationNameExists(ctx, donationCategory.CategoryName, d.DonationCategoryRepo)
    if err != nil {
        return err
    }
    if ok {
        return errors.New(localization.ErrorDonationCategoryNameDuplicated.Code)
    }
    
    if donationCategory.Icon == nil {
        return errors.New(localization.ErrorMissingFile.Code)
    }
    
    url, err := lib.UploadFileToMinio(ctx, d.minio, d.bucketName, donationCategory.Icon, string(constants.Avatar), d.minioEndPoint, d.logger)
    if err != nil {
        return err
    }
    
    result := dto.DonationCategoryResponse{
        CategoryName: donationCategory.CategoryName,
        Icon:         url,
    }

    cpsAction := lib.CpsModelBuilder("", makerData, "", result, string(constants.RequestCreateDonationCategory), constants.CREATE)
    if err := d.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
        return err
    }
    return nil
}
```

#### Helper (`core/helper.go`)
```go
func DonationNameExists(ctx context.Context, categoryName string, donationCategoryRepo storage.DonationCategoryRepository) (bool, error) {
    filter := bson.M{
        "is_deleted": false,
        "category_name": bson.M{
            "$regex":   fmt.Sprintf("^%s$", regexp.QuoteMeta(strings.TrimSpace(categoryName))),
            "$options": "i", // case-insensitive
        },
    }
    // ... implementation
}

func MapToDonationCategory(categoryName, iconURL string) *model.DonationCategory {
    return &model.DonationCategory{
        CategoryName:   categoryName,
        Icon:           iconURL,
        IsDeleted:      false,
        CreatedAt:      time.Now(),
        LastModifiedAt: time.Now(),
    }
}
```

**Rules:**
- Implement business logic in service layer
- Use CPS (Check, Process, Store) pattern for data modifications
- Extract reusable logic to helper functions
- Validate business rules before processing
- Handle file uploads and external service calls
- Use proper error handling with localization codes

### 7. Persistence Layer (`internal/storage/persistance/donation_category/`)

#### Repository (`repository.go`)
```go
type DonationCategoryStorage struct {
    dal    dal.MongoDal[model.DonationCategory, model.DonationCategory]
    client *mongo.Client
    logger utils.Logger
}

func (s *DonationCategoryStorage) Create(ctx context.Context, details *model.DonationCategory) error {
    _, err := s.dal.InsertOne(ctx, *details)
    if err != nil {
        return errors.New(localization.ErrorUnexpectedError.Code)
    }
    return nil
}

func (s *DonationCategoryStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.DonationCategory], error) {
    filter := bson.M{"is_deleted": false}
    searchKeys := bson.M{}
    allowedKeys := []string{"category_name"}
    
    if filterParam.Search != "" {
        searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
        searchKeys["$or"] = []bson.M{{"category_name": searchRegex}}
    }

    projection := bson.M{}
    filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

    data, err := s.dal.FindAllWithPagination(ctx, filter, projection, skip, limit)
    if err != nil {
        return nil, errors.New(localization.ErrorUnexpectedError.Message)
    }

    total, err := s.dal.TotalCount(ctx, filter)
    if err != nil {
        return nil, errors.New(localization.ErrorUnexpectedError.Message)
    }

    meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
    return &types.PaginatedResponse[[]*model.DonationCategory]{
        Data: data,
        Meta: meta,
    }, nil
}
```

#### Projection (`projection.go`)
```go
func DonationCategoryMapper(data model.DonationCategory) bson.M {
    result := bson.M{}
    if data.CategoryName != "" {
        result["category_name"] = data.CategoryName
    }
    if data.Icon != "" {
        result["donation_icon"] = data.Icon
    }
    return result
}
```

**Rules:**
- Use generic DAL (Data Access Layer) for database operations
- Implement proper pagination with search functionality
- Use projection mappers for update operations
- Handle soft deletes with `is_deleted` flag
- Use proper error handling and logging
-there is nothing you modify on storage from now on  , you just consume the existing functions on the service layers
### 8. Initiator Integration

#### Service Registration (`initiator/service.go`)
```go
type ServiceContainer struct {
    // ... other services
    DonationCategoryContainer service.DonationCategoryService
}

func NewServiceContainer(/* ... */) *ServiceContainer {
    // ... initialization
    donationCategoryService := donation_category.NewDonationCategoryService(
        mongoClient,
        persistence.DonationCategoryPersistence,
        cpsActionService,
        logger,
        minioClient,
        minioPubUrl,
        "donation_category_icon",
        cfg,
    )
    
    return &ServiceContainer{
        // ... other services
        DonationCategoryContainer: donationCategoryService,
    }
}
```

#### Handler Registration (`initiator/handler.go`)
```go
type HandlerContainer struct {
    // ... other handlers
    DonationCategoryHandler donation_category.DonationCategoryAdapter
}

func NewHandlerContainer(serviceLayer *ServiceContainer, logger utils.Logger) *HandlerContainer {
    return &HandlerContainer{
        // ... other handlers
        DonationCategoryHandler: donationCategoryHandler.InitDonationCategoryAdapter(
            serviceLayer.DonationCategory,
            logger,
        ),
    }
}
```

#### Route Registration (`initiator/route.go`)
```go
func InitRoutes(r chi.Router, handlerLayer *HandlerContainer, authMiddleware middleware.AuthMiddleware) {
    // ... other routes
    donation_category.Init(r, handlerLayer.DonationCategoryHandler, authMiddleware)
}
```

**Rules:**
- Register services in the service container
- Register handlers in the handler container
- Register routes in the route initialization
- Use dependency injection pattern
- Pass all required dependencies

## Generic Functions and Reusable Patterns

### Generic Functions Used in Donation Category Implementation

Based on the actual donation category implementation, these are the generic functions that can be reused for similar features:

#### File Upload Utilities
```go
// Generic file upload to MinIO (used in donation category)
url, err := lib.UploadFileToMinio(ctx, d.minio, d.bucketName, donationCategory.Icon, string(constants.Avatar), d.minioEndPoint, d.logger)

// Generic multipart form parsing (used in donation category)
req, err := core.ParseRequestFromMultipartForm(r, true) // true for create, false for update

// Generic file validation (used in donation category)
func validateImage(value interface{}) error {
    file, ok := value.(*multipart.FileHeader)
    if !ok {
        return validation.NewError("validation_image_invalid", "invalid image file")
    }
    // Check file size (max 10MB)
    if file.Size > 10*1024*1024 {
        return validation.NewError("validation_image_size", "image file size must not exceed 10MB")
    }
    // Check file extension
    if !hasAllowedExtension(file.Filename, []string{".jpg", ".jpeg", ".png", ".gif"}) {
        return validation.NewError("validation_image_format", "image must be JPG, JPEG, PNG, or GIF")
    }
    return nil
}
```

#### Database Operations
```go
// Generic DAL operations (used in donation category)
dal.NewMongoDal[model.DonationCategory, model.DonationCategory](client, dbName, collection)

// Generic pagination builder (used in donation category)
filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

// Generic pagination metadata (used in donation category)
meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
```

#### CPS Action Utilities
```go
// Generic CPS model builder (used in donation category)
cpsAction := lib.CpsModelBuilder("", makerData, "", result, string(constants.RequestCreateDonationCategory), constants.CREATE)

// Generic action binding (used in donation category)
err := core.BindAction(action.CurrentAction, &donationCPS)
```

#### Validation Utilities
```go
// Generic validation rules (used in donation category)
validation.By(utils.NoSpecialChars)
validation.Length(3, 100)
validation.Required

// Generic validation for update operations
func (d DonationCategoryRequest) ValidateForUpdate() error {
    if d.CategoryName == "" && d.Icon == nil {
        return validation.NewError("validation_at_least_one_field", "at least one field must be provided for update")
    }
    // ... validation logic
}
```

#### Context and User Utilities
```go
// Generic user extraction from context (used in donation category)
makerData := local_util.ExtractUserFromContext(ctx)

// Generic user completeness check (used in donation category)
if incomplet := local_util.IsIncomplete(makerData); incomplet {
    return errors.New(localization.ErrorAccountNumberRequired.Code)
}
```

#### Business Logic Utilities
```go
// Generic name existence check (used in donation category)
func DonationNameExists(ctx context.Context, categoryName string, donationCategoryRepo storage.DonationCategoryRepository) (bool, error) {
    filter := bson.M{
        "is_deleted": false,
        "category_name": bson.M{
            "$regex":   fmt.Sprintf("^%s$", regexp.QuoteMeta(strings.TrimSpace(categoryName))),
            "$options": "i", // case-insensitive
        },
    }
    // ... implementation
}

// Generic model mapping (used in donation category)
func MapToDonationCategory(categoryName, iconURL string) *model.DonationCategory {
    return &model.DonationCategory{
        CategoryName:   categoryName,
        Icon:           iconURL,
        IsDeleted:      false,
        CreatedAt:      time.Now(),
        LastModifiedAt: time.Now(),
    }
}
```

### Image Handling in Donation Category

The donation category implementation handles image uploads as part of the main create/update operations:

#### Image Upload in Create Operation
```go
// Image is handled within CreateDonationCategory
func (d *donationCategoryAdapter) CreateDonationCategory(w http.ResponseWriter, r *http.Request) {
    req, err := core.ParseRequestFromMultipartForm(r, true)
    // ... validation and processing
    // Icon is uploaded to MinIO and URL is stored
}
```

#### Image Upload in Update Operation
```go
// Image is handled within UpdateDonationCategory
func (d *donationCategoryAdapter) UpdateDonationCategory(w http.ResponseWriter, r *http.Request) {
    req, err := core.ParseRequestFromMultipartForm(r, false)
    // ... validation and processing
    // Icon can be updated as part of the update operation
}
```

## Code Patterns and Rules

### 1. Naming Conventions
- **Packages**: Use lowercase, descriptive names (`donation_category`)
- **Types**: Use PascalCase (`DonationCategory`)
- **Methods**: Use PascalCase for public, camelCase for private
- **Files**: Use snake_case (`donation_category.go`)
- **Constants**: Use UPPER_SNAKE_CASE (`RequestCreateDonationCategory`)

### 2. Error Handling
- Use localization codes for all errors
- Log errors with context
- Return meaningful error messages
- Use consistent error response format

### 3. Validation
- Validate at the handler level for input format
- Validate at the service level for business rules
- Use the `ozzo-validation` library
- Create separate validation methods for create/update

### 4. File Upload Handling
- Use `multipart.FileHeader` for HTTP requests
- Implement file size and type validation
- Upload to MinIO and store URLs
- Handle file upload errors gracefully

### 5. CPS Pattern
- All data modifications go through CPS (Check, Process, Store)
- Create CPS actions for audit trail
- Implement authorization methods
- Use proper action types and statuses

### 6. Pagination and Filtering
- Implement consistent pagination across all list endpoints
- Support search functionality
- Use proper filter builders
- Return standardized pagination metadata

## File Structure

```
internal/
├── constants/
│   ├── dto/donation_category/
│   │   ├── request_dto.go
│   │   ├── response_dto.go
│   │   └── validator.go
│   ├── interfaces/donation_category/
│   │   └── donation_category.go
│   ├── model/
│   │   └── donation_category.go
│   └── localization/
│       ├── messages.go
│       └── response_codes.go
├── glue/routing/donation_category/
│   └── donation_category.go
├── handlers/rest/http/donation_category/
│   ├── donation_category.go
│   └── core/
│       └── helper.go
├── service/donation_category/
│   ├── donation_category.go
│   └── core/
│       └── helper.go
└── storage/persistance/donation_category/
    ├── repository.go
    └── projection.go
```

## Implementation Checklist

### Phase 1: Foundation
- [ ] Create model with proper BSON/JSON tags
- [ ] Define DTOs for request/response
- [ ] Implement validation logic
- [ ] Create interfaces for all layers
- [ ] **Identify and document generic functions** that can be reused across similar endpoints

### Phase 2: Data Layer
- [ ] Implement repository with CRUD operations
- [ ] Create projection mappers
- [ ] Add pagination and search functionality
- [ ] Implement soft delete pattern
- [ ] **Use generic DAL functions** for common database operations

### Phase 3: Business Layer
- [ ] Implement service with business logic
- [ ] Add CPS action integration
- [ ] Implement file upload handling
- [ ] Add proper error handling
- [ ] **Leverage generic file upload utilities** from `lib.UploadFileToMinio`

### Phase 4: API Layer
- [ ] Create HTTP handlers
- [ ] Implement request parsing helpers
- [ ] Add proper validation
- [ ] Implement response formatting
- [ ] **Use generic multipart form parsing** from `core.ParseRequestFromMultipartForm`

### Phase 5: Integration
- [ ] Register services in initiator
- [ ] Register handlers in initiator
- [ ] Register routes in initiator
- [ ] Add localization messages and codes
- [ ] **Configure generic image upload settings** (bucket, size limits, allowed types)



## Best Practices

1. **Separation of Concerns**: Keep each layer focused on its responsibility
2. **Dependency Injection**: Use constructor injection for dependencies
3. **Error Handling**: Use consistent error handling patterns
4. **Logging**: Log important operations and errors
5. **Validation**: Validate at appropriate layers
6. **Testing**: Write tests for all business logic
7. **Documentation**: Document complex business rules
8. **Consistency**: Follow established patterns throughout
9. **Security**: Implement proper authentication and authorization
10. **Performance**: Consider pagination and filtering for large datasets

This pattern ensures maintainable, testable, and scalable code that follows the established architecture of the CBE Super App CPS Action system.
