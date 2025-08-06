# Refactored Set-PIN Process Documentation

## Overview

This document describes the refactored implementation of the `/api/v1/cbesuperapp/user/set-pin` endpoint and its related processes: device-lookup and verify-otp. The refactoring focuses on improving code quality, error handling, security, and maintainability.

## Architecture Overview

The refactored implementation follows a clean architecture pattern with clear separation of concerns:

```
HTTP Layer (Adapter) → Application Layer → Domain Layer → Repository Layer
```

## Key Components

### 1. HTTP Layer (`internal/adapter/inbound/http/users/users.go`)

#### SetPin Handler
- **Endpoint**: `POST /api/v1/cbesuperapp/user/set-pin`
- **Authentication**: Required (JWT token)
- **Request Body**:
```json
{
  "new_pin": "123456",
  "otp": "654321",
  "device_uuid": "optional-device-uuid",
  "user_realm": "member",
  "otp_for": "pin_set"
}
```

**Key Improvements**:
- Enhanced input validation with proper error messages
- Better error status code mapping
- Improved request structure validation
- Default value handling for optional fields

#### DeviceLookup Handler
- **Endpoint**: `GET /api/v1/cbesuperapp/user/device-lookup`
- **Authentication**: Not required
- **Headers Required**:
  - `device-uuid`
  - `platform`
  - `app-version`

**Key Improvements**:
- Comprehensive header validation
- Better error handling for missing headers
- Improved response structure
- Enhanced logging

#### VerifyOtp Handler
- **Endpoint**: `POST /api/v1/cbesuperapp/user/sms/verify-otp`
- **Authentication**: Required (JWT token)
- **Request Body**:
```json
{
  "otp_code": "654321",
  "otp_for": "pin_set"
}
```

**Key Improvements**:
- Enhanced input validation
- Better error handling with specific error codes
- Improved security with OTP clearing
- Comprehensive logging

### 2. Application Layer (`internal/application/users/impl.go`)

#### CheckDevice Function
**Purpose**: Validates device registration and handles unverified users

**Process Flow**:
1. **Input Validation**: Validates device UUID, platform, and app version
2. **HQ Data Fetch**: Retrieves system configuration for version comparison
3. **User Lookup**: Finds user by device UUID
4. **Version Check**: Compares app version with latest versions
5. **Token Generation**: Creates temporary token for device lookup
6. **OTP Generation**: For unverified users, generates and sends OTP

**Key Improvements**:
- Separated concerns into helper functions
- Better error handling and logging
- Improved OTP generation process
- Enhanced security with proper token generation

#### VerifyOtp Function
**Purpose**: Verifies OTP and generates temporary access token

**Process Flow**:
1. **Context Validation**: Ensures context is valid
2. **Input Validation**: Validates all required parameters
3. **OTP Encryption**: Encrypts OTP for comparison
4. **User Context**: Retrieves user information from context
5. **OTP Verification**: Verifies OTP with domain service
6. **Token Generation**: Generates temporary access token

**Key Improvements**:
- Enhanced input validation
- Better error handling with specific error messages
- Improved security with OTP clearing
- Comprehensive logging

### 3. Domain Layer (`internal/domain/users/service.go`)

#### SetPin Function
**Purpose**: Sets new PIN after OTP verification

**Process Flow**:
1. **Input Validation**: Validates all input parameters
2. **OTP Verification**: Verifies OTP before proceeding
3. **PIN Validation**: Validates new PIN format and strength
4. **User Fetch**: Retrieves user information
5. **PIN History Check**: Ensures new PIN is not in history
6. **PIN Update**: Updates user's PIN and PIN history

**Key Improvements**:
- Separated validation logic into helper functions
- Enhanced PIN history management
- Better error handling with specific error codes
- Improved security checks

#### ValidatePin Function
**Purpose**: Validates PIN format and strength

**Validation Rules**:
- Must be exactly 6 digits
- Must contain only numeric characters
- Cannot have more than 4 consecutive identical digits
- Cannot be sequential (ascending or descending)

#### VerifyOtp Function
**Purpose**: Verifies OTP validity and expiration

**Process Flow**:
1. **OTP Lookup**: Finds OTP record by user ID and purpose
2. **Expiration Check**: Validates OTP expiration time
3. **OTP Comparison**: Compares provided OTP with stored OTP
4. **OTP Cleanup**: Deletes used OTP record

## Error Handling

### Error Categories

