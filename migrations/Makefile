create-migration: 
	@ goose create $(name) sql

migration-up:
	@ go run main.go up

migration-down:
	@ go run main.go down

build:
	@ go build -o migration

migrate:
	@ ./migration up