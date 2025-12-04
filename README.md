# CBE Super App - CPS Action Module

This is the CPS (Customer and Partner Services) Action module for the CBE Super App. It provides a comprehensive set of APIs for managing various aspects of customer and partner services.

## Features

- **Bank Management**: Create, read, update, and delete bank information
- **User Management**: Manage user accounts and permissions
- **Documentation**: Comprehensive API documentation with Swagger UI
- **Authentication**: JWT-based authentication for secure API access
- **RESTful API**: Well-structured REST endpoints following best practices

## Prerequisites

- Go 1.19 or higher
- MongoDB
- Redis (for caching and sessions)
- MinIO (for file storage)

## Getting Started

### Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd cbe-supper-app-cps-action
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up environment variables:
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

### Running the Application

Start the application:
```bash
go run cmd/main.go
```

The application will be available at `http://localhost:8080`

## API Documentation

### Viewing API Documentation

1. Start the Swagger UI server:
   ```bash
   python scripts/serve_swagger.py
   ```

2. Open your browser and navigate to:
   ```
   http://localhost:8000
   ```

The Swagger UI provides interactive documentation where you can:
- View all available API endpoints
- Try out API calls directly from the browser
- View request/response schemas
- Authenticate using your JWT token

### Authentication

1. Obtain a JWT token by authenticating with the authentication service
2. Click the "Authorize" button in the Swagger UI
3. Enter your JWT token in the format: `Bearer <your-jwt-token>`

## Project Structure

```
.
├── cmd/                 # Main application entry point
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
└── scripts/             # Utility scripts
```

## Development

### Code Generation

To generate gRPC code:
```bash
./scripts/generate_proto.bat
```

### Testing

Run tests:
```bash
go test ./...
```

### Linting

```bash
golangci-lint run
```

## Deployment

### Building the Application

```bash
go build -o bin/server cmd/main.go
```

### Docker

Build the Docker image:
```bash
docker build -t cbe-super-app-cps-action .
```

Run the container:
```bash
docker run -p 8080:8080 cbe-super-app-cps-action
```

## License

This project is proprietary software. All rights reserved.
