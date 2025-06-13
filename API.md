# REST API Documentation

This application provides a simple REST API for user authentication and management.

## Base URL

```
http://localhost:8081/api
```

## Endpoints

### Authentication

#### POST /signup

Register a new user account.

**Request Body:**

```json
{
  "username": "string (min 7 characters)",
  "password": "string"
}
```

**Response (201 Created):**

```json
{
  "tokenType": "Bearer",
  "token": "jwt_token_string",
  "expiresIn": 1234567890
}
```

**Error Response (400 Bad Request):**

```json
{
  "error": "error message"
}
```

#### POST /login

Authenticate an existing user.

**Request Body:**

```json
{
  "username": "string",
  "password": "string"
}
```

**Response (200 OK):**

```json
{
  "tokenType": "Bearer",
  "token": "jwt_token_string",
  "expiresIn": 1234567890
}
```

**Error Response (401 Unauthorized):**

```json
{
  "error": "error message"
}
```

### User Information

#### GET /me

Get current authenticated user information.

**Headers:**

```
Authorization: Bearer <jwt_token>
```

**Response (200 OK):**

```json
{
  "username": "string"
}
```

**Error Response (401 Unauthorized):**

```json
{
  "error": "Unauthorized"
}
```

## Example Usage

### 1. Register a new user

```bash
curl -X POST http://localhost:8081/api/signup \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser", "password": "password123"}'
```

### 2. Login

```bash
curl -X POST http://localhost:8081/api/login \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser", "password": "password123"}'
```

### 3. Get user info

```bash
curl -X GET http://localhost:8081/api/me \
  -H "Authorization: Bearer <your_jwt_token>"
```

## Testing

Run tests with:

```bash
docker compose -f docker-compose.test.yml run --rm test
docker compose -f docker-compose.test.yml rm -fsv
```

## Features

- ✅ Simple REST API (no complex GraphQL)
- ✅ JWT-based authentication
- ✅ CORS support
- ✅ Comprehensive test coverage (81.2%)
- ✅ Simplified testing with automatic schema migration
- ✅ Factory pattern for test data generation
- ✅ Clean separation of concerns
