.PHONY: migrate-up migrate-down lint run

migrate-up:
	migrate -path migrations/sql -database "postgres://postgres:postgres@localhost:5432/usdt_rates?sslmode=disable" up

migrate-down:
	migrate -path migrations/sql -database "postgres://postgres:postgres@localhost:5432/usdt_rates?sslmode=disable" down

lint:
	golangci-lint run ./...

run:
	go run cmd/server/main.go

test:
	go test -v ./...