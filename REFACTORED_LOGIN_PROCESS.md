# Refactored Login Process Documentation

## Overview

This document describes the refactored login process for the CBE Super App, providing secure authentication with comprehensive validation, rate limiting, and proper error handling.

## Architecture

The login process follows a clean architecture pattern with clear separation of concerns:

```
HTTP Layer (Adapter) → Application Layer → Domain Layer → Repository Layer
```

### Components

- **HTTP Adapter**: Handles HTTP requests/responses, validation, and uses `utils.BaseResponseMaker`
- **Application Layer**: Business logic orchestration and validation
- **Domain Layer**: Core business rules and entity management
- **Repository Layer**: Data persistence and retrieval

## Login Flow

### Login Endpoint (`/api/v1/cbesuperapp/user/login`)

**Purpose**: Authenticates users with phone number, device UUID, and PIN

**Request**:
```json
{
  "phone": "1234567890",
  "device_uuid": "device-uuid-12345",
  "pin": "123456"
}
```

**Headers Required**:
- `device-uuid`: Device UUID (must match request body)
- `platform`: android, ios, or web
- `app-version`: Application version

**Process**:
1. Validates input parameters (phone, device UUID, PIN)
2. Checks device UUID consistency between header and body
3. Finds user by phone number
4. Validates user status (not deleted, not blocked)
5. Verifies device is linked to user
6. Checks login attempt limits
7. Validates PIN
8. Generates JWT access token
9. Updates last login time and resets login attempts
10. Returns comprehensive login response

**Response**:
```json
{
  "status": 200,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user_id": "507f1f77bcf86cd799439011",
    "user_code": "12345678",
    "full_name": "John Doe",
    "phone_number": "1234567890",
    "kyc_level": 1,
    "is_verified": true,
    "device_uuid": "device-uuid-12345",
    "platform": "android",
    "app_version": "1.0.0",
    "login_time": "2024-01-15T10:30:00Z",
    "session_expires": "2024-01-16T10:30:00Z",
    "last_login": "2024-01-14T15:45:00Z",
    "login_attempts": 0
  }
}
```

## Data Models

### LoginRequest
```go
type LoginRequest struct {
    Phone      string `json:"phone" validate:"required"`
    DeviceUUID string `json:"device_uuid" validate:"required"`
    Pin        string `json:"pin" validate:"required,min=6,max=6"`
}
```

### LoginResponse
```go
type LoginResponse struct {
    Token           string    `json:"token"`
    UserID          string    `json:"user_id"`
    UserCode        string    `json:"user_code"`
    FullName        string    `json:"full_name"`
    PhoneNumber     string    `json:"phone_number"`
    KYCLevel        uint8     `json:"kyc_level"`
    IsVerified      bool      `json:"is_verified"`
    DeviceUUID      string    `json:"device_uuid"`
    Platform        string    `json:"platform"`
    AppVersion      string    `json:"app_version"`
    LoginTime       time.Time `json:"login_time"`
    SessionExpires  time.Time `json:"session_expires"`
    LastLogin       time.Time `json:"last_login"`
    LoginAttempts   int       `json:"login_attempts"`
}
```

## Error Handling

### Login Errors
- `USER_NOT_FOUND`: User with phone number not found
- `ACCOUNT_DELETED`: User account has been deleted
- `ACCOUNT_BLOCKED`: User account is blocked
- `DEVICE_NOT_LINKED`: Device not linked to user account
- `TOO_MANY_LOGIN_ATTEMPTS`: Too many failed login attempts (rate limited)
- `PIN_NOT_SET`: User has no PIN set
- `INVALID_PIN`: Incorrect PIN provided
- `TOKEN_GENERATION_FAILED`: Failed to generate access token
- `INVALID_PHONE_NUMBER`: Invalid phone number format
- `INVALID_DEVICE_UUID`: Invalid device UUID format

### HTTP Status Codes
- `200`: Login successful
- `400`: Bad request (validation errors)
- `401`: Unauthorized (invalid credentials)
- `403`: Forbidden (account blocked, device not linked)
- `404`: Not found (user not found)
- `410`: Gone (account deleted)
- `429`: Too many requests (rate limited)
- `500`: Internal server error

## Security Features

