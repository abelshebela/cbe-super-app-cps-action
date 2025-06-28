# PIN Reset Process Documentation

## Overview

The PIN reset functionality allows users to securely reset their PIN when they forget it. The process follows a two-step approach:
1. **Forget PIN Send OTP** - User requests OTP for PIN reset
2. **Reset PIN** - User verifies OTP and sets new PIN with restricted access

## Architecture

### Flow Diagram
```
User Forgets PIN
       ↓
Request OTP (ForgetPinSendOtp)
       ↓
Validate User & Device
       ↓
Generate & Send OTP
       ↓
Store PIN Reset Session
       ↓
User Receives OTP
       ↓
Verify OTP & Set New PIN (ResetPin)
       ↓
Update User PIN
       ↓
Grant Access with Restrictions
```

### Components

1. **HTTP Adapter Layer** (`internal/adapter/inbound/http/users/`)
   - `ForgetPinSendOtp` - Handles OTP request
   - `ResetPin` - Handles PIN reset with OTP verification

2. **Application Layer** (`internal/application/users/`)
   - Input validation and business logic coordination
   - Error handling and response formatting

3. **Domain Service Layer** (`internal/domain/users/`)
   - Core business logic for PIN reset
   - Security validation and PIN history management

4. **Repository Layer** (`internal/adapter/outbound/persistence/`)
   - MongoDB operations for PIN reset sessions
   - User data persistence

## API Endpoints

### 1. Forget PIN Send OTP
```
POST /api/v1/cbesuperapp/user/forget-pin/send-otp
```

**Request Body:**
```json
{
  "phone": "1234567890",
  "device_uuid": "device-123"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "OTP sent for PIN reset",
  "data": {
    "phone_number": "1234567890",
    "device_uuid": "device-123",
    "otp_sent": true,
    "otp_expiry_minutes": 10,
    "reset_session_id": "uuid-session-id",
    "next_step": "verify_otp"
  }
}
```

### 2. Reset PIN
```
POST /api/v1/cbesuperapp/user/forget-pin/reset
```

**Request Body:**
```json
{
  "reset_session_id": "uuid-session-id",
  "phone": "1234567890",
  "device_uuid": "device-123",
  "otp": "123456",
  "new_pin": "654321"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "PIN reset completed successfully",
  "data": {
    "user_id": "user-id",
    "user_code": "12345678",
    "full_name": "John Doe",
    "phone_number": "1234567890",
    "pin_reset": true,
    "reset_time": "2024-01-01T12:00:00Z",
    "access_restricted": true,
    "restrictions": ["transfer", "balance_check"],
    "next_step": "login_with_new_pin"
  }
}
```

## Data Models

### PinResetSession
```go
type PinResetSession struct {
    ID              string    `bson:"_id" json:"id"`
    UserID          string    `bson:"user_id" json:"user_id"`
    PhoneNumber     string    `bson:"phone_number" json:"phone_number"`
    DeviceUUID      string    `bson:"device_uuid" json:"device_uuid"`
    OTP             string    `bson:"otp" json:"otp"`
    OTPFor          string    `bson:"otp_for" json:"otp_for"`
    Status          string    `bson:"status" json:"status"` // pending, verified, completed, expired
    ExpiresAt       time.Time `bson:"expires_at" json:"expires_at"`
    CreatedAt       time.Time `bson:"created_at" json:"created_at"`
    Attempts        int       `bson:"attempts" json:"attempts"`
    MaxAttempts     int       `bson:"max_attempts" json:"max_attempts"`
    VerifiedAt      time.Time `bson:"verified_at" json:"verified_at"`
    CompletedAt     time.Time `bson:"completed_at" json:"completed_at"`
    AccessRestricted bool     `bson:"access_restricted" json:"access_restricted"`
    Restrictions    []string  `bson:"restrictions" json:"restrictions"`
}
```

## Security Features

### 1. Input Validation
- Phone number format validation
- Device UUID validation
- OTP format validation (6 digits)
- PIN format validation (6 digits, numeric only)

### 2. User Verification
- User existence check
- Account status validation (not deleted, not blocked)
- Device linking verification
- Phone number ownership verification

### 3. Session Management
- Unique session ID generation
- Session expiration (10 minutes)
- Attempt limiting (max 3 attempts)
- Session status tracking

### 4. PIN Security
- PIN history validation (prevents reuse of last 4 PINs)
- Encrypted OTP storage
- Secure PIN update process

### 5. Access Restrictions
After PIN reset, users get restricted access:
- Transfer operations blocked
- Balance check operations blocked
- Other sensitive operations may be restricted

## Error Handling

