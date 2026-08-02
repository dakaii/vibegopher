.PHONY: build up down test migrate migrate-status migrate-down create-migration run-db create-dev-db drop-dev-db

# Local Docker Postgres lives on the Compose profile "local-db".
# When config/local.conf is present (e.g. Neon), skip that profile so
# backend/migrator do not start postgresql-dev.
USE_LOCAL_DB := $(shell test ! -f config/local.conf && echo yes)
COMPOSE_LOCAL_DB := $(if $(USE_LOCAL_DB),--profile local-db,)

# --- Database migrations (goose) ---
# Same SQL runs in local Docker, tests, and GitHub Actions deploy (Neon).

migrate:
	docker compose $(COMPOSE_LOCAL_DB) run --rm migrator up

migrate-status:
	docker compose $(COMPOSE_LOCAL_DB) run --rm migrator status

migrate-down:
	docker compose $(COMPOSE_LOCAL_DB) run --rm migrator down

# Usage: make create-migration NAME=add_something
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
	docker compose --profile local-db up -d postgresql-dev

up:
ifneq ($(USE_LOCAL_DB),)
	docker compose --profile local-db up -d postgresql-dev
endif
	$(MAKE) migrate
	docker compose $(COMPOSE_LOCAL_DB) up backend && docker compose $(COMPOSE_LOCAL_DB) rm -fsv

down:
	docker compose --profile local-db down --volumes

test:
	docker compose -f docker-compose.test.yml run --rm test
	docker compose -f docker-compose.test.yml rm -fsv

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
