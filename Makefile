ifneq (,$(wildcard ./.env))
    include .env
    export
endif
run:
	air

migrate:
	go install github.com/pressly/goose/v3/cmd/goose@latest
	goose -dir ./db/migrations postgres $(DB_URL) up

migrate-down:
	goose -dir ./db/migrations postgres $(DB_URL) down

migration:
	goose -dir ./db/migrations create $(name) sql