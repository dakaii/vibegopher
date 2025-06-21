.PHONY: build up down

# Atlas migration commands
create_migration:
	@echo "Creating migration: $(NAME)"
	@echo "Make sure database is running: make run-db"
	@echo "Creating temporary dev database if needed..."
	@docker exec vibegopher-postgresql-dev1 psql -U postgres -c "CREATE DATABASE IF NOT EXISTS atlas_dev;" 2>/dev/null || true
	docker compose run --rm atlas-dev migrate diff $(NAME) --env dev --dev-url "postgres://postgres:postgres@vibegopher-postgresql-dev1:5432/atlas_dev?sslmode=disable"

migrate:
	docker compose run --rm atlas-dev migrate apply --env dev

migrate-status:
	docker compose run --rm atlas-dev migrate status --env dev

migrate-test-db:
	echo "Database setup is now handled automatically by the test container"

# Schema diff commands (for development only)
schema-apply-dev:
	@echo "⚠️  WARNING: This applies schema changes directly without migration files"
	@echo "⚠️  Only use this for local development!"
	@read -p "Continue? (y/N): " confirm && [ "$$confirm" = "y" ] || exit 1
	docker compose run --rm atlas-dev schema apply --env dev

schema-diff-preview:
	@echo "Previewing schema differences..."
	docker compose run --rm atlas-dev schema diff --env dev

schema-inspect:
	docker compose run --rm atlas-dev schema inspect --env dev

create-dev-db:
	docker exec -it vibegopher-postgresql-dev1 psql -U postgres -c "CREATE DATABASE vibegopher_development;"

drop-dev-db:
	docker exec -it vibegopher-postgresql-dev1 psql -U postgres -c "DROP DATABASE vibegopher_development;"

build:
	docker compose build

run-db:
	docker compose up -d postgresql-dev

up:
	docker compose up -d postgresql-dev
	docker compose up backend && docker compose rm -fsv

down:
	docker compose down --volumes

test:
	docker compose -f docker-compose.test.yml run --rm test
	docker compose -f docker-compose.test.yml rm -fsv

clear-test:
	docker volume remove vibegopher_postgres_test_data

binary:
	docker compose run --rm backend go build -o build ./cmd/server/main.go

clean-containers:
	docker rm -f $(docker ps -a -q)

clean-images:
	docker image prune

act:
	act -P ubuntu-latest=ghcr.io/catthehacker/ubuntu:act-latest -W .github/workflows/unit-test.yml