### Common Error Codes
- `PIN_RESET_USER_NOT_FOUND` - User not found
- `ACCOUNT_DELETED` - Account has been deleted
- `ACCOUNT_BLOCKED` - Account is blocked
- `PIN_RESET_DEVICE_MISMATCH` - Device not linked to user
- `PIN_RESET_ALREADY_IN_PROGRESS` - Reset already in progress
- `PIN_RESET_SESSION_NOT_FOUND` - Session not found
- `PIN_RESET_SESSION_EXPIRED` - Session expired
- `PIN_RESET_TOO_MANY_ATTEMPTS` - Too many attempts
- `PIN_RESET_OTP_INVALID` - Invalid OTP
- `PIN_RESET_FAILED` - General reset failure

### HTTP Status Codes
- `200` - Success
- `400` - Bad Request (validation errors)
- `401` - Unauthorized
- `403` - Forbidden (account blocked)
- `404` - Not Found (user/session not found)
- `409` - Conflict (reset in progress)
- `410` - Gone (account deleted, session expired)
- `429` - Too Many Requests (rate limited)
- `500` - Internal Server Error

## Database Schema

### pin_reset_sessions Collection
```javascript
{
  "_id": "uuid-session-id",
  "user_id": "user-object-id",
  "phone_number": "1234567890",
  "device_uuid": "device-123",
  "otp": "encrypted-otp-string",
  "otp_for": "pin_reset",
  "status": "pending",
  "expires_at": ISODate("2024-01-01T12:10:00Z"),
  "created_at": ISODate("2024-01-01T12:00:00Z"),
  "attempts": 0,
  "max_attempts": 3,
  "verified_at": null,
  "completed_at": null,
  "access_restricted": true,
  "restrictions": ["transfer", "balance_check"]
}
```

## Testing

### Unit Tests
Comprehensive test coverage including:
- Success scenarios
- User validation errors
- Session management errors
- Input validation errors
- Security validation errors

### Test Cases
1. **ForgetPinSendOtp Tests**
   - Successful OTP generation
   - User not found
   - Account deleted/blocked
   - Device mismatch
   - Invalid input validation

2. **ResetPin Tests**
   - Successful PIN reset
   - Session not found/expired
   - Too many attempts
   - Invalid OTP/PIN format
   - User validation errors

## Usage Examples

### 1. Request OTP for PIN Reset
```bash
curl -X POST http://localhost:8080/api/v1/cbesuperapp/user/forget-pin/send-otp \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "1234567890",
    "device_uuid": "device-123"
  }'
```

### 2. Reset PIN with OTP
```bash
curl -X POST http://localhost:8080/api/v1/cbesuperapp/user/forget-pin/reset \
  -H "Content-Type: application/json" \
  -d '{
    "reset_session_id": "uuid-session-id",
    "phone": "1234567890",
    "device_uuid": "device-123",
    "otp": "123456",
    "new_pin": "654321"
  }'
```

## Configuration

### Environment Variables
- `JWT_SECRET_KEY` - For OTP encryption
- `SMS_API_ENDPOINT` - For OTP delivery
- `MONGODB_URI` - Database connection

### Timeouts
- OTP expiration: 10 minutes
- Session cleanup: Automatic after expiration
- Rate limiting: Configurable per endpoint

## Monitoring & Logging

### Key Metrics
- PIN reset success rate
- OTP delivery success rate
- Session expiration rate
- Failed attempt patterns

### Log Events
- PIN reset initiated
- OTP sent successfully
- PIN reset completed
- Security violations
- Error conditions

## Future Enhancements

### Planned Features
1. **Multi-factor Authentication**
   - Email verification in addition to SMS
   - Biometric verification options

2. **Enhanced Security**
   - Device fingerprinting
   - Behavioral analysis
   - Risk-based restrictions

3. **User Experience**
   - Progressive PIN reset flow
   - Self-service account recovery
   - PIN strength requirements

4. **Administrative Features**
   - PIN reset audit logs
   - Admin-initiated resets
   - Bulk PIN reset operations

## Integration Points

### External Services
- SMS Gateway (OTP delivery)
- Email Service (optional notifications)
- Device Management Service
- Audit Logging Service

### Internal Services
- User Management Service
- Authentication Service
- Notification Service
- Audit Service

## Compliance & Security

### Data Protection
- OTP encryption at rest
- Secure transmission protocols
- Data retention policies
- Privacy compliance

### Audit Requirements
- Complete audit trail
- Session tracking
- Access logging
- Security event monitoring

## Troubleshooting

### Common Issues
1. **OTP Not Received**
   - Check phone number format
   - Verify SMS service status
   - Check user account status

2. **Session Expired**
   - Request new OTP
   - Check system time synchronization
   - Verify session cleanup processes

3. **Device Mismatch**
   - Verify device registration
   - Check device linking status
   - Contact support if needed

### Debug Information
- Session ID tracking
- Request/response logging
- Error code mapping
- Performance metrics 