.PHONY: build up down test migrate migrate-status migrate-down create-migration run-db create-dev-db drop-dev-db

# --- Database migrations (goose) ---
# Same SQL runs in local Docker, tests, and GitHub Actions deploy (Neon).

migrate:
	docker compose run --rm migrator up

migrate-status:
	docker compose run --rm migrator status

migrate-down:
	docker compose run --rm migrator down

# Prefer `docker compose` (v2 plugin). Legacy `docker-compose` is not used in CI.

# Usage: make create-migration NAME=add_google_sub
create-migration:
	@test -n "$(NAME)" || (echo 'Usage: make create-migration NAME=add_something'; exit 1)
	go run github.com/pressly/goose/v3/cmd/goose@v3.24.3 \
		-dir db/migrations create $(NAME) sql

create-dev-db:
	docker exec -it vibegopher-postgresql-dev1 psql -U postgres -c "CREATE DATABASE vibegopher_development;"

drop-dev-db:
	docker exec -it vibegopher-postgresql-dev1 psql -U postgres -c "DROP DATABASE vibegopher_development;"

build:
	docker compose build backend

run-db:
	docker compose up -d postgresql-dev

up:
	docker compose up -d postgresql-dev
	$(MAKE) migrate
	docker compose up backend && docker compose rm -fsv

down:
	docker compose down --volumes

test:
	go test ./internal/... -count=1
	docker compose -f docker-compose.test.yml run --rm test
	docker compose -f docker-compose.test.yml down --volumes

clear-test:
	docker volume remove vibegopher_postgres_test_data

binary:
	docker compose run --rm --no-deps --entrypoint go migrator build -o build ./cmd/server

clean-containers:
	docker rm -f $$(docker ps -a -q)

clean-images:
	docker image prune

act:
	act -P ubuntu-latest=ghcr.io/catthehacker/ubuntu:act-latest -W .github/workflows/unit-test.yml
