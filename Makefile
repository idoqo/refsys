include .env
export

PG_HOST ?= $(shell docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' refsys_postgres_1)
PG_URL ?= "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${PG_HOST}:${POSTGRES_DB_PORT}/${POSTGRES_DB}?sslmode=disable"

migrate-up:
	migrate -database ${PG_URL} -path ./sql/ up

migrate-down:
	migrate -database ${PG_URL} -path ./sql/ down

# creates a new migration in the sql directory
# e.g: `make migration name=create_users_table` creates a new migration named "create_users_table"
migration:
	migrate create -ext sql -dir ./sql/ -seq $(name)
