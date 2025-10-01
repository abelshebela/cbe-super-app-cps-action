# Swagger Documentation Setup Guide

## Overview

This project now includes comprehensive Swagger/OpenAPI documentation for the Account Block API. The documentation is interactive and can be accessed through a web browser.

## Quick Start

### 1. Start the Application

Run your application as usual:
```bash
go run cmd/main.go
```

### 2. Access Swagger UI

Open your browser and navigate to:
```
http://localhost:<YOUR_PORT>/swagger/
```

Or use the shortcut:
```
http://localhost:<YOUR_PORT>/docs
```

### 3. Authenticate

1. Click the **"Authorize"** button (lock icon) in the top right
2. Enter your JWT token in the format: `Bearer <your-token>`
3. Click **"Authorize"** to save
4. Click **"Close"**

### 4. Test Endpoints

1. Expand any endpoint (e.g., GET /account_block/branches)
2. Click **"Try it out"**
3. Fill in any required parameters
4. Click **"Execute"**
5. View the response below

## What's Included

### Files Created

1. **`docs/swagger/account_block.yaml`**
   - Complete OpenAPI 3.0 specification for Account Block API
   - All 16 endpoints documented with request/response schemas
   - Authentication and authorization requirements

2. **`initiator/swagger.go`**
   - Swagger UI handler implementation
   - Serves both UI and YAML specifications
   - Security features (prevents directory traversal)
   - Logging integration

3. **`docs/swagger/README.md`**
   - Comprehensive documentation guide
   - Endpoint descriptions
   - Request/response examples
   - Authentication instructions

4. **`initiator/route.go`** (Updated)
   - Added Swagger UI routes
   - Integrated with existing routing

## API Endpoints Documented

### Account Block API (`/api/v1/cbesuperapp/cps_action/account_block`)

#### Branches
- `GET /branches` - Get all branches
- `GET /branches/{branch_code}` - Get branch by code
- `POST /branches/enable` - Enable branches
- `POST /branches/disable` - Disable branches

#### Regions
- `GET /regions` - Get all regions
- `GET /regions/{region_code}` - Get region by code
- `POST /regions/enable` - Enable regions
- `POST /regions/disable` - Disable regions

#### Districts
- `GET /districts` - Get all districts
- `GET /districts/{district_code}` - Get district by code
- `POST /districts/enable` - Enable districts
- `POST /districts/disable` - Disable districts

#### Cities
- `GET /city` - Get all cities
- `GET /city/{city_code}` - Get city by code
- `POST /city/enable` - Enable cities
- `POST /city/disable` - Disable cities

## Features

### Interactive Documentation
- **Live Testing**: Test endpoints directly from the browser
- **Request Builder**: Automatically builds requests with proper formatting
- **Response Viewer**: View formatted responses with syntax highlighting
- **Schema Explorer**: Browse all data models and their properties

### Security
- **JWT Authentication**: Built-in authentication support
- **Role-Based Access**: Documents required roles for each endpoint
- **Secure File Serving**: Prevents directory traversal attacks

### Developer Experience
- **Auto-completion**: Swagger UI provides auto-completion for requests
- **Validation**: Client-side validation of request parameters
- **Examples**: Pre-filled examples for quick testing
- **Persistent Auth**: Authentication token persists across page refreshes

## Data Models

All DTOs from `internal/constants/dto/account_block/` are documented:

- **BranchResponse** - Branch information with enable/disable status
- **RegionResponse** - Region information with enable/disable status
- **DistrictResponse** - District information with hierarchical data
- **CityResponse** - City information with hierarchical data
- **EnableOrDisableBranches** - Request to enable/disable branches
- **EnableOrDisableRegions** - Request to enable/disable regions
- **EnableOrDisableDistricts** - Request to enable/disable districts
- **EnableOrDisableCities** - Request to enable/disable cities

## Authorization Roles

The documentation includes role requirements:

- **Maker** - Can create and modify resources
- **IFBMaker** - IFB-specific maker role
- **Checker** - Can review and approve resources
- **IFBChecker** - IFB-specific checker role

## Adding More API Documentation

To document additional APIs:

1. Create a new YAML file in `docs/swagger/`
2. Follow the OpenAPI 3.0 specification
3. Use the base path: `/api/v1/cbesuperapp/cps_action`
4. Update the Swagger UI to include multiple specs (if needed)

Example structure:
```yaml
openapi: 3.0.3
info:
  title: Your API Name
  version: 1.0.0
servers:
  - url: /api/v1/cbesuperapp/cps_action
paths:
  /your/endpoint:
    get:
      summary: Your endpoint description
      # ... rest of the specification
```

## Troubleshooting

### Swagger UI Not Loading
- Ensure the application is running
- Check that port is correct in the URL
- Verify `docs/swagger/` directory exists
- Check application logs for errors

### YAML File Not Found
- Verify the YAML file exists in `docs/swagger/`
- Check file permissions
- Ensure file has `.yaml` or `.yml` extension

### Authentication Failing
- Ensure JWT token is valid
- Check token format: `Bearer <token>`
- Verify token hasn't expired
- Check user has required roles

### CORS Issues
- CORS is already configured in the Swagger handler
- Check browser console for specific errors
- Verify middleware configuration in `route.go`

## Best Practices

1. **Keep Documentation Updated**: Update Swagger specs when APIs change
2. **Use Examples**: Provide realistic examples in the documentation
3. **Document Errors**: Include all possible error responses
4. **Version Control**: Track API versions in the documentation
5. **Test Regularly**: Use Swagger UI to test endpoints during development

## Next Steps

1. **Add More APIs**: Document other endpoints (wallet, budget, etc.)
2. **API Versioning**: Implement version-specific documentation
3. **Custom Themes**: Customize Swagger UI appearance
4. **Export Options**: Add Postman/Insomnia collection export
5. **CI/CD Integration**: Validate OpenAPI specs in CI pipeline

## Support

For questions or issues:
- Check `docs/swagger/README.md` for detailed API documentation
- Review the OpenAPI specification files
- Contact the development team

## Resources

- [OpenAPI Specification](https://swagger.io/specification/)
- [Swagger UI Documentation](https://swagger.io/tools/swagger-ui/)
- [OpenAPI Generator](https://openapi-generator.tech/)
