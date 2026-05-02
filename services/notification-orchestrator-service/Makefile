.PHONY: run dev migrate create-migration db-reset

run:
	go run ./cmd/app

dev:
	air

migrate:
	goose -dir=./migrations postgres "postgres://postgres:postgres@localhost:5433/notification_orchestrator?sslmode=disable" up

create-migration:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make create-migration name=<migration_name>"; \
		exit 1; \
	fi
	goose -dir=./migrations create $(name) sql

db-reset:
	goose -dir=./migrations postgres "postgres://postgres:postgres@localhost:5433/notification_orchestrator?sslmode=disable" reset