### 1. Input Validation
- Phone number format validation
- PIN format validation (6 digits, numeric only)
- Device UUID format validation
- Device UUID consistency check between header and body

### 2. Account Security
- Account deletion check
- Account blocking check
- Device linking verification
- PIN validation

### 3. Rate Limiting
- Maximum 5 login attempts
- 15-minute cooldown period after 5 failed attempts
- Automatic reset of attempts after cooldown

### 4. Token Security
- JWT tokens with 24-hour expiration
- Token payload includes user information and permissions
- Secure token generation with proper error handling

### 5. Session Management
- Session expiry tracking
- Last login time updates
- Login attempt tracking and reset

## Response Format

All responses use `utils.BaseResponseMaker` for consistent formatting:

```json
{
  "status": 200,
  "message": "Success message",
  "data": {
    // Response data
  }
}
```

## Testing

Comprehensive test coverage includes:

- Successful login flow
- User not found scenarios
- Account status validation (deleted, blocked)
- Device linking verification
- PIN validation and rate limiting
- Input validation scenarios
- Error handling and status codes
- Token generation and response format

## Usage Examples

### cURL Examples

**1. Successful Login:**
```bash
curl -X POST http://localhost:8080/api/v1/cbesuperapp/user/login \
  -H "Content-Type: application/json" \
  -H "device-uuid: device-123" \
  -H "platform: android" \
  -H "app-version: 1.0.0" \
  -d '{
    "phone": "1234567890",
    "device_uuid": "device-123",
    "pin": "123456"
  }'
```

**2. Invalid PIN:**
```bash
curl -X POST http://localhost:8080/api/v1/cbesuperapp/user/login \
  -H "Content-Type: application/json" \
  -H "device-uuid: device-123" \
  -H "platform: android" \
  -H "app-version: 1.0.0" \
  -d '{
    "phone": "1234567890",
    "device_uuid": "device-123",
    "pin": "654321"
  }'
```

**3. Device Mismatch:**
```bash
curl -X POST http://localhost:8080/api/v1/cbesuperapp/user/login \
  -H "Content-Type: application/json" \
  -H "device-uuid: different-device" \
  -H "platform: android" \
  -H "app-version: 1.0.0" \
  -d '{
    "phone": "1234567890",
    "device_uuid": "device-123",
    "pin": "123456"
  }'
```

## Implementation Notes

### 1. Database Operations
- User lookup by phone number
- Login attempt tracking and updates
- Last login time updates
- Device verification

### 2. Token Generation
- JWT token with 24-hour expiration
- Token payload includes user ID, name, phone, realm, role, permissions
- Secure token generation with proper error handling

### 3. Rate Limiting Logic
- Tracks login attempts in user record
- Implements 15-minute cooldown after 5 failed attempts
- Automatic reset of attempts after successful login or cooldown

### 4. Device Security
- Verifies device UUID matches user's linked device
- Validates device UUID consistency between header and request body
- Prevents login from unauthorized devices

### 5. Error Recovery
- Graceful handling of database errors
- Proper error logging for debugging
- User-friendly error messages

## Future Enhancements

### 1. Multi-factor Authentication
- SMS OTP verification
- Email verification
- Biometric authentication

### 2. Advanced Security
- IP address tracking
- Geographic location validation
- Device fingerprinting
- Suspicious activity detection

### 3. Session Management
- Multiple device sessions
- Session revocation
- Real-time session monitoring

### 4. Analytics and Monitoring
- Login success rate tracking
- Failed login pattern analysis
- Security event logging
- Performance monitoring

### 5. Compliance Features
- Audit trail for login events
- GDPR compliance features
- Data retention policies
- Privacy controls

## Security Best Practices

### 1. Password/PIN Security
- PIN history tracking
- PIN complexity requirements
- Secure PIN storage (encrypted)
- PIN change notifications

### 2. Session Security
- Secure token storage
- Token rotation
- Session timeout
- Concurrent session limits

### 3. Device Security
- Device registration process
- Device verification
- Device deauthorization
- Device activity monitoring

### 4. Network Security
- HTTPS enforcement
- Request rate limiting
- DDoS protection
- API security headers

### 5. Data Protection
- Data encryption at rest
- Data encryption in transit
- Secure logging practices
- Data minimization principles 