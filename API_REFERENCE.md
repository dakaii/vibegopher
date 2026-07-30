# VibegGopher API Reference

VibegGopher is a Twitter-like social media API built with Go. This document outlines all available endpoints.

## Base URL

```
http://localhost:8081/api
```

## Authentication

Most endpoints require authentication via Bearer token. Include the token in the Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

## Endpoints

### Authentication

#### Sign Up

- **POST** `/signup`
- **Description**: Create a new user account
- **Body**:
  ```json
  {
    "username": "string (min 7 characters)",
    "password": "string"
  }
  ```
- **Response**:
  ```json
  {
    "tokenType": "Bearer",
    "token": "string",
    "expiresIn": 3600
  }
  ```

#### Login

- **POST** `/login`
- **Description**: Authenticate an existing user
- **Body**:
  ```json
  {
    "username": "string",
    "password": "string"
  }
  ```
- **Response**:
  ```json
  {
    "tokenType": "Bearer",
    "token": "string",
    "expiresIn": 3600
  }
  ```

#### Get Current User

- **GET** `/me`
- **Description**: Get current authenticated user info
- **Headers**: `Authorization: Bearer <token>`
- **Response**:
  ```json
  {
    "username": "string"
  }
  ```

### Posts

#### Create Post

- **POST** `/posts`
- **Description**: Create a new post
- **Headers**: `Authorization: Bearer <token>`
- **Body**:
  ```json
  {
    "content": "string (1-280 characters)"
  }
  ```
- **Response**:
  ```json
  {
    "id": "uuid",
    "content": "string",
    "user_id": "uuid",
    "created_at": "timestamp",
    "updated_at": "timestamp",
    "user": {
      "id": "uuid",
      "username": "string",
      "created_at": "timestamp",
      "updated_at": "timestamp"
    }
  }
  ```

#### Get All Posts

- **GET** `/posts`
- **Description**: Retrieve all posts (newest first)
- **Response**: Array of post objects

#### Get Post by ID

- **GET** `/posts/{id}`
- **Description**: Retrieve a specific post by ID
- **Parameters**:
  - `id` (path): UUID of the post
- **Response**: Post object

#### Get Posts by User ID

- **GET** `/posts/user/{userId}`
- **Description**: Retrieve all posts by a specific user
- **Parameters**:
  - `userId` (path): UUID of the user
- **Response**: Array of post objects

#### Update Post

- **PATCH** `/posts/{id}`
- **Description**: Update a post (only post owner can update)
- **Headers**: `Authorization: Bearer <token>`
- **Parameters**:
  - `id` (path): UUID of the post
- **Body**:
  ```json
  {
    "content": "string (1-280 characters)"
  }
  ```
- **Response**: Updated post object

#### Delete Post

- **DELETE** `/posts/{id}`
- **Description**: Delete a post (only post owner can delete)
- **Headers**: `Authorization: Bearer <token>`
- **Parameters**:
  - `id` (path): UUID of the post
- **Response**: 204 No Content

### Comments

#### Create Comment

- **POST** `/comments`
- **Description**: Create a new comment on a post
- **Headers**: `Authorization: Bearer <token>`
- **Body**:
  ```json
  {
    "content": "string (1-280 characters)",
    "post_id": "uuid"
  }
  ```
- **Response**:
  ```json
  {
    "id": "uuid",
    "content": "string",
    "user_id": "uuid",
    "post_id": "uuid",
    "created_at": "timestamp",
    "updated_at": "timestamp",
    "user": {
      "id": "uuid",
      "username": "string",
      "created_at": "timestamp",
      "updated_at": "timestamp"
    }
  }
  ```

#### Get Comments by Post ID

- **GET** `/comments/post/{postId}`
- **Description**: Retrieve all comments for a specific post
- **Parameters**:
  - `postId` (path): UUID of the post
- **Response**: Array of comment objects

#### Get Comments by User ID

- **GET** `/comments/user/{userId}`
- **Description**: Retrieve all comments by a specific user
- **Parameters**:
  - `userId` (path): UUID of the user
- **Response**: Array of comment objects

#### Get Comment by ID

- **GET** `/comments/{id}`
- **Description**: Retrieve a specific comment by ID
- **Parameters**:
  - `id` (path): UUID of the comment
- **Response**: Comment object

#### Update Comment

- **PATCH** `/comments/{id}`
- **Description**: Update a comment (only comment owner can update)
- **Headers**: `Authorization: Bearer <token>`
- **Parameters**:
  - `id` (path): UUID of the comment
- **Body**:
  ```json
  {
    "content": "string (1-280 characters)"
  }
  ```
- **Response**: Updated comment object

#### Delete Comment

- **DELETE** `/comments/{id}`
- **Description**: Delete a comment (only comment owner can delete)
- **Headers**: `Authorization: Bearer <token>`
- **Parameters**:
  - `id` (path): UUID of the comment
- **Response**: 204 No Content

## Error Responses

All endpoints may return the following error format:

```json
{
  "error": "string"
}
```

### Common Status Codes

- `200 OK` - Successful request
- `201 Created` - Resource created successfully
- `204 No Content` - Successful deletion
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Authentication required or failed
- `403 Forbidden` - Access denied (e.g., trying to modify another user's content)
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

## Database Setup

1. Make sure PostgreSQL is running
2. Run migrations:
   ```bash
   make migrate
   ```
3. Start the development server:
   ```bash
   make dev
   ```

## Testing

Run the test suite:

```bash
make test
```

## Features

- ✅ User registration and authentication
- ✅ JWT-based authorization
- ✅ Create, read, update, delete posts
- ✅ Create, read, update, delete comments
- ✅ User-specific post and comment retrieval
- ✅ Comprehensive test coverage
- ✅ Proper error handling and validation
- ✅ Database relationships with foreign keys
- ✅ Authorization checks (users can only modify their own content)
