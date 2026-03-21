# Authentication API Documentation

## Overview
This Go backend implements JWT-based authentication with the following features:
- User registration and login
- Password hashing with bcrypt
- JWT token generation and validation
- Protected routes with middleware
- Role-based access control (Regular and Admin roles)

## Environment Setup

Add the following to your `.env` file:
```
JWT_SECRET=your-secure-secret-key-change-this-in-production-min-32-chars
```

⚠️ **Important**: Use a strong, random secret key in production (at least 32 characters).

## API Endpoints

### Public Endpoints

#### 1. Register a New User
**POST** `/auth/register`

Request body:
```json
{
  "email": "user@example.com",
  "password": "securepassword123",
  "role": 0
}
```

- `role`: 0 = Regular user, 1 = Admin (defaults to 0 if not provided)
- `password`: Must be at least 8 characters

Response (201 Created):
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "role": 0
  }
}
```

#### 2. Login
**POST** `/auth/login`

Request body:
```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

Response (200 OK):
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "role": 0
  }
}
```

### Protected Endpoints (Require JWT Token)

All protected endpoints require the `Authorization` header:
```
Authorization: Bearer <your-jwt-token>
```

#### 3. Get Current User
**GET** `/api/me`

Response (200 OK):
```json
{
  "id": 1,
  "email": "user@example.com",
  "role": 0
}
```

### Admin-Only Endpoints (Require JWT Token + Admin Role)

#### 4. Admin Example Route
**GET** `/admin/users`

Response (200 OK):
```json
{
  "message": "Admin only route"
}
```

## Usage Examples

### Using curl

1. **Register a new user:**
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword123"
  }'
```

2. **Login:**
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword123"
  }'
```

3. **Access protected route:**
```bash
curl -X GET http://localhost:8080/api/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
```

### Using JavaScript/Fetch

```javascript
// Register
const registerResponse = await fetch('http://localhost:8080/auth/register', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    email: 'user@example.com',
    password: 'securepassword123'
  })
});
const { token, user } = await registerResponse.json();

// Login
const loginResponse = await fetch('http://localhost:8080/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    email: 'user@example.com',
    password: 'securepassword123'
  })
});
const data = await loginResponse.json();

// Access protected route
const meResponse = await fetch('http://localhost:8080/api/me', {
  headers: { 'Authorization': `Bearer ${token}` }
});
const userData = await meResponse.json();
```

## Implementation Details

### File Structure
```
backend/
├── internal/
│   ├── auth/
│   │   └── auth.go           # JWT & password utilities
│   ├── middleware/
│   │   └── jwt.go            # JWT authentication middleware
│   ├── models/
│   │   └── user.go           # User model with Role
│   └── server/
│       ├── auth_handlers.go  # Register, Login, Me handlers
│       └── routes.go         # Route configuration
```

### Security Features

1. **Password Hashing**: Uses bcrypt with default cost (10)
2. **JWT Tokens**: 
   - HS256 signing algorithm
   - 24-hour expiration
   - Custom claims include user ID, email, and role
3. **Token Validation**:
   - Signature verification
   - Expiration check
   - Issuer validation
   - Custom claims validation
4. **Role-Based Access**: Middleware to restrict admin-only routes

### JWT Token Structure

The JWT token contains these claims:
```json
{
  "user_id": 1,
  "email": "user@example.com",
  "role": 0,
  "exp": 1234567890,
  "iat": 1234567890,
  "nbf": 1234567890,
  "iss": "booklet-api"
}
```

### Middleware Usage

To protect routes, use the middleware in your route definitions:

```go
// Single middleware
protected := app.Group("/api", middleware.JWTProtected())

// Multiple middleware (JWT + Admin check)
admin := app.Group("/admin", middleware.JWTProtected(), middleware.RequireAdmin())
```

Access user information in handlers:
```go
func MyHandler(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)
    email := c.Locals("email").(string)
    role := c.Locals("role").(int)
    // ... your handler logic
}
```

## Error Responses

### 400 Bad Request
```json
{
  "error": "Invalid request body"
}
```

### 401 Unauthorized
```json
{
  "error": "Invalid or expired token"
}
```

### 403 Forbidden
```json
{
  "error": "Admin access required"
}
```

### 409 Conflict
```json
{
  "error": "User with this email already exists"
}
```

### 500 Internal Server Error
```json
{
  "error": "Failed to process registration"
}
```
