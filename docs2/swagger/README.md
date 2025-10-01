# CBE Super App - CPS Action API Documentation

This directory contains OpenAPI 3.0 (Swagger) documentation for the CBE Super App CPS Action APIs.

## Overview

The Swagger documentation provides interactive API documentation that allows you to:
- View all available API endpoints
- Understand request/response schemas
- Test API endpoints directly from the browser
- See authentication requirements
- View example requests and responses

## Base URL

All API endpoints are prefixed with:
```
/api/v1/cbesuperapp/cps_action
```

## Accessing the Documentation

### Swagger UI
Access the interactive Swagger UI at:
```
http://localhost:<PORT>/swagger/
```
or
```
http://localhost:<PORT>/docs
```

The Swagger UI provides:
- Interactive API testing
- Request/response examples
- Schema definitions
- Authentication setup

### API Specifications

Individual API specification files are available at:
```
http://localhost:<PORT>/swagger/<spec-name>.yaml
```

## Available API Specifications

### Account Block API
**File:** `account_block.yaml`  
**Endpoints:** `/api/v1/cbesuperapp/cps_action/account_block/*`

Manages account blocking functionality for:
- **Branches** - Branch management and enable/disable operations
- **Regions** - Region management and enable/disable operations
- **Districts** - District management and enable/disable operations
- **Cities** - City management and enable/disable operations

#### Endpoints:

**Branches:**
- `GET /account_block/branches` - Get all branches
- `GET /account_block/branches/{branch_code}` - Get branch by code
- `POST /account_block/branches/enable` - Enable branches
- `POST /account_block/branches/disable` - Disable branches

**Regions:**
- `GET /account_block/regions` - Get all regions
- `GET /account_block/regions/{region_code}` - Get region by code
- `POST /account_block/regions/enable` - Enable regions
- `POST /account_block/regions/disable` - Disable regions

**Districts:**
- `GET /account_block/districts` - Get all districts
- `GET /account_block/districts/{district_code}` - Get district by code
- `POST /account_block/districts/enable` - Enable districts
- `POST /account_block/districts/disable` - Disable districts

**Cities:**
- `GET /account_block/city` - Get all cities
- `GET /account_block/city/{city_code}` - Get city by code
- `POST /account_block/city/enable` - Enable cities
- `POST /account_block/city/disable` - Disable cities

## Authentication

All endpoints require JWT Bearer token authentication. To use the API:

1. Obtain a JWT token from the authentication service
2. In Swagger UI, click the "Authorize" button
3. Enter your token in the format: `Bearer <your-token>`
4. Click "Authorize" to save

## Authorization Roles

Different endpoints require different user roles:

- **Maker** - Can create and modify resources
- **IFBMaker** - IFB-specific maker role
- **Checker** - Can review and approve resources
- **IFBChecker** - IFB-specific checker role

## Request/Response Examples

### Enable Branches
```bash
POST /api/v1/cbesuperapp/cps_action/account_block/branches/enable
Content-Type: application/json
Authorization: Bearer <your-token>

{
  "branches_code": ["BR001", "BR002", "BR003"]
}
```

### Get All Regions
```bash
GET /api/v1/cbesuperapp/cps_action/account_block/regions
Authorization: Bearer <your-token>
```

## Data Models

### BranchResponse
```json
{
  "id": "string",
  "branch_code": "string",
  "branch_name": "string",
  "branch_address": "string",
  "district_code": "string",
  "district_name": "string",
  "branch_region": "string",
  "record_stat": "string",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "enabled": true
}
```

### RegionResponse
```json
{
  "id": "string",
  "region_code": "string",
  "region_name": "string",
  "region_address": "string",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "enabled": true
}
```

### DistrictResponse
```json
{
  "id": "string",
  "district_code": "string",
  "district_name": "string",
  "district_address": "string",
  "region_id": "string",
  "region_name": "string",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "enabled": true
}
```

### CityResponse
```json
{
  "id": "string",
  "city_code": "string",
  "city_name": "string",
  "city_address": "string",
  "district_id": "string",
  "district_name": "string",
  "region_id": "string",
  "region_name": "string",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "enabled": true
}
```

## Error Responses

All endpoints return consistent error responses:

```json
{
  "error": "Error message describing what went wrong",
  "code": "ERROR_CODE"
}
```

Common HTTP status codes:
- `200` - Success
- `400` - Bad Request (invalid input)
- `401` - Unauthorized (missing or invalid token)
- `403` - Forbidden (insufficient permissions)
- `404` - Not Found (resource doesn't exist)
- `500` - Internal Server Error

## Development

### Adding New API Specifications

1. Create a new YAML file in this directory
2. Follow the OpenAPI 3.0 specification format
3. Use the base path: `/api/v1/cbesuperapp/cps_action`
4. Include proper authentication and authorization requirements
5. Document all request/response schemas
6. Update this README with the new specification details

### Updating Existing Specifications

When updating API specifications:
1. Maintain backward compatibility when possible
2. Document breaking changes clearly
3. Update version numbers appropriately
4. Test changes in Swagger UI before committing

## Testing

Use the Swagger UI to test endpoints:
1. Navigate to http://localhost:<PORT>/swagger/
2. Click "Authorize" and enter your JWT token
3. Select an endpoint to test
4. Click "Try it out"
5. Fill in required parameters
6. Click "Execute" to send the request
7. View the response below

## Support

For issues or questions about the API documentation:
- Contact: CBE Super App Team
- Email: support@cbe.com

## Version History

- **v1.0.0** - Initial release with Account Block API documentation
