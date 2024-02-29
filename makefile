.PHONY: default run e2e

DB_URL=postgres://postgres:docker@localhost:5432/moodfy?sslmode=disable

default: run

run:
	@env DB_URL=$(DB_URL) go run main.go

e2e:
	@env DB_URL=$(DB_URL) go test ./e2e -v
