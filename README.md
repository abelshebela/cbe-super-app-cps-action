<<<<<<< HEAD
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
=======
# CBE Super App Development Rules

This repository provides a clean, maintainable foundation for building robust micro-services using Go. It embraces best practices for development, testing, versioning, and team collaboration.

---
## 📦 Project Overview

The project follows a modular architecture using clear separation of concerns for maintainability and scalability. It uses MongoDB, Cassandra Database & Oracle Database ( as the primary database and `golang-migrate` for schema versioning ) .

---
## 🔄 Branching Strategy

### **Development Branch (`dev`)**
- The `dev` branch is the main working branch for active development.
- Feature branches are merged into `dev` via pull requests after review and testing.

### **UAT Branch (`staging`)**
- Used for User Acceptance Testing (UAT).
- Pull requests to `staging` must originate from `dev`.

### **Production Branch (`main`)**
- Production-ready, stable code lives here.
- Pull requests to `main` must originate from `staging`.

---
## 🔀 Pull Request Guidelines

To ensure consistency and maintain code quality, adhere to the following rules:

1.  **Branch Naming**
- Use clear, descriptive names:
-  `feature/auth-login`, `fix/db-connection`, `chore/lint-cleanup`

2.  **PR Targets**
- Feature PRs → `dev`
- QA/UAT PRs → `staging`
- Release PRs → `main`

1.  **PR Rules**
- ✅ Must include test cases for new features or bug fixes
- ✅ Must pass all lint and test checks
- ✅ Must be peer-reviewed before merging
- ✅ Must avoid introducing breaking changes without documentation
- ❌ Do not commit directly to `main`, `staging`, or `dev`

---
## 🎯 Code Style Guide

| Element                | Convention   | Example                            |
| ---------------------- | ------------ | ---------------------------------- |
| Public Functions/Vars  | `Go Rule`    | `GetUserByID`, `CreateTransaction` |
| Private Functions/Vars | `Go Rule`    | `getUserByID`, `createTransaction` |
| Public Structs/Types   | `PascalCase` | `AuthPayload`, `TransactionInput`  |
| Private Structs/Types  | `PascalCase` | `AuthPayload`, `TransactionInput`  |
| JSON Fields            | `snake_case` | `user_id`, `txn_id`                |
| BSON Fields            | `snake_case` | `user_id`, `txn_id`                |
| Files/Folders          | `snake_case` | `user_handler.go`, `db_conn.go`    |
  
---

## ✅ Development Checklist
- [ ] Feature isolated in its own branch
- [ ] Test cases added
- [ ] Input Validation is Included
- [ ] Code reviewed
- [ ] PR targets correct branch
- [ ] CI checks passing
- [ ] Migration included (if schema is updated)

---
## ⚠️ NOTE 
1. We recommend you to use the shared repository
2. Follow the code structure
3. After you push you code to Git lab, don't forget to send a PR Request for reviews
  
> For any questions or contributions, please contact the project leads.
>>>>>>> origin
