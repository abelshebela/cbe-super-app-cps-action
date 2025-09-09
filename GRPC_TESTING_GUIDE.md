# gRPC Bank Service Testing Guide

## Overview
The gRPC bank service is now properly implemented with service registration. Here are multiple ways to test and access the service.

## Service Status
✅ **Fixed Issues:**
- gRPC server registration (using `pb.RegisterBankServiceServer`)
- Proper protobuf message types implementation
- Service interface compliance with generated protobuf code
- Server startup integration with HTTP server

## Testing Methods

### 1. Using Postman (Recommended for GUI)

**Prerequisites:**
- Postman version 8.0+ with gRPC support
- Server running on `localhost:9090`

**Steps:**
1. Start the application: `.\cbe-super-app-cps-action.exe`
2. Open Postman and create a new gRPC request
3. Enter server URL: `localhost:9090`
4. Import proto file: `proto/bank/bank.proto`
5. Select service method (e.g., `BankService.GetAllBanks`)
6. Send request with appropriate JSON payload

**Sample Requests:**

**GetAllBanks:**
```json
{
  "page": 1,
  "per_page": 10,
  "search": ""
}
```

**GetBank:**
```json
{
  "id": "your-bank-id-here"
}
```

**CreateBank:**
```json
{
  "name": "Test Bank",
  "code": "TB001",
  "bic": "TESTBIC001"
}
```

### 2. Using grpcurl (Command Line)

**Installation:**
```bash
# Windows (using Chocolatey)
choco install grpcurl

# Or download from: https://github.com/fullstorydev/grpcurl/releases
```

**Test Commands:**

**List available services:**
```bash
grpcurl -plaintext localhost:9090 list
```

**Get service methods:**
```bash
grpcurl -plaintext localhost:9090 list bank.BankService
```

**Call GetAllBanks:**
```bash
grpcurl -plaintext -d "{\"page\": 1, \"per_page\": 10}" localhost:9090 bank.BankService/GetAllBanks
```

**Call GetBank:**
```bash
grpcurl -plaintext -d "{\"id\": \"your-bank-id\"}" localhost:9090 bank.BankService/GetBank
```

**Call CreateBank:**
```bash
grpcurl -plaintext -d "{\"name\": \"Test Bank\", \"code\": \"TB001\", \"bic\": \"TESTBIC001\"}" localhost:9090 bank.BankService/CreateBank
```

### 3. Using Go Client Examples

**Simple Connection Test:**
```bash
go run examples/connection_test.go
```

**Full Client Demo:**
```bash
go run examples/grpc_client/main.go
```

**Bank Operations Demo:**
```bash
go run examples/bank_client_demo.go
```

### 4. Using BloomRPC (Alternative GUI)

**Installation:**
- Download from: https://github.com/bloomrpc/bloomrpc/releases
- Import `proto/bank/bank.proto`
- Connect to `localhost:9090`

### 5. Using Evans (Interactive CLI)

**Installation:**
```bash
# Windows
scoop install evans
```

**Usage:**
```bash
evans --host localhost --port 9090 --reflection repl
```

## Troubleshooting

### Common Issues and Solutions

**1. Connection Refused**
- Ensure the application is running: `.\cbe-super-app-cps-action.exe`
- Check if port 9090 is available: `netstat -an | findstr :9090`
- Verify both HTTP (8080) and gRPC (9090) servers are starting

**2. Service Not Found**
- Confirm service registration is working
- Use reflection to list services: `grpcurl -plaintext localhost:9090 list`

**3. Method Not Implemented**
- Ensure all protobuf methods are implemented in `bank_server.go`
- Check that the server embeds `pb.UnimplementedBankServiceServer`

**4. Build Errors**
- Run `go mod tidy` to ensure dependencies are correct
- Regenerate protobuf code: `scripts\generate_proto.bat`

### Verification Steps

**1. Check Server Startup Logs:**
Look for these log messages when starting the application:
```
Starting gRPC server on :9090
gRPC server started successfully
```

**2. Test Basic Connectivity:**
```bash
# Test if port is open
telnet localhost 9090

# Or use PowerShell
Test-NetConnection -ComputerName localhost -Port 9090
```

**3. Verify Service Registration:**
```bash
grpcurl -plaintext localhost:9090 list
# Should show: bank.BankService
```

## Service Methods Available

| Method | Description | Request Type | Response Type |
|--------|-------------|--------------|---------------|
| `GetBank` | Get single bank by ID | `GetBankRequest` | `GetBankResponse` |
| `GetAllBanks` | Get all banks with pagination | `GetAllBanksRequest` | `GetAllBanksResponse` |
| `CreateBank` | Create new bank | `CreateBankRequest` | `CreateBankResponse` |
| `UpdateBank` | Update existing bank | `UpdateBankRequest` | `UpdateBankResponse` |
| `DeleteBank` | Delete bank | `DeleteBankRequest` | `DeleteBankResponse` |
| `EnableDisableBank` | Enable/disable bank | `EnableDisableBankRequest` | `EnableDisableBankResponse` |
| `UpdateBankLogo` | Update bank logo | `UpdateBankLogoRequest` | `UpdateBankLogoResponse` |

## Security Notes

**Current Implementation:**
- No authentication/authorization (development only)
- Plain text connection (no TLS)

**Production Recommendations:**
- Implement JWT/OAuth authentication
- Enable TLS encryption
- Add rate limiting
- Implement proper error handling
- Add request validation middleware

## Next Steps

1. **Test the service** using any of the methods above
2. **Verify all methods work** with sample data
3. **Add authentication** for production use
4. **Implement TLS** for secure communication
5. **Add monitoring** and metrics collection

The gRPC service is now properly configured and should be accessible via Postman and other gRPC clients on `localhost:9090`.
