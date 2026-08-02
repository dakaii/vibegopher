# VibeGopher 🐹

A Go REST API for social media functionality with posts, comments, and user management.

## 🚀 Features

- **User Management**: Registration, login, authentication with JWT
- **Posts**: Create, read, update, delete posts
- **Comments**: Comment on posts with full CRUD operations
- **Database**: PostgreSQL with Atlas migrations for production
- **Testing**: Comprehensive test suite with auto-migration
- **Docker**: Containerized application for easy deployment

## 🛠️ Tech Stack

- **Backend**: Go with Gorilla Mux router
- **Database**: PostgreSQL with GORM ORM
- **Migrations**: Atlas for production, GORM AutoMigrate for testing
- **Authentication**: JWT tokens
- **Containerization**: Docker & Docker Compose
- **Testing**: Go testing with test factories

## 📋 Prerequisites

- **Docker** and **Docker Compose**
- **Make** (for running commands)

## 🚀 Quick Start

### 1. Clone the Repository

```bash
git clone <repository-url>
cd vibegopher
```

### 2. Build the Application

```bash
make build
```

### 3. Start the Database

```bash
make run-db
```

### 4. Set Up the Database

```bash
# Create the development database
make create-dev-db

# Apply the database schema
docker exec -i vibegopher-postgresql-dev1 psql -U postgres -d vibegopher_development < schema.sql
```

### 5. Start the Application

```bash
docker compose up backend
```

The API will be available at `http://localhost:8081`

## 🔧 Available Commands

### Development

```bash
# Build all containers
make build

# Start PostgreSQL database
make run-db

# Create development database
make create-dev-db

# Drop development database (for cleanup)
make drop-dev-db

# Run migrations (when available)
make migrate

# Start the full application
make up
```

### Testing

```bash
# Run all tests
make test

# Clean up test containers
make clear-test
```

### Cleanup

```bash
# Stop all containers and remove volumes
make down

# Clean up containers
make clean-containers

# Clean up images
make clean-images
```

## 📡 API Endpoints

### Authentication

- `POST /api/signup` - User registration
- `POST /api/login` - User login
- `GET /api/me` - Get current user info (requires auth)

### Posts

- `GET /api/posts` - Get all posts
- `POST /api/posts` - Create a new post (requires auth)
- `GET /api/posts/{id}` - Get post by ID
- `GET /api/posts/user/{userId}` - Get posts by user ID
- `PATCH /api/posts/{id}` - Update post (requires auth)
- `DELETE /api/posts/{id}` - Delete post (requires auth)

### Comments

- `POST /api/comments` - Create a comment (requires auth)
- `GET /api/comments/post/{postId}` - Get comments for a post
- `GET /api/comments/user/{userId}` - Get comments by user
- `GET /api/comments/{id}` - Get comment by ID
- `PATCH /api/comments/{id}` - Update comment (requires auth)
- `DELETE /api/comments/{id}` - Delete comment (requires auth)

## 🧪 Testing

The application uses a hybrid approach for database management:

- **Production**: Atlas migrations for versioned, reviewed schema changes
- **Testing**: GORM AutoMigrate for fast, automatic schema setup

Run tests with:

```bash
make test
```

## 🗄️ Database Schema

The application uses three main tables:

### Users

- `id` (UUID, Primary Key)
- `created_at`, `updated_at`, `deleted_at` (Timestamps)
- `username` (String, Unique)
- `password` (String, Hashed)

### Posts

- `id` (UUID, Primary Key)
- `created_at`, `updated_at` (Timestamps)
- `content` (Text)
- `user_id` (UUID, Foreign Key to Users)

### Comments

- `id` (UUID, Primary Key)
- `created_at`, `updated_at` (Timestamps)
- `content` (Text)
- `user_id` (UUID, Foreign Key to Users)
- `post_id` (UUID, Foreign Key to Posts)

## 🔒 Authentication

The API uses JWT (JSON Web Tokens) for authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## 🐳 Docker Configuration

### Development Containers

- **backend**: Go application container
- **postgresql-dev**: PostgreSQL database for development
- **atlas-dev**: Atlas CLI for migrations

### Environment Files

- `config/development.conf` - Default development env (local Docker Postgres)
- `config/test.conf` - Test environment variables
- `config/local.conf.example` - Template for machine-local overrides (e.g. Neon)
- `config/local.conf` - Optional gitignored override; copy from the example when needed

By default, Compose loads `config/development.conf` for the backend. If `config/local.conf` exists, it is loaded after and overrides matching variables (Docker Compose `env_file` with `required: false`). The local Postgres service always uses `development.conf`, so Neon credentials never reconfigure that container.

**Use Neon (or another remote DB) locally:**

```bash
cp config/local.conf.example config/local.conf
# Edit config/local.conf with your credentials (POSTGRES_SSLMODE=require for Neon)
docker compose up backend
```

Without `config/local.conf`, the backend keeps using Docker Postgres on port `5431`.

## 🔧 Configuration

Key environment variables:

```bash
PORT=8081
AUTH_SECRET=secret_key
POSTGRES_HOST=host.docker.internal
POSTGRES_DB=vibegopher_development
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_PORT=5431
POSTGRES_SSLMODE=disable
HASH_COST=14
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite: `make test`
6. Submit a pull request

## 📝 License

This project is licensed under the MIT License.

## 🆘 Troubleshooting

### Common Issues

**Backend container fails to start:**

```bash
# Rebuild the containers
make build
```

**Database connection issues:**

```bash
# Restart the database
make run-db
```

**Schema not applied:**

```bash
# Apply schema manually
docker exec -i vibegopher-postgresql-dev1 psql -U postgres -d vibegopher_development < schema.sql
```

**Clean start:**

```bash
# Complete cleanup and restart
make down
make build
make run-db
make create-dev-db
docker exec -i vibegopher-postgresql-dev1 psql -U postgres -d vibegopher_development < schema.sql
docker compose up backend
```
