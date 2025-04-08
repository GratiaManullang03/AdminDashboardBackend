# Admin Dashboard Backend - Golang

A RESTful API backend service built with Golang and Gin framework for an administrative dashboard application. This service provides authentication, user management, role management, division management, position management, and dashboard statistics.

## Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [API Endpoints](#api-endpoints)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Environment Variables](#environment-variables)
- [Database Schema](#database-schema)
- [Authentication](#authentication)
- [Deployment](#deployment)
- [API Documentation](#api-documentation)

## Features

- **User Authentication**: JWT-based authentication system
- **User Management**: CRUD operations for users with role assignments
- **Role Management**: Define and manage user roles with permission levels
- **Division Management**: Organize users by divisions
- **Position Management**: Define and manage different positions within the organization
- **Dashboard Statistics**: Get organizational statistics and data visualizations
- **Middleware**: Authentication, CORS, Logging, and Error handling

## Tech Stack

- **Language**: Go (Golang)
- **Web Framework**: Gin
- **ORM**: GORM
- **Database**: PostgreSQL
- **Authentication**: JWT (JSON Web Tokens)
- **Deployment**: Railway

## API Endpoints

### Authentication

| Endpoint | Method | Description | Authentication |
|----------|--------|-------------|----------------|
| `/api/auth/login` | POST | Login user | No |
| `/api/auth/profile` | GET | Get current user profile | Yes |

### User Management

| Endpoint | Method | Description | Authentication |
|----------|--------|-------------|----------------|
| `/api/users` | GET | List all users (with pagination) | Yes |
| `/api/users/{id}` | GET | Get user details by ID | Yes |
| `/api/users` | POST | Create new user | Yes |
| `/api/users/{id}` | PUT | Update user | Yes |
| `/api/users/{id}` | DELETE | Delete user | Yes |

### Role Management

| Endpoint | Method | Description | Authentication |
|----------|--------|-------------|----------------|
| `/api/roles` | GET | List all roles (with pagination) | Yes |
| `/api/roles/all` | GET | List all active roles (without pagination) | Yes |
| `/api/roles/{id}` | GET | Get role details by ID | Yes |
| `/api/roles` | POST | Create new role | Yes |
| `/api/roles/{id}` | PUT | Update role | Yes |
| `/api/roles/{id}` | DELETE | Delete role | Yes |

### Division Management

| Endpoint | Method | Description | Authentication |
|----------|--------|-------------|----------------|
| `/api/divisions` | GET | List all divisions (with pagination) | Yes |
| `/api/divisions/all` | GET | List all active divisions (without pagination) | Yes |
| `/api/divisions/{id}` | GET | Get division details by ID | Yes |
| `/api/divisions` | POST | Create new division | Yes |
| `/api/divisions/{id}` | PUT | Update division | Yes |
| `/api/divisions/{id}` | DELETE | Delete division | Yes |

### Position Management

| Endpoint | Method | Description | Authentication |
|----------|--------|-------------|----------------|
| `/api/positions` | GET | List all positions (with pagination) | Yes |
| `/api/positions/all` | GET | List all active positions (without pagination) | Yes |
| `/api/positions/{id}` | GET | Get position details by ID | Yes |
| `/api/positions` | POST | Create new position | Yes |
| `/api/positions/{id}` | PUT | Update position | Yes |
| `/api/positions/{id}` | DELETE | Delete position | Yes |

### Dashboard

| Endpoint | Method | Description | Authentication |
|----------|--------|-------------|----------------|
| `/api/dashboard/statistics` | GET | Get dashboard statistics | Yes |

### Health Check

| Endpoint | Method | Description | Authentication |
|----------|--------|-------------|----------------|
| `/health` | GET | Check server status | No |

## Project Structure

```
admin-dashboard/
├── cmd/
│   └── api/
│       └── main.go       # Application entry point
├── internal/
│   ├── config/           # Configuration and environment settings
│   ├── handlers/         # HTTP request handlers
│   ├── middleware/       # HTTP middleware components
│   ├── models/           # Data models and DTOs
│   ├── repository/       # Database access layer
│   ├── services/         # Business logic layer
│   └── utils/            # Utility functions and helpers
├── Dockerfile            # Docker configuration
├── Procfile              # Railway deployment configuration
├── railway.json          # Railway service configuration
├── go.mod                # Go modules declaration
└── go.sum                # Go modules checksums
```

## Getting Started

### Prerequisites

- Go (version 1.16 or higher)
- PostgreSQL (version 12 or higher)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/admin-dashboard.git
   cd admin-dashboard
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up environment variables (see below)

4. Run the application:
   ```bash
   go run cmd/api/main.go
   ```

### Environment Variables

Create a `.env` file in the root directory with the following variables:

```
# Database Configuration
DB_CONNECTION=postgresql
DB_HOST=localhost
DB_PORT=5432
DB_DATABASE=admin_dashboard
DB_USERNAME=postgres
DB_PASSWORD=your_password
DB_SSL_MODE=disable

# JWT Configuration
JWT_SECRET=your_jwt_secret_key
JWT_EXPIRY=24

# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=3000
```

## Database Schema

The application uses the following database schema:

- **users**: Stores user information and credentials
- **roles**: Defines different roles in the system
- **user_roles**: Links users to their assigned roles (many-to-many)
- **divisions**: Organizational divisions
- **positions**: Job positions within the organization

All tables include audit columns (created_at, created_by, updated_at, updated_by).

## Authentication

The application uses JWT (JSON Web Token) for authentication. To access protected endpoints:

1. Get a token by sending a POST request to `/api/auth/login` with valid credentials
2. Include the token in the `Authorization` header of subsequent requests:
   ```
   Authorization: Bearer <your_token>
   ```

## Deployment

The application is configured for deployment on Railway. The following files are included:

- `Dockerfile`: Docker container configuration
- `Procfile`: Process type declaration for Railway
- `railway.json`: Railway service configuration

## API Documentation

The API is documented using Swagger. Once the application is running, you can access the Swagger UI at:

```
http://localhost:3000/swagger/index.html
```

### Example API Requests

#### Login

```json
POST /api/auth/login
{
  "email": "admin@company.com",
  "password": "Admin123#"
}
```

#### Create User

```json
POST /api/users
{
  "employee_id": "EMP005",
  "name": "John Doe",
  "email": "john.doe@company.com",
  "password": "password123",
  "phone": "+6281234567005",
  "address": "Jl. Jakarta No. 123",
  "birthdate": "1990-01-15",
  "join_date": "2023-01-15",
  "division_id": 1,
  "position_id": 5,
  "is_manager": false,
  "manager_id": null,
  "role_ids": [4]
}
```

#### Create Role

```json
POST /api/roles
{
  "name": "operator",
  "level": 5
}
```

#### Create Division

```json
POST /api/divisions
{
  "code": "MFG",
  "name": "Manufacturing"
}
```

#### Create Position

```json
POST /api/positions
{
  "code": "TECH",
  "name": "Technician"
}
```


## Kontribusi

Kontribusi untuk meningkatkan proyek ini sangat diterima. Silakan fork repositori, buat perubahan, dan kirimkan pull request.

## Lisensi

[MIT License](LICENSE)

## Kontak

Jika Anda memiliki pertanyaan, silakan buka issue di repositori ini atau hubungi pengembang di [gratiamanullang03@gmail.com](mailto:gratiamanullang03@gmail.com).

---

&copy; 2025 Admin Dashboard Laravel. Dibuat dengan 💻 dan ❤️
