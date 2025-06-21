.PHONY: build up down

create_migration:
	docker compose run --rm atlas-dev atlas migrate diff $(NAME) --env dev

migrate:
	docker compose run --rm atlas-dev atlas migrate apply --env dev

migrate-test-db:
	echo "Database setup is now handled automatically by the test container"

schema-inspect:
	docker compose run --rm atlas-dev atlas schema inspect --env dev

schema-apply:
	docker compose run --rm atlas-dev atlas schema apply --env dev

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
	docker compose run --rm atlas-dev atlas migrate apply --env dev
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
