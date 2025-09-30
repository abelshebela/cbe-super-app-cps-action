# CPS Action API Documentation

## Interactive API Documentation

The CPS Action API is documented using the OpenAPI 3.0 specification and can be explored interactively using Swagger UI.

### Accessing Swagger UI

1. **Using the Development Server**:
   ```bash
   # Start the Swagger UI server
   python scripts/serve_swagger.py
   
   # Open in browser
   http://localhost:8000
   ```

2. **With Docker**:
   ```bash
   docker run -p 8000:8000 -v $(pwd)/docs/swagger:/usr/share/nginx/html/swagger -d nginx
   ```
   Then navigate to: `http://localhost:8000/swagger/`

### Features

- **Interactive API Exploration**: Try out API endpoints directly from your browser
- **Request/Response Visualization**: View detailed request and response schemas
- **Authentication**: Test authenticated endpoints using JWT tokens
- **Model Documentation**: View detailed documentation of all data models

### Authentication in Swagger UI

1. Click the "Authorize" button in the top-right corner
2. Enter your JWT token in the format: `Bearer <your-jwt-token>`
3. Click "Authorize" to authenticate

## Project Structure

```
cbe-supper-app-cps-action/
├── cmd/                  # Main application entry point
├── config/              # Configuration management
├── docs/                # API documentation
│   └── swagger/         # Swagger/OpenAPI specifications
├── internal/            # Core application code
│   ├── constants/       # DTOs, enums, and constants
│   ├── glue/            # Routing and API layer
│   ├── handlers/        # HTTP handlers
│   ├── service/         # Business logic layer
│   └── storage/         # Data access layer
├── initiator/           # Application initialization and DI
├── pkg/                 # Shared utility packages
└── platform/            # Infrastructure components
```

## API Documentation

### Authentication
All API endpoints require authentication using a JWT token. Include the token in the `Authorization` header:

```
Authorization: Bearer <your-jwt-token>
```

### Base URL
```
http://localhost:8080/api/v1/cbesuperapp/cps_action
```

### Response Format
All API responses follow a standard format:

```json
{
  "ok": true,
  "status": 200,
  "timestamp": "2023-01-01T00:00:00Z",
  "message": "Success message",
  "data": {}
}
```

### Error Response

```json
{
  "ok": false,
  "status": 400,
  "timestamp": "2023-01-01T00:00:00Z",
  "message": "Error message",
  "error": {
    "code": "INVALID_REQUEST",
    "message": "Detailed error message",
    "status_code": 400,
    "type": "validation_error",
    "field_errors": [
      {
        "field": "field_name",
        "message": "Error message",
        "value": "invalid_value",
        "constraint": "required"
      }
    ]
  }
}
```

## Available Endpoints

### Account Block
- `GET /account-block` - List all account blocks
- `POST /account-block` - Create a new account block
- `GET /account-block/{id}` - Get account block by ID
- `PUT /account-block/{id}` - Update account block
- `DELETE /account-block/{id}` - Delete account block

### Bank
- `GET /bank` - List all banks
- `POST /bank` - Create a new bank
- `GET /bank/{id}` - Get bank by ID
- `PUT /bank/{id}` - Update bank
- `DELETE /bank/{id}` - Delete bank
- `PATCH /bank/{id}/status` - Enable/disable bank

### Donation
- `GET /donation` - List all donations
- `POST /donation` - Create a new donation
- `GET /donation/{id}` - Get donation by ID
- `PUT /donation/{id}` - Update donation
- `DELETE /donation/{id}` - Delete donation
- `PATCH /donation/{id}/status` - Enable/disable donation

### User Management
- `GET /user` - List all users
- `POST /user` - Create a new user
- `GET /user/{id}` - Get user by ID
- `PUT /user/{id}` - Update user
- `DELETE /user/{id}` - Delete user
- `PATCH /user/{id}/status` - Enable/disable user

### Wallet
- `GET /wallet` - List all wallets
- `POST /wallet` - Create a new wallet
- `GET /wallet/{id}` - Get wallet by ID
- `PUT /wallet/{id}` - Update wallet
- `DELETE /wallet/{id}` - Delete wallet

## Generating API Documentation

1. Install Python dependencies:
   ```bash
   pip install -r requirements.txt
   ```

2. Generate Swagger documentation:
   ```bash
   python scripts/generate_swagger.py
   ```

3. View the documentation:
   - Open `docs/swagger/swagger.json` in Swagger UI or any OpenAPI viewer
   - Or use the online Swagger Editor: https://editor.swagger.io/

## Development

### Prerequisites
- Go 1.19+
- MongoDB
- Redis
- MinIO (for file storage)

### Running the Application

1. Set up environment variables (copy from `.env.example` to `.env` and update values)
2. Start the application:
   ```bash
   go run cmd/main.go
   ```

### Testing

Run tests:
```bash
go test ./...
```

## License

This project is proprietary software. All rights reserved.
