# Device Lookup OTP Functionality

## Overview

The device lookup API now includes OTP generation and sending functionality with specific conditions for development and testing environments.

## API Endpoint

- **Method**: `GET`
- **Path**: `/api/v1/cbesuperapp/user/device-lookup`
- **Authentication**: Not required

## Request Headers

- `device-uuid`: Device unique identifier (required)
- `platform`: Platform type - android, ios, web (required)
- `app-version`: Application version (optional)

## Response Structure

```json
{
  "device_uuid": "string",
  "platform": "string",
  "app_version": "string",
  "is_latest": boolean,
  "user_found": boolean,
  "token": "string",
  "token_type": "device_lookup",
  "token_expiry": "timestamp",
  "next_step": "string",
  "otp_code": "string" // Only included under specific conditions
}
```

## OTP Generation Conditions

The OTP code is **only** included in the response when **ALL** of the following conditions are met:

1. **User Found**: A user with the provided device UUID exists in the system
2. **User Not Verified**: The user's `is_verified` field is `false`
3. **Environment**: The application is running in `dev` or `uat` environment

## OTP Behavior by Scenario

### Scenario 1: User Found + Not Verified + Dev/UAT Environment
- ✅ **OTP Generated**: Yes
- ✅ **OTP Included in Response**: Yes
- ✅ **SMS Sent**: Yes
- ✅ **OTP Stored in Database**: Yes

**Response includes:**
```json
{
  "user_found": true,
  "next_step": "login",
  "otp_code": "123456"
}
```

### Scenario 2: User Found + Verified + Dev/UAT Environment
- ❌ **OTP Generated**: No
- ❌ **OTP Included in Response**: No
- ❌ **SMS Sent**: No
- ❌ **OTP Stored in Database**: No

**Response includes:**
```json
{
  "user_found": true,
  "next_step": "login",
  "otp_code": null
}
```

### Scenario 3: User Found + Not Verified + Production Environment
- ❌ **OTP Generated**: No
- ❌ **OTP Included in Response**: No
- ❌ **SMS Sent**: No
- ❌ **OTP Stored in Database**: No

**Response includes:**
```json
{
  "user_found": true,
  "next_step": "login",
  "otp_code": null
}
```

### Scenario 4: No User Found (Any Environment)
- ❌ **OTP Generated**: No
- ❌ **OTP Included in Response**: No
- ❌ **SMS Sent**: No
- ❌ **OTP Stored in Database**: No

**Response includes:**
```json
{
  "user_found": false,
  "next_step": "register",
  "otp_code": null
}
```

## OTP Details

When OTP is generated:

- **Length**: 6 digits
- **Format**: Numeric only
- **Purpose**: `pin_set` (using enum `OTPForPINSet`)
- **Expiration**: Configurable via `OTP_WAITING_TIME` (default: 10 minutes)
- **Storage**: Encrypted in database
- **Delivery**: SMS to user's phone number

## SMS Message Format

```
Your device lookup OTP is: 123456. Valid for 10 minutes.
```

## Security Features

1. **Environment Restriction**: OTP only exposed in development/testing environments
2. **Verification Check**: Only unverified users receive OTP
3. **Encryption**: OTP is encrypted before database storage
4. **Expiration**: OTP has configurable expiration time
5. **Async SMS**: SMS sending is non-blocking

## Configuration

The following environment variables control OTP behavior:

- `GO_ENV`: Environment type (dev, uat, production)
- `OTP_WAITING_TIME`: OTP expiration time in minutes
- `KEY`: Encryption key for OTP
- `IV`: Initialization vector for encryption

## Testing

Use the test file `tests/domain/users/device_lookup_test.go` to verify all scenarios:

```bash
go test -v ./tests/domain/users/ -run TestUserService_DeviceLookup_WithOTP
```

## Implementation Details

### Domain Service Changes

The `DeviceLookup` method in `internal/domain/users/service.go` now includes:

```go
// Generate OTP and include in response for dev and uat environments when user is found and not verified
if userFound && !user.IsVerified && (s.cfg.GoEnv == "dev" || s.cfg.GoEnv == "uat") {
    // OTP generation and sending logic
    otpRecord := OTPRecord{
        UserID:     user.ID.Hex(),
        OTP:        encOtpCode,
        OTPFor:     string(enums.OTPForPINSet), // Using enum for "pin_set"
        UserRealm:  "member",
        ExpiresAt:  time.Now().Add(expirationTime),
        CreatedAt:  time.Now(),
        DeviceUUID: &deviceUUID,
    }
}
```

### DTO Changes

The `DeviceLookupResponse` DTO in `internal/application/dto/user_dto.go` includes:

```go
type DeviceLookupResponse struct {
    // ... existing fields
    OTPCode string `json:"otp_code,omitempty"` // Only included in dev and uat environments
}
```

## Error Handling

- **OTP_ENCRYPTION_FAILED**: Failed to encrypt OTP for storage
- **OTP_CREATION_FAILED**: Failed to create OTP record in database
- **TOKEN_GENERATION_FAILED**: Failed to generate device lookup token

## Logging

The service logs OTP generation events with details:

```
OTP generated and sent for device lookup - Device: device-uuid, User: phone-number, OTP: 123456, IsVerified: false
``` 