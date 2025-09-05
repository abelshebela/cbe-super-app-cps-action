# gRPC Bank Service Setup

This document explains how to use the gRPC bank service that has been set up in your CBE Super App CPS Action module.

## Overview

The gRPC service exposes your existing bank API functionality through gRPC, allowing other services to communicate with your bank module efficiently using binary protocol buffers.

## Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   gRPC Client   │───▶│   gRPC Server    │───▶│   Bank Module   │
│                 │    │   (Port 9090)    │    │   (HTTP Layer)  │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

## Services Available

### BankService

The gRPC service provides the following methods:

1. **GetBank** - Retrieve a single bank by ID
2. **GetAllBanks** - Get all banks with pagination and filtering
3. **CreateBank** - Create a new bank
4. **UpdateBank** - Update an existing bank
5. **DeleteBank** - Delete a bank
6. **EnableDisableBank** - Enable or disable a bank
7. **UpdateBankLogo** - Update a bank's logo

## Server Configuration

### Ports
- **HTTP Server**: Port 8080
- **gRPC Server**: Port 9090

### Starting the Server

The server automatically starts both HTTP and gRPC services when you run:

```bash
go run cmd/main.go
```

You should see logs indicating both servers are starting:
```
🚀 HTTP Server starting on port 8080
🚀 gRPC Server starting on port 9090
```

## Client Usage

### Go Client Example

```go
package main

import (
    "context"
    "log"
    "time"
    
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

func main() {
    // Connect to gRPC server
    conn, err := grpc.Dial("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer conn.Close()

    // Create client and make calls
    // (See examples/grpc_client/main.go for full implementation)
}
```

### Testing the Client

Run the example client:

```bash
cd examples/grpc_client
go run main.go
```

## Protocol Buffer Definitions

The service is defined in `proto/bank/bank.proto` with the following key messages:

### Bank Message
```protobuf
message Bank {
  string id = 1;
  string name = 2;
  string logo = 3;
  string code = 4;
  string bic = 5;
  bool enabled = 6;
  google.protobuf.Timestamp created_at = 7;
  google.protobuf.Timestamp last_modified_at = 8;
}
```

### Request/Response Examples

#### GetAllBanks
```protobuf
message GetAllBanksRequest {
  int32 page = 1;
  int32 per_page = 2;
  string search = 3;
  string sort_by = 4;
  string sort_order = 5;
}

message GetAllBanksResponse {
  repeated Bank banks = 1;
  int32 total = 2;
  int32 page = 3;
  int32 per_page = 4;
  int32 total_pages = 5;
  string message = 6;
}
```

## Implementation Details

### Server Structure
- **`internal/grpc/server/bank_server.go`** - Main gRPC server implementation
- **`internal/grpc/grpc_server.go`** - gRPC server wrapper and startup logic
- **`internal/grpc/service/bank_service.go`** - Service interface definitions

### Key Features
- **Logging Interceptor** - All gRPC calls are logged
- **Error Handling** - Proper gRPC status codes and error messages
- **Graceful Shutdown** - Both HTTP and gRPC servers shutdown gracefully
- **Reflection** - gRPC reflection is enabled for debugging

### Data Flow
1. gRPC client makes a call to the gRPC server (port 9090)
2. gRPC server receives the request and converts it to internal format
3. Server calls the existing bank application service
4. Application service processes the request using existing business logic
5. Response is converted back to protobuf format and returned to client

## Security Considerations

Currently, the gRPC server runs without authentication for simplicity. For production use, consider:

1. **TLS/SSL** - Enable secure connections
2. **Authentication** - Implement token-based or certificate-based auth
3. **Authorization** - Add role-based access control
4. **Rate Limiting** - Prevent abuse

## Troubleshooting

### Common Issues

1. **Port Already in Use**
   ```
   Error: listen tcp :9090: bind: address already in use
   ```
   Solution: Change the gRPC port in `initiator/initiator.go`

2. **Connection Refused**
   ```
   Error: connection refused
   ```
   Solution: Ensure the server is running and accessible

3. **Import Errors**
   ```
   Error: cannot find package
   ```
   Solution: Run `go mod tidy` to update dependencies

### Logs
Check the application logs for gRPC-specific messages:
- `gRPC call: /BankService/GetBank` - Incoming requests
- `gRPC call completed successfully` - Successful responses
- `gRPC call failed:` - Error responses

## Next Steps

1. **Generate Proper Protobuf Code**: Install `protoc` and run the generation script
2. **Add Authentication**: Implement secure authentication mechanisms
3. **Add More Services**: Extend to other modules (account, wallet, etc.)
4. **Performance Testing**: Test with high load scenarios
5. **Monitoring**: Add metrics and monitoring for gRPC endpoints

## Files Created/Modified

### New Files
- `proto/bank/bank.proto` - Protocol buffer definitions
- `internal/grpc/server/bank_server.go` - gRPC server implementation
- `internal/grpc/grpc_server.go` - Server wrapper
- `internal/grpc/service/bank_service.go` - Service interface
- `examples/grpc_client/main.go` - Example client
- `scripts/generate_proto.bat` - Proto generation script

### Modified Files
- `go.mod` - Added gRPC dependencies
- `initiator/initiator.go` - Added gRPC server startup

The gRPC service is now ready to serve your bank API to other clients efficiently!