1. **Validation Errors** (400 Bad Request):
   - `MISSING_REQUIRED_FIELDS`
   - `INVALID_JSON_PAYLOAD`
   - `INVALID_PIN`
   - `MISSING_OTP`

2. **Authentication Errors** (401 Unauthorized):
   - `UNAUTHORIZED`

3. **Not Found Errors** (404 Not Found):
   - `NOT_FOUND`
   - `DEVICE_NOT_FOUND`

4. **Expired Resource Errors** (410 Gone):
   - `EXPIRED_OTP`

5. **Server Errors** (500 Internal Server Error):
   - `DEVICE_LOOKUP_FAILED`
   - `ERROR_SETTING_PIN`
   - `OTP_CREATION_FAILED`

### Error Response Format
```json
{
  "status": 400,
  "message": "Invalid PIN format",
  "data": null
}
```

## Security Improvements

### 1. OTP Security
- OTP is encrypted before storage
- OTP is cleared from memory after use
- OTP has configurable expiration time
- OTP is deleted after successful verification

### 2. PIN Security
- PIN validation against common patterns
- PIN history to prevent reuse
- PIN format validation (6 digits only)
- PIN strength requirements

### 3. Token Security
- Temporary tokens with limited scope
- Token expiration handling
- Secure token generation

## Configuration

### Environment Variables
- `OTP_WAITING_TIME`: OTP expiration time in minutes
- `GO_ENV`: Environment (dev, uat, prod)
- `KEY`: Encryption key for OTP
- `IV`: Initialization vector for encryption

## Testing

### Test Coverage
The refactored code includes comprehensive test coverage:

1. **SetPin Tests**:
   - Successful PIN setting
   - Invalid input parameters
   - Invalid PIN format
   - Same PIN as current
   - PIN in history

2. **VerifyOtp Tests**:
   - Successful OTP verification
   - OTP not found
   - Expired OTP
   - Invalid OTP

3. **ValidatePin Tests**:
   - Valid PIN formats
   - Invalid PIN lengths
   - Non-numeric PINs
   - Sequential PINs
   - Redundant PINs

### Running Tests
```bash
go test ./tests/domain/users/ -v
```

## API Usage Examples

### 1. Device Lookup
```bash
curl -X GET "http://localhost:8080/api/v1/cbesuperapp/user/device-lookup" \
  -H "device-uuid: device-123" \
  -H "platform: android" \
  -H "app-version: 1.0.0"
```

### 2. Verify OTP
```bash
curl -X POST "http://localhost:8080/api/v1/cbesuperapp/user/sms/verify-otp" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "otp_code": "654321",
    "otp_for": "pin_set"
  }'
```

### 3. Set PIN
```bash
curl -X POST "http://localhost:8080/api/v1/cbesuperapp/user/set-pin" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "new_pin": "123456",
    "otp": "654321",
    "otp_for": "pin_set"
  }'
```

## Migration Guide

### Breaking Changes
1. **Request Structure**: Some request fields have been renamed for clarity
2. **Error Responses**: Error response format has been standardized
3. **Headers**: Additional headers are now required for device lookup

### Backward Compatibility
- The core functionality remains the same
- Existing valid requests will continue to work
- New optional fields have been added without breaking existing clients

## Performance Considerations

### Optimizations
1. **Database Queries**: Optimized queries with proper indexing
2. **Caching**: Consider implementing Redis for OTP storage
3. **Async Operations**: SMS sending is handled asynchronously
4. **Connection Pooling**: Proper database connection management

### Monitoring
- Log all critical operations
- Monitor OTP generation and verification rates
- Track PIN change success/failure rates
- Monitor device lookup performance

## Future Enhancements

### Planned Improvements
1. **Rate Limiting**: Implement rate limiting for OTP generation
2. **Audit Logging**: Enhanced audit trail for security events
3. **Multi-factor Authentication**: Support for additional authentication methods
4. **PIN Complexity**: Configurable PIN complexity requirements
5. **Device Management**: Enhanced device registration and management

### Security Enhancements
1. **OTP Throttling**: Prevent OTP brute force attacks
2. **Device Fingerprinting**: Enhanced device identification
3. **Session Management**: Improved session handling
4. **Encryption**: Enhanced data encryption at rest and in transit

## Conclusion

The refactored set-pin process provides a robust, secure, and maintainable solution for user PIN management. The implementation follows best practices for security, error handling, and code organization, making it suitable for production use in a financial application context. 