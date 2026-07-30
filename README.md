# VibeGopher 🐹

A Go REST API for social media functionality with posts, comments, and user management.

## 🚀 Features

- **User Management**: Registration, login, authentication with JWT
- **Posts**: Create, read, update, delete posts
- **Comments**: Comment on posts with full CRUD operations
- **Database**: PostgreSQL (Neon in cloud) with versioned goose migrations
- **Testing**: Same goose migrations as local/prod (no AutoMigrate drift)
- **Docker**: Containerized application for easy deployment

## 🛠️ Tech Stack

- **Backend**: Go with Gorilla Mux router
- **ORM**: GORM (queries only — not schema ownership)
- **Migrations**: [goose](https://github.com/pressly/goose) SQL migrations in `db/migrations`
- **Authentication**: Google Sign-In (primary) + app JWT; password auth kept for legacy/tests
- **Frontend**: Vue 3 + TypeScript 7 SPA in `frontend/` (Biome lint/format)
- **AI**: `@vibe_critic` bot worker (Gemini) on async `bot_jobs`
- **Containerization**: Docker & Docker Compose
- **Cloud / IaC**: Pulumi → GCP (Cloud Run, Artifact Registry, Secret Manager)
- **CI/CD**: GitHub Actions (test, deploy + migrate, destroy)
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

### 4. Migrate & Start

```bash
# Apply versioned migrations (goose)
make migrate

# Or start DB + migrate + API together
make up
```

The API will be available at `http://localhost:8081`

## 🔧 Available Commands

### Development

```bash
# Build all containers
make build

# Start PostgreSQL database
make run-db

# Apply / check migrations
make migrate
make migrate-status

# Create a new migration file
make create-migration NAME=add_something

# Start DB + migrate + API
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

Tests apply the same goose migrations as local/CI (`db/migrations`), then truncate app tables between cases (the `goose_db_version` table is kept).

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

- **backend**: API container (production Dockerfile)
- **postgresql-dev**: PostgreSQL 16 for development
- **migrator**: one-shot `go run ./cmd/migrate` (golang image)

### Environment Files

- `config/development.conf` - Development environment variables
- `config/test.conf` - Test environment variables

## 🔧 Configuration

Key environment variables:

```bash
PORT=8080
AUTH_SECRET=secret_key
POSTGRES_HOST=postgresql-dev
POSTGRES_DB=vibegopher_development
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_PORT=5432
HASH_COST=14
# Cloud / Neon:
# DATABASE_URL=postgres://user:pass@host/db?sslmode=require
```

## ☁️ GCP deploy (Pulumi + GitHub Actions)

Infrastructure lives in [`infra/`](./infra/). Secret placement (GCP Secret Manager vs Pulumi vs GitHub) is documented in [`infra/SECRETS.md`](./infra/SECRETS.md).

| Workflow | Trigger | Behavior |
|----------|---------|----------|
| [deploy.yml](./.github/workflows/deploy.yml) | Push to `main` / manual | Build/push → `pulumi up` → sync secrets → **goose migrate** (direct Neon URL) → Cloud Run bump |
| [destroy.yml](./.github/workflows/destroy.yml) | Manual (`confirm=destroy`) | `pulumi destroy --exclude-protected` — **keeps Secret Manager by default** |

Set `destroy_secrets=true` on the destroy workflow only when you intentionally want Secret Manager secrets removed.

Bootstrap and required GitHub variables/secrets: see [`infra/README.md`](./infra/README.md) and [`docs/DEPLOY.md`](./docs/DEPLOY.md).

## 🖥️ Frontend (Vue 3)

```bash
cd frontend
cp .env.example .env   # VITE_GOOGLE_CLIENT_ID + VITE_API_BASE_URL
npm install
npm run dev            # http://localhost:5173
```

## 🤖 AI critic bot

New posts/comments enqueue `bot_jobs`. With `BOT_WORKER_ENABLED=true` (default in Cloud Run / local compose), the API process polls and `@vibe_critic` replies via Gemini. Standalone worker: `go run ./cmd/bot`.

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
make migrate
make migrate-status
```

**Clean start:**

```bash
make down
make build
make up
```
