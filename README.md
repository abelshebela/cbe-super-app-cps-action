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
